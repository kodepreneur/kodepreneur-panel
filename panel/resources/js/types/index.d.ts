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
    deployment_source?: 'empty' | 'zip' | 'git';
    project_type?: string;
    git_repository?: string | null;
    git_branch?: string | null;
    git_auth_type?: 'none' | 'ssh_key' | 'token';
    git_token_user?: string | null;
    git_ssh_public_key?: string | null;
    last_deployed_at?: string | null;
    domains?: Domain[];
    deployments?: Deployment[];
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

export interface RequestLogEntry {
    timestamp: string;
    client_ip: string;
    method: string;
    path: string;
    protocol: string;
    status_code: number;
    bytes_sent: number;
    referer: string;
    user_agent: string;
    browser: string;
}

export interface TimeSeriesDataPoint {
    timestamp: string;
    label: string;
    requests: number;
    bytes_sent: number;
    success_2xx: number;
    redirect_3xx: number;
    client_err_4xx: number;
    server_err_5xx: number;
}

export interface TopMetricItem {
    key: string;
    count: number;
    bytes_sent?: number;
    percentage: number;
}

export interface WebsiteTrafficSummary {
    domain: string;
    period: 'today' | '24h' | '7d' | '30d' | string;
    total_requests: number;
    total_bytes_sent: number;
    unique_visitors: number;
    success_rate: number;
    status_codes: Record<string, number>;
    status_categories: {
        '2xx': number;
        '3xx': number;
        '4xx': number;
        '5xx': number;
    };
    time_series: TimeSeriesDataPoint[];
    top_paths: TopMetricItem[];
    top_ips: TopMetricItem[];
    top_referrers: TopMetricItem[];
    top_user_agents: TopMetricItem[];
    recent_requests: RequestLogEntry[];
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

export interface TableInfo {
    name: string;
    type: string;
    engine: string;
    rows: number;
    data_length: number;
    index_length: number;
    total_size: number;
    collation: string;
    comment: string;
}

export interface ColumnInfo {
    name: string;
    position: number;
    type: string;
    data_type: string;
    is_nullable: boolean;
    key: string;
    default: string | null;
    extra: string;
    comment: string;
}

export interface IndexInfo {
    name: string;
    column: string;
    non_unique: boolean;
    is_primary: boolean;
    seq_in_index: number;
    type: string;
}

export interface ForeignKeyInfo {
    constraint_name: string;
    column: string;
    referenced_table: string;
    referenced_column: string;
}

export interface TableStructure {
    table_name: string;
    columns: ColumnInfo[];
    indexes: IndexInfo[];
    foreign_keys: ForeignKeyInfo[];
    create_statement: string;
}

export interface TableDataResult {
    table_name: string;
    columns: string[];
    rows: Record<string, any>[];
    total_rows: number;
    page: number;
    per_page: number;
    total_pages: number;
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
    mode_octal?: string;
    owner?: string;
    group?: string;
    modified_at: string;
    item_count?: number;
    mime_type?: string;
    extension?: string;
}

export interface FileDetails {
    name: string;
    path: string;
    is_dir: boolean;
    size_bytes: number;
    permissions: string;
    mode_octal: string;
    owner: string;
    group: string;
    uid: number;
    gid: number;
    modified_at: string;
    created_at: string;
    mime_type: string;
    extension: string;
    item_count: number;
}

export interface DiskUsageInfo {
    path: string;
    total_bytes: number;
    used_bytes: number;
    free_bytes: number;
    usage_percent: number;
    path_size: number;
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

export interface GitCommitInfo {
    hash: string;
    short_hash: string;
    message: string;
    author: string;
    date: string;
    url?: string;
}

export interface UpdateInfo {
    current_version: string;
    current_commit: string;
    current_commit_date?: string;
    branch: string;
    repository: string;
    latest_commit: GitCommitInfo | null;
    has_update: boolean;
    behind_by?: number;
    last_checked_at: string;
    laravel_version: string;
    php_version: string;
    agent_status: 'active' | 'degraded' | 'unreachable' | string;
    recent_commits: GitCommitInfo[];
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
