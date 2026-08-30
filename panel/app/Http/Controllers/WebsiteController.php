<?php

namespace App\Http\Controllers;

use App\Models\ActivityLog;
use App\Models\Database;
use App\Models\DatabaseAccess;
use App\Models\DatabaseUser;
use App\Models\Deployment;
use App\Models\Domain;
use App\Models\SslCertificate;
use App\Models\Website;
use App\Services\Agent\AgentClientInterface;
use Exception;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\RedirectResponse;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Storage;
use Illuminate\Support\Str;
use Inertia\Inertia;
use Inertia\Response;

class WebsiteController extends Controller
{
    public function __construct(
        protected AgentClientInterface $agentClient
    ) {}

    public function index(): Response
    {
        $websites = Website::with(['sslCertificate', 'domains'])
            ->withCount('deployments')
            ->latest()
            ->paginate(15);

        return Inertia::render('Websites/Index', [
            'websites' => $websites,
        ]);
    }

    public function create(): Response
    {
        return Inertia::render('Websites/Create');
    }

    public function generateDeployKey(): JsonResponse
    {
        // 1. Try ssh-keygen for modern ed25519 key
        $tempFile = tempnam(sys_get_temp_dir(), 'kp_key_');
        @unlink($tempFile);
        $cmd = "ssh-keygen -t ed25519 -N '' -C 'kodepreneur-deploy-key' -f " . escapeshellarg($tempFile) . " 2>&1";
        @exec($cmd, $out, $ret);

        if ($ret === 0 && file_exists($tempFile) && file_exists("{$tempFile}.pub")) {
            $privateKey = file_get_contents($tempFile);
            $publicKey = trim(file_get_contents("{$tempFile}.pub"));
            @unlink($tempFile);
            @unlink("{$tempFile}.pub");
            return response()->json([
                'success' => true,
                'public_key' => $publicKey,
                'private_key' => $privateKey,
                'type' => 'ed25519',
            ]);
        }

        // 2. Fallback to OpenSSL RSA-4096 if ssh-keygen is unavailable
        $config = [
            "digest_alg" => "sha512",
            "private_key_bits" => 4096,
            "private_key_type" => OPENSSL_KEYTYPE_RSA,
        ];
        $res = openssl_pkey_new($config);
        if ($res) {
            openssl_pkey_export($res, $privateKey);
            $keyDetails = openssl_pkey_get_details($res);
            $rsa = $keyDetails['rsa'] ?? [];
            if (!empty($rsa['e']) && !empty($rsa['n'])) {
                $buffer = pack('N', 7) . 'ssh-rsa' .
                    self::sshEncodeBuffer($rsa['e']) .
                    self::sshEncodeBuffer($rsa['n']);
                $publicKey = 'ssh-rsa ' . base64_encode($buffer) . ' kodepreneur-deploy-key';
                return response()->json([
                    'success' => true,
                    'public_key' => $publicKey,
                    'private_key' => $privateKey,
                    'type' => 'rsa',
                ]);
            }
        }

        return response()->json([
            'success' => false,
            'message' => 'Unable to generate SSH key pair on this server.',
        ], 500);
    }

    private static function sshEncodeBuffer(string $buffer): string
    {
        $len = strlen($buffer);
        if (ord($buffer[0]) & 0x80) {
            $len++;
            $buffer = "\x00" . $buffer;
        }
        return pack('Na*', $len, $buffer);
    }

