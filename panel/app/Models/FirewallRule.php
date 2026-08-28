<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Factories\HasFactory;
use Illuminate\Database\Eloquent\Model;

class FirewallRule extends Model
{
    use HasFactory;

    protected $fillable = [
        'port',
        'protocol',
        'action',
        'source_ip',
        'description',
    ];

    protected function casts(): array
    {
        return [
            'port' => 'integer',
        ];
    }
}
