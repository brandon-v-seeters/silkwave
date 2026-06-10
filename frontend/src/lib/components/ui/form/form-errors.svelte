<script lang="ts">
	import { fromStore } from 'svelte/store';
	import type { SuperForm } from 'sveltekit-superforms';
	import Icon from '../icon/Icon.svelte';

	let { form, class: className = '' }: { form: SuperForm<any>; class?: string } = $props();
	let errors = $derived(fromStore(form.errors));
	let errorMessages = $derived(errors.current._errors ?? []);
</script>

{#if errorMessages.length}
	<div
		class="flex items-center gap-2 rounded-md bg-rose-400/5 px-3 py-2 text-sm text-rose-400 {className}"
	>
		<Icon icon="alert-circle" class="h-6 w-6 flex-shrink-0 fill-rose-400" />
		{#each errorMessages as error, index (`${index}-${error}`)}
			<p>{error}</p>
		{/each}
	</div>
{/if}
