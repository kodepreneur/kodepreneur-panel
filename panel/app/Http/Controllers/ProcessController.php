<?php

namespace App\Http\Controllers;

use App\Models\ActivityLog;
use App\Models\Process;
use App\Models\Website;
use App\Services\Agent\AgentClientInterface;
use Exception;
use Illuminate\Http\RedirectResponse;
use Illuminate\Http\Request;
use Illuminate\Support\Str;
use Inertia\Inertia;
use Inertia\Response;

class ProcessController extends Controller
{
    public function __construct(
        protected AgentClientInterface $agentClient
    ) {}

    public function index(): Response
    {
        $processes = Process::with('website')
            ->latest()
            ->paginate(15);

        $websites = Website::select('id', 'domain', 'system_user', 'document_root')->orderBy('domain')->get();

        return Inertia::render('Processes/Index', [
            'processes' => $processes,
            'websites' => $websites,
        ]);
    }

    public function store(Request $request): RedirectResponse
    {
        $validated = $request->validate([
            'name' => ['required', 'string', 'max:64'],
            'website_id' => ['nullable', 'exists:websites,id'],
            'command' => ['required', 'string', 'max:1024'],
            'system_user' => ['nullable', 'string', 'max:64'],
            'instances' => ['nullable', 'integer', 'min:1', 'max:10'],
        ]);

        $procId = Str::slug($validated['name']) . '-' . Str::random(5);
        $systemUser = $validated['system_user'] ?? 'kodepreneur';
        $workingDir = '/var/www';

        if (!empty($validated['website_id'])) {
            $site = Website::find($validated['website_id']);
            if ($site) {
                if ($site->system_user) $systemUser = $site->system_user;
                if ($site->document_root) $workingDir = dirname($site->document_root);
            }
        }

        try {
            $agentRes = $this->agentClient->createProcess([
                'id' => $procId,
                'name' => $validated['name'],
                'command' => $validated['command'],
                'system_user' => $systemUser,
                'working_dir' => $workingDir,
                'auto_restart' => true,
            ]);

            $proc = Process::create([
                'name' => $validated['name'],
                'website_id' => $validated['website_id'] ?? null,
                'command' => $validated['command'],
                'system_user' => $systemUser,
                'instances' => $validated['instances'] ?? 1,
                'status' => 'running',
                'unit_file_path' => $agentRes['unit_path'] ?? "/etc/systemd/system/kp-proc-{$procId}.service",
            ]);

            ActivityLog::create([
                'user_id' => $request->user()->id,
                'user_email' => $request->user()->email,
                'ip_address' => $request->ip() ?: '127.0.0.1',
                'user_agent' => $request->userAgent(),
                'action' => 'process.create',
                'resource_type' => 'process',
                'resource_id' => (string) $proc->id,
                'status' => 'success',
                'payload_summary' => [
                    'name' => $proc->name,
                    'command' => $proc->command,
                ],
            ]);

            return redirect()->route('processes.index')->with('success', "Managed process '{$proc->name}' created and started.");
        } catch (Exception $e) {
            return back()->withInput()->with('error', "Failed to create process: " . $e->getMessage());
        }
    }

    public function manage(Request $request, Process $process, string $action): RedirectResponse
    {
        if (!in_array($action, ['start', 'stop', 'restart'])) {
            return back()->with('error', "Invalid action.");
        }

        // Extract ID from unit_file_path or slug name
        $id = basename($process->unit_file_path, '.service');
        $id = str_replace('kp-proc-', '', $id);

        try {
            $this->agentClient->manageProcess($id, $action);

            $newStatus = $action === 'stop' ? 'stopped' : 'running';
            $process->update(['status' => $newStatus]);

            ActivityLog::create([
                'user_id' => $request->user()->id,
                'user_email' => $request->user()->email,
                'ip_address' => $request->ip() ?: '127.0.0.1',
                'user_agent' => $request->userAgent(),
                'action' => "process.{$action}",
                'resource_type' => 'process',
                'resource_id' => (string) $process->id,
                'status' => 'success',
            ]);

            return redirect()->route('processes.index')->with('success', "Action '{$action}' executed on {$process->name}.");
        } catch (Exception $e) {
            return back()->with('error', "Process action failed: " . $e->getMessage());
        }
    }

    public function destroy(Request $request, Process $process): RedirectResponse
    {
        $id = basename($process->unit_file_path, '.service');
        $id = str_replace('kp-proc-', '', $id);

        try {
            $this->agentClient->deleteProcess($id);
        } catch (Exception $e) {
            // Non-fatal
        }

        $process->delete();

        ActivityLog::create([
            'user_id' => $request->user()->id,
            'user_email' => $request->user()->email,
            'ip_address' => $request->ip() ?: '127.0.0.1',
            'user_agent' => $request->userAgent(),
            'action' => 'process.delete',
            'resource_type' => 'process',
            'resource_id' => (string) $process->id,
            'status' => 'success',
        ]);

        return redirect()->route('processes.index')->with('success', "Process removed.");
    }
}
