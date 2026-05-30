# CHAT_issue-3-release-identity

## Project context

Silkwave is aligning Release identity and lifecycle naming with ADR-0001 and ADR-0002. Issue #3 is an internal backend foundation slice: `Hash` becomes `Id`, `releaseDate` becomes `publishAt`, redundant visibility booleans disappear, and invalid Release statuses are removed. This slice is intentionally internal. It must not change public URL shape; issue #4 already handled the route parameter label, and issue #10 handles nested public slug URLs.

## Stack location

Likely fair-game after pre-flight confirms exact blast radius:

- `backend/internal/models/release.go`
- `backend/internal/models/track.go`
- `backend/internal/models/artist.go`
- `backend/internal/repository/release_repository.go`
- `backend/internal/handlers/release_handler.go`
- `backend/internal/handlers/artist_handler.go`
- `backend/internal/database/schema.go`
- `backend/internal/storage/**` only if identifiers are named by field rather than value

Read-only references:

- `CONTEXT.md`
- `docs/adr/0001-release-staging-and-lifecycle.md`
- `docs/adr/0002-identifier-conventions.md`
- `CODING_STANDARDS.md`
- GitHub issue #3: "Internal model field renames (Hash->Id, PublishAt, drop legacy booleans)"

## Scope

Pre-flight first. Do not code until the overseer replies.

Target feature scope after approval:

- Rename `Hash` to `Id` on Release and Track model fields where the field is the stable external UUID.
- Add `Id` to Artist and ensure it is generated at create time.
- Remove `Published` and `IsUploaded` from Release if pre-flight confirms replacement logic is clear.
- Remove `scheduled` and `deleted` Release statuses.
- Rename `ReleaseDate` to `PublishAt`.
- Update AQL and bind variables for the renamed model fields.
- Keep R2 object keys semantically unchanged because they embed the UUID value, not the field name.

## File boundaries

Pre-flight fair-game:

- Read any backend files needed.
- Write only `chat_msgs/CHAT_issue-3-release-identity_msg.md`.

No code edits before reply.

Do not touch:

- `frontend/**`
- generated TypeScript
- Stripe, Follow, Subscription, access-module work
- public URL route shape beyond what already exists
- GitHub issue state

Signature-stability rule: do not change handler or repository public signatures during pre-flight. If implementation needs signature changes, report the exact proposed signatures first.

## Coordination rules

- You are not alone in the repo. Release files are already dirty from completed workers. Work with the current state; do not revert or normalize unrelated edits.
- Communicate through `chat_msgs/CHAT_issue-3-release-identity_msg.md`.
- Read `chat_msgs/CHAT_issue-3-release-identity_reply.md` when told `replied`, and verify the reply message number matches your last sent message.

## Hard rules

- Read `CONTEXT.md`, ADR-0001, ADR-0002, and `CODING_STANDARDS.md` before planning edits.
- No new legacy drift: no inline AQL in handlers, no ad-hoc response shapes, no `fmt.Println`, no `any`, no secret logging.
- Repo before handler. If AQL needs to move out of handlers, say where.
- Generated tygo types over hand-rolled mirrors. Do not hand-edit generated frontend types.
- Do not commit.
- Stop and report if the scope expands into lifecycle module extraction, public URL shape, Follow, Subscription, or access checks.

## Acceptance criteria for final implementation

- Release and Track use `Id` instead of `Hash` for the stable external UUID field.
- Artist has an `Id` field and creation path generates it.
- Release no longer uses `Published` or `IsUploaded` if pre-flight confirms safe replacement.
- ReleaseStatus has exactly `draft`, `published`, and `archived`.
- Release publish gate is `PublishAt`, not `ReleaseDate`.
- No AQL query references `release.hash`, `release.published`, `release.isUploaded`, or `release.releaseDate`.
- `go test ./...` and `go build ./...` pass from `backend`.
- Tygo regeneration is attempted if available; otherwise report the exact missing tool or command blocker.

## Pre-flight report requirements

Message #1 must include:

- Exact fields and constants currently present.
- Every backend file that references `Hash`, `ReleaseDate`, `Published`, `IsUploaded`, `scheduled`, or `deleted`.
- AQL queries and bind variables that need edits.
- Whether storage paths depend on field names or only UUID values.
- Whether tygo can run in this environment.
- A recommended implementation split if this is too large for one safe worker pass.

## Effort estimate

Pre-flight: 30-45 minutes. Implementation likely 2-4 hours if the field changes are concentrated in models, repository, and handlers. If the scan finds lifecycle extraction or public URL shape work, stop. That is not issue #3; that is issue #7 or #10 trying on a fake nose.

## Bidirectional message paths

- Worker to overseer: `chat_msgs/CHAT_issue-3-release-identity_msg.md`
- Overseer to worker: `chat_msgs/CHAT_issue-3-release-identity_reply.md`
