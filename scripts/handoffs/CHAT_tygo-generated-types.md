# CHAT_tygo-generated-types

## Project context

Silkwave uses Go models as the source of truth and generated TypeScript types in `frontend/src/lib/types/generated/models.ts`. A completed subscription-semantics worker found the tygo config still writes to an old path: `../silk-wave/src/lib/types/generated/models.ts`. This makes generated frontend types stale, which is how type systems quietly become decorative wallpaper.

## Stack location

- `backend/tygo.yaml`
- `frontend/src/lib/types/generated/models.ts`
- Read-only references:
  - `backend/internal/models/*.go`
  - `frontend/src/lib/types/generated/index.ts`
  - `CODING_STANDARDS.md`

## Scope

- Fix tygo output configuration so generation targets `frontend/src/lib/types/generated/models.ts`.
- Identify the exact local generation command for this checkout.
- Regenerate the checked-in generated models file only if the tool is available locally.
- If tygo is missing and would require network install, stop and report the exact command/blocker. Do not hand-edit generated output to "look right." That is forgery with syntax highlighting.

## File boundaries

Fair-game:

- `backend/tygo.yaml`
- `frontend/src/lib/types/generated/models.ts` only if produced by tygo generation

Do not touch:

- Backend model definitions
- Backend handlers, repositories, routes
- Frontend app code outside generated types
- Existing dirty feature files

Signature-stability rule: generated types reflect Go models. Do not rename fields manually to satisfy frontend call sites.

## Hard rules

- Read `CONTEXT.md`, `docs/adr/0004-subscription-semantics.md`, and `CODING_STANDARDS.md` before changes.
- Generated tygo types over hand-rolled mirrors.
- No manual edits to generated model contents unless the repo already documents a post-generation transform. Report if one exists.
- Do not fetch GitHub issues; current `gh auth status` reports an invalid token.
- Do not commit.

## Acceptance criteria

- `backend/tygo.yaml` points to the real frontend generated models path.
- The worker reports the exact generation command attempted.
- If generation runs, `frontend/src/lib/types/generated/models.ts` updates from the tool and no unrelated files change.
- If generation cannot run locally, the blocker is precise and reproducible.
- Backend validation such as `go test ./...` or a narrower relevant command is attempted if it does not require unavailable services; report exact blockers.

## Reference reading

- `CONTEXT.md` Subscription and identifiers sections
- `docs/adr/0004-subscription-semantics.md`
- `CODING_STANDARDS.md` Generated tygo types rule
- `backend/tygo.yaml`
- `backend/internal/models/subscription.go`
- `frontend/src/lib/types/generated/models.ts`

## Effort estimate

1-3 hours if tygo is installed or vendored. Longer only if the project lacks a reproducible generation path, in which case stop and report instead of installing random global machinery.

## Bidirectional message paths

- Worker to overseer: `chat_msgs/CHAT_tygo-generated-types_msg.md`
- Overseer to worker: `chat_msgs/CHAT_tygo-generated-types_reply.md`
