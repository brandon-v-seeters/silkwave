<script lang="ts">
	// Vendor
	import { getContext } from 'svelte';
	import { superForm } from 'sveltekit-superforms';

	// Components
	import Icon from '$lib/components/atoms/Icon.svelte';
	import Button from '$lib/components/ui/button/button.svelte';
	import DataTable from '$lib/components/ui/data-table/data-table.svelte';
	import { columns } from './drafts/columns';

	// Stores
	import * as Form from '$lib/components/ui/form';
	import Input from '$lib/components/ui/input/input.svelte';

	// Types
	import type { AppUser } from '$lib/types/generated';
	import { page } from '$app/state';
	import FormErrors from '$lib/components/atoms/form/FormErrors.svelte';

	const userCtx = getContext<{ current: AppUser }>('user');
	const user = $derived(userCtx.current);
	const drafts = $derived(page.data.drafts ?? []);

	const form = superForm(page.data.form);
	const { form: formData, enhance, delayed } = form;
</script>

{#if !user?.artist?._key}
	<div class="mx-auto flex max-w-3xl flex-col gap-2">
		<h1 class="leading-xl text-xl md:text-left md:text-2xl">First things first...</h1>
		<p class="mb-4 text-base font-light text-foreground-muted md:text-left">
			Before you can upload your music, you need to set up your artist name.
		</p>
		<form method="POST" use:enhance class="flex w-full max-w-sm flex-col gap-4">
			<Form.Field {form} name="name">
				<Form.Control>
					{#snippet children({ props })}
						<Form.Label>Artist Name</Form.Label>
						<Input
							{...props}
							placeholder="Your artist name"
							autocomplete="off"
							autocapitalize="words"
							bind:value={$formData.name}
							disabled={$delayed}
						/>
					{/snippet}
				</Form.Control>
				<Form.FieldErrors />
			</Form.Field>
			<FormErrors {form} />
			<Button type="submit" variant="gradient" disabled={$delayed} class="w-fit">
				{#if $delayed}
					<Icon icon="loader-2" class="mr-2 h-4 w-4 animate-spin" />
				{/if}
				Set my Artist Name
			</Button>
		</form>
	</div>
{:else}
	<div class="mx-auto flex max-w-5xl flex-col gap-8">
		<div class="flex flex-col gap-4 md:flex-row md:items-end md:justify-between">
			<div class="flex flex-col gap-2">
				<h1 class="leading-xl text-xl md:text-left md:text-2xl">Upload music</h1>
				<p class="max-w-2xl text-base font-light text-foreground-muted md:text-left">
					Create a release from scratch or pick up a draft.
				</p>
			</div>
			<Button
				variant="gradient"
				href="/upload/release"
				class="w-full justify-center md:w-fit"
			>
				<Icon icon="music-note-2" class="h-5 w-5" />
				Create a new release
			</Button>
		</div>

		<section class="flex flex-col gap-3" aria-labelledby="drafts-heading">
			<div class="flex flex-col gap-1">
				<h2 id="drafts-heading" class="text-lg font-medium text-foreground">Drafts</h2>
				<p class="text-sm text-foreground-muted">
					{drafts.length === 1 ? '1 saved draft' : `${drafts.length} saved drafts`}
				</p>
			</div>

			<DataTable data={drafts} {columns}>
				{#snippet noResults()}
					<div class="flex flex-col items-center justify-center gap-3 p-6">
						<Icon icon="file-text" variant="line" class="h-6 w-6 fill-foreground" />
						<p class="text-sm text-foreground-muted">No drafts yet.</p>
						<Button variant="gradient" href="/upload/release">
							<Icon icon="plus" variant="line" class="h-5 w-5 fill-foreground" />
							Create a new release
						</Button>
					</div>
				{/snippet}
			</DataTable>
		</section>
	</div>
{/if}
