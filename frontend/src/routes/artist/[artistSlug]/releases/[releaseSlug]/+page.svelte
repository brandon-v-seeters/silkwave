<script lang="ts">
	import { browser } from '$app/environment';
	import { resolve } from '$app/paths';
	import { MediaPlayer } from '$lib/components/organisms/media-player';
	import Button from '$lib/components/ui/button/button.svelte';
	import { cart } from '$lib/stores/cart';
	import type { MediaPlayerTrack } from '$lib/components/organisms/media-player';
	import type { Track } from '$lib/types/generated/models.js';
	import {
		ArrowLeft,
		CalendarDays,
		Clock3,
		Disc3,
		Download,
		ListMusic,
		Music2,
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

	let { data, params }: PageProps = $props();

	const release = $derived(data.release);
	const artist = $derived(release.artist ?? data.artist);
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
		release.metadata?.genres?.length ? release.metadata.genres : release.genres ?? []
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
</script>

<svelte:head>
	<title>{release.title} by {artist?.name ?? 'Silkwave'} - Silkwave</title>
	<meta
		name="description"
		content={release.description || `${release.title} on Silkwave`}
	/>
</svelte:head>

<article class="pb-32">
	<div class="mb-10">
		<a
			href={resolve('/artist/[artistSlug]', {
				artistSlug: artist?.slug ?? params.artistSlug
			})}
			class="inline-flex items-center gap-2 text-sm font-medium text-muted-foreground transition hover:text-foreground"
		>
			<ArrowLeft class="h-4 w-4" aria-hidden="true" />
			<span>{artist?.name ?? 'Back to artist'}</span>
		</a>
	</div>

	<section
		class="grid gap-10 lg:grid-cols-[minmax(240px,360px)_minmax(0,1fr)] lg:items-end lg:gap-16"
	>
		<aside class="lg:sticky lg:top-24 lg:self-start">
			<div class="relative isolate aspect-square overflow-hidden rounded-lg bg-card">
				{#if coverArt}
					<img
						src={coverArt}
						alt=""
						class="absolute -inset-8 -z-10 h-[calc(100%+4rem)] w-[calc(100%+4rem)] object-cover opacity-[0.18] blur-2xl saturate-125"
						aria-hidden="true"
					/>
				{/if}
				{#if coverArt}
					<img src={coverArt} alt="{release.title} cover art" class="h-full w-full object-cover" />
				{:else}
					<div class="flex h-full w-full items-center justify-center bg-muted text-muted-foreground">
						<Music2 class="h-16 w-16" aria-hidden="true" />
					</div>
				{/if}
			</div>

			<dl class="mt-6 space-y-3 text-sm">
				<div class="flex items-start justify-between gap-4">
					<dt class="flex items-center gap-2 text-muted-foreground">
						<Disc3 class="h-4 w-4" aria-hidden="true" />
						Type
					</dt>
					<dd class="font-medium">{releaseType}</dd>
				</div>
				{#if publishDate}
					<div class="flex items-start justify-between gap-4">
						<dt class="flex items-center gap-2 text-muted-foreground">
							<CalendarDays class="h-4 w-4" aria-hidden="true" />
							Released
						</dt>
						<dd class="text-right font-medium">{publishDate}</dd>
					</div>
				{/if}
				{#if totalDuration}
					<div class="flex items-start justify-between gap-4">
						<dt class="flex items-center gap-2 text-muted-foreground">
							<Clock3 class="h-4 w-4" aria-hidden="true" />
							Length
						</dt>
						<dd class="font-medium">{formatTotalDuration(totalDuration)}</dd>
					</div>
				{/if}
				<div class="flex items-start justify-between gap-4">
					<dt class="flex items-center gap-2 text-muted-foreground">
						<Download class="h-4 w-4" aria-hidden="true" />
						Access
					</dt>
					<dd class="text-right font-medium">{availabilityLabel}</dd>
				</div>
			</dl>
		</aside>

		<div class="min-w-0">
			<div class="flex flex-wrap items-center gap-x-3 gap-y-1 text-sm text-muted-foreground">
				{#each primaryMeta as item, index (`meta-${item}`)}
					<span>
						{item}
					</span>
					{#if index < primaryMeta.length - 1}
						<span aria-hidden="true">/</span>
					{/if}
				{/each}
			</div>

			<h1 class="mt-5 max-w-4xl text-balance font-serif text-6xl font-light leading-[0.92] tracking-tight sm:text-7xl lg:text-8xl">
				{release.title}
			</h1>

			{#if artist}
				<a
					href={resolve('/artist/[artistSlug]', {
						artistSlug: artist?.slug ?? params.artistSlug
					})}
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
			{/if}

			<div class="mt-8 flex flex-col gap-3 sm:flex-row sm:items-center">
				<Button
					variant="primary"
					class="justify-center rounded-md px-6 sm:w-auto"
					disabled={isInCart}
					onclick={addReleaseToCart}
				>
					<ShoppingBag class="h-4 w-4" aria-hidden="true" />
					{buyLabel}
				</Button>

				{#if hasTracks}
					<Button href="#tracks" variant="secondary" class="justify-center rounded-md px-5">
						<ListMusic class="h-4 w-4" aria-hidden="true" />
						View tracks
					</Button>
				{/if}

				<Button variant="ghost" class="justify-center rounded-md px-4" onclick={shareRelease}>
					<Share2 class="h-4 w-4" aria-hidden="true" />
					Share
				</Button>
			</div>

			{#if release.description}
				<p class="mt-9 max-w-2xl whitespace-pre-line text-pretty text-lg leading-8 text-muted-foreground">
					{release.description}
				</p>
			{/if}

			{#if genres.length}
				<div class="mt-6 flex flex-wrap gap-2">
					{#each genres as genre (`genre-${genre}`)}
						<span class="rounded-md bg-muted px-3 py-1 text-sm text-muted-foreground">
							{genre}
						</span>
					{/each}
				</div>
			{/if}
		</div>
	</section>

	<section
		id="tracks"
		class="mt-18 grid gap-8 lg:grid-cols-[minmax(0,1fr)_minmax(220px,300px)] lg:gap-14"
	>
		<div>
			<div class="flex flex-col gap-2 sm:flex-row sm:items-end sm:justify-between">
				<div>
					<p class="flex items-center gap-2 text-sm font-medium text-muted-foreground">
						<ListMusic class="h-4 w-4" aria-hidden="true" />
						Tracks
					</p>
					<h2 class="mt-2 text-3xl font-semibold tracking-tight">Side A to closing note</h2>
				</div>
				{#if totalDuration}
					<p class="text-sm text-muted-foreground">{formatTotalDuration(totalDuration)}</p>
				{/if}
			</div>

			{#if hasTracks}
				<ol class="mt-6 space-y-1">
					{#if sortedTracks.length}
						{#each sortedTracks as track, index (track.id || track._key || `${index}-${track.title}`)}
							{@const source = resolveTrackSource(track)}
							{@const duration = formatDuration(track.duration, track.durationDisplay)}
							<li class="group grid grid-cols-[2.25rem_minmax(0,1fr)_auto] items-center gap-3 rounded-md px-2 py-3 transition hover:bg-muted/50 sm:grid-cols-[2.25rem_minmax(0,1fr)_auto_auto] sm:px-3">
								<span class="text-sm tabular-nums text-muted-foreground">
									{track.order || index + 1}
								</span>
								<div class="min-w-0">
									<p class="truncate font-medium">{track.title}</p>
									{#if track.isExplicit}
										<p class="mt-1 text-xs text-muted-foreground">Explicit</p>
									{/if}
								</div>
								{#if duration}
									<span class="text-sm tabular-nums text-muted-foreground">{duration}</span>
								{/if}
								<span class="hidden text-sm text-muted-foreground sm:inline">
									{source ? 'Preview in player' : 'No preview'}
								</span>
							</li>
						{/each}
					{:else}
						{#each trackTitles as title, index (`${index}-${title}`)}
							<li class="grid grid-cols-[2.25rem_minmax(0,1fr)] items-center gap-3 rounded-md px-2 py-3 transition hover:bg-muted/50 sm:px-3">
								<span class="text-sm tabular-nums text-muted-foreground">{index + 1}</span>
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

		<aside class="text-sm">
			<h2 class="text-base font-semibold tracking-tight">Release notes</h2>
			<dl class="mt-5 space-y-5">
				<div>
					<dt class="text-muted-foreground">Artist</dt>
					<dd class="mt-1 truncate font-medium">{artist?.name ?? 'Unknown artist'}</dd>
				</div>
				<div>
					<dt class="text-muted-foreground">Availability</dt>
					<dd class="mt-1 font-medium">{availabilityLabel}</dd>
				</div>
				{#if release.metadata?.label}
					<div>
						<dt class="text-muted-foreground">Label</dt>
						<dd class="mt-1 truncate font-medium">{release.metadata.label}</dd>
					</div>
				{/if}
				{#if release.metadata?.catalogNumber}
					<div>
						<dt class="text-muted-foreground">Catalog</dt>
						<dd class="mt-1 truncate font-medium">{release.metadata.catalogNumber}</dd>
					</div>
				{/if}
				{#if release.isExplicit}
					<div>
						<dt class="text-muted-foreground">Content</dt>
						<dd class="mt-1 font-medium">Explicit</dd>
					</div>
				{/if}
			</dl>
		</aside>
	</section>

	{#if hasPlayableTracks}
		<MediaPlayer
			tracks={playableTracks}
			initialTrackId={playableTracks[0]?.id}
			releaseTitle={release.title}
			artistName={artist?.name}
			artworkUrl={coverArt}
			placement="floating"
			label={`${release.title} media player`}
		/>
	{/if}
</article>
