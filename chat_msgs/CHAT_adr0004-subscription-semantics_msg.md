# 2026-05-11 20:08 CEST - message #1

Plan confirmation.

Scope accepted:
- Own only `backend/internal/models/subscription.go`, `backend/internal/database/schema.go`, and subscription wording in `backend/README.md` if stale.
- Do not touch release, auth, frontend route/component, generated type, Stripe, checkout, entitlement, or billing-admin files.
- Keep model names stable: `Subscription`, `Subscriber`, `SubscriptionStatus`.

Pre-flight started:
- Tier 1 done: `CONTEXT.md` and `docs/adr/0004-subscription-semantics.md` confirm Bandcamp-style paid support levels, identical streaming access, no free Subscription tier, no entitlement table, and `Subscriber` status as `active` / `inactive`.
- Tier 3 started: `CODING_STANDARDS.md` confirms backend model/schema work must avoid new legacy drift and generated frontend types should come from tygo, not hand edits.

Plan before code changes:
1. Finish pre-flight tiers 2-5 with `rg` over subscription/follow/status usage, fair-game file reads, schema/index inspection, and generated-type staleness check.
2. Report immediately if compile forces a DTO reference outside fair-game files.
3. Edit only fair-game files to remove tier semantics, add support-level fields, collapse status constants, and replace `artistKey,tier` indexing.
4. Run `gofmt -l`, then `go test ./...` and `go vet ./...` from `backend/`.
5. Report final status, test output, and whether frontend generated types are stale.

No code changes made yet.

# 2026-05-11 20:18 CEST - message #2

Implementation and verification complete.

Changed fair-game files only:
- `backend/internal/models/subscription.go`
  - Removed `SubscriptionTier`, `SubscriptionTierFree`, and `SubscriptionTierPaid`.
  - Collapsed `SubscriptionStatus` to `active` / `inactive`.
  - Kept `Subscription`, `Subscriber`, and `SubscriptionStatus` names stable.
  - Kept support-level shape as `price` + `currency` and added optional `subscriberDiscountPercent`.
  - Removed `tier` from create/update DTOs and `ClientSubscriber`.
  - Added support-level response fields on `ClientSubscriber`: `subscriptionPrice`, `subscriptionCurrency`, `subscriberDiscountPercent`.
- `backend/internal/database/schema.go`
  - Replaced the `artistKey,tier` index with unique `artistKey,currency,price`.
- `backend/README.md`
  - Changed stale DB wording from subscription tiers/user subscriptions to paid support levels/subscription relationships.
  - Note: this file already had an unrelated release-route diff from another chat. I left it alone.

Pre-flight result:
- Tier 1: ADR and `CONTEXT.md` agree: Bandcamp-style paid support levels, no free Subscription, no access tiers, no entitlement table, status active/inactive only.
- Tier 2: `rg` found tier/status drift concentrated in subscription model/schema and generated TS. No handler/repository subscription implementation exists in this slice.
- Tier 3: No new banned backend drift added.
- Tier 4: Generated frontend types are stale. `frontend/src/lib/types/generated/models.ts` still exposes `SubscriptionTier`, old Stripe-like statuses, and `tier` fields. I did not edit generated files.
- Tier 5: DB shape is schema-only for this scope. No subscription repository/migration expansion touched. Index now encodes distinct paid support levels.

Verification:
- `gofmt -l backend/internal/models/subscription.go backend/internal/database/schema.go` clean.
- `go test ./...` from `backend/` passed after sandbox approval. Initial sandbox run hit Go build-cache permission denial under `/Users/brandon/Library/Caches/go-build`.
- `go vet ./...` from `backend/` passed.

Generated type status:
- Documented command is `tygo generate` from `backend/`, but `tygo` is not installed/on PATH in this environment.
- `backend/tygo.yaml` currently outputs to `../silk-wave/src/lib/types/generated/models.ts`, which does not exist from `backend/`.
- Actual generated file in this repo is `frontend/src/lib/types/generated/models.ts`.
- So the precise next step is either install/use tygo and fix `backend/tygo.yaml` output path to the current frontend location, or run the repo’s missing wrapper if one exists elsewhere. The current documented command cannot regenerate the checked-in generated file as configured. The cosmos remains unmoved.
