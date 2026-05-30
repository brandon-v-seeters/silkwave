import type { AppUser } from '$lib/types/generated';

declare global {
	namespace App {
		// interface Error {}
		interface Locals {
			user: AppUser | null;
		}
		// interface PageData {}
		// interface PageState {}
		// interface Platform {}
	}
}

export {};
