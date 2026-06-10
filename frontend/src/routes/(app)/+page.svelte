<script lang="ts">
	import Icon from '$lib/components/ui/icon/Icon.svelte';
	import Button from '$lib/components/ui/button/button.svelte';
	import { resolve } from '$app/paths';
	import { ReleaseCard, releaseRoute, releaseRouteParams } from '$lib/features/catalog';
	import type { PageData } from './$types';

	let { data }: { data: PageData } = $props();

	const discoverHref = resolve('/discover');
	const featuredRelease = $derived(data.featuredRelease);
	const latestReleases = $derived(data.latestReleases);
</script>

<svelte:head>
	<title>Silkwave - New Releases</title>
</svelte:head>

<div class="space-y-9">
	{#if featuredRelease}
		<section class="grid overflow-hidden rounded-3xl bg-muted/35 md:grid-cols-[minmax(0,1fr)_18rem]">
			<div class="flex min-h-80 flex-col justify-between gap-8 p-6 sm:p-8">
				<div>
					<p class="text-[0.68rem] font-semibold uppercase tracking-[0.16em] text-muted-foreground">
						Featured Release
					</p>
					<h1 class="mt-8 max-w-xl font-serif text-4xl font-light leading-tight sm:text-5xl">
						{featuredRelease.title}
					</h1>
					{#if featuredRelease.artist?.name}
						<p class="mt-3 text-sm text-muted-foreground">
							{featuredRelease.artist.name} · {featuredRelease.kind}
						</p>
					{/if}
					{#if featuredRelease.description}
						<p class="mt-5 max-w-lg line-clamp-3 text-sm leading-relaxed text-muted-foreground">
							{featuredRelease.description}
						</p>
					{/if}
				</div>

				<div class="flex flex-wrap items-center gap-3">
					{#if featuredRelease.slug}
						<Button
							href={resolve(releaseRoute, releaseRouteParams(featuredRelease))}
							variant="primary"
							class="rounded-full px-5"
						>
							View Release
						</Button>
					{/if}
					<Button href={discoverHref} variant="secondary" class="rounded-full px-5">
						Explore Catalog
					</Button>
				</div>
			</div>

			<a
				href={resolve(releaseRoute, releaseRouteParams(featuredRelease))}
				class="group relative min-h-72 overflow-hidden bg-foreground/5 md:min-h-full"
				aria-label="Open {featuredRelease.title}"
			>
				{#if featuredRelease.coverArtUrl}
					<img
						src={featuredRelease.coverArtUrl}
						alt={featuredRelease.title}
						class="h-full w-full object-cover transition duration-500 group-hover:scale-105"
					/>
				{:else}
					<div class="flex h-full min-h-72 items-center justify-center text-muted-foreground">
						<Icon icon="music-note-2" class="h-12 w-12 fill-current" />
					</div>
				{/if}
			</a>
		</section>
	{:else}
		<section class="rounded-3xl bg-muted/35 p-8">
			<p class="text-[0.68rem] font-semibold uppercase tracking-[0.16em] text-muted-foreground">
				Featured Releases
			</p>
			<h1 class="mt-8 max-w-xl font-serif text-4xl font-light leading-tight sm:text-5xl">
				Independent music, ready when the catalog is.
			</h1>
			<p class="mt-4 max-w-lg text-sm leading-relaxed text-muted-foreground">
				{data.catalogError ?? 'Published Releases will appear here as artists go live.'}
			</p>
			<Button href={discoverHref} variant="secondary" class="mt-7 rounded-full px-5">
				Browse Discover
			</Button>
		</section>
	{/if}

	<section>
		<div class="mb-4 flex items-end justify-between">
			<h2 class="text-lg font-semibold">Latest Releases</h2>
			<a
				href={discoverHref}
				class="text-[0.72rem] font-semibold text-muted-foreground transition hover:text-foreground"
			>
				View all
			</a>
		</div>

		{#if latestReleases.length}
			<div class="grid grid-cols-2 gap-x-6 gap-y-10 sm:grid-cols-3 xl:grid-cols-4">
				{#each latestReleases as release (release.key)}
					<ReleaseCard {release} />
				{/each}
			</div>
		{:else}
			<div class="rounded-2xl bg-muted/35 p-8 text-sm text-muted-foreground">
				{data.catalogError ?? 'No published Releases yet.'}
			</div>
		{/if}
	</section>
</div>
