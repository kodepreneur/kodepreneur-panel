<?php

namespace Database\Seeders;

use App\Models\Role;
use App\Models\User;
use Illuminate\Database\Seeder;
use Illuminate\Support\Facades\Hash;

class AdminUserSeeder extends Seeder
{
    public function run(): void
    {
        $superAdminRole = Role::where('slug', 'super-admin')->first();

        User::updateOrCreate(
            ['email' => 'admin@kodepreneur.com'],
            [
                'name' => 'Kodepreneur Admin',
                'password' => Hash::make('password'),
                'role_id' => $superAdminRole?->id,
            ]
        );
    }
}
