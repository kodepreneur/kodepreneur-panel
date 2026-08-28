<script setup lang="ts">
import { Link } from '@inertiajs/vue3';
import AppLayout from '@/Layouts/AppLayout.vue';
import { GitBranch, Globe, CheckCircle2, XCircle, Clock, ExternalLink } from 'lucide-vue-next';
import type { Deployment, Website, User } from '@/types';

defineProps<{
    deployments: {
        data: (Deployment & { website: Website; initiated_by?: User | null })[];
        links: any[];
    };
}>();
</script>

<template>
    <AppLayout title="Deployments">
        <div class="space-y-6 max-w-7xl mx-auto">
            <!-- Header -->
            <div>
                <h2 class="text-base font-semibold text-white tracking-tight">Deployment Pipeline</h2>
                <p class="text-xs text-surface-400 mt-0.5">Track real-time Git build logs, execution durations, and release histories</p>
            </div>

            <!-- Deployments Table Card -->
            <div class="rounded-2xl bg-surface-900/60 border border-surface-800/80 shadow-lg overflow-hidden">
                <div v-if="deployments.data.length === 0" class="text-center py-12 px-4">
                    <GitBranch class="w-10 h-10 text-surface-600 mx-auto mb-3" />
                    <h3 class="text-sm font-medium text-surface-200">No deployment runs recorded yet</h3>
                    <p class="text-xs text-surface-400 mt-1 max-w-sm mx-auto">
                        Trigger a deployment from any website view to run the automated build pipeline.
                    </p>
                </div>

                <div v-else class="overflow-x-auto">
                    <table class="w-full text-left text-xs">
                        <thead class="bg-surface-950/40 text-surface-400 uppercase text-[10px] tracking-wider border-b border-surface-800">
                            <tr>
                                <th class="py-3 px-4 font-semibold">Run ID</th>
                                <th class="py-3 px-4 font-semibold">Website</th>
                                <th class="py-3 px-4 font-semibold">Branch</th>
                                <th class="py-3 px-4 font-semibold">Trigger</th>
                                <th class="py-3 px-4 font-semibold">Duration</th>
                                <th class="py-3 px-4 font-semibold">Status</th>
                                <th class="py-3 px-4 font-semibold text-right">Logs</th>
                            </tr>
                        </thead>
                        <tbody class="divide-y divide-surface-800/60">
                            <tr
                                v-for="deploy in deployments.data"
                                :key="deploy.id"
                                class="hover:bg-surface-800/30 transition"
                            >
                                <td class="py-3.5 px-4 font-mono font-medium text-white">
                                    #{{ deploy.id }}
                                </td>
                                <td class="py-3.5 px-4">
                                    <Link
                                        :href="`/websites/${deploy.website.id}`"
                                        class="text-brand-400 hover:underline font-medium"
                                    >
                                        {{ deploy.website.domain }}
                                    </Link>
                                </td>
                                <td class="py-3.5 px-4 font-mono text-[11px] text-surface-300">
                                    <span class="px-2 py-0.5 rounded-md bg-surface-800 border border-surface-700">
                                        {{ deploy.branch }}
                                    </span>
                                </td>
                                <td class="py-3.5 px-4 text-surface-400">
                                    {{ deploy.trigger_source }}
                                </td>
                                <td class="py-3.5 px-4 font-mono text-[11px] text-surface-400">
                                    {{ deploy.duration_seconds ? `${deploy.duration_seconds}s` : '-' }}
                                </td>
                                <td class="py-3.5 px-4">
                                    <span
                                        v-if="deploy.status === 'success'"
                                        class="inline-flex items-center gap-1 text-[11px] text-emerald-400 font-medium"
                                    >
                                        <CheckCircle2 class="w-3.5 h-3.5" />
                                        <span>Success</span>
                                    </span>
                                    <span
                                        v-else
                                        class="inline-flex items-center gap-1 text-[11px] text-rose-400 font-medium"
                                    >
                                        <XCircle class="w-3.5 h-3.5" />
                                        <span>{{ deploy.status }}</span>
                                    </span>
                                </td>
                                <td class="py-3.5 px-4 text-right">
                                    <Link
                                        :href="`/deployments/${deploy.id}`"
                                        class="inline-flex items-center gap-1 px-2.5 py-1 rounded-lg bg-surface-800 hover:bg-surface-700 text-surface-200 text-[11px] font-mono transition"
                                    >
                                        <span>Terminal</span>
                                        <ExternalLink class="w-3 h-3 text-surface-400" />
                                    </Link>
                                </td>
                            </tr>
                        </tbody>
                    </table>
                </div>
            </div>
        </div>
    </AppLayout>
</template>
