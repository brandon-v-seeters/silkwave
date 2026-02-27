<script lang="ts">
	import { enhance } from '$app/forms';
	import { page } from '$app/stores';
	import Icon from '$lib/components/atoms/Icon.svelte';
	import Button from '$lib/components/ui/button/button.svelte';
	import * as Form from '$lib/components/ui/form/index.ts';
	import Input from '$lib/components/ui/input/input.svelte';
	import { superForm } from 'sveltekit-superforms';

	const form = superForm($page.data.form);
	const { form: formData } = form;

	let isLoading = false;
</script>

<div
	class="container relative flex h-screen flex-col items-center justify-center md:px-0 lg:max-w-none lg:grid-cols-2"
>
	<div class="mx-auto w-full max-w-lg px-0 sm:px-20 md:px-20">
		<div class="xs:w-[350px] mx-auto flex w-full flex-col justify-center space-y-6">
			<div class="flex flex-col space-y-2 text-center">
				<h1 class="font-serif text-2xl font-light">Set your Artist Name</h1>
				<p class="text-base text-foreground-muted">
					Enter your Artist Name below to complete your account
				</p>
				<form method="POST" use:enhance class="flex w-full flex-col gap-4">
					<Form.Field name="name" {form}>
						<Form.Control let:attrs>
							<Form.Label class="sr-only" for="name">Artist Name</Form.Label>
							<Input
								class="w-full"
								type="text"
								placeholder="Silkwaver"
								required
								bind:value={$formData.name}
								{...attrs}
							/>
						</Form.Control>
					</Form.Field>
					<Button type="submit" variant="gradient" disabled={isLoading}>
						{#if isLoading}
							<Icon icon="loader-2" class="mr-2 h-4 w-4 animate-spin" />
						{/if}
						Set my Artist Name
					</Button>
				</form>
			</div>
		</div>
	</div>
</div>
