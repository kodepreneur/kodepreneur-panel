<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue';
import { router, Link } from '@inertiajs/vue3';
import axios from 'axios';
import AppLayout from '@/Layouts/AppLayout.vue';
import {
    Database,
    Table,
    Layers,
    Code,
    Search,
    RefreshCw,
    ChevronLeft,
    ChevronRight,
    ChevronsLeft,
    ChevronsRight,
    ArrowUpDown,
    ArrowUp,
    ArrowDown,
    Copy,
    Check,
    Key,
    Lock,
    HardDrive,
    SlidersHorizontal,
    Maximize2,
    X,
    AlertCircle,
    Info,
    ExternalLink,
    Hash,
    Calendar,
    CheckCircle2,
} from 'lucide-vue-next';
import type { DatabaseRecord, TableInfo, TableStructure, TableDataResult } from '@/types';

const props = defineProps<{
    currentDatabase: DatabaseRecord;
    allDatabases: DatabaseRecord[];
    initialTables: TableInfo[];
}>();

// State
const tables = ref<TableInfo[]>(props.initialTables || []);
const isTablesLoading = ref(false);
const tableSearchQuery = ref('');
const selectedTable = ref<string>('');
const activeTab = ref<'data' | 'structure' | 'ddl'>('data');

// Table Structure State
const structure = ref<TableStructure | null>(null);
const isStructureLoading = ref(false);

// Table Data State
const tableData = ref<TableDataResult | null>(null);
const isDataLoading = ref(false);
const currentPage = ref(1);
const perPage = ref(50);
const sortColumn = ref('');
const sortDirection = ref<'asc' | 'desc'>('asc');
const dataSearchQuery = ref('');
const dataSearchColumn = ref('');

// Cell View Modal State
const isCellModalOpen = ref(false);
const activeCellContent = ref<string>('');
const activeCellColumn = ref<string>('');
const isCopied = ref(false);
const isDdlCopied = ref(false);

// Filtered tables in sidebar
const filteredTables = computed(() => {
    if (!tableSearchQuery.value) return tables.value;
    const q = tableSearchQuery.value.toLowerCase();
    return tables.value.filter(t => t.name.toLowerCase().includes(q));
});

// Format byte sizes into readable string
function formatBytes(bytes: number | null | undefined): string {
    if (!bytes || bytes === 0) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
}

// Format numbers with commas
function formatNumber(num: number | null | undefined): string {
    if (num === null || num === undefined) return '0';
    return num.toLocaleString();
}

// Fetch tables list from server
async function fetchTables() {
    isTablesLoading.value = true;
    try {
        const response = await axios.get(`/databases/${props.currentDatabase.id}/tables`);
        if (response.data.success) {
            tables.value = response.data.tables || [];
            if (tables.value.length > 0 && (!selectedTable.value || !tables.value.some(t => t.name === selectedTable.value))) {
                selectTable(tables.value[0].name);
            }
        }
    } catch (err) {
        console.error('Failed to fetch tables:', err);
    } finally {
        isTablesLoading.value = false;
    }
}

// Select a table and load its structure & data
function selectTable(tableName: string) {
    selectedTable.value = tableName;
    currentPage.value = 1;
    sortColumn.value = '';
    sortDirection.value = 'asc';
    dataSearchQuery.value = '';
    dataSearchColumn.value = '';
    loadTableData();
    loadTableStructure();
}

// Load table structure
async function loadTableStructure() {
    if (!selectedTable.value) return;
    isStructureLoading.value = true;
    try {
        const response = await axios.get(`/databases/${props.currentDatabase.id}/tables/${selectedTable.value}/structure`);
        if (response.data.success) {
            structure.value = response.data.structure;
        }
    } catch (err) {
        console.error('Failed to load table structure:', err);
    } finally {
        isStructureLoading.value = false;
    }
}

// Load table paginated data
async function loadTableData() {
    if (!selectedTable.value) return;
    isDataLoading.value = true;
    try {
        const params: Record<string, any> = {
            page: currentPage.value,
            per_page: perPage.value,
        };
        if (sortColumn.value) {
            params.sort = sortColumn.value;
            params.direction = sortDirection.value;
        }
        if (dataSearchQuery.value && dataSearchColumn.value) {
            params.search = dataSearchQuery.value;
            params.search_column = dataSearchColumn.value;
        }

        const response = await axios.get(`/databases/${props.currentDatabase.id}/tables/${selectedTable.value}/data`, { params });
        if (response.data.success) {
            tableData.value = response.data.data;
        }
    } catch (err) {
        console.error('Failed to load table data:', err);
    } finally {
        isDataLoading.value = false;
    }
}

// Sort column clicked
function toggleSort(column: string) {
    if (sortColumn.value === column) {
        if (sortDirection.value === 'asc') {
            sortDirection.value = 'desc';
        } else {
            sortColumn.value = '';
            sortDirection.value = 'asc';
        }
    } else {
        sortColumn.value = column;
        sortDirection.value = 'asc';
    }
    currentPage.value = 1;
    loadTableData();
}

// Pagination navigation
function changePage(page: number) {
    if (!tableData.value) return;
    if (page < 1 || page > tableData.value.total_pages) return;
    currentPage.value = page;
    loadTableData();
}

// Change database switcher
function switchDatabase(e: Event) {
    const target = e.target as HTMLSelectElement;
    const newDbId = target.value;
    if (newDbId && String(newDbId) !== String(props.currentDatabase.id)) {
        router.visit(`/databases/${newDbId}/explorer`);
    }
}

// Cell details viewer
function openCellModal(column: string, val: any) {
    activeCellColumn.value = column;
    if (typeof val === 'object' && val !== null) {
        activeCellContent.value = JSON.stringify(val, null, 2);
    } else {
        activeCellContent.value = String(val ?? '');
    }
    isCellModalOpen.value = true;
}

