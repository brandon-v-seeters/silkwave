import type { AqlQuery } from 'arangojs/aql';
import {
	ARANGO_ROOT_PASSWORD,
	ARANGO_ROOT_USER,
	DATABASE_NAME,
	DATABASE_URL
} from '$env/static/private';
import { Database } from 'arangojs';
import { ArangoError } from 'arangojs/error';

import { SystemError, UniqueConstraintError, UnknownError } from '$lib/types/Error';

// TRANSACTIONS 
//https://arangodb.github.io/arangojs/6.14.1/Reference/Database/Transactions.html

type SingleQuery = {
	query: string;
	bindVars: Record<string, any>; // eslint-disable-line @typescript-eslint/no-explicit-any
};

export type ErrorMessage = {
	body: Record<string, unknown>; // eslint-disable-line @typescript-eslint/no-explicit-any
	message: string;
	error: boolean;
};

export class DatabaseService {
	private static instance: DatabaseService;
	private database: Database;

	private constructor() {
		this.database = new Database({
			url: DATABASE_URL,
			databaseName: DATABASE_NAME,
			auth: {
				username: ARANGO_ROOT_USER,
				password: ARANGO_ROOT_PASSWORD
			}
		});
	}

	public static getInstance(): DatabaseService {
		if (!DatabaseService.instance) {
			DatabaseService.instance = new DatabaseService();
		}

		return DatabaseService.instance;
	}

	public async query(query: AqlQuery) {
		return this.database.query(query);
	}
}

export const query = async <T>(query: SingleQuery): Promise<T | Error> => {
	try {
		const cursor = await DatabaseService.getInstance().query(query);
		const result = await cursor.next();

		if (!result) {
			throw new Error('No result found.');
		}

		return result;
	} catch (err: unknown) {
		return await handleError(err);
	}
};

export const queryMultiple = async <T>(query: SingleQuery): Promise<T[] | Error> => {
	try {
		const cursor = await DatabaseService.getInstance().query(query);

		return await cursor.all();
	} catch (err: unknown) {
		return await handleError(err);
	}
};

export const executeTransaction = async <T>(query: SingleQuery): Promise<T | Error> => {

	try {
		const cursor = await DatabaseService.getInstance().query(query);

		return await cursor.next();
	} catch (err: unknown) {
		return await handleError(err);
	}
};

// eslint-disable-next-line @typescript-eslint/no-explicit-any
const handleError = async (err: any): Promise<Error> => {
	if (err instanceof ArangoError) {
		if (err?.response?.parsedBody?.errorMessage?.includes('unique constraint violated')) {
			const key = handleKeyError(err.response.parsedBody.errorMessage);
			return new UniqueConstraintError(key);
		}
		return err;
	} else if (isSystemError(err)) {
		return new SystemError(err.message);
	} else {
		return new UnknownError('An unknown error occurred.');
	}
};

const handleKeyError = (err: string) => {
	// Regular expression to match content between single quotes
	const regex = /'([^']+)'/g;

	// Extract matches
	const matches = [];
	let match;
	while ((match = regex.exec(err)) !== null) {
		matches.push(match[1]);
	}

	return matches[0];
};

// eslint-disable-next-line @typescript-eslint/no-explicit-any
const isSystemError = (err: any): boolean => {
	return err instanceof Error && err.name === 'SystemError';
};
