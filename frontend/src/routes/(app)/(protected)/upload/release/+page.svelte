<script lang="ts">
	import { goto } from '$app/navigation';
	import Icon from '$lib/components/atoms/Icon.svelte';
	import { ProjectCoverArt, ProjectTrackList } from '$lib/components/organisms/release-editor';
	import { createWizardContext } from '$lib/components/organisms/release-wizard';
	import Button from '$lib/components/ui/button/button.svelte';

	const wizard = createWizardContext();

	const { data } = $props();

	if (data.draft) {
		wizard.initWizard(data.draft);
	}

	let totalDuration = $derived(() => {
		let totalSeconds = 0;
		for (const track of wizard.tracks) {
			const parts = track.duration?.split(':') ?? ['0', '00'];
			totalSeconds += parseInt(parts[0]) * 60 + parseInt(parts[1]);
		}
		const m = Math.floor(totalSeconds / 60);
		const s = totalSeconds % 60;
		return m > 0 ? `${m}m ${s}s` : `${s}s`;
	});
</script>

<div class="relative sm:block sm:pt-8">
	<!-- Header actions -->
	<div class="flex items-center justify-between px-6 py-4 sm:px-16 sm:py-6">
		<Button
			variant="secondary"
			size="icon"
			class="h-11 w-11 shrink-0 rounded-2xl"
			onclick={() => goto('/upload')}
		>
			<Icon icon="chevron-left" variant="line" class="h-4 w-4 fill-foreground" />
		</Button>

		<div class="flex gap-2">
			<Button
				variant="secondary"
				class="rounded-2xl"
				onclick={() => wizard.saveAsDraft()}
				disabled={wizard.isUploading}
			>
				Save Draft
			</Button>
			<Button
				variant="primary"
				class="rounded-2xl"
				onclick={() => wizard.publishRelease()}
				disabled={!wizard.isFormValid || wizard.isUploading}
			>
				Publish
			</Button>
		</div>
	</div>

	<!-- Main content: two-column layout -->
	<div class="flex flex-col items-center gap-4 px-6 sm:gap-12 sm:px-16 sm:pt-0">
		<!-- Left column: Cover art (sticky on desktop) -->
		<div class="flex w-full flex-col items-center gap-6 sm:w-[320px] lg:w-[405px]">
			<ProjectCoverArt
				bind:preview={wizard.coverArtPreview}
				onchange={(file) => wizard.setCoverArt(file)}
			/>
		</div>

		<!-- Right column: Title + Tracks -->
		<div class="w-full pb-24 md:max-w-xl">
			<div class="flex flex-col gap-4">
				<!-- Title + meta + preview -->
				<div class="flex items-start gap-3">
					<div class="flex min-w-0 flex-1 flex-col">
						<input
							type="text"
							placeholder="untitled project"
							bind:value={wizard.releaseTitle}
							class="w-full bg-transparent font-serif text-3xl font-medium tracking-tight text-foreground placeholder:text-foreground-muted/40 focus:outline-none sm:text-4xl"
						/>
						<div class="mt-1 flex items-center gap-1.5 text-sm text-foreground-muted">
							{#if wizard.tracks.length > 0}
								<span>
									{wizard.tracks.length}
									{wizard.tracks.length === 1 ? 'track' : 'tracks'}
								</span>
								<span aria-hidden="true">·</span>
								<span>{totalDuration()}</span>
							{:else}
								<span>No tracks yet</span>
							{/if}
						</div>
					</div>
					<Button
						size="icon"
						class="h-11 w-11 shrink-0 rounded-2xl"
						disabled={wizard.tracks.length === 0}
						aria-label="Preview release"
					>
						<Icon
							icon="play"
							variant="filled"
							class="h-4 w-4 fill-primary-foreground"
						/>
					</Button>
				</div>

				<!-- Track list with add button -->
				<ProjectTrackList
					tracks={wizard.tracks}
					onfiles={(files) => wizard.addBulkTracks(files)}
					onTitleChange={(id, title) => wizard.updateTrackTitle(id, title)}
					onRemove={(id) => wizard.removeTrack(id)}
					onMove={(index, direction) => wizard.moveTrack(index, direction)}
				/>
			</div>
		</div>
	</div>
</div>
