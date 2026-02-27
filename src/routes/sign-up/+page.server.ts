import { POST } from '$lib/api/Api';
import { fail, redirect, isHttpError, type Actions } from '@sveltejs/kit';
import { zod4 as zod } from 'sveltekit-superforms/adapters';
import { setError, superValidate } from 'sveltekit-superforms';
import { z } from 'zod';
import { redirect as flashRedirect } from 'sveltekit-flash-message/server';
import { PUBLIC_API_URL } from '$env/static/public';

const schema = z.object({
	email: z.string().email().default(''),
	password: z
		.string()
		.min(8, { message: 'Password must be at least 8 characters long' })
		.default(''),
	isArtist: z.boolean().default(false)
});

export async function load({ locals, url }) {
	if (locals.user) {
		redirect(307, '/');
	}

	const form = await superValidate(zod(schema));
	const redirectTo = url.searchParams.get('redirectTo');

	return {
		form,
		redirectTo
	};
}

export const actions: Actions = {
	default: async (event) => {
		const { request, fetch, cookies, url } = event;
		const form = await superValidate(request, zod(schema));
		const { email, password, isArtist } = form.data;

		if (!form.valid) {
			return fail(400, { form });
		}

		try {
			const res = await fetch(`${PUBLIC_API_URL}/register`, {
				method: 'POST',
				body: JSON.stringify({ email, password, isArtist }),
				headers: {
					'Content-Type': 'application/json'
				}
			});

			// Forward the session cookie from API
			const setCookie = res.headers.get('set-cookie');
			if (!setCookie) return;

			const cookieMatch = setCookie.match(/session=([^;]+)/);
			if (!cookieMatch) return;

			cookies.set('session', cookieMatch[1], {
				path: '/',
				httpOnly: true,
				sameSite: 'strict',
				secure: import.meta.env.PROD,
				maxAge: 60 * 60 * 24 * 30 // 30 days
			});
		} catch (err) {
			if (isHttpError(err)) {
				setError(form, '', err.body.message ?? 'Registration failed');
			} else {
				setError(form, '', 'An unexpected error occurred');
			}
			return fail(400, { form });
		}

		// Check for redirectTo parameter first
		const redirectTo = url.searchParams.get('redirectTo');
		if (redirectTo) {
			flashRedirect(
				redirectTo,
				{ type: 'success', message: 'Welcome to SilkWave! Your account has been created.' },
				event
			);
		}

		// Check if user has items in cart
		const cartCookie = cookies.get('cart');
		let hasCartItems = false;
		try {
			hasCartItems = cartCookie ? JSON.parse(cartCookie).length > 0 : false;
		} catch {
			hasCartItems = false;
		}

		// Redirect to checkout if cart has items, otherwise to home
		const redirectUrl = hasCartItems ? '/checkout' : '/';

		flashRedirect(
			redirectUrl,
			{
				type: 'success',
				message: 'Welcome to SilkWave! Your account has been created.'
			},
			event
		);
	}
};
