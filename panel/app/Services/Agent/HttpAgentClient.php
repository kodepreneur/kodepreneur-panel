<?php

namespace App\Services\Agent;

use Exception;
use Illuminate\Support\Facades\Http;
use Illuminate\Support\Str;

class HttpAgentClient implements AgentClientInterface
{
    public function __construct(
        protected string $baseUrl,
        protected string $secret
    ) {}

    public function health(): array
    {
        return $this->request('GET', '/api/v1/health');
    }

    public function getSystemInfo(): array
    {
        $res = $this->request('GET', '/api/v1/system/info');
        return $res['data'] ?? [];
    }

    public function getSystemMetrics(): array
    {
        $res = $this->request('GET', '/api/v1/system/metrics');
        return $res['data'] ?? [];
    }

    public function getServiceStatus(string $service): array
    {
        $res = $this->request('GET', "/api/v1/services/{$service}/status");
        return $res['data'] ?? [];
    }

    public function manageService(string $service, string $action): array
    {
        $res = $this->request('POST', "/api/v1/services/{$service}/{$action}");
        return $res['data'] ?? [];
    }

    public function createWebsite(array $payload): array
    {
        $res = $this->request('POST', '/api/v1/websites', $payload, 600);
        return $res['data'] ?? [];
    }

    public function deleteWebsite(string $domain, array $payload = []): array
    {
        $res = $this->request('DELETE', "/api/v1/websites/{$domain}", $payload);
        return $res['data'] ?? [];
    }

    public function switchPhpVersion(string $domain, array $payload): array
    {
        $res = $this->request('PUT', "/api/v1/websites/{$domain}/php", $payload);
        return $res['data'] ?? [];
    }

    public function issueSsl(array $payload): array
    {
        $res = $this->request('POST', '/api/v1/ssl/issue', $payload);
        return $res['data'] ?? [];
    }

    public function getWebsiteLogs(string $domain, string $type = 'access', int $lines = 100): array
    {
        $res = $this->request('GET', "/api/v1/websites/{$domain}/logs/{$type}?lines={$lines}");
        return $res['data'] ?? [];
    }

    public function getWebsiteTraffic(string $domain, string $period = '24h'): array
    {
        $res = $this->request('GET', "/api/v1/websites/{$domain}/traffic?period={$period}");
        return $res['data'] ?? [];
    }

    public function createDatabase(array $payload): array
    {
        $res = $this->request('POST', '/api/v1/databases', $payload);
        return $res['data'] ?? [];
    }

    public function deleteDatabase(string $engine, string $name): array
    {
        $res = $this->request('DELETE', "/api/v1/databases/{$engine}/{$name}");
        return $res ?? [];
    }

    public function createDatabaseUser(array $payload): array
    {
        $res = $this->request('POST', '/api/v1/databases/users', $payload);
        return $res['data'] ?? [];
    }

    public function deleteDatabaseUser(string $engine, string $username, string $host = 'localhost'): array
    {
        $res = $this->request('DELETE', "/api/v1/databases/users/{$engine}/{$username}?host={$host}");
        return $res ?? [];
    }

    public function grantDatabaseAccess(array $payload): array
    {
        $res = $this->request('POST', '/api/v1/databases/grants', $payload);
        return $res ?? [];
    }

    public function resetDatabaseUserPassword(array $payload): array
    {
        $res = $this->request('POST', '/api/v1/databases/users/password', $payload);
        return $res ?? [];
    }

    public function getDatabaseTables(string $engine, string $name): array
    {
        $res = $this->request('GET', "/api/v1/databases/{$engine}/{$name}/tables");
        return $res['data'] ?? [];
    }

    public function getTableStructure(string $engine, string $name, string $table): array
    {
        $res = $this->request('GET', "/api/v1/databases/{$engine}/{$name}/tables/{$table}/structure");
        return $res['data'] ?? [];
    }

    public function getTableData(string $engine, string $name, string $table, array $params = []): array
    {
        $queryString = http_build_query($params);
        $url = "/api/v1/databases/{$engine}/{$name}/tables/{$table}/data" . ($queryString ? "?{$queryString}" : '');
        $res = $this->request('GET', $url);
        return $res['data'] ?? [];
    }

