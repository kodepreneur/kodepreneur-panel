<?php

namespace App\Providers;

use Illuminate\Support\ServiceProvider;

class AppServiceProvider extends ServiceProvider
{
    /**
     * Register any application services.
     */
    public function register(): void
    {
        $this->app->singleton(\App\Services\Agent\AgentClientInterface::class, function () {
            if (config('services.agent.use_mock', env('KODEPRENEUR_AGENT_USE_MOCK', true))) {
                return new \App\Services\Agent\MockAgentClient();
            }

            return new \App\Services\Agent\HttpAgentClient(
                baseUrl: (string) config('services.agent.host', env('KODEPRENEUR_AGENT_HOST', 'http://127.0.0.1:8765')),
                secret: (string) config('services.agent.secret', env('KODEPRENEUR_AGENT_SECRET', 'kodepreneur-dev-secret-key-change-in-production'))
            );
        });
    }

    /**
     * Bootstrap any application services.
     */
    public function boot(): void
    {
        //
    }
}
