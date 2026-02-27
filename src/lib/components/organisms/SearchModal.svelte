<script lang="ts">
	import * as Dialog from '$lib/components/ui/dialog/index.ts';
	import { goto } from '$app/navigation';
	import Icon from '../atoms/Icon.svelte';
	import Input from '../ui/input/input.svelte';
	import { cn } from '$lib/utils/utils';

	let {
		open = $bindable(false),
		searchQuery = $bindable('')
	}: {
		open?: boolean;
		searchQuery?: string;
	} = $props();

	let selectedCategory = $state<'all' | 'genres' | 'artists' | 'releases' | 'albums'>('all');
	let imageErrors = $state<Set<string>>(new Set());

	// Mock data - replace with actual API calls
	const popularGenres = ['Electronic', 'Hip Hop', 'Jazz', 'Rock', 'Ambient', 'House'];
	const trendingArtists = [
		{
			name: 'KOAN Sound',
			avatar: 'https://images.unsplash.com/photo-1493225457124-a3eb161ffa5f?w=100&h=100&fit=crop&q=80'
		},
		{
			name: 'Neon Dreams',
			avatar: 'https://images.unsplash.com/photo-1459749411175-04bf5292ceea?w=100&h=100&fit=crop&q=80'
		},
		{ name: 'Synth Wave', avatar: null }, // No avatar - will use icon
		{
			name: 'City Vibes',
			avatar: 'https://images.unsplash.com/photo-1514320291840-2e0a9bf2a9ae?w=100&h=100&fit=crop&q=80'
		}
	];
	const featuredReleases = [
		{
			title: 'Led by Ancient Light',
			artist: 'KOAN Sound',
			type: 'Album',
			coverArt:
				'https://images.unsplash.com/photo-1493225457124-a3eb161ffa5f?w=200&h=200&fit=crop&q=80'
		},
		{
			title: 'Midnight Echoes',
			artist: 'Neon Dreams',
			type: 'EP',
			coverArt:
				'https://images.unsplash.com/photo-1459749411175-04bf5292ceea?w=200&h=200&fit=crop&q=80'
		},
		{
			title: 'Electric Dreams',
			artist: 'Synth Wave',
			type: 'Album',
			coverArt: null // No cover art - will use icon
		}
	];

	function handleSearch() {
		if (searchQuery.trim()) {
			goto(`/discover?q=${encodeURIComponent(searchQuery.trim())}`);
			open = false;
		}
	}

	function handleKeydown(e: KeyboardEvent) {
		if (e.key === 'Enter') {
			handleSearch();
		} else if (e.key === 'Escape') {
			open = false;
		}
	}

	function selectCategory(category: typeof selectedCategory) {
		selectedCategory = category;
	}

	function selectRecentSearch(search: string) {
		searchQuery = search;
		handleSearch();
	}

	function handleImageError(src: string) {
		imageErrors.add(src);
		imageErrors = imageErrors; // Trigger reactivity
	}

	$effect(() => {
		if (open) {
			// Focus will be handled by autofocus attribute
		}
	});
</script>

