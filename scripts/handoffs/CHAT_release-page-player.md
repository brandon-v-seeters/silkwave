# Release Page Player Worker

Project context: Silkwave needs a public Release page where listeners can play tracks and eventually buy the Release. A Figma mock shows a persistent bottom player with cover art, current track, transport controls, progress, volume, and a buy-song CTA. There is already a reusable `MediaPlayer` component.

Stack location: `frontend/src/lib/components/organisms/media-player/MediaPlayer.svelte`.

## Scope

Own improvements to the existing MediaPlayer so the Release page can use it cleanly.

- Inspect current `MediaPlayer` behavior and props.
- Improve visual fit for the public Release page: calm, compact, dark-mode friendly, bottom-player capable.
- Preserve current props unless a small additive prop is clearly needed.
- Keep `placement="floating"` and `placement="inline"` working.
- Improve disabled/no-source behavior if tracks lack audio URLs.
- Do not wire route-specific data. The route UI worker owns page composition.

## File Boundaries

Your fair-game files:

- `frontend/src/lib/components/organisms/media-player/MediaPlayer.svelte`
- `frontend/src/lib/components/organisms/media-player/index.ts` only if an export needs updating.

DO NOT touch:

- `frontend/src/routes/artist/[artistSlug]/releases/[releaseSlug]/+page.svelte`
- `frontend/src/routes/artist/[artistSlug]/releases/[releaseSlug]/+page.server.ts`
- `backend/**`
- Generated types.

Signature-stability rule: do not remove or rename existing props. Additive props are allowed only if you keep defaults and report them in `CHAT_release-page-player_msg.md`.

## Coordination

Message file: `chat_msgs/CHAT_release-page-player_msg.md`
Reply file: `chat_msgs/CHAT_release-page-player_reply.md`

Use the file pair for all decisions and blockers. On `replied`, read the reply file and verify the reply number.

## Hard Rules

- No MCP servers.
- Runes only. No `svelte/store`.
- No raw route-specific imports.
- No fake audio behavior. If no source exists, controls must communicate disabled state accessibly.
- No new global CSS unless there is no component-local alternative.
- Do not commit.

## Acceptance Criteria

- Existing MediaPlayer remains usable with the current exported `MediaPlayerTrack` shape.
- Floating placement can visually serve the Figma-style bottom player.
- Inline placement still works for future embeds.
- Keyboard/focus states remain visible.
- `npm run check` from `frontend` passes, or you report the exact blocker.

## Reference Reading

- `AGENTS.md`
- `CODING_STANDARDS.md`
- `frontend/src/lib/components/organisms/media-player/MediaPlayer.svelte`
- `frontend/src/routes/artist/[artistSlug]/releases/[releaseSlug]/+page.svelte`

## Effort Estimate

Pre-flight first. Expected scope: small-to-medium. If page-level orchestration is required, stop and report instead of reaching into the route worker's files.

