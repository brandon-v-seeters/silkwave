<script lang="ts">
	// Vendor
	import { enhance } from '$app/forms';
	import { getContext } from 'svelte';
	import { superForm } from 'sveltekit-superforms';

	// Components
	import Icon from '$lib/components/atoms/Icon.svelte';

	// Stores
	import * as Form from '$lib/components/ui/form/index.ts';
	import Input from '$lib/components/ui/input/input.svelte';

	// Types
	import type { AppUser } from '$lib/types/generated';
	import { page } from '$app/state';
	import Button from '$lib/components/ui/button/button.svelte';
	import FormErrors from '$lib/components/atoms/form/FormErrors.svelte';

	const userCtx = getContext<{ current: AppUser }>('user');
	const user = $derived(userCtx.current);

	const form = superForm(page.data.form);
	const { form: formData } = form;

	let isLoading = $state(false);
</script>

{#snippet btn(href: string, title: string, icon: string)}
	<a
		{href}
		class="group col-span-1 flex h-40 cursor-pointer flex-col
		items-center justify-center gap-4 rounded-lg border border-zinc-700 bg-muted
		bg-zinc-900 p-4 transition-all duration-500 hover:border-primary hover:bg-primary/10 hover:text-primary"
	>
		<Icon
			{icon}
			class="transition-fill h-8 w-8 fill-foreground duration-500 group-hover:fill-primary"
		/>
		<h2 class="text-base">{title}</h2>
	</a>
{/snippet}

{#if !user?.artist?._key}
	<div class="mx-auto flex max-w-3xl flex-col gap-2">
		<h1 class="leading-xl text-xl md:text-left md:text-2xl">First things first...</h1>
		<p class="mb-4 text-base font-light text-foreground-muted md:text-left">
			Before you can upload your music, you need to set up your artist name.
		</p>
		<form method="POST" use:enhance class="flex w-full max-w-sm flex-col gap-4">
			<Form.Field {form} name="name">
				<Form.Control let:attrs>
					<Form.Label for="name">Artist Name</Form.Label>
					<Input {...attrs} id="name" name="name" bind:value={$formData.name} />
					<Form.FieldErrors />
				</Form.Control>
			</Form.Field>
			<FormErrors {form} />
			<Button type="submit" variant="primary" disabled={isLoading} class="w-fit">
				{#if isLoading}
					<Icon icon="loader-2" class="mr-2 h-4 w-4 animate-spin" />
				{/if}
				Set my Artist Name
			</Button>
		</form>
	</div>
{:else}
	<div class="mx-auto flex max-w-3xl flex-col gap-2">
		<h1 class="leading-xl text-xl md:text-left md:text-2xl">Manage your music</h1>
		<p class="mb-4 text-base font-light text-foreground-muted md:text-left">
			Start making money from your music today
		</p>
		<div class="gap-base mt-sm grid max-w-[280px] grid-cols-1 md:gap-6">
			<Button variant="primary" href="/upload/release">
				<Icon icon="music-note-2" class="h-6 w-6" />
				Create a new release
			</Button>
			<Button variant="outline" href="/upload/drafts">
				<Icon icon="file-text" variant="line" class="h-6 w-6 fill-foreground" />
				Manage your drafts
			</Button>
		</div>
	</div>
{/if}
