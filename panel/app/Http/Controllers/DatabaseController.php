<?php

namespace App\Http\Controllers;

use App\Models\ActivityLog;
use App\Models\Database;
use App\Models\DatabaseAccess;
use App\Models\DatabaseUser;
use App\Services\Agent\AgentClientInterface;
use Exception;
use Illuminate\Http\RedirectResponse;
use Illuminate\Http\Request;
use Inertia\Inertia;
use Inertia\Response;

class DatabaseController extends Controller
{
    public function __construct(
        protected AgentClientInterface $agentClient
    ) {}

    public function index(): Response
    {
        $databases = Database::with('users')
            ->latest()
            ->paginate(15);

        $databaseUsers = DatabaseUser::with('databases')
            ->latest()
            ->get();

        $allDatabases = Database::select('id', 'name', 'engine')->orderBy('name')->get();

        return Inertia::render('Databases/Index', [
            'databases' => $databases,
            'databaseUsers' => $databaseUsers,
            'allDatabases' => $allDatabases,
        ]);
    }

    public function store(Request $request): RedirectResponse
    {
        $validated = $request->validate([
            'engine' => ['required', 'string', 'in:mysql,postgresql'],
            'name' => ['required', 'string', 'max:64', 'unique:databases,name', 'regex:/^[a-zA-Z0-9_]+$/'],
            'character_set' => ['nullable', 'string', 'max:32'],
            'collation' => ['nullable', 'string', 'max:64'],
        ]);

        $name = strtolower($validated['name']);
        $engine = $validated['engine'];
        $charset = $validated['character_set'] ?: ($engine === 'mysql' ? 'utf8mb4' : 'UTF8');
        $collation = $validated['collation'] ?: ($engine === 'mysql' ? 'utf8mb4_unicode_ci' : 'en_US.UTF-8');

        try {
            // Provision on server via Agent
            $this->agentClient->createDatabase([
                'engine' => $engine,
                'name' => $name,
                'charset' => $charset,
                'collation' => $collation,
            ]);

            $db = Database::create([
                'engine' => $engine,
                'name' => $name,
                'character_set' => $charset,
                'collation' => $collation,
            ]);

            ActivityLog::create([
                'user_id' => $request->user()->id,
                'user_email' => $request->user()->email,
                'ip_address' => $request->ip() ?: '127.0.0.1',
                'user_agent' => $request->userAgent(),
                'action' => 'database.create',
                'resource_type' => 'database',
                'resource_id' => (string) $db->id,
                'status' => 'success',
                'payload_summary' => [
                    'engine' => $db->engine,
                    'name' => $db->name,
                ],
            ]);

            return redirect()->route('databases.index')->with('success', "Database {$db->name} created successfully.");
        } catch (Exception $e) {
            return back()->withInput()->with('error', "Failed to create database: " . $e->getMessage());
        }
    }

    public function storeUser(Request $request): RedirectResponse
    {
        $validated = $request->validate([
            'engine' => ['required', 'string', 'in:mysql,postgresql'],
            'username' => ['required', 'string', 'max:64', 'regex:/^[a-zA-Z0-9_]+$/'],
            'host' => ['nullable', 'string', 'max:64'],
            'password' => ['required', 'string', 'min:8'],
        ]);

        $username = strtolower($validated['username']);
        $engine = $validated['engine'];
        $host = $validated['host'] ?: 'localhost';

        // Check uniqueness for (engine, username, host)
        $exists = DatabaseUser::where('engine', $engine)
            ->where('username', $username)
            ->where('host', $host)
            ->exists();

        if ($exists) {
            return back()->withInput()->withErrors(['username' => "User {$username}@{$host} already exists for {$engine}."]);
        }

        try {
            // Provision user on server via Agent
            $this->agentClient->createDatabaseUser([
                'engine' => $engine,
                'username' => $username,
                'host' => $host,
                'password' => $validated['password'],
            ]);

            $dbUser = DatabaseUser::create([
                'engine' => $engine,
                'username' => $username,
                'host' => $host,
            ]);

            ActivityLog::create([
                'user_id' => $request->user()->id,
                'user_email' => $request->user()->email,
                'ip_address' => $request->ip() ?: '127.0.0.1',
                'user_agent' => $request->userAgent(),
                'action' => 'database_user.create',
                'resource_type' => 'database_user',
                'resource_id' => (string) $dbUser->id,
                'status' => 'success',
                'payload_summary' => [
                    'engine' => $dbUser->engine,
                    'username' => $dbUser->username,
                    'host' => $dbUser->host,
                ],
            ]);

            return redirect()->route('databases.index')->with('success', "Database user {$username}@{$host} created successfully.");
        } catch (Exception $e) {
            return back()->withInput()->with('error', "Failed to create database user: " . $e->getMessage());
        }
    }

