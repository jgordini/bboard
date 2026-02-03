# BlazeBoard Cloud.rc Deployment Checklist

Use this checklist to track your deployment progress.

## Pre-Deployment Preparation

### UAB IT Coordination

- [ ] Request cloud.rc account if you don't have one
- [ ] Use instance IP `138.26.48.197` (or request a domain name when ready)
- [ ] Request DNS A record pointing to floating IP (can do after getting IP)
- [ ] Request SMTP credentials for UAB mail server (optional - can use MailHog)
- [ ] Request SAML IdP certificate (optional - can add later)

### Local Preparation

- [ ] Generate SSH key pair if needed: `ssh-keygen -t ed25519`
- [ ] Review deployment files in `deployment/cloud-rc/`
- [ ] Read `QUICKSTART.md`
- [ ] Have terminal and browser ready

### Access Verification

- [ ] Can login to https://dashboard.cloud.rc.uab.edu
- [ ] On UAB network or connected to UAB VPN
- [ ] Have Duo 2FA set up

## Phase 1: OpenStack Resources (Dashboard)

### Create Volume

- [ ] Navigate to: Project → Volumes → Volumes → Create Volume
- [ ] Name: `blazeboard-postgres-data`
- [ ] Description: "PostgreSQL persistent storage for BlazeBoard"
- [ ] Size: `50` GB
- [ ] Volume Type: `__DEFAULT__`
- [ ] Availability Zone: (leave default)
- [ ] Click "Create Volume"
- [ ] Wait for status: "Available"

### Launch Instance

- [ ] Navigate to: Project → Compute → Instances → Launch Instance
- [ ] **Details Tab:**
  - [ ] Instance Name: `blazeboard-prod`
  - [ ] Description: "BlazeBoard production instance"
  - [ ] Availability Zone: (leave default)
  - [ ] Count: 1
  - [ ] Click Next
- [ ] **Source Tab:**
  - [ ] Select Boot Source: "Image"
  - [ ] Create New Volume: "No"
  - [ ] Select: "Ubuntu 24.04 LTS" (or latest Ubuntu)
  - [ ] Click Next
- [ ] **Flavor Tab:**
  - [ ] Select: `m1.medium` (2 vCPU, 4GB RAM) or larger
  - [ ] Click Next
- [ ] **Networks Tab:**
  - [ ] Select your project network
  - [ ] Click Next
- [ ] **Security Groups Tab:**
  - [ ] Select: `default`
  - [ ] Click Next
- [ ] **Key Pair Tab:**
  - [ ] Select your SSH key
  - [ ] If no key exists, click "Create Key Pair"
  - [ ] Click "Launch Instance"
- [ ] Wait for status: "Active"

### Attach Volume

- [ ] Navigate to: Project → Volumes → Volumes
- [ ] Find: `blazeboard-postgres-data`
- [ ] Actions → Manage Attachments
- [ ] Attach to Instance: Select `blazeboard-prod`
- [ ] Click "Attach Volume"
- [ ] Wait for status: "In-use"

### Configure Security Groups

- [ ] Navigate to: Project → Network → Security Groups
- [ ] Select: `default`
- [ ] Click "Manage Rules"
- [ ] Add Rule: SSH
  - [ ] Rule: SSH
  - [ ] Remote: CIDR
  - [ ] CIDR: `0.0.0.0/0` (or restrict to UAB IPs)
  - [ ] Click "Add"
- [ ] Add Rule: HTTP
  - [ ] Rule: HTTP
  - [ ] Remote: CIDR
  - [ ] CIDR: `0.0.0.0/0`
  - [ ] Click "Add"
- [ ] Add Rule: HTTPS
  - [ ] Rule: HTTPS
  - [ ] Remote: CIDR
  - [ ] CIDR: `0.0.0.0/0`
  - [ ] Click "Add"

### Allocate Floating IP

- [ ] Navigate to: Project → Network → Floating IPs
- [ ] Click "Allocate IP to Project"
- [ ] Pool: (select available pool)
- [ ] Description: "BlazeBoard production"
- [ ] Click "Allocate IP"
- [ ] Note the IP: `_________________________`
- [ ] Actions → Associate
- [ ] Port to be associated: Select `blazeboard-prod`
- [ ] Click "Associate"

### OpenStack Setup Complete!

- [ ] Instance Status: Active
- [ ] Volume Status: In-use
- [ ] Floating IP: Associated
- [ ] Security Groups: Configured
- [ ] Ready for SSH access

## Phase 2: Instance Setup

### SSH Access

- [ ] Test SSH: `ssh ubuntu@<floating-ip>`
- [ ] Connected successfully

### Upload Setup Script

