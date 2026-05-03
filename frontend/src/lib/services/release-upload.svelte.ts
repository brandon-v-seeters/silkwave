import { toast } from 'svelte-sonner';
import type { WizardTrack, UploadProgress, UploadStatus } from '$lib/types/WizardTrack';
import type { CreateDraftRequest, CreateDraftResponse } from '$lib/types/generated/models';
import { POST } from '$lib/api/Api';


const createReleaseUploadService = () => {
    // Reactive state using Svelte 5 runes
    let status = $state<UploadStatus>('idle');
    let releaseHash = $state<string | null>(null);
    let releaseKey = $state<string | null>(null);
    let coverArtProgress = $state<UploadProgress | null>(null);
    let trackProgress = $state<UploadProgress[]>([]);
    let error = $state<string | null>(null);

    // Derived states
    const isUploading = $derived(
        status === 'creating_draft' || status === 'uploading' || status === 'confirming'
    );

    const overallProgress = $derived(() => {
        if (trackProgress.length === 0) return 0;
        const total = trackProgress.reduce((sum, t) => sum + t.progress, 0);
        const coverProgress = coverArtProgress?.progress ?? 0;
        const count = trackProgress.length + (coverArtProgress ? 1 : 0);
        return Math.round((total + coverProgress) / count);
    });

    const completedCount = $derived(
        trackProgress.filter((t) => t.status === 'completed').length +
        (coverArtProgress?.status === 'completed' ? 1 : 0)
    );

    const totalCount = $derived(trackProgress.length + (coverArtProgress ? 1 : 0));

    // Reset store to initial state
    const reset = () => {
        status = 'idle';
        releaseHash = null;
        releaseKey = null;
        coverArtProgress = null;
        trackProgress = [];
        error = null;
    };

    // Upload a single file to R2
    const uploadFileToR2 = async (
        file: File,
        presignedUrl: string,
        onProgress?: (progress: number) => void
    ): Promise<boolean> => {
        console.log('Uploading to presigned URL:', presignedUrl);

        try {
            const response = await fetch(presignedUrl, {
                method: 'PUT',
                body: file,
                headers: {
                    'Content-Type': file.type,
                    'x-amz-acl': 'public-read'
                }
            });

            if (response.ok) {
                onProgress?.(100);
                return true;
            } else {
                console.error('Upload failed:', response.status, response.statusText);
                return false;
            }
        } catch (err) {
            console.error('Upload error:', err);
            return false;
        }
    };

    // Create release draft in backend
    const createDraft = async (
        releaseTitle: string,
        genres: string[],
        tracks: WizardTrack[],
        coverArt: File | null,
        artistKey: string
    ): Promise<CreateDraftResponse | null> => {
        const request: CreateDraftRequest = {
            artistKey: artistKey,
            title: releaseTitle,
            genres,
            tracks: tracks.map((track, index) => ({
                title: track.title,
                fileName: track.file?.name || '',
                fileType: track.file?.type || 'audio/mpeg',
                fileSize: track.file?.size || 0,
                duration: track.duration,
                order: index + 1
            })),
            coverArt: coverArt
                ? {
                    fileName: coverArt.name,
                    fileType: coverArt.type,
                    fileSize: coverArt.size
                }
                : undefined
        };

        try {
            return await POST('/releases/draft', undefined, request) as CreateDraftResponse;
        } catch (err) {
            console.error('Failed to create release draft:', err);
            toast.error(err instanceof Error ? err.message : 'Failed to create release draft');
            throw err;
        }
    };

    // Confirm uploads to backend
    const confirmUploads = async (
        hash: string,
        uploadedTracks: { trackId: string; storagePath: string }[],
        coverArtPath: string | null
    ): Promise<boolean> => {
        try {
            const response = await POST(`/releases/${hash}/confirm`, undefined, { tracks: uploadedTracks, coverArtPath });
            return response.success;
        } catch (err) {
            console.error('Failed to confirm uploads:', err);
            toast.error('Failed to confirm uploads');
            return false;
        }
    };

    // Upload files to R2 and confirm the release
    const uploadFilesAndConfirm = async (
        draftResponse: CreateDraftResponse,
        tracks: WizardTrack[],
        coverArt: File | null
    ): Promise<{ success: boolean; releaseKey: string | null }> => {

        if (!draftResponse) return { success: false, releaseKey: null };

        // Initialize progress tracking
        trackProgress = draftResponse.presignedUrls.tracks.map((t: { hash: string; fileName: string }) => ({
            trackId: t.hash,
            fileName: t.fileName,
            progress: 0,
            status: 'pending' as const
        }));
        coverArtProgress = coverArt
            ? {
                trackId: 'cover-art',
                fileName: coverArt.name,
                progress: 0,
                status: 'pending' as const
            }
            : null;

        toast.dismiss();
        toast.loading('Uploading files...');

        // Upload files to R2
        const uploadPromises: Promise<{ id: string; success: boolean; storagePath: string }>[] = [];

        // Upload cover art
        if (coverArt && draftResponse.presignedUrls.coverArt) {
            const coverPromise = uploadFileToR2(
                coverArt,
                draftResponse.presignedUrls.coverArt,
                (progress) => {
                    coverArtProgress = {
                        ...coverArtProgress!,
                        progress,
                        status: progress === 100 ? 'completed' : 'uploading'
                    };
                }
            ).then((success) => ({
                id: 'cover-art',
                success,
                storagePath: `artist_content/${draftResponse.artistKey}/releases/${draftResponse.releaseHash}/draft/cover.${coverArt.name.split('.').pop()}`
            }));

            uploadPromises.push(coverPromise);
        }

        // Upload tracks
        for (const trackUrl of draftResponse.presignedUrls.tracks) {
            const track = tracks.find(
                (t) => t.file?.name === trackUrl.fileName || t.id === trackUrl.hash
            );
            if (!track?.file) continue;

            const trackFile = track.file;
            const trackPromise = uploadFileToR2(trackFile, trackUrl.presignedUrl, (progress) => {
                trackProgress = trackProgress.map((tp) =>
                    tp.trackId === trackUrl.hash
                        ? { ...tp, progress, status: progress === 100 ? 'completed' : 'uploading' }
                        : tp
                );
            }).then((success) => ({
                id: trackUrl.hash,
                success,
                storagePath: trackUrl.storagePath
            })).catch((err) => {
                console.error('Failed to upload track:', err);
                return { id: trackUrl.hash, success: false, storagePath: trackUrl.storagePath };
            });

            uploadPromises.push(trackPromise);
        }

        const results = await Promise.all(uploadPromises);
        const failedUploads = results.filter((r) => !r.success);

        if (failedUploads.length > 0) {
            status = 'failed';
            error = `Failed to upload ${failedUploads.length} file(s)`;
            trackProgress = trackProgress.map((tp) =>
                failedUploads.some((f) => f.id === tp.trackId) ? { ...tp, status: 'failed' } : tp
            );
            toast.dismiss();
            toast.error(`Failed to upload ${failedUploads.length} file(s)`);
            return { success: false, releaseKey: draftResponse.releaseKey };
        }

        // Confirm uploads
        status = 'confirming';
        toast.dismiss();
        toast.loading('Finalizing...');

        const uploadedTracks = results
            .filter((r) => r.id !== 'cover-art')
            .map((r) => ({ trackId: r.id, storagePath: r.storagePath }));

        const coverArtResult = results.find((r) => r.id === 'cover-art');

        const confirmed = await confirmUploads(
            draftResponse.releaseHash,
            uploadedTracks,
            coverArtResult?.storagePath || null
        );

        if (!confirmed) {
            status = 'failed';
            error = 'Failed to confirm uploads';
            toast.dismiss();
            toast.error('Failed to finalize release');
            return { success: false, releaseKey: draftResponse.releaseKey };
        }

        status = 'completed';
        toast.dismiss();
        toast.success('Release saved as draft!');

        return { success: true, releaseKey: draftResponse.releaseKey };
    };

    // Main upload function
    const saveAsDraft = async (
        releaseTitle: string,
        genres: string[],
        tracks: WizardTrack[],
        coverArt: File | null,
        artistKey: string
    ): Promise<{ success: boolean; releaseKey: string | null }> => {
        if (isUploading) {
            return { success: false, releaseKey: null };
        }

        try {
            // Step 1: Create draft and get presigned URLs
            status = 'creating_draft';
            error = null;
            toast.loading('Creating release draft...');

            const draftResponse = await createDraft(releaseTitle, genres, tracks, coverArt, artistKey);

            if (!draftResponse) {
                console.log('Failed to create draft');
                status = 'failed';
                error = 'Failed to create draft';
                toast.dismiss();
                return { success: false, releaseKey: null };
            }

            // Store release identifiers in memory (persists even if upload fails)
            releaseKey = draftResponse.releaseKey;
            releaseHash = draftResponse.releaseHash;

            // Upload files and confirm
            status = 'uploading';
            return await uploadFilesAndConfirm(draftResponse, tracks, coverArt);
        } catch (err) {
            console.error('Release upload failed:', err);
            status = 'failed';
            error = err instanceof Error ? err.message : 'Unknown error';
            toast.dismiss();
            toast.error('Release upload failed');
            return { success: false, releaseKey: null };
        }
    };

    // Return reactive getters and actions
    return {
        // Reactive state (getters)
        get status() {
            return status;
        },
        get releaseHash() {
            return releaseHash;
        },
        get releaseKey() {
            return releaseKey;
        },
        get coverArtProgress() {
            return coverArtProgress;
        },
        get trackProgress() {
            return trackProgress;
        },
        get error() {
            return error;
        },
        get isUploading() {
            return isUploading;
        },
        get overallProgress() {
            return overallProgress;
        },
        get completedCount() {
            return completedCount;
        },
        get totalCount() {
            return totalCount;
        },

        // Actions
        reset,
        saveAsDraft
    };
};

// Export singleton instance
export const releaseUploadService = createReleaseUploadService();
