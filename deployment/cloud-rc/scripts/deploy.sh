#!/bin/bash
# =============================================================================
# BlazeBoard Deployment Script for UAB Cloud.rc
# =============================================================================
# Deploys or updates the bboard application
#
# Usage:
#   ./deploy.sh          # Deploy with HTTPS (requires SSL certificates)
#   ./deploy.sh http     # Deploy HTTP-only (for initial setup)
#   ./deploy.sh update   # Update to latest version
#
# =============================================================================

set -e  # Exit on error

DEPLOYMENT_DIR="/var/fider"
MODE="${1:-https}"

echo "========================================================================"
echo "BlazeBoard Deployment - Mode: $MODE"
echo "========================================================================"

# Change to deployment directory
cd "$DEPLOYMENT_DIR"

# -----------------------------------------------------------------------------
# Validate Environment
# -----------------------------------------------------------------------------
echo ""
echo "[1/5] Validating environment..."

if [ ! -f .env ]; then
    echo "ERROR: .env file not found."
    echo "Please create .env from .env.template and configure it."
    exit 1
fi

# Source environment for validation
set -a
source .env
set +a

# Check required variables
REQUIRED_VARS=("BASE_URL" "DB_PASSWORD" "JWT_SECRET" "SMTP_HOST")
for var in "${REQUIRED_VARS[@]}"; do
    if [ -z "${!var}" ]; then
        echo "ERROR: Required variable $var is not set in .env"
        exit 1
    fi
done

echo "Environment validated."

# -----------------------------------------------------------------------------
# Configure Nginx Based on Mode
# -----------------------------------------------------------------------------
echo ""
echo "[2/5] Configuring nginx for $MODE mode..."

case "$MODE" in
    http)
        echo "Using HTTP-only configuration..."
        if [ ! -f nginx/nginx-http-only.conf ]; then
            echo "ERROR: nginx/nginx-http-only.conf not found."
            exit 1
        fi
        cp nginx/nginx-http-only.conf nginx/nginx.conf
        ;;
    https)
        echo "Using HTTPS configuration..."
        # Extract domain from BASE_URL
        DOMAIN=$(echo "$BASE_URL" | sed -e 's|^https\?://||' -e 's|/.*$||')

        # Check if SSL certificates exist
        if [ ! -f "/etc/letsencrypt/live/$DOMAIN/fullchain.pem" ]; then
            echo "ERROR: SSL certificates not found for $DOMAIN"
            echo "Run: sudo certbot certonly --webroot -w /var/www/certbot -d $DOMAIN"
            exit 1
        fi
        echo "SSL certificates found for $DOMAIN"
        ;;
    update)
        echo "Update mode - keeping existing nginx configuration..."
        ;;
    *)
        echo "ERROR: Unknown mode '$MODE'"
        echo "Usage: $0 [http|https|update]"
        exit 1
        ;;
esac

# -----------------------------------------------------------------------------
# Pull Latest Images
# -----------------------------------------------------------------------------
echo ""
echo "[3/5] Pulling Docker images..."
docker compose pull

# -----------------------------------------------------------------------------
# Deploy Application
# -----------------------------------------------------------------------------
echo ""
echo "[4/5] Deploying application..."

if [ "$MODE" = "update" ]; then
    echo "Performing rolling update..."
    docker compose up -d --no-deps app
else
    echo "Starting all services..."
    docker compose up -d
fi

# -----------------------------------------------------------------------------
# Verify Deployment
# -----------------------------------------------------------------------------
echo ""
echo "[5/5] Verifying deployment..."

# Wait for services to be healthy
echo "Waiting for services to start..."
sleep 5

# Check container status
echo ""
echo "Container Status:"
docker compose ps

# Check application logs
echo ""
echo "Recent application logs:"
docker compose logs --tail=20 app

# -----------------------------------------------------------------------------
# Deployment Summary
# -----------------------------------------------------------------------------
echo ""
echo "========================================================================"
echo "Deployment Complete!"
echo "========================================================================"
echo ""
echo "Service Status:"
docker compose ps --format "table {{.Name}}\t{{.Status}}\t{{.Ports}}"
echo ""

case "$MODE" in
    http)
        echo "Application accessible at: http://<your-floating-ip>"
        echo ""
        echo "Next Steps:"
        echo "1. Test the application: curl http://<your-floating-ip>"
        echo "2. Get SSL certificate: sudo certbot certonly --webroot -w /var/www/certbot -d <your-domain>"
        echo "3. Re-deploy with HTTPS: ./deploy.sh https"
        ;;
    https)
        echo "Application accessible at: $BASE_URL"
        echo ""
        echo "Next Steps:"
        echo "1. Complete initial setup at $BASE_URL"
        echo "2. Configure SAML if needed (see deployment guide)"
        echo "3. Setup backup cron job (see operations guide)"
        ;;
    update)
        echo "Application updated at: $BASE_URL"
        ;;
esac

echo ""
echo "Useful Commands:"
echo "  View logs:     docker compose logs -f [app|db|nginx]"
echo "  Restart:       docker compose restart [service]"
echo "  Stop:          docker compose down"
echo "  Backup DB:     docker compose exec db pg_dump -U fider fider > backup.sql"
echo ""
