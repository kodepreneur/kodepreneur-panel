<script setup lang="ts">
import { ref } from 'vue';
import { useForm, router } from '@inertiajs/vue3';
import AppLayout from '@/Layouts/AppLayout.vue';
import {
    Shield,
    Plus,
    Trash2,
    Lock,
    Unlock,
    CheckCircle2,
    ShieldAlert,
} from 'lucide-vue-next';
import type { FirewallRule } from '@/types';

const props = defineProps<{
    rules: {
        data: FirewallRule[];
        links: any[];
    };
}>();

const isAddRuleOpen = ref(false);

const ruleForm = useForm({
    port: 80,
    protocol: 'tcp',
    action: 'allow',
    source_ip: '',
    description: '',
});

function submitAddRule() {
    ruleForm.post('/firewall', {
        onSuccess: () => {
            isAddRuleOpen.value = false;
            ruleForm.reset();
        },
    });
}

function deleteRule(rule: FirewallRule) {
    if (confirm(`Delete firewall rule for port ${rule.port}/${rule.protocol}?`)) {
        router.delete(`/firewall/${rule.id}`);
    }
}
</script>

<template>
    <AppLayout title="Firewall">
        <div class="max-w-7xl mx-auto space-y-6">
            <!-- Header -->
            <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
                <div>
                    <h2 class="text-base font-semibold text-white tracking-tight">UFW Firewall Manager</h2>
                    <p class="text-xs text-surface-400 mt-0.5">
                        Manage port access, protocol routing, and restrict incoming connections by IP address
                    </p>
                </div>
                <button
                    @click="isAddRuleOpen = true"
                    class="px-4 py-2 rounded-xl bg-brand-600 hover:bg-brand-500 text-white text-xs font-semibold shadow-lg shadow-brand-600/20 flex items-center gap-1.5 transition"
                >
                    <Plus class="w-4 h-4" />
                    <span>Add Firewall Rule</span>
                </button>
            </div>

            <!-- Table Card -->
            <div class="rounded-2xl bg-surface-900/60 border border-surface-800/80 shadow-lg overflow-hidden">
                <div v-if="rules.data.length === 0" class="text-center py-14 px-4">
                    <Shield class="w-10 h-10 text-surface-600 mx-auto mb-3" />
                    <h3 class="text-sm font-medium text-surface-200">No custom firewall rules</h3>
                    <p class="text-xs text-surface-400 mt-1 max-w-sm mx-auto">
                        Default rules (SSH:22, HTTP:80, HTTPS:443) are active. Add custom port rules for apps and databases.
                    </p>
                    <button
                        @click="isAddRuleOpen = true"
                        class="mt-4 inline-flex items-center gap-1.5 px-4 py-2 rounded-xl bg-brand-600 text-white text-xs font-semibold hover:bg-brand-500 transition shadow-lg shadow-brand-600/20"
                    >
                        <Plus class="w-3.5 h-3.5" />
                        <span>Add Port Rule</span>
                    </button>
                </div>

                <div v-else class="overflow-x-auto">
                    <table class="w-full text-left text-xs">
                        <thead class="bg-surface-950/50 text-surface-400 uppercase text-[10px] tracking-wider border-b border-surface-800">
                            <tr>
                                <th class="py-3.5 px-4 font-semibold">Port</th>
                                <th class="py-3.5 px-4 font-semibold">Protocol</th>
                                <th class="py-3.5 px-4 font-semibold">Action</th>
                                <th class="py-3.5 px-4 font-semibold">Source IP</th>
                                <th class="py-3.5 px-4 font-semibold">Description</th>
                                <th class="py-3.5 px-4 font-semibold text-right">Actions</th>
                            </tr>
                        </thead>
                        <tbody class="divide-y divide-surface-800/60">
                            <tr
                                v-for="r in rules.data"
                                :key="r.id"
                                class="hover:bg-surface-800/30 transition"
                            >
                                <td class="py-3.5 px-4 font-semibold text-white font-mono text-xs">
                                    <div class="flex items-center gap-2">
                                        <div class="p-1.5 rounded-lg bg-emerald-500/10 border border-emerald-500/20 text-emerald-400">
                                            <Shield class="w-3.5 h-3.5" />
                                        </div>
                                        <span>{{ r.port }}</span>
                                    </div>
                                </td>
                                <td class="py-3.5 px-4 font-mono text-[11px] uppercase font-semibold text-surface-300">
                                    {{ r.protocol }}
                                </td>
                                <td class="py-3.5 px-4">
                                    <span
                                        :class="[
                                            'px-2 py-0.5 rounded text-[10px] font-mono uppercase font-semibold',
                                            r.action === 'allow'
                                                ? 'bg-emerald-500/10 border border-emerald-500/30 text-emerald-400'
                                                : 'bg-rose-500/10 border border-rose-500/30 text-rose-400'
                                        ]"
                                    >
                                        {{ r.action }}
                                    </span>
                                </td>
                                <td class="py-3.5 px-4 font-mono text-[11px] text-surface-400">
                                    {{ r.source_ip || 'Anywhere (0.0.0.0/0)' }}
                                </td>
                                <td class="py-3.5 px-4 text-[11px] text-surface-300">
                                    {{ r.description || '—' }}
                                </td>
                                <td class="py-3.5 px-4 text-right">
                                    <button
                                        @click="deleteRule(r)"
                                        class="p-1.5 rounded-lg text-rose-400 hover:bg-rose-500/10 transition"
                                        title="Delete Rule"
                                    >
                                        <Trash2 class="w-3.5 h-3.5" />
                                    </button>
                                </td>
                            </tr>
                        </tbody>
                    </table>
                </div>
            </div>

            <!-- Add Rule Modal -->
            <div
                v-if="isAddRuleOpen"
                class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-surface-950/80 backdrop-blur-sm"
            >
                <div class="w-full max-w-md rounded-2xl bg-surface-900 border border-surface-800 p-6 shadow-2xl space-y-4">
                    <h3 class="text-sm font-semibold text-white">Add Firewall Rule</h3>

                    <form @submit.prevent="submitAddRule" class="space-y-4">
                        <div class="grid grid-cols-2 gap-3">
                            <div>
                                <label class="block text-xs font-medium text-surface-300 mb-1">Port</label>
                                <input
                                    v-model.number="ruleForm.port"
                                    type="number"
                                    required
                                    min="1"
                                    max="65535"
                                    class="w-full px-3.5 py-2.5 rounded-xl bg-surface-950/60 border border-surface-800 text-xs text-white font-mono focus:outline-none focus:ring-2 focus:ring-brand-500/40"
                                />
                            </div>
                            <div>
                                <label class="block text-xs font-medium text-surface-300 mb-1">Protocol</label>
                                <select
                                    v-model="ruleForm.protocol"
                                    class="w-full px-3.5 py-2.5 rounded-xl bg-surface-950/60 border border-surface-800 text-xs text-white focus:outline-none focus:ring-2 focus:ring-brand-500/40"
                                >
                                    <option value="tcp">TCP</option>
                                    <option value="udp">UDP</option>
                                    <option value="both">Both (TCP & UDP)</option>
                                </select>
                            </div>
                        </div>

                        <div>
                            <label class="block text-xs font-medium text-surface-300 mb-1">Action</label>
                            <select
                                v-model="ruleForm.action"
                                class="w-full px-3.5 py-2.5 rounded-xl bg-surface-950/60 border border-surface-800 text-xs text-white focus:outline-none focus:ring-2 focus:ring-brand-500/40"
                            >
                                <option value="allow">ALLOW (Open port)</option>
                                <option value="deny">DENY (Block traffic)</option>
                            </select>
                        </div>

                        <div>
                            <label class="block text-xs font-medium text-surface-300 mb-1">Source IP (Optional)</label>
                            <input
                                v-model="ruleForm.source_ip"
                                type="text"
                                placeholder="e.g. 192.168.1.100 or leave empty for all"
                                class="w-full px-3.5 py-2.5 rounded-xl bg-surface-950/60 border border-surface-800 text-xs text-white placeholder-surface-500 font-mono focus:outline-none focus:ring-2 focus:ring-brand-500/40"
                            />
                        </div>

                        <div>
                            <label class="block text-xs font-medium text-surface-300 mb-1">Description</label>
                            <input
                                v-model="ruleForm.description"
                                type="text"
                                placeholder="Custom API service"
                                class="w-full px-3.5 py-2.5 rounded-xl bg-surface-950/60 border border-surface-800 text-xs text-white placeholder-surface-500 focus:outline-none focus:ring-2 focus:ring-brand-500/40"
                            />
                        </div>

                        <div class="flex items-center justify-end gap-3 pt-3 border-t border-surface-800">
                            <button
                                type="button"
                                @click="isAddRuleOpen = false"
                                class="px-3.5 py-2 rounded-xl text-xs text-surface-400 hover:text-white"
                            >
                                Cancel
                            </button>
                            <button
                                type="submit"
                                :disabled="ruleForm.processing"
                                class="px-4 py-2 rounded-xl bg-brand-600 hover:bg-brand-500 text-white text-xs font-semibold shadow-lg shadow-brand-600/20 disabled:opacity-50"
                            >
                                {{ ruleForm.processing ? 'Saving...' : 'Add Rule' }}
                            </button>
                        </div>
                    </form>
                </div>
            </div>
        </div>
    </AppLayout>
</template>
