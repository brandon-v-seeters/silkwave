# CHAT_adr0004-subscription-semantics

## Project context

Silkwave Subscriptions follow the Bandcamp model from ADR-0004: paid patron-support levels, identical streaming access, no free Subscription tier, no entitlement table, and Subscriber status is only active or inactive.

## Stack location

- `backend/internal/models/subscription.go`
- `backend/internal/database/schema.go`
- `backend/README.md` subscription wording only if stale route/model docs are present
- Generated frontend types only via tygo if a verified project command exists and the backend model changes require it

## Scope

Clean backend subscription model drift:

- Remove `SubscriptionTierFree` / `SubscriptionTierPaid` tier semantics from models and requests.
- Replace tier with support-level fields that match ADR-0004, such as `price`, `currency`, and optional `subscriberDiscountPercent`.
- Collapse `SubscriptionStatus` to `active` / `inactive`.
- Keep `Follow` separate. Do not model free following as a Subscription.
- Update schema indexes that currently use `artistKey, tier`.

Do not implement Stripe webhooks, billing-admin UI, checkout, streaming gates, purchases, or frontend subscription pages in this worker.

## File boundaries

Fair-game:

- `backend/internal/models/subscription.go`
- `backend/internal/database/schema.go`
- `backend/README.md` subscription model wording only

Do not touch:

- Release files
- Auth files
- Frontend route/component files
- Generated frontend types unless the repo has a known generation command and you report before running it
- Stripe integration, because none is in scope here

Signature-stability rule: model names should remain `Subscription`, `Subscriber`, and `SubscriptionStatus` unless pre-flight proves the repo uses a different canonical name.

## Hard rules

- Read `CONTEXT.md` Subscription, Subscriber, Follow sections and `docs/adr/0004-subscription-semantics.md`.
- Do not add `UserReleaseAccess`, entitlement tables, free tiers, exclusive content fields, or Patreon-style gated concepts.
- Keep status domain as active/inactive only.
- No handler or route work unless compile forces a tiny DTO reference update.
- Do not commit.

## Acceptance criteria

- Subscription model no longer has free/paid tier constants.
- Subscriber status only has `active` and `inactive`.
- Requests and client response DTOs no longer expose tier as access semantics.
- Schema no longer indexes `artistKey, tier`.
- `gofmt -l` clean.
- `go test ./...` and `go vet ./...` pass from `backend/`.
- Report whether generated frontend types are now stale and what command should regenerate them.

## Reference reading

- `CONTEXT.md` Subscription, Subscriber, Follow, Relationships
- `docs/adr/0004-subscription-semantics.md`
- `CODING_STANDARDS.md` Backend and generated type sections
- `backend/internal/models/subscription.go`
- `backend/internal/database/schema.go`

## Effort estimate

1-3 hours. This should be model/schema cleanup, not a billing system. If it tries to become Stripe, stop. Stripe is an event horizon for scope.

## Bidirectional message paths

- Worker to overseer: `chat_msgs/CHAT_adr0004-subscription-semantics_msg.md`
- Overseer to worker: `chat_msgs/CHAT_adr0004-subscription-semantics_reply.md`
