<script lang="ts">
	import Icon from '../icon/Icon.svelte';
	import type { IconKey, IconVariant } from '$lib/types/Icon';

	const { title, description, variant = 'info' } = $props();

	const alertClasses: Record<string, string> = {
		info: 'bg-indigo-500/[0.05]',
		warning: 'bg-yellow-500/[0.05]',
		error: 'bg-rose-400/[0.05]',
		success: 'bg-success/[0.03]'
	};

	const iconClasses: Record<string, string> = {
		info: 'fill-indigo-500',
		warning: 'fill-yellow-500',
		danger: 'fill-rose-400',
		success: 'fill-success'
	};

	const iconKeys: Record<string, IconKey> = {
		info: 'info',
		warning: 'alert-circle',
		danger: 'alert-triangle',
		success: 'check-circle'
	};
</script>

<div class="flex w-full gap-3 rounded-md border-border {alertClasses[variant]} p-4">
	<Icon
		icon={iconKeys[variant]}
		variant="filled"
		class="h-6 w-6 flex-shrink-0 {iconClasses[variant]}"
	/>
	<div class="flex flex-col">
		<h5 class="text-base font-medium">{title}</h5>
		<p class="text-md text-foreground-muted">{description}</p>
	</div>
</div>
