import { defineConfig } from 'vite';
import { resolve } from 'path';

export default defineConfig({
    build: {
        outDir: 'static/dist',
        rollupOptions: {
            input: {
                app: resolve(__dirname, 'static/js/app.js')
            },
            output: {
                entryFileNames: 'js/[name].js',
                chunkFileNames: 'js/[name].js',
                assetFileNames: (assetInfo) => {
                    if (/\.(png|jpe?g|gif|svg|ico)$/i.test(assetInfo.name)) {
                        return 'images/[name][extname]';
                    }
                    return 'css/[name][extname]';
                }
            }
        }
    },
    css: {
        plugins: [
            require('tailwindcss'),
            require('autoprefixer')
        ]
    }
});