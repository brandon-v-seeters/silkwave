<script lang="ts">
	import { page } from '$app/state';
</script>

<div class="container">
	<h1 class="text-2xl font-bold">Error</h1>
	<p class="text-lg">{page.data.error}</p>
</div>
