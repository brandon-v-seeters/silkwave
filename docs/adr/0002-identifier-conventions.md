# Identifier conventions for Releases and Artists

Releases and Artists each carry three identifiers with strict role separation: the ArangoDB `_key` is internal-only and never appears in API paths or URLs; a stable external `id` (UUID, immutable for the document's lifetime) is used in every authenticated API path and in object-storage keys; a mutable `slug` is used only in public marketing URLs. The current `Hash` field on Release is renamed to `id` because the value is a UUID, not a hash — the original name lied about content-addressability that doesn't exist.

Slugs are unique **per Artist**, not globally. Public release URLs nest under the Artist: `/artists/:artistSlug/releases/:releaseSlug`. Two different Artists may each have a release slugged `live-at-the-roxy`. Artist slugs are globally unique because they form the top-level path component.

## Consequences

- Renaming the field from `Hash` to `id` does not require migrating R2 object keys — the keys embed the *value* (the UUID), not the field name. The rename is a code-level sweep.
- Public route shape changes from `GET /api/releases/:slug` to `GET /api/artists/:artistSlug/releases/:releaseSlug`. The prior global-slug route is dropped.
- Slug mutations (artist renames a release) should set up a redirect from the old slug to the new one to avoid breaking external links. The implementation of redirects is out of scope for this ADR.
