<script setup lang="ts">
import { ref } from 'vue';
import { useForm, router } from '@inertiajs/vue3';
import AppLayout from '@/Layouts/AppLayout.vue';
import {
    Folder,
    FileText,
    FolderPlus,
    FilePlus,
    ArrowLeft,
    Trash2,
    Edit3,
    Save,
    X,
    Globe,
    HardDrive,
    Code,
    Check,
} from 'lucide-vue-next';
import type { FileEntry, Website } from '@/types';

const props = defineProps<{
    websites: Website[];
    selectedWebsite?: Website | null;
    currentPath: string;
    basePath: string;
    files: FileEntry[];
}>();

// Navigation
function navigateToPath(relPath: string) {
    router.get('/files', {
        website_id: props.selectedWebsite?.id,
        path: relPath,
    });
}

function onWebsiteChange(e: Event) {
    const target = e.target as HTMLSelectElement;
    router.get('/files', {
        website_id: target.value,
        path: '',
    });
}

function goUp() {
    if (!props.currentPath) return;
    const parts = props.currentPath.split('/').filter(Boolean);
    parts.pop();
    navigateToPath(parts.join('/'));
}

// File Editor Modal
const isEditorOpen = ref(false);
const editingFileName = ref('');
const fileEditForm = useForm({
    base_path: props.basePath,
    relative_path: '',
    content: '',
    website_id: props.selectedWebsite?.id,
    current_path: props.currentPath,
});

async function openFileEditor(entry: FileEntry) {
    editingFileName.value = entry.name;
    fileEditForm.base_path = props.basePath;
    fileEditForm.relative_path = entry.path;
    fileEditForm.website_id = props.selectedWebsite?.id;
    fileEditForm.current_path = props.currentPath;

    try {
        const res = await fetch(`/files/read?base_path=${encodeURIComponent(props.basePath)}&relative_path=${encodeURIComponent(entry.path)}`);
        const data = await res.json();
        if (data.success) {
            fileEditForm.content = data.content;
            isEditorOpen.value = true;
        } else {
            alert(data.error || 'Failed to read file');
        }
    } catch (e) {
        alert('Network error while reading file');
    }
}

function saveFileContent() {
    fileEditForm.post('/files/write', {
        onSuccess: () => {
            isEditorOpen.value = false;
        },
    });
}

// File Delete
function deleteFileEntry(entry: FileEntry) {
    if (confirm(`Are you sure you want to delete ${entry.name}?`)) {
        router.post('/files/delete', {
            base_path: props.basePath,
            relative_path: entry.path,
            website_id: props.selectedWebsite?.id,
            current_path: props.currentPath,
        });
    }
}

function formatBytes(bytes: number): string {
    if (bytes === 0) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i];
}
</script>

