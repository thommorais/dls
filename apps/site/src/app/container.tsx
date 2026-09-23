import { createStatsAdapter } from '_/adapters/pocketbase/stats-adapter'
import type { StatsPort } from '_/core/ports/stats'

/**
 * The container is the only place that knows which adapter implements which
 * port. A test swaps the whole thing for fakes; a component never imports an
 * adapter.
 */
export type Container = {
	readonly stats: StatsPort
}

export const createContainer = (): Container => ({
	stats: createStatsAdapter(),
})
