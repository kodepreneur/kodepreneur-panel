<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Factories\HasFactory;
use Illuminate\Database\Eloquent\Model;

class PhpVersion extends Model
{
    use HasFactory;

    protected $fillable = [
        'version',
        'status',
        'is_default',
        'binary_path',
        'fpm_socket_path',
    ];

    protected function casts(): array
    {
        return [
            'is_default' => 'boolean',
        ];
    }
}
