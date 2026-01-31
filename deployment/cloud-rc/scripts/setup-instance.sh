#!/bin/bash
# =============================================================================
# UAB Cloud.rc Instance Setup Script
# =============================================================================
# This script prepares a fresh Ubuntu instance on cloud.rc for running bboard
#
# Run on the cloud.rc instance after first login:
#   chmod +x setup-instance.sh
#   ./setup-instance.sh
#
# =============================================================================

set -e  # Exit on error

echo "========================================================================"
echo "UAB Cloud.rc Instance Setup for BlazeBoard"
echo "========================================================================"

# Check if running as root
if [ "$EUID" -eq 0 ]; then
    echo "ERROR: Do not run this script as root. Run as ubuntu user."
    exit 1
fi

# -----------------------------------------------------------------------------
# 1. System Updates
# -----------------------------------------------------------------------------
echo ""
echo "[1/7] Updating system packages..."
sudo apt update
sudo apt upgrade -y

# -----------------------------------------------------------------------------
# 2. Install Docker
# -----------------------------------------------------------------------------
echo ""
echo "[2/7] Installing Docker..."
if ! command -v docker &> /dev/null; then
    curl -fsSL https://get.docker.com | sudo sh
    sudo usermod -aG docker ubuntu
    echo "Docker installed. You'll need to log out and back in for group membership."
else
    echo "Docker already installed."
fi

# -----------------------------------------------------------------------------
# 3. Install Docker Compose (if not included)
# -----------------------------------------------------------------------------
echo ""
echo "[3/7] Checking Docker Compose..."
if ! docker compose version &> /dev/null; then
    echo "Installing Docker Compose..."
    sudo apt install -y docker-compose-plugin
else
    echo "Docker Compose already available."
fi

# -----------------------------------------------------------------------------
# 4. Setup Persistent Volume
# -----------------------------------------------------------------------------
echo ""
echo "[4/7] Setting up persistent volume for PostgreSQL..."

# Detect volume device (usually /dev/vdb for first attached volume)
VOLUME_DEVICE="/dev/vdb"

if [ ! -b "$VOLUME_DEVICE" ]; then
    echo "WARNING: Volume device $VOLUME_DEVICE not found."
    echo "Available block devices:"
    lsblk
    echo ""
    read -p "Enter the correct device path (e.g., /dev/vdc): " VOLUME_DEVICE
fi

# Check if volume is already formatted
if sudo file -s "$VOLUME_DEVICE" | grep -q "filesystem"; then
    echo "Volume $VOLUME_DEVICE is already formatted."
else
    echo "Formatting $VOLUME_DEVICE as ext4..."
    sudo mkfs.ext4 "$VOLUME_DEVICE"
fi

# Create mount point
sudo mkdir -p /mnt/postgres-data

# Check if already mounted
if mountpoint -q /mnt/postgres-data; then
    echo "/mnt/postgres-data is already mounted."
else
    echo "Mounting $VOLUME_DEVICE to /mnt/postgres-data..."
    sudo mount "$VOLUME_DEVICE" /mnt/postgres-data
fi

# Add to fstab if not already present
if ! grep -q "/mnt/postgres-data" /etc/fstab; then
    echo "Adding mount to /etc/fstab..."
    echo "$VOLUME_DEVICE /mnt/postgres-data ext4 defaults 0 2" | sudo tee -a /etc/fstab
else
    echo "Mount already in /etc/fstab."
fi

# Set permissions for PostgreSQL container (UID 999)
sudo chown -R 999:999 /mnt/postgres-data

# -----------------------------------------------------------------------------
# 5. Create Application Directory
# -----------------------------------------------------------------------------
echo ""
echo "[5/7] Creating application directory..."
sudo mkdir -p /var/fider/{nginx,ssl}
sudo chown -R ubuntu:ubuntu /var/fider

# Create certbot directory
sudo mkdir -p /var/www/certbot

# -----------------------------------------------------------------------------
# 6. Install Certbot for SSL
# -----------------------------------------------------------------------------
echo ""
echo "[6/7] Installing Certbot..."
if ! command -v certbot &> /dev/null; then
    sudo apt install -y certbot
else
    echo "Certbot already installed."
fi

# -----------------------------------------------------------------------------
# 7. System Configuration
# -----------------------------------------------------------------------------
echo ""
echo "[7/7] Configuring system settings..."

# Increase file descriptors for Docker
if ! grep -q "fs.file-max" /etc/sysctl.conf; then
    echo "fs.file-max = 65536" | sudo tee -a /etc/sysctl.conf
    sudo sysctl -p
fi

echo ""
echo "========================================================================"
echo "Setup Complete!"
echo "========================================================================"
echo ""
echo "Next Steps:"
echo "1. Log out and back in for Docker group membership to take effect"
echo "2. Upload deployment files to /var/fider/"
echo "3. Create .env file from .env.template"
echo "4. Follow deployment guide to start application"
echo ""
echo "To verify Docker access after re-login:"
echo "  docker ps"
echo ""
echo "To check volume mount:"
echo "  df -h /mnt/postgres-data"
echo ""
