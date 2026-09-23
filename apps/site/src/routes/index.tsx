import { useQuery } from '@tanstack/react-query'
import { createFileRoute } from '@tanstack/react-router'
import { formatPerEpisode } from '@thom/dls-domain/format'
import { tallySeries } from '@thom/dls-domain/series'
import { overviewQuery } from '_/app/queries'
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
			<PageHeader section='CAFÉ COM CAOS' title='Os números do programa' />

			<div className='border-border bg-card mb-5 rounded-xl border px-5 pt-5 pb-4 shadow-sm'>
				<div className='border-border mb-4 border-b pb-4'>
					<p className='text-muted-foreground mb-1 text-xs font-semibold'>Momentos contados</p>
					<p className='nums text-4xl leading-tight font-bold'>{overview.tally.totalMoments}</p>
					<p className='text-muted-foreground mt-2 text-xs'>
						em {overview.episodes} episódios · {range}
					</p>
				</div>

				<div className='grid grid-cols-2 gap-x-5 gap-y-4 sm:grid-cols-3'>
					<div>
						<p className='text-muted-foreground mb-1 text-xs font-medium'>Mulheres entrevistadas</p>
						<p className='nums text-xl font-bold'>{overview.tally.womenInterviewed}</p>
					</div>
					<div>
						<p className='text-muted-foreground mb-1 text-xs font-medium'>Aberturas no ar</p>
						<p className='nums text-xl font-bold'>{overview.tally.openingsAired}</p>
					</div>
					<div>
						<p className='text-muted-foreground mb-1 text-xs font-medium'>Momentos por episódio</p>
						<p className='nums text-primary text-xl font-bold'>{formatPerEpisode(perEpisode)}</p>
					</div>
				</div>
			</div>

			<div className='mb-5 grid grid-cols-2 gap-3 lg:grid-cols-3'>
				{series.map(point => (
					<StatTile key={point.slug} label={point.label} value={point.count} color={point.color} />
				))}
			</div>

			<div className='border-border bg-card rounded-xl border px-5 pt-5 pb-4 shadow-sm'>
				<p className='mb-3 text-sm font-semibold'>Momentos por tipo</p>
				<TypeBars points={series} />
			</div>
		</>
	)
}
