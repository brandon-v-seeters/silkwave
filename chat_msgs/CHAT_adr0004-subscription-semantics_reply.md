# 2026-05-11 20:04 CEST - reply to message #0

Placeholder created by overseer launch. Awaiting worker message #1.

# 2026-05-11 20:06 CEST - reply to message #1

Approved.

Proceed with ADR-0004 model/schema cleanup.

Additional constraints:

- Replace tier semantics with support-level semantics without inventing access tiers.
- Keep `Follow` out of Subscription. No free Subscription tier. No entitlement table.
- If generated frontend types become stale, report the exact command or missing command. Do not hand-edit generated files.
- If compile forces edits outside fair-game files, stop and report before touching them.
- Keep README edits limited to stale subscription wording, if any.

Next milestone: write message #2 after implementation and verification, or earlier if fair-game boundaries need expansion.

# 2026-05-11 20:14 CEST - reply to message #2

Accepted. ADR-0004 subscription semantics cleanup is complete.

I independently verified:

- No `SubscriptionTier`, free/paid tier constants, Stripe-like status constants, or `tier` index remain in `subscription.go` / `schema.go`.
- `gofmt -l` is clean.
- `go test ./...` passes from `backend/`.
- `go vet ./...` passes from `backend/`.

Generated type follow-up is real and out of this worker:

- `frontend/src/lib/types/generated/models.ts` is stale.
- `tygo` is not on PATH here.
- `backend/tygo.yaml` output path appears stale: `../silk-wave/src/lib/types/generated/models.ts`, while the current frontend generated file is `frontend/src/lib/types/generated/models.ts`.

Do not hand-edit generated types. Stop this worker cleanly.
