<script lang="ts">
	import { page } from '$app/state';
	import { Button } from '$lib/components/ui/button/index.js';
	import Link from '$lib/components/atoms/Link.svelte';
	import SocialAuth from '$lib/components/molecules/SocialAuth.svelte';
	import UserAuthForm from '$lib/components/organisms/UserAuthForm.svelte';

	const redirectTo = page.data.redirectTo;
	const signUpUrl = redirectTo
		? `/sign-up?redirectTo=${encodeURIComponent(redirectTo)}`
		: '/sign-up';
</script>

<div class="relative flex h-screen flex-col items-center justify-center lg:px-0">
	<Button href={signUpUrl} variant="ghost" class="absolute right-4 top-4 md:right-8 md:top-8">
		Sign Up
	</Button>
	<div class="mx-auto w-full max-w-xs">
		<div class="flex w-full flex-col justify-center">
			<div class="mb-6 flex flex-col text-left">
				<h1 class="">Start selling music</h1>
			</div>
			<UserAuthForm />

			<p class="mt-2 text-left text-sm text-foreground-muted">
				By clicking continue, you agree to our
				<Link href="/terms">Terms of Service</Link>
				and
				<Link href="/privacy">Privacy Policy</Link>
				.
			</p>
		</div>
	</div>
</div>
