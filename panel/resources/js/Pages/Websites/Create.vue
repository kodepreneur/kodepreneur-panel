<script setup lang="ts">
import { ref, computed } from 'vue';
import { useForm, Link } from '@inertiajs/vue3';
import axios from 'axios';
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
    Database,
    Key,
    RefreshCw,
    Eye,
    EyeOff,
    Copy,
    Check,
    Terminal,
    Lock,
    Unlock,
    ExternalLink,
    Code,
} from 'lucide-vue-next';

function generateRandomPassword() {
    const chars = 'abcdefghjkmnpqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ23456789!@#$%&*';
    let pass = '';
    for (let i = 0; i < 16; i++) {
        pass += chars.charAt(Math.floor(Math.random() * chars.length));
    }
    return pass;
}

const showDbPassword = ref(false);
const copiedDbPassword = ref(false);

const isPrivateRepo = ref(false);
const gitAuthType = ref<'ssh_key' | 'token'>('ssh_key');
const isGeneratingKey = ref(false);
const showCustomPrivateKey = ref(false);
const showGitToken = ref(false);
const copiedDeployKey = ref(false);
const keyGenError = ref('');

const form = useForm({
    domain: '',
    php_version: '8.3',
    deployment_source: 'empty' as 'empty' | 'zip' | 'git',
    project_type: 'laravel' as 'laravel' | 'generic_php' | 'static',
    document_root: '',
    git_repository: '',
    git_branch: 'main',
    git_auth_type: 'none' as 'none' | 'ssh_key' | 'token',
    git_token: '',
    git_token_user: '',
    git_ssh_private_key: '',
    git_ssh_public_key: '',
    zip_file: null as File | null,
    auto_ssl: false,
    ssl_email: '',

    // Database Creation
    create_database: true,
    db_engine: 'mysql' as 'mysql' | 'postgresql',
    db_name: '',
    db_username: '',
    db_password: generateRandomPassword(),

    // Laravel Automated Post-Setup
    auto_setup_laravel: true,
    setup_env: true,
    run_composer: true,
    run_key_generate: true,
    run_migrations: true,
    run_seeders: false,
    run_npm_build: true,
    run_optimize: true,
});

async function fetchDeployKey() {
    if (form.git_ssh_public_key && form.git_ssh_private_key) return;
    isGeneratingKey.value = true;
    keyGenError.value = '';
    try {
        const res = await axios.post('/websites/deploy-key/generate');
        if (res.data && res.data.public_key) {
            form.git_ssh_public_key = res.data.public_key;
            form.git_ssh_private_key = res.data.private_key;
        }
    } catch (e: any) {
        keyGenError.value = e.response?.data?.message || 'Failed to generate SSH Deploy Key';
    } finally {
        isGeneratingKey.value = false;
    }
}

async function regenerateDeployKey() {
    isGeneratingKey.value = true;
    keyGenError.value = '';
    try {
        const res = await axios.post('/websites/deploy-key/generate');
        if (res.data && res.data.public_key) {
            form.git_ssh_public_key = res.data.public_key;
            form.git_ssh_private_key = res.data.private_key;
        }
    } catch (e: any) {
        keyGenError.value = e.response?.data?.message || 'Failed to generate SSH Deploy Key';
    } finally {
        isGeneratingKey.value = false;
    }
}

function copyDeployKey() {
    if (!form.git_ssh_public_key) return;
    navigator.clipboard.writeText(form.git_ssh_public_key);
    copiedDeployKey.value = true;
    setTimeout(() => {
        copiedDeployKey.value = false;
    }, 2000);
}

function togglePrivateRepo(isPrivate: boolean) {
    isPrivateRepo.value = isPrivate;
    if (isPrivate) {
        form.git_auth_type = gitAuthType.value;
        if (form.git_auth_type === 'ssh_key') {
            fetchDeployKey();
        }
    } else {
        form.git_auth_type = 'none';
    }
}

function setGitAuthType(type: 'ssh_key' | 'token') {
    gitAuthType.value = type;
    form.git_auth_type = type;
    if (type === 'ssh_key') {
        fetchDeployKey();
    }
}

const isDragging = ref(false);
const fileInputRef = ref<HTMLInputElement | null>(null);

function updateDomainDerivedFields() {
    if (!form.domain) return;
    const cleanDomain = form.domain.toLowerCase().trim();
    const slug = cleanDomain.split('.')[0].replace(/[^a-zA-Z0-9_]/g, '_');

    if (form.project_type === 'laravel') {
        form.document_root = `/var/www/${cleanDomain}/public`;
    } else {
        form.document_root = `/var/www/${cleanDomain}`;
    }

    if (!form.db_name || form.db_name.startsWith('db_')) {
        form.db_name = `db_${slug}`;
    }
    if (!form.db_username || form.db_username.startsWith('u_')) {
        form.db_username = `u_${slug}`;
    }
}

