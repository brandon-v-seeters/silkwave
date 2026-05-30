# Silkwave Brand

The single source of truth for how Silkwave behaves and signals.

This file does not define product domain (releases, tracks, subscriptions, etc.). For that, see `/CONTEXT.md`. This file defines what the brand sounds like, looks like, and curates against. Every Silkwave skill at `.claude/skills/silkwave/**` reads this file to stay aligned.

## Brand statement

Silkwave is a brand house with curatorial divisions. A marketplace, community, and label-shaped curator for the alternative electronic music scene. Bandcamp meets Patreon meets Splice in shape. Hyperdub meets Hospital Records in attitude.

Sister product: **Ureka by Silkwave**. A public AI tool that helps producers move past writer's block. Top-of-funnel acquisition for Silkwave.

## Positioning

| | |
|---|---|
| What we are | Curated marketplace for human-originated alternative electronic music. Tracks, samples, subscriptions. |
| What we are not | A general distributor. Not Bandcamp scale. Not DistroKid. Not Spotify. |
| Who we serve | Producers and listeners in the producer-craft side of alt-electronic: bass belt (DnB, dubstep, UKG, neuro, jungle) and cinematic / experimental (ambient, IDM, breakcore, glitch). |
| What separates us | Curatorial discipline, sound-design-forward sensibility, human-only audio policy, and Ureka as the creative-process tool that pulls producers in. |

## Scope

**Launch fronts (year 1)**
- (a) Bass belt: DnB, dubstep, UKG, neurofunk, jungle revival.
- (c) Cinematic / experimental: one signature flag-plant drop in the first 90 days.

**Year 2 expansion**: sub-fronts under the Silkwave umbrella with their own personalities. Silkwave Bass, Silkwave Cinema, eventually Club and Bedroom.

**Out of scope**: mainstream EDM, festival big-room, commercial trance, hardstyle, anything where the production tooling is the entire identity.

## Voice (two-register)

Silkwave runs two voice registers depending on the surface.

### Public brand layer

Used on: homepage, drops, editorial, social, marketing copy, Ureka chatbot.

**Tastemaker / dry / authoritative + technical / craft-forward.**

- Sparse. Earned every word.
- No exclamation marks unless genuinely surprising.
- No emoji in editorial copy. Sparingly in social, only as scene-coded sigils.
- Producer-talk welcome. "Sub at 50Hz, 808s through analog drive."
- "We don't sign that" is preferable to "we'd love to receive your submission."
- Ureka inherits this voice **slightly warmer**: a friend who happens to know sound design at a serious level.

**Do**
- "New from X. Four tracks. Out Friday."
- "A study in sub-low decay."
- "Halftime, 84 BPM, every drum wet."

**Don't**
- "We're SO excited to drop this banger! 🔥🔥🔥"
- "AI-powered platform for music creators."
- "Unlock your potential as an artist."

### Product layer

Used on: artist dashboard, signup, settings, account flows, payment, billing.

**Calm / trusted-partner. Quiet confidence. Premium without flashy.**

This register comes from the original CLAUDE.md design context and stays. Artists managing their catalog need warmth and reassurance, not curatorial dryness.

- Friendly but not chatty.
- Empowering, not patronizing.
- "Your release is live" beats "Hyped that your drop is out."

The two registers share a backbone: no corporate softness, no fake hype, no AI bragging.

## Sonic gravity well

**Deep, dark, cinematic.** Sub-low energy, restrained drums, mood-led, sound-design-forward. Headphones at 2am, not festival mainstage.

When in doubt about a track, sample, or curatorial call, ask: is this a headphones-at-2am piece, or is this a dancefloor-utility piece? Headphones wins by default. Dancefloor utility wins only when the sound design holds up.

Reference touchstones (not exhaustive): Burial, Photek, Ivy Lab, Squarepusher, Amon Tobin, Boards-of-Canada-adjacent textures, autonomic-era DnB, halftime, deep dubstep, restrained 140.

## Visual identity (two-context split)

### Brand layer

Used on: homepage, marketing pages, drop assets, editorial, social.

