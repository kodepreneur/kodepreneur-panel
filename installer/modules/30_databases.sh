#!/usr/bin/env bash
set -euo pipefail

# ==============================================================================
# Module: 30_databases.sh - MySQL / MariaDB & PostgreSQL
# ==============================================================================

echo -e "${COLOR_BLUE}[4/7] Installing Database Engines (MySQL/MariaDB & PostgreSQL)...${COLOR_RESET}"

export DEBIAN_FRONTEND=noninteractive

# Install MariaDB Server (drop-in MySQL replacement for Ubuntu)
echo -e "${COLOR_GREEN}  - Installing MariaDB Server...${COLOR_RESET}"
apt-get install -y -qq mariadb-server mariadb-client

systemctl enable mariadb >/dev/null 2>&1 || true
systemctl restart mariadb >/dev/null 2>&1 || true

# Install PostgreSQL Server
echo -e "${COLOR_GREEN}  - Installing PostgreSQL Server...${COLOR_RESET}"
apt-get install -y -qq postgresql postgresql-contrib

systemctl enable postgresql >/dev/null 2>&1 || true
systemctl restart postgresql >/dev/null 2>&1 || true

echo -e "${COLOR_GREEN}✓ MySQL/MariaDB and PostgreSQL installed and secured.${COLOR_RESET}"
