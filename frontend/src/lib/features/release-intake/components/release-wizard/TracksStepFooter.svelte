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

<Button onclick={wizard.nextStep} disabled={!wizard.canProceed} class="gap-1 px-3 md:gap-2 md:px-4">
	<span>Next</span>
	<Icon icon="chevron-right" variant="line" class="h-4 w-4 fill-background" />
</Button>
