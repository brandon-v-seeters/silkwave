# CHAT_adr0002-release-url-backend

## Project context

Silkwave is a music distribution platform for independent artists. This worker implements the backend side of ADR-0002's public Release URL shape: public Release lookup must move from a global slug route to an Artist-scoped route.

## Stack location

- `backend/internal/repository/release_repository.go`
- `backend/internal/repository/release_repository_test.go`
- `backend/internal/handlers/release_handler.go`
- `backend/internal/routes/routes.go` release public route lines only
- `backend/README.md` route table only if touched by the route rename

## Scope

Implement backend support for:

- `GET /api/artists/:artistSlug/releases/:releaseSlug`
- Drop or stop wiring public `GET /api/releases/:slug`
- Repository lookup must scope by Artist slug plus Release slug
- Preserve existing `GET /api/releases` list behavior
- Preserve prior wave archive/auth changes

Do not rename `Hash` to `id` in this worker. That is a separate code-level sweep and will collide with upload/storage surfaces.

## File boundaries

Fair-game:

- `backend/internal/repository/release_repository.go`
- `backend/internal/repository/release_repository_test.go`
- `backend/internal/handlers/release_handler.go`
- `backend/internal/routes/routes.go` public release route line only
- `backend/README.md` route docs only

Do not touch:

- Frontend files
- `backend/internal/models/release.go` except if a compile break forces a tiny import/type fix
- Auth files
- Subscription files
- Generated frontend types

Signature-stability rule: backend owns the route contract. Frontend worker consumes `GET /api/artists/:artistSlug/releases/:releaseSlug`.

## Hard rules

- Read `CONTEXT.md`, `docs/adr/0002-identifier-conventions.md`, and `CODING_STANDARDS.md` before editing.
- Handlers do not touch ArangoDB. Route lookup belongs in `ReleaseRepository`.
- Preserve `/*aql*/` hints.
- Use `respondOK` / `respondError`.
- No new legacy drift.
- Do not commit.
- Do not rewrite previous-wave work.
- Report if the route change forces a broader identifier sweep.

## Acceptance criteria

- `GET /api/artists/:artistSlug/releases/:releaseSlug` is wired.
- Global public `GET /api/releases/:slug` is no longer wired for public detail lookup.
- Repository lookup joins Artist by slug and filters Release by `artistKey` plus `releaseSlug`.
- Archived releases are not returned.
- Tests cover same release slug under two different Artists resolving correctly.
- `gofmt -l` clean.
- `go test ./...` and `go vet ./...` pass from `backend/`.

## Reference reading

- `CONTEXT.md` Identifiers and URLs
- `docs/adr/0002-identifier-conventions.md`
- `CODING_STANDARDS.md` Backend repository pattern and response envelope
- `backend/internal/repository/release_repository.go`
- `backend/internal/handlers/release_handler.go`
- `backend/internal/routes/routes.go`

## Effort estimate

2-4 hours if the prior release repository work is stable. Stop and report if this expands into the `Hash` to `id` sweep.

## Bidirectional message paths

- Worker to overseer: `chat_msgs/CHAT_adr0002-release-url-backend_msg.md`
- Overseer to worker: `chat_msgs/CHAT_adr0002-release-url-backend_reply.md`
