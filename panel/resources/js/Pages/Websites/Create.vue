<script setup lang="ts">
import { useForm, Link } from '@inertiajs/vue3';
import AppLayout from '@/Layouts/AppLayout.vue';
import {
    Globe,
    ArrowLeft,
    Plus,
    Shield,
    Cpu,
    Folder,
    Sparkles,
    CheckCircle2,
} from 'lucide-vue-next';

const form = useForm({
    domain: '',
    php_version: '8.3',
    document_root: '',
    auto_ssl: true,
    ssl_email: '',
});

function updateDocRoot() {
    if (form.domain && !form.document_root) {
        form.document_root = `/var/www/${form.domain.toLowerCase().trim()}/public`;
    }
}

function submit() {
    form.post('/websites');
}
</script>

<template>
    <AppLayout title="Create Website">
        <div class="max-w-3xl mx-auto space-y-6">
            <!-- Header -->
            <div class="flex items-center gap-3">
                <Link
                    href="/websites"
                    class="p-2 rounded-xl bg-white dark:bg-surface-900 border border-slate-200/80 dark:border-surface-800 text-slate-500 hover:text-slate-900 dark:text-surface-400 dark:hover:text-white transition shadow-sm"
                >
                    <ArrowLeft class="w-4 h-4" />
                </Link>
                <div>
                    <h2 class="text-base font-bold text-slate-900 dark:text-white tracking-tight">Provision New Virtual Host</h2>
                    <p class="text-xs text-slate-500 dark:text-surface-400">
                        Create an isolated system user, dedicated PHP-FPM socket pool, and Nginx virtual host
                    </p>
                </div>
            </div>

            <!-- Form Card -->
            <div class="rounded-2xl bg-white dark:bg-surface-900/60 border border-slate-200/80 dark:border-surface-800/80 p-6 sm:p-8 shadow-sm dark:shadow-xl">
                <form @submit.prevent="submit" class="space-y-6">
                    <!-- Domain Name -->
                    <div>
                        <label class="block text-xs font-semibold text-slate-700 dark:text-surface-200 mb-1.5 flex items-center justify-between">
                            <span>Primary Domain Name <span class="text-rose-500">*</span></span>
                            <span class="text-[11px] text-slate-400 dark:text-surface-500 font-normal">e.g. example.com or app.mydomain.io</span>
                        </label>
                        <div class="relative">
                            <Globe class="w-4 h-4 text-slate-400 dark:text-surface-500 absolute left-3.5 top-1/2 -translate-y-1/2" />
                            <input
                                v-model="form.domain"
                                @input="updateDocRoot"
                                @blur="updateDocRoot"
                                type="text"
                                required
                                placeholder="example.com"
                                class="w-full pl-10 pr-4 py-2.5 rounded-xl bg-slate-50 dark:bg-surface-950/60 border border-slate-200 dark:border-surface-800 text-xs text-slate-900 dark:text-white placeholder-slate-400 dark:placeholder-surface-600 focus:outline-none focus:ring-2 focus:ring-brand-500/40 focus:border-brand-500 transition font-mono"
                            />
                        </div>
                        <p v-if="form.errors.domain" class="text-[11px] text-rose-500 mt-1.5">{{ form.errors.domain }}</p>
                    </div>

                    <!-- PHP Runtime Selector (Card Grid) -->
                    <div>
                        <label class="block text-xs font-semibold text-slate-700 dark:text-surface-200 mb-2">
                            PHP Runtime Environment
                        </label>
                        <div class="grid grid-cols-1 sm:grid-cols-3 gap-3">
                            <!-- PHP 8.3 -->
                            <label
                                :class="[
                                    'cursor-pointer rounded-xl border p-3.5 flex flex-col justify-between transition',
                                    form.php_version === '8.3'
                                        ? 'bg-brand-50/80 border-brand-500 text-slate-900 dark:bg-brand-500/10 dark:text-white ring-1 ring-brand-500'
                                        : 'bg-slate-50 dark:bg-surface-950/40 border-slate-200/80 dark:border-surface-800 text-slate-600 dark:text-surface-400 hover:border-slate-300 dark:hover:border-surface-700'
                                ]"
                            >
                                <div class="flex items-center justify-between">
                                    <span class="text-xs font-bold font-mono text-slate-900 dark:text-white">PHP 8.3 LTS</span>
                                    <input type="radio" v-model="form.php_version" value="8.3" class="sr-only" />
                                    <CheckCircle2 v-if="form.php_version === '8.3'" class="w-4 h-4 text-brand-600 dark:text-brand-400" />
                                </div>
                                <p class="text-[10px] text-slate-500 dark:text-surface-400 mt-2">Recommended default. Long-term support.</p>
                            </label>

                            <!-- PHP 8.4 -->
                            <label
                                :class="[
                                    'cursor-pointer rounded-xl border p-3.5 flex flex-col justify-between transition',
                                    form.php_version === '8.4'
                                        ? 'bg-brand-50/80 border-brand-500 text-slate-900 dark:bg-brand-500/10 dark:text-white ring-1 ring-brand-500'
                                        : 'bg-slate-50 dark:bg-surface-950/40 border-slate-200/80 dark:border-surface-800 text-slate-600 dark:text-surface-400 hover:border-slate-300 dark:hover:border-surface-700'
                                ]"
                            >
                                <div class="flex items-center justify-between">
                                    <span class="text-xs font-bold font-mono text-slate-900 dark:text-white">PHP 8.4 Latest</span>
                                    <input type="radio" v-model="form.php_version" value="8.4" class="sr-only" />
                                    <CheckCircle2 v-if="form.php_version === '8.4'" class="w-4 h-4 text-brand-600 dark:text-brand-400" />
                                </div>
                                <p class="text-[10px] text-slate-500 dark:text-surface-400 mt-2">Latest release with property hooks & JIT.</p>
                            </label>

                            <!-- Static / Reverse Proxy -->
                            <label
                                :class="[
                                    'cursor-pointer rounded-xl border p-3.5 flex flex-col justify-between transition',
                                    form.php_version === 'none'
                                        ? 'bg-brand-50/80 border-brand-500 text-slate-900 dark:bg-brand-500/10 dark:text-white ring-1 ring-brand-500'
                                        : 'bg-slate-50 dark:bg-surface-950/40 border-slate-200/80 dark:border-surface-800 text-slate-600 dark:text-surface-400 hover:border-slate-300 dark:hover:border-surface-700'
                                ]"
                            >
                                <div class="flex items-center justify-between">
                                    <span class="text-xs font-bold text-slate-900 dark:text-white">Static / Proxy</span>
                                    <input type="radio" v-model="form.php_version" value="none" class="sr-only" />
                                    <CheckCircle2 v-if="form.php_version === 'none'" class="w-4 h-4 text-brand-600 dark:text-brand-400" />
                                </div>
                                <p class="text-[10px] text-slate-500 dark:text-surface-400 mt-2">Static HTML, SPA, or Node.js proxy.</p>
                            </label>
                        </div>
                        <p v-if="form.errors.php_version" class="text-[11px] text-rose-500 mt-1.5">{{ form.errors.php_version }}</p>
                    </div>

                    <!-- Document Root -->
                    <div>
                        <label class="block text-xs font-semibold text-slate-700 dark:text-surface-200 mb-1.5">
                            Document Root (Web Directory)
                        </label>
                        <div class="relative">
                            <Folder class="w-4 h-4 text-slate-400 dark:text-surface-500 absolute left-3.5 top-1/2 -translate-y-1/2" />
                            <input
                                v-model="form.document_root"
                                type="text"
                                placeholder="/var/www/example.com/public"
                                class="w-full pl-10 pr-4 py-2.5 rounded-xl bg-slate-50 dark:bg-surface-950/60 border border-slate-200 dark:border-surface-800 text-xs text-slate-900 dark:text-white placeholder-slate-400 dark:placeholder-surface-600 font-mono focus:outline-none focus:ring-2 focus:ring-brand-500/40 focus:border-brand-500 transition"
                            />
                        </div>
                        <p class="text-[11px] text-slate-500 dark:text-surface-500 mt-1.5">
                            Point to Laravel's <code class="text-brand-600 dark:text-brand-400 font-mono">/public</code> directory or custom application webroot.
                        </p>
                        <p v-if="form.errors.document_root" class="text-[11px] text-rose-500 mt-1.5">{{ form.errors.document_root }}</p>
                    </div>

                    <!-- Automated SSL Option -->
                    <div class="p-4 rounded-xl bg-slate-50 dark:bg-surface-950/40 border border-slate-200/80 dark:border-surface-800/80 space-y-3">
                        <label class="flex items-start gap-3 cursor-pointer">
                            <input
                                v-model="form.auto_ssl"
                                type="checkbox"
                                class="mt-0.5 rounded border-slate-300 dark:border-surface-700 bg-white dark:bg-surface-900 text-brand-600 focus:ring-brand-500"
                            />
                            <div>
                                <span class="text-xs font-semibold text-slate-900 dark:text-white flex items-center gap-1.5">
                                    <Shield class="w-3.5 h-3.5 text-emerald-600 dark:text-emerald-400" />
                                    <span>Issue Free Let's Encrypt SSL Certificate</span>
                                </span>
                                <p class="text-[11px] text-slate-500 dark:text-surface-400 mt-0.5">
                                    Automatically executes ACME challenge and enables HTTPS redirection if DNS resolves to this server.
                                </p>
                            </div>
                        </label>

                        <div v-if="form.auto_ssl" class="pt-2">
                            <label class="block text-[11px] font-medium text-slate-600 dark:text-surface-400 mb-1">
                                Let's Encrypt Notification Email (Optional)
                            </label>
                            <input
                                v-model="form.ssl_email"
                                type="email"
                                placeholder="admin@example.com"
                                class="w-full px-3 py-2 rounded-xl bg-white dark:bg-surface-900 border border-slate-200 dark:border-surface-800 text-xs text-slate-900 dark:text-white placeholder-slate-400 dark:placeholder-surface-600 focus:outline-none focus:ring-2 focus:ring-brand-500/40 transition"
                            />
                        </div>
                    </div>

                    <!-- Submit Buttons -->
                    <div class="flex items-center justify-end gap-3 pt-4 border-t border-slate-100 dark:border-surface-800">
                        <Link
                            href="/websites"
                            class="px-4 py-2 rounded-xl text-xs font-medium text-slate-600 hover:text-slate-900 dark:text-surface-400 dark:hover:text-white transition"
                        >
                            Cancel
                        </Link>
                        <button
                            type="submit"
                            :disabled="form.processing"
                            class="px-5 py-2.5 rounded-xl bg-brand-600 hover:bg-brand-500 text-white text-xs font-semibold shadow-md shadow-brand-600/20 flex items-center gap-1.5 transition disabled:opacity-50"
                        >
                            <Plus class="w-4 h-4" />
                            <span>{{ form.processing ? 'Provisioning Server...' : 'Provision Website' }}</span>
                        </button>
                    </div>
                </form>
            </div>
        </div>
    </AppLayout>
</template>
