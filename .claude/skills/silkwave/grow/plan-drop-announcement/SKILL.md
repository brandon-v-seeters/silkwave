---
name: plan-drop-announcement
description: Plan a cross-platform release announcement for a Silkwave drop. Coordinates messaging across IG, X, TikTok, Reddit, Discord, newsletter, and the Silkwave drop page. Maintains public-layer voice and adapts per-platform cadence without breaking brand. Use when launching a new release, EP, or sample pack. Reads BRAND.md for voice and growth strategy.
---

# plan-drop-announcement

Builds the cross-platform promo plan for a Silkwave release. Aligns release-day messaging across the channels Silkwave runs and pre-empts the common platform-tone trap (where the same brand ends up sounding like four different companies on four different feeds).

## Quick start

Inputs needed:

- Release metadata (artist, title, type, drop date, drop URL on Silkwave)
- The release blurb (from `voice/write-release-blurb`)
- Cover art (final, 3000x3000)
- Audio preview clip (15 to 30s, ideally a moment of sound design rather than a drop)
- Artist's social handles
- Optional: launch budget (paid promo or organic only)

Output: a 7-day announce-plan with per-platform copy, scheduling, and asset list.

## Voice and channel rules (from BRAND.md)

Public-layer voice everywhere: tastemaker / dry + technical / craft-forward. Slight per-platform inflection allowed.

| Platform | Inflection | Cadence |
|---|---|---|
| Silkwave drop page | Pure editorial. Full blurb. | Day 0 |
| Newsletter | Editorial + behind-the-scenes detail. | Day 0 morning |
| IG / Threads | Sparse, image-led. Caption short, technical. | Day -3, Day 0, Day +3 |
| X | Drier still. Single line. Link. | Day 0, Day +1 |
| TikTok | Audio-led. Caption sparse. Use the sound-design moment, not the drop. | Day -1, Day 0 |
| Reddit (relevant subs only) | Self-promo rules respected. Editorial framing. | Day 0 |
| Discord (Silkwave + scene servers) | Native, conversational, but still curatorial. | Day 0 |

**Never** use hype emoji clusters, "AI-powered," or "OUT NOW" phrasing. A single emoji is fine when it's a deliberate scene-coded sigil (a wave, a sub-bass character, etc).

## Workflow

1. Load BRAND.md voice and grow sections.
2. Read inputs. Pull the release blurb. Confirm cover art and preview clip exist.
3. Build the per-platform copy:
   - **Silkwave drop page**: blurb + tracklist + sample pack tie-in if relevant.
   - **Newsletter**: blurb + 2 to 3 sentences of behind-the-scenes context.
   - **IG**: 3 posts. Day -3 cover-art reveal. Day 0 release post with 60-second preview. Day +3 detail post (production, gear, scene context).
   - **X**: 2 posts. Day 0 announce. Day +1 production-detail thread.
   - **TikTok**: 2 videos. Day -1 sound-design teaser. Day 0 release video with link in bio.
   - **Reddit**: one post in 1 to 2 relevant subs (e.g., r/DnB, r/dubstep, r/IDM, r/ambient). Editorial framing, not "check out my track."
   - **Discord**: native message in Silkwave's own server plus, with permission, in 1 to 2 scene-friend servers.
4. Schedule the calendar (7 days: -3, -2, -1, 0, +1, +2, +3).
5. Asset list: which image, which clip, which copy variant goes where.
6. Final pass against BRAND.md anti-patterns: hype emoji, "AI-powered," exclamation-mark spam.

## Optional: leverage existing primitives

- Use `clarify` to tighten any caption that drifts long.
- Use `distill` to cut the blurb down for short-form posts.
- Use `delight` cautiously: a small touch of personality is fine, never at the expense of dryness.

## See also

- `.claude/skills/silkwave/BRAND.md` (voice, growth strategy)
- `.claude/skills/silkwave/voice/write-release-blurb` (source blurb)
- Per-platform short-form skill lands in tier 2.
