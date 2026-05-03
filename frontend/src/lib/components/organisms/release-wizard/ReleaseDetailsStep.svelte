<script lang="ts">
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import MultiSelect from '$lib/components/molecules/MultiSelect.svelte';
	import { getWizardContext } from './wizard.svelte';
	import Textarea from '$lib/components/ui/textarea/textarea.svelte';

	const wizard = getWizardContext();

	// Local state bound to inputs, synced with context
	let title = $state(wizard.releaseTitle);
	let selectedGenres = $state(wizard.genres);
	let description = $state(wizard.description);

	setTimeout(() => {
		title = 'Pulse of Desire';
		description =
			'With Pulse of Desire, the focus shifts to the hypnotic power of repetition. Each track strips away the excess of modern Drum & Bass to highlight a lone, driving motif. Expect cold, metallic soundscapes, precision-engineered sub-bass, and a dark, tech-driven atmosphere that stays locked in from start to finish.';
		selectedGenres = ['Drum & Bass'];
	}, 100);

	// Sync back to context on change
	$effect(() => {
		wizard.releaseTitle = title;
	});

	$effect(() => {
		wizard.genres = selectedGenres;
	});

	$effect(() => {
		wizard.description = description;
	});
</script>

<div class="flex flex-col gap-4 md:gap-5">
	<div class="flex flex-col gap-1">
		<Label for="release-title" class="text-sm font-medium">Release Title</Label>
		<Input id="release-title" bind:value={title} placeholder="Enter release title" />
	</div>

	<div class="flex flex-col gap-1">
		<Label class="text-sm font-medium">Genres</Label>
		<MultiSelect
			bind:selected={selectedGenres}
			options={wizard.suggestedGenres}
			placeholder="Select genres..."
			searchPlaceholder="Search genres..."
			allowCustom={true}
		/>
		<p class="text-xs text-foreground-muted">You can add multiple genres</p>
	</div>

	<div class="flex flex-col gap-1">
		<Label class="text-sm font-medium">Description</Label>
		<Textarea
			bind:value={description}
			placeholder="Enter release description"
			rows={3}
			class="resize-none md:resize-y"
		/>
	</div>
</div>
