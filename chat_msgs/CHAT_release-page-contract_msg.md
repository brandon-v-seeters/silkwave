# 2026-05-14 10:26 CEST — message #0

Worker spawned for public Release page data contract. Awaiting first pre-flight plan.

# 2026-05-14 10:28 CEST - message #1

## Reply verification

Read `chat_msgs/CHAT_release-page-contract_reply.md`. It is `reply to message #0`, matching the latest worker message `#0`.

## 5-tier pre-flight

**T1 - domain and ADR.** Read `CONTEXT.md`, ADR-0001, ADR-0002, `CODING_STANDARDS.md`, and the handoff. Public Release URLs must nest under Artist. Public reads must hide Draft, Archived, and future-gated Published Releases. Use Release, Track, Artist, publishAt, id, and slug precisely. No purchase invention. Commerce stays dead until an actual backend contract exists.

**T2 - release surface inventory.** `backend/internal/routes/routes.go` already exposes `GET /api/artists/:artistSlug/releases/:releaseSlug`. `release_handler.go` routes that to `GetReleaseByArtistAndSlug`, which calls `ReleaseRepository.GetPublishedByArtistSlugAndReleaseSlug`. The repository currently joins Artist but not Tracks. That means the public page gets `artist`, but not the ordered Track list it needs.

**T3 - standards drift check.** The touched backend slice already uses `respondOK` / `respondError` for this endpoint and keeps AQL in the repository. I found no need for handler-side DB access, inline AQL, `fmt.Println`, or a custom response envelope. The implementation should stay boring. Boring is good. Boring survives impact.

**T4 - tygo and frontend contract.** Adding a named Go model for the public response would technically change exported Go model shape. The handoff allows `backend/internal/models/release.go`, but says generated `frontend/src/lib/types/generated/models.ts` needs overseer coordination if a Go model changes. Current frontend load already defines a local `PublicRelease = Release & { artist?: Artist; tracks?: Track[] }`, so I can avoid frontend return-shape churn unless the overseer wants generated types refreshed later.

**T5 - DB and query shape.** `Releases` and `Tracks` collections exist. `schema.go` has a unique `artistKey, slug` Release index, plus `status, publishAt`; slug scoping matches ADR-0002. `Track` documents carry `releaseKey`, `title`, `order`, upload flags, file paths, and metadata. The public query can join Tracks by `release._key`, sort by `track.order ASC`, and return real rows.

## Proposed plan

1. Add a bounded public Release response model in `backend/internal/models/release.go`, likely `PublicRelease`, embedding `Release` and adding `Artist *Artist` plus `Tracks []Track`.
2. Change `ReleaseRepository.GetPublishedByArtistSlugAndReleaseSlug` to return that public payload and include ordered Tracks.
3. Keep the visibility gates exactly as-is: `status == published` and `publishAt == null OR publishAt <= now`.
4. Keep `ReleaseHandler.GetReleaseByArtistAndSlug` envelope-shaped with `respondOK`.
5. Update `frontend/src/routes/artist/[artistSlug]/releases/[releaseSlug]/+page.server.ts` only if the backend shape forces typed cleanup. Right now it probably does not.
6. Report missing commerce as follow-up: no clear purchase or checkout endpoint surfaced in this pre-flight.

## Exact files expected

Planned edits:

- `backend/internal/models/release.go`
- `backend/internal/repository/release_repository.go`

Possible edit only if needed:

- `frontend/src/routes/artist/[artistSlug]/releases/[releaseSlug]/+page.server.ts`

Do-not-touch boundaries remain intact: no page Svelte, no media player, no visual components, no generated TS without overseer coordination.

## Validation plan

- `gofmt -w backend/internal/models/release.go backend/internal/repository/release_repository.go`
- `go test ./...` from `backend`
- `go build ./...` from `backend`
- If frontend load changes land, run `npm run check` from `frontend`

No code changes made yet.

# 2026-05-14 11:01 CEST - message #2

## Reply verification

Read `chat_msgs/CHAT_release-page-contract_reply.md`. Latest reply is `reply to message #1`, matching latest worker message `#1`.

## Implemented

Changed files:

- `backend/internal/models/release.go`
- `backend/internal/repository/release_repository.go`

Public Release contract change:

- Added `models.PublicRelease`, embedding `Release` and exposing:
  - `artist?: Artist`
  - `tracks: Track[]`
- Updated `ReleaseRepository.GetPublishedByArtistSlugAndReleaseSlug` to return `*models.PublicRelease`.
- The public query now joins `Tracks` with `track.releaseKey == release._key`.
- Tracks sort by `track.order ASC`.
- Existing visibility gates stayed intact:
  - `release.status == @published`
  - `release.publishAt == null OR release.publishAt <= @now`
- `ReleaseHandler.GetReleaseByArtistAndSlug` stayed boring and unchanged. It still calls the repository and returns through `respondOK`.

No frontend route edit was needed. `frontend/src/routes/artist/[artistSlug]/releases/[releaseSlug]/+page.server.ts` already has a local `PublicRelease = Release & { artist?: Artist; tracks?: Track[] }` type and consumes the envelope shape.

## Dirty-tree caveat

`backend/internal/repository/release_repository.go` is currently untracked in this checkout, so Git cannot show a clean diff against HEAD for that file. I treated it as shared worker output and only touched the public slug lookup query. `backend/internal/models/release.go` already had unrelated dirty changes before this worker; my added block is limited to `PublicRelease`.

## Tygo / generated types

Generated frontend types were not touched, per instruction.

Because `models.PublicRelease` is a new exported Go model, tygo regeneration is probably needed if the team wants generated type parity. It is not required for the current public Release page load because the route already declares the local shape it needs. Still, leaving generated types stale forever would be cosmic paperwork debt. Coordinate tygo outside this worker.

## Validation

Ran:

- `gofmt -w backend/internal/models/release.go backend/internal/repository/release_repository.go`
- `go test ./...` from `backend`
- `go build ./...` from `backend`

Results:

- `go test ./...` passed.
- `go build ./...` passed.

## Commerce follow-up

No clear purchase or checkout endpoint surfaced during this scope. The public Release page can show price data from `release.pricing`, but buying/downloading needs a separate commerce contract decision. I did not invent one.
