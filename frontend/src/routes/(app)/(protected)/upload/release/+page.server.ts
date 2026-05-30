import { GET } from '$lib/api/Api';
import { readApiData } from '$lib/api/envelope';
import type { ReleaseWithTracks } from '$lib/types/generated';
import { redirect } from '@sveltejs/kit';

export const load = async ({ locals, url, fetch }) => {
    const { user } = locals;

    if (!user) {
        throw redirect(307, '/login?redirectTo=/upload/release');
    }

    const draftKey = url.searchParams.get('draftKey');
    let draft: ReleaseWithTracks | null = null;
    if (draftKey) {
        const response = await GET(`/releases/drafts/${draftKey}`, fetch);
        const body = await readApiData<{ draft?: ReleaseWithTracks }>(response);
        draft = body?.draft ?? null;
    }

    return {
        user,
        draft
    };
};
