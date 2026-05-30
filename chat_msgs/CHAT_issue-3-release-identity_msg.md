# 2026-05-12 13:43 CEST - message #0

Placeholder created by overseer. Worker should write message #1 after pre-flight and before code changes.

# 2026-05-12 14:07 CEST - message #1

Pre-flight completed by overseer after the worker stalled and was shut down. No code changes made for issue #3 yet.

## Findings

ADR-0001 and ADR-0002 agree on the target shape:

- Release lifecycle is exactly `draft`, `published`, `archived`.
- There is no `scheduled` status. Scheduling is `published` plus a future `publishAt`.
- There is no `deleted` status. Deletion is a hard delete.
- Release and Artist each need stable external `id` UUIDs.
- Release `Hash` is a lie; the value is a UUID and should be named `id`.
- R2 object keys embed the UUID value, not the field name, so no storage-object migration is needed.

## Exact current fields and constants

`backend/internal/models/release.go` still has:

- `CreateReleaseRequest.ReleaseDate json:"releaseDate"`
- `CreateReleaseRequest.Hash json:"hash"`
- `UpdateReleaseRequest.ReleaseDate json:"releaseDate"`
- `ConfirmDraftReleaseRequest.ReleaseHash json:"releaseHash"`
- `ReleaseStatusScheduled`
- `ReleaseStatusDeleted`
- `ReleaseSchedule.ReleaseDate json:"releaseDate"`
- `Release.Hash json:"hash"`
- `Release.ReleaseDate json:"releaseDate"`
- `Release.Published json:"published"`
- `Release.IsUploaded json:"isUploaded"`
- `CreateDraftResponse.ReleaseHash json:"releaseHash"`
- `TrackUrlDTO.Hash json:"hash"`

`backend/internal/models/track.go` still has:

- `Track.Hash json:"hash"`

`backend/internal/models/artist.go` has no stable external `Id` field.

## Files with relevant references

Core implementation files:

- `backend/internal/models/release.go`
- `backend/internal/models/track.go`
- `backend/internal/models/artist.go`
- `backend/internal/handlers/release_handler.go`
- `backend/internal/handlers/artist_handler.go`
- `backend/internal/repository/release_repository.go`
- `backend/internal/database/schema.go`
- `backend/internal/storage/resolver.go`
- `backend/internal/storage/operations.go`
- `backend/internal/repository/release_repository_test.go`

No frontend or generated TypeScript should be hand-edited. Tygo regeneration is the correct follow-up if the tool exists.

## AQL and index changes needed

`backend/internal/repository/release_repository.go`:

- `FILTER release.published == true` appears in public listing/detail queries.
- `SORT release.releaseDate DESC` appears in public listing queries.
- `FILTER draft.hash == @hash` appears in draft lookup/update/delete queries.
- `UPDATE draft WITH { isUploaded: true, updatedAt: @now }` appears in upload confirmation.
- `UPDATE draft WITH { status: @status, published: @published, updatedAt: @now }` appears in status update.

`backend/internal/database/schema.go`:

- `ReleaseDrafts` has index `hash`.
- `Releases` indexes `artistKey,releaseDate` and `published,releaseDate`.
- `Tracks` indexes `artistKey,releaseDate`, `releaseKey,hash`, and `published,releaseDate`.
- `TrackDrafts` indexes `releaseKey,hash`.

Target index direction:

- `hash` indexes become `id`.
- `releaseDate` indexes become `publishAt` where the concept remains.
- `published,releaseDate` should become lifecycle-aware, likely `status,publishAt`.

## Storage finding

Storage paths use the UUID value in `artist_content/{artistKey}/releases/{releaseHash}/...`.

No object key migration is required. However, code names and comments in `backend/internal/storage/resolver.go` and `backend/internal/storage/operations.go` still use `releaseHash`. Renaming those parameters to `releaseId` is a code-level clarity sweep only. The path segments stay identical because the value stays identical.

## Stale issue-body finding

`Release.IsUploaded` is not unused in the current code.

- `ConfirmDraftRelease` calls `MarkUploadsComplete`.
- `MarkUploadsComplete` writes `isUploaded: true`.
- `PublishRelease` rejects publishing when `!release.IsUploaded`.

Deleting `IsUploaded` without a replacement would remove the current publish-readiness invariant. The likely replacement is repository-owned upload readiness based on `TrackDrafts.uploaded` and the confirmed file paths, but that is a real logic change, not a casual field rename.

## Tygo

`tygo` is not available on PATH:

```text
zsh:1: command not found: tygo
```

`backend/tygo.yaml` has already been corrected by earlier work to point at `../frontend/src/lib/types/generated/models.ts`, but generation cannot run until `tygo` is installed or a repo wrapper appears.

## Recommendation

Do not implement issue #3 as one broad pass.

Split it sequentially:

1. `issue-3a`: mechanical identity/publish naming sweep.
   - Rename Release/Track `Hash` fields to `Id`.
   - Rename `releaseDate` fields and AQL/indexes to `publishAt`.
   - Remove `scheduled` and `deleted` status constants.
   - Add `Artist.Id` and generate it in artist creation.
   - Update repository tests.
   - Rename storage code parameters/comments from `releaseHash` to `releaseId` without changing key shape.
   - Keep `IsUploaded` temporarily and explicitly report it as issue-3b.
2. `issue-3b`: replace upload readiness.
   - Add repository-owned draft upload readiness check from tracks/files.
   - Remove `Release.IsUploaded`.
   - Remove `isUploaded` AQL writes/reads.
   - Update publish path tests once the replacement invariant is explicit.

Reason: `IsUploaded` is live behavior. Bundling its deletion into the naming sweep would turn a mechanical refactor into a hidden lifecycle change. That is how small PRs learn to bite.
