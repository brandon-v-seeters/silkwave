<script lang="ts">
	// Vendor
	import { superForm } from 'sveltekit-superforms';
	import { page } from '$app/state';

	// Components
	import * as Form from '$lib/components/ui/form/index';
	import * as Dialog from '$lib/components/ui/dialog/index';
	import Input from '$lib/components/ui/input/input.svelte';
	import Button from '../ui/button/button.svelte';
	import Icon from '../atoms/Icon.svelte';
	import FormErrors from '../atoms/form/FormErrors.svelte';

	const emailForm = superForm(page.data.form, {
		invalidateAll: 'force',
		onResult: ({ result }) => {
			if (result.type !== 'success') return;

			saveModalOpen = false;
			touched = false;
			$emailFormData.password = '';
		}
	});
	const passwordForm = superForm(page.data.passwordForm, {
		invalidateAll: 'force',
		onResult: ({ result }) => {
			if (result.type !== 'success') return;

			changePasswordModalOpen = false;
			$passwordFormData.oldPassword = '';
			$passwordFormData.newPassword = '';
		}
	});

	const { form: emailFormData, enhance } = emailForm;
	const { form: passwordFormData, enhance: passwordEnhance } = passwordForm;

	let isLoading = $state(false);
	let saveModalOpen = $state(false);
	let touched = $state(false);
	let changePasswordModalOpen = $state(false);

	const openSaveModal = () => {
		saveModalOpen = true;
	};

	const openChangePasswordModal = () => {
		changePasswordModalOpen = true;
	};
</script>

<div class="flex flex-col">
	<p class="mb-4 text-base text-foreground-muted">
		This section is used to change settings that are related to your account.
	</p>
	<div class="flex max-w-sm flex-col gap-2 pb-6">
		<Form.Field form={emailForm} name="email">
			<Form.Control>
				{#snippet children({ props })}
					<Form.Label>Email</Form.Label>
					<Input
						{...props}
						placeholder="email@example.com"
						type="email"
						bind:value={$emailFormData.email}
						oninput={() => (touched = true)}
					/>
				{/snippet}
			</Form.Control>
		</Form.Field>

		<Button variant="primary" disabled={isLoading || !touched} onclick={openSaveModal}>
			Save changes
		</Button>
	</div>

	<div class="flex flex-col border-t border-muted-background pt-6">
		<h4 class="text-base font-semibold leading-none">Password</h4>
		<p class="mt-2 text-base text-foreground-muted">
			Of course, we can't show your password here. But you can change it here.
		</p>
		<Button
			type="button"
			variant="secondary"
			class="mt-3 w-fit"
			disabled={isLoading}
			onclick={openChangePasswordModal}
		>
			<Icon icon="edit" class="!h-5 !w-5 flex-shrink-0  fill-foreground" />
			Change your password
		</Button>
	</div>
</div>

<Dialog.Root bind:open={saveModalOpen}>
	<Dialog.Portal>
		<Dialog.Overlay />
		<Dialog.Content>
			<Dialog.Header>
				<Dialog.Title>Change your email</Dialog.Title>
			</Dialog.Header>
			<Dialog.Description>
				Before we can change your email, we need to verify your password.
			</Dialog.Description>
			<form method="POST" action="?/changeEmail" use:enhance>
				<input type="hidden" name="email" value={$emailFormData.email} />

				<Form.Field form={emailForm} name="password">
					<Form.Control>
						{#snippet children({ props })}
							<Form.Label>Password</Form.Label>
							<Input
								{...props}
								placeholder="***"
								type="password"
								bind:value={$emailFormData.password}
							/>
						{/snippet}
					</Form.Control>
				</Form.Field>

				<FormErrors form={emailForm} class="mt-3" />
				<Dialog.Footer>
					<div class="mt-4 flex gap-2">
						<Button
							size="sm"
							variant="secondary"
							type="button"
							disabled={isLoading}
							onclick={() => (saveModalOpen = false)}>Cancel</Button
						>
						<Button variant="primary" type="submit" size="sm" disabled={isLoading}>
							Save changes
						</Button>
					</div>
				</Dialog.Footer>
			</form>
		</Dialog.Content>
	</Dialog.Portal>
</Dialog.Root>

<Dialog.Root bind:open={changePasswordModalOpen}>
	<Dialog.Portal>
		<Dialog.Overlay />
		<Dialog.Content>
			<Dialog.Header>
				<Dialog.Title>Change your password</Dialog.Title>
			</Dialog.Header>
			<Dialog.Description>
				Simply enter your old password and your new password to change it.
			</Dialog.Description>
			<form method="POST" action="?/changePassword" use:passwordEnhance class="flex flex-col">
				<Form.Field form={passwordForm} name="oldPassword" class="mb-4">
					<Form.Control>
						{#snippet children({ props })}
							<Form.Label>Old password</Form.Label>
							<Input
								{...props}
								type="password"
								bind:value={$passwordFormData.oldPassword}
							/>
						{/snippet}
					</Form.Control>
				</Form.Field>
				<Form.Field form={passwordForm} name="newPassword">
					<Form.Control>
						{#snippet children({ props })}
							<Form.Label>New password</Form.Label>
							<Input
								{...props}
								type="password"
								bind:value={$passwordFormData.newPassword}
							/>
						{/snippet}
					</Form.Control>
				</Form.Field>
				<FormErrors form={passwordForm} />
				<Dialog.Footer>
					<div class="mt-4 flex gap-2">
						<Button
							size="sm"
							variant="secondary"
							type="button"
							onclick={() => (changePasswordModalOpen = false)}
						>
							Cancel
						</Button>
						<Button variant="primary" type="submit" size="sm" disabled={isLoading}>
							Save changes
						</Button>
					</div>
				</Dialog.Footer>
			</form>
		</Dialog.Content>
	</Dialog.Portal>
</Dialog.Root>
