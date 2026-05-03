<script lang="ts">
	import * as Drawer from '$lib/components/ui/drawer';
	import Icon from './Icon.svelte';
	import Button from '../ui/button/button.svelte';
	import { get } from 'svelte/store';
	import { page } from '$app/stores';

	const isActive = (path: string) => get(page).url.pathname === path;
</script>

<div
	class="fixed bottom-0 right-0 mr-4 sm:border-b md:bottom-auto md:left-0 md:top-0 md:flex md:w-full"
>
	<div class="mx-auto flex h-16">
		<div class="ml-8 flex items-center space-x-2">
			<Button
				variant="ghost"
				href="/login"
				class="rounded-full text-base font-medium text-foreground-muted">Login</Button
			>
			<Button variant="gradient" href="/sign-up" class="rounded-full text-base font-medium">
				Join Silk Wave
			</Button>
		</div>
	</div>
</div>