From your local machine:

```bash
cd deployment/cloud-rc
scp scripts/setup-instance.sh ubuntu@<floating-ip>:~/
```

- [ ] File uploaded successfully

### Run Setup Script

On the cloud.rc instance:

```bash
chmod +x setup-instance.sh
./setup-instance.sh
```

- [ ] Script completed without errors
- [ ] Log out: `exit`
- [ ] Log back in: `ssh ubuntu@<floating-ip>`
- [ ] Verify Docker: `docker ps` (should work without sudo)

## Phase 3: Deploy Application Files

### Upload Deployment Files

From your local machine:

```bash
cd deployment/cloud-rc
rsync -avz --exclude='.env' ./ ubuntu@<floating-ip>:/var/fider/
```

Or using scp:

```bash
scp -r docker-compose.yml .env.template nginx scripts ubuntu@<floating-ip>:/var/fider/
```

- [ ] Files uploaded successfully
- [ ] Verify: `ssh ubuntu@<floating-ip> "ls -la /var/fider"`

## Phase 4: Configure Environment

### Generate Secrets

On your local machine:

```bash
# Generate and save these somewhere secure temporarily
openssl rand -base64 32  # DB_PASSWORD
openssl rand -base64 64  # JWT_SECRET
```

- [ ] DB_PASSWORD generated: `___________________________`
- [ ] JWT_SECRET generated: `___________________________`

### Create .env File

On the cloud.rc instance:

```bash
cd /var/fider
cp .env.template .env
nano .env
```

Edit these values:

- [ ] `BASE_URL=http://<floating-ip>` (use your floating IP)
- [ ] `DB_PASSWORD=` (paste generated password)
- [ ] `JWT_SECRET=` (paste generated secret)
- [ ] `EMAIL_NOREPLY=noreply@uab.edu`
- [ ] `SMTP_HOST=mailhog` (for testing)
- [ ] `SMTP_PORT=1025`
- [ ] Leave SMTP_USER and SMTP_PASS empty for MailHog
- [ ] Save file (Ctrl+X, Y, Enter)
- [ ] Secure file: `chmod 600 .env`

## Phase 5: Initial Deployment (HTTP)

### Deploy Application

```bash
cd /var/fider
./scripts/deploy.sh http
```

- [ ] Docker images pulled successfully
- [ ] Containers started
- [ ] No errors in output

### Verify Deployment

```bash
# Check container status
docker compose ps
```

- [ ] All containers show "Up" or "Up (healthy)"

```bash
# Check application logs
docker compose logs app | tail -20
```

- [ ] Look for: "http server started on :3000"
- [ ] No error messages

```bash
# Test local access
curl http://localhost
```

- [ ] Returns HTML (BlazeBoard page)

### Test Browser Access

- [ ] Open browser: `http://<floating-ip>`
- [ ] BlazeBoard setup page loads
- [ ] No errors in browser console

### MailHog Access (Optional)

- [ ] Open browser: `http://<floating-ip>:8025`
- [ ] MailHog web interface loads

## Phase 6: Complete Initial Setup

### BlazeBoard Setup Wizard

In browser at `http://<floating-ip>`:

- [ ] Enter organization/tenant name
- [ ] Enter admin email address
- [ ] Submit form
- [ ] Check MailHog (`http://<floating-ip>:8025`) for confirmation email
- [ ] Click confirmation link
- [ ] Set admin password
- [ ] Basic setup complete

## Phase 7: HTTPS Setup (Production)

### Request DNS Configuration

Contact UAB IT:

- [ ] Request DNS A record
- [ ] Host: `138.26.48.197`
- [ ] IP: `<your-floating-ip>`
- [ ] Wait for confirmation
- [ ] Test: `ping 138.26.48.197` (or `nslookup <domain>` when using a domain)

### Get SSL Certificate

On the cloud.rc instance:

```bash
sudo certbot certonly --webroot \
  -w /var/www/certbot \
  -d 138.26.48.197 \
  --email your-email@uab.edu \
  --agree-tos
```

- [ ] Certificate obtained successfully
- [ ] Files in `/etc/letsencrypt/live/138.26.48.197/`

### Update Configuration

```bash
cd /var/fider
nano .env
```

- [ ] Change `BASE_URL=https://138.26.48.197`
- [ ] Save file (Ctrl+X, Y, Enter)

### Deploy with HTTPS

```bash
./scripts/deploy.sh https
```

- [ ] Deployment successful
- [ ] No errors

### Verify HTTPS

- [ ] Open browser: `https://138.26.48.197`
- [ ] Page loads with valid SSL certificate
- [ ] No browser warnings
- [ ] Lock icon shows in address bar

