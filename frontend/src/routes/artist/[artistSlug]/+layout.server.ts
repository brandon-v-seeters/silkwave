import { GET } from '$lib/api/Api.js';
import { error } from '@sveltejs/kit';

export const load = async ({ params, fetch, parent }) => {
	const res = await GET(`/artists/${params.artistSlug}`, fetch);

	if (!(await res.ok)) {
		return error(404, {
			message: 'Not found'
		});
	}

	const { artist } = await res.json();

	if (!artist) {
		return error(404, {
			message: 'Not found'
		});
	}

	return {
		artist
	};
};
