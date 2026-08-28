<script setup lang="ts">
import { ref, computed } from 'vue';
import { useForm, Link } from '@inertiajs/vue3';
import AppLayout from '@/Layouts/AppLayout.vue';
import {
    Globe,
    ArrowLeft,
    Plus,
    Shield,
    Folder,
    FolderPlus,
    Sparkles,
    CheckCircle2,
    GitBranch,
    UploadCloud,
    FileArchive,
    Trash2,
    Layers,
    Info,
} from 'lucide-vue-next';

const form = useForm({
    domain: '',
    php_version: '8.3',
    deployment_source: 'empty' as 'empty' | 'zip' | 'git',
    project_type: 'laravel' as 'laravel' | 'generic_php' | 'static',
    document_root: '',
    git_repository: '',
    git_branch: 'main',
    zip_file: null as File | null,
    auto_ssl: true,
    ssl_email: '',
});

const isDragging = ref(false);
const fileInputRef = ref<HTMLInputElement | null>(null);

function updateDocRoot() {
    if (!form.domain) return;
    const cleanDomain = form.domain.toLowerCase().trim();
    if (form.project_type === 'laravel') {
        form.document_root = `/var/www/${cleanDomain}/public`;
    } else {
        form.document_root = `/var/www/${cleanDomain}/public`;
    }
}

function setProjectType(type: 'laravel' | 'generic_php' | 'static') {
    form.project_type = type;
    if (type === 'static') {
        form.php_version = 'none';
    } else if (form.php_version === 'none') {
        form.php_version = '8.3';
    }
    updateDocRoot();
}

function handleFileChange(event: Event) {
    const target = event.target as HTMLInputElement;
    if (target.files && target.files.length > 0) {
        form.zip_file = target.files[0];
    }
}

function handleDrop(event: DragEvent) {
    isDragging.value = false;
    if (event.dataTransfer?.files && event.dataTransfer.files.length > 0) {
        const file = event.dataTransfer.files[0];
        if (file.name.endsWith('.zip') || file.type.includes('zip')) {
            form.zip_file = file;
        }
    }
}

function removeFile() {
    form.zip_file = null;
    if (fileInputRef.value) {
        fileInputRef.value.value = '';
    }
}