    public function executeDeployment(array $payload): array
    {
        $res = $this->request('POST', '/api/v1/deployments/execute', $payload, 600);
        return $res['data'] ?? [];
    }

    public function syncCronJobs(string $systemUser, array $jobs): array
    {
        $res = $this->request('POST', '/api/v1/cron/sync', [
            'system_user' => $systemUser,
            'jobs' => $jobs,
        ]);
        return $res ?? [];
    }

    public function createProcess(array $payload): array
    {
        $res = $this->request('POST', '/api/v1/processes', $payload);
        return $res ?? [];
    }

    public function manageProcess(string $id, string $action): array
    {
        $res = $this->request('POST', "/api/v1/processes/{$id}/{$action}");
        return $res ?? [];
    }

    public function deleteProcess(string $id): array
    {
        $res = $this->request('DELETE', "/api/v1/processes/{$id}");
        return $res ?? [];
    }

    public function addFirewallRule(array $payload): array
    {
        $res = $this->request('POST', '/api/v1/firewall/rules', $payload);
        return $res ?? [];
    }

    public function deleteFirewallRule(array $payload): array
    {
        $res = $this->request('DELETE', '/api/v1/firewall/rules', $payload);
        return $res ?? [];
    }

    public function browseFiles(string $basePath, string $relativePath = '', bool $showHidden = false): array
    {
        $res = $this->request('POST', '/api/v1/files/browse', [
            'base_path' => $basePath,
            'relative_path' => $relativePath,
            'show_hidden' => $showHidden,
        ]);
        return $res['data'] ?? [];
    }

    public function readFile(string $basePath, string $relativePath, int $maxBytes = 5242880): string
    {
        $res = $this->request('POST', '/api/v1/files/read', [
            'base_path' => $basePath,
            'relative_path' => $relativePath,
            'max_bytes' => $maxBytes,
        ]);
        return $res['content'] ?? '';
    }

    public function readFileBase64(string $basePath, string $relativePath, int $maxBytes = 52428800): string
    {
        $res = $this->request('POST', '/api/v1/files/read', [
            'base_path' => $basePath,
            'relative_path' => $relativePath,
            'max_bytes' => $maxBytes,
            'as_base64' => true,
        ]);
        return $res['content_base64'] ?? '';
    }

    public function writeFile(string $basePath, string $relativePath, string $content): array
    {
        $res = $this->request('POST', '/api/v1/files/write', [
            'base_path' => $basePath,
            'relative_path' => $relativePath,
            'content' => $content,
        ]);
        return $res ?? [];
    }

    public function writeFileBase64(string $basePath, string $relativePath, string $contentBase64): array
    {
        $res = $this->request('POST', '/api/v1/files/write', [
            'base_path' => $basePath,
            'relative_path' => $relativePath,
            'content_base64' => $contentBase64,
        ]);
        return $res ?? [];
    }

    public function createFile(string $basePath, string $relativePath): array
    {
        $res = $this->request('POST', '/api/v1/files/create', [
            'base_path' => $basePath,
            'relative_path' => $relativePath,
        ]);
        return $res ?? [];
    }

    public function createDirectory(string $basePath, string $relativePath): array
    {
        $res = $this->request('POST', '/api/v1/files/mkdir', [
            'base_path' => $basePath,
            'relative_path' => $relativePath,
        ]);
        return $res ?? [];
    }

    public function deleteFile(string $basePath, string $relativePath): array
    {
        $res = $this->request('POST', '/api/v1/files/delete', [
            'base_path' => $basePath,
            'relative_path' => $relativePath,
        ]);
        return $res ?? [];
    }

    public function renameFile(string $basePath, string $oldPath, string $newPath): array
    {
        $res = $this->request('POST', '/api/v1/files/rename', [
            'base_path' => $basePath,
            'old_path' => $oldPath,
            'new_path' => $newPath,
        ]);
        return $res ?? [];
    }

    public function copyFile(string $basePath, string $srcPath, string $destPath): array
    {
        $res = $this->request('POST', '/api/v1/files/copy', [
            'base_path' => $basePath,
            'src_path' => $srcPath,
            'dest_path' => $destPath,
        ]);
        return $res ?? [];
    }

