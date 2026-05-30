# CHAT_archive-frontend

## Project context

Silkwave is a music distribution platform for independent artists (SvelteKit 5 frontend, Go 1.25 / Gin backend, ArangoDB, Cloudflare R2). This chat updates the **frontend consumers of the `unpublish` → `archive` rename** per ADR-0001. Sister chats: `releases-repo` (owns the backend contract change you depend on) and `auth-envelope` (unrelated to your scope).

## Stack location

- `frontend/src/lib/api/...` — release-related API client calls
- `frontend/src/lib/components/organisms/` — release UI organisms (release wizard, editor, list cards, status menus)
- `frontend/src/lib/components/{atoms,molecules}/` — only if a release-status surface lives there
- `frontend/src/routes/` — any route that calls the unpublish endpoint or renders unpublish copy

## Scope

Two concerns, separate commits, one PR. **The PR holds until `releases-repo` chat's commit 3 (the backend rename) is merged to `main`.**

**Commit 1 — API client and call sites**

- `rg -n 'unpublish|Unpublish' frontend/src/` to find every reference. Build the full list before changing anything.
- Replace endpoint URL: `POST /api/releases/:id/unpublish` → `POST /api/releases/:id/archive`.
- Rename API client function (e.g. `unpublishRelease` → `archiveRelease`).
- Rename any local handler / event handler that wraps the call.
- If the function lives in a Svelte 5 rune module under `lib/state/`, update it there.
- Type signatures derived from `lib/types/generated.ts` may flip once `releases-repo` regenerates tygo. **Do not hand-edit `generated.ts`.** If the new contract requires a type that doesn't exist yet, surface in `chat_msgs/CHAT_archive-frontend_msg.md`.

**Commit 2 — UI copy and labels**

- Replace user-visible "Unpublish" with "Archive" everywhere in the release UI: button labels, dialog titles, confirmation copy, tooltips, status menu items, toast messages.
- Domain language sanity check: "Archive" is a **reversible action that hides a Published Release from the public**. The copy should feel calm and reversible (per `AGENTS.md` design principles — quiet confidence, reversible actions, not alarming). Do not write copy like "Permanently remove" or "Delete from your catalog" — that's a different operation (hard delete) that this chat is not implementing.
- If a confirmation dialog explains what archiving does, use language consistent with `CONTEXT.md`: "Hide this release from public listings. You can re-publish anytime. Purchases and subscriber access remain intact for owners; new buyers won't see it."
- Frontend filename hygiene: if you're renaming or touching files that are still PascalCase, **rename to kebab-case in the same commit** per the migration-on-touch rule in `CODING_STANDARDS.md`. Update imports.

## File boundaries

**Fair-game (you own these):**

- Anything under `frontend/src/lib/api/` that touches the release endpoint.
- Anything under `frontend/src/lib/components/` that surfaces the archive action or its copy.
- Anything under `frontend/src/lib/state/` that holds release-status state and calls the endpoint.
- Anything under `frontend/src/routes/` that calls the endpoint or renders archive UI.

**DO NOT touch:**

- Any backend file (`releases-repo` and `auth-envelope` own the backend).
- `frontend/src/lib/types/generated.ts` — this is tygo output; only regenerated, never hand-edited.
- Frontend code unrelated to the archive concern (no opportunistic refactors of unrelated organisms; one concern per PR).
- Other resource consumers (Artists, Tracks, Subscriptions UI) unless the file is directly entangled with archive — flag if you find one.

**Signature-stability rule:**

You consume the contract `releases-repo` ships: `POST /api/releases/:id/archive` (response shape: envelope per `CODING_STANDARDS.md`). Wait for it. Don't predict variant names.

## Hard rules

- `CONTEXT.md` is the source of truth for domain language. `docs/adr/` for past architectural decisions. `CODING_STANDARDS.md` for code conventions. `AGENTS.md` for brand and design tone.
- No new legacy drift. New code never adds `svelte/store` usage, PascalCase component filenames, `Record<string, any>`, `any` outside vendored shadcn, mixed response handling.
- Runes only on the frontend. No `svelte/store` in app code. Shared reactive state goes in `lib/state/<name>.svelte.ts`.
- Kebab-case component filenames. Folder names too. PascalCase only on the import identifier.
- SuperForms + Zod for any server-bound form. (Not directly in scope here, but if your archive confirmation is a `<form>` posting to a server action, it must use SuperForms.)
- Generated tygo types > hand-rolled mirrors. Do not hand-edit `generated.ts`.
- One concern per PR. The two commits in this chat are the same logical migration (API call + UI copy); fine to ship as one PR.
- **Don't commit unless asked.** Stop at logical commit boundaries and surface to overseer.
- **Don't push the PR until `releases-repo` commit 3 is merged.** Verify the backend endpoint is live before pushing.
- 5-tier pre-flight discipline before estimating effort.
- Surface scope discovery early.
- **REUSE EXISTING PATTERNS.** Before inventing, run:
  ```
  rg -n 'superForm|tv\(|Form.Field|cn\(' frontend/src/lib/components
  rg -n 'unpublish' frontend/src/
  ```
