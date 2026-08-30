<?php

namespace Tests\Feature;

use App\Models\Deployment;
use App\Models\User;
use App\Models\Website;
use Illuminate\Foundation\Testing\RefreshDatabase;
use Illuminate\Http\UploadedFile;
use Tests\TestCase;

class WebsiteTest extends TestCase
{
    use RefreshDatabase;

    protected function setUp(): void
    {
        parent::setUp();
        $this->seed();
    }

    public function test_authenticated_user_can_list_websites(): void
    {
        $user = User::where('email', 'admin@kodepreneur.com')->first();

        Website::create([
            'domain' => 'testsite.com',
            'php_version' => '8.3',
            'document_root' => '/var/www/testsite.com/public',
            'system_user' => 'kp_testsite',
            'status' => 'active',
        ]);

        $response = $this->actingAs($user)->get('/websites');

        $response->assertStatus(200);
        $response->assertInertia(fn ($page) => $page
            ->component('Websites/Index')
            ->has('websites.data', 1)
        );
    }

    public function test_user_can_create_a_website(): void
    {
        $user = User::where('email', 'admin@kodepreneur.com')->first();

        $response = $this->actingAs($user)->post('/websites', [
            'domain' => 'app.example.com',
            'php_version' => '8.4',
            'deployment_source' => 'empty',
            'project_type' => 'laravel',
            'document_root' => '/var/www/app.example.com/public',
        ]);

        $website = Website::where('domain', 'app.example.com')->firstOrFail();
        $response->assertRedirect("/websites/{$website->id}");

        $this->assertDatabaseHas('websites', [
            'domain' => 'app.example.com',
            'php_version' => '8.4',
            'deployment_source' => 'empty',
            'document_root' => '/var/www/app.example.com/public',
        ]);
        $this->assertDatabaseHas('domains', [
            'domain' => 'app.example.com',
            'is_primary' => true,
        ]);
        $this->assertDatabaseHas('activity_logs', [
            'action' => 'website.create',
        ]);
    }

    public function test_user_can_create_website_with_git_repo(): void
    {
        $user = User::where('email', 'admin@kodepreneur.com')->first();

        $response = $this->actingAs($user)->post('/websites', [
            'domain' => 'git-app.com',
            'php_version' => '8.3',
            'deployment_source' => 'git',
            'project_type' => 'laravel',
            'git_repository' => 'https://github.com/example/laravel-starter.git',
            'git_branch' => 'main',
        ]);

        $website = Website::where('domain', 'git-app.com')->firstOrFail();
        $response->assertRedirect("/websites/{$website->id}");

        $this->assertDatabaseHas('websites', [
            'domain' => 'git-app.com',
            'deployment_source' => 'git',
            'git_repository' => 'https://github.com/example/laravel-starter.git',
            'git_branch' => 'main',
            'document_root' => '/var/www/git-app.com/public',
        ]);

        $this->assertDatabaseHas('deployments', [
            'website_id' => $website->id,
            'branch' => 'main',
            'status' => 'success',
        ]);
    }

    public function test_user_can_generate_deploy_key(): void
    {
        $user = User::where('email', 'admin@kodepreneur.com')->first();

        $response = $this->actingAs($user)->postJson('/websites/deploy-key/generate');

        $response->assertStatus(200);
        $response->assertJsonStructure([
            'success',
            'public_key',
            'private_key',
        ]);
        $this->assertTrue($response->json('success'));
        $this->assertStringContainsString('ssh-', $response->json('public_key'));
    }

