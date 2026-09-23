import type { StatsFilter } from '@thom/dls-domain/filter'
import { queryOptions } from '@tanstack/react-query'
import type { Container } from './container'

/**
 * unwrap turns the adapter's Result into what TanStack Query expects. The
 * ports answer with values so a component never sees a throw; the query layer
 * is the one place that wants a rejection, because that is how it tracks a
 * failed fetch.
 */
const unwrap = async <T>(
	read: Promise<{ success: true; data: T } | { success: false; error: { message: string } }>,
) => {
	const result = await read

	if (!result.success) throw new Error(result.error.message)

	return result.data
}

export const overviewQuery = (container: Container, filter: StatsFilter = {}) =>
	queryOptions({
		queryKey: ['overview', filter],
		queryFn: () => unwrap(container.stats.overview(filter)),
	})

export const episodesQuery = (container: Container, filter: StatsFilter = {}) =>
	queryOptions({
		queryKey: ['episodes', filter],
		queryFn: () => unwrap(container.stats.episodes(filter)),
	})

export const rankingsQuery = (container: Container, filter: StatsFilter = {}) =>
	queryOptions({
		queryKey: ['rankings', filter],
		queryFn: () => unwrap(container.stats.rankings(filter)),
	})

export const openingsQuery = (container: Container, filter: StatsFilter = {}) =>
	queryOptions({
		queryKey: ['openings', filter],
		queryFn: () => unwrap(container.stats.openings(filter)),
	})

export const archiveQuery = (container: Container, filter: StatsFilter = {}) =>
	queryOptions({
		queryKey: ['archive', filter],
		queryFn: () => unwrap(container.stats.archive(filter)),
	})

export const episodeQuery = (container: Container, slug: string) =>
	queryOptions({
		queryKey: ['episode', slug],
		queryFn: () => unwrap(container.stats.episode(slug)),
	})

export const archiveEntryQuery = (container: Container, slug: string) =>
	queryOptions({
		queryKey: ['archive-entry', slug],
		queryFn: () => unwrap(container.stats.archiveEntry(slug)),
	})
