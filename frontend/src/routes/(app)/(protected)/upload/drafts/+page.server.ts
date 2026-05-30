import { GET } from '$lib/api/Api';
import { readApiData } from '$lib/api/envelope';
import type { ReleaseWithArtist } from '$lib/types/generated/models';
import { redirect } from '@sveltejs/kit';

export const ssr = true;

export async function load({ locals, fetch }) {
    const { user } = locals;

    if (!user?.artist?._key) {
        throw redirect(307, '/upload');
    }

    const response = await GET('/releases/drafts', fetch, { artistKey: user.artist._key });
    const body = await readApiData<{ drafts?: ReleaseWithArtist[] }>(response);

    return {
        drafts: body?.drafts ?? []
    };
}
