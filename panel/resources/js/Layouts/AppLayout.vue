<script setup lang="ts">
import { ref, computed } from 'vue';
import { Link, usePage, router } from '@inertiajs/vue3';
import {
    LayoutDashboard,
    Globe,
    Database,
    GitBranch,
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
} from 'lucide-vue-next';
import type { PageProps } from '@/types';

defineProps<{
    title?: string;
}>();

const page = usePage<PageProps>();
const user = computed(() => page.props.auth.user);
const flash = computed(() => page.props.flash);

const isMobileMenuOpen = ref(false);
const isUserDropdownOpen = ref(false);

const navigation = [
    { name: 'Dashboard', href: '/dashboard', icon: LayoutDashboard, current: routeMatches('/dashboard') },
    { name: 'Websites', href: '/websites', icon: Globe, current: routeMatches('/websites') },
    { name: 'Databases', href: '/databases', icon: Database, current: routeMatches('/databases') },
    { name: 'Deployments', href: '/deployments', icon: GitBranch, current: routeMatches('/deployments') },
    { name: 'Cron Jobs', href: '/cron', icon: Clock, current: routeMatches('/cron') },
    { name: 'Processes', href: '/processes', icon: Sliders, current: routeMatches('/processes') },
    { name: 'File Manager', href: '/files', icon: Folder, current: routeMatches('/files') },
    { name: 'Firewall', href: '/firewall', icon: Shield, current: routeMatches('/firewall') },
    { name: 'Services', href: '/services', icon: Cpu, current: routeMatches('/services') },
    { name: 'Activity Log', href: '/activity', icon: Activity, current: routeMatches('/activity') },
    { name: 'Settings', href: '/settings', icon: Settings, current: routeMatches('/settings') },
];

function routeMatches(path: string): boolean {
    return window.location.pathname.startsWith(path);
}

function handleLogout() {
    router.post('/logout');
}
</script>

