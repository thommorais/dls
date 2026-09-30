import {
	archiveId,
	episodeId,
	momentId,
	momentTypeId,
	openingId,
	personId,
	songId,
	type Appearance,
	type ArchiveEntry,
	type ArchiveKind,
	type Average,
	type Count,
	type Episode,
	type EpisodeDetail,
	type Gender,
	type Moment,
	type MomentSource,
	type MomentType,
	type Opening,
	type Overview,
	type Person,
	type PersonKind,
	type Rankings,
	type Song,
	type StatRecord,
	type Streak,
	type Tally,
} from '@thom/dls-domain/types'

/**
 * The wire shapes. They are snake_case because the API is, and they live
 * here rather than in the shared package so a rename on either side is one
 * translation to fix.
 */
type WireTally = {
	women_interviewed: number
	openings_aired: number
	moments: Record<string, number>
	total_moments: number
}

type WireMomentType = {
	id: string
	slug: string
	label: string
	description: string
	color: string
	position: number
}

type WireEpisode = {
	id: string
	number: number
	slug: string
	title: string
	youtube_id: string
	published_at: string
	duration_seconds: number
	thumbnail_url: string
	description?: string
	tally: WireTally
}

type WireMoment = {
	id: string
	type: string
	type_label: string
	color: string
	video_timestamp: number
	actor?: string
	summary?: string
	trigger_word?: string
	song?: string
	clip_url?: string
	source?: string
}

type WirePerson = {
	id: string
	slug: string
	name: string
	kind: string
	gender: string
	photo?: string
}

type WireAppearance = { person: string; role: string; is_interview: boolean }

type WireSong = { id: string; title: string; artist?: string; genre?: string; year?: number; link?: string }

type WireOpening = {
	id: string
	title: string
	author_name: string
	author_handle?: string
	genre?: string
	episode?: string
	at_seconds?: number
	aired_at?: string
	media_url?: string
}

type WireArchive = {
	id: string
	slug: string
	title: string
	summary?: string
	body?: string
	kind: string
	tags: string[]
	episode?: string
	person?: string
	moment_type?: string
	source_url?: string
	published_at: string
}

type WireCount = { key: string; label: string; count: number }

export type WireOverview = {
	tally: WireTally
	types: WireMomentType[]
	episodes: number
	first: WireEpisode | null
	latest: WireEpisode | null
}

export type WireEpisodes = { episodes: WireEpisode[] }

export type WireEpisodeDetail = {
	episode: WireEpisode
	moments: WireMoment[]
	appearances: WireAppearance[]
	people: WirePerson[]
	songs: WireSong[]
	openings: WireOpening[]
}

export type WireOpenings = { openings: WireOpening[] }

export type WireArchiveList = { entries: WireArchive[] }

export type WireRankings = {
	types: WireMomentType[]
	songs: WireCount[]
	trigger_words: WireCount[]
	guests: WireCount[]
	actors: Record<string, WireCount[]>
	records: { type: string; count: number; episode: string; title: string }[]
	streaks: { type: string; length: number; from: string; to: string }[]
	averages: { type: string; total: number; episodes: number; per_episode: number }[]
}

/**
 * An absent date is the zero time on the Go side, which arrives as year one
 * rather than as null. Reading it as an invalid Date would print "Invalid
 * Date" in the page, so it becomes the epoch and a caller can test for it.
 */
const toDate = (value: string | undefined): Date => {
	if (!value) return new Date(0)

	const parsed = new Date(value)
	if (Number.isNaN(parsed.getTime()) || parsed.getUTCFullYear() < 1970) return new Date(0)

	return parsed
}

const toOptionalDate = (value: string | undefined): Date | null => {
	if (!value) return null

	const parsed = toDate(value)
	return parsed.getTime() === 0 ? null : parsed
}

const toTally = (wire: WireTally | undefined): Tally => ({
	womenInterviewed: wire?.women_interviewed ?? 0,
	openingsAired: wire?.openings_aired ?? 0,
	moments: wire?.moments ?? {},
	totalMoments: wire?.total_moments ?? 0,
})

export const toMomentType = (wire: WireMomentType): MomentType => ({
	id: momentTypeId(wire.id),
	slug: wire.slug,
	label: wire.label,
	description: wire.description,
	color: wire.color,
	position: wire.position,
})

export const toEpisode = (wire: WireEpisode): Episode => ({
	id: episodeId(wire.id),
	number: wire.number,
	slug: wire.slug,
	title: wire.title,
	youtubeId: wire.youtube_id,
	publishedAt: toDate(wire.published_at),
	durationSeconds: wire.duration_seconds,
	thumbnailUrl: wire.thumbnail_url,
	description: wire.description ?? '',
	tally: toTally(wire.tally),
})

