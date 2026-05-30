import { type ClassValue, clsx } from 'clsx';
import { twMerge } from 'tailwind-merge';
import type { Snippet } from 'svelte';

// Utility types for UI components
export type WithElementRef<T, E extends HTMLElement = HTMLElement> = T & {
	ref?: E | null;
};

export type WithoutChildren<T> = Omit<T, 'children'>;

export type WithoutChild<T> = Omit<T, 'child'>;

export type WithoutChildrenOrChild<T> = Omit<T, 'children' | 'child'>;

export const getTimestamp = () => {
	return Math.floor(Date.now() / 1000);
};

export const capitalizeFirstLetter = (string: string): string => {
	if (!string) return string;
	return string.charAt(0).toUpperCase() + string.slice(1);
};

export function cn(...inputs: ClassValue[]) {
	return twMerge(clsx(inputs));
}

export const resolveImage = (key: string, fileName: string) => {
	if (!key) {
		return ``;
	}
	return `${import.meta.env.VITE_IMAGE_BASE_URL}/images/${key}/${fileName}.png`;
};

export const isBrowser = typeof document !== 'undefined';

export function updateTheme(activeTheme: string, path: string) {
	if (!isBrowser) return;
	document.body.classList.forEach((className) => {
		if (className.match(/^theme.*/)) {
			document.body.classList.remove(className);
		}
	});

	const theme = path === '/themes' ? activeTheme : null;
	if (theme) {
		return document.body.classList.add(`theme-${theme}`);
	}
}

export const generateId = () => {
	return Math.random().toString(36).substring(2, 9);
};
