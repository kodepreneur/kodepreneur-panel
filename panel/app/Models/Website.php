<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Factories\HasFactory;
use Illuminate\Database\Eloquent\Model;
use Illuminate\Database\Eloquent\Relations\HasMany;
use Illuminate\Database\Eloquent\Relations\HasOne;

class Website extends Model
{
    use HasFactory;

    protected $fillable = [
        'domain',
        'aliases',
        'php_version',
        'document_root',
        'system_user',
        'ssl_enabled',
        'force_https',
        'status',
        'custom_nginx_config',
        'deployment_source',
        'project_type',
        'git_repository',
        'git_branch',
        'last_deployed_at',
    ];

    protected function casts(): array
    {
        return [
            'aliases' => 'array',
            'ssl_enabled' => 'boolean',
            'force_https' => 'boolean',
            'last_deployed_at' => 'datetime',
        ];
    }

    public function domains(): HasMany
    {
        return $this->hasMany(Domain::class);
    }

    public function deployments(): HasMany
    {
        return $this->hasMany(Deployment::class)->latest();
    }

    public function latestDeployment(): HasOne
    {
        return $this->hasOne(Deployment::class)->latestOfMany();
    }

    public function processes(): HasMany
    {
        return $this->hasMany(Process::class);
    }

    public function cronJobs(): HasMany
    {
        return $this->hasMany(CronJob::class);
    }

    public function sslCertificate(): HasOne
    {
        return $this->hasOne(SslCertificate::class)->latestOfMany();
    }
}
