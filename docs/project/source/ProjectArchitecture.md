---
layers:
  - name: "Presentation (clients)"
    description: "Web and mobile clients for pet owners, reception, and veterinarians (future or external repos)."
    color: "#6366F1"
    expanded: true
    components:
      - "Customer portal (appointments, pets, medical read-only)"
      - "Staff workstation (schedule, clinical documentation)"
      - "Admin console (users, payments, notifications)"
    responsibilities:
      - "Store JWT access token; refresh before expiry"
      - "Call REST JSON under /api/v2 and /health"
    technologies:
      - "HTTPS"
      - "Bearer JWT"
      - "OpenAPI client (optional)"

  - name: "API gateway & edge"
    description: "TLS termination and routing to the Go API (ALB target group or local Docker port 8000)."
    color: "#10B981"
    expanded: false
    components:
      - "AWS ALB + ACM (planned)"
      - "CORS middleware (gin-contrib/cors)"
    responsibilities:
      - "Terminate TLS and forward to container port 8000"
      - "Health checks on GET /health"
    technologies:
      - "Application Load Balancer"
      - "Let's Encrypt or ACM"

  - name: "Application layer (hexagonal)"
    description: "Gin HTTP adapters delegating to domain services in internal/core; sqlc repositories in internal/infrastructure/persistence."
    color: "#F59E0B"
    expanded: true
    components:
      - "main.go — server lifecycle, middleware stack"
      - "internal/infrastructure/http — handlers, router, DTOs, mappers"
      - "internal/core/* — auth, users, customers, employees, pets, appointments, medical, payments, notifications, addresses"
      - "internal/config — env loading, Redis, Twilio, bootstrap"
    responsibilities:
      - "Enforce RBAC via JWT claims"
      - "Map domain errors to APIResponse envelope"
      - "Run migrations on container start"
    technologies:
      - "Gin"
      - "go-playground/validator"
      - "uber/zap logging"

  - name: "Data & cache"
    description: "PostgreSQL for authoritative state; Redis for token/session storage."
    color: "#EF4444"
    expanded: true
    components:
      - "Amazon RDS PostgreSQL"
      - "Upstash Redis"
      - "sqlc generated Queries"
    responsibilities:
      - "ACID transactions for appointments and payments"
      - "Persist and revoke refresh tokens"
    technologies:
      - "pgx/v5"
      - "go-redis/v9"
      - "golang-migrate"

  - name: "Integrations"
    description: "Outbound email and SMS; no message queue in current codebase."
    color: "#8B5CF6"
    expanded: false
    components:
      - "SMTP (activation, reset, notifications)"
      - "Twilio SMS"
    responsibilities:
      - "Deliver transactional messages when senders are wired"
    technologies:
      - "net/smtp patterns via notification service"
      - "twilio-go SDK"

designPatterns:
  - title: "Hexagonal (ports & adapters)"
    emoji: "⬡"
    description: "Domain packages define interfaces; infrastructure implements repositories and HTTP handlers as driving/driven adapters."
    category: "Structural"
    badge: "Core"
  - title: "CQRS-lite"
    emoji: "↔️"
    description: "Separate CommandService and QueryService types in appointments, payments, and users reduce coupling between reads and writes."
    category: "Behavioral"
    badge: "Domain"
  - title: "Specification pattern"
    emoji: "🔍"
    description: "Appointment and pet list endpoints build specification objects for composable SQL filters."
    category: "Behavioral"
    badge: "Query"
  - title: "Domain service"
    emoji: "⚙️"
    description: "AppointmentDomainService encapsulates overlap and capacity rules that need repository access."
    category: "Domain"
    badge: "Appointments"
  - title: "DTO + mapper"
    emoji: "🗺️"
    description: "HTTP DTOs in handlers/dtos map to domain commands via dedicated mapper structs."
    category: "Structural"
    badge: "HTTP"
  - title: "Middleware pipeline"
    emoji: "🔗"
    description: "Recovery → CORS → AuditLog → RateLimiter → route groups with AuthMiddleware."
    category: "Behavioral"
    badge: "Gin"

scalabilityStrategies:
  - title: "Stateless API tasks"
    description: "Horizontally scale ECS tasks behind ALB; JWT validation requires no sticky sessions."
  - title: "Managed PostgreSQL"
    description: "RDS handles backups and storage growth; connection pool tuned via DATABASE_MAX_* env vars."
  - title: "External Redis"
    description: "Upstash shared cache for tokens across multiple API replicas."
  - title: "Read-heavy clinical queries"
    description: "sqlc queries with pagination; add read replicas on RDS when reporting load grows."

