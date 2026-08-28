<!DOCTYPE html>
<html lang="{{ str_replace('_', '-', app()->getLocale()) }}">
    <head>
        <meta charset="utf-8">
        <meta name="viewport" content="width=device-width, initial-scale=1">
        <title inertia>{{ config('app.name', 'Kodepreneur Panel') }}</title>

        <!-- Initialize Theme (Dark/Light) immediately before render -->
        <script>
            (function() {
                try {
                    const savedTheme = localStorage.getItem('kp_theme');
                    if (savedTheme === 'light') {
                        document.documentElement.classList.remove('dark');
                    } else if (savedTheme === 'dark') {
                        document.documentElement.classList.add('dark');
                    } else if (window.matchMedia('(prefers-color-scheme: dark)').matches) {
                        document.documentElement.classList.add('dark');
                    } else {
                        // Default to dark mode for server control panels
                        document.documentElement.classList.add('dark');
                    }
                } catch (e) {
                    document.documentElement.classList.add('dark');
                }
            })();
        </script>

        <!-- Fonts -->
        <link rel="preconnect" href="https://fonts.googleapis.com">
        <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
        <link href="https://fonts.googleapis.com/css2?family=Inter:wght@300;400;500;600;700&family=JetBrains+Mono:wght@400;500;600&display=swap" rel="stylesheet">

        <!-- Scripts & Styles -->
        @vite(['resources/css/app.css', 'resources/js/app.ts'])
        @inertiaHead
    </head>
    <body class="font-sans antialiased bg-slate-50 text-slate-900 dark:bg-[#070b14] dark:text-slate-100 min-h-screen transition-colors duration-150">
        @inertia
    </body>
</html>
