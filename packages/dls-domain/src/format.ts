/** formatTimestamp renders seconds into a video the way a player shows them. */
export const formatTimestamp = (seconds: number): string => {
	if (!Number.isFinite(seconds) || seconds <= 0) return '0:00'

	// Truncated, not rounded: 59.9 seconds is still the fifty-ninth second,
	// and rounding it up would print a minute that never happened.
	const total = Math.floor(seconds)
	const hours = Math.floor(total / 3600)
	const minutes = Math.floor((total % 3600) / 60)
	const rest = total % 60

	const pad = (value: number) => String(value).padStart(2, '0')

	return hours > 0 ? `${hours}:${pad(minutes)}:${pad(rest)}` : `${minutes}:${pad(rest)}`
}

/**
 * youtubeLink deep links to the moment. An episode with no video id yields
 * nothing, so a caller renders text rather than a link to nowhere.
 */
export const youtubeLink = (youtubeId: string, seconds?: number): string => {
	if (!youtubeId) return ''

	const base = `https://www.youtube.com/watch?v=${youtubeId}`
	if (!seconds || !Number.isFinite(seconds) || seconds <= 0) return base

	return `${base}&t=${Math.floor(seconds)}s`
}

/** formatPerEpisode renders a rate, which only reads as one with decimals. */
export const formatPerEpisode = (value: number): string => (Number.isFinite(value) ? value.toFixed(2) : '0.00')
