<script lang="ts">
	import Button from '$lib/components/ui/button/button.svelte';
	import { getWizardContext } from './wizard.svelte';

	const wizard = getWizardContext();
</script>

<!-- Mobile step indicator - fixed at bottom -->
<div
	class="fixed inset-x-0 bottom-0 z-50 flex items-center justify-center gap-2 bg-background/80 py-4 backdrop-blur-sm md:hidden"
>
	{#each wizard.steps as step, i}
		{@const state = wizard.getStepState(i)}
		<button
			class="flex h-2 w-2 items-center justify-center rounded-full transition-all
				{state === 'completed'
				? 'bg-emerald-500'
				: state === 'current'
					? 'h-2.5 w-2.5 bg-primary'
					: 'bg-neutral-300 dark:bg-neutral-700'}"
			onclick={() => wizard.goToStep(i)}
			disabled={state === 'upcoming'}
			aria-label="Go to step {i + 1}: {step.title}"
		></button>
	{/each}
</div>

<!-- Desktop step indicator - sidebar -->
<div class="absolute left-1/2 top-36 hidden w-full max-w-2xl -translate-x-1/2 lg:block">
	<div class="absolute right-[calc(100%+2rem)] flex flex-col gap-1">
		{#each wizard.steps as step, i}
			{@const state = wizard.getStepState(i)}
			<Button
				variant="ghost"
				class="flex w-full items-center justify-start gap-2 text-left"
				onclick={() => wizard.goToStep(i)}
				disabled={state === 'upcoming'}
			>
				<div
					class="h-1.5 w-1.5 rounded-full transition-colors
						{state === 'completed'
						? 'bg-emerald-500'
						: state === 'current'
							? 'bg-foreground'
							: 'border border-neutral-400 bg-neutral-200'}"
				></div>
				<span class="whitespace-nowrap text-sm font-medium">{step.title}</span>
			</Button>
		{/each}
	</div>
</div>
