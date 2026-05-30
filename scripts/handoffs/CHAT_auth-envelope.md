# CHAT_auth-envelope

## Project context

Silkwave is a music distribution platform for independent artists (SvelteKit 5 frontend, Go 1.25 / Gin backend, ArangoDB, Cloudflare R2). This chat ships the **response envelope foundation** and migrates the auth handler as the canonical migrated-handler example. It's the first chat in a three-chat foundation wave; sister chats are `releases-repo` and `archive-frontend`.

## Stack location

- `backend/internal/handlers/response.go` — **create**
- `backend/internal/handlers/auth_handler.go` — migrate
- `backend/internal/repository/user_repository.go` — accept the inline-AQL `Login` query
- `backend/internal/models/` — auth-related DTOs only (LoginUserRequest, RegisterUserRequest, etc.) for `binding:"..."` struct tags

## Scope

Two concerns, separate commits, one PR is fine.

**Commit 1 — envelope foundation**

Create `internal/handlers/response.go`:

- `type Response[T any]` with `Data *T`, `Error *ErrorBody`, `Message string`
- `type ErrorBody` with `Code`, `Message`, `Details`
- Helper functions `respondOK(c *gin.Context, data any, message string)` and `respondError(c *gin.Context, status int, code, message string)`
- Match the shape spelled out in `CODING_STANDARDS.md` "Response shape: envelope" exactly. Frontend reads `data` for success, `error` for failures.

Commit signature is **stable from this point**. Sister chat `releases-repo` will import these helpers; do not rename or change their signatures after the first push.

**Commit 2 — auth_handler migration**

Migrate `auth_handler.go`:

- **Security: remove password from log line at `auth_handler.go:47`.** This is the fix-on-sight bug from the standards doc.
- Delete `fmt.Println(user)` at `:135`.
- Delete the hand-rolled `contains` and `containsHelper` at `:196-207`. Replace usage with `strings.Contains`.
- Replace manual nil/empty-string validation after `ShouldBindJSON` with `binding:"required,..."` struct tags on the request DTOs.
- Move inline AQL (the `Login` query that calls `database.QueryOne` directly) into a method on `user_repository.go`. Use the `/*aql*/` hint comment.
- Migrate every response to `respondOK` / `respondError`.
- Strip logger calls below the handler. Wrap errors with `%w` instead. Log only at the handler boundary.
- Use Zap structured fields for the handler's log lines. **Email is fine to log; passwords, hashes, JWTs, full request bodies are not.**

## File boundaries

**Fair-game (you own these):**

- `backend/internal/handlers/response.go` (you create it)
- `backend/internal/handlers/auth_handler.go`
- `backend/internal/repository/user_repository.go`
- `backend/internal/models/*.go` — auth-related DTOs only

**DO NOT touch:**

- `backend/internal/handlers/release_handler.go` (`releases-repo` owns this)
- `backend/internal/repository/release_repository.go` (will be created by `releases-repo`)
- `backend/internal/routes/routes.go` (`releases-repo` is touching the release route lines; stay out)
- Any frontend file (`archive-frontend` and the user own those)
- Any other resource handler

**Signature-stability rule:**

After commit 1 ships, `respondOK(c, data, message)` and `respondError(c, status, code, message)` are frozen. If you discover a reason to change them, surface it via `chat_msgs/CHAT_auth-envelope_msg.md` before changing — sister chats will already be consuming the helpers.

## Hard rules

