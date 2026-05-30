# CHAT_issue-4-release-id-param

## Project context

Silkwave is aligning public API language with ADR-0002. Issue #4 is the narrow cosmetic backend slice that renames Release route path parameters from `releaseHash` to `releaseId`. This does not rename the underlying model fields, DTO fields, storage keys, or response JSON. That heavier identity work belongs to issue #3.

## Stack location

- `backend/internal/routes/routes.go`
- `backend/internal/handlers/release_handler.go`

Read-only references:

- `CONTEXT.md`
- `docs/adr/0002-identifier-conventions.md`
- `CODING_STANDARDS.md`
- GitHub issue #4: "Rename API path parameter :releaseHash -> :releaseId"

## Scope

- Change protected Release route definitions from `:releaseHash` to `:releaseId`.
- Update matching `c.Param("releaseHash")` calls to `c.Param("releaseId")`.
- Rename local variables that only hold the route parameter from `releaseHash` to `releaseId`.
- Keep generated UUID creation variables alone when they represent newly created Release UUID values, unless pre-flight proves they are route parameters.
- Keep external response JSON and model fields unchanged unless pre-flight proves issue #4 explicitly requires it. It does not.

## File boundaries

Fair-game:

- `backend/internal/routes/routes.go`
- `backend/internal/handlers/release_handler.go`

Do not touch:

- `backend/internal/models/release.go`
- `backend/internal/models/subscription.go`
- `backend/internal/database/schema.go`
- `backend/internal/repository/**`
- `frontend/**`
- generated TypeScript
- `chat_msgs/` except your own message file

Signature-stability rule: do not change handler method names or constructor signatures. This is parameter naming only.

## Coordination rules

- You are not alone in the repo. Release handler and route files are already dirty from completed workers. Work with the current state; do not revert or normalize unrelated edits.
- Communicate through `chat_msgs/CHAT_issue-4-release-id-param_msg.md`.
- Read `chat_msgs/CHAT_issue-4-release-id-param_reply.md` when told `replied`, and verify the reply message number matches your last sent message.

## Hard rules

- Read `CONTEXT.md`, ADR-0002, and `CODING_STANDARDS.md` before code changes.
- No new legacy drift: no inline AQL in handlers, no ad-hoc response shapes, no `fmt.Println`, no `any`, no secret logging.
- Do not rename model fields from `Hash` to `Id`; that is issue #3.
- Do not alter storage key naming or R2 paths.
- Do not commit.
- Stop and report if pre-flight shows the route rename is entangled with issue #3.

## Acceptance criteria

- `backend/internal/routes/routes.go` has no `:releaseHash` route parameter.
- `backend/internal/handlers/release_handler.go` has no `c.Param("releaseHash")`.
- Local variables that only hold the route path parameter are named `releaseId`.
- `go test ./...` from `backend` is attempted.
- `go build ./...` from `backend` is attempted.
- If a local backend server is not running for curl validation, report that exact blocker instead of inventing a fake curl result.

## Reference reading

- `CONTEXT.md` Identifiers and URLs
- `docs/adr/0002-identifier-conventions.md`
- `CODING_STANDARDS.md`
- `backend/internal/routes/routes.go`
- `backend/internal/handlers/release_handler.go`

## Effort estimate

30-60 minutes if the rename stays cosmetic. Stop and report if route parameter naming pulls on model identity, storage paths, or generated types. That would mean issue #4 is hiding issue #3 under a cheap mask.

## Bidirectional message paths

- Worker to overseer: `chat_msgs/CHAT_issue-4-release-id-param_msg.md`
- Overseer to worker: `chat_msgs/CHAT_issue-4-release-id-param_reply.md`
