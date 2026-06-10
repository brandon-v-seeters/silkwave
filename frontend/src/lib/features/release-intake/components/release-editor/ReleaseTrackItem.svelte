<script lang="ts">
	import { tick } from 'svelte';
	import Icon from '$lib/components/ui/icon/Icon.svelte';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
	import type { WizardTrack } from '$lib/features/release-intake/types';

	interface Props {
		track: WizardTrack;
		index: number;
		total: number;
		canPlay?: boolean;
		onTitleChange?: (title: string) => void;
		onRemove?: () => void;
		onMove?: (direction: 'up' | 'down') => void;
		onPlay?: () => void;
	}

	let {
		track,
		index,
		total,
		canPlay = false,
		onTitleChange,
		onRemove,
		onMove,
		onPlay
	}: Props = $props();

	let editing = $state(false);
	let editValue = $state('');
	let inputEl = $state<HTMLInputElement | null>(null);

	function startEdit() {
		editValue = track.title;
		editing = true;

		setTimeout(async () => {
			await tick();
			inputEl?.focus();
			inputEl?.select();
		});
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

	function playTrack(event: MouseEvent) {
		event.stopPropagation();
		onPlay?.();
	}
</script>

<li
	class="group/row flex w-full select-none items-center rounded-2xl py-2 transition-colors duration-200 hover:bg-foreground/3 sm:px-3"
>
	<div class="relative mr-2 flex w-[30px] shrink-0 items-center justify-center">
		<span
			class="flex items-center justify-center font-mono text-xs text-foreground-muted/70 transition-opacity duration-150 {canPlay
				? 'group-hover/row:opacity-0 group-focus-within/row:opacity-0'
				: ''}"
		>
			{String(index + 1).padStart(2, '0')}
		</span>

		{#if canPlay}
			<button
				type="button"
				onclick={playTrack}
				class="absolute flex h-7 w-7 items-center justify-center rounded-full text-foreground opacity-0 transition duration-150 hover:bg-foreground/8 focus-visible:bg-foreground/8 focus-visible:opacity-100 focus-visible:outline-none group-hover/row:opacity-100 group-focus-within/row:opacity-100"
				aria-label="Play {track.title}"
			>
				<Icon icon="play" variant="filled" class="h-3.5 w-3.5 fill-current" />
			</button>
		{/if}
	</div>

	<div class="flex w-full min-w-0 flex-col">
		{#if editing}
			<input
				bind:this={inputEl}
				type="text"
				class="mr-3 w-full rounded-lg bg-transparent px-2 py-1 text-sm font-semibold text-foreground outline-none ring-1 ring-border"
				bind:value={editValue}
				onkeydown={(e) => {
					if (e.key === 'Enter') confirmEdit();
					if (e.key === 'Escape') cancelEdit();
				}}
				onblur={confirmEdit}
			/>
		{:else}
			<h5 class="line-clamp-1 text-left text-sm font-semibold break-all">{track.title}</h5>
		{/if}
		<div class="line-clamp-1 text-left text-xs break-all text-foreground-muted">
			{track.duration || '0:00'}
		</div>
	</div>

	<div class="ml-3 flex shrink-0 items-center">
		<DropdownMenu.Root>
			<DropdownMenu.Trigger
				class="flex h-8 w-8 items-center justify-center rounded-md text-foreground-muted opacity-100 transition-all duration-150 hover:bg-foreground/5 hover:text-foreground focus-visible:opacity-100 data-[state=open]:opacity-100 sm:opacity-0 sm:group-hover/row:opacity-100"
				aria-label="Track options"
			>
				<Icon icon="more-horiz" variant="line" class="h-4 w-4 fill-current" />
			</DropdownMenu.Trigger>
			<DropdownMenu.Content align="end" class="w-44">
				<DropdownMenu.Item onSelect={startEdit}>
					<Icon icon="edit" variant="line" class="h-4 w-4 fill-current" />
					<span>Rename</span>
				</DropdownMenu.Item>
				{#if onMove}
					<DropdownMenu.Item disabled={index === 0} onSelect={() => onMove?.('up')}>
						<Icon icon="chevron-up" variant="line" class="h-4 w-4 fill-current" />
						<span>Move up</span>
					</DropdownMenu.Item>
					<DropdownMenu.Item
						disabled={index === total - 1}
						onSelect={() => onMove?.('down')}
					>
						<Icon icon="chevron-down" variant="line" class="h-4 w-4 fill-current" />
						<span>Move down</span>
					</DropdownMenu.Item>
				{/if}
				<DropdownMenu.Separator />
				<DropdownMenu.Item
					class="text-destructive data-highlighted:bg-destructive/10 data-highlighted:text-destructive"
					onSelect={onRemove}
				>
					<Icon icon="trash" variant="line" class="h-4 w-4 fill-current" />
					<span>Remove</span>
				</DropdownMenu.Item>
			</DropdownMenu.Content>
		</DropdownMenu.Root>
	</div>
</li>
