<script lang="ts">
	import Icon from '$lib/components/ui/icon/Icon.svelte';
	import { resolve } from '$app/paths';

	let { count = 0 }: { count?: number } = $props();
</script>

<a
	href={resolve('/checkout')}
	class="relative flex items-center justify-center rounded-md p-2 text-foreground-muted transition-colors hover:bg-accent/50 hover:text-foreground"
>
	<Icon
		icon="shopping-cart"
		class="h-5 w-5 fill-foreground-muted transition-colors hover:fill-foreground"
	/>
	{#if count > 0}
		<span
			class="absolute -right-1 -top-1 flex h-5 min-w-5 items-center justify-center rounded-full bg-primary px-1.5 text-xs font-medium text-primary-foreground"
		>
			{count > 99 ? '99+' : count}
		</span>
	{/if}
</a>
