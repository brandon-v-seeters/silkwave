export type ApiEnvelope<T> = {
	data?: T;
	error?: {
		code?: string;
		message?: string;
		details?: unknown;
	};
	message?: string;
};

export async function readApiData<T>(response: Response): Promise<T | undefined> {
	const payload = (await response.json()) as (ApiEnvelope<T> & Partial<T>) | T;

	if (payload && typeof payload === 'object' && 'data' in payload) {
		return (payload as ApiEnvelope<T>).data;
	}

	return payload as T;
}
