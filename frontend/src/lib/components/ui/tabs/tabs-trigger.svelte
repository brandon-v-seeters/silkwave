<script lang="ts">
	import { Tabs as TabsPrimitive } from 'bits-ui';
	import { cn } from '$lib/utils/utils.js';
	import type { ComponentProps } from 'svelte';

	type Props = ComponentProps<typeof TabsPrimitive.Trigger>;

	let {
		class: className,
		value,
		children,
		...restProps
	}: Props = $props();
</script>

<TabsPrimitive.Trigger
	class={cn(
		'mr-2 inline-flex items-center justify-center whitespace-nowrap rounded-sm px-3 py-1 text-base transition-all',
		'text-foreground-muted hover:text-foreground',
		'ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2',
		'border-b-2 border-transparent disabled:pointer-events-none disabled:opacity-50',
		'data-[state=active]:bg-primary/[5%] data-[state=active]:text-primary',
		className
	)}
	{value}
	{...restProps}
>
	{@render children?.()}
</TabsPrimitive.Trigger>
