<script>
	import Google from '../atoms/icon/Google.svelte';
	import { Button } from '../ui/button';
</script>

<div class="relative {$$props.class}">
	<div class="relative flex justify-center text-xs uppercase">
		<Button variant="secondary" class="py-4 rounded-full max-w-[280px] w-full justify-between">
			<div class="flex">
				<Google class="shrink-0" />
			</div>
			Continue with Google
			<div class="flex">
				<Google class="invisible" />
			</div>
		</Button>
	</div>
</div>