    public function moveFile(string $basePath, string $srcPath, string $destPath): array
    {
        $res = $this->request('POST', '/api/v1/files/move', [
            'base_path' => $basePath,
            'src_path' => $srcPath,
            'dest_path' => $destPath,
        ]);
        return $res ?? [];
    }

    public function chmodFile(string $basePath, string $relativePath, string $mode, bool $recursive = false): array
    {
        $res = $this->request('POST', '/api/v1/files/chmod', [
            'base_path' => $basePath,
            'relative_path' => $relativePath,
            'mode' => $mode,
            'recursive' => $recursive,
        ]);
        return $res ?? [];
    }

    public function chownFile(string $basePath, string $relativePath, int $uid, int $gid, bool $recursive = false): array
    {
        $res = $this->request('POST', '/api/v1/files/chown', [
            'base_path' => $basePath,
            'relative_path' => $relativePath,
            'uid' => $uid,
            'gid' => $gid,
            'recursive' => $recursive,
        ]);
        return $res ?? [];
    }

    public function statFile(string $basePath, string $relativePath): array
    {
        $res = $this->request('POST', '/api/v1/files/stat', [
            'base_path' => $basePath,
            'relative_path' => $relativePath,
        ]);
        return $res['data'] ?? [];
    }

    public function compressFiles(string $basePath, array $sources, string $destPath, string $format = 'zip'): array
    {
        $res = $this->request('POST', '/api/v1/files/compress', [
            'base_path' => $basePath,
            'sources' => $sources,
            'dest_path' => $destPath,
            'format' => $format,
        ], 120);
        return $res ?? [];
    }

    public function extractArchive(string $basePath, string $archivePath, string $destPath): array
    {
        $res = $this->request('POST', '/api/v1/files/extract', [
            'base_path' => $basePath,
            'archive_path' => $archivePath,
            'dest_path' => $destPath,
        ], 120);
        return $res ?? [];
    }

    public function searchFiles(string $basePath, string $query, int $maxResults = 100): array
    {
        $res = $this->request('POST', '/api/v1/files/search', [
            'base_path' => $basePath,
            'query' => $query,
            'max_results' => $maxResults,
        ], 60);
        return $res['data'] ?? [];
    }

    public function getDiskUsage(string $basePath): array
    {
        $res = $this->request('POST', '/api/v1/files/disk', [
            'base_path' => $basePath,
        ]);
        return $res['data'] ?? [];
    }

    public function executeSystemUpdate(array $payload = []): array
    {
        $res = $this->request('POST', '/api/v1/system/update', $payload, 600);
        return $res['data'] ?? [];
    }

    protected function request(string $method, string $path, array $data = [], int $timeout = 30): array
    {
        $url = rtrim($this->baseUrl, '/') . $path;
        $timestamp = (string) time();
        $nonce = Str::uuid()->toString();
        $body = !empty($data) ? json_encode($data) : '';

        // Extract path without query string for HMAC signature
        $parsedPath = parse_url($path, PHP_URL_PATH) ?: $path;

        $bodyHash = hash('sha256', (string) $body);
        $message = "{$method}\n{$parsedPath}\n{$timestamp}\n{$nonce}\n{$bodyHash}";
        $signature = hash_hmac('sha256', $message, $this->secret);

        $headers = [
            'Accept' => 'application/json',
            'Content-Type' => 'application/json',
            'X-Kodepreneur-Timestamp' => $timestamp,
            'X-Kodepreneur-Nonce' => $nonce,
            'X-Kodepreneur-Signature' => $signature,
        ];

        try {
            $response = Http::withHeaders($headers)
                ->timeout($timeout)
                ->send($method, $url, [
                    'body' => $body,
                ]);

            if (!$response->successful()) {
                throw new Exception("Agent API error [{$response->status()}]: " . $response->body());
            }

            return $response->json() ?? [];
        } catch (Exception $e) {
            throw new Exception("Failed to communicate with Kodepreneur Agent: " . $e->getMessage(), 0, $e);
        }
    }
}
