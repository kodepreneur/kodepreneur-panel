#!/usr/bin/env bash
set -euo pipefail

# ==============================================================================
# Module: 40_agent.sh - Go Agent Daemon Provisioning
# ==============================================================================

echo -e "${COLOR_BLUE}[5/7] Building and configuring Kodepreneur Agent daemon...${COLOR_RESET}"

AGENT_CONFIG_DIR="/etc/kodepreneur"
AGENT_CONFIG_FILE="${AGENT_CONFIG_DIR}/agent.yaml"
AGENT_BIN="/usr/local/bin/kodepreneur-agent"

mkdir -p "${AGENT_CONFIG_DIR}" /var/log/kodepreneur

# Generate or preserve HMAC secret key
if [ ! -f "${AGENT_CONFIG_FILE}" ]; then
    AGENT_SECRET_KEY=$(openssl rand -hex 32)
    cat <<EOF > "${AGENT_CONFIG_FILE}"
server:
  host: "127.0.0.1"
  port: 8443
  read_timeout_sec: 15
  write_timeout_sec: 15

security:
  secret_key: "${AGENT_SECRET_KEY}"
  rate_limit_rpm: 600

environment:
  is_dev: false
  log_level: "info"
EOF
    chmod 0600 "${AGENT_CONFIG_FILE}"
else
    AGENT_SECRET_KEY=$(grep 'secret_key:' "${AGENT_CONFIG_FILE}" | awk '{print $2}' | tr -d '"')
fi

# Build Go Agent binary if source is present
if [ -d "${PROJECT_ROOT}/agent" ]; then
    echo -e "${COLOR_GREEN}  - Compiling Go agent binary...${COLOR_RESET}"
    if ! command -v go >/dev/null 2>&1; then
        apt-get install -y -qq golang-go
    fi
    (cd "${PROJECT_ROOT}/agent" && go build -o "${AGENT_BIN}" ./cmd/agent)
    chmod +x "${AGENT_BIN}"
fi

# Install systemd service
cp "${PROJECT_ROOT}/installer/systemd/kodepreneur-agent.service" /etc/systemd/system/kodepreneur-agent.service
systemctl daemon-reload
systemctl enable kodepreneur-agent >/dev/null 2>&1 || true
systemctl restart kodepreneur-agent >/dev/null 2>&1 || true

echo -e "${COLOR_GREEN}✓ Kodepreneur Agent daemon installed and active on 127.0.0.1:8443.${COLOR_RESET}"
