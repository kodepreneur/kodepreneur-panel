<?php

namespace App\Services\Agent;

class MockAgentClient implements AgentClientInterface
{
    public function health(): array
    {
        return [
            'status' => 'healthy',
            'version' => '1.0.0-mock',
            'dev' => true,
        ];
    }

    public function getSystemInfo(): array
    {
        return [
            'hostname' => gethostname() ?: 'kodepreneur-dev-vps',
            'os' => 'Ubuntu 24.04 LTS (Noble Numbat)',
            'kernel' => '6.8.0-31-generic',
            'architecture' => 'x86_64',
            'public_ip' => '203.0.113.42',
            'uptime_seconds' => 174200,
        ];
    }

    public function getSystemMetrics(): array
    {
        $cpuPct = round(15.4 + sin(time() / 10) * 8.2, 1);
        if ($cpuPct < 2) $cpuPct = 2.0;

        $totalMem = 2 * 1024 * 1024 * 1024; // 2GB
        $usedMem = 780 * 1024 * 1024 + (int)(sin(time() / 5) * 50 * 1024 * 1024);
        $freeMem = $totalMem - $usedMem;

        $totalDisk = 50 * 1024 * 1024 * 1024; // 50GB
        $usedDisk = 12 * 1024 * 1024 * 1024;  // 12GB
        $freeDisk = $totalDisk - $usedDisk;

        return [
            'cpu' => [
                'usage_percent' => $cpuPct,
                'cores' => 2,
            ],
            'memory' => [
                'total_bytes' => $totalMem,
                'used_bytes' => $usedMem,
                'free_bytes' => $freeMem,
                'usage_percent' => round(($usedMem / $totalMem) * 100, 1),
            ],
            'disk' => [
                'mount' => '/',
                'total_bytes' => $totalDisk,
                'used_bytes' => $usedDisk,
                'free_bytes' => $freeDisk,
                'usage_percent' => round(($usedDisk / $totalDisk) * 100, 1),
            ],
            'load_average' => [
                'load1' => 0.28,
                'load5' => 0.21,
                'load15' => 0.15,
            ],
            'timestamp' => time(),
        ];
    }

    public function getServiceStatus(string $service): array
    {
        return [
            'name' => $service,
            'is_active' => true,
            'is_enabled' => true,
            'status' => "active (running) [mock]",
        ];
    }

    public function manageService(string $service, string $action): array
    {
        return [
            'service' => $service,
            'action' => $action,
            'output' => "Mock executed systemctl {$action} {$service} successfully",
        ];
    }

    public function createWebsite(array $payload): array
    {
        $domain = $payload['domain'] ?? 'example.com';
        $docRoot = $payload['document_root'] ?? "/var/www/{$domain}/public";
        $isLaravel = ($payload['project_type'] ?? '') === 'laravel' || str_ends_with($docRoot, '/public');

        return [
            'domain' => $domain,
            'vhost_path' => "/etc/nginx/sites-available/{$domain}.conf",
            'system_user' => $payload['system_user'] ?? 'kp_user',
            'document_root' => $docRoot,
            'php_version' => $payload['php_version'] ?? '8.3',
            'deployment_source' => $payload['deployment_source'] ?? 'empty',
            'is_laravel' => $isLaravel,
        ];
    }

    public function deleteWebsite(string $domain, array $payload = []): array
    {
        return [
            'domain' => $domain,
            'status' => 'deleted',
        ];
    }

    public function switchPhpVersion(string $domain, array $payload): array
    {
        return [
            'domain' => $domain,
            'php_version' => $payload['new_php_version'] ?? '8.4',
        ];
    }

    public function issueSsl(array $payload): array
    {
        $domain = $payload['domain'] ?? 'example.com';
        return [
            'domain' => $domain,
            'cert_path' => "/etc/letsencrypt/live/{$domain}/fullchain.pem",
            'key_path' => "/etc/letsencrypt/live/{$domain}/privkey.pem",
            'issuer' => "Let's Encrypt (Mock)",
            'valid_from' => date('Y-m-d H:i:s'),
            'valid_until' => date('Y-m-d H:i:s', strtotime('+90 days')),
            'status' => 'valid',
        ];
    }

    public function getWebsiteLogs(string $domain, string $type = 'access', int $lines = 100): array
    {
        $now = date('d/M/Y:H:i:s O');
        $sampleLines = [];
        if ($type === 'access') {
            for ($i = 0; $i < 10; $i++) {
                $sampleLines[] = "127.0.0.1 - - [{$now}] \"GET /api/v1/test HTTP/1.1\" 200 1420 \"https://{$domain}\" \"Mozilla/5.0\"";
            }
        } else {
            $sampleLines[] = date('Y/m/d H:i:s') . " [notice] 1001#1001: using the epoll event method for {$domain}";
            $sampleLines[] = date('Y/m/d H:i:s') . " [notice] 1001#1001: worker process started";
        }

        return [
            'domain' => $domain,
            'type' => $type,
            'lines' => $sampleLines,
        ];
    }

    public function createDatabase(array $payload): array
    {
        return [
            'engine' => $payload['engine'] ?? 'mysql',
            'name' => $payload['name'] ?? 'app_db',
        ];
    }

    public function deleteDatabase(string $engine, string $name): array
    {
        return [
            'success' => true,
            'message' => "Database {$name} dropped successfully",
        ];
    }

    public function createDatabaseUser(array $payload): array
    {
        return [
            'engine' => $payload['engine'] ?? 'mysql',
            'username' => $payload['username'] ?? 'app_user',
            'host' => $payload['host'] ?? 'localhost',
        ];
    }

