<script lang="ts">
	import ReleaseDetailPage from '$lib/features/catalog/components/ReleaseDetailPage.svelte';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();
</script>

<ReleaseDetailPage release={data.release} />
