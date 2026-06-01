# Architecture

## Presentation (clients)

Web and mobile clients for pet owners, reception, and veterinarians (future or external repos).

### Components

- Customer portal (appointments, pets, medical read-only)
- Staff workstation (schedule, clinical documentation)
- Admin console (users, payments, notifications)

### Responsibilities

- Store JWT access token; refresh before expiry
- Call REST JSON under /api/v2 and /health

### Technologies

- HTTPS
- Bearer JWT
- OpenAPI client (optional)

## API gateway & edge

TLS termination and routing to the Go API (ALB target group or local Docker port 8000).

### Components

- AWS ALB + ACM (planned)
- CORS middleware (gin-contrib/cors)

### Responsibilities

- Terminate TLS and forward to container port 8000
- Health checks on GET /health

### Technologies

- Application Load Balancer
- Let's Encrypt or ACM

## Application layer (hexagonal)

Gin HTTP adapters delegating to domain services in internal/core; sqlc repositories in internal/infrastructure/persistence.

### Components

- main.go — server lifecycle, middleware stack
- internal/infrastructure/http — handlers, router, DTOs, mappers
- internal/core/* — auth, users, customers, employees, pets, appointments, medical, payments, notifications, addresses
- internal/config — env loading, Redis, Twilio, bootstrap

### Responsibilities

- Enforce RBAC via JWT claims
- Map domain errors to APIResponse envelope
- Run migrations on container start

### Technologies

- Gin
- go-playground/validator
- uber/zap logging

## Data & cache

PostgreSQL for authoritative state; Redis for token/session storage.

### Components

- Amazon RDS PostgreSQL
- Upstash Redis
- sqlc generated Queries

### Responsibilities

- ACID transactions for appointments and payments
- Persist and revoke refresh tokens

### Technologies

- pgx/v5
- go-redis/v9
- golang-migrate

## Integrations

Outbound email and SMS; no message queue in current codebase.

### Components

- SMTP (activation, reset, notifications)
- Twilio SMS

### Responsibilities

- Deliver transactional messages when senders are wired

### Technologies

- net/smtp patterns via notification service
- twilio-go SDK

## Design patterns

| Pattern | Category | Description |
| --- | --- | --- |
| ⬡ Hexagonal (ports & adapters) | Structural | Domain packages define interfaces; infrastructure implements repositories and HTTP handlers as driving/driven adapters. |
| ↔️ CQRS-lite | Behavioral | Separate CommandService and QueryService types in appointments, payments, and users reduce coupling between reads and writes. |
| 🔍 Specification pattern | Behavioral | Appointment and pet list endpoints build specification objects for composable SQL filters. |
| ⚙️ Domain service | Domain | AppointmentDomainService encapsulates overlap and capacity rules that need repository access. |
| 🗺️ DTO + mapper | Structural | HTTP DTOs in handlers/dtos map to domain commands via dedicated mapper structs. |
| 🔗 Middleware pipeline | Behavioral | Recovery → CORS → AuditLog → RateLimiter → route groups with AuthMiddleware. |

## Scalability strategies

- **Stateless API tasks** — Horizontally scale ECS tasks behind ALB; JWT validation requires no sticky sessions.
- **Managed PostgreSQL** — RDS handles backups and storage growth; connection pool tuned via DATABASE_MAX_* env vars.
- **External Redis** — Upstash shared cache for tokens across multiple API replicas.
- **Read-heavy clinical queries** — sqlc queries with pagination; add read replicas on RDS when reporting load grows.

## Security strategies

- **Bearer JWT on protected routes** — AuthMiddleware parses HS256 access tokens; role checked per route group.
- **Customer data scoping** — Me/* routes resolve customer ID from JWT user ID before returning pets or medical data.
- **Password hashing** — internal/shared/password encoder for stored credentials.
- **Audit logging** — Request ID, method, path, status, and duration logged for traceability.
- **Rate limiting** — Global IP-based limiter when RATE_LIMIT_ENABLED=true.
- **RDS TLS** — database_url.go adds sslmode=require for Amazon RDS hostnames.

## Cache strategies

| Name | TTL | Coverage | Description |
| --- | --- | --- | --- |
| JWT / session tokens | JWT_EXPIRATION_TIME / REFRESH_TOKEN_EXPIRY from env | Login, refresh, logout revocation | TokenService stores access and refresh tokens in Redis when configured |
| Rate limit counters | RATE_LIMIT_WINDOW (default 0 = per-second bucket) | All HTTP routes when enabled | In-memory sliding window per client IP in default middleware |
| PostgreSQL pool | 1h max lifetime default | All sqlc repository calls | pgx pool reuse with CONN_MAX_LIFETIME and idle timeouts |

## Architecture highlights

### 📦 Unified APIResponse envelope

success, message, data, error, meta, timestamp, and request_id on every JSON response.

### 📖 Swagger-first handlers

swag annotations on handlers; docs package generated into docs/docs.go.

### 💚 Health endpoint

GET /health returns status, timestamp, service name, and version for probes.

### 🧩 Domain modules per bounded context

Separate packages for auth, pets, medical, payments—clear boundaries for team ownership.

## Architecture diagram

### Legend

| Type | Label |
| --- | --- |
| client | Client |
| gateway | Gateway |
| service | API service |
| database | Database |
| queue | Cache |
| monitoring | External |

### Nodes

| ID | Label | Type | Status |
| --- | --- | --- | --- |
| clients | Web / mobile clients | client | healthy |
| alb | ALB (TLS) | gateway | healthy |
| api | Vet Clinic API (Go/Gin) | service | healthy |
| rds | RDS PostgreSQL | database | healthy |
| redis | Upstash Redis | queue | healthy |
| smtp | SMTP / Twilio | monitoring | healthy |
| ecr | ECR (images) | monitoring | healthy |

### Connections

| From | To | Label | Protocol |
| --- | --- | --- | --- |
| clients | alb | HTTPS | TLS 1.2+ |
| alb | api | Forward | HTTP |
| api | rds | SQL | PostgreSQL |
| api | redis | Tokens | Redis TLS |
| api | smtp | Notify | SMTP/SMS |
| ecr | api | Deploy | OCI image |

### Mermaid overview

```mermaid
flowchart LR
    clients([Web / mobile clients])
    alb{ALB (TLS)}
    api[Vet Clinic API (Go/Gin)]
    rds[(RDS PostgreSQL)]
    redis[/Upstash Redis/]
    smtp>SMTP / Twilio]
    ecr>ECR (images)]
    clients -->|HTTPS| alb
    alb -->|Forward| api
    api -->|SQL| rds
    api -->|Tokens| redis
    api -->|Notify| smtp
    ecr -->|Deploy| api
```

## Data flow

### Request flow

1. **Client request** — Client sends HTTP request; Bearer token on protected routes except public auth and /health.
2. **Middleware chain** — Recovery, CORS, audit log, optional rate limit, then JWT Authenticate and RequireAnyRole on route groups.
3. **Handler → domain service** — Gin handler validates DTO, maps to command/query, invokes core service.
4. **sqlc repository** — Persistence adapter runs parameterized SQL via pgx against RDS.
5. **APIResponse JSON** — shared/http helpers return success, pagination meta, or typed ApplicationError mapping.

### Event flow

1. **Domain action** — e.g. appointment confirmed or user registered triggers notification service call.
2. **Notification record** — Row inserted in notifications table via sqlc repository.
3. **Channel dispatch (when wired)** — EmailSender or SMSender implementation delivers via SMTP or Twilio.
4. **Client poll** — User reads /api/v2/me/notifications for delivery status in MVP.

## Additional notes

# Architecture

> **Veterinary domain:** Customer-facing routes always scope data to the authenticated pet owner; staff routes require elevated JWT roles once role names are aligned with router guards.

> **Deploy alignment:** Diagram reflects **live** RDS + Upstash and **planned** ALB/ECS/ECR path documented in ProjectInfrastructure.md.

> **Technical debt:** No async worker—notifications are synchronous DB writes until email/SMS adapters are injected. Medical handler nil blocks clinical HTTP despite route definitions. Consider Redis-backed rate limiter for multi-instance ECS.

