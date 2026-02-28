<script lang="ts">
	let { release } = $props();
</script>

<div class="flex flex-col gap-2">
	<a
		href="/{release.artist?.slug}/{release.slug}"
		class="group/image flex aspect-square cursor-pointer flex-col gap-2 overflow-hidden rounded-xl"
	>
		<img
			src={release.coverArt}
			alt={release.title}
			class="mx-auto h-full w-full rounded-xl object-cover transition-transform duration-300 ease-in-out group-hover/image:scale-[1.1]"
		/>
	</a>
	<div class="flex flex-col">
		<h5 class="cursor-pointer truncate text-lg font-medium">{release.title}</h5>
		<a
			href="/{release.artist?.slug}"
			class="group/artist cursor-pointer truncate text-sm text-foreground-muted transition-colors duration-300 hover:text-primary"
		>
			{release.artist?.name}
		</a>
	</div>
</div>
