<script lang="ts">
	import { Toaster as Sonner, type ToasterProps as SonnerProps } from 'svelte-sonner';
	// import { mode } from 'mode-watcher';

	type $$Props = SonnerProps;
</script>

<Sonner
	class="toaster group"
	toastOptions={{
		classes: {
			toast: 'group toast group-[.toaster]:bg-background dark:group-[.toaster]:bg-background group-[.toaster]:text-foreground group-[.toaster]:border-border group-[div]:shadow-lg',
			description: 'group-[.toast]:text-foreground-muted',
			actionButton: 'group-[.toast]:bg-primary group-[.toast]:text-primary-foreground',
			cancelButton: 'group-[.toast]:bg-muted group-[.toast]:text-foreground-muted'
		}
	}}
	{...$$restProps}
/>
