import { GET } from '$lib/api/Api';
import {
	catalogReleasesApiParams,
	readCatalogReleasesResponse
} from '$lib/features/catalog/api';

export async function load({ fetch }) {
	try {
		const response = await GET('/releases', fetch, catalogReleasesApiParams({ limit: 12 }));
		const latestReleases = await readCatalogReleasesResponse(response);

		return {
			catalogError: null,
			featuredRelease: latestReleases[0] ?? null,
			latestReleases
		};
	} catch (error) {
		console.error('Failed to load home Releases', error);

		return {
			catalogError: 'Latest Releases are unavailable right now.',
			featuredRelease: null,
			latestReleases: []
		};
	}
}
