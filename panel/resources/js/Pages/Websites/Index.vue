<script setup lang="ts">
import { ref, computed } from 'vue';
import { Link, router } from '@inertiajs/vue3';
import AppLayout from '@/Layouts/AppLayout.vue';
import {
    Globe,
    Plus,
    Trash2,
    ShieldCheck,
    ShieldAlert,
    ExternalLink,
    Search,
    Filter,
    Layers,
    Lock,
    Sparkles,
    Activity,
} from 'lucide-vue-next';
import type { Website } from '@/types';

const props = defineProps<{
    websites: {
        data: (Website & { deployments_count: number })[];
        links: any[];
    };
}>();

const searchQuery = ref('');
const filterPhp = ref('all');
const filterSsl = ref('all');
const showDeleteModal = ref(false);
const siteToDelete = ref<Website | null>(null);

const filteredWebsites = computed(() => {
    return props.websites.data.filter((site) => {
        const matchesSearch =
            site.domain.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
            site.system_user.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
            site.document_root.toLowerCase().includes(searchQuery.value.toLowerCase());

        const matchesPhp = filterPhp.value === 'all' || site.php_version === filterPhp.value;

        const matchesSsl =
            filterSsl.value === 'all' ||
            (filterSsl.value === 'ssl' && site.ssl_enabled) ||
            (filterSsl.value === 'http' && !site.ssl_enabled);

        return matchesSearch && matchesPhp && matchesSsl;
    });
});

function confirmDelete(site: Website) {
    siteToDelete.value = site;
    showDeleteModal.value = true;
}

function executeDelete() {
    if (!siteToDelete.value) return;
    router.delete(`/websites/${siteToDelete.value.id}`, {
        onFinish: () => {
            showDeleteModal.value = false;
            siteToDelete.value = null;
        },
    });
}
</script>

