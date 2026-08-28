<?php

namespace App\Http\Controllers;

use App\Models\ActivityLog;
use App\Models\Database;
use App\Models\Deployment;
use App\Models\Website;
use App\Services\Agent\AgentClientInterface;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;
use Inertia\Inertia;
use Inertia\Response;

class DashboardController extends Controller
{
    public function __construct(
        protected AgentClientInterface $agentClient
    ) {}

    public function index(): Response
    {
        $systemInfo = $this->agentClient->getSystemInfo();
        $systemMetrics = $this->agentClient->getSystemMetrics();

        $servicesToInspect = ['nginx', 'php8.3-fpm', 'php8.4-fpm', 'mysql', 'postgresql', 'kodepreneur-agent'];
        $services = [];
        foreach ($servicesToInspect as $svc) {
            $services[] = $this->agentClient->getServiceStatus($svc);
        }

        $counts = [
            'websites' => Website::count(),
            'databases' => Database::count(),
            'deployments' => Deployment::count(),
        ];

        $recentActivity = ActivityLog::latest('created_at')
            ->take(5)
            ->get();

        return Inertia::render('Dashboard/Index', [
            'systemInfo' => $systemInfo,
            'systemMetrics' => $systemMetrics,
            'services' => $services,
            'counts' => $counts,
            'recentActivity' => $recentActivity,
        ]);
    }

    public function metrics(): JsonResponse
    {
        $metrics = $this->agentClient->getSystemMetrics();
        return response()->json([
            'success' => true,
            'data' => $metrics,
        ]);
    }
}
