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
        Schema::table('websites', function (Blueprint $table) {
            if (!Schema::hasColumn('websites', 'deployment_source')) {
                $table->string('deployment_source', 32)->default('empty')->after('custom_nginx_config');
            }
            if (!Schema::hasColumn('websites', 'project_type')) {
                $table->string('project_type', 32)->default('laravel')->after('deployment_source');
            }
            if (!Schema::hasColumn('websites', 'git_repository')) {
                $table->string('git_repository')->nullable()->after('project_type');
            }
            if (!Schema::hasColumn('websites', 'git_branch')) {
                $table->string('git_branch', 100)->nullable()->default('main')->after('git_repository');
            }
            if (!Schema::hasColumn('websites', 'last_deployed_at')) {
                $table->timestamp('last_deployed_at')->nullable()->after('git_branch');
            }
        });
    }

    /**
     * Reverse the migrations.
     */
    public function down(): void
    {
        Schema::table('websites', function (Blueprint $table) {
            $columnsToDrop = [];
            foreach (['last_deployed_at', 'git_branch', 'git_repository', 'project_type', 'deployment_source'] as $col) {
                if (Schema::hasColumn('websites', $col)) {
                    $columnsToDrop[] = $col;
                }
            }
            if (!empty($columnsToDrop)) {
                $table->dropColumn($columnsToDrop);
            }
        });
    }
};