    public function test_user_can_create_website_with_private_git_repo_ssh_key(): void
    {
        $user = User::where('email', 'admin@kodepreneur.com')->first();

        $response = $this->actingAs($user)->post('/websites', [
            'domain' => 'private-ssh.com',
            'php_version' => '8.4',
            'deployment_source' => 'git',
            'project_type' => 'laravel',
            'git_repository' => 'git@github.com:myorg/private-repo.git',
            'git_branch' => 'production',
            'git_auth_type' => 'ssh_key',
            'git_ssh_public_key' => 'ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIG... kodepreneur-deploy-key',
            'git_ssh_private_key' => "-----BEGIN OPENSSH PRIVATE KEY-----\ntest\n-----END OPENSSH PRIVATE KEY-----",
        ]);

        $website = Website::where('domain', 'private-ssh.com')->firstOrFail();
        $response->assertRedirect("/websites/{$website->id}");

        $this->assertEquals('ssh_key', $website->git_auth_type);
        $this->assertEquals('git@github.com:myorg/private-repo.git', $website->git_repository);
        $this->assertEquals('production', $website->git_branch);
        $this->assertNotNull($website->git_ssh_private_key);
        $this->assertStringContainsString('ssh-ed25519', $website->git_ssh_public_key);
    }

    public function test_user_can_create_website_with_private_git_repo_token(): void
    {
        $user = User::where('email', 'admin@kodepreneur.com')->first();

        $response = $this->actingAs($user)->post('/websites', [
            'domain' => 'private-token.com',
            'php_version' => '8.3',
            'deployment_source' => 'git',
            'project_type' => 'laravel',
            'git_repository' => 'https://github.com/myorg/token-repo.git',
            'git_branch' => 'main',
            'git_auth_type' => 'token',
            'git_token' => 'ghp_secret_access_token_123',
            'git_token_user' => 'x-access-token',
        ]);

        $website = Website::where('domain', 'private-token.com')->firstOrFail();
        $response->assertRedirect("/websites/{$website->id}");

        $this->assertEquals('token', $website->git_auth_type);
        $this->assertEquals('https://github.com/myorg/token-repo.git', $website->git_repository);
        $this->assertEquals('ghp_secret_access_token_123', $website->git_token);
        $this->assertEquals('x-access-token', $website->git_token_user);
    }

    public function test_user_can_create_website_with_zip_file(): void
    {
        $user = User::where('email', 'admin@kodepreneur.com')->first();

        // Create a fake zip file
        $fakeZipContent = "PK\x03\x04fake zip content";
        $file = UploadedFile::fake()->createWithContent('project.zip', $fakeZipContent);

        $response = $this->actingAs($user)->post('/websites', [
            'domain' => 'zip-app.com',
            'php_version' => '8.3',
            'deployment_source' => 'zip',
            'project_type' => 'laravel',
            'zip_file' => $file,
        ]);

        $website = Website::where('domain', 'zip-app.com')->firstOrFail();
        $response->assertRedirect("/websites/{$website->id}");

        $this->assertDatabaseHas('websites', [
            'domain' => 'zip-app.com',
            'deployment_source' => 'zip',
            'document_root' => '/var/www/zip-app.com/public',
        ]);
    }

    public function test_laravel_project_auto_routes_to_public(): void
    {
        $user = User::where('email', 'admin@kodepreneur.com')->first();

        // When document_root is omitted or set to root /var/www/laravel-site.com
        $response = $this->actingAs($user)->post('/websites', [
            'domain' => 'laravel-site.com',
            'php_version' => '8.3',
            'deployment_source' => 'empty',
            'project_type' => 'laravel',
            'document_root' => '/var/www/laravel-site.com',
        ]);

        $website = Website::where('domain', 'laravel-site.com')->firstOrFail();
        // Should automatically direct to /public
        $this->assertEquals('/var/www/laravel-site.com/public', $website->document_root);
    }

    public function test_user_cannot_create_duplicate_website(): void
    {
        $user = User::where('email', 'admin@kodepreneur.com')->first();

        Website::create([
            'domain' => 'duplicate.com',
            'php_version' => '8.3',
            'document_root' => '/var/www/duplicate.com/public',
            'system_user' => 'kp_duplicate',
        ]);

        $response = $this->actingAs($user)->post('/websites', [
            'domain' => 'duplicate.com',
            'php_version' => '8.3',
        ]);

        $response->assertSessionHasErrors('domain');
    }

