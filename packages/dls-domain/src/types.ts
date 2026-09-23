/**
 * The vocabulary the site shares with the Go API. These are the domain
 * shapes, in camelCase; the wire shapes and their mapping belong to the
 * adapter, so a rename on either side is one translation to fix.
 */

// One shared symbol: a per-file `declare const brand` would make each file's
// Branded<string, 'EpisodeId'> a different type, so ids could not cross
// modules.
declare const brand: unique symbol

export type Branded<T, B extends string> = T & { readonly [brand]: B }

export type EpisodeId = Branded<string, 'EpisodeId'>
export type PersonId = Branded<string, 'PersonId'>
export type MomentId = Branded<string, 'MomentId'>
export type MomentTypeId = Branded<string, 'MomentTypeId'>
export type SongId = Branded<string, 'SongId'>
export type OpeningId = Branded<string, 'OpeningId'>
export type ArchiveId = Branded<string, 'ArchiveId'>

export const episodeId = (value: string): EpisodeId => value as EpisodeId
export const personId = (value: string): PersonId => value as PersonId
export const momentId = (value: string): MomentId => value as MomentId
export const momentTypeId = (value: string): MomentTypeId => value as MomentTypeId
export const songId = (value: string): SongId => value as SongId
export const openingId = (value: string): OpeningId => value as OpeningId
export const archiveId = (value: string): ArchiveId => value as ArchiveId

/**
 * The closed sets. Moment types are deliberately absent: they are rows in a
 * collection, so their slugs cannot be enumerated at build time.
 */
export const ARCHIVE_KINDS = ['about', 'glossary', 'segment', 'trivia', 'bio', 'milestone'] as const
export const OPENING_STATUSES = ['received', 'aired', 'archived'] as const
export const PERSON_KINDS = ['host', 'guest', 'staff'] as const
export const GENDERS = ['woman', 'man', 'nonbinary', 'unknown'] as const
export const APPEARANCE_ROLES = ['host', 'guest', 'remote'] as const
export const MOMENT_SOURCES = ['manual', 'llm', 'algo'] as const

export type ArchiveKind = (typeof ARCHIVE_KINDS)[number]
export type OpeningStatus = (typeof OPENING_STATUSES)[number]
export type PersonKind = (typeof PERSON_KINDS)[number]
export type Gender = (typeof GENDERS)[number]
export type AppearanceRole = (typeof APPEARANCE_ROLES)[number]
export type MomentSource = (typeof MOMENT_SOURCES)[number]

/**
 * isMember turns any of those tuples into a runtime guard, which is what
 * lets a route validate a search param against the same list the types come
 * from instead of a second copy of it.
 */
export const isMember = <T extends readonly string[]>(values: T, value: unknown): value is T[number] =>
	typeof value === 'string' && (values as readonly string[]).includes(value)

export type Tally = {
	readonly womenInterviewed: number
	readonly openingsAired: number
	/** Counted by moment type slug, because the types are rows. */
	readonly moments: Readonly<Record<string, number>>
	readonly totalMoments: number
}

export type MomentType = {
	readonly id: MomentTypeId
	readonly slug: string
	readonly label: string
	readonly description: string
	readonly color: string
	readonly position: number
}

export type Episode = {
	readonly id: EpisodeId
	readonly number: number
	readonly slug: string
	readonly title: string
	readonly youtubeId: string
	readonly publishedAt: Date
	readonly durationSeconds: number
	readonly thumbnailUrl: string
	readonly description: string
	readonly tally: Tally
}

export type Moment = {
	readonly id: MomentId
	readonly type: string
	readonly typeLabel: string
	readonly color: string
	readonly videoTimestamp: number
	readonly actor: PersonId | null
	readonly summary: string
	readonly triggerWord: string
	readonly song: SongId | null
	readonly clipUrl: string
	readonly source: MomentSource | null
}

export type Person = {
	readonly id: PersonId
	readonly slug: string
	readonly name: string
	readonly kind: PersonKind
	readonly gender: Gender
	readonly photo: string
}

export type Appearance = {
	readonly person: PersonId
	readonly role: AppearanceRole
	readonly isInterview: boolean
}

export type Song = {
	readonly id: SongId
	readonly title: string
	readonly artist: string
	readonly genre: string
	readonly year: number
	readonly link: string
}

export type Opening = {
	readonly id: OpeningId
	readonly title: string
	readonly authorName: string
	readonly authorHandle: string
	readonly genre: string
	readonly episode: EpisodeId | null
	readonly atSeconds: number
	readonly sentAt: Date
	readonly airedAt: Date | null
	readonly mediaUrl: string
	readonly status: OpeningStatus
}

export type ArchiveEntry = {
	readonly id: ArchiveId
	readonly slug: string
	readonly title: string
	readonly summary: string
	/** Empty in an index listing, which never ships the bodies. */
	readonly body: string
	readonly kind: ArchiveKind
	readonly tags: readonly string[]
	readonly episode: EpisodeId | null
	readonly person: PersonId | null
	readonly momentType: MomentTypeId | null
	readonly sourceUrl: string
	readonly publishedAt: Date
}

export type Count = {
	readonly key: string
	readonly label: string
	readonly count: number
}

/** Named StatRecord because Record is taken by the language. */
export type StatRecord = {
	readonly type: string
	readonly count: number
	readonly episode: EpisodeId
	readonly title: string
}

export type Streak = {
	readonly type: string
	readonly length: number
	readonly from: EpisodeId
	readonly to: EpisodeId
}

export type Average = {
	readonly type: string
	readonly total: number
	readonly episodes: number
	readonly perEpisode: number
}

export type Overview = {
	readonly tally: Tally
	readonly types: readonly MomentType[]
	readonly episodes: number
	readonly first: Episode | null
	readonly latest: Episode | null
}

export type EpisodeDetail = {
	readonly episode: Episode
	readonly moments: readonly Moment[]
	readonly appearances: readonly Appearance[]
	readonly people: readonly Person[]
	readonly songs: readonly Song[]
	readonly openings: readonly Opening[]
}

export type Rankings = {
	readonly types: readonly MomentType[]
	readonly songs: readonly Count[]
	readonly triggerWords: readonly Count[]
	readonly guests: readonly Count[]
	/** One leaderboard per moment type slug. */
	readonly actors: Readonly<Record<string, readonly Count[]>>
	readonly records: readonly StatRecord[]
	readonly streaks: readonly Streak[]
	readonly averages: readonly Average[]
}