- `CONTEXT.md` is the source of truth for domain language. `docs/adr/` for past architectural decisions. `CODING_STANDARDS.md` for code conventions. Read before guessing.
- No new legacy drift. New code never adds `fmt.Println`, password logging, manual nil/empty checks after `ShouldBindJSON`, hand-rolled `containsHelper`-style utilities, errors logged below the handler, mixed response shapes, inline AQL in handlers, `any` outside vendored shadcn.
- Repo before handler. Handlers do not touch ArangoDB.
- Wrap with `%w` at every layer below the handler. Log only at the handler boundary.
- Tests for new business logic: happy path + obvious failure mode. Don't gold-plate.
- Generated tygo types > hand-rolled mirrors.
- One concern per PR — but two commits in one PR is fine when the second commit consumes the first (envelope foundation → auth handler migration).
- **Don't commit unless asked.** Stop at logical commit boundaries and surface to overseer.
- 5-tier pre-flight discipline before estimating effort.
- Surface scope discovery early. Stop-and-report if effort balloons.
- **REUSE EXISTING PATTERNS.** Before inventing, run:
  ```
  rg -n 'respondOK|respondError|envelope' backend/
  rg -n 'QueryOne|QueryAll' backend/internal/
  ```
  `user_repository.go` is the canonical repo example. Match its shape.
- **NEVER LOG SECRETS.** Passwords, password hashes, JWTs, full request bodies. Email is fine.
- **DOMAIN LANGUAGE FROM `CONTEXT.md`.** "User", "Artist" (a User manages an Artist; not "creator" or "account").

## Acceptance criteria

- `internal/handlers/response.go` exists with the typed envelope and helpers, matching the standards-doc shape.
- `auth_handler.go` after migration:
  - No password in any log line.
  - No `fmt.Println`.
  - No `containsHelper` / `contains`. Uses `strings.Contains`.
  - No manual nil/empty-string checks after `ShouldBindJSON`. Uses `binding:"..."` tags.
  - No inline AQL — moved to a `UserRepository` method.
  - Every response goes through `respondOK` / `respondError`.
  - Errors below the handler are wrapped with `%w`; only handler logs.
- Tests for the new `UserRepository` method (the moved Login query). Happy path + not-found.
- `go vet ./...` clean. `gofmt -l` returns nothing.
- No new legacy drift.

## Reference reading

- `CONTEXT.md` — User and Artist sections.
- `CODING_STANDARDS.md` — Backend → Repository pattern, Request validation, Response shape, Errors and logging, Logging conventions.
- `docs/adr/0003-single-operator-artist-ownership.md` — auth/User context.
- `backend/internal/repository/user_repository.go` — canonical repo example to extend.
- `backend/internal/handlers/auth_handler.go` — migration target.

## Effort estimate

Realistic: 4-8 hours of focused work, two commits, one PR.

Apply 5-tier pre-flight first:

1. `CONTEXT.md` User section + `CODING_STANDARDS.md` Backend sections.
2. `rg -n 'respondOK|envelope' backend/` to confirm no helper already drifted in. Read `user_repository.go` for the canonical repo shape.
3. Cross-check the standards-doc legacy-drift table against `auth_handler.go`. Confirm the seven items still match (password log, `fmt.Println`, manual nil checks, `containsHelper`, inline AQL, mixed response shapes, errors logged below handler).
4. Tygo types are not relevant for this scope (no Go model shape changes).
5. DB shape: `Users` collection. The Login query joins or filters — confirm the existing query before moving it.

Stop-and-report if you discover anything that would push the work past one day.

## Bidirectional message paths

- Worker → overseer: `chat_msgs/CHAT_auth-envelope_msg.md`
- Overseer → worker: `chat_msgs/CHAT_auth-envelope_reply.md`

Convention: see `chat_msgs/README.md`. Header format `# YYYY-MM-DD HH:MM — message #N`. Always read only the latest header onward.

## Sister chats running in parallel

- `releases-repo` — creating Releases repo, migrating release handler, backend `unpublish`→`archive` rename. Will import your `respondOK` / `respondError` once commit 1 lands. Coordinate via overseer if you discover a signature change is needed.
- `archive-frontend` — frontend consumer of the rename. No file overlap with you.

You are the foundation. Land commit 1 quickly; the other chats are unblocked the moment it's pushed.
