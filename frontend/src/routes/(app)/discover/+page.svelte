<script lang="ts">
	import { onMount } from 'svelte';
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import {
		fetchCatalogReleases,
		hasReleaseRoute,
		releaseRoute,
		releaseRouteParams,
		type CatalogRelease
	} from '$lib/features/catalog';
	import Icon from '$lib/components/ui/icon/Icon.svelte';
	import Input from '$lib/components/ui/input/input.svelte';

	let releases = $state.raw<CatalogRelease[]>([]);
	let isLoading = $state(true);
	let error = $state<string | null>(null);
	let searchQuery = $state(page.url.searchParams.get('q') || '');

	async function loadReleases(query: string = '') {
		isLoading = true;
		error = null;

		try {
			releases = await fetchCatalogReleases(fetch, {
				limit: 24,
				query
			});
		} catch (e) {
			error = e instanceof Error ? e.message : 'Failed to load releases';
			console.error(e);
		} finally {
			isLoading = false;
		}
	}

	onMount(() => {
		loadReleases(searchQuery);
	});

	function handleSearch() {
		if (searchQuery.trim()) {
			loadReleases(searchQuery.trim());
			window.history.pushState(
				{},
				'',
				`/discover?q=${encodeURIComponent(searchQuery.trim())}`
			);
		} else {
			loadReleases();
			window.history.pushState({}, '', '/discover');
		}
	}

	function handleSearchKeydown(e: KeyboardEvent) {
		if (e.key === 'Enter') {
			handleSearch();
		}
	}

</script>

<svelte:head>
	<title>Discover - Silk Wave</title>
</svelte:head>

<div class="mx-auto">
	<!-- Header -->
	<div class="mb-8">
		<h1 class="font-serif text-4xl font-light md:text-6xl">Latest Releases</h1>
		<p class="mt-4 text-lg text-foreground-muted">
			Discover the newest music from independent artists
		</p>
	</div>

	<!-- Search Bar -->
	<div class="mb-8">
		<div class="relative max-w-md">
			<Icon
				icon="search"
				class="absolute left-3 top-1/2 h-5 w-5 -translate-y-1/2 text-foreground-muted"
			/>
			<Input
				bind:value={searchQuery}
				type="search"
				placeholder="Search releases, artists..."
				class="w-full pl-10 pr-4"
				onkeydown={handleSearchKeydown}
			/>
		</div>
	</div>

	<!-- Loading State -->
	{#if isLoading}
		<div class="flex items-center justify-center py-12">
			<Icon icon="loader-2" class="h-8 w-8 animate-spin text-foreground-muted" />
		</div>
	{:else if error}
		<!-- Error State -->
		<div class="rounded-lg bg-rose-500/10 p-6 text-center text-rose-400">
			<p>{error}</p>
		</div>
	{:else if releases.length === 0}
		<!-- Empty State -->
		<div class="bg-background rounded-lg p-12 text-center">
			<Icon icon="music-note-2" class="mx-auto h-12 w-12 text-foreground-muted" />
			<p class="mt-4 text-lg text-foreground-muted">
				{searchQuery ? 'No releases found matching your search' : 'No releases found'}
			</p>
		</div>
	{:else}
		<!-- Releases Grid -->
		<div class="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
			{#each releases as release (release.key)}
				{@const coverArt = release.coverArtUrl}
				{#snippet releaseContent()}
					<!-- Cover Art Placeholder -->
					<div class="bg-muted mb-4 aspect-square w-full overflow-hidden rounded-md">
						{#if coverArt}
							<img
								src={coverArt}
								alt={release.title}
								class="h-full w-full object-cover transition-transform group-hover:scale-105"
							/>
						{:else}
							<div
								class="flex h-full w-full items-center justify-center text-foreground-muted"
							>
								<Icon icon="music-note-2" class="h-12 w-12" />
							</div>
						{/if}
					</div>

					<!-- Release Info -->
					<div>
						<h3 class="line-clamp-1 font-medium group-hover:text-primary">
							{release.title}
						</h3>
						{#if release.artist}
							<p class="mt-1 text-base text-foreground-muted">
								{release.artist.name}
							</p>
						{/if}
						<div class="mt-2 flex items-center gap-2 text-xs text-foreground-muted">
							<span class="bg-muted rounded-full px-2 py-1">
								{release.kind}
							</span>
							<span>•</span>
							<span>{release.publishedAtLabel}</span>
						</div>
					</div>
				{/snippet}

				{#if hasReleaseRoute(release)}
					<a
						href={resolve(releaseRoute, releaseRouteParams(release))}
						class="group rounded-lg border border-border bg-card p-4 transition-all hover:border-primary hover:shadow-lg"
					>
						{@render releaseContent()}
					</a>
				{:else}
					<div class="group rounded-lg border border-border bg-card p-4">
						{@render releaseContent()}
					</div>
				{/if}
			{/each}
		</div>
	{/if}
</div>
