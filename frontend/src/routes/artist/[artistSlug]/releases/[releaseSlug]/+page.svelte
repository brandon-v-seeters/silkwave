<script lang="ts">
	import Icon from '$lib/components/atoms/Icon.svelte';
	import { MediaPlayer } from '$lib/components/organisms/media-player';
	import Button from '$lib/components/ui/button/button.svelte';
	import type { MediaPlayerTrack } from '$lib/components/organisms/media-player';
	import type { Track } from '$lib/types/generated/models.js';
	import type { PageProps } from './$types';

	type TrackWithOptionalSource = Track & {
		source?: string | null;
		audioUrl?: string | null;
		previewUrl?: string | null;
		streamUrl?: string | null;
		artworkUrl?: string | null;
	};

	let { data }: PageProps = $props();

	const release = $derived(data.release);
	const artist = $derived(release.artist ?? data.artist);
	const tracks = $derived<TrackWithOptionalSource[]>(
		(release.tracks ?? []) as TrackWithOptionalSource[]
	);
	const trackTitles = $derived(release.trackList ?? []);
	const hasTracks = $derived(tracks.length > 0 || trackTitles.length > 0);
	const artistHref = $derived(artist?.slug ? `/artist/${artist.slug}` : '/discover');
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
		release.trackCount || tracks.length || trackTitles.length || 0
	);
	const totalDuration = $derived(formatDuration(release.totalDuration));
	const buyLabel = $derived(
		release.releaseType?.toLowerCase() === 'album' ? 'Buy Album' : 'Buy Release'
	);
	const playableTracks = $derived(buildPlayableTracks(tracks));
	const hasPlayableTracks = $derived(playableTracks.length > 0);
	const primaryMeta = $derived(
		[
			formatReleaseType(release.releaseType),
			publishDate,
			trackCount ? `${trackCount} ${trackCount === 1 ? 'track' : 'tracks'}` : null,
			totalDuration
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
		return track.source || track.audioUrl || track.previewUrl || track.streamUrl || null;
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
</script>

<svelte:head>
	<title>{release.title} by {artist?.name ?? 'Silk Wave'} - Silk Wave</title>
	<meta
		name="description"
		content={release.description || `${release.title} on Silk Wave`}
	/>
</svelte:head>

<article class="relative isolate -mx-4 -mt-6 min-h-dvh overflow-hidden bg-background text-foreground md:-mx-10">
	<div class="absolute inset-0 -z-20 bg-background"></div>
	{#if coverArt}
		<img
			src={coverArt}
			alt=""
			class="absolute inset-0 -z-20 h-full w-full scale-110 object-cover opacity-70 blur-3xl saturate-125"
			aria-hidden="true"
		/>
	{/if}
	<div class="absolute inset-0 -z-10 bg-[radial-gradient(circle_at_20%_20%,color-mix(in_oklch,var(--primary)_28%,transparent),transparent_34%),linear-gradient(120deg,color-mix(in_oklch,var(--background)_72%,black),color-mix(in_oklch,var(--background)_88%,black)_42%,color-mix(in_oklch,var(--foreground)_9%,var(--background)))] dark:bg-[radial-gradient(circle_at_24%_18%,color-mix(in_oklch,var(--primary)_23%,transparent),transparent_34%),linear-gradient(125deg,rgba(0,0,0,0.76),rgba(0,0,0,0.9)_46%,color-mix(in_oklch,var(--background)_82%,black))]"></div>

	<div class="mx-auto flex min-h-dvh w-full max-w-7xl flex-col px-4 pb-32 pt-6 sm:px-6 md:px-10 lg:px-12">
		<a
			href={artistHref}
			class="mb-8 inline-flex w-fit items-center gap-2 rounded-md border border-foreground/10 bg-background/55 px-3 py-2 text-sm font-medium text-muted-foreground backdrop-blur-xl transition hover:border-primary/35 hover:text-foreground dark:border-white/10 dark:bg-black/20"
		>
			<Icon icon="chevron-left" variant="line" class="h-4 w-4 fill-current" />
			<span>{artist?.name ?? 'Back to Artist'}</span>
		</a>

		<section class="grid flex-1 gap-8 lg:grid-cols-[minmax(260px,440px)_minmax(0,1fr)] lg:items-center lg:gap-14">
			<div class="mx-auto w-full max-w-[440px] lg:mx-0">
				<div class="relative aspect-square overflow-hidden rounded-lg border border-foreground/10 bg-card/55 shadow-2xl shadow-black/35 backdrop-blur-xl dark:border-white/10">
					{#if coverArt}
						<img
							src={coverArt}
							alt="{release.title} cover art"
							class="h-full w-full object-cover"
						/>
					{:else}
						<div class="flex h-full w-full items-center justify-center bg-muted/70 text-muted-foreground">
							<Icon icon="music-note-2" class="h-20 w-20 fill-current" />
						</div>
					{/if}
				</div>
			</div>

			<div class="min-w-0 text-foreground">
				<div class="flex flex-wrap items-center gap-2 text-sm font-medium text-muted-foreground">
					{#each primaryMeta as item (`meta-${item}`)}
						<span class="rounded-full border border-foreground/10 bg-background/50 px-3 py-1 backdrop-blur-xl dark:border-white/10 dark:bg-white/[0.06]">
							{item}
						</span>
					{/each}
				</div>

				<h1 class="mt-5 max-w-4xl text-balance break-words text-5xl font-semibold leading-[0.95] tracking-tight text-foreground sm:text-6xl lg:text-7xl">
					{release.title}
				</h1>

				{#if artist}
					<a
						href={artistHref}
						class="mt-5 inline-flex max-w-full items-center gap-3 text-lg font-medium text-muted-foreground transition hover:text-primary"
					>
						<span class="flex h-11 w-11 shrink-0 items-center justify-center overflow-hidden rounded-md border border-foreground/10 bg-background/60 text-sm text-foreground backdrop-blur-xl dark:border-white/10">
							{artist.name?.slice(0, 1) ?? 'A'}
						</span>
						<span class="min-w-0 truncate">{artist.name}</span>
					</a>
				{/if}

				<div class="mt-8 grid gap-4 rounded-lg border border-foreground/10 bg-background/60 p-4 backdrop-blur-xl dark:border-white/10 dark:bg-black/25 sm:grid-cols-[1fr_auto] sm:items-center">
					<div>
						<p class="text-sm text-muted-foreground">Digital release</p>
						<p class="mt-1 text-2xl font-semibold">{price ?? 'Price not set'}</p>
					</div>

					<div class="flex flex-col gap-2 sm:flex-row">
						<Button
							variant="primary"
							class="h-12 rounded-md px-5"
							disabled
							aria-label={`${buyLabel} is not available yet`}
							title="Purchase flow is not connected yet"
						>
							<Icon icon="shopping-cart" class="h-5 w-5 fill-current" />
							{buyLabel}
						</Button>

						<Button
							variant="secondary"
							class="h-12 rounded-md px-5"
							disabled={!hasPlayableTracks}
							aria-label={hasPlayableTracks ? 'Use the player controls below' : 'Playback is not available yet'}
							title={hasPlayableTracks ? 'Use the player controls below' : 'No playable track source is available yet'}
						>
							<Icon icon="play" variant="filled" class="h-4 w-4 fill-current" />
							Play
						</Button>
					</div>
				</div>

				<div class="mt-4 flex flex-wrap gap-2">
					<Button variant="secondary" class="h-10 rounded-md px-3" disabled>
						<Icon icon="plus" variant="line" class="h-4 w-4 fill-current" />
						Add
					</Button>
					<Button variant="secondary" class="h-10 rounded-md px-3" disabled>
						Share
					</Button>
				</div>

				{#if release.description}
					<p class="mt-8 max-w-2xl whitespace-pre-line text-base leading-7 text-muted-foreground">
						{release.description}
					</p>
				{/if}

				{#if release.metadata?.genres?.length}
					<div class="mt-6 flex flex-wrap gap-2">
						{#each release.metadata.genres as genre}
							<span class="rounded-full border border-foreground/10 bg-background/45 px-3 py-1 text-sm text-muted-foreground backdrop-blur-xl dark:border-white/10">
								{genre}
							</span>
						{/each}
					</div>
				{/if}
			</div>
		</section>

		<section id="tracks" class="mt-14 grid gap-6 lg:grid-cols-[minmax(0,1fr)_minmax(260px,340px)]">
			<div class="rounded-lg border border-foreground/10 bg-background/72 p-4 backdrop-blur-xl dark:border-white/10 dark:bg-black/30 sm:p-6">
				<div class="flex flex-col gap-2 border-b border-border/70 pb-4 sm:flex-row sm:items-end sm:justify-between">
					<div>
						<p class="text-sm font-medium uppercase tracking-[0.16em] text-muted-foreground">
							Track list
						</p>
						<h2 class="mt-1 text-2xl font-semibold">Listen through the Release</h2>
					</div>
					{#if totalDuration}
						<p class="text-sm text-muted-foreground">
							{formatTotalDuration(totalDuration)}
						</p>
					{/if}
				</div>

				{#if hasTracks}
					<ol class="divide-y divide-border/70">
						{#if tracks.length}
							{#each tracks as track, index (track.id || track._key || `${index}-${track.title}`)}
								{@const source = resolveTrackSource(track)}
								{@const duration = formatDuration(track.duration, track.durationDisplay)}
								<li class="grid grid-cols-[2rem_minmax(0,1fr)_auto] items-center gap-3 py-4 sm:grid-cols-[2rem_2.75rem_minmax(0,1fr)_auto_auto]">
									<span class="text-sm tabular-nums text-muted-foreground">
										{track.order || index + 1}
									</span>
									<button
										type="button"
										class="hidden h-10 w-10 items-center justify-center rounded-md border border-border bg-background/55 text-foreground transition hover:border-primary/50 disabled:cursor-not-allowed disabled:opacity-40 sm:inline-flex"
										disabled={!source}
										aria-label={source ? `Play ${track.title}` : `Playback unavailable for ${track.title}`}
										title={source ? `Play ${track.title}` : 'No playable source available'}
									>
										<Icon icon="play" variant="filled" class="h-4 w-4 fill-current" />
									</button>
									<div class="min-w-0">
										<p class="truncate font-medium">{track.title}</p>
										{#if track.isExplicit}
											<p class="mt-1 text-xs text-muted-foreground">Explicit</p>
										{/if}
									</div>
									{#if duration}
										<span class="text-sm tabular-nums text-muted-foreground">{duration}</span>
									{/if}
									<button
										type="button"
										class="hidden rounded-md px-2 py-1 text-sm text-muted-foreground transition hover:bg-foreground/5 hover:text-foreground disabled:cursor-not-allowed disabled:opacity-40 sm:inline-flex"
										disabled
										aria-label={`Add ${track.title}`}
										title="Track purchase flow is not connected yet"
									>
										Add
									</button>
								</li>
							{/each}
						{:else}
							{#each trackTitles as title, index (`${index}-${title}`)}
								<li class="grid grid-cols-[2rem_minmax(0,1fr)] items-center gap-3 py-4">
									<span class="text-sm tabular-nums text-muted-foreground">{index + 1}</span>
									<span class="min-w-0 truncate font-medium">{title}</span>
								</li>
							{/each}
						{/if}
					</ol>
				{:else}
					<div class="py-10 text-muted-foreground">
						No tracks are listed for this Release yet.
					</div>
				{/if}
			</div>

			<aside class="h-fit rounded-lg border border-foreground/10 bg-background/64 p-5 backdrop-blur-xl dark:border-white/10 dark:bg-black/25">
				<p class="text-sm font-medium uppercase tracking-[0.16em] text-muted-foreground">
					Release details
				</p>
				<dl class="mt-4 space-y-4 text-sm">
					<div class="flex justify-between gap-4">
						<dt class="text-muted-foreground">Artist</dt>
						<dd class="min-w-0 truncate font-medium">{artist?.name ?? 'Unknown Artist'}</dd>
					</div>
					<div class="flex justify-between gap-4">
						<dt class="text-muted-foreground">Type</dt>
						<dd class="font-medium">{formatReleaseType(release.releaseType)}</dd>
					</div>
					<div class="flex justify-between gap-4">
						<dt class="text-muted-foreground">Price</dt>
						<dd class="font-medium">{price ?? 'Not set'}</dd>
					</div>
					{#if release.downloadEnabled !== undefined}
						<div class="flex justify-between gap-4">
							<dt class="text-muted-foreground">Downloads</dt>
							<dd class="font-medium">{release.downloadEnabled ? 'Enabled' : 'Disabled'}</dd>
						</div>
					{/if}
					{#if release.streamingEnabled !== undefined}
						<div class="flex justify-between gap-4">
							<dt class="text-muted-foreground">Streaming</dt>
							<dd class="font-medium">{release.streamingEnabled ? 'Enabled' : 'Disabled'}</dd>
						</div>
					{/if}
				</dl>
			</aside>
		</section>
	</div>

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
