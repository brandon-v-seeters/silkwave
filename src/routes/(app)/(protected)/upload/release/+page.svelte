<script lang="ts">
	import { goto } from '$app/navigation';
	import Icon from '$lib/components/atoms/Icon.svelte';
	import {
		createWizardContext,
		WizardCard,
		WizardStepIndicator,
		WizardNav,
		ReleaseDetailsStep,
		CoverArtStep,
		TracksStep,
		TracksStepFooter,
		PricingStep,
		PricingStepFooter
	} from '$lib/components/organisms/release-wizard';
	import Button from '$lib/components/ui/button/button.svelte';

	// Initialize the wizard context - all state lives here
	const wizard = createWizardContext();

	const { data } = $props();

	if (data.draft) {
		wizard.initWizard(data.draft);
	}
</script>

<div class="flex min-h-[calc(100vh-4rem)] flex-col pb-16 md:min-h-[calc(100vh-8rem)] md:pb-0">
	<WizardStepIndicator />
	<Button
		class="absolute right-4 top-4 rounded-full"
		variant="ghost"
		onclick={() => {
			goto('/upload');
		}}
	>
		<Icon icon="cross" variant="line" class="h-4 w-4 fill-foreground md:h-5 md:w-5" />
	</Button>
	<div class="px-4 pb-4 pt-4 md:px-6 md:pb-6 md:pt-8">
		<div class="mx-auto max-w-2xl">
			<h1 class="text-xl font-medium text-neutral-900 dark:text-white md:text-2xl">
				Create Your Release
			</h1>
			<p class="mt-1 text-sm text-foreground-muted md:hidden">
				Step {wizard.currentStep + 1} of {wizard.totalSteps}
			</p>
		</div>
	</div>

	<div class="flex-1 px-4 md:px-6">
		<div class="mx-auto max-w-2xl">
			<WizardCard stepIndex={0}>
				{#snippet children()}
					<ReleaseDetailsStep />
				{/snippet}

				{#snippet footer()}
					<WizardNav showCancel nextDisabled={!wizard.isStepValid(0)} />
				{/snippet}
			</WizardCard>
			<WizardCard stepIndex={1}>
				{#snippet children()}
					<CoverArtStep />
				{/snippet}

				{#snippet footer()}
					<WizardNav showBack nextDisabled={!wizard.isStepValid(1)} />
				{/snippet}
			</WizardCard>
			<WizardCard stepIndex={2}>
				{#snippet children()}
					<TracksStep />
				{/snippet}

				{#snippet footer()}
					<TracksStepFooter />
				{/snippet}
			</WizardCard>
			<WizardCard stepIndex={3}>
				{#snippet children()}
					<PricingStep />
				{/snippet}

				{#snippet footer()}
					<PricingStepFooter />
				{/snippet}
			</WizardCard>
		</div>
	</div>
</div>
