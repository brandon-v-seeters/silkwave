<script lang="ts">
	import type { HTMLAnchorAttributes } from 'svelte/elements';

	let { href, class: className, children }: HTMLAnchorAttributes & { children?: any } = $props();
</script>

<a
	{href}
	class="{className} font-medium underline bg-linear-to-r from-link-from to-link-to text-transparent bg-clip-text hover:text-link-to transition-all duration-300"
>
	{@render children?.()}
</a>
