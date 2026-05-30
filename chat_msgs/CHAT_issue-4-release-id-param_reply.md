# 2026-05-12 12:45 CEST - reply to message #0

Proceed with pre-flight only, then report message #1 before code changes.

Focus on issue #4: route path parameter rename only.

Hard limits:

- Fair-game files: `backend/internal/routes/routes.go`, `backend/internal/handlers/release_handler.go`.
- Do not touch models, repositories, schema, storage, frontend, generated types, or unrelated dirty files.
- Do not rename `Hash` fields or response JSON in this slice.
- Run `rg` for `:releaseHash`, `c.Param("releaseHash")`, `releaseHash`, and `releaseId` before proposing edits.

Next message: append message #1 with pre-flight findings, exact files needed, validation plan, and whether the scope stayed cosmetic.

# 2026-05-12 13:31 CEST - reply to message #1

Approved. Proceed with the cosmetic route-param rename only.

Allowed edits:

- `backend/internal/routes/routes.go`
- `backend/internal/handlers/release_handler.go`

Required change:

- Change protected route params from `:releaseHash` to `:releaseId`.
- Change `c.Param("releaseHash")` to `c.Param("releaseId")` in the four affected handlers.
- Rename only local variables that hold the route parameter to `releaseId`.

Keep unchanged:

- Created UUID variables in `SaveDraftRelease`, `generatePresignedURLs`, and `buildDraftRelease`.
- Repository method names such as `GetDraftByHash`.
- Model fields, DTO fields, response JSON keys, Zap field keys, storage paths, generated TS.

Validation:

- Run the two targeted `rg` checks from your plan.
- Run `gofmt -w backend/internal/routes/routes.go backend/internal/handlers/release_handler.go`.
- Run `go test ./...` from `backend`.
- Run `go build ./...` from `backend`.
- Report curl validation only if a local backend server is already running. Otherwise report the exact blocker.

Next message: append message #2 with changed files, validation result, and whether issue #4 is complete.

# 2026-05-12 13:39 CEST - reply to message #2

Accepted and closed.

The route-param rename stayed cosmetic:

- `:releaseHash` is gone from the protected Release routes.
- `c.Param("releaseHash")` is gone from the affected handlers.
- The route-param local variables are now `releaseId`.
- Model, repository, storage, DTO, response JSON, Zap fields, generated TS, and issue #3 identity work stayed untouched.

Validation is sufficient:

- `gofmt` ran.
- Targeted route-param scans passed.
- `go test ./...` passed after Go build-cache access.
- `go build ./...` passed after Go build-cache access.
- Curl was not run because no local backend server was running.

No further work for `CHAT_issue-4-release-id-param`.
