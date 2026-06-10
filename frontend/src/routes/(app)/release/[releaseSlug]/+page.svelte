<script lang="ts">
	import { browser } from '$app/environment';
	import { resolve } from '$app/paths';
	import { onDestroy } from 'svelte';
	import { MediaPlayer } from '$lib/features/playback';
	import Button from '$lib/components/ui/button/button.svelte';
	import { cart } from '$lib/features/cart';
	import type { MediaPlayerTrack } from '$lib/features/playback';
	import type { Track } from '$lib/types/generated/models.js';
	import {
		ArrowLeft,
		CalendarDays,
		Clock3,
		Disc3,
		Download,
		Heart,
		ListMusic,
		Music2,
		Play,
		Share2,
		ShoppingBag
	} from '@lucide/svelte';
	import { toast } from 'svelte-sonner';
	import type { PageProps } from './$types';

	type TrackWithOptionalSource = Track & {
		source?: string | null;
		audioUrl?: string | null;
		previewUrl?: string | null;
		streamUrl?: string | null;
		artworkUrl?: string | null;
	};

	let { data }: PageProps = $props();
	let previewAudio: HTMLAudioElement | null = null;

	const release = $derived(data.release);
	const artist = $derived(release.artist);
	const tracks = $derived<TrackWithOptionalSource[]>(
		(release.tracks ?? []) as TrackWithOptionalSource[]
	);
	const sortedTracks = $derived.by(() =>
		[...tracks].sort((first, second) => (first.order || 0) - (second.order || 0))
	);
	const trackTitles = $derived(release.trackList ?? []);
	const hasTracks = $derived(sortedTracks.length > 0 || trackTitles.length > 0);
	const coverArt = $derived(
		release.cover ||
			release.assets?.coverArt?.original ||
			release.assets?.coverArt?.medium ||
			release.assets?.coverArt?.thumbnail ||
			null
	);
	const publishDate = $derived(formatDate(release.publishAt || release.schedule?.publishAt));
	const price = $derived(formatPrice(release.pricing?.basePrice));
	const trackCount = $derived(
		release.trackCount || sortedTracks.length || trackTitles.length || 0
	);
	const totalDuration = $derived(formatDuration(release.totalDuration));
	const releaseType = $derived(formatReleaseType(release.releaseType));
	const playableTracks = $derived(buildPlayableTracks(sortedTracks));
	const hasPlayableTracks = $derived(playableTracks.length > 0);
	const genres = $derived(
		release.metadata?.genres?.length ? release.metadata.genres : (release.genres ?? [])
	);
	const cartItemId = $derived(`release:${release.id || release._key || release.slug}`);
	const isInCart = $derived($cart.some((item) => item.id === cartItemId));
	const buyLabel = $derived(
		isInCart
			? 'In cart'
			: price === 'Free'
				? 'Add free release'
				: price
					? `Add to cart - ${price}`
					: 'Add to cart'
	);
	const availabilityLabel = $derived(
		release.streamingEnabled && release.downloadEnabled
			? 'Streaming and downloads'
			: release.streamingEnabled
				? 'Streaming only'
				: release.downloadEnabled
					? 'Downloads included'
					: 'Limited availability'
	);
	const primaryMeta = $derived(
		[
			releaseType,
			publishDate,
			trackCount ? `${trackCount} ${trackCount === 1 ? 'track' : 'tracks'}` : null,
			totalDuration ? formatTotalDuration(totalDuration) : null
		].filter(Boolean)
	);

	function formatReleaseType(value: string | undefined) {
		if (!value) return 'Release';

		return value
			.split(/[-_\s]+/)
			.filter(Boolean)
			.map((part) => part.charAt(0).toUpperCase() + part.slice(1).toLowerCase())
			.join(' ');
	}

	function formatDate(value: string | undefined) {
		if (!value) return null;

		const date = new Date(value);
		if (Number.isNaN(date.getTime())) return null;

		return date.toLocaleDateString('en-US', {
			year: 'numeric',
			month: 'long',
			day: 'numeric'
		});
	}

	function formatPrice(value: number | undefined) {
		if (value === undefined || value === null) return null;
		if (value === 0) return 'Free';

		return new Intl.NumberFormat('en-US', {
			style: 'currency',
			currency: 'USD'
		}).format(value / 100);
	}

	function formatDuration(seconds: number | undefined | null, fallback?: string | null) {
		if (fallback) return fallback;
		if (!seconds || !Number.isFinite(seconds)) return null;

		const safeSeconds = Math.max(0, Math.floor(seconds));
		const minutes = Math.floor(safeSeconds / 60);
		const remainingSeconds = safeSeconds % 60;

		return `${minutes}:${remainingSeconds.toString().padStart(2, '0')}`;
	}

	function formatTotalDuration(value: string | null) {
		if (!value) return null;

		const [minutesPart, secondsPart] = value.split(':');
		const minutes = Number(minutesPart);
		const seconds = Number(secondsPart);

		if (!Number.isFinite(minutes) || !Number.isFinite(seconds)) return value;
		if (minutes < 60) return `${minutes} min ${seconds.toString().padStart(2, '0')} sec`;

		const hours = Math.floor(minutes / 60);
		const remainingMinutes = minutes % 60;

		return `${hours} hr ${remainingMinutes} min`;
	}

	function resolveTrackSource(track: TrackWithOptionalSource) {
		return (
			track.source ||
			track.audioUrl ||
			track.previewUrl ||
			track.streamUrl ||
			track.files?.previewPath ||
			null
		);
	}

	function buildPlayableTracks(items: TrackWithOptionalSource[]): MediaPlayerTrack[] {
		const playerTracks: MediaPlayerTrack[] = [];

		for (const [index, track] of items.entries()) {
			const source = resolveTrackSource(track);
			if (!source) continue;

			playerTracks.push({
				id: track.id || track._key || String(index + 1),
				title: track.title,
				artistName: artist?.name,
				source,
				artworkUrl: track.artworkUrl || coverArt,
				duration: track.duration
			});
		}

		return playerTracks;
	}

	function addReleaseToCart() {
		if (isInCart) return;

		cart.add({
			id: cartItemId,
			type: release.releaseType?.toLowerCase() === 'single' ? 'track' : 'album',
			title: release.title,
			artist: artist?.name ?? 'Unknown artist',
			price: release.pricing?.basePrice ?? 0,
			coverUrl: coverArt ?? undefined
		});

		toast.success('Added to cart');
	}

	function playTrack(track: TrackWithOptionalSource) {
		if (!browser) return;

		const source = resolveTrackSource(track);
		if (!source) return;

		previewAudio ??= new Audio();
		previewAudio.pause();
		previewAudio.src = source;
		previewAudio.currentTime = 0;
		void previewAudio.play().catch(() => {});
	}

	async function shareRelease() {
		if (!browser) return;

		const url = window.location.href;

		try {
			if (navigator.share) {
				await navigator.share({
					title: release.title,
					text: artist?.name ? `${release.title} by ${artist.name}` : release.title,
					url
				});
				return;
			}

			await navigator.clipboard.writeText(url);
			toast.success('Link copied');
		} catch (error) {
			if (error instanceof DOMException && error.name === 'AbortError') return;

			toast.error('Could not share this release');
		}
	}

	onDestroy(() => {
		previewAudio?.pause();
	});
