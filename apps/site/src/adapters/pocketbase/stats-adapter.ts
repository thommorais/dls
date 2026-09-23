import { filterToQuery, type StatsFilter } from '@thom/dls-domain/filter'
import { err, ok, type Safe } from '@thom/safe-return'
import { tryCatch, type ErrorResponse } from '@thom/try-catch'
import type { StatsPort } from '_/core/ports/stats'
import { getPocketBaseClient } from './client'
import {
	toArchive,
	toArchiveOne,
	toEpisodeDetail,
	toEpisodes,
	toOpenings,
	toOverview,
	toRankings,
	type WireArchiveList,
	type WireEpisodeDetail,
	type WireEpisodes,
	type WireOpenings,
	type WireOverview,
	type WireRankings,
} from './mappers'

const failure = (what: string, error: ErrorResponse): Safe<never, ErrorResponse> =>
	err({ ...error, message: `${what}: ${error.message}` })

/**
 * The counters are computed in Go, so every read here goes to /api/dls
 * through the SDK's send rather than to the record API. Reading the
 * collections directly would mean counting in the browser, which is the
 * thing the custom endpoints exist to avoid.
 */
export const createStatsAdapter = (): StatsPort => {
	const client = getPocketBaseClient()

	const read = async <W, T>(
		path: string,
		toDomain: (wire: W) => T,
		what: string,
		filter?: StatsFilter,
	): Promise<Safe<T, ErrorResponse>> => {
		// The mapping runs inside tryCatch on purpose. A payload that does not
		// have the shape claimed here is a failed read like any other, and a
		// port that promises a value must not throw one.
		const result = await tryCatch(
			client.send<W>(`/api/dls${path}`, { method: 'GET', query: filterToQuery(filter ?? {}) }).then(toDomain),
		)

		return result.success ? ok(result.data) : failure(what, result.error)
	}

	return {
		overview: filter =>
			read<WireOverview, ReturnType<typeof toOverview>>(
				'/overview',
				toOverview,
				'Não foi possível carregar os números',
				filter,
			),

		episodes: filter =>
			read<WireEpisodes, ReturnType<typeof toEpisodes>>(
				'/episodes',
				toEpisodes,
				'Não foi possível carregar os episódios',
				filter,
			),

		episode: slug =>
			read<WireEpisodeDetail, ReturnType<typeof toEpisodeDetail>>(
				`/episodes/${encodeURIComponent(slug)}`,
				toEpisodeDetail,
				'Não foi possível carregar o episódio',
			),

		openings: filter =>
			read<WireOpenings, ReturnType<typeof toOpenings>>(
				'/openings',
				toOpenings,
				'Não foi possível carregar as aberturas',
				filter,
			),

		rankings: filter =>
			read<WireRankings, ReturnType<typeof toRankings>>(
				'/rankings',
				toRankings,
				'Não foi possível carregar os rankings',
				filter,
			),

		archive: filter =>
			read<WireArchiveList, ReturnType<typeof toArchive>>(
				'/archive',
				toArchive,
				'Não foi possível carregar o acervo',
				filter,
			),

		archiveEntry: slug =>
			read<Parameters<typeof toArchiveOne>[0], ReturnType<typeof toArchiveOne>>(
				`/archive/${encodeURIComponent(slug)}`,
				toArchiveOne,
				'Não foi possível carregar a página do acervo',
			),
	}
}
