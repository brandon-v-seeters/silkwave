<script lang="ts">
	import MainNavbar from '$lib/components/organisms/MainNavbar.svelte';
	import SearchModal from '$lib/components/organisms/SearchModal.svelte';
	import { searchModalOpen } from '$lib/stores/ui';

	let { children } = $props();
</script>

<div class="flex min-h-screen w-full flex-col bg-background">
	<main class="flex-1">
		<div class="mx-auto w-full px-3 py-4 pb-20 sm:px-4 sm:py-6 md:px-6 md:py-8 md:pb-8">
			{@render children?.()}
		</div>
	</main>
</div>
