// Client-only error types

export class SystemError extends Error {
	constructor(message: string) {
		super(message);
		this.name = 'SystemError';
	}
}

export class UniqueConstraintError extends Error {
	constructor(message: string) {
		super(message);
		this.name = 'UniqueConstraintError';
	}
}

export class UnknownError extends Error {
	constructor(message: string) {
		super(message);
		this.name = 'UnknownError';
	}
}

