import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';

export default defineConfig({
	plugins: [sveltekit()],
	server: {
		watch: {
			usePolling: true,
		},
		host: "0.0.0.0",
		port: 80,
	}
});
