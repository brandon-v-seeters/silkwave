<script lang="ts">
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import Icon from '$lib/components/ui/icon/Icon.svelte';
	import Button from '$lib/components/ui/button/button.svelte';
	import type { IconKey } from '$lib/types/Icon';

	let { mediaPlayerActive = false }: { mediaPlayerActive?: boolean } = $props();

	type NavItem = {
		title: string;
		href: string;
		icon: IconKey;
	};

	const browseLinks: NavItem[] = [
		{ title: 'Home', href: '/', icon: 'home' },
		{ title: 'Albums', href: '/discover', icon: 'grid' },
		{ title: 'Tracks', href: '/discover?type=tracks', icon: 'music-note' },
		{ title: 'Genres', href: '/discover?view=genres', icon: 'playlist' }
	];

	function isActive(href: string) {
		if (href === '/') return page.url.pathname === '/';
		return page.url.pathname === href || page.url.href.includes(href);
	}
</script>

<aside
	class="hidden min-h-0 flex-col overflow-y-auto bg-muted/20 px-6 py-7 transition-[height] duration-200 ease-out lg:flex {mediaPlayerActive
		? 'lg:h-[calc(100vh-5.75rem)]'
		: 'lg:h-screen'}"
>
	<div class="mb-8 flex items-center justify-between">
		<a href={resolve('/')} class="text-sm font-bold tracking-tight">Silkwave</a>
		<button
			type="button"
			class="rounded-full p-1.5 text-muted-foreground transition hover:bg-muted hover:text-foreground"
			aria-label="Collapse navigation"
		>
			<Icon icon="arrow-left" class="h-3.5 w-3.5 fill-current" />
		</button>
	</div>

	<nav class="space-y-8 text-[0.78rem]">
		<div>
			<p class="mb-3 text-[0.68rem] font-semibold text-foreground">Browse Music</p>
			<div class="space-y-1">
				{#each browseLinks as item (item.href)}
					<a
						href={resolve(item.href as '/')}
						class="group flex items-center gap-2.5 rounded-md px-2.5 py-2 transition {isActive(
							item.href
						)
							? 'bg-muted text-foreground shadow-sm'
							: 'text-muted-foreground hover:bg-muted/60 hover:text-foreground'}"
					>
						<Icon icon={item.icon} class="h-3.5 w-3.5 fill-current opacity-70" />
						<span>{item.title}</span>
					</a>
				{/each}
			</div>
		</div>
	</nav>

	<div class="mt-auto rounded-2xl bg-muted/55 p-3">
		<p class="text-[0.7rem] font-semibold text-foreground">For artists</p>
		<p class="mt-1 text-[0.68rem] leading-relaxed text-muted-foreground">
			Upload a release and keep the storefront quiet.
		</p>
		<Button href={resolve('/upload')} variant="primary" size="sm" class="mt-3 w-full text-xs">
			Upload Music
		</Button>
	</div>
</aside>
