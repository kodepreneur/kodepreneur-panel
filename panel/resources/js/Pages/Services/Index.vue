<script setup lang="ts">
import { router } from '@inertiajs/vue3';
import AppLayout from '@/Layouts/AppLayout.vue';
import { Cpu, RefreshCw, Power, RotateCw, CheckCircle2, Shield } from 'lucide-vue-next';
import type { ServiceStatus } from '@/types';

defineProps<{
    services: ServiceStatus[];
}>();

function handleAction(service: string, action: string) {
    if (confirm(`Are you sure you want to execute '${action}' on service '${service}'?`)) {
        router.post(`/services/${service}/${action}`);
    }
}
</script>

<template>
    <AppLayout title="System Services">
        <div class="space-y-6 max-w-7xl mx-auto">
            <!-- Header -->
            <div>
                <h2 class="text-base font-bold text-slate-900 dark:text-white tracking-tight">Systemd Services</h2>
                <p class="text-xs text-slate-500 dark:text-surface-400 mt-0.5">Control web servers, database daemons, and PHP-FPM pools</p>
            </div>

            <!-- Services Grid -->
            <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                <div
                    v-for="svc in services"
                    :key="svc.name"
                    class="rounded-2xl bg-white dark:bg-surface-900/60 border border-slate-200/80 dark:border-surface-800/80 p-5 shadow-sm dark:shadow-xl space-y-4"
                >
                    <div class="flex items-center justify-between">
                        <div class="flex items-center gap-2.5">
                            <div class="w-8 h-8 rounded-xl bg-indigo-50 text-brand-600 dark:bg-surface-800 dark:border dark:border-surface-700 dark:text-brand-400 flex items-center justify-center font-mono text-xs">
                                <Cpu class="w-4 h-4" />
                            </div>
                            <div>
                                <h3 class="text-xs font-semibold text-slate-900 dark:text-white font-mono">{{ svc.name }}</h3>
                                <span class="text-[10px] text-slate-400 dark:text-surface-500 font-mono">systemd unit</span>
                            </div>
                        </div>

                        <span class="flex items-center gap-1.5 px-2 py-0.5 rounded-full text-[10px] font-medium bg-emerald-50 text-emerald-700 border border-emerald-200 dark:bg-emerald-500/10 dark:border-emerald-500/30 dark:text-emerald-400">
                            <span class="w-1.5 h-1.5 rounded-full bg-emerald-500"></span>
                            Active
                        </span>
                    </div>

                    <div class="flex items-center gap-2 pt-2 border-t border-slate-100 dark:border-surface-800/60">
                        <button
                            @click="handleAction(svc.name, 'restart')"
                            class="flex-1 py-1.5 px-3 rounded-xl bg-slate-100 hover:bg-slate-200 dark:bg-surface-800 dark:hover:bg-surface-700 text-slate-700 dark:text-surface-200 text-xs font-medium flex items-center justify-center gap-1.5 transition border border-slate-200/60 dark:border-surface-700"
                        >
                            <RotateCw class="w-3.5 h-3.5" />
                            <span>Restart</span>
                        </button>
                        <button
                            @click="handleAction(svc.name, 'reload')"
                            class="py-1.5 px-3 rounded-xl bg-slate-50 hover:bg-slate-100 dark:bg-surface-800/60 dark:hover:bg-surface-800 text-slate-600 dark:text-surface-300 text-xs font-medium flex items-center justify-center gap-1.5 transition border border-slate-200/60 dark:border-surface-700"
                        >
                            <RefreshCw class="w-3.5 h-3.5" />
                            <span>Reload</span>
                        </button>
                    </div>
                </div>
            </div>
        </div>
    </AppLayout>
</template>
