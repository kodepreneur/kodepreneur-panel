<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Factories\HasFactory;
use Illuminate\Database\Eloquent\Model;
use Illuminate\Database\Eloquent\Relations\BelongsToMany;

class DatabaseUser extends Model
{
    use HasFactory;

    protected $fillable = [
        'engine',
        'username',
        'host',
    ];

    public function databases(): BelongsToMany
    {
        return $this->belongsToMany(Database::class, 'database_access', 'database_user_id', 'database_id')
            ->withPivot('permissions')
            ->withTimestamps();
    }
}
