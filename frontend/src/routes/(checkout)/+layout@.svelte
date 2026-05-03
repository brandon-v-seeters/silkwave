<script lang="ts">
	let { children } = $props();
</script>

<div class="min-h-screen bg-background">
	<main class="container mx-auto px-4 py-8">
		{@render children?.()}
	</main>
</div>
