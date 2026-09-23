import { describe, expect, it } from 'vitest'
import { cumulative, episodeSeries, tallySeries } from './series'
import { episodeId, momentTypeId, type Episode, type MomentType, type Tally } from './types'

const type = (slug: string, label: string, position: number): MomentType => ({
	id: momentTypeId(`t-${slug}`),
	slug,
	label,
	description: '',
	color: '#6366F1',
	position,
})

const types = [type('sing', 'Saiu cantando', 1), type('question', 'Pergunta', 2), type('dance', 'Dançou', 3)]

const tally = (moments: Record<string, number>): Tally => ({
	womenInterviewed: 0,
	openingsAired: 0,
	moments,
	totalMoments: Object.values(moments).reduce((sum, n) => sum + n, 0),
})

const episode = (n: number, moments: Record<string, number>): Episode => ({
	id: episodeId(`e${n}`),
	number: n,
	slug: `ep-${n}`,
	title: `Episódio ${n}`,
	youtubeId: `yt${n}`,
	publishedAt: new Date(Date.UTC(2025, 0, n)),
	durationSeconds: 3600,
	thumbnailUrl: '',
	description: '',
	tally: tally(moments),
})

describe('tallySeries', () => {
	it('follows the order the types came in', () => {
		const series = tallySeries(tally({ question: 3, sing: 5 }), types)

		expect(series.map(point => point.slug)).toEqual(['sing', 'question', 'dance'])
	})

	// A counter the site advertises should not vanish because nobody has done
	// it yet.
	it('keeps a type that has never happened, at zero', () => {
		const series = tallySeries(tally({ sing: 5 }), types)

		expect(series.find(point => point.slug === 'dance')).toEqual({
			slug: 'dance',
			label: 'Dançou',
			color: '#6366F1',
			count: 0,
		})
	})

	// It would have no label and no colour, so there is nothing to draw.
	it('ignores a count whose type is not in the list', () => {
		const series = tallySeries(tally({ sing: 5, deleted: 9 }), types)

		expect(series.map(point => point.slug)).not.toContain('deleted')
	})
})

describe('episodeSeries', () => {
	it('gives one point per episode, silent ones included', () => {
		const series = episodeSeries([episode(1, { sing: 2 }), episode(2, {}), episode(3, { sing: 1 })], 'sing')

		expect(series.map(point => point.value)).toEqual([2, 0, 1])
		expect(series.map(point => point.label)).toEqual(['ep-1', 'ep-2', 'ep-3'])
	})

	it('reads in publish order however the episodes arrived', () => {
		const series = episodeSeries([episode(3, { sing: 1 }), episode(1, { sing: 2 })], 'sing')

		expect(series.map(point => point.label)).toEqual(['ep-1', 'ep-3'])
	})

	it('has nothing to plot for no episodes', () => {
		expect(episodeSeries([], 'sing')).toEqual([])
	})
})

describe('cumulative', () => {
	it('turns a per-episode count into a running total', () => {
		const series = episodeSeries([episode(1, { sing: 2 }), episode(2, {}), episode(3, { sing: 3 })], 'sing')

		expect(cumulative(series).map(point => point.value)).toEqual([2, 2, 5])
	})

	it('keeps the labels it was given', () => {
		const series = episodeSeries([episode(1, { sing: 2 })], 'sing')

		expect(cumulative(series)[0]?.label).toBe('ep-1')
	})
})
