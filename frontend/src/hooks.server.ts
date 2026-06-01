import type { Handle } from '@sveltejs/kit';
import { apiUrl } from '$lib/api/Api.js';

export const handle: Handle = async ({ event, resolve }) => {
	if (event.url.pathname.includes('vite') && event.url.pathname.includes('/errors/no-access')) {
		return resolve(event);
	}

	if (event.url.pathname.startsWith('/api/') || event.url.pathname === '/api' || event.url.pathname === '/health') {
		return resolve(event);
	}

	// get cookies from browser
	const session = event.cookies.get('session');

	// if there is no session load page as normal
	if (!session) {
		event.locals.user = null;
		return await resolve(event);
	}

	try {
		// find the user based on the session
		const res = await event.fetch(apiUrl('/user'), {
			headers: {
				Cookie: `session=${session}`
			}
		});

		if (res.status !== 200) {
			throw new Error('Failed to get user');
		}

		// if `user` exists set `events.locals.user`
		event.locals.user = (await res.json()) || null;
	} catch (e) {
		// if the session is invalid or expired, remove the cookie
		console.log('Session error (token expired or invalid):', e instanceof Error ? e.message : e);
		event.cookies.delete('session', {
			path: '/',
			httpOnly: true,
			sameSite: 'strict',
			secure: import.meta.env.PROD
		});

		// set user to null for invalid/expired sessions
		event.locals.user = null;
	}

	// load page as normal
	return await resolve(event);
};
