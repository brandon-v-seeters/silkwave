<script lang="ts">
	import Icon from '$lib/components/ui/icon/Icon.svelte';
	import { AccountSettings, DeleteAccount } from '$lib/features/account';
	import Button from '$lib/components/ui/button/button.svelte';
	import * as Tabs from '$lib/components/ui/tabs/index.js';
</script>

<svelte:head>
	<title>Settings - Silk Wave</title>
</svelte:head>

<div class="mx-auto max-w-3xl">
	<h1 class="font-serif text-2xl md:text-3xl">Settings</h1>
	<p class="mt-2 text-base text-foreground">Manage your account settings</p>
	<Tabs.Root value="account" class="mt-4">
		<Tabs.List class="mb-4 flex w-full flex-wrap ">
			<Tabs.Trigger value="account">Account</Tabs.Trigger>
			<Tabs.Trigger value="subscriptions">Subscriptions</Tabs.Trigger>
			<Tabs.Trigger value="purchase-history">Purchase History</Tabs.Trigger>
			<Tabs.Trigger value="notifications">Notifications</Tabs.Trigger>
			<Tabs.Trigger value="delete">Delete Account</Tabs.Trigger>
		</Tabs.List>
		<Tabs.Content value="account">
			<AccountSettings />
		</Tabs.Content>
		<Tabs.Content value="subscriptions">View and manage your subscriptions here.</Tabs.Content>
		<Tabs.Content value="purchase-history">See your purchase history here.</Tabs.Content>
		<Tabs.Content value="payment-methods">
			Add, edit, or delete your payment methods here.
		</Tabs.Content>
		<Tabs.Content value="notifications">
			Configure your notification preferences here.
		</Tabs.Content>
		<Tabs.Content value="delete">
			<DeleteAccount />
		</Tabs.Content>
	</Tabs.Root>
</div>