securityStrategies:
  - title: "Bearer JWT on protected routes"
    description: "AuthMiddleware parses HS256 access tokens; role checked per route group."
  - title: "Customer data scoping"
    description: "Me/* routes resolve customer ID from JWT user ID before returning pets or medical data."
  - title: "Password hashing"
    description: "internal/shared/password encoder for stored credentials."
  - title: "Audit logging"
    description: "Request ID, method, path, status, and duration logged for traceability."
  - title: "Rate limiting"
    description: "Global IP-based limiter when RATE_LIMIT_ENABLED=true."
  - title: "RDS TLS"
    description: "database_url.go adds sslmode=require for Amazon RDS hostnames."

cacheStrategies:
  - name: "JWT / session tokens"
    description: "TokenService stores access and refresh tokens in Redis when configured"
    ttl: "JWT_EXPIRATION_TIME / REFRESH_TOKEN_EXPIRY from env"
    coverage: "Login, refresh, logout revocation"
  - name: "Rate limit counters"
    description: "In-memory sliding window per client IP in default middleware"
    ttl: "RATE_LIMIT_WINDOW (default 0 = per-second bucket)"
    coverage: "All HTTP routes when enabled"
  - name: "PostgreSQL pool"
    description: "pgx pool reuse with CONN_MAX_LIFETIME and idle timeouts"
    ttl: "1h max lifetime default"
    coverage: "All sqlc repository calls"

architectureFeatures:
  - title: "Unified APIResponse envelope"
    emoji: "📦"
    description: "success, message, data, error, meta, timestamp, and request_id on every JSON response."
  - title: "Swagger-first handlers"
    emoji: "📖"
    description: "swag annotations on handlers; docs package generated into docs/docs.go."
  - title: "Health endpoint"
    emoji: "💚"
    description: "GET /health returns status, timestamp, service name, and version for probes."
  - title: "Domain modules per bounded context"
    emoji: "🧩"
    description: "Separate packages for auth, pets, medical, payments—clear boundaries for team ownership."

architectureDiagram:
  legendItems:
    - type: "client"
      label: "Client"
      color: "#6366F1"
      icon: "monitor"
    - type: "gateway"
      label: "Gateway"
      color: "#10B981"
      icon: "shield"
    - type: "service"
      label: "API service"
      color: "#F59E0B"
      icon: "server"
    - type: "database"
      label: "Database"
      color: "#EF4444"
      icon: "database"
    - type: "queue"
      label: "Cache"
      color: "#8B5CF6"
      icon: "layers"
    - type: "monitoring"
      label: "External"
      color: "#64748B"
      icon: "activity"

  nodes:
    - id: "clients"
      label: "Web / mobile clients"
      type: "client"
      x: 80
      y: 120
      connections: ["alb"]
      status: "healthy"
      traffic: 100
    - id: "alb"
      label: "ALB (TLS)"
      type: "gateway"
      x: 280
      y: 120
      connections: ["api"]
      status: "healthy"
      traffic: 100
    - id: "api"
      label: "Vet Clinic API (Go/Gin)"
      type: "service"
      x: 480
      y: 120
      connections: ["rds", "redis", "smtp"]
      status: "healthy"
      traffic: 90
    - id: "rds"
      label: "RDS PostgreSQL"
      type: "database"
      x: 700
      y: 60
      connections: []
      status: "healthy"
      traffic: 75
    - id: "redis"
      label: "Upstash Redis"
      type: "queue"
      x: 700
      y: 180
      connections: []
      status: "healthy"
      traffic: 50
    - id: "smtp"
      label: "SMTP / Twilio"
      type: "monitoring"
      x: 480
      y: 260
      connections: []
      status: "healthy"
      traffic: 20
    - id: "ecr"
      label: "ECR (images)"
      type: "monitoring"
      x: 80
      y: 260
      connections: ["api"]
      status: "healthy"
      traffic: 5

  connections:
    - id: "c1"
      from: "clients"
      to: "alb"
      label: "HTTPS"
      protocol: "TLS 1.2+"
      isActive: true
    - id: "c2"
      from: "alb"
      to: "api"
      label: "Forward"
      protocol: "HTTP"
      isActive: true
    - id: "c3"
      from: "api"
      to: "rds"
      label: "SQL"
      protocol: "PostgreSQL"
      isActive: true
    - id: "c4"
      from: "api"
      to: "redis"
      label: "Tokens"
      protocol: "Redis TLS"
      isActive: true
    - id: "c5"
      from: "api"
      to: "smtp"
      label: "Notify"
      protocol: "SMTP/SMS"
      isActive: true
    - id: "c6"
      from: "ecr"
      to: "api"
      label: "Deploy"
      protocol: "OCI image"
      isActive: true

