<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue';
import { Link } from '@inertiajs/vue3';
import AppLayout from '@/Layouts/AppLayout.vue';
import {
    Cpu,
    HardDrive,
    Database as DbIcon,
    Globe,
    Activity,
    Server,
    Clock,
    CheckCircle2,
    RefreshCw,
    Copy,
    Check,
    ChevronRight,
    Rocket,
    Layers,
    Sliders,
    Shield,
    Sparkles,
} from 'lucide-vue-next';
import type { SystemInfo, SystemMetrics, ServiceStatus, ActivityLog } from '@/types';

const props = defineProps<{
    systemInfo: SystemInfo;
    systemMetrics: SystemMetrics;
    services: ServiceStatus[];
    counts: {
        websites: number;
        databases: number;
        deployments: number;
    };
    recentActivity: ActivityLog[];
}>();

const liveMetrics = ref<SystemMetrics>(props.systemMetrics);
const isRefreshing = ref(false);
const copiedIp = ref(false);
let pollingInterval: any = null;

async function fetchLiveMetrics() {
    isRefreshing.value = true;
    try {
        const res = await fetch('/dashboard/metrics');
        const json = await res.json();
        if (json.success && json.data) {
            liveMetrics.value = json.data;
        }
    } catch (e) {
        console.error('Failed to fetch metrics', e);
    } finally {
        isRefreshing.value = false;
    }
}

function copyIp() {
    navigator.clipboard.writeText(props.systemInfo.public_ip);
    copiedIp.value = true;
    setTimeout(() => (copiedIp.value = false), 2000);
}

onMounted(() => {
    pollingInterval = setInterval(fetchLiveMetrics, 5000);
});

onUnmounted(() => {
    if (pollingInterval) clearInterval(pollingInterval);
});

function formatBytes(bytes: number): string {
    if (!bytes || bytes === 0) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i];
}

function formatUptime(seconds: number): string {
    if (!seconds) return '0m';
    const days = Math.floor(seconds / 86400);
    const hours = Math.floor((seconds % 86400) / 3600);
    const minutes = Math.floor((seconds % 3600) / 60);

    if (days > 0) return `${days}d ${hours}h ${minutes}m`;
    if (hours > 0) return `${hours}h ${minutes}m`;
    return `${minutes}m`;
}

function formatTime(dateStr: string): string {
    try {
        const d = new Date(dateStr);
        return d.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit', second: '2-digit', hour12: true });
    } catch {
        return dateStr;
    }
}
</script>

