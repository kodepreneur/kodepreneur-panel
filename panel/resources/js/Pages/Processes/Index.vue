<script setup lang="ts">
import { ref } from 'vue';
import { useForm, router } from '@inertiajs/vue3';
import AppLayout from '@/Layouts/AppLayout.vue';
import {
    Sliders,
    Plus,
    Trash2,
    Play,
    Square,
    RotateCw,
    CheckCircle2,
    AlertCircle,
    Server,
} from 'lucide-vue-next';
import type { Process, Website } from '@/types';

const props = defineProps<{
    processes: {
        data: Process[];
        links: any[];
    };
    websites: Website[];
}>();

const isCreateModalOpen = ref(false);

const procForm = useForm({
    name: '',
    website_id: '',
    command: 'php artisan queue:work --tries=3 --timeout=90',
    system_user: 'kodepreneur',
    instances: 1,
});

function onWebsiteChange() {
    const site = props.websites.find(w => String(w.id) === String(procForm.website_id));
    if (site) {
        procForm.system_user = site.system_user || 'kodepreneur';
        if (!procForm.name) {
            procForm.name = `${site.domain.split('.')[0]}-worker`;
        }
    }
}

function submitCreateProc() {
    procForm.post('/processes', {
        onSuccess: () => {
            isCreateModalOpen.value = false;
            procForm.reset();
        },
    });
}

function manageProc(proc: Process, action: 'start' | 'stop' | 'restart') {
    router.post(`/processes/${proc.id}/${action}`);
}

function deleteProc(proc: Process) {
    if (confirm(`Are you sure you want to stop and delete managed process '${proc.name}'?`)) {
        router.delete(`/processes/${proc.id}`);
    }
}
</script>

