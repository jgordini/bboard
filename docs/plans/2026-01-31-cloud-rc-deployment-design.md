# UAB Cloud.rc Deployment Design

**Date:** 2026-01-31
**Status:** Design Complete
**Target Platform:** UAB Research Computing Cloud (OpenStack)

## Overview

Production deployment of the bboard (Fider) application on UAB's cloud.rc OpenStack platform with persistent PostgreSQL storage, HTTPS/SSL, and UAB SAML authentication.

## Requirements

- **Deployment Type:** Production web service
- **Data Persistence:** Critical - database must survive instance failures
- **Authentication:** UAB SAML/Shibboleth integration
- **Access:** Public-facing, accessible to UAB community
- **Availability:** Single instance with persistent volume backup strategy

## Architecture

### Infrastructure Components

```
┌─────────────────────────────────────────┐
│  OpenStack Instance (cloud.rc VM)      │
│  Ubuntu 24.04 LTS, m1.medium+          │
│                                         │
│  ┌──────────────────────────────────┐  │
│  │ Docker Compose Stack             │  │
│  │                                  │  │
│  │  ┌────────────┐  ┌────────────┐ │  │
│  │  │   Nginx    │  │  bboard    │ │  │
│  │  │ (80/443)   │→ │  (app:3000)│ │  │
│  │  └────────────┘  └──────┬─────┘ │  │
│  │                          │       │  │
│  │                  ┌───────▼─────┐ │  │
│  │                  │ PostgreSQL  │ │  │
│  │                  │   (db:5432) │ │  │
│  │                  └──────┬──────┘ │  │
│  └─────────────────────────┼────────┘  │
│                            │           │
│                   Volume mount:        │
│                   /mnt/postgres-data   │
│                            │           │
│              ┌─────────────▼─────────┐ │
│              │  OpenStack Volume    │ │
│              │  50GB ext4           │ │
│              │  blazeboard-postgres │ │
│              └──────────────────────┘ │
│                                        │
│  Floating IP: XXX.XXX.XXX.XXX         │
└────────────────────────────────────────┘
```

### OpenStack Resources

| Resource | Specification | Purpose |
|----------|--------------|---------|
| Instance | Ubuntu 24.04 LTS, m1.medium (2 vCPU, 4GB RAM) | Application host |
| Volume | 50GB, ext4 filesystem | PostgreSQL persistent storage |
| Floating IP | Public IPv4 | External access |
| Security Groups | TCP 22, 80, 443 | SSH, HTTP, HTTPS access |
| Network | Project network + router | Campus network connectivity |

### Docker Services

1. **nginx** - Reverse proxy with SSL termination
2. **app** - Fider application (getfider/fider:stable)
3. **db** - PostgreSQL 17 with persistent volume
4. **mailhog** (optional) - Testing SMTP server

## Configuration Files

### docker-compose.yml

```yaml
services:
  # PostgreSQL - Data persisted to OpenStack volume
  db:
    restart: always
    image: postgres:17
    volumes:
      - /mnt/postgres-data:/var/lib/postgresql/data
    environment:
      POSTGRES_USER: fider
      POSTGRES_PASSWORD: ${DB_PASSWORD}
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U fider"]
      interval: 10s
      timeout: 5s
      retries: 5

  # Fider Application
  app:
    restart: always
    image: getfider/fider:stable
    depends_on:
      db:
        condition: service_healthy
    volumes:
      - ./ssl:/app/etc:ro
    environment:
      BASE_URL: https://blazeboard.cloud.rc.uab.edu
      DATABASE_URL: postgres://fider:${DB_PASSWORD}@db:5432/fider?sslmode=disable
      JWT_SECRET: ${JWT_SECRET}
      EMAIL_NOREPLY: noreply@uab.edu
      EMAIL_SMTP_HOST: ${SMTP_HOST}
      EMAIL_SMTP_PORT: ${SMTP_PORT}
      EMAIL_SMTP_USERNAME: ${SMTP_USER}
      EMAIL_SMTP_PASSWORD: ${SMTP_PASS}
      SAML_ENTITY_ID: https://blazeboard.cloud.rc.uab.edu/saml/metadata
      SAML_IDP_ENTITY_ID: https://idp.uab.edu/idp/shibboleth
      SAML_IDP_SSO_URL: https://idp.uab.edu/idp/profile/SAML2/Redirect/SSO
      SAML_IDP_CERT: ${SAML_IDP_CERT}
      SAML_SP_CERT_PATH: etc/sp.crt
      SAML_SP_KEY_PATH: etc/sp.key

  # Nginx reverse proxy
  nginx:
    restart: always
    image: nginx:alpine
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx/nginx.conf:/etc/nginx/nginx.conf:ro
      - /etc/letsencrypt:/etc/letsencrypt:ro
      - /var/www/certbot:/var/www/certbot:ro
    depends_on:
      - app

  # MailHog (optional, for testing)
  mailhog:
    image: mailhog/mailhog
    restart: always
    ports:
      - "8025:8025"
```

