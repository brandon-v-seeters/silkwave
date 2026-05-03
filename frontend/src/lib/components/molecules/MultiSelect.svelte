<script lang="ts">
	import Icon from '$lib/components/atoms/Icon.svelte';

	interface Props {
		/** Currently selected values */
		selected?: string[];
		/** Available options to choose from */
		options?: string[];
		/** Placeholder text when empty */
		placeholder?: string;
		/** Placeholder for the search input */
		searchPlaceholder?: string;
		/** Allow adding custom values not in options */
		allowCustom?: boolean;
		/** Callback when selection changes */
		onchange?: (selected: string[]) => void;
	}

	let {
		selected = $bindable([]),
		options = [],
		placeholder = 'Select items...',
		searchPlaceholder = 'Search...',
		allowCustom = true,
		onchange
	}: Props = $props();

	let searchValue = $state('');
	let isOpen = $state(false);
	let inputEl: HTMLInputElement | undefined = $state();
	let containerEl: HTMLDivElement | undefined = $state();

	// Available options (not yet selected)
	const availableOptions = $derived(options.filter((o) => !selected.includes(o)));

	// Filtered options based on search
	const filteredOptions = $derived(
		searchValue
			? availableOptions.filter((o) => o.toLowerCase().includes(searchValue.toLowerCase()))
			: availableOptions
	);

	// Check if search value could be added as custom
	const canAddCustom = $derived(
		allowCustom &&
			searchValue.trim() &&
			!options.some((o) => o.toLowerCase() === searchValue.toLowerCase()) &&
			!selected.some((s) => s.toLowerCase() === searchValue.toLowerCase())
	);

	const addItem = (item: string) => {
		const trimmed = item.trim();
		if (trimmed && !selected.includes(trimmed)) {
			selected = [...selected, trimmed];
			onchange?.(selected);
		}
		searchValue = '';
		inputEl?.focus();
	};

	const removeItem = (item: string, e?: MouseEvent) => {
		e?.stopPropagation();
		selected = selected.filter((s) => s !== item);
		onchange?.(selected);
	};

	const handleKeydown = (e: KeyboardEvent) => {
		if (e.key === 'Enter') {
			e.preventDefault();
			if (filteredOptions.length > 0) {
				addItem(filteredOptions[0]);
			} else if (canAddCustom) {
				addItem(searchValue);
			}
		}
		if (e.key === 'Backspace' && !searchValue && selected.length > 0) {
			selected = selected.slice(0, -1);
			onchange?.(selected);
		}
		if (e.key === 'Escape') {
			isOpen = false;
			inputEl?.blur();
		}
	};

	const toggleDropdown = () => {
		isOpen = !isOpen;
		if (isOpen) {
			// Focus input when opening
			setTimeout(() => inputEl?.focus(), 0);
		}
	};

	const handleClickOutside = (e: MouseEvent) => {
		if (containerEl && !containerEl.contains(e.target as Node)) {
			isOpen = false;
		}
	};

	$effect(() => {
		if (isOpen) {
			document.addEventListener('click', handleClickOutside);
			return () => document.removeEventListener('click', handleClickOutside);
		}
	});
</script>

<!-- svelte-ignore a11y_no_static_element_interactions -->
<!-- svelte-ignore a11y_click_events_have_key_events -->
<div class="relative" bind:this={containerEl}>
	{#if selected.length > 0}
		<button
			class="group absolute -top-6 right-0 flex flex-row items-center gap-1"
			onclick={() => (selected = [])}
		>
			<Icon icon="cross" class="h-4 w-4 fill-foreground-muted group-hover:fill-foreground" />
			<span class="cursor-pointer text-sm text-foreground-muted group-hover:text-foreground">
				Clear all
			</span>
		</button>
	{/if}
	<!-- Trigger -->
	<div
		onclick={toggleDropdown}
		class="flex min-h-12 w-full cursor-pointer items-center justify-between gap-2 rounded-lg border border-input bg-muted-background px-3 py-2 text-left transition-all"
		class:shadow-focus-primary={isOpen}
	>
		<div class="flex flex-1 flex-wrap items-center gap-2">
			{#if selected.length > 0}
				{#each selected as item}
					<span
						class="inline-flex items-center gap-1 rounded-full bg-gradient-to-r from-foreground/20 to-foreground/10 px-3 py-1 text-sm text-foreground"
					>
						{item}
						<button
							type="button"
							onclick={(e) => removeItem(item, e)}
							class="ml-1 rounded-full p-0.5 transition-colors hover:bg-indigo-200/20 hover:fill-foreground"
						>
							<Icon icon="cross" class="h-4 w-4 fill-foreground-muted" />
						</button>
					</span>
				{/each}
			{:else}
				<span class="text-foreground-muted">{placeholder}</span>
			{/if}
		</div>
		<Icon
			icon="chevron-down"
			class="h-4 w-4 flex-shrink-0 fill-foreground-muted transition-transform duration-200 {isOpen
				? 'rotate-180'
				: ''}"
		/>
	</div>

	<!-- Dropdown -->
	{#if isOpen}
		<div
			class="absolute left-0 right-0 top-full z-10 mt-1 overflow-hidden rounded-md border border-input bg-card shadow-lg"
		>
			<!-- Search Input -->
			<div class="border-b border-input p-2">
				<div class="flex items-center gap-2 rounded-md bg-muted-background px-3 py-2">
					<Icon icon="search" class="h-4 w-4 flex-shrink-0 fill-foreground-muted" />
					<input
						bind:this={inputEl}
						type="text"
						bind:value={searchValue}
						onkeydown={handleKeydown}
						placeholder={searchPlaceholder}
						class="w-full bg-transparent text-sm outline-none placeholder:text-foreground-muted"
					/>
				</div>
			</div>

			<!-- Options List -->
			<div class="max-h-48 overflow-y-auto">
				{#if filteredOptions.length > 0}
					{#each filteredOptions as option}
						<button
							type="button"
							onclick={() => addItem(option)}
							class="flex w-full items-center gap-2 px-4 py-2 text-left text-sm transition-colors hover:bg-primary/5 hover:text-primary"
						>
							{option}
						</button>
					{/each}
				{:else if canAddCustom}
					<button
						type="button"
						onclick={() => addItem(searchValue)}
						class="flex w-full items-center gap-2 px-4 py-2 text-left text-sm transition-colors hover:bg-primary/5 hover:text-primary"
					>
						<Icon icon="plus" class="h-4 w-4 fill-primary" />
						<span
							>Add "<span class="font-medium text-primary">{searchValue}</span>"</span
						>
					</button>
				{:else if searchValue}
					<div class="px-4 py-3 text-center text-sm text-foreground-muted">
						No results found
					</div>
				{:else}
					<div class="px-4 py-3 text-center text-sm text-foreground-muted">
						No options available
					</div>
				{/if}
			</div>
		</div>
	{/if}
</div>
