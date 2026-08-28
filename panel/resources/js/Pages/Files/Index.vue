<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, watch } from 'vue';
import { router } from '@inertiajs/vue3';
import AppLayout from '@/Layouts/AppLayout.vue';
import {
    Folder,
    FolderPlus,
    FilePlus,
    Upload,
    Download,
    Copy,
    Scissors,
    Trash2,
    Edit3,
    Save,
    X,
    Globe,
    HardDrive,
    Code,
    Check,
    ArrowLeft,
    ArrowRight,
    CornerLeftUp,
    RefreshCw,
    Search,
    FileText,
    Image as ImageIcon,
    Archive,
    Shield,
    User as UserIcon,
    Info,
    Eye,
    MoreVertical,
    CheckSquare,
    Square,
    Maximize2,
    Minimize2,
    Filter,
    Lock,
    ChevronRight,
    AlertTriangle,
    FileCode,
    Terminal,
    File as GenericFileIcon,
    Layers,
    FileSpreadsheet,
    FileArchive,
    RotateCcw,
    ChevronDown,
} from 'lucide-vue-next';
import type { FileEntry, FileDetails, DiskUsageInfo, Website } from '@/types';

const props = defineProps<{
    websites: Website[];
    selectedWebsite?: Website | null;
    currentPath: string;
    basePath: string;
    files: FileEntry[];
    diskUsage?: DiskUsageInfo | null;
    showHidden?: boolean;
}>();

// Navigation & Location State
const activeWebsite = ref<Website | null>(props.selectedWebsite || props.websites[0] || null);
const currentBasePath = ref<string>(props.basePath);
const currentRelPath = ref<string>(props.currentPath || '');
const fileList = ref<FileEntry[]>(props.files || []);
const diskInfo = ref<DiskUsageInfo | null>(props.diskUsage || null);
const isShowHidden = ref<boolean>(props.showHidden || false);
const isLoading = ref<boolean>(false);
const activeFilter = ref<string>('all'); // all, folders, files, images, code, archives
const searchQuery = ref<string>('');

// History Stack for Back / Forward
const historyStack = ref<string[]>([props.currentPath || '']);
const historyIndex = ref<number>(0);

// Sync with Props changes
watch(() => props.files, (newFiles) => {
    if (newFiles) fileList.value = newFiles;
});
watch(() => props.currentPath, (newPath) => {
    if (newPath !== undefined) currentRelPath.value = newPath;
});
watch(() => props.basePath, (newBase) => {
    if (newBase) currentBasePath.value = newBase;
});
watch(() => props.diskUsage, (newDisk) => {
    if (newDisk) diskInfo.value = newDisk;
});
watch(() => props.showHidden, (newHidden) => {
    if (newHidden !== undefined) isShowHidden.value = newHidden;
});
watch(() => props.selectedWebsite, (newSite) => {
    if (newSite) activeWebsite.value = newSite;
});

// Selection State
const selectedPaths = ref<Set<string>>(new Set());

// Context Menu State
const contextMenu = ref<{
    isOpen: boolean;
    x: number;
    y: number;
    item: FileEntry | null;
}>({
    isOpen: false,
    x: 0,
    y: 0,
    item: null,
});

// Drag & Drop Upload State
const isDraggingOver = ref<boolean>(false);
const isUploading = ref<boolean>(false);
const uploadProgress = ref<number>(0);
const uploadFileName = ref<string>('');
const fileInputRef = ref<HTMLInputElement | null>(null);

// Modal States
const isCreateFileOpen = ref<boolean>(false);
const newFileName = ref<string>('');

const isCreateFolderOpen = ref<boolean>(false);
const newFolderName = ref<string>('');

const isRenameOpen = ref<boolean>(false);
const renameTarget = ref<FileEntry | null>(null);
const renameNewName = ref<string>('');

const isCopyMoveOpen = ref<boolean>(false);
const copyMoveMode = ref<'copy' | 'move'>('copy');
const copyMoveDestPath = ref<string>('');
const copyMoveItems = ref<string[]>([]);

const isDeleteOpen = ref<boolean>(false);
const deleteTargets = ref<string[]>([]);

const isCompressOpen = ref<boolean>(false);
const compressArchiveName = ref<string>('archive.zip');
const compressFormat = ref<'zip' | 'tar.gz'>('zip');
const compressSources = ref<string[]>([]);

const isExtractOpen = ref<boolean>(false);
const extractTarget = ref<FileEntry | null>(null);
const extractDestPath = ref<string>('');

const isPermissionsOpen = ref<boolean>(false);
const permissionsTarget = ref<FileEntry | null>(null);
const permOwnerR = ref<boolean>(true);
const permOwnerW = ref<boolean>(true);
const permOwnerX = ref<boolean>(false);
const permGroupR = ref<boolean>(true);
const permGroupW = ref<boolean>(false);
const permGroupX = ref<boolean>(false);
const permOtherR = ref<boolean>(true);
const permOtherW = ref<boolean>(false);
const permOtherX = ref<boolean>(false);
const permRecursive = ref<boolean>(false);

const isDetailsOpen = ref<boolean>(false);
const detailsData = ref<FileDetails | null>(null);

const isPreviewOpen = ref<boolean>(false);
const previewItem = ref<FileEntry | null>(null);
const previewTextContent = ref<string>('');

// Code Editor State
const isEditorOpen = ref<boolean>(false);
const isEditorFullscreen = ref<boolean>(false);
const isEditorModified = ref<boolean>(false);
const editingFile = ref<FileEntry | null>(null);
const editorContent = ref<string>('');
const initialEditorContent = ref<string>('');
const isEditorSaving = ref<boolean>(false);
const editorWordWrap = ref<boolean>(false);
const isEditorSearchOpen = ref<boolean>(false);
const editorFindText = ref<string>('');
const editorReplaceText = ref<string>('');

// Sort State
const sortBy = ref<'name' | 'size' | 'modified' | 'permissions'>('name');
const sortAsc = ref<boolean>(true);

// ==========================================
// AJAX DIRECTORY BROWSING & NAVIGATION
// ==========================================

async function fetchDirectory(relPath: string, addToHistory: boolean = true) {
    isLoading.value = true;
    selectedPaths.value.clear();
    contextMenu.value.isOpen = false;

    try {
        const res = await fetch('/files/browse', {
            method: 'POST',
            headers: {
                'Accept': 'application/json',
                'Content-Type': 'application/json',
                'X-CSRF-TOKEN': getCsrfToken(),
            },
            body: JSON.stringify({
                base_path: currentBasePath.value,
                relative_path: relPath,
                show_hidden: isShowHidden.value,
                website_id: activeWebsite.value?.id,
            }),
        });
        const data = await res.json();
        if (data.success) {
            fileList.value = data.files || [];
            if (data.disk_usage) {
                diskInfo.value = data.disk_usage;
            }
            currentRelPath.value = relPath;

            if (addToHistory) {
                // If browsing forward from middle of stack, truncate future history
                historyStack.value = historyStack.value.slice(0, historyIndex.value + 1);
                historyStack.value.push(relPath);
                historyIndex.value = historyStack.value.length - 1;
            }
        } else {
            showToast(data.error || 'Failed to load directory', 'error');
        }
    } catch (e: any) {
        showToast('Network error loading files', 'error');
    } finally {
        isLoading.value = false;
    }
}

function navigateBack() {
    if (historyIndex.value > 0) {
        historyIndex.value--;
        const path = historyStack.value[historyIndex.value];
        fetchDirectory(path, false);
    }
}

function navigateForward() {
    if (historyIndex.value < historyStack.value.length - 1) {
        historyIndex.value++;
        const path = historyStack.value[historyIndex.value];
        fetchDirectory(path, false);
    }
}

function goUp() {
    if (!currentRelPath.value) return;
    const parts = currentRelPath.value.split('/').filter(Boolean);
    parts.pop();
    fetchDirectory(parts.join('/'));
}

function refreshDirectory() {
    fetchDirectory(currentRelPath.value, false);
}

function toggleShowHidden() {
    isShowHidden.value = !isShowHidden.value;
    fetchDirectory(currentRelPath.value, false);
}

function onWebsiteChange(e: Event) {
    const target = e.target as HTMLSelectElement;
    const site = props.websites.find(w => w.id === target.value);
    if (site) {
        activeWebsite.value = site;
        currentBasePath.value = site.document_root ? site.document_root.replace(/\/public\/?$/, '') : '/var/www';
        currentRelPath.value = '';
        historyStack.value = [''];
        historyIndex.value = 0;
        fetchDirectory('', false);
    }
}

// ==========================================
// BREADCRUMBS
// ==========================================

const breadcrumbSegments = computed(() => {
    const segments = [];
    // Base root name
    const domain = activeWebsite.value?.domain || 'www';
    segments.push({ name: `/var/www/${domain}`, path: '' });

    if (currentRelPath.value) {
        const parts = currentRelPath.value.split('/').filter(Boolean);
        let accumulated = '';
        for (const part of parts) {
            accumulated = accumulated ? `${accumulated}/${part}` : part;
            segments.push({ name: part, path: accumulated });
        }
    }
    return segments;
});

// ==========================================
// FILTERING & SORTING
// ==========================================

const filteredFiles = computed(() => {
    let result = [...fileList.value];

    // Filter by search query
    if (searchQuery.value.trim()) {
        const q = searchQuery.value.toLowerCase().trim();
        result = result.filter(f => f.name.toLowerCase().includes(q));
    }

    // Filter by category
    if (activeFilter.value === 'folders') {
        result = result.filter(f => f.is_dir);
    } else if (activeFilter.value === 'files') {
        result = result.filter(f => !f.is_dir);
    } else if (activeFilter.value === 'images') {
        result = result.filter(f => isImageFile(f.name));
    } else if (activeFilter.value === 'code') {
        result = result.filter(f => isCodeFile(f.name));
    } else if (activeFilter.value === 'archives') {
        result = result.filter(f => isArchiveFile(f.name));
    }

    // Sort
    result.sort((a, b) => {
        // Folders always first
        if (a.is_dir && !b.is_dir) return -1;
        if (!a.is_dir && b.is_dir) return 1;

        let cmp = 0;
        if (sortBy.value === 'name') {
            cmp = a.name.localeCompare(b.name, undefined, { numeric: true, sensitivity: 'base' });
        } else if (sortBy.value === 'size') {
            cmp = a.size_bytes - b.size_bytes;
        } else if (sortBy.value === 'modified') {
            cmp = new Date(a.modified_at).getTime() - new Date(b.modified_at).getTime();
        } else if (sortBy.value === 'permissions') {
            cmp = (a.permissions || '').localeCompare(b.permissions || '');
        }

        return sortAsc.value ? cmp : -cmp;
    });

    return result;
});

function toggleSort(column: 'name' | 'size' | 'modified' | 'permissions') {
    if (sortBy.value === column) {
        sortAsc.value = !sortAsc.value;
    } else {
        sortBy.value = column;
        sortAsc.value = true;
    }
}

// ==========================================
// SELECTION & BULK OPERATIONS
// ==========================================

const isAllSelected = computed(() => {
    return filteredFiles.value.length > 0 && selectedPaths.value.size === filteredFiles.value.length;
});

function toggleSelectAll() {
    if (isAllSelected.value) {
        selectedPaths.value.clear();
    } else {
        selectedPaths.value.clear();
        filteredFiles.value.forEach(f => selectedPaths.value.add(f.path));
    }
}

function toggleItemSelection(path: string, event?: MouseEvent) {
    if (event && (event.metaKey || event.ctrlKey)) {
        if (selectedPaths.value.has(path)) {
            selectedPaths.value.delete(path);
        } else {
            selectedPaths.value.add(path);
        }
    } else {
        if (selectedPaths.value.has(path) && selectedPaths.value.size === 1) {
            selectedPaths.value.clear();
        } else {
            selectedPaths.value.clear();
            selectedPaths.value.add(path);
        }
    }
}