### .env (template)

```bash
# Database Password
DB_PASSWORD=GENERATE_STRONG_PASSWORD_HERE

# JWT Secret (generate: openssl rand -base64 64)
JWT_SECRET=GENERATE_STRONG_JWT_SECRET_HERE

# Email Configuration
SMTP_HOST=smtp.uab.edu
SMTP_PORT=587
SMTP_USER=your-blazerid@uab.edu
SMTP_PASS=your-smtp-password

# For testing with MailHog instead:
# SMTP_HOST=mailhog
# SMTP_PORT=1025
# SMTP_USER=
# SMTP_PASS=

# SAML Certificate (get from UAB IT)
SAML_IDP_CERT="-----BEGIN CERTIFICATE-----
...
-----END CERTIFICATE-----"
```

### nginx/nginx.conf

```nginx
events {
    worker_connections 1024;
}

http {
    upstream fider {
        server app:3000;
    }

    # HTTP - Redirect to HTTPS
    server {
        listen 80;
        server_name blazeboard.cloud.rc.uab.edu;

        location /.well-known/acme-challenge/ {
            root /var/www/certbot;
        }

        location / {
            return 301 https://$host$request_uri;
        }
    }

    # HTTPS
    server {
        listen 443 ssl http2;
        server_name blazeboard.cloud.rc.uab.edu;

        ssl_certificate /etc/letsencrypt/live/blazeboard.cloud.rc.uab.edu/fullchain.pem;
        ssl_certificate_key /etc/letsencrypt/live/blazeboard.cloud.rc.uab.edu/privkey.pem;

        add_header Strict-Transport-Security "max-age=31536000" always;
        add_header X-Frame-Options "SAMEORIGIN" always;
        add_header X-Content-Type-Options "nosniff" always;

        client_max_body_size 10M;

        location / {
            proxy_pass http://fider;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
        }
    }
}
```

## Deployment Process

### Phase 1: Prerequisites

- [ ] Cloud.rc account access (dashboard.cloud.rc.uab.edu)
- [ ] SSH key pair created and added to OpenStack
- [ ] Domain name decided: `blazeboard.cloud.rc.uab.edu`
- [ ] Contact UAB IT for:
  - [ ] SMTP credentials
  - [ ] SAML IdP certificate
  - [ ] DNS A record
  - [ ] SAML SP registration

### Phase 2: OpenStack Setup (Dashboard)

1. **Create Volume**
   - Navigate: Project → Volumes → Volumes → Create Volume
   - Name: `blazeboard-postgres-data`
   - Size: `50 GB`
   - Type: `__DEFAULT__`

2. **Launch Instance**
   - Navigate: Project → Compute → Instances → Launch Instance
   - Name: `blazeboard-prod`
   - Source: Ubuntu 24.04 LTS
   - Flavor: `m1.medium` or larger
   - Networks: Select project network
   - Key Pair: Your SSH key

3. **Attach Volume**
   - Navigate: Volumes → blazeboard-postgres-data → Manage Attachments
   - Attach to: `blazeboard-prod`

4. **Configure Security Groups**
   - Navigate: Project → Network → Security Groups → default → Manage Rules
   - Add rules:
     - `TCP 22` (SSH) from `0.0.0.0/0`
     - `TCP 80` (HTTP) from `0.0.0.0/0`
     - `TCP 443` (HTTPS) from `0.0.0.0/0`

5. **Allocate Floating IP**
   - Navigate: Project → Network → Floating IPs → Allocate IP
   - Associate with: `blazeboard-prod`
   - Note the IP address

### Phase 3: Instance Configuration

```bash
# SSH into instance
ssh ubuntu@<floating-ip>

# Format and mount persistent volume (FIRST TIME ONLY)
sudo mkfs.ext4 /dev/vdb
sudo mkdir -p /mnt/postgres-data
sudo mount /dev/vdb /mnt/postgres-data
echo '/dev/vdb /mnt/postgres-data ext4 defaults 0 2' | sudo tee -a /etc/fstab
sudo chown -R 999:999 /mnt/postgres-data

# Install Docker
curl -fsSL https://get.docker.com | sudo sh
sudo usermod -aG docker ubuntu
# Log out and back in for group membership

# Create application directory
sudo mkdir -p /var/fider/{nginx,ssl}
sudo chown -R ubuntu:ubuntu /var/fider
mkdir -p /var/www/certbot
```

### Phase 4: Application Deployment

```bash
cd /var/fider

# Generate secrets
openssl rand -base64 32  # DB_PASSWORD
openssl rand -base64 64  # JWT_SECRET

# Create .env file with secrets
nano .env

# Create docker-compose.yml
nano docker-compose.yml

# Create simplified nginx config (HTTP only initially)
nano nginx/nginx.conf

# Start application (HTTP only)
docker compose pull
docker compose up -d

# Verify application
docker compose logs -f app
# Wait for: "http server started on :3000"
```