    public function store(Request $request): RedirectResponse
    {
        @set_time_limit(600);
        @ini_set('max_execution_time', '600');
        @ini_set('memory_limit', '512M');

        $validated = $request->validate([
            'domain' => ['required', 'string', 'max:255', 'unique:websites,domain', 'regex:/^[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/'],
            'php_version' => ['required', 'string', 'in:8.3,8.4,none'],
            'deployment_source' => ['nullable', 'string', 'in:empty,zip,git'],
            'project_type' => ['nullable', 'string', 'in:laravel,generic_php,static,auto'],
            'document_root' => ['nullable', 'string', 'max:255'],
            'git_repository' => ['required_if:deployment_source,git', 'nullable', 'string', 'max:255'],
            'git_branch' => ['nullable', 'string', 'max:100'],
            'git_auth_type' => ['nullable', 'string', 'in:none,ssh_key,token'],
            'git_token' => ['nullable', 'string', 'max:500'],
            'git_token_user' => ['nullable', 'string', 'max:100'],
            'git_ssh_private_key' => ['nullable', 'string'],
            'git_ssh_public_key' => ['nullable', 'string'],
            'zip_file' => ['required_if:deployment_source,zip', 'nullable', 'file', 'mimes:zip', 'max:524288'],
            'aliases' => ['nullable', 'array'],
            'aliases.*' => ['string', 'max:255'],
            'auto_ssl' => ['nullable', 'boolean'],
            'ssl_email' => ['nullable', 'email'],

            // Database Creation
            'create_database' => ['nullable', 'boolean'],
            'db_engine' => ['nullable', 'string', 'in:mysql,postgresql'],
            'db_name' => ['nullable', 'string', 'max:64', 'regex:/^[a-zA-Z0-9_]+$/'],
            'db_username' => ['nullable', 'string', 'max:64', 'regex:/^[a-zA-Z0-9_]+$/'],
            'db_password' => ['nullable', 'string', 'min:8'],

            // Laravel Automated Setup
            'auto_setup_laravel' => ['nullable', 'boolean'],
            'setup_env' => ['nullable', 'boolean'],
            'run_composer' => ['nullable', 'boolean'],
            'run_key_generate' => ['nullable', 'boolean'],
            'run_migrations' => ['nullable', 'boolean'],
            'run_seeders' => ['nullable', 'boolean'],
            'run_npm_build' => ['nullable', 'boolean'],
            'run_optimize' => ['nullable', 'boolean'],
        ]);

        $domain = strtolower($validated['domain']);
        $deploymentSource = $validated['deployment_source'] ?? 'empty';
        $projectType = $validated['project_type'] ?? 'laravel';
        $gitRepo = $validated['git_repository'] ?? null;
        $gitBranch = !empty($validated['git_branch']) ? $validated['git_branch'] : 'main';
        $gitAuthType = $validated['git_auth_type'] ?? 'none';
        $gitToken = $validated['git_token'] ?? null;
        $gitTokenUser = $validated['git_token_user'] ?? null;
        $gitSshPrivateKey = $validated['git_ssh_private_key'] ?? null;
        $gitSshPublicKey = $validated['git_ssh_public_key'] ?? null;
        $domainSlug = Str::slug(explode('.', $domain)[0], '_');
        $systemUser = 'kp_' . $domainSlug;

        $aliases = $request->input('aliases', []);
        if (is_string($aliases)) {
            $aliases = array_filter(array_map('trim', explode(',', $aliases)));
        }
        if (!is_array($aliases)) {
            $aliases = [];
        }
        $aliases = array_values(array_filter(array_map('strval', $aliases)));

        $autoSsl = $request->boolean('auto_ssl', false);
        $createDatabase = $request->boolean('create_database', false);
        $autoSetupLaravel = $request->boolean('auto_setup_laravel', true);

        // Auto document root determination
        if (empty($validated['document_root']) || ($projectType === 'laravel' && rtrim($validated['document_root'], '/') === "/var/www/{$domain}")) {
            if ($projectType === 'laravel') {
                $docRoot = "/var/www/{$domain}/public";
            } else {
                $docRoot = "/var/www/{$domain}";
            }
        } else {
            $docRoot = $validated['document_root'];
        }

        $fullZipPath = null;
        if ($deploymentSource === 'zip' && $request->hasFile('zip_file')) {
            $zipFile = $request->file('zip_file');
            $storedName = 'upload_' . Str::uuid()->toString() . '.zip';
            $tempZipPath = $zipFile->storeAs('temp-zips', $storedName, 'local');
            $fullZipPath = Storage::disk('local')->path($tempZipPath);
        }

        // Database Provisioning variables
        $dbName = null;
        $dbUsername = null;
        $dbPassword = null;
        $dbEngine = $validated['db_engine'] ?? 'mysql';

        try {
            // 1. Provision Database if requested
            if ($createDatabase) {
                $dbName = strtolower($validated['db_name'] ?? ('db_' . $domainSlug));
                $dbUsername = strtolower($validated['db_username'] ?? ('u_' . $domainSlug));
                $dbPassword = $validated['db_password'] ?: Str::password(16);
                $charset = $dbEngine === 'mysql' ? 'utf8mb4' : 'UTF8';
                $collation = $dbEngine === 'mysql' ? 'utf8mb4_unicode_ci' : 'en_US.UTF-8';

                // Provision on server via Agent
                $this->agentClient->createDatabase([
                    'engine' => (string) $dbEngine,
                    'name' => (string) $dbName,
                    'charset' => (string) $charset,
                    'collation' => (string) $collation,
                ]);

                $this->agentClient->createDatabaseUser([
                    'engine' => (string) $dbEngine,
                    'username' => (string) $dbUsername,
                    'host' => 'localhost',
                    'password' => (string) $dbPassword,
                ]);

                $this->agentClient->grantDatabaseAccess([
                    'engine' => (string) $dbEngine,
                    'database' => (string) $dbName,
                    'username' => (string) $dbUsername,
                    'host' => 'localhost',
                    'permissions' => 'all',
                ]);

                // Save to local database records
                $dbRecord = Database::create([
                    'engine' => $dbEngine,
                    'name' => $dbName,
                    'character_set' => $charset,
                    'collation' => $collation,
                ]);

                $dbUserRecord = DatabaseUser::create([
                    'engine' => $dbEngine,
                    'username' => $dbUsername,
                    'host' => 'localhost',
                ]);

                DatabaseAccess::create([
                    'database_id' => $dbRecord->id,
                    'database_user_id' => $dbUserRecord->id,
                    'permissions' => 'all',
                ]);
            }

            // 2. Prepare Laravel Automated Post-Setup configuration
            $laravelSetup = null;
            if ($projectType === 'laravel' && $autoSetupLaravel) {
                $scheme = $autoSsl ? 'https' : 'http';
                $envVars = [
                    'APP_NAME' => ucfirst(explode('.', $domain)[0]),
                    'APP_ENV' => 'production',
                    'APP_DEBUG' => 'false',
                    'APP_URL' => "{$scheme}://{$domain}",
                ];

                if ($createDatabase && !empty($dbName)) {
                    $envVars['DB_CONNECTION'] = $dbEngine === 'postgresql' ? 'pgsql' : 'mysql';
                    $envVars['DB_HOST'] = '127.0.0.1';
                    $envVars['DB_PORT'] = $dbEngine === 'postgresql' ? '5432' : '3306';
                    $envVars['DB_DATABASE'] = (string) $dbName;
                    $envVars['DB_USERNAME'] = (string) $dbUsername;
                    $envVars['DB_PASSWORD'] = (string) $dbPassword;
                }

                $laravelSetup = [
                    'enabled' => true,
                    'setup_env' => $request->boolean('setup_env', true),
                    'env_vars' => $envVars,
                    'run_composer' => $request->boolean('run_composer', true),
                    'run_key_generate' => $request->boolean('run_key_generate', true),
                    'run_migrations' => $request->boolean('run_migrations', true),
                    'run_seeders' => $request->boolean('run_seeders', false),
                    'run_npm_build' => $request->boolean('run_npm_build', true),
                    'run_optimize' => $request->boolean('run_optimize', true),
                ];
            }

            // 3. Provision Virtual Host on server via Agent
            $agentPayload = [
                'domain' => (string) $domain,
                'aliases' => $aliases,
                'php_version' => (string) $validated['php_version'],
                'document_root' => (string) $docRoot,
                'system_user' => (string) $systemUser,
                'ssl_enabled' => false,
                'force_https' => false,
                'deployment_source' => (string) $deploymentSource,
                'project_type' => (string) $projectType,
                'git_repository' => $gitRepo ? (string) $gitRepo : '',
                'git_branch' => (string) $gitBranch,
                'git_auth_type' => (string) $gitAuthType,
                'git_token' => $gitToken ? (string) $gitToken : '',
                'git_token_user' => $gitTokenUser ? (string) $gitTokenUser : '',
                'git_ssh_private_key' => $gitSshPrivateKey ? (string) $gitSshPrivateKey : '',
                'git_ssh_public_key' => $gitSshPublicKey ? (string) $gitSshPublicKey : '',
                'laravel_setup' => $laravelSetup,
            ];

            if (!empty($fullZipPath) && file_exists($fullZipPath)) {
                $agentPayload['zip_path'] = $fullZipPath;
            }

            $agentRes = $this->agentClient->createWebsite($agentPayload);
            if (!empty($agentRes['document_root'])) {
                $docRoot = $agentRes['document_root'];
            }

            // 4. Persist Website record
            $website = Website::create([
                'domain' => $domain,
                'aliases' => $aliases,
                'php_version' => $validated['php_version'],
                'document_root' => $docRoot,
                'system_user' => $systemUser,
                'ssl_enabled' => false,
                'force_https' => false,
                'status' => 'active',
                'deployment_source' => $deploymentSource,
                'project_type' => $projectType,
                'git_repository' => $gitRepo,
                'git_branch' => $gitBranch,
                'git_auth_type' => $gitAuthType,
                'git_token' => $gitToken,
                'git_token_user' => $gitTokenUser,
                'git_ssh_private_key' => $gitSshPrivateKey,
                'git_ssh_public_key' => $gitSshPublicKey,
                'last_deployed_at' => in_array($deploymentSource, ['zip', 'git']) ? now() : null,
            ]);

            // 5. Create Primary Domain record
            Domain::create([
                'website_id' => $website->id,
                'domain' => $domain,
                'is_primary' => true,
                'ssl_status' => 'pending',
            ]);

            // 6. Record Setup / Deployment Output Logs
            $setupResult = $agentRes['setup_result'] ?? null;
            if (!empty($setupResult) && !empty($setupResult['log_output'])) {
                Deployment::create([
                    'website_id' => $website->id,
                    'branch' => $gitBranch,
                    'status' => ($setupResult['success'] ?? true) ? 'success' : 'failed',
                    'trigger_source' => 'manual',
                    'commit_message' => 'Automated initial Laravel provisioning & setup',
                    'log_output' => $setupResult['log_output'],
                    'duration_seconds' => $setupResult['duration_seconds'] ?? 1,
                    'initiated_by_user_id' => $request->user()->id,
                ]);
            } elseif ($deploymentSource === 'git' && !empty($gitRepo)) {
                Deployment::create([
                    'website_id' => $website->id,
                    'branch' => $gitBranch,
                    'status' => 'success',
                    'trigger_source' => 'manual',
                    'commit_message' => 'Initial repository clone upon website creation',
                    'log_output' => "Repository {$gitRepo} ({$gitBranch}) cloned successfully into /var/www/{$domain}.",
                    'duration_seconds' => 1,
                    'initiated_by_user_id' => $request->user()->id,
                ]);
            }

            // 7. Auto-issue SSL if requested
            if ($autoSsl) {
                try {
                    $email = $validated['ssl_email'] ?: $request->user()->email;
                    $sslRes = $this->agentClient->issueSsl([
                        'domain' => $domain,
                        'aliases' => $aliases,
                        'email' => $email,
                        'document_root' => $docRoot,
                        'php_version' => $validated['php_version'],
                        'system_user' => $systemUser,
                        'force_https' => true,
                    ]);

                    $website->update([
                        'ssl_enabled' => true,
                        'force_https' => true,
                    ]);

                    SslCertificate::create([
                        'website_id' => $website->id,
                        'domain' => $domain,
                        'issuer' => $sslRes['issuer'] ?? "Let's Encrypt",
                        'cert_path' => $sslRes['cert_path'] ?? null,
                        'key_path' => $sslRes['key_path'] ?? null,
                        'valid_from' => $sslRes['valid_from'] ?? now(),
                        'valid_until' => $sslRes['valid_until'] ?? now()->addDays(90),
                        'status' => 'valid',
                    ]);
                } catch (Exception $sslEx) {
                    // Non-fatal: website is still created, SSL can be retried later
                }
            }

            ActivityLog::create([
                'user_id' => $request->user()->id,
                'user_email' => $request->user()->email,
                'ip_address' => $request->ip() ?: '127.0.0.1',
                'user_agent' => $request->userAgent(),
                'action' => 'website.create',
                'resource_type' => 'website',
                'resource_id' => (string) $website->id,
                'status' => 'success',
                'payload_summary' => [
                    'domain' => $domain,
                    'php_version' => $validated['php_version'],
                    'deployment_source' => $deploymentSource,
                    'document_root' => $docRoot,
                    'ssl_enabled' => $website->ssl_enabled,
                ],
            ]);

            return redirect()->route('websites.show', $website)->with('success', "Website {$domain} provisioned successfully.");
        } catch (Exception $e) {
            return back()->withInput()->with('error', "Failed to provision website: " . $e->getMessage());
        } finally {
            if (!empty($fullZipPath) && file_exists($fullZipPath)) {
                @unlink($fullZipPath);
            }
        }
    }