function handleRowDoubleClick(entry: FileEntry) {
    if (entry.is_dir) {
        fetchDirectory(entry.path);
    } else if (isTextOrCodeFile(entry.name)) {
        openEditor(entry);
    } else if (isImageFile(entry.name) || isPdfFile(entry.name)) {
        openPreview(entry);
    } else if (isArchiveFile(entry.name)) {
        openExtractModal(entry);
    } else {
        openDetailsModal(entry);
    }
}

// ==========================================
// CONTEXT MENU
// ==========================================

function onContextMenu(event: MouseEvent, entry?: FileEntry) {
    event.preventDefault();
    event.stopPropagation();

    if (entry) {
        if (!selectedPaths.value.has(entry.path)) {
            selectedPaths.value.clear();
            selectedPaths.value.add(entry.path);
        }
        contextMenu.value.item = entry;
    } else {
        contextMenu.value.item = null;
    }

    // Keep menu inside viewport
    const menuWidth = 190;
    const menuHeight = 260;
    const posX = event.clientX + menuWidth > window.innerWidth ? window.innerWidth - menuWidth - 10 : event.clientX;
    const posY = event.clientY + menuHeight > window.innerHeight ? window.innerHeight - menuHeight - 10 : event.clientY;

    contextMenu.value.x = posX;
    contextMenu.value.y = posY;
    contextMenu.value.isOpen = true;
}

function closeContextMenu() {
    contextMenu.value.isOpen = false;
}

// ==========================================
// FILE ACTIONS (CREATE, RENAME, DELETE, ETC)
// ==========================================

// 1. Create File
function openCreateFileDialog() {
    newFileName.value = '';
    isCreateFileOpen.value = true;
}

async function handleCreateFile() {
    if (!newFileName.value.trim()) return;
    const rel = currentRelPath.value ? `${currentRelPath.value}/${newFileName.value.trim()}` : newFileName.value.trim();

    try {
        const res = await fetch('/files/create-file', {
            method: 'POST',
            headers: { 'Accept': 'application/json', 'Content-Type': 'application/json', 'X-CSRF-TOKEN': getCsrfToken() },
            body: JSON.stringify({ base_path: currentBasePath.value, relative_path: rel }),
        });
        const data = await res.json();
        if (data.success) {
            showToast('File created successfully', 'success');
            isCreateFileOpen.value = false;
            fetchDirectory(currentRelPath.value, false);
        } else {
            showToast(data.error || 'Failed to create file', 'error');
        }
    } catch (e) {
        showToast('Failed to create file', 'error');
    }
}

// 2. Create Folder
function openCreateFolderDialog() {
    newFolderName.value = '';
    isCreateFolderOpen.value = true;
}

async function handleCreateFolder() {
    if (!newFolderName.value.trim()) return;
    const rel = currentRelPath.value ? `${currentRelPath.value}/${newFolderName.value.trim()}` : newFolderName.value.trim();

    try {
        const res = await fetch('/files/create-folder', {
            method: 'POST',
            headers: { 'Accept': 'application/json', 'Content-Type': 'application/json', 'X-CSRF-TOKEN': getCsrfToken() },
            body: JSON.stringify({ base_path: currentBasePath.value, relative_path: rel }),
        });
        const data = await res.json();
        if (data.success) {
            showToast('Folder created successfully', 'success');
            isCreateFolderOpen.value = false;
            fetchDirectory(currentRelPath.value, false);
        } else {
            showToast(data.error || 'Failed to create folder', 'error');
        }
    } catch (e) {
        showToast('Failed to create folder', 'error');
    }
}

// 3. Rename
function openRenameModal(entry?: FileEntry) {
    const target = entry || getSingleSelectedItem();
    if (!target) return;
    renameTarget.value = target;
    renameNewName.value = target.name;
    isRenameOpen.value = true;
    closeContextMenu();
}

async function handleRename() {
    if (!renameTarget.value || !renameNewName.value.trim()) return;
    const oldRel = renameTarget.value.path;
    const parentDir = oldRel.includes('/') ? oldRel.substring(0, oldRel.lastIndexOf('/')) : '';
    const newRel = parentDir ? `${parentDir}/${renameNewName.value.trim()}` : renameNewName.value.trim();

    try {
        const res = await fetch('/files/rename', {
            method: 'POST',
            headers: { 'Accept': 'application/json', 'Content-Type': 'application/json', 'X-CSRF-TOKEN': getCsrfToken() },
            body: JSON.stringify({
                base_path: currentBasePath.value,
                old_path: oldRel,
                new_path: newRel,
            }),
        });
        const data = await res.json();
        if (data.success) {
            showToast('Renamed successfully', 'success');
            isRenameOpen.value = false;
            fetchDirectory(currentRelPath.value, false);
        } else {
            showToast(data.error || 'Failed to rename', 'error');
        }
    } catch (e) {
        showToast('Failed to rename', 'error');
    }
}

// 4. Copy / Move
function openCopyMoveModal(mode: 'copy' | 'move', entry?: FileEntry) {
    copyMoveMode.value = mode;
    copyMoveDestPath.value = currentRelPath.value;

    if (entry) {
        copyMoveItems.value = [entry.path];
    } else {
        copyMoveItems.value = Array.from(selectedPaths.value);
    }

    if (copyMoveItems.value.length === 0) return;
    isCopyMoveOpen.value = true;
    closeContextMenu();
}

async function handleCopyMove() {
    const endpoint = copyMoveMode.value === 'copy' ? '/files/copy' : '/files/move';
    try {
        const res = await fetch(endpoint, {
            method: 'POST',
            headers: { 'Accept': 'application/json', 'Content-Type': 'application/json', 'X-CSRF-TOKEN': getCsrfToken() },
            body: JSON.stringify({
                base_path: currentBasePath.value,
                sources: copyMoveItems.value,
                dest_path: copyMoveDestPath.value,
            }),
        });
        const data = await res.json();
        if (data.success) {
            showToast(`${copyMoveMode.value === 'copy' ? 'Copied' : 'Moved'} successfully`, 'success');
            isCopyMoveOpen.value = false;
            fetchDirectory(currentRelPath.value, false);
        } else {
            showToast(data.error || 'Operation failed', 'error');
        }
    } catch (e) {
        showToast('Operation failed', 'error');
    }
}

// 5. Delete
function openDeleteModal(entry?: FileEntry) {
    if (entry) {
        deleteTargets.value = [entry.path];
    } else {
        deleteTargets.value = Array.from(selectedPaths.value);
    }

    if (deleteTargets.value.length === 0) return;
    isDeleteOpen.value = true;
    closeContextMenu();
}

async function handleDelete() {
    try {
        const res = await fetch('/files/delete', {
            method: 'POST',
            headers: { 'Accept': 'application/json', 'Content-Type': 'application/json', 'X-CSRF-TOKEN': getCsrfToken() },
            body: JSON.stringify({
                base_path: currentBasePath.value,
                paths: deleteTargets.value,
            }),
        });
        const data = await res.json();
        if (data.success) {
            showToast('Item(s) deleted successfully', 'success');
            isDeleteOpen.value = false;
            fetchDirectory(currentRelPath.value, false);
        } else {
            showToast(data.error || 'Failed to delete', 'error');
        }
    } catch (e) {
        showToast('Failed to delete items', 'error');
    }
}

// 6. Compress
function openCompressModal(entry?: FileEntry) {
    if (entry) {
        compressSources.value = [entry.path];
        compressArchiveName.value = `${entry.name}.zip`;
    } else {
        compressSources.value = Array.from(selectedPaths.value);
        compressArchiveName.value = `archive-${Date.now().toString().slice(-4)}.zip`;
    }

    if (compressSources.value.length === 0) return;
    isCompressOpen.value = true;
    closeContextMenu();
}

async function handleCompress() {
    const dest = currentRelPath.value ? `${currentRelPath.value}/${compressArchiveName.value}` : compressArchiveName.value;
    try {
        const res = await fetch('/files/compress', {
            method: 'POST',
            headers: { 'Accept': 'application/json', 'Content-Type': 'application/json', 'X-CSRF-TOKEN': getCsrfToken() },
            body: JSON.stringify({
                base_path: currentBasePath.value,
                sources: compressSources.value,
                dest_path: dest,
                format: compressFormat.value,
            }),
        });
        const data = await res.json();
        if (data.success) {
            showToast('Archive created successfully', 'success');
            isCompressOpen.value = false;
            fetchDirectory(currentRelPath.value, false);
        } else {
            showToast(data.error || 'Failed to compress', 'error');
        }
    } catch (e) {
        showToast('Failed to create archive', 'error');
    }
}

// 7. Extract
function openExtractModal(entry?: FileEntry) {
    const target = entry || getSingleSelectedItem();
    if (!target) return;
    extractTarget.value = target;
    extractDestPath.value = currentRelPath.value;
    isExtractOpen.value = true;
    closeContextMenu();
}

async function handleExtract() {
    if (!extractTarget.value) return;
    try {
        const res = await fetch('/files/extract', {
            method: 'POST',
            headers: { 'Accept': 'application/json', 'Content-Type': 'application/json', 'X-CSRF-TOKEN': getCsrfToken() },
            body: JSON.stringify({
                base_path: currentBasePath.value,
                archive_path: extractTarget.value.path,
                dest_path: extractDestPath.value,
            }),
        });
        const data = await res.json();
        if (data.success) {
            showToast('Archive extracted successfully', 'success');
            isExtractOpen.value = false;
            fetchDirectory(currentRelPath.value, false);
        } else {
            showToast(data.error || 'Failed to extract archive', 'error');
        }
    } catch (e) {
        showToast('Failed to extract archive', 'error');
    }
}

// 8. Permissions (chmod)
function openPermissionsModal(entry?: FileEntry) {
    const target = entry || getSingleSelectedItem();
    if (!target) return;
    permissionsTarget.value = target;

    // Parse mode octal or symbolic
    const octal = target.mode_octal || '0644';
    const last3 = octal.slice(-3);
    const o = parseInt(last3[0] || '6', 10);
    const g = parseInt(last3[1] || '4', 10);
    const w = parseInt(last3[2] || '4', 10);

    permOwnerR.value = (o & 4) !== 0;
    permOwnerW.value = (o & 2) !== 0;
    permOwnerX.value = (o & 1) !== 0;

    permGroupR.value = (g & 4) !== 0;
    permGroupW.value = (g & 2) !== 0;
    permGroupX.value = (g & 1) !== 0;

    permOtherR.value = (w & 4) !== 0;
    permOtherW.value = (w & 2) !== 0;
    permOtherX.value = (w & 1) !== 0;

    permRecursive.value = false;
    isPermissionsOpen.value = true;
    closeContextMenu();
}

const calculatedOctal = computed(() => {
    let o = 0;
    if (permOwnerR.value) o += 4;
    if (permOwnerW.value) o += 2;
    if (permOwnerX.value) o += 1;

    let g = 0;
    if (permGroupR.value) g += 4;
    if (permGroupW.value) g += 2;
    if (permGroupX.value) g += 1;

    let w = 0;
    if (permOtherR.value) w += 4;
    if (permOtherW.value) w += 2;
    if (permOtherX.value) w += 1;

    return `0${o}${g}${w}`;
});

async function handleSavePermissions() {
    if (!permissionsTarget.value) return;
    try {
        const res = await fetch('/files/chmod', {
            method: 'POST',
            headers: { 'Accept': 'application/json', 'Content-Type': 'application/json', 'X-CSRF-TOKEN': getCsrfToken() },
            body: JSON.stringify({
                base_path: currentBasePath.value,
                relative_path: permissionsTarget.value.path,
                mode: calculatedOctal.value,
                recursive: permRecursive.value,
            }),
        });
        const data = await res.json();
        if (data.success) {
            showToast('Permissions updated', 'success');
            isPermissionsOpen.value = false;
            fetchDirectory(currentRelPath.value, false);
        } else {
            showToast(data.error || 'Failed to update permissions', 'error');
        }
    } catch (e) {
        showToast('Failed to update permissions', 'error');
    }
}

