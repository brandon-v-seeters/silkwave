<script lang="ts">
	// Vendor
	import { toast } from 'svelte-sonner';
	import { page } from '$app/state';
	import { goto, invalidateAll } from '$app/navigation';
	import { getContext, onMount } from 'svelte';

	// Components
	import Icon from '$lib/components/atoms/Icon.svelte';
	import NavbarCart from '$lib/components/atoms/NavbarCart.svelte';
	import * as NavigationMenu from '$lib/components/ui/navigation-menu/index.js';
	import * as Dialog from '$lib/components/ui/dialog';
	import Button from '$lib/components/ui/button/button.svelte';
	import { UserRound } from '@lucide/svelte';
	import DarkModeToggle from '$lib/components/atoms/DarkModeToggle.svelte';

	// Stores
	import { searchModalOpen } from '$lib/stores/ui';
	import { cartCount } from '$lib/stores/cart';

	// Utils
	import { cn } from '$lib/utils/utils';

	// Types
	import type { HTMLAttributes } from 'svelte/elements';
	import { POST } from '$lib/api/Api';
	import type { AppUser } from '$lib/types/generated';

	let mobileMenuOpen = $state(false);
	let genresExpanded = $state(false);

	const userCtx = getContext<{ current: AppUser }>('user');
	const user = $derived(userCtx.current);

	// Close mobile menu on navigation
	$effect(() => {
		page.url.pathname;
		mobileMenuOpen = false;
	});

	const genres: { title: string; href: string; description: string }[] = [
		{
			title: 'Drum & Bass',
			href: '/drum-and-bass',
			description:
				'Fast-paced breakbeats at 160-180 BPM with heavy basslines. Born from UK rave culture in the early 90s.'
		},
		{
			title: 'Dubstep',
			href: '/dubstep',
			description:
				'Dark, bass-heavy sound with syncopated rhythms and wobble bass. Emerged from South London in the early 2000s.'
		},
		{
			title: 'Electro',
			href: '/electro',
			description:
				'Robotic funk with TR-808 drums and vocoder vocals. Pioneered in early 80s Detroit and New York.'
		},
		{
			title: 'House',
			href: '/house',
			description:
				'Four-on-the-floor grooves built for the dancefloor. Born in Chicago warehouses in the mid-80s.'
		},
		{
			title: 'Techno',
			href: '/techno',
			description:
				'Hypnotic, machine-driven beats and futuristic soundscapes. Created in Detroit in the mid-80s.'
		},
		{
			title: 'Trance',
			href: '/trance',
			description:
				'Euphoric melodies and building arpeggios designed to induce a hypnotic state. Emerged from 90s Germany.'
		},
		{
			title: 'Trap',
			href: '/trap',
			description:
				'Hard-hitting 808s, rolling hi-hats, and Southern hip-hop influences. Originated in Atlanta in the early 2000s.'
		},
		{
			title: 'UK Garage',
			href: '/uk-garage',
			description:
				'Shuffling 2-step rhythms with soulful vocals and skippy bass. Evolved from 90s London club culture.'
		}
	];

	type ListItemProps = HTMLAttributes<HTMLAnchorElement> & {
		title: string;
		href: string;
		content: string;
	};

	function handleSearchbarClick() {
		$searchModalOpen = true;
	}

	function handleKeydown(e: KeyboardEvent) {
		// Cmd+K or Ctrl+K to open search
		if ((e.metaKey || e.ctrlKey) && e.key === 'k') {
			e.preventDefault();
			$searchModalOpen = true;
		}
		// / to open search (when not typing in an input)
		if (e.key === '/' && e.target instanceof HTMLElement && e.target.tagName !== 'INPUT') {
			e.preventDefault();
			$searchModalOpen = true;
		}
	}

	onMount(() => {
		window.addEventListener('keydown', handleKeydown);
		return () => {
			window.removeEventListener('keydown', handleKeydown);
		};
	});
</script>

