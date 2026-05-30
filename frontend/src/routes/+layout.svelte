<script lang="ts">
	import './layout.css';

	// Vendor
	import { getFlash } from 'sveltekit-flash-message';
	import { toast, Toaster } from 'svelte-sonner';
	import { page } from '$app/state';
	import { setContext } from 'svelte';
	import { ModeWatcher, mode } from 'mode-watcher';

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
			return page.data.user;
		}
	});
</script>

<ModeWatcher defaultMode="system" />

<Toaster
	theme={mode.current ?? 'system'}
	toastOptions={{
		classes: {
			toast: 'bg-card border-border text-foreground',
			title: 'text-foreground',
			description: 'text-foreground',
			success: '[&_[data-icon]]:text-success',
			error: '[&_[data-icon]]:text-destructive',
			info: '[&_[data-icon]]:text-info',
			warning: '[&_[data-icon]]:text-warning'
		}
	}}
/>

{@render children?.()}

<style>
	:global(.dark [data-sonner-toast][data-styled='true']) {
		background: var(--card) !important;
		border-color: var(--border) !important;
		color: var(--foreground) !important;
	}

	:global(.dark [data-sonner-toast] [data-title]),
	:global(.dark [data-sonner-toast] [data-description]) {
		color: var(--foreground) !important;
	}
</style>