    public function deleteDatabaseUser(string $engine, string $username, string $host = 'localhost'): array
    {
        return [
            'success' => true,
            'message' => "User {$username} dropped successfully",
        ];
    }

    public function grantDatabaseAccess(array $payload): array
    {
        return [
            'success' => true,
            'message' => "Granted permissions successfully",
        ];
    }

    public function resetDatabaseUserPassword(array $payload): array
    {
        return [
            'success' => true,
            'message' => "Password reset successfully",
        ];
    }

    public function getDatabaseTables(string $engine, string $name): array
    {
        if (strtolower($engine) === 'postgresql' || strtolower($engine) === 'postgres') {
            return [
                ['name' => 'accounts', 'type' => 'BASE TABLE', 'engine' => 'PostgreSQL', 'rows' => 5400, 'data_length' => 524288, 'index_length' => 131072, 'total_size' => 655360, 'collation' => 'UTF8', 'comment' => 'Organization billing accounts'],
                ['name' => 'subscriptions', 'type' => 'BASE TABLE', 'engine' => 'PostgreSQL', 'rows' => 3200, 'data_length' => 262144, 'index_length' => 65536, 'total_size' => 327680, 'collation' => 'UTF8', 'comment' => 'SaaS subscriptions'],
                ['name' => 'invoices', 'type' => 'BASE TABLE', 'engine' => 'PostgreSQL', 'rows' => 15400, 'data_length' => 2097152, 'index_length' => 524288, 'total_size' => 2621440, 'collation' => 'UTF8', 'comment' => 'Customer invoices'],
                ['name' => 'audit_events', 'type' => 'BASE TABLE', 'engine' => 'PostgreSQL', 'rows' => 84000, 'data_length' => 8388608, 'index_length' => 2097152, 'total_size' => 10485760, 'collation' => 'UTF8', 'comment' => 'System audit logs'],
            ];
        }

        return [
            ['name' => 'users', 'type' => 'BASE TABLE', 'engine' => 'InnoDB', 'rows' => 1420, 'data_length' => 163840, 'index_length' => 65536, 'total_size' => 229376, 'collation' => 'utf8mb4_unicode_ci', 'comment' => 'System and app users'],
            ['name' => 'orders', 'type' => 'BASE TABLE', 'engine' => 'InnoDB', 'rows' => 8930, 'data_length' => 1048576, 'index_length' => 393216, 'total_size' => 1441792, 'collation' => 'utf8mb4_unicode_ci', 'comment' => 'Customer orders'],
            ['name' => 'order_items', 'type' => 'BASE TABLE', 'engine' => 'InnoDB', 'rows' => 24500, 'data_length' => 2097152, 'index_length' => 786432, 'total_size' => 2883584, 'collation' => 'utf8mb4_unicode_ci', 'comment' => 'Items per order'],
            ['name' => 'products', 'type' => 'BASE TABLE', 'engine' => 'InnoDB', 'rows' => 340, 'data_length' => 98304, 'index_length' => 32768, 'total_size' => 131072, 'collation' => 'utf8mb4_unicode_ci', 'comment' => 'Product catalog'],
            ['name' => 'settings', 'type' => 'BASE TABLE', 'engine' => 'InnoDB', 'rows' => 28, 'data_length' => 16384, 'index_length' => 16384, 'total_size' => 32768, 'collation' => 'utf8mb4_unicode_ci', 'comment' => 'Configuration settings'],
            ['name' => 'migrations', 'type' => 'BASE TABLE', 'engine' => 'InnoDB', 'rows' => 18, 'data_length' => 16384, 'index_length' => 0, 'total_size' => 16384, 'collation' => 'utf8mb4_unicode_ci', 'comment' => 'Database migration ledger'],
        ];
    }

