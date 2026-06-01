<script lang="ts">
	import Breadcrumbs from '$lib/components/molecules/Breadcrumbs.svelte';

	let { children } = $props();
</script>

<div class="flex min-h-screen w-full flex-col bg-background">
	<main class="flex-1">
		<div class="mx-auto w-full px-3 py-4 pb-20 sm:px-4 sm:py-6 md:px-6 md:py-8 md:pb-8">
			<Breadcrumbs class="mb-5" />
			{@render children?.()}
		</div>
	</main>
</div>
