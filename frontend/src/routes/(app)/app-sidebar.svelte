<script lang="ts">
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import DarkModeToggle from '$lib/components/atoms/DarkModeToggle.svelte';
	import Icon from '$lib/components/atoms/Icon.svelte';
	import NavbarCart from '$lib/components/atoms/NavbarCart.svelte';
	import Button from '$lib/components/ui/button/button.svelte';
	import { cartCount } from '$lib/stores/cart';
	import type { IconKey } from '$lib/types/Icon';
	import { getContext } from 'svelte';

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

	const libraryLinks: NavItem[] = [
		{
			title: 'Recently Played',
			href: '/discover?sort=recent',
			icon: 'clock-plus'
		},
		{
			title: 'Favorite Tracks',
			href: '/discover?filter=favorites',
			icon: 'heart'
		},
		{ title: 'Charts', href: '/discover?view=charts', icon: 'chart-line' },
		{ title: 'Radio', href: '/discover?view=radio', icon: 'radio' }
	];

	function isActive(href: string) {
		if (href === '/') return page.url.pathname === '/';
		return page.url.pathname === href || page.url.href.includes(href);
	}

	const userCtx = getContext('user');
</script>

<aside class="hidden min-h-0 flex-col bg-background px-8 py-7 lg:flex">
	<div class="mb-8 flex items-center justify-between">
		<a href={resolve('/')} class="text-sm font-bold tracking-tight">Silkwave</a>
	</div>

	<nav class="space-y-8">
		<div>
			<p class="mb-3 font-semibold text-foreground">Browse Music</p>
			<div class="space-y-1">
				{#each browseLinks as item (item.href)}
					<a
						href={resolve(item.href as '/')}
						class="group flex items-center gap-2.5 rounded-md px-2.5 py-2 transition {isActive(
							item.href
						)
							? 'bg-linear-to-br from-foreground/15 to-foreground/5 text-foreground shadow-sm'
							: 'text-muted-foreground hover:bg-muted/60 hover:text-foreground'}"
					>
						<Icon icon={item.icon} class="h-3.5 w-3.5 fill-current opacity-70" />
						<span>{item.title}</span>
					</a>
				{/each}
			</div>
		</div>

		<div>
			<p class="mb-3 font-semibold text-foreground">Library</p>
			<div class="space-y-1">
				{#each libraryLinks as item (item.href)}
					<a
						href={resolve(item.href as '/')}
						class="group flex items-center gap-2.5 rounded-md px-2.5 py-2 text-muted-foreground/55 transition hover:bg-muted/50 hover:text-foreground"
					>
						<Icon icon={item.icon} class="h-3.5 w-3.5 fill-current opacity-65" />
						<span>{item.title}</span>
					</a>
				{/each}
			</div>
		</div>
	</nav>

	<div class="mt-auto rounded-2xl bg-muted/55 p-3">
		<p class="font-semibold text-foreground">For artists</p>
		<p class="mt-1 leading-relaxed text-muted-foreground">
			Upload a release and keep the storefront quiet.
		</p>
		<Button href={resolve('/upload')} variant="primary" size="sm" class="mt-3 w-full text-xs">
			Upload Music
		</Button>
	</div>

	<div class="flex items-center justify-between mt-8">
		<div class="flex items-center gap-3">
			<img src="/avatar.jpg" alt="" class="h-8 w-8 rounded-xl object-cover" />
			<div>
				<p class="font-medium leading-none">{JSON.stringify(userCtx)}</p>
			</div>
		</div>
	</div>
</aside>
