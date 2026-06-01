<script lang="ts">
	import { page } from '$app/state';
	import { resolve } from '$app/paths';
	import Icon from '$lib/components/atoms/Icon.svelte';
	import DarkModeToggle from '$lib/components/atoms/DarkModeToggle.svelte';
	import NavbarCart from '$lib/components/atoms/NavbarCart.svelte';
	import SearchModal from '$lib/components/organisms/SearchModal.svelte';
	import Button from '$lib/components/ui/button/button.svelte';
	import { cartCount } from '$lib/stores/cart';
	import { searchModalOpen } from '$lib/stores/ui';
	import type { IconKey } from '$lib/types/Icon';

	let { children } = $props();

	type NavItem = {
		title: string;
		href: string;
		icon: IconKey;
		muted?: boolean;
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
			icon: 'clock-plus',
			muted: true
		},
		{
			title: 'Favorite Tracks',
			href: '/discover?filter=favorites',
			icon: 'heart',
			muted: true
		},
		{ title: 'Charts', href: '/discover?view=charts', icon: 'chart-line', muted: true },
		{ title: 'Radio', href: '/discover?view=radio', icon: 'radio', muted: true }
	];

	const activityItems = [
		{ title: 'Posted by Dungen', meta: '5m', icon: 'folder-music' as IconKey },
		{ title: 'Merchandise', meta: 'New drop', icon: 'shopping-bag' as IconKey }
	];

	const playlists = [
		{ title: 'Indie Sales', date: '2d', color: 'bg-sky-500' },
		{ title: 'Boards of Canada (Full)', date: '4 Oct', color: 'bg-zinc-900' },
		{ title: 'IC122 at Prince Bar', date: '22 Sep', color: 'bg-orange-400' },
		{ title: 'Waves. Playlist', date: '14 Sep', color: 'bg-cyan-500' },
		{ title: 'Library Music', date: '11 Sep', color: 'bg-lime-500' }
	];

	function isActive(href: string) {
		if (href === '/') return page.url.pathname === '/';
		return page.url.pathname === href || page.url.href.includes(href);
	}

	function openSearch() {
		$searchModalOpen = true;
	}
</script>

