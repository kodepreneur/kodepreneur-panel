<script setup lang="ts">
import { useForm, Link } from '@inertiajs/vue3';
import AppLayout from '@/Layouts/AppLayout.vue';
import { Settings, Save, ArrowUpCircle, ChevronRight, Sparkles } from 'lucide-vue-next';

const props = defineProps<{
    settings: Record<string, string>;
}>();

const form = useForm({
    panel_name: props.settings.panel_name || 'Kodepreneur Panel',
    panel_port: props.settings.panel_port || '8443',
    default_php_version: props.settings.default_php_version || '8.3',
});

function submit() {
    form.post('/settings');
}
</script>

<template>
    <AppLayout title="Panel Settings">
        <div class="max-w-3xl mx-auto space-y-6">
            <!-- Header -->
            <div>
                <h2 class="text-base font-bold text-slate-900 dark:text-white tracking-tight">Panel Settings</h2>
                <p class="text-xs text-slate-500 dark:text-surface-400 mt-0.5">Configure system defaults, port configurations, and runtime preferences</p>
            </div>

            <!-- System Updates Quick Card -->
            <Link
                href="/updates"
                class="group block rounded-2xl bg-white dark:bg-gradient-to-r dark:from-brand-950/50 dark:via-surface-900/80 dark:to-surface-900/80 border border-slate-200/80 dark:border-brand-500/20 hover:border-brand-500/40 p-5 shadow-sm dark:shadow-xl transition duration-200"
            >
                <div class="flex items-center justify-between">
                    <div class="flex items-center gap-3.5">
                        <div class="w-10 h-10 rounded-xl bg-brand-50 text-brand-600 dark:bg-brand-500/10 dark:text-brand-400 border border-brand-200 dark:border-brand-500/20 flex items-center justify-center group-hover:scale-105 transition">
                            <ArrowUpCircle class="w-5 h-5" />
                        </div>
                        <div>
                            <div class="flex items-center gap-2">
                                <h3 class="text-xs font-semibold text-slate-900 dark:text-white group-hover:text-brand-600 dark:group-hover:text-brand-300 transition">System Updates & Upgrades</h3>
                                <span class="px-2 py-0.5 rounded-full text-[10px] font-mono bg-brand-50 text-brand-700 border border-brand-200 dark:bg-brand-500/20 dark:text-brand-400 dark:border-brand-500/30">v1.0.0</span>
                            </div>
                            <p class="text-[11px] text-slate-500 dark:text-surface-400 mt-0.5">Check latest GitHub commits, release notes, and run 1-click panel updates.</p>
                        </div>
                    </div>
                    <div class="flex items-center gap-1 text-xs text-brand-600 dark:text-brand-400 font-medium">
                        <span>Manage Updates</span>
                        <ChevronRight class="w-4 h-4 group-hover:translate-x-0.5 transition" />
                    </div>
                </div>
            </Link>

            <!-- Settings Form Card -->
            <div class="rounded-2xl bg-white dark:bg-surface-900/60 border border-slate-200/80 dark:border-surface-800/80 p-6 sm:p-8 shadow-sm dark:shadow-xl">
                <form @submit.prevent="submit" class="space-y-5">
                    <div>
                        <label class="block text-xs font-medium text-slate-700 dark:text-surface-300 mb-1">Control Panel Name</label>
                        <input
                            v-model="form.panel_name"
                            type="text"
                            required
                            class="w-full px-3.5 py-2.5 rounded-xl bg-slate-50 dark:bg-surface-950/60 border border-slate-200 dark:border-surface-800 text-xs text-slate-900 dark:text-white placeholder-slate-400 dark:placeholder-surface-500 focus:outline-none focus:ring-2 focus:ring-brand-500/40 focus:border-brand-500"
                        />
                        <p v-if="form.errors.panel_name" class="text-[11px] text-rose-500 mt-1">{{ form.errors.panel_name }}</p>
                    </div>

                    <div>
                        <label class="block text-xs font-medium text-slate-700 dark:text-surface-300 mb-1">Web Panel HTTPS Port</label>
                        <input
                            v-model="form.panel_port"
                            type="number"
                            required
                            class="w-full px-3.5 py-2.5 rounded-xl bg-slate-50 dark:bg-surface-950/60 border border-slate-200 dark:border-surface-800 text-xs text-slate-900 dark:text-white placeholder-slate-400 dark:placeholder-surface-500 font-mono focus:outline-none focus:ring-2 focus:ring-brand-500/40 focus:border-brand-500"
                        />
                        <p class="text-[11px] text-slate-500 dark:text-surface-500 mt-1">Default is <code class="text-brand-600 dark:text-brand-400 font-mono">8443</code></p>
                        <p v-if="form.errors.panel_port" class="text-[11px] text-rose-500 mt-1">{{ form.errors.panel_port }}</p>
                    </div>

                    <div>
                        <label class="block text-xs font-medium text-slate-700 dark:text-surface-300 mb-1">Default PHP Version</label>
                        <select
                            v-model="form.default_php_version"
                            class="w-full px-3.5 py-2.5 rounded-xl bg-slate-50 dark:bg-surface-950/60 border border-slate-200 dark:border-surface-800 text-xs text-slate-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-brand-500/40 focus:border-brand-500"
                        >
                            <option value="8.3">PHP 8.3 LTS</option>
                            <option value="8.4">PHP 8.4</option>
                        </select>
                        <p v-if="form.errors.default_php_version" class="text-[11px] text-rose-500 mt-1">{{ form.errors.default_php_version }}</p>
                    </div>

                    <div class="flex items-center justify-end pt-4 border-t border-slate-100 dark:border-surface-800">
                        <button
                            type="submit"
                            :disabled="form.processing"
                            class="px-5 py-2.5 rounded-xl bg-brand-600 hover:bg-brand-500 text-white text-xs font-semibold shadow-md shadow-brand-600/20 flex items-center gap-1.5 transition disabled:opacity-50"
                        >
                            <Save class="w-4 h-4" />
                            <span>{{ form.processing ? 'Saving...' : 'Save Settings' }}</span>
                        </button>
                    </div>
                </form>
            </div>
        </div>
    </AppLayout>
</template>
