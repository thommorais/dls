import react from '@vitejs/plugin-react'
import { resolve } from 'node:path'
import { defineConfig } from 'vitest/config'

// The test config is deliberately separate from vite.config.ts: unit tests
// have no use for the Start plugin, and loading it means resolving the router
// entry before a single test can run.
// oxlint-disable-next-line import/no-default-export
export default defineConfig({
	plugins: [react()],
	resolve: {
		alias: {
			_: resolve(import.meta.dirname, './src'),
		},
	},
	test: {
		environment: 'happy-dom',
		include: ['src/**/*.test.{ts,tsx}'],
	},
})
