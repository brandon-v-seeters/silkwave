// Client-only type for the release upload wizard
// This includes the File object which doesn't exist in the Go model

export type WizardTrack = {
	id: string;
	title: string;
	file: File | null;
	duration: string;
	order: number;
};

export type UploadStatus = 'idle' | 'creating_draft' | 'uploading' | 'confirming' | 'completed' | 'failed';

export type UploadProgress = {
	trackId: string;
	fileName: string;
	progress: number;
	status: 'pending' | 'uploading' | 'completed' | 'failed';
};

