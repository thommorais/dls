import { ARCHIVE_KINDS, isMember, type ArchiveKind } from './types'

/**
 * The parameters every list route shares. They live in the URL, so a filtered
 * page can be bookmarked and shared, and the same object is what the adapter
 * turns into a query string.
 */
export type StatsFilter = {
	/** Bounds the episode publish date, as YYYY-MM-DD. */
	readonly from?: string
	readonly to?: string
	/** Narrows the openings library only. */
	readonly genre?: string
	/** Narrows the archive only. */
	readonly kind?: ArchiveKind
	/** Caps each ranking board. */
	readonly limit?: number
}

const DATE = /^\d{4}-\d{2}-\d{2}$/

/**
 * parseDate refuses anything that is not a real calendar day. The pattern
 * alone would accept 2025-13-45, which reads as a valid filter and returns an
 * empty page.
 */
const parseDate = (value: unknown): string | undefined => {
	if (typeof value !== 'string' || !DATE.test(value)) return undefined

	const parsed = new Date(`${value}T00:00:00Z`)
	if (Number.isNaN(parsed.getTime())) return undefined

	return parsed.toISOString().slice(0, 10) === value ? value : undefined
}

const parseLimit = (value: unknown): number | undefined => {
	const parsed = typeof value === 'number' ? value : Number(value)
	if (!Number.isFinite(parsed) || parsed <= 0) return undefined

	return Math.floor(parsed)
}

const parseText = (value: unknown): string | undefined => {
	if (typeof value !== 'string') return undefined

	const trimmed = value.trim()
	return trimmed === '' ? undefined : trimmed
}

/**
 * parseFilter reads whatever is in the URL and keeps only what it
 * understands. A stale or hand-typed parameter is dropped rather than
 * rejected: the page should still render, which is the same rule the Go
 * handler follows.
 */
export const parseFilter = (search: Record<string, unknown>): StatsFilter => {
	const filter: {
		from?: string
		to?: string
		genre?: string
		kind?: ArchiveKind
		limit?: number
	} = {}

	const from = parseDate(search.from)
	if (from) filter.from = from

	const to = parseDate(search.to)
	if (to) filter.to = to

	const genre = parseText(search.genre)
	if (genre) filter.genre = genre

	if (isMember(ARCHIVE_KINDS, search.kind)) filter.kind = search.kind

	const limit = parseLimit(search.limit)
	if (limit) filter.limit = limit

	return filter
}

/** filterToQuery writes the filter back out for a request or a URL. */
export const filterToQuery = (filter: StatsFilter): Record<string, string> => {
	const query: Record<string, string> = {}

	if (filter.from) query.from = filter.from
	if (filter.to) query.to = filter.to
	if (filter.genre) query.genre = filter.genre
	if (filter.kind) query.kind = filter.kind
	if (filter.limit) query.limit = String(filter.limit)

	return query
}

export const isEmptyFilter = (filter: StatsFilter): boolean => Object.keys(filterToQuery(filter)).length === 0
