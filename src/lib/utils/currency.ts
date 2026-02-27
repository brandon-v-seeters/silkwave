// Currency utilities for international support

export type CurrencyCode = 'USD' | 'EUR' | 'GBP' | 'JPY' | 'CAD' | 'AUD' | 'CHF' | 'CNY' | 'INR' | 'BRL';

export const CURRENCIES: Record<CurrencyCode, { symbol: string; name: string; rate: number }> = {
	USD: { symbol: '$', name: 'US Dollar', rate: 1 },
	EUR: { symbol: '€', name: 'Euro', rate: 0.92 },
	GBP: { symbol: '£', name: 'British Pound', rate: 0.79 },
	JPY: { symbol: '¥', name: 'Japanese Yen', rate: 149.50 },
	CAD: { symbol: 'CA$', name: 'Canadian Dollar', rate: 1.36 },
	AUD: { symbol: 'A$', name: 'Australian Dollar', rate: 1.53 },
	CHF: { symbol: 'CHF', name: 'Swiss Franc', rate: 0.88 },
	CNY: { symbol: '¥', name: 'Chinese Yuan', rate: 7.24 },
	INR: { symbol: '₹', name: 'Indian Rupee', rate: 83.12 },
	BRL: { symbol: 'R$', name: 'Brazilian Real', rate: 4.97 }
};

// Detect user's preferred currency based on locale
export function detectCurrency(): CurrencyCode {
	if (typeof navigator === 'undefined') return 'USD';
	
	const locale = navigator.language || 'en-US';
	const region = locale.split('-')[1]?.toUpperCase();
	
	const regionToCurrency: Record<string, CurrencyCode> = {
		US: 'USD',
		GB: 'GBP',
		UK: 'GBP',
		DE: 'EUR',
		FR: 'EUR',
		IT: 'EUR',
		ES: 'EUR',
		NL: 'EUR',
		BE: 'EUR',
		AT: 'EUR',
		IE: 'EUR',
		PT: 'EUR',
		FI: 'EUR',
		GR: 'EUR',
		JP: 'JPY',
		CA: 'CAD',
		AU: 'AUD',
		CH: 'CHF',
		CN: 'CNY',
		IN: 'INR',
		BR: 'BRL'
	};
	
	return regionToCurrency[region] || 'USD';
}

// Get user's locale
export function getLocale(): string {
	if (typeof navigator === 'undefined') return 'en-US';
	return navigator.language || 'en-US';
}

// Format a number as currency
export function formatCurrency(
	amount: number,
	currency: CurrencyCode = 'USD',
	locale?: string
): string {
	const userLocale = locale || getLocale();
	
	return new Intl.NumberFormat(userLocale, {
		style: 'currency',
		currency,
		minimumFractionDigits: currency === 'JPY' ? 0 : 2,
		maximumFractionDigits: currency === 'JPY' ? 0 : 2
	}).format(amount);
}

// Format cents as currency (999 cents = $9.99)
export function formatCents(
	cents: number,
	currency: CurrencyCode = 'USD',
	locale?: string
): string {
	const amount = currency === 'JPY' ? cents : cents / 100;
	return formatCurrency(amount, currency, locale);
}

// Get just the currency symbol for the locale
export function getCurrencySymbol(currency: CurrencyCode = 'USD', locale?: string): string {
	const userLocale = locale || getLocale();
	
	return new Intl.NumberFormat(userLocale, {
		style: 'currency',
		currency,
		minimumFractionDigits: 0,
		maximumFractionDigits: 0
	})
		.format(0)
		.replace(/[\d.,\s]/g, '')
		.trim();
}

// Convert from a currency to USD (returns amount in the same unit - dollars or cents)
export function convertToUSD(amount: number, fromCurrency: CurrencyCode): number {
	const rate = CURRENCIES[fromCurrency].rate;
	return amount / rate;
}

// Convert from USD to another currency
export function convertFromUSD(amount: number, toCurrency: CurrencyCode): number {
	const rate = CURRENCIES[toCurrency].rate;
	return amount * rate;
}

// Convert a display amount (e.g., 9.99) to cents in USD
export function displayToUSDCents(displayAmount: number, fromCurrency: CurrencyCode): number {
	// For JPY, the display amount is already in the smallest unit
	const amountInFromCurrency = fromCurrency === 'JPY' ? displayAmount : displayAmount;
	const amountInUSD = convertToUSD(amountInFromCurrency, fromCurrency);
	// Convert to cents (multiply by 100 for non-JPY currencies)
	return Math.round(amountInUSD * 100);
}

// Convert USD cents to display amount in a currency
export function usdCentsToDisplay(cents: number, toCurrency: CurrencyCode): number {
	// Convert cents to dollars
	const usdAmount = cents / 100;
	// Convert to target currency
	const targetAmount = convertFromUSD(usdAmount, toCurrency);
	// Round to 2 decimal places (or 0 for JPY)
	return toCurrency === 'JPY' 
		? Math.round(targetAmount) 
		: Math.round(targetAmount * 100) / 100;
}

// Parse a currency string back to a number
export function parseCurrencyInput(value: string): number {
	// Remove all non-numeric characters except decimal point and minus
	const cleaned = value.replace(/[^\d.,-]/g, '');
	// Handle European format (comma as decimal separator)
	const normalized = cleaned.replace(',', '.');
	return parseFloat(normalized) || 0;
}
