export function StatTile({
	label,
	value,
	hint,
	emphasis = false,
}: {
	label: string
	value: string | number
	hint?: string
	emphasis?: boolean
}) {
	return (
		<div
			className={
				emphasis
					? 'border-foreground bg-foreground text-background rounded-2xl border p-4'
					: 'border-border bg-card text-card-foreground border-dash rounded-2xl border p-4'
			}
		>
			<p className={`label mb-2 text-[11px] ${emphasis ? 'opacity-70' : 'text-muted-foreground'}`}>{label}</p>
			<p className='nums text-2xl font-bold'>{value}</p>
			{hint ? <p className={`nums mt-1 text-xs ${emphasis ? 'opacity-70' : 'text-muted-foreground'}`}>{hint}</p> : null}
		</div>
	)
}
