<script lang="ts">
	import { onMount } from 'svelte';
	import type { Release } from '$lib/types/generated/models';

type HeroSlide = {
	release: Partial<Release> &
		Pick<Release, 'id' | 'title' | 'slug' | 'artistKey' | 'releaseType' | 'publishAt'> & {
			artist?: { name: string; slug: string };
			coverArt?: string;
		};
	backgroundImage: string;
};

	let currentSlide = $state(0);
	let autoplayInterval: ReturnType<typeof setInterval> | null = null;

	// Sample data - replace with actual API call
	let slides = $state<HeroSlide[]>([
		{
			release: {
				_id: 'hero-1',
				_key: 'hero-1',
				_rev: '1',
				id: 'hero-1',
				title: 'Led by Ancient Light',
				slug: 'led-by-ancient-light',
				artistKey: 'koan-sound',
				releaseType: 'album',
				publishAt: new Date().toISOString(),
				createdAt: new Date().toISOString(),
				updatedAt: new Date().toISOString(),
				coverArt:
					'https://images.unsplash.com/photo-1506905925346-21bda4d32df4?w=1920&h=1080&fit=crop&q=80',
				artist: { name: 'KOAN Sound', slug: 'koan-sound' }
			},
			backgroundImage:
				'https://images.unsplash.com/photo-1506905925346-21bda4d32df4?w=1920&h=1080&fit=crop&q=80'
		},
		{
			release: {
				_id: 'hero-2',
				_key: 'hero-2',
				_rev: '1',
				id: 'hero-2',
				title: 'Midnight Echoes',
				slug: 'midnight-echoes',
				artistKey: 'neon-dreams',
				releaseType: 'ep',
				publishAt: new Date().toISOString(),
				createdAt: new Date().toISOString(),
				updatedAt: new Date().toISOString(),
				coverArt:
					'https://images.unsplash.com/photo-1470229722913-7c0e2dbbafd3?w=1920&h=1080&fit=crop&q=80',
				artist: { name: 'Neon Dreams', slug: 'neon-dreams' }
			},
			backgroundImage:
				'https://images.unsplash.com/photo-1470229722913-7c0e2dbbafd3?w=1920&h=1080&fit=crop&q=80'
		}
	]);

	onMount(() => {
		// Start autoplay
		autoplayInterval = setInterval(() => {
			currentSlide = (currentSlide + 1) % slides.length;
		}, 5000);

		return () => {
			if (autoplayInterval) {
				clearInterval(autoplayInterval);
			}
		};
	});

	function goToSlide(index: number) {
		currentSlide = index;
		if (autoplayInterval) {
			clearInterval(autoplayInterval);
			autoplayInterval = setInterval(() => {
				currentSlide = (currentSlide + 1) % slides.length;
			}, 5000);
		}
	}
</script>

<div
	class="relative flex h-48 w-full items-center justify-center overflow-hidden rounded-lg sm:h-64 sm:rounded-xl md:h-80 md:rounded-2xl lg:h-96 lg:rounded-3xl"
>
	{#each slides as slide, index}
		<div
			class="absolute inset-0 transition-opacity duration-700 ease-in-out {currentSlide ===
			index
				? 'opacity-100'
				: 'opacity-0'}"
		>
			<!-- Background Image -->
			<div
				class="absolute inset-0 bg-gradient-to-br from-orange-900 via-red-900 to-purple-900 bg-cover bg-center bg-no-repeat"
				style="background-image: url('{slide.backgroundImage ||
					slide.release.coverArt ||
					''}');"
			>
				<!-- Gradient Overlay -->
				<div
					class="absolute inset-0 bg-gradient-to-t from-black/80 via-black/40 to-transparent"
				></div>
			</div>

			<!-- Content -->
			<a
				href="/{slide.release.artist?.slug || ''}/{slide.release.slug}"
				class="relative z-10 flex h-full flex-col justify-end p-4 transition-opacity hover:opacity-95 sm:p-6 md:p-8"
			>
				<!-- Text Content (Bottom Left) -->
				<div class="max-w-2xl">
					{#if slide.release.artist}
						<p
							class=" font-light tracking-wide text-white sm:text-base md:text-base lg:text-lg"
						>
							{slide.release.artist.name}
						</p>
					{/if}
					<h2
						class="text-lg font-semibold !leading-tight text-white sm:text-lg sm:!leading-none md:text-xl lg:text-xl xl:text-xl"
					>
						{slide.release.title}
					</h2>
				</div>
			</a>
		</div>
	{/each}

	<!-- Carousel Indicators (Top Right) -->
	{#if slides.length > 1}
		<div class="absolute top-3 z-20 mx-auto flex gap-1.5 sm:top-4 sm:gap-2 md:top-6 lg:top-12">
			{#each slides as _, index}
				<button
					type="button"
					onclick={() => goToSlide(index)}
					class="h-1 rounded-full transition-all duration-300 sm:h-1.5 {currentSlide ===
					index
						? 'w-6 bg-white sm:w-8'
						: 'w-1 bg-white/50 hover:bg-white/75 sm:w-1.5'}"
					aria-label="Go to slide {index + 1}"
				></button>
			{/each}
		</div>
	{/if}
</div>
