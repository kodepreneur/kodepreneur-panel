<?php

namespace Tests\Feature;

use App\Models\Database;
use App\Models\DatabaseAccess;
use App\Models\DatabaseUser;
use App\Models\User;
use Illuminate\Foundation\Testing\RefreshDatabase;
use Tests\TestCase;

class DatabaseTest extends TestCase
{
    use RefreshDatabase;

    protected function setUp(): void
    {
        parent::setUp();
        $this->seed();
    }

    public function test_authenticated_user_can_list_databases(): void
    {
        $user = User::where('email', 'admin@kodepreneur.com')->first();

        Database::create([
            'engine' => 'mysql',
            'name' => 'test_db',
            'character_set' => 'utf8mb4',
            'collation' => 'utf8mb4_unicode_ci',
        ]);

        $response = $this->actingAs($user)->get('/databases');

        $response->assertStatus(200);
        $response->assertInertia(fn ($page) => $page
            ->component('Databases/Index')
            ->has('databases.data', 1)
        );
    }

    public function test_user_can_create_a_database(): void
    {
        $user = User::where('email', 'admin@kodepreneur.com')->first();

        $response = $this->actingAs($user)->post('/databases', [
            'engine' => 'mysql',
            'name' => 'shop_production',
            'character_set' => 'utf8mb4',
            'collation' => 'utf8mb4_unicode_ci',
        ]);

        $response->assertRedirect('/databases');
        $this->assertDatabaseHas('databases', [
            'name' => 'shop_production',
            'engine' => 'mysql',
        ]);
        $this->assertDatabaseHas('activity_logs', [
            'action' => 'database.create',
        ]);
    }

    public function test_user_can_create_a_database_user(): void
    {
        $user = User::where('email', 'admin@kodepreneur.com')->first();

        $response = $this->actingAs($user)->post('/databases/users', [
            'engine' => 'mysql',
            'username' => 'shop_user',
            'host' => 'localhost',
            'password' => 'SuperSecretPass123!',
        ]);

        $response->assertRedirect('/databases');
        $this->assertDatabaseHas('database_users', [
            'username' => 'shop_user',
            'engine' => 'mysql',
            'host' => 'localhost',
        ]);
        $this->assertDatabaseHas('activity_logs', [
            'action' => 'database_user.create',
        ]);
    }

    public function test_user_can_grant_database_access(): void
    {
        $user = User::where('email', 'admin@kodepreneur.com')->first();

        $db = Database::create([
            'engine' => 'mysql',
            'name' => 'crm_db',
            'character_set' => 'utf8mb4',
            'collation' => 'utf8mb4_unicode_ci',
        ]);

        $dbUser = DatabaseUser::create([
            'engine' => 'mysql',
            'username' => 'crm_user',
            'host' => 'localhost',
        ]);

        $response = $this->actingAs($user)->post('/databases/grants', [
            'database_id' => $db->id,
            'database_user_id' => $dbUser->id,
            'permissions' => 'all',
        ]);

        $response->assertRedirect('/databases');
        $this->assertDatabaseHas('database_access', [
            'database_id' => $db->id,
            'database_user_id' => $dbUser->id,
            'permissions' => 'all',
        ]);
        $this->assertDatabaseHas('activity_logs', [
            'action' => 'database.grant_access',
        ]);
    }

    public function test_user_can_reset_database_user_password(): void
    {
        $user = User::where('email', 'admin@kodepreneur.com')->first();

        $dbUser = DatabaseUser::create([
            'engine' => 'mysql',
            'username' => 'api_user',
            'host' => 'localhost',
        ]);

        $response = $this->actingAs($user)->post("/databases/users/{$dbUser->id}/password", [
            'password' => 'NewSecurePassword123!',
        ]);

        $response->assertRedirect('/databases');
        $this->assertDatabaseHas('activity_logs', [
            'action' => 'database_user.reset_password',
        ]);
    }

