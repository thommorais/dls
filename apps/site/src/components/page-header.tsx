export function PageHeader({ section, title, action }: { section: string; title: string; action?: React.ReactNode }) {
	return (
		<div className='mb-6 flex flex-col gap-3 sm:flex-row sm:items-end sm:justify-between'>
			<div>
				<p className='text-muted-foreground mb-0.5 text-xs font-medium'>{section}</p>
				<h1 className='text-foreground text-2xl font-semibold tracking-tight'>{title}</h1>
			</div>
			{action}
		</div>
	)
}
