<script lang="ts">
	import SearchModal from '$lib/components/organisms/SearchModal.svelte';
	import MainNavbar from '$lib/components/organisms/MainNavbar.svelte';
	import { searchModalOpen } from '$lib/stores/ui';

	let { children } = $props();
</script>

{#if import.meta.env.PROD}
	<div class="flex min-h-screen w-full flex-col bg-background">
		{@render children?.()}
	</div>
{:else}
	<div class="flex min-h-screen w-full flex-col bg-background">
		<header
			class="fixed bottom-0 z-40 flex h-14 w-full items-center gap-3 border-b border-border/50 bg-background/95 px-3 backdrop-blur supports-[backdrop-filter]:bg-background sm:h-16 sm:gap-4 sm:px-4 md:sticky md:top-0"
		>
			<MainNavbar />
		</header>
		<SearchModal bind:open={$searchModalOpen} />
		<main class="flex-1">
			<div
				class="mx-auto w-full max-w-7xl px-3 py-4 pb-20 sm:px-4 sm:py-6 md:px-6 md:py-8 md:pb-8"
			>
				{@render children?.()}
			</div>
		</main>
	</div>
{/if}
