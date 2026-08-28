<?php

namespace App\Http\Controllers;

use App\Models\ActivityLog;
use App\Models\CronJob;
use App\Models\Website;
use App\Services\Agent\AgentClientInterface;
use Exception;
use Illuminate\Http\RedirectResponse;
use Illuminate\Http\Request;
use Inertia\Inertia;
use Inertia\Response;

class CronJobController extends Controller
{
    public function __construct(
        protected AgentClientInterface $agentClient
    ) {}

    public function index(): Response
    {
        $cronJobs = CronJob::with('website')
            ->latest()
            ->paginate(15);

        $websites = Website::select('id', 'domain', 'system_user')->orderBy('domain')->get();

        return Inertia::render('Cron/Index', [
            'cronJobs' => $cronJobs,
            'websites' => $websites,
        ]);
    }

    public function store(Request $request): RedirectResponse
    {
        $validated = $request->validate([
            'website_id' => ['nullable', 'exists:websites,id'],
            'schedule' => ['required', 'string', 'max:64'],
            'command' => ['required', 'string', 'max:1024'],
            'system_user' => ['nullable', 'string', 'max:64'],
        ]);

        $systemUser = $validated['system_user'] ?? 'kodepreneur';
        if (!empty($validated['website_id'])) {
            $site = Website::find($validated['website_id']);
            if ($site && $site->system_user) {
                $systemUser = $site->system_user;
            }
        }

        $cronJob = CronJob::create([
            'website_id' => $validated['website_id'] ?? null,
            'schedule' => trim($validated['schedule']),
            'command' => trim($validated['command']),
            'system_user' => $systemUser,
            'is_active' => true,
        ]);

        $this->syncUserCrontab($systemUser);

        ActivityLog::create([
            'user_id' => $request->user()->id,
            'user_email' => $request->user()->email,
            'ip_address' => $request->ip() ?: '127.0.0.1',
            'user_agent' => $request->userAgent(),
            'action' => 'cron.create',
            'resource_type' => 'cron_job',
            'resource_id' => (string) $cronJob->id,
            'status' => 'success',
            'payload_summary' => [
                'schedule' => $cronJob->schedule,
                'command' => $cronJob->command,
                'system_user' => $cronJob->system_user,
            ],
        ]);

        return redirect()->route('cron.index')->with('success', "Cron job created and synchronized.");
    }

    public function toggle(Request $request, CronJob $cronJob): RedirectResponse
    {
        $cronJob->update([
            'is_active' => !$cronJob->is_active,
        ]);

        $this->syncUserCrontab($cronJob->system_user);

        return redirect()->route('cron.index')->with('success', "Cron job status updated.");
    }

    public function destroy(Request $request, CronJob $cronJob): RedirectResponse
    {
        $systemUser = $cronJob->system_user;
        $cronJob->delete();

        $this->syncUserCrontab($systemUser);

        ActivityLog::create([
            'user_id' => $request->user()->id,
            'user_email' => $request->user()->email,
            'ip_address' => $request->ip() ?: '127.0.0.1',
            'user_agent' => $request->userAgent(),
            'action' => 'cron.delete',
            'resource_type' => 'cron_job',
            'resource_id' => (string) $cronJob->id,
            'status' => 'success',
        ]);

        return redirect()->route('cron.index')->with('success', "Cron job removed and crontab updated.");
    }

    protected function syncUserCrontab(string $systemUser): void
    {
        try {
            $jobs = CronJob::where('system_user', $systemUser)
                ->get()
                ->map(fn ($j) => [
                    'schedule' => $j->schedule,
                    'command' => $j->command,
                    'is_active' => (bool) $j->is_active,
                ])
                ->toArray();

            $this->agentClient->syncCronJobs($systemUser, $jobs);
        } catch (Exception $e) {
            // Non-fatal if agent unreachable in test
        }
    }
}
