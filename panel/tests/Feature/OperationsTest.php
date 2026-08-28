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

        // 1. Initial page render
        $browseRes = $this->actingAs($user)->get("/files?website_id={$website->id}");
        $browseRes->assertStatus(200);

        // 2. AJAX Browse
        $ajaxBrowse = $this->actingAs($user)->postJson('/files/browse', [
            'base_path' => '/var/www/filetest.com',
            'relative_path' => '',
            'show_hidden' => true,
        ]);
        $ajaxBrowse->assertStatus(200);
        $ajaxBrowse->assertJsonStructure(['success', 'files', 'disk_usage']);

        // 3. Read file
        $readRes = $this->actingAs($user)->get("/files/read?base_path=/var/www/filetest.com&relative_path=.env");
        $readRes->assertStatus(200);
        $readRes->assertJsonStructure(['success', 'content']);

        // 4. Write file
        $writeRes = $this->actingAs($user)->postJson('/files/write', [
            'base_path' => '/var/www/filetest.com',
            'relative_path' => '.env',
            'content' => 'APP_NAME=Testing',
        ]);
        $writeRes->assertStatus(200);

        // 5. Create File & Folder
        $createFileRes = $this->actingAs($user)->postJson('/files/create-file', [
            'base_path' => '/var/www/filetest.com',
            'relative_path' => 'new-file.txt',
        ]);
        $createFileRes->assertStatus(200);

        $createFolderRes = $this->actingAs($user)->postJson('/files/create-folder', [
            'base_path' => '/var/www/filetest.com',
            'relative_path' => 'new-folder',
        ]);
        $createFolderRes->assertStatus(200);

        // 6. Rename
        $renameRes = $this->actingAs($user)->postJson('/files/rename', [
            'base_path' => '/var/www/filetest.com',
            'old_path' => 'new-file.txt',
            'new_path' => 'renamed-file.txt',
        ]);
        $renameRes->assertStatus(200);

        // 7. Copy & Move
        $copyRes = $this->actingAs($user)->postJson('/files/copy', [
            'base_path' => '/var/www/filetest.com',
            'src_path' => 'renamed-file.txt',
            'dest_path' => 'new-folder',
        ]);
        $copyRes->assertStatus(200);

        $moveRes = $this->actingAs($user)->postJson('/files/move', [
            'base_path' => '/var/www/filetest.com',
            'src_path' => 'renamed-file.txt',
            'dest_path' => 'new-folder/moved.txt',
        ]);
        $moveRes->assertStatus(200);

        // 8. Chmod & Chown & Stat
        $chmodRes = $this->actingAs($user)->postJson('/files/chmod', [
            'base_path' => '/var/www/filetest.com',
            'relative_path' => 'new-folder',
            'mode' => '0755',
            'recursive' => false,
        ]);
        $chmodRes->assertStatus(200);

        $chownRes = $this->actingAs($user)->postJson('/files/chown', [
            'base_path' => '/var/www/filetest.com',
            'relative_path' => 'new-folder',
            'uid' => 33,
            'gid' => 33,
            'recursive' => false,
        ]);
        $chownRes->assertStatus(200);

        $statRes = $this->actingAs($user)->postJson('/files/stat', [
            'base_path' => '/var/www/filetest.com',
            'relative_path' => 'new-folder',
        ]);
        $statRes->assertStatus(200);
        $statRes->assertJsonStructure(['success', 'data']);

        // 9. Compress & Extract
        $compressRes = $this->actingAs($user)->postJson('/files/compress', [
            'base_path' => '/var/www/filetest.com',
            'sources' => ['new-folder'],
            'dest_path' => 'archive.zip',
            'format' => 'zip',
        ]);
        $compressRes->assertStatus(200);

        $extractRes = $this->actingAs($user)->postJson('/files/extract', [
            'base_path' => '/var/www/filetest.com',
            'archive_path' => 'archive.zip',
            'dest_path' => 'extracted',
        ]);
        $extractRes->assertStatus(200);

        // 10. Search
        $searchRes = $this->actingAs($user)->postJson('/files/search', [
            'base_path' => '/var/www/filetest.com',
            'query' => 'index',
        ]);
        $searchRes->assertStatus(200);
        $searchRes->assertJsonStructure(['success', 'data']);

        // 11. Upload file
        $fakeFile = \Illuminate\Http\UploadedFile::fake()->create('avatar.png', 10, 'image/png');
        $uploadRes = $this->actingAs($user)->post('/files/upload', [
            'base_path' => '/var/www/filetest.com',
            'relative_path' => 'public/images',
            'files' => [$fakeFile],
        ], ['Accept' => 'application/json']);
        $uploadRes->assertStatus(200);
        $uploadRes->assertJsonStructure(['success', 'uploaded']);

        // 12. Preview file
        $previewRes = $this->actingAs($user)->get('/files/preview?base_path=/var/www/filetest.com&relative_path=.env');
        $previewRes->assertStatus(200);

        // 13. Download file
        $downloadRes = $this->actingAs($user)->get('/files/download?base_path=/var/www/filetest.com&relative_path=.env');
        $downloadRes->assertStatus(200);

        // 14. Delete
        $delRes = $this->actingAs($user)->postJson('/files/delete', [
            'base_path' => '/var/www/filetest.com',
            'relative_path' => 'archive.zip',
        ]);
        $delRes->assertStatus(200);

        // 15. Verify Activity Logs were recorded
        $this->assertDatabaseHas('activity_logs', [
            'action' => 'file.write',
            'resource_type' => 'file',
        ]);
        $this->assertDatabaseHas('activity_logs', [
            'action' => 'file.create',
            'resource_type' => 'file',
        ]);
        $this->assertDatabaseHas('activity_logs', [
            'action' => 'file.upload',
            'resource_type' => 'file',
        ]);
        $this->assertDatabaseHas('activity_logs', [
            'action' => 'file.delete',
            'resource_type' => 'file',
        ]);
    }
}