</script>

<svelte:head>
	<title>{release.title} by {artist?.name ?? 'Silkwave'} - Silkwave</title>
	<meta name="description" content={release.description || `${release.title} on Silkwave`} />
</svelte:head>

<article class="pb-32">
	<section
		class="grid gap-10 lg:grid-cols-[minmax(240px,360px)_minmax(0,1fr)] lg:items-end lg:gap-16"
	>
		<aside class="lg:sticky lg:top-24 lg:self-start">
			<div class="relative isolate aspect-square overflow-hidden rounded-lg bg-card">
				{#if coverArt}
					<img
						src={coverArt}
						alt="{release.title} cover art"
						class="h-full w-full object-cover"
					/>
				{:else}
					<div
						class="flex h-full w-full items-center justify-center bg-muted text-muted-foreground"
					>
						<Music2 class="h-16 w-16" aria-hidden="true" />
					</div>
				{/if}
			</div>
		</aside>

		<div class="min-w-0">
			<h1
				class="mt-5 max-w-4xl text-balance font-serif text-6xl font-bold leading-[0.92] tracking-tight sm:text-7xl lg:text-8xl"
			>
				{release.title}
			</h1>

			{#if artist?.slug}
				<a
					href={resolve('/(app)/artist/[artistSlug]', { artistSlug: artist.slug })}
					class="mt-5 inline-flex max-w-full items-center gap-3 text-base font-medium text-muted-foreground transition hover:text-primary"
				>
					<span
						class="flex h-10 w-10 shrink-0 items-center justify-center rounded-full bg-card text-sm font-semibold text-card-foreground"
						aria-hidden="true"
					>
						{artist.name?.slice(0, 1).toUpperCase() ?? 'A'}
					</span>
					<span class="min-w-0 truncate">{artist.name}</span>
				</a>
			{:else if artist}
				<div
					class="mt-5 inline-flex max-w-full items-center gap-3 text-base font-medium text-muted-foreground"
				>
					<span
						class="flex h-10 w-10 shrink-0 items-center justify-center rounded-full bg-card text-sm font-semibold text-card-foreground"
						aria-hidden="true"
					>
						{artist.name?.slice(0, 1).toUpperCase() ?? 'A'}
					</span>
					<span class="min-w-0 truncate">{artist.name}</span>
				</div>
			{/if}

			<div class="mt-8 flex flex-col gap-3 sm:flex-row sm:items-center">
				<Button
					class="justify-center sm:w-auto"
					disabled={isInCart}
					size="lg"
					onclick={addReleaseToCart}
				>
					<ShoppingBag class="h-4 w-4" aria-hidden="true" />
					{buyLabel}
				</Button>
				<Button
					variant="ghost"
					class="justify-center rounded-md px-4"
					onclick={shareRelease}
				>
					<Heart class="h-4 w-4" aria-hidden="true" />
				</Button>
				<Button
					variant="ghost"
					class="justify-center rounded-md px-4"
					onclick={shareRelease}
				>
					<Share2 class="h-4 w-4" aria-hidden="true" />
				</Button>
			</div>

			{#if release.description}
				<p
					class="mt-9 max-w-2xl whitespace-pre-line text-pretty text-lg leading-8 text-muted-foreground"
				>
					{release.description}
				</p>
			{/if}
			<dl class="mt-6 space-y-3 text-sm">
				{#if publishDate}
					<div class="flex items-center gap-4">
						<dd class="text-right font-light text-muted-foreground">{publishDate}</dd>
					</div>
				{/if}
			</dl>

			{#if genres.length}
				<div class="mt-6 flex flex-wrap gap-2">
					{#each genres as genre (`genre-${genre}`)}
						<span
							class="rounded-full leading-0 bg-foreground/5 px-6 py-6 text-sm text-muted-foreground hover:bg-foreground/10 hover:text-foreground transition duration-200 cursor-pointer"
						>
							{genre}
						</span>
					{/each}
				</div>
			{/if}
		</div>
	</section>

	<section id="tracks" class="mt-18 grid gap-8 lg:grid-cols-[minmax(0,1fr)] lg:gap-14">
		<div>
			<div class="flex flex-col gap-2 sm:flex-row sm:items-end sm:justify-between">
				<div>
					<div class="flex items-end gap-2 text-sm font-medium text-muted-foreground">
						<ListMusic class="h-4 w-4" aria-hidden="true" />
						{trackCount} Tracks
						<span class="text-xs text-muted-foreground leading-0">
							{formatTotalDuration(totalDuration)}
						</span>
					</div>
				</div>
			</div>

			{#if hasTracks}
				<ol class="mt-6 space-y-1">
					{#if sortedTracks.length}
						{#each sortedTracks as track, index (track.id || track._key || `${index}-${track.title}`)}
							{@const source = resolveTrackSource(track)}
							{@const duration = formatDuration(
								track.duration,
								track.durationDisplay
							)}
							<li
								class="group grid grid-cols-[2.25rem_minmax(0,1fr)_auto] items-center gap-3 rounded-md px-2 py-3 transition hover:bg-muted/50 sm:grid-cols-[2.25rem_minmax(0,1fr)_auto_auto] sm:px-3"
							>
								<div class="relative flex h-7 items-center">
									<span
										class="text-sm tabular-nums text-muted-foreground transition-opacity duration-150 {source
											? 'group-hover:opacity-0 group-focus-within:opacity-0'
											: ''}"
									>
										{track.order || index + 1}
									</span>

									{#if source}
										<button
											type="button"
											onclick={() => playTrack(track)}
											class="absolute -left-1 flex h-7 w-7 items-center justify-center rounded-full text-foreground opacity-0 transition duration-150 hover:bg-foreground/8 focus-visible:bg-foreground/8 focus-visible:opacity-100 focus-visible:outline-none group-hover:opacity-100 group-focus-within:opacity-100"
											aria-label="Play {track.title}"
										>
											<Play
												class="h-3.5 w-3.5 fill-current"
												aria-hidden="true"
											/>
										</button>
									{/if}
								</div>
								<div class="min-w-0">
									<p class="truncate font-medium">{track.title}</p>
									{#if track.isExplicit}
										<p class="mt-1 text-xs text-muted-foreground">Explicit</p>
									{/if}
								</div>
								{#if duration}
									<span class="text-sm tabular-nums text-muted-foreground"
										>{duration}</span
									>
								{/if}
							</li>
						{/each}
					{:else}
						{#each trackTitles as title, index (`${index}-${title}`)}
							<li
								class="grid grid-cols-[2.25rem_minmax(0,1fr)] items-center gap-3 rounded-md px-2 py-3 transition hover:bg-muted/50 sm:px-3"
							>
								<span class="text-sm tabular-nums text-muted-foreground"
									>{index + 1}</span
								>
								<span class="min-w-0 truncate font-medium">{title}</span>
							</li>
						{/each}
					{/if}
				</ol>
			{:else}
				<div class="mt-6 rounded-lg bg-card px-5 py-10 text-muted-foreground">
					The track list is still being finished.
				</div>
			{/if}
		</div>
	</section>
</article>
