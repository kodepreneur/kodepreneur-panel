<script setup lang="ts">
import { ref } from 'vue';
import { useForm, router } from '@inertiajs/vue3';
import AppLayout from '@/Layouts/AppLayout.vue';
import {
    Clock,
    Plus,
    Trash2,
    CheckCircle2,
    XCircle,
    Globe,
    Terminal,
    Sparkles,
} from 'lucide-vue-next';
import type { CronJob, Website } from '@/types';

const props = defineProps<{
    cronJobs: {
        data: CronJob[];
        links: any[];
    };
    websites: Website[];
}>();

const isCreateModalOpen = ref(false);

const cronForm = useForm({
    website_id: '',
    schedule: '* * * * *',
    command: 'php artisan schedule:run',
    system_user: 'kodepreneur',
});

const presets = [
    { label: 'Every minute', value: '* * * * *' },
    { label: 'Every 5 min', value: '*/5 * * * *' },
    { label: 'Every 15 min', value: '*/15 * * * *' },
    { label: 'Hourly', value: '0 * * * *' },
    { label: 'Daily (Midnight)', value: '0 0 * * *' },
    { label: 'Weekly (Sunday)', value: '0 0 * * 0' },
];

function selectPreset(val: string) {
    cronForm.schedule = val;
}

function onWebsiteChange() {
    const site = props.websites.find(w => String(w.id) === String(cronForm.website_id));
    if (site) {
        cronForm.system_user = site.system_user || 'kodepreneur';
        cronForm.command = `cd /var/www/${site.domain} && php artisan schedule:run >> /dev/null 2>&1`;
    }
}

function submitCreateCron() {
    cronForm.post('/cron', {
        onSuccess: () => {
            isCreateModalOpen.value = false;
            cronForm.reset();
        },
    });
}

function toggleJob(job: CronJob) {
    router.post(`/cron/${job.id}/toggle`);
}

function deleteJob(job: CronJob) {
    if (confirm(`Are you sure you want to delete this cron job?`)) {
        router.delete(`/cron/${job.id}`);
    }
}
</script>