    public function show(Website $website): Response
    {
        $website->load(['domains', 'deployments' => fn ($q) => $q->take(10), 'sslCertificate']);

        return Inertia::render('Websites/Show', [
            'website' => $website,
        ]);
    }

    public function update(Request $request, Website $website): RedirectResponse
    {
        $validated = $request->validate([
            'document_root' => ['required', 'string', 'max:255'],
            'force_https' => ['nullable', 'boolean'],
        ]);

        $website->update($validated);

        ActivityLog::create([
            'user_id' => $request->user()->id,
            'user_email' => $request->user()->email,
            'ip_address' => $request->ip() ?: '127.0.0.1',
            'user_agent' => $request->userAgent(),
            'action' => 'website.update',
            'resource_type' => 'website',
            'resource_id' => (string) $website->id,
            'status' => 'success',
            'payload_summary' => $validated,
        ]);

        return back()->with('success', "Website settings updated successfully.");
    }

    public function switchPhp(Request $request, Website $website): RedirectResponse
    {
        $validated = $request->validate([
            'php_version' => ['required', 'string', 'in:8.3,8.4,none'],
        ]);

        $oldVersion = $website->php_version;
        $newVersion = $validated['php_version'];

        if ($oldVersion === $newVersion) {
            return back()->with('info', "Website is already using PHP {$newVersion}.");
        }

        $aliases = $website->aliases ?? [];
        if (!is_array($aliases)) {
            $aliases = [];
        }
        $aliases = array_values(array_filter(array_map('strval', $aliases)));

        try {
            $this->agentClient->switchPhpVersion($website->domain, [
                'new_php_version' => (string) $newVersion,
                'old_php_version' => (string) $oldVersion,
                'system_user' => (string) $website->system_user,
                'document_root' => (string) $website->document_root,
                'aliases' => $aliases,
                'ssl_enabled' => (bool) $website->ssl_enabled,
                'force_https' => (bool) $website->force_https,
                'cert_path' => (string) ($website->sslCertificate?->cert_path ?? ''),
                'key_path' => (string) ($website->sslCertificate?->key_path ?? ''),
            ]);

            $website->update(['php_version' => $newVersion]);

            ActivityLog::create([
                'user_id' => $request->user()->id,
                'user_email' => $request->user()->email,
                'ip_address' => $request->ip() ?: '127.0.0.1',
                'user_agent' => $request->userAgent(),
                'action' => 'website.switch_php',
                'resource_type' => 'website',
                'resource_id' => (string) $website->id,
                'status' => 'success',
                'payload_summary' => [
                    'domain' => $website->domain,
                    'from' => $oldVersion,
                    'to' => $newVersion,
                ],
            ]);

            return back()->with('success', "Switched runtime to PHP {$newVersion} successfully.");
        } catch (Exception $e) {
            return back()->with('error', "Failed to switch PHP version: " . $e->getMessage());
        }
    }

