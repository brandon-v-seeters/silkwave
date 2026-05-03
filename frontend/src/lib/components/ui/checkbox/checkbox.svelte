<script lang="ts">
	import { Checkbox as CheckboxPrimitive, type WithoutChildrenOrChild } from 'bits-ui';
	import { cn } from '$lib/utils/utils.js';
	import Icon from '$lib/components/atoms/Icon.svelte';

	let {
		ref = $bindable(null),
		checked = $bindable(false),
		indeterminate = $bindable(false),
		class: className,
		...restProps
	}: WithoutChildrenOrChild<CheckboxPrimitive.RootProps> = $props();
</script>

<CheckboxPrimitive.Root
	bind:ref
	class={cn(
		'peer box-content size-4 shrink-0 rounded-[4px] border border-primary',
		'ring-offset-background focus-visible:outline-none focus-visible:ring-2',
		'focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed',
		'disabled:opacity-50 data-[disabled=true]:cursor-not-allowed data-[state=checked]:bg-primary',
		'data-[disabled=true]:text-foreground-muted data-[state=checked]:text-primary-foreground data-[disabled=true]:opacity-50',
		className
	)}
	bind:checked
	bind:indeterminate
	{...restProps}
>
	{#snippet children({ checked, indeterminate })}
		<div class="flex size-4 items-center justify-center text-current">
			<Icon
				icon={indeterminate ? 'minus' : 'check'}
				class={cn('h-5 w-5', !checked && 'fill-transparent')}
			/>
		</div>
	{/snippet}
</CheckboxPrimitive.Root>
