import type { StatsFilter } from '@thom/dls-domain/filter'
import type { ArchiveEntry, Episode, EpisodeDetail, Opening, Overview, Rankings } from '@thom/dls-domain/types'
import type { Safe } from '@thom/safe-return'
import type { ErrorResponse } from '@thom/try-catch'

type Read<T> = Promise<Safe<T, ErrorResponse>>

/**
 * Everything the site reads. The counters are computed in Go, so this port
 * is one call per page rather than a query builder.
 */
export type StatsPort = {
	readonly overview: (filter?: StatsFilter) => Read<Overview>
	readonly episodes: (filter?: StatsFilter) => Read<readonly Episode[]>
	readonly episode: (slug: string) => Read<EpisodeDetail>
	readonly openings: (filter?: StatsFilter) => Read<readonly Opening[]>
	readonly rankings: (filter?: StatsFilter) => Read<Rankings>
	readonly archive: (filter?: StatsFilter) => Read<readonly ArchiveEntry[]>
	readonly archiveEntry: (slug: string) => Read<ArchiveEntry>
}