<template>
    <div class="min-h-screen bg-surface-950 text-surface-100 flex flex-col md:flex-row antialiased selection:bg-brand-500 selection:text-white">
        <!-- Sidebar Navigation (Desktop) -->
        <aside class="hidden md:flex flex-col w-64 border-r border-surface-800 bg-surface-900/60 backdrop-blur-xl shrink-0">
            <!-- Brand Logo -->
            <div class="h-16 flex items-center px-6 gap-3 border-b border-surface-800/80">
                <div class="w-8 h-8 rounded-lg bg-gradient-to-tr from-brand-600 to-brand-400 flex items-center justify-center shadow-lg shadow-brand-500/20 ring-1 ring-white/20">
                    <Server class="w-4 h-4 text-white" />
                </div>
                <div>
                    <h1 class="text-sm font-semibold tracking-tight text-white flex items-center gap-1.5">
                        Kodepreneur <span class="text-brand-400 text-xs font-mono font-normal">v1.0</span>
                    </h1>
                    <p class="text-[11px] text-surface-400">Server Management</p>
                </div>
            </div>

            <!-- Navigation Links -->
            <nav class="flex-1 px-3 py-4 space-y-1 overflow-y-auto">
                <Link
                    v-for="item in navigation"
                    :key="item.name"
                    :href="item.href"
                    :class="[
                        item.current
                            ? 'bg-brand-600/15 text-brand-400 font-medium border-brand-500/30'
                            : 'text-surface-400 hover:text-surface-100 hover:bg-surface-800/50 border-transparent',
                        'group flex items-center gap-3 px-3 py-2 text-xs rounded-lg transition-all duration-150 border'
                    ]"
                >
                    <component
                        :is="item.icon"
                        :class="[
                            item.current ? 'text-brand-400' : 'text-surface-500 group-hover:text-surface-300',
                            'w-4 h-4 shrink-0 transition-colors'
                        ]"
                    />
                    <span>{{ item.name }}</span>
                </Link>
            </nav>

            <!-- Agent Status Footer -->
            <div class="p-3 border-t border-surface-800/80 bg-surface-950/40">
                <div class="flex items-center justify-between px-3 py-2 rounded-lg bg-surface-900/80 border border-surface-800/60">
                    <div class="flex items-center gap-2">
                        <div class="w-2 h-2 rounded-full bg-emerald-500 animate-pulse"></div>
                        <span class="text-[11px] font-mono text-surface-300">Agent Active</span>
                    </div>
                    <span class="text-[10px] text-surface-500 font-mono">root:0660</span>
                </div>
            </div>
        </aside>

        <!-- Main Content Area -->
        <div class="flex-1 flex flex-col min-w-0">
            <!-- Top Header Bar -->
            <header class="h-16 border-b border-surface-800/80 bg-surface-900/40 backdrop-blur-md px-4 sm:px-6 flex items-center justify-between sticky top-0 z-30">
                <!-- Mobile Menu Button -->
                <button
                    @click="isMobileMenuOpen = !isMobileMenuOpen"
                    class="md:hidden p-2 rounded-lg text-surface-400 hover:text-white hover:bg-surface-800"
                >
                    <Menu v-if="!isMobileMenuOpen" class="w-5 h-5" />
                    <X v-else class="w-5 h-5" />
                </button>

                <!-- Breadcrumb / Page Title -->
                <div class="flex items-center gap-2 text-xs text-surface-400">
                    <span class="hidden sm:inline">Kodepreneur Panel</span>
                    <span class="hidden sm:inline text-surface-600">/</span>
                    <span class="text-surface-100 font-medium">{{ title || 'Overview' }}</span>
                </div>

                <!-- User Controls -->
                <div class="flex items-center gap-3">
                    <div class="relative">
                        <button
                            @click="isUserDropdownOpen = !isUserDropdownOpen"
                            class="flex items-center gap-2.5 px-3 py-1.5 rounded-lg bg-surface-800/40 border border-surface-700/50 hover:bg-surface-800 transition text-xs"
                        >
                            <div class="w-6 h-6 rounded-full bg-brand-600/30 border border-brand-500/40 text-brand-400 flex items-center justify-center font-medium text-[11px]">
                                {{ user?.name ? user.name.charAt(0).toUpperCase() : 'A' }}
                            </div>
                            <span class="text-surface-200 font-medium hidden sm:inline">{{ user?.name || 'Administrator' }}</span>
                            <ChevronDown class="w-3.5 h-3.5 text-surface-400" />
                        </button>

                        <!-- Dropdown Menu -->
                        <div
                            v-if="isUserDropdownOpen"
                            class="absolute right-0 mt-2 w-48 rounded-xl bg-surface-900 border border-surface-800 shadow-2xl py-1 text-xs z-50 animate-in fade-in zoom-in-95 duration-100"
                        >
                            <div class="px-3 py-2 border-b border-surface-800/80">
                                <p class="text-white font-medium truncate">{{ user?.name }}</p>
                                <p class="text-surface-400 text-[11px] truncate">{{ user?.email }}</p>
                            </div>
                            <Link
                                href="/settings"
                                class="flex items-center gap-2 px-3 py-2 text-surface-300 hover:text-white hover:bg-surface-800/70"
                                @click="isUserDropdownOpen = false"
                            >
                                <Settings class="w-3.5 h-3.5" />
                                <span>Settings</span>
                            </Link>
                            <button
                                @click="handleLogout"
                                class="w-full flex items-center gap-2 px-3 py-2 text-rose-400 hover:bg-rose-500/10 transition text-left"
                            >
                                <LogOut class="w-3.5 h-3.5" />
                                <span>Sign out</span>
                            </button>
                        </div>
                    </div>
                </div>
            </header>

            <!-- Flash Messages -->
            <div v-if="flash?.success" class="m-4 sm:mx-6 mb-0 p-3.5 rounded-xl bg-emerald-500/10 border border-emerald-500/30 text-emerald-400 flex items-center gap-2.5 text-xs">
                <CheckCircle2 class="w-4 h-4 shrink-0" />
                <span>{{ flash.success }}</span>
            </div>
            <div v-if="flash?.error" class="m-4 sm:mx-6 mb-0 p-3.5 rounded-xl bg-rose-500/10 border border-rose-500/30 text-rose-400 flex items-center gap-2.5 text-xs">
                <AlertCircle class="w-4 h-4 shrink-0" />
                <span>{{ flash.error }}</span>
            </div>

            <!-- Page Body -->
            <main class="flex-1 p-4 sm:p-6 md:p-8 overflow-y-auto">
                <slot />
            </main>
        </div>
    </div>
</template>
