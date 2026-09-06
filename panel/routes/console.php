<?php

use App\Models\Role;
use App\Models\User;
use Illuminate\Foundation\Inspiring;
use Illuminate\Support\Facades\Artisan;
use Illuminate\Support\Facades\Hash;
use Illuminate\Support\Str;

Artisan::command('inspire', function () {
    $this->comment(Inspiring::quote());
})->purpose('Display an inspiring quote');

Artisan::command('panel:reset-admin {email? : The administrator email address} {--password= : The new password (leave empty to prompt or auto-generate)}', function (?string $email = null) {
    $this->info('=========================================');
    $this->info('  Kodepreneur Panel - Admin Account Reset');
    $this->info('=========================================');

    if (! $email) {
        $users = User::all(['id', 'name', 'email']);
        if ($users->isNotEmpty()) {
            $this->comment('Existing accounts found:');
            foreach ($users as $u) {
                $this->line("  • [ID: {$u->id}] {$u->name} ({$u->email})");
            }
            $this->newLine();
        }

        $email = $this->ask('Enter admin email address to reset');
    }

    if (! filter_var($email, FILTER_VALIDATE_EMAIL)) {
        $this->error("Invalid email address: {$email}");
        return 1;
    }

    $password = $this->option('password');
    if (! $password) {
        $password = $this->secret('Enter new password (leave blank to auto-generate)');
        if (! $password) {
            $password = Str::random(16);
            $this->warn("Auto-generated random password: {$password}");
        }
    }

    $superAdminRole = Role::firstOrCreate(
        ['slug' => 'super-admin'],
        ['name' => 'Super Administrator', 'permissions' => ['*']]
    );

    $user = User::where('email', $email)->first();

    if ($user) {
        $user->password = Hash::make($password);
        $user->save();
        $this->info("✓ Password successfully updated for user: {$user->email}");
    } else {
        $name = $this->ask('User not found. Enter name to create new account', 'Administrator');
        $user = User::create([
            'name' => $name,
            'email' => $email,
            'password' => Hash::make($password),
            'role_id' => $superAdminRole->id,
        ]);
        $this->info("✓ New administrator account created: {$user->email}");
    }

    $this->newLine();
    $this->line("Email:    <comment>{$user->email}</comment>");
    $this->line("Password: <comment>{$password}</comment>");
    $this->newLine();
    $this->info('You can now log in to Kodepreneur Panel using these credentials.');

    return 0;
})->purpose('Reset or create an administrator account for Kodepreneur Panel');
