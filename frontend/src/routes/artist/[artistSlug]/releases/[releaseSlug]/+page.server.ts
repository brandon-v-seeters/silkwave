import { GET } from '$lib/api/Api.js';
import type { PublicRelease } from '$lib/types/generated/models.js';
import { error, isHttpError } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';

type ApiEnvelope<T> = {
	data?: T;
	error?: {
		message?: string;
	};
};

export const load: PageServerLoad = async ({ params, fetch }) => {
	try {
		const res = await GET(
			`/artists/${params.artistSlug}/releases/${params.releaseSlug}`,
			fetch
		);
		const body = (await res.json()) as ApiEnvelope<PublicRelease>;

		if (!body.data) {
			error(404, 'Release not found');
		}

		return {
			release: body.data
		};
	} catch (caught) {
		if (isHttpError(caught) && caught.status === 404) {
			error(404, 'Release not found');
		}

		throw caught;
	}
};
