<script lang="ts">
	import { Google } from '$lib/components/ui/icon';
	import { Button } from '$lib/components/ui/button';

	let { class: className }: { class?: string } = $props();
</script>

<div class="relative {className ?? ''}">
	<div class="relative flex items-center gap-3">
		<div class="h-px flex-1 bg-border"></div>
		<span class="text-xs text-foreground-muted">or</span>
		<div class="h-px flex-1 bg-border"></div>
	</div>
	<div class="mt-4 flex justify-center">
		<Button variant="secondary" class="rounded-full py-4 max-w-[280px] w-full">
			<Google class="shrink-0" />
			Continue with Google
		</Button>
	</div>
</div>
