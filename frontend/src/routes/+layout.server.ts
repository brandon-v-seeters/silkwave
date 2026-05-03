import { redirect } from '@sveltejs/kit';
import { loadFlash } from 'sveltekit-flash-message/server';

const PUBLIC_ROUTES = ['/login', '/sign-up'];

export const load = loadFlash(async ({ locals, url, cookies }) => {
	if (import.meta.env.PROD && url.pathname !== '/') {
		throw redirect(307, '/');
	}

	const { user } = locals;
	const theme = cookies.get('theme');

	if (!user && (url.pathname === '/' || PUBLIC_ROUTES.includes(url.pathname))) {
		return {};
	}

	return { user, theme };


	// switch (user?.mode) {
	// 	case 'member':
	// 		return redirect(307, '/home');
	// 	case 'artist':
	// 		return redirect(307, '/' + user.artist.slug);
	// 	default:
	// 		return redirect(307, '/home');
	// }
});