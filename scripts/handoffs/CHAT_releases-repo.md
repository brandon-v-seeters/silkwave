# CHAT_releases-repo

## Project context

Silkwave is a music distribution platform for independent artists (SvelteKit 5 frontend, Go 1.25 / Gin backend, ArangoDB, Cloudflare R2). This chat creates the **Releases repository**, ports `release_handler.go` to use it, and bundles the backend `unpublish`→`archive` rename per ADR-0001. Sister chats: `auth-envelope` (response envelope foundation, your dependency) and `archive-frontend` (consumes your renamed contract).

## Stack location

- `backend/internal/repository/release_repository.go` — **create**
- `backend/internal/handlers/release_handler.go` — migrate (no direct DB, envelope responses, archive rename)
- `backend/internal/routes/routes.go` — release routes only (route rename `/unpublish` → `/archive`); leave auth and other resource routes alone
- `backend/internal/models/release.go` — touch only what the rename and repo work require (status constants, method renames)

## Scope

Three concerns, separate commits, one PR.

**Commit 1 — Releases repository (parallel-safe with `auth-envelope`)**

Create `internal/repository/release_repository.go`:

- One method per query `release_handler.go` currently calls. Audit the handler first; produce the method list before writing any code.
- Use the `database.QueryOne[T]` / `database.QueryAll[T]` helpers, matching the shape of `user_repository.go`.
- Keep the `/*aql*/` hint comment on every query string. The AQL extension uses it.
- Per-Artist slug scoping on lookups (`FILTER r.artistKey == @artistKey AND r.slug == @slug`) — slugs are unique per Artist, not globally. See `CONTEXT.md` "Identifiers and URLs" and ADR-0002.
- Tests in a sibling `_test.go`: happy path + the obvious failure (not-found, duplicate-slug, archived-not-listed) per query type. Don't gold-plate coverage.

This commit is independent of `auth-envelope`. Land it while the other chat works on `response.go`.

**Commit 2 — handler migration (waits for `auth-envelope` commit 1)**

Once `internal/handlers/response.go` exists on `main`, migrate `release_handler.go`:

- Strip every `database.QueryOne` / `database.QueryAll` call. Handlers do not touch ArangoDB.
- Replace with calls to the repo from commit 1.
- Migrate every response to `respondOK` / `respondError` from `response.go`.
- Replace manual nil/empty-string validation after `ShouldBindJSON` with `binding:"required,..."` struct tags on the request DTOs.
- Strip logger calls below the handler. Wrap with `%w`. Handler logs structured fields with Zap.
- Custom rules struct tags can't express (cross-field constraints, business rules) live in a service or the repo, not the handler.

**Commit 3 — `unpublish` → `archive` rename per ADR-0001**

Per ADR-0001 / `CONTEXT.md` "Flagged ambiguities":

- The operation transitions a Release from **Published** to **Archived**, not back to Draft. Draft means *never Published*.
- Rename the handler method `UnpublishRelease` → `ArchiveRelease`.
- Rename the route registration in `routes/routes.go` from `POST /api/releases/:id/unpublish` to `POST /api/releases/:id/archive`. Keep the path style (no trailing slash, ID in path).
- If the status transition logic still sets status to `draft`, change it to `archived`.
- This is a **breaking contract change**. The frontend consumer (`archive-frontend` chat) holds its PR until this lands.

**OUT OF SCOPE — surface in message #1 if discovered, but do not bundle:**

- Full ADR-0002 `Hash` → `id` field rename across `models/release.go`.
- Removing the redundant `Published` boolean from the Release model.
- Renaming `releaseDate` → `publishAt`.
- Dropping the `scheduled` and `deleted` status constants.
- Collapsing `Releases` / `ReleaseDrafts` collections.

These are bigger model migrations that ripple into tygo regen, frontend types, and AQL. They are their own future chat. **Stop-and-report if any of them block the commit-3 rename.**

## File boundaries

**Fair-game (you own these):**

- `backend/internal/repository/release_repository.go` (you create it)
- `backend/internal/handlers/release_handler.go`
- `backend/internal/routes/routes.go` — release-route lines only
- `backend/internal/models/release.go` — only changes the repo or rename require

**DO NOT touch:**

- `backend/internal/handlers/response.go` (`auth-envelope` owns it; you read and import it)
- `backend/internal/handlers/auth_handler.go` and `internal/repository/user_repository.go` (`auth-envelope` owns these)
- `backend/internal/routes/routes.go` — non-release lines (auth, other resources). Leave them untouched.
- Any frontend file (`archive-frontend` owns the consumer side)
- Any other resource handler (Artists, Tracks, Subscriptions, Subscribers — they need repos too, but not in this chat)

**Signature-stability rule:**

After commit 3 pushes, the URL `POST /api/releases/:id/archive` is the contract `archive-frontend` consumes. Don't change the path or the method name after the first push. If the contract needs to change, surface via `chat_msgs/CHAT_releases-repo_msg.md`.

## Hard rules

