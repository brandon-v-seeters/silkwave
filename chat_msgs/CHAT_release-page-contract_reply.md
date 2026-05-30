# 2026-05-14 10:26 CEST — reply to message #0

Read `scripts/handoffs/CHAT_release-page-contract.md`, run the 5-tier pre-flight, then write message #1 with your proposed plan before code changes.

# 2026-05-14 10:35 CEST — reply to message #1

Approved with one constraint.

Proceed with the bounded contract fix: include ordered Tracks on the public Release response while preserving the existing public visibility gates and envelope shape.

Preferred path:

- Add a named backend model only if it keeps the repository signature honest.
- Update `ReleaseRepository.GetPublishedByArtistSlugAndReleaseSlug` to join `Tracks` by `release._key`, sort by `track.order ASC`, and return `tracks`.
- Keep `ReleaseHandler.GetReleaseByArtistAndSlug` boring. No handler-side query logic.
- Do not touch generated frontend types in this worker. If a Go model change requires tygo regeneration, report it in message #2 and let overseer coordinate.
- Only touch `+page.server.ts` if TypeScript genuinely needs a local type cleanup. The UI worker already consumes `release.tracks`.

Important dirty-tree note: `backend/internal/handlers/release_handler.go` and `backend/internal/models/release.go` already have unrelated modifications in this checkout. Read before editing and do not erase them. If the model file is too risky to touch, prefer returning an anonymous merged shape from AQL that fits the existing local `PublicRelease` type and report the tradeoff.

Validation:

- Run `gofmt` only on files you edit.
- Run `go test ./...` from `backend` if backend code changes.
- Run `go build ./...` from `backend` if tests pass or if tests are unavailable.

After implementation, write message #2 with changed files, validation result, and whether tygo regeneration is needed.

# 2026-05-14 11:03 CEST — reply to message #2

Accepted.

The contract change is the right shape: public slug lookup now returns `artist` plus ordered `tracks`, and it preserves the Published plus `publishAt` gates. Handler stayed boring. Good. Boring code is bedrock.

Known follow-up:

- `models.PublicRelease` is exported Go surface. Generated frontend types were intentionally not regenerated in this worker. Track this as tygo parity debt, not a blocker for the current page because the route already has a local `PublicRelease` type.
- Commerce remains out of scope. No checkout or purchase endpoint was found.

Status: contract workstream complete pending overseer integrated validation.
