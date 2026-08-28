<?php

namespace Tests\Feature;

use App\Models\ActivityLog;
use App\Models\User;
use Illuminate\Foundation\Testing\RefreshDatabase;
use Tests\TestCase;

class UpdateTest extends TestCase
{
    use RefreshDatabase;

    protected function setUp(): void
    {
        parent::setUp();
        $this->seed();
    }

    public function test_updates_screen_can_be_rendered(): void
    {
        $user = User::where('email', 'admin@kodepreneur.com')->first();

        $response = $this->actingAs($user)->get('/updates');
        $response->assertStatus(200);
        $response->assertInertia(fn ($page) => $page
            ->component('Updates/Index')
            ->has('updateInfo')
            ->where('updateInfo.current_version', 'v1.0.0')
        );
    }

    public function test_updates_check_endpoint_returns_json(): void
    {
        $user = User::where('email', 'admin@kodepreneur.com')->first();

        $response = $this->actingAs($user)->postJson('/updates/check');
        $response->assertStatus(200)
            ->assertJsonStructure([
                'success',
                'updateInfo' => [
                    'current_version',
                    'current_commit',
                    'repository',
                    'latest_commit',
                    'has_update',
                    'agent_status',
                ],
            ]);
    }

    public function test_updates_execute_triggers_system_update(): void
    {
        $user = User::where('email', 'admin@kodepreneur.com')->first();

        $response = $this->actingAs($user)->postJson('/updates/execute', [
            'repository' => 'https://github.com/kodepreneur/kodepreneur-panel.git',
            'branch' => 'main',
        ]);

        $response->assertStatus(200)
            ->assertJson([
                'success' => true,
                'data' => [
                    'success' => true,
                    'commit_hash' => '8706c6c49832',
                ],
            ]);

        $this->assertDatabaseHas('activity_logs', [
            'action' => 'system.update',
            'status' => 'success',
        ]);
    }
}