<template>
    <AppLayout title="Websites">
        <div class="space-y-6 max-w-7xl mx-auto">
            <!-- Header with Actions -->
            <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
                <div>
                    <h2 class="text-base font-bold text-slate-900 dark:text-white tracking-tight">Managed Websites</h2>
                    <p class="text-xs text-slate-500 dark:text-surface-400 mt-0.5">
                        High-performance Nginx virtual hosts with isolated PHP-FPM pools and Let's Encrypt SSL
                    </p>
                </div>
                <Link
                    href="/websites/create"
                    class="px-4 py-2 rounded-xl bg-brand-600 hover:bg-brand-500 text-white text-xs font-semibold shadow-md shadow-brand-600/20 flex items-center gap-1.5 transition self-start sm:self-auto"
                >
                    <Plus class="w-4 h-4" />
                    <span>Create Website</span>
                </Link>
            </div>

            <!-- Filters Bar -->
            <div class="grid grid-cols-1 sm:grid-cols-3 gap-3">
                <div class="relative">
                    <Search class="w-4 h-4 text-slate-400 dark:text-surface-400 absolute left-3.5 top-1/2 -translate-y-1/2" />
                    <input
                        v-model="searchQuery"
                        type="text"
                        placeholder="Search domains or users..."
                        class="w-full pl-9 pr-4 py-2 rounded-xl bg-white dark:bg-surface-900/60 border border-slate-200/80 dark:border-surface-800 text-xs text-slate-900 dark:text-white placeholder-slate-400 dark:placeholder-surface-500 focus:outline-none focus:ring-2 focus:ring-brand-500/40 focus:border-brand-500 transition"
                    />
                </div>

                <div class="flex items-center gap-2">
                    <select
                        v-model="filterPhp"
                        class="w-full px-3 py-2 rounded-xl bg-white dark:bg-surface-900/60 border border-slate-200/80 dark:border-surface-800 text-xs text-slate-700 dark:text-surface-300 focus:outline-none focus:ring-2 focus:ring-brand-500/40 transition"
                    >
                        <option value="all">All PHP Runtimes</option>
                        <option value="8.3">PHP 8.3 LTS</option>
                        <option value="8.4">PHP 8.4</option>
                        <option value="none">Static (No PHP)</option>
                    </select>
                </div>

                <div class="flex items-center gap-2">
                    <select
                        v-model="filterSsl"
                        class="w-full px-3 py-2 rounded-xl bg-white dark:bg-surface-900/60 border border-slate-200/80 dark:border-surface-800 text-xs text-slate-700 dark:text-surface-300 focus:outline-none focus:ring-2 focus:ring-brand-500/40 transition"
                    >
                        <option value="all">All SSL States</option>
                        <option value="ssl">HTTPS Active</option>
                        <option value="http">HTTP Only</option>
                    </select>
                </div>
            </div>

            <!-- Websites Table Card -->
            <div class="rounded-2xl bg-white dark:bg-surface-900/60 border border-slate-200/80 dark:border-surface-800/80 shadow-sm dark:shadow-xl overflow-hidden">
                <div v-if="filteredWebsites.length === 0" class="text-center py-14 px-4">
                    <Globe class="w-10 h-10 text-slate-300 dark:text-surface-600 mx-auto mb-3" />
                    <h3 class="text-sm font-medium text-slate-700 dark:text-surface-200">No websites found</h3>
                    <p class="text-xs text-slate-500 dark:text-surface-400 mt-1 max-w-sm mx-auto">
                        {{ searchQuery ? 'Try adjusting your search filters.' : 'Create your first website to deploy applications with PHP-FPM and SSL.' }}
                    </p>
                    <Link
                        v-if="!searchQuery"
                        href="/websites/create"
                        class="mt-4 inline-flex items-center gap-1.5 px-3.5 py-2 rounded-xl bg-brand-600 text-white text-xs font-semibold hover:bg-brand-500 transition shadow-md shadow-brand-600/20"
                    >
                        <Plus class="w-3.5 h-3.5" />
                        <span>Provision Website</span>
                    </Link>
                </div>

                <div v-else class="overflow-x-auto">
                    <table class="w-full text-left text-xs">
                        <thead class="bg-slate-50 dark:bg-surface-950/50 text-slate-500 dark:text-surface-400 uppercase text-[10px] tracking-wider border-b border-slate-200/80 dark:border-surface-800">
                            <tr>
                                <th class="py-3.5 px-4 font-semibold">Domain & Isolation</th>
                                <th class="py-3.5 px-4 font-semibold">PHP Runtime</th>
                                <th class="py-3.5 px-4 font-semibold">Document Root</th>
                                <th class="py-3.5 px-4 font-semibold">SSL Certificate</th>
                                <th class="py-3.5 px-4 font-semibold">Status</th>
                                <th class="py-3.5 px-4 font-semibold text-right">Actions</th>
                            </tr>
                        </thead>
                        <tbody class="divide-y divide-slate-100 dark:divide-surface-800/60">
                            <tr
                                v-for="site in filteredWebsites"
                                :key="site.id"
                                class="hover:bg-slate-50/80 dark:hover:bg-surface-800/30 transition group"
                            >
                                <td class="py-3.5 px-4">
                                    <div class="flex items-center gap-2.5">
                                        <div class="p-2 rounded-xl bg-indigo-50 text-brand-600 dark:bg-brand-500/10 dark:border dark:border-brand-500/20 dark:text-brand-400 shrink-0">
                                            <Globe class="w-4 h-4" />
                                        </div>
                                        <div>
                                            <Link
                                                :href="`/websites/${site.id}`"
                                                class="font-semibold text-slate-900 dark:text-surface-100 hover:text-brand-600 dark:hover:text-brand-400 transition"
                                            >
                                                {{ site.domain }}
                                            </Link>
                                            <span class="block text-[10px] text-slate-400 dark:text-surface-500 font-mono">
                                                user: {{ site.system_user }}
                                            </span>
                                        </div>
                                    </div>
                                </td>
                                <td class="py-3.5 px-4">
                                    <span
                                        v-if="site.php_version !== 'none'"
                                        class="px-2 py-0.5 rounded-md bg-slate-100 dark:bg-surface-800 border border-slate-200 dark:border-surface-700 text-slate-700 dark:text-surface-200 font-mono text-[11px]"
                                    >
                                        PHP {{ site.php_version }}
                                    </span>
                                    <span
                                        v-else
                                        class="px-2 py-0.5 rounded-md bg-slate-100 dark:bg-surface-800/50 border border-slate-200 dark:border-surface-700/50 text-slate-500 dark:text-surface-400 font-mono text-[11px]"
                                    >
                                        Static
                                    </span>
                                </td>
                                <td class="py-3.5 px-4 font-mono text-[11px] text-slate-500 dark:text-surface-400 truncate max-w-xs">
                                    {{ site.document_root }}
                                </td>
                                <td class="py-3.5 px-4">
                                    <span
                                        v-if="site.ssl_enabled"
                                        class="inline-flex items-center gap-1 text-[11px] text-emerald-600 dark:text-emerald-400 font-medium"
                                    >
                                        <ShieldCheck class="w-3.5 h-3.5" />
                                        <span>Let's Encrypt</span>
                                    </span>
                                    <span
                                        v-else
                                        class="inline-flex items-center gap-1 text-[11px] text-amber-600 dark:text-amber-400 font-medium"
                                    >
                                        <ShieldAlert class="w-3.5 h-3.5" />
                                        <span>HTTP Only</span>
                                    </span>
                                </td>
                                <td class="py-3.5 px-4">
                                    <span class="px-2.5 py-0.5 rounded-full text-[10px] font-medium bg-emerald-50 text-emerald-600 border border-emerald-200 dark:bg-emerald-500/10 dark:border-emerald-500/30 dark:text-emerald-400">
                                        {{ site.status }}
                                    </span>
                                </td>
                                <td class="py-3.5 px-4 text-right">
                                    <div class="flex items-center justify-end gap-1.5">
                                        <Link
                                            :href="`/websites/${site.id}`"
                                            class="px-2.5 py-1.5 rounded-lg bg-slate-100 hover:bg-slate-200 dark:bg-surface-800 dark:hover:bg-surface-700 text-slate-700 dark:text-surface-200 text-xs font-medium flex items-center gap-1 transition"
                                            title="Manage Website"
                                        >
                                            <Activity class="w-3 h-3 text-indigo-500 dark:text-brand-400" />
                                            <span>Manage</span>
                                            <ExternalLink class="w-3 h-3 text-slate-400 dark:text-surface-400" />
                                        </Link>
                                        <button
                                            @click="confirmDelete(site)"
                                            class="p-1.5 rounded-lg text-rose-500 hover:bg-rose-50 dark:hover:bg-rose-500/10 transition"
                                            title="Delete Website"
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

            <!-- Delete Confirmation Modal -->
            <div
                v-if="showDeleteModal"
                class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/60 dark:bg-surface-950/80 backdrop-blur-sm"
            >
                <div class="w-full max-w-md rounded-2xl bg-white dark:bg-surface-900 border border-slate-200/90 dark:border-surface-800 p-6 shadow-2xl space-y-4">
                    <div class="flex items-center gap-3">
                        <div class="p-2.5 rounded-xl bg-rose-50 dark:bg-rose-500/10 border border-rose-200 dark:border-rose-500/20 text-rose-600 dark:text-rose-400">
                            <Trash2 class="w-5 h-5" />
                        </div>
                        <div>
                            <h3 class="text-sm font-semibold text-slate-900 dark:text-white">Delete Website</h3>
                            <p class="text-xs text-slate-500 dark:text-surface-400 mt-0.5">
                                Are you sure you want to delete <strong class="text-slate-900 dark:text-white">{{ siteToDelete?.domain }}</strong>?
                            </p>
                        </div>
                    </div>

                    <p class="text-xs text-slate-600 dark:text-surface-400 bg-slate-50 dark:bg-surface-950/60 p-3 rounded-xl border border-slate-200/80 dark:border-surface-800/80 leading-relaxed">
                        This will remove the Nginx virtual host, unbind the PHP-FPM socket pool, and remove the configuration from the server.
                    </p>

                    <div class="flex items-center justify-end gap-2.5 pt-2">
                        <button
                            @click="showDeleteModal = false"
                            class="px-4 py-2 rounded-xl text-xs font-medium text-slate-600 hover:text-slate-900 dark:text-surface-400 dark:hover:text-white transition"
                        >
                            Cancel
                        </button>
                        <button
                            @click="executeDelete"
                            class="px-4 py-2 rounded-xl bg-rose-600 hover:bg-rose-500 text-white text-xs font-semibold shadow-md shadow-rose-600/20 transition"
                        >
                            Delete Website
                        </button>
                    </div>
                </div>
            </div>
        </div>
    </AppLayout>
</template>
