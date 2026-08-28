#!/usr/bin/env bash
set -euo pipefail

# ==============================================================================
# Module: 00_system.sh - Base Packages & Firewall
# ==============================================================================

echo -e "${COLOR_BLUE}[1/7] Configuring system packages and security firewall...${COLOR_RESET}"

export DEBIAN_FRONTEND=noninteractive

# Update repository lists
apt-get update -y -qq

# Install essential base tools
apt-get install -y -qq \
    curl \
    wget \
    git \
    unzip \
    zip \
    tar \
    software-properties-common \
    ca-certificates \
    gnupg \
    lsb-release \
    jq \
    htop \
    ufw \
    fail2ban \
    build-essential

# Configure UFW Firewall
echo -e "${COLOR_GREEN}  - Configuring UFW rules (SSH:22, HTTP:80, HTTPS:443, Panel:${PANEL_PORT})...${COLOR_RESET}"
ufw default deny incoming >/dev/null 2>&1 || true
ufw default allow outgoing >/dev/null 2>&1 || true
ufw allow 22/tcp comment 'SSH' >/dev/null 2>&1 || true
ufw allow 80/tcp comment 'HTTP' >/dev/null 2>&1 || true
ufw allow 443/tcp comment 'HTTPS' >/dev/null 2>&1 || true
ufw allow "${PANEL_PORT}/tcp" comment 'Kodepreneur Panel' >/dev/null 2>&1 || true
ufw --force enable >/dev/null 2>&1 || true

# Setup Swap if memory is constrained (< 2GB) and no swap exists
SWAP_TOTAL=$(free -m 2>/dev/null | awk '/^Swap:/ {print $2}' || echo "0")
MEM_TOTAL=$(free -m 2>/dev/null | awk '/^Mem:/ {print $2}' || echo "2048")
if [ "${SWAP_TOTAL:-0}" -eq 0 ] && [ "${MEM_TOTAL:-2048}" -lt 2048 ]; then
    echo -e "${COLOR_GREEN}  - Configuring 1GB swap space for system stability...${COLOR_RESET}"
    fallocate -l 1G /swapfile 2>/dev/null || dd if=/dev/zero of=/swapfile bs=1M count=1024 2>/dev/null || true
    chmod 600 /swapfile 2>/dev/null || true
    mkswap /swapfile >/dev/null 2>&1 || true
    swapon /swapfile >/dev/null 2>&1 || true
    if ! grep -q "/swapfile" /etc/fstab 2>/dev/null; then
        echo "/swapfile swap swap defaults 0 0" >> /etc/fstab
    fi
fi

# Start & enable Fail2ban
systemctl enable fail2ban >/dev/null 2>&1 || true
systemctl restart fail2ban >/dev/null 2>&1 || true

echo -e "${COLOR_GREEN}✓ System packages, UFW firewall, and Fail2ban configured.${COLOR_RESET}"
