import { redirect } from "@sveltejs/kit";
import { draftsService } from "$lib/services/drafts";
import { GET } from "$lib/api/Api";
import type { Release } from "$lib/types/generated";

export const load = async ({ locals, url, fetch }) => {
    const { user } = locals;

    if (!user) {
        throw redirect(307, '/login?redirectTo=/upload/release');
    }

    const draftKey = url.searchParams.get('draftKey');
    let draft: Release | null = null;
    if (draftKey) {
        const response = await GET(`/releases/drafts/${draftKey}`, fetch);
        draft = response.draft as Release;
    }

    return {
        user,
        draft
    };
};
