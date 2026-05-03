import type { Route, UserMode } from '$lib/types/Route';

export const SIDEBAR_ROUTES: Record<UserMode, Route[]> = {
	artist: [
		{
			title: 'Home',
			href: '/',
			icon: 'home'
		},
		{
			title: 'Courses',
			href: '/courses',
			icon: 'book'
		},
		{
			title: 'Subscriptions',
			href: '/subscriptions',
			icon: 'book'
		},
		{
			title: 'Fans',
			href: '/fans',
			icon: 'users'
		},
		{
			title: 'Trends',
			href: '/trends',
			icon: 'chart-line'
		},
		{
			title: 'Earnings',
			href: '/earnings',
			icon: 'cash'
		},
		{
			title: 'Community',
			href: '/community',
			icon: 'comments-2'
		},
		{
			title: 'Settings',
			href: '/settings',
			icon: 'gear'
		}
	],
	member: [
		{
			title: 'Home',
			href: '/home',
			icon: 'home'
		},
		{
			title: 'Discover',
			href: '/discover',
			icon: 'search'
		},
		{
			title: 'Community',
			href: '/community',
			icon: 'comments-2'
		},
		{
			title: 'Settings',
			href: '/settings',
			icon: 'gear'
		}
	]
};

export const ARTIST_PROFILE_ROUTES = [
	{
		href: '',
		title: 'Home',
		icon: 'home'
	},
	{
		href: '/posts',
		title: 'Posts',
		icon: 'users'
	},
	{
		href: '/archives',
		title: 'Archives',
		icon: 'chart-line'
	},
	{
		href: '/chat',
		title: 'Chat',
		icon: 'cash'
	},
	{
		href: '/subscription',
		title: 'Subscription',
		icon: 'comments-2'
	},
	{
		href: '/about',
		title: 'About',
		icon: 'gear'
	}
];

export const PROTECTED_ROUTES = [
	'/about-us',
	'/discover',
	'/profile',
	'/artist',
	'/login',
	'/register'
];
