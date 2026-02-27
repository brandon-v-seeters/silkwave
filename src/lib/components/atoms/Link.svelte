<script lang="ts">
	import type { HTMLAnchorAttributes } from 'svelte/elements';

	let {
		href,
		class: className,
		children
	}: HTMLAnchorAttributes & { children?: any } = $props();
</script>

<a {href} class="{className} text-link hover:text-link-hover transition-all duration-300">
	{@render children?.()}
</a>
