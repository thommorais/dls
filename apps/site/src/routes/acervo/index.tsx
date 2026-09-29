import { useQuery } from '@tanstack/react-query'
import { createFileRoute, Link } from '@tanstack/react-router'
import { parseFilter } from '@thom/dls-domain/filter'
import { ARCHIVE_KINDS, type ArchiveKind } from '@thom/dls-domain/types'
import { archiveQuery } from '_/app/queries'
import { Empty } from '_/components/card'
import { PageHeader } from '_/components/page-header'

// oxlint-disable-next-line import/no-default-export -- route
export const Route = createFileRoute('/acervo/')({
	validateSearch: (search: Record<string, unknown>) => parseFilter(search),
	loaderDeps: ({ search }) => search,
	loader: ({ context, deps }) => context.queryClient.ensureQueryData(archiveQuery(context.container, deps)),
	component: Archive,
})

const KIND_LABELS: Record<ArchiveKind, string> = {
	about: 'sobre',
	glossary: 'glossário',
	segment: 'quadros',
	trivia: 'curiosidades',
	bio: 'gente',
	milestone: 'marcos',
}

function Archive() {
	const { container } = Route.useRouteContext()
	const search = Route.useSearch()
	const initialData = Route.useLoaderData()

	const { data: entries } = useQuery({ ...archiveQuery(container, search), initialData })

	return (
		<>
			<PageHeader section='Acervo' title='O que é cada coisa aqui' />

			<div className='mb-5 flex flex-wrap items-center gap-1.5'>
				<Link
					to='/acervo'
					search={{}}
					className='border-border border-dash hover:border-foreground label rounded-full border px-3 py-1 text-[11px]'
					activeProps={{ className: 'bg-foreground text-background border-foreground border-solid' }}
					activeOptions={{ exact: true, includeSearch: true }}
				>
					tudo
				</Link>
				{ARCHIVE_KINDS.map(kind => (
					<Link
						key={kind}
						to='/acervo'
						search={{ kind }}
						className='border-border border-dash hover:border-foreground label rounded-full border px-3 py-1 text-[11px]'
						activeProps={{ className: 'bg-foreground text-background border-foreground border-solid' }}
						activeOptions={{ includeSearch: true }}
					>
						{KIND_LABELS[kind]}
					</Link>
				))}
			</div>

			{entries.length === 0 ? (
				<Empty>Nada no acervo com esse filtro.</Empty>
			) : (
				<div className='grid gap-3 sm:grid-cols-2'>
					{entries.map(entry => (
						<Link
							key={entry.id}
							to='/acervo/$slug'
							params={{ slug: entry.slug }}
							className='border-border bg-card border-dash hover:border-foreground rounded-2xl border p-4 transition-colors'
						>
							<p className='label text-muted-foreground mb-1 text-[11px]'>{KIND_LABELS[entry.kind] ?? entry.kind}</p>
							<p className='text-sm font-semibold'>{entry.title}</p>
							<p className='text-muted-foreground mt-1 text-xs'>{entry.summary}</p>
						</Link>
					))}
				</div>
			)}
		</>
	)
}
