import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { createRouter } from '@tanstack/react-router'
import { createContainer, type Container } from '_/app/container'
import { routeTree } from './routeTree.gen'

export type RouterContext = {
	readonly queryClient: QueryClient
	readonly container: Container
}

export function getRouter() {
	const queryClient = new QueryClient({
		defaultOptions: {
			queries: {
				// The page ships with numbers baked in at build time. Zero
				// staleness is what makes the browser go and check them the
				// moment it takes over.
				staleTime: 0,
				retry: 1,
			},
		},
	})

	const container = createContainer()

	return createRouter({
		routeTree,
		context: { queryClient, container },
		scrollRestoration: true,
		defaultPreload: 'intent',
		Wrap: ({ children }) => <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>,
	})
}

declare module '@tanstack/react-router' {
	interface Register {
		router: ReturnType<typeof getRouter>
	}
}
