<script lang="ts">
	import ImageDropzone from '../inputs/ImageDropzone.svelte';
	import { getWizardContext } from './wizard.svelte';

	const wizard = getWizardContext();

	let preview = $state(wizard.coverArtPreview);

	$effect(() => {
		wizard.coverArtPreview = preview;
	});

	function handleChange(file: File) {
		wizard.setCoverArt(file);
	}
</script>

<div class="flex flex-col gap-3">
	<ImageDropzone
		bind:preview
		onchange={handleChange}
		placeholder="Drop image here"
		hint="or click to browse"
	/>
	<p class="text-center text-xs text-neutral-500">Recommended: 3000x3000px, JPG or PNG</p>
</div>
