export function StatTile({
	label,
	value,
	hint,
	color,
}: {
	label: string
	value: string | number
	hint?: string
	color?: string
}) {
	return (
		<div className='border-border bg-card rounded-xl border p-4 shadow-sm'>
			<div className='mb-2 flex items-center gap-2'>
				{color ? <span className='size-2.5 shrink-0 rounded-full' style={{ backgroundColor: color }} /> : null}
				<p className='text-muted-foreground text-xs font-medium'>{label}</p>
			</div>
			<p className='nums text-2xl font-bold'>{value}</p>
			{hint ? <p className='text-muted-foreground nums mt-1 text-xs'>{hint}</p> : null}
		</div>
	)
}
