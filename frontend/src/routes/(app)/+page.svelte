<script lang="ts">
	import Icon from '$lib/components/atoms/Icon.svelte';
	import Button from '$lib/components/ui/button/button.svelte';
	import { resolve } from '$app/paths';

	const releaseCards = [
		{
			title: 'Vocal Studies and Upro...',
			artist: 'Prefuse 73',
			image: 'https://images.unsplash.com/photo-1493225457124-a3eb161ffa5f?w=500&h=500&fit=crop&q=80',
			tint: 'linear-gradient(135deg, #b9f5ff, #5b84e5)'
		},
		{
			title: 'Temples',
			artist: 'Lone',
			image: 'https://images.unsplash.com/photo-1470225620780-dba8ba36b745?w=500&h=500&fit=crop&q=80',
			tint: 'linear-gradient(135deg, #1c2d59, #8aa6ff)'
		},
		{
			title: 'Earth Tones',
			artist: 'Lenzman',
			image: 'https://images.unsplash.com/photo-1500530855697-b586d89ba3ee?w=500&h=500&fit=crop&q=80',
			tint: 'linear-gradient(135deg, #c8d8dd, #3f8d9e)'
		},
		{
			title: 'Kollections 06',
			artist: 'VA',
			image: 'https://images.unsplash.com/photo-1516280440614-37939bbacd81?w=500&h=500&fit=crop&q=80',
			tint: 'linear-gradient(135deg, #0c1118, #f4762f)'
		},
		{
			title: 'Dance It Away',
			artist: 'Donna',
			image: 'https://images.unsplash.com/photo-1511379938547-c1f69419868d?w=500&h=500&fit=crop&q=80',
			tint: 'linear-gradient(135deg, #171717, #6f6b65)'
		},
		{
			title: 'Moon Circuit',
			artist: 'Nujabes Study',
			image: 'https://images.unsplash.com/photo-1511671782779-c97d3d27a1d4?w=500&h=500&fit=crop&q=80',
			tint: 'linear-gradient(135deg, #111827, #f8fafc)'
		},
		{
			title: 'Late Checkout',
			artist: 'Mara TK',
			image: 'https://images.unsplash.com/photo-1508700115892-45ecd05ae2ad?w=500&h=500&fit=crop&q=80',
			tint: 'linear-gradient(135deg, #f2c4b3, #4338ca)'
		},
		{
			title: 'Soft Static',
			artist: 'Caldera',
			image: 'https://images.unsplash.com/photo-1500534314209-a25ddb2bd429?w=500&h=500&fit=crop&q=80',
			tint: 'linear-gradient(135deg, #2dd4bf, #111827)'
		}
	];

	const collectionStats = ['50,056 Likes', '213 Songs', '13 hr 7 min'];

	const discoverHref = resolve('/discover');
</script>

<svelte:head>
	<title>Silkwave - New Releases</title>
</svelte:head>

<div class="space-y-8">
	<div class="mt-4 flex items-center gap-4 overflow-x-auto">
		<a href={resolve('/')} class="shrink-0 font-semibold text-foreground">New Releases</a>
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
	<section
		class="relative isolate overflow-hidden rounded-[1.55rem] px-6 py-7 text-white shadow-[0_18px_50px_rgba(244,116,72,0.25)] sm:min-h-72 sm:px-8 sm:py-8"
	>
		<div
			class="absolute inset-y-0 right-0 -z-10 hidden w-[58%] bg-cover bg-center opacity-85 mix-blend-multiply md:block"
			style="background-image: linear-gradient(90deg, rgba(244,116,72,0), rgba(244,116,72,0.12)), url('/producer.jpg');"
		></div>
		<div
			class="absolute inset-0 -z-10 bg-[radial-gradient(circle_at_90%_50%,rgba(255,255,255,0.26),transparent_26%),linear-gradient(90deg,rgba(244,116,72,0.98),rgba(244,116,72,0.84)_44%,rgba(244,116,72,0.2))]"
		></div>

		<div class="flex h-full max-w-lg flex-col justify-between gap-10">
			<div>
				<p class="text-[0.66rem] font-bold uppercase tracking-[0.14em] text-white/82">
					Curated playlist
				</p>
				<h1 class="mt-8 max-w-sm text-4xl font-black tracking-tight text-white sm:text-5xl">
					R&B Hits
				</h1>
				<p class="mt-3 max-w-xs text-sm leading-relaxed text-white/82">
					Hot shot, confessions, deep cuts and studio-born grooves for late uploads.
				</p>
			</div>

			<div class="flex flex-wrap items-center gap-2 text-[0.72rem] text-white/82">
				<span class="flex h-6 w-6 items-center justify-center rounded-md bg-white/18">
					<Icon icon="heart" class="h-3.5 w-3.5 fill-current" />
				</span>
				{#each collectionStats as stat, index (stat)}
					<span>{stat}</span>
					{#if index < collectionStats.length - 1}
						<span class="text-white/45">•</span>
					{/if}
				{/each}
			</div>
		</div>

		<div class="mt-6 flex items-center gap-2 lg:absolute lg:bottom-6 lg:right-6 lg:mt-0">
			<Button variant="primary" class="rounded-full px-5">Play now</Button>
		</div>
	</section>

	<section>
		<div class="mb-4 flex items-end justify-between">
			<a
				href={discoverHref}
				class="text-[0.72rem] font-semibold text-muted-foreground transition hover:text-foreground"
			>
				View all
			</a>
		</div>

		<div class="grid grid-cols-2 gap-12 sm:grid-cols-3 xl:grid-cols-4">
			{#each releaseCards as release (release.title)}
				<a href={discoverHref} class="group min-w-0">
					<div
						class="relative aspect-square overflow-hidden rounded-4xl shadow-[0_16px_34px_rgba(15,23,42,0.1)] transition duration-300 group-hover:-translate-y-1 group-hover:shadow-[0_18px_42px_rgba(15,23,42,0.16)]"
						style="background: {release.tint};"
					>
						<img
							src={release.image}
							alt=""
							class="h-full w-full object-cover opacity-82 mix-blend-luminosity transition duration-500 group-hover:scale-105 group-hover:opacity-95"
						/>
						<div
							class="absolute inset-0 bg-gradient-to-t from-black/32 to-transparent opacity-0 transition group-hover:opacity-100"
						></div>
					</div>
					<div class="mt-4 min-w-0">
						<p class="truncate text-[0.78rem] font-semibold">{release.title}</p>
						<p class="mt-1 truncate text-[0.68rem] text-muted-foreground">
							{release.artist}
						</p>
					</div>
				</a>
			{/each}
		</div>
	</section>
</div>
