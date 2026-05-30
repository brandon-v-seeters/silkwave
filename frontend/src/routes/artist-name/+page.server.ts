import { POST } from '$lib/api/Api.js';
import { fail, type Actions } from '@sveltejs/kit';
import { zod4 as zod } from 'sveltekit-superforms/adapters';
import { setError, superValidate } from 'sveltekit-superforms';
import { z } from 'zod';
import { setFlash, redirect } from 'sveltekit-flash-message/server';

const schema = z.object({
	name: z.string().min(3, { message: 'Artist name must be at least 3 characters long.' }).default('')
});

export async function load({ locals }) {
	if (!locals.user || !locals.user.invalidArtistName) {
		return redirect(307, '/');
	}

	const form = await superValidate(zod(schema));

	return {
		form
	};
}

export const actions: Actions = {
	default: async (event) => {
		const { request, fetch } = event;
		const form = await superValidate(request, zod(schema));
		const { name } = form.data;

		if (!form.valid) {
			return fail(400, { form });
		}

		const res = await POST('/register/artist-name', fetch, { name }).catch((e) => e);

		if (res.status !== 200) {
			setFlash({ type: 'error', message: res.body.message }, event.cookies);
			return setError(form, 'name', res.body.message);
		}

		throw redirect(
			'/home',
			{
				type: 'success',
				message: 'Your new artist name has been set!'
			},
			event
		);
	}
};