<template>
    <AppLayout title="Managed Processes">
        <div class="max-w-7xl mx-auto space-y-6">
            <!-- Header -->
            <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
                <div>
                    <h2 class="text-base font-semibold text-white tracking-tight">Process Supervisor</h2>
                    <p class="text-xs text-surface-400 mt-0.5">
                        Manage background application daemons, queue workers, and Node.js microservices via Systemd
                    </p>
                </div>
                <button
                    @click="isCreateModalOpen = true"
                    class="px-4 py-2 rounded-xl bg-brand-600 hover:bg-brand-500 text-white text-xs font-semibold shadow-lg shadow-brand-600/20 flex items-center gap-1.5 transition"
                >
                    <Plus class="w-4 h-4" />
                    <span>Create Managed Process</span>
                </button>
            </div>

            <!-- Table Card -->
            <div class="rounded-2xl bg-surface-900/60 border border-surface-800/80 shadow-lg overflow-hidden">
                <div v-if="processes.data.length === 0" class="text-center py-14 px-4">
                    <Sliders class="w-10 h-10 text-surface-600 mx-auto mb-3" />
                    <h3 class="text-sm font-medium text-surface-200">No managed processes</h3>
                    <p class="text-xs text-surface-400 mt-1 max-w-sm mx-auto">
                        Deploy automated queue workers or persistent background services with auto-restart policies.
                    </p>
                    <button
                        @click="isCreateModalOpen = true"
                        class="mt-4 inline-flex items-center gap-1.5 px-4 py-2 rounded-xl bg-brand-600 text-white text-xs font-semibold hover:bg-brand-500 transition shadow-lg shadow-brand-600/20"
                    >
                        <Plus class="w-3.5 h-3.5" />
                        <span>Create First Process</span>
                    </button>
                </div>

                <div v-else class="overflow-x-auto">
                    <table class="w-full text-left text-xs">
                        <thead class="bg-surface-950/50 text-surface-400 uppercase text-[10px] tracking-wider border-b border-surface-800">
                            <tr>
                                <th class="py-3.5 px-4 font-semibold">Process Name</th>
                                <th class="py-3.5 px-4 font-semibold">Command</th>
                                <th class="py-3.5 px-4 font-semibold">User / Scope</th>
                                <th class="py-3.5 px-4 font-semibold">Status</th>
                                <th class="py-3.5 px-4 font-semibold text-right">Controls</th>
                            </tr>
                        </thead>
                        <tbody class="divide-y divide-surface-800/60">
                            <tr
                                v-for="p in processes.data"
                                :key="p.id"
                                class="hover:bg-surface-800/30 transition"
                            >
                                <td class="py-3.5 px-4 font-semibold text-white font-mono">
                                    <div class="flex items-center gap-2.5">
                                        <div class="p-1.5 rounded-lg bg-brand-500/10 border border-brand-500/20 text-brand-400">
                                            <Sliders class="w-3.5 h-3.5" />
                                        </div>
                                        <span>{{ p.name }}</span>
                                    </div>
                                </td>
                                <td class="py-3.5 px-4 font-mono text-[11px] text-surface-300 max-w-md truncate">
                                    {{ p.command }}
                                </td>
                                <td class="py-3.5 px-4">
                                    <div class="text-[11px] text-white">
                                        <span v-if="p.website" class="font-mono text-brand-400 font-semibold">{{ p.website.domain }}</span>
                                        <span v-else class="text-surface-400">System</span>
                                    </div>
                                    <div class="text-[10px] font-mono text-surface-500">User: {{ p.system_user }}</div>
                                </td>
                                <td class="py-3.5 px-4">
                                    <span
                                        :class="[
                                            'px-2.5 py-0.5 rounded-full text-[10px] font-semibold uppercase tracking-wider',
                                            p.status === 'running'
                                                ? 'bg-emerald-500/10 border border-emerald-500/30 text-emerald-400'
                                                : 'bg-rose-500/10 border border-rose-500/30 text-rose-400'
                                        ]"
                                    >
                                        {{ p.status }}
                                    </span>
                                </td>
                                <td class="py-3.5 px-4 text-right">
                                    <div class="flex items-center justify-end gap-1.5">
                                        <button
                                            v-if="p.status === 'stopped'"
                                            @click="manageProc(p, 'start')"
                                            class="p-1.5 rounded-lg text-emerald-400 hover:bg-emerald-500/10 transition"
                                            title="Start Process"
                                        >
                                            <Play class="w-3.5 h-3.5" />
                                        </button>
                                        <button
                                            v-if="p.status === 'running'"
                                            @click="manageProc(p, 'restart')"
                                            class="p-1.5 rounded-lg text-amber-400 hover:bg-amber-500/10 transition"
                                            title="Restart Process"
                                        >
                                            <RotateCw class="w-3.5 h-3.5" />
                                        </button>
                                        <button
                                            v-if="p.status === 'running'"
                                            @click="manageProc(p, 'stop')"
                                            class="p-1.5 rounded-lg text-rose-400 hover:bg-rose-500/10 transition"
                                            title="Stop Process"
                                        >
                                            <Square class="w-3.5 h-3.5" />
                                        </button>
                                        <button
                                            @click="deleteProc(p)"
                                            class="p-1.5 rounded-lg text-surface-500 hover:text-rose-400 hover:bg-rose-500/10 transition ml-1"
                                            title="Delete Process"
                                        >
                                            <Trash2 class="w-3.5 h-3.5" />
                                        </button>
                                    </div>
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
                    <h3 class="text-sm font-semibold text-white">Create Managed Process</h3>

                    <form @submit.prevent="submitCreateProc" class="space-y-4">
                        <div>
                            <label class="block text-xs font-medium text-surface-300 mb-1">Process Name</label>
                            <input
                                v-model="procForm.name"
                                type="text"
                                required
                                placeholder="laravel-worker"
                                class="w-full px-3.5 py-2.5 rounded-xl bg-surface-950/60 border border-surface-800 text-xs text-white placeholder-surface-500 font-mono focus:outline-none focus:ring-2 focus:ring-brand-500/40"
                            />
                        </div>

                        <div>
                            <label class="block text-xs font-medium text-surface-300 mb-1">Target Website (Optional)</label>
                            <select
                                v-model="procForm.website_id"
                                @change="onWebsiteChange"
                                class="w-full px-3.5 py-2.5 rounded-xl bg-surface-950/60 border border-surface-800 text-xs text-white focus:outline-none focus:ring-2 focus:ring-brand-500/40"
                            >
                                <option value="">System Default (/var/www)</option>
                                <option v-for="w in websites" :key="w.id" :value="w.id">
                                    {{ w.domain }} ({{ w.system_user }})
                                </option>
                            </select>
                        </div>

                        <div>
                            <label class="block text-xs font-medium text-surface-300 mb-1">Command</label>
                            <input
                                v-model="procForm.command"
                                type="text"
                                required
                                placeholder="php artisan queue:work --tries=3"
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
                                :disabled="procForm.processing"
                                class="px-4 py-2 rounded-xl bg-brand-600 hover:bg-brand-500 text-white text-xs font-semibold shadow-lg shadow-brand-600/20 disabled:opacity-50"
                            >
                                {{ procForm.processing ? 'Creating...' : 'Create & Start' }}
                            </button>
                        </div>
                    </form>
                </div>
            </div>
        </div>
    </AppLayout>
</template>
