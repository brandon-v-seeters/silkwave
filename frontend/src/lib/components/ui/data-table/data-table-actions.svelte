<script lang="ts">
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import Icon from '$lib/components/atoms/Icon.svelte';
	import type { IconKey } from '$lib/types/Icon';

	type Action = {
		label: string;
		onclick: () => void;
		icon?: IconKey;
		variant?: 'default' | 'destructive';
	};

	let {
		actions = []
	}: {
		actions: Action[];
	} = $props();
</script>

<div class="flex justify-end">
	<DropdownMenu.Root>
		<DropdownMenu.Trigger>
			{#snippet child({ props })}
				<Button {...props} variant="ghost" size="icon" class="h-8 w-8">
					<Icon icon="more-vert" class="h-5 w-5 fill-foreground" />
					<span class="sr-only">Open menu</span>
				</Button>
			{/snippet}
		</DropdownMenu.Trigger>
		<DropdownMenu.Content align="end">
			{#each actions as action}
				<DropdownMenu.Item
					onclick={action.onclick}
					class={action.variant === 'destructive' ? 'text-destructive-foreground' : ''}
				>
					{#if action.icon}
						<Icon icon={action.icon} class="mr-2 h-4 w-4 fill-current" />
					{/if}
					{action.label}
				</DropdownMenu.Item>
			{/each}
		</DropdownMenu.Content>
	</DropdownMenu.Root>
</div>
