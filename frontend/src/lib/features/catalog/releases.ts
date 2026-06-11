import type { Artist, Release, ReleaseWithArtist } from '$lib/types/generated/models';

export type CatalogArtist = Pick<Artist, 'name' | 'slug'>;

export type FlatCatalogRelease = Release & {
	artist?: CatalogArtist;
	coverArt?: string | null;
	type?: string;
};

export type CatalogReleaseRow = FlatCatalogRelease | ReleaseWithArtist;

export type CatalogRelease = FlatCatalogRelease & {
	artist?: CatalogArtist;
	coverArtUrl: string | null;
	key: string;
	kind: string;
	publishedAtLabel: string;
};

export const artistRoute = '/(app)/artist/[artistSlug]';
export const releaseRoute = '/(app)/artists/[artistSlug]/releases/[releaseSlug]';

type LinkableRelease = Pick<FlatCatalogRelease, 'slug'> & {
	artist: CatalogArtist;
};

function hasNestedRelease(row: CatalogReleaseRow): row is ReleaseWithArtist {
	return 'Release' in row;
}

function compactArtist(artist: Artist | CatalogArtist | undefined): CatalogArtist | undefined {
	if (!artist) return undefined;

	return {
		name: artist.name,
		slug: artist.slug
	};
}

export function coverArtFor(release: FlatCatalogRelease) {
	return (
		release.coverArt ||
		release.cover ||
		release.assets?.coverArt?.medium ||
		release.assets?.coverArt?.original ||
		release.assets?.coverArt?.thumbnail ||
		null
	);
}

export function releaseKind(release: FlatCatalogRelease) {
	return release.type ?? release.releaseType ?? 'Release';
}

export function releaseKey(release: FlatCatalogRelease) {
	return release.id || release._key || `${release.artist?.slug ?? 'artist'}-${release.slug}`;
}

export function hasReleaseRoute(
	release: Pick<FlatCatalogRelease, 'slug'> & { artist?: CatalogArtist }
): release is LinkableRelease {
	return Boolean(release.slug && release.artist?.slug);
}

export function releaseRouteParams(release: LinkableRelease) {
	return {
		artistSlug: release.artist.slug,
		releaseSlug: release.slug
	};
}

export function artistRouteParams(artist: CatalogArtist) {
	return { artistSlug: artist.slug };
}

export function formatReleaseDate(value: string | number | undefined) {
	if (!value) return 'Unknown';

	const date = new Date(value);
	if (Number.isNaN(date.getTime())) return 'Unknown';

	return date.toLocaleDateString('en-US', {
		year: 'numeric',
		month: 'long',
		day: 'numeric'
	});
}

export function normalizeCatalogRelease(row: CatalogReleaseRow): CatalogRelease {
	const release: FlatCatalogRelease = hasNestedRelease(row)
		? {
				...row.Release,
				artist: compactArtist(row.artist)
			}
		: {
				...row,
				artist: compactArtist(row.artist)
			};

	return {
		...release,
		coverArtUrl: coverArtFor(release),
		key: releaseKey(release),
		kind: releaseKind(release),
		publishedAtLabel: formatReleaseDate(release.publishAt)
	};
}

export function normalizeCatalogReleases(rows: CatalogReleaseRow[] = []) {
	return rows.map(normalizeCatalogRelease);
}
