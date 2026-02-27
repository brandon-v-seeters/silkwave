<script lang="ts">
	import { onMount } from 'svelte';
	import type { Release } from '$lib/types/generated/models';
	import Icon from '$lib/components/atoms/Icon.svelte';
	import Button from '$lib/components/ui/button/button.svelte';
	import HeroCarousel from '$lib/components/organisms/HeroCarousel.svelte';
	import ReleasesGrid from '$lib/components/organisms/ReleasesGrid.svelte';

	type ReleaseWithArtist = Release & { artist?: { name: string; slug: string } };

	const releasesMock = [
		{
			_id: '1',
			_key: '1',
			_rev: '1',
			title: 'Led by Ancient Light',
			slug: 'led-by-ancient-light',
			artistKey: 'koan-sound',
			type: 'album',
			releaseDate: Date.now(),
			createdAt: Date.now(),
			updatedAt: Date.now(),
			published: true,
			artist: { name: 'KOAN Sound', slug: 'koan-sound' },
			coverArt:
				'https://images.unsplash.com/photo-1493225457124-a3eb161ffa5f?w=600&h=600&fit=crop&q=80'
		},
		{
			_id: '2',
			_key: '2',
			_rev: '2',
			title: 'Midnight Echoes',
			slug: 'midnight-echoes',
			artistKey: 'neon-dreams',
			type: 'ep',
			releaseDate: Date.now(),
			createdAt: Date.now(),
			updatedAt: Date.now(),
			published: true,
			artist: { name: 'Neon Dreams', slug: 'neon-dreams' },
			coverArt:
				'https://images.unsplash.com/photo-1459749411175-04bf5292ceea?w=600&h=600&fit=crop&q=80'
		},
		{
			_id: '3',
			_key: '3',
			_rev: '3',
			title: 'Electric Dreams',
			slug: 'electric-dreams',
			artistKey: 'synth-wave',
			type: 'album',
			releaseDate: Date.now(),
			createdAt: Date.now(),
			updatedAt: Date.now(),
			published: true,
			artist: { name: 'Synth Wave', slug: 'synth-wave' },
			coverArt:
				'https://images.unsplash.com/photo-1514525253161-7a46d19cd819?w=600&h=600&fit=crop&q=80'
		},
		{
			_id: '4',
			_key: '4',
			_rev: '4',
			title: 'Urban Nights',
			slug: 'urban-nights',
			artistKey: 'city-vibes',
			type: 'single',
			releaseDate: Date.now(),
			createdAt: Date.now(),
			updatedAt: Date.now(),
			published: true,
			artist: { name: 'City Vibes', slug: 'city-vibes' },
			coverArt:
				'https://images.unsplash.com/photo-1493225457124-a3eb161ffa5f?w=600&h=600&fit=crop&q=80'
		},
		{
			_id: '5',
			_key: '5',
			_rev: '5',
			title: 'Ocean Breeze',
			slug: 'ocean-breeze',
			artistKey: 'coastal-sounds',
			type: 'ep',
			releaseDate: Date.now(),
			createdAt: Date.now(),
			updatedAt: Date.now(),
			published: true,
			artist: { name: 'Coastal Sounds', slug: 'coastal-sounds' },
			coverArt:
				'https://images.unsplash.com/photo-1514320291840-2e0a9bf2a9ae?w=600&h=600&fit=crop&q=80'
		},
		{
			_id: '6',
			_key: '6',
			_rev: '6',
			title: 'Cosmic Journey',
			slug: 'cosmic-journey',
			artistKey: 'space-explorer',
			type: 'album',
			releaseDate: Date.now(),
			createdAt: Date.now(),
			updatedAt: Date.now(),
			published: true,
			artist: { name: 'Space Explorer', slug: 'space-explorer' },
			coverArt:
				'https://images.unsplash.com/photo-1446776653964-20c1d3a81b06?w=600&h=600&fit=crop&q=80'
		}
	];

	let releases = $state<ReleaseWithArtist[]>(releasesMock as ReleaseWithArtist[]);
	let isLoading = $state(true);
	let error = $state<string | null>(null);

	onMount(async () => {
		// try {
		// 	const response = await fetch('/api/releases?limit=12');
		// 	const data = await response.json();
		// 	if (data.error) {
		// 		error = data.error;
		// 	} else {
		// 		releases = data.releases || releasesMock;
		// 	}
		// } catch (e) {
		// 	error = 'Failed to load releases';
		// 	console.error(e);
		// } finally {
		// 	isLoading = false;
		// }
		releases = releasesMock as ReleaseWithArtist[];
		isLoading = false;
	});

	function formatDate(timestamp: number) {
		return new Date(timestamp).toLocaleDateString('en-US', {
			year: 'numeric',
			month: 'long',
			day: 'numeric'
		});
	}
</script>

<svelte:head>
	<title>Silk Wave - Latest Releases</title>
</svelte:head>

{#if import.meta.env.PROD}
	<div class="flex h-screen w-full flex-1 flex-col items-center justify-center bg-background">
		<h1 class="text-4xl font-bold">Silk Wave</h1>
		<p class="text-lg text-foreground-muted">Coming soon...</p>
	</div>
{:else}
	<div class="space-y-6 sm:space-y-8 md:space-y-12">
		<!-- Hero Carousel -->
		<div class="mb-4 sm:mb-6 md:mb-8">
			<HeroCarousel />
		</div>

		<!-- Fresh Picks Section -->
		<div class="mb-6 sm:mb-8 md:mb-12">
			<ReleasesGrid {releases} title="Fresh picks for you" subtitle="Based on your taste" />
		</div>

		<!-- Artists You Might Like Section -->
		<div class="mb-6 sm:mb-8 md:mb-12">
			<ReleasesGrid
				{releases}
				title="Artists you might like"
				subtitle="Matching your favorites"
			/>
		</div>

		<!-- Load More Button -->
		{#if releases.length >= 12}
			<div class="mt-4 text-center sm:mt-6 md:mt-8">
				<Button href="/discover" variant="outline" class="w-full sm:w-auto">
					View More Releases
				</Button>
			</div>
		{/if}
	</div>
{/if}
