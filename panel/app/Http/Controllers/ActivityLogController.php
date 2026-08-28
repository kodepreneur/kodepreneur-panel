<?php

namespace App\Http\Controllers;

use App\Models\ActivityLog;
use Inertia\Inertia;
use Inertia\Response;

class ActivityLogController extends Controller
{
    public function index(): Response
    {
        $logs = ActivityLog::latest('created_at')
            ->paginate(30);

        return Inertia::render('Activity/Index', [
            'logs' => $logs,
        ]);
    }
}