<template>
    <AppLayout title="Server Dashboard">
        <div class="space-y-6 max-w-7xl mx-auto">
            <!-- Server Header Banner -->
            <div
                class="rounded-2xl bg-white dark:bg-gradient-to-r dark:from-[#1e1b4b]/80 dark:via-[#1e293b]/70 dark:to-[#0f172a]/90 border border-slate-200/80 dark:border-indigo-500/20 p-6 shadow-sm dark:shadow-xl transition-all duration-200"
            >
                <div class="flex flex-col md:flex-row md:items-center justify-between gap-5">
                    <!-- Left: Hostname & OS -->
                    <div class="flex items-center gap-4">
                        <div class="w-12 h-12 rounded-2xl bg-gradient-to-tr from-brand-600 via-indigo-500 to-indigo-400 text-white flex items-center justify-center shadow-lg shadow-brand-500/25 ring-1 ring-white/20 shrink-0">
                            <Layers class="w-6 h-6" />
                        </div>
                        <div>
                            <div class="flex items-center gap-2.5">
                                <h2 class="text-lg font-bold text-slate-900 dark:text-white tracking-tight">{{ systemInfo.hostname }}</h2>
                                <span class="px-2.5 py-0.5 rounded-full text-xs font-semibold bg-emerald-50 text-emerald-600 border border-emerald-200 dark:bg-emerald-500/10 dark:border-emerald-500/30 dark:text-emerald-400 flex items-center gap-1.5">
                                    <span class="w-1.5 h-1.5 rounded-full bg-emerald-500 animate-pulse"></span>
                                    Online
                                </span>
                            </div>
                            <p class="text-xs text-slate-500 dark:text-surface-400 mt-1 font-mono">
                                {{ systemInfo.os }} • {{ systemInfo.kernel }} ({{ systemInfo.architecture }})
                            </p>
                        </div>
                    </div>

                    <!-- Right: IP, Uptime, Refresh -->
                    <div class="flex flex-wrap items-center gap-4 sm:gap-7 border-t md:border-t-0 pt-3 md:pt-0 border-slate-100 dark:border-surface-800">
                        <!-- Public IP with copy -->
                        <div>
                            <p class="text-[11px] font-medium text-slate-400 dark:text-surface-400">Public IP</p>
                            <button
                                @click="copyIp"
                                class="group flex items-center gap-1.5 font-mono text-xs font-semibold text-slate-800 dark:text-surface-200 hover:text-brand-600 dark:hover:text-brand-400 transition"
                                title="Click to copy IP"
                            >
                                <span>{{ systemInfo.public_ip }}</span>
                                <Check v-if="copiedIp" class="w-3.5 h-3.5 text-emerald-500" />
                                <Copy v-else class="w-3.5 h-3.5 text-slate-400 group-hover:text-brand-500 opacity-60 group-hover:opacity-100 transition" />
                            </button>
                        </div>

                        <!-- Uptime -->
                        <div>
                            <p class="text-[11px] font-medium text-slate-400 dark:text-surface-400">Uptime</p>
                            <p class="font-mono text-xs font-semibold text-slate-800 dark:text-surface-200">{{ formatUptime(systemInfo.uptime_seconds) }}</p>
                        </div>

                        <!-- Refresh telemetry button -->
                        <button
                            @click="fetchLiveMetrics"
                            :disabled="isRefreshing"
                            class="p-2.5 rounded-xl bg-slate-100 hover:bg-slate-200/80 dark:bg-surface-800/80 dark:hover:bg-surface-700 text-slate-600 dark:text-surface-300 border border-slate-200/60 dark:border-surface-700/60 transition shadow-sm"
                            title="Refresh telemetry"
                        >
                            <RefreshCw :class="['w-4 h-4 text-brand-600 dark:text-brand-400', isRefreshing ? 'animate-spin' : '']" />
                        </button>
                    </div>
                </div>
            </div>

            <!-- Real-time Hardware Telemetry Gauges (4 Grid) -->
            <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
                <!-- CPU Usage -->
                <div class="rounded-2xl bg-white dark:bg-surface-900/70 border border-slate-200/80 dark:border-surface-800/80 p-5 shadow-sm dark:shadow-xl relative overflow-hidden transition-colors">
                    <div class="flex items-center justify-between">
                        <div class="w-10 h-10 rounded-xl bg-indigo-50 text-indigo-600 dark:bg-indigo-500/10 dark:text-indigo-400 border border-indigo-100 dark:border-indigo-500/20 flex items-center justify-center shadow-sm">
                            <Cpu class="w-5 h-5" />
                        </div>
                        <span class="text-xs font-semibold text-slate-500 dark:text-surface-400">CPU Usage</span>
                    </div>
                    <div class="mt-4 flex items-baseline justify-between">
                        <span class="text-2xl font-extrabold text-slate-900 dark:text-white tracking-tight font-mono">
                            {{ liveMetrics.cpu.usage_percent }}%
                        </span>
                        <span class="text-xs text-slate-500 dark:text-surface-400 font-mono">{{ liveMetrics.cpu.cores }} Cores</span>
                    </div>
                    <!-- Bottom Progress Bar -->
                    <div class="mt-3.5 h-1.5 w-full bg-slate-100 dark:bg-surface-800 rounded-full overflow-hidden">
                        <div
                            class="h-full bg-indigo-500 dark:bg-brand-500 rounded-full transition-all duration-500"
                            :style="{ width: `${Math.min(liveMetrics.cpu.usage_percent, 100)}%` }"
                        ></div>
                    </div>
                </div>

                <!-- Memory Usage -->
                <div class="rounded-2xl bg-white dark:bg-surface-900/70 border border-slate-200/80 dark:border-surface-800/80 p-5 shadow-sm dark:shadow-xl relative overflow-hidden transition-colors">
                    <div class="flex items-center justify-between">
                        <div class="w-10 h-10 rounded-xl bg-emerald-50 text-emerald-600 dark:bg-emerald-500/10 dark:text-emerald-400 border border-emerald-100 dark:border-emerald-500/20 flex items-center justify-center shadow-sm">
                            <Sliders class="w-5 h-5" />
                        </div>
                        <span class="text-xs font-semibold text-slate-500 dark:text-surface-400">Memory (RAM)</span>
                    </div>
                    <div class="mt-4 flex items-baseline justify-between">
                        <span class="text-2xl font-extrabold text-slate-900 dark:text-white tracking-tight font-mono">
                            {{ liveMetrics.memory.usage_percent }}%
                        </span>
                        <span class="text-xs text-slate-500 dark:text-surface-400 font-mono">
                            {{ formatBytes(liveMetrics.memory.used_bytes) }} / {{ formatBytes(liveMetrics.memory.total_bytes) }}
                        </span>
                    </div>
                    <div class="mt-3.5 h-1.5 w-full bg-slate-100 dark:bg-surface-800 rounded-full overflow-hidden">
                        <div
                            class="h-full bg-emerald-500 rounded-full transition-all duration-500"
                            :style="{ width: `${Math.min(liveMetrics.memory.usage_percent, 100)}%` }"
                        ></div>
                    </div>
                </div>

                <!-- Disk Usage -->
                <div class="rounded-2xl bg-white dark:bg-surface-900/70 border border-slate-200/80 dark:border-surface-800/80 p-5 shadow-sm dark:shadow-xl relative overflow-hidden transition-colors">
                    <div class="flex items-center justify-between">
                        <div class="w-10 h-10 rounded-xl bg-sky-50 text-sky-600 dark:bg-sky-500/10 dark:text-sky-400 border border-sky-100 dark:border-sky-500/20 flex items-center justify-center shadow-sm">
                            <HardDrive class="w-5 h-5" />
                        </div>
                        <span class="text-xs font-semibold text-slate-500 dark:text-surface-400">Disk Space</span>
                    </div>
                    <div class="mt-4 flex items-baseline justify-between">
                        <span class="text-2xl font-extrabold text-slate-900 dark:text-white tracking-tight font-mono">
                            {{ liveMetrics.disk.usage_percent }}%
                        </span>
                        <span class="text-xs text-slate-500 dark:text-surface-400 font-mono">
                            {{ formatBytes(liveMetrics.disk.used_bytes) }} / {{ formatBytes(liveMetrics.disk.total_bytes) }}
                        </span>
                    </div>
                    <div class="mt-3.5 h-1.5 w-full bg-slate-100 dark:bg-surface-800 rounded-full overflow-hidden">
                        <div
                            class="h-full bg-sky-500 rounded-full transition-all duration-500"
                            :style="{ width: `${Math.min(liveMetrics.disk.usage_percent, 100)}%` }"
                        ></div>
                    </div>
                </div>

                <!-- Load Average -->
                <div class="rounded-2xl bg-white dark:bg-surface-900/70 border border-slate-200/80 dark:border-surface-800/80 p-5 shadow-sm dark:shadow-xl relative overflow-hidden transition-colors">
                    <div class="flex items-center justify-between">
                        <div class="w-10 h-10 rounded-xl bg-amber-50 text-amber-600 dark:bg-amber-500/10 dark:text-amber-400 border border-amber-100 dark:border-amber-500/20 flex items-center justify-center shadow-sm">
                            <Clock class="w-5 h-5" />
                        </div>
                        <span class="text-xs font-semibold text-slate-500 dark:text-surface-400">Load Average</span>
                    </div>
                    <div class="mt-4 flex items-baseline justify-between">
                        <span class="text-2xl font-extrabold text-slate-900 dark:text-white tracking-tight font-mono">
                            {{ liveMetrics.load_average.load1 }}
                        </span>
                        <span class="text-xs text-slate-500 dark:text-surface-400 font-mono">
                            {{ liveMetrics.load_average.load5 }} • {{ liveMetrics.load_average.load15 }}
                        </span>
                    </div>
                    <p class="text-[11px] text-slate-400 dark:text-surface-500 mt-3 font-mono">
                        1m: {{ liveMetrics.load_average.load1 }}  5m: {{ liveMetrics.load_average.load5 }}  15m: {{ liveMetrics.load_average.load15 }}
                    </p>
                </div>
            </div>

            <!-- Quick Resource Counters & Services Grid -->
            <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
                <!-- Managed Resource Stats (1 Col) -->
                <div class="rounded-2xl bg-white dark:bg-surface-900/70 border border-slate-200/80 dark:border-surface-800/80 p-6 shadow-sm dark:shadow-xl space-y-4 transition-colors">
                    <div class="flex items-center gap-2">
                        <Server class="w-4 h-4 text-brand-600 dark:text-brand-400" />
                        <h3 class="text-sm font-semibold text-slate-900 dark:text-white">Managed Resources</h3>
                    </div>
                    <div class="grid grid-cols-3 gap-3">
                        <!-- Websites -->
                        <Link
                            href="/websites"
                            class="p-4 rounded-2xl bg-slate-50 hover:bg-slate-100 dark:bg-surface-950/60 dark:hover:bg-surface-800/60 border border-slate-200/80 dark:border-surface-800/80 hover:border-brand-500/40 transition group flex flex-col justify-between"
                        >
                            <Globe class="w-5 h-5 text-indigo-500 dark:text-brand-400 mb-2 group-hover:scale-110 transition" />
                            <div>
                                <p class="text-xl font-bold text-slate-900 dark:text-white font-mono">{{ counts.websites }}</p>
                                <p class="text-[11px] font-medium text-slate-500 dark:text-surface-400 mt-0.5">Websites</p>
                            </div>
                        </Link>

                        <!-- Databases -->
                        <Link
                            href="/databases"
                            class="p-4 rounded-2xl bg-slate-50 hover:bg-slate-100 dark:bg-surface-950/60 dark:hover:bg-surface-800/60 border border-slate-200/80 dark:border-surface-800/80 hover:border-emerald-500/40 transition group flex flex-col justify-between"
                        >
                            <DbIcon class="w-5 h-5 text-emerald-500 dark:text-emerald-400 mb-2 group-hover:scale-110 transition" />
                            <div>
                                <p class="text-xl font-bold text-slate-900 dark:text-white font-mono">{{ counts.databases }}</p>
                                <p class="text-[11px] font-medium text-slate-500 dark:text-surface-400 mt-0.5">Databases</p>
                            </div>
                        </Link>

                        <!-- Deployments -->
                        <Link
                            href="/deployments"
                            class="p-4 rounded-2xl bg-slate-50 hover:bg-slate-100 dark:bg-surface-950/60 dark:hover:bg-surface-800/60 border border-slate-200/80 dark:border-surface-800/80 hover:border-purple-500/40 transition group flex flex-col justify-between"
                        >
                            <Rocket class="w-5 h-5 text-purple-500 dark:text-purple-400 mb-2 group-hover:scale-110 transition" />
                            <div>
                                <p class="text-xl font-bold text-slate-900 dark:text-white font-mono">{{ counts.deployments }}</p>
                                <p class="text-[11px] font-medium text-slate-500 dark:text-surface-400 mt-0.5">Deployments</p>
                            </div>
                        </Link>
                    </div>
                </div>

                <!-- Core System Services Status (2 Cols) -->
                <div class="lg:col-span-2 rounded-2xl bg-white dark:bg-surface-900/70 border border-slate-200/80 dark:border-surface-800/80 p-6 shadow-sm dark:shadow-xl transition-colors">
                    <div class="flex items-center justify-between mb-4">
                        <div class="flex items-center gap-2">
                            <Cpu class="w-4 h-4 text-brand-600 dark:text-brand-400" />
                            <h3 class="text-sm font-semibold text-slate-900 dark:text-white">System Services</h3>
                        </div>
                        <Link href="/services" class="text-xs text-brand-600 dark:text-brand-400 hover:underline font-semibold flex items-center gap-1">
                            <span>Manage Services</span>
                            <ChevronRight class="w-3.5 h-3.5" />
                        </Link>
                    </div>
                    <div class="grid grid-cols-2 sm:grid-cols-3 gap-3">
                        <div
                            v-for="svc in services"
                            :key="svc.name"
                            class="p-3.5 rounded-2xl bg-slate-50 dark:bg-surface-950/60 border border-slate-200/80 dark:border-surface-800/80 flex items-center justify-between transition hover:border-slate-300 dark:hover:border-surface-700"
                        >
                            <div class="min-w-0 pr-2">
                                <p class="text-xs font-semibold text-slate-800 dark:text-surface-200 font-mono truncate">{{ svc.name }}</p>
                                <span class="text-[10px] text-slate-400 dark:text-surface-500">systemd</span>
                            </div>
                            <span class="flex items-center gap-1.5 text-xs font-medium text-emerald-600 dark:text-emerald-400 shrink-0">
                                <span class="w-2 h-2 rounded-full bg-emerald-500 animate-pulse"></span>
                                Active
                            </span>
                        </div>
                    </div>
                </div>
            </div>

            <!-- Recent Activity Audit Feed (Full Width) -->
            <div class="rounded-2xl bg-white dark:bg-surface-900/70 border border-slate-200/80 dark:border-surface-800/80 p-6 shadow-sm dark:shadow-xl transition-colors">
                <div class="flex items-center justify-between mb-4">
                    <div class="flex items-center gap-2">
                        <Activity class="w-4 h-4 text-brand-600 dark:text-brand-400" />
                        <h3 class="text-sm font-semibold text-slate-900 dark:text-white">Recent Activity Ledger</h3>
                    </div>
                    <Link href="/activity" class="text-xs text-brand-600 dark:text-brand-400 hover:underline font-semibold flex items-center gap-1">
                        <span>View All Audit Logs</span>
                        <ChevronRight class="w-3.5 h-3.5" />
                    </Link>
                </div>
                <div v-if="recentActivity.length === 0" class="text-center py-8 text-xs text-slate-400 dark:text-surface-500">
                    No recent administrative activities logged yet.
                </div>
                <div v-else class="divide-y divide-slate-100 dark:divide-surface-800/60">
                    <div
                        v-for="log in recentActivity"
                        :key="log.id"
                        class="py-3.5 flex items-center justify-between text-xs"
                    >
                        <div class="flex items-center gap-3.5 min-w-0 pr-4">
                            <span class="px-2 py-0.5 rounded-md text-[10px] font-mono font-semibold bg-brand-50 text-brand-600 border border-brand-200/70 dark:bg-brand-500/20 dark:text-brand-300 dark:border-brand-500/30 shrink-0">
                                {{ log.action.split('.')[0] || 'audit' }}
                            </span>
                            <div class="min-w-0">
                                <p class="text-slate-900 dark:text-surface-200 font-semibold font-mono truncate">{{ log.action }}</p>
                                <p class="text-[11px] text-slate-500 dark:text-surface-400 truncate">
                                    {{ log.user_email || 'system' }} • {{ log.ip_address }}
                                </p>
                            </div>
                        </div>
                        <span class="text-[11px] text-slate-400 dark:text-surface-400 font-mono shrink-0">
                            {{ formatTime(log.created_at) }}
                        </span>
                    </div>
                </div>
            </div>
        </div>
    </AppLayout>
</template>
