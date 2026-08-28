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
        return [
            'domain' => $payload['domain'] ?? 'example.com',
            'vhost_path' => "/etc/nginx/sites-available/" . ($payload['domain'] ?? 'example.com') . ".conf",
            'system_user' => $payload['system_user'] ?? 'kp_user',
            'document_root' => $payload['document_root'] ?? '/var/www/example.com/public',
            'php_version' => $payload['php_version'] ?? '8.3',
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

    public function browseFiles(string $basePath, string $relativePath = ''): array
    {
        return [
            [
                'name' => 'public',
                'path' => 'public',
                'is_dir' => true,
                'size_bytes' => 4096,
                'permissions' => 'drwxr-xr-x',
                'modified_at' => now()->toISOString(),
            ],
            [
                'name' => '.env',
                'path' => '.env',
                'is_dir' => false,
                'size_bytes' => 1240,
                'permissions' => '-rw-r--r--',
                'modified_at' => now()->toISOString(),
            ],
            [
                'name' => 'composer.json',
                'path' => 'composer.json',
                'is_dir' => false,
                'size_bytes' => 3120,
                'permissions' => '-rw-r--r--',
                'modified_at' => now()->toISOString(),
            ],
        ];
    }

    public function readFile(string $basePath, string $relativePath): string
    {
        return "APP_NAME=KodepreneurApp\nAPP_ENV=production\nAPP_DEBUG=false\nAPP_URL=https://example.com\n";
    }

    public function writeFile(string $basePath, string $relativePath, string $content): array
    {
        return [
            'success' => true,
            'message' => "File written successfully",
        ];
    }

    public function deleteFile(string $basePath, string $relativePath): array
    {
        return [
            'success' => true,
            'message' => "Entry deleted successfully",
        ];
    }
}
