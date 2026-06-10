import { toast } from 'svelte-sonner';
import { readApiData } from '$lib/api/envelope';
import type { ReleaseWithTracks } from '$lib/types/generated/models';
import { DELETE, GET } from '$lib/api/Api';


const createDraftsService = () => {
    const getDraft = async (draftKey: string): Promise<ReleaseWithTracks | null> => {
        try {
            const response = await GET(`/releases/drafts/${draftKey}`);
            const body = await readApiData<{ draft?: ReleaseWithTracks }>(response);
            return body?.draft ?? null;
        } catch (err) {
            console.error('Failed to get draft:', err);
            return null;
        }
    };

    const removeDraft = async (draftKey: string) => {
        try {
            const response = await DELETE(`/releases/drafts/${draftKey}`, undefined);
            return response.ok;
        } catch (err) {
            console.error('Failed to remove draft:', err);
            toast.error('Failed to remove draft');
            return false;
        }
    };

    return {
        getDraft,
        removeDraft
    };
};

// Export singleton instance
export const draftsService = createDraftsService();
