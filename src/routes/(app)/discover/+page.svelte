<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import type { Release } from '$lib/types/generated/models';
	import Icon from '$lib/components/atoms/Icon.svelte';
	import Input from '$lib/components/ui/input/input.svelte';

	type ReleaseWithArtist = Release & { artist?: { name: string; slug: string } };

	let releases = $state<ReleaseWithArtist[]>([]);
	let isLoading = $state(true);
	let error = $state<string | null>(null);
	let searchQuery = $state($page.url.searchParams.get('q') || '');

	async function loadReleases(query: string = '') {
		isLoading = true;
		error = null;

		try {
			const url = query
				? `/api/releases?q=${encodeURIComponent(query)}`
				: '/api/releases?limit=24';
			const response = await fetch(url);
			const data = await response.json();

			if (data.error) {
				error = data.error;
			} else {
				releases = data.releases || [];
			}
		} catch (e) {
			error = 'Failed to load releases';
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

	function formatDate(timestamp: number) {
		return new Date(timestamp).toLocaleDateString('en-US', {
			year: 'numeric',
			month: 'long',
			day: 'numeric'
		});
	}
</script>

<svelte:head>
	<title>Discover - Silk Wave</title>
</svelte:head>

<div class="mx-auto max-w-7xl">
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
		<div class="bg-muted rounded-lg p-12 text-center">
			<Icon icon="music-note-2" class="mx-auto h-12 w-12 text-foreground-muted" />
			<p class="mt-4 text-lg text-foreground-muted">
				{searchQuery ? 'No releases found matching your search' : 'No releases found'}
			</p>
		</div>
	{:else}
		<!-- Releases Grid -->
		<div class="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
			{#each releases as release}
				<a
					href="/{release.artist?.slug}/{release.slug}"
					class="group rounded-lg border border-border bg-card p-4 transition-all hover:border-primary hover:shadow-lg"
				>
					<!-- Cover Art Placeholder -->
					<div class="bg-muted mb-4 aspect-square w-full overflow-hidden rounded-md">
						{#if release.coverArt}
							<img
								src={release.coverArt}
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
								{release.type}
							</span>
							<span>•</span>
							<span>{formatDate(release.releaseDate)}</span>
						</div>
					</div>
				</a>
			{/each}
		</div>
	{/if}
</div>
