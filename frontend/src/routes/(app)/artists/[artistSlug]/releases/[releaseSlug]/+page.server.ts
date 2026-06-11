import { loadReleaseByArtistSlug } from '$lib/features/catalog/public-release.server';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ params, fetch }) => {
	return {
		release: await loadReleaseByArtistSlug(params.artistSlug, params.releaseSlug, fetch)
	};
};
