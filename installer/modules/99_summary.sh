#!/usr/bin/env bash
set -euo pipefail

# ==============================================================================
# Module: 99_summary.sh - Installation Summary & Credentials Banner
# ==============================================================================

echo ""
echo -e "${COLOR_GREEN}==============================================================================${COLOR_RESET}"
echo -e "${COLOR_GREEN}           🎉 Kodepreneur Panel Installation Completed Successfully!          ${COLOR_RESET}"
echo -e "${COLOR_GREEN}==============================================================================${COLOR_RESET}"
echo ""
echo -e "  🌐 ${COLOR_WHITE}Panel URL:${COLOR_RESET}         ${COLOR_CYAN}http://${SERVER_IP}:${PANEL_PORT}${COLOR_RESET}"
echo -e "  👤 ${COLOR_WHITE}Admin Email:${COLOR_RESET}       ${COLOR_YELLOW}${ADMIN_EMAIL}${COLOR_RESET}"
echo -e "  🔑 ${COLOR_WHITE}Admin Password:${COLOR_RESET}    ${COLOR_YELLOW}${ADMIN_PASSWORD}${COLOR_RESET}"
echo ""
echo -e "  🛡️  ${COLOR_WHITE}Agent Daemon:${COLOR_RESET}      ${COLOR_GREEN}Active on 127.0.0.1:8443 (HMAC Protected)${COLOR_RESET}"
echo -e "  🚀 ${COLOR_WHITE}PHP Runtimes:${COLOR_RESET}      ${COLOR_GREEN}PHP 8.3 LTS & PHP 8.4 Latest${COLOR_RESET}"
echo -e "  🗄️  ${COLOR_WHITE}Database Engines:${COLOR_RESET}  ${COLOR_GREEN}MariaDB (MySQL) & PostgreSQL${COLOR_RESET}"
echo ""
echo -e "${COLOR_BLUE}------------------------------------------------------------------------------${COLOR_RESET}"
echo -e "  Useful Management Commands:"
echo -e "  - View Agent status:      ${COLOR_WHITE}systemctl status kodepreneur-agent${COLOR_RESET}"
echo -e "  - View Agent logs:        ${COLOR_WHITE}journalctl -u kodepreneur-agent -f${COLOR_RESET}"
echo -e "  - Restart Nginx:          ${COLOR_WHITE}systemctl restart nginx${COLOR_RESET}"
echo -e "  - Panel CLI:              ${COLOR_WHITE}cd /var/www/kodepreneur-panel && php artisan${COLOR_RESET}"
echo -e "${COLOR_BLUE}------------------------------------------------------------------------------${COLOR_RESET}"
echo ""
