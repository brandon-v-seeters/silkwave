<script lang="ts">
	import { toggleMode, mode } from 'mode-watcher';
	import { Sun, Moon } from '@lucide/svelte';

	let { class: className = '' }: { class?: string } = $props();
</script>

<button
	type="button"
	onclick={toggleMode}
	class="flex h-10 w-10 items-center justify-center rounded-lg text-foreground-muted transition-colors hover:bg-accent hover:text-foreground {className}"
	aria-label="Toggle dark mode"
>
	{#if mode.current === 'light'}
		<Moon class="h-5 w-5" />
	{:else}
		<Sun class="h-5 w-5" />
	{/if}
</button>
