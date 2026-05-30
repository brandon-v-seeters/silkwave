---
name: write-release-blurb
description: Generate a Silkwave-flavored release announcement blurb in public-layer voice (tastemaker / dry, technical / craft-forward). Use when announcing a new release on Silkwave's homepage, drop page, newsletter, or social. Reads BRAND.md for voice rules.
---

# write-release-blurb

Writes the editorial blurb that accompanies every Silkwave release. The blurb shows up on the drop page, in the weekly newsletter, on social, and (in shortened form) on the homepage feed.

## Quick start

Inputs needed:

- Artist name and any relevant aliases
- Release title
- Release type (single, EP, album, remix, sample pack)
- Track count and total runtime
- BPM range and key signature(s) if relevant
- Sub-genre / scene placement (e.g., halftime DnB, deep dubstep, autonomic, IDM)
- 2 to 4 sentences from the artist describing the work, optional but useful
- Cover art reference, optional

Output: a 60 to 110 word editorial blurb in Silkwave public-layer voice.

## Voice rules (from BRAND.md)

**Tastemaker / dry / authoritative + technical / craft-forward.**

- Sparse. Earned every word.
- No exclamation marks unless genuinely surprising.
- No emoji.
- Producer-talk welcome: "Sub at 50Hz." "Halftime, 84 BPM." "Every drum wet."
- Open with a fact, not hype. "Four tracks. Halftime throughout. Out Friday."
- Close with a hook into listener experience or production craft, not a CTA.

**Do**

- "X returns with four halftime studies in sub-low decay. Drums dry, reverbs long, kicks tuned to the room."
- "A reissue, two new versions, one cut from the original 2014 tape."

**Don't**

- "We're SO HYPED to bring you this absolute banger from X."
- "Check out this AI-driven sonic exploration."
- "Out now on all platforms."

## Workflow

1. Load BRAND.md voice section if not already loaded.
2. Read inputs. If artist sentences are provided, distill them into one sentence of craft observation.
3. Open with a fact: format, count, BPM, scene placement.
4. Add one to two sentences of editorial color: sound design, mood, scene context, lineage.
5. Close with one sentence about the listener experience or a release detail (date, format, special edition, sample pack tie-in).
6. Cut to 60 to 110 words.
7. Final pass: any exclamation marks? Any emoji? Any "AI-powered"? Any hype words ("banger," "fire," "insane")? Strip them.

## Optional: leverage existing primitives

- Use `distill` to compress over-long drafts.
- Use `clarify` if a sentence is muddled.
- Use `edit-article` for longer-form drop essays (over 200 words).

## See also

- `.claude/skills/silkwave/BRAND.md` (voice rules, anti-patterns)
- Sister voice skills (rejection, newsletter, social) land in tier 2.
