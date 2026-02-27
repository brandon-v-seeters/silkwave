<script lang="ts">
	import type { HTMLInputAttributes } from 'svelte/elements';
	import type { WithElementRef } from 'bits-ui';
	import { cn } from '$lib/utils/utils';

	let {
		ref = $bindable(null),
		value = $bindable(),
		class: className,
		prefix,
		...restProps
	}: WithElementRef<HTMLInputAttributes> & { prefix?: string } = $props();
</script>

{#if prefix}
	<div class="relative">
		<span
			class="pointer-events-none absolute left-3 top-1/2 -translate-y-1/2 text-foreground-muted"
		>
			{prefix}
		</span>
		<input
			bind:this={ref}
			class={cn(
				'h-12 bg-muted-background disabled:bg-muted-background',
				'flex w-full rounded-lg border border-input py-2 pr-3 text-base',
				'file:border-0 file:bg-transparent file:text-base file:font-medium placeholder:text-foreground-muted',
				'focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-primary/60 disabled:cursor-not-allowed disabled:opacity-50',
				'focus-visible:shadow-focus-primary',
				'aria-[invalid=true]:border-rose-400 aria-[invalid=true]:focus-visible:shadow-none',
				'aria-[invalid=true]:focus-visible:ring-2 aria-[invalid=true]:focus-visible:ring-destructive/50',
				'pl-8',
				className
			)}
			bind:value
			{...restProps}
		/>
	</div>
{:else}
	<input
		bind:this={ref}
		class={cn(
			'h-12 bg-muted-background disabled:bg-muted-background',
			'flex w-full rounded-lg border border-input px-3 py-2 text-base',
			'file:border-0 file:bg-transparent file:text-base file:font-medium placeholder:text-foreground-muted',
			'focus-visible:outline-none  focus-visible:ring-1 focus-visible:ring-primary/60 disabled:cursor-not-allowed disabled:opacity-50',
			'focus-visible:shadow-focus-primary',
			'aria-[invalid=true]:border-rose-400 aria-[invalid=true]:focus-visible:shadow-none',
			'aria-[invalid=true]:focus-visible:ring-2 aria-[invalid=true]:focus-visible:ring-destructive/50',
			className
		)}
		bind:value
		{...restProps}
	/>
{/if}
