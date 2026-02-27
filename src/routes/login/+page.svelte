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
	class="container relative flex h-screen flex-col items-center justify-center md:grid lg:max-w-none lg:grid-cols-2 lg:px-0"
>
	<Button href={signUpUrl} variant="ghost" class="absolute right-4 top-4 md:right-8 md:top-8">
		Sign Up
	</Button>
	<div class="relative hidden h-full flex-col bg-muted p-10 text-white dark:border-r lg:flex">
		<div
			class="absolute inset-0 bg-cover bg-center bg-no-repeat opacity-90 dark:opacity-75"
			style="
				background-image:
					url('/producer.jpg');"
		></div>
		<a class="relative z-20 flex items-center font-serif" href="/">Silk Wave</a>
		<div class="relative z-20 mt-auto">
			<blockquote class="space-y-1">
				<p class="text-lg">
					"Music is the divine way to tell beautiful, poetic things to the heart."
				</p>
				<footer class="text-base">AI Generated Quote</footer>
			</blockquote>
		</div>
	</div>
	<div class="mx-auto w-full max-w-lg px-0 sm:px-20 md:px-20">
		<div class="xs:w-[350px] mx-auto flex w-full flex-col justify-center">
			<div class="mb-6 flex flex-col text-center">
				<h1 class="font-serif text-2xl">Welcome back!</h1>
				<p class="text-base text-foreground-muted">Enter your details to sign back in.</p>
			</div>
			<UserAuthForm />
			<p class="mt-3 text-center text-sm text-foreground-muted">
				Forgot your password?
				<Link href="/forgot-password">Reset it</Link>
			</p>
			<p class="mt-1 text-center text-sm text-foreground-muted">
				No account?
				<Link href={signUpUrl}>Sign up</Link>
			</p>

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
