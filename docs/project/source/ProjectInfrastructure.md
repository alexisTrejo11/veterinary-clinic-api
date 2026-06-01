---
metrics:
  - label: "Container port"
    value: "8000"
    icon: "server"
    description: "Gin HTTP server binds SERVER_HOST:SERVER_PORT (default 0.0.0.0:8000)"
  - label: "Local Postgres port"
    value: "5431"
    icon: "database"
    description: "Host mapping POSTGRES_PUBLISHED_PORT in compose.local.yml (container 5432)"
  - label: "Connection pool"
    value: "25 open / 5 idle"
    icon: "layers"
    description: "DATABASE_MAX_OPEN_CONNS and DATABASE_MAX_IDLE_CONNS from .env.example"
  - label: "Migrations"
    value: "9 up files"
    icon: "migrate"
    description: "golang-migrate runs from entrypoint unless SKIP_MIGRATIONS=true"

cloudServices:
  - name: "Amazon RDS (PostgreSQL)"
    purpose: "Primary transactional store for users, customers, employees, pets, appointments, medical sessions, payments, notifications"
    icon: "aws-rds"
    cost: "~$25–80/mo (db.t4g.micro placeholder; scale with storage and Multi-AZ)"
  - name: "Upstash Redis"
    purpose: "JWT/session token storage, revocation lists, and future rate-limit backend; REDIS_URL supports rediss:// TLS"
    icon: "redis"
    cost: "Free tier or ~$10–20/mo for production serverless"
  - name: "AWS ECS Fargate (planned)"
    purpose: "Run docker/Dockerfile image without managing EC2 patches; target orchestrator for production deploy"
    icon: "aws-ecs"
    cost: "~$30–60/mo (0.25 vCPU / 0.5 GB task placeholder)"
  - name: "AWS Application Load Balancer (planned)"
    purpose: "TLS termination, health checks on GET /health, route to ECS service"
    icon: "aws-alb"
    cost: "~$18/mo + LCU usage"
  - name: "Amazon ECR (planned)"
    purpose: "Private container registry for CI-built API images"
    icon: "aws-ecr"
    cost: "~$1/mo storage + pull bandwidth"
  - name: "Twilio"
    purpose: "SMS notifications (appointment reminders, activation codes when wired)"
    icon: "twilio"
    cost: "Per SMS segment"
  - name: "SMTP (Gmail / Amazon SES)"
    purpose: "Account activation, password reset, notification emails"
    icon: "mail"
    cost: "SES low volume or Gmail app password (dev only)"

deploymentLayers:
  - name: "Clients"
    color: "#4F46E5"
    components:
      - name: "Pet owner portal (future)"
        icon: "smartphone"
        description: "JWT customer role—book appointments, view pets and read-only medical summaries"
      - name: "Clinic staff web app (future)"
        icon: "layout"
        description: "Receptionist and veterinarian workflows against /api/v2"
      - name: "Admin / integrators"
        icon: "terminal"
        description: "User management, reporting, and operational scripts"

  - name: "Edge & compute (AWS — planned)"
    color: "#059669"
    components:
      - name: "ALB + ACM certificate"
        icon: "shield"
        description: "HTTPS to target group; health check path /health"
      - name: "ECS service (clinical_vet_api)"
        icon: "docker"
        description: "Task from docker/Dockerfile; env from SSM/Secrets Manager"
      - name: "Optional: EC2 + Docker Compose"
        icon: "server"
        description: "Alternative to ECS—compose.dev.yml with cloud DATABASE_URL and REDIS_URL"

  - name: "Data plane (live)"
    color: "#DC2626"
    components:
      - name: "RDS PostgreSQL"
        icon: "database"
        description: "DATABASE_URL with sslmode=require auto-added for *.rds.amazonaws.com"
      - name: "Upstash Redis"
        icon: "redis"
        description: "REDIS_URL rediss:// for TLS; used by token service"
      - name: "Local dev volumes"
        icon: "folder"
        description: "pgdata and redis_data in compose.local.yml only"

  - name: "Observability & messaging"
    color: "#D97706"
    components:
      - name: "Zap structured logs"
        icon: "file-text"
        description: "Audit middleware + lumberjack rotation under logs/"
      - name: "Twilio + SMTP"
        icon: "message"
        description: "Notification service accepts nil senders until configured"
      - name: "Sentry (optional placeholder)"
        icon: "activity"
        description: "Not wired in code yet—add DSN via env when needed"

dockerFiles:
  - service: "compose.local.yml (full stack)"
    description: "Postgres 12, Redis 7, and API on vet_network. Started via ./docker/up-local.sh which rewrites DATABASE_URL/REDIS_URL hostnames for in-network DNS."
    content: |
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

  - service: "compose.dev.yml (API only)"
    description: "Single API container pointing at cloud RDS and Upstash; SKIP_MIGRATIONS and SKIP_DB_WAIT default true for faster dev iteration."
    content: |
      name: vet-clinic-dev
      services:
        api:
          build: { context: .., dockerfile: docker/Dockerfile }
          env_file: ../.env
          ports: ["8000:8000"]
          environment:
            SKIP_MIGRATIONS: ${SKIP_MIGRATIONS:-true}
            SKIP_DB_WAIT: ${SKIP_DB_WAIT:-true}

  - service: "Dockerfile (multi-stage)"
    description: "Builder: Go 1.24 Alpine, compile main, install migrate CLI. Runtime: Alpine 3.18, non-root /root, EXPOSE 8000, entrypoint runs migrations then ./main."
    content: |
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

  - service: "entrypoint.sh"
    description: "Loads .env, parses DATABASE_URL via parse-database-url.sh, waits for Postgres (unless SKIP_DB_WAIT), runs migrate up, execs the API binary."
    content: |
      # Requires DATABASE_URL + DATABASE_USER + DATABASE_PASSWORD
      # migrate -path ./db/migrations -database "$DATABASE_URL" up
      # exec "$@"  → ./main
---

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
