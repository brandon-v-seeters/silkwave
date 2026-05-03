// Vendor 
import Cookies from 'js-cookie';
import { writable, derived, get } from 'svelte/store';
import { browser } from '$app/environment';

export type CartItem = {
	id: string;
	type: 'track' | 'album';
	title: string;
	artist: string;
	price: number;
	coverUrl?: string;
};

const CART_COOKIE_NAME = 'cart';
const CART_EXPIRY_DAYS = 30;
const CART_EXPIRY_MAX_AGE = 60 * 60 * 24 * CART_EXPIRY_DAYS;

// Load initial cart from localStorage
const loadCart = (): CartItem[] => {
	if (!browser) return [];

	try {
		const cookieValue = Cookies.get(CART_COOKIE_NAME);
		if (!cookieValue) return [];
		return JSON.parse(cookieValue);
	} catch (error) {
		console.error('Failed to load cart from cookie:', error);
		return [];
	}
}

// Sync just IDs to cookie (for server-side access)
const syncToCookie = (items: CartItem[]) => {
	if (!browser) return;

	try {
		const ids = items.map(item => item.id);
		Cookies.set(CART_COOKIE_NAME, JSON.stringify(ids), {
			expires: CART_EXPIRY_MAX_AGE,
			path: '/',
			sameSite: 'strict'
		});
	} catch (error) {
		console.error('Failed to sync cart to cookie:', error);
	}
}

// Clear the cookie
const clearCookie = () => {
	if (!browser) return;

	try {
		Cookies.remove(CART_COOKIE_NAME, {
			path: '/',
			sameSite: 'strict'
		});
	} catch (error) {
		console.error('Failed to clear cart cookie:', error);
	}
}

const createCartStore = () => {
	const { subscribe, set, update } = writable<CartItem[]>(loadCart());

	// Persist changes to both localStorage and cookie
	function persist(items: CartItem[]) {
		syncToCookie(items);
	}

	return {
		subscribe,

		add: (item: CartItem) => {
			update(items => {
				try {
					// Prevent duplicates
					if (items.some(i => i.id === item.id)) return items;
					const updated = [...items, item];
					persist(updated);
					return updated;
				} catch (error) {
					console.error('Failed to add item to cart:', error);
					return items;
				}
			});
		},

		remove: (id: string) => {
			update(items => {
				try {
					const updated = items.filter(i => i.id !== id);
					persist(updated);
					return updated;
				} catch (error) {
					console.error('Failed to remove item from cart:', error);
					return items;
				}
			});
		},

		clear: () => {
			try {
				set([]);
				persist([]);
				clearCookie();
			} catch (error) {
				console.error('Failed to clear cart:', error);
			}
		},

		// Check if item is in cart
		has: (id: string): boolean => {
			try {
				return get({ subscribe }).some(item => item.id === id);
			} catch (error) {
				console.error('Failed to check cart:', error);
				return false;
			}
		},

		// Get current cart (useful outside reactive context)
		get: (): CartItem[] => {
			try {
				return get({ subscribe });
			} catch (error) {
				console.error('Failed to get cart:', error);
				return [];
			}
		}
	};
}

export const cart = createCartStore();

// Derived stores for common values
export const cartCount = derived(cart, $cart => $cart.length);
export const cartTotal = derived(cart, $cart =>
	$cart.reduce((sum, item) => sum + item.price, 0)
);
