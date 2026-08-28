<?php

namespace App\Http\Controllers;

use App\Models\ActivityLog;
use App\Services\Agent\AgentClientInterface;
use Illuminate\Http\RedirectResponse;
use Illuminate\Http\Request;
use Inertia\Inertia;
use Inertia\Response;

class ServiceController extends Controller
{
    public function __construct(
        protected AgentClientInterface $agentClient
    ) {}

    public function index(): Response
    {
        $servicesToInspect = [
            'nginx',
            'php8.3-fpm',
            'php8.4-fpm',
            'mysql',
            'postgresql',
            'kodepreneur-agent',
            'kodepreneur-worker',
            'ufw',
            'fail2ban'
        ];

        $services = [];
        foreach ($servicesToInspect as $svc) {
            $services[] = $this->agentClient->getServiceStatus($svc);
        }

        return Inertia::render('Services/Index', [
            'services' => $services,
        ]);
    }

    public function manage(Request $request, string $service, string $action): RedirectResponse
    {
        $request->validate([
            'service' => ['string'],
            'action' => ['string', 'in:start,stop,restart,reload'],
        ]);

        $res = $this->agentClient->manageService($service, $action);

        ActivityLog::create([
            'user_id' => $request->user()->id,
            'user_email' => $request->user()->email,
            'ip_address' => $request->ip() ?: '127.0.0.1',
            'user_agent' => $request->userAgent(),
            'action' => "service.{$action}",
            'resource_type' => 'service',
            'resource_id' => $service,
            'status' => 'success',
            'payload_summary' => ['output' => $res['output'] ?? ''],
        ]);

        return redirect()->back()->with('success', "Service {$service} {$action} command dispatched.");
    }
}
