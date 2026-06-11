import { loadReleaseByLegacySlug } from '$lib/features/catalog/public-release.server';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ params, fetch }) => {
	return {
		release: await loadReleaseByLegacySlug(params.releaseSlug, fetch)
	};
};
