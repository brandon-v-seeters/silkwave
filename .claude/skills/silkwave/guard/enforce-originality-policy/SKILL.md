---
name: enforce-originality-policy
description: Run the three-layer originality enforcement (attestation review + curatorial check + AI-detection signal interpretation) on a Silkwave marketplace submission. Decides whether a submission is approved, flagged for human review, or rejected. Use when reviewing any uploaded audio (track, sample, stem, sub-tier content) before it goes public. Reads BRAND.md for the policy text.
---

# enforce-originality-policy

Runs the three enforcement layers from BRAND.md against a marketplace submission and produces a verdict: **approve**, **flag for human**, or **reject**.

This skill does not redefine the policy. The policy text lives in `BRAND.md` (single source of truth). This skill operates against that text.

## Quick start

Inputs needed:

- Submission audio (or URL / storage ref)
- Artist attestation (true / false on the human-origin checkbox)
- Submission metadata (title, BPM, key, sub-genre, artist's stated tools / DAWs)
- Optional: AI-detection signal output (Ircam Amplify, Pex, or whatever detector is wired in)
- Optional: artist's submission notes / tools-used disclosure

Output: a verdict with reasoning, plus a next-action recommendation.

## Policy summary (from BRAND.md)

Silkwave is a marketplace for human-originated music. AI is allowed as a production assistant (mastering, mixing, pitch correction, stem separation, sound design synthesis, denoise). AI is prohibited as a generator (Suno, Udio, equivalents). Failure to disclose AI origin equals removal and ban.

## Three-layer enforcement

### Layer 1: Attestation

- Checkbox unchecked → **reject** with notice to artist that attestation is required.
- Checkbox checked but submission notes contradict it (e.g., artist says "made with Suno") → **reject** with policy citation, ban account.
- Otherwise → continue to Layer 2.

### Layer 2: Curatorial review

In year 1, every submission passes through Brandon or a trusted A&R sub-agent. This skill prepares the curatorial review packet:

- Attestation status (passed)
- Submission metadata
- Artist tool disclosure
- Audio characteristics summary (BPM, key, structure, sound design fingerprint)
- Detection signal (if available)
- Curatorial fit notes (sonic gravity well match)

The curator (Brandon) makes the call. Year 2+: this skill may be allowed to auto-approve A-tier signal, escalating only B and C.

### Layer 3: AI-detection signals (best effort)

If a detector is wired in:

- High-confidence Suno / Udio signal → **reject** with citation, ban account.
- Medium-confidence signal → **flag for human**, do not auto-reject. Detection is imperfect.
- Low or no signal → no action; proceed on attestation + curation alone.

Detection alone is never the verdict.

## Workflow

1. Load BRAND.md originality section.
2. Read submission. Check attestation.
3. Attestation fails → reject, notify artist, log incident.
4. Attestation passes → run detection signal if available, prepare curatorial packet.
5. Surface to curator (Brandon).
6. Apply curator decision.
7. On reject for AI origin: ban account per policy. On reject for curatorial fit (not AI): polite rejection, no ban (use the rejection-writing skill once it exists).

## Output format

```
SUBMISSION VERDICT

Submission: [title by artist]
Attestation: [pass | fail]
Detection: [signal level or N/A]
Curatorial fit: [pending | pass | fail]
Verdict: [approve | flag for human | reject]
Reason: [one paragraph]
Next action: [publish | route to Brandon | notify artist | ban]
```

## Optional: leverage existing primitives

- Use `audit` for systematic packet preparation.

## See also

- `.claude/skills/silkwave/BRAND.md` (full originality policy, AI disclosure stance, anti-patterns)
- Rejection-writing skill lands in tier 2.
