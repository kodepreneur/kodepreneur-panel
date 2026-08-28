<?php

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

return new class extends Migration
{
    /**
     * Run the migrations.
     */
    public function up(): void
    {
        // 1. Websites
        Schema::create('websites', function (Blueprint $table) {
            $table->id();
            $table->string('domain')->unique();
            $table->json('aliases')->nullable();
            $table->string('php_version', 10)->default('8.3');
            $table->string('document_root');
            $table->string('system_user', 64);
            $table->boolean('ssl_enabled')->default(false);
            $table->boolean('force_https')->default(false);
            $table->string('status', 32)->default('active'); // active, suspended, provisioning, error
            $table->text('custom_nginx_config')->nullable();
            $table->string('deployment_source', 32)->default('empty'); // empty, zip, git
            $table->string('project_type', 32)->default('laravel'); // laravel, generic_php, static, auto
            $table->string('git_repository')->nullable();
            $table->string('git_branch', 100)->nullable()->default('main');
            $table->timestamp('last_deployed_at')->nullable();
            $table->timestamps();
        });

        // 2. Domains
        Schema::create('domains', function (Blueprint $table) {
            $table->id();
            $table->foreignId('website_id')->constrained('websites')->cascadeOnDelete();
            $table->string('domain')->unique();
            $table->boolean('is_primary')->default(false);
            $table->string('ssl_status', 32)->default('pending');
            $table->timestamps();
        });

        // 3. Databases
        Schema::create('databases', function (Blueprint $table) {
            $table->id();
            $table->string('engine', 32)->default('mysql'); // mysql, postgresql
            $table->string('name', 64)->unique();
            $table->string('character_set', 32)->default('utf8mb4');
            $table->string('collation', 64)->default('utf8mb4_unicode_ci');
            $table->unsignedBigInteger('size_bytes')->nullable();
            $table->timestamps();
        });

        // 4. Database Users
        Schema::create('database_users', function (Blueprint $table) {
            $table->id();
            $table->string('engine', 32)->default('mysql');
            $table->string('username', 64);
            $table->string('host', 64)->default('localhost');
            $table->timestamps();
            $table->unique(['engine', 'username', 'host']);
        });

        // 5. Database Access / Permissions
        Schema::create('database_access', function (Blueprint $table) {
            $table->id();
            $table->foreignId('database_id')->constrained('databases')->cascadeOnDelete();
            $table->foreignId('database_user_id')->constrained('database_users')->cascadeOnDelete();
            $table->string('permissions', 32)->default('all');
            $table->timestamps();
            $table->unique(['database_id', 'database_user_id']);
        });

        // 6. PHP Versions
        Schema::create('php_versions', function (Blueprint $table) {
            $table->id();
            $table->string('version', 10)->unique();
            $table->string('status', 32)->default('installed');
            $table->boolean('is_default')->default(false);
            $table->string('binary_path')->nullable();
            $table->string('fpm_socket_path')->nullable();
            $table->timestamps();
        });

        // 7. Deployments
        Schema::create('deployments', function (Blueprint $table) {
            $table->id();
            $table->foreignId('website_id')->constrained('websites')->cascadeOnDelete();
            $table->string('commit_hash', 40)->nullable();
            $table->text('commit_message')->nullable();
            $table->string('branch', 100)->default('main');
            $table->string('status', 32)->default('pending'); // pending, running, success, failed, cancelled
            $table->string('trigger_source', 32)->default('manual'); // manual, webhook, rollback
            $table->mediumText('log_output')->nullable();
            $table->unsignedInteger('duration_seconds')->nullable();
            $table->foreignId('initiated_by_user_id')->nullable()->constrained('users')->nullOnDelete();
            $table->timestamps();
        });

        // 8. Managed Processes
        Schema::create('processes', function (Blueprint $table) {
            $table->id();
            $table->string('name', 100);
            $table->foreignId('website_id')->nullable()->constrained('websites')->nullOnDelete();
            $table->text('command');
            $table->string('system_user', 64)->default('kodepreneur');
            $table->unsignedSmallInteger('instances')->default(1);
            $table->string('status', 32)->default('stopped');
            $table->string('unit_file_path')->nullable();
            $table->timestamps();
        });

        // 9. Cron Jobs
        Schema::create('cron_jobs', function (Blueprint $table) {
            $table->id();
            $table->foreignId('website_id')->nullable()->constrained('websites')->nullOnDelete();
            $table->string('schedule', 100); // e.g. * * * * *
            $table->text('command');
            $table->string('system_user', 64)->default('kodepreneur');
            $table->boolean('is_active')->default(true);
            $table->timestamp('last_run_at')->nullable();
            $table->timestamp('next_run_at')->nullable();
            $table->timestamps();
        });

        // 10. SSL Certificates
        Schema::create('ssl_certificates', function (Blueprint $table) {
            $table->id();
            $table->foreignId('website_id')->constrained('websites')->cascadeOnDelete();
            $table->string('domain');
            $table->string('issuer', 100)->default('letsencrypt');
            $table->string('cert_path')->nullable();
            $table->string('key_path')->nullable();
            $table->timestamp('valid_from')->nullable();
            $table->timestamp('valid_until')->nullable();
            $table->boolean('auto_renew')->default(true);
            $table->string('status', 32)->default('valid');
            $table->timestamps();
        });

        // 11. Firewall Rules
        Schema::create('firewall_rules', function (Blueprint $table) {
            $table->id();
            $table->unsignedInteger('port');
            $table->string('protocol', 10)->default('tcp'); // tcp, udp, both
            $table->string('action', 10)->default('allow'); // allow, deny, reject
            $table->string('source_ip', 45)->nullable()->default('0.0.0.0/0');
            $table->string('description', 255)->nullable();
            $table->timestamps();
        });

        // 12. Activity Logs
        Schema::create('activity_logs', function (Blueprint $table) {
            $table->id();
            $table->foreignId('user_id')->nullable()->constrained('users')->nullOnDelete();
            $table->string('user_email', 255)->nullable();
            $table->string('ip_address', 45);
            $table->text('user_agent')->nullable();
            $table->string('action', 100);
            $table->string('resource_type', 50)->nullable();
            $table->string('resource_id', 64)->nullable();
            $table->string('status', 20)->default('success');
            $table->json('payload_summary')->nullable();
            $table->timestamp('created_at')->useCurrent();
        });

        // 13. Server Settings
        Schema::create('server_settings', function (Blueprint $table) {
            $table->id();
            $table->string('key', 100)->unique();
            $table->text('value')->nullable();
            $table->string('type', 32)->default('string'); // string, integer, boolean, json
            $table->string('description', 255)->nullable();
            $table->timestamps();
        });
    }

    /**
     * Reverse the migrations.
     */
    public function down(): void
    {
        Schema::dropIfExists('server_settings');
        Schema::dropIfExists('activity_logs');
        Schema::dropIfExists('firewall_rules');
        Schema::dropIfExists('ssl_certificates');
        Schema::dropIfExists('cron_jobs');
        Schema::dropIfExists('processes');
        Schema::dropIfExists('deployments');
        Schema::dropIfExists('php_versions');
        Schema::dropIfExists('database_access');
        Schema::dropIfExists('database_users');
        Schema::dropIfExists('databases');
        Schema::dropIfExists('domains');
        Schema::dropIfExists('websites');
    }
};
