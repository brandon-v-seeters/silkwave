<script lang="ts">
	import { handleFileUpload, handleGetPresignedUrl } from '$lib/tools/filemanagement.ts';
	import { createEventDispatcher, onMount } from 'svelte';

	interface Props {
		path?: string;
	}

	let { path = 'pages' }: Props = $props();

	const dispatch = createEventDispatcher();

	let dropZone: HTMLDivElement = $state();
	let uploadedImages: string[] = $state([]);

	const handleDragOver = (event: DragEvent) => {
		event.preventDefault();
		event.dataTransfer!.dropEffect = 'copy';
	};

	const handleDrop = (event: DragEvent) => {
		event.preventDefault();
		const files = event.dataTransfer!.files;

		if (!files) return;

		handleFile(files[0]);
	};

	const openFileManager = () => {
		const input = document.createElement('input');
		input.type = 'file';
		input.accept = 'image/*';
		input.multiple = true;

		input.onchange = (e) => {
			if (!input.files) return;

			handleFile(input.files[0]);
		};

		input.click();
	};

	const handleFile = (file: File) => {
		const reader = new FileReader();
		reader.onload = async (e) => {
			uploadedImages = [...uploadedImages, e.target!.result as string];
			const url = await handleGetPresignedUrl(file.name, path);
			const res = await handleFileUpload(file, url);

			dispatch('fileUploaded', { url, path, name: `/${path}/${file.name}` });
		};
		reader.readAsDataURL(file);
	};

	onMount(() => {
		dropZone.addEventListener('dragover', handleDragOver);
		dropZone.addEventListener('drop', handleDrop);
	});
</script>

<!-- svelte-ignore a11y_no_static_element_interactions -->
<!-- svelte-ignore a11y_click_events_have_key_events -->
<div
	bind:this={dropZone}
	class="drop-zone text-foreground-muted"
	onclick={() => {
		openFileManager();
	}}
>
	Drag & Drop your files here or click to upload
</div>

<div class="uploaded-images">
	{#each uploadedImages as image, i}
		<img src={image} alt="Uploaded {i}" />
	{/each}
</div>

<style>
	.drop-zone {
		border: 2px dashed #ccc;
		border-radius: 4px;
		padding: 20px;
		text-align: center;
		cursor: pointer;
	}
	.uploaded-images {
		display: flex;
		flex-wrap: wrap;
		margin-top: 20px;
	}
	.uploaded-images img {
		max-width: 100px;
		margin: 10px;
	}
</style>
