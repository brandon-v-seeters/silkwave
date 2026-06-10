<script lang="ts">
	import Icon from '$lib/components/atoms/Icon.svelte';
	import { generateId } from '$lib/utils/utils';

	interface Props {
		preview?: string | null;
		accept?: string;
		aspectRatio?: 'square' | 'video' | 'banner';
		placeholder?: string;
		hint?: string;
		onchange?: (file: File) => void;
	}

	let {
		preview = $bindable(null),
		accept = 'image/*',
		aspectRatio = 'square',
		placeholder = 'Drop image here',
		hint = 'or click to browse',
		onchange
	}: Props = $props();

	const inputId = `image-dropzone-${generateId()}`;

	const aspectClasses = {
		square: 'aspect-square',
		video: 'aspect-video',
		banner: 'aspect-[3/1]'
	};

	const processFile = (file: File) => {
		const reader = new FileReader();
		reader.onload = (e) => {
			preview = e.target?.result as string;
		};
		reader.readAsDataURL(file);
		onchange?.(file);
	};

	const handleDrop = (e: DragEvent) => {
		e.preventDefault();
		const files = e.dataTransfer?.files;
		if (files && files.length > 0) {
			const file = files[0];
			if (file.type.startsWith('image/')) {
				processFile(file);
			}
		}
	};

	const handleSelect = (e: Event) => {
		const input = e.target as HTMLInputElement;
		if (input.files && input.files.length > 0) {
			const file = input.files[0];
			if (file.type.startsWith('image/')) {
				processFile(file);
			}
		}
	};

	const openFilePicker = () => {
		document.getElementById(inputId)?.click();
	};
</script>

<!-- svelte-ignore a11y_no_static_element_interactions -->
<div
	class="group relative {aspectClasses[
		aspectRatio
	]} w-full cursor-pointer overflow-hidden rounded-lg transition-all duration-300 {preview
		? ''
		: 'border-2 border-dashed border-foreground/20 bg-muted-background hover:border-primary hover:bg-primary/5'}"
	ondragover={(e) => e.preventDefault()}
	ondrop={handleDrop}
	onclick={openFilePicker}
	onkeydown={(e) => e.key === 'Enter' && openFilePicker()}
	role="button"
	tabindex="0"
>
	{#if preview}
		<img src={preview} alt="Preview" class="h-full w-full object-cover" />
		<div
			class="absolute inset-0 flex items-center justify-center bg-black/60 opacity-0 transition-opacity duration-300 group-hover:opacity-100"
		>
			<span class="text-sm text-white">Click to change</span>
		</div>
	{:else}
		<div class="flex h-full flex-col items-center justify-center gap-3 p-6 text-center">
			<div class="rounded-full bg-zinc-800 p-4 transition-colors group-hover:bg-primary/20">
				<Icon
					icon="image"
					variant="filled"
					class="h-8 w-8 fill-foreground-muted group-hover:fill-primary"
				/>
			</div>
			<div>
				<p class="text-sm font-medium text-foreground">{placeholder}</p>
				<p class="text-xs text-foreground-muted">{hint}</p>
			</div>
		</div>
	{/if}
</div>
<input id={inputId} type="file" {accept} class="hidden" onchange={handleSelect} />
