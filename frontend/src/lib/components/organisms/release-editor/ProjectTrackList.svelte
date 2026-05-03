<script lang="ts">
	import Icon from '$lib/components/atoms/Icon.svelte';
	import type { WizardTrack } from '$lib/types/WizardTrack';
	import ProjectTrackItem from './ProjectTrackItem.svelte';

	interface Props {
		tracks: WizardTrack[];
		onfiles?: (files: File[]) => void;
		onTitleChange?: (id: string, title: string) => void;
		onRemove?: (id: string) => void;
		onMove?: (index: number, direction: 'up' | 'down') => void;
	}

	let { tracks, onfiles, onTitleChange, onRemove, onMove }: Props = $props();

	const openFilePicker = () => {
		const input = document.createElement('input');
		input.type = 'file';
		input.accept = 'audio/*';
		input.multiple = true;
		input.onchange = (e) => {
			const files = (e.target as HTMLInputElement).files;
			if (files) {
				const audioFiles = Array.from(files).filter((file) =>
					file.type.startsWith('audio/')
				);
				if (audioFiles.length > 0) {
					onfiles?.(audioFiles);
				}
			}
		};
		input.click();
	};
</script>

<div class="flex w-full flex-col gap-4">
	<!-- Add tracks button -->
	<button
		onclick={openFilePicker}
		class="flex w-full cursor-pointer items-center justify-center gap-2 rounded-xl bg-foreground/5 px-4 py-3 text-sm font-medium text-foreground transition-opacity duration-150 select-none hover:opacity-80"
	>
		<Icon icon="plus" variant="line" class="h-4 w-4 fill-foreground" />
		Add tracks
	</button>

	<!-- Track list -->
	{#if tracks.length > 0}
		<ul class="flex flex-col">
			{#each tracks as track, index (track.id)}
				<ProjectTrackItem
					{track}
					{index}
					total={tracks.length}
					onTitleChange={(title) => onTitleChange?.(track.id, title)}
					onRemove={() => onRemove?.(track.id)}
					onMove={onMove ? (direction) => onMove?.(index, direction) : undefined}
				/>
			{/each}
		</ul>
	{:else}{/if}
</div>
