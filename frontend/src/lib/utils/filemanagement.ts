import { toast } from 'svelte-sonner';

type PATHS = 'images/' | 'audio/';

const IMAGE_PATHS: Record<string, PATHS> = {
	image: 'images/',
	audio: 'audio/'
};

export const handleFileUpload = async (file: File, url: string) => {
	if (!file) {
		toast.error('No file selected');
		return;
	}

	const res = await fetch(url, {
		method: 'PUT',
		headers: {
			'Content-Type': file.type
		},
		body: file
	});

	if (!res.ok) {
		console.error('Failed to upload file to R2');
	}

	return res;
};

export const handleGetPresignedUrl = async (key: string) => {
	const type = getFileType(key).split('/')[0];

	if (!(type in IMAGE_PATHS)) {
		toast.error('Could not get correct path due to invalid file type');
		return;
	}

	const path = IMAGE_PATHS[type];

	const res = await fetch(`/api/upload`, {
		method: 'POST',
		body: JSON.stringify({
			fileName: `${key}`,
			fileType: type,
			path
		}),
		headers: {
			'Content-Type': 'application/json'
		}
	});

	if (!res.ok) {
		toast.error('Could not get presigned URL');
		return;
	}

	const { presignedUrl } = await res.json();
	return presignedUrl;
};

const getFileType = (key: string) => {
	const ext = key.split('.').pop();

	switch (ext) {
		case 'jpg':
		case 'jpeg':
			return 'image/jpeg';
		case 'png':
			return 'image/png';
		case 'gif':
			return 'image/gif';
		case 'svg':
			return 'image/svg+xml';
		default:
			return 'application/octet-stream';
	}
};


export const generateHash = () => {
	return window.crypto.randomUUID();
}