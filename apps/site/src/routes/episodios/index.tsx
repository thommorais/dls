import { useQuery } from '@tanstack/react-query'
import { createFileRoute, Link } from '@tanstack/react-router'
import { cumulative, episodeSeries } from '@thom/dls-domain/series'
import type { Episode, MomentType } from '@thom/dls-domain/types'
import { episodesQuery, overviewQuery } from '_/app/queries'
import { Card, Empty } from '_/components/card'
import { PageHeader } from '_/components/page-header'
import { grayShades } from '_/lib/palette'
import { Area, AreaChart, ResponsiveContainer, Tooltip, XAxis, YAxis } from 'recharts'

// oxlint-disable-next-line import/no-default-export -- route
export const Route = createFileRoute('/episodios/')({
	loader: async ({ context }) => {
		const [episodes, overview] = await Promise.all([
			context.queryClient.ensureQueryData(episodesQuery(context.container)),
			context.queryClient.ensureQueryData(overviewQuery(context.container)),
		])
		return { episodes, overview }
	},
	component: Episodes,
})

const day = new Intl.DateTimeFormat('pt-BR', { day: '2-digit', month: '2-digit', year: '2-digit', timeZone: 'UTC' })

const tooltipStyle: React.CSSProperties = {
	background: 'var(--card)',
	color: 'var(--foreground)',
	border: '1px solid var(--border)',
	borderRadius: '0.75rem',
	fontSize: '12px',
}

function Episodes() {
	const { container } = Route.useRouteContext()
	const initial = Route.useLoaderData()

	const { data: episodes } = useQuery({ ...episodesQuery(container), initialData: initial.episodes })
	const { data: overview } = useQuery({ ...overviewQuery(container), initialData: initial.overview })

	return (
		<>
			<PageHeader section='Episódios' title={`${episodes.length} episódios`} />

			<div className='mb-5'>
				<Card title='Momentos acumulados' subtitle='Total somado episódio a episódio, por tipo'>
					<Cumulative episodes={episodes} types={overview.types} />
				</Card>
			</div>

			{episodes.length === 0 ? (
				<Empty>Nenhum episódio no período.</Empty>
			) : (
				<div className='flex flex-col gap-2'>
					{[...episodes]
						.sort((a, b) => b.number - a.number)
						.map(episode => (
							<EpisodeRow key={episode.id} episode={episode} types={overview.types} />
						))}
				</div>
			)}
		</>
	)
}

function EpisodeRow({ episode, types }: { episode: Episode; types: readonly MomentType[] }) {
	return (
		<Link
			to='/episodios/$slug'
			params={{ slug: episode.slug }}
			className='border-border bg-card hover:border-foreground flex flex-wrap items-center gap-x-4 gap-y-2 rounded-xl border px-4 py-3 transition-colors'
		>
			<span className='text-muted-foreground nums w-10 text-xs font-semibold'>#{episode.number}</span>
			<span className='min-w-0 flex-1 truncate text-sm font-medium'>{episode.title}</span>
			<span className='text-muted-foreground nums text-xs'>{day.format(episode.publishedAt)}</span>
			<span className='flex items-center gap-1.5'>
				{types.map(type => {
					const count = episode.tally.moments[type.slug] ?? 0
					return (
						<span
							key={type.slug}
							title={`${type.label}: ${count}`}
							className={
								count > 0
									? 'bg-foreground/8 text-foreground nums rounded px-1.5 py-0.5 text-[11px] font-semibold'
									: 'text-muted-foreground/50 nums rounded px-1.5 py-0.5 text-[11px] font-semibold'
							}
						>
							{count}
						</span>
					)
				})}
			</span>
		</Link>
	)
}

function Cumulative({ episodes, types }: { episodes: readonly Episode[]; types: readonly MomentType[] }) {
	if (episodes.length === 0) return <Empty>Nada para desenhar.</Empty>

	const shades = grayShades(types.length)

	// One row per episode, one column per type, so the areas stack on a
	// shared x axis without recomputing the sort for each series.
	const byLabel = new Map<string, Record<string, number | string>>()
	for (const type of types) {
		for (const point of cumulative(episodeSeries(episodes, type.slug))) {
			const row = byLabel.get(point.label) ?? { label: point.label }
			row[type.slug] = point.value
			byLabel.set(point.label, row)
		}
	}

	return (
		<div className='h-[300px]'>
			<ResponsiveContainer width='100%' height='100%'>
				<AreaChart data={[...byLabel.values()]} margin={{ top: 8, right: 8, left: 0, bottom: 0 }}>
					<XAxis
						dataKey='label'
						tick={{ fontSize: 10, fill: 'var(--muted-foreground)' }}
						axisLine={false}
						tickLine={false}
						interval='preserveStartEnd'
					/>
					<YAxis
						tick={{ fontSize: 10, fill: 'var(--muted-foreground)' }}
						axisLine={false}
						tickLine={false}
						width={40}
						tickCount={5}
					/>
					<Tooltip contentStyle={tooltipStyle} />
					{types.map((type, index) => (
						<Area
							key={type.slug}
							type='monotone'
							dataKey={type.slug}
							name={type.label}
							stackId='total'
							stroke={shades[index]}
							fill={shades[index]}
							fillOpacity={0.55}
							strokeWidth={1.5}
							dot={false}
						/>
					))}
				</AreaChart>
			</ResponsiveContainer>
		</div>
	)
}