- **DOMAIN LANGUAGE FROM `CONTEXT.md`:**
  - "Release" not "album".
  - "Archive" / "Archived" not "Unpublish" / "Unpublished".
  - "Archived" is a reversible state. Copy should not feel destructive.
  - "Artist" not "creator", "account", or "profile".

## Acceptance criteria

- `rg -n 'unpublish|Unpublish' frontend/src/` returns nothing in app code (test fixtures or generated files that mirror old upstream are acceptable, flag them in your message if you find any).
- API client uses the new `/archive` endpoint URL and method name.
- All release-status UI surfaces use "Archive" language; copy aligns with `AGENTS.md` tone (calm, reversible, not alarming).
- Type-check passes (`pnpm check` or whatever the project uses; verify before reporting).
- No new legacy drift. Files you renamed are kebab-case; their imports are updated.
- PR is held until `releases-repo` commit 3 is merged.

## Reference reading

- `CONTEXT.md` — Release lifecycle (especially the `Published → Archived` transition and the resolved `UnpublishRelease` ambiguity).
- `docs/adr/0001-release-staging-and-lifecycle.md`.
- `CODING_STANDARDS.md` — Frontend sections (Files and folders, Component naming, Reactivity: runes only).
- `AGENTS.md` — Design Principles (especially "Quiet confidence over loud decoration" and "Progressive disclosure" — the archive confirmation should feel reversible, not alarming).
- `frontend/src/lib/components/organisms/release-wizard/wizard.svelte.ts` — canonical organism-with-state pattern; reference if you need to add archive state.
- `frontend/src/lib/components/ui/button/button.svelte` — canonical `tv()` variant pattern if you need a destructive-but-not-too-destructive button variant.

## Effort estimate

Realistic: 4-6 hours of active work, plus idle time waiting on `releases-repo` to merge.

Apply 5-tier pre-flight first:

1. `CONTEXT.md` Release lifecycle + ADR-0001. Internalize that Archive is reversible, not destructive. `AGENTS.md` Design Principles for the copy tone.
2. `rg -n 'unpublish' frontend/src/` for the full call-site map. Read at least one canonical organism (`release-wizard/`) to know the existing patterns.
3. Cross-check standards-doc legacy-drift table against the files you'll touch. If a touched file is PascalCase or imports from `lib/stores/`, plan the migration-on-touch as a separate commit in your plan.
4. Generated tygo types: when `releases-repo` regenerates, the request DTO type for the archive endpoint may shift. Confirm what type you'll be passing into the API client and whether it's stable upstream.
5. DB shape: not directly relevant for frontend work, but confirm slug-per-Artist routing on any URL you generate.

Stop-and-report if any of the following surface:
- The `unpublish` references are entangled with other release-status concerns (e.g. a generic `setStatus` function that handles archive, publish, and delete in one path) — that may demand a refactor that's a separate concern.
- A touched file is so deep in legacy drift (PascalCase + `svelte/store` + `any` + Record<string, any>) that the migration-on-touch becomes its own scope.
- The frontend has no API client for the release endpoint at all yet (raw `fetch` everywhere) — that's a separate scope; add the client only if your scope already touches the call site.

## Bidirectional message paths

- Worker → overseer: `chat_msgs/CHAT_archive-frontend_msg.md`
- Overseer → worker: `chat_msgs/CHAT_archive-frontend_reply.md`

Convention: see `chat_msgs/README.md`. Header format `# YYYY-MM-DD HH:MM — message #N`.

## Sister chats running in parallel

- `releases-repo` — owns the backend contract you consume (`POST /api/releases/:id/archive`). Your PR holds until their commit 3 lands on `main`.
- `auth-envelope` — unrelated to your scope. No file overlap, no contract dependency.

You are the smallest of the three foundation chats by code volume but the most schedule-coupled. Use the wait time productively (write the API client and copy changes against the planned contract, hold the push). Surface to overseer if the backend contract slips.