<template>
    <AppLayout title="File Manager">
        <div class="max-w-7xl mx-auto space-y-6">
            <!-- Header & Root Switcher -->
            <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
                <div>
                    <h2 class="text-base font-semibold text-white tracking-tight">Path-Jailed File Manager</h2>
                    <p class="text-xs text-surface-400 mt-0.5">
                        Securely browse, edit configuration files, and manage webroot assets
                    </p>
                </div>
                <div class="flex items-center gap-3">
                    <div class="flex items-center gap-2 bg-surface-900 border border-surface-800 rounded-xl px-3 py-1.5">
                        <Globe class="w-3.5 h-3.5 text-brand-400" />
                        <select
                            :value="selectedWebsite?.id"
                            @change="onWebsiteChange"
                            class="bg-transparent text-xs text-white focus:outline-none font-mono"
                        >
                            <option v-for="w in websites" :key="w.id" :value="w.id" class="bg-surface-900 text-white">
                                {{ w.domain }}
                            </option>
                        </select>
                    </div>
                </div>
            </div>

            <!-- Explorer Card -->
            <div class="rounded-2xl bg-surface-900/60 border border-surface-800/80 shadow-lg overflow-hidden">
                <!-- Navigation Bar -->
                <div class="p-3.5 bg-surface-950/60 border-b border-surface-800 flex items-center justify-between gap-3 text-xs">
                    <div class="flex items-center gap-2 font-mono text-surface-300 overflow-x-auto">
                        <button
                            v-if="currentPath"
                            @click="goUp"
                            class="p-1 rounded-lg bg-surface-800 hover:bg-surface-700 text-surface-300 mr-1"
                            title="Go Up"
                        >
                            <ArrowLeft class="w-3.5 h-3.5" />
                        </button>
                        <button @click="navigateToPath('')" class="text-brand-400 hover:underline">
                            /var/www/{{ selectedWebsite?.domain || 'root' }}
                        </button>
                        <span v-if="currentPath" class="text-surface-600">/</span>
                        <span v-if="currentPath" class="text-white">{{ currentPath }}</span>
                    </div>
                </div>

                <!-- Files Table -->
                <div v-if="files.length === 0" class="text-center py-14 px-4">
                    <Folder class="w-10 h-10 text-surface-600 mx-auto mb-3" />
                    <h3 class="text-sm font-medium text-surface-200">Directory is empty</h3>
                    <p class="text-xs text-surface-400 mt-1 max-w-sm mx-auto">
                        No files or folders found in this directory.
                    </p>
                </div>

                <div v-else class="overflow-x-auto">
                    <table class="w-full text-left text-xs">
                        <thead class="bg-surface-950/50 text-surface-400 uppercase text-[10px] tracking-wider border-b border-surface-800">
                            <tr>
                                <th class="py-3.5 px-4 font-semibold">Name</th>
                                <th class="py-3.5 px-4 font-semibold">Size</th>
                                <th class="py-3.5 px-4 font-semibold">Permissions</th>
                                <th class="py-3.5 px-4 font-semibold">Modified</th>
                                <th class="py-3.5 px-4 font-semibold text-right">Actions</th>
                            </tr>
                        </thead>
                        <tbody class="divide-y divide-surface-800/60">
                            <tr
                                v-for="file in files"
                                :key="file.name"
                                class="hover:bg-surface-800/30 transition"
                            >
                                <td class="py-3 px-4 font-mono font-medium text-white">
                                    <div class="flex items-center gap-2.5">
                                        <div
                                            :class="[
                                                'p-1.5 rounded-lg shrink-0',
                                                file.is_dir
                                                    ? 'bg-amber-500/10 text-amber-400 border border-amber-500/20'
                                                    : 'bg-brand-500/10 text-brand-400 border border-brand-500/20'
                                            ]"
                                        >
                                            <Folder v-if="file.is_dir" class="w-3.5 h-3.5" />
                                            <FileText v-else class="w-3.5 h-3.5" />
                                        </div>
                                        <button
                                            v-if="file.is_dir"
                                            @click="navigateToPath(file.path)"
                                            class="text-amber-300 hover:underline font-semibold text-left"
                                        >
                                            {{ file.name }}
                                        </button>
                                        <span v-else class="text-surface-200">{{ file.name }}</span>
                                    </div>
                                </td>
                                <td class="py-3 px-4 font-mono text-[11px] text-surface-400">
                                    {{ file.is_dir ? '—' : formatBytes(file.size_bytes) }}
                                </td>
                                <td class="py-3 px-4 font-mono text-[10px] text-surface-500">
                                    {{ file.permissions }}
                                </td>
                                <td class="py-3 px-4 font-mono text-[11px] text-surface-500">
                                    {{ new Date(file.modified_at).toLocaleDateString() }}
                                </td>
                                <td class="py-3 px-4 text-right">
                                    <div class="flex items-center justify-end gap-1.5">
                                        <button
                                            v-if="!file.is_dir"
                                            @click="openFileEditor(file)"
                                            class="p-1.5 rounded-lg text-surface-300 hover:text-white hover:bg-surface-800 transition"
                                            title="Edit File"
                                        >
                                            <Edit3 class="w-3.5 h-3.5" />
                                        </button>
                                        <button
                                            @click="deleteFileEntry(file)"
                                            class="p-1.5 rounded-lg text-rose-400 hover:bg-rose-500/10 transition"
                                            title="Delete"
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

            <!-- Code Editor Modal -->
            <div
                v-if="isEditorOpen"
                class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-surface-950/80 backdrop-blur-sm"
            >
                <div class="w-full max-w-4xl h-[80vh] flex flex-col rounded-2xl bg-surface-900 border border-surface-800 shadow-2xl overflow-hidden">
                    <!-- Editor Titlebar -->
                    <div class="h-12 bg-surface-950 border-b border-surface-800 px-4 flex items-center justify-between shrink-0">
                        <div class="flex items-center gap-2 font-mono text-xs text-white">
                            <Code class="w-4 h-4 text-brand-400" />
                            <span>{{ editingFileName }}</span>
                        </div>
                        <div class="flex items-center gap-2">
                            <button
                                @click="saveFileContent"
                                :disabled="fileEditForm.processing"
                                class="px-3.5 py-1.5 rounded-xl bg-emerald-600 hover:bg-emerald-500 text-white text-xs font-semibold flex items-center gap-1.5 transition shadow-lg shadow-emerald-600/20 disabled:opacity-50"
                            >
                                <Save class="w-3.5 h-3.5" />
                                <span>{{ fileEditForm.processing ? 'Saving...' : 'Save File' }}</span>
                            </button>
                            <button
                                @click="isEditorOpen = false"
                                class="p-1.5 rounded-xl text-surface-400 hover:text-white hover:bg-surface-800 transition"
                            >
                                <X class="w-4 h-4" />
                            </button>
                        </div>
                    </div>

                    <!-- Editor Textarea -->
                    <div class="flex-1 bg-surface-950 p-4 font-mono text-xs text-surface-200 overflow-auto">
                        <textarea
                            v-model="fileEditForm.content"
                            spellcheck="false"
                            class="w-full h-full bg-transparent border-0 resize-none font-mono text-xs text-emerald-300 focus:outline-none leading-relaxed selection:bg-brand-500 selection:text-white"
                        ></textarea>
                    </div>
                </div>
            </div>
        </div>
    </AppLayout>
</template>