    public function storeGrant(Request $request): RedirectResponse
    {
        $validated = $request->validate([
            'database_id' => ['required', 'exists:databases,id'],
            'database_user_id' => ['required', 'exists:database_users,id'],
            'permissions' => ['nullable', 'string', 'max:32'],
        ]);

        $db = Database::findOrFail($validated['database_id']);
        $dbUser = DatabaseUser::findOrFail($validated['database_user_id']);
        $permissions = $validated['permissions'] ?: 'all';

        if ($db->engine !== $dbUser->engine) {
            return back()->with('error', "Cannot grant access: Database engine ({$db->engine}) and user engine ({$dbUser->engine}) must match.");
        }

        try {
            $this->agentClient->grantDatabaseAccess([
                'engine' => $db->engine,
                'database' => $db->name,
                'username' => $dbUser->username,
                'host' => $dbUser->host,
                'permissions' => $permissions,
            ]);

            DatabaseAccess::updateOrCreate(
                [
                    'database_id' => $db->id,
                    'database_user_id' => $dbUser->id,
                ],
                [
                    'permissions' => $permissions,
                ]
            );

            ActivityLog::create([
                'user_id' => $request->user()->id,
                'user_email' => $request->user()->email,
                'ip_address' => $request->ip() ?: '127.0.0.1',
                'user_agent' => $request->userAgent(),
                'action' => 'database.grant_access',
                'resource_type' => 'database_access',
                'resource_id' => "{$db->id}-{$dbUser->id}",
                'status' => 'success',
                'payload_summary' => [
                    'database' => $db->name,
                    'user' => $dbUser->username,
                    'permissions' => $permissions,
                ],
            ]);

            return redirect()->route('databases.index')->with('success', "Granted access on {$db->name} to {$dbUser->username}.");
        } catch (Exception $e) {
            return back()->with('error', "Failed to grant database access: " . $e->getMessage());
        }
    }

    public function resetPassword(Request $request, DatabaseUser $databaseUser): RedirectResponse
    {
        $validated = $request->validate([
            'password' => ['required', 'string', 'min:8'],
        ]);

        try {
            $this->agentClient->resetDatabaseUserPassword([
                'engine' => $databaseUser->engine,
                'username' => $databaseUser->username,
                'host' => $databaseUser->host,
                'new_password' => $validated['password'],
            ]);

            ActivityLog::create([
                'user_id' => $request->user()->id,
                'user_email' => $request->user()->email,
                'ip_address' => $request->ip() ?: '127.0.0.1',
                'user_agent' => $request->userAgent(),
                'action' => 'database_user.reset_password',
                'resource_type' => 'database_user',
                'resource_id' => (string) $databaseUser->id,
                'status' => 'success',
                'payload_summary' => [
                    'username' => $databaseUser->username,
                    'host' => $databaseUser->host,
                ],
            ]);

            return redirect()->route('databases.index')->with('success', "Password for {$databaseUser->username} updated successfully.");
        } catch (Exception $e) {
            return back()->with('error', "Failed to reset password: " . $e->getMessage());
        }
    }

    public function revokeGrant(Request $request, DatabaseAccess $grant): RedirectResponse
    {
        $grant->delete();

        ActivityLog::create([
            'user_id' => $request->user()->id,
            'user_email' => $request->user()->email,
            'ip_address' => $request->ip() ?: '127.0.0.1',
            'user_agent' => $request->userAgent(),
            'action' => 'database.revoke_access',
            'resource_type' => 'database_access',
            'resource_id' => (string) $grant->id,
            'status' => 'success',
        ]);

        return redirect()->route('databases.index')->with('success', "Revoked database access grant.");
    }

    public function destroyUser(Request $request, DatabaseUser $databaseUser): RedirectResponse
    {
        $username = $databaseUser->username;
        $engine = $databaseUser->engine;
        $host = $databaseUser->host;

        try {
            $this->agentClient->deleteDatabaseUser($engine, $username, $host);
        } catch (Exception $e) {
            // Non-fatal if server cleanup fails
        }

        $databaseUser->delete();

        ActivityLog::create([
            'user_id' => $request->user()->id,
            'user_email' => $request->user()->email,
            'ip_address' => $request->ip() ?: '127.0.0.1',
            'user_agent' => $request->userAgent(),
            'action' => 'database_user.delete',
            'resource_type' => 'database_user',
            'resource_id' => (string) $databaseUser->id,
            'status' => 'success',
            'payload_summary' => [
                'engine' => $engine,
                'username' => $username,
            ],
        ]);

        return redirect()->route('databases.index')->with('success', "User {$username} dropped successfully.");
    }

    public function destroy(Request $request, Database $database): RedirectResponse
    {
        $name = $database->name;
        $engine = $database->engine;

        try {
            $this->agentClient->deleteDatabase($engine, $name);
        } catch (Exception $e) {
            // Non-fatal if server cleanup fails
        }

        $database->delete();

        ActivityLog::create([
            'user_id' => $request->user()->id,
            'user_email' => $request->user()->email,
            'ip_address' => $request->ip() ?: '127.0.0.1',
            'user_agent' => $request->userAgent(),
            'action' => 'database.delete',
            'resource_type' => 'database',
            'resource_id' => (string) $database->id,
            'status' => 'success',
            'payload_summary' => [
                'engine' => $engine,
                'name' => $name,
            ],
        ]);

        return redirect()->route('databases.index')->with('success', "Database {$name} dropped successfully.");
    }
}
