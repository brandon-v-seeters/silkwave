---
name: silkwave-overseer
description: Skill mirror of the Silkwave Overseer agent. Loads BRAND.md, routes brand work to the correct lieutenant skill domain (Voice, Curate, Design, Grow, Guard), and reviews output for brand consistency. Use when invoking Silkwave brand management as a function rather than as a conversational agent persona, or when calling from another skill or agent.
---

# Silkwave Overseer (skill form)

The function-call mirror of `.claude/agents/silkwave-overseer.md`. Same brain, different surface. Use this skill when:

- A Claude Code workflow needs to delegate brand work programmatically.
- Another skill or agent needs Silkwave brand management without spawning a full sub-agent.
- You want a stateless one-shot brand check on a piece of output.

For interactive, multi-turn brand work, prefer the Overseer agent.

## Quick start

1. **Read BRAND.md** at `.claude/skills/silkwave/BRAND.md`. This is the source of truth.
2. **Identify domain**: Voice, Curate, Design, Grow, or Guard.
3. **Invoke the lieutenant skill** at `.claude/skills/silkwave/<domain>/<skill-name>/SKILL.md`.
4. **Review output** against BRAND.md's hard rules and anti-patterns.

## Domains

| Domain | What it covers | Tier-1 skills |
|---|---|---|
| Voice | Written / spoken output, any surface | `voice/write-release-blurb` |
| Curate | Taste, A&R, calendar, submissions | `curate/shortlist-collaborators` |
| Design | Visual identity across surfaces | `design/brief-cover-art` |
| Grow | Acquisition, promo, funnel | `grow/plan-drop-announcement` |
| Guard | Trust, safety, originality | `guard/enforce-originality-policy` |

More skills under each domain land in tier 2 (rejections, newsletter, drop calendar, sample-pack audit, visual consistency audit, Ureka SEO articles, takedowns, etc.).

## Routing heuristics

- "Write a..." → Voice
- "Should we sign / pick / curate..." → Curate
- "Design / brief / mock / visual..." → Design
- "Promote / announce / launch / acquire..." → Grow
- "Is this AI / verify / takedown / policy..." → Guard

When a request crosses domains (e.g., a release announcement = Voice + Grow + Design), invoke each domain's skill and reconcile the output against BRAND.md's voice-register split.

## Hard checks before returning any output

- Public-layer copy reads dry / curator, not warm / partner.
- No "AI-powered" or equivalent anywhere in public copy.
- No hype emoji clusters in editorial copy.
- Visual decisions respect the two-context split (cyan-petrol for brand layer, existing yellow-green / indigo for product layer).
- Originality policy honored.

## See also

- `.claude/agents/silkwave-overseer.md` for the conversational persona.
- `.claude/skills/silkwave/BRAND.md` for the source of truth.
