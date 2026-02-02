#!/bin/bash
# =============================================================================
# Update BlazeBoard Docker containers on UAB Cloud.rc
# =============================================================================
# Run this from your local machine (not on the RC cloud instance).
# Syncs deployment files to the instance and runs deploy.sh update.
#
# Usage:
#   ./update-rc-cloud.sh ubuntu@blazeboard.cloud.rc.uab.edu
#   RC_HOST=ubuntu@138.26.48.197 RC_SSH_KEY=~/.ssh/cloud_key ./update-rc-cloud.sh
#
# Optional: RC_SSH_KEY path to SSH private key (e.g. ~/.ssh/cloud_key)
#
# Prerequisites:
#   - SSH access to the RC cloud instance (key or agent)
#   - rsync installed locally
# =============================================================================

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
DEPLOYMENT_DIR="$(dirname "$SCRIPT_DIR")"
REMOTE_DIR="/var/fider"

RC_HOST="${RC_HOST:-$1}"
RC_SSH_KEY="${RC_SSH_KEY:-}"
if [ -z "$RC_HOST" ]; then
    echo "Usage: $0 <user@host>"
    echo "   or: RC_HOST=ubuntu@blazeboard.cloud.rc.uab.edu $0"
    echo ""
    echo "Examples:"
    echo "  $0 ubuntu@blazeboard.cloud.rc.uab.edu"
    echo "  RC_HOST=ubuntu@138.26.48.197 RC_SSH_KEY=~/.ssh/cloud_key $0"
    exit 1
fi

if [ -n "$RC_SSH_KEY" ]; then
    RC_SSH_KEY="$(eval echo "$RC_SSH_KEY")"
    RSYNC_RSH="ssh -i $RC_SSH_KEY"
    SSH_OPTS=(-i "$RC_SSH_KEY")
else
    RSYNC_RSH="ssh"
    SSH_OPTS=()
fi

echo "========================================================================"
echo "Updating BlazeBoard on RC Cloud: $RC_HOST"
echo "========================================================================"
echo ""
echo "[1/3] Syncing deployment files to $RC_HOST:$REMOTE_DIR ..."
RSYNC_RSH="$RSYNC_RSH" rsync -avz -e "$RSYNC_RSH" --exclude='.env' --exclude='.git' --exclude='nginx/nginx.conf' \
    "$DEPLOYMENT_DIR/" "$RC_HOST:$REMOTE_DIR/"

echo ""
echo "[2/3] Running deploy.sh update on remote..."
ssh "${SSH_OPTS[@]}" "$RC_HOST" "cd $REMOTE_DIR && ./scripts/deploy.sh update"

echo ""
echo "[3/3] Done. Containers on RC cloud have been updated."
echo ""
