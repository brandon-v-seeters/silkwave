<script lang="ts">
	import { page } from '$app/state';
	import { Button } from '$lib/components/ui/button/index.js';
	import UserAuthForm from '$lib/components/organisms/UserAuthForm.svelte';
	import Link from '$lib/components/atoms/Link.svelte';
	import SocialAuth from '$lib/components/molecules/SocialAuth.svelte';

	const redirectTo = page.data.redirectTo;
	const signUpUrl = redirectTo
		? `/sign-up?redirectTo=${encodeURIComponent(redirectTo)}`
		: '/sign-up';
</script>

<div
	class="container relative flex h-screen flex-col items-center justify-center md:grid lg:max-w-none lg:px-0"
>
	<Button href={signUpUrl} variant="ghost" class="absolute right-4 top-4 md:right-8 md:top-8">
		Sign Up
	</Button>
	<div class="mx-auto w-full max-w-lg px-0 sm:px-20 md:px-20">
		<div class="xs:w-[350px] mx-auto flex w-full flex-col justify-center">
			<div class="mb-6 flex flex-col text-center">
				<h1 class="max-w-[280px] mx-auto">Giving the power back to the artists</h1>
			</div>
			<!-- <UserAuthForm /> -->

			<SocialAuth class="my-6" />
			<p class="px-8 text-center text-sm text-foreground-muted">
				By clicking continue, you agree to our
				<Link href="/terms">Terms of Service</Link>
				and
				<Link href="/privacy">Privacy Policy</Link>
				.
			</p>
		</div>
	</div>
</div>
