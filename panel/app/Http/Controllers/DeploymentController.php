<?php

namespace App\Http\Controllers;

use App\Models\ActivityLog;
use App\Models\Deployment;
use App\Models\Website;
use App\Services\Agent\AgentClientInterface;
use Exception;
use Illuminate\Http\RedirectResponse;
use Illuminate\Http\Request;
use Inertia\Inertia;
use Inertia\Response;

class DeploymentController extends Controller
{
    public function __construct(
        protected AgentClientInterface $agentClient
    ) {}

    public function index(): Response
    {
        $deployments = Deployment::with(['website', 'initiatedBy'])
            ->latest()
            ->paginate(20);

        return Inertia::render('Deployments/Index', [
            'deployments' => $deployments,
        ]);
    }

    public function show(Deployment $deployment): Response
    {
        $deployment->load(['website', 'initiatedBy']);

        return Inertia::render('Deployments/Show', [
            'deployment' => $deployment,
        ]);
    }

    public function trigger(Request $request, Website $website): RedirectResponse
    {
        $branch = $website->git_branch ?: 'main';
        $workingDir = $website->document_root ? dirname($website->document_root) : "/var/www/{$website->domain}";
        $user = $website->system_user ?: 'kodepreneur';

        $commands = [
            "git pull origin {$branch}",
        ];

        if ($website->php_version && $website->php_version !== 'none') {
            $commands[] = 'if [ -f composer.json ]; then composer install --no-dev --prefer-dist --optimize-autoloader --no-interaction; fi';
            $commands[] = 'if [ -f artisan ]; then php artisan migrate --force; php artisan optimize:clear; php artisan config:cache; php artisan route:cache; php artisan view:cache; fi';
        }

        $commands[] = 'if [ -f package.json ]; then if [ ! -d node_modules/vite ] && [ -f node_modules/.bin/vite ]; then rm -rf node_modules package-lock.json; fi; NODE_ENV=development npm install --include=dev --no-audit; npm run build; fi';

        try {
            $result = $this->agentClient->executeDeployment([
                'system_user' => $user,
                'working_dir' => $workingDir,
                'repository' => $website->git_repository ?: '',
                'branch' => $branch,
                'commands' => $commands,
                'timeout_sec' => 600,
            ]);

            $isSuccess = ($result['success'] ?? false) && ($result['exit_code'] ?? 1) === 0;

            $deployment = Deployment::create([
                'website_id' => $website->id,
                'commit_hash' => $result['commit_hash'] ?? null,
                'commit_message' => $result['commit_message'] ?? null,
                'branch' => $branch,
                'status' => $isSuccess ? 'success' : 'failed',
                'trigger_source' => 'manual',
                'log_output' => $result['log_output'] ?? 'No output returned.',
                'duration_seconds' => $result['duration_seconds'] ?? 1,
                'initiated_by_user_id' => $request->user()->id,
            ]);

            $website->update([
                'last_deployed_at' => now(),
            ]);

            ActivityLog::create([
                'user_id' => $request->user()->id,
                'user_email' => $request->user()->email,
                'ip_address' => $request->ip() ?: '127.0.0.1',
                'user_agent' => $request->userAgent(),
                'action' => 'deployment.trigger',
                'resource_type' => 'deployment',
                'resource_id' => (string) $deployment->id,
                'status' => $isSuccess ? 'success' : 'failure',
                'payload_summary' => [
                    'website' => $website->domain,
                    'branch' => $branch,
                    'status' => $deployment->status,
                ],
            ]);

            return redirect()->route('deployments.show', $deployment)->with('success', "Deployment #{$deployment->id} completed with status: {$deployment->status}.");
        } catch (Exception $e) {
            $deployment = Deployment::create([
                'website_id' => $website->id,
                'branch' => $branch,
                'status' => 'failed',
                'trigger_source' => 'manual',
                'log_output' => "Agent deployment execution failed: " . $e->getMessage(),
                'duration_seconds' => 1,
                'initiated_by_user_id' => $request->user()->id,
            ]);

            return redirect()->route('deployments.show', $deployment)->with('error', "Deployment failed: " . $e->getMessage());
        }
    }

    public function rollback(Request $request, Deployment $deployment): RedirectResponse
    {
        $website = $deployment->website;
        if (!$deployment->commit_hash) {
            return back()->with('error', "Cannot rollback: No commit hash recorded for deployment #{$deployment->id}.");
        }

        $workingDir = $website->document_root ? dirname($website->document_root) : "/var/www/{$website->domain}";
        $user = $website->system_user ?: 'kodepreneur';

        $commands = [
            "git checkout {$deployment->commit_hash}",
            'if [ -f composer.json ]; then composer install --no-dev --prefer-dist --optimize-autoloader --no-interaction; fi',
            'if [ -f artisan ]; then php artisan optimize:clear; php artisan config:cache; fi',
        ];

        try {
            $result = $this->agentClient->executeDeployment([
                'system_user' => $user,
                'working_dir' => $workingDir,
                'branch' => $deployment->branch,
                'commands' => $commands,
                'timeout_sec' => 600,
            ]);

            $newDeploy = Deployment::create([
                'website_id' => $website->id,
                'commit_hash' => $deployment->commit_hash,
                'commit_message' => "Rollback to #{$deployment->id} ({$deployment->commit_hash})",
                'branch' => $deployment->branch,
                'status' => 'success',
                'trigger_source' => 'rollback',
                'log_output' => $result['log_output'] ?? 'Rollback complete.',
                'duration_seconds' => $result['duration_seconds'] ?? 1,
                'initiated_by_user_id' => $request->user()->id,
            ]);

            return redirect()->route('deployments.show', $newDeploy)->with('success', "Rollback deployment #{$newDeploy->id} completed successfully.");
        } catch (Exception $e) {
            return back()->with('error', "Rollback failed: " . $e->getMessage());
        }
    }
}
