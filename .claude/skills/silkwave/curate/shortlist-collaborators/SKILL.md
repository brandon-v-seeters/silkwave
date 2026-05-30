---
name: shortlist-collaborators
description: Generate a curated shortlist of 8 to 12 alt-electronic artists for Silkwave's borrowed-audience launch. Filters against the bass-belt + cinematic-flag-plant scope and the producer-craft sonic gravity well. Use when planning the year-1 launch roster, identifying potential featured artists, or expanding curatorial reach. Reads BRAND.md for scope and curatorial rules.
---

# shortlist-collaborators

Builds the year-1 launch roster: 8 to 12 hand-picked alt-electronic artists who fit Silkwave's curatorial filter and have audiences worth borrowing. Critical for cold-start move 5 (borrowed audience).

## Quick start

Inputs needed:

- Brandon's stated network (artists he has direct contact with, prioritized)
- Scene tags to focus on (default: bass belt: DnB, dubstep, UKG, neurofunk, jungle revival)
- Optional flag-plant slot for one cinematic / experimental artist
- Existing release catalog references (URLs, Spotify, Bandcamp, SoundCloud)
- Audience-size band (e.g., 2k to 50k followers; medium not megastar)

Output: a ranked shortlist with rationale per artist.

## Curatorial filter (from BRAND.md)

An artist fits Silkwave if:

1. **Sonic fit**: their work sits in the deep / dark / cinematic gravity well. Sound-design-forward, headphones-leaning, not festival or dancefloor-utility-only.
2. **Scene fit**: bass belt or cinematic / experimental. Not mainstream EDM, not commercial trance, not big-room.
3. **Audience momentum**: medium-tier following (roughly 2k to 50k). Big enough to bring people. Small enough that Silkwave is meaningful to them.
4. **Producer credibility**: makes the music themselves. No A&R-built ghosting. Patron-of-the-craft signal.
5. **Brand-fit personality**: serious about the work, not chasing virality. Would not say "AI-powered" or post fire emojis on every release.

Anti-fit signals:

- Catalog dominated by remixes for major-label pop.
- Heavy use of generative AI (Suno-style outputs in their own catalog).
- Audience built primarily on virality / TikTok rather than scene credibility.
- Public branding contradicts Silkwave's values.

## Workflow

1. Load BRAND.md curatorial sections (Scope, Sonic gravity well, Voice).
2. Take input list of candidate names. If none provided, prompt Brandon for his bass-belt network.
3. For each candidate, check:
   - Sonic fit against the gravity well (sample recent releases from last 18 months).
   - Scene placement.
   - Audience-size band.
   - Producer credibility (do they self-produce?).
   - Brand-fit personality (look at socials, release copy, recent collabs).
4. Rank candidates: A-tier (clear fit + meaningful audience pull), B-tier (good fit, smaller pull), C-tier (interesting but borderline).
5. Reserve one slot for a cinematic / experimental flag-plant if scope allows.
6. Output: 8 to 12 names with one-paragraph rationale each, plus 2 to 3 alternates.

## Output format

```
LAUNCH ROSTER SHORTLIST

A-TIER (must-have, high pull)
- Artist Name 1
  Why: [one paragraph: sonic fit, scene placement, audience signal, brand fit, contact channel]
  Audience: [follower estimate, scene visibility]
  Risk: [any anti-fit signal worth flagging]

B-TIER (strong fit, support roster)
- ...

CINEMATIC FLAG-PLANT
- One artist for the year-1 cinematic / experimental drop.

ALTERNATES
- 2 to 3 names to swap in if A-tier declines.
```

## Optional: leverage existing primitives

- Use `audit` to scan an artist's catalog systematically.
- Use `critique` to evaluate brand-fit on social presence.

## See also

- `.claude/skills/silkwave/BRAND.md` (curatorial scope, anti-patterns)
- Rejection-writing skill (when an artist doesn't fit) lands in tier 2.
