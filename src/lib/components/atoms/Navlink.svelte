<script lang="ts">
	import { cn } from '$lib/utils/utils';
	import Icon from './Icon.svelte';
	import type { IconKey } from '$lib/types/Icon';
	import { page } from '$app/stores';

	let { href, icon, label }: { href: string; icon: IconKey; label: string } = $props();

	const isActive = (path: string) => $page.url.pathname === path;
</script>

<a
	{href}
	class={cn(
		'flex items-center gap-2 rounded-md px-3 py-2 text-base font-medium transition-colors',
		isActive(href)
			? 'bg-accent/50 text-foreground'
			: 'text-foreground-muted hover:bg-accent/50 hover:text-foreground'
	)}
>
	<Icon
		{icon}
		class={cn('h-5 w-5', isActive(href) ? 'fill-foreground' : 'fill-foreground-muted')}
	/>
	<span>{label}</span>
</a>
