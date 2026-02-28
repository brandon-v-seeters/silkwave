<script lang="ts">
	import './layout.css';

	// Vendor
	import { getFlash } from 'sveltekit-flash-message';
	import { toast, Toaster } from 'svelte-sonner';
	import { page } from '$app/state';
	import { setContext } from 'svelte';
	import { ModeWatcher } from 'mode-watcher';

	const flash = getFlash(page);

	let { children } = $props();

	$effect(() => {
		if (!$flash) return;

		switch ($flash.type) {
			case 'success':
				toast.success($flash.message, { duration: 5000 });
				break;
			case 'error':
				toast.error($flash.message, { duration: 5000 });
				break;
			default:
				toast($flash.message, { duration: 5000 });
		}

		$flash = undefined;
	});

	// Pass a getter so context consumers get reactive updates
	setContext('user', {
		get current() {
			console.log(page.data.user);
			return page.data.user;
		}
	});
</script>

<ModeWatcher defaultMode="system" />

<Toaster
	toastOptions={{
		classes: {
			toast: 'bg-card border-foreground/10 text-foreground',
			title: 'text-foreground',
			description: 'text-muted-foreground',
			success: '[&_[data-icon]]:text-success',
			error: '[&_[data-icon]]:text-destructive',
			info: '[&_[data-icon]]:text-info',
			warning: '[&_[data-icon]]:text-warning'
		}
	}}
/>

{@render children?.()}
