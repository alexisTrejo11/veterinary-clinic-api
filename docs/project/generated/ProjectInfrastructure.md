# Infrastructure

## Metrics

| Label | Value | Description |
| --- | --- | --- |
| Container port | 8000 | Gin HTTP server binds SERVER_HOST:SERVER_PORT (default 0.0.0.0:8000) |
| Local Postgres port | 5431 | Host mapping POSTGRES_PUBLISHED_PORT in compose.local.yml (container 5432) |
| Connection pool | 25 open / 5 idle | DATABASE_MAX_OPEN_CONNS and DATABASE_MAX_IDLE_CONNS from .env.example |
| Migrations | 9 up files | golang-migrate runs from entrypoint unless SKIP_MIGRATIONS=true |

## Cloud services

| Service | Purpose | Est. cost |
| --- | --- | --- |
| Amazon RDS (PostgreSQL) | Primary transactional store for users, customers, employees, pets, appointments, medical sessions, payments, notifications | ~$25–80/mo (db.t4g.micro placeholder; scale with storage and Multi-AZ) |
| Upstash Redis | JWT/session token storage, revocation lists, and future rate-limit backend; REDIS_URL supports rediss:// TLS | Free tier or ~$10–20/mo for production serverless |
| AWS ECS Fargate (planned) | Run docker/Dockerfile image without managing EC2 patches; target orchestrator for production deploy | ~$30–60/mo (0.25 vCPU / 0.5 GB task placeholder) |
| AWS Application Load Balancer (planned) | TLS termination, health checks on GET /health, route to ECS service | ~$18/mo + LCU usage |
| Amazon ECR (planned) | Private container registry for CI-built API images | ~$1/mo storage + pull bandwidth |
| Twilio | SMS notifications (appointment reminders, activation codes when wired) | Per SMS segment |
| SMTP (Gmail / Amazon SES) | Account activation, password reset, notification emails | SES low volume or Gmail app password (dev only) |

## Deployment layers

### Clients

- **Pet owner portal (future)** — JWT customer role—book appointments, view pets and read-only medical summaries
- **Clinic staff web app (future)** — Receptionist and veterinarian workflows against /api/v2
- **Admin / integrators** — User management, reporting, and operational scripts

### Edge & compute (AWS — planned)

- **ALB + ACM certificate** — HTTPS to target group; health check path /health
- **ECS service (clinical_vet_api)** — Task from docker/Dockerfile; env from SSM/Secrets Manager
- **Optional: EC2 + Docker Compose** — Alternative to ECS—compose.dev.yml with cloud DATABASE_URL and REDIS_URL

### Data plane (live)

- **RDS PostgreSQL** — DATABASE_URL with sslmode=require auto-added for *.rds.amazonaws.com
- **Upstash Redis** — REDIS_URL rediss:// for TLS; used by token service
- **Local dev volumes** — pgdata and redis_data in compose.local.yml only

### Observability & messaging

- **Zap structured logs** — Audit middleware + lumberjack rotation under logs/
- **Twilio + SMTP** — Notification service accepts nil senders until configured
- **Sentry (optional placeholder)** — Not wired in code yet—add DSN via env when needed

## Docker configuration

### compose.local.yml (full stack)

Postgres 12, Redis 7, and API on vet_network. Started via ./docker/up-local.sh which rewrites DATABASE_URL/REDIS_URL hostnames for in-network DNS.

```yaml
name: vet-clinic-local
services:
  postgres12:
    image: postgres:12-alpine
    ports: ["5431:5432"]
  redis:
    image: redis:7-alpine
    command: redis-server --requirepass ${REDIS_PASSWORD}
  api:
    build: { context: .., dockerfile: docker/Dockerfile }
    ports: ["8000:8000"]
    environment:
      DATABASE_URL: ${DOCKER_DATABASE_URL}
      REDIS_URL: ${DOCKER_REDIS_URL}
    depends_on: [postgres12, redis]
```

### compose.dev.yml (API only)

Single API container pointing at cloud RDS and Upstash; SKIP_MIGRATIONS and SKIP_DB_WAIT default true for faster dev iteration.

```yaml
name: vet-clinic-dev
services:
  api:
    build: { context: .., dockerfile: docker/Dockerfile }
    env_file: ../.env
    ports: ["8000:8000"]
    environment:
      SKIP_MIGRATIONS: ${SKIP_MIGRATIONS:-true}
      SKIP_DB_WAIT: ${SKIP_DB_WAIT:-true}
```

### Dockerfile (multi-stage)

Builder: Go 1.24 Alpine, compile main, install migrate CLI. Runtime: Alpine 3.18, non-root /root, EXPOSE 8000, entrypoint runs migrations then ./main.

```yaml
FROM golang:1.24.0-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download && CGO_ENABLED=0 go build -o main .
FROM alpine:3.18
COPY --from=builder /app/main .
COPY db/migrations ./db/migrations/
COPY scripts/entrypoint.sh /entrypoint.sh
EXPOSE 8000
ENTRYPOINT ["/entrypoint.sh"]
CMD ["./main"]
```

### entrypoint.sh

Loads .env, parses DATABASE_URL via parse-database-url.sh, waits for Postgres (unless SKIP_DB_WAIT), runs migrate up, execs the API binary.

```yaml
# Requires DATABASE_URL + DATABASE_USER + DATABASE_PASSWORD
# migrate -path ./db/migrations -database "$DATABASE_URL" up
# exec "$@"  → ./main
```

## Additional notes

# Infrastructure

> **Deploy story (target):** Build and push image to ECR → ECS Fargate service behind ALB → environment variables from AWS Secrets Manager (`DATABASE_URL`, `REDIS_URL`, `JWT_SECRET`, Twilio, SMTP). RDS and Upstash are **already** the production data tier; this repo documents the API container path only.

> **Local parity:** `./docker/up-local.sh --build` mirrors production dependencies (Postgres + Redis + API) on localhost:8000 with migrations on startup.

> **AWS checklist (when you deploy):**
> - Security group: allow 443 on ALB; restrict RDS/Redis to VPC or Upstash IP allowlist.
> - Set `ENVIRONMENT=production`, `ENABLE_SWAGGER=false`, strong `JWT_SECRET` (32+ chars).
> - Run migrations once per release (`SKIP_MIGRATIONS=false` on first deploy, then true if you migrate in CI).
> - Map `CORS_ALLOW_ORIGINS` to real frontend domains—not `*` in production.

> **Dangerous:** `000007_insert_demo_data.up.sql` seeds demo rows—disable or skip in production databases. Never expose Redis without password/TLS on the public internet.

> **Missing in repo:** Terraform/CloudFormation, GitHub Actions deploy workflow, and ECS task definition—these are the next artifacts to add for one-click AWS rollout.