### Setup Certificate Auto-Renewal

```bash
sudo crontab -e
```

Add line:

```
0 3 * * * certbot renew --quiet && docker compose -f /var/fider/docker-compose.yml restart nginx
```

- [ ] Cron job added
- [ ] Test renewal: `sudo certbot renew --dry-run`

## Phase 8: SAML Setup (Optional)

### Generate SAML Certificates

```bash
cd /var/fider
./scripts/generate-saml-certs.sh 138.26.48.197
```

- [ ] Certificates generated
- [ ] Files in `/var/fider/ssl/`

### Configure SAML in .env

```bash
nano .env
```

Add SAML configuration:

- [ ] Get SAML_IDP_CERT from UAB IT
- [ ] Uncomment SAML variables
- [ ] Add SAML_IDP_CERT value
- [ ] Save file

### Update docker-compose.yml

```bash
nano docker-compose.yml
```

- [ ] Uncomment SAML environment variables in app service
- [ ] Save file

### Redeploy

```bash
./scripts/deploy.sh https
```

- [ ] Deployment successful

### Download SP Metadata

```bash
curl https://138.26.48.197/saml/metadata > sp-metadata.xml
cat ssl/sp.crt
```

- [ ] sp-metadata.xml downloaded
- [ ] sp.crt contents copied

### Register with UAB IT

Send to UAB IT:

- [ ] sp-metadata.xml file
- [ ] sp.crt certificate
- [ ] Request SAML SP registration
- [ ] Wait for confirmation

### Test SAML Login

- [ ] Visit: `https://138.26.48.197`
- [ ] "Sign in with UAB" button visible
- [ ] Click button
- [ ] Redirects to UAB login
- [ ] Enter BlazerID credentials
- [ ] Duo authentication
- [ ] Redirects back to BlazeBoard
- [ ] Logged in successfully

## Phase 9: Production Readiness

### Setup Backups

```bash
cd /var/fider
nano scripts/backup.sh
```

Create backup script (see README.md)

- [ ] Backup script created
- [ ] Script executable: `chmod +x scripts/backup.sh`
- [ ] Test backup: `./scripts/backup.sh`
- [ ] Add to crontab for daily backups

### Create Volume Snapshot

In OpenStack Dashboard:

- [ ] Navigate to: Volumes → blazeboard-postgres-data
- [ ] Create Snapshot
- [ ] Name: `blazeboard-postgres-initial-<date>`
- [ ] Snapshot created

### Document Access

Create operations document:

- [ ] Floating IP address
- [ ] Domain name
- [ ] Admin email
- [ ] Backup location
- [ ] Contact information

### Test Recovery Procedures

- [ ] Test database backup/restore
- [ ] Test container restart
- [ ] Test volume snapshot restore
- [ ] Document findings

## Phase 10: Monitoring & Maintenance

### Setup Monitoring

- [ ] Create health check script
- [ ] Test email notifications
- [ ] Document monitoring procedures

### Security Review

- [ ] Verify .env has 600 permissions
- [ ] Verify no secrets in git
- [ ] Review security group rules
- [ ] Review nginx security headers
- [ ] Update system packages: `sudo apt update && sudo apt upgrade`

### Documentation

- [ ] Document custom configurations
- [ ] Create runbook for operations
- [ ] Train team members
- [ ] Update contact information

### Go Live

- [ ] Announce to users
- [ ] Monitor for issues
- [ ] Collect feedback
- [ ] Plan improvements

## Post-Deployment Checklist

### Regular Maintenance (Weekly)

- [ ] Check application logs
- [ ] Monitor disk usage
- [ ] Review backup status
- [ ] Check for updates

### Regular Maintenance (Monthly)

- [ ] System updates
- [ ] Database optimization
- [ ] Volume snapshot
- [ ] Review security

## Troubleshooting Reference

If issues occur, see:

- [ ] `README.md` - Troubleshooting section
- [ ] `docker compose logs [service]` - Container logs
- [ ] `docs/plans/2026-01-31-cloud-rc-deployment-design.md` - Architecture
- [ ] UAB RC Support: https://docs.rc.uab.edu/help/support

## Deployment Complete! 🎉

- [ ] Application accessible at production URL
- [ ] HTTPS working with valid certificate
- [ ] SAML authentication functional (if configured)
- [ ] Backups configured
- [ ] Monitoring in place
- [ ] Documentation updated
- [ ] Team trained
- [ ] Users notified

**Deployment Date:** **********\_\_\_**********
**Deployed By:** **********\_\_\_**********
**URL:** https://138.26.48.197
**Status:** ☐ Production ☐ Staging ☐ Testing
