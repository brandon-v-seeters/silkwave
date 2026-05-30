# Release Page Contract Worker

Project context: Silkwave is a music distribution platform where independent Artists publish Releases with Tracks, pricing, cover art, and public URLs. The user wants the public Release page, where a listener lands and can eventually buy the Release, pushed toward a Figma direction: dark immersive cover-art background, prominent cover art, Release title, Artist link, price CTA, playable track list, recommendations, and a bottom player.

Stack location: backend public Release endpoint plus the SvelteKit server load for `frontend/src/routes/artist/[artistSlug]/releases/[releaseSlug]/+page.server.ts`.

## Scope

Own the data contract for the public Release page.

- Verify `GET /api/artists/:artistSlug/releases/:releaseSlug` returns the data the page needs.
- If safe and bounded, make the backend return a public Release payload including Artist and ordered Tracks.
- Keep visibility gates correct: only Published Releases whose `publishAt` is absent or in the past.
- Keep the response envelope shape through `respondOK` / `respondError`.
- Update the SvelteKit server load only if it needs typed envelope handling or mapping cleanup.
- Do not implement purchases yet unless the existing backend already has a clear purchase endpoint. If not, report the missing commerce contract as a follow-up decision.

## File Boundaries

Your fair-game files:

- `backend/internal/repository/release_repository.go`
- `backend/internal/handlers/release_handler.go`
- `backend/internal/models/release.go`
- `frontend/src/routes/artist/[artistSlug]/releases/[releaseSlug]/+page.server.ts`
- Backend tests only if you add or change repository behavior.

DO NOT touch:

- `frontend/src/routes/artist/[artistSlug]/releases/[releaseSlug]/+page.svelte`
- `frontend/src/lib/components/organisms/media-player/MediaPlayer.svelte`
- New public Release page visual components. The UI worker owns those.
- Generated `frontend/src/lib/types/generated/models.ts` unless a Go model changes and the overseer explicitly coordinates tygo regeneration.

Signature-stability rule: if you change exported Go model shapes or frontend load return shape, write a message before broadening the change. The UI worker will depend on stable `release`, `artist`, `tracks`, `coverArt`, `price`, and `publishAt` fields.

## Coordination

Message file: `chat_msgs/CHAT_release-page-contract_msg.md`
Reply file: `chat_msgs/CHAT_release-page-contract_reply.md`

Write milestone decisions to the message file only. On `replied`, read the reply file and verify the header says `reply to message #N` for your latest message.

## Hard Rules

- `CONTEXT.md` is the source of truth: say Release, Track, Artist, Follow, Subscription. Avoid using "album" as a synonym for Release.
- Public release URLs nest under Artist: `/artist/:artistSlug/releases/:releaseSlug` in frontend, `/api/artists/:artistSlug/releases/:releaseSlug` in backend.
- No new legacy drift: no inline AQL in handlers, no `any`, no `fmt.Println`, no ad-hoc response shapes, no `Hash`, no `releaseDate`, no `unpublish`.
- Repo before handler. Handlers never touch ArangoDB directly.
- Preserve `/*aql*/` comments.
- No MCP servers. Use local files, `rg`, and repo scripts.
- Do not commit.

## Acceptance Criteria

- Public Release endpoint supplies enough data for the page shell to render the hero and ordered track list without fake track titles when real tracks exist.
- Archived, Draft, and future-gated Published Releases remain hidden from public reads.
- API responses stay envelope-shaped.
- If purchase/checkout data is missing, report the exact missing contract instead of inventing a fake purchase system.
- Run focused validation. At minimum, run Go tests/build for touched backend packages if backend changes land. If frontend load changes land, run `npm run check` from `frontend` if dependencies are usable.

## Reference Reading

- `CONTEXT.md`
- `docs/adr/0001-release-staging-and-lifecycle.md`
- `docs/adr/0002-identifier-conventions.md`
- `CODING_STANDARDS.md`
- `backend/internal/repository/release_repository.go`
- `backend/internal/handlers/release_handler.go`
- `frontend/src/routes/artist/[artistSlug]/releases/[releaseSlug]/+page.server.ts`

## Effort Estimate

Pre-flight first. Expected scope: small-to-medium if tracks are already joinable from `Tracks`; larger only if public streaming URLs or commerce shape are missing. Stop and report before turning this into a purchase architecture project.