function setProjectType(type: 'laravel' | 'generic_php' | 'static') {
    form.project_type = type;
    if (type === 'static') {
        form.php_version = 'none';
        form.create_database = false;
        form.auto_setup_laravel = false;
    } else {
        if (form.php_version === 'none') {
            form.php_version = '8.3';
        }
        if (type === 'laravel') {
            form.create_database = true;
            form.auto_setup_laravel = true;
        }
    }
    updateDomainDerivedFields();
}

function regeneratePassword() {
    form.db_password = generateRandomPassword();
}

function copyDbPassword() {
    if (!form.db_password) return;
    navigator.clipboard.writeText(form.db_password);
    copiedDbPassword.value = true;
    setTimeout(() => {
        copiedDbPassword.value = false;
    }, 2000);
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
                                @input="updateDomainDerivedFields"
                                @blur="updateDomainDerivedFields"
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
                    <div v-if="form.deployment_source === 'git'" class="space-y-4 p-4 sm:p-5 rounded-2xl bg-slate-50 dark:bg-surface-950/50 border border-slate-200/80 dark:border-surface-800/80">
                        <!-- Repository Access Selector: Public vs Private -->
                        <div>
                            <label class="block text-xs font-semibold text-slate-700 dark:text-surface-200 mb-2 flex items-center justify-between">
                                <span>Repository Access & Privacy <span class="text-rose-500">*</span></span>
                                <span class="text-[11px] text-slate-400 dark:text-surface-500 font-normal">Supports GitHub, GitLab, Bitbucket, & Custom Git</span>
                            </label>
                            <div class="grid grid-cols-1 sm:grid-cols-2 gap-2.5">
                                <button
                                    type="button"
                                    @click="togglePrivateRepo(false)"
                                    :class="[
                                        'px-3.5 py-2.5 rounded-xl border text-left flex items-center justify-between transition',
                                        !isPrivateRepo
                                            ? 'bg-white dark:bg-surface-900 border-brand-500 ring-1 ring-brand-500 text-slate-900 dark:text-white shadow-sm'
                                            : 'bg-white/60 dark:bg-surface-900/40 border-slate-200 dark:border-surface-800 text-slate-600 dark:text-surface-400 hover:border-slate-300 dark:hover:border-surface-700'
                                    ]"
                                >
                                    <div class="flex items-center gap-2.5">
                                        <div class="p-1.5 rounded-lg bg-emerald-50 dark:bg-emerald-500/10 text-emerald-600 dark:text-emerald-400">
                                            <Unlock class="w-4 h-4" />
                                        </div>
                                        <div>
                                            <p class="text-xs font-bold">Public Repository</p>
                                            <p class="text-[10px] text-slate-500 dark:text-surface-500">Open-source / publicly accessible repo</p>
                                        </div>
                                    </div>
                                    <CheckCircle2 v-if="!isPrivateRepo" class="w-4 h-4 text-brand-600 dark:text-brand-400" />
                                </button>

                                <button
                                    type="button"
                                    @click="togglePrivateRepo(true)"
                                    :class="[
                                        'px-3.5 py-2.5 rounded-xl border text-left flex items-center justify-between transition',
                                        isPrivateRepo
                                            ? 'bg-white dark:bg-surface-900 border-brand-500 ring-1 ring-brand-500 text-slate-900 dark:text-white shadow-sm'
                                            : 'bg-white/60 dark:bg-surface-900/40 border-slate-200 dark:border-surface-800 text-slate-600 dark:text-surface-400 hover:border-slate-300 dark:hover:border-surface-700'
                                    ]"
                                >
                                    <div class="flex items-center gap-2.5">
                                        <div class="p-1.5 rounded-lg bg-amber-50 dark:bg-amber-500/10 text-amber-600 dark:text-amber-400">
                                            <Lock class="w-4 h-4" />
                                        </div>
                                        <div>
                                            <p class="text-xs font-bold">Private Repository</p>
                                            <p class="text-[10px] text-slate-500 dark:text-surface-500">Authenticate with Deploy Key or Token</p>
                                        </div>
                                    </div>
                                    <CheckCircle2 v-if="isPrivateRepo" class="w-4 h-4 text-brand-600 dark:text-brand-400" />
                                </button>
                            </div>
                        </div>

                        <!-- Private Repo Authentication Sub-Options -->
                        <div v-if="isPrivateRepo" class="space-y-3 pt-1 border-t border-slate-200/80 dark:border-surface-800/80">
                            <label class="block text-xs font-semibold text-slate-700 dark:text-surface-200 mb-1 flex items-center justify-between">
                                <span>Authentication Method</span>
                                <span class="text-[11px] text-brand-600 dark:text-brand-400 font-medium">SSH Deploy Key recommended</span>
                            </label>
                            <div class="grid grid-cols-1 sm:grid-cols-2 gap-2.5">
                                <label
                                    :class="[
                                        'cursor-pointer rounded-xl border p-3 flex items-center justify-between transition',
                                        form.git_auth_type === 'ssh_key'
                                            ? 'bg-brand-50/80 border-brand-500 text-slate-900 dark:bg-brand-500/10 dark:text-white ring-1 ring-brand-500'
                                            : 'bg-white dark:bg-surface-900 border-slate-200 dark:border-surface-800 text-slate-600 dark:text-surface-400 hover:border-slate-300 dark:hover:border-surface-700'
                                    ]"
                                >
                                    <div class="flex items-center gap-2.5">
                                        <Key class="w-4 h-4 text-brand-600 dark:text-brand-400" />
                                        <div>
                                            <p class="text-xs font-bold text-slate-900 dark:text-white">SSH Deploy Key</p>
                                            <p class="text-[10px] text-slate-500 dark:text-surface-400">Add public key to repository</p>
                                        </div>
                                    </div>
                                    <input type="radio" :checked="form.git_auth_type === 'ssh_key'" @change="setGitAuthType('ssh_key')" class="sr-only" />
                                    <CheckCircle2 v-if="form.git_auth_type === 'ssh_key'" class="w-4 h-4 text-brand-600 dark:text-brand-400" />
                                </label>

                                <label
                                    :class="[
                                        'cursor-pointer rounded-xl border p-3 flex items-center justify-between transition',
                                        form.git_auth_type === 'token'
                                            ? 'bg-brand-50/80 border-brand-500 text-slate-900 dark:bg-brand-500/10 dark:text-white ring-1 ring-brand-500'
                                            : 'bg-white dark:bg-surface-900 border-slate-200 dark:border-surface-800 text-slate-600 dark:text-surface-400 hover:border-slate-300 dark:hover:border-surface-700'
                                    ]"
                                >
                                    <div class="flex items-center gap-2.5">
                                        <Code class="w-4 h-4 text-brand-600 dark:text-brand-400" />
                                        <div>
                                            <p class="text-xs font-bold text-slate-900 dark:text-white">Personal Access Token</p>
                                            <p class="text-[10px] text-slate-500 dark:text-surface-400">HTTPS Token / PAT</p>
                                        </div>
                                    </div>
                                    <input type="radio" :checked="form.git_auth_type === 'token'" @change="setGitAuthType('token')" class="sr-only" />
                                    <CheckCircle2 v-if="form.git_auth_type === 'token'" class="w-4 h-4 text-brand-600 dark:text-brand-400" />
                                </label>
                            </div>
                        </div>

                        <!-- Repo URL & Branch Inputs -->
                        <div class="grid grid-cols-1 sm:grid-cols-3 gap-3">
                            <div class="sm:col-span-2">
                                <label class="block text-xs font-semibold text-slate-700 dark:text-surface-200 mb-1.5 flex items-center justify-between">
                                    <span>Repository URL <span class="text-rose-500">*</span></span>
                                    <span v-if="isPrivateRepo && form.git_auth_type === 'ssh_key'" class="text-[10px] font-mono text-brand-600 dark:text-brand-400">
                                        git@...
                                    </span>
                                    <span v-else class="text-[10px] font-mono text-slate-400 dark:text-surface-500">
                                        https://...
                                    </span>
                                </label>
                                <div class="relative">
                                    <GitBranch class="w-4 h-4 text-slate-400 dark:text-surface-500 absolute left-3.5 top-1/2 -translate-y-1/2" />
                                    <input
                                        v-model="form.git_repository"
                                        type="text"
                                        required
                                        :placeholder="isPrivateRepo && form.git_auth_type === 'ssh_key' ? 'git@github.com:organization/private-repo.git' : 'https://github.com/organization/repo.git'"
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

                        <!-- SSH Deploy Key Box -->
                        <div v-if="isPrivateRepo && form.git_auth_type === 'ssh_key'" class="space-y-3 p-4 rounded-xl bg-white dark:bg-surface-900 border border-slate-200 dark:border-surface-700/80 shadow-sm">
                            <div class="flex items-center justify-between flex-wrap gap-2">
                                <div class="flex items-center gap-2">
                                    <div class="p-1 rounded-md bg-brand-50 dark:bg-brand-500/10 text-brand-600 dark:text-brand-400">
                                        <Key class="w-3.5 h-3.5" />
                                    </div>
                                    <span class="text-xs font-bold text-slate-900 dark:text-white">Generated Public Deploy Key</span>
                                    <span class="text-[10px] px-2 py-0.5 rounded-full bg-emerald-50 dark:bg-emerald-500/10 text-emerald-600 dark:text-emerald-400 font-medium border border-emerald-200/60 dark:border-emerald-500/20">
                                        Read-Only Access
                                    </span>
                                </div>
                                <div class="flex items-center gap-2">
                                    <button
                                        type="button"
                                        @click="regenerateDeployKey"
                                        :disabled="isGeneratingKey"
                                        class="px-2.5 py-1 text-[11px] font-medium rounded-lg text-slate-600 dark:text-surface-400 hover:text-slate-900 dark:hover:text-white hover:bg-slate-100 dark:hover:bg-surface-800 transition flex items-center gap-1"
                                    >
                                        <RefreshCw :class="['w-3 h-3', isGeneratingKey ? 'animate-spin' : '']" />
                                        <span>Regenerate</span>
                                    </button>
                                    <button
                                        type="button"
                                        @click="copyDeployKey"
                                        class="px-3 py-1 text-[11px] font-semibold rounded-lg bg-brand-600 hover:bg-brand-500 text-white transition flex items-center gap-1.5 shadow-sm"
                                    >
                                        <Check v-if="copiedDeployKey" class="w-3 h-3 text-emerald-200" />
                                        <Copy v-else class="w-3 h-3" />
                                        <span>{{ copiedDeployKey ? 'Copied!' : 'Copy Public Key' }}</span>
                                    </button>
                                </div>
                            </div>

                            <div class="relative">
                                <textarea
                                    :value="form.git_ssh_public_key || (isGeneratingKey ? 'Generating OpenSSH key pair...' : 'No public key generated yet.')"
                                    readonly
                                    rows="3"
                                    class="w-full px-3 py-2 rounded-lg bg-slate-900 text-slate-200 font-mono text-[11px] border border-slate-800 focus:outline-none select-all resize-none leading-relaxed"
                                    @click="($event.target as HTMLTextAreaElement).select()"
                                ></textarea>
                            </div>

                            <!-- Provider Guide Badges -->
                            <div class="p-3 rounded-lg bg-slate-50 dark:bg-surface-950/60 border border-slate-200/80 dark:border-surface-800 space-y-1.5 text-[11px] text-slate-600 dark:text-surface-400">
                                <p class="font-semibold text-slate-800 dark:text-surface-200 flex items-center gap-1.5">
                                    <Info class="w-3.5 h-3.5 text-brand-600 dark:text-brand-400" />
                                    How to add this deploy key to your Git provider:
                                </p>
                                <ul class="list-disc list-inside space-y-0.5 text-[10.5px] pl-1">
                                    <li><strong>GitHub:</strong> Repository <span class="text-slate-700 dark:text-surface-300 font-medium">Settings &rarr; Deploy Keys &rarr; Add deploy key</span> (paste key, write access not needed).</li>
                                    <li><strong>GitLab:</strong> Repository <span class="text-slate-700 dark:text-surface-300 font-medium">Settings &rarr; Repository &rarr; Deploy keys &rarr; Add key</span>.</li>
                                    <li><strong>Bitbucket:</strong> Repository <span class="text-slate-700 dark:text-surface-300 font-medium">Repository settings &rarr; Access keys &rarr; Add key</span>.</li>
                                </ul>
                            </div>

                            <!-- Advanced: Custom Private Key Toggle -->
                            <div class="pt-1">
                                <button
                                    type="button"
                                    @click="showCustomPrivateKey = !showCustomPrivateKey"
                                    class="text-[11px] text-slate-500 hover:text-brand-600 dark:text-surface-400 dark:hover:text-brand-400 font-medium flex items-center gap-1"
                                >
                                    <span>{{ showCustomPrivateKey ? 'Hide custom private key' : 'Paste custom private key instead (advanced)' }}</span>
                                </button>

                                <div v-if="showCustomPrivateKey" class="mt-2 space-y-1.5">
                                    <label class="block text-[11px] font-semibold text-slate-700 dark:text-surface-300">
                                        Custom SSH Private Key (OpenSSH format)
                                    </label>
                                    <textarea
                                        v-model="form.git_ssh_private_key"
                                        rows="4"
                                        placeholder="-----BEGIN OPENSSH PRIVATE KEY-----&#10;...&#10;-----END OPENSSH PRIVATE KEY-----"
                                        class="w-full px-3 py-2 rounded-lg bg-white dark:bg-surface-900 border border-slate-200 dark:border-surface-700 font-mono text-[11px] text-slate-900 dark:text-white placeholder-slate-400 dark:placeholder-surface-600 focus:outline-none focus:ring-2 focus:ring-brand-500/40"
                                    ></textarea>
                                </div>
                            </div>
                        </div>

                        <!-- Personal Access Token Box -->
                        <div v-if="isPrivateRepo && form.git_auth_type === 'token'" class="space-y-3 p-4 rounded-xl bg-white dark:bg-surface-900 border border-slate-200 dark:border-surface-700/80 shadow-sm">
                            <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
                                <div>
                                    <label class="block text-xs font-semibold text-slate-700 dark:text-surface-200 mb-1.5 flex items-center justify-between">
                                        <span>Personal Access Token / PAT <span class="text-rose-500">*</span></span>
                                        <button
                                            type="button"
                                            @click="showGitToken = !showGitToken"
                                            class="text-[11px] text-slate-400 hover:text-slate-600 dark:hover:text-surface-300 flex items-center gap-1"
                                        >
                                            <component :is="showGitToken ? EyeOff : Eye" class="w-3 h-3" />
                                            <span>{{ showGitToken ? 'Hide' : 'Show' }}</span>
                                        </button>
                                    </label>
                                    <input
                                        v-model="form.git_token"
                                        :type="showGitToken ? 'text' : 'password'"
                                        required
                                        placeholder="ghp_xxxxxxxxxxxx or glpat-xxxxxxxxxxxx"
                                        class="w-full px-3.5 py-2.5 rounded-xl bg-slate-50 dark:bg-surface-950/60 border border-slate-200 dark:border-surface-800 text-xs text-slate-900 dark:text-white font-mono focus:outline-none focus:ring-2 focus:ring-brand-500/40"
                                    />
                                    <p v-if="form.errors.git_token" class="text-[11px] text-rose-500 mt-1.5">{{ form.errors.git_token }}</p>
                                </div>

                                <div>
                                    <label class="block text-xs font-semibold text-slate-700 dark:text-surface-200 mb-1.5 flex items-center justify-between">
                                        <span>Token Username</span>
                                        <span class="text-[10px] text-slate-400 font-normal">Optional</span>
                                    </label>
                                    <input
                                        v-model="form.git_token_user"
                                        type="text"
                                        placeholder="e.g. x-access-token or username"
                                        class="w-full px-3.5 py-2.5 rounded-xl bg-slate-50 dark:bg-surface-950/60 border border-slate-200 dark:border-surface-800 text-xs text-slate-900 dark:text-white font-mono focus:outline-none focus:ring-2 focus:ring-brand-500/40"
                                    />
                                    <p class="text-[10px] text-slate-400 dark:text-surface-500 mt-1">
                                        Default is auto-detected (GitHub: <code class="font-mono">x-access-token</code>, GitLab: <code class="font-mono">oauth2</code>).
                                    </p>
                                </div>
                            </div>

                            <div class="p-3 rounded-lg bg-slate-50 dark:bg-surface-950/60 border border-slate-200/80 dark:border-surface-800 text-[10.5px] text-slate-500 dark:text-surface-400 space-y-1">
                                <p class="font-semibold text-slate-700 dark:text-surface-300">Required Token Scopes / Permissions:</p>
                                <p>Ensure your token has <strong class="text-slate-800 dark:text-surface-200">repo</strong> (GitHub) or <strong class="text-slate-800 dark:text-surface-200">read_repository</strong> (GitLab) permission to clone private repositories.</p>
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

                    <!-- Database Provisioning Card -->
                    <div class="rounded-2xl border border-slate-200/80 dark:border-surface-800 bg-slate-50/50 dark:bg-surface-950/40 p-5 space-y-4">
                        <div class="flex items-start justify-between">
                            <div class="flex items-start gap-3">
                                <div class="p-2 rounded-xl bg-brand-50 dark:bg-brand-500/10 border border-brand-200 dark:border-brand-500/20 text-brand-600 dark:text-brand-400 mt-0.5">
                                    <Database class="w-4 h-4" />
                                </div>
                                <div>
                                    <h3 class="text-xs font-bold text-slate-900 dark:text-white flex items-center gap-2">
                                        <span>Create & Link Database</span>
                                        <span class="px-2 py-0.5 rounded-full text-[10px] font-semibold bg-brand-100 dark:bg-brand-500/20 text-brand-700 dark:text-brand-300">
                                            Automated
                                        </span>
                                    </h3>
                                    <p class="text-[11px] text-slate-500 dark:text-surface-400 mt-0.5">
                                        Provisions a dedicated database and database user with full grant access automatically.
                                    </p>
                                </div>
                            </div>
                            <label class="relative inline-flex items-center cursor-pointer">
                                <input type="checkbox" v-model="form.create_database" class="sr-only peer" />
                                <div class="w-9 h-5 bg-slate-200 peer-focus:outline-none rounded-full peer dark:bg-surface-700 peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-slate-300 after:border after:rounded-full after:h-4 after:w-4 after:transition-all dark:border-surface-600 peer-checked:bg-brand-600"></div>
                            </label>
                        </div>

                        <!-- Database Fields -->
                        <div v-if="form.create_database" class="pt-3 border-t border-slate-200/60 dark:border-surface-800/60 space-y-4">
                            <!-- Engine Selector -->
                            <div>
                                <label class="block text-[11px] font-semibold text-slate-700 dark:text-surface-300 mb-1.5">
                                    Database Engine
                                </label>
                                <div class="grid grid-cols-2 gap-3">
                                    <label
                                        :class="[
                                            'cursor-pointer rounded-xl border p-3 flex items-center justify-between transition',
                                            form.db_engine === 'mysql'
                                                ? 'bg-white dark:bg-surface-900 border-brand-500 text-slate-900 dark:text-white ring-1 ring-brand-500 shadow-xs'
                                                : 'bg-white/60 dark:bg-surface-900/40 border-slate-200 dark:border-surface-800 text-slate-600 dark:text-surface-400'
                                        ]"
                                    >
                                        <div class="flex items-center gap-2">
                                            <div class="w-2 h-2 rounded-full bg-emerald-500"></div>
                                            <div>
                                                <div class="text-xs font-bold font-mono">MySQL 8.0</div>
                                                <div class="text-[10px] text-slate-400">Port 3306 (utf8mb4)</div>
                                            </div>
                                        </div>
                                        <input type="radio" v-model="form.db_engine" value="mysql" class="sr-only" />
                                        <CheckCircle2 v-if="form.db_engine === 'mysql'" class="w-3.5 h-3.5 text-brand-600 dark:text-brand-400" />
                                    </label>

                                    <label
                                        :class="[
                                            'cursor-pointer rounded-xl border p-3 flex items-center justify-between transition',
                                            form.db_engine === 'postgresql'
                                                ? 'bg-white dark:bg-surface-900 border-brand-500 text-slate-900 dark:text-white ring-1 ring-brand-500 shadow-xs'
                                                : 'bg-white/60 dark:bg-surface-900/40 border-slate-200 dark:border-surface-800 text-slate-600 dark:text-surface-400'
                                        ]"
                                    >
                                        <div class="flex items-center gap-2">
                                            <div class="w-2 h-2 rounded-full bg-sky-500"></div>
                                            <div>
                                                <div class="text-xs font-bold font-mono">PostgreSQL 16</div>
                                                <div class="text-[10px] text-slate-400">Port 5432 (UTF8)</div>
                                            </div>
                                        </div>
                                        <input type="radio" v-model="form.db_engine" value="postgresql" class="sr-only" />
                                        <CheckCircle2 v-if="form.db_engine === 'postgresql'" class="w-3.5 h-3.5 text-brand-600 dark:text-brand-400" />
                                    </label>
                                </div>
                            </div>

                            <!-- Name & User Grid -->
                            <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
                                <div>
                                    <label class="block text-[11px] font-semibold text-slate-700 dark:text-surface-300 mb-1">
                                        Database Name
                                    </label>
                                    <input
                                        v-model="form.db_name"
                                        type="text"
                                        required
                                        placeholder="db_myapp"
                                        class="w-full px-3 py-2 rounded-xl bg-white dark:bg-surface-900 border border-slate-200 dark:border-surface-800 text-xs text-slate-900 dark:text-white placeholder-slate-400 dark:placeholder-surface-600 font-mono focus:outline-none focus:ring-2 focus:ring-brand-500/40 transition"
                                    />
                                    <p v-if="form.errors.db_name" class="text-[10px] text-rose-500 mt-1">{{ form.errors.db_name }}</p>
                                </div>

                                <div>
                                    <label class="block text-[11px] font-semibold text-slate-700 dark:text-surface-300 mb-1">
                                        Database Username
                                    </label>
                                    <input
                                        v-model="form.db_username"
                                        type="text"
                                        required
                                        placeholder="u_myapp"
                                        class="w-full px-3 py-2 rounded-xl bg-white dark:bg-surface-900 border border-slate-200 dark:border-surface-800 text-xs text-slate-900 dark:text-white placeholder-slate-400 dark:placeholder-surface-600 font-mono focus:outline-none focus:ring-2 focus:ring-brand-500/40 transition"
                                    />
                                    <p v-if="form.errors.db_username" class="text-[10px] text-rose-500 mt-1">{{ form.errors.db_username }}</p>
                                </div>
                            </div>

                            <!-- Password with Generator -->
                            <div>
                                <label class="block text-[11px] font-semibold text-slate-700 dark:text-surface-300 mb-1 flex items-center justify-between">
                                    <span>Database Password</span>
                                    <div class="flex items-center gap-2">
                                        <button
                                            type="button"
                                            @click="regeneratePassword"
                                            class="text-[10px] text-brand-600 dark:text-brand-400 hover:underline flex items-center gap-1 font-medium"
                                        >
                                            <RefreshCw class="w-3 h-3" />
                                            <span>Generate New</span>
                                        </button>
                                        <span class="text-slate-300 dark:text-surface-700">•</span>
                                        <button
                                            type="button"
                                            @click="copyDbPassword"
                                            class="text-[10px] text-slate-600 dark:text-surface-400 hover:text-slate-900 dark:hover:text-white flex items-center gap-1 font-medium"
                                        >
                                            <component :is="copiedDbPassword ? Check : Copy" class="w-3 h-3 text-emerald-500" />
                                            <span>{{ copiedDbPassword ? 'Copied!' : 'Copy' }}</span>
                                        </button>
                                    </div>
                                </label>
                                <div class="relative">
                                    <Key class="w-4 h-4 text-slate-400 dark:text-surface-500 absolute left-3 top-1/2 -translate-y-1/2" />
                                    <input
                                        v-model="form.db_password"
                                        :type="showDbPassword ? 'text' : 'password'"
                                        required
                                        class="w-full pl-9 pr-10 py-2 rounded-xl bg-white dark:bg-surface-900 border border-slate-200 dark:border-surface-800 text-xs text-slate-900 dark:text-white font-mono placeholder-slate-400 dark:placeholder-surface-600 focus:outline-none focus:ring-2 focus:ring-brand-500/40 transition"
                                    />
                                    <button
                                        type="button"
                                        @click="showDbPassword = !showDbPassword"
                                        class="absolute right-3 top-1/2 -translate-y-1/2 text-slate-400 hover:text-slate-600 dark:hover:text-white"
                                    >
                                        <component :is="showDbPassword ? EyeOff : Eye" class="w-3.5 h-3.5" />
                                    </button>
                                </div>
                                <p class="text-[10px] text-slate-500 dark:text-surface-400 mt-1">
                                    These credentials will be injected directly into your Laravel <code class="font-mono text-brand-600 dark:text-brand-400 font-semibold">.env</code> configuration.
                                </p>
                            </div>
                        </div>
                    </div>

                    <!-- Laravel Automated Post-Setup Pipeline (Visible for Laravel) -->
                    <div
                        v-if="form.project_type === 'laravel'"
                        class="rounded-2xl border border-rose-200/80 dark:border-rose-900/30 bg-rose-50/40 dark:bg-rose-950/10 p-5 space-y-4"
                    >
                        <div class="flex items-start justify-between">
                            <div class="flex items-start gap-3">
                                <div class="p-2 rounded-xl bg-rose-100 dark:bg-rose-500/20 border border-rose-200 dark:border-rose-500/30 text-rose-600 dark:text-rose-400 mt-0.5">
                                    <Sparkles class="w-4 h-4" />
                                </div>
                                <div>
                                    <h3 class="text-xs font-bold text-slate-900 dark:text-white flex items-center gap-2">
                                        <span>Laravel Automated Zero-SSH Setup</span>
                                        <span class="px-2 py-0.5 rounded-full text-[10px] font-semibold bg-rose-100 dark:bg-rose-500/20 text-rose-700 dark:text-rose-300">
                                            1-Click Provisioning
                                        </span>
                                    </h3>
                                    <p class="text-[11px] text-slate-500 dark:text-surface-400 mt-0.5">
                                        Automatically executes all necessary build, migration, and dependency tasks so your app is immediately live without touching the terminal.
                                    </p>
                                </div>
                            </div>
                            <label class="relative inline-flex items-center cursor-pointer">
                                <input type="checkbox" v-model="form.auto_setup_laravel" class="sr-only peer" />
                                <div class="w-9 h-5 bg-slate-200 peer-focus:outline-none rounded-full peer dark:bg-surface-700 peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-slate-300 after:border after:rounded-full after:h-4 after:w-4 after:transition-all dark:border-surface-600 peer-checked:bg-rose-600"></div>
                            </label>
                        </div>

                        <!-- Pipeline Options Checklist -->
                        <div v-if="form.auto_setup_laravel" class="pt-3 border-t border-rose-200/60 dark:border-rose-900/20 grid grid-cols-1 sm:grid-cols-2 gap-2.5">
                            <!-- Option 1: .env configuration -->
                            <label class="flex items-start gap-2.5 p-2.5 rounded-xl bg-white dark:bg-surface-900/80 border border-slate-200/80 dark:border-surface-800 cursor-pointer hover:border-slate-300 dark:hover:border-surface-700 transition">
                                <input type="checkbox" v-model="form.setup_env" class="mt-0.5 rounded border-slate-300 text-rose-600 focus:ring-rose-500" />
                                <div>
                                    <span class="text-xs font-semibold text-slate-800 dark:text-surface-200 block">Configure .env</span>
                                    <span class="text-[10px] text-slate-500 dark:text-surface-400">Sync database credentials, APP_URL, and production settings.</span>
                                </div>
                            </label>

                            <!-- Option 2: Composer Install -->
                            <label class="flex items-start gap-2.5 p-2.5 rounded-xl bg-white dark:bg-surface-900/80 border border-slate-200/80 dark:border-surface-800 cursor-pointer hover:border-slate-300 dark:hover:border-surface-700 transition">
                                <input type="checkbox" v-model="form.run_composer" class="mt-0.5 rounded border-slate-300 text-rose-600 focus:ring-rose-500" />
                                <div>
                                    <span class="text-xs font-semibold text-slate-800 dark:text-surface-200 block">composer install</span>
                                    <span class="text-[10px] text-slate-500 dark:text-surface-400">Install PHP packages with optimized autoloader (--no-dev).</span>
                                </div>
                            </label>

                            <!-- Option 3: Key Generate -->
                            <label class="flex items-start gap-2.5 p-2.5 rounded-xl bg-white dark:bg-surface-900/80 border border-slate-200/80 dark:border-surface-800 cursor-pointer hover:border-slate-300 dark:hover:border-surface-700 transition">
                                <input type="checkbox" v-model="form.run_key_generate" class="mt-0.5 rounded border-slate-300 text-rose-600 focus:ring-rose-500" />
                                <div>
                                    <span class="text-xs font-semibold text-slate-800 dark:text-surface-200 block">php artisan key:generate</span>
                                    <span class="text-[10px] text-slate-500 dark:text-surface-400">Generate secure Laravel APP_KEY encryption key.</span>
                                </div>
                            </label>

                            <!-- Option 4: Migrate -->
                            <label class="flex items-start gap-2.5 p-2.5 rounded-xl bg-white dark:bg-surface-900/80 border border-slate-200/80 dark:border-surface-800 cursor-pointer hover:border-slate-300 dark:hover:border-surface-700 transition">
                                <input type="checkbox" v-model="form.run_migrations" class="mt-0.5 rounded border-slate-300 text-rose-600 focus:ring-rose-500" />
                                <div>
                                    <span class="text-xs font-semibold text-slate-800 dark:text-surface-200 block">php artisan migrate</span>
                                    <span class="text-[10px] text-slate-500 dark:text-surface-400">Run database migrations on the newly created database.</span>
                                </div>
                            </label>

                            <!-- Option 5: Seeders -->
                            <label class="flex items-start gap-2.5 p-2.5 rounded-xl bg-white dark:bg-surface-900/80 border border-slate-200/80 dark:border-surface-800 cursor-pointer hover:border-slate-300 dark:hover:border-surface-700 transition">
                                <input type="checkbox" v-model="form.run_seeders" class="mt-0.5 rounded border-slate-300 text-rose-600 focus:ring-rose-500" />
                                <div>
                                    <span class="text-xs font-semibold text-slate-800 dark:text-surface-200 block">php artisan db:seed</span>
                                    <span class="text-[10px] text-slate-500 dark:text-surface-400">Seed database with initial data / demo records.</span>
                                </div>
                            </label>

                            <!-- Option 6: NPM Build -->
                            <label class="flex items-start gap-2.5 p-2.5 rounded-xl bg-white dark:bg-surface-900/80 border border-slate-200/80 dark:border-surface-800 cursor-pointer hover:border-slate-300 dark:hover:border-surface-700 transition">
                                <input type="checkbox" v-model="form.run_npm_build" class="mt-0.5 rounded border-slate-300 text-rose-600 focus:ring-rose-500" />
                                <div>
                                    <span class="text-xs font-semibold text-slate-800 dark:text-surface-200 block">npm install & build</span>
                                    <span class="text-[10px] text-slate-500 dark:text-surface-400">Compile frontend assets (Vite / Tailwind / Vue / React).</span>
                                </div>
                            </label>

                            <!-- Option 7: Optimize -->
                            <label class="flex items-start gap-2.5 p-2.5 rounded-xl bg-white dark:bg-surface-900/80 border border-slate-200/80 dark:border-surface-800 cursor-pointer hover:border-slate-300 dark:hover:border-surface-700 transition sm:col-span-2">
                                <input type="checkbox" v-model="form.run_optimize" class="mt-0.5 rounded border-slate-300 text-rose-600 focus:ring-rose-500" />
                                <div>
                                    <span class="text-xs font-semibold text-slate-800 dark:text-surface-200 block">php artisan optimize</span>
                                    <span class="text-[10px] text-slate-500 dark:text-surface-400">Cache configuration, routes, and views for production speed.</span>
                                </div>
                            </label>
                        </div>
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
