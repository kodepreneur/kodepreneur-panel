#!/usr/bin/env bash
set -euo pipefail

# ==============================================================================
# Module: 50_panel.sh - Laravel Control Plane Deployment
# ==============================================================================

echo -e "${COLOR_BLUE}[6/7] Deploying Kodepreneur Control Plane...${COLOR_RESET}"

PANEL_DIR="/var/www/kodepreneur-panel"
mkdir -p "${PANEL_DIR}"

# Copy panel codebase
if [ -d "${PROJECT_ROOT}/panel" ] && [ "${PROJECT_ROOT}/panel" != "${PANEL_DIR}" ]; then
    echo -e "${COLOR_GREEN}  - Copying application files to ${PANEL_DIR}...${COLOR_RESET}"
    cp -r "${PROJECT_ROOT}/panel/." "${PANEL_DIR}/"
fi

# Ensure installer update scripts are preserved locally
if [ -d "${PROJECT_ROOT}/installer" ]; then
    mkdir -p "${PANEL_DIR}/installer" /etc/kodepreneur
    cp -r "${PROJECT_ROOT}/installer/." "${PANEL_DIR}/installer/"
    cp -f "${PROJECT_ROOT}/installer/update.sh" /etc/kodepreneur/update.sh
    chmod +x /etc/kodepreneur/update.sh "${PANEL_DIR}/installer/update.sh"
fi

cd "${PANEL_DIR}"

# Create SQLite database file
mkdir -p "${PANEL_DIR}/database"
touch "${PANEL_DIR}/database/database.sqlite"

# Create .env if missing
if [ ! -f "${PANEL_DIR}/.env" ]; then
    echo -e "${COLOR_GREEN}  - Generating production environment configuration...${COLOR_RESET}"
    cat <<EOF > "${PANEL_DIR}/.env"
APP_NAME=Kodepreneur
APP_ENV=production
APP_KEY=
APP_DEBUG=false
APP_TIMEZONE=UTC
APP_URL=http://${SERVER_IP}:${PANEL_PORT}

APP_LOCALE=en
APP_FALLBACK_LOCALE=en
APP_FAKER_LOCALE=en_US

APP_MAINTENANCE_DRIVER=file

BCRYPT_ROUNDS=12

LOG_CHANNEL=stack
LOG_STACK=single
LOG_DEPRECATIONS_CHANNEL=null
LOG_LEVEL=info

DB_CONNECTION=sqlite
DB_DATABASE=${PANEL_DIR}/database/database.sqlite

SESSION_DRIVER=file
SESSION_LIFETIME=120
SESSION_ENCRYPT=false
SESSION_PATH=/
SESSION_DOMAIN=null

BROADCAST_CONNECTION=log
FILESYSTEM_DISK=local
QUEUE_CONNECTION=database

CACHE_STORE=file

AGENT_BASE_URL=http://127.0.0.1:8443
AGENT_SECRET_KEY=${AGENT_SECRET_KEY}
AGENT_MOCK=false
EOF
fi

# Ensure storage & cache directories exist
mkdir -p "${PANEL_DIR}/storage/framework/sessions" \
         "${PANEL_DIR}/storage/framework/views" \
         "${PANEL_DIR}/storage/framework/cache/data" \
         "${PANEL_DIR}/storage/logs" \
         "${PANEL_DIR}/bootstrap/cache" \
         "${PANEL_DIR}/database"
chmod -R 775 "${PANEL_DIR}/storage" "${PANEL_DIR}/bootstrap/cache" "${PANEL_DIR}/database"

# Install PHP dependencies
echo -e "${COLOR_GREEN}  - Installing Composer production dependencies...${COLOR_RESET}"
composer install --no-dev --prefer-dist --optimize-autoloader --no-interaction >/dev/null 2>&1

# Generate APP_KEY if empty
if ! grep -q "APP_KEY=base64:" "${PANEL_DIR}/.env"; then
    php artisan key:generate --force >/dev/null 2>&1
fi

# Run database migrations and seed default records
echo -e "${COLOR_GREEN}  - Running database migrations and default seeders...${COLOR_RESET}"
php artisan migrate --force --seed >/dev/null 2>&1

# Provision initial admin user
php artisan tinker --execute="
\$role = \App\Models\Role::firstOrCreate(['slug' => 'super-admin'], ['name' => 'Super Administrator', 'permissions' => ['*']]);
\$u = \App\Models\User::updateOrCreate(
    ['email' => getenv('ADMIN_EMAIL')],
    [
        'name' => 'Administrator',
        'password' => \Illuminate\Support\Facades\Hash::make(getenv('ADMIN_PASSWORD')),
        'role_id' => \$role->id,
    ]
);
" >/dev/null 2>&1

# Build frontend assets if npm is installed and node_modules not present
if [ ! -d "${PANEL_DIR}/public/build" ] || [ ! -f "${PANEL_DIR}/public/build/manifest.json" ]; then
    echo -e "${COLOR_GREEN}  - Compiling frontend assets via Vite...${COLOR_RESET}"
    npm install --silent >/dev/null 2>&1
    npm run build >/dev/null 2>&1
fi

# Cache Laravel configuration & routes
php artisan optimize:clear >/dev/null 2>&1
php artisan config:cache >/dev/null 2>&1
php artisan route:cache >/dev/null 2>&1
php artisan view:cache >/dev/null 2>&1

# Deploy Nginx Virtual Host for the Panel
echo -e "${COLOR_GREEN}  - Configuring Nginx panel virtual host on port ${PANEL_PORT}...${COLOR_RESET}"
sed "s/{{PANEL_PORT}}/${PANEL_PORT}/g" "${PROJECT_ROOT}/installer/nginx/kodepreneur-panel.conf" > /etc/nginx/sites-available/kodepreneur-panel.conf
ln -sf /etc/nginx/sites-available/kodepreneur-panel.conf /etc/nginx/sites-enabled/kodepreneur-panel.conf

# Set permissions
chown -R www-data:www-data "${PANEL_DIR}"
chmod -R 775 "${PANEL_DIR}/storage" "${PANEL_DIR}/bootstrap/cache" "${PANEL_DIR}/database"

# Restart Nginx
nginx -t >/dev/null 2>&1
systemctl restart nginx

echo -e "${COLOR_GREEN}✓ Kodepreneur Control Plane successfully deployed.${COLOR_RESET}"
