// Client-only type - routes are frontend-specific
import type { IconKey } from './Icon';

export type UserMode = 'member' | 'artist';

export type Route = {
	title: string;
	href: string;
	icon?: IconKey;
	modes?: UserMode[];
	children?: Route[];
};
