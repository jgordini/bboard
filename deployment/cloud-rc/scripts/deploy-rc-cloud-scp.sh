#!/bin/bash
# =============================================================================
# Deploy BlazeBoard to RC Cloud using SCP (build locally, copy image)
# =============================================================================
# Builds the Docker image locally, saves it, syncs deployment files and image
# to the instance via rsync/scp, then loads the image and runs deploy.sh update.
#
# Usage:
#   ./deploy-rc-cloud-scp.sh ubuntu@138.26.48.197
#   RC_HOST=ubuntu@138.26.48.197 RC_SSH_KEY=~/.ssh/cloud_key ./deploy-rc-cloud-scp.sh
#
# Optional: RC_SSH_KEY path to SSH private key (e.g. ~/.ssh/cloud_key)
#
# Prerequisites:
#   - Docker installed locally (to build image)
#   - SSH access to the RC cloud instance
#   - rsync installed locally
# =============================================================================

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
DEPLOYMENT_DIR="$(dirname "$SCRIPT_DIR")"
REPO_ROOT="$(cd "$DEPLOYMENT_DIR/../.." && pwd)"
REMOTE_DIR="/var/fider"
IMAGE_NAME="blazeboard:latest"
IMAGE_TAR="blazeboard-latest.tar"

RC_HOST="${RC_HOST:-$1}"
RC_SSH_KEY="${RC_SSH_KEY:-}"
if [ -z "$RC_HOST" ]; then
    echo "Usage: $0 <user@host>"
    echo "   or: RC_HOST=ubuntu@138.26.48.197 RC_SSH_KEY=~/.ssh/cloud_key $0"
    echo ""
    echo "Examples:"
    echo "  $0 ubuntu@138.26.48.197"
    echo "  RC_HOST=ubuntu@blazeboard.cloud.rc.uab.edu RC_SSH_KEY=~/.ssh/cloud_key $0"
    exit 1
fi

if [ -n "$RC_SSH_KEY" ]; then
    RC_SSH_KEY="$(eval echo "$RC_SSH_KEY")"
    RSYNC_RSH="ssh -i $RC_SSH_KEY"
    SCP_OPTS=(-i "$RC_SSH_KEY")
    SSH_OPTS=(-i "$RC_SSH_KEY")
else
    RSYNC_RSH="ssh"
    SCP_OPTS=()
    SSH_OPTS=()
fi

echo "========================================================================"
echo "Deploying BlazeBoard to RC Cloud (build + SCP): $RC_HOST"
echo "========================================================================"
echo ""

echo "[1/5] Building Docker image locally ($IMAGE_NAME) for linux/amd64..."
cd "$REPO_ROOT"
docker build --platform linux/amd64 -t "$IMAGE_NAME" .

echo ""
echo "[2/5] Saving image to $IMAGE_TAR ..."
docker save "$IMAGE_NAME" -o "$DEPLOYMENT_DIR/$IMAGE_TAR"

echo ""
echo "[3/5] Syncing deployment files and image to $RC_HOST:$REMOTE_DIR ..."
RSYNC_RSH="$RSYNC_RSH" rsync -avz -e "$RSYNC_RSH" --exclude='.env' --exclude='.git' \
    "$DEPLOYMENT_DIR/" "$RC_HOST:$REMOTE_DIR/"

echo ""
echo "[4/5] Loading image and running deploy.sh update on remote..."
ssh "${SSH_OPTS[@]}" "$RC_HOST" "cd $REMOTE_DIR && docker load -i $IMAGE_TAR && rm -f $IMAGE_TAR && ./scripts/deploy.sh update"

# Remove local tar to avoid committing
rm -f "$DEPLOYMENT_DIR/$IMAGE_TAR"

echo ""
echo "[5/5] Done. BlazeBoard has been deployed to RC cloud."
echo ""
