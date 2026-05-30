# 2026-05-12 13:43 CEST - reply to message #0

Proceed with pre-flight only. No code changes.

Focus on issue #3: internal Release identity and publish-field cleanup.

Hard limits:

- Write only `chat_msgs/CHAT_issue-3-release-identity_msg.md`.
- Do not touch frontend, generated TypeScript, GitHub issue state, storage paths, public URL shape, Subscription, Follow, access-module, or lifecycle extraction code.
- Do not rename anything yet.

Required pre-flight:

- Read `CONTEXT.md`, ADR-0001, ADR-0002, `CODING_STANDARDS.md`, and the handoff.
- Scan backend references for `Hash`, `ReleaseDate`, `Published`, `IsUploaded`, `ReleaseStatusScheduled`, `ReleaseStatusDeleted`, `release.hash`, `release.releaseDate`, `release.published`, and `release.isUploaded`.
- Identify exact implementation files and whether this can be one worker pass.
- Check whether `tygo` can run.

Next message: append message #1 with findings, exact files needed, risk split, validation plan, and whether implementation should proceed as one pass or be decomposed.
