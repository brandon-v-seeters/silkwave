<script lang="ts">
	import Icon from '../icon/Icon.svelte';

	let { onRemove }: { onRemove?: () => void } = $props();
</script>

<button
	class="
			flex flex-col rounded-md fill-foreground p-2 text-muted-foreground transition-colors
			hover:bg-muted-background hover:fill-foreground hover:text-destructive-foreground"
>
	<Icon icon="chevron-up" class="h-5 w-5" />
	<Icon icon="chevron-down" class="h-5 w-5" />
</button>
<button
	type="button"
	onclick={onRemove}
	class="rounded-md fill-foreground p-2 text-muted-foreground transition-colors hover:bg-destructive-foreground/10 hover:fill-destructive-foreground hover:text-destructive-foreground"
>
	<Icon icon="trash" class="h-4 w-4" />
</button>
