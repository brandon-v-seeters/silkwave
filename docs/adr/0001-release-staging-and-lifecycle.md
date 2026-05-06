# Release identity, staging, and lifecycle

A Release is a single document with a single canonical identity from initial draft through every subsequent edit. Staged changes to a Published Release live in an embedded `pending` field on the Release itself (not in a sibling row, sibling collection, or separate "draft" entity), and that blob carries partial overrides for both Release-level fields and individual Track-level fields keyed by track id. Track *membership* is fixed at first publish — adding or removing tracks after publish is not allowed; that case is handled by a future Reissue concept (a new linked Release), not by mutating the original.

The Release lifecycle has three states: `Draft → Published → Archived`. There is no `scheduled` state (a future-dated `publishAt` on a Published Release covers the scheduling case), no `deleted` status (deletion is a hard delete, not a soft-delete row), and no path from Published back to Draft (unpublishing transitions to Archived; Draft means *never been Published*). Discarded Pending Edits are gone forever — there is no version history, no rollback, no edit log; recovering from a bad publish is a backup-restore problem, not a model concern.

## Considered options

- **Two collections (`Releases` + `ReleaseDrafts`)** — current code has the schema split but the publish flow doesn't actually move records across it, and the split offers no benefit once we collapse staging into an embedded blob.
- **Sibling draft row in one collection** — gives the same shape on disk but introduces a "Release identity vs Release version" distinction that's not load-bearing once track membership is locked at publish.
- **Separate `ReleasePendingEdit` document linked by `releaseKey`** — clean separation but adds a fourth document type and a join on every artist edit-view read.
- **Per-Track `pending` blob (in addition to Release `pending`)** — symmetric but spreads dirty-state across N+1 documents and creates concurrency hazards on multi-field saves.
- **Versioned releases with full history** — over-investment for users we don't have; can be added later as an append-only `releaseVersions` log if it becomes a feature.

## Consequences

- The Release document grows during a Pending Edit (carries both the live values and an embedded partial-override of those values plus track partials). This is acceptable at expected scale.
- The R2 object layout (`artist_content/{artistKey}/releases/{releaseId}/draft/...` vs `.../releases/{releaseId}/...`) was originally designed around the now-collapsed two-collection split. During a Pending Edit, the only file that can change is cover art (track audio is locked by the membership rule), so the `/draft/` subdirectory's role narrows from "home of all draft files" to "home of pending-cover-art." The directory structure stays for now; the publish-the-edit operation moves the staged cover up into the published path.
- Adding tracks to a Published release is structurally impossible. Future "deluxe edition" use cases require a separate Reissue Release that references the original.
