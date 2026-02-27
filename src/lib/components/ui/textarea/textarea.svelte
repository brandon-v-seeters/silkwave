<script lang="ts">
	import { cn, type WithElementRef, type WithoutChildren } from '$lib/utils/utils.js';
	import type { HTMLTextareaAttributes } from 'svelte/elements';

	let {
		ref = $bindable(null),
		value = $bindable(),
		class: className,
		'data-slot': dataSlot = 'textarea',
		...restProps
	}: WithoutChildren<WithElementRef<HTMLTextareaAttributes>> = $props();
</script>

<textarea
	bind:this={ref}
	data-slot={dataSlot}
	class={cn(
		'field-sizing-content flex min-h-16 w-full rounded-lg border border-input bg-muted-background px-3 py-2 text-base',
		'placeholder:text-foreground-muted disabled:cursor-not-allowed disabled:bg-muted-background disabled:opacity-50',
		'focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-primary/60',
		'auto-size:resize focus-visible:shadow-focus-primary',
		'aria-[invalid=true]:border-destructive-foreground aria-[invalid=true]:focus-visible:shadow-none',
		'aria-[invalid=true]:focus-visible:ring-2 aria-[invalid=true]:focus-visible:ring-destructive/50',
		className
	)}
	bind:value
	{...restProps}
></textarea>
