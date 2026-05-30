<script lang="ts">
	import type { Artist, Release } from '$lib/types/generated/models';

	type CardRelease = Release & {
		artist?: Pick<Artist, 'name' | 'slug'>;
		coverArt?: string | null;
	};

	let { release }: { release: CardRelease } = $props();

	const releaseHref = $derived(
		release.artist?.slug ? `/artist/${release.artist.slug}/releases/${release.slug}` : '#'
	);
	const artistHref = $derived(release.artist?.slug ? `/artist/${release.artist.slug}` : '#');
	const coverArt = $derived(
		release.coverArt ||
			release.cover ||
			release.assets?.coverArt?.medium ||
			release.assets?.coverArt?.original ||
			release.assets?.coverArt?.thumbnail ||
			null
	);
</script>

<div class="flex flex-col gap-2">
	<a
		href={releaseHref}
		class="group/image flex aspect-square cursor-pointer flex-col gap-2 overflow-hidden rounded-xl"
	>
		{#if coverArt}
			<img
				src={coverArt}
				alt={release.title}
				class="mx-auto h-full w-full rounded-xl object-cover transition-transform duration-300 ease-in-out group-hover/image:scale-[1.1]"
			/>
		{:else}
			<div class="flex h-full w-full items-center justify-center rounded-xl bg-muted text-foreground-muted">
				{release.title}
			</div>
		{/if}
	</a>
	<div class="flex flex-col">
		<h5 class="cursor-pointer truncate text-lg font-medium">{release.title}</h5>
		<a
			href={artistHref}
			class="group/artist cursor-pointer truncate text-sm text-foreground-muted transition-colors duration-300 hover:text-primary"
		>
			{release.artist?.name}
		</a>
	</div>
</div>
