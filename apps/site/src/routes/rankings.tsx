import { useQuery } from '@tanstack/react-query'
import { createFileRoute } from '@tanstack/react-router'
import { formatPerEpisode } from '@thom/dls-domain/format'
import { rankingsQuery } from '_/app/queries'
import { Card, Empty } from '_/components/card'
import { CountTable } from '_/components/count-table'
import { PageHeader } from '_/components/page-header'

// oxlint-disable-next-line import/no-default-export -- route
export const Route = createFileRoute('/rankings')({
	loader: ({ context }) => context.queryClient.ensureQueryData(rankingsQuery(context.container, { limit: 10 })),
	component: Rankings,
})

function Rankings() {
	const { container } = Route.useRouteContext()
	const initialData = Route.useLoaderData()

	const { data: rankings } = useQuery({ ...rankingsQuery(container, { limit: 10 }), initialData })

	const records = new Map(rankings.records.map(record => [record.type, record]))
	const streaks = new Map(rankings.streaks.map(streak => [streak.type, streak]))
	const averages = new Map(rankings.averages.map(average => [average.type, average]))

	return (
		<>
			<PageHeader section='Rankings' title='Quem faz mais, e quando' />

			<div className='mb-5 grid gap-3 sm:grid-cols-2 lg:grid-cols-3'>
				{rankings.types.map(type => {
					const record = records.get(type.slug)
					const streak = streaks.get(type.slug)
					const average = averages.get(type.slug)

					return (
						<div key={type.slug} className='border-border bg-card border-dash rounded-2xl border p-4'>
							<p className='label mb-3 text-xs'>{type.label}</p>

							<p className='text-muted-foreground text-xs'>Recorde num episódio</p>
							<p className='nums text-2xl font-bold'>{record?.count ?? 0}</p>
							<p className='text-muted-foreground mt-0.5 truncate text-xs'>{record?.title ?? 'nunca aconteceu'}</p>

							<div className='border-border mt-3 grid grid-cols-2 gap-2 border-t pt-3'>
								<div>
									<p className='text-muted-foreground text-xs'>Sequência</p>
									<p className='nums text-sm font-semibold'>{streak?.length ?? 0} eps</p>
								</div>
								<div>
									<p className='text-muted-foreground text-xs'>Por episódio</p>
									<p className='nums text-sm font-semibold'>{formatPerEpisode(average?.perEpisode ?? 0)}</p>
								</div>
							</div>
						</div>
					)
				})}
			</div>

			<div className='mb-5 grid gap-5 lg:grid-cols-3'>
				<Card title='Músicas mais emendadas'>
					<CountTable rows={rankings.songs} />
				</Card>
				<Card title='Palavras que puxam música'>
					<CountTable rows={rankings.triggerWords} />
				</Card>
				<Card title='Convidados que mais voltaram'>
					<CountTable rows={rankings.guests} />
				</Card>
			</div>

			<div className='grid gap-5 lg:grid-cols-2'>
				{rankings.types.map(type => {
					const board = rankings.actors[type.slug] ?? []
					return (
						<Card key={type.slug} title={type.label} subtitle='Quem mais fez'>
							{board.length === 0 ? <Empty>Ainda não aconteceu.</Empty> : <CountTable rows={board} />}
						</Card>
					)
				})}
			</div>
		</>
	)
}
