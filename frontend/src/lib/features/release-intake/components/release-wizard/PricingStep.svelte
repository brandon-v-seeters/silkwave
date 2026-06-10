<script lang="ts">
	import { untrack } from 'svelte';
	import { Label } from '$lib/components/ui/label/index.js';
	import { getWizardContext } from './wizard.svelte';
	import Button from '$lib/components/ui/button/button.svelte';
	import Icon from '$lib/components/ui/icon/Icon.svelte';
	import Input from '$lib/components/ui/input/input.svelte';
	import {
		detectCurrency,
		getCurrencySymbol,
		CURRENCIES,
		displayToUSDCents,
		formatCents,
		type CurrencyCode
	} from '$lib/utils/currency';
	import * as Select from '$lib/components/ui/select';
	import { Checkbox } from '$lib/components/ui/checkbox';

	const wizard = getWizardContext();

	// Auto-detect user's currency on mount
	let selectedCurrency = $state<CurrencyCode>(detectCurrency());
	let currencySymbol = $derived(getCurrencySymbol(selectedCurrency));
	let isUSD = $derived(selectedCurrency === 'USD');

	// Local display values (in user's currency)
	let basePriceDisplay = $state('');
	let minimumPriceDisplay = $state('');
	let trackPriceDisplay = $state('');

	// Pay what you want toggle
	let enablePayWhatYouWant = $state(false);
	let enableTrackPrice = $state(false);

	// Convert display values to USD cents and sync with wizard
	// Using untrack to prevent reading wizard.pricing from creating a dependency loop
	$effect(() => {
		const baseAmount = parseFloat(basePriceDisplay) || 0;
		const basePriceCents = displayToUSDCents(baseAmount, selectedCurrency);
		untrack(() => {
			wizard.pricing = {
				...wizard.pricing,
				basePrice: basePriceCents
			};
		});
	});

	$effect(() => {
		const minAmount = parseFloat(minimumPriceDisplay) || 0;
		const minPriceCents = enablePayWhatYouWant
			? displayToUSDCents(minAmount, selectedCurrency)
			: null;
		untrack(() => {
			wizard.pricing = {
				...wizard.pricing,
				minimumPrice: minPriceCents
			};
		});
	});

	$effect(() => {
		const trackAmount = parseFloat(trackPriceDisplay) || 0;
		const trackPriceCents = enableTrackPrice
			? displayToUSDCents(trackAmount, selectedCurrency)
			: null;
		untrack(() => {
			wizard.pricing = {
				...wizard.pricing,
				trackPrice: trackPriceCents
			};
		});
	});

	// USD equivalent displays
	let basePriceUSD = $derived(formatCents(wizard.pricing.basePrice, 'USD'));
	let minimumPriceUSD = $derived(
		wizard.pricing.minimumPrice !== null
			? formatCents(wizard.pricing.minimumPrice, 'USD')
			: null
	);
	let trackPriceUSD = $derived(
		wizard.pricing.trackPrice !== null ? formatCents(wizard.pricing.trackPrice, 'USD') : null
	);

	const currencyOptions = Object.entries(CURRENCIES) as [
		CurrencyCode,
		{ symbol: string; name: string; rate: number }
	][];
</script>