{#snippet ListItem({ title, content, href, class: className, ...restProps }: ListItemProps)}
	<li>
		<NavigationMenu.Link>
			{#snippet child()}
				<a
					{href}
					class={cn(
						'block select-none space-y-1.5 rounded-md p-3 leading-none no-underline outline-none transition-colors hover:bg-accent hover:text-accent-foreground focus:bg-accent focus:text-accent-foreground',
						className
					)}
					{...restProps}
				>
					<div class="text-base font-medium leading-none">{title}</div>
					<!-- <p class="line-clamp-2 text-xs leading-relaxed text-foreground-muted">
						{content}
					</p> -->
				</a>
			{/snippet}
		</NavigationMenu.Link>
	</li>
{/snippet}

<nav class="flex w-full items-center justify-between gap-3 md:hidden">
	<a href="/" class="flex-shrink-0 font-serif text-lg text-foreground hover:opacity-80">
		silkwave
	</a>

	<!-- Right side: Search, Cart, Menu -->
	<div class="flex items-center gap-2">
		<!-- Search Button -->
		<button
			type="button"
			onclick={handleSearchbarClick}
			class="flex h-10 w-10 items-center justify-center rounded-lg text-foreground-muted transition-colors hover:bg-accent hover:text-foreground"
			aria-label="Search"
		>
			<Icon icon="search" class="h-5 w-5 fill-current" />
		</button>

		<!-- Dark Mode Toggle -->
		<DarkModeToggle />

		<!-- Cart -->
		<NavbarCart count={$cartCount} />

		<!-- Hamburger Menu -->
		<Dialog.Root bind:open={mobileMenuOpen}>
			<Dialog.Trigger
				class="flex h-10 w-10 items-center justify-center rounded-lg text-foreground-muted transition-colors hover:bg-accent hover:text-foreground"
			>
				<Icon icon="menu" class="h-6 w-6 fill-current" />
				<span class="sr-only">Open menu</span>
			</Dialog.Trigger>
			<Dialog.Content
				class="flex h-[100dvh] w-full max-w-full flex-col overflow-hidden rounded-none border-0 p-0 sm:max-w-full"
				showCloseButton={false}
			>
				<Dialog.Header
					class="flex-shrink-0 items-center justify-center border-b border-border px-3 py-4"
				>
					<Dialog.Title class="font-serif text-lg font-light">silkwave</Dialog.Title>
					<Dialog.Close
						class="absolute right-4 top-3 rounded-md p-2 transition-colors hover:bg-accent"
					>
						<Icon icon="cross" class="h-5 w-5 fill-current" />
						<span class="sr-only">Close</span>
					</Dialog.Close>
				</Dialog.Header>

				<div class="flex min-h-0 flex-1 flex-col overflow-y-auto">
					<!-- User Section -->
					<div class="flex flex-col">
						<a
							href="/settings"
							class="flex items-center gap-3 px-4 py-3 text-base text-foreground transition-colors hover:bg-accent"
						>
							<Icon icon="gear" class="h-4 w-4 fill-current" />
							Settings
						</a>
						{#if user}
							<button
								type="button"
								onclick={async () => {
									await POST('/logout', fetch);
									mobileMenuOpen = false;
									await invalidateAll();
									toast.success('You have been logged out.');
									goto('/');
								}}
								class="flex items-center gap-3 px-4 py-3 text-left text-base text-rose-300 transition-colors hover:bg-accent hover:text-rose-400"
							>
								<Icon icon="log-out" class="h-4 w-4 fill-current" />
								Logout
							</button>
						{:else}
							<a
								href="/login"
								class="flex items-center gap-3 px-4 py-3 text-base text-foreground transition-colors hover:bg-accent"
							>
								<Icon icon="log-in" class="h-4 w-4 fill-current" />
								Login
							</a>
						{/if}
					</div>

					<div class="mt-auto border-b border-t border-border">
						<button
							type="button"
							onclick={() => (genresExpanded = !genresExpanded)}
							class="flex w-full items-center justify-between px-4 py-3 text-left text-base font-medium text-foreground transition-colors hover:bg-accent"
						>
							Genres
							<Icon
								icon="chevron-down"
								class={cn(
									'h-4 w-4 fill-current transition-transform duration-200',
									genresExpanded && 'rotate-180'
								)}
							/>
						</button>
						{#if genresExpanded}
							<ul class="border-t border-border/50 bg-accent/30 pb-2">
								{#each genres as genre}
									<li>
										<a
											href={genre.href}
											class="block px-6 py-2.5 text-base text-foreground-muted transition-colors hover:bg-accent hover:text-foreground"
										>
											{genre.title}
										</a>
									</li>
								{/each}
							</ul>
						{/if}
					</div>

					<div class="border-t border-border p-4">
						<Button variant="primary" class="w-full" href="/upload">Upload Music</Button
						>
					</div>
				</div>
			</Dialog.Content>
		</Dialog.Root>
	</div>
</nav>

<!-- Desktop Navigation (hidden on mobile, shown on md+) -->
<div
	class="mx-auto hidden w-full max-w-7xl items-center justify-between gap-1 bg-background px-6 md:flex"
>
	<!-- Logo/Brand -->
	<a href="/" class="flex-shrink-0 font-serif text-lg text-foreground hover:opacity-80">
		silkwave
	</a>

	<!-- Genres Navigation Menu -->
	<NavigationMenu.Root viewport={false}>
		<NavigationMenu.List>
			<NavigationMenu.Item>
				<NavigationMenu.Trigger>Genres</NavigationMenu.Trigger>
				<NavigationMenu.Content>
					<ul class="grid w-[400px] gap-2 p-2 md:w-[500px] md:grid-cols-2 lg:w-[600px]">
						{#each genres as genre, i (i)}
							{@render ListItem({
								href: genre.href,
								title: genre.title,
								content: genre.description
							})}
						{/each}
					</ul>
				</NavigationMenu.Content>
			</NavigationMenu.Item>
		</NavigationMenu.List>
	</NavigationMenu.Root>

	<!-- Search Bar -->
	<button
		type="button"
		onclick={handleSearchbarClick}
		class="relative flex min-w-0 flex-1 items-center rounded-lg border border-border bg-background px-3 py-2 text-left text-base text-foreground-muted shadow-sm transition-colors hover:bg-accent/50 focus:outline-none focus:ring-2 focus:ring-ring sm:max-w-lg"
	>
		<Icon icon="search" class="mr-2 h-4 w-4 flex-shrink-0 fill-foreground-muted sm:mr-3" />
		<span class="truncate text-xs sm:text-base">Search genres, artists...</span>
		<div class="ml-auto hidden flex-shrink-0 items-center gap-1 sm:flex">
			<kbd
				class="pointer-events-none flex h-7 select-none items-center gap-1 rounded-sm border bg-muted px-1.5 font-mono text-base font-medium opacity-100"
			>
				<span class="text-xs">⌘</span>K
			</kbd>
		</div>
	</button>

	<!-- Right Side Actions -->
	<div class="flex flex-shrink-0 items-center gap-2">
		<Button variant="primary" size="sm" href="/upload">Upload Music</Button>
		<DarkModeToggle />
		<NavbarCart count={$cartCount} />

		<!-- Profile Navigation Menu -->
		<NavigationMenu.Root viewport={false}>
			<NavigationMenu.List>
				<NavigationMenu.Item>
					<NavigationMenu.Trigger class="p-2">
						<UserRound />
					</NavigationMenu.Trigger>
					<NavigationMenu.Content class="!end-0 !start-auto !w-48">
						<ul class="grid gap-1 p-2">
							{#if user}
								<li>
									<NavigationMenu.Link
										href="/settings"
										class="flex-row items-center gap-2"
									>
										<Icon icon="gear" class="h-4 w-4 fill-foreground" />
										Settings
									</NavigationMenu.Link>
								</li>
								<li class="group/logout">
									<NavigationMenu.Link
										href="#"
										class="flex-row items-center gap-2 text-rose-300 group-hover/logout:text-rose-400"
										onclick={async (e: MouseEvent) => {
											e.preventDefault();
											await POST('/logout', fetch);
											await invalidateAll();
											goto('/login');
										}}
									>
										<Icon
											icon="log-out"
											class="h-4 w-4 fill-rose-300 group-hover/logout:fill-rose-400"
										/>
										Logout
									</NavigationMenu.Link>
								</li>
							{:else}
								<li>
									<NavigationMenu.Link
										href="/login"
										class="flex-row items-center gap-2"
									>
										<Icon icon="log-in" class="h-4 w-4 fill-foreground" />
										Login
									</NavigationMenu.Link>
								</li>
							{/if}
						</ul>
					</NavigationMenu.Content>
				</NavigationMenu.Item>
			</NavigationMenu.List>
		</NavigationMenu.Root>
	</div>
</div>
