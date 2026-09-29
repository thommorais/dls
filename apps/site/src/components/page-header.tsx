export function PageHeader({ section, title, action }: { section: string; title: string; action?: React.ReactNode }) {
	return (
		<div className='mb-6 flex flex-col gap-3 sm:flex-row sm:items-end sm:justify-between'>
			<div>
				<p className='label text-muted-foreground mb-1 text-xs'>{section}</p>
				<h1 className='font-display text-foreground text-2xl leading-tight uppercase sm:text-3xl'>{title}</h1>
			</div>
			{action}
		</div>
	)
}