<div class="flex flex-col gap-4 md:gap-6">
	<!-- Summary of previous steps -->
	<div class="flex items-center gap-3 rounded-lg bg-muted-background p-3 md:gap-4 md:p-4">
		{#if wizard.coverArtPreview}
			<img
				src={wizard.coverArtPreview}
				alt="Cover art"
				class="h-12 w-12 rounded-md object-cover md:h-16 md:w-16"
			/>
		{/if}
		<div class="min-w-0 flex-1">
			<p
				class="truncate text-sm font-medium text-neutral-900 dark:text-neutral-50 md:text-base"
			>
				{wizard.releaseTitle || 'Untitled Release'}
			</p>
			<p class="truncate text-xs text-neutral-500 md:text-sm">
				{wizard.genres.length > 0 ? wizard.genres.join(', ') : 'No genres selected'}
			</p>
		</div>
		<Button variant="ghost" onclick={() => wizard.goToStep(0)} class="shrink-0 p-2">
			<Icon icon="edit" class="h-5 w-5 fill-foreground md:h-6 md:w-6" />
		</Button>
	</div>

	<div class="flex flex-col gap-4 md:gap-5">
		<!-- Currency Selector -->
		<div class="flex flex-col gap-1">
			<Label class="text-sm font-medium">Your Currency</Label>
			<Select.Root type="single" name="currency" bind:value={selectedCurrency}>
				<Select.Trigger>
					<span class="truncate"
						>{selectedCurrency} - {CURRENCIES[selectedCurrency].name}</span
					>
				</Select.Trigger>
				<Select.Content>
					<Select.Group>
						{#each currencyOptions as [code, { name }] (code)}
							<Select.Item value={code} label={`${code} - ${name}`} />
						{/each}
					</Select.Group>
				</Select.Content>
			</Select.Root>
			<p class="text-xs text-foreground-muted">
				Enter prices in your currency. We'll convert to USD.
			</p>
		</div>

		<!-- Base Price -->
		<div class="flex flex-col gap-1">
			<Label class="text-sm font-medium md:text-base">Release Price</Label>
			<Input
				disabled={enablePayWhatYouWant}
				type="number"
				id="base-price"
				prefix={currencySymbol}
				bind:value={basePriceDisplay}
				placeholder="9.99"
				step="0.01"
				min="0"
			/>
			{#if !isUSD && basePriceDisplay}
				<p class="text-xs text-foreground-muted">≈ {basePriceUSD} USD</p>
			{:else}
				<p class="text-xs text-foreground-muted">Price for the complete release</p>
			{/if}
		</div>

		<!-- Pay What You Want -->
		<div class="flex flex-col gap-2 md:gap-3">
			<label class="flex cursor-pointer items-center gap-2">
				<Checkbox id="pwyw" bind:checked={enablePayWhatYouWant} />
				<span class="text-sm font-medium">Enable "Pay What You Want"</span>
			</label>
			{#if enablePayWhatYouWant}
				<div class="flex flex-col gap-1">
					<Label class="text-sm">Minimum Price</Label>
					<Input
						type="number"
						id="min-price"
						prefix={currencySymbol}
						bind:value={minimumPriceDisplay}
						placeholder="0.00"
						step="0.01"
						min="0"
					/>
					{#if !isUSD && minimumPriceDisplay}
						<p class="text-xs text-foreground-muted">≈ {minimumPriceUSD} USD</p>
					{:else}
						<p class="text-xs text-foreground-muted">Set to 0 for "name your price"</p>
					{/if}
				</div>
			{/if}
		</div>

		<!-- Individual Track Price -->
		<div class="flex flex-col gap-2 md:gap-3">
			<label
				class="flex items-center gap-2
					{enablePayWhatYouWant ? 'cursor-not-allowed' : 'cursor-pointer'}"
			>
				<Checkbox
					id="track-pricing"
					bind:checked={enableTrackPrice}
					disabled={enablePayWhatYouWant}
				/>
				<span
					class="text-sm font-medium
						{enablePayWhatYouWant ? 'text-foreground-muted opacity-50' : ''}"
				>
					Allow individual track purchases
				</span>
			</label>
			{#if enableTrackPrice}
				<div class="flex flex-col gap-1">
					<Label class="text-sm">Price per Track</Label>
					<Input
						type="number"
						id="track-price"
						prefix={currencySymbol}
						bind:value={trackPriceDisplay}
						placeholder="0.99"
						step="0.01"
						min="0"
					/>
					{#if !isUSD && trackPriceDisplay}
						<p class="text-xs text-foreground-muted">
							Converted: ≈ {trackPriceUSD} USD
						</p>
					{:else}
						<p class="text-xs text-foreground-muted">Price for each individual track</p>
					{/if}
				</div>
			{/if}
		</div>
	</div>
</div>
