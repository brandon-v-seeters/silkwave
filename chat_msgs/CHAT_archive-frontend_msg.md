# 2026-05-11 10:28 CEST - message #0

Placeholder created by overseer launch. Worker must write message #1 with plan confirmation before code changes.

# 2026-05-11 10:29 CEST - message #1

Plan confirmation before code changes.

Scope accepted: I own only frontend consumers of the `unpublish` to `archive` rename. Fair game remains release-related files under `frontend/src/lib/api/`, `frontend/src/lib/components/`, `frontend/src/lib/state/`, and `frontend/src/routes/` when they call the endpoint or surface archive UI. I will not touch backend files, `frontend/src/lib/types/generated.ts`, unrelated resource consumers, or sister-chat files. The cosmos can survive without opportunistic refactors.

Pre-flight plan:

1. Read `CONTEXT.md`, ADR-0001, frontend standards in `CODING_STANDARDS.md`, and relevant `AGENTS.md` design guidance before estimating or editing.
2. Build the full frontend reference map with `rg -n 'unpublish|Unpublish' frontend/src/`, plus the required pattern checks for `superForm`, `tv(`, `Form.Field`, and `cn(`.
3. Cross-check touched files for legacy drift before edits: PascalCase filenames, `svelte/store`, `any`, `Record<string, any>`, and response-shape drift.
4. Inspect generated type usage but do not edit generated tygo output. If the archive contract needs a type not present yet, I will stop and report.
5. Confirm no frontend route work creates backend or DB assumptions. DB shape is sister-chat territory unless URL generation exposes a frontend issue.

Implementation plan after pre-flight:

- Commit-boundary 1 equivalent: rename API client and call-site handlers from unpublish to archive, switching the endpoint to `POST /api/releases/:id/archive`.
- Commit-boundary 2 equivalent: update user-facing release UI copy from Unpublish to Archive with calm reversible language.
- Hold PR/push until `releases-repo` commit 3 lands on `main`; no commits unless explicitly asked.

Blockers: none yet. Dependency noted: backend archive endpoint must land before final integration confidence.

# 2026-05-11 10:32 CEST - message #2

Pre-flight and reference map complete. Scope contracted to no frontend code patch.

Findings:

- `rg -n 'unpublish|Unpublish' frontend/src` returns no matches.
- A wider fair-game scan found release publish/draft code, generated `ReleaseStatusArchived`, and generic archive icons/routes, but no release archive/unpublish API consumer, handler, confirmation dialog, toast, tooltip, or visible release-status action to rename.
- `frontend/src/lib/state` does not exist in this checkout. The handoff mentioned it as a possible location; this repo currently uses `frontend/src/lib/stores` legacy state in unrelated surfaces. I did not touch it. Entropy keeps receipts.
- Generated types live under `frontend/src/lib/types/generated/index.ts` and `frontend/src/lib/types/generated/models.ts`, not `frontend/src/lib/types/generated.ts`. I inspected only enough to confirm `ReleaseStatusArchived` exists and did not edit generated output.
- Existing legacy drift exists in unrelated frontend files (`lib/stores` imports, PascalCase component files, `any` in API helpers), but none is on a touched archive/unpublish path because there is no touched path.

Verification:

- `pnpm check` from `frontend/` did not run because dependencies are not installed: `svelte-kit: command not found`, with pnpm warning that `node_modules` is missing.
- No app-code files were edited.

Decision ask:

Should this worker stop as complete/no-op for now, or should I wait for `releases-repo` commit 3 and re-scan after the backend rename lands? I do not recommend inventing a new archive UI or API wrapper unless the overseer expands the scope. That would be feature creation wearing a rename mask, which is how small tasks become fossil strata.
