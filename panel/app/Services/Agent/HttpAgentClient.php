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
        $res = $this->request('POST', '/api/v1/websites', $payload);
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

    public function executeDeployment(array $payload): array
    {
        $res = $this->request('POST', '/api/v1/deployments/execute', $payload);
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

    public function browseFiles(string $basePath, string $relativePath = ''): array
    {
        $res = $this->request('POST', '/api/v1/files/browse', [
            'base_path' => $basePath,
            'relative_path' => $relativePath,
        ]);
        return $res['data'] ?? [];
    }

    public function readFile(string $basePath, string $relativePath): string
    {
        $res = $this->request('POST', '/api/v1/files/read', [
            'base_path' => $basePath,
            'relative_path' => $relativePath,
        ]);
        return $res['content'] ?? '';
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

    public function deleteFile(string $basePath, string $relativePath): array
    {
        $res = $this->request('POST', '/api/v1/files/delete', [
            'base_path' => $basePath,
            'relative_path' => $relativePath,
        ]);
        return $res ?? [];
    }

    public function executeSystemUpdate(array $payload = []): array
    {
        $res = $this->request('POST', '/api/v1/system/update', $payload);
        return $res['data'] ?? [];
    }

    protected function request(string $method, string $path, array $data = []): array
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
                ->timeout(20)
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
