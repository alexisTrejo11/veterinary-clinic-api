---
codeExamples:
  - id: "appointment-overlap-rules"
    title: "Employee appointment overlap validation"
    description: "AppointmentDomainService prevents double-booking the same veterinarian within a 30-minute window, excluding cancelled and no-show visits."
    category: "domain"
    duration: "4 min read"
    views: 0
    tags:
      - "appointments"
      - "validation"
      - "domain-service"
    files:
      - name: "domain_service.go"
        path: "internal/core/appointments/domain_service.go"
        language: "go"
        highlighted: true
        explanation: "Repository-backed rule: ±30 minutes around scheduledDate for the same employee_id."
        content: |
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

  - id: "jwt-middleware"
    title: "Stateless JWT authentication middleware"
    description: "Parses Bearer tokens, validates claims, and attaches UserContext to Gin without database round-trips per request."
    category: "security"
    duration: "3 min read"
    views: 0
    tags:
      - "jwt"
      - "gin"
      - "middleware"
    files:
      - name: "auth_middleware.go"
        path: "internal/middleware/auth_middleware.go"
        language: "go"
        highlighted: true
        explanation: "Uses github.com/golang-jwt/jwt/v5 MapClaims for user id, email, and role."
        content: |
          func (am *AuthMiddleware) Authenticate() gin.HandlerFunc {
              return func(c *gin.Context) {
                  raw := c.GetHeader("Authorization")
                  tokenStr := strings.TrimPrefix(raw, bearerPrefix)
                  claims, err := am.parseAccessToken(tokenStr)
                  // sets userID, userEmail, userRole on context
                  c.Next()
              }
          }

  - id: "api-response-envelope"
    title: "Standardized JSON API responses"
    description: "All handlers use shared/http helpers returning a consistent success and error shape for frontends."
    category: "api"
    duration: "2 min read"
    views: 0
    tags:
      - "http"
      - "gin"
    files:
      - name: "response.go"
        path: "internal/shared/http/response.go"
        language: "go"
        highlighted: true
        explanation: "APIResponse includes Success, Data, Message, Error, Meta, Timestamp, RequestID."
        content: |
          type APIResponse struct {
              Success   bool       `json:"success"`
              Data      any        `json:"data,omitempty"`
              Message   string     `json:"message,omitempty"`
              Error     *ErrorInfo `json:"error,omitempty"`
              Meta      any        `json:"meta,omitempty"`
              Timestamp time.Time  `json:"timestamp"`
              RequestID string     `json:"request_id,omitempty"`
          }

  - id: "database-url-rds"
    title: "RDS-safe DATABASE_URL parsing"
    description: "Normalizes jdbc:postgresql:// URLs and enforces sslmode=require for Amazon RDS endpoints."
    category: "infrastructure"
    duration: "3 min read"
    views: 0
    tags:
      - "postgres"
      - "aws"
      - "config"
    files:
      - name: "database_url.go"
        path: "internal/config/database_url.go"
        language: "go"
        highlighted: true
        explanation: "Used by app config and parse-database-url.sh in Docker entrypoint."
        content: |
          // Strips jdbc: prefix, merges user/password from env,
          // adds sslmode=require when host contains .rds.amazonaws.com

  - id: "hexagonal-bootstrap"
    title: "Wiring repositories to HTTP router"
    description: "Bootstrap() constructs sqlc repositories, domain services, handlers, and registers routes via APIRouter."
    category: "architecture"
    duration: "5 min read"
    views: 0
    tags:
      - "bootstrap"
      - "dependency-injection"
    files:
      - name: "depencies.go"
        path: "internal/infrastructure/http/depencies.go"
        language: "go"
        highlighted: true
        explanation: "Single composition root; auth requires Redis + JWT_SECRET; medical handler intentionally nil until wired."
        content: |
          func Bootstrap(engine *gin.Engine, queries *sqlc.Queries, ...) error {
              petSvc := pets.NewPetService(petRepo, petCustomerRepo)
              apiRouter, _ := NewAPIRouter(appHandlers, config)
              apiRouter.RegisterRoutes()
              return nil
          }

  - id: "docker-entrypoint-migrate"
    title: "Migrate-then-run container entrypoint"
    description: "Production-friendly startup: wait for Postgres, run golang-migrate up, then exec the API binary."
    category: "infrastructure"
    duration: "2 min read"
    views: 0
    tags:
      - "docker"
      - "migrations"
    files:
      - name: "entrypoint.sh"
        path: "scripts/entrypoint.sh"
        language: "shell"
        highlighted: true
        explanation: "SKIP_MIGRATIONS and SKIP_DB_WAIT escape hatches for cloud dev profile."
        content: |
          . /scripts/parse-database-url.sh
          until pg_isready -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER"; do sleep 2; done
          migrate -path ./db/migrations -database "$DATABASE_URL" up
          exec "$@"
---

# Code Showcase

> Snippets are abbreviated from the repository; open the referenced paths for full error handling, tests, and sqlc types.

> **Recommended reading order:** database URL parsing → bootstrap wiring → JWT middleware → appointment domain rules → API response envelope → Docker entrypoint.

> **Dangerous:** `000007_insert_demo_data.up.sql` should not run against production RDS—review migration list before first deploy.

> **Next contributor task:** Add `petRoutes()`, `paymentRoutes()`, and `addressRoutes()` in `router.go` mirroring existing customer/employee patterns, then wire `MedicalHandler` in `depencies.go`.
