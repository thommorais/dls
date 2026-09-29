import { useQuery } from '@tanstack/react-query'
import { createFileRoute } from '@tanstack/react-router'
import { formatPerEpisode } from '@thom/dls-domain/format'
import { tallySeries } from '@thom/dls-domain/series'
import { overviewQuery } from '_/app/queries'
import { Card } from '_/components/card'
import { PageHeader } from '_/components/page-header'
import { StatTile } from '_/components/stat-tile'
import { TypeBars } from '_/components/type-bars'

// oxlint-disable-next-line import/no-default-export -- route
export const Route = createFileRoute('/')({
	// Runs at build time, which is what bakes the numbers into the HTML, and
	// again on a client navigation.
	loader: ({ context }) => context.queryClient.ensureQueryData(overviewQuery(context.container)),
	component: Home,
})

const day = new Intl.DateTimeFormat('pt-BR', { day: '2-digit', month: 'short', year: 'numeric', timeZone: 'UTC' })

function Home() {
	const { container } = Route.useRouteContext()
	const initialData = Route.useLoaderData()

	// The prerendered markup is the initial data, so the first paint matches
	// the HTML exactly; the browser then refetches and the numbers move if
	// the show happened again in the meantime.
	const { data: overview } = useQuery({ ...overviewQuery(container), initialData })

	const series = tallySeries(overview.tally, overview.types)
	const perEpisode = overview.episodes > 0 ? overview.tally.totalMoments / overview.episodes : 0

	const range =
		overview.first && overview.latest
			? `${day.format(overview.first.publishedAt)} a ${day.format(overview.latest.publishedAt)}`
			: 'sem episódios no período'

	return (
		<>
			<PageHeader section='Café com Caos' title='Os números do programa' />

			<div className='border-foreground bg-foreground text-background mb-5 rounded-2xl border px-5 pt-5 pb-4'>
				<div className='mb-4 border-b border-white/15 pb-4'>
					<p className='label mb-1 text-xs opacity-70'>Momentos contados</p>
					<p className='nums text-4xl leading-tight font-bold'>{overview.tally.totalMoments}</p>
					<p className='mt-2 text-xs opacity-70'>
						em {overview.episodes} episódios · {range}
					</p>
				</div>

				<div className='grid grid-cols-2 gap-x-5 gap-y-4 sm:grid-cols-3'>
					<div>
						<p className='label mb-1 text-[11px] opacity-70'>Mulheres entrevistadas</p>
						<p className='nums text-xl font-bold'>{overview.tally.womenInterviewed}</p>
					</div>
					<div>
						<p className='label mb-1 text-[11px] opacity-70'>Aberturas no ar</p>
						<p className='nums text-xl font-bold'>{overview.tally.openingsAired}</p>
					</div>
					<div>
						<p className='label mb-1 text-[11px] opacity-70'>Momentos por episódio</p>
						<p className='nums text-xl font-bold'>{formatPerEpisode(perEpisode)}</p>
					</div>
				</div>
			</div>

			<div className='mb-5 grid grid-cols-2 gap-3 lg:grid-cols-3'>
				{series.map(point => (
					<StatTile key={point.slug} label={point.label} value={point.count} />
				))}
			</div>

			<Card title='Momentos por tipo'>
				<TypeBars points={series} />
			</Card>
		</>
	)
}
