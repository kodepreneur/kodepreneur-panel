<?php

namespace App\Http\Controllers;

use App\Models\ActivityLog;
use App\Models\FirewallRule;
use App\Services\Agent\AgentClientInterface;
use Exception;
use Illuminate\Http\RedirectResponse;
use Illuminate\Http\Request;
use Inertia\Inertia;
use Inertia\Response;

class FirewallController extends Controller
{
    public function __construct(
        protected AgentClientInterface $agentClient
    ) {}

    public function index(): Response
    {
        $rules = FirewallRule::latest()->paginate(20);

        return Inertia::render('Firewall/Index', [
            'rules' => $rules,
        ]);
    }

    public function store(Request $request): RedirectResponse
    {
        $validated = $request->validate([
            'port' => ['required', 'integer', 'min:1', 'max:65535'],
            'protocol' => ['required', 'string', 'in:tcp,udp,both'],
            'action' => ['required', 'string', 'in:allow,deny'],
            'source_ip' => ['nullable', 'string', 'max:64'],
            'description' => ['nullable', 'string', 'max:255'],
        ]);

        try {
            $this->agentClient->addFirewallRule([
                'port' => $validated['port'],
                'protocol' => $validated['protocol'],
                'action' => $validated['action'],
                'source_ip' => $validated['source_ip'] ?? '',
                'description' => $validated['description'] ?? '',
            ]);

            $rule = FirewallRule::create([
                'port' => $validated['port'],
                'protocol' => $validated['protocol'],
                'action' => $validated['action'],
                'source_ip' => $validated['source_ip'] ?? null,
                'description' => $validated['description'] ?? null,
            ]);

            ActivityLog::create([
                'user_id' => $request->user()->id,
                'user_email' => $request->user()->email,
                'ip_address' => $request->ip() ?: '127.0.0.1',
                'user_agent' => $request->userAgent(),
                'action' => 'firewall.add_rule',
                'resource_type' => 'firewall_rule',
                'resource_id' => (string) $rule->id,
                'status' => 'success',
                'payload_summary' => [
                    'port' => $rule->port,
                    'protocol' => $rule->protocol,
                    'action' => $rule->action,
                ],
            ]);

            return redirect()->route('firewall.index')->with('success', "Firewall rule for port {$rule->port}/{$rule->protocol} added.");
        } catch (Exception $e) {
            return back()->withInput()->with('error', "Failed to add firewall rule: " . $e->getMessage());
        }
    }

    public function destroy(Request $request, FirewallRule $firewallRule): RedirectResponse
    {
        try {
            $this->agentClient->deleteFirewallRule([
                'port' => $firewallRule->port,
                'protocol' => $firewallRule->protocol,
                'action' => $firewallRule->action,
            ]);
        } catch (Exception $e) {
            // Non-fatal
        }

        $port = $firewallRule->port;
        $protocol = $firewallRule->protocol;
        $firewallRule->delete();

        ActivityLog::create([
            'user_id' => $request->user()->id,
            'user_email' => $request->user()->email,
            'ip_address' => $request->ip() ?: '127.0.0.1',
            'user_agent' => $request->userAgent(),
            'action' => 'firewall.delete_rule',
            'resource_type' => 'firewall_rule',
            'resource_id' => (string) $firewallRule->id,
            'status' => 'success',
        ]);

        return redirect()->route('firewall.index')->with('success', "Firewall rule for port {$port}/{$protocol} deleted.");
    }
}
