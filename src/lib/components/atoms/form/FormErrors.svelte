<script lang="ts">
	import type { SuperForm } from 'sveltekit-superforms';
	import Icon from '../Icon.svelte';

	let { form, class: className = '' }: { form: SuperForm<any>; class?: string } = $props();
	const errors = $derived(form.errors);
</script>

{#if $errors._errors?.length}
	<div
		class="flex items-center gap-2 rounded-md bg-rose-400/5 px-3 py-2 text-sm text-rose-400 {className}"
	>
		<Icon icon="alert-circle" class="h-6 w-6 flex-shrink-0 fill-rose-400" />
		{#each $errors._errors as error}
			<p>{error}</p>
		{/each}
	</div>
{/if}