    public function getTableStructure(string $engine, string $name, string $table): array
    {
        if (strtolower($engine) === 'postgresql' || strtolower($engine) === 'postgres') {
            return [
                'table_name' => $table,
                'columns' => [
                    ['name' => 'id', 'position' => 1, 'type' => 'uuid', 'data_type' => 'uuid', 'is_nullable' => false, 'key' => 'PRI', 'default' => null, 'extra' => '', 'comment' => 'Primary Key UUID'],
                    ['name' => 'account_id', 'position' => 2, 'type' => 'uuid', 'data_type' => 'uuid', 'is_nullable' => false, 'key' => 'MUL', 'default' => null, 'extra' => '', 'comment' => 'Parent account FK'],
                    ['name' => 'status', 'position' => 3, 'type' => 'character varying(32)', 'data_type' => 'varchar', 'is_nullable' => false, 'key' => '', 'default' => "'active'::text", 'extra' => '', 'comment' => 'Current status'],
                    ['name' => 'metadata', 'position' => 4, 'type' => 'jsonb', 'data_type' => 'jsonb', 'is_nullable' => true, 'key' => '', 'default' => null, 'extra' => '', 'comment' => 'Custom JSON attributes'],
                    ['name' => 'created_at', 'position' => 5, 'type' => 'timestamp with time zone', 'data_type' => 'timestamp', 'is_nullable' => false, 'key' => '', 'default' => 'now()', 'extra' => '', 'comment' => 'Creation timestamp'],
                ],
                'indexes' => [
                    ['name' => "{$table}_pkey", 'column' => 'id', 'non_unique' => false, 'is_primary' => true, 'seq_in_index' => 1, 'type' => 'btree'],
                    ['name' => "{$table}_account_id_idx", 'column' => 'account_id', 'non_unique' => true, 'is_primary' => false, 'seq_in_index' => 1, 'type' => 'btree'],
                ],
                'foreign_keys' => [
                    ['constraint_name' => "{$table}_account_id_fkey", 'column' => 'account_id', 'referenced_table' => 'accounts', 'referenced_column' => 'id'],
                ],
                'create_statement' => "CREATE TABLE public.{$table} (\n    id uuid NOT NULL,\n    account_id uuid NOT NULL,\n    status character varying(32) NOT NULL DEFAULT 'active'::text,\n    metadata jsonb,\n    created_at timestamp with time zone NOT NULL DEFAULT now(),\n    CONSTRAINT {$table}_pkey PRIMARY KEY (id)\n);",
            ];
        }

        if ($table === 'orders') {
            return [
                'table_name' => 'orders',
                'columns' => [
                    ['name' => 'id', 'position' => 1, 'type' => 'bigint unsigned', 'data_type' => 'bigint', 'is_nullable' => false, 'key' => 'PRI', 'default' => null, 'extra' => 'auto_increment', 'comment' => 'Primary Key'],
                    ['name' => 'user_id', 'position' => 2, 'type' => 'bigint unsigned', 'data_type' => 'bigint', 'is_nullable' => false, 'key' => 'MUL', 'default' => null, 'extra' => '', 'comment' => 'Customer FK'],
                    ['name' => 'order_number', 'position' => 3, 'type' => 'varchar(64)', 'data_type' => 'varchar', 'is_nullable' => false, 'key' => 'UNI', 'default' => null, 'extra' => '', 'comment' => 'Invoice Ref'],
                    ['name' => 'total_amount', 'position' => 4, 'type' => 'decimal(12,2)', 'data_type' => 'decimal', 'is_nullable' => false, 'key' => '', 'default' => '0.00', 'extra' => '', 'comment' => 'Total in USD'],
                    ['name' => 'status', 'position' => 5, 'type' => 'varchar(32)', 'data_type' => 'varchar', 'is_nullable' => false, 'key' => '', 'default' => null, 'extra' => '', 'comment' => 'Order Status'],
                    ['name' => 'created_at', 'position' => 6, 'type' => 'timestamp', 'data_type' => 'timestamp', 'is_nullable' => true, 'key' => '', 'default' => 'CURRENT_TIMESTAMP', 'extra' => '', 'comment' => ''],
                    ['name' => 'updated_at', 'position' => 7, 'type' => 'timestamp', 'data_type' => 'timestamp', 'is_nullable' => true, 'key' => '', 'default' => 'CURRENT_TIMESTAMP', 'extra' => 'on update CURRENT_TIMESTAMP', 'comment' => ''],
                ],
                'indexes' => [
                    ['name' => 'PRIMARY', 'column' => 'id', 'non_unique' => false, 'is_primary' => true, 'seq_in_index' => 1, 'type' => 'BTREE'],
                    ['name' => 'orders_order_number_unique', 'column' => 'order_number', 'non_unique' => false, 'is_primary' => false, 'seq_in_index' => 1, 'type' => 'BTREE'],
                    ['name' => 'orders_user_id_index', 'column' => 'user_id', 'non_unique' => true, 'is_primary' => false, 'seq_in_index' => 1, 'type' => 'BTREE'],
                ],
                'foreign_keys' => [
                    ['constraint_name' => 'orders_user_id_foreign', 'column' => 'user_id', 'referenced_table' => 'users', 'referenced_column' => 'id'],
                ],
                'create_statement' => "CREATE TABLE `orders` (\n  `id` bigint unsigned NOT NULL AUTO_INCREMENT,\n  `user_id` bigint unsigned NOT NULL,\n  `order_number` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL,\n  `total_amount` decimal(12,2) NOT NULL DEFAULT '0.00',\n  `status` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL,\n  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,\n  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,\n  PRIMARY KEY (`id`),\n  UNIQUE KEY `orders_order_number_unique` (`order_number`),\n  KEY `orders_user_id_index` (`user_id`),\n  CONSTRAINT `orders_user_id_foreign` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE\n) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;",
            ];
        }

        return [
            'table_name' => $table,
            'columns' => [
                ['name' => 'id', 'position' => 1, 'type' => 'bigint unsigned', 'data_type' => 'bigint', 'is_nullable' => false, 'key' => 'PRI', 'default' => null, 'extra' => 'auto_increment', 'comment' => 'Unique Identifier'],
                ['name' => 'name', 'position' => 2, 'type' => 'varchar(255)', 'data_type' => 'varchar', 'is_nullable' => false, 'key' => '', 'default' => null, 'extra' => '', 'comment' => 'Name field'],
                ['name' => 'email', 'position' => 3, 'type' => 'varchar(255)', 'data_type' => 'varchar', 'is_nullable' => false, 'key' => 'UNI', 'default' => null, 'extra' => '', 'comment' => 'Email Address'],
                ['name' => 'is_active', 'position' => 4, 'type' => 'tinyint(1)', 'data_type' => 'tinyint', 'is_nullable' => false, 'key' => '', 'default' => '1', 'extra' => '', 'comment' => 'Active Flag'],
                ['name' => 'created_at', 'position' => 5, 'type' => 'timestamp', 'data_type' => 'timestamp', 'is_nullable' => true, 'key' => '', 'default' => 'CURRENT_TIMESTAMP', 'extra' => '', 'comment' => ''],
            ],
            'indexes' => [
                ['name' => 'PRIMARY', 'column' => 'id', 'non_unique' => false, 'is_primary' => true, 'seq_in_index' => 1, 'type' => 'BTREE'],
                ['name' => 'users_email_unique', 'column' => 'email', 'non_unique' => false, 'is_primary' => false, 'seq_in_index' => 1, 'type' => 'BTREE'],
            ],
            'foreign_keys' => [],
            'create_statement' => "CREATE TABLE `{$table}` (\n  `id` bigint unsigned NOT NULL AUTO_INCREMENT,\n  `name` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,\n  `email` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,\n  `is_active` tinyint(1) NOT NULL DEFAULT '1',\n  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,\n  PRIMARY KEY (`id`),\n  UNIQUE KEY `users_email_unique` (`email`)\n) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;",
        ];
    }

