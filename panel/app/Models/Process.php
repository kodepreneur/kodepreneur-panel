<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Factories\HasFactory;
use Illuminate\Database\Eloquent\Model;
use Illuminate\Database\Eloquent\Relations\BelongsTo;

class Process extends Model
{
    use HasFactory;

    protected $fillable = [
        'name',
        'website_id',
        'command',
        'system_user',
        'instances',
        'status',
        'unit_file_path',
    ];

    public function website(): BelongsTo
    {
        return $this->belongsTo(Website::class);
    }
}
