# BlazeBoard Deployment on UAB Cloud.rc

Complete deployment package for running BlazeBoard (Fider) on UAB Research Computing Cloud OpenStack infrastructure.

## Overview

This deployment uses:
- **Docker Compose** for container orchestration
- **PostgreSQL 17** with persistent OpenStack volume storage
- **Nginx** reverse proxy with SSL/TLS termination
- **Let's Encrypt** for SSL certificates
- **UAB SAML** authentication (optional)

## Prerequisites

### On Cloud.rc Dashboard

- [ ] Cloud.rc account with access to https://dashboard.cloud.rc.uab.edu
- [ ] OpenStack volume created (50GB+) named `blazeboard-postgres-data`
- [ ] Ubuntu 24.04 instance launched (m1.medium or larger)
- [ ] Volume attached to instance
- [ ] Security groups configured (TCP 22, 80, 443)
- [ ] Floating IP allocated and associated

### From UAB IT

- [ ] Domain name configured (e.g., `blazeboard.cloud.rc.uab.edu`)
- [ ] DNS A record pointing to floating IP
- [ ] SMTP credentials (optional - can use MailHog for testing)
- [ ] SAML IdP certificate (optional - for UAB authentication)

## Quick Start

### 1. Setup the Instance

SSH into your cloud.rc instance:

```bash
ssh ubuntu@<your-floating-ip>
```

Upload and run the setup script:

```bash
# Upload setup-instance.sh to the instance
# Then run:
chmod +x setup-instance.sh
./setup-instance.sh

# Log out and back in for Docker group membership
exit
ssh ubuntu@<your-floating-ip>

# Verify Docker access
docker ps
```

### 2. Deploy Application Files

Upload all files from this directory to `/var/fider/` on the instance:

```bash
# From your local machine
cd deployment/cloud-rc
scp -r * ubuntu@<your-floating-ip>:/var/fider/

# Or use rsync for updates
rsync -avz --exclude='.env' ./ ubuntu@<your-floating-ip>:/var/fider/
```

### 3. Configure Environment

On the cloud.rc instance:

```bash
cd /var/fider

# Create .env from template
cp .env.template .env

# Generate secrets
DB_PASSWORD=$(openssl rand -base64 32)
JWT_SECRET=$(openssl rand -base64 64)

# Edit .env and add the generated secrets
nano .env
# Update BASE_URL, DB_PASSWORD, JWT_SECRET, and email settings
# Save and exit (Ctrl+X, Y, Enter)

# Secure the file
chmod 600 .env
```

### 4. Initial Deployment (HTTP)

Start with HTTP to verify everything works:

```bash
cd /var/fider
./scripts/deploy.sh http
```

Verify the application is running:

```bash
# Check containers
docker compose ps

# Check logs
docker compose logs -f app

# Test HTTP access
curl http://localhost
```

### 5. Get SSL Certificate

Once HTTP is working and DNS is configured:

```bash
# Get Let's Encrypt certificate
sudo certbot certonly --webroot \
  -w /var/www/certbot \
  -d blazeboard.cloud.rc.uab.edu \
  --email your-email@uab.edu \
  --agree-tos
```

### 6. Deploy with HTTPS

Update to HTTPS configuration:

```bash
cd /var/fider

# Update BASE_URL in .env to https://
nano .env

# Deploy with HTTPS
./scripts/deploy.sh https
```

Test HTTPS access:

```bash
curl https://blazeboard.cloud.rc.uab.edu
```

### 7. Initial Application Setup

Visit `https://blazeboard.cloud.rc.uab.edu` in your browser and complete the setup wizard:

1. Create tenant (organization) name
2. Set admin email
3. Configure initial settings

### 8. Configure SAML (Optional)

To enable UAB SAML authentication:

```bash
cd /var/fider

# Generate SAML certificates
./scripts/generate-saml-certs.sh blazeboard.cloud.rc.uab.edu

# Get SAML IdP certificate from UAB IT
# Add to .env file as SAML_IDP_CERT

# Update .env with SAML configuration
nano .env
# Uncomment SAML_* variables and fill in values

# Update docker-compose.yml
nano docker-compose.yml
# Uncomment SAML environment variables in app service

# Redeploy
./scripts/deploy.sh https

# Download SP metadata
curl https://blazeboard.cloud.rc.uab.edu/saml/metadata > sp-metadata.xml

# Send sp-metadata.xml and ssl/sp.crt to UAB IT for IdP registration
```