// 9. Details / Properties
async function openDetailsModal(entry?: FileEntry) {
    const target = entry || getSingleSelectedItem();
    if (!target) return;
    closeContextMenu();

    try {
        const res = await fetch('/files/stat', {
            method: 'POST',
            headers: { 'Accept': 'application/json', 'Content-Type': 'application/json', 'X-CSRF-TOKEN': getCsrfToken() },
            body: JSON.stringify({ base_path: currentBasePath.value, relative_path: target.path }),
        });
        const data = await res.json();
        if (data.success) {
            detailsData.value = data.data;
            isDetailsOpen.value = true;
        } else {
            showToast(data.error || 'Failed to get file details', 'error');
        }
    } catch (e) {
        showToast('Failed to get file details', 'error');
    }
}

// 10. Download (Single & Multi / Zip Stream)
function handleDownload(entry?: FileEntry) {
    let paths: string[] = [];
    if (entry) {
        paths = [entry.path];
    } else if (selectedPaths.value.size > 0) {
        paths = Array.from(selectedPaths.value);
    } else if (currentRelPath.value) {
        paths = [currentRelPath.value];
    }

    if (paths.length === 0) return;
    closeContextMenu();

    const queryParams = new URLSearchParams();
    queryParams.set('base_path', currentBasePath.value);
    paths.forEach(p => queryParams.append('paths[]', p));

    window.location.href = `/files/download?${queryParams.toString()}`;
}

// 11. File Upload (Drag & Drop + File Picker)
function triggerFileInput() {
    fileInputRef.value?.click();
}

function handleFileInputChange(e: Event) {
    const target = e.target as HTMLInputElement;
    if (target.files && target.files.length > 0) {
        uploadFiles(Array.from(target.files));
    }
}

function onDragOver(e: DragEvent) {
    e.preventDefault();
    isDraggingOver.value = true;
}

function onDragLeave(e: DragEvent) {
    e.preventDefault();
    isDraggingOver.value = false;
}

function onDrop(e: DragEvent) {
    e.preventDefault();
    isDraggingOver.value = false;
    if (e.dataTransfer?.files && e.dataTransfer.files.length > 0) {
        uploadFiles(Array.from(e.dataTransfer.files));
    }
}

async function uploadFiles(files: File[]) {
    if (files.length === 0) return;
    isUploading.value = true;
    uploadProgress.value = 0;
    uploadFileName.value = files.length === 1 ? files[0].name : `${files.length} files`;

    const formData = new FormData();
    formData.append('base_path', currentBasePath.value);
    formData.append('relative_path', currentRelPath.value);
    files.forEach(f => formData.append('files[]', f));

    try {
        const xhr = new XMLHttpRequest();
        xhr.open('POST', '/files/upload', true);
        xhr.setRequestHeader('Accept', 'application/json');
        xhr.setRequestHeader('X-CSRF-TOKEN', getCsrfToken());

        xhr.upload.onprogress = (e) => {
            if (e.lengthComputable) {
                uploadProgress.value = Math.round((e.loaded / e.total) * 100);
            }
        };

        xhr.onload = () => {
            isUploading.value = false;
            if (xhr.status === 200) {
                showToast('Files uploaded successfully', 'success');
                fetchDirectory(currentRelPath.value, false);
            } else {
                try {
                    const res = JSON.parse(xhr.responseText);
                    showToast(res.message || 'Upload failed', 'error');
                } catch {
                    showToast('Upload failed', 'error');
                }
            }
        };

        xhr.onerror = () => {
            isUploading.value = false;
            showToast('Upload network error', 'error');
        };

        xhr.send(formData);
    } catch (e) {
        isUploading.value = false;
        showToast('Upload failed', 'error');
    }
}

// 12. File Preview
async function openPreview(entry: FileEntry) {
    previewItem.value = entry;
    previewTextContent.value = '';
    closeContextMenu();

    if (isTextOrCodeFile(entry.name)) {
        try {
            const res = await fetch(`/files/read?base_path=${encodeURIComponent(currentBasePath.value)}&relative_path=${encodeURIComponent(entry.path)}`, {
                headers: { 'Accept': 'application/json' },
            });
            const data = await res.json();
            if (data.success) {
                previewTextContent.value = data.content;
            }
        } catch {}
    }

    isPreviewOpen.value = true;
}

const previewUrl = computed(() => {
    if (!previewItem.value) return '';
    return `/files/preview?base_path=${encodeURIComponent(currentBasePath.value)}&relative_path=${encodeURIComponent(previewItem.value.path)}`;
});

// 13. Professional Code Editor
async function openEditor(entry: FileEntry) {
    editingFile.value = entry;
    isEditorSaving.value = false;
    isEditorModified.value = false;
    closeContextMenu();

    try {
        const res = await fetch(`/files/read?base_path=${encodeURIComponent(currentBasePath.value)}&relative_path=${encodeURIComponent(entry.path)}`, {
            headers: { 'Accept': 'application/json' },
        });
        const data = await res.json();
        if (data.success) {
            editorContent.value = data.content;
            initialEditorContent.value = data.content;
            isEditorOpen.value = true;
        } else {
            showToast(data.error || 'Failed to read file', 'error');
        }
    } catch (e) {
        showToast('Failed to load file content', 'error');
    }
}

watch(editorContent, (newVal) => {
    isEditorModified.value = newVal !== initialEditorContent.value;
});

async function saveEditorContent() {
    if (!editingFile.value || isEditorSaving.value) return;
    isEditorSaving.value = true;

    try {
        const res = await fetch('/files/write', {
            method: 'POST',
            headers: { 'Accept': 'application/json', 'Content-Type': 'application/json', 'X-CSRF-TOKEN': getCsrfToken() },
            body: JSON.stringify({
                base_path: currentBasePath.value,
                relative_path: editingFile.value.path,
                content: editorContent.value,
            }),
        });
        const data = await res.json();
        if (data.success) {
            initialEditorContent.value = editorContent.value;
            isEditorModified.value = false;
            showToast('File saved successfully', 'success');
        } else {
            showToast(data.error || 'Failed to save file', 'error');
        }
    } catch (e) {
        showToast('Failed to save file', 'error');
    } finally {
        isEditorSaving.value = false;
    }
}

function closeEditor() {
    if (isEditorModified.value) {
        if (!confirm('You have unsaved changes. Are you sure you want to exit without saving?')) {
            return;
        }
    }
    isEditorOpen.value = false;
}

const editorLineCount = computed(() => {
    return Math.max(1, (editorContent.value.match(/\n/g) || []).length + 1);
});

function handleEditorReplaceAll() {
    if (!editorFindText.value) return;
    editorContent.value = editorContent.value.split(editorFindText.value).join(editorReplaceText.value);
}

// ==========================================
// KEYBOARD SHORTCUTS & CLICK OUTSIDE
// ==========================================

function handleKeyDown(e: KeyboardEvent) {
    // Ctrl+S or Cmd+S to save editor
    if ((e.ctrlKey || e.metaKey) && e.key === 's') {
        if (isEditorOpen.value) {
            e.preventDefault();
            saveEditorContent();
        }
    }
    // Escape to close modals/menus
    if (e.key === 'Escape') {
        if (contextMenu.value.isOpen) {
            closeContextMenu();
        }
    }
}

onMounted(() => {
    window.addEventListener('keydown', handleKeyDown);
    window.addEventListener('click', closeContextMenu);
});

onUnmounted(() => {
    window.removeEventListener('keydown', handleKeyDown);
    window.removeEventListener('click', closeContextMenu);
});

// ==========================================
// HELPERS
// ==========================================

function getSingleSelectedItem(): FileEntry | null {
    if (selectedPaths.value.size === 1) {
        const path = Array.from(selectedPaths.value)[0];
        return fileList.value.find(f => f.path === path) || null;
    }
    return null;
}

function getCsrfToken(): string {
    const meta = document.querySelector('meta[name="csrf-token"]') as HTMLMetaElement;
    return meta?.content || '';
}

function showToast(message: string, type: 'success' | 'error' = 'success') {
    // Custom non-blocking lightweight alert / notification
    const div = document.createElement('div');
    div.className = `fixed bottom-5 right-5 z-[9999] px-4 py-3 rounded-xl shadow-2xl text-xs font-semibold flex items-center gap-2 transition-all transform duration-300 translate-y-0 ${
        type === 'success'
            ? 'bg-emerald-600 text-white shadow-emerald-600/30'
            : 'bg-rose-600 text-white shadow-rose-600/30'
    }`;
    div.textContent = message;
    document.body.appendChild(div);
    setTimeout(() => {
        div.style.opacity = '0';
        div.style.transform = 'translateY(10px)';
        setTimeout(() => div.remove(), 300);
    }, 3000);
}

function formatBytes(bytes: number): string {
    if (!bytes || bytes === 0) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i];
}

function isImageFile(name: string): boolean {
    return /\.(jpg|jpeg|png|gif|webp|svg|ico)$/i.test(name);
}

function isCodeFile(name: string): boolean {
    return /\.(php|js|ts|vue|html|css|scss|json|yaml|yml|env|sql|sh|bash|py|go|rs|c|cpp|h)$/i.test(name);
}

function isArchiveFile(name: string): boolean {
    return /\.(zip|tar|gz|tgz|bz2|xz|rar|7z)$/i.test(name);
}

function isPdfFile(name: string): boolean {
    return /\.pdf$/i.test(name);
}

function isTextOrCodeFile(name: string): boolean {
    return isCodeFile(name) || /\.(txt|log|md|ini|conf|htaccess|gitignore|lock)$/i.test(name);
}

function getFileIcon(name: string, isDir: boolean) {
    if (isDir) return Folder;
    if (isImageFile(name)) return ImageIcon;
    if (isArchiveFile(name)) return FileArchive;
    if (/\.(php|js|ts|vue|html|css|json|yaml|sql)$/i.test(name)) return FileCode;
    if (/\.(sh|bash)$/i.test(name)) return Terminal;
    if (/\.(csv|xlsx|xls)$/i.test(name)) return FileSpreadsheet;
    return FileText;
}

function getFileIconColorClass(name: string, isDir: boolean) {
    if (isDir) return 'text-amber-500 bg-amber-500/10 border-amber-500/20';
    if (isImageFile(name)) return 'text-sky-500 bg-sky-500/10 border-sky-500/20';
    if (isArchiveFile(name)) return 'text-rose-500 bg-rose-500/10 border-rose-500/20';
    if (/\.php$/i.test(name)) return 'text-indigo-500 bg-indigo-500/10 border-indigo-500/20';
    if (/\.(js|ts|vue)$/i.test(name)) return 'text-yellow-500 bg-yellow-500/10 border-yellow-500/20';
    if (/\.(json|yaml|yml|env)$/i.test(name)) return 'text-emerald-500 bg-emerald-500/10 border-emerald-500/20';
    return 'text-slate-500 bg-slate-500/10 border-slate-500/20';
}
</script>

