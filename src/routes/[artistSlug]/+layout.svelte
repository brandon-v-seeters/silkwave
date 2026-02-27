<script lang="ts">
	import { page } from '$app/stores';
	import { Button } from '$lib/components/ui/button/index';
	import { crossfade } from 'svelte/transition';
	import { cn } from '$lib/utils/utils.js';
	import { cubicInOut } from 'svelte/easing';
	import type { WithElementRef } from 'bits-ui';
	import type { HTMLInputAttributes } from 'svelte/elements';
	import { ARTIST_PROFILE_ROUTES } from '$lib/constants/routes.ts';
	import Icon from '$lib/components/atoms/Icon.svelte';

	let { children, class: className }: WithElementRef<HTMLInputAttributes> = $props();

	const [send, receive] = crossfade({
		duration: 250,
		easing: cubicInOut
	});

	const { artist } = $page.data.user ?? {};
	const isUserArtist = $derived(artist !== undefined && artist !== null);
</script>

<div class="relative" data-name="a">
	<div class="relative flex h-full max-h-96 w-full">
		{#if isUserArtist}
			<div
				class="group absolute right-0 top-0 flex h-full w-full cursor-pointer items-center justify-center transition-colors duration-300 hover:bg-black/20"
			>
				<Icon
					icon="edit-2"
					class="z-10 h-8 w-8 fill-white opacity-0 transition-opacity duration-300 group-hover:opacity-100"
				/>
			</div>
		{/if}
		<img src="/banner.jpg" alt="Banner" class="object-cover" />
	</div>
	<div class="grid grid-cols-1 px-4 md:gap-6 md:px-10 lg:grid-cols-6">
		<div class="relative col-span-1 flex flex-col gap-4 pt-[75px] md:col-span-2">
			<div class="absolute -top-20 flex flex-row">
				<div class="relative h-40 w-40 overflow-hidden rounded-lg">
					{#if isUserArtist}
						<div
							class="group absolute left-0 top-0 flex h-full w-full cursor-pointer items-center justify-center transition-colors duration-300 hover:bg-black/20"
						>
							<Icon
								icon="edit-2"
								class="z-10 h-8 w-8 fill-white opacity-0 transition-opacity duration-300 group-hover:opacity-100"
							/>
						</div>
					{/if}
					<img class="h-full w-full object-cover" src="/avatar.jpg" alt="Avatar" />
				</div>
				<div class="ml-4 mt-auto flex flex-row space-x-2 md:hidden">
					<Button variant="ghost">
						<Icon icon="comment-2-plus" class="h-5 w-5 fill-foreground" />
					</Button>
				</div>
			</div>
			<div class="mt-4">
				<h1 class="text-2xl font-semibold">{artist?.name}</h1>
				<p class="text-foreground-muted mt-1 font-medium leading-relaxed">
					{artist?.bio || 'Creating music production tutorials & resources'}
				</p>
				{#if !isUserArtist}
					<Button class="mt-4" variant="gradient">
						Join {artist?.name}
					</Button>
				{/if}
				<Button variant="ghost" class="mt-4 hidden md:block">
					<Icon icon="comment-2-plus" class="h-6 w-6 fill-foreground" />
				</Button>
			</div>
		</div>
		<div class="overflow-hidden pt-6 md:col-span-4">
			<nav
				class={cn(
					'no-scrollbar flex flex-row items-center space-x-1 overflow-x-auto lg:space-x-2 lg:space-y-1',
					className
				)}
			>
				{#each ARTIST_PROFILE_ROUTES as item}
					{@const href = `/${artist?.slug}${item.href}`}
					{@const isActive = $page.url.pathname === href}
					<Button
						{href}
						variant="ghost"
						class={cn(
							isActive ? 'bg-muted' : '',
							'relative justify-start hover:bg-muted'
						)}
						data-sveltekit-noscroll
					>
						{#if isActive}
							<div
								class="absolute inset-0 rounded-md bg-muted"
								in:send={{ key: 'active-sidebar-tab' }}
								out:receive={{ key: 'active-sidebar-tab' }}
							></div>
						{/if}
						<div class="relative">
							{item.title}
						</div>
					</Button>
				{/each}
			</nav>
			<div class="mt-6">
				{@render children?.()}
			</div>
		</div>
	</div>
</div>

<style>
	/* Hide scrollbar for Chrome, Safari and Opera */
	.no-scrollbar::-webkit-scrollbar {
		display: none;
	}

	/* Hide scrollbar for IE, Edge and Firefox */
	.no-scrollbar {
		-ms-overflow-style: none; /* IE and Edge */
		scrollbar-width: none; /* Firefox */
	}
</style>
