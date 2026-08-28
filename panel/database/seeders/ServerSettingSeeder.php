<?php

namespace Database\Seeders;

use App\Models\PhpVersion;
use App\Models\ServerSetting;
use Illuminate\Database\Seeder;

class ServerSettingSeeder extends Seeder
{
    public function run(): void
    {
        ServerSetting::set('panel_name', 'Kodepreneur Panel', 'string', 'Display name of the control panel');
        ServerSetting::set('panel_port', '8443', 'integer', 'HTTPS port for the control panel web interface');
        ServerSetting::set('auto_update_agent', 'true', 'boolean', 'Automatically check for agent daemon updates');
        ServerSetting::set('default_php_version', '8.3', 'string', 'Default PHP runtime version for new websites');

        PhpVersion::updateOrCreate(
            ['version' => '8.3'],
            [
                'status' => 'installed',
                'is_default' => true,
                'binary_path' => '/usr/bin/php8.3',
                'fpm_socket_path' => '/run/php/php8.3-fpm.sock',
            ]
        );

        PhpVersion::updateOrCreate(
            ['version' => '8.4'],
            [
                'status' => 'installed',
                'is_default' => false,
                'binary_path' => '/usr/bin/php8.4',
                'fpm_socket_path' => '/run/php/php8.4-fpm.sock',
            ]
        );
    }
}