<template>
    <AppLayout title="Cron Jobs">
        <div class="max-w-7xl mx-auto space-y-6">
            <!-- Header -->
            <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
                <div>
                    <h2 class="text-base font-semibold text-white tracking-tight">Cron Schedule Manager</h2>
                    <p class="text-xs text-surface-400 mt-0.5">
                        Automate periodic server tasks, Laravel scheduler runs, and maintenance scripts
                    </p>
                </div>
                <button
                    @click="isCreateModalOpen = true"
                    class="px-4 py-2 rounded-xl bg-brand-600 hover:bg-brand-500 text-white text-xs font-semibold shadow-lg shadow-brand-600/20 flex items-center gap-1.5 transition"
                >
                    <Plus class="w-4 h-4" />
                    <span>Create Cron Job</span>
                </button>
            </div>

            <!-- Table Card -->
            <div class="rounded-2xl bg-surface-900/60 border border-surface-800/80 shadow-lg overflow-hidden">
                <div v-if="cronJobs.data.length === 0" class="text-center py-14 px-4">
                    <Clock class="w-10 h-10 text-surface-600 mx-auto mb-3" />
                    <h3 class="text-sm font-medium text-surface-200">No scheduled cron jobs</h3>
                    <p class="text-xs text-surface-400 mt-1 max-w-sm mx-auto">
                        Add scheduled background tasks with automated Linux crontab synchronization.
                    </p>
                    <button
                        @click="isCreateModalOpen = true"
                        class="mt-4 inline-flex items-center gap-1.5 px-4 py-2 rounded-xl bg-brand-600 text-white text-xs font-semibold hover:bg-brand-500 transition shadow-lg shadow-brand-600/20"
                    >
                        <Plus class="w-3.5 h-3.5" />
                        <span>Add First Cron Job</span>
                    </button>
                </div>

                <div v-else class="overflow-x-auto">
                    <table class="w-full text-left text-xs">
                        <thead class="bg-surface-950/50 text-surface-400 uppercase text-[10px] tracking-wider border-b border-surface-800">
                            <tr>
                                <th class="py-3.5 px-4 font-semibold">Schedule</th>
                                <th class="py-3.5 px-4 font-semibold">Command</th>
                                <th class="py-3.5 px-4 font-semibold">Target / User</th>
                                <th class="py-3.5 px-4 font-semibold">Status</th>
                                <th class="py-3.5 px-4 font-semibold text-right">Actions</th>
                            </tr>
                        </thead>
                        <tbody class="divide-y divide-surface-800/60">
                            <tr
                                v-for="job in cronJobs.data"
                                :key="job.id"
                                class="hover:bg-surface-800/30 transition"
                            >
                                <td class="py-3.5 px-4 font-mono font-semibold text-emerald-400 text-xs">
                                    <div class="flex items-center gap-2">
                                        <div class="p-1.5 rounded-lg bg-emerald-500/10 border border-emerald-500/20 text-emerald-400">
                                            <Clock class="w-3.5 h-3.5" />
                                        </div>
                                        <span>{{ job.schedule }}</span>
                                    </div>
                                </td>
                                <td class="py-3.5 px-4 font-mono text-[11px] text-surface-300 max-w-md truncate">
                                    {{ job.command }}
                                </td>
                                <td class="py-3.5 px-4">
                                    <div class="text-[11px] text-white">
                                        <span v-if="job.website" class="font-mono text-brand-400 font-semibold">{{ job.website.domain }}</span>
                                        <span v-else class="text-surface-400">System Global</span>
                                    </div>
                                    <div class="text-[10px] font-mono text-surface-500">User: {{ job.system_user }}</div>
                                </td>
                                <td class="py-3.5 px-4">
                                    <button
                                        @click="toggleJob(job)"
                                        :class="[
                                            'px-2.5 py-0.5 rounded-full text-[10px] font-semibold uppercase tracking-wider transition',
                                            job.is_active
                                                ? 'bg-emerald-500/10 border border-emerald-500/30 text-emerald-400 hover:bg-emerald-500/20'
                                                : 'bg-surface-800 border border-surface-700 text-surface-500 hover:bg-surface-700'
                                        ]"
                                    >
                                        {{ job.is_active ? 'Active' : 'Disabled' }}
                                    </button>
                                </td>
                                <td class="py-3.5 px-4 text-right">
                                    <button
                                        @click="deleteJob(job)"
                                        class="p-1.5 rounded-lg text-rose-400 hover:bg-rose-500/10 transition"
                                        title="Delete Cron Job"
                                    >
                                        <Trash2 class="w-3.5 h-3.5" />
                                    </button>
                                </td>
                            </tr>
                        </tbody>
                    </table>
                </div>
            </div>

            <!-- Create Modal -->
            <div
                v-if="isCreateModalOpen"
                class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-surface-950/80 backdrop-blur-sm"
            >
                <div class="w-full max-w-lg rounded-2xl bg-surface-900 border border-surface-800 p-6 shadow-2xl space-y-4">
                    <h3 class="text-sm font-semibold text-white">Create New Cron Job</h3>

                    <form @submit.prevent="submitCreateCron" class="space-y-4">
                        <div>
                            <label class="block text-xs font-medium text-surface-300 mb-1">Target Website (Optional)</label>
                            <select
                                v-model="cronForm.website_id"
                                @change="onWebsiteChange"
                                class="w-full px-3.5 py-2.5 rounded-xl bg-surface-950/60 border border-surface-800 text-xs text-white focus:outline-none focus:ring-2 focus:ring-brand-500/40"
                            >
                                <option value="">System Default</option>
                                <option v-for="w in websites" :key="w.id" :value="w.id">
                                    {{ w.domain }} ({{ w.system_user }})
                                </option>
                            </select>
                        </div>

                        <div>
                            <div class="flex items-center justify-between mb-1">
                                <label class="text-xs font-medium text-surface-300">Cron Schedule (5-field)</label>
                                <span class="text-[10px] font-mono text-surface-500">* * * * *</span>
                            </div>
                            <div class="flex flex-wrap gap-1.5 mb-2">
                                <button
                                    v-for="p in presets"
                                    :key="p.value"
                                    type="button"
                                    @click="selectPreset(p.value)"
                                    :class="[
                                        'px-2 py-0.5 rounded text-[10px] font-mono transition',
                                        cronForm.schedule === p.value
                                            ? 'bg-brand-500/20 border border-brand-500/40 text-brand-300'
                                            : 'bg-surface-800 text-surface-400 hover:text-white'
                                    ]"
                                >
                                    {{ p.label }}
                                </button>
                            </div>
                            <input
                                v-model="cronForm.schedule"
                                type="text"
                                required
                                placeholder="* * * * *"
                                class="w-full px-3.5 py-2.5 rounded-xl bg-surface-950/60 border border-surface-800 text-xs text-white placeholder-surface-500 font-mono focus:outline-none focus:ring-2 focus:ring-brand-500/40"
                            />
                        </div>

                        <div>
                            <label class="block text-xs font-medium text-surface-300 mb-1">Command to Execute</label>
                            <input
                                v-model="cronForm.command"
                                type="text"
                                required
                                placeholder="php /var/www/app/artisan schedule:run >> /dev/null 2>&1"
                                class="w-full px-3.5 py-2.5 rounded-xl bg-surface-950/60 border border-surface-800 text-xs text-white placeholder-surface-500 font-mono focus:outline-none focus:ring-2 focus:ring-brand-500/40"
                            />
                        </div>

                        <div class="flex items-center justify-end gap-3 pt-3 border-t border-surface-800">
                            <button
                                type="button"
                                @click="isCreateModalOpen = false"
                                class="px-3.5 py-2 rounded-xl text-xs text-surface-400 hover:text-white"
                            >
                                Cancel
                            </button>
                            <button
                                type="submit"
                                :disabled="cronForm.processing"
                                class="px-4 py-2 rounded-xl bg-brand-600 hover:bg-brand-500 text-white text-xs font-semibold shadow-lg shadow-brand-600/20 disabled:opacity-50"
                            >
                                {{ cronForm.processing ? 'Creating...' : 'Create Cron Job' }}
                            </button>
                        </div>
                    </form>
                </div>
            </div>
        </div>
    </AppLayout>
</template>
