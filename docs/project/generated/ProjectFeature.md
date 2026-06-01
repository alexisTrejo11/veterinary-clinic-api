# Project Features

## JWT authentication & RBAC

Email/password registration with role-specific flows (admin, customer, veterinarian, receptionist), account activation, login, refresh, logout, 2FA endpoints, and password reset. Stateless Gin middleware reads role from JWT claims.

| Property | Value |
| --- | --- |
| ID | jwt-auth-rbac |
| Category | authentication |
| Status | stable |
| Icon | shield-lock |

### Highlights

- Public routes under /api/v2/auth (register, login, activate)
- Protected /auth/* for refresh, logout, 2FA, reset-password
- MAX_LOGIN_ATTEMPTS and LOCKOUT_DURATION from env

### Tech stack

- github.com/golang-jwt/jwt/v5
- internal/core/auth
- internal/middleware/auth_middleware.go

### Metrics

| Label | Value | Trend |
| --- | --- | --- |
| Access token TTL | 24h default | stable |
| Refresh token TTL | 168h default | stable |

### Code snippet

_internal/middleware/auth_middleware.go_

```go
func (am *AuthMiddleware) RequireAnyRole(roles ...string) gin.HandlerFunc {
    userRole, _ := GetUserRoleFromContext(c)
    for _, r := range roles {
        if userRole == r { c.Next(); return }
    }
    sharedhttp.Forbidden(c, ...)
}
```

## Appointment scheduling & lifecycle

Customers request appointments on /api/v2/me/appointments; employees manage assigned visits; admins search clinic-wide calendar with confirm, complete, cancel, reschedule, and not-attend actions.

| Property | Value |
| --- | --- |
| ID | appointments-scheduling |
| Category | api |
| Status | stable |
| Icon | calendar-check |

### Highlights

- 30-minute employee overlap guard
- Max 5 active appointments per hour (clinic capacity)
- Nested lists by customer, employee, and pet IDs

### Tech stack

- internal/core/appointments
- internal/infrastructure/http/router.go

### Metrics

| Label | Value | Trend |
| --- | --- | --- |
| Capacity limit | 5 / hour | stable |

### Code snippet

_internal/core/appointments/domain_service.go_

```go
const maxConcurrentAppointments = 5
func (s *appointmentDomainService) ValidateCapacity(...) error {
    if activeCount >= maxConcurrentAppointments {
        return fmt.Errorf("maximum capacity reached ...")
    }
    return nil
}
```

## Customers & employees CRUD

Manager/admin search, create, update, soft-delete, and restore for clinic customers and veterinary staff profiles linked to user accounts.

| Property | Value |
| --- | --- |
| ID | customers-employees |
| Category | api |
| Status | stable |
| Icon | users |

### Highlights

- /api/v2/customers and /api/v2/employees
- Pagination via shared page.Page and X-Total-Count header
- Employee specialties enum (surgery, dentistry, emergency, etc.)

### Tech stack

- internal/core/customers
- internal/core/employees
- sqlc

## Pet registry (domain ready)

Pet CRUD with species/gender validation, customer scoping, and my-pets endpoints for pet owners—handlers implemented; HTTP routes pending registration in router.go.

| Property | Value |
| --- | --- |
| ID | pets-registry |
| Category | database |
| Status | beta |
| Icon | paw |

### Highlights

- PetService with customer existence checks
- Planned paths: /api/v2/pets, /api/v2/me/pets
- Linked to appointments and medical sessions

### Tech stack

- internal/core/pets
- db/migrations/000004_pets_related.up.sql

## Clinical sessions & catalogs

Medical sessions with vaccinations, surgeries, prescriptions, attachments, and service lines plus vaccine/medication/service catalogs. Router and handlers exist; bootstrap sets MedicalHandler to nil until repositories are wired.

| Property | Value |
| --- | --- |
| ID | medical-sessions |
| Category | api |
| Status | beta |
| Icon | file-medical |

### Highlights

- Customer read-only /api/v2/me/medical/*
- Staff write /api/v2/medical/*
- Soft and hard delete for sessions

### Tech stack

- internal/core/medical
- db/migrations/000005_appointments_med_sessions.up.sql

## Payments ledger (domain ready)

Create, update, process, cancel, and refund payments with search by specification and customer—PaymentHandler complete; routes not yet mounted on Gin router.

| Property | Value |
| --- | --- |
| ID | payments-ledger |
| Category | integration |
| Status | beta |
| Icon | credit-card |

### Highlights

- Payment status enum and command/query services
- Planned /api/v2/payments/* (admin/manager)
- No Stripe dependency—clinic payment records in PostgreSQL

### Tech stack

- internal/core/payments
- sqlc/payments.sql

## In-app notifications

Users read /api/v2/me/notifications; staff query by type/channel, view summary, and POST manual notifications. EmailSender and SMSender nil until SMTP/Twilio adapters are injected.

| Property | Value |
| --- | --- |
| ID | notifications |
| Category | messaging |
| Status | beta |
| Icon | bell |

### Highlights

- Notification repository via sqlc
- Twilio and SMTP env vars in .env.example
- DB table from 000009_notification_table.up.sql

### Tech stack

- internal/core/notifications
- github.com/twilio/twilio-go

## Global rate limit & audit logging

Optional Gin rate limiter (in-memory sliding window per IP) and audit middleware that logs request/response metadata with UUID request IDs.

| Property | Value |
| --- | --- |
| ID | rate-limiting-audit |
| Category | security |
| Status | stable |
| Icon | speedometer |

### Highlights

- RATE_LIMIT_ENABLED, RATE_LIMIT_RPS, RATE_LIMIT_BURST
- Expose X-Rate-Limit-Remaining header
- AuditLog middleware on all routes

### Tech stack

- internal/middleware/rate_limiter.go
- internal/middleware/audit_log.go

### Code snippet

_main.go_

```go
if settings.RateLimit.Enabled {
    router.Use(middleware.RateLimiter(rateLimitConfig))
}
router.Use(middleware.AuditLog())
```

## Type-safe PostgreSQL access

sqlc generates Go from SQL queries; pgx connection pool with configurable max connections and RDS-friendly SSL mode parsing in database_url.go.

| Property | Value |
| --- | --- |
| ID | sqlc-postgres |
| Category | database |
| Status | stable |
| Icon | database |

### Highlights

- Migrations via golang-migrate in Docker entrypoint
- jdbc:postgresql:// prefix supported in DATABASE_URL
- 9 migration versions including demo seed

### Tech stack

- github.com/jackc/pgx/v5
- sqlc
- github.com/golang-migrate/migrate/v4

## Docker local & cloud dev profiles

Multi-stage Alpine image, compose.local full stack, compose.dev API-only against RDS/Upstash, and shell helpers to parse connection URLs inside the container.

| Property | Value |
| --- | --- |
| ID | docker-tooling |
| Category | integration |
| Status | stable |
| Icon | docker |

### Highlights

- ./docker/up-local.sh
- scripts/parse-database-url.sh and parse-redis-url.sh
- cmd/migrate for host-side migrations

### Tech stack

- docker/Dockerfile
- docker/compose.local.yml
- docker/compose.dev.yml

## Additional notes

# Project Features

> **Stable:** Auth middleware, appointments, customers, employees, users, profile, notifications (read/send paths registered), health check, Docker tooling.

> **Beta:** Pets, payments, addresses (handlers only), medical (routes guarded but handler nil), email/SMS delivery (nil senders).

> **Before AWS go-live:** Fix JWT role names vs router guards; implement SessionRepository; register pet/payment/address routes; wire MedicalHandler; tighten CORS; disable demo migration on RDS.

> **Highlight for portfolio:** Hexagonal layout (`internal/core` domain, `internal/infrastructure` adapters) with a single deployable binary—well suited to ECS without microservice overhead.

