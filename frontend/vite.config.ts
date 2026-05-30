import tailwindcss from '@tailwindcss/vite';
import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';

export default defineConfig({
	plugins: [tailwindcss(), sveltekit()],
	server: {
		host: '0.0.0.0',
		port: 4173,
		proxy: {
			'/api': {
				target: 'http://backend:8080',
				changeOrigin: true
			},
			'/health': {
				target: 'http://backend:8080',
				changeOrigin: true
			}
		}
	},
	optimizeDeps: {
		include: ['svelte-sonner']
	},
	ssr: {
		noExternal: ['sveltekit-superforms']
	}
});
