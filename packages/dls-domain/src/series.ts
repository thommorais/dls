import type { Episode, MomentType, Tally } from './types'

/** A counter row, and the bars of the chart under it. */
export type TallyPoint = {
	readonly slug: string
	readonly label: string
	readonly color: string
	readonly count: number
}

/** One episode's worth of a single stat, for a line across the season. */
export type EpisodePoint = {
	readonly slug: string
	readonly label: string
	readonly at: Date
	readonly value: number
}

/**
 * tallySeries walks the types rather than the counts, so the order is the one
 * the API chose and a type nobody has done yet still gets a row. A count
 * whose type is missing is dropped: it has no label and no colour, so there
 * is nothing to draw.
 */
export const tallySeries = (tally: Tally, types: readonly MomentType[]): readonly TallyPoint[] =>
	types.map(type => ({
		slug: type.slug,
		label: type.label,
		color: type.color,
		count: tally.moments[type.slug] ?? 0,
	}))

/**
 * episodeSeries is one point per episode in publish order, silent episodes
 * included: a line that skipped them would compress the gaps and lie about
 * the pace.
 */
export const episodeSeries = (episodes: readonly Episode[], slug: string): readonly EpisodePoint[] =>
	[...episodes]
		.sort((a, b) => a.publishedAt.getTime() - b.publishedAt.getTime())
		.map(episode => ({
			slug: episode.slug,
			label: episode.slug,
			at: episode.publishedAt,
			value: episode.tally.moments[slug] ?? 0,
		}))

/** cumulative turns a per-episode count into the running total. */
export const cumulative = (points: readonly EpisodePoint[]): readonly EpisodePoint[] => {
	let total = 0

	return points.map(point => {
		total += point.value
		return { ...point, value: total }
	})
}
