#!/usr/bin/env bash
set -euo pipefail

# ==============================================================================
# Kodepreneur Panel - Automated 1-Command Updater
# ==============================================================================

COLOR_GREEN='\033[0;32m'
COLOR_BLUE='\033[0;34m'
COLOR_RED='\033[0;31m'
COLOR_RESET='\033[0m'

if [ "${EUID:-$(id -u)}" -ne 0 ]; then
    echo -e "${COLOR_RED}Error: This update script must be run as root (use sudo).${COLOR_RESET}"
    exit 1
fi

REPO_GIT_URL="https://github.com/kodepreneur/kodepreneur-panel.git"
PANEL_DIR="/var/www/kodepreneur-panel"
AGENT_BIN="/usr/local/bin/kodepreneur-agent"

# Determine PROJECT_ROOT and INSTALLER_DIR
if [ -n "${BASH_SOURCE[0]:-}" ] && [ -f "${BASH_SOURCE[0]}" ]; then
    INSTALLER_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
    PROJECT_ROOT="$(cd "${INSTALLER_DIR}/.." && pwd)"
else
    INSTALLER_DIR="$(pwd)/installer"
    PROJECT_ROOT="$(pwd)"
fi

echo -e "${COLOR_BLUE}🚀 Starting Kodepreneur Panel update...${COLOR_RESET}"

# If executed via curl | bash or missing local sources, clone latest from GitHub
if [ ! -d "${PROJECT_ROOT}/panel" ] || [ ! -d "${PROJECT_ROOT}/agent" ]; then
    echo -e "${COLOR_BLUE}Fetching latest release from ${REPO_GIT_URL}...${COLOR_RESET}"
    TMP_UPDATE_DIR="/tmp/kodepreneur-panel-update"
    rm -rf "${TMP_UPDATE_DIR}"
    git clone --depth 1 "${REPO_GIT_URL}" "${TMP_UPDATE_DIR}"
    PROJECT_ROOT="${TMP_UPDATE_DIR}"
elif [ -d "${PROJECT_ROOT}/.git" ]; then
    echo -e "${COLOR_BLUE}Pulling latest changes from Git repository...${COLOR_RESET}"
    (cd "${PROJECT_ROOT}" && git pull --ff-only origin main 2>/dev/null || git pull 2>/dev/null || true)
fi


# 1. Update Laravel Control Plane
if [ -d "${PANEL_DIR}" ]; then
    echo -e "${COLOR_BLUE}[1/2] Updating Laravel Control Plane at ${PANEL_DIR}...${COLOR_RESET}"
    cd "${PANEL_DIR}"
    
    if [ -d "${PROJECT_ROOT}/panel" ] && [ "${PROJECT_ROOT}/panel" != "${PANEL_DIR}" ]; then
        cp -r "${PROJECT_ROOT}/panel/." "${PANEL_DIR}/"
    fi

    # Update dependencies & run migrations
    composer install --no-dev --prefer-dist --optimize-autoloader --no-interaction >/dev/null 2>&1 || true
    php artisan migrate --force >/dev/null 2>&1 || true
    
    # Rebuild frontend if npm exists
    if command -v npm >/dev/null 2>&1; then
        npm install --silent >/dev/null 2>&1 || true
        npm run build >/dev/null 2>&1 || true
    fi

    # Optimize caches
    php artisan optimize:clear >/dev/null 2>&1
    php artisan config:cache >/dev/null 2>&1
    php artisan route:cache >/dev/null 2>&1
    php artisan view:cache >/dev/null 2>&1

    chown -R www-data:www-data "${PANEL_DIR}"
    chmod -R 775 "${PANEL_DIR}/storage" "${PANEL_DIR}/bootstrap/cache" "${PANEL_DIR}/database"
    echo -e "${COLOR_GREEN}✓ Control plane updated and cached.${COLOR_RESET}"
fi

# 2. Update Go Agent Daemon
if [ -d "${PROJECT_ROOT}/agent" ]; then
    echo -e "${COLOR_BLUE}[2/2] Rebuilding Go Agent daemon...${COLOR_RESET}"
    if command -v go >/dev/null 2>&1; then
        (cd "${PROJECT_ROOT}/agent" && go build -o "${AGENT_BIN}" ./cmd/agent)
        chmod +x "${AGENT_BIN}"
        systemctl restart kodepreneur-agent
        echo -e "${COLOR_GREEN}✓ Go Agent recompiled and restarted.${COLOR_RESET}"
    fi
fi

# 3. Reload Nginx
nginx -t >/dev/null 2>&1 && systemctl reload nginx

echo ""
echo -e "${COLOR_GREEN}==============================================================================${COLOR_RESET}"
echo -e "${COLOR_GREEN}          🎉 Kodepreneur Panel has been updated successfully!                 ${COLOR_RESET}"
echo -e "${COLOR_GREEN}==============================================================================${COLOR_RESET}"
echo ""
