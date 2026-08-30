<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue';
import { Link, router, useForm } from '@inertiajs/vue3';
import AppLayout from '@/Layouts/AppLayout.vue';
import {
    Globe,
    ArrowLeft,
    ShieldCheck,
    ShieldAlert,
    Cpu,
    Folder,
    User as UserIcon,
    GitBranch,
    Play,
    Terminal,
    RefreshCw,
    ExternalLink,
    Lock,
    Unlock,
    Settings,
    Check,
    FileText,
    Server,
    Key,
    Copy,
    Activity,
    BarChart3,
    TrendingUp,
    Zap,
    Users,
    HardDrive,
    Search,
    Compass,
    Laptop,
    Radio,
    Clock,
    AlertTriangle,
    CheckCircle2,
    ArrowUpRight,
} from 'lucide-vue-next';
import type { Website, Deployment, SslCertificate, Domain, WebsiteTrafficSummary, TimeSeriesDataPoint, RequestLogEntry } from '@/types';

const props = defineProps<{
    website: Website & {
        domains: Domain[];
        deployments: Deployment[];
        ssl_certificate?: SslCertificate | null;
    };
}>();

const activeTab = ref<'overview' | 'traffic' | 'php' | 'ssl' | 'logs' | 'deployments'>('overview');

const copiedShowDeployKey = ref(false);
function copyShowDeployKey() {
    if (!props.website.git_ssh_public_key) return;
    navigator.clipboard.writeText(props.website.git_ssh_public_key);
    copiedShowDeployKey.value = true;
    setTimeout(() => {
        copiedShowDeployKey.value = false;
    }, 2000);
}

// PHP Switch Form
const phpForm = useForm({
    php_version: props.website.php_version,
});

function switchPhp() {
    phpForm.put(`/websites/${props.website.id}/php`, {
        preserveScroll: true,
    });
}

// SSL Issue Form
const sslForm = useForm({
    email: '',
    force_https: true,
});

function issueSsl() {
    sslForm.post(`/websites/${props.website.id}/ssl`, {
        preserveScroll: true,
    });
}

// Settings Update Form
const settingsForm = useForm({
    document_root: props.website.document_root,
    force_https: props.website.force_https,
});

function updateSettings() {
    settingsForm.put(`/websites/${props.website.id}`, {
        preserveScroll: true,
    });
}

// Live Logs State
const logType = ref<'access' | 'error'>('access');
const logLines = ref(100);
const logs = ref<string[]>([]);
const isFetchingLogs = ref(false);

async function fetchLogs() {
    isFetchingLogs.value = true;
    try {
        const res = await fetch(`/websites/${props.website.id}/logs/${logType.value}?lines=${logLines.value}`);
        const data = await res.json();
        if (data.success && Array.isArray(data.lines)) {
            logs.value = data.lines;
        } else {
            logs.value = ['No logs available or failed to load.'];
        }
    } catch (e: any) {
        logs.value = [`Error fetching logs: ${e.message}`];
    } finally {
        isFetchingLogs.value = false;
    }
}

// Traffic Monitor State
const trafficPeriod = ref<'24h' | 'today' | '7d' | '30d'>('24h');
const trafficData = ref<WebsiteTrafficSummary | null>(null);
const isFetchingTraffic = ref(false);
const autoPollTraffic = ref(false);
const trafficSearchQuery = ref('');
const activeInsightsTab = ref<'paths' | 'ips' | 'referrers' | 'browsers'>('paths');
const chartMetricType = ref<'requests' | 'bandwidth'>('requests');
const hoveredPointIndex = ref<number | null>(null);

let trafficPollingInterval: any = null;

async function fetchTraffic(showLoading = true) {
    if (showLoading) isFetchingTraffic.value = true;
    try {
        const res = await fetch(`/websites/${props.website.id}/traffic?period=${trafficPeriod.value}`);
        const data = await res.json();
        if (data.success && data.data) {
            trafficData.value = data.data;
        }
    } catch (e: any) {
        console.error('Failed to fetch traffic stats', e);
    } finally {
        if (showLoading) isFetchingTraffic.value = false;
    }
}

function setTrafficPeriod(period: '24h' | 'today' | '7d' | '30d') {
    trafficPeriod.value = period;
    fetchTraffic(true);
}

function toggleAutoPollTraffic() {
    autoPollTraffic.value = !autoPollTraffic.value;
    if (autoPollTraffic.value) {
        trafficPollingInterval = setInterval(() => {
            if (activeTab.value === 'traffic') {
                fetchTraffic(false);
            }
        }, 10000);
    } else {
        if (trafficPollingInterval) {
            clearInterval(trafficPollingInterval);
            trafficPollingInterval = null;
        }
    }
}

onUnmounted(() => {
    if (trafficPollingInterval) {
        clearInterval(trafficPollingInterval);
    }
});

function formatBytes(bytes: number): string {
    if (!bytes || bytes === 0) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i];
}

function formatNumber(num: number): string {
    return new Intl.NumberFormat().format(num || 0);
}

function timeAgo(dateStr: string): string {
    try {
        const d = new Date(dateStr);
        const now = new Date();
        const diffSec = Math.floor((now.getTime() - d.getTime()) / 1000);
        if (diffSec < 10) return 'Just now';
        if (diffSec < 60) return `${diffSec}s ago`;
        const diffMin = Math.floor(diffSec / 60);
        if (diffMin < 60) return `${diffMin}m ago`;
        const diffHr = Math.floor(diffMin / 60);
        if (diffHr < 24) return `${diffHr}h ago`;
        return d.toLocaleDateString();
    } catch {
        return dateStr;
    }
}

// SVG Chart Path Generator
const chartPoints = computed(() => {
    if (!trafficData.value?.time_series || trafficData.value.time_series.length === 0) return [];
    return trafficData.value.time_series;
});

const maxChartValue = computed(() => {
    if (chartPoints.value.length === 0) return 100;
    const values = chartPoints.value.map(p => chartMetricType.value === 'requests' ? p.requests : p.bytes_sent);
    const max = Math.max(...values, 1);
    return Math.ceil(max * 1.15); // Add headroom
});

const svgCoordinates = computed(() => {
    const points = chartPoints.value;
    if (points.length === 0) return [];
    const width = 800;
    const height = 200;
    const padding = 20;
    const maxVal = maxChartValue.value;

    return points.map((p, index) => {
        const x = padding + (index / Math.max(1, points.length - 1)) * (width - padding * 2);
        const val = chartMetricType.value === 'requests' ? p.requests : p.bytes_sent;
        const y = height - padding - (val / maxVal) * (height - padding * 2);
        return { x, y, point: p, val };
    });
});

const svgLinePath = computed(() => {
    const coords = svgCoordinates.value;
    if (coords.length === 0) return '';
    if (coords.length === 1) return `M ${coords[0].x} ${coords[0].y}`;

    let path = `M ${coords[0].x} ${coords[0].y}`;
    for (let i = 0; i < coords.length - 1; i++) {
        const p0 = coords[i];
        const p1 = coords[i + 1];
        const cpX = (p0.x + p1.x) / 2;
        path += ` C ${cpX} ${p0.y}, ${cpX} ${p1.y}, ${p1.x} ${p1.y}`;
    }
    return path;
});

const svgAreaPath = computed(() => {
    const coords = svgCoordinates.value;
    if (coords.length === 0) return '';
    const height = 200;
    const padding = 20;
    const line = svgLinePath.value;
    const last = coords[coords.length - 1];
    const first = coords[0];
    return `${line} L ${last.x} ${height - padding} L ${first.x} ${height - padding} Z`;
});

const filteredRecentRequests = computed(() => {
    if (!trafficData.value?.recent_requests) return [];
    const q = trafficSearchQuery.value.trim().toLowerCase();
    if (!q) return trafficData.value.recent_requests;
    return trafficData.value.recent_requests.filter(
        r => r.path.toLowerCase().includes(q) || r.client_ip.toLowerCase().includes(q) || r.method.toLowerCase().includes(q)
    );
});

function triggerDeploy() {
    if (confirm(`Trigger a deployment build for ${props.website.domain}?`)) {
        router.post(`/websites/${props.website.id}/deploy`);
    }
}
</script>

