<script lang="ts">
	import CheckIcon from "@lucide/svelte/icons/check";
	import { Select as SelectPrimitive } from "bits-ui";
	import { cn, type WithoutChild } from "$lib/utils/utils.js";

	let {
		ref = $bindable(null),
		class: className,
		value,
		label,
		children: childrenProp,
		...restProps
	}: WithoutChild<SelectPrimitive.ItemProps> = $props();
</script>

<SelectPrimitive.Item
	bind:ref
	{value}
	data-slot="select-item"
	class={cn(
		'relative flex w-full cursor-default items-center gap-2 rounded-md py-2 ps-3 pe-8 text-base outline-hidden select-none',
		'data-[highlighted]:bg-primary/10 data-[highlighted]:text-foreground',
		'data-[disabled]:pointer-events-none data-[disabled]:opacity-50',
		'[&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*="size-"])]:size-4',
		'[&_svg:not([class*="text-"])]:text-foreground-muted',
		'*:[span]:last:flex *:[span]:last:items-center *:[span]:last:gap-2',
		className
	)}
	{...restProps}
>
	{#snippet children({ selected, highlighted })}
		<span class="absolute end-2 flex size-4 items-center justify-center">
			{#if selected}
				<CheckIcon class="size-4 text-primary" />
			{/if}
		</span>
		{#if childrenProp}
			{@render childrenProp({ selected, highlighted })}
		{:else}
			{label || value}
		{/if}
	{/snippet}
</SelectPrimitive.Item>
