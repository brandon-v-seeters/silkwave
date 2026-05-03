import { toast } from 'svelte-sonner';
import type { WizardTrack, UploadProgress, UploadStatus } from '$lib/types/WizardTrack';
import type { CreateDraftRequest, CreateDraftResponse, Release } from '$lib/types/generated/models';
import { DELETE, GET, POST } from '$lib/api/Api';


const createDraftsService = () => {
    const getDraft = async (draftKey: string): Promise<Release | null> => {
        try {
            const response = await GET(`/releases/drafts/${draftKey}`);
            return response as Release;
        } catch (err) {
            console.error('Failed to get draft:', err);
            return null;
        }
    };

    const removeDraft = async (draftKey: string) => {
        try {
            const response = await DELETE(`/releases/drafts/${draftKey}`, undefined);
            return response.success;
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
