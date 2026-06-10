<script lang="ts">
	import Icon from '$lib/components/ui/icon/Icon.svelte';
	import type { IconKey } from '$lib/types/Icon';

	interface Props {
		accept?: string;
		multiple?: boolean;
		icon?: IconKey;
		placeholder?: string;
		hint?: string;
		onfiles?: (files: File[]) => void;
	}

	let {
		accept = 'audio/*',
		multiple = true,
		icon = 'music-note-2',
		placeholder = 'Drop audio files here',
		hint = 'or click to browse • MP3, WAV, FLAC, AIFF supported',
		onfiles
	}: Props = $props();

	const handleDrop = (e: DragEvent) => {
		e.preventDefault();
		const files = e.dataTransfer?.files;
		if (files) {
			const audioFiles = Array.from(files).filter((file) => file.type.startsWith('audio/'));
			if (audioFiles.length > 0) {
				onfiles?.(audioFiles);
			}
		}
	};

	const openFilePicker = () => {
		const input = document.createElement('input');
		input.type = 'file';
		input.accept = accept;
		input.multiple = multiple;
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

<!-- svelte-ignore a11y_no_static_element_interactions -->
<div
	class="group flex cursor-pointer flex-col items-center justify-center gap-3
	rounded-lg border-2 border-dashed border-foreground/20 bg-muted-background p-8 transition-all
	duration-300 hover:border-primary hover:bg-primary/5"
	ondragover={(e) => e.preventDefault()}
	ondrop={handleDrop}
	onclick={openFilePicker}
	onkeydown={(e) => e.key === 'Enter' && openFilePicker()}
	role="button"
	tabindex="0"
>
	<div class="rounded-full bg-zinc-800 p-4 transition-colors group-hover:bg-primary/20">
		<Icon {icon} class="h-8 w-8 fill-foreground-muted group-hover:fill-primary" />
	</div>
	<div class="text-center">
		<p class="text-sm font-medium text-foreground">{placeholder}</p>
		<p class="text-xs text-foreground-muted">{hint}</p>
	</div>
</div>
