<script setup lang="ts">
import { ref } from 'vue';
import { useForm } from '@inertiajs/vue3';
import { Server, Lock, Mail, Eye, EyeOff, ArrowRight } from 'lucide-vue-next';

const showPassword = ref(false);

const form = useForm({
    email: 'admin@kodepreneur.com',
    password: 'password',
    remember: true,
});

function submit() {
    form.post('/login', {
        onFinish: () => form.reset('password'),
    });
}
</script>

<template>
    <div class="min-h-screen bg-surface-950 flex flex-col items-center justify-center p-4 selection:bg-brand-500 selection:text-white">
        <!-- Glow / Backdrop Accents -->
        <div class="absolute inset-0 overflow-hidden pointer-events-none flex items-center justify-center">
            <div class="w-[500px] h-[500px] bg-brand-600/10 blur-[120px] rounded-full"></div>
        </div>

        <div class="w-full max-w-md relative z-10">
            <!-- Brand Badge -->
            <div class="flex flex-col items-center text-center mb-8">
                <div class="w-12 h-12 rounded-2xl bg-gradient-to-tr from-brand-600 to-brand-400 flex items-center justify-center shadow-xl shadow-brand-500/25 ring-1 ring-white/20 mb-4">
                    <Server class="w-6 h-6 text-white" />
                </div>
                <h1 class="text-xl font-bold text-white tracking-tight">Kodepreneur Panel</h1>
                <p class="text-xs text-surface-400 mt-1">Simple server management without unnecessary complexity</p>
            </div>

            <!-- Login Card -->
            <div class="rounded-2xl bg-surface-900/70 border border-surface-800/80 shadow-2xl backdrop-blur-xl p-6 sm:p-8">
                <form @submit.prevent="submit" class="space-y-4">
                    <!-- Email Field -->
                    <div>
                        <label class="block text-xs font-medium text-surface-300 mb-1.5">Email Address</label>
                        <div class="relative">
                            <Mail class="w-4 h-4 text-surface-500 absolute left-3.5 top-1/2 -translate-y-1/2 pointer-events-none" />
                            <input
                                v-model="form.email"
                                type="email"
                                required
                                autocomplete="email"
                                placeholder="admin@kodepreneur.com"
                                class="w-full pl-10 pr-4 py-2.5 rounded-xl bg-surface-950/60 border border-surface-800 text-xs text-white placeholder-surface-500 focus:outline-none focus:ring-2 focus:ring-brand-500/40 focus:border-brand-500 transition"
                            />
                        </div>
                        <p v-if="form.errors.email" class="text-[11px] text-rose-400 mt-1">{{ form.errors.email }}</p>
                    </div>

                    <!-- Password Field -->
                    <div>
                        <label class="block text-xs font-medium text-surface-300 mb-1.5">Password</label>
                        <div class="relative">
                            <Lock class="w-4 h-4 text-surface-500 absolute left-3.5 top-1/2 -translate-y-1/2 pointer-events-none" />
                            <input
                                v-model="form.password"
                                :type="showPassword ? 'text' : 'password'"
                                required
                                autocomplete="current-password"
                                placeholder="••••••••"
                                class="w-full pl-10 pr-10 py-2.5 rounded-xl bg-surface-950/60 border border-surface-800 text-xs text-white placeholder-surface-500 focus:outline-none focus:ring-2 focus:ring-brand-500/40 focus:border-brand-500 transition"
                            />
                            <button
                                type="button"
                                @click="showPassword = !showPassword"
                                class="absolute right-3 top-1/2 -translate-y-1/2 text-surface-500 hover:text-surface-300 p-1"
                            >
                                <Eye v-if="!showPassword" class="w-3.5 h-3.5" />
                                <EyeOff v-else class="w-3.5 h-3.5" />
                            </button>
                        </div>
                        <p v-if="form.errors.password" class="text-[11px] text-rose-400 mt-1">{{ form.errors.password }}</p>
                    </div>

                    <!-- Remember Me -->
                    <div class="flex items-center justify-between pt-1">
                        <label class="flex items-center gap-2 cursor-pointer">
                            <input
                                v-model="form.remember"
                                type="checkbox"
                                class="rounded border-surface-700 bg-surface-950 text-brand-600 focus:ring-brand-500/40 focus:ring-offset-0 w-3.5 h-3.5"
                            />
                            <span class="text-xs text-surface-400">Remember this session</span>
                        </label>
                    </div>

                    <!-- Submit Button -->
                    <button
                        type="submit"
                        :disabled="form.processing"
                        class="w-full mt-2 py-2.5 px-4 rounded-xl bg-brand-600 hover:bg-brand-500 text-white text-xs font-semibold shadow-lg shadow-brand-600/25 flex items-center justify-center gap-2 transition disabled:opacity-50"
                    >
                        <span v-if="form.processing">Authenticating...</span>
                        <span v-else class="flex items-center gap-1.5">
                            Sign In to Server <ArrowRight class="w-3.5 h-3.5" />
                        </span>
                    </button>
                </form>

                <!-- Dev Helper Box -->
                <div class="mt-6 p-3 rounded-xl bg-surface-950/40 border border-surface-800/60 text-[11px] text-surface-400">
                    <p class="text-surface-300 font-medium mb-1">Default Admin Credentials:</p>
                    <div class="font-mono text-[10px] text-surface-400 space-y-0.5">
                        <p>Email: <span class="text-brand-400">admin@kodepreneur.com</span></p>
                        <p>Pass: <span class="text-brand-400">password</span></p>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>
