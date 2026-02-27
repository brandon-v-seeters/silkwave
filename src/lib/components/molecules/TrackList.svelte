<script lang="ts">
	import Icon from '$lib/components/atoms/Icon.svelte';
	import TrackListItem from './TrackListItem.svelte';
	import type { WizardTrack } from '$lib/types/WizardTrack';

	interface Props {
		tracks: WizardTrack[];
		onTitleChange?: (id: string, title: string) => void;
		onFileChange?: (id: string, file: File) => void;
		onMove?: (index: number, direction: 'up' | 'down') => void;
		onRemove?: (id: string) => void;
	}

	let { tracks, onTitleChange, onFileChange, onMove, onRemove }: Props = $props();
</script>

{#if tracks.length > 0}
	<div class="flex flex-col gap-3">
		{#each tracks as track, index (track.id)}
			<TrackListItem
				{track}
				{index}
				isFirst={index === 0}
				isLast={index === tracks.length - 1}
				onTitleChange={(title) => onTitleChange?.(track.id, title)}
				onFileChange={(file) => onFileChange?.(track.id, file)}
				onMoveUp={() => onMove?.(index, 'up')}
				onMoveDown={() => onMove?.(index, 'down')}
				onRemove={() => onRemove?.(track.id)}
			/>
		{/each}
	</div>
{:else}
	<div
		class="flex flex-col items-center justify-center gap-2 rounded-lg border border-dashed border-zinc-700 bg-zinc-900/20 py-12 text-center"
	>
		<Icon icon="music-note-2" class="h-12 w-12 fill-foreground-muted/50" />
		<p class="text-sm text-foreground-muted">No tracks added yet</p>
		<p class="text-xs text-foreground-muted/70">
			Upload files above or click "Add Track" to get started
		</p>
	</div>
{/if}