dataFlow:
  requestFlow:
    - number: 1
      title: "Client request"
      description: "Client sends HTTP request; Bearer token on protected routes except public auth and /health."
      icon: "send"
    - number: 2
      title: "Middleware chain"
      description: "Recovery, CORS, audit log, optional rate limit, then JWT Authenticate and RequireAnyRole on route groups."
      icon: "filter"
    - number: 3
      title: "Handler → domain service"
      description: "Gin handler validates DTO, maps to command/query, invokes core service."
      icon: "cog"
    - number: 4
      title: "sqlc repository"
      description: "Persistence adapter runs parameterized SQL via pgx against RDS."
      icon: "database"
    - number: 5
      title: "APIResponse JSON"
      description: "shared/http helpers return success, pagination meta, or typed ApplicationError mapping."
      icon: "reply"

  eventFlow:
    - number: 1
      title: "Domain action"
      description: "e.g. appointment confirmed or user registered triggers notification service call."
      icon: "zap"
    - number: 2
      title: "Notification record"
      description: "Row inserted in notifications table via sqlc repository."
      icon: "inbox"
    - number: 3
      title: "Channel dispatch (when wired)"
      description: "EmailSender or SMSender implementation delivers via SMTP or Twilio."
      icon: "mail"
    - number: 4
      title: "Client poll"
      description: "User reads /api/v2/me/notifications for delivery status in MVP."
      icon: "eye"

techDecisions:
  - title: "Go monolith vs microservices"
    problem: "Solo/small team needs one deployable unit and strong typing without network chatter between services."
    solution: "Single module clinic-vet-api with hexagonal internal packages; scale ECS tasks horizontally."
    alternatives:
      - "Microservices per domain"
      - "Serverless Lambda per route"
    outcome: "Faster iteration; sqlc keeps SQL explicit and reviewable."
    icon: "layers"
  - title: "sqlc + pgx instead of ORM"
    problem: "Complex clinical queries and performance predictability for appointment search."
    solution: "SQL files in db/queries generate type-safe Go; migrations in db/migrations."
    alternatives:
      - "GORM"
      - "Ent"
    outcome: "Compile-time query checks; easy RDS tuning with raw SQL."
    icon: "database"
  - title: "Gin over net/http alone"
    problem: "Need routing groups, middleware, and fast JSON binding."
    solution: "Gin with validator integration and swagger-friendly handler structure."
    alternatives:
      - "Echo"
      - "Chi"
    outcome: "Familiar ecosystem; middleware matches portfolio docs patterns."
    icon: "router"
  - title: "Upstash Redis for tokens"
    problem: "Multiple API replicas must share token revocation state."
    solution: "REDIS_URL with TLS (rediss://) for TokenService; local Redis in compose.local."
    alternatives:
      - "PostgreSQL session table only"
      - "ElastiCache in-VPC"
    outcome: "Matches serverless-friendly dev/prod parity; ElastiCache optional later in VPC."
    icon: "redis"
  - title: "Docker multi-stage Alpine"
    problem: "Small attack surface image for ECS with migrations bundled."
    solution: "Builder compiles static binary; runtime runs entrypoint.sh + migrate + main."
    alternatives:
      - "Distroless without shell"
      - "Deploy bare binary on EC2"
    outcome: "~minimal image with migration tooling included."
    icon: "docker"
---

# Architecture

> **Veterinary domain:** Customer-facing routes always scope data to the authenticated pet owner; staff routes require elevated JWT roles once role names are aligned with router guards.

> **Deploy alignment:** Diagram reflects **live** RDS + Upstash and **planned** ALB/ECS/ECR path documented in ProjectInfrastructure.md.

> **Technical debt:** No async worker—notifications are synchronous DB writes until email/SMS adapters are injected. Medical handler nil blocks clinical HTTP despite route definitions. Consider Redis-backed rate limiter for multi-instance ECS.
