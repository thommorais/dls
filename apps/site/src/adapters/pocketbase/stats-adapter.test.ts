import { beforeEach, describe, expect, it, vi } from 'vitest'
import { createStatsAdapter } from './stats-adapter'

const jsonResponse = (body: unknown, status = 200) =>
	new Response(JSON.stringify(body), { status, headers: { 'Content-Type': 'application/json' } })

const fetchMock = vi.fn()

beforeEach(() => {
	fetchMock.mockReset()
	vi.stubGlobal('fetch', fetchMock)
})

const urlOf = (call: number) => String((fetchMock.mock.calls[call] as [string, RequestInit])[0])

const overviewBody = {
	tally: { women_interviewed: 20, openings_aired: 30, moments: { sing: 109 }, total_moments: 514 },
	types: [{ id: 't1', slug: 'sing', label: 'Saiu cantando', description: '', color: '#10B981', position: 1 }],
	episodes: 40,
	first: null,
	latest: null,
}

describe('createStatsAdapter', () => {
	// Counting happens in Go. A read that went to the record API would be the
	// browser recounting five hundred rows it should never have downloaded.
	it('reads the computed endpoint rather than the record API', async () => {
		fetchMock.mockResolvedValue(jsonResponse(overviewBody))

		await createStatsAdapter().overview()

		expect(urlOf(0)).toContain('/api/dls/overview')
		expect(urlOf(0)).not.toContain('/api/collections')
	})

	it('carries the filter as a query string', async () => {
		fetchMock.mockResolvedValue(jsonResponse(overviewBody))

		await createStatsAdapter().overview({ from: '2025-07-01', limit: 5 })

		expect(urlOf(0)).toContain('from=2025-07-01')
		expect(urlOf(0)).toContain('limit=5')
	})

	it('sends no query at all when nothing is filtered', async () => {
		fetchMock.mockResolvedValue(jsonResponse(overviewBody))

		await createStatsAdapter().overview()

		expect(urlOf(0)).not.toContain('?')
	})

	it('maps the wire shape into the domain', async () => {
		fetchMock.mockResolvedValue(jsonResponse(overviewBody))

		const result = await createStatsAdapter().overview()

		expect(result.success).toBe(true)
		if (!result.success) return
		expect(result.data.tally.womenInterviewed).toBe(20)
		expect(result.data.tally.moments.sing).toBe(109)
		expect(result.data.types[0]?.label).toBe('Saiu cantando')
	})

	it('escapes a slug rather than pasting it into the path', async () => {
		fetchMock.mockResolvedValue(jsonResponse({ episode: { tally: {} }, moments: [] }))

		await createStatsAdapter().episode('ep 1/2')

		expect(urlOf(0)).toContain('ep%201%2F2')
	})

	// A board the API has not filled in yet is an empty board, not a broken
	// page: every list in the payload is optional on purpose.
	it('reads a payload with missing lists as empty ones', async () => {
		fetchMock.mockResolvedValue(jsonResponse({ nothing: 'useful' }))

		const result = await createStatsAdapter().rankings()

		expect(result.success).toBe(true)
		if (!result.success) return
		expect(result.data.songs).toEqual([])
		expect(result.data.types).toEqual([])
	})

	// What is not optional is the subject itself. The mapping runs inside the
	// guarded region, so a payload with no episode in it is a failed read
	// rather than an exception escaping the port.
	it('reports a payload missing its subject as a failure instead of throwing', async () => {
		fetchMock.mockResolvedValue(jsonResponse({ moments: [] }))

		const result = await createStatsAdapter().episode('ep-001')

		expect(result.success).toBe(false)
		if (result.success) return
		expect(result.error.message).toContain('Não foi possível carregar o episódio')
	})

	// A page renders a message; it never sees a throw.
	it('reports a failure as a value, named for what was being read', async () => {
		fetchMock.mockResolvedValue(jsonResponse({ message: 'not found' }, 404))

		const result = await createStatsAdapter().episode('nope')

		expect(result.success).toBe(false)
		if (result.success) return
		expect(result.error.message).toContain('Não foi possível carregar o episódio')
	})
})
