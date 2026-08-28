/// <reference types="vite/client" />

export interface User {
    id: string;
    name: string;
    email: string;
    role?: {
        id: string;
        name: string;
        slug: string;
        permissions: string[];
    };
    created_at?: string;
}

export interface SystemMetrics {
    cpu: {
        usage_percent: number;
        cores: number;
    };
    memory: {
        total_bytes: number;
        used_bytes: number;
        free_bytes: number;
        usage_percent: number;
    };
    disk: {
        mount: string;
        total_bytes: number;
        used_bytes: number;
        free_bytes: number;
        usage_percent: number;
    };
    load_average: {
        load1: number;
        load5: number;
        load15: number;
    };
    timestamp: number;
}

export interface SystemInfo {
    hostname: string;
    os: string;
    kernel: string;
    architecture: string;
    public_ip: string;
    uptime_seconds: number;
}

export interface ServiceStatus {
    name: string;
    is_active: boolean;
    is_enabled: boolean;
    status: string;
}

export interface Domain {
    id: string;
    website_id: string;
    domain: string;
    is_primary: boolean;
    ssl_status: string;
    created_at: string;
}

export interface Website {
    id: string;
    domain: string;
    aliases: string[] | null;
    php_version: string;
    document_root: string;
    system_user: string;
    ssl_enabled: boolean;
    force_https: boolean;
    status: 'active' | 'suspended' | 'provisioning' | 'error';
    domains?: Domain[];
    ssl_certificate?: SslCertificate | null;
    created_at: string;
    updated_at: string;
}

export interface Deployment {
    id: number | string;
    website_id: number | string;
    commit_hash: string | null;
    commit_message: string | null;
    branch: string;
    status: 'pending' | 'running' | 'success' | 'failed' | 'cancelled';
    trigger_source: string;
    log_output: string | null;
    duration_seconds: number | null;
    initiated_by_user_id: number | string | null;
    created_at: string;
    updated_at: string;
}

export interface SslCertificate {
    id: string;
    website_id: string;
    domain: string;
    issuer: string;
    cert_path: string | null;
    key_path: string | null;
    valid_from: string | null;
    valid_until: string | null;
    auto_renew: boolean;
    status: string;
    created_at: string;
}

export interface DatabaseUser {
    id: string | number;
    engine: 'mysql' | 'postgresql';
    username: string;
    host: string;
    databases?: DatabaseRecord[];
    created_at: string;
}

export interface DatabaseAccess {
    id: string | number;
    database_id: string | number;
    database_user_id: string | number;
    permissions: string;
    database?: DatabaseRecord;
    database_user?: DatabaseUser;
    created_at: string;
}

export interface DatabaseRecord {
    id: string | number;
    engine: 'mysql' | 'postgresql';
    name: string;
    character_set: string;
    collation: string;
    size_bytes: number | null;
    users?: DatabaseUser[];
    created_at: string;
}

export interface CronJob {
    id: string | number;
    website_id: string | number | null;
    schedule: string;
    command: string;
    system_user: string;
    is_active: boolean;
    last_run_at: string | null;
    next_run_at: string | null;
    website?: Website;
    created_at: string;
}

export interface Process {
    id: string | number;
    website_id: string | number | null;
    name: string;
    command: string;
    system_user: string;
    instances: number;
    status: 'running' | 'stopped' | 'failed';
    unit_file_path: string;
    website?: Website;
    created_at: string;
}

export interface FirewallRule {
    id: string | number;
    port: number;
    protocol: 'tcp' | 'udp' | 'both';
    action: 'allow' | 'deny';
    source_ip: string | null;
    description: string | null;
    created_at: string;
}

export interface FileEntry {
    name: string;
    path: string;
    is_dir: boolean;
    size_bytes: number;
    permissions: string;
    modified_at: string;
}

export interface ActivityLog {
    id: string;
    user_email: string;
    ip_address: string;
    action: string;
    resource_type: string;
    resource_id: string | null;
    status: 'success' | 'failure';
    payload_summary: Record<string, any> | null;
    created_at: string;
}

export type PageProps<T extends Record<string, unknown> = Record<string, unknown>> = T & {
    auth: {
        user: User | null;
    };
    flash: {
        success?: string;
        error?: string;
    };
    errors: Record<string, string>;
};