<template>
    <AppLayout :title="website.domain">
        <div class="max-w-6xl mx-auto space-y-6">
            <!-- Header -->
            <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
                <div class="flex items-center gap-3">
                    <Link
                        href="/websites"
                        class="p-2 rounded-xl bg-white dark:bg-surface-900 border border-slate-200/80 dark:border-surface-800 text-slate-500 hover:text-slate-900 dark:text-surface-400 dark:hover:text-white transition shadow-sm"
                    >
                        <ArrowLeft class="w-4 h-4" />
                    </Link>
                    <div>
                        <div class="flex items-center gap-2.5">
                            <h2 class="text-base font-bold text-slate-900 dark:text-white tracking-tight">{{ website.domain }}</h2>
                            <span class="px-2 py-0.5 rounded-full text-[10px] font-medium bg-emerald-50 text-emerald-600 border border-emerald-200 dark:bg-emerald-500/10 dark:border-emerald-500/30 dark:text-emerald-400">
                                {{ website.status }}
                            </span>
                        </div>
                        <div class="flex items-center gap-3 text-xs text-slate-500 dark:text-surface-400 font-mono mt-0.5">
                            <span>user: {{ website.system_user }}</span>
                            <span>•</span>
                            <a
                                :href="`http://${website.domain}`"
                                target="_blank"
                                class="text-brand-600 dark:text-brand-400 hover:underline flex items-center gap-1"
                            >
                                <span>Visit Site</span>
                                <ExternalLink class="w-3 h-3" />
                            </a>
                        </div>
                    </div>
                </div>

                <div class="flex items-center gap-2">
                    <button
                        @click="triggerDeploy"
                        class="px-4 py-2 rounded-xl bg-brand-600 hover:bg-brand-500 text-white text-xs font-semibold shadow-md shadow-brand-600/20 flex items-center gap-1.5 transition"
                    >
                        <Play class="w-3.5 h-3.5 fill-current" />
                        <span>Deploy Now</span>
                    </button>
                </div>
            </div>

            <!-- Navigation Tabs -->
            <div class="flex items-center gap-2 border-b border-slate-200/80 dark:border-surface-800 pb-2 overflow-x-auto text-xs">
                <button
                    @click="activeTab = 'overview'"
                    :class="[
                        'px-3.5 py-2 rounded-xl font-semibold flex items-center gap-1.5 transition whitespace-nowrap',
                        activeTab === 'overview'
                            ? 'bg-brand-50 text-brand-600 border border-brand-200 dark:bg-brand-500/10 dark:text-brand-400 dark:border-brand-500/30'
                            : 'text-slate-600 hover:text-slate-900 dark:text-surface-400 dark:hover:text-white'
                    ]"
                >
                    <Globe class="w-3.5 h-3.5" />
                    <span>Overview</span>
                </button>

                <button
                    @click="{ activeTab = 'traffic'; fetchTraffic(); }"
                    :class="[
                        'px-3.5 py-2 rounded-xl font-semibold flex items-center gap-1.5 transition whitespace-nowrap',
                        activeTab === 'traffic'
                            ? 'bg-brand-50 text-brand-600 border border-brand-200 dark:bg-brand-500/10 dark:text-brand-400 dark:border-brand-500/30'
                            : 'text-slate-600 hover:text-slate-900 dark:text-surface-400 dark:hover:text-white'
                    ]"
                >
                    <Activity class="w-3.5 h-3.5 text-indigo-500 dark:text-brand-400" />
                    <span>Traffic Monitor</span>
                    <span class="w-2 h-2 rounded-full bg-emerald-500 animate-pulse ml-0.5"></span>
                </button>

                <button
                    @click="activeTab = 'php'"
                    :class="[
                        'px-3.5 py-2 rounded-xl font-semibold flex items-center gap-1.5 transition whitespace-nowrap',
                        activeTab === 'php'
                            ? 'bg-brand-50 text-brand-600 border border-brand-200 dark:bg-brand-500/10 dark:text-brand-400 dark:border-brand-500/30'
                            : 'text-slate-600 hover:text-slate-900 dark:text-surface-400 dark:hover:text-white'
                    ]"
                >
                    <Cpu class="w-3.5 h-3.5" />
                    <span>PHP Runtime</span>
                </button>

                <button
                    @click="activeTab = 'ssl'"
                    :class="[
                        'px-3.5 py-2 rounded-xl font-semibold flex items-center gap-1.5 transition whitespace-nowrap',
                        activeTab === 'ssl'
                            ? 'bg-brand-50 text-brand-600 border border-brand-200 dark:bg-brand-500/10 dark:text-brand-400 dark:border-brand-500/30'
                            : 'text-slate-600 hover:text-slate-900 dark:text-surface-400 dark:hover:text-white'
                    ]"
                >
                    <ShieldCheck class="w-3.5 h-3.5" />
                    <span>SSL / HTTPS</span>
                </button>

                <button
                    @click="{ activeTab = 'logs'; fetchLogs(); }"
                    :class="[
                        'px-3.5 py-2 rounded-xl font-semibold flex items-center gap-1.5 transition whitespace-nowrap',
                        activeTab === 'logs'
                            ? 'bg-brand-50 text-brand-600 border border-brand-200 dark:bg-brand-500/10 dark:text-brand-400 dark:border-brand-500/30'
                            : 'text-slate-600 hover:text-slate-900 dark:text-surface-400 dark:hover:text-white'
                    ]"
                >
                    <FileText class="w-3.5 h-3.5" />
                    <span>Nginx Logs</span>
                </button>

                <button
                    @click="activeTab = 'deployments'"
                    :class="[
                        'px-3.5 py-2 rounded-xl font-semibold flex items-center gap-1.5 transition whitespace-nowrap',
                        activeTab === 'deployments'
                            ? 'bg-brand-50 text-brand-600 border border-brand-200 dark:bg-brand-500/10 dark:text-brand-400 dark:border-brand-500/30'
                            : 'text-slate-600 hover:text-slate-900 dark:text-surface-400 dark:hover:text-white'
                    ]"
                >
                    <GitBranch class="w-3.5 h-3.5" />
                    <span>Deployments ({{ website.deployments.length }})</span>
                </button>
            </div>

            <!-- Tab 1: Overview -->
            <div v-if="activeTab === 'overview'" class="space-y-6">
                <!-- Summary Metrics Cards (4 Grid) -->
                <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
                    <div class="rounded-2xl bg-white dark:bg-surface-900/60 border border-slate-200/80 dark:border-surface-800/80 p-5 shadow-sm dark:shadow-xl space-y-2">
                        <div class="flex items-center gap-2 text-xs font-semibold text-slate-500 dark:text-surface-400">
                            <Cpu class="w-4 h-4 text-brand-600 dark:text-brand-400" />
                            <span>PHP Runtime</span>
                        </div>
                        <p class="text-base font-bold text-slate-900 dark:text-white font-mono">
                            {{ website.php_version !== 'none' ? `PHP ${website.php_version}` : 'Static HTML' }}
                        </p>
                        <p class="text-[11px] text-slate-400 dark:text-surface-500 font-mono truncate">
                            Socket: /run/php/php{{ website.php_version }}-fpm-{{ website.system_user }}.sock
                        </p>
                    </div>

                    <div class="rounded-2xl bg-white dark:bg-surface-900/60 border border-slate-200/80 dark:border-surface-800/80 p-5 shadow-sm dark:shadow-xl space-y-2">
                        <div class="flex items-center gap-2 text-xs font-semibold text-slate-500 dark:text-surface-400">
                            <Folder class="w-4 h-4 text-sky-500 dark:text-cyan-400" />
                            <span>Document Root</span>
                        </div>
                        <p class="text-xs font-bold text-slate-900 dark:text-white font-mono truncate">{{ website.document_root }}</p>
                        <p class="text-[11px] text-slate-400 dark:text-surface-500">Physical web root</p>
                    </div>

                    <div class="rounded-2xl bg-white dark:bg-surface-900/60 border border-slate-200/80 dark:border-surface-800/80 p-5 shadow-sm dark:shadow-xl space-y-2">
                        <div class="flex items-center gap-2 text-xs font-semibold text-slate-500 dark:text-surface-400">
                            <ShieldCheck v-if="website.ssl_enabled" class="w-4 h-4 text-emerald-600 dark:text-emerald-400" />
                            <ShieldAlert v-else class="w-4 h-4 text-amber-600 dark:text-amber-400" />
                            <span>SSL Security</span>
                        </div>
                        <p class="text-base font-bold text-slate-900 dark:text-white font-mono">
                            {{ website.ssl_enabled ? "Let's Encrypt" : 'HTTP Only' }}
                        </p>
                        <p class="text-[11px] text-slate-400 dark:text-surface-500">
                            {{ website.ssl_enabled ? 'Auto-renew active' : 'No certificate issued' }}
                        </p>
                    </div>

                    <!-- Traffic Quick Glance Card -->
                    <button
                        @click="{ activeTab = 'traffic'; fetchTraffic(); }"
                        class="rounded-2xl bg-white dark:bg-surface-900/60 border border-slate-200/80 dark:border-surface-800/80 hover:border-brand-500/40 p-5 shadow-sm dark:shadow-xl space-y-2 text-left transition group"
                    >
                        <div class="flex items-center justify-between">
                            <div class="flex items-center gap-2 text-xs font-semibold text-slate-500 dark:text-surface-400">
                                <Activity class="w-4 h-4 text-indigo-500 dark:text-brand-400" />
                                <span>Traffic Monitor</span>
                            </div>
                            <ArrowUpRight class="w-3.5 h-3.5 text-slate-400 group-hover:text-brand-500 group-hover:translate-x-0.5 group-hover:-translate-y-0.5 transition" />
                        </div>
                        <p class="text-base font-bold text-slate-900 dark:text-white font-mono flex items-center gap-2">
                            <span>Live Analytics</span>
                            <span class="w-2 h-2 rounded-full bg-emerald-500 animate-pulse"></span>
                        </p>
                        <p class="text-[11px] text-brand-600 dark:text-brand-400 font-medium">Click to inspect traffic & hits &rarr;</p>
                    </button>
                </div>

                <!-- Git Repository Card if Git Deploy Source -->
                <div v-if="website.git_repository || website.deployment_source === 'git'" class="rounded-2xl bg-white dark:bg-surface-900/60 border border-slate-200/80 dark:border-surface-800/80 p-6 shadow-sm dark:shadow-xl space-y-4">
                    <div class="flex items-center justify-between flex-wrap gap-2">
                        <h3 class="text-xs font-semibold text-slate-900 dark:text-white uppercase tracking-wider flex items-center gap-2">
                            <GitBranch class="w-4 h-4 text-purple-600 dark:text-purple-400" />
                            <span>Git Repository & Continuous Deployment</span>
                        </h3>
                        <div class="flex items-center gap-2">
                            <span v-if="website.git_auth_type === 'ssh_key'" class="px-2.5 py-0.5 rounded-full text-[10px] font-semibold bg-purple-50 text-purple-700 dark:bg-purple-500/10 dark:text-purple-300 border border-purple-200/60 dark:border-purple-500/30 flex items-center gap-1">
                                <Lock class="w-3 h-3" />
                                Private (SSH Deploy Key)
                            </span>
                            <span v-else-if="website.git_auth_type === 'token'" class="px-2.5 py-0.5 rounded-full text-[10px] font-semibold bg-indigo-50 text-indigo-700 dark:bg-indigo-500/10 dark:text-indigo-300 border border-indigo-200/60 dark:border-indigo-500/30 flex items-center gap-1">
                                <Lock class="w-3 h-3" />
                                Private (Token Auth)
                            </span>
                            <span v-else class="px-2.5 py-0.5 rounded-full text-[10px] font-semibold bg-emerald-50 text-emerald-700 dark:bg-emerald-500/10 dark:text-emerald-300 border border-emerald-200/60 dark:border-emerald-500/30 flex items-center gap-1">
                                <Unlock class="w-3 h-3" />
                                Public Repository
                            </span>
                        </div>
                    </div>

                    <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
                        <div class="space-y-1">
                            <p class="text-[11px] text-slate-400 dark:text-surface-500 font-medium">Repository URL</p>
                            <p class="text-xs font-bold text-slate-900 dark:text-white font-mono break-all">{{ website.git_repository || 'Not configured' }}</p>
                        </div>
                        <div class="space-y-1">
                            <p class="text-[11px] text-slate-400 dark:text-surface-500 font-medium">Target Branch</p>
                            <p class="text-xs font-bold text-slate-900 dark:text-white font-mono">{{ website.git_branch || 'main' }}</p>
                        </div>
                    </div>

                    <!-- SSH Public Key Box if available -->
                    <div v-if="website.git_auth_type === 'ssh_key' && website.git_ssh_public_key" class="p-3.5 rounded-xl bg-slate-50 dark:bg-surface-950/60 border border-slate-200/80 dark:border-surface-800 space-y-2">
                        <div class="flex items-center justify-between">
                            <div class="flex items-center gap-1.5 text-xs font-semibold text-slate-700 dark:text-surface-300">
                                <Key class="w-3.5 h-3.5 text-brand-600 dark:text-brand-400" />
                                <span>Deploy Key (Public Key)</span>
                            </div>
                            <button
                                type="button"
                                @click="copyShowDeployKey"
                                class="px-2.5 py-1 text-[11px] font-semibold rounded-lg bg-white dark:bg-surface-800 border border-slate-200 dark:border-surface-700 text-slate-700 dark:text-surface-200 hover:text-brand-600 dark:hover:text-brand-400 transition flex items-center gap-1 shadow-sm"
                            >
                                <Check v-if="copiedShowDeployKey" class="w-3 h-3 text-emerald-500" />
                                <Copy v-else class="w-3 h-3" />
                                <span>{{ copiedShowDeployKey ? 'Copied' : 'Copy Key' }}</span>
                            </button>
                        </div>
                        <textarea
                            :value="website.git_ssh_public_key"
                            readonly
                            rows="2"
                            class="w-full px-3 py-1.5 rounded-lg bg-slate-900 text-slate-200 font-mono text-[10.5px] border border-slate-800 focus:outline-none select-all resize-none leading-relaxed"
                            @click="($event.target as HTMLTextAreaElement).select()"
                        ></textarea>
                    </div>
                </div>

                <!-- Webroot Settings Card -->
                <div class="rounded-2xl bg-white dark:bg-surface-900/60 border border-slate-200/80 dark:border-surface-800/80 p-6 shadow-sm dark:shadow-xl space-y-4">
                    <h3 class="text-xs font-semibold text-slate-900 dark:text-white uppercase tracking-wider flex items-center gap-2">
                        <Settings class="w-4 h-4 text-brand-600 dark:text-brand-400" />
                        <span>Directory & Security Settings</span>
                    </h3>

                    <form @submit.prevent="updateSettings" class="space-y-4 max-w-xl">
                        <div>
                            <label class="block text-xs font-medium text-slate-700 dark:text-surface-300 mb-1.5">Document Root</label>
                            <input
                                v-model="settingsForm.document_root"
                                type="text"
                                class="w-full px-3.5 py-2.5 rounded-xl bg-slate-50 dark:bg-surface-950/60 border border-slate-200 dark:border-surface-800 text-xs text-slate-900 dark:text-white font-mono focus:ring-2 focus:ring-brand-500/40 focus:outline-none transition"
                            />
                        </div>

                        <label class="flex items-center gap-2.5 cursor-pointer">
                            <input
                                v-model="settingsForm.force_https"
                                type="checkbox"
                                class="rounded border-slate-300 dark:border-surface-700 bg-white dark:bg-surface-900 text-brand-600 focus:ring-brand-500"
                            />
                            <span class="text-xs text-slate-700 dark:text-surface-300">Force HTTPS redirect (301)</span>
                        </label>

                        <button
                            type="submit"
                            :disabled="settingsForm.processing"
                            class="px-4 py-2 rounded-xl bg-slate-900 hover:bg-slate-800 dark:bg-surface-800 dark:hover:bg-surface-700 text-white text-xs font-semibold transition disabled:opacity-50"
                        >
                            {{ settingsForm.processing ? 'Saving...' : 'Save Settings' }}
                        </button>
                    </form>
                </div>
            </div>

            <!-- Tab 2: Traffic Monitor (NEW) -->
            <div v-if="activeTab === 'traffic'" class="space-y-6">
                <!-- Controls & Period Switcher Bar -->
                <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 rounded-2xl bg-white dark:bg-surface-900/60 border border-slate-200/80 dark:border-surface-800/80 p-4 shadow-sm dark:shadow-xl">
                    <div class="flex items-center gap-3">
                        <div class="w-9 h-9 rounded-xl bg-indigo-50 dark:bg-indigo-500/10 text-indigo-600 dark:text-indigo-400 border border-indigo-100 dark:border-indigo-500/20 flex items-center justify-center shadow-sm">
                            <Activity class="w-4 h-4" />
                        </div>
                        <div>
                            <div class="flex items-center gap-2">
                                <h3 class="text-xs font-bold text-slate-900 dark:text-white uppercase tracking-wider">Website Traffic & Access Telemetry</h3>
                                <span class="px-2 py-0.5 rounded-full text-[10px] font-semibold bg-emerald-50 text-emerald-600 border border-emerald-200 dark:bg-emerald-500/10 dark:border-emerald-500/30 dark:text-emerald-400 flex items-center gap-1">
                                    <span class="w-1.5 h-1.5 rounded-full bg-emerald-500 animate-pulse"></span>
                                    Nginx Stream
                                </span>
                            </div>
                            <p class="text-[11px] text-slate-500 dark:text-surface-400">Aggregated from /var/log/nginx/{{ website.domain }}.access.log</p>
                        </div>
                    </div>

                    <div class="flex flex-wrap items-center gap-2">
                        <!-- Period Selector -->
                        <div class="flex items-center bg-slate-100 dark:bg-surface-950 p-1 rounded-xl border border-slate-200/80 dark:border-surface-800 text-xs">
                            <button
                                @click="setTrafficPeriod('today')"
                                :class="[
                                    'px-2.5 py-1 rounded-lg font-semibold transition',
                                    trafficPeriod === 'today'
                                        ? 'bg-white dark:bg-surface-800 text-brand-600 dark:text-brand-400 shadow-sm'
                                        : 'text-slate-600 dark:text-surface-400 hover:text-slate-900 dark:hover:text-white'
                                ]"
                            >
                                Today
                            </button>
                            <button
                                @click="setTrafficPeriod('24h')"
                                :class="[
                                    'px-2.5 py-1 rounded-lg font-semibold transition',
                                    trafficPeriod === '24h'
                                        ? 'bg-white dark:bg-surface-800 text-brand-600 dark:text-brand-400 shadow-sm'
                                        : 'text-slate-600 dark:text-surface-400 hover:text-slate-900 dark:hover:text-white'
                                ]"
                            >
                                24 Hours
                            </button>
                            <button
                                @click="setTrafficPeriod('7d')"
                                :class="[
                                    'px-2.5 py-1 rounded-lg font-semibold transition',
                                    trafficPeriod === '7d'
                                        ? 'bg-white dark:bg-surface-800 text-brand-600 dark:text-brand-400 shadow-sm'
                                        : 'text-slate-600 dark:text-surface-400 hover:text-slate-900 dark:hover:text-white'
                                ]"
                            >
                                7 Days
                            </button>
                            <button
                                @click="setTrafficPeriod('30d')"
                                :class="[
                                    'px-2.5 py-1 rounded-lg font-semibold transition',
                                    trafficPeriod === '30d'
                                        ? 'bg-white dark:bg-surface-800 text-brand-600 dark:text-brand-400 shadow-sm'
                                        : 'text-slate-600 dark:text-surface-400 hover:text-slate-900 dark:hover:text-white'
                                ]"
                            >
                                30 Days
                            </button>
                        </div>

                        <!-- Auto Poll Toggle -->
                        <button
                            @click="toggleAutoPollTraffic"
                            :class="[
                                'px-3 py-1.5 rounded-xl border text-xs font-semibold flex items-center gap-1.5 transition shadow-sm',
                                autoPollTraffic
                                    ? 'bg-emerald-50 text-emerald-700 border-emerald-300 dark:bg-emerald-500/10 dark:text-emerald-400 dark:border-emerald-500/30'
                                    : 'bg-white dark:bg-surface-800 text-slate-600 dark:text-surface-300 border-slate-200 dark:border-surface-700 hover:text-slate-900'
                            ]"
                            title="Toggle live telemetry auto-polling (every 10s)"
                        >
                            <Radio :class="['w-3.5 h-3.5', autoPollTraffic ? 'animate-pulse text-emerald-500' : '']" />
                            <span>{{ autoPollTraffic ? 'Live: On' : 'Live: Off' }}</span>
                        </button>

                        <!-- Manual Refresh Button -->
                        <button
                            @click="fetchTraffic(true)"
                            :disabled="isFetchingTraffic"
                            class="p-2 rounded-xl bg-white dark:bg-surface-800 hover:bg-slate-50 dark:hover:bg-surface-700 text-slate-600 dark:text-surface-300 border border-slate-200 dark:border-surface-700 transition shadow-sm"
                            title="Refresh Traffic Data"
                        >
                            <RefreshCw :class="['w-3.5 h-3.5', isFetchingTraffic ? 'animate-spin text-brand-500' : '']" />
                        </button>
                    </div>
                </div>

                <!-- 4 KPI Summary Cards -->
                <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
                    <!-- Total Requests -->
                    <div class="rounded-2xl bg-white dark:bg-surface-900/70 border border-slate-200/80 dark:border-surface-800/80 p-5 shadow-sm dark:shadow-xl relative overflow-hidden transition-all">
                        <div class="flex items-center justify-between">
                            <div class="w-10 h-10 rounded-xl bg-indigo-50 text-indigo-600 dark:bg-indigo-500/10 dark:text-indigo-400 border border-indigo-100 dark:border-indigo-500/20 flex items-center justify-center shadow-sm">
                                <TrendingUp class="w-5 h-5" />
                            </div>
                            <span class="text-xs font-semibold text-slate-500 dark:text-surface-400">Total Requests</span>
                        </div>
                        <div class="mt-4 flex items-baseline justify-between">
                            <span class="text-2xl font-extrabold text-slate-900 dark:text-white tracking-tight font-mono">
                                {{ formatNumber(trafficData?.total_requests ?? 0) }}
                            </span>
                            <span class="text-[11px] text-slate-500 dark:text-surface-400 font-mono">Hits</span>
                        </div>
                        <div class="mt-3 text-[11px] text-slate-400 dark:text-surface-500 flex items-center gap-1.5">
                            <span class="w-1.5 h-1.5 rounded-full bg-indigo-500"></span>
                            <span>Period: {{ trafficPeriod }}</span>
                        </div>
                    </div>

                    <!-- Total Bandwidth -->
                    <div class="rounded-2xl bg-white dark:bg-surface-900/70 border border-slate-200/80 dark:border-surface-800/80 p-5 shadow-sm dark:shadow-xl relative overflow-hidden transition-all">
                        <div class="flex items-center justify-between">
                            <div class="w-10 h-10 rounded-xl bg-sky-50 text-sky-600 dark:bg-sky-500/10 dark:text-sky-400 border border-sky-100 dark:border-sky-500/20 flex items-center justify-center shadow-sm">
                                <HardDrive class="w-5 h-5" />
                            </div>
                            <span class="text-xs font-semibold text-slate-500 dark:text-surface-400">Bandwidth Out</span>
                        </div>
                        <div class="mt-4 flex items-baseline justify-between">
                            <span class="text-2xl font-extrabold text-slate-900 dark:text-white tracking-tight font-mono">
                                {{ formatBytes(trafficData?.total_bytes_sent ?? 0) }}
                            </span>
                            <span class="text-[11px] text-slate-500 dark:text-surface-400 font-mono">Transferred</span>
                        </div>
                        <div class="mt-3 text-[11px] text-slate-400 dark:text-surface-500 flex items-center gap-1.5">
                            <span class="w-1.5 h-1.5 rounded-full bg-sky-500"></span>
                            <span>Outbound payload</span>
                        </div>
                    </div>

                    <!-- Unique Visitors -->
                    <div class="rounded-2xl bg-white dark:bg-surface-900/70 border border-slate-200/80 dark:border-surface-800/80 p-5 shadow-sm dark:shadow-xl relative overflow-hidden transition-all">
                        <div class="flex items-center justify-between">
                            <div class="w-10 h-10 rounded-xl bg-emerald-50 text-emerald-600 dark:bg-emerald-500/10 dark:text-emerald-400 border border-emerald-100 dark:border-emerald-500/20 flex items-center justify-center shadow-sm">
                                <Users class="w-5 h-5" />
                            </div>
                            <span class="text-xs font-semibold text-slate-500 dark:text-surface-400">Unique Visitors</span>
                        </div>
                        <div class="mt-4 flex items-baseline justify-between">
                            <span class="text-2xl font-extrabold text-slate-900 dark:text-white tracking-tight font-mono">
                                {{ formatNumber(trafficData?.unique_visitors ?? 0) }}
                            </span>
                            <span class="text-[11px] text-slate-500 dark:text-surface-400 font-mono">Unique IPs</span>
                        </div>
                        <div class="mt-3 text-[11px] text-slate-400 dark:text-surface-500 flex items-center gap-1.5">
                            <span class="w-1.5 h-1.5 rounded-full bg-emerald-500"></span>
                            <span>Distinct client addresses</span>
                        </div>
                    </div>

                    <!-- Success Rate / Health -->
                    <div class="rounded-2xl bg-white dark:bg-surface-900/70 border border-slate-200/80 dark:border-surface-800/80 p-5 shadow-sm dark:shadow-xl relative overflow-hidden transition-all">
                        <div class="flex items-center justify-between">
                            <div class="w-10 h-10 rounded-xl bg-purple-50 text-purple-600 dark:bg-purple-500/10 dark:text-purple-400 border border-purple-100 dark:border-purple-500/20 flex items-center justify-center shadow-sm">
                                <CheckCircle2 class="w-5 h-5" />
                            </div>
                            <span class="text-xs font-semibold text-slate-500 dark:text-surface-400">HTTP Health</span>
                        </div>
                        <div class="mt-4 flex items-baseline justify-between">
                            <span class="text-2xl font-extrabold text-slate-900 dark:text-white tracking-tight font-mono">
                                {{ (trafficData?.success_rate ?? 100).toFixed(1) }}%
                            </span>
                            <span class="text-[11px] text-slate-500 dark:text-surface-400 font-mono">2xx & 3xx</span>
                        </div>
                        <div class="mt-3 text-[11px] text-slate-400 dark:text-surface-500 flex items-center gap-1.5">
                            <span class="w-1.5 h-1.5 rounded-full bg-purple-500"></span>
                            <span>Errors: {{ (trafficData?.status_categories?.['4xx'] ?? 0) + (trafficData?.status_categories?.['5xx'] ?? 0) }}</span>
                        </div>
                    </div>
                </div>

                <!-- Interactive Traffic Timeline Chart -->
                <div class="rounded-2xl bg-white dark:bg-surface-900/60 border border-slate-200/80 dark:border-surface-800/80 p-6 shadow-sm dark:shadow-xl space-y-4">
                    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3">
                        <div>
                            <h3 class="text-xs font-semibold text-slate-900 dark:text-white uppercase tracking-wider flex items-center gap-2">
                                <BarChart3 class="w-4 h-4 text-brand-600 dark:text-brand-400" />
                                <span>Traffic Volume Over Time ({{ trafficPeriod }})</span>
                            </h3>
                            <p class="text-[11px] text-slate-500 dark:text-surface-400 mt-0.5">
                                Hourly / daily aggregated request and data transfer curve
                            </p>
                        </div>

                        <!-- Toggle Requests vs Bandwidth -->
                        <div class="flex items-center bg-slate-100 dark:bg-surface-950 p-1 rounded-xl border border-slate-200/80 dark:border-surface-800 text-xs">
                            <button
                                @click="chartMetricType = 'requests'"
                                :class="[
                                    'px-3 py-1 rounded-lg font-semibold transition',
                                    chartMetricType === 'requests'
                                        ? 'bg-white dark:bg-surface-800 text-brand-600 dark:text-brand-400 shadow-sm'
                                        : 'text-slate-600 dark:text-surface-400 hover:text-slate-900 dark:hover:text-white'
                                ]"
                            >
                                Requests (Hits)
                            </button>
                            <button
                                @click="chartMetricType = 'bandwidth'"
                                :class="[
                                    'px-3 py-1 rounded-lg font-semibold transition',
                                    chartMetricType === 'bandwidth'
                                        ? 'bg-white dark:bg-surface-800 text-brand-600 dark:text-brand-400 shadow-sm'
                                        : 'text-slate-600 dark:text-surface-400 hover:text-slate-900 dark:hover:text-white'
                                ]"
                            >
                                Bandwidth (Bytes)
                            </button>
                        </div>
                    </div>

                    <!-- SVG Chart Container -->
                    <div class="relative w-full h-64 bg-slate-50/50 dark:bg-surface-950/40 rounded-xl p-3 border border-slate-100 dark:border-surface-800/60 overflow-hidden">
                        <svg class="w-full h-full" viewBox="0 0 800 200" preserveAspectRatio="none">
                            <defs>
                                <linearGradient id="chartGradient" x1="0" y1="0" x2="0" y2="1">
                                    <stop offset="0%" stop-color="#6366f1" stop-opacity="0.35" />
                                    <stop offset="100%" stop-color="#6366f1" stop-opacity="0.00" />
                                </linearGradient>
                                <linearGradient id="lineGradient" x1="0" y1="0" x2="1" y2="0">
                                    <stop offset="0%" stop-color="#818cf8" />
                                    <stop offset="100%" stop-color="#4f46e5" />
                                </linearGradient>
                            </defs>

                            <!-- Horizontal Grid Lines -->
                            <line x1="20" y1="40" x2="780" y2="40" stroke="currentColor" class="text-slate-200 dark:text-surface-800/60" stroke-dasharray="3 3" />
                            <line x1="20" y1="90" x2="780" y2="90" stroke="currentColor" class="text-slate-200 dark:text-surface-800/60" stroke-dasharray="3 3" />
                            <line x1="20" y1="140" x2="780" y2="140" stroke="currentColor" class="text-slate-200 dark:text-surface-800/60" stroke-dasharray="3 3" />

                            <!-- Gradient Area Fill -->
                            <path :d="svgAreaPath" fill="url(#chartGradient)" />

                            <!-- Smooth Curve Line -->
                            <path :d="svgLinePath" fill="none" stroke="url(#lineGradient)" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round" />

                            <!-- Data Point Circles -->
                            <g v-for="(item, idx) in svgCoordinates" :key="idx">
                                <circle
                                    :cx="item.x"
                                    :cy="item.y"
                                    r="4"
                                    class="fill-white dark:fill-surface-900 stroke-indigo-600 dark:stroke-indigo-400 stroke-2 hover:r-6 transition-all cursor-pointer"
                                    @mouseenter="hoveredPointIndex = idx"
                                    @mouseleave="hoveredPointIndex = null"
                                />
                            </g>
                        </svg>

                        <!-- Tooltip on Hover -->
                        <div
                            v-if="hoveredPointIndex !== null && svgCoordinates[hoveredPointIndex]"
                            class="absolute z-20 pointer-events-none transform -translate-x-1/2 -translate-y-full px-3 py-2 rounded-xl bg-slate-900/95 dark:bg-surface-800 text-white text-[11px] shadow-xl border border-slate-700 space-y-1 font-mono"
                            :style="{
                                left: `${(svgCoordinates[hoveredPointIndex].x / 800) * 100}%`,
                                top: `${Math.max(10, (svgCoordinates[hoveredPointIndex].y / 200) * 100 - 8)}%`
                            }"
                        >
                            <div class="flex items-center gap-2 font-bold text-indigo-300">
                                <Clock class="w-3 h-3" />
                                <span>{{ svgCoordinates[hoveredPointIndex].point.label }}</span>
                            </div>
                            <div>Requests: <span class="font-bold text-white">{{ formatNumber(svgCoordinates[hoveredPointIndex].point.requests) }}</span></div>
                            <div>Bandwidth: <span class="font-bold text-sky-300">{{ formatBytes(svgCoordinates[hoveredPointIndex].point.bytes_sent) }}</span></div>
                            <div class="text-[10px] text-slate-400">
                                2xx: {{ svgCoordinates[hoveredPointIndex].point.success_2xx }} • 4xx: {{ svgCoordinates[hoveredPointIndex].point.client_err_4xx }} • 5xx: {{ svgCoordinates[hoveredPointIndex].point.server_err_5xx }}
                            </div>
                        </div>

                        <!-- X Axis Labels -->
                        <div class="absolute bottom-1 left-4 right-4 flex justify-between text-[10px] text-slate-400 dark:text-surface-500 font-mono select-none">
                            <span v-if="chartPoints[0]">{{ chartPoints[0].label }}</span>
                            <span v-if="chartPoints[Math.floor(chartPoints.length / 2)]">{{ chartPoints[Math.floor(chartPoints.length / 2)].label }}</span>
                            <span v-if="chartPoints[chartPoints.length - 1]">{{ chartPoints[chartPoints.length - 1].label }}</span>
                        </div>
                    </div>
                </div>

                <!-- HTTP Status Code Distribution & Health Breakdown -->
                <div class="rounded-2xl bg-white dark:bg-surface-900/60 border border-slate-200/80 dark:border-surface-800/80 p-6 shadow-sm dark:shadow-xl space-y-4">
                    <h3 class="text-xs font-semibold text-slate-900 dark:text-white uppercase tracking-wider flex items-center gap-2">
                        <CheckCircle2 class="w-4 h-4 text-emerald-600 dark:text-emerald-400" />
                        <span>HTTP Status Code Distribution</span>
                    </h3>

                    <!-- Progress Bar Breakdown -->
                    <div class="h-3 w-full bg-slate-100 dark:bg-surface-800 rounded-full overflow-hidden flex shadow-inner">
                        <div
                            class="h-full bg-emerald-500 transition-all duration-500"
                            :style="{ width: `${trafficData?.total_requests ? ((trafficData.status_categories['2xx'] || 0) / trafficData.total_requests) * 100 : 100}%` }"
                            title="2xx Success"
                        ></div>
                        <div
                            class="h-full bg-sky-500 transition-all duration-500"
                            :style="{ width: `${trafficData?.total_requests ? ((trafficData.status_categories['3xx'] || 0) / trafficData.total_requests) * 100 : 0}%` }"
                            title="3xx Redirect"
                        ></div>
                        <div
                            class="h-full bg-amber-500 transition-all duration-500"
                            :style="{ width: `${trafficData?.total_requests ? ((trafficData.status_categories['4xx'] || 0) / trafficData.total_requests) * 100 : 0}%` }"
                            title="4xx Client Error"
                        ></div>
                        <div
                            class="h-full bg-rose-500 transition-all duration-500"
                            :style="{ width: `${trafficData?.total_requests ? ((trafficData.status_categories['5xx'] || 0) / trafficData.total_requests) * 100 : 0}%` }"
                            title="5xx Server Error"
                        ></div>
                    </div>

                    <!-- 4 Category Cards -->
                    <div class="grid grid-cols-2 sm:grid-cols-4 gap-3 pt-2">
                        <!-- 2xx -->
                        <div class="p-3.5 rounded-xl bg-emerald-50/70 dark:bg-emerald-500/10 border border-emerald-200/70 dark:border-emerald-500/20 space-y-1">
                            <span class="text-[11px] font-semibold text-emerald-700 dark:text-emerald-300 flex items-center gap-1.5">
                                <span class="w-2 h-2 rounded-full bg-emerald-500"></span>
                                2xx Success
                            </span>
                            <p class="text-lg font-bold text-slate-900 dark:text-white font-mono">
                                {{ formatNumber(trafficData?.status_categories?.['2xx'] ?? 0) }}
                            </p>
                            <p class="text-[10.5px] text-slate-500 dark:text-surface-400">
                                {{ trafficData?.total_requests ? (((trafficData.status_categories['2xx'] || 0) / trafficData.total_requests) * 100).toFixed(1) : '100' }}% of traffic
                            </p>
                        </div>

                        <!-- 3xx -->
                        <div class="p-3.5 rounded-xl bg-sky-50/70 dark:bg-sky-500/10 border border-sky-200/70 dark:border-sky-500/20 space-y-1">
                            <span class="text-[11px] font-semibold text-sky-700 dark:text-sky-300 flex items-center gap-1.5">
                                <span class="w-2 h-2 rounded-full bg-sky-500"></span>
                                3xx Redirect
                            </span>
                            <p class="text-lg font-bold text-slate-900 dark:text-white font-mono">
                                {{ formatNumber(trafficData?.status_categories?.['3xx'] ?? 0) }}
                            </p>
                            <p class="text-[10.5px] text-slate-500 dark:text-surface-400">
                                {{ trafficData?.total_requests ? (((trafficData.status_categories['3xx'] || 0) / trafficData.total_requests) * 100).toFixed(1) : '0' }}% of traffic
                            </p>
                        </div>

                        <!-- 4xx -->
                        <div class="p-3.5 rounded-xl bg-amber-50/70 dark:bg-amber-500/10 border border-amber-200/70 dark:border-amber-500/20 space-y-1">
                            <span class="text-[11px] font-semibold text-amber-700 dark:text-amber-300 flex items-center gap-1.5">
                                <span class="w-2 h-2 rounded-full bg-amber-500"></span>
                                4xx Client Error
                            </span>
                            <p class="text-lg font-bold text-slate-900 dark:text-white font-mono">
                                {{ formatNumber(trafficData?.status_categories?.['4xx'] ?? 0) }}
                            </p>
                            <p class="text-[10.5px] text-slate-500 dark:text-surface-400">
                                {{ trafficData?.total_requests ? (((trafficData.status_categories['4xx'] || 0) / trafficData.total_requests) * 100).toFixed(1) : '0' }}% of traffic
                            </p>
                        </div>

                        <!-- 5xx -->
                        <div class="p-3.5 rounded-xl bg-rose-50/70 dark:bg-rose-500/10 border border-rose-200/70 dark:border-rose-500/20 space-y-1">
                            <span class="text-[11px] font-semibold text-rose-700 dark:text-rose-300 flex items-center gap-1.5">
                                <span class="w-2 h-2 rounded-full bg-rose-500"></span>
                                5xx Server Error
                            </span>
                            <p class="text-lg font-bold text-slate-900 dark:text-white font-mono">
                                {{ formatNumber(trafficData?.status_categories?.['5xx'] ?? 0) }}
                            </p>
                            <p class="text-[10.5px] text-slate-500 dark:text-surface-400">
                                {{ trafficData?.total_requests ? (((trafficData.status_categories['5xx'] || 0) / trafficData.total_requests) * 100).toFixed(1) : '0' }}% of traffic
                            </p>
                        </div>
                    </div>
                </div>

                <!-- Analytical Insights Tabs (Top Paths, IPs, Referrers, Browsers) -->
                <div class="rounded-2xl bg-white dark:bg-surface-900/60 border border-slate-200/80 dark:border-surface-800/80 p-6 shadow-sm dark:shadow-xl space-y-4">
                    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3 border-b border-slate-200/80 dark:border-surface-800 pb-3">
                        <div class="flex items-center gap-2 text-xs">
                            <button
                                @click="activeInsightsTab = 'paths'"
                                :class="[
                                    'px-3 py-1.5 rounded-xl font-semibold transition',
                                    activeInsightsTab === 'paths'
                                        ? 'bg-brand-50 text-brand-600 border border-brand-200 dark:bg-brand-500/10 dark:text-brand-400 dark:border-brand-500/30'
                                        : 'text-slate-600 hover:text-slate-900 dark:text-surface-400 dark:hover:text-white'
                                ]"
                            >
                                Top Visited Paths
                            </button>
                            <button
                                @click="activeInsightsTab = 'ips'"
                                :class="[
                                    'px-3 py-1.5 rounded-xl font-semibold transition',
                                    activeInsightsTab === 'ips'
                                        ? 'bg-brand-50 text-brand-600 border border-brand-200 dark:bg-brand-500/10 dark:text-brand-400 dark:border-brand-500/30'
                                        : 'text-slate-600 hover:text-slate-900 dark:text-surface-400 dark:hover:text-white'
                                ]"
                            >
                                Top Client IPs
                            </button>
                            <button
                                @click="activeInsightsTab = 'referrers'"
                                :class="[
                                    'px-3 py-1.5 rounded-xl font-semibold transition',
                                    activeInsightsTab === 'referrers'
                                        ? 'bg-brand-50 text-brand-600 border border-brand-200 dark:bg-brand-500/10 dark:text-brand-400 dark:border-brand-500/30'
                                        : 'text-slate-600 hover:text-slate-900 dark:text-surface-400 dark:hover:text-white'
                                ]"
                            >
                                Top Referrers
                            </button>
                            <button
                                @click="activeInsightsTab = 'browsers'"
                                :class="[
                                    'px-3 py-1.5 rounded-xl font-semibold transition',
                                    activeInsightsTab === 'browsers'
                                        ? 'bg-brand-50 text-brand-600 border border-brand-200 dark:bg-brand-500/10 dark:text-brand-400 dark:border-brand-500/30'
                                        : 'text-slate-600 hover:text-slate-900 dark:text-surface-400 dark:hover:text-white'
                                ]"
                            >
                                Browsers & Devices
                            </button>
                        </div>
                    </div>

                    <!-- Tab Content: Paths -->
                    <div v-if="activeInsightsTab === 'paths'" class="space-y-3">
                        <div v-if="!trafficData?.top_paths?.length" class="text-center py-8 text-xs text-slate-400">
                            No visited path metrics available.
                        </div>
                        <div v-for="(pathItem, idx) in trafficData?.top_paths" :key="idx" class="space-y-1">
                            <div class="flex items-center justify-between text-xs font-mono">
                                <span class="font-semibold text-slate-800 dark:text-slate-200 truncate max-w-md">{{ pathItem.key }}</span>
                                <div class="flex items-center gap-3 text-slate-500 dark:text-surface-400 shrink-0">
                                    <span v-if="pathItem.bytes_sent">{{ formatBytes(pathItem.bytes_sent) }}</span>
                                    <span class="font-bold text-slate-900 dark:text-white">{{ formatNumber(pathItem.count) }} hits</span>
                                    <span class="text-brand-600 dark:text-brand-400 font-semibold w-12 text-right">{{ pathItem.percentage.toFixed(1) }}%</span>
                                </div>
                            </div>
                            <div class="h-1.5 w-full bg-slate-100 dark:bg-surface-800 rounded-full overflow-hidden">
                                <div class="h-full bg-indigo-500 rounded-full transition-all duration-300" :style="{ width: `${pathItem.percentage}%` }"></div>
                            </div>
                        </div>
                    </div>

                    <!-- Tab Content: Client IPs -->
                    <div v-if="activeInsightsTab === 'ips'" class="space-y-3">
                        <div v-if="!trafficData?.top_ips?.length" class="text-center py-8 text-xs text-slate-400">
                            No visitor IP metrics available.
                        </div>
                        <div v-for="(ipItem, idx) in trafficData?.top_ips" :key="idx" class="space-y-1">
                            <div class="flex items-center justify-between text-xs font-mono">
                                <span class="font-semibold text-slate-800 dark:text-slate-200">{{ ipItem.key }}</span>
                                <div class="flex items-center gap-3 text-slate-500 dark:text-surface-400 shrink-0">
                                    <span v-if="ipItem.bytes_sent">{{ formatBytes(ipItem.bytes_sent) }}</span>
                                    <span class="font-bold text-slate-900 dark:text-white">{{ formatNumber(ipItem.count) }} hits</span>
                                    <span class="text-emerald-600 dark:text-emerald-400 font-semibold w-12 text-right">{{ ipItem.percentage.toFixed(1) }}%</span>
                                </div>
                            </div>
                            <div class="h-1.5 w-full bg-slate-100 dark:bg-surface-800 rounded-full overflow-hidden">
                                <div class="h-full bg-emerald-500 rounded-full transition-all duration-300" :style="{ width: `${ipItem.percentage}%` }"></div>
                            </div>
                        </div>
                    </div>

                    <!-- Tab Content: Referrers -->
                    <div v-if="activeInsightsTab === 'referrers'" class="space-y-3">
                        <div v-if="!trafficData?.top_referrers?.length" class="text-center py-8 text-xs text-slate-400">
                            No referrer metrics available.
                        </div>
                        <div v-for="(refItem, idx) in trafficData?.top_referrers" :key="idx" class="space-y-1">
                            <div class="flex items-center justify-between text-xs font-mono">
                                <span class="font-semibold text-slate-800 dark:text-slate-200">{{ refItem.key }}</span>
                                <div class="flex items-center gap-3 text-slate-500 dark:text-surface-400 shrink-0">
                                    <span class="font-bold text-slate-900 dark:text-white">{{ formatNumber(refItem.count) }} hits</span>
                                    <span class="text-sky-600 dark:text-sky-400 font-semibold w-12 text-right">{{ refItem.percentage.toFixed(1) }}%</span>
                                </div>
                            </div>
                            <div class="h-1.5 w-full bg-slate-100 dark:bg-surface-800 rounded-full overflow-hidden">
                                <div class="h-full bg-sky-500 rounded-full transition-all duration-300" :style="{ width: `${refItem.percentage}%` }"></div>
                            </div>
                        </div>
                    </div>

                    <!-- Tab Content: Browsers -->
                    <div v-if="activeInsightsTab === 'browsers'" class="space-y-3">
                        <div v-if="!trafficData?.top_user_agents?.length" class="text-center py-8 text-xs text-slate-400">
                            No browser metrics available.
                        </div>
                        <div v-for="(uaItem, idx) in trafficData?.top_user_agents" :key="idx" class="space-y-1">
                            <div class="flex items-center justify-between text-xs font-mono">
                                <span class="font-semibold text-slate-800 dark:text-slate-200 flex items-center gap-2">
                                    <Laptop class="w-3.5 h-3.5 text-slate-400" />
                                    <span>{{ uaItem.key }}</span>
                                </span>
                                <div class="flex items-center gap-3 text-slate-500 dark:text-surface-400 shrink-0">
                                    <span class="font-bold text-slate-900 dark:text-white">{{ formatNumber(uaItem.count) }} hits</span>
                                    <span class="text-purple-600 dark:text-purple-400 font-semibold w-12 text-right">{{ uaItem.percentage.toFixed(1) }}%</span>
                                </div>
                            </div>
                            <div class="h-1.5 w-full bg-slate-100 dark:bg-surface-800 rounded-full overflow-hidden">
                                <div class="h-full bg-purple-500 rounded-full transition-all duration-300" :style="{ width: `${uaItem.percentage}%` }"></div>
                            </div>
                        </div>
                    </div>
                </div>

                <!-- Live Request Stream Table -->
                <div class="rounded-2xl bg-white dark:bg-surface-900/60 border border-slate-200/80 dark:border-surface-800/80 p-6 shadow-sm dark:shadow-xl space-y-4">
                    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3">
                        <div>
                            <h3 class="text-xs font-semibold text-slate-900 dark:text-white uppercase tracking-wider flex items-center gap-2">
                                <Radio class="w-4 h-4 text-emerald-500 animate-pulse" />
                                <span>Recent Access Stream (Latest 50 Requests)</span>
                            </h3>
                            <p class="text-[11px] text-slate-500 dark:text-surface-400 mt-0.5">Live incoming HTTP requests</p>
                        </div>

                        <!-- Filter Input -->
                        <div class="relative w-full sm:w-64">
                            <Search class="w-3.5 h-3.5 absolute left-3 top-1/2 -translate-y-1/2 text-slate-400" />
                            <input
                                v-model="trafficSearchQuery"
                                type="text"
                                placeholder="Filter path, IP, or method..."
                                class="w-full pl-8 pr-3 py-1.5 rounded-xl bg-slate-50 dark:bg-surface-950 border border-slate-200 dark:border-surface-800 text-xs text-slate-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-brand-500/40"
                            />
                        </div>
                    </div>

                    <div class="overflow-x-auto">
                        <table class="w-full text-left border-collapse text-xs font-mono">
                            <thead>
                                <tr class="border-b border-slate-200/80 dark:border-surface-800 text-slate-400 dark:text-surface-500 text-[11px]">
                                    <th class="py-2.5 px-3">Time</th>
                                    <th class="py-2.5 px-3">Method</th>
                                    <th class="py-2.5 px-3">Status</th>
                                    <th class="py-2.5 px-3">Path</th>
                                    <th class="py-2.5 px-3">Size</th>
                                    <th class="py-2.5 px-3">Client IP</th>
                                    <th class="py-2.5 px-3">User Agent</th>
                                </tr>
                            </thead>
                            <tbody class="divide-y divide-slate-100 dark:divide-surface-800/60">
                                <tr v-if="filteredRecentRequests.length === 0">
                                    <td colspan="7" class="py-8 text-center text-slate-400 text-xs italic">
                                        No recent requests matching filter.
                                    </td>
                                </tr>
                                <tr
                                    v-for="(req, idx) in filteredRecentRequests"
                                    :key="idx"
                                    class="hover:bg-slate-50/60 dark:hover:bg-surface-800/40 transition"
                                >
                                    <td class="py-2.5 px-3 text-slate-500 dark:text-surface-400 whitespace-nowrap text-[11px]">
                                        {{ timeAgo(req.timestamp) }}
                                    </td>
                                    <td class="py-2.5 px-3 whitespace-nowrap">
                                        <span
                                            :class="[
                                                'px-2 py-0.5 rounded text-[10px] font-bold uppercase',
                                                req.method === 'GET' ? 'bg-emerald-50 text-emerald-700 dark:bg-emerald-500/10 dark:text-emerald-400 border border-emerald-200 dark:border-emerald-500/20' :
                                                req.method === 'POST' ? 'bg-indigo-50 text-indigo-700 dark:bg-indigo-500/10 dark:text-indigo-400 border border-indigo-200 dark:border-indigo-500/20' :
                                                req.method === 'PUT' ? 'bg-amber-50 text-amber-700 dark:bg-amber-500/10 dark:text-amber-400 border border-amber-200 dark:border-amber-500/20' :
                                                'bg-rose-50 text-rose-700 dark:bg-rose-500/10 dark:text-rose-400 border border-rose-200 dark:border-rose-500/20'
                                            ]"
                                        >
                                            {{ req.method }}
                                        </span>
                                    </td>
                                    <td class="py-2.5 px-3 whitespace-nowrap">
                                        <span
                                            :class="[
                                                'px-2 py-0.5 rounded-full text-[10.5px] font-semibold',
                                                req.status_code >= 200 && req.status_code < 300 ? 'bg-emerald-100/80 text-emerald-800 dark:bg-emerald-500/20 dark:text-emerald-300' :
                                                req.status_code >= 300 && req.status_code < 400 ? 'bg-sky-100/80 text-sky-800 dark:bg-sky-500/20 dark:text-sky-300' :
                                                req.status_code >= 400 && req.status_code < 500 ? 'bg-amber-100/80 text-amber-800 dark:bg-amber-500/20 dark:text-amber-300' :
                                                'bg-rose-100/80 text-rose-800 dark:bg-rose-500/20 dark:text-rose-300'
                                            ]"
                                        >
                                            {{ req.status_code }}
                                        </span>
                                    </td>
                                    <td class="py-2.5 px-3 text-slate-800 dark:text-slate-200 truncate max-w-xs font-semibold" :title="req.path">
                                        {{ req.path }}
                                    </td>
                                    <td class="py-2.5 px-3 text-slate-500 dark:text-surface-400 whitespace-nowrap text-[11px]">
                                        {{ formatBytes(req.bytes_sent) }}
                                    </td>
                                    <td class="py-2.5 px-3 text-slate-600 dark:text-surface-300 whitespace-nowrap text-[11px]">
                                        {{ req.client_ip }}
                                    </td>
                                    <td class="py-2.5 px-3 text-slate-400 dark:text-surface-500 truncate max-w-xs text-[10.5px]" :title="req.user_agent">
                                        {{ req.browser }}
                                    </td>
                                </tr>
                            </tbody>
                        </table>
                    </div>
                </div>
            </div>

            <!-- Tab 3: PHP Runtime -->
            <div v-if="activeTab === 'php'" class="rounded-2xl bg-white dark:bg-surface-900/60 border border-slate-200/80 dark:border-surface-800/80 p-6 shadow-sm dark:shadow-xl space-y-4">
                <h3 class="text-xs font-semibold text-slate-900 dark:text-white uppercase tracking-wider flex items-center gap-2">
                    <Cpu class="w-4 h-4 text-brand-600 dark:text-brand-400" />
                    <span>PHP-FPM Runtime Version</span>
                </h3>

                <p class="text-xs text-slate-500 dark:text-surface-400">
                    Select the PHP engine version for this virtual host. Kodepreneur automatically spins up and manages dedicated, isolated PHP-FPM pools for each website.
                </p>

                <form @submit.prevent="switchPhp" class="space-y-4 max-w-md">
                    <div class="space-y-2">
                        <label class="flex items-center gap-3 p-3.5 rounded-xl border border-slate-200 dark:border-surface-800 cursor-pointer hover:bg-slate-50 dark:hover:bg-surface-800/50 transition">
                            <input
                                v-model="phpForm.php_version"
                                type="radio"
                                value="8.3"
                                class="text-brand-600 focus:ring-brand-500"
                            />
                            <div>
                                <p class="text-xs font-semibold text-slate-900 dark:text-white">PHP 8.3 FPM</p>
                                <p class="text-[11px] text-slate-500 dark:text-surface-400">Standard stable production release</p>
                            </div>
                        </label>

                        <label class="flex items-center gap-3 p-3.5 rounded-xl border border-slate-200 dark:border-surface-800 cursor-pointer hover:bg-slate-50 dark:hover:bg-surface-800/50 transition">
                            <input
                                v-model="phpForm.php_version"
                                type="radio"
                                value="8.4"
                                class="text-brand-600 focus:ring-brand-500"
                            />
                            <div>
                                <p class="text-xs font-semibold text-slate-900 dark:text-white">PHP 8.4 FPM</p>
                                <p class="text-[11px] text-slate-500 dark:text-surface-400">Latest PHP 8.4 engine</p>
                            </div>
                        </label>

                        <label class="flex items-center gap-3 p-3.5 rounded-xl border border-slate-200 dark:border-surface-800 cursor-pointer hover:bg-slate-50 dark:hover:bg-surface-800/50 transition">
                            <input
                                v-model="phpForm.php_version"
                                type="radio"
                                value="none"
                                class="text-brand-600 focus:ring-brand-500"
                            />
                            <div>
                                <p class="text-xs font-semibold text-slate-900 dark:text-white">Static Site (No PHP)</p>
                                <p class="text-[11px] text-slate-500 dark:text-surface-400">Serves static HTML/JS/CSS without FastCGI backend</p>
                            </div>
                        </label>
                    </div>

                    <button
                        type="submit"
                        :disabled="phpForm.processing || phpForm.php_version === website.php_version"
                        class="px-5 py-2.5 rounded-xl bg-brand-600 hover:bg-brand-500 text-white text-xs font-semibold shadow-md shadow-brand-600/20 transition disabled:opacity-50"
                    >
                        {{ phpForm.processing ? 'Switching Runtime...' : 'Apply PHP Version' }}
                    </button>
                </form>
            </div>

            <!-- Tab 4: SSL / HTTPS -->
            <div v-if="activeTab === 'ssl'" class="rounded-2xl bg-white dark:bg-surface-900/60 border border-slate-200/80 dark:border-surface-800/80 p-6 shadow-sm dark:shadow-xl space-y-6">
                <div>
                    <h3 class="text-xs font-semibold text-slate-900 dark:text-white uppercase tracking-wider flex items-center gap-2">
                        <ShieldCheck class="w-4 h-4 text-emerald-600 dark:text-emerald-400" />
                        <span>SSL / TLS Certificate Management</span>
                    </h3>
                    <p class="text-xs text-slate-500 dark:text-surface-400 mt-1">
                        Secure your website traffic with automated Let's Encrypt SSL/TLS certificates over HTTP-01 challenge.
                    </p>
                </div>

                <div v-if="website.ssl_certificate" class="p-4 rounded-xl bg-emerald-50/70 dark:bg-emerald-500/10 border border-emerald-200 dark:border-emerald-500/20 space-y-3">
                    <div class="flex items-center justify-between">
                        <div class="flex items-center gap-2">
                            <ShieldCheck class="w-5 h-5 text-emerald-600 dark:text-emerald-400" />
                            <span class="text-xs font-bold text-emerald-900 dark:text-emerald-200">Active Certificate Issued</span>
                        </div>
                        <span class="px-2.5 py-0.5 rounded-full text-[10px] font-semibold bg-emerald-100 text-emerald-700 dark:bg-emerald-500/20 dark:text-emerald-300">
                            {{ website.ssl_certificate.status }}
                        </span>
                    </div>

                    <div class="grid grid-cols-1 sm:grid-cols-3 gap-3 text-xs">
                        <div>
                            <p class="text-slate-500 dark:text-surface-400 font-medium">Issuer</p>
                            <p class="font-mono font-bold text-slate-900 dark:text-white">{{ website.ssl_certificate.issuer }}</p>
                        </div>
                        <div>
                            <p class="text-slate-500 dark:text-surface-400 font-medium">Valid From</p>
                            <p class="font-mono font-bold text-slate-900 dark:text-white">{{ website.ssl_certificate.valid_from }}</p>
                        </div>
                        <div>
                            <p class="text-slate-500 dark:text-surface-400 font-medium">Valid Until</p>
                            <p class="font-mono font-bold text-slate-900 dark:text-white">{{ website.ssl_certificate.valid_until }}</p>
                        </div>
                    </div>
                </div>

                <!-- Issue New SSL Form -->
                <div class="space-y-4 max-w-md pt-2 border-t border-slate-200/80 dark:border-surface-800">
                    <h4 class="text-xs font-semibold text-slate-900 dark:text-white">Issue / Re-issue Let's Encrypt Certificate</h4>

                    <form @submit.prevent="issueSsl" class="space-y-4">
                        <div>
                            <label class="block text-xs font-medium text-slate-700 dark:text-surface-300 mb-1.5">Notification Email</label>
                            <input
                                v-model="sslForm.email"
                                type="email"
                                placeholder="admin@example.com"
                                class="w-full px-3.5 py-2.5 rounded-xl bg-slate-50 dark:bg-surface-950/60 border border-slate-200 dark:border-surface-800 text-xs text-slate-900 dark:text-white focus:ring-2 focus:ring-brand-500/40 focus:outline-none transition"
                            />
                        </div>

                        <label class="flex items-center gap-2.5 cursor-pointer">
                            <input
                                v-model="sslForm.force_https"
                                type="checkbox"
                                class="rounded border-slate-300 dark:border-surface-700 bg-white dark:bg-surface-900 text-brand-600 focus:ring-brand-500"
                            />
                            <span class="text-xs text-slate-700 dark:text-surface-300">Automatically redirect all HTTP traffic to HTTPS (301)</span>
                        </label>

                        <button
                            type="submit"
                            :disabled="sslForm.processing"
                            class="px-5 py-2.5 rounded-xl bg-brand-600 hover:bg-brand-500 text-white text-xs font-semibold shadow-md shadow-brand-600/20 flex items-center gap-2 transition disabled:opacity-50"
                        >
                            <ShieldCheck class="w-4 h-4" />
                            <span>{{ sslForm.processing ? 'Requesting Let\'s Encrypt...' : 'Issue SSL Certificate' }}</span>
                        </button>
                    </form>
                </div>
            </div>

            <!-- Tab 5: Nginx Logs -->
            <div v-if="activeTab === 'logs'" class="rounded-2xl bg-white dark:bg-surface-900/60 border border-slate-200/80 dark:border-surface-800/80 p-6 shadow-sm dark:shadow-xl space-y-4">
                <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3">
                    <div class="flex items-center gap-2">
                        <button
                            @click="{ logType = 'access'; fetchLogs(); }"
                            :class="[
                                'px-3 py-1.5 rounded-lg text-xs font-medium transition',
                                logType === 'access' ? 'bg-slate-200 dark:bg-surface-800 text-slate-900 dark:text-white font-semibold' : 'text-slate-500 hover:text-slate-900 dark:text-surface-400 dark:hover:text-white'
                            ]"
                        >
                            Access Log
                        </button>
                        <button
                            @click="{ logType = 'error'; fetchLogs(); }"
                            :class="[
                                'px-3 py-1.5 rounded-lg text-xs font-medium transition',
                                logType === 'error' ? 'bg-slate-200 dark:bg-surface-800 text-slate-900 dark:text-white font-semibold' : 'text-slate-500 hover:text-slate-900 dark:text-surface-400 dark:hover:text-white'
                            ]"
                        >
                            Error Log
                        </button>
                    </div>

                    <div class="flex items-center gap-2">
                        <select
                            v-model="logLines"
                            @change="fetchLogs"
                            class="px-2.5 py-1.5 rounded-lg bg-slate-50 dark:bg-surface-950 border border-slate-200 dark:border-surface-800 text-xs text-slate-700 dark:text-surface-300"
                        >
                            <option :value="50">Last 50 lines</option>
                            <option :value="100">Last 100 lines</option>
                            <option :value="250">Last 250 lines</option>
                        </select>
                        <button
                            @click="fetchLogs"
                            :disabled="isFetchingLogs"
                            class="p-2 rounded-lg bg-slate-100 hover:bg-slate-200 dark:bg-surface-800 dark:hover:bg-surface-700 text-slate-600 dark:text-surface-300 transition"
                            title="Refresh Logs"
                        >
                            <RefreshCw :class="['w-3.5 h-3.5', isFetchingLogs ? 'animate-spin' : '']" />
                        </button>
                    </div>
                </div>

                <!-- Log Output Terminal -->
                <div class="p-4 rounded-xl bg-slate-950 border border-slate-800 font-mono text-[11px] text-slate-300 h-96 overflow-y-auto space-y-1 select-text">
                    <div v-if="logs.length === 0" class="text-slate-500 italic">
                        No log entries available.
                    </div>
                    <div v-for="(line, idx) in logs" :key="idx" class="leading-relaxed hover:bg-slate-900 px-1 py-0.5 rounded">
                        {{ line }}
                    </div>
                </div>
            </div>

            <!-- Tab 6: Deployments -->
            <div v-if="activeTab === 'deployments'" class="rounded-2xl bg-white dark:bg-surface-900/60 border border-slate-200/80 dark:border-surface-800/80 p-6 shadow-sm dark:shadow-xl space-y-4">
                <div class="flex items-center justify-between">
                    <h3 class="text-xs font-semibold text-slate-700 dark:text-surface-300 uppercase tracking-wider flex items-center gap-2">
                        <GitBranch class="w-4 h-4 text-purple-500 dark:text-purple-400" />
                        <span>Deployment Pipeline History</span>
                    </h3>
                    <button
                        @click="triggerDeploy"
                        class="px-3 py-1.5 rounded-lg bg-brand-600 hover:bg-brand-500 text-white text-xs font-medium transition"
                    >
                        Trigger Build
                    </button>
                </div>

                <div v-if="website.deployments.length === 0" class="text-center py-12 text-xs text-slate-400 dark:text-surface-500">
                    No deployments executed for this website yet. Click "Deploy Now" to trigger a Git pull and build hook execution.
                </div>

                <div v-else class="divide-y divide-slate-100 dark:divide-surface-800/60">
                    <div
                        v-for="deploy in website.deployments"
                        :key="deploy.id"
                        class="py-3.5 flex items-center justify-between text-xs"
                    >
                        <div class="flex items-center gap-3">
                            <span
                                :class="[
                                    'w-2 h-2 rounded-full',
                                    deploy.status === 'success' ? 'bg-emerald-500' : deploy.status === 'running' ? 'bg-amber-400 animate-ping' : 'bg-rose-500'
                                ]"
                            ></span>
                            <div>
                                <p class="text-slate-900 dark:text-surface-100 font-mono font-medium">
                                    Deployment #{{ deploy.id }} • {{ deploy.branch }}
                                </p>
                                <p class="text-[11px] text-slate-500 dark:text-surface-400">
                                    Trigger: {{ deploy.trigger_source }} • {{ deploy.duration_seconds ?? 0 }}s
                                </p>
                            </div>
                        </div>
                        <Link
                            :href="`/deployments/${deploy.id}`"
                            class="px-3 py-1.5 rounded-lg bg-slate-100 hover:bg-slate-200 dark:bg-surface-800 dark:hover:bg-surface-700 text-slate-700 dark:text-surface-200 text-xs font-mono transition"
                        >
                            View Terminal Logs
                        </Link>
                    </div>
                </div>
            </div>
        </div>
    </AppLayout>
</template>
