# AGENTS.md

High-signal notes for working in this repo. Skip if the answer is obvious from filenames or framework defaults.

## Start the app

Docker is the only officially supported way to run the system. Three containers: `api` (Go), `client` (React/Vite), `database` (Postgres 17) plus a `migrate` one-shot.

```sh
cp .env.example .env   # set DB_PASSWORD, ADMIN_PASSWORD, JWT_SECRET, WS_ORIGIN
# Dev (exposes client:3000, api:8080):
docker compose -f docker-compose.yml -f docker-compose.dev.yml up
# Prod (no exposed ports — bring your own reverse proxy):
docker compose up
```

`WS_ORIGIN` is matched against the WebSocket upgrade Origin header. Use the URL the browser will hit (e.g. `http://localhost:3000` for dev).

## Repo layout

- `api/` — Go 1.26, module `classdir/api`. HTTP + WebSocket server. Entry: `api/main.go`. Packages: `auth/`, `db/`, `hub/`, `presentation/`, `shared/{cfg,response,sanitize,validate}/`.
- `client/` — React 19 + Vite 8 + TypeScript, Mantine v9. Entry: `client/src/main.tsx`. Routes defined in `client/src/shared/cfg/routes.ts`. Mirror API contract in `client/src/presentation/types.ts` and `client/src/presentation/cfg.ts`.
- `database/` — Go 1.26, module `classdir/db-migrate`. Custom migration runner; embeds `database/migrations/*.sql` and tracks versions in a `schema_migrations` table. Run automatically by the `migrate` service before the API starts.

## Verify a change

```sh
# Backend
cd api && go build ./... && go test ./...

# Frontend
cd client && pnpm install   # only if deps changed
pnpm lint                   # oxlint
pnpm fmt:check              # oxfmt --check
pnpm build                  # tsc -b && vite build

# Migrations
ls database/migrations      # files are applied in lexical order; versions are NOT timestamped
```

The client uses `oxlint` and `oxfmt` (not ESLint/Prettier). Configs: `client/.oxlintrc.json`, `client/.oxfmtrc.json`. Build runs `tsc -b` over the `tsconfig.app.json` + `tsconfig.node.json` project references.

## Contracts and cross-container changes

When a task touches more than one container, update the shared contract (WebSocket commands/events and REST endpoints) first, then implement each side in isolation. The single source of truth is:

- `docs/todo/COMMANDS.md` — WebSocket command/event spec (which are `[x]` done, which `[ ]` pending).
- `docs/todo/ENDPOINTS.md` — REST endpoint spec, request/response shapes.
- API constant names: `api/internal/hub/hub.go` (commands + events) and `api/internal/hub/annotation_types.go` (annotation types).
- Client mirror: `client/src/presentation/cfg.ts` (commands, events, annotation types) and `client/src/presentation/types.ts` (zod schemas).

Treat each container as isolated — do not copy patterns from one into the other without reason. Don't modify the specs without an explicit instruction.

## Key quirks an agent will trip on

- **Auth is single-password.** Login (`POST /api/v1/auth/login` with `{password}`) issues a JWT in an `HttpOnly`, `Secure`, `SameSite=Strict` cookie named `token` (see `api/internal/shared/cfg/cfg.go`). The same cookie authorizes the WebSocket upgrade at `GET /ws/v1`. No user accounts.
- **Presentation IDs must be UUID v7.** Validated in `api/internal/shared/validate`. Hand-crafted IDs in tests/clients must use v7.
- **Annotations are not persisted.** They live in-memory per room, per slide. Server replays them on join via an `annotations_batch` event. Lost on restart.
- **Slide splitting.** The server splits `presentations.content` into slides using the regex `^---+\s*$` (line-based). The client just renders whatever HTML the server returns.
- **WS_ORIGIN.** Must match the browser's Origin header exactly. Mismatch silently rejects the WebSocket.
- **VITE_API_URL.** Build arg in `client/Dockerfile`; defaults to `""` in compose so the client uses relative paths and the nginx config proxies `/api/` and `/ws/` to `api:8080`. Local `vite dev` proxies the same paths to `localhost:8080` (see `client/vite.config.ts`).
- **JWT in dev over HTTP.** The cookie is `Secure`, so `WS_ORIGIN=http://...` only works behind a proxy that terminates TLS, or in local dev over plain HTTP where the browser will still reject `Secure` cookies on `http://`. For local dev, accept that login may not persist a cookie and use the API directly for testing.
- **Rate limiting** is per WebSocket connection via `golang.org/x/time/rate`; defaults are defined in the hub package. Tests use a `mockRateLimitProvider` to control behavior.
- **Go modules are two.** `api/go.mod` and `database/go.mod` are independent. Don't `go work` them — `docker compose` builds them separately.

## Domain rules (don't break)

From `ARCHITECTURE.md`:

- Slide content and student registration are locked once a presentation starts. Teachers can still add/remove pre-registered students from the spin wheel.
- Annotations are slide-scoped. Changing slides clears the canvas.
- `docs/troubleshooting.md` is for human users — do not read it as ground truth.

## Where to look first

- `ARCHITECTURE.md` — system overview, mermaid diagrams, response/error envelope shapes.
- `api/internal/hub/` — WebSocket command dispatch, room lifecycle, broadcast logic. Tests in `hub_test.go` are the clearest spec of expected behavior.
- `api/internal/presentation/handler.go` + `store.go` — REST handler and pgx-backed store. `Store` is an interface, so handlers are testable without a DB (see `presentation_test.go`).
- `client/src/presentation/` — hooks and components for present/control/configure views.
