#!/usr/bin/env bash
set -euo pipefail

# ==============================================================================
# Kodepreneur Panel - Master Production Installer
# Supported OS: Ubuntu 24.04 LTS, Ubuntu 22.04 LTS
# ==============================================================================

# ANSI Color Codes
COLOR_RED='\033[0;31m'
COLOR_GREEN='\033[0;32m'
COLOR_YELLOW='\033[1;33m'
COLOR_BLUE='\033[0;34m'
COLOR_CYAN='\033[0;36m'
COLOR_WHITE='\033[1;37m'
COLOR_RESET='\033[0m'

export COLOR_RED COLOR_GREEN COLOR_YELLOW COLOR_BLUE COLOR_CYAN COLOR_WHITE COLOR_RESET

# Banner
clear || true
echo -e "${COLOR_CYAN}"
cat << 'EOF'
  _  ______  _____  ______ _____  _____  ______ _   _ ______ _    _ _____  
 | |/ / __ \|  __ \|  ____|  __ \|  __ \|  ____| \ | |  ____| |  | |  __ \ 
 | ' / |  | | |  | | |__  | |__) | |__) | |__  |  \| | |__  | |  | | |__) |
 |  <| |  | | |  | |  __| |  ___/|  _  /|  __| | . ` |  __| | |  | |  _  / 
 | . \ |__| | |__| | |____| |    | | \ \| |____| |\  | |____| |__| | | \ \ 
 |_|\_\____/|_____/|______|_|    |_|  \_\______|_| \_|______|\____/|_|  \_\
                                                                            
             High-Performance Cloud Web Server Management Panel
EOF
echo -e "${COLOR_RESET}"

# Verify Root Privileges
if [ "${EUID:-$(id -u)}" -ne 0 ]; then
    echo -e "${COLOR_RED}Error: This installer must be run as root (use sudo).${COLOR_RESET}"
    exit 1
fi

# Detect OS & Version
if [ ! -f /etc/os-release ]; then
    echo -e "${COLOR_RED}Error: Cannot detect Linux distribution (/etc/os-release missing).${COLOR_RESET}"
    exit 1
fi

. /etc/os-release
if [ "${ID:-}" != "ubuntu" ]; then
    echo -e "${COLOR_YELLOW}Warning: Kodepreneur Panel is tailored for Ubuntu. Detected: ${ID:-unknown}.${COLOR_RESET}"
fi

# Default Variables
ADMIN_EMAIL="admin@kodepreneur.com"
ADMIN_PASSWORD=""
PANEL_PORT="8080"
UNATTENDED=false

# Determine PROJECT_ROOT
INSTALLER_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${INSTALLER_DIR}/.." && pwd)"
export PROJECT_ROOT

# Parse Arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --email)
            ADMIN_EMAIL="$2"
            shift 2
            ;;
        --password)
            ADMIN_PASSWORD="$2"
            shift 2
            ;;
        --port)
            PANEL_PORT="$2"
            shift 2
            ;;
        --unattended|-y)
            UNATTENDED=true
            shift
            ;;
        --help|-h)
            echo "Usage: ./install.sh [OPTIONS]"
            echo ""
            echo "Options:"
            echo "  --email <email>       Admin email address (default: admin@kodepreneur.com)"
            echo "  --password <pass>     Admin password (auto-generated if omitted)"
            echo "  --port <port>         Panel HTTP web port (default: 8080)"
            echo "  --unattended, -y      Run non-interactively"
            echo "  --help, -h            Show this help dialog"
            exit 0
            ;;
        *)
            echo -e "${COLOR_RED}Unknown option: $1${COLOR_RESET}"
            exit 1
            ;;
    esac
done

# Detect Public / Server IP
SERVER_IP=$(curl -s4 --max-time 3 ifconfig.me || hostname -I | awk '{print $1}' || echo "127.0.0.1")
export SERVER_IP

# Interactive prompts if not unattended
if [ "$UNATTENDED" = false ]; then
    echo -e "${COLOR_WHITE}Installer Configuration:${COLOR_RESET}"
    read -r -p "  Enter Admin Email [${ADMIN_EMAIL}]: " input_email
    if [ -n "$input_email" ]; then ADMIN_EMAIL="$input_email"; fi

    read -r -p "  Enter Panel Port [${PANEL_PORT}]: " input_port
    if [ -n "$input_port" ]; then PANEL_PORT="$input_port"; fi

    if [ -z "$ADMIN_PASSWORD" ]; then
        read -r -s -p "  Enter Admin Password (leave blank to auto-generate): " input_pass
        echo ""
        if [ -n "$input_pass" ]; then
            ADMIN_PASSWORD="$input_pass"
        fi
    fi
fi

# Auto-generate password if still empty
if [ -z "$ADMIN_PASSWORD" ]; then
    ADMIN_PASSWORD=$(openssl rand -base64 12 | tr -dc 'a-zA-Z0-9!@#$%^&*()_+' | head -c 16)
fi

export ADMIN_EMAIL ADMIN_PASSWORD PANEL_PORT

echo ""
echo -e "${COLOR_BLUE}Starting installation with:${COLOR_RESET}"
echo -e "  - Server IP:    ${COLOR_WHITE}${SERVER_IP}${COLOR_RESET}"
echo -e "  - Panel Port:   ${COLOR_WHITE}${PANEL_PORT}${COLOR_RESET}"
echo -e "  - Admin Email:  ${COLOR_WHITE}${ADMIN_EMAIL}${COLOR_RESET}"
echo ""

# Execution Trap
trap 'echo -e "\n${COLOR_RED}❌ Installation failed on line $LINENO. Check output above for details.${COLOR_RESET}"; exit 1' ERR

# Run Modules Sequentially
source "${INSTALLER_DIR}/modules/00_system.sh"
source "${INSTALLER_DIR}/modules/10_php.sh"
source "${INSTALLER_DIR}/modules/20_nginx.sh"
source "${INSTALLER_DIR}/modules/30_databases.sh"
source "${INSTALLER_DIR}/modules/40_agent.sh"
source "${INSTALLER_DIR}/modules/50_panel.sh"
source "${INSTALLER_DIR}/modules/99_summary.sh"
