<script lang="ts" module>
	export type MediaPlayerTrack = {
		id: string;
		title: string;
		artistName?: string | null;
		source?: string | null;
		audioUrl?: string | null;
		previewUrl?: string | null;
		artworkUrl?: string | null;
		duration?: number | null;
	};

	export type MediaPlayerPlacement = 'floating' | 'inline';
</script>

<script lang="ts">
	import { tick } from 'svelte';
	import Icon from '$lib/components/atoms/Icon.svelte';
	import { cn } from '$lib/utils/utils.js';

	interface Props {
		tracks?: MediaPlayerTrack[];
		initialTrackId?: string | null;
		releaseTitle?: string;
		artistName?: string;
		artworkUrl?: string | null;
		placement?: MediaPlayerPlacement;
		label?: string;
		class?: string;
	}

	let {
		tracks = [],
		initialTrackId = null,
		releaseTitle = 'Untitled release',
		artistName = 'Unknown artist',
		artworkUrl = null,
		placement = 'floating',
		label = 'Media controls',
		class: className = ''
	}: Props = $props();

	let audioEl: HTMLAudioElement | null = $state(null);
	let currentTrackId = $state<string | null>(null);
	let currentTime = $state(0);
	let measuredDuration = $state(0);
	let isPlaying = $state(false);
	let volume = $state(0.85);
	let isMuted = $state(false);

	let currentIndex = $derived(tracks.findIndex((track) => track.id === currentTrackId));
	let currentTrack = $derived(tracks[currentIndex >= 0 ? currentIndex : 0] ?? null);
	let source = $derived(resolveTrackSource(currentTrack));
	let resolvedArtwork = $derived(currentTrack?.artworkUrl || artworkUrl || null);
	let duration = $derived(
		normaliseDuration(measuredDuration) || normaliseDuration(currentTrack?.duration) || 0
	);
	let progress = $derived(duration > 0 ? Math.min((currentTime / duration) * 100, 100) : 0);
	let volumeProgress = $derived(isMuted ? 0 : volume * 100);
	let canPlay = $derived(Boolean(source));
	let canSkip = $derived(tracks.length > 1);
	let unavailableMessage = $derived(
		tracks.length === 0
			? 'No tracks are available yet.'
			: 'Audio preview unavailable for this track.'
	);
	let displayTitle = $derived(currentTrack?.title || releaseTitle);
	let displayArtist = $derived(currentTrack?.artistName || artistName || 'Unknown artist');

	$effect(() => {
		const fallbackId = tracks[0]?.id ?? null;
		const requestedId = initialTrackId ?? fallbackId;
		const hasCurrentTrack = tracks.some((track) => track.id === currentTrackId);

		if (!hasCurrentTrack) {
			currentTrackId = requestedId;
		}
	});

	$effect(() => {
		if (!audioEl) return;

		audioEl.volume = isMuted ? 0 : volume;
	});

	$effect(() => {
		currentTrackId;
		currentTime = 0;
		measuredDuration = 0;
	});

	function resolveTrackSource(track: MediaPlayerTrack | null) {
		return track?.source || track?.audioUrl || track?.previewUrl || '';
	}

	function normaliseDuration(value: number | null | undefined) {
		if (!value || !Number.isFinite(value)) return 0;

		return value;
	}

	function formatTime(seconds: number) {
		const safeSeconds = Math.max(0, Math.floor(seconds || 0));
		const minutes = Math.floor(safeSeconds / 60);
		const remainingSeconds = safeSeconds % 60;

		return `${minutes}:${remainingSeconds.toString().padStart(2, '0')}`;
	}

	async function togglePlayback() {
		if (!audioEl || !canPlay) return;

		if (isPlaying) {
			audioEl.pause();
			isPlaying = false;
			return;
		}

		await playCurrent();
	}

	async function playCurrent() {
		if (!audioEl || !canPlay) return;

		try {
			await audioEl.play();
			isPlaying = true;
		} catch {
			isPlaying = false;
		}
	}

	async function selectTrack(index: number, shouldPlay = isPlaying) {
		const nextTrack = tracks[index];
		if (!nextTrack) return;

		currentTrackId = nextTrack.id;
		await tick();

		if (shouldPlay) {
			await playCurrent();
		}
	}

	function previousTrack() {
		if (!canSkip) return;

		const nextIndex = currentIndex <= 0 ? tracks.length - 1 : currentIndex - 1;
		void selectTrack(nextIndex);
	}

	function nextTrack() {
		if (!canSkip) return;

		const nextIndex =
			currentIndex >= tracks.length - 1 || currentIndex < 0 ? 0 : currentIndex + 1;
		void selectTrack(nextIndex);
	}

	function handleLoadedMetadata() {
		if (!audioEl) return;

		measuredDuration = normaliseDuration(audioEl.duration);
	}

	function handleTimeUpdate() {
		if (!audioEl) return;

		currentTime = audioEl.currentTime;
	}

	function handleSeek(event: Event) {
		const nextTime = Number((event.currentTarget as HTMLInputElement).value);
		currentTime = nextTime;

		if (audioEl) {
			audioEl.currentTime = nextTime;
		}
	}

	function handleVolume(event: Event) {
		const nextVolume = Number((event.currentTarget as HTMLInputElement).value);
		volume = Math.max(0, Math.min(nextVolume, 1));
		isMuted = volume === 0;
	}

	function toggleMute() {
		isMuted = !isMuted;
	}

	function handleEnded() {
		if (currentIndex >= 0 && currentIndex < tracks.length - 1) {
			void selectTrack(currentIndex + 1, true);
			return;
		}

		isPlaying = false;
		currentTime = 0;
	}
