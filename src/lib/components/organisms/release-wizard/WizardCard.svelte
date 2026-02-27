<script lang="ts">
	import type { Snippet } from 'svelte';
	import { getWizardContext } from './wizard.svelte';

	type Props = {
		stepIndex: number;
		children: Snippet;
		footer: Snippet;
	};

	let { stepIndex, children, footer }: Props = $props();
	const wizard = getWizardContext();

	// Check if this card is the current step
	const isCurrentStep = $derived(stepIndex === wizard.currentStep);

	const step = $derived(wizard.steps[stepIndex]);
</script>

<!-- Only show current step with fade transition -->
{#if isCurrentStep}
	<div
		class="animate-in fade-in slide-in-from-bottom-2 duration-300"
	>
		<div class="overflow-hidden rounded-xl border bg-background">
			<!-- Card content -->
			<div class="flex flex-col gap-4 p-4 md:gap-6 md:p-6">
				<div class="flex flex-col gap-1">
					<h2 class="text-base font-medium text-neutral-900 dark:text-neutral-50 md:text-lg">
						{step.title}
					</h2>
					<p class="text-xs text-neutral-500 dark:text-neutral-400 md:text-sm">
						{step.description}
					</p>
				</div>
				{@render children()}
			</div>

			<!-- Card footer -->
			<div
				class="flex items-center justify-between gap-2 border-t border-neutral-100 p-3 dark:border-neutral-800 md:p-4"
			>
				{@render footer()}
			</div>
		</div>
	</div>
{/if}
