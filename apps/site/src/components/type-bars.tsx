import type { TallyPoint } from '@thom/dls-domain/series'
import { Bar, BarChart, Cell, ResponsiveContainer, Tooltip, XAxis, YAxis } from 'recharts'

const tooltipStyle: React.CSSProperties = {
	background: 'var(--card)',
	color: 'var(--foreground)',
	border: '1px solid var(--border)',
	borderRadius: '0.75rem',
	fontSize: '12px',
}

/**
 * The colour of every bar comes from the moment type itself, served by the
 * API, so a stat keeps its colour on every page it appears on.
 */
export function TypeBars({ points }: { points: readonly TallyPoint[] }) {
	if (points.length === 0) {
		return <p className='text-muted-foreground py-16 text-center text-sm'>Nada contado ainda.</p>
	}

	return (
		<div className='h-[280px]'>
			<ResponsiveContainer width='100%' height='100%'>
				<BarChart data={[...points]} margin={{ top: 8, right: 8, left: 0, bottom: 0 }}>
					<XAxis
						dataKey='label'
						tick={{ fontSize: 10, fill: 'var(--muted-foreground)' }}
						axisLine={false}
						tickLine={false}
						interval={0}
						height={48}
						angle={-18}
						textAnchor='end'
					/>
					<YAxis
						tick={{ fontSize: 10, fill: 'var(--muted-foreground)' }}
						axisLine={false}
						tickLine={false}
						width={40}
						tickCount={5}
					/>
					<Tooltip contentStyle={tooltipStyle} cursor={{ fill: 'var(--muted)', opacity: 0.4 }} />
					<Bar dataKey='count' name='momentos' radius={[6, 6, 0, 0]} maxBarSize={64}>
						{points.map(point => (
							<Cell key={point.slug} fill={point.color} />
						))}
					</Bar>
				</BarChart>
			</ResponsiveContainer>
		</div>
	)
}
