export function Card({ title, subtitle, children }: { title?: string; subtitle?: string; children: React.ReactNode }) {
	return (
		<section className='border-border bg-card rounded-xl border px-5 pt-5 pb-4 shadow-sm'>
			{title ? (
				<header className='mb-3'>
					<h2 className='text-sm font-semibold'>{title}</h2>
					{subtitle ? <p className='text-muted-foreground mt-0.5 text-xs'>{subtitle}</p> : null}
				</header>
			) : null}
			{children}
		</section>
	)
}

export function Empty({ children }: { children: React.ReactNode }) {
	return <p className='text-muted-foreground py-16 text-center text-sm'>{children}</p>
}