## File Structure

```
/var/fider/
├── docker-compose.yml           # Container orchestration
├── .env                        # Environment secrets (create from template)
├── .env.template               # Environment template
├── nginx/
│   ├── nginx.conf             # HTTPS configuration
│   └── nginx-http-only.conf   # HTTP-only configuration
├── ssl/
│   ├── sp.key                 # SAML SP private key (generated)
│   └── sp.crt                 # SAML SP certificate (generated)
└── scripts/
    ├── setup-instance.sh      # Initial instance setup
    ├── deploy.sh              # Deployment script
    └── generate-saml-certs.sh # SAML certificate generation
```

## Operations

### Viewing Logs

```bash
cd /var/fider

# All logs
docker compose logs -f

# Specific service
docker compose logs -f app
docker compose logs -f db
docker compose logs -f nginx
```

### Restarting Services

```bash
cd /var/fider

# Restart all services
docker compose restart

# Restart specific service
docker compose restart app
docker compose restart nginx
```

### Updating Application

**Deploy with SCP** (build image locally, copy to instance — no registry):

```bash
cd deployment/cloud-rc
RC_HOST=ubuntu@138.26.48.197 RC_SSH_KEY=~/.ssh/cloud_key ./scripts/deploy-rc-cloud-scp.sh
```

This builds `blazeboard:latest` from the repo, saves it, syncs deployment files and the image to `/var/fider`, loads the image on the instance, and runs `deploy.sh update`.

**Sync config only** (no image; use when instance already has `blazeboard:latest` or you build on instance):

```bash
cd deployment/cloud-rc
RC_HOST=ubuntu@138.26.48.197 RC_SSH_KEY=~/.ssh/cloud_key ./scripts/update-rc-cloud.sh
```

**CI/CD (deploy on push to main):**

1. In GitHub: **Settings → Secrets and variables → Actions**, add:
   - **RC_SSH_KEY**: contents of your private key (e.g. paste the full content of `~/.ssh/cloud_key`).
   - **RC_HOST**: `ubuntu@138.26.48.197` (or your instance hostname).

2. Push to `main`; the workflow **Deploy to RC Cloud** (`.github/workflows/deploy-rc-cloud.yml`) will build the image, sync deployment files and the image to the instance, then run `deploy.sh update` there.

**On the RC cloud instance** (if you're already SSH'd in):

```bash
cd /var/fider

# Pull latest Fider version and restart
./scripts/deploy.sh update

# Or manually
docker compose pull app
docker compose up -d app
```

### Database Backup

```bash
cd /var/fider

# Create backup
docker compose exec db pg_dump -U fider fider > backup-$(date +%Y%m%d).sql

# Compress backup
gzip backup-$(date +%Y%m%d).sql
```

### Database Restore

```bash
cd /var/fider

# Restore from backup
cat backup-20260131.sql | docker compose exec -T db psql -U fider fider
```

### SSL Certificate Renewal

Certificates auto-renew if cron job is configured:

```bash
# Manual renewal test
sudo certbot renew --dry-run

# Manual renewal
sudo certbot renew
docker compose restart nginx
```

### Monitoring

```bash
# Container status
docker compose ps

# Resource usage
docker stats

# Disk usage
df -h /mnt/postgres-data

# Check application health
curl -f https://blazeboard.cloud.rc.uab.edu || echo "Down"
```

## Troubleshooting

### Container Won't Start

```bash
# Check logs
docker compose logs [service]

# Check configuration
docker compose config

# Verify environment variables
cat .env
```

### Database Connection Errors

```bash
# Check database is healthy
docker compose exec db pg_isready -U fider

# Check volume is mounted
df -h /mnt/postgres-data

# Check database logs
docker compose logs db

# Restart database
docker compose restart db
```

