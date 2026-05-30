<script lang="ts">
	import { goto } from '$app/navigation';
	import * as Dialog from '$lib/components/ui/dialog/index';
	import type { IconKey } from '$lib/types/Icon';
	import type { Artist, Release, ReleaseWithArtist } from '$lib/types/generated/models';
	import { cn } from '$lib/utils/utils';
	import Icon from '../atoms/Icon.svelte';
	import Input from '../ui/input/input.svelte';

	type SearchCategory = 'all' | 'genres' | 'artists' | 'releases' | 'albums';

	type ReleaseResult = Release & {
		artist?: Pick<Artist, 'name' | 'slug'>;
		coverArt?: string | null;
		type?: string;
	};

	type ApiEnvelope<T> = {
		data?: T;
		error?: { message?: string } | string;
	};

	type SearchReleaseRow = ReleaseResult | ReleaseWithArtist;
	type SearchReleasesPayload = {
		releases?: SearchReleaseRow[];
	};

	let {
		open = $bindable(false),
		searchQuery = $bindable('')
	}: {
		open?: boolean;
		searchQuery?: string;
	} = $props();

	const categories: { value: SearchCategory; label: string; icon: IconKey }[] = [
		{ value: 'all', label: 'All', icon: 'home' },
		{ value: 'genres', label: 'Genres', icon: 'folder' },
		{ value: 'artists', label: 'Artists', icon: 'user' },
		{ value: 'releases', label: 'Releases', icon: 'music-note-3' },
		{ value: 'albums', label: 'Albums', icon: 'music-note-2' }
	];

	let selectedCategory = $state<SearchCategory>('all');
	let releases = $state.raw<ReleaseResult[]>([]);
	let isLoading = $state(false);
	let error = $state<string | null>(null);
	let imageErrors = $state<string[]>([]);

	let trimmedQuery = $derived(searchQuery.trim());
	let canSearchReleases = $derived(
		selectedCategory === 'all' || selectedCategory === 'releases'
	);
	let shouldSearch = $derived(open && trimmedQuery.length > 0 && canSearchReleases);

	function handleSearch() {
		if (trimmedQuery) {
			goto(`/discover?q=${encodeURIComponent(trimmedQuery)}`);
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

	function selectCategory(category: SearchCategory) {
		selectedCategory = category;
		error = null;
	}

	function handleImageError(src: string) {
		if (!imageErrors.includes(src)) {
			imageErrors = [...imageErrors, src];
		}
	}

	function coverArtFor(release: ReleaseResult) {
		return (
			release.coverArt ||
			release.cover ||
			release.assets?.coverArt?.medium ||
			release.assets?.coverArt?.original ||
			release.assets?.coverArt?.thumbnail ||
			null
		);
	}

	function releaseHref(release: ReleaseResult) {
		if (!release.artist?.slug || !release.slug) return null;

		return `/artist/${release.artist.slug}/releases/${release.slug}`;
	}

	function releaseKind(release: ReleaseResult) {
		return release.type ?? release.releaseType ?? 'Release';
	}

	function resultKey(release: ReleaseResult) {
		return release.id || release._key || `${release.artist?.slug ?? 'artist'}-${release.slug}`;
	}

	function parseError(errorValue: ApiEnvelope<SearchReleasesPayload>['error']) {
		if (!errorValue) return 'Failed to search releases';
		if (typeof errorValue === 'string') return errorValue;

		return errorValue.message ?? 'Failed to search releases';
	}

	function isGeneratedReleaseWithArtist(row: SearchReleaseRow): row is ReleaseWithArtist {
		return 'Release' in row;
	}

	function normalizeReleaseResult(row: SearchReleaseRow): ReleaseResult {
		if (isGeneratedReleaseWithArtist(row)) {
			return {
				...row.Release,
				artist: row.artist
			};
		}

		return row;
	}

	async function loadReleaseResults(query: string, signal: AbortSignal) {
		isLoading = true;
		error = null;

		try {
			const response = await fetch(
				`/api/releases?q=${encodeURIComponent(query)}&limit=8`,
				{ signal }
			);
			const payload = (await response.json()) as ApiEnvelope<SearchReleasesPayload> &
				SearchReleasesPayload;

			if (!response.ok || payload.error) {
				error = parseError(payload.error);
				releases = [];
				return;
			}

			releases = (payload.data?.releases ?? payload.releases ?? []).map(
				normalizeReleaseResult
			);
		} catch (e) {
			if (e instanceof DOMException && e.name === 'AbortError') return;

			error = 'Failed to search releases';
			releases = [];
		} finally {
			if (!signal.aborted) {
				isLoading = false;
			}
		}
	}

	$effect(() => {
		if (!shouldSearch) {
			releases = [];
			isLoading = false;
			error = null;
			return;
		}

		const controller = new AbortController();
		const timeout = window.setTimeout(() => {
			loadReleaseResults(trimmedQuery, controller.signal);
		}, 180);

		return () => {
			window.clearTimeout(timeout);
			controller.abort();
		};
	});
</script>

<Dialog.Root bind:open>
	<Dialog.Portal>
		<Dialog.Overlay />
		<Dialog.Content>
			<div class="flex min-h-[28rem] flex-row">
				<div class="border-r border-border px-2 py-4">
					<div class="flex flex-col gap-1">
						{#each categories as category (category.value)}
							<button
								type="button"
								onclick={() => selectCategory(category.value)}
								class={cn(
									'flex w-full items-center gap-2 rounded-md px-3 py-2 text-left text-base font-medium transition-colors',
									selectedCategory === category.value
										? 'bg-foreground/5 fill-foreground text-foreground'
										: 'text-foreground-muted hover:bg-foreground/5 hover:text-foreground'
								)}
							>
								<Icon icon={category.icon} class="h-5 w-5 fill-current" />
								<span>{category.label}</span>
							</button>
						{/each}
					</div>
				</div>

				<div class="flex min-w-0 flex-1 flex-col">
					<div class="border-b border-border p-4">
						<div class="relative">
							<Icon
								icon="search"
								class="absolute left-3 top-1/2 h-5 w-5 -translate-y-1/2 fill-foreground-muted"
							/>
							<Input
								bind:value={searchQuery}
								type="search"
								placeholder="Search releases..."
								class="w-full pl-10 pr-4 text-lg"
								onkeydown={handleKeydown}
								autofocus
							/>
						</div>
					</div>

					<div class="max-h-[60vh] min-h-80 overflow-y-auto p-4">
						{#if !trimmedQuery}
							<div class="flex h-full min-h-64 flex-col items-center justify-center text-center">
								<div class="mb-4 flex h-14 w-14 items-center justify-center rounded-lg bg-primary/10">
									<Icon icon="music-note-2" class="h-7 w-7 fill-primary" />
								</div>
								<h3 class="text-lg font-medium text-foreground">Search releases</h3>
								<p class="mt-2 max-w-sm text-base text-foreground-muted">
									Search by release title or description.
								</p>
							</div>
						{:else if !canSearchReleases}
							<div class="flex h-full min-h-64 flex-col items-center justify-center text-center">
								<div class="mb-4 flex h-14 w-14 items-center justify-center rounded-lg bg-foreground/5">
									<Icon icon="search" class="h-7 w-7 fill-foreground-muted" />
								</div>
								<h3 class="text-lg font-medium text-foreground">
									{categories.find((category) => category.value === selectedCategory)?.label}
									search is not available yet.
								</h3>
								<p class="mt-2 max-w-sm text-base text-foreground-muted">
									Try Releases for now.
								</p>
							</div>
						{:else if isLoading}
							<div class="flex min-h-64 items-center justify-center">
								<Icon icon="loader-2" class="h-8 w-8 animate-spin fill-foreground-muted" />
							</div>
						{:else if error}
							<div class="rounded-lg bg-rose-500/10 p-6 text-center text-rose-400">
								<p>{error}</p>
							</div>
						{:else if releases.length === 0}
							<div class="flex min-h-64 flex-col items-center justify-center text-center">
								<div class="mb-4 flex h-14 w-14 items-center justify-center rounded-lg bg-foreground/5">
									<Icon icon="music-note-2" class="h-7 w-7 fill-foreground-muted" />
								</div>
								<h3 class="text-lg font-medium text-foreground">No releases found</h3>
								<p class="mt-2 max-w-sm text-base text-foreground-muted">
									No published Releases matched "{trimmedQuery}".
								</p>
							</div>
						{:else}
							<div class="space-y-2">
								{#each releases as release (resultKey(release))}
									{@const coverArt = coverArtFor(release)}
									{@const href = releaseHref(release)}
									{#if href}
										<a
											href={href}
											onclick={() => (open = false)}
											class="group flex w-full items-center gap-3 rounded-lg border border-border bg-background p-3 text-left transition-colors hover:bg-accent"
										>
											<div class="h-12 w-12 shrink-0 overflow-hidden rounded-md bg-primary/10">
												{#if coverArt && !imageErrors.includes(coverArt)}
													<img
														src={coverArt}
														alt={release.title}
														class="h-full w-full object-cover"
														onerror={() => handleImageError(coverArt)}
													/>
												{:else}
													<div class="flex h-full w-full items-center justify-center">
														<Icon icon="music-note-2" class="h-6 w-6 fill-primary" />
													</div>
												{/if}
											</div>

											<div class="min-w-0 flex-1">
												<div class="truncate font-medium group-hover:text-primary">
													{release.title}
												</div>
												<div class="mt-1 truncate text-base text-foreground-muted">
													{release.artist?.name ?? 'Unknown Artist'} · {releaseKind(release)}
												</div>
											</div>
										</a>
									{/if}
								{/each}
							</div>
						{/if}
					</div>
				</div>
			</div>

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
							<span>to search</span>
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