</script>

<div
	class={cn(
		placement === 'floating'
			? 'pointer-events-none fixed inset-x-0 bottom-0 z-50 flex justify-center px-4 pb-4 sm:pb-8'
			: 'w-full',
		className
	)}
>
	<section
		aria-label={label}
		class={cn(
			'pointer-events-auto w-full rounded-lg border border-border/60 bg-background/76 text-foreground shadow-2xl shadow-foreground/10 backdrop-blur-xl dark:border-white/10 dark:bg-card/70 dark:shadow-black/45',
			placement === 'floating' ? 'max-w-240 p-3 sm:p-4' : 'p-4'
		)}
		style={`--media-progress: ${progress}%; --volume-progress: ${volumeProgress}%;`}
	>
		<audio
			bind:this={audioEl}
			src={source || undefined}
			preload="metadata"
			onloadedmetadata={handleLoadedMetadata}
			ontimeupdate={handleTimeUpdate}
			onended={handleEnded}
			onpause={() => (isPlaying = false)}
			onplay={() => (isPlaying = true)}
		></audio>

		<div
			class="grid gap-3 md:grid-cols-[minmax(0,260px)_minmax(180px,1fr)_auto] md:items-center"
		>
			<div class="flex min-w-0 items-center gap-3">
				<div
					class="relative h-12 w-12 shrink-0 overflow-hidden rounded-md border border-border/80 bg-muted shadow-sm sm:h-14 sm:w-14"
					aria-hidden="true"
				>
					{#if resolvedArtwork}
						<img src={resolvedArtwork} alt="" class="h-full w-full object-cover" />
					{:else}
						<div class="flex h-full w-full items-center justify-center">
							<Icon icon="music-note-2" class="h-5 w-5 fill-muted-foreground" />
						</div>
					{/if}
				</div>

				<div class="min-w-0">
					<p class="truncate text-sm font-semibold leading-5 text-foreground">
						{displayTitle}
					</p>
					<p class="truncate text-xs leading-5 text-muted-foreground">
						{displayArtist}
					</p>
				</div>
			</div>

			<div class="hidden items-center gap-3 md:flex">
				<span class="w-10 text-right text-xs tabular-nums text-muted-foreground">
					{formatTime(currentTime)}
				</span>
				<input
					type="range"
					min="0"
					max={duration || 0}
					step="0.1"
					value={Math.min(currentTime, duration || 0)}
					disabled={!canPlay || duration <= 0}
					aria-label="Track progress"
					class="media-player-range"
					oninput={handleSeek}
				/>
				<span class="w-10 text-xs tabular-nums text-muted-foreground">
					{formatTime(duration)}
				</span>
			</div>

			<div class="flex items-center justify-between gap-2 md:justify-end">
				<div class="flex items-center gap-1">
					<button
						type="button"
						class="media-player-button"
						aria-label="Previous track"
						disabled={!canSkip}
						onclick={previousTrack}
					>
						<Icon icon="chevron-left" class="h-5 w-5 fill-current" />
					</button>
					<button
						type="button"
						class="media-player-button media-player-button-primary"
						aria-label={canPlay ? (isPlaying ? 'Pause' : 'Play') : unavailableMessage}
						aria-pressed={isPlaying}
						title={!canPlay ? unavailableMessage : undefined}
						disabled={!canPlay}
						onclick={togglePlayback}
					>
						<Icon
							icon={isPlaying ? 'pause' : 'play'}
							variant="filled"
							class="h-5 w-5 fill-current"
						/>
					</button>
					<button
						type="button"
						class="media-player-button"
						aria-label="Next track"
						disabled={!canSkip}
						onclick={nextTrack}
					>
						<Icon icon="chevron-right" class="h-5 w-5 fill-current" />
					</button>
				</div>

				<div class="hidden items-center gap-2 sm:flex">
					<button
						type="button"
						class="media-player-button"
						aria-label={isMuted ? 'Unmute' : 'Mute'}
						aria-pressed={isMuted}
						disabled={!canPlay}
						onclick={toggleMute}
					>
						<Icon
							icon={isMuted ? 'volume-cross' : 'volume-high'}
							class="h-5 w-5 fill-current"
						/>
					</button>
					<input
						type="range"
						min="0"
						max="1"
						step="0.01"
						value={isMuted ? 0 : volume}
						aria-label="Volume"
						class="media-player-range media-player-volume"
						disabled={!canPlay}
						oninput={handleVolume}
					/>
				</div>
			</div>
		</div>

		<div class="mt-3 flex items-center gap-3 md:hidden">
			<span class="w-9 text-right text-xs tabular-nums text-muted-foreground">
				{formatTime(currentTime)}
			</span>
			<input
				type="range"
				min="0"
				max={duration || 0}
				step="0.1"
				value={Math.min(currentTime, duration || 0)}
				disabled={!canPlay || duration <= 0}
				aria-label="Track progress"
				class="media-player-range"
				oninput={handleSeek}
			/>
			<span class="w-9 text-xs tabular-nums text-muted-foreground">
				{formatTime(duration)}
			</span>
		</div>
	</section>
</div>

<style>
	.media-player-button {
		display: inline-flex;
		height: 2.5rem;
		width: 2.5rem;
		align-items: center;
		justify-content: center;
		border-radius: calc(var(--radius) - 2px);
		color: var(--foreground);
		transition:
			background-color 160ms ease,
			color 160ms ease,
			opacity 160ms ease,
			transform 160ms ease;
	}

	.media-player-button:hover:not(:disabled) {
		background: color-mix(in oklch, var(--foreground) 8%, transparent);
		transform: translateY(-1px);
	}

	.media-player-button:focus-visible {
		outline: 2px solid var(--ring);
		outline-offset: 2px;
	}

	.media-player-button:disabled {
		cursor: not-allowed;
		opacity: 0.38;
	}

	.media-player-button-primary {
		background: color-mix(in oklch, var(--primary) 88%, var(--foreground));
		color: var(--primary-foreground);
		box-shadow: inset 0 1px 1px color-mix(in oklch, var(--background) 24%, transparent);
	}

	.media-player-button-primary:hover:not(:disabled) {
		background: color-mix(in oklch, var(--primary) 78%, var(--foreground));
	}

	.media-player-range {
		--range-fill: var(--media-progress);
		--range-rest: color-mix(in oklch, var(--foreground) 14%, transparent);

		height: 0.375rem;
		width: 100%;
		cursor: pointer;
		appearance: none;
		border-radius: calc(var(--radius) - 4px);
		background: linear-gradient(
			to right,
			var(--primary) 0 var(--range-fill),
			var(--range-rest) var(--range-fill) 100%
		);
	}

	.media-player-volume {
		--range-fill: var(--volume-progress);

		width: 5.5rem;
	}

	.media-player-range:disabled {
		cursor: not-allowed;
		opacity: 0.45;
	}

	.media-player-range:focus-visible {
		outline: 2px solid var(--ring);
		outline-offset: 3px;
	}

	.media-player-range::-webkit-slider-thumb {
		height: 0.875rem;
		width: 0.5rem;
		appearance: none;
		border-radius: calc(var(--radius) - 4px);
		background: var(--foreground);
		box-shadow: 0 1px 6px color-mix(in oklch, var(--foreground) 22%, transparent);
	}

	.media-player-range::-moz-range-thumb {
		height: 0.875rem;
		width: 0.5rem;
		border: 0;
		border-radius: calc(var(--radius) - 4px);
		background: var(--foreground);
		box-shadow: 0 1px 6px color-mix(in oklch, var(--foreground) 22%, transparent);
	}
</style>
