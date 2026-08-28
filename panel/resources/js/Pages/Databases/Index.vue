<script setup lang="ts">
import { ref } from 'vue';
import { useForm, router, Link } from '@inertiajs/vue3';
import AppLayout from '@/Layouts/AppLayout.vue';
import {
    Database,
    Table,
    Plus,
    Trash2,
    Shield,
    HardDrive,
    UserCheck,
    Users,
    Key,
    Lock,
    Unlock,
    Sparkles,
    Check,
    AlertTriangle,
} from 'lucide-vue-next';
import type { DatabaseRecord, DatabaseUser } from '@/types';

const props = defineProps<{
    databases: {
        data: (DatabaseRecord & { users?: DatabaseUser[] })[];
        links: any[];
    };
    databaseUsers: (DatabaseUser & { databases?: DatabaseRecord[] })[];
    allDatabases: { id: string | number; name: string; engine: string }[];
}>();

const activeTab = ref<'databases' | 'users'>('databases');

// Create Database Modal State & Form
const isCreateDbOpen = ref(false);
const dbForm = useForm({
    engine: 'mysql',
    name: '',
    character_set: 'utf8mb4',
    collation: 'utf8mb4_unicode_ci',
});

function submitCreateDb() {
    dbForm.post('/databases', {
        onSuccess: () => {
            isCreateDbOpen.value = false;
            dbForm.reset();
        },
    });
}

// Create User Modal State & Form
const isCreateUserOpen = ref(false);
const userForm = useForm({
    engine: 'mysql',
    username: '',
    host: 'localhost',
    password: '',
});

function generateRandomPassword() {
    const chars = 'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()_+';
    let pass = '';
    for (let i = 0; i < 18; i++) {
        pass += chars.charAt(Math.floor(Math.random() * chars.length));
    }
    userForm.password = pass;
}

function submitCreateUser() {
    userForm.post('/databases/users', {
        onSuccess: () => {
            isCreateUserOpen.value = false;
            userForm.reset();
        },
    });
}

// Grant Privileges Modal State & Form
const isGrantModalOpen = ref(false);
const grantForm = useForm({
    database_id: '',
    database_user_id: '',
    permissions: 'all',
});

function openGrantModal(dbId?: string | number, userId?: string | number) {
    if (dbId) grantForm.database_id = String(dbId);
    if (userId) grantForm.database_user_id = String(userId);
    isGrantModalOpen.value = true;
}

function submitGrant() {
    grantForm.post('/databases/grants', {
        onSuccess: () => {
            isGrantModalOpen.value = false;
            grantForm.reset();
        },
    });
}

// Reset Password Modal State & Form
const isResetPassOpen = ref(false);
const selectedUserForReset = ref<DatabaseUser | null>(null);
const resetPassForm = useForm({
    password: '',
});

function openResetPassModal(u: DatabaseUser) {
    selectedUserForReset.value = u;
    resetPassForm.password = '';
    isResetPassOpen.value = true;
}

function submitResetPass() {
    if (!selectedUserForReset.value) return;
    resetPassForm.post(`/databases/users/${selectedUserForReset.value.id}/password`, {
        onSuccess: () => {
            isResetPassOpen.value = false;
            resetPassForm.reset();
            selectedUserForReset.value = null;
        },
    });
}

// Delete Database
function deleteDatabase(db: DatabaseRecord) {
    if (confirm(`Are you sure you want to drop database ${db.name}? ALL TABLES AND DATA WILL BE PERMANENTLY DELETED.`)) {
        router.delete(`/databases/${db.id}`);
    }
}

// Delete User
function deleteUser(u: DatabaseUser) {
    if (confirm(`Are you sure you want to drop database user ${u.username}@${u.host}?`)) {
        router.delete(`/databases/users/${u.id}`);
    }
}
</script>

