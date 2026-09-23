import { describe, expect, it } from 'vitest'
import { filterToQuery, isEmptyFilter, parseFilter } from './filter'

describe('parseFilter', () => {
	it('keeps the parameters a route actually uses', () => {
		const filter = parseFilter({ from: '2025-07-01', to: '2025-12-31', genre: 'samba', kind: 'glossary', limit: '5' })

		expect(filter).toEqual({ from: '2025-07-01', to: '2025-12-31', genre: 'samba', kind: 'glossary', limit: 5 })
	})

	it('reads an empty search as no filter at all', () => {
		expect(parseFilter({})).toEqual({})
		expect(isEmptyFilter(parseFilter({}))).toBe(true)
	})

	// A bookmarked URL with a stale parameter should still show the page,
	// which is the same rule the Go handler follows.
	it('drops a date it cannot read instead of throwing', () => {
		expect(parseFilter({ from: 'ontem', to: '2025-13-45' })).toEqual({})
	})

	it('drops a kind that is not one of the archive shelves', () => {
		expect(parseFilter({ kind: 'receitas' })).toEqual({})
	})

	it('drops a limit that cannot cap anything', () => {
		expect(parseFilter({ limit: '0' })).toEqual({})
		expect(parseFilter({ limit: '-3' })).toEqual({})
		expect(parseFilter({ limit: 'muitos' })).toEqual({})
	})

	it('takes a number for the limit as readily as a string', () => {
		expect(parseFilter({ limit: 12 })).toEqual({ limit: 12 })
	})

	it('trims a genre and ignores one that is only spaces', () => {
		expect(parseFilter({ genre: '  forró  ' })).toEqual({ genre: 'forró' })
		expect(parseFilter({ genre: '   ' })).toEqual({})
	})

	it('survives values of the wrong shape entirely', () => {
		expect(parseFilter({ from: null, to: [], genre: 42, kind: {}, limit: undefined })).toEqual({})
	})
})

describe('filterToQuery', () => {
	it('writes only what is set', () => {
		expect(filterToQuery({ from: '2025-07-01', limit: 5 })).toEqual({ from: '2025-07-01', limit: '5' })
	})

	it('writes nothing for an empty filter', () => {
		expect(filterToQuery({})).toEqual({})
	})

	it('round-trips through parseFilter', () => {
		const filter = { from: '2025-07-01', to: '2025-12-31', genre: 'samba', kind: 'trivia', limit: 7 } as const

		expect(parseFilter(filterToQuery(filter))).toEqual(filter)
	})
})
