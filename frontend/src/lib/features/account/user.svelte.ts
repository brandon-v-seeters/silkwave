import { getContext } from "svelte";
import type { AppUser } from '$lib/types/generated';

// Helper functions to check user state
export function getUser(): AppUser | null {
	const ctx = getContext<{ current: AppUser | null }>('user');
	return ctx?.current ?? null;
}

export function isLoggedIn(): boolean {
	const u = getUser();
	return u !== null && u !== undefined;
}

export function userIsArtist(): boolean {
	const u = getUser();
	return u?.artist !== undefined && u?.artist !== null;
}

export function userIsMember(): boolean {
	const u = getUser();
	return u !== null && u !== undefined && !userIsArtist();
}