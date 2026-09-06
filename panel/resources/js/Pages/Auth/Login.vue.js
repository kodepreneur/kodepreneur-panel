/// <reference types="../../../../node_modules/.vue-global-types/vue_3.5_0_0_0.d.ts" />
import { ref, computed } from 'vue';
import { useForm, usePage } from '@inertiajs/vue3';
import { Lock, Mail, Eye, EyeOff, ArrowRight, Layers, Sun, Moon, ShieldCheck, AlertCircle, Loader2, CheckCircle2, } from 'lucide-vue-next';
import { useTheme } from '@/composables/useTheme';
const { theme, toggleTheme } = useTheme();
const page = usePage();
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
debugger; /* PartiallyEnd: #3632/scriptSetup.vue */
const __VLS_ctx = {};
let __VLS_components;
let __VLS_directives;
__VLS_asFunctionalElement(__VLS_intrinsicElements.div, __VLS_intrinsicElements.div)({
    ...{ class: "min-h-screen bg-slate-50 dark:bg-[#070b14] text-slate-900 dark:text-slate-100 flex flex-col justify-between p-4 sm:p-6 selection:bg-brand-500 selection:text-white relative overflow-hidden transition-colors duration-200" },
});
__VLS_asFunctionalElement(__VLS_intrinsicElements.div, __VLS_intrinsicElements.div)({
    ...{ class: "absolute inset-0 pointer-events-none overflow-hidden" },
});
__VLS_asFunctionalElement(__VLS_intrinsicElements.div, __VLS_intrinsicElements.div)({
    ...{ class: "absolute -top-40 left-1/2 -translate-x-1/2 w-[700px] h-[500px] bg-gradient-to-b from-brand-500/15 via-indigo-600/10 to-transparent blur-[140px] rounded-full" },
});
__VLS_asFunctionalElement(__VLS_intrinsicElements.div, __VLS_intrinsicElements.div)({
    ...{ class: "absolute -bottom-40 right-10 w-[500px] h-[400px] bg-gradient-to-t from-indigo-500/10 via-purple-600/5 to-transparent blur-[130px] rounded-full" },
});
__VLS_asFunctionalElement(__VLS_intrinsicElements.div, __VLS_intrinsicElements.div)({
    ...{ class: "absolute inset-0 bg-[linear-gradient(to_right,#e2e8f015_1px,transparent_1px),linear-gradient(to_bottom,#e2e8f015_1px,transparent_1px)] dark:bg-[linear-gradient(to_right,#1e293b20_1px,transparent_1px),linear-gradient(to_bottom,#1e293b20_1px,transparent_1px)] bg-[size:4rem_4rem] [mask-image:radial-gradient(ellipse_60%_50%_at_50%_40%,#000_70%,transparent_100%)]" },
});
__VLS_asFunctionalElement(__VLS_intrinsicElements.header, __VLS_intrinsicElements.header)({
    ...{ class: "w-full max-w-5xl mx-auto flex items-center justify-between relative z-20 py-2" },
});
__VLS_asFunctionalElement(__VLS_intrinsicElements.div, __VLS_intrinsicElements.div)({
    ...{ class: "flex items-center gap-2.5" },
});
__VLS_asFunctionalElement(__VLS_intrinsicElements.div, __VLS_intrinsicElements.div)({
    ...{ class: "w-9 h-9 rounded-xl bg-gradient-to-tr from-brand-600 via-indigo-500 to-indigo-400 flex items-center justify-center shadow-lg shadow-brand-500/20 ring-1 ring-white/20 dark:ring-white/10" },
});
const __VLS_0 = {}.Layers;
/** @type {[typeof __VLS_components.Layers, ]} */ ;
// @ts-ignore
const __VLS_1 = __VLS_asFunctionalComponent(__VLS_0, new __VLS_0({
    ...{ class: "w-4 h-4 text-white" },
}));
const __VLS_2 = __VLS_1({
    ...{ class: "w-4 h-4 text-white" },
}, ...__VLS_functionalComponentArgsRest(__VLS_1));
__VLS_asFunctionalElement(__VLS_intrinsicElements.div, __VLS_intrinsicElements.div)({});
__VLS_asFunctionalElement(__VLS_intrinsicElements.span, __VLS_intrinsicElements.span)({
    ...{ class: "text-sm font-semibold tracking-tight text-slate-900 dark:text-white flex items-center gap-1.5" },
});
__VLS_asFunctionalElement(__VLS_intrinsicElements.span, __VLS_intrinsicElements.span)({
    ...{ class: "px-1.5 py-0.5 rounded text-[10px] font-mono font-medium bg-brand-50 text-brand-600 border border-brand-200 dark:bg-brand-500/10 dark:text-brand-400 dark:border-brand-500/20" },
});
__VLS_asFunctionalElement(__VLS_intrinsicElements.div, __VLS_intrinsicElements.div)({
    ...{ class: "flex items-center gap-3" },
});
__VLS_asFunctionalElement(__VLS_intrinsicElements.div, __VLS_intrinsicElements.div)({
    ...{ class: "hidden sm:flex items-center gap-2 px-2.5 py-1 rounded-full bg-white/80 dark:bg-surface-900/80 border border-slate-200/80 dark:border-surface-800 backdrop-blur-md text-[11px] font-medium text-slate-600 dark:text-surface-300 shadow-sm" },
});
__VLS_asFunctionalElement(__VLS_intrinsicElements.span, __VLS_intrinsicElements.span)({
    ...{ class: "relative flex h-2 w-2" },
});
__VLS_asFunctionalElement(__VLS_intrinsicElements.span, __VLS_intrinsicElements.span)({
    ...{ class: "animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-400 opacity-75" },
});
__VLS_asFunctionalElement(__VLS_intrinsicElements.span, __VLS_intrinsicElements.span)({
    ...{ class: "relative inline-flex rounded-full h-2 w-2 bg-emerald-500" },
});
__VLS_asFunctionalElement(__VLS_intrinsicElements.span, __VLS_intrinsicElements.span)({});
__VLS_asFunctionalElement(__VLS_intrinsicElements.button, __VLS_intrinsicElements.button)({
    ...{ onClick: (__VLS_ctx.toggleTheme) },
    type: "button",
    ...{ class: "p-2 rounded-xl bg-white/90 dark:bg-surface-900/90 border border-slate-200/80 dark:border-surface-800 text-slate-600 dark:text-surface-300 hover:text-slate-900 dark:hover:text-white hover:border-slate-300 dark:hover:border-surface-700 transition shadow-sm backdrop-blur-sm focus:outline-none focus:ring-2 focus:ring-brand-500/30" },
    title: (__VLS_ctx.theme === 'dark' ? 'Switch to Light Mode' : 'Switch to Dark Mode'),
    'aria-label': "Toggle color theme",
});
if (__VLS_ctx.theme === 'dark') {
    const __VLS_4 = {}.Sun;
    /** @type {[typeof __VLS_components.Sun, ]} */ ;
    // @ts-ignore
    const __VLS_5 = __VLS_asFunctionalComponent(__VLS_4, new __VLS_4({
        ...{ class: "w-4 h-4 text-amber-400" },
    }));
    const __VLS_6 = __VLS_5({
        ...{ class: "w-4 h-4 text-amber-400" },
    }, ...__VLS_functionalComponentArgsRest(__VLS_5));
}
else {
    const __VLS_8 = {}.Moon;
    /** @type {[typeof __VLS_components.Moon, ]} */ ;
    // @ts-ignore
    const __VLS_9 = __VLS_asFunctionalComponent(__VLS_8, new __VLS_8({
        ...{ class: "w-4 h-4 text-indigo-600" },
    }));
    const __VLS_10 = __VLS_9({
        ...{ class: "w-4 h-4 text-indigo-600" },
    }, ...__VLS_functionalComponentArgsRest(__VLS_9));
}
__VLS_asFunctionalElement(__VLS_intrinsicElements.main, __VLS_intrinsicElements.main)({
    ...{ class: "w-full max-w-[420px] mx-auto my-auto relative z-10 py-6" },
});
__VLS_asFunctionalElement(__VLS_intrinsicElements.div, __VLS_intrinsicElements.div)({
    ...{ class: "relative rounded-2xl bg-white/90 dark:bg-surface-900/85 border border-slate-200/90 dark:border-surface-800 shadow-2xl shadow-slate-900/5 dark:shadow-black/60 backdrop-blur-2xl p-6 sm:p-8 transition-all" },
});
__VLS_asFunctionalElement(__VLS_intrinsicElements.div, __VLS_intrinsicElements.div)({
    ...{ class: "absolute top-0 inset-x-0 h-[2px] bg-gradient-to-r from-transparent via-brand-500 to-transparent" },
});
__VLS_asFunctionalElement(__VLS_intrinsicElements.div, __VLS_intrinsicElements.div)({
    ...{ class: "flex flex-col items-center text-center mb-6" },
});
__VLS_asFunctionalElement(__VLS_intrinsicElements.div, __VLS_intrinsicElements.div)({
    ...{ class: "w-12 h-12 rounded-2xl bg-brand-50 dark:bg-brand-500/10 border border-brand-200/80 dark:border-brand-500/20 flex items-center justify-center mb-3.5 shadow-inner" },
});
const __VLS_12 = {}.Lock;
/** @type {[typeof __VLS_components.Lock, ]} */ ;
// @ts-ignore
const __VLS_13 = __VLS_asFunctionalComponent(__VLS_12, new __VLS_12({
    ...{ class: "w-5 h-5 text-brand-600 dark:text-brand-400" },
}));
const __VLS_14 = __VLS_13({
    ...{ class: "w-5 h-5 text-brand-600 dark:text-brand-400" },
}, ...__VLS_functionalComponentArgsRest(__VLS_13));
__VLS_asFunctionalElement(__VLS_intrinsicElements.h2, __VLS_intrinsicElements.h2)({
    ...{ class: "text-lg font-bold tracking-tight text-slate-900 dark:text-white" },
});
__VLS_asFunctionalElement(__VLS_intrinsicElements.p, __VLS_intrinsicElements.p)({
    ...{ class: "text-xs text-slate-500 dark:text-surface-400 mt-1 max-w-xs" },
});
if (__VLS_ctx.flash?.error) {
    __VLS_asFunctionalElement(__VLS_intrinsicElements.div, __VLS_intrinsicElements.div)({
        ...{ class: "mb-4 p-3 rounded-xl bg-rose-50 dark:bg-rose-500/10 border border-rose-200 dark:border-rose-500/20 flex items-start gap-2.5 text-xs text-rose-700 dark:text-rose-300" },
    });
    const __VLS_16 = {}.AlertCircle;
    /** @type {[typeof __VLS_components.AlertCircle, ]} */ ;
    // @ts-ignore
    const __VLS_17 = __VLS_asFunctionalComponent(__VLS_16, new __VLS_16({
        ...{ class: "w-4 h-4 text-rose-500 shrink-0 mt-0.5" },
    }));
    const __VLS_18 = __VLS_17({
        ...{ class: "w-4 h-4 text-rose-500 shrink-0 mt-0.5" },
    }, ...__VLS_functionalComponentArgsRest(__VLS_17));
    __VLS_asFunctionalElement(__VLS_intrinsicElements.span, __VLS_intrinsicElements.span)({});
    (__VLS_ctx.flash.error);
}
if (__VLS_ctx.flash?.success) {
    __VLS_asFunctionalElement(__VLS_intrinsicElements.div, __VLS_intrinsicElements.div)({
        ...{ class: "mb-4 p-3 rounded-xl bg-emerald-50 dark:bg-emerald-500/10 border border-emerald-200 dark:border-emerald-500/20 flex items-start gap-2.5 text-xs text-emerald-700 dark:text-emerald-300" },
    });
    const __VLS_20 = {}.CheckCircle2;
    /** @type {[typeof __VLS_components.CheckCircle2, ]} */ ;
    // @ts-ignore
    const __VLS_21 = __VLS_asFunctionalComponent(__VLS_20, new __VLS_20({
        ...{ class: "w-4 h-4 text-emerald-500 shrink-0 mt-0.5" },
    }));
    const __VLS_22 = __VLS_21({
        ...{ class: "w-4 h-4 text-emerald-500 shrink-0 mt-0.5" },
    }, ...__VLS_functionalComponentArgsRest(__VLS_21));
    __VLS_asFunctionalElement(__VLS_intrinsicElements.span, __VLS_intrinsicElements.span)({});
    (__VLS_ctx.flash.success);
}
__VLS_asFunctionalElement(__VLS_intrinsicElements.form, __VLS_intrinsicElements.form)({
    ...{ onSubmit: (__VLS_ctx.submit) },
    ...{ class: "space-y-4" },
});
__VLS_asFunctionalElement(__VLS_intrinsicElements.div, __VLS_intrinsicElements.div)({});
__VLS_asFunctionalElement(__VLS_intrinsicElements.label, __VLS_intrinsicElements.label)({
    for: "email",
    ...{ class: "block text-xs font-medium text-slate-700 dark:text-surface-200 mb-1.5" },
});
__VLS_asFunctionalElement(__VLS_intrinsicElements.div, __VLS_intrinsicElements.div)({
    ...{ class: "relative" },
});
__VLS_asFunctionalElement(__VLS_intrinsicElements.div, __VLS_intrinsicElements.div)({
    ...{ class: "absolute inset-y-0 left-0 pl-3.5 flex items-center pointer-events-none text-slate-400 dark:text-surface-500" },
});
const __VLS_24 = {}.Mail;
/** @type {[typeof __VLS_components.Mail, ]} */ ;
// @ts-ignore
const __VLS_25 = __VLS_asFunctionalComponent(__VLS_24, new __VLS_24({
    ...{ class: "w-4 h-4" },
}));
const __VLS_26 = __VLS_25({
    ...{ class: "w-4 h-4" },
}, ...__VLS_functionalComponentArgsRest(__VLS_25));
__VLS_asFunctionalElement(__VLS_intrinsicElements.input)({
    id: "email",
    type: "email",
    required: true,
    autofocus: true,
    autocomplete: "email",
    placeholder: "admin@kodepreneur.com",
    ...{ class: ([
            'w-full pl-10 pr-4 py-2.5 rounded-xl text-xs bg-slate-50 dark:bg-surface-950/60 border text-slate-900 dark:text-white placeholder-slate-400 dark:placeholder-surface-500 focus:outline-none focus:ring-2 focus:ring-brand-500/30 transition duration-150',
            __VLS_ctx.form.errors.email
                ? 'border-rose-300 dark:border-rose-500/50 focus:border-rose-500 focus:ring-rose-500/20'
                : 'border-slate-200 dark:border-surface-800 focus:border-brand-500'
        ]) },
});
(__VLS_ctx.form.email);
if (__VLS_ctx.form.errors.email) {
    __VLS_asFunctionalElement(__VLS_intrinsicElements.div, __VLS_intrinsicElements.div)({
        ...{ class: "flex items-center gap-1.5 text-[11px] text-rose-500 mt-1.5" },
    });
    const __VLS_28 = {}.AlertCircle;
    /** @type {[typeof __VLS_components.AlertCircle, ]} */ ;
    // @ts-ignore
    const __VLS_29 = __VLS_asFunctionalComponent(__VLS_28, new __VLS_28({
        ...{ class: "w-3.5 h-3.5 shrink-0" },
    }));
    const __VLS_30 = __VLS_29({
        ...{ class: "w-3.5 h-3.5 shrink-0" },
    }, ...__VLS_functionalComponentArgsRest(__VLS_29));
    __VLS_asFunctionalElement(__VLS_intrinsicElements.span, __VLS_intrinsicElements.span)({});
    (__VLS_ctx.form.errors.email);
}
__VLS_asFunctionalElement(__VLS_intrinsicElements.div, __VLS_intrinsicElements.div)({});
__VLS_asFunctionalElement(__VLS_intrinsicElements.label, __VLS_intrinsicElements.label)({
    for: "password",
    ...{ class: "block text-xs font-medium text-slate-700 dark:text-surface-200 mb-1.5" },
});
__VLS_asFunctionalElement(__VLS_intrinsicElements.div, __VLS_intrinsicElements.div)({
    ...{ class: "relative" },
});
__VLS_asFunctionalElement(__VLS_intrinsicElements.div, __VLS_intrinsicElements.div)({
    ...{ class: "absolute inset-y-0 left-0 pl-3.5 flex items-center pointer-events-none text-slate-400 dark:text-surface-500" },
});
const __VLS_32 = {}.Lock;
/** @type {[typeof __VLS_components.Lock, ]} */ ;
// @ts-ignore
const __VLS_33 = __VLS_asFunctionalComponent(__VLS_32, new __VLS_32({
    ...{ class: "w-4 h-4" },
}));
const __VLS_34 = __VLS_33({
    ...{ class: "w-4 h-4" },
}, ...__VLS_functionalComponentArgsRest(__VLS_33));
__VLS_asFunctionalElement(__VLS_intrinsicElements.input)({
    id: "password",
    type: (__VLS_ctx.showPassword ? 'text' : 'password'),
    required: true,
    autocomplete: "current-password",
    placeholder: "••••••••••••",
    ...{ class: ([
            'w-full pl-10 pr-10 py-2.5 rounded-xl text-xs bg-slate-50 dark:bg-surface-950/60 border text-slate-900 dark:text-white placeholder-slate-400 dark:placeholder-surface-500 focus:outline-none focus:ring-2 focus:ring-brand-500/30 transition duration-150',
            __VLS_ctx.form.errors.password
                ? 'border-rose-300 dark:border-rose-500/50 focus:border-rose-500 focus:ring-rose-500/20'
                : 'border-slate-200 dark:border-surface-800 focus:border-brand-500'
        ]) },
});
(__VLS_ctx.form.password);
__VLS_asFunctionalElement(__VLS_intrinsicElements.button, __VLS_intrinsicElements.button)({
    ...{ onClick: (...[$event]) => {
            __VLS_ctx.showPassword = !__VLS_ctx.showPassword;
        } },
    type: "button",
    ...{ class: "absolute inset-y-0 right-0 pr-3 flex items-center text-slate-400 hover:text-slate-600 dark:text-surface-500 dark:hover:text-surface-300 focus:outline-none transition" },
    'aria-label': (__VLS_ctx.showPassword ? 'Hide password' : 'Show password'),
});
if (__VLS_ctx.showPassword) {
    const __VLS_36 = {}.EyeOff;
    /** @type {[typeof __VLS_components.EyeOff, ]} */ ;
    // @ts-ignore
    const __VLS_37 = __VLS_asFunctionalComponent(__VLS_36, new __VLS_36({
        ...{ class: "w-4 h-4" },
    }));
    const __VLS_38 = __VLS_37({
        ...{ class: "w-4 h-4" },
    }, ...__VLS_functionalComponentArgsRest(__VLS_37));
}
else {
    const __VLS_40 = {}.Eye;
    /** @type {[typeof __VLS_components.Eye, ]} */ ;
    // @ts-ignore
    const __VLS_41 = __VLS_asFunctionalComponent(__VLS_40, new __VLS_40({
        ...{ class: "w-4 h-4" },
    }));
    const __VLS_42 = __VLS_41({
        ...{ class: "w-4 h-4" },
    }, ...__VLS_functionalComponentArgsRest(__VLS_41));
}
if (__VLS_ctx.form.errors.password) {
    __VLS_asFunctionalElement(__VLS_intrinsicElements.div, __VLS_intrinsicElements.div)({
        ...{ class: "flex items-center gap-1.5 text-[11px] text-rose-500 mt-1.5" },
    });
    const __VLS_44 = {}.AlertCircle;
    /** @type {[typeof __VLS_components.AlertCircle, ]} */ ;
    // @ts-ignore
    const __VLS_45 = __VLS_asFunctionalComponent(__VLS_44, new __VLS_44({
        ...{ class: "w-3.5 h-3.5 shrink-0" },
    }));
    const __VLS_46 = __VLS_45({
        ...{ class: "w-3.5 h-3.5 shrink-0" },
    }, ...__VLS_functionalComponentArgsRest(__VLS_45));
    __VLS_asFunctionalElement(__VLS_intrinsicElements.span, __VLS_intrinsicElements.span)({});
    (__VLS_ctx.form.errors.password);
}
__VLS_asFunctionalElement(__VLS_intrinsicElements.div, __VLS_intrinsicElements.div)({
    ...{ class: "flex items-center justify-between pt-0.5" },
});
__VLS_asFunctionalElement(__VLS_intrinsicElements.label, __VLS_intrinsicElements.label)({
    ...{ class: "flex items-center gap-2 cursor-pointer select-none group" },
});
__VLS_asFunctionalElement(__VLS_intrinsicElements.input)({
    type: "checkbox",
    ...{ class: "rounded border-slate-300 dark:border-surface-700 bg-white dark:bg-surface-950 text-brand-600 focus:ring-brand-500/30 focus:ring-offset-0 w-3.5 h-3.5 cursor-pointer" },
});
(__VLS_ctx.form.remember);
__VLS_asFunctionalElement(__VLS_intrinsicElements.span, __VLS_intrinsicElements.span)({
    ...{ class: "text-xs text-slate-600 dark:text-surface-300 group-hover:text-slate-900 dark:group-hover:text-white transition" },
});
__VLS_asFunctionalElement(__VLS_intrinsicElements.button, __VLS_intrinsicElements.button)({
    type: "submit",
    disabled: (__VLS_ctx.form.processing),
    ...{ class: "w-full mt-2 py-2.5 px-4 rounded-xl bg-gradient-to-r from-brand-600 via-indigo-600 to-brand-500 hover:from-brand-500 hover:to-indigo-500 text-white text-xs font-semibold shadow-md shadow-brand-500/25 active:scale-[0.99] flex items-center justify-center gap-2 transition duration-150 disabled:opacity-60 disabled:cursor-not-allowed cursor-pointer group" },
});
if (__VLS_ctx.form.processing) {
    const __VLS_48 = {}.Loader2;
    /** @type {[typeof __VLS_components.Loader2, ]} */ ;
    // @ts-ignore
    const __VLS_49 = __VLS_asFunctionalComponent(__VLS_48, new __VLS_48({
        ...{ class: "w-4 h-4 animate-spin" },
    }));
    const __VLS_50 = __VLS_49({
        ...{ class: "w-4 h-4 animate-spin" },
    }, ...__VLS_functionalComponentArgsRest(__VLS_49));
}
if (__VLS_ctx.form.processing) {
    __VLS_asFunctionalElement(__VLS_intrinsicElements.span, __VLS_intrinsicElements.span)({});
}
else {
    __VLS_asFunctionalElement(__VLS_intrinsicElements.span, __VLS_intrinsicElements.span)({});
    const __VLS_52 = {}.ArrowRight;
    /** @type {[typeof __VLS_components.ArrowRight, ]} */ ;
    // @ts-ignore
    const __VLS_53 = __VLS_asFunctionalComponent(__VLS_52, new __VLS_52({
        ...{ class: "w-3.5 h-3.5 group-hover:translate-x-0.5 transition-transform" },
    }));
    const __VLS_54 = __VLS_53({
        ...{ class: "w-3.5 h-3.5 group-hover:translate-x-0.5 transition-transform" },
    }, ...__VLS_functionalComponentArgsRest(__VLS_53));
}
__VLS_asFunctionalElement(__VLS_intrinsicElements.div, __VLS_intrinsicElements.div)({
    ...{ class: "mt-6 pt-5 border-t border-slate-200/70 dark:border-surface-800/70 flex items-center justify-center gap-2 text-[11px] text-slate-400 dark:text-surface-500" },
});
const __VLS_56 = {}.ShieldCheck;
/** @type {[typeof __VLS_components.ShieldCheck, ]} */ ;
// @ts-ignore
const __VLS_57 = __VLS_asFunctionalComponent(__VLS_56, new __VLS_56({
    ...{ class: "w-3.5 h-3.5 text-emerald-500 shrink-0" },
}));
const __VLS_58 = __VLS_57({
    ...{ class: "w-3.5 h-3.5 text-emerald-500 shrink-0" },
}, ...__VLS_functionalComponentArgsRest(__VLS_57));
__VLS_asFunctionalElement(__VLS_intrinsicElements.span, __VLS_intrinsicElements.span)({});
__VLS_asFunctionalElement(__VLS_intrinsicElements.footer, __VLS_intrinsicElements.footer)({
    ...{ class: "w-full max-w-5xl mx-auto flex flex-col sm:flex-row items-center justify-between gap-2 text-center text-[11px] text-slate-400 dark:text-surface-500 relative z-20 py-2" },
});
__VLS_asFunctionalElement(__VLS_intrinsicElements.div, __VLS_intrinsicElements.div)({});
(new Date().getFullYear());
__VLS_asFunctionalElement(__VLS_intrinsicElements.div, __VLS_intrinsicElements.div)({
    ...{ class: "flex items-center gap-4" },
});
__VLS_asFunctionalElement(__VLS_intrinsicElements.span, __VLS_intrinsicElements.span)({});
__VLS_asFunctionalElement(__VLS_intrinsicElements.span, __VLS_intrinsicElements.span)({});
__VLS_asFunctionalElement(__VLS_intrinsicElements.span, __VLS_intrinsicElements.span)({});
__VLS_asFunctionalElement(__VLS_intrinsicElements.span, __VLS_intrinsicElements.span)({});
__VLS_asFunctionalElement(__VLS_intrinsicElements.span, __VLS_intrinsicElements.span)({});
/** @type {__VLS_StyleScopedClasses['min-h-screen']} */ ;
/** @type {__VLS_StyleScopedClasses['bg-slate-50']} */ ;
/** @type {__VLS_StyleScopedClasses['dark:bg-[#070b14]']} */ ;
/** @type {__VLS_StyleScopedClasses['text-slate-900']} */ ;
/** @type {__VLS_StyleScopedClasses['dark:text-slate-100']} */ ;
/** @type {__VLS_StyleScopedClasses['flex']} */ ;
/** @type {__VLS_StyleScopedClasses['flex-col']} */ ;
/** @type {__VLS_StyleScopedClasses['justify-between']} */ ;
/** @type {__VLS_StyleScopedClasses['p-4']} */ ;
/** @type {__VLS_StyleScopedClasses['sm:p-6']} */ ;
/** @type {__VLS_StyleScopedClasses['selection:bg-brand-500']} */ ;
/** @type {__VLS_StyleScopedClasses['selection:text-white']} */ ;
/** @type {__VLS_StyleScopedClasses['relative']} */ ;
/** @type {__VLS_StyleScopedClasses['overflow-hidden']} */ ;
/** @type {__VLS_StyleScopedClasses['transition-colors']} */ ;
/** @type {__VLS_StyleScopedClasses['duration-200']} */ ;
/** @type {__VLS_StyleScopedClasses['absolute']} */ ;
/** @type {__VLS_StyleScopedClasses['inset-0']} */ ;
/** @type {__VLS_StyleScopedClasses['pointer-events-none']} */ ;
/** @type {__VLS_StyleScopedClasses['overflow-hidden']} */ ;
/** @type {__VLS_StyleScopedClasses['absolute']} */ ;
/** @type {__VLS_StyleScopedClasses['-top-40']} */ ;
/** @type {__VLS_StyleScopedClasses['left-1/2']} */ ;
/** @type {__VLS_StyleScopedClasses['-translate-x-1/2']} */ ;
/** @type {__VLS_StyleScopedClasses['w-[700px]']} */ ;
/** @type {__VLS_StyleScopedClasses['h-[500px]']} */ ;
/** @type {__VLS_StyleScopedClasses['bg-gradient-to-b']} */ ;
/** @type {__VLS_StyleScopedClasses['from-brand-500/15']} */ ;
/** @type {__VLS_StyleScopedClasses['via-indigo-600/10']} */ ;
/** @type {__VLS_StyleScopedClasses['to-transparent']} */ ;
/** @type {__VLS_StyleScopedClasses['blur-[140px]']} */ ;
/** @type {__VLS_StyleScopedClasses['rounded-full']} */ ;
/** @type {__VLS_StyleScopedClasses['absolute']} */ ;
/** @type {__VLS_StyleScopedClasses['-bottom-40']} */ ;
/** @type {__VLS_StyleScopedClasses['right-10']} */ ;
/** @type {__VLS_StyleScopedClasses['w-[500px]']} */ ;
/** @type {__VLS_StyleScopedClasses['h-[400px]']} */ ;
/** @type {__VLS_StyleScopedClasses['bg-gradient-to-t']} */ ;
/** @type {__VLS_StyleScopedClasses['from-indigo-500/10']} */ ;
/** @type {__VLS_StyleScopedClasses['via-purple-600/5']} */ ;
/** @type {__VLS_StyleScopedClasses['to-transparent']} */ ;
/** @type {__VLS_StyleScopedClasses['blur-[130px]']} */ ;
/** @type {__VLS_StyleScopedClasses['rounded-full']} */ ;
/** @type {__VLS_StyleScopedClasses['absolute']} */ ;
/** @type {__VLS_StyleScopedClasses['inset-0']} */ ;
/** @type {__VLS_StyleScopedClasses['bg-[linear-gradient(to_right,#e2e8f015_1px,transparent_1px),linear-gradient(to_bottom,#e2e8f015_1px,transparent_1px)]']} */ ;
/** @type {__VLS_StyleScopedClasses['dark:bg-[linear-gradient(to_right,#1e293b20_1px,transparent_1px),linear-gradient(to_bottom,#1e293b20_1px,transparent_1px)]']} */ ;
/** @type {__VLS_StyleScopedClasses['bg-[size:4rem_4rem]']} */ ;
/** @type {__VLS_StyleScopedClasses['[mask-image:radial-gradient(ellipse_60%_50%_at_50%_40%,#000_70%,transparent_100%)]']} */ ;
/** @type {__VLS_StyleScopedClasses['w-full']} */ ;
/** @type {__VLS_StyleScopedClasses['max-w-5xl']} */ ;
/** @type {__VLS_StyleScopedClasses['mx-auto']} */ ;
/** @type {__VLS_StyleScopedClasses['flex']} */ ;
/** @type {__VLS_StyleScopedClasses['items-center']} */ ;
/** @type {__VLS_StyleScopedClasses['justify-between']} */ ;
/** @type {__VLS_StyleScopedClasses['relative']} */ ;
/** @type {__VLS_StyleScopedClasses['z-20']} */ ;
/** @type {__VLS_StyleScopedClasses['py-2']} */ ;
/** @type {__VLS_StyleScopedClasses['flex']} */ ;
/** @type {__VLS_StyleScopedClasses['items-center']} */ ;
/** @type {__VLS_StyleScopedClasses['gap-2.5']} */ ;
/** @type {__VLS_StyleScopedClasses['w-9']} */ ;
/** @type {__VLS_StyleScopedClasses['h-9']} */ ;
/** @type {__VLS_StyleScopedClasses['rounded-xl']} */ ;
/** @type {__VLS_StyleScopedClasses['bg-gradient-to-tr']} */ ;
/** @type {__VLS_StyleScopedClasses['from-brand-600']} */ ;
/** @type {__VLS_StyleScopedClasses['via-indigo-500']} */ ;
/** @type {__VLS_StyleScopedClasses['to-indigo-400']} */ ;
/** @type {__VLS_StyleScopedClasses['flex']} */ ;
/** @type {__VLS_StyleScopedClasses['items-center']} */ ;
/** @type {__VLS_StyleScopedClasses['justify-center']} */ ;
/** @type {__VLS_StyleScopedClasses['shadow-lg']} */ ;
/** @type {__VLS_StyleScopedClasses['shadow-brand-500/20']} */ ;
/** @type {__VLS_StyleScopedClasses['ring-1']} */ ;
/** @type {__VLS_StyleScopedClasses['ring-white/20']} */ ;
/** @type {__VLS_StyleScopedClasses['dark:ring-white/10']} */ ;
/** @type {__VLS_StyleScopedClasses['w-4']} */ ;
/** @type {__VLS_StyleScopedClasses['h-4']} */ ;
/** @type {__VLS_StyleScopedClasses['text-white']} */ ;
/** @type {__VLS_StyleScopedClasses['text-sm']} */ ;
/** @type {__VLS_StyleScopedClasses['font-semibold']} */ ;
/** @type {__VLS_StyleScopedClasses['tracking-tight']} */ ;
/** @type {__VLS_StyleScopedClasses['text-slate-900']} */ ;
/** @type {__VLS_StyleScopedClasses['dark:text-white']} */ ;
/** @type {__VLS_StyleScopedClasses['flex']} */ ;
/** @type {__VLS_StyleScopedClasses['items-center']} */ ;
/** @type {__VLS_StyleScopedClasses['gap-1.5']} */ ;
/** @type {__VLS_StyleScopedClasses['px-1.5']} */ ;
/** @type {__VLS_StyleScopedClasses['py-0.5']} */ ;
/** @type {__VLS_StyleScopedClasses['rounded']} */ ;
/** @type {__VLS_StyleScopedClasses['text-[10px]']} */ ;
/** @type {__VLS_StyleScopedClasses['font-mono']} */ ;
/** @type {__VLS_StyleScopedClasses['font-medium']} */ ;
/** @type {__VLS_StyleScopedClasses['bg-brand-50']} */ ;
/** @type {__VLS_StyleScopedClasses['text-brand-600']} */ ;
/** @type {__VLS_StyleScopedClasses['border']} */ ;
/** @type {__VLS_StyleScopedClasses['border-brand-200']} */ ;
/** @type {__VLS_StyleScopedClasses['dark:bg-brand-500/10']} */ ;
/** @type {__VLS_StyleScopedClasses['dark:text-brand-400']} */ ;
/** @type {__VLS_StyleScopedClasses['dark:border-brand-500/20']} */ ;
/** @type {__VLS_StyleScopedClasses['flex']} */ ;
/** @type {__VLS_StyleScopedClasses['items-center']} */ ;
/** @type {__VLS_StyleScopedClasses['gap-3']} */ ;
/** @type {__VLS_StyleScopedClasses['hidden']} */ ;
/** @type {__VLS_StyleScopedClasses['sm:flex']} */ ;
/** @type {__VLS_StyleScopedClasses['items-center']} */ ;
/** @type {__VLS_StyleScopedClasses['gap-2']} */ ;
/** @type {__VLS_StyleScopedClasses['px-2.5']} */ ;
/** @type {__VLS_StyleScopedClasses['py-1']} */ ;
/** @type {__VLS_StyleScopedClasses['rounded-full']} */ ;
/** @type {__VLS_StyleScopedClasses['bg-white/80']} */ ;
/** @type {__VLS_StyleScopedClasses['dark:bg-surface-900/80']} */ ;
/** @type {__VLS_StyleScopedClasses['border']} */ ;
/** @type {__VLS_StyleScopedClasses['border-slate-200/80']} */ ;
/** @type {__VLS_StyleScopedClasses['dark:border-surface-800']} */ ;
/** @type {__VLS_StyleScopedClasses['backdrop-blur-md']} */ ;
/** @type {__VLS_StyleScopedClasses['text-[11px]']} */ ;
/** @type {__VLS_StyleScopedClasses['font-medium']} */ ;
/** @type {__VLS_StyleScopedClasses['text-slate-600']} */ ;
/** @type {__VLS_StyleScopedClasses['dark:text-surface-300']} */ ;
/** @type {__VLS_StyleScopedClasses['shadow-sm']} */ ;
/** @type {__VLS_StyleScopedClasses['relative']} */ ;
/** @type {__VLS_StyleScopedClasses['flex']} */ ;
/** @type {__VLS_StyleScopedClasses['h-2']} */ ;
/** @type {__VLS_StyleScopedClasses['w-2']} */ ;
/** @type {__VLS_StyleScopedClasses['animate-ping']} */ ;
/** @type {__VLS_StyleScopedClasses['absolute']} */ ;
/** @type {__VLS_StyleScopedClasses['inline-flex']} */ ;
/** @type {__VLS_StyleScopedClasses['h-full']} */ ;
/** @type {__VLS_StyleScopedClasses['w-full']} */ ;
/** @type {__VLS_StyleScopedClasses['rounded-full']} */ ;
/** @type {__VLS_StyleScopedClasses['bg-emerald-400']} */ ;
/** @type {__VLS_StyleScopedClasses['opacity-75']} */ ;
/** @type {__VLS_StyleScopedClasses['relative']} */ ;
/** @type {__VLS_StyleScopedClasses['inline-flex']} */ ;
/** @type {__VLS_StyleScopedClasses['rounded-full']} */ ;
/** @type {__VLS_StyleScopedClasses['h-2']} */ ;
/** @type {__VLS_StyleScopedClasses['w-2']} */ ;
/** @type {__VLS_StyleScopedClasses['bg-emerald-500']} */ ;
/** @type {__VLS_StyleScopedClasses['p-2']} */ ;
/** @type {__VLS_StyleScopedClasses['rounded-xl']} */ ;
/** @type {__VLS_StyleScopedClasses['bg-white/90']} */ ;
/** @type {__VLS_StyleScopedClasses['dark:bg-surface-900/90']} */ ;
/** @type {__VLS_StyleScopedClasses['border']} */ ;
/** @type {__VLS_StyleScopedClasses['border-slate-200/80']} */ ;
/** @type {__VLS_StyleScopedClasses['dark:border-surface-800']} */ ;
/** @type {__VLS_StyleScopedClasses['text-slate-600']} */ ;
/** @type {__VLS_StyleScopedClasses['dark:text-surface-300']} */ ;
/** @type {__VLS_StyleScopedClasses['hover:text-slate-900']} */ ;
/** @type {__VLS_StyleScopedClasses['dark:hover:text-white']} */ ;
/** @type {__VLS_StyleScopedClasses['hover:border-slate-300']} */ ;
/** @type {__VLS_StyleScopedClasses['dark:hover:border-surface-700']} */ ;
/** @type {__VLS_StyleScopedClasses['transition']} */ ;
/** @type {__VLS_StyleScopedClasses['shadow-sm']} */ ;
/** @type {__VLS_StyleScopedClasses['backdrop-blur-sm']} */ ;
/** @type {__VLS_StyleScopedClasses['focus:outline-none']} */ ;
/** @type {__VLS_StyleScopedClasses['focus:ring-2']} */ ;
/** @type {__VLS_StyleScopedClasses['focus:ring-brand-500/30']} */ ;
/** @type {__VLS_StyleScopedClasses['w-4']} */ ;
/** @type {__VLS_StyleScopedClasses['h-4']} */ ;
/** @type {__VLS_StyleScopedClasses['text-amber-400']} */ ;
/** @type {__VLS_StyleScopedClasses['w-4']} */ ;
/** @type {__VLS_StyleScopedClasses['h-4']} */ ;
/** @type {__VLS_StyleScopedClasses['text-indigo-600']} */ ;
/** @type {__VLS_StyleScopedClasses['w-full']} */ ;
/** @type {__VLS_StyleScopedClasses['max-w-[420px]']} */ ;
/** @type {__VLS_StyleScopedClasses['mx-auto']} */ ;
/** @type {__VLS_StyleScopedClasses['my-auto']} */ ;
/** @type {__VLS_StyleScopedClasses['relative']} */ ;
/** @type {__VLS_StyleScopedClasses['z-10']} */ ;
/** @type {__VLS_StyleScopedClasses['py-6']} */ ;
/** @type {__VLS_StyleScopedClasses['relative']} */ ;
/** @type {__VLS_StyleScopedClasses['rounded-2xl']} */ ;
/** @type {__VLS_StyleScopedClasses['bg-white/90']} */ ;
/** @type {__VLS_StyleScopedClasses['dark:bg-surface-900/85']} */ ;
/** @type {__VLS_StyleScopedClasses['border']} */ ;
/** @type {__VLS_StyleScopedClasses['border-slate-200/90']} */ ;
/** @type {__VLS_StyleScopedClasses['dark:border-surface-800']} */ ;
/** @type {__VLS_StyleScopedClasses['shadow-2xl']} */ ;
/** @type {__VLS_StyleScopedClasses['shadow-slate-900/5']} */ ;
/** @type {__VLS_StyleScopedClasses['dark:shadow-black/60']} */ ;
/** @type {__VLS_StyleScopedClasses['backdrop-blur-2xl']} */ ;
/** @type {__VLS_StyleScopedClasses['p-6']} */ ;
/** @type {__VLS_StyleScopedClasses['sm:p-8']} */ ;
/** @type {__VLS_StyleScopedClasses['transition-all']} */ ;
/** @type {__VLS_StyleScopedClasses['absolute']} */ ;
/** @type {__VLS_StyleScopedClasses['top-0']} */ ;
/** @type {__VLS_StyleScopedClasses['inset-x-0']} */ ;
/** @type {__VLS_StyleScopedClasses['h-[2px]']} */ ;
/** @type {__VLS_StyleScopedClasses['bg-gradient-to-r']} */ ;
/** @type {__VLS_StyleScopedClasses['from-transparent']} */ ;
/** @type {__VLS_StyleScopedClasses['via-brand-500']} */ ;
/** @type {__VLS_StyleScopedClasses['to-transparent']} */ ;
/** @type {__VLS_StyleScopedClasses['flex']} */ ;
/** @type {__VLS_StyleScopedClasses['flex-col']} */ ;
/** @type {__VLS_StyleScopedClasses['items-center']} */ ;
/** @type {__VLS_StyleScopedClasses['text-center']} */ ;
/** @type {__VLS_StyleScopedClasses['mb-6']} */ ;
/** @type {__VLS_StyleScopedClasses['w-12']} */ ;
/** @type {__VLS_StyleScopedClasses['h-12']} */ ;
/** @type {__VLS_StyleScopedClasses['rounded-2xl']} */ ;
/** @type {__VLS_StyleScopedClasses['bg-brand-50']} */ ;
/** @type {__VLS_StyleScopedClasses['dark:bg-brand-500/10']} */ ;
/** @type {__VLS_StyleScopedClasses['border']} */ ;
/** @type {__VLS_StyleScopedClasses['border-brand-200/80']} */ ;
/** @type {__VLS_StyleScopedClasses['dark:border-brand-500/20']} */ ;
/** @type {__VLS_StyleScopedClasses['flex']} */ ;
/** @type {__VLS_StyleScopedClasses['items-center']} */ ;
/** @type {__VLS_StyleScopedClasses['justify-center']} */ ;
/** @type {__VLS_StyleScopedClasses['mb-3.5']} */ ;
/** @type {__VLS_StyleScopedClasses['shadow-inner']} */ ;
/** @type {__VLS_StyleScopedClasses['w-5']} */ ;
/** @type {__VLS_StyleScopedClasses['h-5']} */ ;
/** @type {__VLS_StyleScopedClasses['text-brand-600']} */ ;
/** @type {__VLS_StyleScopedClasses['dark:text-brand-400']} */ ;
/** @type {__VLS_StyleScopedClasses['text-lg']} */ ;
/** @type {__VLS_StyleScopedClasses['font-bold']} */ ;
/** @type {__VLS_StyleScopedClasses['tracking-tight']} */ ;
/** @type {__VLS_StyleScopedClasses['text-slate-900']} */ ;
/** @type {__VLS_StyleScopedClasses['dark:text-white']} */ ;
/** @type {__VLS_StyleScopedClasses['text-xs']} */ ;
/** @type {__VLS_StyleScopedClasses['text-slate-500']} */ ;
/** @type {__VLS_StyleScopedClasses['dark:text-surface-400']} */ ;
/** @type {__VLS_StyleScopedClasses['mt-1']} */ ;
/** @type {__VLS_StyleScopedClasses['max-w-xs']} */ ;
/** @type {__VLS_StyleScopedClasses['mb-4']} */ ;
/** @type {__VLS_StyleScopedClasses['p-3']} */ ;
/** @type {__VLS_StyleScopedClasses['rounded-xl']} */ ;
/** @type {__VLS_StyleScopedClasses['bg-rose-50']} */ ;
/** @type {__VLS_StyleScopedClasses['dark:bg-rose-500/10']} */ ;
/** @type {__VLS_StyleScopedClasses['border']} */ ;
/** @type {__VLS_StyleScopedClasses['border-rose-200']} */ ;
/** @type {__VLS_StyleScopedClasses['dark:border-rose-500/20']} */ ;
/** @type {__VLS_StyleScopedClasses['flex']} */ ;
/** @type {__VLS_StyleScopedClasses['items-start']} */ ;
/** @type {__VLS_StyleScopedClasses['gap-2.5']} */ ;
/** @type {__VLS_StyleScopedClasses['text-xs']} */ ;
/** @type {__VLS_StyleScopedClasses['text-rose-700']} */ ;
/** @type {__VLS_StyleScopedClasses['dark:text-rose-300']} */ ;
/** @type {__VLS_StyleScopedClasses['w-4']} */ ;
/** @type {__VLS_StyleScopedClasses['h-4']} */ ;
/** @type {__VLS_StyleScopedClasses['text-rose-500']} */ ;
/** @type {__VLS_StyleScopedClasses['shrink-0']} */ ;
/** @type {__VLS_StyleScopedClasses['mt-0.5']} */ ;
/** @type {__VLS_StyleScopedClasses['mb-4']} */ ;
/** @type {__VLS_StyleScopedClasses['p-3']} */ ;
/** @type {__VLS_StyleScopedClasses['rounded-xl']} */ ;
/** @type {__VLS_StyleScopedClasses['bg-emerald-50']} */ ;
/** @type {__VLS_StyleScopedClasses['dark:bg-emerald-500/10']} */ ;
/** @type {__VLS_StyleScopedClasses['border']} */ ;
/** @type {__VLS_StyleScopedClasses['border-emerald-200']} */ ;
/** @type {__VLS_StyleScopedClasses['dark:border-emerald-500/20']} */ ;
/** @type {__VLS_StyleScopedClasses['flex']} */ ;
/** @type {__VLS_StyleScopedClasses['items-start']} */ ;
/** @type {__VLS_StyleScopedClasses['gap-2.5']} */ ;
/** @type {__VLS_StyleScopedClasses['text-xs']} */ ;
/** @type {__VLS_StyleScopedClasses['text-emerald-700']} */ ;
/** @type {__VLS_StyleScopedClasses['dark:text-emerald-300']} */ ;
/** @type {__VLS_StyleScopedClasses['w-4']} */ ;
/** @type {__VLS_StyleScopedClasses['h-4']} */ ;
/** @type {__VLS_StyleScopedClasses['text-emerald-500']} */ ;
/** @type {__VLS_StyleScopedClasses['shrink-0']} */ ;
/** @type {__VLS_StyleScopedClasses['mt-0.5']} */ ;
/** @type {__VLS_StyleScopedClasses['space-y-4']} */ ;
/** @type {__VLS_StyleScopedClasses['block']} */ ;
/** @type {__VLS_StyleScopedClasses['text-xs']} */ ;
/** @type {__VLS_StyleScopedClasses['font-medium']} */ ;
/** @type {__VLS_StyleScopedClasses['text-slate-700']} */ ;
/** @type {__VLS_StyleScopedClasses['dark:text-surface-200']} */ ;
/** @type {__VLS_StyleScopedClasses['mb-1.5']} */ ;
/** @type {__VLS_StyleScopedClasses['relative']} */ ;
/** @type {__VLS_StyleScopedClasses['absolute']} */ ;
/** @type {__VLS_StyleScopedClasses['inset-y-0']} */ ;
/** @type {__VLS_StyleScopedClasses['left-0']} */ ;
/** @type {__VLS_StyleScopedClasses['pl-3.5']} */ ;
/** @type {__VLS_StyleScopedClasses['flex']} */ ;
/** @type {__VLS_StyleScopedClasses['items-center']} */ ;
/** @type {__VLS_StyleScopedClasses['pointer-events-none']} */ ;
/** @type {__VLS_StyleScopedClasses['text-slate-400']} */ ;
/** @type {__VLS_StyleScopedClasses['dark:text-surface-500']} */ ;
/** @type {__VLS_StyleScopedClasses['w-4']} */ ;
/** @type {__VLS_StyleScopedClasses['h-4']} */ ;
/** @type {__VLS_StyleScopedClasses['flex']} */ ;
/** @type {__VLS_StyleScopedClasses['items-center']} */ ;
/** @type {__VLS_StyleScopedClasses['gap-1.5']} */ ;
/** @type {__VLS_StyleScopedClasses['text-[11px]']} */ ;
/** @type {__VLS_StyleScopedClasses['text-rose-500']} */ ;
/** @type {__VLS_StyleScopedClasses['mt-1.5']} */ ;
/** @type {__VLS_StyleScopedClasses['w-3.5']} */ ;
/** @type {__VLS_StyleScopedClasses['h-3.5']} */ ;
/** @type {__VLS_StyleScopedClasses['shrink-0']} */ ;
/** @type {__VLS_StyleScopedClasses['block']} */ ;
/** @type {__VLS_StyleScopedClasses['text-xs']} */ ;
/** @type {__VLS_StyleScopedClasses['font-medium']} */ ;
/** @type {__VLS_StyleScopedClasses['text-slate-700']} */ ;
/** @type {__VLS_StyleScopedClasses['dark:text-surface-200']} */ ;
/** @type {__VLS_StyleScopedClasses['mb-1.5']} */ ;
/** @type {__VLS_StyleScopedClasses['relative']} */ ;
/** @type {__VLS_StyleScopedClasses['absolute']} */ ;
/** @type {__VLS_StyleScopedClasses['inset-y-0']} */ ;
/** @type {__VLS_StyleScopedClasses['left-0']} */ ;
/** @type {__VLS_StyleScopedClasses['pl-3.5']} */ ;
/** @type {__VLS_StyleScopedClasses['flex']} */ ;
/** @type {__VLS_StyleScopedClasses['items-center']} */ ;
/** @type {__VLS_StyleScopedClasses['pointer-events-none']} */ ;
/** @type {__VLS_StyleScopedClasses['text-slate-400']} */ ;
/** @type {__VLS_StyleScopedClasses['dark:text-surface-500']} */ ;
/** @type {__VLS_StyleScopedClasses['w-4']} */ ;
/** @type {__VLS_StyleScopedClasses['h-4']} */ ;
/** @type {__VLS_StyleScopedClasses['absolute']} */ ;
/** @type {__VLS_StyleScopedClasses['inset-y-0']} */ ;
/** @type {__VLS_StyleScopedClasses['right-0']} */ ;
/** @type {__VLS_StyleScopedClasses['pr-3']} */ ;
/** @type {__VLS_StyleScopedClasses['flex']} */ ;
/** @type {__VLS_StyleScopedClasses['items-center']} */ ;
/** @type {__VLS_StyleScopedClasses['text-slate-400']} */ ;
/** @type {__VLS_StyleScopedClasses['hover:text-slate-600']} */ ;
/** @type {__VLS_StyleScopedClasses['dark:text-surface-500']} */ ;
/** @type {__VLS_StyleScopedClasses['dark:hover:text-surface-300']} */ ;
/** @type {__VLS_StyleScopedClasses['focus:outline-none']} */ ;
/** @type {__VLS_StyleScopedClasses['transition']} */ ;
/** @type {__VLS_StyleScopedClasses['w-4']} */ ;
/** @type {__VLS_StyleScopedClasses['h-4']} */ ;
/** @type {__VLS_StyleScopedClasses['w-4']} */ ;
/** @type {__VLS_StyleScopedClasses['h-4']} */ ;
/** @type {__VLS_StyleScopedClasses['flex']} */ ;
/** @type {__VLS_StyleScopedClasses['items-center']} */ ;
/** @type {__VLS_StyleScopedClasses['gap-1.5']} */ ;
/** @type {__VLS_StyleScopedClasses['text-[11px]']} */ ;
/** @type {__VLS_StyleScopedClasses['text-rose-500']} */ ;
/** @type {__VLS_StyleScopedClasses['mt-1.5']} */ ;
/** @type {__VLS_StyleScopedClasses['w-3.5']} */ ;
/** @type {__VLS_StyleScopedClasses['h-3.5']} */ ;
/** @type {__VLS_StyleScopedClasses['shrink-0']} */ ;
/** @type {__VLS_StyleScopedClasses['flex']} */ ;
/** @type {__VLS_StyleScopedClasses['items-center']} */ ;
/** @type {__VLS_StyleScopedClasses['justify-between']} */ ;
/** @type {__VLS_StyleScopedClasses['pt-0.5']} */ ;
/** @type {__VLS_StyleScopedClasses['flex']} */ ;
/** @type {__VLS_StyleScopedClasses['items-center']} */ ;
/** @type {__VLS_StyleScopedClasses['gap-2']} */ ;
/** @type {__VLS_StyleScopedClasses['cursor-pointer']} */ ;
/** @type {__VLS_StyleScopedClasses['select-none']} */ ;
/** @type {__VLS_StyleScopedClasses['group']} */ ;
/** @type {__VLS_StyleScopedClasses['rounded']} */ ;
/** @type {__VLS_StyleScopedClasses['border-slate-300']} */ ;
/** @type {__VLS_StyleScopedClasses['dark:border-surface-700']} */ ;
/** @type {__VLS_StyleScopedClasses['bg-white']} */ ;
/** @type {__VLS_StyleScopedClasses['dark:bg-surface-950']} */ ;
/** @type {__VLS_StyleScopedClasses['text-brand-600']} */ ;
/** @type {__VLS_StyleScopedClasses['focus:ring-brand-500/30']} */ ;
/** @type {__VLS_StyleScopedClasses['focus:ring-offset-0']} */ ;
/** @type {__VLS_StyleScopedClasses['w-3.5']} */ ;
/** @type {__VLS_StyleScopedClasses['h-3.5']} */ ;
/** @type {__VLS_StyleScopedClasses['cursor-pointer']} */ ;
/** @type {__VLS_StyleScopedClasses['text-xs']} */ ;
/** @type {__VLS_StyleScopedClasses['text-slate-600']} */ ;
/** @type {__VLS_StyleScopedClasses['dark:text-surface-300']} */ ;
/** @type {__VLS_StyleScopedClasses['group-hover:text-slate-900']} */ ;
/** @type {__VLS_StyleScopedClasses['dark:group-hover:text-white']} */ ;
/** @type {__VLS_StyleScopedClasses['transition']} */ ;
/** @type {__VLS_StyleScopedClasses['w-full']} */ ;
/** @type {__VLS_StyleScopedClasses['mt-2']} */ ;
/** @type {__VLS_StyleScopedClasses['py-2.5']} */ ;
/** @type {__VLS_StyleScopedClasses['px-4']} */ ;
/** @type {__VLS_StyleScopedClasses['rounded-xl']} */ ;
/** @type {__VLS_StyleScopedClasses['bg-gradient-to-r']} */ ;
/** @type {__VLS_StyleScopedClasses['from-brand-600']} */ ;
/** @type {__VLS_StyleScopedClasses['via-indigo-600']} */ ;
/** @type {__VLS_StyleScopedClasses['to-brand-500']} */ ;
/** @type {__VLS_StyleScopedClasses['hover:from-brand-500']} */ ;
/** @type {__VLS_StyleScopedClasses['hover:to-indigo-500']} */ ;
/** @type {__VLS_StyleScopedClasses['text-white']} */ ;
/** @type {__VLS_StyleScopedClasses['text-xs']} */ ;
/** @type {__VLS_StyleScopedClasses['font-semibold']} */ ;
/** @type {__VLS_StyleScopedClasses['shadow-md']} */ ;
/** @type {__VLS_StyleScopedClasses['shadow-brand-500/25']} */ ;
/** @type {__VLS_StyleScopedClasses['active:scale-[0.99]']} */ ;
/** @type {__VLS_StyleScopedClasses['flex']} */ ;
/** @type {__VLS_StyleScopedClasses['items-center']} */ ;
/** @type {__VLS_StyleScopedClasses['justify-center']} */ ;
/** @type {__VLS_StyleScopedClasses['gap-2']} */ ;
/** @type {__VLS_StyleScopedClasses['transition']} */ ;
/** @type {__VLS_StyleScopedClasses['duration-150']} */ ;
/** @type {__VLS_StyleScopedClasses['disabled:opacity-60']} */ ;
/** @type {__VLS_StyleScopedClasses['disabled:cursor-not-allowed']} */ ;
/** @type {__VLS_StyleScopedClasses['cursor-pointer']} */ ;
/** @type {__VLS_StyleScopedClasses['group']} */ ;
/** @type {__VLS_StyleScopedClasses['w-4']} */ ;
/** @type {__VLS_StyleScopedClasses['h-4']} */ ;
/** @type {__VLS_StyleScopedClasses['animate-spin']} */ ;
/** @type {__VLS_StyleScopedClasses['w-3.5']} */ ;
/** @type {__VLS_StyleScopedClasses['h-3.5']} */ ;
/** @type {__VLS_StyleScopedClasses['group-hover:translate-x-0.5']} */ ;
/** @type {__VLS_StyleScopedClasses['transition-transform']} */ ;
/** @type {__VLS_StyleScopedClasses['mt-6']} */ ;
/** @type {__VLS_StyleScopedClasses['pt-5']} */ ;
/** @type {__VLS_StyleScopedClasses['border-t']} */ ;
/** @type {__VLS_StyleScopedClasses['border-slate-200/70']} */ ;
/** @type {__VLS_StyleScopedClasses['dark:border-surface-800/70']} */ ;
/** @type {__VLS_StyleScopedClasses['flex']} */ ;
/** @type {__VLS_StyleScopedClasses['items-center']} */ ;
/** @type {__VLS_StyleScopedClasses['justify-center']} */ ;
/** @type {__VLS_StyleScopedClasses['gap-2']} */ ;
/** @type {__VLS_StyleScopedClasses['text-[11px]']} */ ;
/** @type {__VLS_StyleScopedClasses['text-slate-400']} */ ;
/** @type {__VLS_StyleScopedClasses['dark:text-surface-500']} */ ;
/** @type {__VLS_StyleScopedClasses['w-3.5']} */ ;
/** @type {__VLS_StyleScopedClasses['h-3.5']} */ ;
/** @type {__VLS_StyleScopedClasses['text-emerald-500']} */ ;
/** @type {__VLS_StyleScopedClasses['shrink-0']} */ ;
/** @type {__VLS_StyleScopedClasses['w-full']} */ ;
/** @type {__VLS_StyleScopedClasses['max-w-5xl']} */ ;
/** @type {__VLS_StyleScopedClasses['mx-auto']} */ ;
/** @type {__VLS_StyleScopedClasses['flex']} */ ;
/** @type {__VLS_StyleScopedClasses['flex-col']} */ ;
/** @type {__VLS_StyleScopedClasses['sm:flex-row']} */ ;
/** @type {__VLS_StyleScopedClasses['items-center']} */ ;
/** @type {__VLS_StyleScopedClasses['justify-between']} */ ;
/** @type {__VLS_StyleScopedClasses['gap-2']} */ ;
/** @type {__VLS_StyleScopedClasses['text-center']} */ ;
/** @type {__VLS_StyleScopedClasses['text-[11px]']} */ ;
/** @type {__VLS_StyleScopedClasses['text-slate-400']} */ ;
/** @type {__VLS_StyleScopedClasses['dark:text-surface-500']} */ ;
/** @type {__VLS_StyleScopedClasses['relative']} */ ;
/** @type {__VLS_StyleScopedClasses['z-20']} */ ;
/** @type {__VLS_StyleScopedClasses['py-2']} */ ;
/** @type {__VLS_StyleScopedClasses['flex']} */ ;
/** @type {__VLS_StyleScopedClasses['items-center']} */ ;
/** @type {__VLS_StyleScopedClasses['gap-4']} */ ;
var __VLS_dollars;
const __VLS_self = (await import('vue')).defineComponent({
    setup() {
        return {
            Lock: Lock,
            Mail: Mail,
            Eye: Eye,
            EyeOff: EyeOff,
            ArrowRight: ArrowRight,
            Layers: Layers,
            Sun: Sun,
            Moon: Moon,
            ShieldCheck: ShieldCheck,
            AlertCircle: AlertCircle,
            Loader2: Loader2,
            CheckCircle2: CheckCircle2,
            theme: theme,
            toggleTheme: toggleTheme,
            showPassword: showPassword,
            form: form,
            flash: flash,
            submit: submit,
        };
    },
});
export default (await import('vue')).defineComponent({
    setup() {
        return {};
    },
});
; /* PartiallyEnd: #4569/main.vue */
//# sourceMappingURL=Login.vue.js.map