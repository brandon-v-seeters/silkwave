<script lang="ts">
	import type { LayoutProps } from './$types';

	let { children }: LayoutProps = $props();
</script>

<div class="min-h-dvh bg-background text-foreground">
	<main class="mx-auto w-full max-w-7xl px-4 lg:py-10">
		{@render children?.()}
	</main>
</div>
