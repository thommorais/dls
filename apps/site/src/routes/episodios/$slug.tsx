import { useQuery } from '@tanstack/react-query'
import { createFileRoute, Link } from '@tanstack/react-router'
import { formatTimestamp, youtubeLink } from '@thom/dls-domain/format'
import type { EpisodeDetail, Moment } from '@thom/dls-domain/types'
import { episodeQuery } from '_/app/queries'
import { Card, Empty } from '_/components/card'
import { PageHeader } from '_/components/page-header'
import { StatTile } from '_/components/stat-tile'

// oxlint-disable-next-line import/no-default-export -- route
export const Route = createFileRoute('/episodios/$slug')({
	loader: ({ context, params }) => context.queryClient.ensureQueryData(episodeQuery(context.container, params.slug)),
	component: EpisodePage,
})

const day = new Intl.DateTimeFormat('pt-BR', { day: '2-digit', month: 'long', year: 'numeric', timeZone: 'UTC' })

function EpisodePage() {
	const { container } = Route.useRouteContext()
	const { slug } = Route.useParams()
	const initialData = Route.useLoaderData()

	const { data } = useQuery({ ...episodeQuery(container, slug), initialData })
	const { episode } = data

	const names = new Map(data.people.map(person => [person.id, person.name]))
	const songs = new Map(data.songs.map(song => [song.id, song]))

	return (
		<>
			<PageHeader
				section={`Episódio ${episode.number} · ${day.format(episode.publishedAt)}`}
				title={episode.title}
				action={
					<Link to='/episodios' className='text-muted-foreground hover:text-foreground text-xs'>
						voltar aos episódios
					</Link>
				}
			/>

			<div className='mb-5 grid gap-3 sm:grid-cols-3'>
				<StatTile label='Momentos' value={episode.tally.totalMoments} />
				<StatTile label='Convidadas entrevistadas' value={episode.tally.womenInterviewed} />
				<StatTile label='Duração' value={formatTimestamp(episode.durationSeconds)} />
			</div>

			<div className='grid gap-5 lg:grid-cols-[2fr_1fr]'>
				<Card title='Linha do tempo' subtitle='Cada momento, no segundo em que aconteceu'>
					<Timeline detail={data} names={names} songs={songs} youtubeId={episode.youtubeId} />
				</Card>

				<div className='flex flex-col gap-5'>
					<Card title='No estúdio'>
						{data.appearances.length === 0 ? (
							<Empty>Ninguém registrado.</Empty>
						) : (
							<ul className='flex flex-col gap-1.5'>
								{data.appearances.map(appearance => (
									<li key={appearance.person} className='flex items-center justify-between gap-2 text-sm'>
										<span className='truncate'>{names.get(appearance.person) ?? 'desconhecido'}</span>
										<span className='text-muted-foreground text-xs'>
											{appearance.role === 'guest'
												? appearance.isInterview
													? 'entrevistado'
													: 'passou por lá'
												: appearance.role === 'host'
													? 'apresentação'
													: 'do corredor'}
										</span>
									</li>
								))}
							</ul>
						)}
					</Card>

					<Card title='Aberturas' subtitle='Mandadas pelo público'>
						{data.openings.length === 0 ? (
							<Empty>Nenhuma abertura neste episódio.</Empty>
						) : (
							<ul className='flex flex-col gap-2'>
								{data.openings.map(opening => (
									<li key={opening.id} className='text-sm'>
										<p className='font-medium'>{opening.title}</p>
										<p className='text-muted-foreground text-xs'>
											{opening.authorName} · {opening.genre}
										</p>
									</li>
								))}
							</ul>
						)}
					</Card>
				</div>
			</div>
		</>
	)
}

function Timeline({
	detail,
	names,
	songs,
	youtubeId,
}: {
	detail: EpisodeDetail
	names: Map<string, string>
	songs: Map<string, { title: string; artist: string }>
	youtubeId: string
}) {
	if (detail.moments.length === 0) return <Empty>Nada contado neste episódio.</Empty>

	const ordered = [...detail.moments].sort((a, b) => a.videoTimestamp - b.videoTimestamp)

	return (
		<ol className='flex flex-col'>
			{ordered.map(moment => (
				<li key={moment.id} className='border-border flex gap-3 border-b py-2.5 last:border-0'>
					<Stamp moment={moment} youtubeId={youtubeId} />
					<div className='min-w-0 flex-1'>
						<p className='label text-muted-foreground text-[11px]'>{moment.typeLabel}</p>
						<p className='text-sm'>{moment.summary}</p>
						<p className='text-muted-foreground mt-0.5 text-xs'>
							{moment.actor ? names.get(moment.actor) : null}
							{moment.song && songs.has(moment.song)
								? ` · ${songs.get(moment.song)?.title} (${songs.get(moment.song)?.artist})`
								: null}
						</p>
					</div>
				</li>
			))}
		</ol>
	)
}

function Stamp({ moment, youtubeId }: { moment: Moment; youtubeId: string }) {
	const label = formatTimestamp(moment.videoTimestamp)
	const href = youtubeLink(youtubeId, moment.videoTimestamp)

	if (!href) return <span className='nums text-muted-foreground w-16 shrink-0 pt-0.5 text-xs'>{label}</span>

	return (
		<a
			href={href}
			target='_blank'
			rel='noreferrer'
			className='nums text-muted-foreground hover:text-primary w-16 shrink-0 pt-0.5 text-xs underline-offset-2 hover:underline'
		>
			{label}
		</a>
	)
}
