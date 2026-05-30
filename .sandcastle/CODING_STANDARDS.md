# Silkwave Coding Standards

Rules for writing code in this repo. The list is short on purpose. If a rule is missing, fall back to the official SvelteKit / Svelte 5 docs for the frontend and Effective Go + the Go style guide for the backend.

This is a living spec. The codebase has not caught up yet. See **Legacy and migration** at the bottom for the gap.

Domain language lives in `CONTEXT.md`. Past architectural decisions live in `docs/adr/`. Brand and design context lives in `AGENTS.md`. Code conventions live here.

---

## Frontend

### Stack

- SvelteKit, Svelte 5 (runes only)
- TypeScript, strict
- Tailwind CSS 4 with `@theme inline`
- Bits UI for headless primitives, shadcn-svelte derivatives in `lib/components/ui/`
- `sveltekit-superforms` + Zod for forms
- `tailwind-variants` for component variants

### Files and folders

```
frontend/src/lib/
├── api/             # network helpers, fetch wrappers
├── assets/          # static assets imported by code
├── components/
│   ├── ui/          # shadcn-svelte primitives, kebab-case
│   ├── atoms/       # smallest custom units (Icon, Alert, Link)
│   ├── molecules/   # composed atoms (UserAvatar, ImageDropzone)
│   └── organisms/   # page-level composites (MainNavbar, release-wizard/)
├── constants/
├── hooks/
├── services/        # cross-cutting client-side services
├── state/           # *.svelte.ts rune modules (replaces lib/stores/)
├── schemas/         # Zod schemas, only when shared by 2+ routes
├── types/           # generated.ts from tygo, plus hand-written types
└── utils/
```

The four-bucket component layout is the rule:

