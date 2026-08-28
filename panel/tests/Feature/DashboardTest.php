<?php

namespace Tests\Feature;

use App\Models\User;
use Illuminate\Foundation\Testing\RefreshDatabase;
use Tests\TestCase;

class DashboardTest extends TestCase
{
    use RefreshDatabase;

    protected function setUp(): void
    {
        parent::setUp();
        $this->seed();
    }

    public function test_guests_are_redirected_to_login(): void
    {
        $response = $this->get('/dashboard');
        $response->assertRedirect('/login');
    }

    public function test_authenticated_user_can_view_dashboard(): void
    {
        $user = User::where('email', 'admin@kodepreneur.com')->first();

        $response = $this->actingAs($user)->get('/dashboard');

        $response->assertStatus(200);
        $response->assertInertia(fn ($page) => $page
            ->component('Dashboard/Index')
            ->has('systemInfo')
            ->has('systemMetrics')
            ->has('services')
            ->has('counts')
        );
    }

    public function test_dashboard_metrics_endpoint_returns_json(): void
    {
        $user = User::where('email', 'admin@kodepreneur.com')->first();

        $response = $this->actingAs($user)->get('/dashboard/metrics');

        $response->assertStatus(200);
        $response->assertJsonStructure([
            'success',
            'data' => [
                'cpu',
                'memory',
                'disk',
                'load_average',
            ],
        ]);
    }
}
