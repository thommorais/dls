import { useQuery } from '@tanstack/react-query'
import { createFileRoute, Link } from '@tanstack/react-router'
import { archiveEntryQuery } from '_/app/queries'
import { PageHeader } from '_/components/page-header'

// oxlint-disable-next-line import/no-default-export -- route
export const Route = createFileRoute('/acervo/$slug')({
	loader: ({ context, params }) =>
		context.queryClient.ensureQueryData(archiveEntryQuery(context.container, params.slug)),
	component: ArchiveEntryPage,
})

function ArchiveEntryPage() {
	const { container } = Route.useRouteContext()
	const { slug } = Route.useParams()
	const initialData = Route.useLoaderData()

	const { data: entry } = useQuery({ ...archiveEntryQuery(container, slug), initialData })

	return (
		<>
			<PageHeader
				section='Acervo'
				title={entry.title}
				action={
					<Link to='/acervo' className='text-muted-foreground hover:text-foreground text-xs'>
						voltar ao acervo
					</Link>
				}
			/>

			<article className='border-border bg-card rounded-2xl border px-6 py-6'>
				{entry.summary ? <p className='text-muted-foreground mb-4 text-sm'>{entry.summary}</p> : null}

				{/* The body is plain text with blank lines between paragraphs, so it
				    is split rather than parsed. */}
				<div className='flex flex-col gap-3 text-sm leading-relaxed'>
					{entry.body
						.split('\n\n')
						.filter(Boolean)
						.map(paragraph => (
							<p key={paragraph.slice(0, 40)}>{paragraph}</p>
						))}
				</div>

				{entry.tags.length > 0 ? (
					<div className='border-border mt-5 flex flex-wrap gap-1.5 border-t pt-4'>
						{entry.tags.map(tag => (
							<span key={tag} className='bg-secondary text-muted-foreground rounded-full px-2.5 py-0.5 text-xs'>
								{tag}
							</span>
						))}
					</div>
				) : null}
			</article>
		</>
	)
}
