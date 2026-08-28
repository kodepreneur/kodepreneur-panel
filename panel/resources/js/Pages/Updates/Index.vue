<script setup lang="ts">
import { ref } from 'vue';
import AppLayout from '@/Layouts/AppLayout.vue';
import {
    ArrowUpCircle,
    RefreshCw,
    CheckCircle2,
    AlertCircle,
    GitCommit,
    GitBranch,
    Server,
    Cpu,
    Terminal,
    Copy,
    Check,
    ExternalLink,
    ShieldCheck,
    DownloadCloud,
    RotateCcw,
    Zap,
} from 'lucide-vue-next';
import type { UpdateInfo } from '@/types';

const props = defineProps<{
    updateInfo: UpdateInfo;
}>();

const info = ref<UpdateInfo>(props.updateInfo);
const isChecking = ref(false);
const isUpdating = ref(false);
const showConfirmModal = ref(false);
const showTerminalModal = ref(false);
const copiedCli = ref(false);
const copiedLogs = ref(false);

const updateLogs = ref<string>('');
const updateSuccess = ref<boolean | null>(null);
const updateDuration = ref<number>(0);
const updateError = ref<string>('');

const currentStep = ref(1);
const totalSteps = 6;
const stepLabels = [
    'Fetching latest release from GitHub',
    'Updating Composer dependencies',
    'Applying database migrations',
    'Compiling frontend assets via Vite',
    'Optimizing Laravel configuration & routes',
    'Recompiling Go Agent & reloading services',
];

const cliCommand = 'curl -fsSL https://raw.githubusercontent.com/kodepreneur/kodepreneur-panel/main/installer/update.sh | sudo bash';

function getCsrfToken(): string {
    const meta = document.querySelector('meta[name="csrf-token"]') as HTMLMetaElement;
    if (meta?.content) {
        return meta.content;
    }
    const match = document.cookie.match(/XSRF-TOKEN=([^;]+)/);
    if (match) {
        return decodeURIComponent(match[1]);
    }
    return '';
}

async function checkForUpdates() {
    if (isChecking.value) return;
    isChecking.value = true;
    try {
        const token = getCsrfToken();
        const res = await fetch('/updates/check', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'X-Requested-With': 'XMLHttpRequest',
                'X-CSRF-TOKEN': token,
                'X-XSRF-TOKEN': token,
            },
        });
        const data = await res.json();
        if (data.success && data.updateInfo) {
            info.value = data.updateInfo;
        }
    } catch (e) {
        console.error('Failed to check for updates:', e);
    } finally {
        isChecking.value = false;
    }
}

function startUpdate() {
    showConfirmModal.value = false;
    showTerminalModal.value = true;
    isUpdating.value = true;
    updateLogs.value = `[${new Date().toLocaleTimeString()}] 🚀 Initiating Kodepreneur Panel Update Pipeline...\n[${new Date().toLocaleTimeString()}] Connecting to root Go Agent daemon...\n`;
    updateSuccess.value = null;
    updateError.value = '';
    currentStep.value = 1;

    const timer = setInterval(() => {
        if (currentStep.value < totalSteps && isUpdating.value) {
            currentStep.value++;
        }
    }, 1200);

    const token = getCsrfToken();
    fetch('/updates/execute', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'X-Requested-With': 'XMLHttpRequest',
            'X-CSRF-TOKEN': token,
            'X-XSRF-TOKEN': token,
        },
        body: JSON.stringify({
            repository: info.value.repository,
            branch: info.value.branch,
        }),
    })
        .then(async (res) => {
            clearInterval(timer);
            currentStep.value = totalSteps;
            const json = await res.json();

            if (res.ok && json.success) {
                updateSuccess.value = true;
                updateLogs.value = json.data?.log_output || updateLogs.value + '\n🎉 Update finished successfully!';
                updateDuration.value = json.data?.duration_seconds || 1;
            } else {
                updateSuccess.value = false;
                updateError.value = json.message || 'Update failed during execution.';
                updateLogs.value += `\n[error] ${updateError.value}`;
            }
        })
        .catch((err) => {
            clearInterval(timer);
            updateSuccess.value = false;
            updateError.value = err.message || 'Network communication error';
            updateLogs.value += `\n[error] ${updateError.value}`;
        })
        .finally(() => {
            isUpdating.value = false;
        });
}