    public function getTableData(string $engine, string $name, string $table, array $params = []): array
    {
        $page = max(1, (int) ($params['page'] ?? 1));
        $perPage = min(500, max(1, (int) ($params['per_page'] ?? 50)));

        if (strtolower($engine) === 'postgresql' || strtolower($engine) === 'postgres') {
            $cols = ['id', 'account_id', 'status', 'metadata', 'created_at'];
            $rows = [
                ['id' => 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'account_id' => 'c1f9b3e1-22e4-4e9b-b6d8-111122223333', 'status' => 'active', 'metadata' => '{"plan": "pro", "seats": 10}', 'created_at' => '2026-08-28T08:00:00Z'],
                ['id' => 'b1ffcd88-8d0a-4fe7-aa5c-5aa8ac270b22', 'account_id' => 'c1f9b3e1-22e4-4e9b-b6d8-111122223333', 'status' => 'active', 'metadata' => '{"plan": "enterprise", "seats": 50}', 'created_at' => '2026-08-28T09:15:00Z'],
                ['id' => 'c2eedd77-7e09-4fd6-994b-4bb7bd160c33', 'account_id' => 'd2e8a4f2-33f5-4f0c-a7e9-444455556666', 'status' => 'trialing', 'metadata' => '{"plan": "starter", "seats": 2}', 'created_at' => '2026-08-28T10:30:00Z'],
            ];

            return [
                'table_name' => $table,
                'columns' => $cols,
                'rows' => $rows,
                'total_rows' => 3200,
                'page' => $page,
                'per_page' => $perPage,
                'total_pages' => 64,
            ];
        }

        if ($table === 'orders') {
            $cols = ['id', 'user_id', 'order_number', 'total_amount', 'status', 'created_at'];
            $rows = [
                ['id' => '1001', 'user_id' => '42', 'order_number' => 'ORD-2026-001', 'total_amount' => '149.99', 'status' => 'completed', 'created_at' => '2026-08-28 10:15:00'],
                ['id' => '1002', 'user_id' => '19', 'order_number' => 'ORD-2026-002', 'total_amount' => '89.50', 'status' => 'processing', 'created_at' => '2026-08-28 11:20:00'],
                ['id' => '1003', 'user_id' => '77', 'order_number' => 'ORD-2026-003', 'total_amount' => '320.00', 'status' => 'completed', 'created_at' => '2026-08-28 12:05:00'],
                ['id' => '1004', 'user_id' => '105', 'order_number' => 'ORD-2026-004', 'total_amount' => '24.95', 'status' => 'pending', 'created_at' => '2026-08-28 13:45:00'],
                ['id' => '1005', 'user_id' => '31', 'order_number' => 'ORD-2026-005', 'total_amount' => '512.10', 'status' => 'completed', 'created_at' => '2026-08-28 14:10:00'],
            ];
            return [
                'table_name' => $table,
                'columns' => $cols,
                'rows' => $rows,
                'total_rows' => 8930,
                'page' => $page,
                'per_page' => $perPage,
                'total_pages' => 179,
            ];
        }

        $cols = ['id', 'name', 'email', 'is_active', 'created_at'];
        $rows = [
            ['id' => '1', 'name' => 'Admin User', 'email' => 'admin@kodepreneur.com', 'is_active' => '1', 'created_at' => '2026-01-01 08:00:00'],
            ['id' => '2', 'name' => 'Sarah Connor', 'email' => 'sarah@cyberdyne.io', 'is_active' => '1', 'created_at' => '2026-01-15 09:30:00'],
            ['id' => '3', 'name' => 'Alex Murphy', 'email' => 'alex.murphy@ocp.corp', 'is_active' => '1', 'created_at' => '2026-02-10 14:22:00'],
            ['id' => '4', 'name' => 'Thomas Anderson', 'email' => 'neo@matrix.net', 'is_active' => '0', 'created_at' => '2026-03-01 19:45:00'],
            ['id' => '5', 'name' => 'Ellen Ripley', 'email' => 'ripley@weyland.org', 'is_active' => '1', 'created_at' => '2026-04-12 11:15:00'],
        ];

        return [
            'table_name' => $table,
            'columns' => $cols,
            'rows' => $rows,
            'total_rows' => 1420,
            'page' => $page,
            'per_page' => $perPage,
            'total_pages' => 29,
        ];
    }

    public function executeDeployment(array $payload): array
    {
        return [
            'success' => true,
            'exit_code' => 0,
            'log_output' => "[git] Checked out branch " . ($payload['branch'] ?? 'main') . "\n[composer] Installing dependencies (no-dev)\n[artisan] Migrating database: nothing to migrate.\n[artisan] Optimizing cached configuration and routes.\nDeployment completed successfully in 3s.",
            'duration_seconds' => 3,
            'commit_hash' => '8706c6c49832',
            'commit_message' => 'Release update: v1.0.0-prod',
        ];
    }

