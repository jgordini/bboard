# BlazeBoard Cloud.rc Quick Start

Get BlazeBoard running on UAB Cloud.rc in 6 steps.

## Before You Start

- [ ] Cloud.rc account access
- [ ] SSH key added to OpenStack
- [ ] Host or IP (e.g., `138.26.48.197`)

## Step 1: Create OpenStack Resources (5 minutes)

In the cloud.rc dashboard (https://dashboard.cloud.rc.uab.edu):

1. **Create Volume**: Project → Volumes → Create Volume

   - Name: `blazeboard-postgres-data`
   - Size: 50 GB

2. **Launch Instance**: Project → Compute → Launch Instance

   - Name: `blazeboard-prod`
   - Source: Ubuntu 24.04 LTS
   - Flavor: m1.medium (2 vCPU, 4GB RAM)
   - Network: Select your project network
   - Key Pair: Your SSH key

3. **Attach Volume**: Volumes → blazeboard-postgres-data → Manage Attachments

   - Attach to: `blazeboard-prod`

4. **Security Groups**: Project → Network → Security Groups → default → Manage Rules

   - Add: TCP 22 (SSH)
   - Add: TCP 80 (HTTP)
   - Add: TCP 443 (HTTPS)

5. **Get Floating IP**: Network → Floating IPs → Allocate IP
   - Associate with: `blazeboard-prod`
   - Note the IP address: `_________________`

## Step 2: Setup Instance (10 minutes)

```bash
# SSH into instance
ssh ubuntu@<floating-ip>

# Upload setup script (from your local machine in another terminal)
# In deployment/cloud-rc directory:
scp scripts/setup-instance.sh ubuntu@<floating-ip>:~/

# Back on the instance, run setup
chmod +x setup-instance.sh
./setup-instance.sh

# Log out and back in
exit
ssh ubuntu@<floating-ip>

# Verify Docker
docker ps
```

## Step 3: Deploy Files (2 minutes)

From your local machine:

```bash
cd deployment/cloud-rc

# Upload all deployment files
rsync -avz --exclude='.env' ./ ubuntu@<floating-ip>:/var/fider/

# Or use scp
scp -r docker-compose.yml .env.template nginx scripts ubuntu@<floating-ip>:/var/fider/
```

## Step 4: Configure Environment (3 minutes)

On the cloud.rc instance:

```bash
cd /var/fider

# Create .env from template
cp .env.template .env

# Generate secrets and save them
openssl rand -base64 32  # DB_PASSWORD - copy this
openssl rand -base64 64  # JWT_SECRET - copy this

# Edit .env
nano .env

# Update these required values:
# - BASE_URL=http://<your-floating-ip>  (use http for now)
# - DB_PASSWORD=<paste generated password>
# - JWT_SECRET=<paste generated secret>
# - SMTP_HOST=mailhog  (for testing)
# - SMTP_PORT=1025

# Save (Ctrl+X, Y, Enter)

# Secure the file
chmod 600 .env
```

## Step 5: Start Application (2 minutes)

```bash
cd /var/fider

# Deploy with HTTP (for initial setup)
./scripts/deploy.sh http

# Wait 30 seconds for containers to start

# Verify it's running
docker compose ps
docker compose logs app | tail -20

# Test access
curl http://localhost
```

## Step 6: Access BlazeBoard

Open your browser and go to: `http://<floating-ip>`

You should see the BlazeBoard setup wizard. Complete it:

1. Enter organization name
2. Enter admin email
3. Check email for confirmation link (if using MailHog, go to `http://<floating-ip>:8025`)

## Next Steps (Optional)

### Enable HTTPS

1. **Configure DNS**: Get UAB IT to create A record pointing to your floating IP
2. **Get SSL certificate**:
   ```bash
   sudo certbot certonly --webroot \
     -w /var/www/certbot \
     -d 138.26.48.197 \
     --email your-email@uab.edu \
     --agree-tos
   ```
3. **Update .env**: Change `BASE_URL` to `https://138.26.48.197`
4. **Redeploy**: `./scripts/deploy.sh https`

### Enable UAB SAML

1. **Generate certificates**: `./scripts/generate-saml-certs.sh 138.26.48.197`
2. **Get IdP cert from UAB IT**
3. **Update .env**: Add SAML configuration
4. **Update docker-compose.yml**: Uncomment SAML environment variables
5. **Redeploy**: `./scripts/deploy.sh https`
6. **Register with UAB**: Send SP metadata to UAB IT

## Troubleshooting

**Can't SSH into instance?**

- Check security group has port 22 open
- Verify you're using correct SSH key
- Try from UAB VPN if required

**Docker permission denied?**

- Did you log out and back in after setup?
- Check: `groups` should show `docker`

**Application not starting?**

- Check logs: `docker compose logs app`
- Verify .env has all required values
- Check database: `docker compose logs db`

**Can't access via browser?**

- Check security group has port 80 open
- Verify containers are running: `docker compose ps`
- Check nginx logs: `docker compose logs nginx`

## Useful Commands

```bash
# View logs
docker compose logs -f app

# Restart application
docker compose restart app

# Stop everything
docker compose down

# Start everything
docker compose up -d

# Check status
docker compose ps

# Backup database
docker compose exec db pg_dump -U fider fider > backup.sql
```

## Support

Need help?

- Full documentation: See `README.md` in this directory
- Deployment design: See `docs/plans/2026-01-31-cloud-rc-deployment-design.md`
- UAB RC Support: https://docs.rc.uab.edu/help/support
