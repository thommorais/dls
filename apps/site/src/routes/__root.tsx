import { createRootRouteWithContext, HeadContent, Link, Outlet, Scripts } from '@tanstack/react-router'
import type { RouterContext } from '_/router'
import appCss from '_/styles/app.css?url'

const NAV = [
	{ to: '/', label: 'Números' },
	{ to: '/episodios', label: 'Episódios' },
	{ to: '/aberturas', label: 'Aberturas' },
	{ to: '/rankings', label: 'Rankings' },
	{ to: '/acervo', label: 'Acervo' },
] as const

// oxlint-disable-next-line import/no-default-export -- route
export const Route = createRootRouteWithContext<RouterContext>()({
	head: () => ({
		meta: [
			{ charSet: 'utf-8' },
			{ name: 'viewport', content: 'width=device-width, initial-scale=1' },
			{ title: 'DESCE A LETRA SHOW · números' },
			{
				name: 'description',
				content: 'As estatísticas do Desce a Letra Show: cantorias, músicas e as aberturas do público.',
			},
		],
		links: [{ rel: 'stylesheet', href: appCss }],
	}),
	component: RootComponent,
	notFoundComponent: () => (
		<Shell>
			<p className='text-muted-foreground py-24 text-center text-sm'>Essa página não existe.</p>
		</Shell>
	),
})

function RootComponent() {
	return (
		<html lang='pt-BR'>
			<head>
				<HeadContent />
			</head>
			<body>
				<Shell>
					<Outlet />
				</Shell>
				<Scripts />
			</body>
		</html>
	)
}

function Shell({ children }: { children: React.ReactNode }) {
	return (
		<div className='min-h-screen'>
			<header className='border-border bg-background/90 sticky top-0 z-40 border-b backdrop-blur'>
				<div className='mx-auto flex max-w-6xl flex-wrap items-center gap-x-6 gap-y-2 px-5 py-3'>
					<Link to='/' className='font-display text-sm uppercase'>
						Desce a Letra Show
					</Link>
					<nav className='label flex flex-wrap items-center gap-4 text-[11px]'>
						{NAV.map(item => (
							<Link
								key={item.to}
								to={item.to}
								className='text-muted-foreground hover:text-foreground border-b border-transparent pb-0.5 transition-colors'
								activeProps={{ className: 'text-foreground border-foreground' }}
								activeOptions={{ exact: item.to === '/' }}
							>
								{item.label}
							</Link>
						))}
					</nav>
				</div>
			</header>
			<main className='mx-auto max-w-6xl px-5 py-8'>{children}</main>
		</div>
	)
}