    public function syncCronJobs(string $systemUser, array $jobs): array
    {
        return [
            'success' => true,
            'message' => "Crontab for {$systemUser} synchronized successfully",
        ];
    }

    public function createProcess(array $payload): array
    {
        return [
            'success' => true,
            'unit_path' => "/etc/systemd/system/kp-proc-{$payload['id']}.service",
        ];
    }

    public function manageProcess(string $id, string $action): array
    {
        return [
            'success' => true,
            'output' => "Executed {$action} on kp-proc-{$id}.service",
        ];
    }

    public function deleteProcess(string $id): array
    {
        return [
            'success' => true,
            'message' => "Process {$id} deleted",
        ];
    }

    public function addFirewallRule(array $payload): array
    {
        return [
            'success' => true,
            'message' => "Firewall rule created",
        ];
    }

    public function deleteFirewallRule(array $payload): array
    {
        return [
            'success' => true,
            'message' => "Firewall rule deleted",
        ];
    }

    public function browseFiles(string $basePath, string $relativePath = '', bool $showHidden = false): array
    {
        $cleanRel = trim($relativePath, '/');

        if ($cleanRel === '') {
            $files = [
                [
                    'name' => 'app',
                    'path' => 'app',
                    'is_dir' => true,
                    'size_bytes' => 4096,
                    'permissions' => 'drwxr-xr-x',
                    'mode_octal' => '0755',
                    'owner' => 'www-data',
                    'group' => 'www-data',
                    'modified_at' => now()->subHours(2)->toISOString(),
                    'item_count' => 3,
                    'mime_type' => 'inode/directory',
                    'extension' => '',
                ],
                [
                    'name' => 'public',
                    'path' => 'public',
                    'is_dir' => true,
                    'size_bytes' => 4096,
                    'permissions' => 'drwxr-xr-x',
                    'mode_octal' => '0755',
                    'owner' => 'www-data',
                    'group' => 'www-data',
                    'modified_at' => now()->toISOString(),
                    'item_count' => 4,
                    'mime_type' => 'inode/directory',
                    'extension' => '',
                ],
                [
                    'name' => 'storage',
                    'path' => 'storage',
                    'is_dir' => true,
                    'size_bytes' => 4096,
                    'permissions' => 'drwxrwxr-x',
                    'mode_octal' => '0775',
                    'owner' => 'www-data',
                    'group' => 'www-data',
                    'modified_at' => now()->subDay()->toISOString(),
                    'item_count' => 3,
                    'mime_type' => 'inode/directory',
                    'extension' => '',
                ],
                [
                    'name' => 'index.php',
                    'path' => 'index.php',
                    'is_dir' => false,
                    'size_bytes' => 2450,
                    'permissions' => '-rw-r--r--',
                    'mode_octal' => '0644',
                    'owner' => 'www-data',
                    'group' => 'www-data',
                    'modified_at' => now()->toISOString(),
                    'item_count' => 0,
                    'mime_type' => 'application/x-httpd-php',
                    'extension' => 'php',
                ],
                [
                    'name' => 'composer.json',
                    'path' => 'composer.json',
                    'is_dir' => false,
                    'size_bytes' => 3120,
                    'permissions' => '-rw-r--r--',
                    'mode_octal' => '0644',
                    'owner' => 'www-data',
                    'group' => 'www-data',
                    'modified_at' => now()->toISOString(),
                    'item_count' => 0,
                    'mime_type' => 'application/json',
                    'extension' => 'json',
                ],
            ];

            if ($showHidden) {
                $files[] = [
                    'name' => '.env',
                    'path' => '.env',
                    'is_dir' => false,
                    'size_bytes' => 1240,
                    'permissions' => '-rw-r--r--',
                    'mode_octal' => '0644',
                    'owner' => 'www-data',
                    'group' => 'www-data',
                    'modified_at' => now()->toISOString(),
                    'item_count' => 0,
                    'mime_type' => 'text/plain',
                    'extension' => 'env',
                ];
                $files[] = [
                    'name' => '.gitignore',
                    'path' => '.gitignore',
                    'is_dir' => false,
                    'size_bytes' => 240,
                    'permissions' => '-rw-r--r--',
                    'mode_octal' => '0644',
                    'owner' => 'www-data',
                    'group' => 'www-data',
                    'modified_at' => now()->toISOString(),
                    'item_count' => 0,
                    'mime_type' => 'text/plain',
                    'extension' => 'gitignore',
                ];
            }

            return $files;
        }

        if ($cleanRel === 'public') {
            $files = [
                [
                    'name' => 'css',
                    'path' => 'public/css',
                    'is_dir' => true,
                    'size_bytes' => 4096,
                    'permissions' => 'drwxr-xr-x',
                    'mode_octal' => '0755',
                    'owner' => 'www-data',
                    'group' => 'www-data',
                    'modified_at' => now()->toISOString(),
                    'item_count' => 2,
                    'mime_type' => 'inode/directory',
                    'extension' => '',
                ],
                [
                    'name' => 'js',
                    'path' => 'public/js',
                    'is_dir' => true,
                    'size_bytes' => 4096,
                    'permissions' => 'drwxr-xr-x',
                    'mode_octal' => '0755',
                    'owner' => 'www-data',
                    'group' => 'www-data',
                    'modified_at' => now()->toISOString(),
                    'item_count' => 2,
                    'mime_type' => 'inode/directory',
                    'extension' => '',
                ],
                [
                    'name' => 'images',
                    'path' => 'public/images',
                    'is_dir' => true,
                    'size_bytes' => 4096,
                    'permissions' => 'drwxr-xr-x',
                    'mode_octal' => '0755',
                    'owner' => 'www-data',
                    'group' => 'www-data',
                    'modified_at' => now()->toISOString(),
                    'item_count' => 3,
                    'mime_type' => 'inode/directory',
                    'extension' => '',
                ],
                [
                    'name' => 'index.php',
                    'path' => 'public/index.php',
                    'is_dir' => false,
                    'size_bytes' => 1820,
                    'permissions' => '-rw-r--r--',
                    'mode_octal' => '0644',
                    'owner' => 'www-data',
                    'group' => 'www-data',
                    'modified_at' => now()->toISOString(),
                    'item_count' => 0,
                    'mime_type' => 'application/x-httpd-php',
                    'extension' => 'php',
                ],
                [
                    'name' => 'robots.txt',
                    'path' => 'public/robots.txt',
                    'is_dir' => false,
                    'size_bytes' => 120,
                    'permissions' => '-rw-r--r--',
                    'mode_octal' => '0644',
                    'owner' => 'www-data',
                    'group' => 'www-data',
                    'modified_at' => now()->toISOString(),
                    'item_count' => 0,
                    'mime_type' => 'text/plain',
                    'extension' => 'txt',
                ],
                [
                    'name' => 'favicon.ico',
                    'path' => 'public/favicon.ico',
                    'is_dir' => false,
                    'size_bytes' => 5400,
                    'permissions' => '-rw-r--r--',
                    'mode_octal' => '0644',
                    'owner' => 'www-data',
                    'group' => 'www-data',
                    'modified_at' => now()->toISOString(),
                    'item_count' => 0,
                    'mime_type' => 'image/x-icon',
                    'extension' => 'ico',
                ],
            ];

            if ($showHidden) {
                $files[] = [
                    'name' => '.htaccess',
                    'path' => 'public/.htaccess',
                    'is_dir' => false,
                    'size_bytes' => 450,
                    'permissions' => '-rw-r--r--',
                    'mode_octal' => '0644',
                    'owner' => 'www-data',
                    'group' => 'www-data',
                    'modified_at' => now()->toISOString(),
                    'item_count' => 0,
                    'mime_type' => 'text/plain',
                    'extension' => 'htaccess',
                ];
            }

            return $files;
        }

        if ($cleanRel === 'public/css') {
            return [
                [
                    'name' => 'app.css',
                    'path' => 'public/css/app.css',
                    'is_dir' => false,
                    'size_bytes' => 14200,
                    'permissions' => '-rw-r--r--',
                    'mode_octal' => '0644',
                    'owner' => 'www-data',
                    'group' => 'www-data',
                    'modified_at' => now()->toISOString(),
                    'item_count' => 0,
                    'mime_type' => 'text/css',
                    'extension' => 'css',
                ],
                [
                    'name' => 'custom.css',
                    'path' => 'public/css/custom.css',
                    'is_dir' => false,
                    'size_bytes' => 3100,
                    'permissions' => '-rw-r--r--',
                    'mode_octal' => '0644',
                    'owner' => 'www-data',
                    'group' => 'www-data',
                    'modified_at' => now()->toISOString(),
                    'item_count' => 0,
                    'mime_type' => 'text/css',
                    'extension' => 'css',
                ],
            ];
        }

        if ($cleanRel === 'public/js') {
            return [
                [
                    'name' => 'app.js',
                    'path' => 'public/js/app.js',
                    'is_dir' => false,
                    'size_bytes' => 38400,
                    'permissions' => '-rw-r--r--',
                    'mode_octal' => '0644',
                    'owner' => 'www-data',
                    'group' => 'www-data',
                    'modified_at' => now()->toISOString(),
                    'item_count' => 0,
                    'mime_type' => 'application/javascript',
                    'extension' => 'js',
                ],
                [
                    'name' => 'vendor.js',
                    'path' => 'public/js/vendor.js',
                    'is_dir' => false,
                    'size_bytes' => 124000,
                    'permissions' => '-rw-r--r--',
                    'mode_octal' => '0644',
                    'owner' => 'www-data',
                    'group' => 'www-data',
                    'modified_at' => now()->toISOString(),
                    'item_count' => 0,
                    'mime_type' => 'application/javascript',
                    'extension' => 'js',
                ],
            ];
        }

        if ($cleanRel === 'public/images') {
            return [
                [
                    'name' => 'logo.png',
                    'path' => 'public/images/logo.png',
                    'is_dir' => false,
                    'size_bytes' => 8400,
                    'permissions' => '-rw-r--r--',
                    'mode_octal' => '0644',
                    'owner' => 'www-data',
                    'group' => 'www-data',
                    'modified_at' => now()->toISOString(),
                    'item_count' => 0,
                    'mime_type' => 'image/png',
                    'extension' => 'png',
                ],
                [
                    'name' => 'banner.jpg',
                    'path' => 'public/images/banner.jpg',
                    'is_dir' => false,
                    'size_bytes' => 45200,
                    'permissions' => '-rw-r--r--',
                    'mode_octal' => '0644',
                    'owner' => 'www-data',
                    'group' => 'www-data',
                    'modified_at' => now()->toISOString(),
                    'item_count' => 0,
                    'mime_type' => 'image/jpeg',
                    'extension' => 'jpg',
                ],
                [
                    'name' => 'icon.svg',
                    'path' => 'public/images/icon.svg',
                    'is_dir' => false,
                    'size_bytes' => 1200,
                    'permissions' => '-rw-r--r--',
                    'mode_octal' => '0644',
                    'owner' => 'www-data',
                    'group' => 'www-data',
                    'modified_at' => now()->toISOString(),
                    'item_count' => 0,
                    'mime_type' => 'image/svg+xml',
                    'extension' => 'svg',
                ],
            ];
        }

        if ($cleanRel === 'storage') {
            return [
                [
                    'name' => 'app',
                    'path' => 'storage/app',
                    'is_dir' => true,
                    'size_bytes' => 4096,
                    'permissions' => 'drwxrwxr-x',
                    'mode_octal' => '0775',
                    'owner' => 'www-data',
                    'group' => 'www-data',
                    'modified_at' => now()->toISOString(),
                    'item_count' => 2,
                    'mime_type' => 'inode/directory',
                    'extension' => '',
                ],
                [
                    'name' => 'framework',
                    'path' => 'storage/framework',
                    'is_dir' => true,
                    'size_bytes' => 4096,
                    'permissions' => 'drwxrwxr-x',
                    'mode_octal' => '0775',
                    'owner' => 'www-data',
                    'group' => 'www-data',
                    'modified_at' => now()->toISOString(),
                    'item_count' => 3,
                    'mime_type' => 'inode/directory',
                    'extension' => '',
                ],
                [
                    'name' => 'logs',
                    'path' => 'storage/logs',
                    'is_dir' => true,
                    'size_bytes' => 4096,
                    'permissions' => 'drwxrwxr-x',
                    'mode_octal' => '0775',
                    'owner' => 'www-data',
                    'group' => 'www-data',
                    'modified_at' => now()->toISOString(),
                    'item_count' => 2,
                    'mime_type' => 'inode/directory',
                    'extension' => '',
                ],
            ];
        }

        if ($cleanRel === 'storage/logs') {
            return [
                [
                    'name' => 'laravel.log',
                    'path' => 'storage/logs/laravel.log',
                    'is_dir' => false,
                    'size_bytes' => 48200,
                    'permissions' => '-rw-rw-r--',
                    'mode_octal' => '0664',
                    'owner' => 'www-data',
                    'group' => 'www-data',
                    'modified_at' => now()->toISOString(),
                    'item_count' => 0,
                    'mime_type' => 'text/plain',
                    'extension' => 'log',
                ],
                [
                    'name' => 'worker.log',
                    'path' => 'storage/logs/worker.log',
                    'is_dir' => false,
                    'size_bytes' => 12400,
                    'permissions' => '-rw-rw-r--',
                    'mode_octal' => '0664',
                    'owner' => 'www-data',
                    'group' => 'www-data',
                    'modified_at' => now()->toISOString(),
                    'item_count' => 0,
                    'mime_type' => 'text/plain',
                    'extension' => 'log',
                ],
            ];
        }

        if ($cleanRel === 'app') {
            return [
                [
                    'name' => 'Http',
                    'path' => 'app/Http',
                    'is_dir' => true,
                    'size_bytes' => 4096,
                    'permissions' => 'drwxr-xr-x',
                    'mode_octal' => '0755',
                    'owner' => 'www-data',
                    'group' => 'www-data',
                    'modified_at' => now()->toISOString(),
                    'item_count' => 2,
                    'mime_type' => 'inode/directory',
                    'extension' => '',
                ],
                [
                    'name' => 'Models',
                    'path' => 'app/Models',
                    'is_dir' => true,
                    'size_bytes' => 4096,
                    'permissions' => 'drwxr-xr-x',
                    'mode_octal' => '0755',
                    'owner' => 'www-data',
                    'group' => 'www-data',
                    'modified_at' => now()->toISOString(),
                    'item_count' => 2,
                    'mime_type' => 'inode/directory',
                    'extension' => '',
                ],
            ];
        }

        // Generic mock fallback for any other subfolder
        return [
            [
                'name' => 'subfile1.txt',
                'path' => "{$cleanRel}/subfile1.txt",
                'is_dir' => false,
                'size_bytes' => 1024,
                'permissions' => '-rw-r--r--',
                'mode_octal' => '0644',
                'owner' => 'www-data',
                'group' => 'www-data',
                'modified_at' => now()->toISOString(),
                'item_count' => 0,
                'mime_type' => 'text/plain',
                'extension' => 'txt',
            ],
            [
                'name' => 'data.json',
                'path' => "{$cleanRel}/data.json",
                'is_dir' => false,
                'size_bytes' => 2048,
                'permissions' => '-rw-r--r--',
                'mode_octal' => '0644',
                'owner' => 'www-data',
                'group' => 'www-data',
                'modified_at' => now()->toISOString(),
                'item_count' => 0,
                'mime_type' => 'application/json',
                'extension' => 'json',
            ],
        ];
    }

