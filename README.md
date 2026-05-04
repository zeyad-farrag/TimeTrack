# Team App

Standalone Go and Next.js team capacity app for Multica-adjacent time tracking.

## Local Development

Copy `.env.example` to `.env`, then run the stack:

```sh
docker compose up
```

The API listens on `:8080`, the frontend listens on `:3000`, and Postgres uses the `team_app` database. The API exposes `GET /healthz`.

Run backend tests:

```sh
go test ./...
```

Run the frontend locally:

```sh
cd frontend
pnpm install
pnpm dev
```

Run frontend checks:

```sh
cd frontend
pnpm lint
pnpm build
```

## Nginx

`nginx/team.conf` routes `team.multica.uittai.com`: `/api/` and `/gates/` proxy to the Go API, and all other paths proxy to the Next.js frontend.

## Companion Multica Fork PRs

Stories 1.2, 1.3, and 1.4 land the Multica-side changes as M-PR#1, M-PR#2, and M-PR#3. For Multica migrations, inspect the multica repo at PR-open time with `git log --oneline -- server/migrations/ | head` against `origin/main`; never assume migration prefix `057`.
