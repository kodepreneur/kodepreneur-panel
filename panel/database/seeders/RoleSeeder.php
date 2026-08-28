<?php

namespace Database\Seeders;

use App\Models\Role;
use Illuminate\Database\Seeder;

class RoleSeeder extends Seeder
{
    public function run(): void
    {
        Role::updateOrCreate(
            ['slug' => 'super-admin'],
            [
                'name' => 'Super Administrator',
                'permissions' => ['*'],
            ]
        );

        Role::updateOrCreate(
            ['slug' => 'operator'],
            [
                'name' => 'Server Operator',
                'permissions' => [
                    'system.metrics.view',
                    'system.info.view',
                    'websites.view',
                    'websites.create_edit',
                    'deployments.trigger',
                    'databases.manage',
                    'services.restart',
                    'audit.view',
                ],
            ]
        );

        Role::updateOrCreate(
            ['slug' => 'developer'],
            [
                'name' => 'Developer',
                'permissions' => [
                    'system.metrics.view',
                    'websites.view',
                    'deployments.trigger',
                ],
            ]
        );
    }
}
