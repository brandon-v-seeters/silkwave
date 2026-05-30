---
name: brief-cover-art
description: Generate a cover-art brief for a Silkwave release that fits the deep / dark / cinematic sonic gravity well and the editorial-minimal brand-layer visual identity. Use when commissioning external designers, generating AI-image prompts for moodboards, or auditing existing cover art for brand fit. Reads BRAND.md for visual rules.
---

# brief-cover-art

Produces a written brief for a Silkwave release cover. Used to commission external designers, generate moodboard images, or check that existing cover art fits the brand.

## Quick start

Inputs needed:

- Release metadata (artist, title, type, BPM range, sub-genre, key mood)
- 2 to 4 sentences from the artist describing the work
- Constraints (color preferences, deal-breakers, format requirements: square 3000x3000 default)
- Reference touchstones or anti-references from the artist if any

Output: a one-page brief with mood direction, type guidance, color guidance, composition guidance, and 3 to 6 visual reference touchstones.

## Visual rules (from BRAND.md)

This is brand-layer visual work.

- **Aesthetic**: editorial-minimal, Bauhaus-inflected. Hard grid. Type-led. Generous negative space.
- **Color base**: near-black backgrounds, off-white type. The single cold-saturated brand accent is **cyan-petrol** in the `oklch(0.55 0.10 210)` to `oklch(0.65 0.12 215)` range. Use sparingly, as accent only.
- **Type**: Cormorant Garamond for editorial display, Satoshi for technical / metadata text.
- **Photography or illustration**: muted, atmospheric, color-graded toward cool / cold. Avoid warm gold / orange dominance unless deliberately subverting the brand for a specific drop.
- **Forbidden**: hype emoji, gradient blurs that read as generic SaaS, "club flyer" aesthetics, AI-generated images that look like Midjourney defaults.

## Workflow

1. Load BRAND.md visual section.
2. Read release inputs. Distill the artist's description into a single sentence of mood.
3. Choose a primary visual register:
   - Type-led editorial cover (default for most releases).
   - Photographic atmospheric cover (when the music has strong place / setting).
   - Abstract / textural cover (for cinematic / experimental drops).
4. Write the brief with the following sections:
   - **Mood and concept** (one paragraph)
   - **Color guidance** (cyan-petrol accent allowed, base near-black, off-white type)
   - **Type guidance** (which face for what, hierarchy)
   - **Composition guidance** (grid, scale, negative space)
   - **Reference touchstones** (3 to 6 cover-art examples from existing scene)
   - **Anti-references** (1 to 3 things this cover should NOT look like)
   - **Deliverables** (square 3000x3000 master plus 1500x1500 web variant, plus any platform-specific sizes)
5. Send for review. Pass through `audit-visual-consistency` once that skill exists.

## Optional: leverage existing primitives

- Use `frontend-design` for layout / composition reasoning if mocking up directly.
- Use `colorize` for the accent treatment.
- Use `typeset` for type hierarchy.
- Use `arrange` for negative-space tuning.

## See also

- `.claude/skills/silkwave/BRAND.md` (visual identity, accent color, type)
- Visual consistency audit skill lands in tier 2.
