import { GET } from '$lib/api/Api';
import { redirect } from '@sveltejs/kit';

export const ssr = true;

export async function load({ locals, fetch }) {
    const { user } = locals;

    if (!user?.artist?._key) {
        throw redirect(307, '/upload');
    }

    const { drafts } = await GET('/releases/drafts', fetch, { artistKey: user.artist._key });

    return {
        drafts: drafts || []
    };
}