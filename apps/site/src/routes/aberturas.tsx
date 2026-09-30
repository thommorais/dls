import { useQuery } from '@tanstack/react-query'
import { createFileRoute, Link } from '@tanstack/react-router'
import { parseFilter } from '@thom/dls-domain/filter'
import { openingsQuery } from '_/app/queries'
import { Card, Empty } from '_/components/card'
import { PageHeader } from '_/components/page-header'
import { StatTile } from '_/components/stat-tile'

// oxlint-disable-next-line import/no-default-export -- route
export const Route = createFileRoute('/aberturas')({
	// The filter lives in the URL, so a filtered page can be shared, and the
	// same codec validates it here and builds the request.
	validateSearch: (search: Record<string, unknown>) => parseFilter(search),
	loaderDeps: ({ search }) => search,
	loader: ({ context, deps }) => context.queryClient.ensureQueryData(openingsQuery(context.container, deps)),
	component: Openings,
})

const day = new Intl.DateTimeFormat('pt-BR', { day: '2-digit', month: '2-digit', year: 'numeric', timeZone: 'UTC' })

function Openings() {
	const { container } = Route.useRouteContext()
	const search = Route.useSearch()
	const initialData = Route.useLoaderData()

	const { data: openings } = useQuery({ ...openingsQuery(container, search), initialData })

	// The genre list comes from what is actually there, so a filter never
	// offers an option that returns nothing.
	const genres = [...new Set(openings.map(opening => opening.genre).filter(Boolean))].sort()

	return (
		<>
			<PageHeader section='Aberturas' title='Mandadas pelo público' />

			<div className='mb-5 flex flex-wrap items-center gap-1.5'>
				<Link
					to='/aberturas'
					search={{}}
					className='border-border border-dash hover:border-foreground label rounded-full border px-3 py-1 text-[11px]'
					activeProps={{ className: 'bg-foreground text-background border-foreground border-solid' }}
					activeOptions={{ exact: true, includeSearch: true }}
				>
					todas
				</Link>
				{genres.map(genre => (
					<Link
						key={genre}
						to='/aberturas'
						search={{ genre }}
						className='border-border border-dash hover:border-foreground label rounded-full border px-3 py-1 text-[11px]'
						activeProps={{ className: 'bg-foreground text-background border-foreground border-solid' }}
						activeOptions={{ includeSearch: true }}
					>
						{genre}
					</Link>
				))}
			</div>

			<div className='mb-5 grid gap-3'>
				<StatTile label='Já foram ao ar' value={openings.length} />
			</div>

			<div className='grid gap-5'>
				<Card title='No ar' subtitle='Com a data em que tocaram'>
					{openings.length === 0 ? (
						<Empty>Nenhuma abertura no ar com esse filtro.</Empty>
					) : (
						<ul className='flex flex-col gap-2.5'>
							{openings.map(opening => (
								<li key={opening.id} className='border-border border-b pb-2.5 last:border-0 last:pb-0'>
									<div className='flex items-baseline justify-between gap-3'>
										<p className='truncate text-sm font-medium'>{opening.title}</p>
										<span className='text-muted-foreground nums shrink-0 text-xs'>
											{opening.airedAt ? day.format(opening.airedAt) : null}
										</span>
									</div>
									<p className='text-muted-foreground text-xs'>
										{opening.authorName} {opening.authorHandle} · {opening.genre}
									</p>
								</li>
							))}
						</ul>
					)}
				</Card>
			</div>
		</>
	)
}
