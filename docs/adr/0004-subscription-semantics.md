# Subscription semantics

Subscriptions follow the Bandcamp model, not the Patreon model: an active Subscriber gets unlimited streaming of the Artist's *entire current and future Release catalog*. There is no exclusive content, no Patreon-style tier-gated posts, no tier-locked Releases. An Artist may offer **multiple Subscriptions at distinct price points**, but all of them grant identical access — the tiers exist as patron-support levels for fans who want to pay more, not as different access tiers.

**Downloads always require a separate purchase.** A Subscription does not grant permanent ownership of any Release — when a Subscription ends, streaming access stops cleanly and the Subscriber retains only Releases they purchased outright. Subscriptions may carry an optional Subscriber-only purchase discount (a numeric percentage off any Release purchase), which is the value-add for paying Subscribers who want owned downloads. Free following ("get notified about new releases without paying") is modeled as a separate `Follow` edge between User and Artist — *not* as a free tier of Subscription, because conflating them forces every "is this a Subscriber?" check to also remember to filter by tier.

`SubscriptionStatus` collapses to **`active` / `inactive`** in Silkwave's domain. Stripe is the source of truth for the underlying billing-lifecycle states (`past_due`, `canceled`, `paused`, `trialing`, etc.); webhooks update the active/inactive flag, and the billing-admin UI queries Stripe directly when it needs the detailed reason. Mirroring Stripe's full 8-state enum into our DB would put us on the hook for keeping the mirror in sync, with no upside — `active` vs `inactive` is the only thing the access-check path needs.

## Consequences

- No `UserReleaseAccess` entity, no entitlement table, no per-User "library of Releases acquired via subscription." Catalog access is computed from the active Subscriber edge alone.
- "Subscribe for one month, bulk-download everything, cancel" is structurally impossible because subscription does not grant downloads.
- Pre-orders are out of v1 scope (see ADR-0001 — `publishAt` is a single field, not split into visibility-gate plus marketing-date). When pre-orders eventually land, they are a Release-purchase concept, not a Subscription concept.
