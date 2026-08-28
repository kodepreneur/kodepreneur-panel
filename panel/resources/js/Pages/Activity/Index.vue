<script setup lang="ts">
import AppLayout from '@/Layouts/AppLayout.vue';
import { Activity, Shield, User, Clock } from 'lucide-vue-next';
import type { ActivityLog } from '@/types';

defineProps<{
    logs: {
        data: ActivityLog[];
        links: any[];
    };
}>();
</script>

<template>
    <AppLayout title="Activity Logs">
        <div class="space-y-6 max-w-7xl mx-auto">
            <!-- Header -->
            <div>
                <h2 class="text-base font-semibold text-white tracking-tight">Audit & Activity Log</h2>
                <p class="text-xs text-surface-400 mt-0.5">Immutable audit trail of all server operations and administrator actions</p>
            </div>

            <!-- Activity Table Card -->
            <div class="rounded-2xl bg-surface-900/60 border border-surface-800/80 shadow-lg overflow-hidden">
                <div v-if="logs.data.length === 0" class="text-center py-12 px-4">
                    <Activity class="w-10 h-10 text-surface-600 mx-auto mb-3" />
                    <h3 class="text-sm font-medium text-surface-200">No activity logs recorded yet</h3>
                </div>

                <div v-else class="overflow-x-auto">
                    <table class="w-full text-left text-xs">
                        <thead class="bg-surface-950/40 text-surface-400 uppercase text-[10px] tracking-wider border-b border-surface-800">
                            <tr>
                                <th class="py-3 px-4 font-semibold">Action</th>
                                <th class="py-3 px-4 font-semibold">User</th>
                                <th class="py-3 px-4 font-semibold">IP Address</th>
                                <th class="py-3 px-4 font-semibold">Status</th>
                                <th class="py-3 px-4 font-semibold text-right">Timestamp</th>
                            </tr>
                        </thead>
                        <tbody class="divide-y divide-surface-800/60 font-mono text-[11px]">
                            <tr
                                v-for="log in logs.data"
                                :key="log.id"
                                class="hover:bg-surface-800/30 transition"
                            >
                                <td class="py-3.5 px-4 font-medium text-surface-200">
                                    {{ log.action }}
                                </td>
                                <td class="py-3.5 px-4 text-surface-400 font-sans">
                                    {{ log.user_email || 'system' }}
                                </td>
                                <td class="py-3.5 px-4 text-surface-500">
                                    {{ log.ip_address }}
                                </td>
                                <td class="py-3.5 px-4 font-sans">
                                    <span
                                        :class="[
                                            log.status === 'success'
                                                ? 'bg-emerald-500/10 border-emerald-500/30 text-emerald-400'
                                                : 'bg-rose-500/10 border-rose-500/30 text-rose-400',
                                            'px-2 py-0.5 rounded-full text-[10px] font-medium border'
                                        ]"
                                    >
                                        {{ log.status }}
                                    </span>
                                </td>
                                <td class="py-3.5 px-4 text-surface-500 text-right">
                                    {{ new Date(log.created_at).toLocaleString() }}
                                </td>
                            </tr>
                        </tbody>
                    </table>
                </div>
            </div>
        </div>
    </AppLayout>
</template>
