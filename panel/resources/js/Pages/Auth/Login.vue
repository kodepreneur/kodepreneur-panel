<script setup lang="ts">
import { ref, computed } from 'vue';
import { useForm, usePage } from '@inertiajs/vue3';
import {
    Lock,
    Mail,
    Eye,
    EyeOff,
    ArrowRight,
    Layers,
    Sun,
    Moon,
    ShieldCheck,
    AlertCircle,
    Loader2,
    CheckCircle2,
} from 'lucide-vue-next';
import { useTheme } from '@/composables/useTheme';
import type { PageProps } from '@/types';

const { theme, toggleTheme } = useTheme();
const page = usePage<PageProps>();

const showPassword = ref(false);

const form = useForm({
    email: '',
    password: '',
    remember: false,
});

const flash = computed(() => page.props.flash);

function submit() {
    form.post('/login', {
        onFinish: () => form.reset('password'),
    });
}
</script>

<template>
    <div class="min-h-screen bg-slate-50 dark:bg-[#070b14] text-slate-900 dark:text-slate-100 flex flex-col justify-between p-4 sm:p-6 selection:bg-brand-500 selection:text-white relative overflow-hidden transition-colors duration-200">
        <!-- Ambient High-Tech Backdrop -->
        <div class="absolute inset-0 pointer-events-none overflow-hidden">
            <!-- Radial gradient mesh -->
            <div class="absolute -top-40 left-1/2 -translate-x-1/2 w-[700px] h-[500px] bg-gradient-to-b from-brand-500/15 via-indigo-600/10 to-transparent blur-[140px] rounded-full"></div>
            <div class="absolute -bottom-40 right-10 w-[500px] h-[400px] bg-gradient-to-t from-indigo-500/10 via-purple-600/5 to-transparent blur-[130px] rounded-full"></div>

            <!-- Subtle background grid pattern -->
            <div class="absolute inset-0 bg-[linear-gradient(to_right,#e2e8f015_1px,transparent_1px),linear-gradient(to_bottom,#e2e8f015_1px,transparent_1px)] dark:bg-[linear-gradient(to_right,#1e293b20_1px,transparent_1px),linear-gradient(to_bottom,#1e293b20_1px,transparent_1px)] bg-[size:4rem_4rem] [mask-image:radial-gradient(ellipse_60%_50%_at_50%_40%,#000_70%,transparent_100%)]"></div>
        </div>

        <!-- Top Navigation Bar -->
        <header class="w-full max-w-5xl mx-auto flex items-center justify-between relative z-20 py-2">
            <!-- Brand Badge -->
            <div class="flex items-center gap-2.5">
                <div class="w-9 h-9 rounded-xl bg-gradient-to-tr from-brand-600 via-indigo-500 to-indigo-400 flex items-center justify-center shadow-lg shadow-brand-500/20 ring-1 ring-white/20 dark:ring-white/10">
                    <Layers class="w-4 h-4 text-white" />
                </div>
                <div>
                    <span class="text-sm font-semibold tracking-tight text-slate-900 dark:text-white flex items-center gap-1.5">
                        Kodepreneur
                        <span class="px-1.5 py-0.5 rounded text-[10px] font-mono font-medium bg-brand-50 text-brand-600 border border-brand-200 dark:bg-brand-500/10 dark:text-brand-400 dark:border-brand-500/20">
                            v1.0
                        </span>
                    </span>
                </div>
            </div>

            <!-- Header Controls: Server Status & Theme Switcher -->
            <div class="flex items-center gap-3">
                <!-- Status Pill -->
                <div class="hidden sm:flex items-center gap-2 px-2.5 py-1 rounded-full bg-white/80 dark:bg-surface-900/80 border border-slate-200/80 dark:border-surface-800 backdrop-blur-md text-[11px] font-medium text-slate-600 dark:text-surface-300 shadow-sm">
                    <span class="relative flex h-2 w-2">
                        <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-400 opacity-75"></span>
                        <span class="relative inline-flex rounded-full h-2 w-2 bg-emerald-500"></span>
                    </span>
                    <span>System Operational</span>
                </div>

                <!-- Theme Toggle Button -->
                <button
                    @click="toggleTheme"
                    type="button"
                    class="p-2 rounded-xl bg-white/90 dark:bg-surface-900/90 border border-slate-200/80 dark:border-surface-800 text-slate-600 dark:text-surface-300 hover:text-slate-900 dark:hover:text-white hover:border-slate-300 dark:hover:border-surface-700 transition shadow-sm backdrop-blur-sm focus:outline-none focus:ring-2 focus:ring-brand-500/30"
                    :title="theme === 'dark' ? 'Switch to Light Mode' : 'Switch to Dark Mode'"
                    aria-label="Toggle color theme"
                >
                    <Sun v-if="theme === 'dark'" class="w-4 h-4 text-amber-400" />
                    <Moon v-else class="w-4 h-4 text-indigo-600" />
                </button>
            </div>
        </header>

        <!-- Main Authentication Section -->
        <main class="w-full max-w-[420px] mx-auto my-auto relative z-10 py-6">
            <!-- Glassmorphic Login Card -->
            <div class="relative rounded-2xl bg-white/90 dark:bg-surface-900/85 border border-slate-200/90 dark:border-surface-800 shadow-2xl shadow-slate-900/5 dark:shadow-black/60 backdrop-blur-2xl p-6 sm:p-8 transition-all">
                <!-- Glowing Top Accent Line -->
                <div class="absolute top-0 inset-x-0 h-[2px] bg-gradient-to-r from-transparent via-brand-500 to-transparent"></div>

                <!-- Card Header -->
                <div class="flex flex-col items-center text-center mb-6">
                    <div class="w-12 h-12 rounded-2xl bg-brand-50 dark:bg-brand-500/10 border border-brand-200/80 dark:border-brand-500/20 flex items-center justify-center mb-3.5 shadow-inner">
                        <Lock class="w-5 h-5 text-brand-600 dark:text-brand-400" />
                    </div>
                    <h2 class="text-lg font-bold tracking-tight text-slate-900 dark:text-white">
                        Welcome Back
                    </h2>
                    <p class="text-xs text-slate-500 dark:text-surface-400 mt-1 max-w-xs">
                        Sign in with your administrator credentials to access your server infrastructure
                    </p>
                </div>

                <!-- Flash Notifications -->
                <div v-if="flash?.error" class="mb-4 p-3 rounded-xl bg-rose-50 dark:bg-rose-500/10 border border-rose-200 dark:border-rose-500/20 flex items-start gap-2.5 text-xs text-rose-700 dark:text-rose-300">
                    <AlertCircle class="w-4 h-4 text-rose-500 shrink-0 mt-0.5" />
                    <span>{{ flash.error }}</span>
                </div>
                <div v-if="flash?.success" class="mb-4 p-3 rounded-xl bg-emerald-50 dark:bg-emerald-500/10 border border-emerald-200 dark:border-emerald-500/20 flex items-start gap-2.5 text-xs text-emerald-700 dark:text-emerald-300">
                    <CheckCircle2 class="w-4 h-4 text-emerald-500 shrink-0 mt-0.5" />
                    <span>{{ flash.success }}</span>
                </div>

                <!-- Form -->
                <form @submit.prevent="submit" class="space-y-4">
                    <!-- Email Field -->
                    <div>
                        <label for="email" class="block text-xs font-medium text-slate-700 dark:text-surface-200 mb-1.5">
                            Email Address
                        </label>
                        <div class="relative">
                            <div class="absolute inset-y-0 left-0 pl-3.5 flex items-center pointer-events-none text-slate-400 dark:text-surface-500">
                                <Mail class="w-4 h-4" />
                            </div>
                            <input
                                id="email"
                                v-model="form.email"
                                type="email"
                                required
                                autofocus
                                autocomplete="email"
                                placeholder="admin@kodepreneur.com"
                                :class="[
                                    'w-full pl-10 pr-4 py-2.5 rounded-xl text-xs bg-slate-50 dark:bg-surface-950/60 border text-slate-900 dark:text-white placeholder-slate-400 dark:placeholder-surface-500 focus:outline-none focus:ring-2 focus:ring-brand-500/30 transition duration-150',
                                    form.errors.email
                                        ? 'border-rose-300 dark:border-rose-500/50 focus:border-rose-500 focus:ring-rose-500/20'
                                        : 'border-slate-200 dark:border-surface-800 focus:border-brand-500'
                                ]"
                            />
                        </div>
                        <div v-if="form.errors.email" class="flex items-center gap-1.5 text-[11px] text-rose-500 mt-1.5">
                            <AlertCircle class="w-3.5 h-3.5 shrink-0" />
                            <span>{{ form.errors.email }}</span>
                        </div>
                    </div>

                    <!-- Password Field -->
                    <div>
                        <label for="password" class="block text-xs font-medium text-slate-700 dark:text-surface-200 mb-1.5">
                            Password
                        </label>
                        <div class="relative">
                            <div class="absolute inset-y-0 left-0 pl-3.5 flex items-center pointer-events-none text-slate-400 dark:text-surface-500">
                                <Lock class="w-4 h-4" />
                            </div>
                            <input
                                id="password"
                                v-model="form.password"
                                :type="showPassword ? 'text' : 'password'"
                                required
                                autocomplete="current-password"
                                placeholder="••••••••••••"
                                :class="[
                                    'w-full pl-10 pr-10 py-2.5 rounded-xl text-xs bg-slate-50 dark:bg-surface-950/60 border text-slate-900 dark:text-white placeholder-slate-400 dark:placeholder-surface-500 focus:outline-none focus:ring-2 focus:ring-brand-500/30 transition duration-150',
                                    form.errors.password
                                        ? 'border-rose-300 dark:border-rose-500/50 focus:border-rose-500 focus:ring-rose-500/20'
                                        : 'border-slate-200 dark:border-surface-800 focus:border-brand-500'
                                ]"
                            />
                            <button
                                type="button"
                                @click="showPassword = !showPassword"
                                class="absolute inset-y-0 right-0 pr-3 flex items-center text-slate-400 hover:text-slate-600 dark:text-surface-500 dark:hover:text-surface-300 focus:outline-none transition"
                                :aria-label="showPassword ? 'Hide password' : 'Show password'"
                            >
                                <EyeOff v-if="showPassword" class="w-4 h-4" />
                                <Eye v-else class="w-4 h-4" />
                            </button>
                        </div>
                        <div v-if="form.errors.password" class="flex items-center gap-1.5 text-[11px] text-rose-500 mt-1.5">
                            <AlertCircle class="w-3.5 h-3.5 shrink-0" />
                            <span>{{ form.errors.password }}</span>
                        </div>
                    </div>

                    <!-- Remember Me & Security Note -->
                    <div class="flex items-center justify-between pt-0.5">
                        <label class="flex items-center gap-2 cursor-pointer select-none group">
                            <input
                                v-model="form.remember"
                                type="checkbox"
                                class="rounded border-slate-300 dark:border-surface-700 bg-white dark:bg-surface-950 text-brand-600 focus:ring-brand-500/30 focus:ring-offset-0 w-3.5 h-3.5 cursor-pointer"
                            />
                            <span class="text-xs text-slate-600 dark:text-surface-300 group-hover:text-slate-900 dark:group-hover:text-white transition">
                                Remember this session
                            </span>
                        </label>
                    </div>

                    <!-- Submit Button -->
                    <button
                        type="submit"
                        :disabled="form.processing"
                        class="w-full mt-2 py-2.5 px-4 rounded-xl bg-gradient-to-r from-brand-600 via-indigo-600 to-brand-500 hover:from-brand-500 hover:to-indigo-500 text-white text-xs font-semibold shadow-md shadow-brand-500/25 active:scale-[0.99] flex items-center justify-center gap-2 transition duration-150 disabled:opacity-60 disabled:cursor-not-allowed cursor-pointer group"
                    >
                        <Loader2 v-if="form.processing" class="w-4 h-4 animate-spin" />
                        <span v-if="form.processing">Authenticating...</span>
                        <template v-else>
                            <span>Sign In to Panel</span>
                            <ArrowRight class="w-3.5 h-3.5 group-hover:translate-x-0.5 transition-transform" />
                        </template>
                    </button>
                </form>

                <!-- Security Assurance Footer -->
                <div class="mt-6 pt-5 border-t border-slate-200/70 dark:border-surface-800/70 flex items-center justify-center gap-2 text-[11px] text-slate-400 dark:text-surface-500">
                    <ShieldCheck class="w-3.5 h-3.5 text-emerald-500 shrink-0" />
                    <span>256-bit Encrypted Management Session</span>
                </div>
            </div>
        </main>

        <!-- Minimalist Page Footer -->
        <footer class="w-full max-w-5xl mx-auto flex flex-col sm:flex-row items-center justify-between gap-2 text-center text-[11px] text-slate-400 dark:text-surface-500 relative z-20 py-2">
            <div>
                © {{ new Date().getFullYear() }} Kodepreneur Panel. All rights reserved.
            </div>
            <div class="flex items-center gap-4">
                <span>Enterprise Security</span>
                <span>•</span>
                <span>Automated Deployments</span>
                <span>•</span>
                <span>Zero Latency</span>
            </div>
        </footer>
    </div>
</template>
