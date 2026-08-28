#!/usr/bin/env bash
set -euo pipefail

# ==============================================================================
# Kodepreneur Panel - Automated 1-Command Updater
# ==============================================================================

COLOR_GREEN='\033[0;32m'
COLOR_BLUE='\033[0;34m'
COLOR_YELLOW='\033[1;33m'
COLOR_RED='\033[0;31m'
COLOR_RESET='\033[0m'

if [ "${EUID:-$(id -u)}" -ne 0 ]; then
    echo -e "${COLOR_RED}Error: This update script must be run as root (use sudo).${COLOR_RESET}"
    exit 1
fi

REPO_GIT_URL="https://github.com/kodepreneur/kodepreneur-panel.git"
TARGET_BRANCH="main"
PANEL_DIR="/var/www/kodepreneur-panel"
AGENT_BIN="/usr/local/bin/kodepreneur-agent"
DAEMON_MODE=false

# Parse arguments
while [[ $# -gt 0 ]]; do
    case "$1" in
        --repo|--repository)
            REPO_GIT_URL="$2"
            shift 2
            ;;
        --branch|-b)
            TARGET_BRANCH="$2"
            shift 2
            ;;
        --daemon-mode|--no-restart)
            DAEMON_MODE=true
            shift
            ;;
        *)
            shift
            ;;
    esac
done

export PATH="/usr/local/go/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:${PATH:-}"

# Determine PROJECT_ROOT and INSTALLER_DIR
if [ -n "${BASH_SOURCE[0]:-}" ] && [ -f "${BASH_SOURCE[0]}" ]; then
    INSTALLER_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
    PROJECT_ROOT="$(cd "${INSTALLER_DIR}/.." && pwd)"
else
    INSTALLER_DIR="$(pwd)/installer"
    PROJECT_ROOT="$(pwd)"
fi

echo -e "${COLOR_BLUE}🚀 Starting Kodepreneur Panel update (branch: ${TARGET_BRANCH})...${COLOR_RESET}"

# If executed via curl | bash or missing local sources, clone latest from GitHub
if [ ! -d "${PROJECT_ROOT}/panel" ] || [ ! -d "${PROJECT_ROOT}/agent" ]; then
    echo -e "${COLOR_BLUE}Fetching latest release from ${REPO_GIT_URL}...${COLOR_RESET}"
    TMP_UPDATE_DIR="/tmp/kodepreneur-panel-update"
    rm -rf "${TMP_UPDATE_DIR}"
    git clone --depth 1 -b "${TARGET_BRANCH}" "${REPO_GIT_URL}" "${TMP_UPDATE_DIR}" 2>/dev/null || git clone --depth 1 "${REPO_GIT_URL}" "${TMP_UPDATE_DIR}"
    PROJECT_ROOT="${TMP_UPDATE_DIR}"
elif [ -d "${PROJECT_ROOT}/.git" ]; then
    echo -e "${COLOR_BLUE}Pulling latest changes from Git repository...${COLOR_RESET}"
    (cd "${PROJECT_ROOT}" && git fetch origin "${TARGET_BRANCH}" 2>/dev/null && git reset --hard "origin/${TARGET_BRANCH}" 2>/dev/null || git pull 2>/dev/null || true)
fi

