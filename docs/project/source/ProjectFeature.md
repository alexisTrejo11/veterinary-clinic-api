---
features:
  - id: "jwt-auth-rbac"
    title: "JWT authentication & RBAC"
    description: "Email/password registration with role-specific flows (admin, customer, veterinarian, receptionist), account activation, login, refresh, logout, 2FA endpoints, and password reset. Stateless Gin middleware reads role from JWT claims."
    icon: "shield-lock"
    category: "authentication"
    status: "stable"
    highlights:
      - "Public routes under /api/v2/auth (register, login, activate)"
      - "Protected /auth/* for refresh, logout, 2FA, reset-password"
      - "MAX_LOGIN_ATTEMPTS and LOCKOUT_DURATION from env"
    techStack:
      - "github.com/golang-jwt/jwt/v5"
      - "internal/core/auth"
      - "internal/middleware/auth_middleware.go"
    metrics:
      - label: "Access token TTL"
        value: "24h default"
        trend: "stable"
        icon: "clock"
      - label: "Refresh token TTL"
        value: "168h default"
        trend: "stable"
        icon: "refresh"
    codeSnippet:
      language: "go"
      filename: "internal/middleware/auth_middleware.go"
      code: |
        func (am *AuthMiddleware) RequireAnyRole(roles ...string) gin.HandlerFunc {
            userRole, _ := GetUserRoleFromContext(c)
            for _, r := range roles {
                if userRole == r { c.Next(); return }
            }
            sharedhttp.Forbidden(c, ...)
        }

  - id: "appointments-scheduling"
    title: "Appointment scheduling & lifecycle"
    description: "Customers request appointments on /api/v2/me/appointments; employees manage assigned visits; admins search clinic-wide calendar with confirm, complete, cancel, reschedule, and not-attend actions."
    icon: "calendar-check"
    category: "api"
    status: "stable"
    highlights:
      - "30-minute employee overlap guard"
      - "Max 5 active appointments per hour (clinic capacity)"
      - "Nested lists by customer, employee, and pet IDs"
    techStack:
      - "internal/core/appointments"
      - "internal/infrastructure/http/router.go"
    metrics:
      - label: "Capacity limit"
        value: "5 / hour"
        trend: "stable"
        icon: "gauge"
    codeSnippet:
      language: "go"
      filename: "internal/core/appointments/domain_service.go"
      code: |
        const maxConcurrentAppointments = 5
        func (s *appointmentDomainService) ValidateCapacity(...) error {
            if activeCount >= maxConcurrentAppointments {
                return fmt.Errorf("maximum capacity reached ...")
            }
            return nil
        }

  - id: "customers-employees"
    title: "Customers & employees CRUD"
    description: "Manager/admin search, create, update, soft-delete, and restore for clinic customers and veterinary staff profiles linked to user accounts."
    icon: "users"
    category: "api"
    status: "stable"
    highlights:
      - "/api/v2/customers and /api/v2/employees"
      - "Pagination via shared page.Page and X-Total-Count header"
      - "Employee specialties enum (surgery, dentistry, emergency, etc.)"
    techStack:
      - "internal/core/customers"
      - "internal/core/employees"
      - "sqlc"

  - id: "pets-registry"
    title: "Pet registry (domain ready)"
    description: "Pet CRUD with species/gender validation, customer scoping, and my-pets endpoints for pet owners—handlers implemented; HTTP routes pending registration in router.go."
    icon: "paw"
    category: "database"
    status: "beta"
    highlights:
      - "PetService with customer existence checks"
      - "Planned paths: /api/v2/pets, /api/v2/me/pets"
      - "Linked to appointments and medical sessions"
    techStack:
      - "internal/core/pets"
      - "db/migrations/000004_pets_related.up.sql"

  - id: "medical-sessions"
    title: "Clinical sessions & catalogs"
    description: "Medical sessions with vaccinations, surgeries, prescriptions, attachments, and service lines plus vaccine/medication/service catalogs. Router and handlers exist; bootstrap sets MedicalHandler to nil until repositories are wired."
    icon: "file-medical"
    category: "api"
    status: "beta"
    highlights:
      - "Customer read-only /api/v2/me/medical/*"
      - "Staff write /api/v2/medical/*"
      - "Soft and hard delete for sessions"
    techStack:
      - "internal/core/medical"
      - "db/migrations/000005_appointments_med_sessions.up.sql"

  - id: "payments-ledger"
    title: "Payments ledger (domain ready)"
    description: "Create, update, process, cancel, and refund payments with search by specification and customer—PaymentHandler complete; routes not yet mounted on Gin router."
    icon: "credit-card"
    category: "integration"
    status: "beta"
    highlights:
      - "Payment status enum and command/query services"
      - "Planned /api/v2/payments/* (admin/manager)"
      - "No Stripe dependency—clinic payment records in PostgreSQL"
    techStack:
      - "internal/core/payments"
      - "sqlc/payments.sql"

  - id: "notifications"
    title: "In-app notifications"
    description: "Users read /api/v2/me/notifications; staff query by type/channel, view summary, and POST manual notifications. EmailSender and SMSender nil until SMTP/Twilio adapters are injected."
    icon: "bell"
    category: "messaging"
    status: "beta"
    highlights:
      - "Notification repository via sqlc"
      - "Twilio and SMTP env vars in .env.example"
      - "DB table from 000009_notification_table.up.sql"
    techStack:
      - "internal/core/notifications"
      - "github.com/twilio/twilio-go"

  - id: "rate-limiting-audit"
    title: "Global rate limit & audit logging"
    description: "Optional Gin rate limiter (in-memory sliding window per IP) and audit middleware that logs request/response metadata with UUID request IDs."
    icon: "speedometer"
    category: "security"
    status: "stable"
    highlights:
      - "RATE_LIMIT_ENABLED, RATE_LIMIT_RPS, RATE_LIMIT_BURST"
      - "Expose X-Rate-Limit-Remaining header"
      - "AuditLog middleware on all routes"
    techStack:
      - "internal/middleware/rate_limiter.go"
      - "internal/middleware/audit_log.go"
    codeSnippet:
      language: "go"
      filename: "main.go"
      code: |
        if settings.RateLimit.Enabled {
            router.Use(middleware.RateLimiter(rateLimitConfig))
        }
        router.Use(middleware.AuditLog())

  - id: "sqlc-postgres"
    title: "Type-safe PostgreSQL access"
    description: "sqlc generates Go from SQL queries; pgx connection pool with configurable max connections and RDS-friendly SSL mode parsing in database_url.go."
    icon: "database"
    category: "database"
    status: "stable"
    highlights:
      - "Migrations via golang-migrate in Docker entrypoint"
      - "jdbc:postgresql:// prefix supported in DATABASE_URL"
      - "9 migration versions including demo seed"
    techStack:
      - "github.com/jackc/pgx/v5"
      - "sqlc"
      - "github.com/golang-migrate/migrate/v4"

  - id: "docker-tooling"
    title: "Docker local & cloud dev profiles"
    description: "Multi-stage Alpine image, compose.local full stack, compose.dev API-only against RDS/Upstash, and shell helpers to parse connection URLs inside the container."
    icon: "docker"
    category: "integration"
    status: "stable"
    highlights:
      - "./docker/up-local.sh"
      - "scripts/parse-database-url.sh and parse-redis-url.sh"
      - "cmd/migrate for host-side migrations"
    techStack:
      - "docker/Dockerfile"
      - "docker/compose.local.yml"
      - "docker/compose.dev.yml"
---

# Project Features

> **Stable:** Auth middleware, appointments, customers, employees, users, profile, notifications (read/send paths registered), health check, Docker tooling.

> **Beta:** Pets, payments, addresses (handlers only), medical (routes guarded but handler nil), email/SMS delivery (nil senders).

> **Before AWS go-live:** Fix JWT role names vs router guards; implement SessionRepository; register pet/payment/address routes; wire MedicalHandler; tighten CORS; disable demo migration on RDS.

> **Highlight for portfolio:** Hexagonal layout (`internal/core` domain, `internal/infrastructure` adapters) with a single deployable binary—well suited to ECS without microservice overhead.
