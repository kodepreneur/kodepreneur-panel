<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Factories\HasFactory;
use Illuminate\Database\Eloquent\Model;
use Illuminate\Database\Eloquent\Relations\BelongsToMany;
use Illuminate\Database\Eloquent\Relations\HasMany;

class Database extends Model
{
    use HasFactory;

    protected $table = 'databases';

    protected $fillable = [
        'engine',
        'name',
        'character_set',
        'collation',
        'size_bytes',
    ];

    public function users(): BelongsToMany
    {
        return $this->belongsToMany(DatabaseUser::class, 'database_access', 'database_id', 'database_user_id')
            ->withPivot('permissions')
            ->withTimestamps();
    }

    public function accesses(): HasMany
    {
        return $this->hasMany(DatabaseAccess::class);
    }
}
