<?php

use App\Http\Controllers\ActivityLogController;
use App\Http\Controllers\Auth\AuthController;
use App\Http\Controllers\CronJobController;
use App\Http\Controllers\DashboardController;
use App\Http\Controllers\DatabaseController;
use App\Http\Controllers\DeploymentController;
use App\Http\Controllers\FileManagerController;
use App\Http\Controllers\FirewallController;
use App\Http\Controllers\ProcessController;
use App\Http\Controllers\ServiceController;
use App\Http\Controllers\SettingController;
use App\Http\Controllers\WebsiteController;
use Illuminate\Support\Facades\Route;

// Authentication Routes
Route::middleware('guest')->group(function () {
    Route::get('/login', [AuthController::class, 'showLogin'])->name('login');
    Route::post('/login', [AuthController::class, 'login']);
});

Route::middleware('auth')->group(function () {
    Route::post('/logout', [AuthController::class, 'logout'])->name('logout');

    Route::get('/', function () {
        return redirect()->route('dashboard');
    });

    // Dashboard
    Route::get('/dashboard', [DashboardController::class, 'index'])->name('dashboard');
    Route::get('/dashboard/metrics', [DashboardController::class, 'metrics'])->name('dashboard.metrics');

    // Websites
    Route::resource('websites', WebsiteController::class);
    Route::put('/websites/{website}/php', [WebsiteController::class, 'switchPhp'])->name('websites.php');
    Route::post('/websites/{website}/ssl', [WebsiteController::class, 'issueSsl'])->name('websites.ssl');
    Route::get('/websites/{website}/logs/{type}', [WebsiteController::class, 'logs'])->name('websites.logs');
    Route::post('/websites/{website}/deploy', [DeploymentController::class, 'trigger'])->name('websites.deploy');

    // Databases
    Route::resource('databases', DatabaseController::class)->only(['index', 'store', 'destroy']);
    Route::post('/databases/users', [DatabaseController::class, 'storeUser'])->name('databases.users.store');
    Route::delete('/databases/users/{databaseUser}', [DatabaseController::class, 'destroyUser'])->name('databases.users.destroy');
    Route::post('/databases/users/{databaseUser}/password', [DatabaseController::class, 'resetPassword'])->name('databases.users.password');
    Route::post('/databases/grants', [DatabaseController::class, 'storeGrant'])->name('databases.grants.store');
    Route::delete('/databases/grants/{grant}', [DatabaseController::class, 'revokeGrant'])->name('databases.grants.destroy');

    // Deployments
    Route::resource('deployments', DeploymentController::class)->only(['index', 'show']);
    Route::post('/deployments/{deployment}/rollback', [DeploymentController::class, 'rollback'])->name('deployments.rollback');

    // Cron Jobs (Phase 6)
    Route::get('/cron', [CronJobController::class, 'index'])->name('cron.index');
    Route::post('/cron', [CronJobController::class, 'store'])->name('cron.store');
    Route::post('/cron/{cronJob}/toggle', [CronJobController::class, 'toggle'])->name('cron.toggle');
    Route::delete('/cron/{cronJob}', [CronJobController::class, 'destroy'])->name('cron.destroy');

    // Managed Processes (Phase 6)
    Route::get('/processes', [ProcessController::class, 'index'])->name('processes.index');
    Route::post('/processes', [ProcessController::class, 'store'])->name('processes.store');
    Route::post('/processes/{process}/{action}', [ProcessController::class, 'manage'])->name('processes.manage');
    Route::delete('/processes/{process}', [ProcessController::class, 'destroy'])->name('processes.destroy');

    // Firewall (Phase 6)
    Route::get('/firewall', [FirewallController::class, 'index'])->name('firewall.index');
    Route::post('/firewall', [FirewallController::class, 'store'])->name('firewall.store');
    Route::delete('/firewall/{firewallRule}', [FirewallController::class, 'destroy'])->name('firewall.destroy');

    // File Manager (Phase 6)
    Route::get('/files', [FileManagerController::class, 'index'])->name('files.index');
    Route::get('/files/read', [FileManagerController::class, 'read'])->name('files.read');
    Route::post('/files/write', [FileManagerController::class, 'write'])->name('files.write');
    Route::post('/files/delete', [FileManagerController::class, 'delete'])->name('files.delete');

    // Services
    Route::get('/services', [ServiceController::class, 'index'])->name('services.index');
    Route::post('/services/{service}/{action}', [ServiceController::class, 'manage'])->name('services.manage');

    // Activity Logs
    Route::get('/activity', [ActivityLogController::class, 'index'])->name('activity.index');

    // Settings
    Route::get('/settings', [SettingController::class, 'index'])->name('settings.index');
    Route::post('/settings', [SettingController::class, 'update'])->name('settings.update');
});
