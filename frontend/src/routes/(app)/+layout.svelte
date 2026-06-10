<script lang="ts">
	import { resolve } from '$app/paths';
	import Icon from '$lib/components/ui/icon/Icon.svelte';
	import DarkModeToggle from '$lib/components/ui/dark-mode-toggle/DarkModeToggle.svelte';
	import { NavbarCart } from '$lib/features/cart';
	import Breadcrumbs from './breadcrumbs.svelte';
	import { SearchModal } from '$lib/features/catalog';
	import { cartCount } from '$lib/features/cart';
	import { searchModalOpen } from '$lib/features/catalog';
	import AppSidebar from './app-sidebar.svelte';

	let { children } = $props();

	type ActiveTrack = {
		title: string;
		artistName: string;
		artworkUrl: string;
	};

	let activeTrack = $state<ActiveTrack | null>({
		title: 'Lady Magnolia',
		artistName: 'Piero Umiliani',
		artworkUrl: '/banner.jpg'
	});

	let mediaPlayerActive = $derived(activeTrack !== null);

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
		<div class="flex h-full w-full flex-col overflow-hidden bg-background">
			<SearchModal bind:open={$searchModalOpen} />

			<div
				class="flex items-center justify-between gap-3 bg-muted/25 px-4 py-3 lg:hidden sticky top-0"
			>
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
				<AppSidebar {mediaPlayerActive} />

				<main class="min-h-0 overflow-y-auto {mediaPlayerActive ? 'pb-32' : 'pb-8'}">
					<div
						class="mb-5 flex flex-col gap-4 md:flex-row md:items-center md:justify-between sticky top-0 bg-background z-20 py-4 pr-4"
					>
						<div class="flex min-w-0 items-center gap-3 text-[0.72rem]">
							<button
								type="button"
								onclick={openSearch}
								class="flex h-8 w-8 shrink-0 items-center justify-center rounded-full text-muted-foreground transition hover:bg-muted hover:text-foreground"
								aria-label="Search"
							>
								<Icon icon="search" class="h-3.5 w-3.5 fill-current" />
							</button>
							<Breadcrumbs class="min-w-0" />
						</div>

						<div class="hidden items-center gap-2 md:flex">
							<DarkModeToggle />
						</div>
					</div>

					<div class="pr-6">
						{@render children?.()}
					</div>
				</main>
			</div>

			{#if activeTrack}
				<div
					class="absolute inset-x-0 bottom-0 isolate z-30 h-23 overflow-hidden bg-zinc-950/72 px-5 text-white shadow-[0_-24px_70px_rgba(0,0,0,0.24)] backdrop-blur-2xl sm:px-8 dark:bg-black/48"
				>
					<div
						class="relative z-10 grid h-full grid-cols-[minmax(0,1fr)_auto] items-center gap-4 lg:grid-cols-[18rem_minmax(0,1fr)_18rem]"
					>
						<div class="flex min-w-0 items-center gap-3">
							<img
								src={activeTrack.artworkUrl}
								alt=""
								class="h-12 w-12 rounded-md object-cover shadow-[0_12px_30px_rgba(0,0,0,0.4)]"
							/>
							<div class="min-w-0">
								<p
									class="truncate text-[0.82rem] font-semibold leading-tight text-white"
								>
									{activeTrack.title}
								</p>
								<p class="mt-1 truncate text-[0.68rem] text-white/48">
									{activeTrack.artistName}
								</p>
							</div>
						</div>

						<div
							class="absolute inset-x-5 bottom-4 z-0 hidden h-px bg-white/18 sm:block lg:inset-x-[24rem]"
						>
							<div class="h-px w-[32%] bg-white/76"></div>
						</div>

						<div
							class="hidden min-w-0 flex-col items-center justify-center gap-3 lg:flex"
						>
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
			{/if}
		</div>
	</div>
{/if}
