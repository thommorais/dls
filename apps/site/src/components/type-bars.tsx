import type { TallyPoint } from '@thom/dls-domain/series'
import { Bar, BarChart, ResponsiveContainer, Tooltip, XAxis, YAxis } from 'recharts'

const tooltipStyle: React.CSSProperties = {
	background: 'var(--card)',
	color: 'var(--foreground)',
	border: '1px solid var(--border)',
	borderRadius: '0.75rem',
	fontSize: '12px',
}

/**
 * Monochrome by design: every bar is the same fill, so a stat reads by
 * height alone, not by which moment type happens to be which hue.
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
					<Tooltip contentStyle={tooltipStyle} cursor={{ fill: 'var(--muted)', opacity: 0.6 }} />
					<Bar dataKey='count' name='momentos' fill='var(--foreground)' radius={[2, 2, 0, 0]} maxBarSize={64} />
				</BarChart>
			</ResponsiveContainer>
		</div>
	)
}
