<script lang="ts">
	import { goto } from '$app/navigation';

	// Vendor
	import { page } from '$app/state';
	import Icon from '$lib/components/atoms/Icon.svelte';
	import Button from '$lib/components/ui/button/button.svelte';

	// Components
	import DataTable from '$lib/components/ui/data-table/data-table.svelte';
	import { columns } from './columns';

	const { drafts } = page.data;
</script>

<div>
	<div class="mx-auto max-w-4xl px-4 pb-4 pt-4 md:px-6 md:pb-6 md:pt-8">
		<div class="mb-4">
			<h1 class="text-xl font-medium text-neutral-900 dark:text-white md:text-2xl">
				Your Drafts
			</h1>
			<p class="mt-1 text-sm text-foreground-muted">
				Manage your drafts, edit them, and publish them when you're ready.
			</p>
		</div>
		<DataTable data={drafts} {columns}>
			{#snippet noResults()}
				<div class="flex flex-col items-center justify-center gap-2 p-4">
					<Icon icon="file-text" variant="line" class="h-6 w-6 fill-foreground" />
					<p class="text-sm text-foreground-muted">You don't have any drafts yet.</p>

					<Button variant="outline" onclick={() => goto('/upload/release')}>
						<Icon icon="plus" variant="line" class="h-6 w-6 fill-foreground" />
						Create your first draft
					</Button>
				</div>
			{/snippet}
		</DataTable>
	</div>
</div>
