# CHAT_adr0002-release-url-frontend

## Project context

Silkwave's public Release URLs must follow ADR-0002: `/artists/:artistSlug/releases/:releaseSlug`. This worker implements the frontend side after the backend route contract is available.

## Stack location

- `frontend/src/routes/artist/[artistSlug]/releases/[releaseSlug]/+page.server.ts` create
- `frontend/src/routes/artist/[artistSlug]/releases/[releaseSlug]/+page.svelte` create
- `frontend/src/routes/(app)/discover/+page.svelte`
- `frontend/src/lib/components/molecules/ReleaseCard.svelte`
- Small related public link helpers only if a local canonical helper already exists

## Scope

- Add a public nested Release detail route under `artist/[artistSlug]/releases/[releaseSlug]`.
- Fetch from `GET /artists/:artistSlug/releases/:releaseSlug` through the existing `GET` helper.
- Update public Release links from `/{artistSlug}/{releaseSlug}` to `/artist/{artistSlug}/releases/{releaseSlug}` unless the route tree proves a different existing convention.
- Keep the page useful but restrained: title, artist link, cover art, basic metadata, tracks if present in the response.
- Do not create purchase, subscribe, checkout, or player behavior in this worker.

## File boundaries

Fair-game:

- `frontend/src/routes/artist/[artistSlug]/releases/[releaseSlug]/+page.server.ts`
- `frontend/src/routes/artist/[artistSlug]/releases/[releaseSlug]/+page.svelte`
- `frontend/src/routes/(app)/discover/+page.svelte`
- `frontend/src/lib/components/molecules/ReleaseCard.svelte`

Do not touch:

- Backend files
- Generated types
- `frontend/src/lib/stores/*` unless a touched file already imports it and a tiny compatibility fix is required
- Subscription UI
- Auth/upload forms

Signature-stability rule: consume the backend route from `CHAT_adr0002-release-url-backend`; do not invent a different endpoint.

## Hard rules

- Use Svelte 5 runes. No new `svelte/store`.
- No `onMount` for route data if a server load can fetch it.
- Use generated types from `frontend/src/lib/types/generated/models.ts`.
- No `any` unless inside existing vendored or explicitly exempt helper code.
- Kebab-case for any new component files. Route filenames are SvelteKit conventions.
- Validate with local repo scripts only. Do not use MCP servers or `npx @sveltejs/mcp`.
- Do not commit.

## Acceptance criteria

- Public Release detail page exists at `/artist/:artistSlug/releases/:releaseSlug`.
- Discover and ReleaseCard links point to that route.
- Page fetches backend detail through server load.
- 404 behavior is clean when backend returns missing Release.
- No generated type edits.
- `pnpm check` passes if dependencies are installed; otherwise report exact blocker.
- Local SvelteKit/TypeScript validation attempted through repo scripts if dependencies are installed.

## Reference reading

- `CONTEXT.md` Identifiers and URLs
- `docs/adr/0002-identifier-conventions.md`
- `CODING_STANDARDS.md` Frontend section
- `frontend/src/routes/artist/[artistSlug]/+layout.server.ts`
- `frontend/src/routes/(app)/discover/+page.svelte`
- `frontend/src/lib/components/molecules/ReleaseCard.svelte`

## Effort estimate

2-5 hours. Wait for or coordinate with the backend worker if the endpoint is not yet wired.

## Bidirectional message paths

- Worker to overseer: `chat_msgs/CHAT_adr0002-release-url-frontend_msg.md`
- Overseer to worker: `chat_msgs/CHAT_adr0002-release-url-frontend_reply.md`
