import type { Count } from '@thom/dls-domain/types'
import { Empty } from './card'

/**
 * A leaderboard. The bar behind each row is drawn against the leader, which
 * is what makes the gap between first and second readable without axes.
 */
export function CountTable({ rows, color = 'var(--primary)' }: { rows: readonly Count[]; color?: string }) {
	if (rows.length === 0) return <Empty>Ainda não aconteceu.</Empty>

	const top = Math.max(...rows.map(row => row.count), 1)

	return (
		<ol className='flex flex-col gap-1'>
			{rows.map((row, index) => (
				<li key={row.key} className='relative flex items-center gap-3 overflow-hidden rounded-md px-2 py-1.5'>
					<span
						aria-hidden
						className='absolute inset-y-0 left-0 rounded-md opacity-15'
						style={{ width: `${(row.count / top) * 100}%`, backgroundColor: color }}
					/>
					<span className='text-muted-foreground nums relative w-5 text-right text-xs'>{index + 1}</span>
					<span className='relative flex-1 truncate text-sm'>{row.label}</span>
					<span className='nums relative text-sm font-semibold'>{row.count}</span>
				</li>
			))}
		</ol>
	)
}