const toMoment = (wire: WireMoment): Moment => ({
	id: momentId(wire.id),
	type: wire.type,
	typeLabel: wire.type_label,
	color: wire.color,
	videoTimestamp: wire.video_timestamp,
	actor: wire.actor ? personId(wire.actor) : null,
	summary: wire.summary ?? '',
	triggerWord: wire.trigger_word ?? '',
	song: wire.song ? songId(wire.song) : null,
	clipUrl: wire.clip_url ?? '',
	source: (wire.source as MomentSource | undefined) ?? null,
})

const toPerson = (wire: WirePerson): Person => ({
	id: personId(wire.id),
	slug: wire.slug,
	name: wire.name,
	kind: wire.kind as PersonKind,
	gender: wire.gender as Gender,
	photo: wire.photo ?? '',
})

const toAppearance = (wire: WireAppearance): Appearance => ({
	person: personId(wire.person),
	role: wire.role as Appearance['role'],
	isInterview: wire.is_interview,
})

const toSong = (wire: WireSong): Song => ({
	id: songId(wire.id),
	title: wire.title,
	artist: wire.artist ?? '',
	genre: wire.genre ?? '',
	year: wire.year ?? 0,
	link: wire.link ?? '',
})

const toOpening = (wire: WireOpening): Opening => ({
	id: openingId(wire.id),
	title: wire.title,
	authorName: wire.author_name,
	authorHandle: wire.author_handle ?? '',
	genre: wire.genre ?? '',
	episode: wire.episode ? episodeId(wire.episode) : null,
	atSeconds: wire.at_seconds ?? 0,
	airedAt: toOptionalDate(wire.aired_at),
	mediaUrl: wire.media_url ?? '',
})

const toArchiveEntry = (wire: WireArchive): ArchiveEntry => ({
	id: archiveId(wire.id),
	slug: wire.slug,
	title: wire.title,
	summary: wire.summary ?? '',
	body: wire.body ?? '',
	kind: wire.kind as ArchiveKind,
	tags: wire.tags ?? [],
	episode: wire.episode ? episodeId(wire.episode) : null,
	person: wire.person ? personId(wire.person) : null,
	momentType: wire.moment_type ? momentTypeId(wire.moment_type) : null,
	sourceUrl: wire.source_url ?? '',
	publishedAt: toDate(wire.published_at),
})

const toCount = (wire: WireCount): Count => ({ key: wire.key, label: wire.label, count: wire.count })

export const toOverview = (wire: WireOverview): Overview => ({
	tally: toTally(wire.tally),
	types: (wire.types ?? []).map(toMomentType),
	episodes: wire.episodes,
	first: wire.first ? toEpisode(wire.first) : null,
	latest: wire.latest ? toEpisode(wire.latest) : null,
})

export const toEpisodes = (wire: WireEpisodes): readonly Episode[] => (wire.episodes ?? []).map(toEpisode)

export const toEpisodeDetail = (wire: WireEpisodeDetail): EpisodeDetail => ({
	episode: toEpisode(wire.episode),
	moments: (wire.moments ?? []).map(toMoment),
	appearances: (wire.appearances ?? []).map(toAppearance),
	people: (wire.people ?? []).map(toPerson),
	songs: (wire.songs ?? []).map(toSong),
	openings: (wire.openings ?? []).map(toOpening),
})

export const toOpenings = (wire: WireOpenings): readonly Opening[] => (wire.openings ?? []).map(toOpening)

export const toArchive = (wire: WireArchiveList): readonly ArchiveEntry[] => (wire.entries ?? []).map(toArchiveEntry)

export const toArchiveOne = toArchiveEntry

export const toRankings = (wire: WireRankings): Rankings => ({
	types: (wire.types ?? []).map(toMomentType),
	songs: (wire.songs ?? []).map(toCount),
	triggerWords: (wire.trigger_words ?? []).map(toCount),
	guests: (wire.guests ?? []).map(toCount),
	actors: Object.fromEntries(Object.entries(wire.actors ?? {}).map(([slug, board]) => [slug, board.map(toCount)])),
	records: (wire.records ?? []).map((record): StatRecord => ({
		type: record.type,
		count: record.count,
		episode: episodeId(record.episode),
		title: record.title,
	})),
	streaks: (wire.streaks ?? []).map((streak): Streak => ({
		type: streak.type,
		length: streak.length,
		from: episodeId(streak.from),
		to: episodeId(streak.to),
	})),
	averages: (wire.averages ?? []).map((average): Average => ({
		type: average.type,
		total: average.total,
		episodes: average.episodes,
		perEpisode: average.per_episode,
	})),
})
