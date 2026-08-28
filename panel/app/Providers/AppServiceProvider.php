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
            $useMock = config('services.agent.use_mock', false);

            if ($useMock) {
                return new \App\Services\Agent\MockAgentClient();
            }

            return new \App\Services\Agent\HttpAgentClient(
                baseUrl: (string) config('services.agent.host', 'http://127.0.0.1:8443'),
                secret: (string) config('services.agent.secret', '')
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
