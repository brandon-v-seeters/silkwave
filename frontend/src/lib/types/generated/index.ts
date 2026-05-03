// Auto-generated types from Go backend
// Run `npm run generate-types` to regenerate

export * from './models';

// Re-export with flattened types for frontend convenience
// Go struct embedding creates nested types, but we want flat access
import type { ClientUser as GeneratedClientUser, User as GeneratedUser, Artist } from './models';

/**
 * Flattened ClientUser type for frontend use
 * Combines User properties with ClientUser relations
 */
export type AppUser = GeneratedUser & {
	subscribedTo: GeneratedClientUser['subscribedTo'];
	artist?: Artist;
};

