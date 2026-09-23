import { describe, expect, it } from 'vitest'
import { formatPerEpisode, formatTimestamp, youtubeLink } from './format'

describe('formatTimestamp', () => {
	it('reads as a clock', () => {
		expect(formatTimestamp(0)).toBe('0:00')
		expect(formatTimestamp(7)).toBe('0:07')
		expect(formatTimestamp(61)).toBe('1:01')
		expect(formatTimestamp(600)).toBe('10:00')
	})

	it('adds the hour only once there is one', () => {
		expect(formatTimestamp(3599)).toBe('59:59')
		expect(formatTimestamp(3600)).toBe('1:00:00')
		expect(formatTimestamp(3661)).toBe('1:01:01')
	})

	// A timestamp comes from a form, so it can be anything.
	it('reads a nonsense value as the start of the video', () => {
		expect(formatTimestamp(-5)).toBe('0:00')
		expect(formatTimestamp(Number.NaN)).toBe('0:00')
		expect(formatTimestamp(Number.POSITIVE_INFINITY)).toBe('0:00')
	})

	it('drops a fractional second rather than rounding up past the minute', () => {
		expect(formatTimestamp(59.9)).toBe('0:59')
	})
})

describe('youtubeLink', () => {
	it('deep links to the second', () => {
		expect(youtubeLink('abc123', 90)).toBe('https://www.youtube.com/watch?v=abc123&t=90s')
	})

	it('links to the video when there is no timestamp worth naming', () => {
		expect(youtubeLink('abc123')).toBe('https://www.youtube.com/watch?v=abc123')
		expect(youtubeLink('abc123', 0)).toBe('https://www.youtube.com/watch?v=abc123')
	})

	// An episode with no youtube id yet must not render a link to nowhere.
	it('gives nothing back when there is no video', () => {
		expect(youtubeLink('')).toBe('')
		expect(youtubeLink('', 90)).toBe('')
	})
})

describe('formatPerEpisode', () => {
	it('keeps two decimals, which is how a rate reads', () => {
		expect(formatPerEpisode(4.7234)).toBe('4.72')
		expect(formatPerEpisode(0)).toBe('0.00')
	})

	it('reads a nonsense value as zero', () => {
		expect(formatPerEpisode(Number.NaN)).toBe('0.00')
	})
})
