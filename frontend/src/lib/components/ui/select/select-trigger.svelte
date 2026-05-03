<script lang="ts">
	import { Select as SelectPrimitive } from "bits-ui";
	import ChevronDownIcon from "@lucide/svelte/icons/chevron-down";
	import { cn, type WithoutChild } from "$lib/utils/utils.js";

	let {
		ref = $bindable(null),
		class: className,
		children,
		...restProps
	}: WithoutChild<SelectPrimitive.TriggerProps> = $props();
</script>

<SelectPrimitive.Trigger
	bind:ref
	data-slot="select-trigger"
	class={cn(
		// Base styles matching input
		'h-12 bg-muted-background disabled:bg-muted-background',
		'flex w-full items-center justify-between gap-2 rounded-lg border border-input px-3 py-2 text-base',
		'placeholder:text-foreground-muted',
		// Focus styles matching input
		'focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-primary/60',
		'focus-visible:shadow-focus-primary',
		// Disabled & invalid states
		'disabled:cursor-not-allowed disabled:opacity-50',
		'aria-invalid:border-rose-400 aria-invalid:focus-visible:shadow-none',
		'aria-invalid:focus-visible:ring-2 aria-invalid:focus-visible:ring-destructive/50',
		// Select-specific
		'select-none whitespace-nowrap transition-colors',
		'data-[placeholder]:text-foreground-muted',
		'*:data-[slot=select-value]:line-clamp-1 *:data-[slot=select-value]:flex *:data-[slot=select-value]:items-center *:data-[slot=select-value]:gap-2',
		'[&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*="size-"])]:size-4',
		className
	)}
	{...restProps}
>
	{@render children?.()}
	<ChevronDownIcon class="size-4 opacity-50" />
</SelectPrimitive.Trigger>
