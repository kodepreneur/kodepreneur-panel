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

# Start & enable Fail2ban
systemctl enable fail2ban >/dev/null 2>&1 || true
systemctl restart fail2ban >/dev/null 2>&1 || true

echo -e "${COLOR_GREEN}✓ System packages, UFW firewall, and Fail2ban configured.${COLOR_RESET}"
