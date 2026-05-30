# CHAT_issue-5-userartist-role

## Project context

Silkwave models Artist ownership as a binary User to Artist relationship. ADR-0003 rejects team roles in v1. Issue #5 removes the stale `Role` field from the `UserArtist` edge so future code does not accidentally build role-aware behavior into a product that explicitly does not support it.

## Stack location

- Primary file: `backend/internal/models/edges.go`
- Read-only references:
  - `CONTEXT.md`
  - `docs/adr/0003-single-operator-artist-ownership.md`
  - `CODING_STANDARDS.md`
  - `backend/internal/middleware/auth.go`
  - `backend/internal/repository/user_repository.go`
  - `backend/internal/database/schema.go`

## Scope

- Remove `Role` from `models.UserArtist`.
- Search for AQL or Go references to `userArtist.role`, `Role` on `UserArtist`, or role-aware artist ownership logic.
- Keep `middleware.UserManagesArtist` binary: user manages artist, or does not.
- Add or adjust focused tests only if existing tests expose a direct seam for this behavior.
- Do not broaden into Follow, Subscription, Release lifecycle, slug routing, or generated type work.

## File boundaries

Fair-game:

- `backend/internal/models/edges.go`

Conditionally fair-game after pre-flight report or obvious direct need:

- `backend/internal/middleware/auth.go`
- `backend/internal/repository/user_repository.go`
- Existing focused test files under `backend/internal/**`

Do not touch:

- `backend/internal/database/schema.go`
- `backend/internal/models/subscription.go`
- `backend/internal/models/release.go`
- `backend/internal/handlers/release_handler.go`
- `backend/internal/routes/routes.go`
- Any frontend files
- Generated TypeScript
- Existing dirty files outside the exact role-removal path

Signature-stability rule: do not change public middleware or repository method signatures unless pre-flight proves the old signature encodes `Role`. Report first if that happens.

## Coordination rules

- You are not alone in the repo. Other workers have dirty changes in release, subscription, generated-type, and search surfaces.
- Do not revert, normalize, format, or "clean up" files outside this handoff.
- Communicate through `chat_msgs/CHAT_issue-5-userartist-role_msg.md`.
- Read `chat_msgs/CHAT_issue-5-userartist-role_reply.md` when told `replied`, and verify the reply message number matches your last sent message.

## Hard rules

- Read `CONTEXT.md`, ADR-0003, and `CODING_STANDARDS.md` before code changes.
- No new legacy drift: no inline AQL in handlers, no ad-hoc response shapes, no `fmt.Println`, no `any`, no secret logging.
- Handlers do not touch ArangoDB directly.
- Generated tygo types over hand-rolled mirrors. Do not hand-edit generated frontend types.
- No MCP servers, MCP resource discovery, `svelte-autofixer`, or `npx @sveltejs/mcp`.
- Do not commit.
- Stop and report if scope expands beyond removing the stale edge field and direct references.

## Acceptance criteria

- `models.UserArtist` has no `Role` field.
- No AQL query references `role` on a `UserArtist` document.
- `middleware.UserManagesArtist` remains binary and still compiles.
- `go test ./...` from `backend` is attempted.
- `go build ./...` from `backend` is attempted.
- Any validation blocker reports the exact command and output.

## Reference reading

- `CONTEXT.md` Ownership and Relationships sections
- `docs/adr/0003-single-operator-artist-ownership.md`
- `CODING_STANDARDS.md`
- Issue #5: "Drop unused Role field from UserArtist edge"
- `backend/internal/models/edges.go`
- `backend/internal/middleware/auth.go`
- `backend/internal/repository/user_repository.go`

## Effort estimate

30-60 minutes if the field is truly unused. Stop and report if pre-flight finds persisted role-dependent logic. That would mean the issue body is stale, and stale requirements deserve daylight, not stealth surgery.

## Bidirectional message paths

- Worker to overseer: `chat_msgs/CHAT_issue-5-userartist-role_msg.md`
- Overseer to worker: `chat_msgs/CHAT_issue-5-userartist-role_reply.md`
