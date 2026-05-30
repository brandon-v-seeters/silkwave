import { z } from 'zod';
import { error, type RequestEvent } from '@sveltejs/kit';
import { env } from '$env/dynamic/public';
import { dev } from '$app/environment';

const DEV_API_URL = 'http://localhost:8080/api';

export function apiUrl(path: string) {
	const baseUrl = env.PUBLIC_API_URL || (dev ? DEV_API_URL : '');

	if (!baseUrl) {
		throw new Error('PUBLIC_API_URL is not configured');
	}

	return `${baseUrl}${path}`;
}

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

	let res: Response;
	try {
		res = await fetch(apiUrl(path), opts);
	} catch {
		throw error(500, { message: 'Could not connect to server.' });
	}

	if (!res.ok) {
		throw error(res.status, await readErrorBody(res));
	}

	return res;
}

async function readErrorBody(res: Response) {
	const text = await res.text();
	if (!text) {
		return { message: res.statusText || 'Request failed' };
	}

	try {
		return JSON.parse(text);
	} catch {
		return { message: text };
	}
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
		for (let i = 0; i < res.issues.length; i++) {
			const error = res.issues[i];
			const field = error.path[0];
			if (typeof field === 'string') {
				errors[field] = error.message;
			}
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
