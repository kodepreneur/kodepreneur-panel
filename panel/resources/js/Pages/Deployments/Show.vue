<script setup lang="ts">
import { ref } from 'vue';
import { Link, router } from '@inertiajs/vue3';
import AppLayout from '@/Layouts/AppLayout.vue';
import {
    Terminal,
    ArrowLeft,
    CheckCircle2,
    XCircle,
    Clock,
    GitBranch,
    Globe,
    RotateCcw,
    Copy,
    Check,
    Play,
} from 'lucide-vue-next';
import type { Deployment, Website, User } from '@/types';

const props = defineProps<{
    deployment: Deployment & {
        website: Website;
        initiated_by?: User | null;
    };
}>();

const copied = ref(false);

function copyLogs() {
    if (!props.deployment.log_output) return;
    navigator.clipboard.writeText(props.deployment.log_output);
    copied.value = true;
    setTimeout(() => {
        copied.value = false;
    }, 2000);
}

function rerunDeployment() {
    if (confirm(`Trigger a new deployment for ${props.deployment.website.domain}?`)) {
        router.post(`/websites/${props.deployment.website.id}/deploy`);
    }
}

function rollbackToThis() {
    if (confirm(`Rollback ${props.deployment.website.domain} to deployment #${props.deployment.id} (${props.deployment.commit_hash || 'this commit'})?`)) {
        router.post(`/deployments/${props.deployment.id}/rollback`);
    }
}
</script>

<template>
    <AppLayout :title="`Deployment #${deployment.id}`">
        <div class="max-w-6xl mx-auto space-y-6">
            <!-- Header -->
            <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
                <div class="flex items-center gap-3">
                    <Link
                        href="/deployments"
                        class="p-2 rounded-xl bg-white dark:bg-surface-900 border border-slate-200/80 dark:border-surface-800 text-slate-500 hover:text-slate-900 dark:text-surface-400 dark:hover:text-white transition shadow-sm"
                    >
                        <ArrowLeft class="w-4 h-4" />
                    </Link>
                    <div>
                        <div class="flex items-center gap-2.5">
                            <h2 class="text-base font-bold text-slate-900 dark:text-white tracking-tight">
                                Deployment Run #{{ deployment.id }}
                            </h2>
                            <span
                                :class="[
                                    'px-2.5 py-0.5 rounded-full text-[10px] font-semibold uppercase font-mono flex items-center gap-1',
                                    deployment.status === 'success'
                                        ? 'bg-emerald-50 text-emerald-700 border border-emerald-200 dark:bg-emerald-500/10 dark:border-emerald-500/30 dark:text-emerald-400'
                                        : 'bg-rose-50 text-rose-700 border border-rose-200 dark:bg-rose-500/10 dark:border-rose-500/30 dark:text-rose-400'
                                ]"
                            >
                                <CheckCircle2 v-if="deployment.status === 'success'" class="w-3 h-3" />
                                <XCircle v-else class="w-3 h-3" />
                                <span>{{ deployment.status }}</span>
                            </span>
                        </div>
                        <p class="text-xs text-slate-500 dark:text-surface-400 mt-0.5">
                            Target:
                            <Link :href="`/websites/${deployment.website.id}`" class="text-brand-600 dark:text-brand-400 hover:underline font-mono">
                                {{ deployment.website.domain }}
                            </Link>
                            • Branch: <span class="font-mono text-slate-700 dark:text-surface-300">{{ deployment.branch }}</span>
                            <span v-if="deployment.commit_hash"> • Commit: <code class="text-emerald-600 dark:text-emerald-400 font-mono">{{ deployment.commit_hash }}</code></span>
                        </p>
                    </div>
                </div>

                <div class="flex items-center gap-2">
                    <button
                        @click="copyLogs"
                        class="px-3 py-1.5 rounded-xl bg-slate-100 hover:bg-slate-200 dark:bg-surface-800 dark:hover:bg-surface-700 text-slate-700 dark:text-surface-300 text-xs font-semibold flex items-center gap-1.5 transition border border-slate-200/80 dark:border-surface-700"
                    >
                        <Check v-if="copied" class="w-3.5 h-3.5 text-emerald-600 dark:text-emerald-400" />
                        <Copy v-else class="w-3.5 h-3.5" />
                        <span>{{ copied ? 'Copied' : 'Copy Logs' }}</span>
                    </button>
                    <button
                        v-if="deployment.commit_hash && deployment.status === 'success'"
                        @click="rollbackToThis"
                        class="px-3.5 py-1.5 rounded-xl bg-amber-50 hover:bg-amber-100 dark:bg-amber-500/10 dark:hover:bg-amber-500/20 border border-amber-200 dark:border-amber-500/30 text-amber-700 dark:text-amber-300 text-xs font-semibold flex items-center gap-1.5 transition"
                    >
                        <RotateCcw class="w-3.5 h-3.5" />
                        <span>Rollback to Here</span>
                    </button>
                    <button
                        @click="rerunDeployment"
                        class="px-3.5 py-1.5 rounded-xl bg-brand-600 hover:bg-brand-500 text-white text-xs font-semibold flex items-center gap-1.5 transition shadow-md shadow-brand-600/20"
                    >
                        <Play class="w-3.5 h-3.5" />
                        <span>Re-run Pipeline</span>
                    </button>
                </div>
            </div>

            <!-- Terminal Output Console -->
            <div class="rounded-2xl bg-slate-950 border border-slate-800 shadow-2xl overflow-hidden">
                <!-- Terminal Titlebar -->
                <div class="h-10 bg-slate-900 border-b border-slate-800 px-4 flex items-center justify-between">
                    <div class="flex items-center gap-2">
                        <div class="w-3 h-3 rounded-full bg-rose-500/80"></div>
                        <div class="w-3 h-3 rounded-full bg-amber-500/80"></div>
                        <div class="w-3 h-3 rounded-full bg-emerald-500/80"></div>
                        <span class="text-xs text-slate-400 font-mono ml-2">runner@kodepreneur (~/{{ deployment.website.domain }})</span>
                    </div>
                    <span class="text-[11px] text-slate-400 font-mono flex items-center gap-1">
                        <Clock class="w-3.5 h-3.5" />
                        <span>{{ deployment.duration_seconds || 1 }}s runtime</span>
                    </span>
                </div>

                <!-- Console Stream -->
                <div class="p-4 sm:p-6 font-mono text-xs text-emerald-400 leading-relaxed overflow-x-auto selection:bg-brand-500 selection:text-white bg-black/70 min-h-[300px]">
                    <pre class="whitespace-pre-wrap font-mono">{{ deployment.log_output || 'No output recorded for this deployment run.' }}</pre>
                </div>
            </div>
        </div>
    </AppLayout>
</template>
