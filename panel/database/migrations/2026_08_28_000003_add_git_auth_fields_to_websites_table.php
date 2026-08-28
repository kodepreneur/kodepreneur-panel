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
            if (!Schema::hasColumn('websites', 'git_auth_type')) {
                $table->string('git_auth_type', 32)->default('none')->after('git_branch');
            }
            if (!Schema::hasColumn('websites', 'git_token')) {
                $table->text('git_token')->nullable()->after('git_auth_type');
            }
            if (!Schema::hasColumn('websites', 'git_token_user')) {
                $table->string('git_token_user', 100)->nullable()->after('git_token');
            }
            if (!Schema::hasColumn('websites', 'git_ssh_private_key')) {
                $table->text('git_ssh_private_key')->nullable()->after('git_token_user');
            }
            if (!Schema::hasColumn('websites', 'git_ssh_public_key')) {
                $table->text('git_ssh_public_key')->nullable()->after('git_ssh_private_key');
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
            foreach (['git_ssh_public_key', 'git_ssh_private_key', 'git_token_user', 'git_token', 'git_auth_type'] as $col) {
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
