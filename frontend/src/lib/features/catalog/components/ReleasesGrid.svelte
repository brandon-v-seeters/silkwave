<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import Icon from '$lib/components/ui/icon/Icon.svelte';
	import ReleaseCard from './ReleaseCard.svelte';

	let { releases, title, subtitle = '' } = $props();

	let container: HTMLDivElement;
	let sliderContainer: HTMLDivElement;
	let startIndex = $state(0);
	let visibleCount = $state(1);
	let resizeObserver: ResizeObserver | null = null;
	let cardWidth = $state(176);

	const gap = 24;

	function calculateVisibleCount() {
		if (!container) return;

		const containerWidth = container.offsetWidth;
		const minCardWidth = 176;

		// Calculate how many cards fit with minimum width
		const possibleCount = Math.floor((containerWidth + gap) / (minCardWidth + gap));
		visibleCount = Math.max(1, possibleCount);

		// Calculate actual card width to fill container exactly
		if (visibleCount > 0) {
			const totalGapWidth = (visibleCount - 1) * gap;
			const totalCardWidth = containerWidth - totalGapWidth;
			cardWidth = totalCardWidth / visibleCount;

			// Ensure card width is at least minCardWidth
			if (cardWidth < minCardWidth) {
				visibleCount = Math.max(1, visibleCount - 1);
				const newTotalGapWidth = (visibleCount - 1) * gap;
				const newTotalCardWidth = containerWidth - newTotalGapWidth;
				cardWidth = newTotalCardWidth / visibleCount;
			}
		}
	}

	function nextPage() {
		const maxStart = Math.max(0, releases.length - visibleCount);
		if (startIndex < maxStart) {
			startIndex = Math.min(startIndex + 1, maxStart);
		}
	}

	function prevPage() {
		if (startIndex > 0) {
			startIndex = Math.max(0, startIndex - 1);
		}
	}

	const visibleReleases = $derived(releases.slice(startIndex, startIndex + visibleCount));
	const canGoNext = $derived(startIndex + visibleCount < releases.length);
	const canGoPrev = $derived(startIndex > 0);

	onMount(() => {
		if (typeof window !== 'undefined' && container) {
			// Use requestAnimationFrame to ensure container is fully rendered
			requestAnimationFrame(() => {
				calculateVisibleCount();
			});

			resizeObserver = new ResizeObserver(() => {
				calculateVisibleCount();
			});

			resizeObserver.observe(container);
		}

		return () => {
			if (resizeObserver) {
				resizeObserver.disconnect();
			}
		};
	});

	onDestroy(() => {
		if (resizeObserver) {
			resizeObserver.disconnect();
		}
	});
</script>

<div class="relative">
	<div class="mb-4 flex flex-col gap-3 sm:flex-row sm:items-end sm:justify-between">
		<div class="min-w-0">
			{#if subtitle}
				<span class="text-base text-foreground-muted sm:text-base">{subtitle}</span>
			{/if}
			<h3 class="text-lg font-semibold">{title}</h3>
		</div>
		<div class="flex flex-shrink-0 items-center justify-end gap-2">
			<button
				type="button"
				onclick={prevPage}
				disabled={!canGoPrev}
				class="flex h-9 w-9 items-center justify-center rounded-full border border-border bg-background transition-colors hover:bg-accent disabled:cursor-not-allowed disabled:bg-transparent disabled:opacity-30 sm:h-10 sm:w-10"
				aria-label="Previous releases"
			>
				<Icon icon="chevron-left" class="h-4 w-4 fill-current sm:h-5 sm:w-5" />
			</button>
			<button
				type="button"
				onclick={nextPage}
				disabled={!canGoNext}
				class="flex h-9 w-9 items-center justify-center rounded-full border border-border bg-background transition-colors hover:bg-accent disabled:cursor-not-allowed disabled:bg-transparent disabled:opacity-30 sm:h-10 sm:w-10"
				aria-label="Next releases"
			>
				<Icon icon="chevron-right" class="h-4 w-4 fill-current sm:h-5 sm:w-5" />
			</button>
		</div>
	</div>

	<div bind:this={container} class="w-full overflow-hidden">
		<div
			bind:this={sliderContainer}
			class="flex gap-6 transition-transform duration-500 ease-in-out"
			style="transform: translateX({-startIndex * (cardWidth + gap)}px);"
		>
			{#each releases as release (release._key || release._id)}
				<div
					class="flex-shrink-0 transition-opacity duration-300 {visibleReleases.includes(
						release
					)
						? 'opacity-100'
						: 'pointer-events-none opacity-0'}"
					style="width: {cardWidth}px;"
				>
					<ReleaseCard {release} />
				</div>
			{/each}
		</div>
	</div>
</div>