<template>
    <AppLayout title="Databases">
        <div class="space-y-6 max-w-7xl mx-auto">
            <!-- Header -->
            <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
                <div>
                    <h2 class="text-base font-bold text-slate-900 dark:text-white tracking-tight">Database Management</h2>
                    <p class="text-xs text-slate-500 dark:text-surface-400 mt-0.5">
                        Provision MySQL/MariaDB & PostgreSQL databases, manage users, and assign access permissions
                    </p>
                </div>
                <div class="flex items-center gap-2">
                    <button
                        @click="isCreateUserOpen = true"
                        class="px-3.5 py-2 rounded-xl bg-slate-100 hover:bg-slate-200 dark:bg-surface-800 dark:hover:bg-surface-700 text-slate-700 dark:text-surface-200 text-xs font-semibold flex items-center gap-1.5 transition border border-slate-200/80 dark:border-surface-700"
                    >
                        <Users class="w-4 h-4 text-brand-600 dark:text-brand-400" />
                        <span>New User</span>
                    </button>
                    <button
                        @click="isCreateDbOpen = true"
                        class="px-4 py-2 rounded-xl bg-emerald-600 hover:bg-emerald-500 text-white text-xs font-semibold shadow-md shadow-emerald-600/20 flex items-center gap-1.5 transition"
                    >
                        <Plus class="w-4 h-4" />
                        <span>Create Database</span>
                    </button>
                </div>
            </div>

            <!-- Tab Navigation -->
            <div class="flex items-center gap-2 border-b border-slate-200/80 dark:border-surface-800 pb-2 text-xs">
                <button
                    @click="activeTab = 'databases'"
                    :class="[
                        'px-4 py-2 rounded-xl font-semibold flex items-center gap-1.5 transition',
                        activeTab === 'databases'
                            ? 'bg-emerald-50 text-emerald-700 border border-emerald-200 dark:bg-emerald-500/10 dark:text-emerald-400 dark:border-emerald-500/30'
                            : 'text-slate-600 hover:text-slate-900 dark:text-surface-400 dark:hover:text-white'
                    ]"
                >
                    <Database class="w-4 h-4" />
                    <span>Databases ({{ databases.data.length }})</span>
                </button>

                <button
                    @click="activeTab = 'users'"
                    :class="[
                        'px-4 py-2 rounded-xl font-semibold flex items-center gap-1.5 transition',
                        activeTab === 'users'
                            ? 'bg-emerald-50 text-emerald-700 border border-emerald-200 dark:bg-emerald-500/10 dark:text-emerald-400 dark:border-emerald-500/30'
                            : 'text-slate-600 hover:text-slate-900 dark:text-surface-400 dark:hover:text-white'
                    ]"
                >
                    <Users class="w-4 h-4" />
                    <span>Database Users ({{ databaseUsers.length }})</span>
                </button>
            </div>

            <!-- Tab 1: Databases Table -->
            <div v-if="activeTab === 'databases'" class="rounded-2xl bg-white dark:bg-surface-900/60 border border-slate-200/80 dark:border-surface-800/80 shadow-sm dark:shadow-xl overflow-hidden">
                <div v-if="databases.data.length === 0" class="text-center py-14 px-4">
                    <Database class="w-10 h-10 text-slate-300 dark:text-surface-600 mx-auto mb-3" />
                    <h3 class="text-sm font-medium text-slate-700 dark:text-surface-200">No databases created yet</h3>
                    <p class="text-xs text-slate-500 dark:text-surface-400 mt-1 max-w-sm mx-auto">
                        Create MySQL or PostgreSQL databases with automated charset and privilege management.
                    </p>
                    <button
                        @click="isCreateDbOpen = true"
                        class="mt-4 inline-flex items-center gap-1.5 px-4 py-2 rounded-xl bg-emerald-600 text-white text-xs font-semibold hover:bg-emerald-500 transition shadow-md shadow-emerald-600/20"
                    >
                        <Plus class="w-3.5 h-3.5" />
                        <span>Create First Database</span>
                    </button>
                </div>

                <div v-else class="overflow-x-auto">
                    <table class="w-full text-left text-xs">
                        <thead class="bg-slate-50 dark:bg-surface-950/50 text-slate-500 dark:text-surface-400 uppercase text-[10px] tracking-wider border-b border-slate-200/80 dark:border-surface-800">
                            <tr>
                                <th class="py-3.5 px-4 font-semibold">Database Name</th>
                                <th class="py-3.5 px-4 font-semibold">Engine</th>
                                <th class="py-3.5 px-4 font-semibold">Charset & Collation</th>
                                <th class="py-3.5 px-4 font-semibold">Assigned Users</th>
                                <th class="py-3.5 px-4 font-semibold">Created</th>
                                <th class="py-3.5 px-4 font-semibold text-right">Actions</th>
                            </tr>
                        </thead>
                        <tbody class="divide-y divide-slate-100 dark:divide-surface-800/60">
                            <tr
                                v-for="db in databases.data"
                                :key="db.id"
                                class="hover:bg-slate-50/80 dark:hover:bg-surface-800/30 transition"
                            >
                                <td class="py-3.5 px-4 font-semibold text-slate-900 dark:text-white font-mono">
                                    <div class="flex items-center gap-2.5">
                                        <div class="p-2 rounded-xl bg-emerald-50 text-emerald-600 dark:bg-emerald-500/10 dark:border dark:border-emerald-500/20 dark:text-emerald-400 shrink-0">
                                            <Database class="w-4 h-4" />
                                        </div>
                                        <span>{{ db.name }}</span>
                                    </div>
                                </td>
                                <td class="py-3.5 px-4">
                                    <span
                                        :class="[
                                            'px-2 py-0.5 rounded-md font-mono text-[11px] uppercase font-semibold',
                                            db.engine === 'mysql'
                                                ? 'bg-blue-50 text-blue-700 border border-blue-200 dark:bg-blue-500/10 dark:border-blue-500/30 dark:text-blue-400'
                                                : 'bg-indigo-50 text-indigo-700 border border-indigo-200 dark:bg-indigo-500/10 dark:border-indigo-500/30 dark:text-indigo-400'
                                        ]"
                                    >
                                        {{ db.engine }}
                                    </span>
                                </td>
                                <td class="py-3.5 px-4 font-mono text-[11px] text-slate-500 dark:text-surface-400">
                                    {{ db.character_set }} <span class="text-slate-300 dark:text-surface-600">•</span> {{ db.collation }}
                                </td>
                                <td class="py-3.5 px-4">
                                    <div v-if="db.users && db.users.length > 0" class="flex flex-wrap gap-1.5">
                                        <span
                                            v-for="u in db.users"
                                            :key="u.id"
                                            class="px-2 py-0.5 rounded bg-slate-100 dark:bg-surface-800 text-slate-700 dark:text-surface-300 font-mono text-[10px] flex items-center gap-1 border border-slate-200 dark:border-surface-700"
                                        >
                                            <UserCheck class="w-3 h-3 text-emerald-600 dark:text-emerald-400" />
                                            <span>{{ u.username }}</span>
                                        </span>
                                    </div>
                                    <button
                                        v-else
                                        @click="openGrantModal(db.id)"
                                        class="text-[11px] text-brand-600 dark:text-brand-400 hover:underline flex items-center gap-1 font-medium"
                                    >
                                        <Plus class="w-3 h-3" />
                                        <span>Assign User</span>
                                    </button>
                                </td>
                                <td class="py-3.5 px-4 font-mono text-[11px] text-slate-400 dark:text-surface-500">
                                    {{ new Date(db.created_at).toLocaleDateString() }}
                                </td>
                                <td class="py-3.5 px-4 text-right">
                                    <div class="flex items-center justify-end gap-1.5">
                                        <Link
                                            :href="`/databases/${db.id}/explorer`"
                                            class="px-2.5 py-1 rounded-lg bg-emerald-50 hover:bg-emerald-100 dark:bg-emerald-500/10 dark:hover:bg-emerald-500/20 text-emerald-700 dark:text-emerald-400 text-[11px] font-semibold flex items-center gap-1 transition border border-emerald-200/80 dark:border-emerald-500/30"
                                            title="Explore Tables and Data"
                                        >
                                            <Table class="w-3.5 h-3.5" />
                                            <span>Explore</span>
                                        </Link>
                                        <button
                                            @click="openGrantModal(db.id)"
                                            class="px-2.5 py-1 rounded-lg bg-slate-100 hover:bg-slate-200 dark:bg-surface-800 dark:hover:bg-surface-700 text-slate-700 dark:text-surface-300 text-[11px] transition"
                                            title="Assign User"
                                        >
                                            Assign User
                                        </button>
                                        <button
                                            @click="deleteDatabase(db)"
                                            class="p-1.5 rounded-lg text-rose-500 hover:bg-rose-50 dark:hover:bg-rose-500/10 transition"
                                            title="Drop Database"
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

            <!-- Tab 2: Database Users Table -->
            <div v-if="activeTab === 'users'" class="rounded-2xl bg-white dark:bg-surface-900/60 border border-slate-200/80 dark:border-surface-800/80 shadow-sm dark:shadow-xl overflow-hidden">
                <div v-if="databaseUsers.length === 0" class="text-center py-14 px-4">
                    <Users class="w-10 h-10 text-slate-300 dark:text-surface-600 mx-auto mb-3" />
                    <h3 class="text-sm font-medium text-slate-700 dark:text-surface-200">No database users created yet</h3>
                    <p class="text-xs text-slate-500 dark:text-surface-400 mt-1 max-w-sm mx-auto">
                        Create dedicated database users with host restrictions and passwords.
                    </p>
                    <button
                        @click="isCreateUserOpen = true"
                        class="mt-4 inline-flex items-center gap-1.5 px-4 py-2 rounded-xl bg-brand-600 text-white text-xs font-semibold hover:bg-brand-500 transition shadow-md shadow-brand-600/20"
                    >
                        <Plus class="w-3.5 h-3.5" />
                        <span>Create Database User</span>
                    </button>
                </div>

                <div v-else class="overflow-x-auto">
                    <table class="w-full text-left text-xs">
                        <thead class="bg-slate-50 dark:bg-surface-950/50 text-slate-500 dark:text-surface-400 uppercase text-[10px] tracking-wider border-b border-slate-200/80 dark:border-surface-800">
                            <tr>
                                <th class="py-3.5 px-4 font-semibold">Username & Host</th>
                                <th class="py-3.5 px-4 font-semibold">Engine</th>
                                <th class="py-3.5 px-4 font-semibold">Authorized Databases</th>
                                <th class="py-3.5 px-4 font-semibold">Created</th>
                                <th class="py-3.5 px-4 font-semibold text-right">Actions</th>
                            </tr>
                        </thead>
                        <tbody class="divide-y divide-slate-100 dark:divide-surface-800/60">
                            <tr
                                v-for="u in databaseUsers"
                                :key="u.id"
                                class="hover:bg-slate-50/80 dark:hover:bg-surface-800/30 transition"
                            >
                                <td class="py-3.5 px-4 font-semibold text-slate-900 dark:text-white font-mono">
                                    <div class="flex items-center gap-2">
                                        <Users class="w-4 h-4 text-brand-600 dark:text-brand-400 shrink-0" />
                                        <span>{{ u.username }}</span>
                                        <span class="text-slate-400 dark:text-surface-500 text-[10px]">@{{ u.host }}</span>
                                    </div>
                                </td>
                                <td class="py-3.5 px-4">
                                    <span
                                        :class="[
                                            'px-2 py-0.5 rounded-md font-mono text-[11px] uppercase font-semibold',
                                            u.engine === 'mysql'
                                                ? 'bg-blue-50 text-blue-700 border border-blue-200 dark:bg-blue-500/10 dark:border-blue-500/30 dark:text-blue-400'
                                                : 'bg-indigo-50 text-indigo-700 border border-indigo-200 dark:bg-indigo-500/10 dark:border-indigo-500/30 dark:text-indigo-400'
                                        ]"
                                    >
                                        {{ u.engine }}
                                    </span>
                                </td>
                                <td class="py-3.5 px-4">
                                    <div v-if="u.databases && u.databases.length > 0" class="flex flex-wrap gap-1.5">
                                        <span
                                            v-for="d in u.databases"
                                            :key="d.id"
                                            class="px-2 py-0.5 rounded bg-slate-100 dark:bg-surface-800 text-slate-700 dark:text-surface-300 font-mono text-[10px] border border-slate-200 dark:border-surface-700"
                                        >
                                            {{ d.name }}
                                        </span>
                                    </div>
                                    <span v-else class="text-slate-400 dark:text-surface-500 text-[11px] italic">No databases assigned</span>
                                </td>
                                <td class="py-3.5 px-4 font-mono text-[11px] text-slate-400 dark:text-surface-500">
                                    {{ new Date(u.created_at).toLocaleDateString() }}
                                </td>
                                <td class="py-3.5 px-4 text-right">
                                    <div class="flex items-center justify-end gap-1.5">
                                        <button
                                            @click="openResetPassModal(u)"
                                            class="px-2.5 py-1 rounded-lg bg-slate-100 hover:bg-slate-200 dark:bg-surface-800 dark:hover:bg-surface-700 text-slate-700 dark:text-surface-300 text-[11px] flex items-center gap-1 transition"
                                            title="Change Password"
                                        >
                                            <Key class="w-3 h-3 text-amber-500" />
                                            <span>Reset Pass</span>
                                        </button>
                                        <button
                                            @click="deleteUser(u)"
                                            class="p-1.5 rounded-lg text-rose-500 hover:bg-rose-50 dark:hover:bg-rose-500/10 transition"
                                            title="Drop User"
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

            <!-- Create Database Modal -->
            <div
                v-if="isCreateDbOpen"
                class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/60 dark:bg-surface-950/80 backdrop-blur-sm"
            >
                <div class="w-full max-w-md rounded-2xl bg-white dark:bg-surface-900 border border-slate-200/90 dark:border-surface-800 p-6 shadow-2xl space-y-4">
                    <h3 class="text-sm font-semibold text-slate-900 dark:text-white">Create New Database</h3>

                    <form @submit.prevent="submitCreateDb" class="space-y-4">
                        <div>
                            <label class="block text-xs font-medium text-slate-700 dark:text-surface-300 mb-1">Database Engine</label>
                            <select
                                v-model="dbForm.engine"
                                class="w-full px-3.5 py-2.5 rounded-xl bg-slate-50 dark:bg-surface-950/60 border border-slate-200 dark:border-surface-800 text-xs text-slate-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-emerald-500/40"
                            >
                                <option value="mysql">MySQL / MariaDB</option>
                                <option value="postgresql">PostgreSQL</option>
                            </select>
                        </div>

                        <div>
                            <label class="block text-xs font-medium text-slate-700 dark:text-surface-300 mb-1">Database Name</label>
                            <input
                                v-model="dbForm.name"
                                type="text"
                                required
                                placeholder="app_production"
                                class="w-full px-3.5 py-2.5 rounded-xl bg-slate-50 dark:bg-surface-950/60 border border-slate-200 dark:border-surface-800 text-xs text-slate-900 dark:text-white placeholder-slate-400 dark:placeholder-surface-500 font-mono focus:outline-none focus:ring-2 focus:ring-emerald-500/40"
                            />
                            <p v-if="dbForm.errors.name" class="text-[11px] text-rose-500 mt-1">{{ dbForm.errors.name }}</p>
                        </div>

                        <div class="flex items-center justify-end gap-3 pt-3 border-t border-slate-100 dark:border-surface-800">
                            <button
                                type="button"
                                @click="isCreateDbOpen = false"
                                class="px-3.5 py-2 rounded-xl text-xs text-slate-600 hover:text-slate-900 dark:text-surface-400 dark:hover:text-white"
                            >
                                Cancel
                            </button>
                            <button
                                type="submit"
                                :disabled="dbForm.processing"
                                class="px-4 py-2 rounded-xl bg-emerald-600 hover:bg-emerald-500 text-white text-xs font-semibold shadow-md shadow-emerald-600/20 disabled:opacity-50"
                            >
                                {{ dbForm.processing ? 'Creating...' : 'Create Database' }}
                            </button>
                        </div>
                    </form>
                </div>
            </div>

            <!-- Create User Modal -->
            <div
                v-if="isCreateUserOpen"
                class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/60 dark:bg-surface-950/80 backdrop-blur-sm"
            >
                <div class="w-full max-w-md rounded-2xl bg-white dark:bg-surface-900 border border-slate-200/90 dark:border-surface-800 p-6 shadow-2xl space-y-4">
                    <h3 class="text-sm font-semibold text-slate-900 dark:text-white">Create Database User</h3>

                    <form @submit.prevent="submitCreateUser" class="space-y-4">
                        <div>
                            <label class="block text-xs font-medium text-slate-700 dark:text-surface-300 mb-1">Database Engine</label>
                            <select
                                v-model="userForm.engine"
                                class="w-full px-3.5 py-2.5 rounded-xl bg-slate-50 dark:bg-surface-950/60 border border-slate-200 dark:border-surface-800 text-xs text-slate-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-brand-500/40"
                            >
                                <option value="mysql">MySQL / MariaDB</option>
                                <option value="postgresql">PostgreSQL</option>
                            </select>
                        </div>

                        <div>
                            <label class="block text-xs font-medium text-slate-700 dark:text-surface-300 mb-1">Username</label>
                            <input
                                v-model="userForm.username"
                                type="text"
                                required
                                placeholder="app_dbuser"
                                class="w-full px-3.5 py-2.5 rounded-xl bg-slate-50 dark:bg-surface-950/60 border border-slate-200 dark:border-surface-800 text-xs text-slate-900 dark:text-white placeholder-slate-400 dark:placeholder-surface-500 font-mono focus:outline-none focus:ring-2 focus:ring-brand-500/40"
                            />
                            <p v-if="userForm.errors.username" class="text-[11px] text-rose-500 mt-1">{{ userForm.errors.username }}</p>
                        </div>

                        <div v-if="userForm.engine === 'mysql'">
                            <label class="block text-xs font-medium text-slate-700 dark:text-surface-300 mb-1">Host Restriction</label>
                            <input
                                v-model="userForm.host"
                                type="text"
                                placeholder="localhost"
                                class="w-full px-3.5 py-2.5 rounded-xl bg-slate-50 dark:bg-surface-950/60 border border-slate-200 dark:border-surface-800 text-xs text-slate-900 dark:text-white font-mono focus:outline-none focus:ring-2 focus:ring-brand-500/40"
                            />
                            <p class="text-[10px] text-slate-400 dark:text-surface-500 mt-1">Use <code class="text-brand-600 dark:text-brand-400 font-mono">localhost</code> for local connections or <code class="text-brand-600 dark:text-brand-400 font-mono">%</code> for any remote host.</p>
                        </div>

                        <div>
                            <div class="flex items-center justify-between mb-1">
                                <label class="text-xs font-medium text-slate-700 dark:text-surface-300">Password</label>
                                <button
                                    type="button"
                                    @click="generateRandomPassword"
                                    class="text-[11px] text-brand-600 dark:text-brand-400 hover:underline flex items-center gap-1 font-medium"
                                >
                                    <Sparkles class="w-3 h-3" />
                                    <span>Generate Strong</span>
                                </button>
                            </div>
                            <input
                                v-model="userForm.password"
                                type="text"
                                required
                                placeholder="Enter or generate password"
                                class="w-full px-3.5 py-2.5 rounded-xl bg-slate-50 dark:bg-surface-950/60 border border-slate-200 dark:border-surface-800 text-xs text-slate-900 dark:text-white placeholder-slate-400 dark:placeholder-surface-500 font-mono focus:outline-none focus:ring-2 focus:ring-brand-500/40"
                            />
                            <p v-if="userForm.errors.password" class="text-[11px] text-rose-500 mt-1">{{ userForm.errors.password }}</p>
                        </div>

                        <div class="flex items-center justify-end gap-3 pt-3 border-t border-slate-100 dark:border-surface-800">
                            <button
                                type="button"
                                @click="isCreateUserOpen = false"
                                class="px-3.5 py-2 rounded-xl text-xs text-slate-600 hover:text-slate-900 dark:text-surface-400 dark:hover:text-white"
                            >
                                Cancel
                            </button>
                            <button
                                type="submit"
                                :disabled="userForm.processing"
                                class="px-4 py-2 rounded-xl bg-brand-600 hover:bg-brand-500 text-white text-xs font-semibold shadow-md shadow-brand-600/20 disabled:opacity-50"
                            >
                                {{ userForm.processing ? 'Creating User...' : 'Create User' }}
                            </button>
                        </div>
                    </form>
                </div>
            </div>

            <!-- Grant Privileges Modal -->
            <div
                v-if="isGrantModalOpen"
                class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/60 dark:bg-surface-950/80 backdrop-blur-sm"
            >
                <div class="w-full max-w-md rounded-2xl bg-white dark:bg-surface-900 border border-slate-200/90 dark:border-surface-800 p-6 shadow-2xl space-y-4">
                    <h3 class="text-sm font-semibold text-slate-900 dark:text-white">Assign User to Database</h3>

                    <form @submit.prevent="submitGrant" class="space-y-4">
                        <div>
                            <label class="block text-xs font-medium text-slate-700 dark:text-surface-300 mb-1">Target Database</label>
                            <select
                                v-model="grantForm.database_id"
                                required
                                class="w-full px-3.5 py-2.5 rounded-xl bg-slate-50 dark:bg-surface-950/60 border border-slate-200 dark:border-surface-800 text-xs text-slate-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-brand-500/40"
                            >
                                <option value="" disabled>Select a database</option>
                                <option v-for="d in allDatabases" :key="d.id" :value="d.id">
                                    {{ d.name }} ({{ d.engine.toUpperCase() }})
                                </option>
                            </select>
                        </div>

                        <div>
                            <label class="block text-xs font-medium text-slate-700 dark:text-surface-300 mb-1">Database User</label>
                            <select
                                v-model="grantForm.database_user_id"
                                required
                                class="w-full px-3.5 py-2.5 rounded-xl bg-slate-50 dark:bg-surface-950/60 border border-slate-200 dark:border-surface-800 text-xs text-slate-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-brand-500/40"
                            >
                                <option value="" disabled>Select a user</option>
                                <option v-for="u in databaseUsers" :key="u.id" :value="u.id">
                                    {{ u.username }} ({{ u.engine.toUpperCase() }}@{{ u.host }})
                                </option>
                            </select>
                        </div>

                        <div>
                            <label class="block text-xs font-medium text-slate-700 dark:text-surface-300 mb-1">Privileges</label>
                            <select
                                v-model="grantForm.permissions"
                                class="w-full px-3.5 py-2.5 rounded-xl bg-slate-50 dark:bg-surface-950/60 border border-slate-200 dark:border-surface-800 text-xs text-slate-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-brand-500/40"
                            >
                                <option value="all">ALL PRIVILEGES (Read, Write, DDL)</option>
                                <option value="read_only">READ ONLY (SELECT)</option>
                            </select>
                        </div>

                        <div class="flex items-center justify-end gap-3 pt-3 border-t border-slate-100 dark:border-surface-800">
                            <button
                                type="button"
                                @click="isGrantModalOpen = false"
                                class="px-3.5 py-2 rounded-xl text-xs text-slate-600 hover:text-slate-900 dark:text-surface-400 dark:hover:text-white"
                            >
                                Cancel
                            </button>
                            <button
                                type="submit"
                                :disabled="grantForm.processing"
                                class="px-4 py-2 rounded-xl bg-emerald-600 hover:bg-emerald-500 text-white text-xs font-semibold shadow-md shadow-emerald-600/20 disabled:opacity-50"
                            >
                                {{ grantForm.processing ? 'Granting...' : 'Grant Access' }}
                            </button>
                        </div>
                    </form>
                </div>
            </div>

            <!-- Reset Password Modal -->
            <div
                v-if="isResetPassOpen"
                class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/60 dark:bg-surface-950/80 backdrop-blur-sm"
            >
                <div class="w-full max-w-md rounded-2xl bg-white dark:bg-surface-900 border border-slate-200/90 dark:border-surface-800 p-6 shadow-2xl space-y-4">
                    <h3 class="text-sm font-semibold text-slate-900 dark:text-white">
                        Reset Password for <strong class="text-emerald-600 dark:text-emerald-400 font-mono">{{ selectedUserForReset?.username }}</strong>
                    </h3>

                    <form @submit.prevent="submitResetPass" class="space-y-4">
                        <div>
                            <label class="block text-xs font-medium text-slate-700 dark:text-surface-300 mb-1">New Password</label>
                            <input
                                v-model="resetPassForm.password"
                                type="text"
                                required
                                placeholder="Enter new password"
                                class="w-full px-3.5 py-2.5 rounded-xl bg-slate-50 dark:bg-surface-950/60 border border-slate-200 dark:border-surface-800 text-xs text-slate-900 dark:text-white placeholder-slate-400 dark:placeholder-surface-500 font-mono focus:outline-none focus:ring-2 focus:ring-amber-500/40"
                            />
                        </div>

                        <div class="flex items-center justify-end gap-3 pt-3 border-t border-slate-100 dark:border-surface-800">
                            <button
                                type="button"
                                @click="isResetPassOpen = false"
                                class="px-3.5 py-2 rounded-xl text-xs text-slate-600 hover:text-slate-900 dark:text-surface-400 dark:hover:text-white"
                            >
                                Cancel
                            </button>
                            <button
                                type="submit"
                                :disabled="resetPassForm.processing"
                                class="px-4 py-2 rounded-xl bg-amber-600 hover:bg-amber-500 text-white text-xs font-semibold shadow-md shadow-amber-600/20 disabled:opacity-50"
                            >
                                {{ resetPassForm.processing ? 'Updating...' : 'Update Password' }}
                            </button>
                        </div>
                    </form>
                </div>
            </div>
        </div>
    </AppLayout>
</template>
