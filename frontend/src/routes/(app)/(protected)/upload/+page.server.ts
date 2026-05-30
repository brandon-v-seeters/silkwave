import { GET, POST } from '$lib/api/Api';
import { readApiData } from '$lib/api/envelope';
import type { ReleaseWithArtist } from '$lib/types/generated/models';
import { redirect } from '@sveltejs/kit';
import { fail, setError, superValidate } from 'sveltekit-superforms';
import { zod4 } from 'sveltekit-superforms/adapters';
import { z } from 'zod';

const schema = z.object({
	name: z.string().min(2, { message: 'Artist name must be at least 2 characters.' })
});

type DraftsResponse = {
	drafts?: ReleaseWithArtist[];
};

export const load = async ({ locals, fetch }) => {
	const form = await superValidate(zod4(schema));
	const { user } = locals;

	if (!user) {
		throw redirect(307, '/login');
	}

	let drafts: ReleaseWithArtist[] = [];
	if (user.artist?._key) {
		const response = await GET('/releases/drafts', fetch, { artistKey: user.artist._key });
		const body = await readApiData<DraftsResponse>(response);
		drafts = body?.drafts ?? [];
	}

	return {
		form,
		user,
		drafts
	};
};

export const actions = {
	default: async ({ request, fetch }) => {
		const form = await superValidate(request, zod4(schema));
		const { name } = form.data;

		if (!form.valid) {
			return fail(400, { form });
		}

		try {
			await POST('/register/artist-name', fetch, { name });
		} catch (error) {
			const message =
				(error as { body?: { message?: string } })?.body?.message ??
				'Could not set your artist name. Please try again.';
			return setError(form, 'name', message);
		}

		throw redirect(303, '/upload');
	}
};
