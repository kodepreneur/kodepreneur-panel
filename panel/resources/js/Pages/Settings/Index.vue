<script setup lang="ts">
import { useForm } from '@inertiajs/vue3';
import AppLayout from '@/Layouts/AppLayout.vue';
import { Settings, Save, Shield } from 'lucide-vue-next';

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
    <AppLayout title="Settings">
        <div class="max-w-3xl mx-auto space-y-6">
            <!-- Header -->
            <div>
                <h2 class="text-base font-semibold text-white tracking-tight">Panel Settings</h2>
                <p class="text-xs text-surface-400 mt-0.5">Configure system defaults, port configurations, and runtime preferences</p>
            </div>

            <!-- Settings Form Card -->
            <div class="rounded-2xl bg-surface-900/60 border border-surface-800/80 p-6 sm:p-8 shadow-xl">
                <form @submit.prevent="submit" class="space-y-5">
                    <div>
                        <label class="block text-xs font-medium text-surface-300 mb-1">Control Panel Name</label>
                        <input
                            v-model="form.panel_name"
                            type="text"
                            required
                            class="w-full px-3.5 py-2.5 rounded-xl bg-surface-950/60 border border-surface-800 text-xs text-white placeholder-surface-500 focus:outline-none focus:ring-2 focus:ring-brand-500/40 focus:border-brand-500"
                        />
                        <p v-if="form.errors.panel_name" class="text-[11px] text-rose-400 mt-1">{{ form.errors.panel_name }}</p>
                    </div>

                    <div>
                        <label class="block text-xs font-medium text-surface-300 mb-1">Web Panel HTTPS Port</label>
                        <input
                            v-model="form.panel_port"
                            type="number"
                            required
                            class="w-full px-3.5 py-2.5 rounded-xl bg-surface-950/60 border border-surface-800 text-xs text-white placeholder-surface-500 font-mono focus:outline-none focus:ring-2 focus:ring-brand-500/40 focus:border-brand-500"
                        />
                        <p class="text-[11px] text-surface-500 mt-1">Default is <code class="text-brand-400">8443</code></p>
                        <p v-if="form.errors.panel_port" class="text-[11px] text-rose-400 mt-1">{{ form.errors.panel_port }}</p>
                    </div>

                    <div>
                        <label class="block text-xs font-medium text-surface-300 mb-1">Default PHP Version</label>
                        <select
                            v-model="form.default_php_version"
                            class="w-full px-3.5 py-2.5 rounded-xl bg-surface-950/60 border border-surface-800 text-xs text-white focus:outline-none focus:ring-2 focus:ring-brand-500/40 focus:border-brand-500"
                        >
                            <option value="8.3">PHP 8.3 LTS</option>
                            <option value="8.4">PHP 8.4</option>
                        </select>
                        <p v-if="form.errors.default_php_version" class="text-[11px] text-rose-400 mt-1">{{ form.errors.default_php_version }}</p>
                    </div>

                    <div class="flex items-center justify-end pt-4 border-t border-surface-800">
                        <button
                            type="submit"
                            :disabled="form.processing"
                            class="px-5 py-2.5 rounded-xl bg-brand-600 hover:bg-brand-500 text-white text-xs font-semibold shadow-lg shadow-brand-600/20 flex items-center gap-1.5 transition disabled:opacity-50"
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