function copyCli() {
    navigator.clipboard.writeText(cliCommand);
    copiedCli.value = true;
    setTimeout(() => (copiedCli.value = false), 2000);
}

function copyLogs() {
    navigator.clipboard.writeText(updateLogs.value);
    copiedLogs.value = true;
    setTimeout(() => (copiedLogs.value = false), 2000);
}

function reloadPage() {
    window.location.reload();
}

function formatDate(dateStr?: string): string {
    if (!dateStr) return 'Just now';
    try {
        const d = new Date(dateStr);
        return d.toLocaleDateString(undefined, {
            month: 'short',
            day: 'numeric',
            year: 'numeric',
            hour: '2-digit',
            minute: '2-digit',
        });
    } catch {
        return dateStr;
    }
}
</script>

<template>
    <AppLayout title="System Updates">
        <div class="max-w-6xl mx-auto space-y-6">
            <!-- Top Header Bar -->
            <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
                <div>
                    <div class="flex items-center gap-2.5">
                        <div class="w-8 h-8 rounded-xl bg-brand-50 text-brand-600 dark:bg-brand-500/10 dark:text-brand-400 border border-brand-200 dark:border-brand-500/20 flex items-center justify-center">
                            <ArrowUpCircle class="w-4 h-4" />
                        </div>
                        <h2 class="text-base font-bold text-slate-900 dark:text-white tracking-tight">System & Panel Updates</h2>
                    </div>
                    <p class="text-xs text-slate-500 dark:text-surface-400 mt-1">
                        Manage software versioning, inspect GitHub releases, and perform zero-downtime control plane upgrades.
                    </p>
                </div>

                <!-- Check for Updates Button -->
                <div class="flex items-center gap-2.5">
                    <button
                        @click="checkForUpdates"
                        :disabled="isChecking || isUpdating"
                        class="px-3.5 py-2 rounded-xl bg-white dark:bg-surface-900 border border-slate-200/80 dark:border-surface-800 hover:bg-slate-50 dark:hover:bg-surface-800 text-slate-700 dark:text-surface-200 text-xs font-semibold flex items-center gap-2 transition disabled:opacity-50 shadow-sm"
                    >
                        <RefreshCw :class="['w-3.5 h-3.5 text-brand-600 dark:text-brand-400', isChecking ? 'animate-spin' : '']" />
                        <span>{{ isChecking ? 'Checking GitHub...' : 'Check for Updates' }}</span>
                    </button>

                    <button
                        @click="showConfirmModal = true"
                        :disabled="isUpdating"
                        class="px-4 py-2 rounded-xl bg-gradient-to-r from-brand-600 to-indigo-500 hover:from-brand-500 hover:to-indigo-400 text-white text-xs font-semibold flex items-center gap-1.5 shadow-md shadow-brand-500/20 transition disabled:opacity-50"
                    >
                        <DownloadCloud class="w-3.5 h-3.5" />
                        <span>Update Panel Now</span>
                    </button>
                </div>
            </div>

            <!-- Version Status Hero Card -->
            <div
                :class="[
                    'relative overflow-hidden rounded-2xl border p-6 sm:p-7 transition shadow-sm dark:shadow-xl',
                    info.has_update
                        ? 'bg-gradient-to-r from-amber-500/10 via-slate-50 to-white dark:from-brand-950/60 dark:via-surface-900/90 dark:to-surface-900/90 border-amber-300 dark:border-brand-500/40'
                        : 'bg-white dark:bg-surface-900/60 border-slate-200/80 dark:border-surface-800/80'
                ]"
            >
                <div class="relative z-10 flex flex-col md:flex-row md:items-center justify-between gap-6">
                    <div class="space-y-3">
                        <div class="flex flex-wrap items-center gap-2.5">
                            <span
                                :class="[
                                    'px-2.5 py-0.5 rounded-full text-xs font-mono font-semibold flex items-center gap-1.5 border',
                                    info.has_update
                                        ? 'bg-amber-50 text-amber-700 border-amber-200 dark:bg-amber-500/10 dark:border-amber-500/30 dark:text-amber-300 animate-pulse'
                                        : 'bg-emerald-50 text-emerald-700 border-emerald-200 dark:bg-emerald-500/10 dark:border-emerald-500/30 dark:text-emerald-400'
                                ]"
                            >
                                <span class="w-2 h-2 rounded-full" :class="info.has_update ? 'bg-amber-500' : 'bg-emerald-500'"></span>
                                {{ info.has_update ? 'Update Available' : 'Up to Date' }}
                            </span>
                            <span class="text-slate-300 dark:text-surface-600 text-xs">•</span>
                            <span class="text-xs text-slate-600 dark:text-surface-300 font-mono">Current: <span class="text-slate-900 dark:text-white font-bold">{{ info.current_version }}</span> ({{ info.current_commit }})</span>
                            <span class="text-slate-300 dark:text-surface-600 text-xs">•</span>
                            <span class="text-xs text-slate-500 dark:text-surface-400 font-mono">Branch: <span class="text-slate-800 dark:text-surface-200">{{ info.branch }}</span></span>
                        </div>

                        <div>
                            <h3 class="text-base font-bold text-slate-900 dark:text-white">
                                {{ info.has_update ? 'A newer version is ready on GitHub!' : 'Your Kodepreneur Panel is running the latest release' }}
                            </h3>
                            <p class="text-xs text-slate-500 dark:text-surface-400 mt-1 max-w-2xl leading-relaxed">
                                {{ info.has_update
                                    ? 'Upgrades include the latest control plane features, security patches, PHP runtime managers, and Go Agent daemon improvements.'
                                    : 'All panel components, systemd daemons, and database migrations are synchronized with the official repository.'
                                }}
                            </p>
                        </div>
                    </div>

                    <!-- Action Banner Button -->
                    <div class="flex items-center gap-3 shrink-0">
                        <button
                            @click="showConfirmModal = true"
                            :disabled="isUpdating"
                            :class="[
                                'px-5 py-2.5 rounded-xl font-semibold text-xs flex items-center gap-2 transition shadow-md',
                                info.has_update
                                    ? 'bg-brand-600 hover:bg-brand-500 text-white shadow-brand-500/25 ring-2 ring-brand-400/30'
                                    : 'bg-slate-100 hover:bg-slate-200 dark:bg-surface-800 dark:hover:bg-surface-700 text-slate-800 dark:text-surface-200 border border-slate-200/80 dark:border-surface-700'
                            ]"
                        >
                            <Zap class="w-3.5 h-3.5 text-amber-400" v-if="info.has_update" />
                            <RotateCcw class="w-3.5 h-3.5" v-else />
                            <span>{{ info.has_update ? 'Upgrade to Latest' : 'Reinstall / Sync Latest' }}</span>
                        </button>
                    </div>
                </div>
            </div>

            <!-- System Environment Metadata Grid -->
            <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
                <div class="p-4 rounded-2xl bg-white dark:bg-surface-900/60 border border-slate-200/80 dark:border-surface-800/80 flex items-center gap-3.5 shadow-sm">
                    <div class="w-10 h-10 rounded-xl bg-indigo-50 text-indigo-600 dark:bg-surface-800 dark:text-brand-400 border border-indigo-100 dark:border-surface-700/60 flex items-center justify-center shrink-0">
                        <Server class="w-5 h-5" />
                    </div>
                    <div class="min-w-0">
                        <p class="text-[11px] font-medium text-slate-400 dark:text-surface-400">Control Plane</p>
                        <p class="text-xs font-semibold text-slate-900 dark:text-white font-mono truncate">Laravel {{ info.laravel_version }}</p>
                    </div>
                </div>

                <div class="p-4 rounded-2xl bg-white dark:bg-surface-900/60 border border-slate-200/80 dark:border-surface-800/80 flex items-center gap-3.5 shadow-sm">
                    <div class="w-10 h-10 rounded-xl bg-purple-50 text-purple-600 dark:bg-surface-800 dark:text-purple-400 border border-purple-100 dark:border-surface-700/60 flex items-center justify-center shrink-0">
                        <Cpu class="w-5 h-5" />
                    </div>
                    <div class="min-w-0">
                        <p class="text-[11px] font-medium text-slate-400 dark:text-surface-400">PHP Runtime</p>
                        <p class="text-xs font-semibold text-slate-900 dark:text-white font-mono truncate">PHP {{ info.php_version }}</p>
                    </div>
                </div>

                <div class="p-4 rounded-2xl bg-white dark:bg-surface-900/60 border border-slate-200/80 dark:border-surface-800/80 flex items-center gap-3.5 shadow-sm">
                    <div class="w-10 h-10 rounded-xl bg-emerald-50 text-emerald-600 dark:bg-surface-800 dark:text-emerald-400 border border-emerald-100 dark:border-surface-700/60 flex items-center justify-center shrink-0">
                        <ShieldCheck class="w-5 h-5" />
                    </div>
                    <div class="min-w-0">
                        <p class="text-[11px] font-medium text-slate-400 dark:text-surface-400">Go Agent Daemon</p>
                        <div class="flex items-center gap-1.5">
                            <span class="w-2 h-2 rounded-full bg-emerald-500 animate-pulse"></span>
                            <span class="text-xs font-semibold text-emerald-600 dark:text-emerald-400 font-mono uppercase">{{ info.agent_status }}</span>
                        </div>
                    </div>
                </div>

                <div class="p-4 rounded-2xl bg-white dark:bg-surface-900/60 border border-slate-200/80 dark:border-surface-800/80 flex items-center gap-3.5 shadow-sm">
                    <div class="w-10 h-10 rounded-xl bg-sky-50 text-sky-600 dark:bg-surface-800 dark:text-cyan-400 border border-sky-100 dark:border-surface-700/60 flex items-center justify-center shrink-0">
                        <GitCommit class="w-5 h-5" />
                    </div>
                    <div class="min-w-0">
                        <p class="text-[11px] font-medium text-slate-400 dark:text-surface-400">Installed Commit</p>
                        <p class="text-xs font-semibold text-slate-900 dark:text-white font-mono truncate">{{ info.current_commit }}</p>
                    </div>
                </div>
            </div>

            <!-- Main Content: GitHub Commits & CLI Updater -->
            <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
                <!-- Recent Commits / Release Changelog (2 Cols) -->
                <div class="lg:col-span-2 space-y-4">
                    <div class="flex items-center justify-between">
                        <div class="flex items-center gap-2">
                            <GitBranch class="w-4 h-4 text-brand-600 dark:text-brand-400" />
                            <h3 class="text-sm font-semibold text-slate-900 dark:text-white">Latest GitHub Commits</h3>
                        </div>
                        <a
                            :href="`${info.repository}/commits/main`"
                            target="_blank"
                            class="text-xs text-brand-600 dark:text-brand-400 hover:underline flex items-center gap-1 font-medium"
                        >
                            <span>View on GitHub</span>
                            <ExternalLink class="w-3.5 h-3.5" />
                        </a>
                    </div>

                    <!-- Commits Timeline List -->
                    <div class="rounded-2xl bg-white dark:bg-surface-900/60 border border-slate-200/80 dark:border-surface-800/80 divide-y divide-slate-100 dark:divide-surface-800/60 overflow-hidden shadow-sm dark:shadow-xl">
                        <div
                            v-for="(commit, idx) in info.recent_commits"
                            :key="commit.hash"
                            class="p-4 sm:p-5 hover:bg-slate-50 dark:hover:bg-surface-800/30 transition flex items-start justify-between gap-4"
                        >
                            <div class="space-y-1.5 min-w-0">
                                <div class="flex items-center gap-2">
                                    <span
                                        v-if="idx === 0"
                                        class="px-2 py-0.5 rounded text-[10px] font-semibold bg-brand-50 text-brand-700 border border-brand-200 dark:bg-brand-500/20 dark:text-brand-300 dark:border-brand-500/30"
                                    >
                                        LATEST RELEASE
                                    </span>
                                    <a
                                        :href="commit.url"
                                        target="_blank"
                                        class="text-xs font-mono font-medium text-brand-600 dark:text-brand-400 hover:underline"
                                    >
                                        #{{ commit.short_hash }}
                                    </a>
                                    <span class="text-slate-300 dark:text-surface-600">•</span>
                                    <span class="text-[11px] text-slate-400 dark:text-surface-400">{{ formatDate(commit.date) }}</span>
                                </div>
                                <p class="text-xs text-slate-900 dark:text-white font-medium leading-snug break-words">
                                    {{ commit.message }}
                                </p>
                                <p class="text-[11px] text-slate-500 dark:text-surface-400">
                                    By <span class="text-slate-700 dark:text-surface-300 font-medium">{{ commit.author }}</span>
                                </p>
                            </div>

                            <a
                                :href="commit.url"
                                target="_blank"
                                class="p-2 rounded-xl bg-slate-100 hover:bg-slate-200 dark:bg-surface-800/50 dark:hover:bg-surface-800 text-slate-500 hover:text-slate-900 dark:text-surface-400 dark:hover:text-white transition shrink-0"
                                title="View Commit Diff"
                            >
                                <ExternalLink class="w-3.5 h-3.5" />
                            </a>
                        </div>

                        <div v-if="!info.recent_commits || info.recent_commits.length === 0" class="p-8 text-center text-xs text-slate-400 dark:text-surface-500">
                            No recent commit history retrieved.
                        </div>
                    </div>
                </div>

                <!-- Right Column: Safety Features & Terminal 1-Liner -->
                <div class="space-y-6">
                    <!-- Update Safety Guarantee Card -->
                    <div class="rounded-2xl bg-white dark:bg-surface-900/60 border border-slate-200/80 dark:border-surface-800/80 p-5 space-y-3.5 shadow-sm dark:shadow-xl">
                        <div class="flex items-center gap-2 text-slate-900 dark:text-white text-xs font-semibold">
                            <ShieldCheck class="w-4 h-4 text-emerald-600 dark:text-emerald-400" />
                            <span>Zero-Downtime Guarantee</span>
                        </div>
                        <ul class="space-y-2 text-[11px] text-slate-500 dark:text-surface-400">
                            <li class="flex items-start gap-2">
                                <Check class="w-3.5 h-3.5 text-emerald-600 dark:text-emerald-400 shrink-0 mt-0.5" />
                                <span>Preserves `.env`, SQLite database & storage directories.</span>
                            </li>
                            <li class="flex items-start gap-2">
                                <Check class="w-3.5 h-3.5 text-emerald-600 dark:text-emerald-400 shrink-0 mt-0.5" />
                                <span>Atomic Nginx validation (`nginx -t`) before configuration reload.</span>
                            </li>
                            <li class="flex items-start gap-2">
                                <Check class="w-3.5 h-3.5 text-emerald-600 dark:text-emerald-400 shrink-0 mt-0.5" />
                                <span>Go Agent compiles and restarts via managed Systemd unit.</span>
                            </li>
                        </ul>
                    </div>

                    <!-- Terminal One-Liner Alternative Card -->
                    <div class="rounded-2xl bg-white dark:bg-surface-900/60 border border-slate-200/80 dark:border-surface-800/80 p-5 space-y-3 shadow-sm dark:shadow-xl">
                        <div class="flex items-center justify-between">
                            <div class="flex items-center gap-2 text-slate-900 dark:text-white text-xs font-semibold">
                                <Terminal class="w-4 h-4 text-brand-600 dark:text-brand-400" />
                                <span>Terminal 1-Line Updater</span>
                            </div>
                            <button
                                @click="copyCli"
                                class="text-[11px] text-brand-600 dark:text-brand-400 hover:underline flex items-center gap-1 font-medium transition"
                            >
                                <Check v-if="copiedCli" class="w-3 h-3 text-emerald-500" />
                                <Copy v-else class="w-3 h-3" />
                                <span>{{ copiedCli ? 'Copied' : 'Copy' }}</span>
                            </button>
                        </div>
                        <p class="text-[11px] text-slate-500 dark:text-surface-400 leading-relaxed">
                            You can also upgrade your server directly from your terminal or SSH session:
                        </p>
                        <div class="p-3 rounded-xl bg-slate-900 dark:bg-surface-950 font-mono text-[11px] text-emerald-400 break-all select-all shadow-inner">
                            {{ cliCommand }}
                        </div>
                    </div>
                </div>
            </div>

            <!-- Confirmation Modal -->
            <div
                v-if="showConfirmModal"
                class="fixed inset-0 bg-slate-900/60 dark:bg-surface-950/80 backdrop-blur-md flex items-center justify-center p-4 z-50 animate-in fade-in duration-150"
            >
                <div class="w-full max-w-md rounded-2xl bg-white dark:bg-surface-900 border border-slate-200/90 dark:border-surface-800 p-6 space-y-5 shadow-2xl">
                    <div class="flex items-center gap-3">
                        <div class="w-10 h-10 rounded-xl bg-brand-50 text-brand-600 dark:bg-brand-500/10 dark:border-brand-500/20 dark:text-brand-400 flex items-center justify-center shrink-0">
                            <DownloadCloud class="w-5 h-5" />
                        </div>
                        <div>
                            <h3 class="text-sm font-semibold text-slate-900 dark:text-white">Upgrade Kodepreneur Panel?</h3>
                            <p class="text-xs text-slate-500 dark:text-surface-400">Target branch: <code class="text-brand-600 dark:text-brand-400 font-mono">{{ info.branch }}</code></p>
                        </div>
                    </div>

                    <p class="text-xs text-slate-600 dark:text-surface-300 leading-relaxed">
                        This operation will fetch the latest commits, update dependencies, run pending database migrations, build the frontend, and restart the Go Agent daemon.
                    </p>

                    <div class="p-3 rounded-xl bg-slate-50 dark:bg-surface-950/60 border border-slate-200/80 dark:border-surface-800/80 text-[11px] text-slate-600 dark:text-surface-400 space-y-1.5 font-mono">
                        <div class="flex justify-between">
                            <span>Repository:</span>
                            <span class="text-slate-900 dark:text-surface-200">kodepreneur/kodepreneur-panel</span>
                        </div>
                        <div class="flex justify-between">
                            <span>Target Commit:</span>
                            <span class="text-emerald-600 dark:text-emerald-400 font-semibold">{{ info.latest_commit?.short_hash || 'Latest' }}</span>
                        </div>
                    </div>

                    <div class="flex items-center justify-end gap-2.5 pt-2">
                        <button
                            @click="showConfirmModal = false"
                            class="px-4 py-2 rounded-xl bg-slate-100 hover:bg-slate-200 dark:bg-surface-800 dark:hover:bg-surface-700 text-slate-700 dark:text-surface-300 text-xs font-semibold transition"
                        >
                            Cancel
                        </button>
                        <button
                            @click="startUpdate"
                            class="px-4 py-2 rounded-xl bg-brand-600 hover:bg-brand-500 text-white text-xs font-semibold flex items-center gap-1.5 shadow-md shadow-brand-500/20 transition"
                        >
                            <Zap class="w-3.5 h-3.5" />
                            <span>Confirm & Upgrade</span>
                        </button>
                    </div>
                </div>
            </div>

            <!-- Live Upgrade Terminal / Progress Modal -->
            <div
                v-if="showTerminalModal"
                class="fixed inset-0 bg-slate-900/60 dark:bg-surface-950/85 backdrop-blur-md flex items-center justify-center p-4 z-50 animate-in fade-in duration-150"
            >
                <div class="w-full max-w-3xl rounded-2xl bg-white dark:bg-surface-900 border border-slate-200/90 dark:border-surface-800 shadow-2xl flex flex-col max-h-[90vh] overflow-hidden">
                    <!-- Modal Header -->
                    <div class="px-6 py-4 border-b border-slate-100 dark:border-surface-800/80 flex items-center justify-between">
                        <div class="flex items-center gap-3">
                            <div
                                :class="[
                                    'w-8 h-8 rounded-lg flex items-center justify-center',
                                    isUpdating ? 'bg-brand-50 text-brand-600 dark:bg-brand-500/10 dark:text-brand-400 border border-brand-200 dark:border-brand-500/20' : '',
                                    updateSuccess === true ? 'bg-emerald-50 text-emerald-600 dark:bg-emerald-500/10 dark:text-emerald-400 border border-emerald-200 dark:border-emerald-500/20' : '',
                                    updateSuccess === false ? 'bg-rose-50 text-rose-600 dark:bg-rose-500/10 dark:text-rose-400 border border-rose-200 dark:border-rose-500/20' : ''
                                ]"
                            >
                                <RefreshCw v-if="isUpdating" class="w-4 h-4 animate-spin" />
                                <CheckCircle2 v-else-if="updateSuccess === true" class="w-4 h-4" />
                                <AlertCircle v-else class="w-4 h-4" />
                            </div>
                            <div>
                                <h3 class="text-sm font-semibold text-slate-900 dark:text-white">
                                    {{ isUpdating ? 'Upgrading Kodepreneur Panel...' : (updateSuccess ? 'Upgrade Completed Successfully!' : 'Upgrade Failed') }}
                                </h3>
                                <p class="text-[11px] text-slate-500 dark:text-surface-400">
                                    {{ isUpdating ? `Step ${currentStep} of ${totalSteps}: ${stepLabels[currentStep - 1]}` : (updateSuccess ? `Executed in ${updateDuration}s` : 'See error logs below') }}
                                </p>
                            </div>
                        </div>

                        <div class="flex items-center gap-2">
                            <button
                                @click="copyLogs"
                                class="px-3 py-1.5 rounded-lg bg-slate-100 hover:bg-slate-200 dark:bg-surface-800 dark:hover:bg-surface-700 text-slate-700 dark:text-surface-300 text-xs font-medium flex items-center gap-1.5 transition"
                            >
                                <Check v-if="copiedLogs" class="w-3.5 h-3.5 text-emerald-500" />
                                <Copy v-else class="w-3.5 h-3.5" />
                                <span>{{ copiedLogs ? 'Copied' : 'Copy Logs' }}</span>
                            </button>
                        </div>
                    </div>

                    <!-- Progress Step Bar (While Running) -->
                    <div v-if="isUpdating" class="px-6 py-3 bg-slate-50 dark:bg-surface-950/40 border-b border-slate-100 dark:border-surface-800/60">
                        <div class="flex items-center justify-between text-[11px] text-slate-600 dark:text-surface-400 mb-1.5">
                            <span>Pipeline Progress</span>
                            <span class="font-mono text-brand-600 dark:text-brand-400 font-semibold">{{ Math.round((currentStep / totalSteps) * 100) }}%</span>
                        </div>
                        <div class="w-full h-1.5 rounded-full bg-slate-200 dark:bg-surface-800 overflow-hidden">
                            <div
                                class="h-full bg-gradient-to-r from-brand-600 to-indigo-500 transition-all duration-300 rounded-full"
                                :style="{ width: `${(currentStep / totalSteps) * 100}%` }"
                            ></div>
                        </div>
                    </div>

                    <!-- Terminal Output Console -->
                    <div class="flex-1 p-5 bg-slate-950 font-mono text-xs text-slate-200 overflow-y-auto max-h-[420px] whitespace-pre-wrap leading-relaxed select-text border-b border-slate-800">
                        {{ updateLogs }}
                    </div>

                    <!-- Modal Footer -->
                    <div class="px-6 py-4 bg-slate-50 dark:bg-surface-900/90 flex items-center justify-between">
                        <div class="text-[11px] text-slate-500 dark:text-surface-400 font-mono">
                            <span v-if="isUpdating" class="flex items-center gap-1.5 text-brand-600 dark:text-brand-400 font-medium">
                                <span class="w-2 h-2 rounded-full bg-brand-500 animate-ping"></span>
                                Live Daemon Runner Active
                            </span>
                            <span v-else-if="updateSuccess" class="text-emerald-600 dark:text-emerald-400 flex items-center gap-1.5 font-medium">
                                <CheckCircle2 class="w-3.5 h-3.5" />
                                All services synchronized
                            </span>
                            <span v-else class="text-rose-600 dark:text-rose-400 flex items-center gap-1.5 font-medium">
                                <AlertCircle class="w-3.5 h-3.5" />
                                Error occurred during upgrade
                            </span>
                        </div>

                        <div class="flex items-center gap-2">
                            <button
                                v-if="!isUpdating"
                                @click="showTerminalModal = false"
                                class="px-4 py-2 rounded-xl bg-slate-200 hover:bg-slate-300 dark:bg-surface-800 dark:hover:bg-surface-700 text-slate-800 dark:text-surface-300 text-xs font-semibold transition"
                            >
                                Close
                            </button>
                            <button
                                v-if="updateSuccess"
                                @click="reloadPage"
                                class="px-4 py-2 rounded-xl bg-emerald-600 hover:bg-emerald-500 text-white text-xs font-semibold flex items-center gap-1.5 shadow-md shadow-emerald-600/20 transition"
                            >
                                <RotateCcw class="w-3.5 h-3.5" />
                                <span>Reload Panel</span>
                            </button>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </AppLayout>
</template>
