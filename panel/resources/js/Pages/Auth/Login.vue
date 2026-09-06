<script setup lang="ts">
import { ref, computed } from 'vue';
import { useForm, usePage } from '@inertiajs/vue3';
import {
    Lock,
    Mail,
    Eye,
    EyeOff,
    ArrowRight,
    Zap,
    ShieldCheck,
    Sliders,
    Sun,
    Moon,
    AlertCircle,
    Loader2,
    CheckCircle2,
    X,
    Copy,
    Check,
    Terminal,
} from 'lucide-vue-next';
import { useTheme } from '@/composables/useTheme';
import type { PageProps } from '@/types';

const { theme, toggleTheme } = useTheme();
const page = usePage<PageProps>();

const showPassword = ref(false);
const showForgotModal = ref(false);
const isCopied = ref(false);

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

function copyResetCommand() {
    const cmd = 'sudo -u www-data php artisan panel:reset-admin';
    navigator.clipboard.writeText(cmd);
    isCopied.value = true;
    setTimeout(() => {
        isCopied.value = false;
    }, 2000);
}
</script>

<template>
    <div class="min-h-screen flex bg-slate-50 dark:bg-[#070b14] text-slate-900 dark:text-slate-100 selection:bg-blue-500 selection:text-white font-sans relative overflow-hidden transition-colors duration-200">
        <!-- ============================================================= -->
        <!-- LEFT BRAND & FEATURE SHOWCASE (DESKTOP)                      -->
        <!-- ============================================================= -->
        <aside class="hidden lg:flex flex-col justify-between w-1/2 min-h-screen p-10 xl:p-14 text-white relative overflow-hidden bg-gradient-to-br from-[#020b1e] via-[#04153c] to-[#07245c] shrink-0">
            <!-- Atmospheric 3D Server Background Overlay -->
            <div class="absolute inset-0 pointer-events-none overflow-hidden">
                <img
                    src="/images/auth-banner.jpg"
                    alt="Server Infrastructure"
                    class="w-full h-full object-cover object-center opacity-30 mix-blend-screen scale-105"
                />
                <!-- Ambient Glow Orbs -->
                <div class="absolute -top-24 -left-24 w-96 h-96 bg-blue-500/25 blur-[120px] rounded-full"></div>
                <div class="absolute bottom-10 right-10 w-[450px] h-[450px] bg-cyan-500/20 blur-[130px] rounded-full"></div>
                <!-- Subtle Tech Grid Lines -->
                <div class="absolute inset-0 bg-[linear-gradient(to_right,#38bdf80f_1px,transparent_1px),linear-gradient(to_bottom,#38bdf80f_1px,transparent_1px)] bg-[size:3.5rem_3.5rem] [mask-image:radial-gradient(ellipse_70%_70%_at_50%_50%,#000_60%,transparent_100%)]"></div>
            </div>

            <!-- Top: Kodepreneur Brand Logo -->
            <div class="relative z-10 flex items-center">
                <img
                    src="/logo-white.png"
                    alt="Kodepreneur Panel"
                    class="h-10 xl:h-11 w-auto object-contain drop-shadow-md"
                />
            </div>

            <!-- Center: Value Proposition & Feature Highlights -->
            <div class="relative z-10 max-w-xl my-auto py-10 space-y-8">
                <!-- Headline -->
                <div>
                    <h1 class="text-4xl xl:text-[46px] font-extrabold tracking-tight text-white leading-[1.18]">
                        Power Your Ideas, <br />
                        <span class="text-transparent bg-clip-text bg-gradient-to-r from-sky-400 via-cyan-300 to-blue-400">
                            Deploy Without Limits
                        </span>
                    </h1>
                    <p class="text-sm xl:text-base text-slate-300/90 mt-4 leading-relaxed max-w-lg font-normal">
                        Kelola server VPS Anda dengan mudah, cepat, dan aman melalui Kodepreneur Panel.
                    </p>
                </div>

                <!-- 3 Key Value Features -->
                <div class="space-y-4 pt-2">
                    <!-- Feature 1: Deploy Cepat -->
                    <div class="flex items-start gap-4 p-3.5 rounded-2xl bg-white/[0.05] border border-white/[0.08] backdrop-blur-md hover:bg-white/[0.08] transition duration-200 group">
                        <div class="w-11 h-11 rounded-xl bg-blue-500/15 border border-blue-400/20 flex items-center justify-center text-cyan-400 shrink-0 group-hover:scale-105 transition-transform">
                            <Zap class="w-5 h-5 fill-cyan-400/20" />
                        </div>
                        <div>
                            <h2 class="text-sm font-semibold text-white tracking-tight">Deploy Cepat</h2>
                            <p class="text-xs text-slate-300/80 mt-0.5">Aktivasi server dalam hitungan menit</p>
                        </div>
                    </div>

                    <!-- Feature 2: Aman & Stabil -->
                    <div class="flex items-start gap-4 p-3.5 rounded-2xl bg-white/[0.05] border border-white/[0.08] backdrop-blur-md hover:bg-white/[0.08] transition duration-200 group">
                        <div class="w-11 h-11 rounded-xl bg-blue-500/15 border border-blue-400/20 flex items-center justify-center text-cyan-400 shrink-0 group-hover:scale-105 transition-transform">
                            <ShieldCheck class="w-5 h-5" />
                        </div>
                        <div>
                            <h2 class="text-sm font-semibold text-white tracking-tight">Aman & Stabil</h2>
                            <p class="text-xs text-slate-300/80 mt-0.5">Infrastruktur andal dengan uptime tinggi</p>
                        </div>
                    </div>

                    <!-- Feature 3: Mudah Dikelola -->
                    <div class="flex items-start gap-4 p-3.5 rounded-2xl bg-white/[0.05] border border-white/[0.08] backdrop-blur-md hover:bg-white/[0.08] transition duration-200 group">
                        <div class="w-11 h-11 rounded-xl bg-blue-500/15 border border-blue-400/20 flex items-center justify-center text-cyan-400 shrink-0 group-hover:scale-105 transition-transform">
                            <Sliders class="w-5 h-5" />
                        </div>
                        <div>
                            <h2 class="text-sm font-semibold text-white tracking-tight">Mudah Dikelola</h2>
                            <p class="text-xs text-slate-300/80 mt-0.5">Kelola server, domain, dan layanan langsung dari satu panel</p>
                        </div>
                    </div>
                </div>
            </div>

            <!-- Bottom: Handwritten Quote Banner -->
            <div class="relative z-10 pt-4">
                <div class="inline-block">
                    <p class="text-2xl font-serif italic text-white/95 font-medium tracking-wide">
                        Build. Deploy. Grow Together.
                    </p>
                    <div class="h-1 w-36 bg-gradient-to-r from-cyan-400 to-transparent rounded-full mt-1.5"></div>
                </div>
            </div>
        </aside>

        <!-- ============================================================= -->
        <!-- RIGHT AUTHENTICATION SECTION                                  -->
        <!-- ============================================================= -->
        <main class="w-full lg:w-1/2 min-h-screen flex flex-col justify-between p-6 sm:p-10 lg:p-12 xl:p-16 relative z-10">
            <!-- Top Controls (Language & Theme Switcher) -->
            <div class="flex items-center justify-between sm:justify-end gap-3 w-full max-w-[440px] mx-auto">
                <!-- Mobile Only Logo (when left sidebar is hidden) -->
                <div class="lg:hidden flex items-center">
                    <img
                        src="/logo.png"
                        alt="Kodepreneur Panel"
                        class="h-7 w-auto object-contain dark:hidden"
                    />
                    <img
                        src="/logo-white.png"
                        alt="Kodepreneur Panel"
                        class="h-7 w-auto object-contain hidden dark:block"
                    />
                </div>

                <div class="flex items-center gap-2.5">
                    <!-- Language Indicator Pill -->
                    <div class="flex items-center gap-1.5 px-3 py-1.5 rounded-full bg-white dark:bg-surface-900 border border-slate-200/90 dark:border-surface-800 text-xs font-medium text-slate-700 dark:text-surface-200 shadow-sm">
                        <span class="text-sm leading-none">🇮🇩</span>
                        <span>Bahasa Indonesia</span>
                    </div>

                    <!-- Theme Toggle -->
                    <button
                        @click="toggleTheme"
                        type="button"
                        class="p-2 rounded-full bg-white dark:bg-surface-900 border border-slate-200/90 dark:border-surface-800 text-slate-600 dark:text-surface-300 hover:text-slate-900 dark:hover:text-white transition shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500/30"
                        :title="theme === 'dark' ? 'Switch to Light Mode' : 'Switch to Dark Mode'"
                        aria-label="Toggle Theme"
                    >
                        <Sun v-if="theme === 'dark'" class="w-4 h-4 text-amber-400" />
                        <Moon v-else class="w-4 h-4 text-blue-600" />
                    </button>
                </div>
            </div>

            <!-- Login Card Container -->
            <div class="w-full max-w-[440px] mx-auto my-auto py-8">
                <div class="rounded-3xl bg-white dark:bg-surface-900/90 border border-slate-200/80 dark:border-surface-800 shadow-2xl shadow-slate-900/5 dark:shadow-black/50 p-8 sm:p-10 backdrop-blur-xl transition-all">
                    <!-- Center Brand Logo -->
                    <div class="flex justify-center mb-6">
                        <img
                            src="/logo.png"
                            alt="Kodepreneur Panel"
                            class="h-10 sm:h-11 w-auto object-contain dark:hidden"
                        />
                        <img
                            src="/logo-white.png"
                            alt="Kodepreneur Panel"
                            class="h-10 sm:h-11 w-auto object-contain hidden dark:block"
                        />
                    </div>

                    <!-- Card Header Title & Subtitle -->
                    <div class="text-center mb-7">
                        <h2 class="text-2xl font-bold tracking-tight text-slate-900 dark:text-white">
                            Masuk ke Panel
                        </h2>
                        <p class="text-xs text-slate-500 dark:text-surface-400 mt-1.5">
                            Kelola VPS dan layanan Anda dengan mudah.
                        </p>
                    </div>

                    <!-- Flash Notifications -->
                    <div v-if="flash?.error" class="mb-5 p-3 rounded-xl bg-rose-50 dark:bg-rose-500/10 border border-rose-200 dark:border-rose-500/20 flex items-start gap-2.5 text-xs text-rose-700 dark:text-rose-300">
                        <AlertCircle class="w-4 h-4 text-rose-500 shrink-0 mt-0.5" />
                        <span>{{ flash.error }}</span>
                    </div>
                    <div v-if="flash?.success" class="mb-5 p-3 rounded-xl bg-emerald-50 dark:bg-emerald-500/10 border border-emerald-200 dark:border-emerald-500/20 flex items-start gap-2.5 text-xs text-emerald-700 dark:text-emerald-300">
                        <CheckCircle2 class="w-4 h-4 text-emerald-500 shrink-0 mt-0.5" />
                        <span>{{ flash.success }}</span>
                    </div>

                    <!-- Form -->
                    <form @submit.prevent="submit" class="space-y-4">
                        <!-- Email Field -->
                        <div>
                            <label for="email" class="block text-xs font-medium text-slate-700 dark:text-surface-200 mb-1.5">
                                Email
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
                                    placeholder="name@contoh.com"
                                    :class="[
                                        'w-full pl-10 pr-4 py-2.5 rounded-xl text-xs bg-white dark:bg-surface-950/60 border text-slate-900 dark:text-white placeholder-slate-400 dark:placeholder-surface-500 focus:outline-none focus:ring-2 focus:ring-blue-500/30 transition duration-150',
                                        form.errors.email
                                            ? 'border-rose-300 dark:border-rose-500/50 focus:border-rose-500 focus:ring-rose-500/20'
                                            : 'border-slate-200 dark:border-surface-800 focus:border-blue-500'
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
                                    placeholder="Masukkan password Anda"
                                    :class="[
                                        'w-full pl-10 pr-10 py-2.5 rounded-xl text-xs bg-white dark:bg-surface-950/60 border text-slate-900 dark:text-white placeholder-slate-400 dark:placeholder-surface-500 focus:outline-none focus:ring-2 focus:ring-blue-500/30 transition duration-150',
                                        form.errors.password
                                            ? 'border-rose-300 dark:border-rose-500/50 focus:border-rose-500 focus:ring-rose-500/20'
                                            : 'border-slate-200 dark:border-surface-800 focus:border-blue-500'
                                    ]"
                                />
                                <button
                                    type="button"
                                    @click="showPassword = !showPassword"
                                    class="absolute inset-y-0 right-0 pr-3.5 flex items-center text-slate-400 hover:text-slate-600 dark:text-surface-500 dark:hover:text-surface-300 focus:outline-none transition"
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

                        <!-- Remember Me & Forgot Password Link -->
                        <div class="flex items-center justify-between pt-1">
                            <label class="flex items-center gap-2 cursor-pointer select-none group">
                                <input
                                    v-model="form.remember"
                                    type="checkbox"
                                    class="rounded border-slate-300 dark:border-surface-700 bg-white dark:bg-surface-950 text-blue-600 focus:ring-blue-500/30 focus:ring-offset-0 w-4 h-4 cursor-pointer"
                                />
                                <span class="text-xs text-slate-600 dark:text-surface-300 group-hover:text-slate-900 dark:group-hover:text-white transition">
                                    Ingat saya
                                </span>
                            </label>

                            <button
                                type="button"
                                @click="showForgotModal = true"
                                class="text-xs font-medium text-blue-600 hover:text-blue-700 dark:text-blue-400 dark:hover:text-blue-300 hover:underline transition"
                            >
                                Lupa password?
                            </button>
                        </div>

                        <!-- Submit Button -->
                        <button
                            type="submit"
                            :disabled="form.processing"
                            class="w-full mt-3 py-3 px-4 rounded-xl bg-gradient-to-r from-[#0062ff] to-[#0095ff] hover:from-[#0052e0] hover:to-[#0080e6] text-white text-xs font-semibold shadow-md shadow-blue-500/25 active:scale-[0.99] flex items-center justify-center gap-2 transition duration-150 disabled:opacity-60 disabled:cursor-not-allowed cursor-pointer group"
                        >
                            <Loader2 v-if="form.processing" class="w-4 h-4 animate-spin" />
                            <span v-if="form.processing">Memproses...</span>
                            <template v-else>
                                <span>Masuk</span>
                                <ArrowRight class="w-4 h-4 group-hover:translate-x-0.5 transition-transform" />
                            </template>
                        </button>
                    </form>

                    <!-- Support Help Footer -->
                    <div class="mt-7 text-center text-xs text-slate-500 dark:text-surface-400">
                        Belum punya akun?
                        <a
                            href="mailto:support@kodepreneur.com"
                            class="text-blue-600 hover:text-blue-700 dark:text-blue-400 dark:hover:text-blue-300 font-semibold hover:underline transition ml-1"
                        >
                            Hubungi Tim Support
                        </a>
                    </div>
                </div>
            </div>

            <!-- Page Bottom Copyright -->
            <footer class="text-center text-xs text-slate-400 dark:text-surface-500 py-2">
                © {{ new Date().getFullYear() }} Kodepreneur Panel. All rights reserved.
            </footer>
        </main>

        <!-- ============================================================= -->
        <!-- FORGOT PASSWORD MODAL (CLI RECOVERY HELPER)                   -->
        <!-- ============================================================= -->
        <div
            v-if="showForgotModal"
            class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/60 dark:bg-black/80 backdrop-blur-sm"
        >
            <div class="w-full max-w-md rounded-2xl bg-white dark:bg-surface-900 border border-slate-200 dark:border-surface-800 shadow-2xl p-6 sm:p-7 relative animate-in fade-in zoom-in-95 duration-150">
                <!-- Close Button -->
                <button
                    @click="showForgotModal = false"
                    class="absolute top-4 right-4 p-1.5 rounded-lg text-slate-400 hover:text-slate-600 dark:hover:text-white hover:bg-slate-100 dark:hover:bg-surface-800 transition"
                >
                    <X class="w-4 h-4" />
                </button>

                <div class="flex items-center gap-3 mb-4">
                    <div class="w-10 h-10 rounded-xl bg-blue-50 dark:bg-blue-500/10 border border-blue-200 dark:border-blue-500/20 flex items-center justify-center text-blue-600 dark:text-blue-400">
                        <Terminal class="w-5 h-5" />
                    </div>
                    <div>
                        <h3 class="text-sm font-bold text-slate-900 dark:text-white">
                            Reset Password Administrator
                        </h3>
                        <p class="text-xs text-slate-500 dark:text-surface-400">
                            Pemulihan akun via Server Terminal (SSH)
                        </p>
                    </div>
                </div>

                <p class="text-xs text-slate-600 dark:text-surface-300 leading-relaxed mb-3">
                    Kodepreneur Panel adalah sistem kontrol server mandiri. Untuk keamanan infrastruktur Anda, reset akun administrator dilakukan langsung melalui terminal server:
                </p>

                <!-- Command Box -->
                <div class="relative p-3 rounded-xl bg-slate-900 text-slate-100 font-mono text-xs border border-slate-800 mb-4 flex items-center justify-between">
                    <span class="text-emerald-400 text-[11px] overflow-x-auto pr-2">
                        sudo -u www-data php artisan panel:reset-admin
                    </span>
                    <button
                        type="button"
                        @click="copyResetCommand"
                        class="px-2 py-1 rounded bg-slate-800 hover:bg-slate-700 text-[11px] text-slate-200 flex items-center gap-1 shrink-0 transition"
                    >
                        <Check v-if="isCopied" class="w-3 h-3 text-emerald-400" />
                        <Copy v-else class="w-3 h-3" />
                        <span>{{ isCopied ? 'Tersalin' : 'Salin' }}</span>
                    </button>
                </div>

                <div class="p-3 rounded-xl bg-blue-50 dark:bg-blue-500/10 border border-blue-200/80 dark:border-blue-500/20 text-[11px] text-blue-800 dark:text-blue-300">
                    💡 Perintah ini akan menampilkan daftar akun dan memandu Anda membuat password baru secara interaktif.
                </div>

                <div class="mt-5 flex justify-end">
                    <button
                        type="button"
                        @click="showForgotModal = false"
                        class="px-4 py-2 rounded-xl bg-slate-100 hover:bg-slate-200 dark:bg-surface-800 dark:hover:bg-surface-700 text-xs font-semibold text-slate-700 dark:text-surface-200 transition"
                    >
                        Tutup
                    </button>
                </div>
            </div>
        </div>
    </div>
</template>
