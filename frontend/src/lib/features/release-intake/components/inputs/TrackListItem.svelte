<script lang="ts">
	import Icon from '$lib/components/atoms/Icon.svelte';
	import { Input } from '$lib/components/ui/input/index.js';
	import type { WizardTrack } from '$lib/features/release-intake/types';

	interface Props {
		track: WizardTrack;
		index: number;
		isFirst: boolean;
		isLast: boolean;
		onTitleChange?: (title: string) => void;
		onFileChange?: (file: File) => void;
		onMoveUp?: () => void;
		onMoveDown?: () => void;
		onRemove?: () => void;
		editingTrack?: boolean;
	}

	let { track, index, onTitleChange, onFileChange, onRemove, editingTrack }: Props = $props();

	const handleDrop = (e: DragEvent) => {
		e.preventDefault();
		const files = e.dataTransfer?.files;
		if (files && files.length > 0) {
			const file = files[0];
			if (file.type.startsWith('audio/')) {
				onFileChange?.(file);
			}
		}
	};

	const openFilePicker = () => {
		const input = document.createElement('input');
		input.type = 'file';
		input.accept = 'audio/*';
		input.onchange = (e) => {
			const files = (e.target as HTMLInputElement).files;
			if (files && files.length > 0) {
				onFileChange?.(files[0]);
			}
		};
		input.click();
	};
</script>

<div
	class="group flex items-center gap-4 rounded-lg border border-border
	 bg-background p-4 transition-all duration-200 hover:border-muted-background hover:bg-foreground/[0.02]"
>
	<div class="flex w-full items-center justify-between gap-4">
		<div
			class="flex h-10 w-10 flex-shrink-0 items-center justify-center rounded-md
		bg-muted-background text-sm font-medium text-muted-foreground"
		>
			{index + 1}
		</div>

		<div class="flex w-full items-center gap-2">
			{#if editingTrack}
				<Input
					value={track.title}
					oninput={(e) => onTitleChange?.((e.target as HTMLInputElement).value)}
					onkeydown={(e) => e.key === 'Enter' && (editingTrack = false)}
					placeholder="Track title"
					class="h-10"
				/>
				<button
					onclick={() => (editingTrack = false)}
					class="rounded-md fill-foreground p-2 text-muted-foreground transition-colors hover:bg-muted-background hover:fill-foreground hover:text-destructive-foreground"
				>
					<Icon icon="check" class="h-5 w-5 fill-foreground" />
				</button>
			{:else}
				<span class="text-sm font-medium text-foreground">{track.title}</span>
				<button
					onclick={() => (editingTrack = true)}
					class="rounded-md fill-foreground p-2 text-muted-foreground transition-colors hover:bg-muted-background hover:fill-foreground hover:text-destructive-foreground"
				>
					<Icon icon="edit" class="h-5 w-5 fill-foreground" />
				</button>
			{/if}
		</div>

		<!-- Actions -->
		<div
			class="ml-auto flex flex-shrink-0 items-center gap-1 transition-opacity group-hover:opacity-100"
		>
			<button
				class="
			flex flex-col rounded-md fill-foreground p-2 text-muted-foreground transition-colors
			hover:bg-muted-background hover:fill-foreground hover:text-destructive-foreground"
			>
				<Icon icon="chevron-up" class="h-5 w-5" />
				<Icon icon="chevron-down" class="h-5 w-5" />
			</button>
			<button
				type="button"
				onclick={onRemove}
				class="rounded-md fill-foreground p-2 text-muted-foreground transition-colors hover:bg-destructive-foreground/10 hover:fill-destructive-foreground hover:text-destructive-foreground"
			>
				<Icon icon="trash" class="h-4 w-4" />
			</button>
		</div>
	</div>
</div>
