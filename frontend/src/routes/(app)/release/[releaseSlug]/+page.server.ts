import { GET } from '$lib/api/Api.js';
import type { PublicRelease, ReleaseWithArtist } from '$lib/types/generated/models.js';
import { error, isHttpError } from '@sveltejs/kit';
import type { RequestEvent } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';

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
	const res = await GET(path, fetch);
	const body = (await res.json()) as ApiEnvelope<PublicRelease>;

	return body.data ?? null;
}

async function loadReleaseThroughLegacyPath(releaseSlug: string, fetch: RequestEvent['fetch']) {
	const res = await GET('/releases', fetch, { limit: 100 });
	const body = (await res.json()) as ApiEnvelope<ReleasesPayload> & ReleasesPayload;
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

export const load: PageServerLoad = async ({ params, fetch }) => {
	try {
		let release: PublicRelease | null = null;

		try {
			release = await loadPublicRelease(`/release/${params.releaseSlug}`, fetch);
		} catch (caught) {
			if (!isHttpError(caught) || caught.status !== 404) {
				throw caught;
			}

			release = await loadReleaseThroughLegacyPath(params.releaseSlug, fetch);
		}

		if (!release) {
			error(404, 'Release not found');
		}

		return {
			release
		};
	} catch (caught) {
		if (isHttpError(caught) && caught.status === 404) {
			error(404, 'Release not found');
		}

		throw caught;
	}
};
