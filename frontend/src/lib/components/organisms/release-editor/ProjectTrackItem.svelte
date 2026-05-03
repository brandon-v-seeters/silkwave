<script lang="ts">
	import Icon from '$lib/components/atoms/Icon.svelte';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
	import type { WizardTrack } from '$lib/types/WizardTrack';

	interface Props {
		track: WizardTrack;
		index: number;
		total: number;
		onTitleChange?: (title: string) => void;
		onRemove?: () => void;
		onMove?: (direction: 'up' | 'down') => void;
	}

	let { track, index, total, onTitleChange, onRemove, onMove }: Props = $props();

	let editing = $state(false);
	let editValue = $state('');
	let inputEl = $state<HTMLInputElement | null>(null);

	function startEdit() {
		editValue = track.title;
		editing = true;
		queueMicrotask(() => inputEl?.select());
	}

	function confirmEdit() {
		if (editValue.trim()) {
			onTitleChange?.(editValue.trim());
		}
		editing = false;
	}

	function cancelEdit() {
		editing = false;
	}
</script>

<li
	class="group/row flex w-full select-none items-center rounded-2xl py-2 transition-colors duration-200 hover:bg-foreground/3 sm:px-3"
>
	<div class="mr-2 w-[30px] shrink-0 text-center font-mono text-xs text-foreground-muted/70">
		<span class="flex items-center justify-center">{String(index + 1).padStart(2, '0')}</span>
	</div>

	<div class="flex w-full min-w-0 flex-col">
		{#if editing}
			<!-- svelte-ignore a11y_autofocus -->
			<input
				bind:this={inputEl}
				class="mr-3 w-full rounded-lg bg-transparent px-2 py-1 text-sm font-semibold text-foreground outline-none ring-1 ring-border"
				bind:value={editValue}
				autofocus
				onkeydown={(e) => {
					if (e.key === 'Enter') confirmEdit();
					if (e.key === 'Escape') cancelEdit();
				}}
				onblur={confirmEdit}
			/>
		{:else}
			<h3 class="line-clamp-1 text-left text-sm font-semibold break-all">{track.title}</h3>
		{/if}
		<div class="line-clamp-1 text-left text-xs break-all text-foreground-muted">
			{track.duration || '0:00'}
		</div>
	</div>

	<div class="ml-3 flex shrink-0 items-center">
		<DropdownMenu.Root>
			<DropdownMenu.Trigger
				class="flex h-8 w-8 items-center justify-center rounded-md text-foreground-muted opacity-0 transition-all duration-150 hover:bg-foreground/5 hover:text-foreground focus-visible:opacity-100 data-[state=open]:opacity-100 group-hover/row:opacity-100 sm:opacity-0"
				aria-label="Track options"
			>
				<Icon icon="more-horiz" variant="line" class="h-4 w-4 fill-current" />
			</DropdownMenu.Trigger>
			<DropdownMenu.Content align="end" class="w-44">
				<DropdownMenu.Item onclick={startEdit}>
					<Icon icon="edit" variant="line" class="h-4 w-4 fill-current" />
					<span>Rename</span>
				</DropdownMenu.Item>
				{#if onMove}
					<DropdownMenu.Item disabled={index === 0} onclick={() => onMove?.('up')}>
						<Icon icon="chevron-up" variant="line" class="h-4 w-4 fill-current" />
						<span>Move up</span>
					</DropdownMenu.Item>
					<DropdownMenu.Item
						disabled={index === total - 1}
						onclick={() => onMove?.('down')}
					>
						<Icon icon="chevron-down" variant="line" class="h-4 w-4 fill-current" />
						<span>Move down</span>
					</DropdownMenu.Item>
				{/if}
				<DropdownMenu.Separator />
				<DropdownMenu.Item
					class="text-destructive data-highlighted:bg-destructive/10 data-highlighted:text-destructive"
					onclick={onRemove}
				>
					<Icon icon="trash" variant="line" class="h-4 w-4 fill-current" />
					<span>Remove</span>
				</DropdownMenu.Item>
			</DropdownMenu.Content>
		</DropdownMenu.Root>
	</div>
</li>
