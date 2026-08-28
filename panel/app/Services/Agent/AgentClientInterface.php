<?php

namespace App\Services\Agent;

interface AgentClientInterface
{
    /**
     * Check agent daemon health.
     */
    public function health(): array;

    /**
     * Get system OS, kernel, hostname, and uptime metadata.
     */
    public function getSystemInfo(): array;

    /**
     * Get real-time hardware metrics (CPU, RAM, Disk, Load).
     */
    public function getSystemMetrics(): array;

    /**
     * Get status of a managed systemd service.
     */
    public function getServiceStatus(string $service): array;

    /**
     * Perform an action (start, stop, restart, reload) on a managed systemd service.
     */
    public function manageService(string $service, string $action): array;

    /**
     * Create an isolated website, user, FPM pool, and Nginx vhost.
     */
    public function createWebsite(array $payload): array;

    /**
     * Remove website vhost and PHP pool.
     */
    public function deleteWebsite(string $domain, array $payload = []): array;

    /**
     * Switch PHP runtime version for a website.
     */
    public function switchPhpVersion(string $domain, array $payload): array;

    /**
     * Issue or renew Let's Encrypt SSL certificate via ACME.
     */
    public function issueSsl(array $payload): array;

    /**
     * Retrieve access or error log lines for a website.
     */
    public function getWebsiteLogs(string $domain, string $type = 'access', int $lines = 100): array;

    /**
     * Create a MySQL or PostgreSQL database.
     */
    public function createDatabase(array $payload): array;

    /**
     * Drop a MySQL or PostgreSQL database.
     */
    public function deleteDatabase(string $engine, string $name): array;

    /**
     * Create a MySQL or PostgreSQL database user.
     */
    public function createDatabaseUser(array $payload): array;

    /**
     * Drop a MySQL or PostgreSQL database user.
     */
    public function deleteDatabaseUser(string $engine, string $username, string $host = 'localhost'): array;

    /**
     * Grant privileges on a database to a user.
     */
    public function grantDatabaseAccess(array $payload): array;

    /**
     * Reset a database user password.
     */
    public function resetDatabaseUserPassword(array $payload): array;

    /**
     * Execute an unprivileged git deployment and build pipeline.
     */
    public function executeDeployment(array $payload): array;

    /**
     * Sync user crontab jobs.
     */
    public function syncCronJobs(string $systemUser, array $jobs): array;

    /**
     * Create a managed Systemd application unit.
     */
    public function createProcess(array $payload): array;

    /**
     * Start, stop, restart a managed process unit.
     */
    public function manageProcess(string $id, string $action): array;

    /**
     * Delete a managed process unit.
     */
    public function deleteProcess(string $id): array;

    /**
     * Add a UFW firewall rule.
     */
    public function addFirewallRule(array $payload): array;

    /**
     * Delete a UFW firewall rule.
     */
    public function deleteFirewallRule(array $payload): array;

    /**
     * Browse sandboxed directory.
     */
    public function browseFiles(string $basePath, string $relativePath = ''): array;

    /**
     * Read file content safely.
     */
    public function readFile(string $basePath, string $relativePath): string;

    /**
     * Write file content safely.
     */
    public function writeFile(string $basePath, string $relativePath, string $content): array;

    /**
     * Delete file or folder safely.
     */
    public function deleteFile(string $basePath, string $relativePath): array;

    /**
     * Trigger root-level automated system update pipeline.
     */
    public function executeSystemUpdate(array $payload = []): array;
}
