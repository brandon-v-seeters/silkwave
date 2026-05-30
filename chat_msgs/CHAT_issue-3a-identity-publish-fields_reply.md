# 2026-05-12 14:12 CEST - reply to message #0

Proceed with pre-flight only, then report message #1 before code changes.

Scope is the mechanical `issue-3a` naming slice:

- Release/Track `Hash` -> `Id`
- `releaseDate` -> `publishAt`
- remove `scheduled` / `deleted` status constants
- add generated Artist `id`
- AQL/schema/test/storage naming updates

Hard stop:

- Do not remove `Release.IsUploaded`.
- Do not change upload readiness logic.
- Do not collapse collections or change public URL shape.

Next message: append message #1 with exact implementation plan, files, risky lines, and whether you can complete 3a in one pass.

# 2026-05-12 14:20 CEST - reply to message #0

Worker closed by overseer after two nudges with no message #1 and no useful status.

No implementation approval was granted. No issue-3a code changes should be attributed to this chat.

# 2026-05-12 16:27 CEST - reply to message #1

Issue-3a is complete via overseer local implementation.

Backend checks passed:

- `go test ./...`
- `go build ./...`

Frontend check still fails repo-wide:

- `pnpm check` reports 60 errors and 3 warnings.
- The remaining failures are unrelated existing debt: `Locals.user`, env exports, Arango frontend imports, Bits UI type drift, and old `.ts` import paths.

Known follow-up:

- Install or wire `tygo` so generated TS can be produced from Go instead of patched locally.
- Split issue-3b for replacing upload readiness and then removing `IsUploaded`.
