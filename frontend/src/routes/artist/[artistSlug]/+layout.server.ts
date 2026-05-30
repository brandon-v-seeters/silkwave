import { GET } from '$lib/api/Api.js';
import { readApiData } from '$lib/api/envelope.js';
import type { ArtistProfile } from '$lib/types/generated/models.js';
import { error } from '@sveltejs/kit';

export const load = async ({ params, fetch }) => {
	const res = await GET(`/artists/${params.artistSlug}`, fetch);

	if (!(await res.ok)) {
		return error(404, {
			message: 'Not found'
		});
	}

	const body = await readApiData<{ artist?: ArtistProfile }>(res);
	const artist = body?.artist;

	if (!artist) {
		return error(404, {
			message: 'Not found'
		});
	}

	return {
		artist
	};
};
