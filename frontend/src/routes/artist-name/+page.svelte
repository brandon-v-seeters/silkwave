<script lang="ts">
	import { page } from '$app/state';
	import Icon from '$lib/components/atoms/Icon.svelte';
	import FormErrors from '$lib/components/atoms/form/FormErrors.svelte';
	import Button from '$lib/components/ui/button/button.svelte';
	import * as Form from '$lib/components/ui/form/index';
	import Input from '$lib/components/ui/input/input.svelte';
	import { superForm } from 'sveltekit-superforms';

	const form = superForm(page.data.form);
	const { form: formData, enhance, delayed } = form;
</script>

<div class="relative flex h-screen flex-col items-center justify-center lg:px-0">
	<div class="mx-auto w-full max-w-xs">
		<div class="mb-6 flex flex-col space-y-2 text-center">
			<h1 class="font-serif text-3xl font-extralight">Set your artist name</h1>
			<p class="text-sm text-foreground-muted">
				This is how fans will find you on Silkwave. You can change it later.
			</p>
		</div>

		<form method="POST" use:enhance class="flex w-full flex-col gap-4">
			<Form.Field {form} name="name">
				<Form.Control>
					{#snippet children({ props })}
						<Form.Label class="sr-only">Artist name</Form.Label>
						<Input
							{...props}
							type="text"
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

			<Button type="submit" variant="gradient" disabled={$delayed}>
				{#if $delayed}
					<Icon icon="loader-2" class="mr-2 h-4 w-4 animate-spin" />
				{/if}
				Set my artist name
			</Button>
		</form>
	</div>
</div>