- **ui/**: third-party-derived primitives. Treat as vendored. Don't refactor for taste.
- **atoms/**: one job, no children of substance. `icon`, `alert`, `link`.
- **molecules/**: a small composition with internal state but no page-level concerns.
- **organisms/**: page sections. Owns its own context, its own `.svelte.ts` state module if it needs one.

Promotion rule: if an atom grows props beyond styling and one slot, it's a molecule. If a molecule starts orchestrating data fetching, it's an organism. Move it.

### Component naming

**Files: kebab-case. Always.** `release-card.svelte`, not `ReleaseCard.svelte`. Multi-file components live in their own folder: `release-wizard/index.ts`, `release-wizard/wizard.svelte.ts`.

The Svelte component identifier on import stays PascalCase. That's a Svelte requirement, not a style choice:

```ts
import ReleaseCard from "$lib/components/molecules/release-card.svelte";
```

Folder names are also kebab-case: `release-wizard/`, `release-editor/`.

### Reactivity: runes only

No `writable`, `readable`, or `derived` from `svelte/store` in app code. The Svelte 5 idiom is rune modules.

Shared reactive state lives in `lib/state/<name>.svelte.ts` and exposes a factory or singleton:

```ts
// lib/state/cart.svelte.ts
function createCartState() {
  let items = $state<CartItem[]>([]);
  const total = $derived(items.reduce((sum, i) => sum + i.price, 0));

  function add(item: CartItem) {
    items = [...items, item];
  }
  function remove(id: string) {
    items = items.filter((i) => i.id !== id);
  }

  return {
    get items() {
      return items;
    },
    get total() {
      return total;
    },
    add,
    remove,
  };
}

export const cartState = createCartState();
```

For per-tree state use `setContext` / `getContext` with a typed factory. See `lib/components/organisms/release-wizard/wizard.svelte.ts` for the canonical pattern.

Inside components: `$state`, `$derived`, `$props`, `$effect`. Never `$:` reactive statements (that's Svelte 4).

### Forms: SuperForms is mandatory for any server-bound form

Any `<form>` that POSTs to a `+page.server.ts` action or to the API uses `sveltekit-superforms` plus a Zod schema. No exceptions on the happy path. Pure-client UI (search filter, toggle, sort dropdown) does not need SuperForms.

Schemas live inline in the route's `+page.server.ts` by default. The moment a second route imports the same schema, move it to `lib/schemas/<name>.ts` and import from both. Don't pre-extract.

Server action skeleton:

```ts
// +page.server.ts
import { z } from "zod";
import { superValidate, fail } from "sveltekit-superforms";
import { zod } from "sveltekit-superforms/adapters";

const schema = z.object({
  email: z.string().email(),
  password: z.string().min(8),
});

export const load = async () => ({ form: await superValidate(zod(schema)) });

export const actions = {
  default: async ({ request }) => {
    const form = await superValidate(request, zod(schema));
    if (!form.valid) return fail(400, { form });
    // ...
    return { form };
  },
};
```

Client form skeleton: see `lib/components/organisms/UserAuthForm.svelte` for the `superForm` + `Form.Field` + `Form.Control` pattern. New forms follow that shape.

### TypeScript

`strict: true`. `any` is banned. If the type is genuinely unknown, use `unknown` and narrow before use. If the type is "I don't want to write it," write it.

Two narrow exceptions:

- Vendored shadcn-svelte code in `lib/components/ui/`. Don't fight it.
- Generic helper utilities where `any` is the only honest answer (e.g., `Component<any>` in flex-render).

Prefer generated types from `lib/types/generated.ts` (tygo output) over hand-rolled mirrors of backend models.

### Styling

- Tailwind utilities first.
- For components with variants (size, intent, state) use `tailwind-variants`. See `lib/components/ui/button/button.svelte` for the canonical `tv()` setup.
- Use `cn()` from `lib/utils/utils.ts` to merge classes from props. Always pass `className` last so the consumer wins.
- Theme tokens live in `@theme inline` blocks. Use semantic tokens (`bg-foreground`, `text-primary`) over raw OKLch values.
- Dark mode is a class toggle, not a media query. Both modes ship.

### Imports and aliases

`$lib/...` for everything inside `frontend/src/lib`. No relative paths longer than two segments (`../../`) crossing folder boundaries. If you reach for `../../../`, the import is wrong, not the path.

Barrel files (`index.ts`) only when a folder has a stable public surface (`lib/components/ui/<primitive>/index.ts`). Don't barrel atoms or organisms; import from the file directly.

---

## Backend

### Stack

- Go 1.25
- Gin
- ArangoDB (`arangodb/go-driver`)
- Cloudflare R2 via S3-compatible client
- JWT + bcrypt
- Zap for logging
- `go-playground/validator` (via Gin's binding)
- tygo for TS type generation

### Package layout

Standard layout under `backend/internal/`:

```
auth/         JWT, password, session services
config/       env loading
database/     ArangoDB client, query helpers, migrations
handlers/     HTTP request handlers, one file per resource
logger/       Zap setup
middleware/   auth and request middleware
models/       data structs and DTOs (input/output shapes)
repository/   data access. ALL DB queries live here.
routes/       wiring
storage/      R2 client
```

`cmd/server/main.go` and `cmd/migrate/main.go` are the two binaries. New entrypoints get their own `cmd/<name>/`.

### Repository pattern

**Handlers do not touch ArangoDB. Period.**

All queries, mutations, and AQL live in `repository/<entity>_repository.go`. Handlers receive a repo via DI in `routes.Setup`. If the resource has no repo yet, create one before writing the handler.

Repo method shape:

```go
func (r *ReleaseRepository) GetBySlug(ctx context.Context, artistKey, slug string) (*models.Release, error) {
    q := /*aql*/ `
        FOR r IN Releases
        FILTER r.artistKey == @artistKey AND r.slug == @slug
        LIMIT 1
        RETURN r
    `
    return database.QueryOne[models.Release](ctx, r.db, q, map[string]interface{}{
        "artistKey": artistKey,
        "slug":      slug,
    })
}
```

The `/*aql*/` comment is a hint for the AQL extension; keep it.

### Request validation

Use struct tags + Gin's `ShouldBindJSON`. Never hand-roll empty-string checks after binding.

```go
type LoginUserRequest struct {
    Email    string `json:"email"    binding:"required,email"`
    Password string `json:"password" binding:"required,min=8"`
}

func (h *AuthHandler) Login(c *gin.Context) {
    var req models.LoginUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        respondError(c, http.StatusBadRequest, "invalid_request", err.Error())
        return
    }
    // ...
}
```

Custom rules that struct tags can't express (cross-field constraints, business rules) live in the service or repo, not the handler. Validation that has to hit the DB (uniqueness, existence) is a service-layer concern.

### Response shape: envelope

Every JSON response uses one shape:

```go
type Response[T any] struct {
    Data    *T          `json:"data,omitempty"`
    Error   *ErrorBody  `json:"error,omitempty"`
    Message string      `json:"message,omitempty"`
}