<template>
    <AppLayout title="File Manager">
        <div class="max-w-7xl mx-auto space-y-4" @dragover="onDragOver" @dragleave="onDragLeave" @drop="onDrop">
            <!-- Hidden file upload input -->
            <input
                ref="fileInputRef"
                type="file"
                multiple
                class="hidden"
                @change="handleFileInputChange"
            />

            <!-- Top Header & Website Root Switcher -->
            <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3 bg-white dark:bg-surface-900 border border-slate-200/80 dark:border-surface-800 rounded-2xl p-4 shadow-sm">
                <div class="flex items-center gap-3">
                    <div class="w-10 h-10 rounded-xl bg-brand-500/10 text-brand-600 dark:text-brand-400 border border-brand-500/20 flex items-center justify-center">
                        <HardDrive class="w-5 h-5" />
                    </div>
                    <div>
                        <div class="flex items-center gap-2">
                            <h2 class="text-sm font-bold text-slate-900 dark:text-white tracking-tight">Kodepreneur File Manager</h2>
                            <span class="px-2 py-0.5 rounded-full text-[10px] font-mono font-medium bg-emerald-50 text-emerald-600 border border-emerald-200 dark:bg-emerald-500/10 dark:text-emerald-400 dark:border-emerald-500/20">
                                Jailed Sandbox
                            </span>
                        </div>
                        <p class="text-xs text-slate-500 dark:text-surface-400 mt-0.5">
                            Manage webroot assets, configurations, permissions, and archives safely.
                        </p>
                    </div>
                </div>

                <div class="flex items-center gap-3">
                    <!-- Website Picker -->
                    <div class="flex items-center gap-2 bg-slate-50 dark:bg-surface-950 border border-slate-200/80 dark:border-surface-800 rounded-xl px-3 py-1.5 shadow-sm">
                        <Globe class="w-3.5 h-3.5 text-brand-600 dark:text-brand-400" />
                        <select
                            :value="activeWebsite?.id"
                            @change="onWebsiteChange"
                            class="bg-transparent text-xs text-slate-800 dark:text-white focus:outline-none font-mono font-medium"
                        >
                            <option v-for="w in websites" :key="w.id" :value="w.id" class="bg-white dark:bg-surface-900 text-slate-900 dark:text-white">
                                {{ w.domain }}
                            </option>
                        </select>
                    </div>
                </div>
            </div>

            <!-- Main Explorer Card -->
            <div class="rounded-2xl bg-white dark:bg-surface-900 border border-slate-200/80 dark:border-surface-800 shadow-sm overflow-hidden flex flex-col">
                <!-- Navigation & Breadcrumb Bar -->
                <div class="p-3 bg-slate-50/80 dark:bg-surface-950/80 border-b border-slate-200/80 dark:border-surface-800 flex flex-wrap items-center justify-between gap-3 text-xs">
                    <!-- Nav buttons & Breadcrumbs -->
                    <div class="flex items-center gap-1.5 overflow-x-auto py-1 max-w-full">
                        <div class="flex items-center gap-1 shrink-0 bg-white dark:bg-surface-900 border border-slate-200 dark:border-surface-800 rounded-xl p-0.5">
                            <button
                                @click="navigateBack"
                                :disabled="historyIndex <= 0"
                                class="p-1.5 rounded-lg hover:bg-slate-100 dark:hover:bg-surface-800 text-slate-600 dark:text-surface-300 disabled:opacity-30 disabled:hover:bg-transparent transition"
                                title="Back"
                            >
                                <ArrowLeft class="w-3.5 h-3.5" />
                            </button>
                            <button
                                @click="navigateForward"
                                :disabled="historyIndex >= historyStack.length - 1"
                                class="p-1.5 rounded-lg hover:bg-slate-100 dark:hover:bg-surface-800 text-slate-600 dark:text-surface-300 disabled:opacity-30 disabled:hover:bg-transparent transition"
                                title="Forward"
                            >
                                <ArrowRight class="w-3.5 h-3.5" />
                            </button>
                            <button
                                @click="goUp"
                                :disabled="!currentRelPath"
                                class="p-1.5 rounded-lg hover:bg-slate-100 dark:hover:bg-surface-800 text-slate-600 dark:text-surface-300 disabled:opacity-30 disabled:hover:bg-transparent transition"
                                title="Up One Level"
                            >
                                <CornerLeftUp class="w-3.5 h-3.5" />
                            </button>
                            <button
                                @click="refreshDirectory"
                                class="p-1.5 rounded-lg hover:bg-slate-100 dark:hover:bg-surface-800 text-slate-600 dark:text-surface-300 transition"
                                title="Refresh"
                            >
                                <RefreshCw :class="['w-3.5 h-3.5', isLoading ? 'animate-spin text-brand-500' : '']" />
                            </button>
                        </div>

                        <!-- Breadcrumb Path Trail -->
                        <div class="flex items-center gap-1 bg-white dark:bg-surface-900 border border-slate-200 dark:border-surface-800 rounded-xl px-2.5 py-1 font-mono text-xs shadow-inner shrink-0">
                            <template v-for="(seg, idx) in breadcrumbSegments" :key="seg.path">
                                <span v-if="idx > 0" class="text-slate-300 dark:text-surface-700">/</span>
                                <button
                                    @click="fetchDirectory(seg.path)"
                                    :class="[
                                        'hover:underline transition px-1 py-0.5 rounded',
                                        idx === breadcrumbSegments.length - 1
                                            ? 'font-bold text-slate-900 dark:text-white'
                                            : 'text-brand-600 dark:text-brand-400 font-medium'
                                    ]"
                                >
                                    {{ seg.name }}
                                </button>
                            </template>
                        </div>
                    </div>

                    <!-- Search & Hidden files toggle -->
                    <div class="flex items-center gap-2 ml-auto shrink-0">
                        <div class="relative flex items-center">
                            <Search class="w-3.5 h-3.5 text-slate-400 absolute left-2.5 pointer-events-none" />
                            <input
                                v-model="searchQuery"
                                type="text"
                                placeholder="Search current view..."
                                class="bg-white dark:bg-surface-900 border border-slate-200 dark:border-surface-800 rounded-xl pl-8 pr-3 py-1.5 text-xs text-slate-900 dark:text-white placeholder-slate-400 focus:outline-none focus:ring-1 focus:ring-brand-500 w-44 sm:w-56"
                            />
                            <button
                                v-if="searchQuery"
                                @click="searchQuery = ''"
                                class="absolute right-2 text-slate-400 hover:text-slate-600 dark:hover:text-surface-200"
                            >
                                <X class="w-3 h-3" />
                            </button>
                        </div>

                        <button
                            @click="toggleShowHidden"
                            :class="[
                                'px-2.5 py-1.5 rounded-xl border text-xs font-semibold flex items-center gap-1.5 transition',
                                isShowHidden
                                    ? 'bg-brand-50 border-brand-300 text-brand-600 dark:bg-brand-500/10 dark:border-brand-500/30 dark:text-brand-400'
                                    : 'bg-white dark:bg-surface-900 border-slate-200 dark:border-surface-800 text-slate-600 dark:text-surface-400 hover:bg-slate-50 dark:hover:bg-surface-800'
                            ]"
                            title="Toggle Hidden Files (.env, .git, etc.)"
                        >
                            <Eye class="w-3.5 h-3.5" />
                            <span class="hidden sm:inline">Hidden</span>
                        </button>
                    </div>
                </div>

                <!-- Main Action Toolbar -->
                <div class="p-2.5 bg-white dark:bg-surface-900 border-b border-slate-200/80 dark:border-surface-800 flex flex-wrap items-center justify-between gap-2 text-xs">
                    <!-- Left Operation Buttons -->
                    <div class="flex items-center gap-1.5 flex-wrap">
                        <button
                            @click="openCreateFileDialog"
                            class="px-3 py-1.5 rounded-xl bg-brand-600 hover:bg-brand-500 text-white font-semibold flex items-center gap-1.5 transition shadow-sm"
                        >
                            <FilePlus class="w-3.5 h-3.5" />
                            <span>New File</span>
                        </button>
                        <button
                            @click="openCreateFolderDialog"
                            class="px-3 py-1.5 rounded-xl bg-slate-100 dark:bg-surface-800 hover:bg-slate-200 dark:hover:bg-surface-700 text-slate-800 dark:text-surface-100 font-semibold flex items-center gap-1.5 transition"
                        >
                            <FolderPlus class="w-3.5 h-3.5 text-amber-500" />
                            <span>New Folder</span>
                        </button>
                        <button
                            @click="triggerFileInput"
                            class="px-3 py-1.5 rounded-xl bg-slate-100 dark:bg-surface-800 hover:bg-slate-200 dark:hover:bg-surface-700 text-slate-800 dark:text-surface-100 font-semibold flex items-center gap-1.5 transition"
                        >
                            <Upload class="w-3.5 h-3.5 text-sky-500" />
                            <span>Upload</span>
                        </button>
                        <button
                            @click="handleDownload()"
                            class="px-3 py-1.5 rounded-xl bg-slate-100 dark:bg-surface-800 hover:bg-slate-200 dark:hover:bg-surface-700 text-slate-800 dark:text-surface-100 font-semibold flex items-center gap-1.5 transition"
                        >
                            <Download class="w-3.5 h-3.5 text-emerald-500" />
                            <span>Download</span>
                        </button>

                        <div class="h-4 w-px bg-slate-200 dark:bg-surface-800 mx-1 hidden sm:block"></div>

                        <!-- Actions when item(s) selected -->
                        <template v-if="selectedPaths.size > 0">
                            <button
                                @click="openCopyMoveModal('copy')"
                                class="px-2.5 py-1.5 rounded-xl bg-slate-100 dark:bg-surface-800 hover:bg-slate-200 dark:hover:bg-surface-700 text-slate-700 dark:text-surface-200 font-medium flex items-center gap-1 transition"
                            >
                                <Copy class="w-3.5 h-3.5" />
                                <span>Copy</span>
                            </button>
                            <button
                                @click="openCopyMoveModal('move')"
                                class="px-2.5 py-1.5 rounded-xl bg-slate-100 dark:bg-surface-800 hover:bg-slate-200 dark:hover:bg-surface-700 text-slate-700 dark:text-surface-200 font-medium flex items-center gap-1 transition"
                            >
                                <Scissors class="w-3.5 h-3.5" />
                                <span>Move</span>
                            </button>
                            <button
                                v-if="selectedPaths.size === 1"
                                @click="openRenameModal()"
                                class="px-2.5 py-1.5 rounded-xl bg-slate-100 dark:bg-surface-800 hover:bg-slate-200 dark:hover:bg-surface-700 text-slate-700 dark:text-surface-200 font-medium flex items-center gap-1 transition"
                            >
                                <Edit3 class="w-3.5 h-3.5" />
                                <span>Rename</span>
                            </button>
                            <button
                                @click="openCompressModal()"
                                class="px-2.5 py-1.5 rounded-xl bg-slate-100 dark:bg-surface-800 hover:bg-slate-200 dark:hover:bg-surface-700 text-slate-700 dark:text-surface-200 font-medium flex items-center gap-1 transition"
                            >
                                <Archive class="w-3.5 h-3.5 text-rose-500" />
                                <span>Compress</span>
                            </button>
                            <button
                                v-if="selectedPaths.size === 1"
                                @click="openPermissionsModal()"
                                class="px-2.5 py-1.5 rounded-xl bg-slate-100 dark:bg-surface-800 hover:bg-slate-200 dark:hover:bg-surface-700 text-slate-700 dark:text-surface-200 font-medium flex items-center gap-1 transition"
                            >
                                <Shield class="w-3.5 h-3.5" />
                                <span>Permissions</span>
                            </button>
                            <button
                                @click="openDeleteModal()"
                                class="px-2.5 py-1.5 rounded-xl bg-rose-500/10 hover:bg-rose-500/20 text-rose-600 dark:text-rose-400 font-medium flex items-center gap-1 transition"
                            >
                                <Trash2 class="w-3.5 h-3.5" />
                                <span>Delete ({{ selectedPaths.size }})</span>
                            </button>
                        </template>
                    </div>

                    <!-- Category Filter Pills -->
                    <div class="flex items-center gap-1 overflow-x-auto">
                        <button
                            v-for="filter in ['all', 'folders', 'files', 'code', 'images', 'archives']"
                            :key="filter"
                            @click="activeFilter = filter"
                            :class="[
                                'px-2.5 py-1 rounded-lg capitalize font-medium transition text-[11px]',
                                activeFilter === filter
                                    ? 'bg-slate-900 text-white dark:bg-white dark:text-slate-900'
                                    : 'text-slate-500 hover:bg-slate-100 dark:text-surface-400 dark:hover:bg-surface-800'
                            ]"
                        >
                            {{ filter }}
                        </button>
                    </div>
                </div>

                <!-- Drag & Drop Upload Overlay -->
                <div
                    v-if="isDraggingOver"
                    class="p-10 border-2 border-dashed border-brand-500 bg-brand-500/10 rounded-2xl m-4 flex flex-col items-center justify-center text-center transition animate-pulse"
                >
                    <Upload class="w-12 h-12 text-brand-600 dark:text-brand-400 mb-2" />
                    <p class="text-sm font-bold text-slate-900 dark:text-white">Drop files here to upload</p>
                    <p class="text-xs text-slate-500 dark:text-surface-400 mt-0.5">Files will be placed in {{ currentRelPath || 'root' }}</p>
                </div>

                <!-- Upload Progress Bar Indicator -->
                <div v-if="isUploading" class="p-3 bg-brand-500/10 border-b border-brand-500/20 flex items-center justify-between gap-4 text-xs font-semibold">
                    <div class="flex items-center gap-2 text-brand-700 dark:text-brand-300">
                        <Upload class="w-4 h-4 animate-bounce" />
                        <span>Uploading {{ uploadFileName }} ({{ uploadProgress }}%)...</span>
                    </div>
                    <div class="w-48 bg-slate-200 dark:bg-surface-800 rounded-full h-2 overflow-hidden">
                        <div class="bg-brand-600 h-full transition-all duration-150" :style="{ width: `${uploadProgress}%` }"></div>
                    </div>
                </div>

                <!-- File Table View -->
                <div
                    class="flex-1 overflow-x-auto min-h-[360px]"
                    @contextmenu="onContextMenu($event)"
                >
                    <div v-if="filteredFiles.length === 0" class="text-center py-16 px-4">
                        <Folder class="w-12 h-12 text-slate-300 dark:text-surface-700 mx-auto mb-3" />
                        <h3 class="text-sm font-semibold text-slate-700 dark:text-surface-200">No items found</h3>
                        <p class="text-xs text-slate-500 dark:text-surface-400 mt-1 max-w-sm mx-auto">
                            {{ searchQuery ? 'No files match your search criteria.' : 'This directory is empty. Upload or create new files to get started.' }}
                        </p>
                    </div>

                    <table v-else class="w-full text-left text-xs">
                        <thead class="bg-slate-50/90 dark:bg-surface-950/50 text-slate-500 dark:text-surface-400 uppercase text-[10px] tracking-wider border-b border-slate-200/80 dark:border-surface-800 select-none">
                            <tr>
                                <th class="py-3 px-3 w-8 text-center">
                                    <input
                                        type="checkbox"
                                        :checked="isAllSelected"
                                        @change="toggleSelectAll"
                                        class="rounded border-slate-300 text-brand-600 focus:ring-brand-500 cursor-pointer"
                                    />
                                </th>
                                <th
                                    @click="toggleSort('name')"
                                    class="py-3 px-3 font-semibold cursor-pointer hover:text-slate-900 dark:hover:text-white"
                                >
                                    <div class="flex items-center gap-1">
                                        <span>Name</span>
                                        <span v-if="sortBy === 'name'">{{ sortAsc ? '↑' : '↓' }}</span>
                                    </div>
                                </th>
                                <th
                                    @click="toggleSort('size')"
                                    class="py-3 px-3 font-semibold cursor-pointer hover:text-slate-900 dark:hover:text-white"
                                >
                                    <div class="flex items-center gap-1">
                                        <span>Size</span>
                                        <span v-if="sortBy === 'size'">{{ sortAsc ? '↑' : '↓' }}</span>
                                    </div>
                                </th>
                                <th
                                    @click="toggleSort('permissions')"
                                    class="py-3 px-3 font-semibold cursor-pointer hover:text-slate-900 dark:hover:text-white"
                                >
                                    <div class="flex items-center gap-1">
                                        <span>Permissions</span>
                                        <span v-if="sortBy === 'permissions'">{{ sortAsc ? '↑' : '↓' }}</span>
                                    </div>
                                </th>
                                <th class="py-3 px-3 font-semibold hidden md:table-cell">Owner:Group</th>
                                <th
                                    @click="toggleSort('modified')"
                                    class="py-3 px-3 font-semibold cursor-pointer hover:text-slate-900 dark:hover:text-white"
                                >
                                    <div class="flex items-center gap-1">
                                        <span>Modified</span>
                                        <span v-if="sortBy === 'modified'">{{ sortAsc ? '↑' : '↓' }}</span>
                                    </div>
                                </th>
                                <th class="py-3 px-3 font-semibold text-right">Actions</th>
                            </tr>
                        </thead>
                        <tbody class="divide-y divide-slate-100 dark:divide-surface-800/60">
                            <tr
                                v-for="file in filteredFiles"
                                :key="file.path"
                                @click="toggleItemSelection(file.path, $event)"
                                @dblclick="handleRowDoubleClick(file)"
                                @contextmenu="onContextMenu($event, file)"
                                :class="[
                                    'group cursor-pointer select-none transition-colors duration-100',
                                    selectedPaths.has(file.path)
                                        ? 'bg-brand-50/70 dark:bg-brand-500/10'
                                        : 'hover:bg-slate-50/80 dark:hover:bg-surface-800/30'
                                ]"
                            >
                                <td class="py-2.5 px-3 text-center" @click.stop>
                                    <input
                                        type="checkbox"
                                        :checked="selectedPaths.has(file.path)"
                                        @change="toggleItemSelection(file.path)"
                                        class="rounded border-slate-300 text-brand-600 focus:ring-brand-500 cursor-pointer"
                                    />
                                </td>
                                <td class="py-2.5 px-3 font-mono font-medium text-slate-900 dark:text-white">
                                    <div class="flex items-center gap-2.5">
                                        <div
                                            @click.stop="file.is_dir ? fetchDirectory(file.path) : null"
                                            :class="['p-1.5 rounded-lg shrink-0 border transition-transform', file.is_dir ? 'cursor-pointer hover:scale-105 hover:shadow-sm' : '', getFileIconColorClass(file.name, file.is_dir)]"
                                            :title="file.is_dir ? 'Open folder' : ''"
                                        >
                                            <component :is="getFileIcon(file.name, file.is_dir)" class="w-4 h-4" />
                                        </div>
                                        <button
                                            v-if="file.is_dir"
                                            type="button"
                                            @click.stop="fetchDirectory(file.path)"
                                            class="text-amber-700 dark:text-amber-300 hover:underline font-semibold text-left truncate max-w-xs sm:max-w-md focus:outline-none"
                                            title="Open folder"
                                        >
                                            {{ file.name }}
                                        </button>
                                        <span v-else class="text-slate-800 dark:text-surface-200 truncate max-w-xs sm:max-w-md">
                                            {{ file.name }}
                                        </span>
                                    </div>
                                </td>
                                <td class="py-2.5 px-3 font-mono text-[11px] text-slate-500 dark:text-surface-400">
                                    {{ file.is_dir ? `${file.item_count ?? '—'} items` : formatBytes(file.size_bytes) }}
                                </td>
                                <td class="py-2.5 px-3 font-mono text-[11px] text-slate-400 dark:text-surface-500">
                                    <span class="px-1.5 py-0.5 rounded bg-slate-100 dark:bg-surface-800 text-slate-600 dark:text-surface-300 font-semibold">
                                        {{ file.mode_octal || '0644' }}
                                    </span>
                                </td>
                                <td class="py-2.5 px-3 font-mono text-[11px] text-slate-400 dark:text-surface-500 hidden md:table-cell">
                                    {{ file.owner || 'www-data' }}:{{ file.group || 'www-data' }}
                                </td>
                                <td class="py-2.5 px-3 font-mono text-[11px] text-slate-400 dark:text-surface-500 whitespace-nowrap">
                                    {{ new Date(file.modified_at).toLocaleString() }}
                                </td>
                                <td class="py-2.5 px-3 text-right" @click.stop>
                                    <div class="flex items-center justify-end gap-1 opacity-80 group-hover:opacity-100">
                                        <button
                                            v-if="!file.is_dir && isTextOrCodeFile(file.name)"
                                            @click="openEditor(file)"
                                            class="p-1 rounded-lg text-slate-500 hover:text-slate-900 hover:bg-slate-100 dark:text-surface-300 dark:hover:text-white dark:hover:bg-surface-800 transition"
                                            title="Edit File"
                                        >
                                            <Edit3 class="w-3.5 h-3.5" />
                                        </button>
                                        <button
                                            v-if="!file.is_dir && (isImageFile(file.name) || isPdfFile(file.name))"
                                            @click="openPreview(file)"
                                            class="p-1 rounded-lg text-slate-500 hover:text-slate-900 hover:bg-slate-100 dark:text-surface-300 dark:hover:text-white dark:hover:bg-surface-800 transition"
                                            title="Preview"
                                        >
                                            <Eye class="w-3.5 h-3.5" />
                                        </button>
                                        <button
                                            @click="handleDownload(file)"
                                            class="p-1 rounded-lg text-slate-500 hover:text-slate-900 hover:bg-slate-100 dark:text-surface-300 dark:hover:text-white dark:hover:bg-surface-800 transition"
                                            title="Download"
                                        >
                                            <Download class="w-3.5 h-3.5" />
                                        </button>
                                        <button
                                            @click="openDeleteModal(file)"
                                            class="p-1 rounded-lg text-rose-500 hover:bg-rose-50 dark:hover:bg-rose-500/10 transition"
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

                <!-- Footer Storage & Metadata Info Bar -->
                <div class="p-3 bg-slate-50 dark:bg-surface-950 border-t border-slate-200/80 dark:border-surface-800 flex flex-col sm:flex-row sm:items-center justify-between gap-3 text-xs text-slate-500 dark:text-surface-400">
                    <div class="flex items-center gap-4">
                        <span><strong>{{ filteredFiles.length }}</strong> item(s)</span>
                        <span v-if="selectedPaths.size > 0" class="text-brand-600 dark:text-brand-400 font-semibold">
                            {{ selectedPaths.size }} selected
                        </span>
                        <span v-if="diskInfo?.path_size">Folder size: {{ formatBytes(diskInfo.path_size) }}</span>
                    </div>

                    <!-- Disk Usage Partition Bar -->
                    <div v-if="diskInfo" class="flex items-center gap-3">
                        <div class="flex items-center gap-2">
                            <span>Disk Usage:</span>
                            <span class="font-mono text-slate-900 dark:text-white font-semibold">
                                {{ formatBytes(diskInfo.used_bytes) }} / {{ formatBytes(diskInfo.total_bytes) }}
                            </span>
                        </div>
                        <div class="w-28 bg-slate-200 dark:bg-surface-800 rounded-full h-1.5 overflow-hidden">
                            <div
                                :class="[
                                    'h-full transition-all duration-300',
                                    diskInfo.usage_percent > 85 ? 'bg-rose-500' : 'bg-brand-500'
                                ]"
                                :style="{ width: `${diskInfo.usage_percent}%` }"
                            ></div>
                        </div>
                    </div>
                </div>
            </div>

            <!-- Custom Right-Click Context Menu -->
            <div
                v-if="contextMenu.isOpen"
                :style="{ top: `${contextMenu.y}px`, left: `${contextMenu.x}px` }"
                class="fixed z-50 w-48 rounded-xl bg-white dark:bg-surface-900 border border-slate-200 dark:border-surface-800 shadow-2xl py-1 text-xs text-slate-700 dark:text-surface-200 divide-y divide-slate-100 dark:divide-surface-800"
                @click.stop
            >
                <div class="py-1">
                    <button
                        v-if="contextMenu.item?.is_dir"
                        @click="fetchDirectory(contextMenu.item.path)"
                        class="w-full px-3 py-1.5 text-left hover:bg-slate-100 dark:hover:bg-surface-800 flex items-center gap-2 font-medium"
                    >
                        <Folder class="w-3.5 h-3.5 text-amber-500" />
                        <span>Open Folder</span>
                    </button>
                    <button
                        v-if="contextMenu.item && !contextMenu.item.is_dir && isTextOrCodeFile(contextMenu.item.name)"
                        @click="openEditor(contextMenu.item)"
                        class="w-full px-3 py-1.5 text-left hover:bg-slate-100 dark:hover:bg-surface-800 flex items-center gap-2 font-medium"
                    >
                        <Edit3 class="w-3.5 h-3.5 text-brand-500" />
                        <span>Edit File</span>
                    </button>
                    <button
                        v-if="contextMenu.item && !contextMenu.item.is_dir && (isImageFile(contextMenu.item.name) || isPdfFile(contextMenu.item.name))"
                        @click="openPreview(contextMenu.item)"
                        class="w-full px-3 py-1.5 text-left hover:bg-slate-100 dark:hover:bg-surface-800 flex items-center gap-2 font-medium"
                    >
                        <Eye class="w-3.5 h-3.5 text-sky-500" />
                        <span>Preview</span>
                    </button>
                    <button
                        v-if="contextMenu.item"
                        @click="handleDownload(contextMenu.item)"
                        class="w-full px-3 py-1.5 text-left hover:bg-slate-100 dark:hover:bg-surface-800 flex items-center gap-2"
                    >
                        <Download class="w-3.5 h-3.5 text-emerald-500" />
                        <span>Download</span>
                    </button>
                </div>

                <div class="py-1">
                    <button
                        @click="openCopyMoveModal('copy', contextMenu.item || undefined)"
                        class="w-full px-3 py-1.5 text-left hover:bg-slate-100 dark:hover:bg-surface-800 flex items-center gap-2"
                    >
                        <Copy class="w-3.5 h-3.5" />
                        <span>Copy</span>
                    </button>
                    <button
                        @click="openCopyMoveModal('move', contextMenu.item || undefined)"
                        class="w-full px-3 py-1.5 text-left hover:bg-slate-100 dark:hover:bg-surface-800 flex items-center gap-2"
                    >
                        <Scissors class="w-3.5 h-3.5" />
                        <span>Move</span>
                    </button>
                    <button
                        v-if="contextMenu.item"
                        @click="openRenameModal(contextMenu.item)"
                        class="w-full px-3 py-1.5 text-left hover:bg-slate-100 dark:hover:bg-surface-800 flex items-center gap-2"
                    >
                        <Edit3 class="w-3.5 h-3.5" />
                        <span>Rename</span>
                    </button>
                </div>

                <div class="py-1">
                    <button
                        @click="openCompressModal(contextMenu.item || undefined)"
                        class="w-full px-3 py-1.5 text-left hover:bg-slate-100 dark:hover:bg-surface-800 flex items-center gap-2"
                    >
                        <Archive class="w-3.5 h-3.5 text-rose-500" />
                        <span>Compress</span>
                    </button>
                    <button
                        v-if="contextMenu.item && isArchiveFile(contextMenu.item.name)"
                        @click="openExtractModal(contextMenu.item)"
                        class="w-full px-3 py-1.5 text-left hover:bg-slate-100 dark:hover:bg-surface-800 flex items-center gap-2 text-indigo-500"
                    >
                        <Layers class="w-3.5 h-3.5" />
                        <span>Extract Archive</span>
                    </button>
                    <button
                        v-if="contextMenu.item"
                        @click="openPermissionsModal(contextMenu.item)"
                        class="w-full px-3 py-1.5 text-left hover:bg-slate-100 dark:hover:bg-surface-800 flex items-center gap-2"
                    >
                        <Shield class="w-3.5 h-3.5" />
                        <span>Permissions</span>
                    </button>
                    <button
                        v-if="contextMenu.item"
                        @click="openDetailsModal(contextMenu.item)"
                        class="w-full px-3 py-1.5 text-left hover:bg-slate-100 dark:hover:bg-surface-800 flex items-center gap-2"
                    >
                        <Info class="w-3.5 h-3.5" />
                        <span>Properties</span>
                    </button>
                </div>

                <div class="py-1">
                    <button
                        @click="openDeleteModal(contextMenu.item || undefined)"
                        class="w-full px-3 py-1.5 text-left hover:bg-rose-500/10 text-rose-600 dark:text-rose-400 flex items-center gap-2 font-medium"
                    >
                        <Trash2 class="w-3.5 h-3.5" />
                        <span>Delete</span>
                    </button>
                </div>
            </div>

            <!-- ========================================== -->
            <!-- MODALS & DIALOGS -->
            <!-- ========================================== -->

            <!-- 1. Create File Modal -->
            <div v-if="isCreateFileOpen" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/60 backdrop-blur-sm">
                <div class="w-full max-w-md bg-white dark:bg-surface-900 rounded-2xl border border-slate-200 dark:border-surface-800 shadow-2xl p-5 space-y-4">
                    <div class="flex items-center justify-between">
                        <h3 class="text-sm font-bold text-slate-900 dark:text-white flex items-center gap-2">
                            <FilePlus class="w-4 h-4 text-brand-600" />
                            <span>Create New File</span>
                        </h3>
                        <button @click="isCreateFileOpen = false" class="p-1 rounded-lg text-slate-400 hover:text-slate-600 dark:hover:text-white">
                            <X class="w-4 h-4" />
                        </button>
                    </div>
                    <div>
                        <label class="block text-xs font-semibold text-slate-700 dark:text-surface-300 mb-1.5">File Name</label>
                        <input
                            v-model="newFileName"
                            type="text"
                            placeholder="e.g. index.php, config.json"
                            autofocus
                            @keyup.enter="handleCreateFile"
                            class="w-full bg-slate-50 dark:bg-surface-950 border border-slate-200 dark:border-surface-800 rounded-xl px-3 py-2 text-xs text-slate-900 dark:text-white font-mono focus:outline-none focus:ring-2 focus:ring-brand-500"
                        />
                    </div>
                    <div class="flex justify-end gap-2">
                        <button @click="isCreateFileOpen = false" class="px-4 py-2 rounded-xl bg-slate-100 dark:bg-surface-800 text-xs font-semibold text-slate-700 dark:text-surface-300 hover:bg-slate-200 transition">
                            Cancel
                        </button>
                        <button @click="handleCreateFile" class="px-4 py-2 rounded-xl bg-brand-600 hover:bg-brand-500 text-white text-xs font-semibold transition shadow-md shadow-brand-600/20">
                            Create File
                        </button>
                    </div>
                </div>
            </div>

            <!-- 2. Create Folder Modal -->
            <div v-if="isCreateFolderOpen" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/60 backdrop-blur-sm">
                <div class="w-full max-w-md bg-white dark:bg-surface-900 rounded-2xl border border-slate-200 dark:border-surface-800 shadow-2xl p-5 space-y-4">
                    <div class="flex items-center justify-between">
                        <h3 class="text-sm font-bold text-slate-900 dark:text-white flex items-center gap-2">
                            <FolderPlus class="w-4 h-4 text-amber-500" />
                            <span>Create New Folder</span>
                        </h3>
                        <button @click="isCreateFolderOpen = false" class="p-1 rounded-lg text-slate-400 hover:text-slate-600 dark:hover:text-white">
                            <X class="w-4 h-4" />
                        </button>
                    </div>
                    <div>
                        <label class="block text-xs font-semibold text-slate-700 dark:text-surface-300 mb-1.5">Folder Name</label>
                        <input
                            v-model="newFolderName"
                            type="text"
                            placeholder="e.g. assets, uploads"
                            autofocus
                            @keyup.enter="handleCreateFolder"
                            class="w-full bg-slate-50 dark:bg-surface-950 border border-slate-200 dark:border-surface-800 rounded-xl px-3 py-2 text-xs text-slate-900 dark:text-white font-mono focus:outline-none focus:ring-2 focus:ring-brand-500"
                        />
                    </div>
                    <div class="flex justify-end gap-2">
                        <button @click="isCreateFolderOpen = false" class="px-4 py-2 rounded-xl bg-slate-100 dark:bg-surface-800 text-xs font-semibold text-slate-700 dark:text-surface-300 hover:bg-slate-200 transition">
                            Cancel
                        </button>
                        <button @click="handleCreateFolder" class="px-4 py-2 rounded-xl bg-brand-600 hover:bg-brand-500 text-white text-xs font-semibold transition shadow-md shadow-brand-600/20">
                            Create Folder
                        </button>
                    </div>
                </div>
            </div>

            <!-- 3. Rename Modal -->
            <div v-if="isRenameOpen" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/60 backdrop-blur-sm">
                <div class="w-full max-w-md bg-white dark:bg-surface-900 rounded-2xl border border-slate-200 dark:border-surface-800 shadow-2xl p-5 space-y-4">
                    <div class="flex items-center justify-between">
                        <h3 class="text-sm font-bold text-slate-900 dark:text-white flex items-center gap-2">
                            <Edit3 class="w-4 h-4 text-brand-600" />
                            <span>Rename Item</span>
                        </h3>
                        <button @click="isRenameOpen = false" class="p-1 rounded-lg text-slate-400 hover:text-slate-600 dark:hover:text-white">
                            <X class="w-4 h-4" />
                        </button>
                    </div>
                    <div>
                        <label class="block text-xs font-semibold text-slate-700 dark:text-surface-300 mb-1.5">New Name</label>
                        <input
                            v-model="renameNewName"
                            type="text"
                            autofocus
                            @keyup.enter="handleRename"
                            class="w-full bg-slate-50 dark:bg-surface-950 border border-slate-200 dark:border-surface-800 rounded-xl px-3 py-2 text-xs text-slate-900 dark:text-white font-mono focus:outline-none focus:ring-2 focus:ring-brand-500"
                        />
                    </div>
                    <div class="flex justify-end gap-2">
                        <button @click="isRenameOpen = false" class="px-4 py-2 rounded-xl bg-slate-100 dark:bg-surface-800 text-xs font-semibold text-slate-700 dark:text-surface-300 hover:bg-slate-200 transition">
                            Cancel
                        </button>
                        <button @click="handleRename" class="px-4 py-2 rounded-xl bg-brand-600 hover:bg-brand-500 text-white text-xs font-semibold transition shadow-md shadow-brand-600/20">
                            Rename
                        </button>
                    </div>
                </div>
            </div>

            <!-- 4. Copy / Move Modal -->
            <div v-if="isCopyMoveOpen" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/60 backdrop-blur-sm">
                <div class="w-full max-w-lg bg-white dark:bg-surface-900 rounded-2xl border border-slate-200 dark:border-surface-800 shadow-2xl p-5 space-y-4">
                    <div class="flex items-center justify-between">
                        <h3 class="text-sm font-bold text-slate-900 dark:text-white capitalize flex items-center gap-2">
                            <component :is="copyMoveMode === 'copy' ? Copy : Scissors" class="w-4 h-4 text-brand-600" />
                            <span>{{ copyMoveMode }} {{ copyMoveItems.length }} Item(s)</span>
                        </h3>
                        <button @click="isCopyMoveOpen = false" class="p-1 rounded-lg text-slate-400 hover:text-slate-600 dark:hover:text-white">
                            <X class="w-4 h-4" />
                        </button>
                    </div>
                    <div>
                        <label class="block text-xs font-semibold text-slate-700 dark:text-surface-300 mb-1.5">Destination Directory (Relative Path)</label>
                        <input
                            v-model="copyMoveDestPath"
                            type="text"
                            placeholder="e.g. public/assets or leave blank for root"
                            class="w-full bg-slate-50 dark:bg-surface-950 border border-slate-200 dark:border-surface-800 rounded-xl px-3 py-2 text-xs text-slate-900 dark:text-white font-mono focus:outline-none focus:ring-2 focus:ring-brand-500"
                        />
                    </div>
                    <div class="flex justify-end gap-2">
                        <button @click="isCopyMoveOpen = false" class="px-4 py-2 rounded-xl bg-slate-100 dark:bg-surface-800 text-xs font-semibold text-slate-700 dark:text-surface-300 hover:bg-slate-200 transition">
                            Cancel
                        </button>
                        <button @click="handleCopyMove" class="px-4 py-2 rounded-xl bg-brand-600 hover:bg-brand-500 text-white text-xs font-semibold transition shadow-md shadow-brand-600/20 capitalize">
                            {{ copyMoveMode }} Items
                        </button>
                    </div>
                </div>
            </div>

            <!-- 5. Delete Confirmation Modal -->
            <div v-if="isDeleteOpen" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/60 backdrop-blur-sm">
                <div class="w-full max-w-md bg-white dark:bg-surface-900 rounded-2xl border border-slate-200 dark:border-surface-800 shadow-2xl p-5 space-y-4">
                    <div class="flex items-center gap-3 text-rose-600 dark:text-rose-400">
                        <div class="p-2 rounded-xl bg-rose-500/10 border border-rose-500/20">
                            <AlertTriangle class="w-6 h-6" />
                        </div>
                        <div>
                            <h3 class="text-sm font-bold text-slate-900 dark:text-white">Delete {{ deleteTargets.length }} Item(s)?</h3>
                            <p class="text-xs text-slate-500 dark:text-surface-400">This action cannot be undone.</p>
                        </div>
                    </div>
                    <div class="max-h-36 overflow-y-auto bg-slate-50 dark:bg-surface-950 p-2.5 rounded-xl border border-slate-200 dark:border-surface-800 font-mono text-xs space-y-1">
                        <div v-for="t in deleteTargets" :key="t" class="text-slate-800 dark:text-surface-200 truncate">
                            • {{ t }}
                        </div>
                    </div>
                    <div class="flex justify-end gap-2">
                        <button @click="isDeleteOpen = false" class="px-4 py-2 rounded-xl bg-slate-100 dark:bg-surface-800 text-xs font-semibold text-slate-700 dark:text-surface-300 hover:bg-slate-200 transition">
                            Cancel
                        </button>
                        <button @click="handleDelete" class="px-4 py-2 rounded-xl bg-rose-600 hover:bg-rose-500 text-white text-xs font-semibold transition shadow-md shadow-rose-600/20">
                            Delete Permanently
                        </button>
                    </div>
                </div>
            </div>

            <!-- 6. Compress Modal -->
            <div v-if="isCompressOpen" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/60 backdrop-blur-sm">
                <div class="w-full max-w-md bg-white dark:bg-surface-900 rounded-2xl border border-slate-200 dark:border-surface-800 shadow-2xl p-5 space-y-4">
                    <div class="flex items-center justify-between">
                        <h3 class="text-sm font-bold text-slate-900 dark:text-white flex items-center gap-2">
                            <Archive class="w-4 h-4 text-rose-500" />
                            <span>Compress {{ compressSources.length }} Item(s)</span>
                        </h3>
                        <button @click="isCompressOpen = false" class="p-1 rounded-lg text-slate-400 hover:text-slate-600 dark:hover:text-white">
                            <X class="w-4 h-4" />
                        </button>
                    </div>
                    <div>
                        <label class="block text-xs font-semibold text-slate-700 dark:text-surface-300 mb-1.5">Archive Name</label>
                        <input
                            v-model="compressArchiveName"
                            type="text"
                            class="w-full bg-slate-50 dark:bg-surface-950 border border-slate-200 dark:border-surface-800 rounded-xl px-3 py-2 text-xs text-slate-900 dark:text-white font-mono focus:outline-none focus:ring-2 focus:ring-brand-500"
                        />
                    </div>
                    <div>
                        <label class="block text-xs font-semibold text-slate-700 dark:text-surface-300 mb-1.5">Compression Format</label>
                        <div class="grid grid-cols-2 gap-2">
                            <button
                                type="button"
                                @click="compressFormat = 'zip'"
                                :class="[
                                    'py-2 px-3 rounded-xl border text-xs font-bold text-center transition',
                                    compressFormat === 'zip'
                                        ? 'bg-brand-50 border-brand-500 text-brand-600 dark:bg-brand-500/10 dark:text-brand-400'
                                        : 'border-slate-200 dark:border-surface-800 text-slate-600 dark:text-surface-300'
                                ]"
                            >
                                ZIP (.zip)
                            </button>
                            <button
                                type="button"
                                @click="compressFormat = 'tar.gz'"
                                :class="[
                                    'py-2 px-3 rounded-xl border text-xs font-bold text-center transition',
                                    compressFormat === 'tar.gz'
                                        ? 'bg-brand-50 border-brand-500 text-brand-600 dark:bg-brand-500/10 dark:text-brand-400'
                                        : 'border-slate-200 dark:border-surface-800 text-slate-600 dark:text-surface-300'
                                ]"
                            >
                                TAR.GZ (.tar.gz)
                            </button>
                        </div>
                    </div>
                    <div class="flex justify-end gap-2">
                        <button @click="isCompressOpen = false" class="px-4 py-2 rounded-xl bg-slate-100 dark:bg-surface-800 text-xs font-semibold text-slate-700 dark:text-surface-300 hover:bg-slate-200 transition">
                            Cancel
                        </button>
                        <button @click="handleCompress" class="px-4 py-2 rounded-xl bg-brand-600 hover:bg-brand-500 text-white text-xs font-semibold transition shadow-md shadow-brand-600/20">
                            Compress
                        </button>
                    </div>
                </div>
            </div>

            <!-- 7. Extract Modal -->
            <div v-if="isExtractOpen" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/60 backdrop-blur-sm">
                <div class="w-full max-w-md bg-white dark:bg-surface-900 rounded-2xl border border-slate-200 dark:border-surface-800 shadow-2xl p-5 space-y-4">
                    <div class="flex items-center justify-between">
                        <h3 class="text-sm font-bold text-slate-900 dark:text-white flex items-center gap-2">
                            <Layers class="w-4 h-4 text-indigo-500" />
                            <span>Extract Archive</span>
                        </h3>
                        <button @click="isExtractOpen = false" class="p-1 rounded-lg text-slate-400 hover:text-slate-600 dark:hover:text-white">
                            <X class="w-4 h-4" />
                        </button>
                    </div>
                    <div>
                        <label class="block text-xs font-semibold text-slate-700 dark:text-surface-300 mb-1.5">Destination Directory</label>
                        <input
                            v-model="extractDestPath"
                            type="text"
                            placeholder="Leave blank for current directory"
                            class="w-full bg-slate-50 dark:bg-surface-950 border border-slate-200 dark:border-surface-800 rounded-xl px-3 py-2 text-xs text-slate-900 dark:text-white font-mono focus:outline-none focus:ring-2 focus:ring-brand-500"
                        />
                    </div>
                    <div class="flex justify-end gap-2">
                        <button @click="isExtractOpen = false" class="px-4 py-2 rounded-xl bg-slate-100 dark:bg-surface-800 text-xs font-semibold text-slate-700 dark:text-surface-300 hover:bg-slate-200 transition">
                            Cancel
                        </button>
                        <button @click="handleExtract" class="px-4 py-2 rounded-xl bg-brand-600 hover:bg-brand-500 text-white text-xs font-semibold transition shadow-md shadow-brand-600/20">
                            Extract
                        </button>
                    </div>
                </div>
            </div>

            <!-- 8. Permissions (chmod) Matrix Modal -->
            <div v-if="isPermissionsOpen" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/60 backdrop-blur-sm">
                <div class="w-full max-w-md bg-white dark:bg-surface-900 rounded-2xl border border-slate-200 dark:border-surface-800 shadow-2xl p-5 space-y-4">
                    <div class="flex items-center justify-between">
                        <h3 class="text-sm font-bold text-slate-900 dark:text-white flex items-center gap-2">
                            <Shield class="w-4 h-4 text-brand-600" />
                            <span>Edit Permissions</span>
                        </h3>
                        <button @click="isPermissionsOpen = false" class="p-1 rounded-lg text-slate-400 hover:text-slate-600 dark:hover:text-white">
                            <X class="w-4 h-4" />
                        </button>
                    </div>

                    <!-- 3x3 Permission Table -->
                    <div class="border border-slate-200 dark:border-surface-800 rounded-xl overflow-hidden text-xs">
                        <table class="w-full text-center">
                            <thead class="bg-slate-50 dark:bg-surface-950 text-slate-500 dark:text-surface-400 font-semibold text-[11px]">
                                <tr>
                                    <th class="py-2 px-3 text-left">Role</th>
                                    <th class="py-2 px-3">Read</th>
                                    <th class="py-2 px-3">Write</th>
                                    <th class="py-2 px-3">Execute</th>
                                </tr>
                            </thead>
                            <tbody class="divide-y divide-slate-100 dark:divide-surface-800">
                                <tr>
                                    <td class="py-2 px-3 text-left font-semibold text-slate-900 dark:text-white">Owner</td>
                                    <td class="py-2 px-3"><input type="checkbox" v-model="permOwnerR" class="rounded text-brand-600 focus:ring-brand-500" /></td>
                                    <td class="py-2 px-3"><input type="checkbox" v-model="permOwnerW" class="rounded text-brand-600 focus:ring-brand-500" /></td>
                                    <td class="py-2 px-3"><input type="checkbox" v-model="permOwnerX" class="rounded text-brand-600 focus:ring-brand-500" /></td>
                                </tr>
                                <tr>
                                    <td class="py-2 px-3 text-left font-semibold text-slate-900 dark:text-white">Group</td>
                                    <td class="py-2 px-3"><input type="checkbox" v-model="permGroupR" class="rounded text-brand-600 focus:ring-brand-500" /></td>
                                    <td class="py-2 px-3"><input type="checkbox" v-model="permGroupW" class="rounded text-brand-600 focus:ring-brand-500" /></td>
                                    <td class="py-2 px-3"><input type="checkbox" v-model="permGroupX" class="rounded text-brand-600 focus:ring-brand-500" /></td>
                                </tr>
                                <tr>
                                    <td class="py-2 px-3 text-left font-semibold text-slate-900 dark:text-white">Others</td>
                                    <td class="py-2 px-3"><input type="checkbox" v-model="permOtherR" class="rounded text-brand-600 focus:ring-brand-500" /></td>
                                    <td class="py-2 px-3"><input type="checkbox" v-model="permOtherW" class="rounded text-brand-600 focus:ring-brand-500" /></td>
                                    <td class="py-2 px-3"><input type="checkbox" v-model="permOtherX" class="rounded text-brand-600 focus:ring-brand-500" /></td>
                                </tr>
                            </tbody>
                        </table>
                    </div>

                    <div class="flex items-center justify-between">
                        <div class="flex items-center gap-2">
                            <span class="text-xs font-semibold text-slate-700 dark:text-surface-300">Octal:</span>
                            <span class="px-2.5 py-1 rounded-lg bg-slate-100 dark:bg-surface-950 font-mono font-bold text-sm text-slate-900 dark:text-white border border-slate-200 dark:border-surface-800">
                                {{ calculatedOctal }}
                            </span>
                        </div>
                        <label v-if="permissionsTarget?.is_dir" class="flex items-center gap-2 text-xs text-slate-600 dark:text-surface-300 cursor-pointer">
                            <input type="checkbox" v-model="permRecursive" class="rounded text-brand-600 focus:ring-brand-500" />
                            <span>Apply recursively</span>
                        </label>
                    </div>

                    <div class="flex justify-end gap-2">
                        <button @click="isPermissionsOpen = false" class="px-4 py-2 rounded-xl bg-slate-100 dark:bg-surface-800 text-xs font-semibold text-slate-700 dark:text-surface-300 hover:bg-slate-200 transition">
                            Cancel
                        </button>
                        <button @click="handleSavePermissions" class="px-4 py-2 rounded-xl bg-brand-600 hover:bg-brand-500 text-white text-xs font-semibold transition shadow-md shadow-brand-600/20">
                            Save Permissions
                        </button>
                    </div>
                </div>
            </div>

            <!-- 9. File Details / Properties Modal -->
            <div v-if="isDetailsOpen && detailsData" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/60 backdrop-blur-sm">
                <div class="w-full max-w-lg bg-white dark:bg-surface-900 rounded-2xl border border-slate-200 dark:border-surface-800 shadow-2xl p-5 space-y-4">
                    <div class="flex items-center justify-between">
                        <h3 class="text-sm font-bold text-slate-900 dark:text-white flex items-center gap-2">
                            <Info class="w-4 h-4 text-brand-600" />
                            <span>File Details</span>
                        </h3>
                        <button @click="isDetailsOpen = false" class="p-1 rounded-lg text-slate-400 hover:text-slate-600 dark:hover:text-white">
                            <X class="w-4 h-4" />
                        </button>
                    </div>

                    <div class="space-y-2 text-xs font-mono">
                        <div class="flex justify-between py-1 border-b border-slate-100 dark:border-surface-800">
                            <span class="text-slate-500 dark:text-surface-400">Name:</span>
                            <span class="font-bold text-slate-900 dark:text-white">{{ detailsData.name }}</span>
                        </div>
                        <div class="flex justify-between py-1 border-b border-slate-100 dark:border-surface-800">
                            <span class="text-slate-500 dark:text-surface-400">Path:</span>
                            <span class="text-slate-800 dark:text-surface-200">{{ detailsData.path || '/' }}</span>
                        </div>
                        <div class="flex justify-between py-1 border-b border-slate-100 dark:border-surface-800">
                            <span class="text-slate-500 dark:text-surface-400">Type:</span>
                            <span class="text-slate-800 dark:text-surface-200">{{ detailsData.is_dir ? 'Directory' : (detailsData.mime_type || 'File') }}</span>
                        </div>
                        <div class="flex justify-between py-1 border-b border-slate-100 dark:border-surface-800">
                            <span class="text-slate-500 dark:text-surface-400">Size:</span>
                            <span class="text-slate-800 dark:text-surface-200">{{ detailsData.is_dir ? `${detailsData.item_count} items` : formatBytes(detailsData.size_bytes) }}</span>
                        </div>
                        <div class="flex justify-between py-1 border-b border-slate-100 dark:border-surface-800">
                            <span class="text-slate-500 dark:text-surface-400">Permissions:</span>
                            <span class="text-slate-800 dark:text-surface-200">{{ detailsData.permissions }} ({{ detailsData.mode_octal }})</span>
                        </div>
                        <div class="flex justify-between py-1 border-b border-slate-100 dark:border-surface-800">
                            <span class="text-slate-500 dark:text-surface-400">Owner / Group:</span>
                            <span class="text-slate-800 dark:text-surface-200">{{ detailsData.owner }} ({{ detailsData.uid }}) : {{ detailsData.group }} ({{ detailsData.gid }})</span>
                        </div>
                        <div class="flex justify-between py-1 border-b border-slate-100 dark:border-surface-800">
                            <span class="text-slate-500 dark:text-surface-400">Last Modified:</span>
                            <span class="text-slate-800 dark:text-surface-200">{{ new Date(detailsData.modified_at).toLocaleString() }}</span>
                        </div>
                    </div>

                    <div class="flex justify-end">
                        <button @click="isDetailsOpen = false" class="px-4 py-2 rounded-xl bg-slate-100 dark:bg-surface-800 text-xs font-semibold text-slate-700 dark:text-surface-300 hover:bg-slate-200 transition">
                            Close
                        </button>
                    </div>
                </div>
            </div>

            <!-- 10. File Preview Modal -->
            <div v-if="isPreviewOpen && previewItem" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/70 backdrop-blur-sm">
                <div class="w-full max-w-4xl max-h-[85vh] flex flex-col bg-white dark:bg-surface-900 rounded-2xl border border-slate-200 dark:border-surface-800 shadow-2xl overflow-hidden">
                    <div class="h-12 bg-slate-100 dark:bg-surface-950 border-b border-slate-200 dark:border-surface-800 px-4 flex items-center justify-between shrink-0">
                        <div class="flex items-center gap-2 font-mono text-xs text-slate-900 dark:text-white font-semibold">
                            <Eye class="w-4 h-4 text-brand-600" />
                            <span>{{ previewItem.name }}</span>
                        </div>
                        <div class="flex items-center gap-2">
                            <button @click="handleDownload(previewItem)" class="px-3 py-1 rounded-lg bg-slate-200 dark:bg-surface-800 hover:bg-slate-300 text-xs font-medium flex items-center gap-1">
                                <Download class="w-3.5 h-3.5" />
                                <span>Download</span>
                            </button>
                            <button @click="isPreviewOpen = false" class="p-1 rounded-lg text-slate-400 hover:text-slate-600 dark:hover:text-white">
                                <X class="w-4 h-4" />
                            </button>
                        </div>
                    </div>

                    <div class="flex-1 overflow-auto p-4 flex items-center justify-center bg-slate-950">
                        <!-- Image Preview -->
                        <img
                            v-if="isImageFile(previewItem.name)"
                            :src="previewUrl"
                            :alt="previewItem.name"
                            class="max-w-full max-h-[70vh] object-contain rounded-lg shadow-lg"
                        />
                        <!-- PDF Preview -->
                        <iframe
                            v-else-if="isPdfFile(previewItem.name)"
                            :src="previewUrl"
                            class="w-full h-[70vh] rounded-lg border border-surface-800"
                        ></iframe>
                        <!-- Text Preview -->
                        <pre
                            v-else
                            class="w-full h-full font-mono text-xs text-emerald-400 p-4 whitespace-pre-wrap overflow-auto"
                        >{{ previewTextContent }}</pre>
                    </div>
                </div>
            </div>

            <!-- 11. Full-Featured Code Editor -->
            <div
                v-if="isEditorOpen && editingFile"
                class="fixed inset-0 z-50 flex items-center justify-center p-2 sm:p-4 bg-slate-900/70 backdrop-blur-sm"
            >
                <div
                    :class="[
                        'flex flex-col bg-white dark:bg-surface-900 border border-slate-200 dark:border-surface-800 shadow-2xl overflow-hidden transition-all',
                        isEditorFullscreen ? 'fixed inset-0 rounded-none' : 'w-full max-w-5xl h-[85vh] rounded-2xl'
                    ]"
                >
                    <!-- Editor Header Toolbar -->
                    <div class="h-12 bg-slate-100 dark:bg-surface-950 border-b border-slate-200 dark:border-surface-800 px-4 flex items-center justify-between shrink-0">
                        <div class="flex items-center gap-3">
                            <div class="flex items-center gap-2 font-mono text-xs text-slate-900 dark:text-white font-semibold">
                                <Code class="w-4 h-4 text-brand-600" />
                                <span>{{ editingFile.name }}</span>
                            </div>
                            <!-- Modified Badge Indicator -->
                            <span
                                :class="[
                                    'px-2 py-0.5 rounded-full text-[10px] font-mono font-medium transition',
                                    isEditorModified
                                        ? 'bg-amber-50 text-amber-600 border border-amber-200 dark:bg-amber-500/10 dark:text-amber-400 dark:border-amber-500/20'
                                        : 'bg-emerald-50 text-emerald-600 border border-emerald-200 dark:bg-emerald-500/10 dark:text-emerald-400 dark:border-emerald-500/20'
                                ]"
                            >
                                {{ isEditorModified ? '● Modified' : '✓ Saved' }}
                            </span>
                        </div>

                        <!-- Editor Controls -->
                        <div class="flex items-center gap-2">
                            <!-- Search & Replace Toggle -->
                            <button
                                @click="isEditorSearchOpen = !isEditorSearchOpen"
                                :class="[
                                    'p-1.5 rounded-lg border text-xs transition',
                                    isEditorSearchOpen
                                        ? 'bg-brand-50 border-brand-400 text-brand-600 dark:bg-brand-500/10 dark:text-brand-400'
                                        : 'border-slate-200 dark:border-surface-800 text-slate-500 hover:bg-slate-200 dark:hover:bg-surface-800'
                                ]"
                                title="Search & Replace (Ctrl+F)"
                            >
                                <Search class="w-3.5 h-3.5" />
                            </button>

                            <!-- Word Wrap Toggle -->
                            <button
                                @click="editorWordWrap = !editorWordWrap"
                                :class="[
                                    'px-2.5 py-1 rounded-lg border text-[11px] font-medium transition',
                                    editorWordWrap
                                        ? 'bg-brand-50 border-brand-400 text-brand-600 dark:bg-brand-500/10 dark:text-brand-400'
                                        : 'border-slate-200 dark:border-surface-800 text-slate-500 hover:bg-slate-200 dark:hover:bg-surface-800'
                                ]"
                            >
                                Wrap
                            </button>

                            <!-- Fullscreen Toggle -->
                            <button
                                @click="isEditorFullscreen = !isEditorFullscreen"
                                class="p-1.5 rounded-lg border border-slate-200 dark:border-surface-800 text-slate-500 hover:bg-slate-200 dark:hover:bg-surface-800 transition"
                                title="Toggle Fullscreen"
                            >
                                <component :is="isEditorFullscreen ? Minimize2 : Maximize2" class="w-3.5 h-3.5" />
                            </button>

                            <!-- Save Button -->
                            <button
                                @click="saveEditorContent"
                                :disabled="isEditorSaving || !isEditorModified"
                                class="px-3.5 py-1.5 rounded-xl bg-emerald-600 hover:bg-emerald-500 text-white text-xs font-semibold flex items-center gap-1.5 transition shadow-sm disabled:opacity-40"
                            >
                                <Save class="w-3.5 h-3.5" />
                                <span>{{ isEditorSaving ? 'Saving...' : 'Save' }}</span>
                            </button>

                            <!-- Close Button -->
                            <button
                                @click="closeEditor"
                                class="p-1.5 rounded-xl text-slate-400 hover:text-slate-600 dark:hover:text-white transition"
                            >
                                <X class="w-4 h-4" />
                            </button>
                        </div>
                    </div>

                    <!-- Search & Replace Bar -->
                    <div
                        v-if="isEditorSearchOpen"
                        class="px-4 py-2 bg-slate-200/70 dark:bg-surface-950 border-b border-slate-200 dark:border-surface-800 flex flex-wrap items-center gap-2 text-xs"
                    >
                        <input
                            v-model="editorFindText"
                            type="text"
                            placeholder="Find..."
                            class="bg-white dark:bg-surface-900 border border-slate-300 dark:border-surface-700 rounded-lg px-2.5 py-1 text-xs text-slate-900 dark:text-white font-mono focus:outline-none focus:ring-1 focus:ring-brand-500"
                        />
                        <input
                            v-model="editorReplaceText"
                            type="text"
                            placeholder="Replace with..."
                            class="bg-white dark:bg-surface-900 border border-slate-300 dark:border-surface-700 rounded-lg px-2.5 py-1 text-xs text-slate-900 dark:text-white font-mono focus:outline-none focus:ring-1 focus:ring-brand-500"
                        />
                        <button
                            @click="handleEditorReplaceAll"
                            class="px-3 py-1 bg-brand-600 text-white rounded-lg font-semibold hover:bg-brand-500 transition"
                        >
                            Replace All
                        </button>
                    </div>

                    <!-- Code Editor Text Area with Line Numbers -->
                    <div class="flex-1 flex bg-[#0d1117] text-slate-100 overflow-hidden font-mono text-xs">
                        <!-- Line Numbers Column -->
                        <div class="w-12 bg-[#090d13] text-slate-600 dark:text-surface-600 select-none py-3 px-2 text-right border-r border-[#1e2633] overflow-hidden leading-6 shrink-0">
                            <div v-for="n in editorLineCount" :key="n">{{ n }}</div>
                        </div>

                        <!-- Editor Textarea -->
                        <div class="flex-1 p-3 overflow-auto">
                            <textarea
                                v-model="editorContent"
                                spellcheck="false"
                                :class="[
                                    'w-full h-full bg-transparent border-0 resize-none font-mono text-xs text-emerald-400 focus:outline-none leading-6 selection:bg-brand-500 selection:text-white',
                                    editorWordWrap ? 'whitespace-pre-wrap' : 'whitespace-pre'
                                ]"
                            ></textarea>
                        </div>
                    </div>

                    <!-- Editor Statusbar Footer -->
                    <div class="h-7 bg-[#090d13] border-t border-[#1e2633] px-4 flex items-center justify-between text-[11px] font-mono text-slate-500">
                        <div class="flex items-center gap-3">
                            <span>Lines: {{ editorLineCount }}</span>
                            <span>Characters: {{ editorContent.length }}</span>
                        </div>
                        <div class="flex items-center gap-3">
                            <span>UTF-8</span>
                            <span>{{ editingFile.name.split('.').pop()?.toUpperCase() || 'TXT' }}</span>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </AppLayout>
</template>
