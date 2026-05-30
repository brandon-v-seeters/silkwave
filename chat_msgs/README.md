# Chat messages

Async coordination channel between worker chats and the overseer. The user is a relay, not an addressee.

## File pair convention

Per active worker chat:

| File                              | Direction         | Contents                                          |
| --------------------------------- | ----------------- | ------------------------------------------------- |
| `CHAT_<scope-slug>_msg.md`        | worker → overseer | Status, decision asks, plans-to-confirm, blockers |
| `CHAT_<scope-slug>_reply.md`      | overseer → worker | Drafted response, ready for the worker to act on  |

`<scope-slug>` is short and kebab-case. Examples: `releases-repo`, `archive-rename`, `auth-envelope`.

## Header format

Every message and every reply opens with a freshness header. The reply header references the message it answers.

```markdown
# 2026-05-10 14:32 — message #3

[worker's status / decision ask / plan-to-confirm]
```

```markdown
# 2026-05-10 14:35 — reply to message #3

[overseer's reply]
```

The "reply to message #N" line lets the worker verify the reply matches their last sent message. If `N` is older or the file is empty, the worker flags it instead of acting on stale context.

## Workflow

```
Worker writes msg.md (message #N)
  ↓
User → overseer: "chat <slug> reported"
  ↓ overseer reads the latest msg, drafts reply
  ↓ overseer writes reply.md (reply to message #N)
  ↓ overseer tells user: "reply written"
User → worker: "replied"
  ↓ worker reads reply.md
  ↓ worker verifies reply.message_id == worker.last_sent_id
  ↓ worker processes reply, continues work
```

## Reading discipline

Files accumulate every prior message over the session. Both worker and overseer should read only the latest header onward, never re-read from the top.

```bash
awk '/^# 2026/{print NR": "}' chat_msgs/CHAT_<slug>_msg.md | tail -3
# then Read with offset=<latest-header-line> and limit=<delta-to-next-header-or-EOF>
```

## Worker output discipline

Worker chats communicate with the overseer **only** through the file pair. User-facing output is a 1-2 paragraph milestone summary that ends with:

> Wrote message #N to CHAT_<slug>_msg.md, awaiting reply.

Workers never ask the user a question directly; every decision is routed to the overseer via `_msg.md`.
