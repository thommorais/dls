/**
 * The API base. Empty means same origin, which is how it is served in
 * production: one Go binary answers /api/dls and serves the built site.
 */
const configured = import.meta.env.VITE_API_URL ?? 'http://127.0.0.1:8090'

export const ENVS = {
	API_URL: configured === '' ? '/' : configured,
	IS_DEV: import.meta.env.DEV,
} as const
