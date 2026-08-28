<?php

namespace Tests\Feature;

use App\Models\CronJob;
use App\Models\Deployment;
use App\Models\FirewallRule;
use App\Models\Process;
use App\Models\User;
use App\Models\Website;
use Illuminate\Foundation\Testing\RefreshDatabase;
use Tests\TestCase;

class OperationsTest extends TestCase
{
    use RefreshDatabase;

    protected function setUp(): void
    {
        parent::setUp();
        $this->seed();
    }

    public function test_deployment_trigger_and_show(): void
    {
        $user = User::where('email', 'admin@kodepreneur.com')->first();

        $website = Website::create([
            'domain' => 'deploytest.com',
            'system_user' => 'kp_deploytest',
            'php_version' => '8.3',
            'document_root' => '/var/www/deploytest.com/public',
            'git_repository' => 'https://github.com/example/repo.git',
            'git_branch' => 'main',
            'status' => 'active',
        ]);

        $response = $this->actingAs($user)->post("/websites/{$website->id}/deploy");
        $response->assertRedirect();

        $deployment = Deployment::where('website_id', $website->id)->first();
        $this->assertNotNull($deployment);
        $this->assertEquals('success', $deployment->status);

        $showResponse = $this->actingAs($user)->get("/deployments/{$deployment->id}");
        $showResponse->assertStatus(200);
        $showResponse->assertInertia(fn ($page) => $page
            ->component('Deployments/Show')
            ->has('deployment')
        );
    }

    public function test_deployment_rollback(): void
    {
        $user = User::where('email', 'admin@kodepreneur.com')->first();

        $website = Website::create([
            'domain' => 'rollbacktest.com',
            'system_user' => 'kp_rollback',
            'php_version' => '8.3',
            'document_root' => '/var/www/rollbacktest.com/public',
            'status' => 'active',
        ]);

        $deployment = Deployment::create([
            'website_id' => $website->id,
            'commit_hash' => 'a1b2c3d4',
            'branch' => 'main',
            'status' => 'success',
            'trigger_source' => 'manual',
            'duration_seconds' => 2,
            'initiated_by_user_id' => $user->id,
        ]);

        $response = $this->actingAs($user)->post("/deployments/{$deployment->id}/rollback");
        $response->assertRedirect();

        $this->assertDatabaseHas('deployments', [
            'website_id' => $website->id,
            'trigger_source' => 'rollback',
            'commit_hash' => 'a1b2c3d4',
        ]);
    }

    public function test_cron_job_lifecycle(): void
    {
        $user = User::where('email', 'admin@kodepreneur.com')->first();

        $website = Website::create([
            'domain' => 'crontest.com',
            'system_user' => 'kp_cron',
            'php_version' => '8.3',
            'document_root' => '/var/www/crontest.com/public',
            'status' => 'active',
        ]);

        // 1. Create cron job
        $res = $this->actingAs($user)->post('/cron', [
            'website_id' => $website->id,
            'schedule' => '*/5 * * * *',
            'command' => 'php /var/www/crontest.com/artisan schedule:run',
        ]);
        $res->assertRedirect('/cron');

        $this->assertDatabaseHas('cron_jobs', [
            'schedule' => '*/5 * * * *',
            'is_active' => true,
        ]);

        $job = CronJob::first();

        // 2. Toggle active
        $this->actingAs($user)->post("/cron/{$job->id}/toggle");
        $this->assertDatabaseHas('cron_jobs', [
            'id' => $job->id,
            'is_active' => false,
        ]);

        // 3. Delete cron job
        $this->actingAs($user)->delete("/cron/{$job->id}");
        $this->assertDatabaseMissing('cron_jobs', [
            'id' => $job->id,
        ]);
    }

    public function test_managed_process_lifecycle(): void
    {
        $user = User::where('email', 'admin@kodepreneur.com')->first();

        // 1. Create process
        $res = $this->actingAs($user)->post('/processes', [
            'name' => 'laravel-queue',
            'command' => 'php artisan queue:work',
            'system_user' => 'kodepreneur',
            'instances' => 1,
        ]);
        $res->assertRedirect('/processes');

        $this->assertDatabaseHas('processes', [
            'name' => 'laravel-queue',
            'status' => 'running',
        ]);

        $proc = Process::first();

        // 2. Manage process action (stop, start, restart)
        $this->actingAs($user)->post("/processes/{$proc->id}/stop");
        $this->assertDatabaseHas('processes', [
            'id' => $proc->id,
            'status' => 'stopped',
        ]);

        // 3. Delete process
        $this->actingAs($user)->delete("/processes/{$proc->id}");
        $this->assertDatabaseMissing('processes', [
            'id' => $proc->id,
        ]);
    }

    public function test_firewall_rule_lifecycle(): void
    {
        $user = User::where('email', 'admin@kodepreneur.com')->first();

        // 1. Add rule
        $res = $this->actingAs($user)->post('/firewall', [
            'port' => 6379,
            'protocol' => 'tcp',
            'action' => 'allow',
            'source_ip' => '127.0.0.1',
            'description' => 'Redis localhost only',
        ]);
        $res->assertRedirect('/firewall');

        $this->assertDatabaseHas('firewall_rules', [
            'port' => 6379,
            'protocol' => 'tcp',
        ]);

        $rule = FirewallRule::first();

        // 2. Delete rule
        $this->actingAs($user)->delete("/firewall/{$rule->id}");
        $this->assertDatabaseMissing('firewall_rules', [
            'id' => $rule->id,
        ]);
    }

    public function test_file_manager_endpoints(): void
    {
        $user = User::where('email', 'admin@kodepreneur.com')->first();

        $website = Website::create([
            'domain' => 'filetest.com',
            'system_user' => 'kp_file',
            'php_version' => '8.3',
            'document_root' => '/var/www/filetest.com/public',
            'status' => 'active',
        ]);

        // 1. Browse files
        $browseRes = $this->actingAs($user)->get("/files?website_id={$website->id}");
        $browseRes->assertStatus(200);

        // 2. Read file
        $readRes = $this->actingAs($user)->get("/files/read?base_path=/var/www/filetest.com&relative_path=.env");
        $readRes->assertStatus(200);
        $readRes->assertJsonStructure(['success', 'content']);

        // 3. Write file
        $writeRes = $this->actingAs($user)->post('/files/write', [
            'base_path' => '/var/www/filetest.com',
            'relative_path' => '.env',
            'content' => 'APP_NAME=Testing',
            'website_id' => $website->id,
        ]);
        $writeRes->assertRedirect();

        // 4. Delete file
        $delRes = $this->actingAs($user)->post('/files/delete', [
            'base_path' => '/var/www/filetest.com',
            'relative_path' => 'scratch.log',
            'website_id' => $website->id,
        ]);
        $delRes->assertRedirect();
    }
}