function formatBytes(bytes: number): string {
    if (bytes === 0) return '0 Bytes';
    const k = 1024;
    const sizes = ['Bytes', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
}

function submit() {
    form.post('/websites', {
        forceFormData: true,
    });
}
</script>

<template>
    <AppLayout title="Create Website">
        <div class="max-w-3xl mx-auto space-y-6 pb-12">
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
                        Create an isolated project directory, dedicated PHP-FPM socket pool, and Nginx virtual host
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

                    <!-- Deployment Source Selector -->
                    <div>
                        <label class="block text-xs font-semibold text-slate-700 dark:text-surface-200 mb-2 flex items-center justify-between">
                            <span>Deployment Source <span class="text-rose-500">*</span></span>
                            <span class="text-[11px] text-slate-400 dark:text-surface-500 font-normal">Choose how project files are provisioned</span>
                        </label>
                        <div class="grid grid-cols-1 sm:grid-cols-3 gap-3">
                            <!-- 1. Start Empty -->
                            <label
                                :class="[
                                    'cursor-pointer rounded-xl border p-3.5 flex flex-col justify-between transition relative overflow-hidden',
                                    form.deployment_source === 'empty'
                                        ? 'bg-brand-50/80 border-brand-500 text-slate-900 dark:bg-brand-500/10 dark:text-white ring-1 ring-brand-500 shadow-sm'
                                        : 'bg-slate-50 dark:bg-surface-950/40 border-slate-200/80 dark:border-surface-800 text-slate-600 dark:text-surface-400 hover:border-slate-300 dark:hover:border-surface-700'
                                ]"
                            >
                                <div>
                                    <div class="flex items-center justify-between">
                                        <div class="flex items-center gap-2">
                                            <div class="p-1.5 rounded-lg bg-white dark:bg-surface-900 border border-slate-200/80 dark:border-surface-800 text-brand-600 dark:text-brand-400">
                                                <FolderPlus class="w-4 h-4" />
                                            </div>
                                            <span class="text-xs font-bold text-slate-900 dark:text-white">Start Empty</span>
                                        </div>
                                        <input type="radio" v-model="form.deployment_source" value="empty" class="sr-only" />
                                        <CheckCircle2 v-if="form.deployment_source === 'empty'" class="w-4 h-4 text-brand-600 dark:text-brand-400" />
                                    </div>
                                    <p class="text-[11px] text-slate-500 dark:text-surface-400 mt-2.5 leading-relaxed">
                                        Auto-creates directory <code class="font-mono text-slate-700 dark:text-surface-300">/var/www/{{ form.domain || 'domain' }}</code> with default index.
                                    </p>
                                </div>
                            </label>

                            <!-- 2. Upload ZIP -->
                            <label
                                :class="[
                                    'cursor-pointer rounded-xl border p-3.5 flex flex-col justify-between transition relative overflow-hidden',
                                    form.deployment_source === 'zip'
                                        ? 'bg-brand-50/80 border-brand-500 text-slate-900 dark:bg-brand-500/10 dark:text-white ring-1 ring-brand-500 shadow-sm'
                                        : 'bg-slate-50 dark:bg-surface-950/40 border-slate-200/80 dark:border-surface-800 text-slate-600 dark:text-surface-400 hover:border-slate-300 dark:hover:border-surface-700'
                                ]"
                            >
                                <div>
                                    <div class="flex items-center justify-between">
                                        <div class="flex items-center gap-2">
                                            <div class="p-1.5 rounded-lg bg-white dark:bg-surface-900 border border-slate-200/80 dark:border-surface-800 text-brand-600 dark:text-brand-400">
                                                <UploadCloud class="w-4 h-4" />
                                            </div>
                                            <span class="text-xs font-bold text-slate-900 dark:text-white">Upload ZIP</span>
                                        </div>
                                        <input type="radio" v-model="form.deployment_source" value="zip" class="sr-only" />
                                        <CheckCircle2 v-if="form.deployment_source === 'zip'" class="w-4 h-4 text-brand-600 dark:text-brand-400" />
                                    </div>
                                    <p class="text-[11px] text-slate-500 dark:text-surface-400 mt-2.5 leading-relaxed">
                                        Upload and automatically extract project archive into website root folder.
                                    </p>
                                </div>
                            </label>

                            <!-- 3. Git Repo -->
                            <label
                                :class="[
                                    'cursor-pointer rounded-xl border p-3.5 flex flex-col justify-between transition relative overflow-hidden',
                                    form.deployment_source === 'git'
                                        ? 'bg-brand-50/80 border-brand-500 text-slate-900 dark:bg-brand-500/10 dark:text-white ring-1 ring-brand-500 shadow-sm'
                                        : 'bg-slate-50 dark:bg-surface-950/40 border-slate-200/80 dark:border-surface-800 text-slate-600 dark:text-surface-400 hover:border-slate-300 dark:hover:border-surface-700'
                                ]"
                            >
                                <div>
                                    <div class="flex items-center justify-between">
                                        <div class="flex items-center gap-2">
                                            <div class="p-1.5 rounded-lg bg-white dark:bg-surface-900 border border-slate-200/80 dark:border-surface-800 text-brand-600 dark:text-brand-400">
                                                <GitBranch class="w-4 h-4" />
                                            </div>
                                            <span class="text-xs font-bold text-slate-900 dark:text-white">Git Repo</span>
                                        </div>
                                        <input type="radio" v-model="form.deployment_source" value="git" class="sr-only" />
                                        <CheckCircle2 v-if="form.deployment_source === 'git'" class="w-4 h-4 text-brand-600 dark:text-brand-400" />
                                    </div>
                                    <p class="text-[11px] text-slate-500 dark:text-surface-400 mt-2.5 leading-relaxed">
                                        Clone directly from GitHub, GitLab, or remote repository on specified branch.
                                    </p>
                                </div>
                            </label>
                        </div>
                        <p v-if="form.errors.deployment_source" class="text-[11px] text-rose-500 mt-1.5">{{ form.errors.deployment_source }}</p>
                    </div>

                    <!-- Dynamic Details for Upload ZIP -->
                    <div v-if="form.deployment_source === 'zip'" class="space-y-3 p-4 rounded-xl bg-slate-50 dark:bg-surface-950/50 border border-slate-200/80 dark:border-surface-800/80">
                        <label class="block text-xs font-semibold text-slate-700 dark:text-surface-200 flex items-center justify-between">
                            <span>Project Archive (.ZIP) <span class="text-rose-500">*</span></span>
                            <span class="text-[11px] text-slate-400 dark:text-surface-500 font-normal">Max size: 500MB</span>
                        </label>

                        <div
                            v-if="!form.zip_file"
                            @dragover.prevent="isDragging = true"
                            @dragleave.prevent="isDragging = false"
                            @drop.prevent="handleDrop"
                            :class="[
                                'border-2 border-dashed rounded-xl p-6 text-center transition cursor-pointer flex flex-col items-center justify-center gap-2',
                                isDragging
                                    ? 'border-brand-500 bg-brand-50/50 dark:bg-brand-500/10'
                                    : 'border-slate-300 dark:border-surface-700 hover:border-slate-400 dark:hover:border-surface-600 bg-white dark:bg-surface-900/50'
                            ]"
                            @click="fileInputRef?.click()"
                        >
                            <input
                                ref="fileInputRef"
                                type="file"
                                accept=".zip,application/zip"
                                class="hidden"
                                @change="handleFileChange"
                            />
                            <div class="p-3 rounded-full bg-brand-50 dark:bg-brand-500/10 text-brand-600 dark:text-brand-400">
                                <UploadCloud class="w-6 h-6" />
                            </div>
                            <div class="text-xs font-semibold text-slate-800 dark:text-white">
                                Drop your .zip file here, or <span class="text-brand-600 dark:text-brand-400 underline">browse</span>
                            </div>
                            <p class="text-[11px] text-slate-400 dark:text-surface-500">
                                Supports Laravel, WordPress, Vue/React dist, or static web archives
                            </p>
                        </div>

                        <!-- Selected File Preview -->
                        <div v-else class="flex items-center justify-between p-3.5 rounded-xl bg-white dark:bg-surface-900 border border-slate-200 dark:border-surface-700">
                            <div class="flex items-center gap-3 overflow-hidden">
                                <div class="p-2 rounded-lg bg-brand-50 dark:bg-brand-500/10 text-brand-600 dark:text-brand-400 shrink-0">
                                    <FileArchive class="w-5 h-5" />
                                </div>
                                <div class="truncate">
                                    <p class="text-xs font-bold text-slate-900 dark:text-white truncate">{{ form.zip_file.name }}</p>
                                    <p class="text-[11px] text-slate-400 dark:text-surface-400 font-mono">{{ formatBytes(form.zip_file.size) }}</p>
                                </div>
                            </div>
                            <button
                                type="button"
                                @click="removeFile"
                                class="p-1.5 rounded-lg text-slate-400 hover:text-rose-500 hover:bg-rose-50 dark:hover:bg-rose-500/10 transition shrink-0"
                                title="Remove file"
                            >
                                <Trash2 class="w-4 h-4" />
                            </button>
                        </div>
                        <p v-if="form.errors.zip_file" class="text-[11px] text-rose-500 mt-1.5">{{ form.errors.zip_file }}</p>
                    </div>

                    <!-- Dynamic Details for Git Repo -->
                    <div v-if="form.deployment_source === 'git'" class="space-y-4 p-4 rounded-xl bg-slate-50 dark:bg-surface-950/50 border border-slate-200/80 dark:border-surface-800/80">
                        <div class="grid grid-cols-1 sm:grid-cols-3 gap-3">
                            <div class="sm:col-span-2">
                                <label class="block text-xs font-semibold text-slate-700 dark:text-surface-200 mb-1.5">
                                    Repository URL <span class="text-rose-500">*</span>
                                </label>
                                <div class="relative">
                                    <GitBranch class="w-4 h-4 text-slate-400 dark:text-surface-500 absolute left-3.5 top-1/2 -translate-y-1/2" />
                                    <input
                                        v-model="form.git_repository"
                                        type="text"
                                        required
                                        placeholder="https://github.com/organization/repo.git"
                                        class="w-full pl-10 pr-4 py-2.5 rounded-xl bg-white dark:bg-surface-900 border border-slate-200 dark:border-surface-700 text-xs text-slate-900 dark:text-white placeholder-slate-400 dark:placeholder-surface-600 focus:outline-none focus:ring-2 focus:ring-brand-500/40 font-mono transition"
                                    />
                                </div>
                                <p v-if="form.errors.git_repository" class="text-[11px] text-rose-500 mt-1.5">{{ form.errors.git_repository }}</p>
                            </div>

                            <div>
                                <label class="block text-xs font-semibold text-slate-700 dark:text-surface-200 mb-1.5">
                                    Branch
                                </label>
                                <input
                                    v-model="form.git_branch"
                                    type="text"
                                    placeholder="main"
                                    class="w-full px-3.5 py-2.5 rounded-xl bg-white dark:bg-surface-900 border border-slate-200 dark:border-surface-700 text-xs text-slate-900 dark:text-white placeholder-slate-400 dark:placeholder-surface-600 focus:outline-none focus:ring-2 focus:ring-brand-500/40 font-mono transition"
                                />
                                <p v-if="form.errors.git_branch" class="text-[11px] text-rose-500 mt-1.5">{{ form.errors.git_branch }}</p>
                            </div>
                        </div>
                    </div>

                    <!-- Project Framework & Routing Preset -->
                    <div>
                        <label class="block text-xs font-semibold text-slate-700 dark:text-surface-200 mb-2 flex items-center justify-between">
                            <span>Framework & Routing Preset</span>
                            <span class="text-[11px] text-brand-600 dark:text-brand-400 font-medium flex items-center gap-1">
                                <Sparkles class="w-3.5 h-3.5" />
                                Auto-routing enabled
                            </span>
                        </label>
                        <div class="grid grid-cols-1 sm:grid-cols-3 gap-3">
                            <!-- Laravel -->
                            <button
                                type="button"
                                @click="setProjectType('laravel')"
                                :class="[
                                    'text-left rounded-xl border p-3.5 flex flex-col justify-between transition relative',
                                    form.project_type === 'laravel'
                                        ? 'bg-rose-50/80 border-rose-500 text-slate-900 dark:bg-rose-500/10 dark:text-white ring-1 ring-rose-500 shadow-sm'
                                        : 'bg-slate-50 dark:bg-surface-950/40 border-slate-200/80 dark:border-surface-800 text-slate-600 dark:text-surface-400 hover:border-slate-300 dark:hover:border-surface-700'
                                ]"
                            >
                                <div class="flex items-center justify-between w-full">
                                    <span class="text-xs font-bold text-slate-900 dark:text-white flex items-center gap-1.5">
                                        <span class="inline-block w-2 h-2 rounded-full bg-rose-500"></span>
                                        Laravel / Modern PHP
                                    </span>
                                    <CheckCircle2 v-if="form.project_type === 'laravel'" class="w-4 h-4 text-rose-600 dark:text-rose-400" />
                                </div>
                                <p class="text-[10px] text-slate-500 dark:text-surface-400 mt-2">
                                    Root route <code class="font-mono text-rose-600 dark:text-rose-400 font-bold">/</code> auto-directs to <code class="font-mono text-rose-600 dark:text-rose-400">/public</code> folder.
                                </p>
                            </button>

                            <!-- Generic PHP -->
                            <button
                                type="button"
                                @click="setProjectType('generic_php')"
                                :class="[
                                    'text-left rounded-xl border p-3.5 flex flex-col justify-between transition relative',
                                    form.project_type === 'generic_php'
                                        ? 'bg-indigo-50/80 border-indigo-500 text-slate-900 dark:bg-indigo-500/10 dark:text-white ring-1 ring-indigo-500 shadow-sm'
                                        : 'bg-slate-50 dark:bg-surface-950/40 border-slate-200/80 dark:border-surface-800 text-slate-600 dark:text-surface-400 hover:border-slate-300 dark:hover:border-surface-700'
                                ]"
                            >
                                <div class="flex items-center justify-between w-full">
                                    <span class="text-xs font-bold text-slate-900 dark:text-white flex items-center gap-1.5">
                                        <span class="inline-block w-2 h-2 rounded-full bg-indigo-500"></span>
                                        Generic PHP / CMS
                                    </span>
                                    <CheckCircle2 v-if="form.project_type === 'generic_php'" class="w-4 h-4 text-indigo-600 dark:text-indigo-400" />
                                </div>
                                <p class="text-[10px] text-slate-500 dark:text-surface-400 mt-2">
                                    Standard PHP application, WordPress, or custom framework.
                                </p>
                            </button>

                            <!-- Static HTML / SPA -->
                            <button
                                type="button"
                                @click="setProjectType('static')"
                                :class="[
                                    'text-left rounded-xl border p-3.5 flex flex-col justify-between transition relative',
                                    form.project_type === 'static'
                                        ? 'bg-cyan-50/80 border-cyan-500 text-slate-900 dark:bg-cyan-500/10 dark:text-white ring-1 ring-cyan-500 shadow-sm'
                                        : 'bg-slate-50 dark:bg-surface-950/40 border-slate-200/80 dark:border-surface-800 text-slate-600 dark:text-surface-400 hover:border-slate-300 dark:hover:border-surface-700'
                                ]"
                            >
                                <div class="flex items-center justify-between w-full">
                                    <span class="text-xs font-bold text-slate-900 dark:text-white flex items-center gap-1.5">
                                        <span class="inline-block w-2 h-2 rounded-full bg-cyan-500"></span>
                                        Static / SPA
                                    </span>
                                    <CheckCircle2 v-if="form.project_type === 'static'" class="w-4 h-4 text-cyan-600 dark:text-cyan-400" />
                                </div>
                                <p class="text-[10px] text-slate-500 dark:text-surface-400 mt-2">
                                    HTML/CSS/JS frontend or client-side single page app.
                                </p>
                            </button>
                        </div>

                        <!-- Routing Notice -->
                        <div v-if="form.project_type === 'laravel'" class="mt-2.5 p-3 rounded-xl bg-rose-50/50 dark:bg-rose-950/20 border border-rose-200/60 dark:border-rose-900/40 flex items-start gap-2.5">
                            <Info class="w-4 h-4 text-rose-600 dark:text-rose-400 shrink-0 mt-0.5" />
                            <p class="text-[11px] text-rose-900 dark:text-rose-200 leading-relaxed">
                                <span class="font-semibold">Laravel Auto-Routing Active:</span> Nginx virtual host document root is automatically mapped to <code class="font-mono bg-rose-100 dark:bg-rose-900/50 px-1 py-0.5 rounded text-rose-950 dark:text-rose-100">/var/www/{{ form.domain || 'example.com' }}/public</code>. Root route requests seamlessly execute Laravel's <code class="font-mono text-rose-950 dark:text-rose-100">public/index.php</code>.
                            </p>
                        </div>
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
                                        ? 'bg-brand-50/80 border-brand-500 text-slate-900 dark:bg-brand-500/10 dark:text-white ring-1 ring-brand-500 shadow-sm'
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
                                        ? 'bg-brand-50/80 border-brand-500 text-slate-900 dark:bg-brand-500/10 dark:text-white ring-1 ring-brand-500 shadow-sm'
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
                                        ? 'bg-brand-50/80 border-brand-500 text-slate-900 dark:bg-brand-500/10 dark:text-white ring-1 ring-brand-500 shadow-sm'
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
                        <label class="block text-xs font-semibold text-slate-700 dark:text-surface-200 mb-1.5 flex items-center justify-between">
                            <span>Document Root (Web Directory)</span>
                            <span class="text-[11px] text-slate-400 dark:text-surface-500 font-mono">
                                Base: /var/www/{{ form.domain || 'example.com' }}
                            </span>
                        </label>
                        <div class="relative">
                            <Folder class="w-4 h-4 text-slate-400 dark:text-surface-500 absolute left-3.5 top-1/2 -translate-y-1/2" />
                            <input
                                v-model="form.document_root"
                                type="text"
                                :placeholder="`/var/www/${form.domain || 'example.com'}/public`"
                                class="w-full pl-10 pr-4 py-2.5 rounded-xl bg-slate-50 dark:bg-surface-950/60 border border-slate-200 dark:border-surface-800 text-xs text-slate-900 dark:text-white placeholder-slate-400 dark:placeholder-surface-600 font-mono focus:outline-none focus:ring-2 focus:ring-brand-500/40 focus:border-brand-500 transition"
                            />
                        </div>
                        <p class="text-[11px] text-slate-500 dark:text-surface-400 mt-1.5">
                            Points to Laravel's <code class="text-brand-600 dark:text-brand-400 font-mono font-semibold">/public</code> folder or custom webroot directory.
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
