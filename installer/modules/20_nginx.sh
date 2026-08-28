#!/usr/bin/env bash
set -euo pipefail

# ==============================================================================
# Module: 20_nginx.sh - Nginx Web Server & Certbot
# ==============================================================================

echo -e "${COLOR_BLUE}[3/7] Installing Nginx and Certbot...${COLOR_RESET}"

export DEBIAN_FRONTEND=noninteractive

# Install Nginx and Certbot
apt-get install -y -qq \
    nginx \
    certbot \
    python3-certbot-nginx

# Ensure standard directory hierarchy
mkdir -p /etc/nginx/sites-available /etc/nginx/sites-enabled /var/www /etc/nginx/ssl

# Enable sites-enabled in nginx.conf if not already included
if ! grep -q "sites-enabled" /etc/nginx/nginx.conf; then
    sed -i '/http {/a \    include /etc/nginx/sites-enabled/*;' /etc/nginx/nginx.conf
fi

# Remove default site if present
if [ -f /etc/nginx/sites-enabled/default ]; then
    rm -f /etc/nginx/sites-enabled/default
fi

# Validate Nginx syntax & restart
nginx -t >/dev/null 2>&1
systemctl enable nginx >/dev/null 2>&1 || true
systemctl restart nginx >/dev/null 2>&1 || true

echo -e "${COLOR_GREEN}✓ Nginx and Certbot installed and running.${COLOR_RESET}"
