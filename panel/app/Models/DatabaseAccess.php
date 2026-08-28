<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Factories\HasFactory;
use Illuminate\Database\Eloquent\Model;
use Illuminate\Database\Eloquent\Relations\BelongsTo;

class DatabaseAccess extends Model
{
    use HasFactory;

    protected $table = 'database_access';

    protected $fillable = [
        'database_id',
        'database_user_id',
        'permissions',
    ];

    public function database(): BelongsTo
    {
        return $this->belongsTo(Database::class);
    }

    public function databaseUser(): BelongsTo
    {
        return $this->belongsTo(DatabaseUser::class);
    }
}
