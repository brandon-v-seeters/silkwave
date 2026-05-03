import { z } from 'zod';
import { error, type RequestEvent } from '@sveltejs/kit';
import { PUBLIC_API_URL } from '$env/static/public';

async function useFetch(
	fetch: RequestEvent['fetch'] = globalThis.fetch,
	{
		method,
		path,
		data,
		token
	}: {
		method: string;
		path: string;
		data?: any;
		token?: string;
	}
) {
	const opts: RequestInit = {
		method,
		credentials: 'include',
		headers: {
			Authorization: `Token ${token}`,
			'Content-Type': 'application/json'
		}
	};

	if (data) {
		opts.body = JSON.stringify(data);
	}

	if (!token) {
		delete opts.headers?.['Authorization' as keyof HeadersInit];
	}

	const res = await fetch(`${PUBLIC_API_URL}${path}`, opts).catch((err) => err);

	if (res.message === 'fetch failed') {
		throw error(500, 'Could not connect to server...');
	}

	if (!res.ok) {
		throw error(res.status, await res.json());
	}

	return res;
}

export function GET(path: string, fetch?: RequestEvent['fetch'], params?: Record<string, string | number | boolean>, token?: string) {
	if (params) {
		const searchParams = new URLSearchParams();

		for (const [key, value] of Object.entries(params)) {
			if (value !== undefined && value !== null) {
				searchParams.append(key, String(value));
			}
		}

		const queryString = searchParams.toString();

		if (queryString) {
			path = `${path}?${queryString}`;
		}
	}

	return useFetch(fetch, { method: 'GET', path, token });
}

export async function DELETE(path: string, fetch?: RequestEvent['fetch'], token?: string) {
	return await useFetch(fetch, { method: 'DELETE', path, token });
}

export function POST(path: string, fetch?: RequestEvent['fetch'], data?: any, token?: string) {
	return useFetch(fetch, { method: 'POST', path, data, token });
}

export function PUT(path: string, fetch?: RequestEvent['fetch'], data?: any, token?: string) {
	return useFetch(fetch, { method: 'PUT', path, data, token });
}

export const handleError = (res: any, title: string, errors?: Record<string, string>) => {
	let errtext: string = '';
	let fieldCount = 0;

	if (res instanceof z.ZodError && typeof errors === 'object') {
		for (let i = 0; i < res.errors.length; i++) {
			const error = res.errors[i];
			errors[error.path[0]] = error.message;
		}
		errors = errors;
		return;
	}

	switch (res.status) {
		case 500: {
			break;
		}
		case 404: {
			break;
		}
		case 400: {
			if (res.body.message.includes("'required' tag")) {
				if (!errors) {
					return;
				}
				const field = res.body.message.split('Error:')[1].split("'")[1].split("'")[0].toLowerCase();
				errors[field] = `Please fill in a valid ${field}`;
				errors = errors;
			}
			break;
		}
		default:
	}
	return;
};
