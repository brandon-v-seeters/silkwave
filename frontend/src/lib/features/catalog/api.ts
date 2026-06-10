import type { ApiEnvelope } from '$lib/api/envelope';
import { normalizeCatalogReleases, type CatalogReleaseRow } from './releases';

type CatalogApiError = ApiEnvelope<ReleasesPayload>['error'] | string;

export type ReleasesPayload = {
	releases?: CatalogReleaseRow[];
};

export type CatalogReleasesOptions = {
	limit?: number;
	offset?: number;
	query?: string;
	signal?: AbortSignal;
};

type CatalogFetch = typeof fetch;

export function catalogReleasesApiParams(options: CatalogReleasesOptions = {}) {
	const params: Record<string, string | number> = {
		limit: options.limit ?? 24,
		offset: options.offset ?? 0
	};

	const query = options.query?.trim();
	if (query) {
		params.q = query;
	}

	return params;
}

export function catalogReleasesApiPath(options: CatalogReleasesOptions = {}) {
	const searchParams = new URLSearchParams();

	for (const [key, value] of Object.entries(catalogReleasesApiParams(options))) {
		searchParams.set(key, String(value));
	}

	return `/api/releases?${searchParams.toString()}`;
}

function parseCatalogApiError(error: CatalogApiError | undefined, fallback: string) {
	if (!error) return fallback;
	if (typeof error === 'string') return error;

	return error.message ?? fallback;
}

export async function readCatalogReleasesResponse(response: Response) {
	const payload = (await response.json()) as ApiEnvelope<ReleasesPayload> & ReleasesPayload;

	if (!response.ok || payload.error) {
		throw new Error(parseCatalogApiError(payload.error, 'Failed to load Releases'));
	}

	return normalizeCatalogReleases(payload.data?.releases ?? payload.releases ?? []);
}

export async function fetchCatalogReleases(
	fetcher: CatalogFetch,
	options: CatalogReleasesOptions = {}
) {
	const response = await fetcher(catalogReleasesApiPath(options), {
		signal: options.signal
	});

	return readCatalogReleasesResponse(response);
}