    public function readFile(string $basePath, string $relativePath, int $maxBytes = 5242880): string
    {
        return "APP_NAME=KodepreneurApp\nAPP_ENV=production\nAPP_DEBUG=false\nAPP_URL=https://example.com\n";
    }

    public function readFileBase64(string $basePath, string $relativePath, int $maxBytes = 52428800): string
    {
        return base64_encode("APP_NAME=KodepreneurApp\nAPP_ENV=production\nAPP_DEBUG=false\nAPP_URL=https://example.com\n");
    }

    public function writeFile(string $basePath, string $relativePath, string $content): array
    {
        return [
            'success' => true,
            'message' => "File written successfully",
        ];
    }

    public function writeFileBase64(string $basePath, string $relativePath, string $contentBase64): array
    {
        return [
            'success' => true,
            'message' => "File written successfully",
        ];
    }

    public function createFile(string $basePath, string $relativePath): array
    {
        return [
            'success' => true,
            'message' => "File created successfully",
        ];
    }

    public function createDirectory(string $basePath, string $relativePath): array
    {
        return [
            'success' => true,
            'message' => "Directory created successfully",
        ];
    }

    public function deleteFile(string $basePath, string $relativePath): array
    {
        return [
            'success' => true,
            'message' => "Entry deleted successfully",
        ];
    }