<Dialog.Root bind:open>
	<Dialog.Portal>
		<Dialog.Overlay />
		<Dialog.Content>
			<div class="flex flex-row">
				<!-- Categories -->
				<div class="border-r border-border px-2 py-4">
					<div class="flex flex-col gap-1">
						<button
							type="button"
							onclick={() => selectCategory('all')}
							class={cn(
								'flex w-full items-center gap-2 rounded-md px-3 py-2 text-left text-base font-medium transition-colors',
								selectedCategory === 'all'
									? 'bg-foreground/5 fill-foreground text-foreground'
									: 'text-foreground-muted hover:bg-foreground/5 hover:text-foreground'
							)}
						>
							<Icon icon="home" class="h-5 w-5 fill-current" />
							<span>All</span>
						</button>
						<button
							type="button"
							onclick={() => selectCategory('genres')}
							class={cn(
								'flex w-full items-center gap-2 rounded-md px-3 py-2 text-left text-base font-medium transition-colors',
								selectedCategory === 'genres'
									? 'bg-foreground/5 fill-foreground text-foreground'
									: 'text-foreground-muted hover:bg-foreground/5 hover:text-foreground'
							)}
						>
							<Icon icon="folder" class="h-6 w-6 fill-current" />
							<span>Genres</span>
						</button>
						<button
							type="button"
							onclick={() => selectCategory('artists')}
							class={cn(
								'flex w-full items-center gap-2 rounded-md px-3 py-2 text-left text-base font-medium transition-colors',
								selectedCategory === 'artists'
									? 'bg-foreground/5 fill-foreground text-foreground'
									: 'text-foreground-muted hover:bg-foreground/5 hover:text-foreground'
							)}
						>
							<Icon icon="user" class="h-5 w-5 fill-current" />
							<span>Artists</span>
						</button>
						<button
							type="button"
							onclick={() => selectCategory('releases')}
							class={cn(
								'flex w-full items-center gap-2 rounded-md px-3 py-2 text-left text-base font-medium transition-colors',
								selectedCategory === 'releases'
									? 'bg-foreground/5 fill-foreground text-foreground'
									: 'text-foreground-muted hover:bg-foreground/5 hover:text-foreground'
							)}
						>
							<Icon icon="music-note-3" class="h-5 w-5 fill-current" />
							<span>Releases</span>
						</button>
						<button
							type="button"
							onclick={() => selectCategory('albums')}
							class={cn(
								'flex w-full items-center gap-2 rounded-md px-3 py-2 text-left text-base font-medium transition-colors',
								selectedCategory === 'albums'
									? 'bg-foreground/5 fill-foreground text-foreground'
									: 'text-foreground-muted hover:bg-foreground/5 hover:text-foreground'
							)}
						>
							<Icon icon="music-note-2" class="h-5 w-5 fill-current" />
							<span>Albums</span>
						</button>
					</div>
				</div>

				<div class="flex flex-col">
					<!-- Search Input -->
					<div class="border-b border-border p-4">
						<div class="relative">
							<Icon
								icon="search"
								class="absolute left-3 top-1/2 h-5 w-5 -translate-y-1/2 fill-foreground-muted"
							/>
							<Input
								bind:value={searchQuery}
								type="search"
								placeholder="Search genres, artists, releases..."
								class="w-full pl-10 pr-4 text-lg"
								onkeydown={handleKeydown}
								autofocus
							/>
						</div>
					</div>

					<!-- Content Area -->
					<div class="max-h-[60vh] overflow-y-auto p-4">
						{#if !searchQuery.trim()}
							<!-- Popular Genres -->
							<div class="mb-6">
								<h3 class="mb-3 text-base font-medium text-foreground-muted">
									Popular Genres
								</h3>
								<div class="flex flex-wrap gap-2">
									{#each popularGenres as genre}
										<button
											type="button"
											onclick={() => {
												searchQuery = genre;
												selectedCategory = 'genres';
											}}
											class="rounded-lg border border-border bg-background px-4 py-2 text-base transition-colors hover:bg-accent"
										>
											{genre}
										</button>
									{/each}
								</div>
							</div>

							<!-- Trending Artists -->
							<div class="mb-6">
								<h3 class="mb-3 text-base font-medium text-foreground-muted">
									Trending Artists
								</h3>
								<div class="space-y-2">
									{#each trendingArtists as artist}
										<button
											type="button"
											onclick={() => {
												searchQuery = artist.name;
												selectedCategory = 'artists';
											}}
											class="flex w-full items-center gap-3 rounded-lg border border-border bg-background p-3 text-left transition-colors hover:bg-accent"
										>
											{#if artist.avatar && !imageErrors.has(artist.avatar)}
												<img
													src={artist.avatar}
													alt={artist.name}
													class="h-10 w-10 rounded-full object-cover"
													onerror={() => handleImageError(artist.avatar!)}
												/>
											{:else}
												<div
													class="flex h-10 w-10 items-center justify-center rounded-full bg-primary/10"
												>
													<Icon
														icon="user"
														class="h-5 w-5 fill-primary"
													/>
												</div>
											{/if}
											<span class="font-medium">{artist.name}</span>
										</button>
									{/each}
								</div>
							</div>

							<!-- Featured Releases -->
							<div>
								<h3 class="mb-3 text-base font-medium text-foreground-muted">
									Featured Releases
								</h3>
								<div class="space-y-2">
									{#each featuredReleases as release}
										<button
											type="button"
											onclick={() => {
												searchQuery = release.title;
												selectedCategory = 'releases';
											}}
											class="flex w-full items-center gap-3 rounded-lg border border-border bg-background p-3 text-left transition-colors hover:bg-accent"
										>
											{#if release.coverArt && !imageErrors.has(release.coverArt)}
												<img
													src={release.coverArt}
													alt={release.title}
													class="h-12 w-12 rounded-lg object-cover"
													onerror={() =>
														handleImageError(release.coverArt!)}
												/>
											{:else}
												<div
													class="flex h-12 w-12 items-center justify-center rounded-lg bg-primary/10"
												>
													<Icon
														icon="file"
														class="h-6 w-6 fill-primary"
													/>
												</div>
											{/if}
											<div class="flex-1">
												<div class="font-medium">{release.title}</div>
												<div class="text-base text-foreground-muted">
													{release.artist} • {release.type}
												</div>
											</div>
										</button>
									{/each}
								</div>
							</div>
						{:else}
							<!-- Search Results -->
							<div class="space-y-4">
								<!-- Results would go here based on searchQuery and selectedCategory -->
								<div class="text-center text-base text-foreground-muted">
									Searching for "{searchQuery}" in {selectedCategory}...
								</div>
								<!-- TODO: Add actual search results here -->
							</div>
						{/if}
					</div>
				</div>
			</div>

			<!-- Footer -->
			<div class="border-t border-border px-4 py-3">
				<div class="flex items-center justify-between text-xs text-foreground-muted">
					<div class="flex items-center gap-4">
						<span class="flex items-center gap-1">
							<kbd
								class="bg-muted pointer-events-none inline-flex h-5 select-none items-center gap-1 rounded border px-1.5 font-mono text-[10px] font-medium opacity-100"
							>
								<span class="text-xs">⌘</span>K
							</kbd>
							<span>to open</span>
						</span>
						<span class="flex items-center gap-1">
							<kbd
								class="bg-muted pointer-events-none inline-flex h-5 select-none items-center gap-1 rounded border px-1.5 font-mono text-[10px] font-medium opacity-100"
							>
								↵
							</kbd>
							<span>to select</span>
						</span>
					</div>
					<Dialog.Close
						class="text-foreground-muted transition-colors hover:text-foreground"
					>
						<kbd
							class="bg-muted pointer-events-none inline-flex h-5 select-none items-center gap-1 rounded border px-1.5 font-mono text-[10px] font-medium opacity-100"
						>
							ESC
						</kbd>
						to close
					</Dialog.Close>
				</div>
			</div>
		</Dialog.Content>
	</Dialog.Portal>
</Dialog.Root>
