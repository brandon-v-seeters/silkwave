<script lang="ts">
	import { Button } from '$lib/components/ui/button/index.js';
	import Icon from '$lib/components/ui/icon/Icon.svelte';
	import { getWizardContext } from './wizard.svelte';

	const wizard = getWizardContext();
</script>

<Button variant="ghost" onclick={wizard.prevStep} class="gap-1 px-2 md:gap-2 md:px-4">
	<Icon icon="chevron-left" variant="line" class="h-4 w-4 fill-foreground md:h-5 md:w-5" />
	<span class="hidden sm:inline">Back</span>
</Button>

<div class="flex gap-2 md:gap-3">
	<Button
		variant="secondary"
		onclick={wizard.saveAsDraft}
		disabled={wizard.isUploading}
		class="gap-1 px-3 text-sm md:gap-2 md:px-4 md:text-base"
	>
		{#if wizard.isUploading}
			<Icon icon="loader" variant="line" class="h-4 w-4 animate-spin fill-foreground" />
			<span class="hidden sm:inline">
				{wizard.uploadStatus === 'creating_draft'
					? 'Creating...'
					: wizard.uploadStatus === 'uploading'
						? 'Uploading...'
						: 'Saving...'}
			</span>
		{:else}
			<Icon icon="diskette" variant="line" class="h-4 w-4 fill-foreground" />
			<span class="hidden sm:inline">Save Draft</span>
		{/if}
	</Button>

	<Button
		onclick={wizard.publishRelease}
		disabled={!wizard.isFormValid || wizard.isUploading}
		class="gap-1 px-3 text-sm md:gap-2 md:px-4 md:text-base"
		variant="primary"
	>
		<Icon icon="rocket" variant="filled" class="h-4 w-4 fill-background" />
		<span>Publish</span>
	</Button>
</div>
