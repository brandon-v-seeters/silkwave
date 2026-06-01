<script lang="ts">
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import DarkModeToggle from '$lib/components/atoms/DarkModeToggle.svelte';
	import Button from '$lib/components/ui/button/button.svelte';
	import { cn } from '$lib/utils/utils.js';
	import { Search } from '@lucide/svelte';
	import type { LayoutProps } from './$types';

	let { data, children, params }: LayoutProps = $props();

	const artist = $derived(data.artist);
	const artistHref = $derived(
		resolve('/artist/[artistSlug]', { artistSlug: artist?.slug ?? params.artistSlug ?? '' })
	);
	const discoverHref = resolve('/discover');
	const isArtistHome = $derived(page.url.pathname === artistHref);
	const artistInitial = $derived(artist?.name?.trim().slice(0, 1).toUpperCase() || 'S');
</script>

<div class="min-h-dvh bg-background text-foreground">
	<header
		class="sticky top-0 z-40 border-b border-border/70 bg-background/90 backdrop-blur-xl"
	>
		<div
			class="mx-auto flex h-16 w-full max-w-7xl items-center justify-between gap-3 px-4 sm:px-6 lg:px-8"
		>
			<div class="flex min-w-0 items-center gap-3">
				<a
					href={resolve('/')}
					class="shrink-0 text-sm font-semibold tracking-tight text-foreground transition hover:text-primary"
				>
					Silkwave
				</a>
				<span class="hidden h-4 w-px bg-border sm:block" aria-hidden="true"></span>
				<a
					href={resolve('/artist/[artistSlug]', {
						artistSlug: artist?.slug ?? params.artistSlug ?? ''
					})}
					class="group flex min-w-0 items-center gap-2 rounded-md px-1.5 py-1 transition hover:bg-muted"
				>
					<span
						class="flex h-8 w-8 shrink-0 items-center justify-center rounded-md border border-border bg-card text-xs font-semibold text-card-foreground"
						aria-hidden="true"
					>
						{artistInitial}
					</span>
					<span class="min-w-0 truncate text-sm font-medium text-muted-foreground group-hover:text-foreground">
						{artist?.name ?? 'Artist'}
					</span>
				</a>
			</div>

			<div class="flex shrink-0 items-center gap-1.5">
				<nav class="hidden items-center gap-1 sm:flex" aria-label="Artist navigation">
					<a
						href={resolve('/artist/[artistSlug]', {
							artistSlug: artist?.slug ?? params.artistSlug ?? ''
						})}
						class={cn(
							'rounded-md px-3 py-2 text-sm font-medium transition',
							isArtistHome
								? 'bg-muted text-foreground'
								: 'text-muted-foreground hover:bg-muted hover:text-foreground'
						)}
						aria-current={isArtistHome ? 'page' : undefined}
					>
						Profile
					</a>
					<a
						href={discoverHref}
						class="rounded-md px-3 py-2 text-sm font-medium text-muted-foreground transition hover:bg-muted hover:text-foreground"
					>
						Discover
					</a>
				</nav>
				<Button
					href={discoverHref}
					variant="ghost"
					size="icon"
					class="rounded-md"
					aria-label="Search music"
				>
					<Search class="h-4 w-4" aria-hidden="true" />
				</Button>
				<DarkModeToggle class="rounded-md" />
			</div>
		</div>
	</header>

	<main class="mx-auto w-full max-w-7xl px-4 py-6 sm:px-6 lg:px-8 lg:py-10">
		{@render children?.()}
	</main>
</div>