    public function issueSsl(Request $request, Website $website): RedirectResponse
    {
        $validated = $request->validate([
            'email' => ['nullable', 'email'],
            'force_https' => ['nullable', 'boolean'],
        ]);

        $email = $validated['email'] ?: $request->user()->email;
        $forceHttps = $request->boolean('force_https', true);
        $aliases = $website->aliases ?? [];
        if (!is_array($aliases)) {
            $aliases = [];
        }
        $aliases = array_values(array_filter(array_map('strval', $aliases)));

        try {
            $sslRes = $this->agentClient->issueSsl([
                'domain' => (string) $website->domain,
                'aliases' => $aliases,
                'email' => (string) $email,
                'document_root' => (string) $website->document_root,
                'php_version' => (string) $website->php_version,
                'system_user' => (string) $website->system_user,
                'force_https' => $forceHttps,
            ]);

            $website->update([
                'ssl_enabled' => true,
                'force_https' => $forceHttps,
            ]);

            SslCertificate::updateOrCreate(
                ['website_id' => $website->id, 'domain' => $website->domain],
                [
                    'issuer' => $sslRes['issuer'] ?? "Let's Encrypt",
                    'cert_path' => $sslRes['cert_path'] ?? null,
                    'key_path' => $sslRes['key_path'] ?? null,
                    'valid_from' => $sslRes['valid_from'] ?? now(),
                    'valid_until' => $sslRes['valid_until'] ?? now()->addDays(90),
                    'status' => 'valid',
                    'auto_renew' => true,
                ]
            );

            ActivityLog::create([
                'user_id' => $request->user()->id,
                'user_email' => $request->user()->email,
                'ip_address' => $request->ip() ?: '127.0.0.1',
                'user_agent' => $request->userAgent(),
                'action' => 'website.issue_ssl',
                'resource_type' => 'website',
                'resource_id' => (string) $website->id,
                'status' => 'success',
                'payload_summary' => ['domain' => $website->domain],
            ]);

            return back()->with('success', "SSL Certificate issued and applied successfully.");
        } catch (Exception $e) {
            return back()->with('error', "Failed to issue SSL: " . $e->getMessage());
        }
    }

