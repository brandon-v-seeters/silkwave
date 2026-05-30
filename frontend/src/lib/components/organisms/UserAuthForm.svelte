<script lang="ts">
	import { cn } from '$lib/utils/utils';
	import Icon from '../atoms/Icon.svelte';
	import Button from '../ui/button/button.svelte';
	import Input from '../ui/input/input.svelte';
	import Checkbox from '../ui/checkbox/checkbox.svelte';
	import * as Form from '$lib/components/ui/form/index';
	import { superForm } from 'sveltekit-superforms';
	import { page } from '$app/state';
	import { toast } from 'svelte-sonner';
	import FormErrors from '../atoms/form/FormErrors.svelte';

	let {
		class: className,
		...restProps
	}: { class?: string | null | undefined } & Record<string, any> = $props();

	let isLoading = false;

	const form = superForm(page.data.form, {
		onUpdate({ form }) {
			if (!form.valid) {
				toast.error(form.message);
				return;
			}

			toast.success(form.message);
		},
		onError(error) {
			console.error(error);
		}
	});

	const { form: formData, enhance } = form;

	// Handle checkbox change - update form data directly
	function handleArtistChange(checked: boolean) {
		if ('isArtist' in $formData) {
			$formData.isArtist = checked;
		}
	}
</script>

<div class={cn('grid', className)} {...restProps}>
	<form method="POST" use:enhance>
		<div class="grid gap-4">
			<div class="grid">
				<Form.Field {form} name="email">
					<Form.Control>
						{#snippet children({ props })}
							<Form.Label>Email</Form.Label>
							<Input
								{...props}
								placeholder="email@example.com"
								type="email"
								bind:value={$formData.email}
								disabled={isLoading}
							/>
						{/snippet}
					</Form.Control>
					<Form.FieldErrors />
				</Form.Field>
			</div>
			<div class="grid">
				<Form.Field {form} name="password">
					<Form.Control>
						{#snippet children({ props })}
							<Form.Label>Password</Form.Label>
							<Input
								{...props}
								placeholder="***"
								type="password"
								bind:value={$formData.password}
								disabled={isLoading}
							/>
						{/snippet}
					</Form.Control>
					<Form.FieldErrors />
				</Form.Field>

				{#if 'isArtist' in $formData}
					<Form.Field {form} name="isArtist" class="hidden">
						<Form.Control>
							{#snippet children({ props })}
								<Form.Label class="sr-only">Is Artist</Form.Label>
								<Checkbox
									{...props}
									checked={$formData.isArtist === true}
									onCheckedChange={handleArtistChange}
									disabled={isLoading}
								/>
							{/snippet}
						</Form.Control>
					</Form.Field>
				{/if}
			</div>

			<FormErrors {form} />

			<Button type="submit" variant="primary" disabled={isLoading}>
				{#if isLoading}
					<Icon icon="loader-2" class="mr-2 h-4 w-4 animate-spin" />
				{/if}
				{page.url.pathname === '/login' ? 'Login' : 'Sign Up'}
			</Button>
		</div>
	</form>
</div>