type ErrorBody struct {
    Code    string `json:"code"`              // stable machine string: "invalid_request", "unauthorized", "not_found"
    Message string `json:"message"`           // human-readable
    Details any    `json:"details,omitempty"` // optional per-field validation map
}
```

Helpers in `internal/handlers/response.go` (create if it doesn't exist):

```go
respondOK(c, data, "Login successful.")
respondError(c, http.StatusBadRequest, "invalid_request", "Email is required.")
```

Frontend always reads `data` for success and `error` for failures. No more "sometimes a `User`, sometimes `{message: ...}`."

### Errors and logging

- Wrap with `%w` at every layer below the handler:
  ```go
  return fmt.Errorf("create user: %w", err)
  ```
- Log at the handler boundary only. Repo and service code returns errors. They do not call `logger.Error`.
- Use Zap's structured fields, never string concatenation:
  ```go
  logger.Error("login failed", err, zap.String("email", req.Email))
  ```
- `fmt.Println`, `log.Println`, and bare `println` are banned in shipped code. Tests and one-shot scripts in `cmd/` are fine.
- Don't reinvent stdlib: use `strings.Contains`, `errors.Is`, `errors.As`. Hand-rolled string helpers (looking at `containsHelper` in `auth_handler.go`) are deletion targets.

### Logging conventions

- One log line per failed request, at the handler.
- Levels: `Error` for things that broke, `Warn` for expected-but-noteworthy (failed login attempts, rate limits), `Info` for lifecycle events (server start, migration applied), `Debug` is opt-in via env.
- Never log secrets, password hashes, JWTs, or full request bodies. Email is fine. Passwords are not (yes, `auth_handler.go:47` currently logs the password; that's a security bug, not a style point).

---

## Cross-cutting

### Tests

Required for new business logic:

- Repository methods: Go `testing` with a real ArangoDB (testcontainers or a dedicated test DB). No mocked drivers.
- Service / auth / pricing logic: Go `testing`, table-driven where it fits.
- Frontend `.svelte.ts` rune modules: Vitest. Test the factory's public surface.

Not yet required (until the infra exists):

- HTTP handler integration tests
- Component tests for `.svelte` files
- E2E (Playwright)

A test for the happy path plus the obvious failure mode is enough to merge. Don't gold-plate coverage.

### Comments

Default to none. Names carry the meaning. Add a comment only when the _why_ is non-obvious: a hidden constraint, a workaround for a specific bug, an invariant that would surprise a reader six months from now.

Never write comments that restate the code or reference a ticket.

### Commit messages

Conventional Commits: `feat:`, `fix:`, `refactor:`, `docs:`, `chore:`, `test:`. Keep the subject under 72 chars. Body explains _why_, not _what_.

### Pull requests

One concern per PR. A bug fix doesn't bring a refactor along for the ride. If you find yourself touching a legacy violation while doing something else, fix it in a separate commit so the PR is reviewable.

---

## Legacy and migration

The current codebase predates these standards. The list below is the known drift. Don't add new instances of any of it. Migrate when you touch the file.

### Frontend

| Drift                                        | Where                                                                          | Action when you touch it                                                    |
| -------------------------------------------- | ------------------------------------------------------------------------------ | --------------------------------------------------------------------------- |
| PascalCase component filenames               | `lib/components/atoms/`, `molecules/`, `organisms/` (most files)               | Rename to kebab-case. Update imports.                                       |
| `svelte/store` based state                   | `lib/stores/cart.ts`, `lib/stores/ui.ts`, `lib/components/atoms/Navbar.svelte` | Port to a `.svelte.ts` rune module under `lib/state/`. Delete the old file. |
| `lib/stores/` directory                      | repo root of frontend lib                                                      | Move `user.svelte.ts` to `lib/state/`. Kill the `stores/` folder.           |
| `lib/api/Api.ts`, `Database.ts` (PascalCase) | `lib/api/`                                                                     | Rename to `api.ts`, `database.ts` when next touched.                        |
| `Record<string, any>`                        | `lib/components/organisms/UserAuthForm.svelte:17`                              | Replace with the actual props type.                                         |
| `any` outside vendored shadcn code           | `lib/components/atoms/Link.svelte:4`, others                                   | Replace with the right type.                                                |

### Backend

| Drift                                                     | Where                                                                                    | Action when you touch it                                                                                           |
| --------------------------------------------------------- | ---------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------ |
| Inline AQL in handlers                                    | `internal/handlers/auth_handler.go` (Login uses inline AQL via `database.QueryOne`)      | Move the query to a repository method.                                                                             |
| Missing repositories                                      | Releases, Artists, Tracks, Subscriptions, Subscribers (only `user_repository.go` exists) | Create the repo before adding/touching the handler.                                                                |
| Manual nil/empty-string validation after `ShouldBindJSON` | `auth_handler.go`, others                                                                | Replace with `binding:"required,..."` struct tags on the request DTOs.                                             |
| `fmt.Println` debug output                                | `auth_handler.go:135` (`fmt.Println(user)`)                                              | Delete.                                                                                                            |
| Hand-rolled string utilities                              | `auth_handler.go:196-207` (`contains`, `containsHelper`)                                 | Replace with `strings.Contains`.                                                                                   |
| Logging passwords                                         | `auth_handler.go:47` (logs `req.Password`)                                               | Remove the password field from the log. This is a security bug, fix on sight.                                      |
| Mixed response shapes                                     | most handlers                                                                            | Migrate to the envelope helpers when touching the handler.                                                         |
| Errors logged at every layer                              | mixed                                                                                    | Strip logger calls below the handler. Wrap with `%w` instead.                                                      |
| `unpublish` semantics in routes                           | `routes.go:103`, `release_handler.go`                                                    | Per ADR-0001 / `CONTEXT.md`, this should become `archive`. Coordinate with the rename when the handler is touched. |
