<script lang="ts">
	import { Avatar, AvatarImage, AvatarFallback } from '$lib/components/ui/avatar/index';
	import { resolveImage } from '$lib/utils/utils';
	import type { User, Artist } from '$lib/types/generated/models';

	const { user }: { user: User | Artist } = $props();

	// Get display name - User has username, Artist has name
	const displayName = $derived(
		'username' in user ? user.username : 'name' in user ? user.name : 'UN'
	);
</script>

<Avatar class="mr-2 h-8 w-8">
	<AvatarImage src={resolveImage(user?._key ?? '', 'avatar')} alt="Avatar" class="h-8 w-8 " />
	<AvatarFallback class="bg-primary-gradient font-semibold text-background">
		{displayName?.slice(0, 2) || 'UN'}
	</AvatarFallback>
</Avatar>
