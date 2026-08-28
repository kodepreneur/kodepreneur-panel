import { ref, onMounted } from 'vue';
const theme = ref('dark');
export function useTheme() {
    function initTheme() {
        if (typeof window === 'undefined')
            return;
        const saved = localStorage.getItem('kp_theme');
        if (saved === 'light' || saved === 'dark') {
            setTheme(saved);
        }
        else if (window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches) {
            setTheme('dark');
        }
        else {
            setTheme('dark');
        }
    }
    function setTheme(newTheme) {
        theme.value = newTheme;
        if (typeof document !== 'undefined') {
            if (newTheme === 'dark') {
                document.documentElement.classList.add('dark');
            }
            else {
                document.documentElement.classList.remove('dark');
            }
            try {
                localStorage.setItem('kp_theme', newTheme);
            }
            catch (e) {
                // Ignore storage error
            }
        }
    }
    function toggleTheme() {
        setTheme(theme.value === 'dark' ? 'light' : 'dark');
    }
    onMounted(() => {
        initTheme();
    });
    return {
        theme,
        setTheme,
        toggleTheme,
        initTheme,
    };
}
//# sourceMappingURL=useTheme.js.map