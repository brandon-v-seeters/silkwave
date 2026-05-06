# Single-operator Artist ownership

Each Artist is managed by exactly one User. There are no team / multi-user roles on an Artist — the `role` field on the `UserArtist` edge (`owner`, `manager`, `member`) is dropped because no code reads it and the indie-self-distribution product persona does not have collaborator access patterns. A single User may still manage multiple Artists (e.g. a main band plus a side project).

The `UserArtist` edge structure is kept rather than collapsed into an `Artist.userKey` foreign key, because it preserves a clean future migration path to multi-user team management or explicit Artist ownership transfer if either becomes a real product need. Ownership transfer (handing an Artist record from one User to another, e.g. for catalog hand-offs or estate handovers) is not implemented in v1; if a real user case arises before the feature lands, it can be handled with a one-off DB edit (rewrite the edge's `_from`).