- `CONTEXT.md` is the source of truth for domain language. `docs/adr/` for past architectural decisions. `CODING_STANDARDS.md` for code conventions.
- No new legacy drift. New code never adds inline AQL in handlers, manual nil/empty checks after `ShouldBindJSON`, mixed response shapes, errors logged below the handler, `fmt.Println`, hand-rolled string utilities.
- Repo before handler. Handlers do not touch ArangoDB.
- Envelope response shape via `respondOK` / `respondError`. No ad-hoc shapes.
- Wrap errors with `%w` below the handler. Log only at the handler boundary, with Zap structured fields.
- Generated tygo types > hand-rolled mirrors.
- One concern per PR — three commits in one PR is fine here because they sequence into the same logical migration. If commits 2 or 3 grow into separate concerns, split into separate PRs.
- **Don't commit unless asked.** Stop at logical commit boundaries and surface to overseer.
- 5-tier pre-flight discipline before estimating effort.
- Surface scope discovery early. Stop-and-report if effort balloons.
- **REUSE EXISTING PATTERNS.** Before inventing, run:
  ```
  rg -n 'QueryOne|QueryAll' backend/internal/
  rg -n 'release|Release' backend/internal/{handlers,repository,models}
  ```
  `user_repository.go` is the canonical repo example. Match its shape.
- **NEVER LOG SECRETS.** No password hashes, JWTs, full request bodies in logs.
- **DOMAIN LANGUAGE FROM `CONTEXT.md`:**
  - "Release" not "album" (album is a release type, not a synonym).
  - "Archived" not "Draft" for a Release that was Published and taken down.
  - "Pending Edit" for staged changes to a Published Release; "Draft" only for *never-Published*.
  - "Artist" not "creator" or "account".
  - "publishAt" not "releaseDate" (when you encounter it; full rename is out of scope but don't introduce new uses of `releaseDate`).
  - "id" not "Hash" (when you encounter it; full rename is out of scope).
  - Slugs are unique **per Artist**, not globally.

## Acceptance criteria

- `release_repository.go` exists. Every query `release_handler.go` calls is implemented as a repo method. Tests for happy path + obvious failure mode per method.
- `release_handler.go` after migration:
  - No `database.QueryOne` / `database.QueryAll` calls (no direct DB).
  - Every response goes through `respondOK` / `respondError`.
  - No manual nil/empty checks after `ShouldBindJSON`.
  - No `fmt.Println`. No hand-rolled string utilities.
  - Errors below the handler are wrapped with `%w`.
- `unpublish` → `archive` rename:
  - Route is `POST /api/releases/:id/archive`.
  - Handler method is `ArchiveRelease`.
  - Status transition is `Published → Archived` (not Draft).
- `go vet ./...` clean. `gofmt -l` returns nothing.
- No new legacy drift introduced. The status constants you do touch reflect the canonical state machine; if you discover the model still has `scheduled` or `deleted` constants, **note it in `chat_msgs/CHAT_releases-repo_msg.md` but do not delete them in this chat** — surface as out-of-scope.

## Reference reading

- `CONTEXT.md` — Releases and tracks; Release lifecycle; Identifiers and URLs; Flagged ambiguities (especially `UnpublishRelease`).
- `docs/adr/0001-release-staging-and-lifecycle.md` — the canonical state machine and the rename rationale.
- `docs/adr/0002-identifier-conventions.md` — slug-per-Artist scoping; `id` vs `_key`.
- `CODING_STANDARDS.md` — Backend sections (Repository pattern, Request validation, Response shape, Errors and logging).
- `backend/internal/repository/user_repository.go` — canonical repo pattern to mirror.
- `backend/internal/handlers/release_handler.go` — migration target. Read end-to-end before planning.
- `frontend/src/lib/components/organisms/release-wizard/wizard.svelte.ts` — read-only context for downstream consumer; do not edit.

## Effort estimate

Realistic: 1-2 days, three commits, one PR.

Apply 5-tier pre-flight first:

1. `CONTEXT.md` Releases sections + ADR-0001 + ADR-0002. Internalize the state machine and slug semantics before writing AQL.
2. `rg -n 'release' backend/internal/{handlers,repository,models}` to find every existing release surface. Read `user_repository.go` for the canonical pattern. Confirm `database.QueryOne[T]` is the helper signature.
3. Cross-check standards-doc legacy-drift table against `release_handler.go`. Confirm the items still apply.
4. Tygo: any change to `models/release.go` field names propagates to `frontend/src/lib/types/generated.ts`. If your scope (rename only, no field renames) doesn't touch field names, tygo regen is not in scope.
5. DB shape: confirm `Releases` collection exists. Confirm slug uniqueness is enforced per-Artist (not globally) — this is a constraint your repo's slug-lookup methods must respect.

Stop-and-report if any of the following surface:
- Repo method count exceeds 8 (handler is bigger than expected; consider splitting).
- The Release model has both `Hash` and `id` fields, both `Published` boolean and a `Status` field, or a `releaseDate` reference your work touches — these are model-rename concerns that escalate scope.
- The `archive` transition logic exists in a service layer that's not yet introduced — you may need to introduce a release service to hold the business rule.

## Bidirectional message paths

- Worker → overseer: `chat_msgs/CHAT_releases-repo_msg.md`
- Overseer → worker: `chat_msgs/CHAT_releases-repo_reply.md`

Convention: see `chat_msgs/README.md`. Header format `# YYYY-MM-DD HH:MM — message #N`.

## Sister chats running in parallel

- `auth-envelope` — ships `response.go` (your dependency for commits 2 and 3) and migrates auth handler. Watch for commit 1 of that chat to land before starting your commit 2.
- `archive-frontend` — frontend consumer of your `archive` rename. Holds its PR until your commit 3 ships. No file overlap.

You are the largest of the three foundation chats. Surface scope expansion early.
