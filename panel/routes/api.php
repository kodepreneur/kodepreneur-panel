<?php

use App\Models\User;
use App\Services\Agent\AgentClientInterface;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Hash;
use Illuminate\Support\Facades\Route;

// Public Health & API Auth
Route::get('/health', function () {
    return response()->json([
        'status' => 'healthy',
        'panel' => 'Kodepreneur Control Plane',
        'version' => '1.0.0',
    ]);
});

Route::post('/v1/auth/login', function (Request $request) {
    $request->validate([
        'email' => 'required|email',
        'password' => 'required',
    ]);

    $user = User::where('email', $request->email)->first();

    if (!$user || !Hash::check($request->password, $user->password)) {
        return response()->json([
            'success' => false,
            'error' => [
                'code' => 'INVALID_CREDENTIALS',
                'message' => 'The provided credentials are incorrect.',
            ],
        ], 401);
    }

    $token = $user->createToken('api-token')->plainTextToken;

    return response()->json([
        'success' => true,
        'data' => [
            'token' => $token,
            'user' => [
                'id' => $user->id,
                'name' => $user->name,
                'email' => $user->email,
            ],
        ],
    ]);
});

// Protected Control Plane API
Route::middleware('auth:sanctum')->prefix('v1')->group(function () {
    Route::get('/auth/me', function (Request $request) {
        return response()->json([
            'success' => true,
            'data' => $request->user(),
        ]);
    });

    Route::get('/system/info', function (AgentClientInterface $agent) {
        return response()->json([
            'success' => true,
            'data' => $agent->getSystemInfo(),
        ]);
    });

    Route::get('/system/metrics', function (AgentClientInterface $agent) {
        return response()->json([
            'success' => true,
            'data' => $agent->getSystemMetrics(),
        ]);
    });
});
