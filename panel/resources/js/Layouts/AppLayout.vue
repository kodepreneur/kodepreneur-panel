<script setup lang="ts">
import { ref, computed } from 'vue';
import { Link, usePage, router } from '@inertiajs/vue3';
import {
    LayoutDashboard,
    Globe,
    Database,
    Rocket,
    Cpu,
    Activity,
    Settings,
    LogOut,
    Server,
    Shield,
    Terminal,
    ChevronDown,
    Menu,
    X,
    CheckCircle2,
    AlertCircle,
    Clock,
    Folder,
    Sliders,
    ArrowUpCircle,
    Sun,
    Moon,
    Bell,
    Search,
    Layers,
} from 'lucide-vue-next';
import { useTheme } from '@/composables/useTheme';
import type { PageProps } from '@/types';

defineProps<{
    title?: string;
}>();

const page = usePage<PageProps>();
const user = computed(() => page.props.auth.user);
const flash = computed(() => page.props.flash);

const { theme, toggleTheme } = useTheme();

const isMobileMenuOpen = ref(false);
const isUserDropdownOpen = ref(false);

const navigation = [
    { name: 'Dashboard', href: '/dashboard', icon: LayoutDashboard, current: routeMatches('/dashboard') },
    { name: 'Websites', href: '/websites', icon: Globe, current: routeMatches('/websites') },
    { name: 'Databases', href: '/databases', icon: Database, current: routeMatches('/databases') },
    { name: 'Deployments', href: '/deployments', icon: Rocket, current: routeMatches('/deployments') },
    { name: 'Cron Jobs', href: '/cron', icon: Clock, current: routeMatches('/cron') },
    { name: 'Processes', href: '/processes', icon: Sliders, current: routeMatches('/processes') },
    { name: 'File Manager', href: '/files', icon: Folder, current: routeMatches('/files') },
    { name: 'Firewall', href: '/firewall', icon: Shield, current: routeMatches('/firewall') },
    { name: 'Services', href: '/services', icon: Cpu, current: routeMatches('/services') },
    { name: 'Activity Log', href: '/activity', icon: Activity, current: routeMatches('/activity') },
    { name: 'Settings', href: '/settings', icon: Settings, current: routeMatches('/settings') },
    { name: 'Updates', href: '/updates', icon: ArrowUpCircle, current: routeMatches('/updates') },
];

function routeMatches(path: string): boolean {
    return window.location.pathname.startsWith(path);
}

function handleLogout() {
    router.post('/logout');
}
</script>