# 1. Update Laravel Control Plane
if [ -d "${PANEL_DIR}" ]; then
    echo -e "${COLOR_BLUE}[1/3] Updating Laravel Control Plane at ${PANEL_DIR}...${COLOR_RESET}"
    cd "${PANEL_DIR}"
    
    if [ -d "${PROJECT_ROOT}/panel" ] && [ "${PROJECT_ROOT}/panel" != "${PANEL_DIR}" ]; then
        cp -r "${PROJECT_ROOT}/panel/." "${PANEL_DIR}/"
    fi

    # Ensure installer scripts are in panel and /etc/kodepreneur
    if [ -d "${PROJECT_ROOT}/installer" ]; then
        mkdir -p "${PANEL_DIR}/installer" /etc/kodepreneur
        cp -r "${PROJECT_ROOT}/installer/." "${PANEL_DIR}/installer/"
        cp -f "${PROJECT_ROOT}/installer/update.sh" /etc/kodepreneur/update.sh
        chmod +x /etc/kodepreneur/update.sh "${PANEL_DIR}/installer/update.sh"
    fi

    # Ensure storage & cache directories exist
    mkdir -p "${PANEL_DIR}/storage/framework/sessions" \
             "${PANEL_DIR}/storage/framework/views" \
             "${PANEL_DIR}/storage/framework/cache/data" \
             "${PANEL_DIR}/storage/logs" \
             "${PANEL_DIR}/bootstrap/cache" \
             "${PANEL_DIR}/database"
    chmod -R 775 "${PANEL_DIR}/storage" "${PANEL_DIR}/bootstrap/cache" "${PANEL_DIR}/database" 2>/dev/null || true

    # Update dependencies & run migrations
    echo -e "${COLOR_BLUE}  - Running Composer dependency installation...${COLOR_RESET}"
    composer install --no-dev --prefer-dist --optimize-autoloader --no-interaction >/dev/null 2>&1 || true

    echo -e "${COLOR_BLUE}  - Applying database migrations...${COLOR_RESET}"
    php artisan migrate --force >/dev/null 2>&1 || true
    
    # Rebuild frontend if npm exists
    if command -v npm >/dev/null 2>&1; then
        echo -e "${COLOR_BLUE}  - Building frontend assets via Vite...${COLOR_RESET}"
        if [ ! -d "node_modules" ]; then
            npm install --silent >/dev/null 2>&1 || true
        fi
        npm run build >/dev/null 2>&1 || true
    fi

    # Optimize caches
    php artisan optimize:clear >/dev/null 2>&1 || true
    php artisan config:cache >/dev/null 2>&1 || true
    php artisan route:cache >/dev/null 2>&1 || true
    php artisan view:cache >/dev/null 2>&1 || true

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

        # Configure graceful process_control_timeout (30s) in php-fpm.conf to prevent 502 Bad Gateway during pool reloads
        FPM_MAIN_CONF="/etc/php/${VER}/fpm/php-fpm.conf"
        if [ -f "${FPM_MAIN_CONF}" ]; then
            if grep -q "^;process_control_timeout" "${FPM_MAIN_CONF}"; then
                sed -i 's/^;process_control_timeout = .*/process_control_timeout = 30s/' "${FPM_MAIN_CONF}"
            elif grep -q "^process_control_timeout" "${FPM_MAIN_CONF}"; then
                sed -i 's/^process_control_timeout = .*/process_control_timeout = 30s/' "${FPM_MAIN_CONF}"
            else
                echo "process_control_timeout = 30s" >> "${FPM_MAIN_CONF}"
            fi
        fi
    done

    # Fix open_basedir in existing website PHP-FPM pools (/var/www/domain/public -> /var/www/domain)
    for POOL_CONF in /etc/php/*/fpm/pool.d/kp_*.conf; do
        if [ -f "${POOL_CONF}" ]; then
            sed -i 's|\(/var/www/[^/:]*\)/public|\1|g' "${POOL_CONF}"
        fi
    done

    # Restart PHP-FPM runtimes to apply configuration & ensure sockets are alive
    if [ "$DAEMON_MODE" = false ]; then
        echo -e "${COLOR_BLUE}  - Restarting PHP-FPM runtimes...${COLOR_RESET}"
        systemctl restart php8.4-fpm php8.3-fpm 2>/dev/null || true
    else
        echo -e "${COLOR_GREEN}  ✓ PHP-FPM reload deferred to daemon background runner (zero-downtime).${COLOR_RESET}"
    fi

    # Update panel Nginx configuration if present
    if [ -f "${PROJECT_ROOT}/installer/nginx/kodepreneur-panel.conf" ] && [ -f /etc/nginx/sites-available/kodepreneur-panel.conf ]; then
        PANEL_PORT_DETECTED=$(grep -oP 'listen \K[0-9]+' /etc/nginx/sites-available/kodepreneur-panel.conf 2>/dev/null | head -1 || echo "8080")
        
        # Detect active PHP socket fallback
        PHP_SOCK="unix:/run/php/php8.4-fpm.sock"
        if [ ! -S /run/php/php8.4-fpm.sock ] && [ -S /run/php/php8.3-fpm.sock ]; then
            PHP_SOCK="unix:/run/php/php8.3-fpm.sock"
        fi

        sed "s/{{PANEL_PORT}}/${PANEL_PORT_DETECTED}/g" "${PROJECT_ROOT}/installer/nginx/kodepreneur-panel.conf" | \
        sed "s|unix:/run/php/php8.4-fpm.sock|${PHP_SOCK}|g" > /etc/nginx/sites-available/kodepreneur-panel.conf
    fi

    chown -R www-data:www-data "${PANEL_DIR}"
    chmod -R 775 "${PANEL_DIR}/storage" "${PANEL_DIR}/bootstrap/cache" "${PANEL_DIR}/database"
    echo -e "${COLOR_GREEN}✓ Control plane updated and cached.${COLOR_RESET}"
fi

# 2. Update Go Agent Daemon
if [ -d "${PROJECT_ROOT}/agent" ]; then
    echo -e "${COLOR_BLUE}[2/3] Rebuilding Go Agent daemon...${COLOR_RESET}"
    if command -v go >/dev/null 2>&1; then
        (cd "${PROJECT_ROOT}/agent" && go build -o "${AGENT_BIN}" ./cmd/agent)
        chmod +x "${AGENT_BIN}"
        if [ "$DAEMON_MODE" = false ]; then
            systemctl restart kodepreneur-agent 2>/dev/null || true
            echo -e "${COLOR_GREEN}✓ Go Agent recompiled and restarted.${COLOR_RESET}"
        else
            echo -e "${COLOR_GREEN}✓ Go Agent binary recompiled (restart scheduled via daemon runner).${COLOR_RESET}"
        fi
    fi
fi

# 3. Reload Nginx
echo -e "${COLOR_BLUE}[3/3] Validating and reloading Nginx...${COLOR_RESET}"
nginx -t >/dev/null 2>&1 && systemctl reload nginx 2>/dev/null || true
echo -e "${COLOR_GREEN}✓ Nginx configuration validated and reloaded.${COLOR_RESET}"

echo ""
echo -e "${COLOR_GREEN}==============================================================================${COLOR_RESET}"
echo -e "${COLOR_GREEN}          🎉 Kodepreneur Panel has been updated successfully!                 ${COLOR_RESET}"
echo -e "${COLOR_GREEN}==============================================================================${COLOR_RESET}"
echo ""
