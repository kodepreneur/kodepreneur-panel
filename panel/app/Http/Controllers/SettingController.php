<?php

namespace App\Http\Controllers;

use App\Models\ActivityLog;
use App\Models\ServerSetting;
use Illuminate\Http\RedirectResponse;
use Illuminate\Http\Request;
use Inertia\Inertia;
use Inertia\Response;

class SettingController extends Controller
{
    public function index(): Response
    {
        $settings = ServerSetting::all()->pluck('value', 'key');

        return Inertia::render('Settings/Index', [
            'settings' => $settings,
        ]);
    }

    public function update(Request $request): RedirectResponse
    {
        $validated = $request->validate([
            'panel_name' => ['required', 'string', 'max:100'],
            'panel_port' => ['required', 'numeric', 'min:1', 'max:65535'],
            'default_php_version' => ['required', 'string', 'in:8.3,8.4'],
        ]);

        foreach ($validated as $key => $value) {
            ServerSetting::set($key, $value);
        }

        ActivityLog::create([
            'user_id' => $request->user()->id,
            'user_email' => $request->user()->email,
            'ip_address' => $request->ip() ?: '127.0.0.1',
            'user_agent' => $request->userAgent(),
            'action' => 'settings.update',
            'resource_type' => 'settings',
            'status' => 'success',
            'payload_summary' => $validated,
        ]);

        return redirect()->back()->with('success', 'Server settings updated successfully.');
    }
}