    public function renameFile(string $basePath, string $oldPath, string $newPath): array
    {
        return [
            'success' => true,
            'message' => "Renamed successfully",
        ];
    }

    public function copyFile(string $basePath, string $srcPath, string $destPath): array
    {
        return [
            'success' => true,
            'message' => "Copied successfully",
        ];
    }

    public function moveFile(string $basePath, string $srcPath, string $destPath): array
    {
        return [
            'success' => true,
            'message' => "Moved successfully",
        ];
    }

    public function chmodFile(string $basePath, string $relativePath, string $mode, bool $recursive = false): array
    {
        return [
            'success' => true,
            'message' => "Permissions updated successfully",
        ];
    }

    public function chownFile(string $basePath, string $relativePath, int $uid, int $gid, bool $recursive = false): array
    {
        return [
            'success' => true,
            'message' => "Ownership updated successfully",
        ];
    }

    public function statFile(string $basePath, string $relativePath): array
    {
        return [
            'name' => basename($relativePath),
            'path' => $relativePath,
            'is_dir' => false,
            'size_bytes' => 1240,
            'permissions' => '-rw-r--r--',
            'mode_octal' => '0644',
            'owner' => 'www-data',
            'group' => 'www-data',
            'uid' => 33,
            'gid' => 33,
            'modified_at' => now()->toISOString(),
            'created_at' => now()->toISOString(),
            'mime_type' => 'text/plain',
            'extension' => 'txt',
            'item_count' => 0,
        ];
    }

