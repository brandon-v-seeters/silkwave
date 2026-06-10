<script lang="ts">
	import type { Snippet } from 'svelte';
	import Icon from '$lib/components/ui/icon/Icon.svelte';

	let {
		children,
		label,
		sortOrder,
		onclick
	}: {
		children?: Snippet;
		label?: string;
		sortOrder?: 'asc' | 'desc' | false;
		onclick?: (e: MouseEvent) => void;
	} = $props();
</script>

<button
	{onclick}
	type="button"
	class="fit-content -mx-2 flex w-fit cursor-pointer items-center gap-2 rounded-sm px-2 hover:bg-muted-background data-[state=asc]:bg-muted-background data-[state=desc]:bg-muted-background"
>
	{#if label}
		{label}
	{:else if children}
		{@render children()}
	{/if}
	{#if sortOrder === 'asc'}
		<Icon icon="sort-text" class="h-6 w-6 fill-foreground" />
	{:else if sortOrder === 'desc'}
		<Icon icon="sort-text-reverse" class="h-6 w-6 fill-foreground" />
	{/if}
</button>
