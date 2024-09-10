import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';

export default defineConfig({
    server: {
        headers: {
            'Cross-Origin-Opener-Policy': 'same-origin',
            'Cross-Origin-Embedder-Policy': 'require-corp',
        },
        hmr: false,
    },
    esbuild: {
        target: 'esnext'
    },
    plugins: [react()],
    root: 'src/',
    base: './',
    build: {
        // outDir: '../build/',
        outDir: '../dist/',
        emptyOutDir: true,
    },
});
