<script lang="ts">
	import { Label } from '$lib/components/ui/label/index.js';
	import AudioDropzone from '../inputs/AudioDropzone.svelte';
	import TrackList from '../inputs/TrackList.svelte';
	import { getWizardContext } from './wizard.svelte';
	import Button from '$lib/components/ui/button/button.svelte';
	import Icon from '$lib/components/ui/icon/Icon.svelte';

	const wizard = getWizardContext();
</script>

<div class="flex flex-col gap-4 md:gap-6">
	<!-- Summary of previous steps -->
	<div class="flex items-center gap-3 rounded-lg bg-muted-background p-3 md:gap-4 md:p-4">
		{#if wizard.coverArtPreview}
			<img
				src={wizard.coverArtPreview}
				alt="Cover art"
				class="h-12 w-12 rounded-md object-cover md:h-16 md:w-16"
			/>
		{/if}
		<div class="min-w-0 flex-1">
			<p class="truncate text-sm font-medium text-neutral-900 dark:text-neutral-50 md:text-base">
				{wizard.releaseTitle || 'Untitled Release'}
			</p>
			<p class="truncate text-xs text-neutral-500 md:text-sm">
				{wizard.genres.length > 0 ? wizard.genres.join(', ') : 'No genres selected'}
			</p>
		</div>
		<Button variant="ghost" onclick={() => wizard.goToStep(0)} class="shrink-0 p-2">
			<Icon icon="edit" class="h-5 w-5 fill-foreground md:h-6 md:w-6" />
		</Button>
	</div>

	<!-- Track List -->
	<div class="flex flex-col gap-3 md:gap-4">
		<div class="flex items-center justify-between">
			<Label class="text-sm font-medium">Tracks ({wizard.tracks.length})</Label>
		</div>

		<TrackList
			tracks={wizard.tracks}
			onTitleChange={wizard.updateTrackTitle}
			onFileChange={wizard.updateTrackFile}
			onMove={wizard.moveTrack}
			onRemove={wizard.removeTrack}
		/>

		<AudioDropzone onfiles={wizard.addBulkTracks} />
	</div>
</div>
