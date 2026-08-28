#!/usr/bin/env bash
set -euo pipefail

# ==============================================================================
# Module: 10_php.sh - PHP 8.3 & PHP 8.4 Runtimes + Node.js & Composer
# ==============================================================================

echo -e "${COLOR_BLUE}[2/7] Installing PHP runtimes, Composer, and Node.js...${COLOR_RESET}"

export DEBIAN_FRONTEND=noninteractive

# Add Ondřej Surý PHP PPA repository
if ! grep -rq "ondrej/php" /etc/apt/sources.list /etc/apt/sources.list.d/ 2>/dev/null; then
    echo -e "${COLOR_GREEN}  - Adding Ondrej PHP PPA repository...${COLOR_RESET}"
    if ! LC_ALL=C.UTF-8 add-apt-repository -y ppa:ondrej/php >/dev/null 2>&1; then
        echo -e "${COLOR_YELLOW}    Warning: add-apt-repository returned non-zero, configuring GPG keyring via HTTPS...${COLOR_RESET}"
        mkdir -p /etc/apt/keyrings
        curl -fsSL "https://keyserver.ubuntu.com/pks/lookup?op=get&search=0x14AA40EC0831756756D7F66C4F4EA0AAE5267A6C" | gpg --dearmor -o /etc/apt/keyrings/ondrej-php.gpg --yes 2>/dev/null || true
        UBUNTU_CODENAME="$(. /etc/os-release && echo "${VERSION_CODENAME:-noble}")"
        echo "deb [signed-by=/etc/apt/keyrings/ondrej-php.gpg] https://ppa.launchpadcontent.net/ondrej/php/ubuntu ${UBUNTU_CODENAME} main" > /etc/apt/sources.list.d/ondrej-php.list
    fi
    apt-get update -y -qq || true
fi

# PHP Extensions list
PHP_EXTENSIONS="fpm cli common mysql pgsql mbstring xml curl zip bcmath intl gd sqlite3 readline opcache"

# Install PHP 8.3 (LTS)
echo -e "${COLOR_GREEN}  - Installing PHP 8.3 runtime...${COLOR_RESET}"
apt-get install -y -qq \
    php8.3 \
    php8.3-fpm \
    php8.3-cli \
    php8.3-common \
    php8.3-mysql \
    php8.3-pgsql \
    php8.3-mbstring \
    php8.3-xml \
    php8.3-curl \
    php8.3-zip \
    php8.3-bcmath \
    php8.3-intl \
    php8.3-gd \
    php8.3-sqlite3 \
    php8.3-readline \
    php8.3-opcache

# Install PHP 8.4 (Latest)
echo -e "${COLOR_GREEN}  - Installing PHP 8.4 runtime...${COLOR_RESET}"
apt-get install -y -qq \
    php8.4 \
    php8.4-fpm \
    php8.4-cli \
    php8.4-common \
    php8.4-mysql \
    php8.4-pgsql \
    php8.4-mbstring \
    php8.4-xml \
    php8.4-curl \
    php8.4-zip \
    php8.4-bcmath \
    php8.4-intl \
    php8.4-gd \
    php8.4-sqlite3 \
    php8.4-readline \
    php8.4-opcache

# Configure high-performance upload & execution limits in php.ini
for VER in 8.3 8.4; do
    for CONF in /etc/php/${VER}/fpm/php.ini /etc/php/${VER}/cli/php.ini; do
        if [ -f "${CONF}" ]; then
            sed -i 's/^upload_max_filesize = .*/upload_max_filesize = 512M/' "${CONF}"
            sed -i 's/^post_max_size = .*/post_max_size = 512M/' "${CONF}"
            sed -i 's/^memory_limit = .*/memory_limit = 512M/' "${CONF}"
            sed -i 's/^max_execution_time = .*/max_execution_time = 600/' "${CONF}"
            sed -i 's/^max_input_time = .*/max_input_time = 600/' "${CONF}"
        fi
    done
done

# Ensure FPM services are running
systemctl enable php8.3-fpm php8.4-fpm >/dev/null 2>&1 || true
systemctl restart php8.3-fpm php8.4-fpm >/dev/null 2>&1 || true

# Install Composer globally
if ! command -v composer >/dev/null 2>&1; then
    echo -e "${COLOR_GREEN}  - Installing Composer globally...${COLOR_RESET}"
    curl -sS https://getcomposer.org/installer | php -- --install-dir=/usr/local/bin --filename=composer >/dev/null 2>&1
fi

# Install Node.js 22.x LTS (NodeSource)
if ! command -v node >/dev/null 2>&1; then
    echo -e "${COLOR_GREEN}  - Installing Node.js LTS...${COLOR_RESET}"
    curl -fsSL https://deb.nodesource.com/setup_22.x | bash - >/dev/null 2>&1
    apt-get install -y -qq nodejs
fi

echo -e "${COLOR_GREEN}✓ PHP 8.3, PHP 8.4, Composer, and Node.js installed.${COLOR_RESET}"
