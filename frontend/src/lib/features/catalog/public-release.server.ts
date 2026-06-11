import { GET } from '$lib/api/Api.js';
import type { PublicRelease, ReleaseWithArtist } from '$lib/types/generated/models.js';
import { error, isHttpError, type RequestEvent } from '@sveltejs/kit';

type ApiEnvelope<T> = {
	data?: T;
	error?: {
		message?: string;
	};
};

type ReleasesPayload = {
	releases?: ListedRelease[];
};

type ListedRelease = Partial<ReleaseWithArtist> & {
	slug?: string;
	artist?: {
		slug?: string;
	};
};

async function loadPublicRelease(path: string, fetch: RequestEvent['fetch']) {
	const response = await GET(path, fetch);
	const body = (await response.json()) as ApiEnvelope<PublicRelease>;

	return body.data ?? null;
}

function notFound(): never {
	error(404, 'Release not found');
}

export async function loadReleaseByArtistSlug(
	artistSlug: string,
	releaseSlug: string,
	fetch: RequestEvent['fetch']
) {
	try {
		const release = await loadPublicRelease(
			`/artists/${artistSlug}/releases/${releaseSlug}`,
			fetch
		);

		return release ?? notFound();
	} catch (caught) {
		if (isHttpError(caught) && caught.status === 404) {
			notFound();
		}

		throw caught;
	}
}

async function loadReleaseThroughArtistScopedFallback(
	releaseSlug: string,
	fetch: RequestEvent['fetch']
) {
	const response = await GET('/releases', fetch, { limit: 100 });
	const body = (await response.json()) as ApiEnvelope<ReleasesPayload> & ReleasesPayload;
	const releases = body.data?.releases ?? body.releases ?? [];
	const release = releases.find((item) => {
		const itemSlug = item.slug ?? item.Release?.slug;

		return itemSlug === releaseSlug && item.artist?.slug;
	});

	if (!release?.artist?.slug) {
		return null;
	}

	return loadPublicRelease(`/artists/${release.artist.slug}/releases/${releaseSlug}`, fetch);
}

export async function loadReleaseByLegacySlug(
	releaseSlug: string,
	fetch: RequestEvent['fetch']
) {
	try {
		let release: PublicRelease | null = null;

		try {
			release = await loadPublicRelease(`/release/${releaseSlug}`, fetch);
		} catch (caught) {
			if (!isHttpError(caught) || caught.status !== 404) {
				throw caught;
			}

			release = await loadReleaseThroughArtistScopedFallback(releaseSlug, fetch);
		}

		return release ?? notFound();
	} catch (caught) {
		if (isHttpError(caught) && caught.status === 404) {
			notFound();
		}

		throw caught;
	}
}
