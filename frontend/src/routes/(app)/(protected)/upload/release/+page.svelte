<script lang="ts">
	import { onMount } from 'svelte';
	import { resolve } from '$app/paths';
	import Icon from '$lib/components/ui/icon/Icon.svelte';
	import {
		createWizardContext,
		ReleaseCoverArt,
		ReleaseTrackList,
		type WizardTrack
	} from '$lib/features/release-intake';
	import Button from '$lib/components/ui/button/button.svelte';

	const wizard = createWizardContext();

	const { data } = $props();
	let previewAudio: HTMLAudioElement | null = null;
	let previewUrl: string | null = null;

	onMount(() => {
		previewAudio = new Audio();

		if (data.draft) {
			wizard.initWizard(data.draft);
		}

		return () => {
			previewAudio?.pause();
			if (previewUrl) {
				URL.revokeObjectURL(previewUrl);
			}
		};
	});

	let totalDuration = $derived.by(() => {
		let totalSeconds = 0;
		for (const track of wizard.tracks) {
			const [minutes = '0', seconds = '0'] = track.duration?.split(':') ?? [];
			totalSeconds += (Number.parseInt(minutes, 10) || 0) * 60;
			totalSeconds += Number.parseInt(seconds, 10) || 0;
		}
		const m = Math.floor(totalSeconds / 60);
		const s = totalSeconds % 60;
		return m > 0 ? `${m}m ${s}s` : `${s}s`;
	});

	let trackSummary = $derived(
		wizard.tracks.length > 0
			? `${wizard.tracks.length} ${wizard.tracks.length === 1 ? 'track' : 'tracks'}`
			: 'No tracks yet'
	);

	let releaseState = $derived(wizard.isFormValid ? 'Ready' : 'Draft');

	function playTrack(track: WizardTrack) {
		if (!track.file || track.file.size === 0 || !previewAudio) return;

		previewAudio.pause();

		if (previewUrl) {
			URL.revokeObjectURL(previewUrl);
		}

		previewUrl = URL.createObjectURL(track.file);
		previewAudio.src = previewUrl;
		void previewAudio.play().catch(() => {});
	}
</script>

<svelte:head>
	<title>New release - Silkwave</title>
</svelte:head>

<div class="mx-auto flex w-full max-w-6xl flex-col gap-7 pb-10">
	<div class="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
		<Button
			href={resolve('/upload')}
			variant="ghost"
			size="sm"
			class="w-fit rounded-full px-2 text-sm text-muted-foreground hover:text-foreground"
		>
			<Icon icon="chevron-left" variant="line" class="h-4 w-4 fill-current" />
			Upload
		</Button>

		<div class="flex w-full gap-2 sm:w-auto">
			<Button
				variant="secondary"
				class="min-w-0 flex-1 rounded-full px-4 sm:flex-none"
				onclick={() => wizard.saveAsDraft()}
				disabled={wizard.isUploading}
			>
				<Icon icon="file-text" variant="line" class="h-4 w-4 fill-current" />
				Save draft
			</Button>
			<Button
				variant="primary"
				class="min-w-0 flex-1 rounded-full px-5 sm:flex-none"
				onclick={() => wizard.publishRelease()}
				disabled={!wizard.isFormValid || wizard.isUploading}
			>
				<Icon icon="check-circle" variant="line" class="h-4 w-4 fill-current" />
				Publish
			</Button>
		</div>
	</div>

	<div class="grid min-w-0 gap-8 lg:grid-cols-[minmax(15rem,24rem)_minmax(0,1fr)] xl:gap-12">
		<aside class="min-w-0 space-y-4">
			<ReleaseCoverArt
				bind:preview={wizard.coverArtPreview}
				onchange={(file) => wizard.setCoverArt(file)}
			/>

			<div class="rounded-2xl border border-border/60 bg-muted/25 p-4 shadow-sm">
				<div class="grid grid-cols-2 gap-4 text-sm">
					<div class="min-w-0">
						<p class="text-xs font-medium text-muted-foreground">State</p>
						<p class="mt-1 truncate font-semibold text-foreground">{releaseState}</p>
					</div>
					<div class="min-w-0">
						<p class="text-xs font-medium text-muted-foreground">Length</p>
						<p class="mt-1 truncate font-semibold text-foreground">{totalDuration}</p>
					</div>
				</div>
			</div>
		</aside>

		<section class="min-w-0 space-y-7">
			<div class="border-b border-border/70 pb-6">
				<p class="text-xs font-semibold uppercase tracking-[0.14em] text-muted-foreground">
					New release
				</p>

				<div class="mt-4 flex items-start gap-3 sm:gap-5">
					<div class="min-w-0 flex-1">
						<input
							type="text"
							aria-label="Release title"
							placeholder="untitled project"
							bind:value={wizard.releaseTitle}
							class="w-full bg-transparent font-serif text-4xl font-light leading-none tracking-tight text-foreground placeholder:text-muted-foreground/35 focus:outline-none sm:text-5xl"
						/>

						<div
							class="mt-4 flex flex-wrap items-center gap-2 text-sm text-muted-foreground"
						>
							<span>{trackSummary}</span>
							{#if wizard.tracks.length > 0}
								<span aria-hidden="true" class="text-border">/</span>
								<span>{totalDuration}</span>
							{/if}
						</div>
					</div>

					<Button
						size="icon"
						class="mt-1 !h-11 !w-11 shrink-0 rounded-full"
						disabled={wizard.tracks.length === 0}
						aria-label="Preview release"
					>
						<Icon
							icon="play"
							variant="filled"
							class="h-4 w-4 fill-primary-foreground"
						/>
					</Button>
				</div>
			</div>

			<div class="space-y-4">
				<div class="flex items-end justify-between gap-4">
					<div>
						<h2 class="text-xl font-semibold tracking-tight text-foreground">Tracks</h2>
						<p class="mt-1 text-sm text-muted-foreground">{trackSummary}</p>
					</div>
				</div>

				<ReleaseTrackList
					tracks={wizard.tracks}
					onfiles={(files) => wizard.addBulkTracks(files)}
					onTitleChange={(id, title) => wizard.updateTrackTitle(id, title)}
					onRemove={(id) => wizard.removeTrack(id)}
					onMove={(index, direction) => wizard.moveTrack(index, direction)}
					onPlay={playTrack}
				/>
			</div>
		</section>
	</div>
</div>