    public function compressFiles(string $basePath, array $sources, string $destPath, string $format = 'zip'): array
    {
        return [
            'success' => true,
            'message' => "Archive created successfully",
        ];
    }

    public function extractArchive(string $basePath, string $archivePath, string $destPath): array
    {
        return [
            'success' => true,
            'message' => "Archive extracted successfully",
        ];
    }

    public function searchFiles(string $basePath, string $query, int $maxResults = 100): array
    {
        return [
            [
                'name' => 'index.php',
                'path' => 'public/index.php',
                'is_dir' => false,
                'size_bytes' => 2450,
                'permissions' => '-rw-r--r--',
                'mode_octal' => '0644',
                'owner' => 'www-data',
                'group' => 'www-data',
                'modified_at' => now()->toISOString(),
                'mime_type' => 'application/x-httpd-php',
                'extension' => 'php',
            ]
        ];
    }

    public function getDiskUsage(string $basePath): array
    {
        return [
            'path' => $basePath,
            'total_bytes' => 50000000000,
            'used_bytes' => 12000000000,
            'free_bytes' => 38000000000,
            'usage_percent' => 24.0,
            'path_size' => 1800000000,
        ];
    }

    public function executeSystemUpdate(array $payload = []): array
    {
        $repo = $payload['repository'] ?? 'https://github.com/kodepreneur/kodepreneur-panel.git';
        $branch = $payload['branch'] ?? 'main';

        $logs = [];
        $logs[] = "[" . date('H:i:s') . "] 🚀 Starting Kodepreneur Panel Automated Update Pipeline...";
        $logs[] = "[" . date('H:i:s') . "] [1/6] Fetching latest release from {$repo} (branch: {$branch})...";
        $logs[] = "  - Cloned latest commits successfully (HEAD -> 8706c6c)";
        $logs[] = "[" . date('H:i:s') . "] [2/6] Synchronizing Control Plane files and updating Composer packages...";
        $logs[] = "  - Installing dependencies (no-dev, --optimize-autoloader)";
        $logs[] = "  - Generated optimized autoload files containing 3,420 classes";
        $logs[] = "[" . date('H:i:s') . "] [3/6] Running database migrations...";
        $logs[] = "  - Nothing to migrate / Migrations applied up to date";
        $logs[] = "[" . date('H:i:s') . "] [4/6] Compiling frontend assets via Vite...";
        $logs[] = "  - vite v6.4.3 building for production...";
        $logs[] = "  - ✓ built in 1.45s";
        $logs[] = "[" . date('H:i:s') . "] [5/6] Optimizing Laravel application caches and permissions...";
        $logs[] = "  - Configuration cache cleared and rebuilt.";
        $logs[] = "  - Route cache optimized.";
        $logs[] = "  - Set ownership to www-data:www-data";
        $logs[] = "[" . date('H:i:s') . "] [6/6] Recompiling Go Agent daemon and reloading Nginx...";
        $logs[] = "  - kodepreneur-agent binary updated: /usr/local/bin/kodepreneur-agent";
        $logs[] = "  - Restarted kodepreneur-agent.service";
        $logs[] = "  - Nginx configuration valid. Reloaded nginx.service";
        $logs[] = "[" . date('H:i:s') . "] 🎉 Kodepreneur Panel successfully updated in 1s!";

        return [
            'success' => true,
            'exit_code' => 0,
            'log_output' => implode("\n", $logs),
            'duration_seconds' => 1,
            'commit_hash' => '8706c6c49832',
            'commit_message' => 'feat(updates): add automated in-panel GitHub upgrade manager',
        ];
    }
}
