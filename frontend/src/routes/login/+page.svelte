<script lang="ts">
	import { page } from '$app/state';
	import Link from '$lib/components/ui/link/Link.svelte';
	import { SocialAuth, UserAuthForm } from '$lib/features/account';
	import { Button } from '$lib/components/ui/button/index.js';

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
		<div class="mb-6 flex flex-col space-y-2 text-center">
			<h1 class="font-serif text-3xl font-extralight">Welcome back</h1>
		</div>
		<UserAuthForm />
		<div class="mt-3 flex flex-col space-y-2 text-center">
			<p class="text-center text-sm text-foreground-muted">
				Don't have an account?
				<Link href={signUpUrl}>Sign Up</Link>
			</p>
		</div>
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