    public function test_user_can_switch_php_version(): void
    {
        $user = User::where('email', 'admin@kodepreneur.com')->first();

        $website = Website::create([
            'domain' => 'php-switch.com',
            'php_version' => '8.3',
            'document_root' => '/var/www/php-switch.com/public',
            'system_user' => 'kp_phpswitch',
        ]);

        $response = $this->actingAs($user)->put("/websites/{$website->id}/php", [
            'php_version' => '8.4',
        ]);

        $response->assertSessionHasNoErrors();
        $this->assertEquals('8.4', $website->fresh()->php_version);
        $this->assertDatabaseHas('activity_logs', [
            'action' => 'website.switch_php',
        ]);
    }

    public function test_user_can_issue_ssl_certificate(): void
    {
        $user = User::where('email', 'admin@kodepreneur.com')->first();

        $website = Website::create([
            'domain' => 'ssl-site.com',
            'php_version' => '8.3',
            'document_root' => '/var/www/ssl-site.com/public',
            'system_user' => 'kp_sslsite',
            'ssl_enabled' => false,
        ]);

        $response = $this->actingAs($user)->post("/websites/{$website->id}/ssl", [
            'email' => 'admin@ssl-site.com',
            'force_https' => true,
        ]);

        $response->assertSessionHasNoErrors();
        $this->assertTrue($website->fresh()->ssl_enabled);
        $this->assertDatabaseHas('ssl_certificates', [
            'website_id' => $website->id,
            'domain' => 'ssl-site.com',
        ]);
    }

    public function test_user_can_fetch_website_logs(): void
    {
        $user = User::where('email', 'admin@kodepreneur.com')->first();

        $website = Website::create([
            'domain' => 'logs-site.com',
            'php_version' => '8.3',
            'document_root' => '/var/www/logs-site.com/public',
            'system_user' => 'kp_logssite',
        ]);

        $response = $this->actingAs($user)->getJson("/websites/{$website->id}/logs/access?lines=20");

        $response->assertStatus(200);
        $response->assertJsonStructure([
            'success',
            'domain',
            'type',
            'lines',
        ]);
    }

    public function test_user_can_fetch_website_traffic(): void
    {
        $user = User::where('email', 'admin@kodepreneur.com')->first();

        $website = Website::create([
            'domain' => 'traffic-site.com',
            'php_version' => '8.3',
            'document_root' => '/var/www/traffic-site.com/public',
            'system_user' => 'kp_trafficsite',
        ]);

        $response = $this->actingAs($user)->getJson("/websites/{$website->id}/traffic?period=24h");

        $response->assertStatus(200);
        $response->assertJsonStructure([
            'success',
            'data' => [
                'domain',
                'period',
                'total_requests',
                'total_bytes_sent',
                'unique_visitors',
                'success_rate',
                'status_codes',
                'status_categories',
                'time_series',
                'top_paths',
                'top_ips',
                'top_referrers',
                'top_user_agents',
                'recent_requests',
            ],
        ]);
        $this->assertEquals('traffic-site.com', $response->json('data.domain'));
    }

    public function test_user_can_delete_a_website(): void
    {
        $user = User::where('email', 'admin@kodepreneur.com')->first();

        $website = Website::create([
            'domain' => 'delete-me.com',
            'php_version' => '8.3',
            'document_root' => '/var/www/delete-me.com/public',
            'system_user' => 'kp_deleteme',
        ]);

        $response = $this->actingAs($user)->delete("/websites/{$website->id}");

        $response->assertRedirect('/websites');
        $this->assertDatabaseMissing('websites', [
            'id' => $website->id,
        ]);
        $this->assertDatabaseHas('activity_logs', [
            'action' => 'website.delete',
        ]);
    }
}