<template>
    <div class="min-h-screen bg-slate-50 dark:bg-[#070b14] text-slate-900 dark:text-slate-100 flex flex-col md:flex-row antialiased selection:bg-brand-500 selection:text-white transition-colors duration-150">
        <!-- Sidebar Navigation (Desktop) -->
        <aside class="hidden md:flex flex-col w-64 border-r border-slate-200/80 dark:border-surface-800/80 bg-white dark:bg-[#0a0e1a] shrink-0 transition-colors duration-150">
            <!-- Brand Logo Header -->
            <div class="h-16 flex items-center px-6 gap-3 border-b border-slate-200/80 dark:border-surface-800/80">
                <div class="w-8 h-8 rounded-xl bg-gradient-to-tr from-brand-600 via-indigo-500 to-indigo-400 flex items-center justify-center shadow-md shadow-brand-500/25 ring-1 ring-white/20">
                    <Layers class="w-4 h-4 text-white" />
                </div>
                <div>
                    <h1 class="text-sm font-semibold tracking-tight text-slate-900 dark:text-white flex items-center gap-1.5">
                        Kodepreneur <span class="px-1.5 py-0.5 rounded text-[10px] font-mono font-medium bg-brand-50 text-brand-600 border border-brand-200 dark:bg-brand-500/10 dark:text-brand-400 dark:border-brand-500/20">v1.0</span>
                    </h1>
                    <p class="text-[11px] text-slate-500 dark:text-surface-400 font-medium">Server Management</p>
                </div>
            </div>

            <!-- Navigation Links List -->
            <nav class="flex-1 px-3.5 py-4 space-y-1 overflow-y-auto">
                <Link
                    v-for="item in navigation"
                    :key="item.name"
                    :href="item.href"
                    :class="[
                        item.current
                            ? 'bg-brand-50 text-brand-600 font-semibold dark:bg-brand-600 dark:text-white dark:shadow-lg dark:shadow-brand-600/30'
                            : 'text-slate-600 hover:text-slate-900 hover:bg-slate-100/80 dark:text-surface-400 dark:hover:text-white dark:hover:bg-surface-800/50 font-medium',
                        'group flex items-center gap-3 px-3 py-2 text-xs rounded-xl transition-all duration-150'
                    ]"
                >
                    <component
                        :is="item.icon"
                        :class="[
                            item.current
                                ? 'text-brand-600 dark:text-white'
                                : 'text-slate-400 group-hover:text-slate-600 dark:text-surface-500 dark:group-hover:text-surface-300',
                            'w-4 h-4 shrink-0 transition-colors'
                        ]"
                    />
                    <span>{{ item.name }}</span>
                </Link>
            </nav>

            <!-- Agent Status Footer Card -->
            <div class="p-3 border-t border-slate-200/80 dark:border-surface-800/80 bg-slate-50/50 dark:bg-surface-950/40">
                <div class="p-3 rounded-xl bg-white dark:bg-surface-900/80 border border-slate-200/80 dark:border-surface-800/60 shadow-sm space-y-1.5">
                    <div class="flex items-center justify-between">
                        <div class="flex items-center gap-2">
                            <span class="relative flex h-2 w-2">
                                <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-400 opacity-75"></span>
                                <span class="relative inline-flex rounded-full h-2 w-2 bg-emerald-500"></span>
                            </span>
                            <span class="text-[11px] font-medium text-emerald-600 dark:text-emerald-400">Agent Active</span>
                        </div>
                        <span class="text-[10px] text-slate-500 dark:text-surface-400 font-mono px-1.5 py-0.5 rounded bg-slate-100 dark:bg-surface-800">root:0660</span>
                    </div>
                    <p class="text-[10px] text-slate-500 dark:text-surface-400 font-mono truncate">root@kodepreneur</p>
                </div>
            </div>
        </aside>

        <!-- Main Content Area -->
        <div class="flex-1 flex flex-col min-w-0">
            <!-- Top Header Bar -->
            <header class="h-16 border-b border-slate-200/80 dark:border-surface-800/80 bg-white/80 dark:bg-[#0a0e1a]/80 backdrop-blur-md px-4 sm:px-6 flex items-center justify-between sticky top-0 z-30 transition-colors duration-150">
                <!-- Left: Mobile Menu Button & Breadcrumb -->
                <div class="flex items-center gap-3">
                    <button
                        @click="isMobileMenuOpen = !isMobileMenuOpen"
                        class="md:hidden p-2 rounded-xl text-slate-500 hover:text-slate-800 hover:bg-slate-100 dark:text-surface-400 dark:hover:text-white dark:hover:bg-surface-800 transition"
                    >
                        <Menu v-if="!isMobileMenuOpen" class="w-5 h-5" />
                        <X v-else class="w-5 h-5" />
                    </button>

                    <div class="flex items-center gap-2 text-xs text-slate-500 dark:text-surface-400 font-medium">
                        <span class="hidden sm:inline">Kodepreneur Panel</span>
                        <span class="hidden sm:inline text-slate-300 dark:text-surface-600">/</span>
                        <span class="text-slate-900 dark:text-white font-semibold">{{ title || 'Server Dashboard' }}</span>
                    </div>
                </div>

                <!-- Right Controls: Theme Toggle, Notifications, User Menu -->
                <div class="flex items-center gap-2.5">
                    <!-- Light / Dark Theme Toggle Button -->
                    <button
                        @click="toggleTheme"
                        type="button"
                        class="p-2 rounded-xl text-slate-500 hover:text-slate-800 hover:bg-slate-100 dark:text-surface-400 dark:hover:text-white dark:hover:bg-surface-800 border border-transparent hover:border-slate-200 dark:hover:border-surface-700 transition"
                        :title="theme === 'dark' ? 'Switch to Light Mode' : 'Switch to Dark Mode'"
                        aria-label="Toggle Theme"
                    >
                        <Sun v-if="theme === 'dark'" class="w-4 h-4 text-amber-300" />
                        <Moon v-else class="w-4 h-4 text-indigo-600" />
                    </button>

                    <!-- User Profile Dropdown -->
                    <div class="relative">
                        <button
                            @click="isUserDropdownOpen = !isUserDropdownOpen"
                            class="flex items-center gap-2.5 px-3 py-1.5 rounded-xl bg-slate-100/80 dark:bg-surface-800/40 border border-slate-200/80 dark:border-surface-700/50 hover:bg-slate-200/60 dark:hover:bg-surface-800 transition text-xs"
                        >
                            <div class="w-6 h-6 rounded-full bg-brand-600 text-white flex items-center justify-center font-medium text-[11px] shadow-sm">
                                {{ user?.name ? user.name.charAt(0).toUpperCase() : 'A' }}
                            </div>
                            <span class="text-slate-800 dark:text-surface-200 font-semibold hidden sm:inline">{{ user?.name || 'Administrator' }}</span>
                            <ChevronDown class="w-3.5 h-3.5 text-slate-400 dark:text-surface-400" />
                        </button>

                        <!-- Dropdown Menu -->
                        <div
                            v-if="isUserDropdownOpen"
                            class="absolute right-0 mt-2 w-52 rounded-2xl bg-white dark:bg-surface-900 border border-slate-200/90 dark:border-surface-800 shadow-2xl py-1.5 text-xs z-50 animate-in fade-in zoom-in-95 duration-100"
                        >
                            <div class="px-3.5 py-2.5 border-b border-slate-100 dark:border-surface-800/80">
                                <p class="text-slate-900 dark:text-white font-semibold truncate">{{ user?.name || 'Administrator' }}</p>
                                <p class="text-slate-500 dark:text-surface-400 text-[11px] truncate">{{ user?.email || 'admin@kodepreneur.com' }}</p>
                            </div>
                            <Link
                                href="/settings"
                                class="flex items-center gap-2.5 px-3.5 py-2 text-slate-700 dark:text-surface-300 hover:text-slate-900 dark:hover:text-white hover:bg-slate-100 dark:hover:bg-surface-800/70 transition"
                                @click="isUserDropdownOpen = false"
                            >
                                <Settings class="w-3.5 h-3.5 text-slate-400 dark:text-surface-400" />
                                <span>Panel Settings</span>
                            </Link>
                            <Link
                                href="/updates"
                                class="flex items-center gap-2.5 px-3.5 py-2 text-slate-700 dark:text-surface-300 hover:text-slate-900 dark:hover:text-white hover:bg-slate-100 dark:hover:bg-surface-800/70 transition"
                                @click="isUserDropdownOpen = false"
                            >
                                <ArrowUpCircle class="w-3.5 h-3.5 text-brand-500" />
                                <span>System Updates</span>
                            </Link>
                            <button
                                @click="handleLogout"
                                class="w-full flex items-center gap-2.5 px-3.5 py-2 text-rose-500 hover:bg-rose-50 dark:hover:bg-rose-500/10 transition text-left"
                            >
                                <LogOut class="w-3.5 h-3.5" />
                                <span>Sign out</span>
                            </button>
                        </div>
                    </div>
                </div>
            </header>

            <!-- Mobile Navigation Drawer -->
            <div v-if="isMobileMenuOpen" class="md:hidden border-b border-slate-200 dark:border-surface-800 bg-white/95 dark:bg-surface-900/95 backdrop-blur-xl px-4 py-3 space-y-1">
                <Link
                    v-for="item in navigation"
                    :key="item.name"
                    :href="item.href"
                    @click="isMobileMenuOpen = false"
                    :class="[
                        item.current
                            ? 'bg-brand-50 text-brand-600 font-semibold dark:bg-brand-600 dark:text-white'
                            : 'text-slate-600 hover:text-slate-900 hover:bg-slate-100 dark:text-surface-400 dark:hover:text-white dark:hover:bg-surface-800/50',
                        'group flex items-center gap-3 px-3 py-2 text-xs rounded-xl transition-all duration-150'
                    ]"
                >
                    <component
                        :is="item.icon"
                        :class="[
                            item.current ? 'text-brand-600 dark:text-white' : 'text-slate-400 dark:text-surface-500',
                            'w-4 h-4 shrink-0 transition-colors'
                        ]"
                    />
                    <span>{{ item.name }}</span>
                </Link>
            </div>

            <!-- Flash Alerts -->
            <div v-if="flash?.success" class="m-4 sm:mx-6 mb-0 p-3.5 rounded-2xl bg-emerald-50 dark:bg-emerald-500/10 border border-emerald-200 dark:border-emerald-500/30 text-emerald-700 dark:text-emerald-400 flex items-center gap-2.5 text-xs shadow-sm">
                <CheckCircle2 class="w-4 h-4 shrink-0 text-emerald-500" />
                <span>{{ flash.success }}</span>
            </div>
            <div v-if="flash?.error" class="m-4 sm:mx-6 mb-0 p-3.5 rounded-2xl bg-rose-50 dark:bg-rose-500/10 border border-rose-200 dark:border-rose-500/30 text-rose-700 dark:text-rose-400 flex items-center gap-2.5 text-xs shadow-sm">
                <AlertCircle class="w-4 h-4 shrink-0 text-rose-500" />
                <span>{{ flash.error }}</span>
            </div>

            <!-- Page Main Body -->
            <main class="flex-1 p-4 sm:p-6 md:p-8 overflow-y-auto">
                <slot />
            </main>
        </div>
    </div>
</template>
