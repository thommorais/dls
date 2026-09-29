const VARIANT: Record<'outline' | 'dashed' | 'solid', string> = {
	outline: 'border-border bg-card text-card-foreground border',
	dashed: 'border-border bg-card text-card-foreground border border-dash',
	solid: 'border-foreground bg-foreground text-background border',
}

export function Card({
	title,
	subtitle,
	variant = 'outline',
	children,
}: {
	title?: string
	subtitle?: string
	variant?: 'outline' | 'dashed' | 'solid'
	children: React.ReactNode
}) {
	return (
		<section className={`rounded-2xl px-5 pt-5 pb-4 ${VARIANT[variant]}`}>
			{title ? (
				<header className='mb-3'>
					<h2 className='label text-xs'>{title}</h2>
					{subtitle ? <p className='mt-0.5 text-xs opacity-60'>{subtitle}</p> : null}
				</header>
			) : null}
			{children}
		</section>
	)
}

export function Empty({ children }: { children: React.ReactNode }) {
	return <p className='text-muted-foreground py-16 text-center text-sm'>{children}</p>
}