    public function test_user_can_delete_a_database_user(): void
    {
        $user = User::where('email', 'admin@kodepreneur.com')->first();

        $dbUser = DatabaseUser::create([
            'engine' => 'mysql',
            'username' => 'delete_user',
            'host' => 'localhost',
        ]);

        $response = $this->actingAs($user)->delete("/databases/users/{$dbUser->id}");

        $response->assertRedirect('/databases');
        $this->assertDatabaseMissing('database_users', [
            'id' => $dbUser->id,
        ]);
        $this->assertDatabaseHas('activity_logs', [
            'action' => 'database_user.delete',
        ]);
    }

    public function test_user_can_delete_a_database(): void
    {
        $user = User::where('email', 'admin@kodepreneur.com')->first();

        $db = Database::create([
            'engine' => 'mysql',
            'name' => 'drop_me_db',
            'character_set' => 'utf8mb4',
            'collation' => 'utf8mb4_unicode_ci',
        ]);

        $response = $this->actingAs($user)->delete("/databases/{$db->id}");

        $response->assertRedirect('/databases');
        $this->assertDatabaseMissing('databases', [
            'id' => $db->id,
        ]);
        $this->assertDatabaseHas('activity_logs', [
            'action' => 'database.delete',
        ]);
    }

    public function test_user_can_view_database_explorer(): void
    {
        $user = User::where('email', 'admin@kodepreneur.com')->first();

        $db = Database::create([
            'engine' => 'mysql',
            'name' => 'shop_analytics',
            'character_set' => 'utf8mb4',
            'collation' => 'utf8mb4_unicode_ci',
        ]);

        $response = $this->actingAs($user)->get("/databases/{$db->id}/explorer");

        $response->assertStatus(200);
        $response->assertInertia(fn ($page) => $page
            ->component('Databases/Explorer')
            ->has('currentDatabase')
            ->has('allDatabases')
            ->has('initialTables')
        );
    }

    public function test_user_can_fetch_database_tables_json(): void
    {
        $user = User::where('email', 'admin@kodepreneur.com')->first();

        $db = Database::create([
            'engine' => 'mysql',
            'name' => 'shop_analytics',
            'character_set' => 'utf8mb4',
            'collation' => 'utf8mb4_unicode_ci',
        ]);

        $response = $this->actingAs($user)->getJson("/databases/{$db->id}/tables");

        $response->assertStatus(200);
        $response->assertJsonStructure([
            'success',
            'tables' => [
                '*' => ['name', 'type', 'engine', 'rows', 'total_size'],
            ],
        ]);
    }

    public function test_user_can_fetch_table_structure_json(): void
    {
        $user = User::where('email', 'admin@kodepreneur.com')->first();

        $db = Database::create([
            'engine' => 'mysql',
            'name' => 'shop_analytics',
            'character_set' => 'utf8mb4',
            'collation' => 'utf8mb4_unicode_ci',
        ]);

        $response = $this->actingAs($user)->getJson("/databases/{$db->id}/tables/orders/structure");

        $response->assertStatus(200);
        $response->assertJsonStructure([
            'success',
            'structure' => [
                'table_name',
                'columns' => [
                    '*' => ['name', 'type', 'is_nullable', 'key'],
                ],
                'indexes',
                'foreign_keys',
                'create_statement',
            ],
        ]);
    }

    public function test_user_can_fetch_table_data_json(): void
    {
        $user = User::where('email', 'admin@kodepreneur.com')->first();

        $db = Database::create([
            'engine' => 'mysql',
            'name' => 'shop_analytics',
            'character_set' => 'utf8mb4',
            'collation' => 'utf8mb4_unicode_ci',
        ]);

        $response = $this->actingAs($user)->getJson("/databases/{$db->id}/tables/orders/data?page=1&per_page=10");

        $response->assertStatus(200);
        $response->assertJsonStructure([
            'success',
            'data' => [
                'table_name',
                'columns',
                'rows',
                'total_rows',
                'page',
                'per_page',
                'total_pages',
            ],
        ]);
    }
}