    public function logs(Request $request, Website $website, string $type): JsonResponse
    {
        if ($type !== 'access' && $type !== 'error') {
            return response()->json(['error' => 'Invalid log type'], 400);
        }

        $lines = (int) $request->query('lines', 100);

        try {
            $logData = $this->agentClient->getWebsiteLogs($website->domain, $type, $lines);
            return response()->json([
                'success' => true,
                'domain' => $website->domain,
                'type' => $type,
                'lines' => $logData['lines'] ?? [],
            ]);
        } catch (Exception $e) {
            return response()->json([
                'success' => false,
                'error' => $e->getMessage(),
            ], 500);
        }
    }

    public function traffic(Request $request, Website $website): JsonResponse
    {
        $period = $request->query('period', '24h');
        if (!in_array($period, ['today', '24h', '7d', '30d'])) {
            $period = '24h';
        }

        try {
            $trafficData = $this->agentClient->getWebsiteTraffic($website->domain, $period);
            return response()->json([
                'success' => true,
                'data' => $trafficData,
            ]);
        } catch (Exception $e) {
            return response()->json([
                'success' => false,
                'error' => $e->getMessage(),
            ], 500);
        }
    }

    public function destroy(Request $request, Website $website): RedirectResponse
    {
        $domain = $website->domain;

        try {
            $this->agentClient->deleteWebsite($domain, [
                'php_version' => $website->php_version,
                'system_user' => $website->system_user,
            ]);
        } catch (Exception $e) {
            // Log warning but continue local cleanup
        }

        $website->delete();

        ActivityLog::create([
            'user_id' => $request->user()->id,
            'user_email' => $request->user()->email,
            'ip_address' => $request->ip() ?: '127.0.0.1',
            'user_agent' => $request->userAgent(),
            'action' => 'website.delete',
            'resource_type' => 'website',
            'resource_id' => (string) $website->id,
            'status' => 'success',
            'payload_summary' => ['domain' => $domain],
        ]);

        return redirect()->route('websites.index')->with('success', "Website {$domain} deleted successfully.");
    }
}
