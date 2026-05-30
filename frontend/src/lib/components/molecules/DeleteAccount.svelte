<script lang="ts">
	// Vendor
	import { superForm } from 'sveltekit-superforms';
	import { page } from '$app/state';

	// Components
	import Button from '$lib/components/ui/button/button.svelte';
	import Icon from '$lib/components/atoms/Icon.svelte';
	import * as Dialog from '$lib/components/ui/dialog/index';
	import * as Form from '$lib/components/ui/form/index';
	import Input from '$lib/components/ui/input/input.svelte';

	const { form: formData, enhance, errors } = superForm(page.data.deleteForm);

	let { open = $bindable(false) } = $props();

	let deleteAccountModalOpen = $state(false);
</script>

<div class="mt-4 flex flex-col gap-4">
	<p class="text-base text-foreground-muted">
		Don't worry, you will have to confirm it by filling in your password to delete your account.
	</p>
	<Button
		variant="destructive"
		class="w-full md:w-fit"
		onclick={() => (deleteAccountModalOpen = true)}
	>
		<Icon icon="trash" class="h-6 w-6 fill-foreground" />
		Delete Account
	</Button>
</div>

<Dialog.Root bind:open={deleteAccountModalOpen}>
	<Dialog.Portal>
		<Dialog.Overlay />
		<Dialog.Content>
			<Dialog.Header>
				<Dialog.Title>Delete Your Account</Dialog.Title>
			</Dialog.Header>
			<Dialog.Description>
				If you've truly sailed your last SilkWave, please confirm by filling in your
				password below. We'll be over here eating ice cream and listening to sad songs.
			</Dialog.Description>

			<form method="POST" action="?/deleteAccount" use:enhance class="flex flex-col gap-4">
				<div class="flex flex-col gap-2">
					<label for="delete-password" class="text-sm font-medium">Password</label>
					<Input
						id="delete-password"
						name="password"
						type="password"
						bind:value={$formData.password}
					/>
					{#if $errors.password}
						<p class="text-sm text-destructive">{$errors.password}</p>
					{/if}
				</div>

				<Dialog.Footer>
					<Button
						type="button"
						variant="secondary"
						size="sm"
						onclick={() => (deleteAccountModalOpen = false)}
					>
						Keep My Account
					</Button>
					<Button type="submit" variant="destructive" size="sm">
						<Icon icon="trash" class="h-6 w-6 fill-foreground" />
						Delete My Account
					</Button>
				</Dialog.Footer>
			</form>
		</Dialog.Content>
	</Dialog.Portal>
</Dialog.Root>
