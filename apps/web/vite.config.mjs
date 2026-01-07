import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';
import path from 'path';

export default defineConfig({
    root: './src/app',
    envDir: path.resolve(__dirname),
    publicDir: path.resolve(__dirname, 'public'),
    plugins: [react()],
    base: './',
    optimizeDeps: {
        exclude: ['@ffmpeg/ffmpeg', '@ffmpeg/util'],
    },
    build: {
        outDir: '../../dist',
        emptyOutDir: true,
        minify: 'terser',
        target: 'esnext',
        sourcemap: false,
        rollupOptions: {
            input: {
                main: path.resolve(__dirname, './src/app/index.html'),
            },
        },
    },
    resolve: {
        alias: {
            '@': path.resolve(__dirname, './src/'),
        },
    },
    server: {
        proxy: {
            '/api': {
                target: 'http://localhost:8080',
                changeOrigin: true,
                rewrite: (path) => path.replace(/^\/api/, ''),
            },
        },
    },
});
