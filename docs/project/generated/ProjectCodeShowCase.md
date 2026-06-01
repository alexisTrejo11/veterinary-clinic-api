# Code Showcase

## Employee appointment overlap validation

AppointmentDomainService prevents double-booking the same veterinarian within a 30-minute window, excluding cancelled and no-show visits.

**Category:** domain | **Duration:** 4 min read | **Tags:** appointments, validation, domain-service

### domain_service.go

**Path:** `internal/core/appointments/domain_service.go`

Repository-backed rule: ±30 minutes around scheduledDate for the same employee_id.

```go
func (s *appointmentDomainService) ValidateNoOverlapping(ctx context.Context, appointment *Appointment) error {
    if appointment.EmployeeID == nil {
        return nil
    }
    startTime := appointment.ScheduledDate.Add(-30 * time.Minute)
    endTime := appointment.ScheduledDate.Add(30 * time.Minute)
    spec := NewAppointmentSpecification().
        WithEmployeeID(*appointment.EmployeeID).
        WithDateRange(startTime, endTime)
    page, err := s.repository.Find(ctx, spec)
    // ... skip same id, cancelled; fail if within 30 minutes
    return nil
}
```

## Stateless JWT authentication middleware

Parses Bearer tokens, validates claims, and attaches UserContext to Gin without database round-trips per request.

**Category:** security | **Duration:** 3 min read | **Tags:** jwt, gin, middleware

### auth_middleware.go

**Path:** `internal/middleware/auth_middleware.go`

Uses github.com/golang-jwt/jwt/v5 MapClaims for user id, email, and role.

```go
func (am *AuthMiddleware) Authenticate() gin.HandlerFunc {
    return func(c *gin.Context) {
        raw := c.GetHeader("Authorization")
        tokenStr := strings.TrimPrefix(raw, bearerPrefix)
        claims, err := am.parseAccessToken(tokenStr)
        // sets userID, userEmail, userRole on context
        c.Next()
    }
}
```

## Standardized JSON API responses

All handlers use shared/http helpers returning a consistent success and error shape for frontends.

**Category:** api | **Duration:** 2 min read | **Tags:** http, gin

### response.go

**Path:** `internal/shared/http/response.go`

APIResponse includes Success, Data, Message, Error, Meta, Timestamp, RequestID.

```go
type APIResponse struct {
    Success   bool       `json:"success"`
    Data      any        `json:"data,omitempty"`
    Message   string     `json:"message,omitempty"`
    Error     *ErrorInfo `json:"error,omitempty"`
    Meta      any        `json:"meta,omitempty"`
    Timestamp time.Time  `json:"timestamp"`
    RequestID string     `json:"request_id,omitempty"`
}
```

## RDS-safe DATABASE_URL parsing

Normalizes jdbc:postgresql:// URLs and enforces sslmode=require for Amazon RDS endpoints.

**Category:** infrastructure | **Duration:** 3 min read | **Tags:** postgres, aws, config

### database_url.go

**Path:** `internal/config/database_url.go`

Used by app config and parse-database-url.sh in Docker entrypoint.

```go
// Strips jdbc: prefix, merges user/password from env,
// adds sslmode=require when host contains .rds.amazonaws.com
```

## Wiring repositories to HTTP router

Bootstrap() constructs sqlc repositories, domain services, handlers, and registers routes via APIRouter.

**Category:** architecture | **Duration:** 5 min read | **Tags:** bootstrap, dependency-injection

### depencies.go

**Path:** `internal/infrastructure/http/depencies.go`

Single composition root; auth requires Redis + JWT_SECRET; medical handler intentionally nil until wired.

```go
func Bootstrap(engine *gin.Engine, queries *sqlc.Queries, ...) error {
    petSvc := pets.NewPetService(petRepo, petCustomerRepo)
    apiRouter, _ := NewAPIRouter(appHandlers, config)
    apiRouter.RegisterRoutes()
    return nil
}
```

## Migrate-then-run container entrypoint

Production-friendly startup: wait for Postgres, run golang-migrate up, then exec the API binary.

**Category:** infrastructure | **Duration:** 2 min read | **Tags:** docker, migrations

### entrypoint.sh

**Path:** `scripts/entrypoint.sh`

SKIP_MIGRATIONS and SKIP_DB_WAIT escape hatches for cloud dev profile.

```shell
. /scripts/parse-database-url.sh
until pg_isready -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER"; do sleep 2; done
migrate -path ./db/migrations -database "$DATABASE_URL" up
exec "$@"
```

## Additional notes

# Code Showcase

> Snippets are abbreviated from the repository; open the referenced paths for full error handling, tests, and sqlc types.

> **Recommended reading order:** database URL parsing → bootstrap wiring → JWT middleware → appointment domain rules → API response envelope → Docker entrypoint.

> **Dangerous:** `000007_insert_demo_data.up.sql` should not run against production RDS—review migration list before first deploy.

> **Next contributor task:** Add `petRoutes()`, `paymentRoutes()`, and `addressRoutes()` in `router.go` mirroring existing customer/employee patterns, then wire `MedicalHandler` in `depencies.go`.

