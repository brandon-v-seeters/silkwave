<script lang="ts">
	import Icon from '$lib/components/atoms/Icon.svelte';
	import { Button } from '$lib/components/ui/button/index.js';
	import { getWizardContext } from './wizard.svelte';

	type Props = {
		showBack?: boolean;
		showCancel?: boolean;
		showNext?: boolean;
		nextDisabled?: boolean;
	};

	let {
		showBack = false,
		showCancel = false,
		showNext = true,
		nextDisabled = false
	}: Props = $props();

	const wizard = getWizardContext();
</script>

{#if showCancel}
	<Button variant="ghost" href="/upload" class="gap-1 px-2 md:gap-2 md:px-4">
		<Icon icon="chevron-left" variant="line" class="h-4 w-4 fill-foreground md:h-5 md:w-5" />
		<span class="hidden sm:inline">Cancel</span>
	</Button>
{:else if showBack}
	<Button variant="ghost" onclick={wizard.prevStep} class="gap-1 px-2 md:gap-2 md:px-4">
		<Icon icon="chevron-left" variant="line" class="h-4 w-4 fill-foreground md:h-5 md:w-5" />
		<span class="hidden sm:inline">Back</span>
	</Button>
{:else}
	<div></div>
{/if}

{#if showNext}
	<Button onclick={wizard.nextStep} disabled={nextDisabled} class="gap-1 px-3 md:gap-2 md:px-4">
		<span>Next</span>
		<Icon icon="chevron-right" variant="line" class="h-4 w-4 fill-background" />
	</Button>
{/if}
