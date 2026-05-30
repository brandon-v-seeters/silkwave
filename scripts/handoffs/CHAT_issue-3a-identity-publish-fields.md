# CHAT_issue-3a-identity-publish-fields

## Project context

Silkwave issue #3 is too broad for one safe pass because `Release.IsUploaded` is live publish-readiness behavior. This worker handles only the mechanical identity and publish-date naming slice. A later `issue-3b` slice will replace upload readiness and remove `IsUploaded`.

## Stack location

Fair-game:

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

Do not touch:

- `frontend/**`
- generated TypeScript
- Subscription, Follow, access-module, Stripe, public URL shape, lifecycle extraction
- `Release.IsUploaded` removal or replacement logic
- GitHub issue state

## Scope

Implement the mechanical naming slice:

- Rename Release stable UUID field from `Hash` to `Id` with JSON `id`.
- Rename Track stable UUID field from `Hash` to `Id` with JSON `id`.
- Rename relevant `releaseHash` code variables/DTO fields to `releaseId` where they refer to Release UUID identity.
- Rename `releaseDate` / `ReleaseDate` to `publishAt` / `PublishAt` where it is the publish visibility gate.
- Remove `ReleaseStatusScheduled` and `ReleaseStatusDeleted`.
- Add `Artist.Id` and generate UUID value when creating an Artist.
- Update repository AQL and schema indexes for `id`, `publishAt`, and lifecycle status.
- Rename storage parameter names/comments from `releaseHash` to `releaseId` without changing object key shape.
- Update repository tests for renamed fields.

Explicitly keep:

- `Release.IsUploaded` field.
- `isUploaded` write in upload confirmation.
- `PublishRelease` readiness check using `release.IsUploaded`.

That is intentional. Issue #3b owns replacing that invariant.

## Coordination rules

- You are not alone in the repo. Work with the current dirty files; do not revert or normalize unrelated edits.
- Communicate through `chat_msgs/CHAT_issue-3a-identity-publish-fields_msg.md`.
- Read `chat_msgs/CHAT_issue-3a-identity-publish-fields_reply.md` when told `replied`, and verify the reply message number matches your last sent message.

## Hard rules

- Read `CONTEXT.md`, ADR-0001, ADR-0002, `CODING_STANDARDS.md`, and `chat_msgs/CHAT_issue-3-release-identity_msg.md` message #1 before editing.
- No new legacy drift: no inline AQL in handlers, no ad-hoc response shapes, no `fmt.Println`, no `any`, no secret logging.
- Do not collapse `ReleaseDrafts` / `TrackDrafts`; that is issue #7.
- Do not change public URL route shape.
- Do not hand-edit generated TS.
- Do not commit.
- Stop and report if removing `Published` requires a lifecycle-module extraction.

## Acceptance criteria

- Backend model fields use `Id` / `json:"id"` for Release and Track stable UUIDs.
- Artist model has `Id` / `json:"id"`, and artist registration inserts a UUID.
- Publish visibility fields use `PublishAt` / `json:"publishAt"`.
- `ReleaseStatus` has no `scheduled` or `deleted` constants.
- Public Release listing/detail AQL filters by lifecycle status and publish gate, not `published == true`.
- Draft lookup/update/delete AQL uses `id`, not `hash`.
- Schema indexes use `id`, `publishAt`, and lifecycle status where applicable.
- Storage object keys remain value-identical; only parameter names/comments change.
- `go test ./...` and `go build ./...` pass from `backend`.
- `tygo generate` is attempted and reports exact blocker if unavailable.

## Validation notes

- First sandboxed Go commands may fail on `/Users/brandon/Library/Caches/go-build`; rerun with approved Go build-cache access if needed.
- `tygo` was missing during overseer pre-flight. Try it anyway after code changes and report exact output.

## Bidirectional message paths

- Worker to overseer: `chat_msgs/CHAT_issue-3a-identity-publish-fields_msg.md`
- Overseer to worker: `chat_msgs/CHAT_issue-3a-identity-publish-fields_reply.md`