- **Aesthetic**: editorial-minimal, Bauhaus-inflected. Hard grid. Type-led. Generous negative space.
- **Color**: near-black base (around `oklch(0.12 0 0)`), off-white type, **single cold-saturated accent: cyan-petrol** in the `oklch(0.55 0.10 210)` to `oklch(0.65 0.12 215)` range. Use sparingly.
- **Type**: Cormorant Garamond for editorial display, Satoshi for technical / metadata.

### Product layer

Used on: artist dashboard, signup, settings, account, billing.

- **Aesthetic**: producer-toolkit functional. Clean UI. Calm. Per the original CLAUDE.md design context.
- **Color**: existing palette. Yellow-green primary in dark mode, deep indigo primary in light mode. Muted, sophisticated, OKLch.
- **Type**: same Cormorant + Satoshi pairing.

The cyan-petrol accent **bridges** the two layers by appearing at brand touchpoints inside the product (publish-success states, drop highlights, primary CTA on release pages) without taking over the daily UI.

## Cold-start strategy (year 1)

Stack **5 + 1 + 4** in priority order.

1. **Borrowed audience**. Brandon's existing scene network. Launch with own releases plus 8 to 12 hand-curated bass-belt collaborators. Each brings a following.
2. **Curator-led seeding**. First 50 drops are hand-picked by Brandon. No self-serve until year 2. Submissions queue exists, heavily gated.
3. **Ureka-led tool funnel**. Ureka launches in parallel as a free public AI tool. SEO-friendly content. Quietly funnels producers into Silkwave membership over months.

**Year-1 success metric**: editorial trust. Did Silkwave become a thing scene people mention without prompting? Not artist count. Not GMV.

## Originality policy

Silkwave is a marketplace for **human-originated music**. All audio submitted (tracks, samples, stems, sub-tier content) must be substantially originated by human creative decision-making and human performance or programming.

**Allowed as production assistants**: mastering (LANDR, Ozone), mixing (iZotope), pitch correction (Antares), stem separation, sound design synthesis, denoise.

**Prohibited**: generative text-to-music platforms (Suno, Udio, and equivalents) and their outputs.

**Enforcement, three layers**

1. **Attestation at upload**. Required checkbox confirming human origin per this policy. Legally meaningful.
2. **Curatorial review**. Year 1: every release passes through Brandon or a trusted A&R sub-agent. Year 2+: open submissions still pass curatorial review.
3. **Best-effort AI detection signals**. Tools like Ircam Amplify, Pex, etc. One signal among many. Never the verdict alone.

Failure to disclose AI origin results in removal and ban.

## AI disclosure stance

| Surface | Disclosure |
|---|---|
| Ureka, the public chatbot | Loud and proud. "AI sparring partner for producers." |
| Silkwave's internal tooling (release copy, scheduling, ops) | Silent. Standard 2026 practice, not deception. |
| AI-generated music submitted to marketplace | Banned. See originality policy. |

**Never use** in public copy: "AI-powered," "powered by AI," "AI-driven," "smart" (when meaning AI), "intelligent" (when meaning AI). Silkwave does not market AI. Ureka markets Ureka.

## Tagline

> **Ureka helps you find the idea. Silkwave helps you sell it. The music in between is yours.**

The public bridge between the two products. Works on a t-shirt, in a Spotify ad, and in a dispute. Encodes the originality policy in eight words.

## Anti-patterns (never do)

- Hype emoji clusters in editorial copy.
- "AI-powered" anywhere on Silkwave's brand surfaces.
- Race-to-the-bottom fee marketing ("lower fees than Bandcamp").
- Generic "for all electronic music" positioning.
- Self-serve uploads in year 1.
- Yellow-green or indigo as the brand-layer accent (those belong to the product layer).
- Cyan-petrol used as a primary in the product layer (it's a brand touchpoint, not a daily-use color).
- Calling artists "creators," "users" (in artist-facing copy), "accounts," or "profiles."

## How skills use this file

Every Silkwave skill at `.claude/skills/silkwave/**` references this file as the source of truth. The Overseer agent at `.claude/agents/silkwave-overseer.md` is the only entity that proposes edits, with Brandon's approval. When a skill output feels off, the answer is here.