// Copy cell content to clipboard
function copyToClipboard(text: string, type: 'cell' | 'ddl' = 'cell') {
    navigator.clipboard.writeText(text);
    if (type === 'ddl') {
        isDdlCopied.value = true;
        setTimeout(() => (isDdlCopied.value = false), 2000);
    } else {
        isCopied.value = true;
        setTimeout(() => (isCopied.value = false), 2000);
    }
}

// Check if value is large or JSON
function isExpandable(val: any): boolean {
    if (val === null || val === undefined) return false;
    if (typeof val === 'object') return true;
    const str = String(val);
    return str.length > 60 || str.includes('\n');
}

// Selected table object
const currentTableInfo = computed(() => {
    return tables.value.find(t => t.name === selectedTable.value);
});

// Mounted lifecycle
onMounted(() => {
    if (tables.value.length > 0) {
        selectTable(tables.value[0].name);
    } else {
        fetchTables();
    }
});

// Watch perPage change
watch(perPage, () => {
    currentPage.value = 1;
    loadTableData();
});
</script>

<template>
    <AppLayout :title="`Database Explorer - ${currentDatabase.name}`">
        <div class="space-y-4 max-w-[98%] mx-auto">
            <!-- Top Navigation & Database Header Card -->
            <div class="rounded-2xl bg-white dark:bg-surface-900/90 border border-slate-200/80 dark:border-surface-800/80 p-4 sm:p-5 shadow-sm dark:shadow-xl">
                <div class="flex flex-col lg:flex-row lg:items-center justify-between gap-4">
                    <!-- Left info: Breadcrumb & Title -->
                    <div class="flex items-center gap-3.5">
                        <Link
                            href="/databases"
                            class="p-2.5 rounded-xl bg-slate-100 hover:bg-slate-200 dark:bg-surface-800 dark:hover:bg-surface-700 text-slate-600 dark:text-surface-300 transition flex items-center justify-center shrink-0"
                            title="Back to Databases"
                        >
                            <ChevronLeft class="w-4 h-4" />
                        </Link>

                        <div>
                            <div class="flex items-center gap-2 flex-wrap">
                                <span class="text-xs text-slate-500 dark:text-surface-400 font-medium">Database Explorer</span>
                                <span class="text-slate-300 dark:text-surface-600">•</span>
                                <span
                                    :class="[
                                        'px-2 py-0.5 rounded-md font-mono text-[10px] uppercase font-bold tracking-wider',
                                        currentDatabase.engine === 'mysql'
                                            ? 'bg-blue-50 text-blue-700 border border-blue-200 dark:bg-blue-500/10 dark:border-blue-500/30 dark:text-blue-400'
                                            : 'bg-indigo-50 text-indigo-700 border border-indigo-200 dark:bg-indigo-500/10 dark:border-indigo-500/30 dark:text-indigo-400'
                                    ]"
                                >
                                    {{ currentDatabase.engine }}
                                </span>
                            </div>
                            <div class="flex items-center gap-2.5 mt-0.5">
                                <h1 class="text-lg font-bold text-slate-900 dark:text-white font-mono flex items-center gap-2">
                                    <Database class="w-5 h-5 text-emerald-500 shrink-0" />
                                    <span>{{ currentDatabase.name }}</span>
                                </h1>
                                <span class="text-xs text-slate-400 dark:text-surface-500 font-mono hidden sm:inline">
                                    ({{ currentDatabase.character_set }} / {{ currentDatabase.collation }})
                                </span>
                            </div>
                        </div>
                    </div>

                    <!-- Right Controls: Switch Database & Refresh -->
                    <div class="flex items-center gap-2.5 flex-wrap">
                        <!-- Switch Database Selector -->
                        <div class="flex items-center gap-2">
                            <label class="text-xs text-slate-500 dark:text-surface-400 font-medium hidden sm:inline">Database:</label>
                            <select
                                :value="currentDatabase.id"
                                @change="switchDatabase"
                                class="px-3 py-2 rounded-xl bg-slate-50 dark:bg-surface-950/80 border border-slate-200 dark:border-surface-800 text-xs font-mono text-slate-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-emerald-500/40"
                            >
                                <option
                                    v-for="db in allDatabases"
                                    :key="db.id"
                                    :value="db.id"
                                >
                                    {{ db.name }} ({{ db.engine.toUpperCase() }})
                                </option>
                            </select>
                        </div>

                        <!-- Refresh Tables Button -->
                        <button
                            @click="fetchTables"
                            :disabled="isTablesLoading"
                            class="px-3.5 py-2 rounded-xl bg-slate-100 hover:bg-slate-200 dark:bg-surface-800 dark:hover:bg-surface-700 text-slate-700 dark:text-surface-200 text-xs font-semibold flex items-center gap-1.5 transition border border-slate-200/80 dark:border-surface-700 disabled:opacity-50"
                            title="Refresh Tables"
                        >
                            <RefreshCw :class="['w-3.5 h-3.5', isTablesLoading ? 'animate-spin text-emerald-500' : 'text-slate-500 dark:text-surface-400']" />
                            <span class="hidden sm:inline">Refresh</span>
                        </button>
                    </div>
                </div>
            </div>

            <!-- Explorer Main Body (Sidebar + Content Workspace) -->
            <div class="grid grid-cols-1 lg:grid-cols-12 gap-4 items-start">
                <!-- Left Sidebar: Tables List (3 cols on lg) -->
                <div class="lg:col-span-3 rounded-2xl bg-white dark:bg-surface-900/80 border border-slate-200/80 dark:border-surface-800/80 shadow-sm overflow-hidden flex flex-col max-h-[calc(100vh-13rem)] sticky top-4">
                    <!-- Search Input & Table Count Header -->
                    <div class="p-3 border-b border-slate-200/80 dark:border-surface-800 space-y-2.5 bg-slate-50/50 dark:bg-surface-950/40">
                        <div class="flex items-center justify-between">
                            <span class="text-xs font-bold uppercase tracking-wider text-slate-600 dark:text-surface-300 flex items-center gap-1.5">
                                <Table class="w-3.5 h-3.5 text-emerald-500" />
                                <span>Tables ({{ tables.length }})</span>
                            </span>
                        </div>
                        <div class="relative">
                            <Search class="w-3.5 h-3.5 absolute left-3 top-1/2 -translate-y-1/2 text-slate-400 dark:text-surface-500" />
                            <input
                                v-model="tableSearchQuery"
                                type="text"
                                placeholder="Filter tables..."
                                class="w-full pl-8 pr-3 py-1.5 rounded-lg bg-white dark:bg-surface-900 border border-slate-200 dark:border-surface-800 text-xs text-slate-900 dark:text-white placeholder-slate-400 dark:placeholder-surface-500 font-mono focus:outline-none focus:ring-2 focus:ring-emerald-500/40"
                            />
                        </div>
                    </div>

                    <!-- Scrollable Tables List -->
                    <div class="overflow-y-auto flex-1 divide-y divide-slate-100 dark:divide-surface-800/50 p-1.5 space-y-0.5">
                        <div v-if="isTablesLoading" class="p-6 text-center text-xs text-slate-400 dark:text-surface-500">
                            <RefreshCw class="w-5 h-5 animate-spin mx-auto mb-2 text-emerald-500" />
                            <span>Loading tables...</span>
                        </div>

                        <div v-else-if="filteredTables.length === 0" class="p-6 text-center text-xs text-slate-400 dark:text-surface-500">
                            <Table class="w-8 h-8 mx-auto mb-2 text-slate-300 dark:text-surface-600" />
                            <p class="font-medium text-slate-600 dark:text-surface-400">No tables found</p>
                            <p class="text-[11px] text-slate-400 dark:text-surface-500 mt-0.5">
                                {{ tableSearchQuery ? 'No match for search filter' : 'This database has no tables yet' }}
                            </p>
                        </div>

                        <button
                            v-for="t in filteredTables"
                            :key="t.name"
                            @click="selectTable(t.name)"
                            :class="[
                                'w-full text-left px-3 py-2.5 rounded-xl text-xs transition flex items-center justify-between gap-2 group',
                                selectedTable === t.name
                                    ? 'bg-emerald-50 text-emerald-900 dark:bg-emerald-500/15 dark:text-emerald-300 font-bold border border-emerald-200/80 dark:border-emerald-500/30 shadow-xs'
                                    : 'text-slate-700 dark:text-surface-300 hover:bg-slate-50 dark:hover:bg-surface-800/60 font-medium'
                            ]"
                        >
                            <div class="flex items-center gap-2 truncate">
                                <Table :class="['w-3.5 h-3.5 shrink-0', selectedTable === t.name ? 'text-emerald-600 dark:text-emerald-400' : 'text-slate-400 dark:text-surface-500 group-hover:text-slate-600 dark:group-hover:text-surface-300']" />
                                <span class="truncate font-mono">{{ t.name }}</span>
                            </div>
                            <div class="flex items-center gap-1.5 shrink-0 font-mono text-[10px]">
                                <span :class="['px-1.5 py-0.5 rounded', selectedTable === t.name ? 'bg-emerald-100/80 text-emerald-800 dark:bg-emerald-500/20 dark:text-emerald-300' : 'bg-slate-100 text-slate-500 dark:bg-surface-800 dark:text-surface-400']">
                                    {{ formatNumber(t.rows) }}
                                </span>
                            </div>
                        </button>
                    </div>
                </div>

                <!-- Right Area: Table Workspace & Tabs (9 cols on lg) -->
                <div class="lg:col-span-9 space-y-4">
                    <!-- Workspace Card -->
                    <div v-if="selectedTable" class="rounded-2xl bg-white dark:bg-surface-900/80 border border-slate-200/80 dark:border-surface-800/80 shadow-sm dark:shadow-xl overflow-hidden">
                        <!-- Table Metadata Header Bar -->
                        <div class="p-4 sm:p-5 border-b border-slate-200/80 dark:border-surface-800 bg-slate-50/60 dark:bg-surface-950/40">
                            <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3">
                                <div>
                                    <div class="flex items-center gap-2.5">
                                        <h2 class="text-base font-bold text-slate-900 dark:text-white font-mono flex items-center gap-2">
                                            <Table class="w-4 h-4 text-emerald-500" />
                                            <span>{{ selectedTable }}</span>
                                        </h2>
                                        <span v-if="currentTableInfo?.type" class="px-2 py-0.5 rounded text-[10px] font-mono uppercase bg-slate-100 text-slate-600 dark:bg-surface-800 dark:text-surface-400 border border-slate-200 dark:border-surface-700">
                                            {{ currentTableInfo.type }}
                                        </span>
                                    </div>
                                    <div class="flex items-center gap-3 text-xs text-slate-500 dark:text-surface-400 mt-1 font-mono flex-wrap">
                                        <span>Estimated Rows: <strong class="text-slate-700 dark:text-surface-200">{{ formatNumber(currentTableInfo?.rows) }}</strong></span>
                                        <span class="text-slate-300 dark:text-surface-600">•</span>
                                        <span>Total Size: <strong class="text-slate-700 dark:text-surface-200">{{ formatBytes(currentTableInfo?.total_size) }}</strong></span>
                                        <span v-if="currentTableInfo?.data_length" class="text-slate-300 dark:text-surface-600">•</span>
                                        <span v-if="currentTableInfo?.data_length">Data: <strong class="text-slate-700 dark:text-surface-200">{{ formatBytes(currentTableInfo?.data_length) }}</strong></span>
                                        <span v-if="currentTableInfo?.index_length" class="text-slate-300 dark:text-surface-600">•</span>
                                        <span v-if="currentTableInfo?.index_length">Index: <strong class="text-slate-700 dark:text-surface-200">{{ formatBytes(currentTableInfo?.index_length) }}</strong></span>
                                    </div>
                                </div>

                                <!-- Tab Buttons -->
                                <div class="flex items-center gap-1.5 p-1 rounded-xl bg-slate-200/70 dark:bg-surface-800/80 self-start sm:self-auto text-xs">
                                    <button
                                        @click="activeTab = 'data'"
                                        :class="[
                                            'px-3 py-1.5 rounded-lg font-semibold flex items-center gap-1.5 transition',
                                            activeTab === 'data'
                                                ? 'bg-white dark:bg-surface-900 text-emerald-600 dark:text-emerald-400 shadow-xs'
                                                : 'text-slate-600 dark:text-surface-400 hover:text-slate-900 dark:hover:text-white'
                                        ]"
                                    >
                                        <Table class="w-3.5 h-3.5" />
                                        <span>Browse Data</span>
                                    </button>

                                    <button
                                        @click="activeTab = 'structure'"
                                        :class="[
                                            'px-3 py-1.5 rounded-lg font-semibold flex items-center gap-1.5 transition',
                                            activeTab === 'structure'
                                                ? 'bg-white dark:bg-surface-900 text-emerald-600 dark:text-emerald-400 shadow-xs'
                                                : 'text-slate-600 dark:text-surface-400 hover:text-slate-900 dark:hover:text-white'
                                        ]"
                                    >
                                        <Layers class="w-3.5 h-3.5" />
                                        <span>Structure</span>
                                    </button>

                                    <button
                                        @click="activeTab = 'ddl'"
                                        :class="[
                                            'px-3 py-1.5 rounded-lg font-semibold flex items-center gap-1.5 transition',
                                            activeTab === 'ddl'
                                                ? 'bg-white dark:bg-surface-900 text-emerald-600 dark:text-emerald-400 shadow-xs'
                                                : 'text-slate-600 dark:text-surface-400 hover:text-slate-900 dark:hover:text-white'
                                        ]"
                                    >
                                        <Code class="w-3.5 h-3.5" />
                                        <span>SQL DDL</span>
                                    </button>
                                </div>
                            </div>
                        </div>

                        <!-- TAB 1: DATA BROWSER -->
                        <div v-if="activeTab === 'data'" class="space-y-0">
                            <!-- Data Query Controls Bar -->
                            <div class="p-3 sm:px-4 border-b border-slate-200/80 dark:border-surface-800 flex flex-col sm:flex-row items-stretch sm:items-center justify-between gap-3 bg-white dark:bg-surface-900">
                                <!-- Search Filter -->
                                <div class="flex items-center gap-2 flex-1 max-w-md">
                                    <select
                                        v-model="dataSearchColumn"
                                        class="px-2.5 py-1.5 rounded-lg bg-slate-50 dark:bg-surface-950/80 border border-slate-200 dark:border-surface-800 text-xs font-mono text-slate-700 dark:text-surface-300 focus:outline-none focus:ring-2 focus:ring-emerald-500/40"
                                    >
                                        <option value="">Search column...</option>
                                        <option v-for="col in tableData?.columns || []" :key="col" :value="col">
                                            {{ col }}
                                        </option>
                                    </select>
                                    <div class="relative flex-1">
                                        <Search class="w-3.5 h-3.5 absolute left-2.5 top-1/2 -translate-y-1/2 text-slate-400 dark:text-surface-500" />
                                        <input
                                            v-model="dataSearchQuery"
                                            @keydown.enter="loadTableData"
                                            type="text"
                                            placeholder="Search value..."
                                            class="w-full pl-8 pr-3 py-1.5 rounded-lg bg-slate-50 dark:bg-surface-950/80 border border-slate-200 dark:border-surface-800 text-xs text-slate-900 dark:text-white font-mono focus:outline-none focus:ring-2 focus:ring-emerald-500/40"
                                        />
                                    </div>
                                    <button
                                        @click="loadTableData"
                                        class="px-3 py-1.5 rounded-lg bg-emerald-600 hover:bg-emerald-500 text-white text-xs font-semibold transition"
                                    >
                                        Filter
                                    </button>
                                </div>

                                <!-- Rows Per Page & Reload -->
                                <div class="flex items-center gap-2 justify-end">
                                    <div class="flex items-center gap-1.5 text-xs text-slate-500 dark:text-surface-400">
                                        <span>Show:</span>
                                        <select
                                            v-model="perPage"
                                            class="px-2 py-1 rounded-lg bg-slate-50 dark:bg-surface-950/80 border border-slate-200 dark:border-surface-800 text-xs font-mono text-slate-700 dark:text-surface-300 focus:outline-none"
                                        >
                                            <option :value="25">25</option>
                                            <option :value="50">50</option>
                                            <option :value="100">100</option>
                                            <option :value="200">200</option>
                                        </select>
                                    </div>

                                    <button
                                        @click="loadTableData"
                                        :disabled="isDataLoading"
                                        class="p-1.5 rounded-lg text-slate-500 hover:bg-slate-100 dark:text-surface-400 dark:hover:bg-surface-800 transition"
                                        title="Reload Data"
                                    >
                                        <RefreshCw :class="['w-4 h-4', isDataLoading ? 'animate-spin text-emerald-500' : '']" />
                                    </button>
                                </div>
                            </div>

                            <!-- Data Table Content -->
                            <div class="overflow-x-auto max-h-[600px] relative">
                                <div v-if="isDataLoading" class="absolute inset-0 z-20 bg-white/70 dark:bg-surface-950/70 backdrop-blur-[1px] flex items-center justify-center">
                                    <div class="flex items-center gap-2 px-4 py-2 rounded-xl bg-white dark:bg-surface-900 border border-slate-200 dark:border-surface-800 shadow-lg text-xs font-medium text-slate-700 dark:text-surface-200">
                                        <RefreshCw class="w-4 h-4 animate-spin text-emerald-500" />
                                        <span>Fetching records...</span>
                                    </div>
                                </div>

                                <div v-if="!tableData || tableData.rows.length === 0" class="py-16 px-4 text-center">
                                    <Table class="w-10 h-10 text-slate-300 dark:text-surface-600 mx-auto mb-3" />
                                    <h3 class="text-sm font-medium text-slate-700 dark:text-surface-200">No records found</h3>
                                    <p class="text-xs text-slate-500 dark:text-surface-400 mt-1 max-w-sm mx-auto">
                                        This table is currently empty or no rows match your search criteria.
                                    </p>
                                </div>

                                <table v-else class="w-full text-left text-xs border-collapse">
                                    <thead class="bg-slate-50 dark:bg-surface-950 text-slate-600 dark:text-surface-300 text-[11px] font-mono sticky top-0 z-10 border-b border-slate-200 dark:border-surface-800 shadow-xs">
                                        <tr>
                                            <th class="py-2.5 px-3 font-semibold text-slate-400 dark:text-surface-500 w-12 text-center border-r border-slate-200/60 dark:border-surface-800/60">
                                                #
                                            </th>
                                            <th
                                                v-for="col in tableData.columns"
                                                :key="col"
                                                @click="toggleSort(col)"
                                                class="py-2.5 px-3 font-semibold whitespace-nowrap cursor-pointer hover:bg-slate-100 dark:hover:bg-surface-900 transition border-r border-slate-200/60 dark:border-surface-800/60 select-none group"
                                            >
                                                <div class="flex items-center justify-between gap-2">
                                                    <span>{{ col }}</span>
                                                    <span class="text-slate-400 dark:text-surface-500">
                                                        <ArrowUp v-if="sortColumn === col && sortDirection === 'asc'" class="w-3 h-3 text-emerald-500 font-bold" />
                                                        <ArrowDown v-else-if="sortColumn === col && sortDirection === 'desc'" class="w-3 h-3 text-emerald-500 font-bold" />
                                                        <ArrowUpDown v-else class="w-3 h-3 opacity-0 group-hover:opacity-100 transition" />
                                                    </span>
                                                </div>
                                            </th>
                                        </tr>
                                    </thead>
                                    <tbody class="divide-y divide-slate-100 dark:divide-surface-800/60 font-mono text-[11px]">
                                        <tr
                                            v-for="(row, rIdx) in tableData.rows"
                                            :key="rIdx"
                                            class="hover:bg-slate-50/80 dark:hover:bg-surface-800/40 transition group"
                                        >
                                            <td class="py-2 px-3 text-center text-slate-400 dark:text-surface-500 bg-slate-50/40 dark:bg-surface-950/20 border-r border-slate-200/60 dark:border-surface-800/60">
                                                {{ (currentPage - 1) * perPage + rIdx + 1 }}
                                            </td>
                                            <td
                                                v-for="col in tableData.columns"
                                                :key="col"
                                                class="py-2 px-3 border-r border-slate-200/40 dark:border-surface-800/40 max-w-xs truncate"
                                            >
                                                <!-- NULL value -->
                                                <span v-if="row[col] === null || row[col] === undefined" class="px-1.5 py-0.5 rounded text-[10px] bg-slate-100 dark:bg-surface-800 text-slate-400 dark:text-surface-500 italic">
                                                    NULL
                                                </span>

                                                <!-- Expandable JSON or long text -->
                                                <div v-else-if="isExpandable(row[col])" class="flex items-center justify-between gap-1.5">
                                                    <span class="truncate text-slate-700 dark:text-surface-200">{{ String(row[col]) }}</span>
                                                    <button
                                                        @click="openCellModal(col, row[col])"
                                                        class="px-1.5 py-0.5 rounded bg-slate-100 hover:bg-slate-200 dark:bg-surface-800 dark:hover:bg-surface-700 text-brand-600 dark:text-brand-400 text-[10px] shrink-0 font-medium"
                                                        title="View full content"
                                                    >
                                                        Expand
                                                    </button>
                                                </div>

                                                <!-- Boolean value -->
                                                <span v-else-if="row[col] === true || row[col] === 'true' || row[col] === 1 && (col.startsWith('is_') || col.startsWith('has_'))" class="px-1.5 py-0.5 rounded text-[10px] font-semibold bg-emerald-50 text-emerald-700 border border-emerald-200 dark:bg-emerald-500/10 dark:border-emerald-500/30 dark:text-emerald-400">
                                                    TRUE
                                                </span>
                                                <span v-else-if="row[col] === false || row[col] === 'false' || row[col] === 0 && (col.startsWith('is_') || col.startsWith('has_'))" class="px-1.5 py-0.5 rounded text-[10px] font-semibold bg-slate-100 text-slate-600 dark:bg-surface-800 dark:text-surface-400">
                                                    FALSE
                                                </span>

                                                <!-- Standard text / number value -->
                                                <span v-else class="text-slate-800 dark:text-surface-200">
                                                    {{ row[col] }}
                                                </span>
                                            </td>
                                        </tr>
                                    </tbody>
                                </table>
                            </div>

                            <!-- Data Pagination Bar -->
                            <div v-if="tableData && tableData.total_rows > 0" class="p-3 sm:px-4 border-t border-slate-200/80 dark:border-surface-800 flex flex-col sm:flex-row items-center justify-between gap-3 text-xs bg-slate-50/50 dark:bg-surface-950/40">
                                <div class="text-slate-500 dark:text-surface-400 font-mono text-[11px]">
                                    Showing <strong class="text-slate-900 dark:text-white">{{ (currentPage - 1) * perPage + 1 }}</strong> to <strong class="text-slate-900 dark:text-white">{{ Math.min(currentPage * perPage, tableData.total_rows) }}</strong> of <strong class="text-slate-900 dark:text-white">{{ formatNumber(tableData.total_rows) }}</strong> records
                                </div>

                                <div class="flex items-center gap-1.5 font-mono">
                                    <!-- First Page -->
                                    <button
                                        @click="changePage(1)"
                                        :disabled="currentPage === 1"
                                        class="p-1.5 rounded-lg bg-white dark:bg-surface-900 border border-slate-200 dark:border-surface-800 text-slate-600 dark:text-surface-300 disabled:opacity-40 hover:bg-slate-50 dark:hover:bg-surface-800 transition"
                                        title="First Page"
                                    >
                                        <ChevronsLeft class="w-3.5 h-3.5" />
                                    </button>

                                    <!-- Previous Page -->
                                    <button
                                        @click="changePage(currentPage - 1)"
                                        :disabled="currentPage === 1"
                                        class="p-1.5 rounded-lg bg-white dark:bg-surface-900 border border-slate-200 dark:border-surface-800 text-slate-600 dark:text-surface-300 disabled:opacity-40 hover:bg-slate-50 dark:hover:bg-surface-800 transition"
                                        title="Previous Page"
                                    >
                                        <ChevronLeft class="w-3.5 h-3.5" />
                                    </button>

                                    <!-- Current Page Indicator -->
                                    <span class="px-3 py-1 rounded-lg bg-white dark:bg-surface-900 border border-slate-200 dark:border-surface-800 text-slate-900 dark:text-white font-semibold text-[11px]">
                                        Page {{ currentPage }} of {{ tableData.total_pages || 1 }}
                                    </span>

                                    <!-- Next Page -->
                                    <button
                                        @click="changePage(currentPage + 1)"
                                        :disabled="currentPage >= tableData.total_pages"
                                        class="p-1.5 rounded-lg bg-white dark:bg-surface-900 border border-slate-200 dark:border-surface-800 text-slate-600 dark:text-surface-300 disabled:opacity-40 hover:bg-slate-50 dark:hover:bg-surface-800 transition"
                                        title="Next Page"
                                    >
                                        <ChevronRight class="w-3.5 h-3.5" />
                                    </button>

                                    <!-- Last Page -->
                                    <button
                                        @click="changePage(tableData.total_pages)"
                                        :disabled="currentPage >= tableData.total_pages"
                                        class="p-1.5 rounded-lg bg-white dark:bg-surface-900 border border-slate-200 dark:border-surface-800 text-slate-600 dark:text-surface-300 disabled:opacity-40 hover:bg-slate-50 dark:hover:bg-surface-800 transition"
                                        title="Last Page"
                                    >
                                        <ChevronsRight class="w-3.5 h-3.5" />
                                    </button>
                                </div>
                            </div>
                        </div>

                        <!-- TAB 2: STRUCTURE & SCHEMA -->
                        <div v-if="activeTab === 'structure'" class="p-4 sm:p-5 space-y-6">
                            <!-- Columns Section -->
                            <div class="space-y-3">
                                <div class="flex items-center justify-between">
                                    <h3 class="text-xs font-bold uppercase tracking-wider text-slate-700 dark:text-surface-300 flex items-center gap-1.5">
                                        <Layers class="w-4 h-4 text-emerald-500" />
                                        <span>Columns & Attributes ({{ structure?.columns.length || 0 }})</span>
                                    </h3>
                                </div>

                                <div class="rounded-xl border border-slate-200/80 dark:border-surface-800 overflow-x-auto shadow-xs">
                                    <table class="w-full text-left text-xs">
                                        <thead class="bg-slate-50 dark:bg-surface-950/60 text-slate-500 dark:text-surface-400 uppercase text-[10px] tracking-wider border-b border-slate-200/80 dark:border-surface-800">
                                            <tr>
                                                <th class="py-2.5 px-3 font-semibold w-10">#</th>
                                                <th class="py-2.5 px-3 font-semibold">Column Name</th>
                                                <th class="py-2.5 px-3 font-semibold">Data Type</th>
                                                <th class="py-2.5 px-3 font-semibold">Nullable</th>
                                                <th class="py-2.5 px-3 font-semibold">Key / Index</th>
                                                <th class="py-2.5 px-3 font-semibold">Default</th>
                                                <th class="py-2.5 px-3 font-semibold">Extra</th>
                                                <th class="py-2.5 px-3 font-semibold">Comment</th>
                                            </tr>
                                        </thead>
                                        <tbody class="divide-y divide-slate-100 dark:divide-surface-800/60 font-mono text-[11px]">
                                            <tr
                                                v-for="col in structure?.columns || []"
                                                :key="col.name"
                                                class="hover:bg-slate-50/80 dark:hover:bg-surface-800/30 transition"
                                            >
                                                <td class="py-2 px-3 text-slate-400 dark:text-surface-500">{{ col.position }}</td>
                                                <td class="py-2 px-3 font-semibold text-slate-900 dark:text-white flex items-center gap-1.5">
                                                    <Key v-if="col.key === 'PRI'" class="w-3 h-3 text-amber-500 shrink-0" />
                                                    <span>{{ col.name }}</span>
                                                </td>
                                                <td class="py-2 px-3 text-emerald-600 dark:text-emerald-400">{{ col.type }}</td>
                                                <td class="py-2 px-3">
                                                    <span
                                                        :class="[
                                                            'px-1.5 py-0.5 rounded text-[10px] font-semibold',
                                                            col.is_nullable
                                                                ? 'bg-amber-50 text-amber-700 dark:bg-amber-500/10 dark:text-amber-400'
                                                                : 'bg-slate-100 text-slate-600 dark:bg-surface-800 dark:text-surface-400'
                                                        ]"
                                                    >
                                                        {{ col.is_nullable ? 'YES' : 'NO' }}
                                                    </span>
                                                </td>
                                                <td class="py-2 px-3">
                                                    <span
                                                        v-if="col.key"
                                                        :class="[
                                                            'px-1.5 py-0.5 rounded text-[10px] font-bold',
                                                            col.key === 'PRI'
                                                                ? 'bg-amber-50 text-amber-700 border border-amber-200 dark:bg-amber-500/10 dark:border-amber-500/30 dark:text-amber-400'
                                                                : col.key === 'UNI'
                                                                ? 'bg-blue-50 text-blue-700 border border-blue-200 dark:bg-blue-500/10 dark:border-blue-500/30 dark:text-blue-400'
                                                                : 'bg-slate-100 text-slate-700 dark:bg-surface-800 dark:text-surface-300'
                                                        ]"
                                                    >
                                                        {{ col.key }}
                                                    </span>
                                                    <span v-else class="text-slate-400 dark:text-surface-600">-</span>
                                                </td>
                                                <td class="py-2 px-3 text-slate-600 dark:text-surface-400">
                                                    {{ col.default !== null ? col.default : 'NULL' }}
                                                </td>
                                                <td class="py-2 px-3 text-slate-500 dark:text-surface-400">{{ col.extra || '-' }}</td>
                                                <td class="py-2 px-3 text-slate-400 dark:text-surface-500">{{ col.comment || '-' }}</td>
                                            </tr>
                                        </tbody>
                                    </table>
                                </div>
                            </div>

                            <!-- Indexes Section -->
                            <div v-if="structure?.indexes && structure.indexes.length > 0" class="space-y-3">
                                <h3 class="text-xs font-bold uppercase tracking-wider text-slate-700 dark:text-surface-300 flex items-center gap-1.5">
                                    <Key class="w-4 h-4 text-amber-500" />
                                    <span>Indexes & Keys ({{ structure.indexes.length }})</span>
                                </h3>

                                <div class="rounded-xl border border-slate-200/80 dark:border-surface-800 overflow-x-auto shadow-xs">
                                    <table class="w-full text-left text-xs">
                                        <thead class="bg-slate-50 dark:bg-surface-950/60 text-slate-500 dark:text-surface-400 uppercase text-[10px] tracking-wider border-b border-slate-200/80 dark:border-surface-800">
                                            <tr>
                                                <th class="py-2.5 px-3 font-semibold">Index Name</th>
                                                <th class="py-2.5 px-3 font-semibold">Column</th>
                                                <th class="py-2.5 px-3 font-semibold">Type</th>
                                                <th class="py-2.5 px-3 font-semibold">Unique</th>
                                                <th class="py-2.5 px-3 font-semibold">Primary</th>
                                            </tr>
                                        </thead>
                                        <tbody class="divide-y divide-slate-100 dark:divide-surface-800/60 font-mono text-[11px]">
                                            <tr
                                                v-for="idx in structure.indexes"
                                                :key="idx.name + idx.column"
                                                class="hover:bg-slate-50/80 dark:hover:bg-surface-800/30 transition"
                                            >
                                                <td class="py-2 px-3 font-semibold text-slate-900 dark:text-white">{{ idx.name }}</td>
                                                <td class="py-2 px-3 text-emerald-600 dark:text-emerald-400">{{ idx.column }}</td>
                                                <td class="py-2 px-3 text-slate-600 dark:text-surface-400">{{ idx.type }}</td>
                                                <td class="py-2 px-3">
                                                    <span :class="['px-1.5 py-0.5 rounded text-[10px]', !idx.non_unique ? 'bg-emerald-50 text-emerald-700 dark:bg-emerald-500/10 dark:text-emerald-400 font-semibold' : 'text-slate-400 dark:text-surface-500']">
                                                        {{ !idx.non_unique ? 'YES' : 'NO' }}
                                                    </span>
                                                </td>
                                                <td class="py-2 px-3">
                                                    <span :class="['px-1.5 py-0.5 rounded text-[10px]', idx.is_primary ? 'bg-amber-50 text-amber-700 dark:bg-amber-500/10 dark:text-amber-400 font-semibold' : 'text-slate-400 dark:text-surface-500']">
                                                        {{ idx.is_primary ? 'YES' : 'NO' }}
                                                    </span>
                                                </td>
                                            </tr>
                                        </tbody>
                                    </table>
                                </div>
                            </div>

                            <!-- Foreign Keys Section -->
                            <div v-if="structure?.foreign_keys && structure.foreign_keys.length > 0" class="space-y-3">
                                <h3 class="text-xs font-bold uppercase tracking-wider text-slate-700 dark:text-surface-300 flex items-center gap-1.5">
                                    <ExternalLink class="w-4 h-4 text-blue-500" />
                                    <span>Foreign Key Constraints ({{ structure.foreign_keys.length }})</span>
                                </h3>

                                <div class="rounded-xl border border-slate-200/80 dark:border-surface-800 overflow-x-auto shadow-xs">
                                    <table class="w-full text-left text-xs">
                                        <thead class="bg-slate-50 dark:bg-surface-950/60 text-slate-500 dark:text-surface-400 uppercase text-[10px] tracking-wider border-b border-slate-200/80 dark:border-surface-800">
                                            <tr>
                                                <th class="py-2.5 px-3 font-semibold">Constraint Name</th>
                                                <th class="py-2.5 px-3 font-semibold">Local Column</th>
                                                <th class="py-2.5 px-3 font-semibold">Target (Referenced Table & Column)</th>
                                            </tr>
                                        </thead>
                                        <tbody class="divide-y divide-slate-100 dark:divide-surface-800/60 font-mono text-[11px]">
                                            <tr
                                                v-for="fk in structure.foreign_keys"
                                                :key="fk.constraint_name"
                                                class="hover:bg-slate-50/80 dark:hover:bg-surface-800/30 transition"
                                            >
                                                <td class="py-2 px-3 font-semibold text-slate-900 dark:text-white">{{ fk.constraint_name }}</td>
                                                <td class="py-2 px-3 text-emerald-600 dark:text-emerald-400">{{ fk.column }}</td>
                                                <td class="py-2 px-3 text-blue-600 dark:text-blue-400">
                                                    {{ fk.referenced_table }}.{{ fk.referenced_column }}
                                                </td>
                                            </tr>
                                        </tbody>
                                    </table>
                                </div>
                            </div>
                        </div>

                        <!-- TAB 3: SQL DDL -->
                        <div v-if="activeTab === 'ddl'" class="p-4 sm:p-5 space-y-3">
                            <div class="flex items-center justify-between">
                                <h3 class="text-xs font-bold uppercase tracking-wider text-slate-700 dark:text-surface-300 flex items-center gap-1.5">
                                    <Code class="w-4 h-4 text-emerald-500" />
                                    <span>Create Table Statement</span>
                                </h3>

                                <button
                                    @click="copyToClipboard(structure?.create_statement || '', 'ddl')"
                                    class="px-3 py-1.5 rounded-lg bg-slate-100 hover:bg-slate-200 dark:bg-surface-800 dark:hover:bg-surface-700 text-slate-700 dark:text-surface-200 text-xs font-semibold flex items-center gap-1.5 transition"
                                >
                                    <Check v-if="isDdlCopied" class="w-3.5 h-3.5 text-emerald-500" />
                                    <Copy v-else class="w-3.5 h-3.5 text-slate-500 dark:text-surface-400" />
                                    <span>{{ isDdlCopied ? 'Copied!' : 'Copy SQL' }}</span>
                                </button>
                            </div>

                            <div class="p-4 rounded-xl bg-slate-950 text-emerald-400 font-mono text-xs overflow-x-auto leading-relaxed border border-slate-800 shadow-inner">
                                <pre>{{ structure?.create_statement || '-- No DDL available' }}</pre>
                            </div>
                        </div>
                    </div>

                    <!-- No Table Selected Placeholder -->
                    <div v-else class="rounded-2xl bg-white dark:bg-surface-900/80 border border-slate-200/80 dark:border-surface-800/80 p-12 text-center shadow-sm">
                        <Table class="w-12 h-12 text-slate-300 dark:text-surface-600 mx-auto mb-3" />
                        <h3 class="text-base font-bold text-slate-800 dark:text-surface-100">Select a table to inspect</h3>
                        <p class="text-xs text-slate-500 dark:text-surface-400 mt-1 max-w-sm mx-auto">
                            Choose any table from the sidebar to view data rows, browse schema structure, columns, indexes, and SQL DDL.
                        </p>
                    </div>
                </div>
            </div>

            <!-- Cell Details / JSON Modal -->
            <div
                v-if="isCellModalOpen"
                class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/60 dark:bg-surface-950/80 backdrop-blur-sm"
            >
                <div class="w-full max-w-2xl rounded-2xl bg-white dark:bg-surface-900 border border-slate-200/90 dark:border-surface-800 p-5 sm:p-6 shadow-2xl space-y-4">
                    <div class="flex items-center justify-between border-b border-slate-100 dark:border-surface-800 pb-3">
                        <h3 class="text-sm font-bold text-slate-900 dark:text-white font-mono flex items-center gap-2">
                            <Info class="w-4 h-4 text-emerald-500" />
                            <span>Column: {{ activeCellColumn }}</span>
                        </h3>
                        <button
                            @click="isCellModalOpen = false"
                            class="p-1 rounded-lg text-slate-400 hover:text-slate-600 dark:text-surface-500 dark:hover:text-surface-300"
                        >
                            <X class="w-4 h-4" />
                        </button>
                    </div>

                    <div class="p-3.5 rounded-xl bg-slate-950 text-slate-200 font-mono text-xs overflow-auto max-h-96 leading-relaxed border border-slate-800">
                        <pre class="whitespace-pre-wrap">{{ activeCellContent }}</pre>
                    </div>

                    <div class="flex items-center justify-end gap-2 pt-2 border-t border-slate-100 dark:border-surface-800">
                        <button
                            type="button"
                            @click="copyToClipboard(activeCellContent, 'cell')"
                            class="px-3.5 py-2 rounded-xl bg-slate-100 hover:bg-slate-200 dark:bg-surface-800 dark:hover:bg-surface-700 text-slate-700 dark:text-surface-200 text-xs font-semibold flex items-center gap-1.5 transition"
                        >
                            <Check v-if="isCopied" class="w-3.5 h-3.5 text-emerald-500" />
                            <Copy v-else class="w-3.5 h-3.5 text-slate-500 dark:text-surface-400" />
                            <span>{{ isCopied ? 'Copied' : 'Copy' }}</span>
                        </button>
                        <button
                            type="button"
                            @click="isCellModalOpen = false"
                            class="px-4 py-2 rounded-xl bg-emerald-600 hover:bg-emerald-500 text-white text-xs font-semibold transition"
                        >
                            Close
                        </button>
                    </div>
                </div>
            </div>
        </div>
    </AppLayout>
</template>
