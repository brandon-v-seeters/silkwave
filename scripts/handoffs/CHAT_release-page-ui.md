# Release Page UI Worker

Project context: Silkwave is a music distribution platform where independent Artists publish Releases for listeners to stream and buy. The user provided a Figma direction for the public Release page: a dark, immersive page using cover art as a blurred atmospheric backdrop, prominent cover art, large Release title, Artist identity, price and buy CTA, play CTA, track table, add/share actions, recommendations, and a persistent bottom player.

Stack location: frontend public Release route at `frontend/src/routes/artist/[artistSlug]/releases/[releaseSlug]/`.

## Scope

Own the visual implementation of the public Release page shell.

- Redesign `+page.svelte` toward the Figma direction while staying Silkwave: calm, spacious, premium, content-first.
- Reuse existing page data from `data.release`, `release.artist`, `release.tracks`, and route layout artist data.
- Use the existing `MediaPlayer` component if helpful; do not edit it.
- Build a useful track list with play affordances, duration, title, and small actions.
- Add a buy CTA that can safely route to existing checkout if there is an obvious path; otherwise keep it as an inert/disabled button or `href="/checkout"` only if that matches existing app routes, and report the commerce gap.
- Make mobile responsive. The Figma is desktop-heavy. Mobile must still feel intentional.
- Use semantic tokens and existing components where they fit. Avoid raw one-off color soup.

## File Boundaries

Your fair-game files:

- `frontend/src/routes/artist/[artistSlug]/releases/[releaseSlug]/+page.svelte`
- New files under `frontend/src/lib/components/organisms/public-release/` if extraction is truly useful.

DO NOT touch:

- `frontend/src/routes/artist/[artistSlug]/releases/[releaseSlug]/+page.server.ts`
- `backend/**`
- `frontend/src/lib/components/organisms/media-player/MediaPlayer.svelte`
- `frontend/src/lib/types/generated/models.ts`

Signature-stability rule: consume the existing load shape. If you need a new field from the contract worker, write a message instead of editing their files.

## Figma Intent To Preserve

- First viewport: cover art and Release identity dominate.
- Background: use cover art as a blurred/dimmed page atmosphere when available. Keep readability brutal and obvious.
- Header cluster: cover art left, title and Artist right on desktop; stacked on mobile.
- CTA cluster: price plus buy button, and a clear play button.
- Track list: dense enough to scan, not a dashboard landfill. Use columns where viewport allows, collapse cleanly on mobile.
- Recommendations: "Others also bought" can use available local data only if present. Do not fake API-backed recommendations. A polished empty/omitted state beats counterfeit content.
- Copy: use Release language in code and product copy. If the mock says "Buy Album", adapt to "Buy Release" unless the `releaseType` specifically equals `album`, where "Buy Album" is acceptable.

## Coordination

Message file: `chat_msgs/CHAT_release-page-ui_msg.md`
Reply file: `chat_msgs/CHAT_release-page-ui_reply.md`

Write plan, blockers, and scope changes to the message file only. On `replied`, read the reply file and verify it matches your latest message number.

## Hard Rules

- No MCP servers. Use local files, existing Svelte patterns, `rg`, and `npm run check`.
- Runes only. No `svelte/store`.
- Kebab-case filenames for any new components.
- No fake backend contracts. If data is missing, render gracefully and report the gap.
- No landing-page filler. This is the actual Release page.
- No em dashes in user-facing copy.
- Do not commit.

## Acceptance Criteria

- Public Release page visually tracks the Figma direction without breaking the current route.
- Dark mode reads correctly. Light mode should not be wrecked.
- Page works with and without cover art, tracks, price, description, and artist image.
- Text does not overlap or overflow on mobile.
- `npm run check` from `frontend` passes, or you report the exact blocker.

## Reference Reading

- `AGENTS.md`
- `CONTEXT.md`
- `CODING_STANDARDS.md`
- `frontend/src/routes/artist/[artistSlug]/releases/[releaseSlug]/+page.svelte`
- `frontend/src/routes/(app)/(protected)/upload/release/+page.svelte`
- `frontend/src/lib/components/organisms/media-player/MediaPlayer.svelte`
- `frontend/src/lib/components/ui/button/button.svelte`
- `frontend/src/lib/components/molecules/ReleaseCard.svelte`

## Effort Estimate

Pre-flight first. Expected scope: medium. The trap is inventing commerce or recommendations. Do not. Render what exists, report what does not.

