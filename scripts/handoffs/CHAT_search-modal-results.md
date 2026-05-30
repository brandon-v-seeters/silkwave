# CHAT_search-modal-results

## Project context

Silkwave has a global search modal in the SvelteKit frontend. It currently contains mock genre, artist, and Release data plus a placeholder for actual results. This worker turns the modal into a real, bounded product surface without inventing a new backend search architecture in the shadows.

## Stack location

- `frontend/src/lib/components/organisms/SearchModal.svelte`
- Read-only references:
  - `frontend/src/routes/(app)/discover/+page.svelte`
  - `frontend/src/lib/components/organisms/MainNavbar.svelte`
  - `frontend/src/lib/api/Api.ts`
  - `backend/internal/handlers/release_handler.go`
  - `backend/internal/repository/release_repository.go`

## Scope

- Replace the mock "Searching for..." placeholder with useful Release search results.
- Reuse the existing public `/api/releases?q=...` capability if it supports the needed data.
- Keep category filtering honest:
  - `all` and `releases` may show Release results from the existing releases endpoint.
  - `artists`, `genres`, and `albums` should not pretend to have first-class backend search unless pre-flight proves they do.
- Remove fake trending artists and fake featured releases if they cannot be backed by real data.
- Preserve the existing keyboard flow: Enter navigates to `/discover?q=...`; Escape closes.
- Keep the modal small and ergonomic. This is search, not a cathedral.

## File boundaries

Fair-game:

- `frontend/src/lib/components/organisms/SearchModal.svelte`

Do not touch:

- Backend files
- Generated types
- `frontend/src/routes/(app)/discover/+page.svelte` unless the worker first reports a strict need and gets a reply
- Existing dirty Release URL files
- `frontend/src/lib/stores/ui.ts` migration unless this file must be touched to keep the modal working

Signature-stability rule: do not change the `SearchModal` exported props (`open`, `searchQuery`) without reporting first.

## Hard rules

- Read `CONTEXT.md`, `docs/adr/0002-identifier-conventions.md`, and `CODING_STANDARDS.md` before code changes.
- Do not use MCP servers, MCP resource discovery, or `npx @sveltejs/mcp`. Validate through local files, existing repo patterns, and repo scripts only.
- Runes only. No new `svelte/store`.
- No `any`.
- No mock production data after this change unless explicitly labeled as an empty-state suggestion.
- Use generated types from `frontend/src/lib/types/generated/models.ts`.
- No backend search API creation in this worker.
- Do not commit.

## Acceptance criteria

- SearchModal no longer shows fake "Searching for..." placeholder content.
- Typing a query shows real Release results or a precise empty state.
- Result links use `/artist/:artistSlug/releases/:releaseSlug`.
- Image fallbacks stay resilient.
- Category controls do not lie about unsupported result types.
- `pnpm check` from `frontend` is attempted. If dependencies are still missing, report the exact blocker.

## Reference reading

- `CONTEXT.md` Identifiers and URLs
- `docs/adr/0002-identifier-conventions.md`
- `CODING_STANDARDS.md` Frontend, Reactivity, TypeScript, Legacy and migration
- `frontend/src/lib/components/organisms/SearchModal.svelte`
- `frontend/src/routes/(app)/discover/+page.svelte`
- `backend/internal/repository/release_repository.go`

## Effort estimate

2-4 hours if the existing releases endpoint is enough. Stop and report if real artist or genre search requires backend work; that is a separate feature, not a hidden dependency.

## Bidirectional message paths

- Worker to overseer: `chat_msgs/CHAT_search-modal-results_msg.md`
- Overseer to worker: `chat_msgs/CHAT_search-modal-results_reply.md`
