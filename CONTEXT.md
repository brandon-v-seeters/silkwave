# Silkwave

A music distribution platform where independent musicians upload, manage, price, and sell their releases directly to fans. Solo or small-team artists — Bandcamp/Patreon territory, not major-label distribution.

## Language

### Releases and tracks

**Release**: A piece of music an artist publishes — album, EP, single, compilation, remix, or live recording — carrying its own pricing, metadata, cover art, and Tracks.
_Avoid_: "album" as a synonym for Release (an album is a Release with `releaseType: "album"`).

**Track**: An individual recording belonging to a Release. Membership is fixed at first publish.

**Pending Edit**: Staged changes to a Published Release, saved but not yet live; visible only to the artist and their managers. Limited to metadata and track-level field overrides — never track membership or audio replacement.
_Avoid_: "draft" for an already-Published Release.

**Draft**: A Release that has never been Published. Once Published, a Release leaves Draft permanently.

**Reissue**: A new Release that references a prior Release as its source — the answer to "I want to add a track to a live album." Not yet implemented; flagged here so future requests get the right shape.

### Release lifecycle

A Release moves through three states:

- **Draft** — has never been Published. Hidden from the public.
- **Published** — has been Published at least once. Public-visible when `publishAt` is absent or in the past; future-dated `publishAt` keeps it scheduled-but-hidden until that moment. May carry a Pending Edit.
- **Archived** — was Published; the artist took it down. Hidden from public listings; identity and URLs preserved. Re-publishable back to Published.

Valid transitions:

- Draft → Published (first publish)
- Published → Archived (take down)
- Archived → Published (re-list)
- *Any* → *hard delete* (record removed; not a status)

There is no `scheduled` state — a future-dated publish is a **Published** Release with a future publish-time gate. There is no `deleted` status — deletion is a hard delete, not a soft-delete row.

### Artist and user

**Artist**: The creative entity that owns Releases — a public-facing identity with name, slug, bio, and catalog.
_Avoid_: "creator", "account", "profile".

**User**: An authenticated person. May manage zero or more Artists (one User can have multiple Artist projects, e.g. a main band and a side project); may subscribe to zero or more Artists. Each Artist has exactly one User who manages it — Silkwave is single-operator; there is no team / multi-user access on an Artist.
_Avoid_: "account".

**Subscription**: A paid plan an Artist offers their fans. An Artist may offer **multiple Subscriptions at different prices** — these are *patron-support levels*, not access tiers; access to the catalog is identical across all of an Artist's Subscriptions. While a Subscriber's Subscription is active, the Subscriber may **stream** any Release in the Artist's *entire current and future catalog* at no per-Release cost. Subscriptions do not unlock exclusive content — there are no Patreon-style gated posts or tier-locked Releases.

**Subscriber**: A User's active relationship to one of an Artist's Subscriptions. While the Subscription is active, the Subscriber gets unlimited streaming of the Artist's catalog. **Downloads always require purchase** — a Subscription does not grant permanent ownership of any Release. A Subscription may carry an optional Subscriber-only purchase discount (e.g. 20% off any Release). When the Subscription ends, streaming access stops cleanly; any Releases the Subscriber purchased separately remain owned.

A Subscriber's status is **active** or **inactive** — those are the only two states Silkwave models. Stripe is the source of truth for the underlying billing reasons (`past_due`, `canceled`, `paused`, etc.); webhooks update the active/inactive flag, and the billing-admin UI queries Stripe directly when it needs the detailed reason.

**Follow**: A free, opt-in relationship between a User and an Artist. A Follower receives release notifications and updates from the Artist but pays no money and gets no catalog access. Distinct from Subscriber — Follow is a separate concept, not a "free tier" of Subscription.
_Avoid_: "free subscriber" (a Follow is not a Subscription).

### Identifiers and URLs

- A **Release** has a stable external **id** (UUID, used in API paths) and a mutable **slug** (used in public URLs, scoped per Artist).
- An **Artist** has the same shape: stable external **id** plus mutable **slug**.
- Slugs are unique **per Artist**, not globally — two Artists may each have a Release slugged `live-at-the-roxy`. Public release URLs nest under the Artist: `/artists/:artistSlug/releases/:releaseSlug`.
- Internal database keys exist on every document but are never exposed in API paths or URLs.

## Relationships

- An **Artist** owns zero or more **Releases**
- A **Release** has one or more **Tracks** (membership fixed at first publish)
- A **Release** has zero or one **Pending Edit**
- A **User** manages zero or more **Artists** via a `UserArtist` edge; an Artist is managed by exactly one User. There are no roles — the relationship is binary (you manage this Artist, or you don't).
- An **Artist** offers zero or more **Subscriptions** at distinct price points; each **Subscriber** is a `(User, Subscription)` pair
- A **User** may **Follow** zero or more **Artists** — a separate edge, independent of any Subscription

## Example dialogue

> **Dev:** "An artist wants to add a bonus track to their EP that's been live for two months — can they edit the track list?"
> **Domain expert:** "No. Track membership is fixed at publish. They put out a deluxe **Reissue** that references the original."
>
> **Dev:** "What about a typo in a track title?"
> **Domain expert:** "That's a **Pending Edit**. They change it, the old title stays live to the public until they publish the edit."

## Flagged ambiguities

- `ReleaseStatus` currently has `draft`, `scheduled`, `published`, `archived`, `deleted`, plus a separate `Published` boolean and an `IsUploaded` boolean. **Resolved:** the canonical state machine is `Draft → Published → Archived`; `scheduled` and `deleted` are dropped (see Release lifecycle above); the `Published` boolean is redundant with the status and will be removed; the `Releases`/`ReleaseDrafts` collection split collapses into a single `Releases` collection now that Pending Edits are embedded.
- `UnpublishRelease` currently transitions a Release back to `Draft`. **Resolved:** the operation is renamed `ArchiveRelease` and transitions to **Archived**, not Draft — Draft means *never Published*.
- `Hash` vs `Slug` vs `_key` on a Release — three identifiers without clear separation of purpose. **Resolved:** the `Hash` field is renamed `id` (it's a UUID, not a hash; the old name lied), `_key` is internal-only and never appears in API paths, and `slug` is per-Artist-scoped for public URLs. The same shape applies to Artist. R2 object keys are unchanged because they embed the UUID *value*, not the field name.
- The current `releaseDate` field both gates visibility and serves as the displayed marketing date. **Resolved:** rename to `publishAt`; treat absent-or-past as "live now," future as scheduled. No separate marketing-date field; pre-orders are deliberately out of v1 scope.
- `SubscriptionStatus` mirrors Stripe's 8-state lifecycle. **Resolved:** collapse to `active` / `inactive` in Silkwave's model; Stripe remains the source of truth for billing-reason detail, queried on demand.