### Phase 5: SSL/HTTPS Setup

```bash
# Install certbot
sudo apt update && sudo apt install -y certbot

# Get SSL certificate (DNS must point to floating IP first)
sudo certbot certonly --webroot \
  -w /var/www/certbot \
  -d blazeboard.cloud.rc.uab.edu \
  --email your-email@uab.edu \
  --agree-tos

# Update nginx config to full HTTPS version
nano nginx/nginx.conf

# Restart nginx
docker compose restart nginx

# Setup auto-renewal
(sudo crontab -l 2>/dev/null; echo "0 3 * * * certbot renew --quiet && docker compose -f /var/fider/docker-compose.yml restart nginx") | sudo crontab -
```

### Phase 6: SAML Configuration

```bash
cd /var/fider/ssl

# Generate SAML SP certificates
openssl genrsa -out sp.key 2048
openssl req -new -key sp.key -out sp.csr \
  -subj "/C=US/ST=Alabama/L=Birmingham/O=UAB/CN=blazeboard.cloud.rc.uab.edu"
openssl x509 -req -days 3650 -in sp.csr -signkey sp.key -out sp.crt
chmod 600 sp.key
chmod 644 sp.crt
rm sp.csr

# Update .env with SAML_IDP_CERT from UAB IT
nano /var/fider/.env

# Update docker-compose.yml with SAML environment variables
nano /var/fider/docker-compose.yml

# Restart application
docker compose down
docker compose up -d

# Download SP metadata
curl https://blazeboard.cloud.rc.uab.edu/saml/metadata > sp-metadata.xml

# Send sp-metadata.xml and sp.crt to UAB IT for IdP registration
```

## Operations

### Daily Operations

```bash
# View logs
docker compose logs -f [app|db|nginx]

# Restart services
docker compose restart [service]

# Update Fider to latest version
docker compose pull app
docker compose up -d app
```

### Backup & Restore

```bash
# Backup database
docker compose exec db pg_dump -U fider fider > backup-$(date +%Y%m%d).sql

# Restore database
cat backup-20260131.sql | docker compose exec -T db psql -U fider fider

# Backup entire volume (snapshot via OpenStack dashboard)
# Navigate: Volumes → blazeboard-postgres-data → Create Snapshot
```

### Monitoring

```bash
# Disk usage
df -h /mnt/postgres-data

# Container health
docker compose ps

# Resource usage
docker stats

# Application health
curl -f https://blazeboard.cloud.rc.uab.edu || echo "Down"
```

### Troubleshooting

**Database connection errors:**
```bash
# Check database health
docker compose exec db pg_isready -U fider

# Check volume mount
df -h /mnt/postgres-data

# Restart database
docker compose restart db
```

**SSL certificate issues:**
```bash
# Test certificate renewal
sudo certbot renew --dry-run

# Check certificate expiry
sudo certbot certificates
```

**SAML authentication failures:**
```bash
# Check SAML metadata endpoint
curl https://blazeboard.cloud.rc.uab.edu/saml/metadata

# Verify certificates exist
ls -la /var/fider/ssl/

# Check application logs for SAML errors
docker compose logs app | grep -i saml
```

## Security Considerations

1. **Secrets Management**
   - Never commit `.env` file to git
   - Use strong passwords (32+ characters)
   - Rotate JWT_SECRET periodically

2. **Network Security**
   - Only expose ports 22, 80, 443
   - Consider restricting SSH to UAB IP ranges
   - Enable HTTPS-only (redirect HTTP)

3. **Data Protection**
   - Regular database backups (automated)
   - Volume snapshots weekly
   - Test restore procedures

4. **Updates**
   - Monitor Fider releases
   - Test updates in staging before production
   - Keep base OS updated: `sudo apt update && sudo apt upgrade`

## Cost & Resource Management

- **Instance**: m1.medium (~2 vCPU, 4GB RAM)
- **Volume**: 50GB (expandable)
- **Bandwidth**: Minimal cost on UAB network
- **Free up resources**: Delete instance when not needed, keep volume for data

## Future Enhancements

- [ ] High availability: Multi-instance with load balancer
- [ ] Separate database server for performance
- [ ] Redis caching layer
- [ ] Automated backups to object storage
- [ ] Monitoring with Prometheus/Grafana
- [ ] Container orchestration with Kubernetes

## References

- [Fider Docker Documentation](https://docs.fider.io/self-hosted/hosting-on-docker)
- [UAB Cloud.rc Documentation](https://docs.rc.uab.edu/uab_cloud/)
- [Let's Encrypt Documentation](https://letsencrypt.org/docs/)
- [UAB SAML/Shibboleth](https://www.uab.edu/it/)
