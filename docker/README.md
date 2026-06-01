# Docker setup

Container files live in this directory. The API image is built from the **repository root**.

## Profiles

| Profile | How to start | What runs |
|---------|----------------|-----------|
| **Local** | `docker compose --env-file .env -f docker/compose.local.yml up --build` | Postgres, Redis, API |
| **Dev** | `docker compose --env-file .env -f docker/compose.dev.yml up --build` | API only (cloud DB/Redis via `.env`) |

Shortcut for local: `./docker/up-local.sh --build`

## Prerequisites

```bash
cp .env.example .env
```

Run commands from the **repository root**.

---

## Local profile (full stack)

`compose.local.yml` **overrides** `DATABASE_URL` and `REDIS_URL` inside the API container so it always connects to `postgres12` and `redis` on the Docker network. Your `.env` can still point at `localhost` or AWS for `go run` on the host — the local profile does not change those files.

### Required in `.env` for local stack

```bash
DATABASE_USER=postgres
DATABASE_PASSWORD=your_password
DATABASE_NAME=vet_database

REDIS_PASSWORD=your_redis_password
REDIS_DB=0
```

Postgres and Redis containers use the same `DATABASE_*` / `REDIS_PASSWORD` values.

### Start

```bash
docker compose --env-file .env -f docker/compose.local.yml up --build
# or
./docker/up-local.sh --build
```

### What the API container receives (automatic)

| Variable | Value (example) |
|----------|-----------------|
| `DATABASE_URL` | `postgresql://postgres:pass@postgres12:5432/vet_database?sslmode=disable` |
| `REDIS_URL` | `redis://:pass@redis:6379/0` |

### Host access (from your machine)

| Service | URL |
|---------|-----|
| API | http://localhost:8000 |
| Swagger | http://localhost:8000/swagger/index.html |
| Postgres | `localhost:5431` (`POSTGRES_PUBLISHED_PORT`) |
| Redis | `localhost:6379` (`REDIS_PUBLISHED_PORT`) |

Migrations run on startup unless `SKIP_MIGRATIONS=true`.

---

## Dev profile (API only)

No overrides — the API uses `DATABASE_URL`, `DATABASE_USER`, `DATABASE_PASSWORD`, and `REDIS_URL` from `.env` as-is (RDS, Upstash, etc.).

```bash
DATABASE_URL=jdbc:postgresql://your-rds-host:5432/vet_database
DATABASE_USER=postgres
DATABASE_PASSWORD=your_password
REDIS_URL=rediss://:token@your-redis.upstash.io:6379
```

```bash
docker compose --env-file .env -f docker/compose.dev.yml up --build
```

Defaults: `SKIP_MIGRATIONS=true`, `SKIP_DB_WAIT=true`.

---

## Migrations without Docker

```bash
export $(grep -v '^#' .env | xargs)  # or use direnv
go run ./cmd/migrate
```

---

## File layout

```
docker/
├── Dockerfile
├── compose.local.yml    # local profile + in-network URL overrides
├── compose.dev.yml      # API only, no overrides
├── up-local.sh          # thin wrapper around compose.local.yml
└── README.md
```

---

## Troubleshooting

| Issue | Fix |
|-------|-----|
| `DATABASE_URL is required` | Set `DATABASE_USER`, `DATABASE_PASSWORD`, `DATABASE_NAME` in `.env` |
| Redis healthcheck fails | Set `REDIS_PASSWORD` in `.env` (same value local API and redis container use) |
| API still hits RDS in local profile | Use `compose.local.yml`, not `compose.dev.yml` |
| JWT too short | Longer `JWT_SECRET` (32+ chars) |
