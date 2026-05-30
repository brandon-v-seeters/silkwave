# 2026-05-12 14:12 CEST - message #0

Placeholder created by overseer. Worker should write message #1 after pre-flight and before code changes.

# 2026-05-12 16:27 CEST - message #1

Overseer took over local implementation after the worker stalled and was closed.

Implemented issue-3a:

- Renamed Release and Track external identity fields from `hash` to `id`.
- Renamed Release visibility gate fields from `releaseDate` to `publishAt`.
- Removed `scheduled` and `deleted` Release status constants.
- Added stable external `id` to Artist creation.
- Updated Release repository AQL to filter by `status == published` and `publishAt == null || publishAt <= now`.
- Updated draft lookup, upload confirmation, status update, and delete paths to use `releaseId`.
- Updated storage helper parameter names and comments from `releaseHash` to `releaseId` without changing object key shape.
- Updated backend schema indexes from `hash` / `releaseDate` / `published` to `id` / `publishAt` / `status`.
- Updated generated frontend model types, release upload service, search result keys, discover/date reads, and draft wizard typing.

Deliberately not done:

- Did not remove `Release.IsUploaded`; it still gates publish.
- Did not collapse `Releases` and `ReleaseDrafts`.
- Did not change public slug URL shape.

Validation:

- `go test ./...` passed after rerun with host Go build-cache access.
- `go build ./...` passed after rerun with host Go build-cache access.
- `tygo generate` failed because `tygo` is not installed on PATH.
- `pnpm check` still fails repo-wide with 60 errors and 3 warnings. Remaining errors are existing project-level type/env/UI-library debt, not issue-3a files.