{#if import.meta.env.PROD}
	<div class="flex h-screen w-full flex-1 flex-col items-center justify-center bg-background">
		<h1 class="text-4xl font-bold">Silk Wave</h1>
		<p class="text-lg text-muted-foreground">Coming soon...</p>
	</div>
{:else}
	<div class="h-screen w-full overflow-hidden bg-background text-foreground">
			<div class="relative flex h-full w-full flex-col overflow-hidden bg-background">
			<SearchModal bind:open={$searchModalOpen} />

			<div class="flex items-center justify-between gap-3 bg-muted/25 px-4 py-3 lg:hidden">
				<a href={resolve('/')} class="text-sm font-semibold tracking-tight">Silkwave</a>
				<div class="flex items-center gap-1">
					<button
						type="button"
						onclick={openSearch}
						class="flex h-9 w-9 items-center justify-center rounded-full text-muted-foreground transition hover:bg-muted hover:text-foreground"
						aria-label="Search"
					>
						<Icon icon="search" class="h-4 w-4 fill-current" />
					</button>
					<DarkModeToggle />
					<NavbarCart count={$cartCount} />
				</div>
			</div>

			<div
				class="grid min-h-0 flex-1 grid-cols-1 lg:grid-cols-[15rem_minmax(0,1fr)] xl:grid-cols-[15rem_minmax(0,1fr)]"
			>
				<aside class="hidden min-h-0 flex-col bg-muted/20 px-6 py-7 lg:flex">
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
							<p class="mb-3 text-[0.68rem] font-semibold text-foreground">
								Browse Music
							</p>
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
										<Icon
											icon={item.icon}
											class="h-3.5 w-3.5 fill-current opacity-70"
										/>
										<span>{item.title}</span>
									</a>
								{/each}
							</div>
						</div>

						<div>
							<p class="mb-3 text-[0.68rem] font-semibold text-foreground">Library</p>
							<div class="space-y-1">
								{#each libraryLinks as item (item.href)}
									<a
										href={resolve(item.href as '/')}
										class="group flex items-center gap-2.5 rounded-md px-2.5 py-2 text-muted-foreground/55 transition hover:bg-muted/50 hover:text-foreground"
									>
										<Icon
											icon={item.icon}
											class="h-3.5 w-3.5 fill-current opacity-65"
										/>
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
						<Button
							href={resolve('/upload')}
							variant="primary"
							size="sm"
							class="mt-3 w-full text-xs"
						>
							Upload Music
						</Button>
					</div>
				</aside>

				<main class="min-h-0 overflow-y-auto px-4 py-5 pb-32 sm:px-6 lg:px-8">
					<div
						class="mb-5 flex flex-col gap-4 md:flex-row md:items-center md:justify-between"
					>
						<div
							class="flex min-w-0 items-center gap-3 text-[0.72rem] text-muted-foreground"
						>
							<button
								type="button"
								onclick={openSearch}
								class="flex h-8 w-8 items-center justify-center rounded-full transition hover:bg-muted hover:text-foreground"
								aria-label="Search"
							>
								<Icon icon="search" class="h-3.5 w-3.5 fill-current" />
							</button>
							<span class="hidden sm:inline">Artists</span>
							<span class="hidden text-border sm:inline">›</span>
							<span class="truncate font-medium text-foreground"
								>Featured Releases</span
							>
						</div>

						<div class="flex items-center gap-4 overflow-x-auto text-[0.72rem]">
							<a href={resolve('/')} class="shrink-0 font-semibold text-foreground"
								>New Releases</a
							>
							<a
								href={resolve('/discover')}
								class="shrink-0 text-muted-foreground transition hover:text-foreground"
							>
								News Feed
							</a>
							<a
								href={resolve('/discover?shuffle=true' as '/')}
								class="shrink-0 text-muted-foreground transition hover:text-foreground"
							>
								Shuffle Play
							</a>
						</div>
					</div>

					{@render children?.()}
				</main>
			</div>

				<div
					class="absolute inset-x-0 bottom-0 isolate z-30 h-23 overflow-hidden bg-zinc-950/72 px-5 text-white shadow-[0_-24px_70px_rgba(0,0,0,0.24)] backdrop-blur-2xl sm:px-8 dark:bg-black/48"
				>
				<div
					class="relative z-10 grid h-full grid-cols-[minmax(0,1fr)_auto] items-center gap-4 lg:grid-cols-[18rem_minmax(0,1fr)_18rem]"
				>
					<div class="flex min-w-0 items-center gap-3">
						<img
							src="/banner.jpg"
							alt=""
							class="h-12 w-12 rounded-md object-cover shadow-[0_12px_30px_rgba(0,0,0,0.4)]"
						/>
						<div class="min-w-0">
							<p
								class="truncate text-[0.82rem] font-semibold leading-tight text-white"
							>
								Lady Magnolia
							</p>
							<p class="mt-1 truncate text-[0.68rem] text-white/48">Piero Umiliani</p>
						</div>
					</div>

					<div
						class="absolute inset-x-5 bottom-4 z-0 hidden h-px bg-white/18 sm:block lg:inset-x-[24rem]"
					>
						<div class="h-px w-[32%] bg-white/76"></div>
					</div>

					<div class="hidden min-w-0 flex-col items-center justify-center gap-3 lg:flex">
						<div class="flex items-center gap-7 text-white/54">
							<button
								type="button"
								class="transition hover:text-white"
								aria-label="Shuffle"
							>
								<Icon icon="shuffle" class="h-3.5 w-3.5 fill-current" />
							</button>
							<button
								type="button"
								class="transition hover:text-white"
								aria-label="Previous track"
							>
								<Icon icon="play-previous" class="h-4 w-4 fill-current" />
							</button>
							<button
								type="button"
								class="flex h-11 w-11 items-center justify-center rounded-full bg-white text-zinc-950 shadow-[0_10px_34px_rgba(255,255,255,0.2)] transition hover:scale-105"
								aria-label="Pause"
							>
								<Icon icon="pause" class="h-4 w-4 fill-current" />
							</button>
							<button
								type="button"
								class="transition hover:text-white"
								aria-label="Next track"
							>
								<Icon icon="play-next" class="h-4 w-4 fill-current" />
							</button>
							<button
								type="button"
								class="transition hover:text-white"
								aria-label="Queue"
							>
								<Icon icon="grid" class="h-3.5 w-3.5 fill-current" />
							</button>
						</div>
					</div>

					<div class="ml-auto flex items-center gap-5 text-white/48">
						<button
							type="button"
							class="hidden transition hover:text-white sm:block"
							aria-label="Lyrics"
						>
							<Icon icon="music-note" class="h-4 w-4 fill-current" />
						</button>
						<button
							type="button"
							class="hidden transition hover:text-white sm:block"
							aria-label="Mini player"
						>
							<Icon icon="miniplayer" class="h-4 w-4 fill-current" />
						</button>
						<Icon icon="soundwave" class="hidden h-4 w-4 fill-current md:block" />
						<div class="hidden h-px w-20 bg-white/18 sm:block">
							<div class="h-px w-3/5 bg-white/72"></div>
						</div>
						<button
							type="button"
							class="transition hover:text-white"
							aria-label="Player menu"
						>
							<Icon icon="menu-2" class="h-4 w-4 fill-current" />
						</button>
					</div>
				</div>
			</div>
		</div>
	</div>
{/if}
