<script lang="ts">
	import { resolve } from '$app/paths';
	import {
		artistRoute,
		artistRouteParams,
		hasReleaseRoute,
		releaseRoute,
		releaseRouteParams,
		type CatalogRelease
	} from '../releases';

	let { release }: { release: CatalogRelease } = $props();
	const coverArt = $derived(release.coverArtUrl);
</script>

<div class="flex flex-col gap-2">
	{#if hasReleaseRoute(release)}
		<a
			href={resolve(releaseRoute, releaseRouteParams(release))}
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
	{:else}
		<div class="flex aspect-square flex-col gap-2 overflow-hidden rounded-xl">
			{#if coverArt}
				<img
					src={coverArt}
					alt={release.title}
					class="mx-auto h-full w-full rounded-xl object-cover"
				/>
			{:else}
				<div class="flex h-full w-full items-center justify-center rounded-xl bg-muted text-foreground-muted">
					{release.title}
				</div>
			{/if}
		</div>
	{/if}
	<div class="flex flex-col">
		<h5 class="cursor-pointer truncate text-lg font-medium">{release.title}</h5>
		{#if release.artist?.slug}
			<a
				href={resolve(artistRoute, artistRouteParams(release.artist))}
				class="group/artist cursor-pointer truncate text-sm text-foreground-muted transition-colors duration-300 hover:text-primary"
			>
				{release.artist?.name}
			</a>
		{:else}
			<span class="truncate text-sm text-foreground-muted">{release.artist?.name}</span>
		{/if}
	</div>
</div>
