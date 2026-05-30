# Handoffs

Per-chat scope packages. One markdown file per worker chat at spawn time.

Filename: `CHAT_<scope-slug>.md`. The slug matches the chat's `chat_msgs/CHAT_<slug>_msg.md` and `_reply.md` pair.

## Required sections

- **Project context** — one paragraph: Silkwave (SvelteKit 5 + Go 1.25 / ArangoDB / R2) and the specific feature.
- **Stack location** — exact paths the chat will touch (`backend/internal/...`, `frontend/src/lib/...`).
- **Scope** — what to do; what sister chats are doing in parallel.
- **File boundaries** — fair-game / do-not-touch / signature-stability rule.
- **Hard rules** — the standard project-wide constants list.
- **Acceptance criteria** — tests required, lint clean, no new legacy drift, envelope responses, etc.
- **Reference reading** — `CONTEXT.md` sections, relevant ADRs, canonical source files.
- **Effort estimate** — with the 5-tier pre-flight discipline note.
- **Bidirectional msg/reply paths** — `chat_msgs/CHAT_<slug>_msg.md` and `_reply.md`.

The kickoff prompt the user pastes into the new chat lives next to the handoff or inline in the spawn message; either way it references this file.
