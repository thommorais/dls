import tailwindcss from '@tailwindcss/vite'
import { tanstackStart } from '@tanstack/react-start/plugin/vite'
import react from '@vitejs/plugin-react'
import { resolve } from 'node:path'
import { defineConfig } from 'vite'

// oxlint-disable-next-line import/no-default-export
export default defineConfig({
	plugins: [
		// Start wraps the same file-based router plugin folio uses, and adds
		// the prerender pass: the loaders below run at build time and each
		// route is written out as HTML.
		tanstackStart({
			prerender: {
				enabled: true,
				crawlLinks: true,
				concurrency: 8,
				// Canonical pages only. A filter link is a client navigation,
				// not a page of its own, and crawling them would write one
				// static file per genre for markup the browser already has.
				filter: ({ path }) => !path.includes('?'),
			},
		}),
		react(),
		tailwindcss(),
	],
	build: {
		sourcemap: true,
		rolldownOptions: {
			output: {
				// The framework core stays cached across app deploys; route
				// chunks are already split on demand by the router plugin.
				advancedChunks: {
					groups: [
						{
							name: 'react-vendor',
							test: /[\\/]node_modules[\\/](react|react-dom|scheduler|@tanstack[\\/]react-router)[\\/]/,
						},
					],
				},
			},
		},
	},
	resolve: {
		alias: {
			_: resolve(import.meta.dirname, './src'),
		},
	},
})
