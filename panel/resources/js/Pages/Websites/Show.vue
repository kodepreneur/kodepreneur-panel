<script setup lang="ts">
import { ref, onMounted } from 'vue';
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
} from 'lucide-vue-next';
import type { Website, Deployment, SslCertificate, Domain } from '@/types';

const props = defineProps<{
    website: Website & {
        domains: Domain[];
        deployments: Deployment[];
        ssl_certificate?: SslCertificate | null;
    };
}>();

const activeTab = ref<'overview' | 'php' | 'ssl' | 'logs' | 'deployments'>('overview');

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
                        class="p-2 rounded-xl bg-surface-900 border border-surface-800 text-surface-400 hover:text-white transition"
                    >
                        <ArrowLeft class="w-4 h-4" />
                    </Link>
                    <div>
                        <div class="flex items-center gap-2.5">
                            <h2 class="text-base font-semibold text-white tracking-tight">{{ website.domain }}</h2>
                            <span class="px-2 py-0.5 rounded-full text-[10px] font-medium bg-emerald-500/10 border border-emerald-500/30 text-emerald-400">
                                {{ website.status }}
                            </span>
                        </div>
                        <div class="flex items-center gap-3 text-xs text-surface-400 font-mono mt-0.5">
                            <span>user: {{ website.system_user }}</span>
                            <span>•</span>
                            <a
                                :href="`http://${website.domain}`"
                                target="_blank"
                                class="text-brand-400 hover:underline flex items-center gap-1"
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
                        class="px-4 py-2 rounded-xl bg-brand-600 hover:bg-brand-500 text-white text-xs font-semibold shadow-lg shadow-brand-600/20 flex items-center gap-1.5 transition"
                    >
                        <Play class="w-3.5 h-3.5 fill-current" />
                        <span>Deploy Now</span>
                    </button>
                </div>
            </div>

            <!-- Navigation Tabs -->
            <div class="flex items-center gap-2 border-b border-surface-800 pb-2 overflow-x-auto text-xs">
                <button
                    @click="activeTab = 'overview'"
                    :class="[
                        'px-3.5 py-2 rounded-xl font-semibold flex items-center gap-1.5 transition whitespace-nowrap',
                        activeTab === 'overview'
                            ? 'bg-brand-500/10 text-brand-400 border border-brand-500/30'
                            : 'text-surface-400 hover:text-white'
                    ]"
                >
                    <Globe class="w-3.5 h-3.5" />
                    <span>Overview</span>
                </button>

                <button
                    @click="activeTab = 'php'"
                    :class="[
                        'px-3.5 py-2 rounded-xl font-semibold flex items-center gap-1.5 transition whitespace-nowrap',
                        activeTab === 'php'
                            ? 'bg-brand-500/10 text-brand-400 border border-brand-500/30'
                            : 'text-surface-400 hover:text-white'
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
                            ? 'bg-brand-500/10 text-brand-400 border border-brand-500/30'
                            : 'text-surface-400 hover:text-white'
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
                            ? 'bg-brand-500/10 text-brand-400 border border-brand-500/30'
                            : 'text-surface-400 hover:text-white'
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
                            ? 'bg-brand-500/10 text-brand-400 border border-brand-500/30'
                            : 'text-surface-400 hover:text-white'
                    ]"
                >
                    <GitBranch class="w-3.5 h-3.5" />
                    <span>Deployments ({{ website.deployments.length }})</span>
                </button>
            </div>

            <!-- Tab 1: Overview -->
            <div v-if="activeTab === 'overview'" class="space-y-6">
                <!-- Summary Metrics Cards -->
                <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
                    <div class="rounded-2xl bg-surface-900/60 border border-surface-800/80 p-5 shadow-lg space-y-2">
                        <div class="flex items-center gap-2 text-xs font-medium text-surface-400">
                            <Cpu class="w-4 h-4 text-brand-400" />
                            <span>PHP Runtime</span>
                        </div>
                        <p class="text-base font-bold text-white font-mono">
                            {{ website.php_version !== 'none' ? `PHP ${website.php_version}` : 'Static HTML' }}
                        </p>
                        <p class="text-[11px] text-surface-500 font-mono">
                            Socket: /run/php/php{{ website.php_version }}-fpm-{{ website.system_user }}.sock
                        </p>
                    </div>

                    <div class="rounded-2xl bg-surface-900/60 border border-surface-800/80 p-5 shadow-lg space-y-2">
                        <div class="flex items-center gap-2 text-xs font-medium text-surface-400">
                            <Folder class="w-4 h-4 text-cyan-400" />
                            <span>Document Root</span>
                        </div>
                        <p class="text-xs font-bold text-white font-mono truncate">{{ website.document_root }}</p>
                        <p class="text-[11px] text-surface-500">Physical document root</p>
                    </div>

                    <div class="rounded-2xl bg-surface-900/60 border border-surface-800/80 p-5 shadow-lg space-y-2">
                        <div class="flex items-center gap-2 text-xs font-medium text-surface-400">
                            <ShieldCheck v-if="website.ssl_enabled" class="w-4 h-4 text-emerald-400" />
                            <ShieldAlert v-else class="w-4 h-4 text-amber-400" />
                            <span>SSL Security</span>
                        </div>
                        <p class="text-base font-bold text-white font-mono">
                            {{ website.ssl_enabled ? "Let's Encrypt" : 'HTTP Only' }}
                        </p>
                        <p class="text-[11px] text-surface-500">
                            {{ website.ssl_enabled ? 'Auto-renew active' : 'No certificate issued' }}
                        </p>
                    </div>
                </div>

                <!-- Webroot Settings Card -->
                <div class="rounded-2xl bg-surface-900/60 border border-surface-800/80 p-6 shadow-lg space-y-4">
                    <h3 class="text-xs font-semibold text-white uppercase tracking-wider flex items-center gap-2">
                        <Settings class="w-4 h-4 text-brand-400" />
                        <span>Directory & Security Settings</span>
                    </h3>

                    <form @submit.prevent="updateSettings" class="space-y-4 max-w-xl">
                        <div>
                            <label class="block text-xs font-medium text-surface-300 mb-1.5">Document Root</label>
                            <input
                                v-model="settingsForm.document_root"
                                type="text"
                                class="w-full px-3.5 py-2.5 rounded-xl bg-surface-950/60 border border-surface-800 text-xs text-white font-mono focus:ring-2 focus:ring-brand-500/40 focus:outline-none transition"
                            />
                        </div>

                        <label class="flex items-center gap-2.5 cursor-pointer">
                            <input
                                v-model="settingsForm.force_https"
                                type="checkbox"
                                class="rounded border-surface-700 bg-surface-900 text-brand-600 focus:ring-brand-500"
                            />
                            <span class="text-xs text-surface-300">Force HTTP to HTTPS redirect</span>
                        </label>

                        <button
                            type="submit"
                            :disabled="settingsForm.processing"
                            class="px-4 py-2 rounded-xl bg-brand-600 hover:bg-brand-500 text-white text-xs font-semibold shadow-md transition disabled:opacity-50"
                        >
                            {{ settingsForm.processing ? 'Saving...' : 'Save Settings' }}
                        </button>
                    </form>
                </div>
            </div>

            <!-- Tab 2: PHP Runtime Switcher -->
            <div v-if="activeTab === 'php'" class="rounded-2xl bg-surface-900/60 border border-surface-800/80 p-6 shadow-lg space-y-6">
                <div>
                    <h3 class="text-sm font-semibold text-white">PHP Runtime & Socket Pool</h3>
                    <p class="text-xs text-surface-400 mt-0.5">
                        Switch PHP versions with zero-downtime pool recreation and automated Nginx socket rebinding.
                    </p>
                </div>

                <div class="grid grid-cols-1 sm:grid-cols-3 gap-3">
                    <label
                        :class="[
                            'cursor-pointer rounded-xl border p-4 flex flex-col justify-between transition',
                            phpForm.php_version === '8.3'
                                ? 'bg-brand-500/10 border-brand-500 text-white ring-1 ring-brand-500'
                                : 'bg-surface-950/40 border-surface-800 text-surface-400 hover:border-surface-700'
                        ]"
                    >
                        <div class="flex items-center justify-between">
                            <span class="text-xs font-bold font-mono">PHP 8.3 LTS</span>
                            <input type="radio" v-model="phpForm.php_version" value="8.3" class="sr-only" />
                            <Check v-if="phpForm.php_version === '8.3'" class="w-4 h-4 text-brand-400" />
                        </div>
                        <p class="text-[11px] text-surface-400 mt-2">Default LTS runtime for maximum stability.</p>
                    </label>

                    <label
                        :class="[
                            'cursor-pointer rounded-xl border p-4 flex flex-col justify-between transition',
                            phpForm.php_version === '8.4'
                                ? 'bg-brand-500/10 border-brand-500 text-white ring-1 ring-brand-500'
                                : 'bg-surface-950/40 border-surface-800 text-surface-400 hover:border-surface-700'
                        ]"
                    >
                        <div class="flex items-center justify-between">
                            <span class="text-xs font-bold font-mono">PHP 8.4 Latest</span>
                            <input type="radio" v-model="phpForm.php_version" value="8.4" class="sr-only" />
                            <Check v-if="phpForm.php_version === '8.4'" class="w-4 h-4 text-brand-400" />
                        </div>
                        <p class="text-[11px] text-surface-400 mt-2">Latest release with improved JIT performance.</p>
                    </label>

                    <label
                        :class="[
                            'cursor-pointer rounded-xl border p-4 flex flex-col justify-between transition',
                            phpForm.php_version === 'none'
                                ? 'bg-brand-500/10 border-brand-500 text-white ring-1 ring-brand-500'
                                : 'bg-surface-950/40 border-surface-800 text-surface-400 hover:border-surface-700'
                        ]"
                    >
                        <div class="flex items-center justify-between">
                            <span class="text-xs font-bold">Static / Proxy</span>
                            <input type="radio" v-model="phpForm.php_version" value="none" class="sr-only" />
                            <Check v-if="phpForm.php_version === 'none'" class="w-4 h-4 text-brand-400" />
                        </div>
                        <p class="text-[11px] text-surface-400 mt-2">Static assets without PHP processing.</p>
                    </label>
                </div>

                <div class="pt-4 border-t border-surface-800 flex items-center justify-between">
                    <span class="text-xs text-surface-400">
                        Current Active Runtime: <strong class="text-white font-mono">PHP {{ website.php_version }}</strong>
                    </span>
                    <button
                        @click="switchPhp"
                        :disabled="phpForm.processing || phpForm.php_version === website.php_version"
                        class="px-4 py-2 rounded-xl bg-brand-600 hover:bg-brand-500 text-white text-xs font-semibold shadow-lg shadow-brand-600/20 transition disabled:opacity-50"
                    >
                        {{ phpForm.processing ? 'Switching Runtime...' : 'Apply PHP Switch' }}
                    </button>
                </div>
            </div>

            <!-- Tab 3: SSL / HTTPS -->
            <div v-if="activeTab === 'ssl'" class="space-y-6">
                <!-- Certificate Status Card -->
                <div class="rounded-2xl bg-surface-900/60 border border-surface-800/80 p-6 shadow-lg space-y-4">
                    <div class="flex items-center justify-between">
                        <h3 class="text-sm font-semibold text-white flex items-center gap-2">
                            <ShieldCheck class="w-4 h-4 text-emerald-400" />
                            <span>SSL Certificate Status</span>
                        </h3>
                        <span
                            :class="[
                                'px-2.5 py-0.5 rounded-full text-[10px] font-semibold',
                                website.ssl_enabled
                                    ? 'bg-emerald-500/10 text-emerald-400 border border-emerald-500/30'
                                    : 'bg-amber-500/10 text-amber-400 border border-amber-500/30'
                            ]"
                        >
                            {{ website.ssl_enabled ? 'Active HTTPS' : 'HTTP Only' }}
                        </span>
                    </div>

                    <div v-if="website.ssl_certificate" class="grid grid-cols-1 sm:grid-cols-3 gap-4 pt-2">
                        <div class="p-3.5 rounded-xl bg-surface-950/40 border border-surface-800">
                            <span class="text-[10px] text-surface-400 uppercase font-semibold">Certificate Issuer</span>
                            <p class="text-xs font-bold text-white mt-1">{{ website.ssl_certificate.issuer }}</p>
                        </div>
                        <div class="p-3.5 rounded-xl bg-surface-950/40 border border-surface-800">
                            <span class="text-[10px] text-surface-400 uppercase font-semibold">Valid Until</span>
                            <p class="text-xs font-bold text-white mt-1 font-mono">
                                {{ new Date(website.ssl_certificate.valid_until || '').toLocaleDateString() }}
                            </p>
                        </div>
                        <div class="p-3.5 rounded-xl bg-surface-950/40 border border-surface-800">
                            <span class="text-[10px] text-surface-400 uppercase font-semibold">Auto-Renewal</span>
                            <p class="text-xs font-bold text-emerald-400 mt-1">Enabled (via Cron)</p>
                        </div>
                    </div>

                    <div v-else class="text-xs text-surface-400 bg-surface-950/40 p-4 rounded-xl border border-surface-800">
                        No SSL certificate is currently active for <strong>{{ website.domain }}</strong>. Use the form below to request a free Let's Encrypt certificate.
                    </div>
                </div>

                <!-- Issue Certificate Form Card -->
                <div class="rounded-2xl bg-surface-900/60 border border-surface-800/80 p-6 shadow-lg space-y-4">
                    <h3 class="text-sm font-semibold text-white">
                        {{ website.ssl_enabled ? 'Renew or Re-issue Certificate' : 'Issue Let\'s Encrypt Certificate' }}
                    </h3>
                    <p class="text-xs text-surface-400">
                        The agent will execute an ACME HTTP-01 challenge against <code class="text-brand-400">{{ website.domain }}</code> and automatically inject SSL directives into Nginx.
                    </p>

                    <form @submit.prevent="issueSsl" class="space-y-4 max-w-lg">
                        <div>
                            <label class="block text-xs font-medium text-surface-300 mb-1.5">Notification Email</label>
                            <input
                                v-model="sslForm.email"
                                type="email"
                                placeholder="admin@example.com"
                                class="w-full px-3.5 py-2.5 rounded-xl bg-surface-950/60 border border-surface-800 text-xs text-white focus:ring-2 focus:ring-brand-500/40 focus:outline-none transition"
                            />
                        </div>

                        <label class="flex items-center gap-2.5 cursor-pointer">
                            <input
                                v-model="sslForm.force_https"
                                type="checkbox"
                                class="rounded border-surface-700 bg-surface-900 text-brand-600 focus:ring-brand-500"
                            />
                            <span class="text-xs text-surface-300">Automatically redirect all HTTP traffic to HTTPS (301)</span>
                        </label>

                        <button
                            type="submit"
                            :disabled="sslForm.processing"
                            class="px-5 py-2.5 rounded-xl bg-brand-600 hover:bg-brand-500 text-white text-xs font-semibold shadow-lg shadow-brand-600/20 flex items-center gap-2 transition disabled:opacity-50"
                        >
                            <ShieldCheck class="w-4 h-4" />
                            <span>{{ sslForm.processing ? 'Requesting Let\'s Encrypt...' : 'Issue SSL Certificate' }}</span>
                        </button>
                    </form>
                </div>
            </div>

            <!-- Tab 4: Nginx Logs -->
            <div v-if="activeTab === 'logs'" class="rounded-2xl bg-surface-900/60 border border-surface-800/80 p-6 shadow-lg space-y-4">
                <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3">
                    <div class="flex items-center gap-2">
                        <button
                            @click="{ logType = 'access'; fetchLogs(); }"
                            :class="[
                                'px-3 py-1.5 rounded-lg text-xs font-medium transition',
                                logType === 'access' ? 'bg-surface-800 text-white' : 'text-surface-400 hover:text-white'
                            ]"
                        >
                            Access Log
                        </button>
                        <button
                            @click="{ logType = 'error'; fetchLogs(); }"
                            :class="[
                                'px-3 py-1.5 rounded-lg text-xs font-medium transition',
                                logType === 'error' ? 'bg-surface-800 text-white' : 'text-surface-400 hover:text-white'
                            ]"
                        >
                            Error Log
                        </button>
                    </div>

                    <div class="flex items-center gap-2">
                        <select
                            v-model="logLines"
                            @change="fetchLogs"
                            class="px-2.5 py-1.5 rounded-lg bg-surface-950 border border-surface-800 text-xs text-surface-300"
                        >
                            <option :value="50">Last 50 lines</option>
                            <option :value="100">Last 100 lines</option>
                            <option :value="250">Last 250 lines</option>
                        </select>
                        <button
                            @click="fetchLogs"
                            :disabled="isFetchingLogs"
                            class="p-2 rounded-lg bg-surface-800 hover:bg-surface-700 text-surface-300 transition"
                            title="Refresh Logs"
                        >
                            <RefreshCw :class="['w-3.5 h-3.5', isFetchingLogs ? 'animate-spin' : '']" />
                        </button>
                    </div>
                </div>

                <!-- Log Output Terminal -->
                <div class="p-4 rounded-xl bg-surface-950 border border-surface-800 font-mono text-[11px] text-surface-300 h-96 overflow-y-auto space-y-1 select-text">
                    <div v-if="logs.length === 0" class="text-surface-600 italic">
                        No log entries available.
                    </div>
                    <div v-for="(line, idx) in logs" :key="idx" class="leading-relaxed hover:bg-surface-900/50 px-1 py-0.5 rounded">
                        {{ line }}
                    </div>
                </div>
            </div>

            <!-- Tab 5: Deployments -->
            <div v-if="activeTab === 'deployments'" class="rounded-2xl bg-surface-900/60 border border-surface-800/80 p-6 shadow-lg space-y-4">
                <div class="flex items-center justify-between">
                    <h3 class="text-xs font-semibold text-surface-300 uppercase tracking-wider flex items-center gap-2">
                        <GitBranch class="w-4 h-4 text-purple-400" />
                        <span>Deployment Pipeline History</span>
                    </h3>
                    <button
                        @click="triggerDeploy"
                        class="px-3 py-1.5 rounded-lg bg-brand-600 hover:bg-brand-500 text-white text-xs font-medium transition"
                    >
                        Trigger Build
                    </button>
                </div>

                <div v-if="website.deployments.length === 0" class="text-center py-12 text-xs text-surface-500">
                    No deployments executed for this website yet. Click "Deploy Now" to trigger a Git pull and build hook execution.
                </div>

                <div v-else class="divide-y divide-surface-800/60">
                    <div
                        v-for="deploy in website.deployments"
                        :key="deploy.id"
                        class="py-3.5 flex items-center justify-between text-xs"
                    >
                        <div class="flex items-center gap-3">
                            <span
                                :class="[
                                    'w-2 h-2 rounded-full',
                                    deploy.status === 'success' ? 'bg-emerald-400' : deploy.status === 'running' ? 'bg-amber-400 animate-ping' : 'bg-rose-400'
                                ]"
                            ></span>
                            <div>
                                <p class="text-surface-100 font-mono font-medium">
                                    Deployment #{{ deploy.id }} • {{ deploy.branch }}
                                </p>
                                <p class="text-[11px] text-surface-400">
                                    Trigger: {{ deploy.trigger_source }} • {{ deploy.duration_seconds ?? 0 }}s
                                </p>
                            </div>
                        </div>
                        <Link
                            :href="`/deployments/${deploy.id}`"
                            class="px-3 py-1.5 rounded-lg bg-surface-800 hover:bg-surface-700 text-surface-200 text-xs font-mono transition"
                        >
                            View Terminal Logs
                        </Link>
                    </div>
                </div>
            </div>
        </div>
    </AppLayout>
</template>
