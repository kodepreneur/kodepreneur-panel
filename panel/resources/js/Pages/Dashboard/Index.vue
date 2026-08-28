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
    Wifi,
    CheckCircle2,
    XCircle,
    ArrowUpRight,
    RefreshCw,
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
</script>

<template>
    <AppLayout title="Server Dashboard">
        <div class="space-y-6 max-w-7xl mx-auto">
            <!-- Server Header Banner -->
            <div class="rounded-2xl bg-gradient-to-r from-surface-900 via-surface-900/90 to-brand-950/40 border border-surface-800/80 p-5 sm:p-6 shadow-xl">
                <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
                    <div class="flex items-center gap-4">
                        <div class="w-12 h-12 rounded-xl bg-brand-600/20 border border-brand-500/30 text-brand-400 flex items-center justify-center shrink-0">
                            <Server class="w-6 h-6" />
                        </div>
                        <div>
                            <div class="flex items-center gap-2">
                                <h2 class="text-base font-semibold text-white tracking-tight">{{ systemInfo.hostname }}</h2>
                                <span class="px-2 py-0.5 rounded-full text-[10px] font-medium bg-emerald-500/10 border border-emerald-500/30 text-emerald-400">
                                    Online
                                </span>
                            </div>
                            <p class="text-xs text-surface-400 mt-0.5 font-mono">
                                {{ systemInfo.os }} • {{ systemInfo.kernel }} ({{ systemInfo.architecture }})
                            </p>
                        </div>
                    </div>

                    <div class="flex items-center gap-4 sm:gap-6 text-xs text-surface-400 border-t sm:border-t-0 pt-3 sm:pt-0 border-surface-800">
                        <div>
                            <p class="text-[11px] text-surface-500">Public IP</p>
                            <p class="font-mono text-surface-200 font-medium">{{ systemInfo.public_ip }}</p>
                        </div>
                        <div>
                            <p class="text-[11px] text-surface-500">Uptime</p>
                            <p class="font-mono text-surface-200 font-medium">{{ formatUptime(systemInfo.uptime_seconds) }}</p>
                        </div>
                        <button
                            @click="fetchLiveMetrics"
                            :disabled="isRefreshing"
                            class="p-2 rounded-lg bg-surface-800 hover:bg-surface-700 text-surface-300 transition"
                            title="Refresh telemetry"
                        >
                            <RefreshCw :class="['w-4 h-4', isRefreshing ? 'animate-spin text-brand-400' : '']" />
                        </button>
                    </div>
                </div>
            </div>

            <!-- Real-time Hardware Telemetry Gauges -->
            <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
                <!-- CPU Usage -->
                <div class="rounded-2xl bg-surface-900/60 border border-surface-800/80 p-5 shadow-lg relative overflow-hidden">
                    <div class="flex items-center justify-between">
                        <span class="text-xs font-medium text-surface-400">CPU Usage</span>
                        <Cpu class="w-4 h-4 text-brand-400" />
                    </div>
                    <div class="mt-3 flex items-baseline justify-between">
                        <span class="text-2xl font-bold text-white tracking-tight font-mono">
                            {{ liveMetrics.cpu.usage_percent }}%
                        </span>
                        <span class="text-xs text-surface-500 font-mono">{{ liveMetrics.cpu.cores }} Cores</span>
                    </div>
                    <!-- Progress Bar -->
                    <div class="mt-3 h-1.5 w-full bg-surface-800 rounded-full overflow-hidden">
                        <div
                            class="h-full bg-brand-500 rounded-full transition-all duration-500"
                            :style="{ width: `${Math.min(liveMetrics.cpu.usage_percent, 100)}%` }"
                        ></div>
                    </div>
                </div>

                <!-- Memory Usage -->
                <div class="rounded-2xl bg-surface-900/60 border border-surface-800/80 p-5 shadow-lg relative overflow-hidden">
                    <div class="flex items-center justify-between">
                        <span class="text-xs font-medium text-surface-400">Memory (RAM)</span>
                        <Activity class="w-4 h-4 text-emerald-400" />
                    </div>
                    <div class="mt-3 flex items-baseline justify-between">
                        <span class="text-2xl font-bold text-white tracking-tight font-mono">
                            {{ liveMetrics.memory.usage_percent }}%
                        </span>
                        <span class="text-xs text-surface-500 font-mono">
                            {{ formatBytes(liveMetrics.memory.used_bytes) }} / {{ formatBytes(liveMetrics.memory.total_bytes) }}
                        </span>
                    </div>
                    <div class="mt-3 h-1.5 w-full bg-surface-800 rounded-full overflow-hidden">
                        <div
                            class="h-full bg-emerald-500 rounded-full transition-all duration-500"
                            :style="{ width: `${Math.min(liveMetrics.memory.usage_percent, 100)}%` }"
                        ></div>
                    </div>
                </div>

                <!-- Disk Usage -->
                <div class="rounded-2xl bg-surface-900/60 border border-surface-800/80 p-5 shadow-lg relative overflow-hidden">
                    <div class="flex items-center justify-between">
                        <span class="text-xs font-medium text-surface-400">Disk Space</span>
                        <HardDrive class="w-4 h-4 text-cyan-400" />
                    </div>
                    <div class="mt-3 flex items-baseline justify-between">
                        <span class="text-2xl font-bold text-white tracking-tight font-mono">
                            {{ liveMetrics.disk.usage_percent }}%
                        </span>
                        <span class="text-xs text-surface-500 font-mono">
                            {{ formatBytes(liveMetrics.disk.used_bytes) }} / {{ formatBytes(liveMetrics.disk.total_bytes) }}
                        </span>
                    </div>
                    <div class="mt-3 h-1.5 w-full bg-surface-800 rounded-full overflow-hidden">
                        <div
                            class="h-full bg-cyan-500 rounded-full transition-all duration-500"
                            :style="{ width: `${Math.min(liveMetrics.disk.usage_percent, 100)}%` }"
                        ></div>
                    </div>
                </div>

                <!-- Load Average -->
                <div class="rounded-2xl bg-surface-900/60 border border-surface-800/80 p-5 shadow-lg relative overflow-hidden">
                    <div class="flex items-center justify-between">
                        <span class="text-xs font-medium text-surface-400">Load Average</span>
                        <Clock class="w-4 h-4 text-amber-400" />
                    </div>
                    <div class="mt-3 flex items-baseline justify-between">
                        <span class="text-2xl font-bold text-white tracking-tight font-mono">
                            {{ liveMetrics.load_average.load1 }}
                        </span>
                        <span class="text-xs text-surface-400 font-mono">
                            {{ liveMetrics.load_average.load5 }} • {{ liveMetrics.load_average.load15 }}
                        </span>
                    </div>
                    <p class="text-[11px] text-surface-500 mt-3 font-mono">1m, 5m, 15m metrics</p>
                </div>
            </div>

            <!-- Quick Resource Counters & Services Grid -->
            <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
                <!-- Managed Resource Stats -->
                <div class="rounded-2xl bg-surface-900/60 border border-surface-800/80 p-5 shadow-lg space-y-4">
                    <h3 class="text-xs font-semibold text-surface-300 uppercase tracking-wider">Managed Resources</h3>
                    <div class="grid grid-cols-3 gap-3">
                        <Link
                            href="/websites"
                            class="p-4 rounded-xl bg-surface-950/60 border border-surface-800/80 hover:border-brand-500/40 transition group"
                        >
                            <Globe class="w-5 h-5 text-brand-400 mb-2 group-hover:scale-110 transition" />
                            <p class="text-xl font-bold text-white font-mono">{{ counts.websites }}</p>
                            <p class="text-[11px] text-surface-400">Websites</p>
                        </Link>

                        <Link
                            href="/databases"
                            class="p-4 rounded-xl bg-surface-950/60 border border-surface-800/80 hover:border-emerald-500/40 transition group"
                        >
                            <DbIcon class="w-5 h-5 text-emerald-400 mb-2 group-hover:scale-110 transition" />
                            <p class="text-xl font-bold text-white font-mono">{{ counts.databases }}</p>
                            <p class="text-[11px] text-surface-400">Databases</p>
                        </Link>

                        <Link
                            href="/deployments"
                            class="p-4 rounded-xl bg-surface-950/60 border border-surface-800/80 hover:border-purple-500/40 transition group"
                        >
                            <ArrowUpRight class="w-5 h-5 text-purple-400 mb-2 group-hover:scale-110 transition" />
                            <p class="text-xl font-bold text-white font-mono">{{ counts.deployments }}</p>
                            <p class="text-[11px] text-surface-400">Deployments</p>
                        </Link>
                    </div>
                </div>

                <!-- Core System Services Status -->
                <div class="lg:col-span-2 rounded-2xl bg-surface-900/60 border border-surface-800/80 p-5 shadow-lg">
                    <div class="flex items-center justify-between mb-4">
                        <h3 class="text-xs font-semibold text-surface-300 uppercase tracking-wider">System Services</h3>
                        <Link href="/services" class="text-xs text-brand-400 hover:text-brand-300 font-medium">Manage Services →</Link>
                    </div>
                    <div class="grid grid-cols-2 sm:grid-cols-3 gap-3">
                        <div
                            v-for="svc in services"
                            :key="svc.name"
                            class="p-3 rounded-xl bg-surface-950/60 border border-surface-800 flex items-center justify-between"
                        >
                            <div>
                                <p class="text-xs font-medium text-surface-200 font-mono">{{ svc.name }}</p>
                                <span class="text-[10px] text-surface-500">systemd</span>
                            </div>
                            <span class="flex items-center gap-1.5 text-[11px] font-medium text-emerald-400">
                                <span class="w-1.5 h-1.5 rounded-full bg-emerald-500"></span>
                                Active
                            </span>
                        </div>
                    </div>
                </div>
            </div>

            <!-- Recent Activity Audit Feed -->
            <div class="rounded-2xl bg-surface-900/60 border border-surface-800/80 p-5 shadow-lg">
                <div class="flex items-center justify-between mb-4">
                    <h3 class="text-xs font-semibold text-surface-300 uppercase tracking-wider">Recent Activity Ledger</h3>
                    <Link href="/activity" class="text-xs text-brand-400 hover:text-brand-300 font-medium">View All Audit Logs →</Link>
                </div>
                <div v-if="recentActivity.length === 0" class="text-center py-6 text-xs text-surface-500">
                    No recent administrative activities logged yet.
                </div>
                <div v-else class="divide-y divide-surface-800/60">
                    <div
                        v-for="log in recentActivity"
                        :key="log.id"
                        class="py-3 flex items-center justify-between text-xs"
                    >
                        <div class="flex items-center gap-3">
                            <div class="w-7 h-7 rounded-lg bg-surface-800 text-surface-300 flex items-center justify-center font-mono text-[11px]">
                                {{ log.action.split('.')[0] }}
                            </div>
                            <div>
                                <p class="text-surface-200 font-medium font-mono">{{ log.action }}</p>
                                <p class="text-[11px] text-surface-400">{{ log.user_email || 'system' }} • {{ log.ip_address }}</p>
                            </div>
                        </div>
                        <span class="text-[11px] text-surface-500 font-mono">
                            {{ new Date(log.created_at).toLocaleTimeString() }}
                        </span>
                    </div>
                </div>
            </div>
        </div>
    </AppLayout>
</template>