### SSL Certificate Issues

```bash
# Check certificate exists
sudo ls -la /etc/letsencrypt/live/blazeboard.cloud.rc.uab.edu/

# Check certificate expiry
sudo certbot certificates

# Test renewal
sudo certbot renew --dry-run

# Check nginx configuration
docker compose exec nginx nginx -t
```

### SAML Authentication Not Working

```bash
# Verify SAML certificates exist
ls -la /var/fider/ssl/

# Check SAML metadata endpoint
curl https://blazeboard.cloud.rc.uab.edu/saml/metadata

# Check application logs for SAML errors
docker compose logs app | grep -i saml

# Verify SAML environment variables
docker compose exec app env | grep SAML
```

### Application Not Accessible

```bash
# Check containers are running
docker compose ps

# Check nginx logs
docker compose logs nginx

# Test internal connectivity
docker compose exec nginx wget -O- http://app:3000

# Check firewall/security groups
# Ensure ports 80 and 443 are open in OpenStack security groups
```

## Backup Strategy

### Automated Daily Backups

Create a backup script:

```bash
cat > /var/fider/scripts/backup.sh << 'EOF'
#!/bin/bash
BACKUP_DIR="/var/fider/backups"
DATE=$(date +%Y%m%d_%H%M%S)

mkdir -p "$BACKUP_DIR"

# Backup database
docker compose -f /var/fider/docker-compose.yml exec -T db \
  pg_dump -U fider fider | gzip > "$BACKUP_DIR/db-$DATE.sql.gz"

# Keep only last 7 days of backups
find "$BACKUP_DIR" -name "db-*.sql.gz" -mtime +7 -delete

echo "Backup completed: $BACKUP_DIR/db-$DATE.sql.gz"
EOF

chmod +x /var/fider/scripts/backup.sh
```

Add to crontab:

```bash
crontab -e
# Add this line for daily 2 AM backups:
0 2 * * * /var/fider/scripts/backup.sh >> /var/fider/backups/backup.log 2>&1
```

### Volume Snapshots

Create OpenStack volume snapshots weekly via dashboard:
1. Navigate to Volumes → blazeboard-postgres-data
2. Create Snapshot
3. Name: `blazeboard-postgres-YYYYMMDD`

## Security Best Practices

- [ ] Strong passwords in `.env` (32+ characters)
- [ ] `.env` file has 600 permissions
- [ ] Regular security updates: `sudo apt update && sudo apt upgrade`
- [ ] Monitor logs for suspicious activity
- [ ] Regular database backups
- [ ] SSL/TLS enabled (HTTPS only)
- [ ] Security headers configured in nginx
- [ ] SAML authentication enabled for production

## Performance Tuning

### Increase Instance Resources

If performance is slow:
1. Create volume snapshot
2. Resize instance flavor in OpenStack (e.g., m1.large)
3. Restart instance
4. No data loss (volume persists)

### Database Optimization

```bash
# Check database size
docker compose exec db psql -U fider -c "SELECT pg_size_pretty(pg_database_size('fider'));"

# Vacuum database
docker compose exec db psql -U fider -c "VACUUM ANALYZE;"
```

## Disaster Recovery

### Complete Instance Failure

1. Create new instance from Ubuntu image
2. Run `setup-instance.sh`
3. Attach existing volume (`blazeboard-postgres-data`)
4. Deploy application files
5. Start services - data is preserved

### Data Corruption

1. Stop application: `docker compose down`
2. Restore from backup or volume snapshot
3. Restart application: `docker compose up -d`

## Support

- **UAB Research Computing**: https://docs.rc.uab.edu/help/support
- **Fider Documentation**: https://docs.fider.io
- **Deployment Design**: See `docs/plans/2026-01-31-cloud-rc-deployment-design.md`

## Additional Resources

- [UAB Cloud.rc Documentation](https://docs.rc.uab.edu/uab_cloud/)
- [Fider Self-Hosted Guide](https://docs.fider.io/self-hosted/)
- [Docker Compose Reference](https://docs.docker.com/compose/)
- [Let's Encrypt Documentation](https://letsencrypt.org/docs/)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
