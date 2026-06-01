# Project Overview

## Disconnected tools for small and mid-size vet clinics

Veterinary practices often run reception desks, clinical records, pet owner portals, and billing on separate spreadsheets or legacy software. That leads to double bookings, lost vaccination history, weak access control on pet medical data, and no single API for a future mobile app or admin dashboard.

### Pain points

- Staff cannot see a unified schedule across veterinarians and receptionists
- Pet owners lack self-service booking and read-only access to their pets' clinical history
- Clinical data (sessions, vaccines, prescriptions) scattered across paper or ad-hoc files
- Payments and appointment status not linked in one transactional system
- No standardized REST contract for integrators or a future client portal

## One Go API for end-to-end clinic workflows

- **JWT auth with role-aware routes** — Register, activate, login, refresh, logout, 2FA hooks, and password reset with tokens stored in Redis; Gin middleware enforces Bearer JWT on protected routes.
- **Appointments with domain rules** — Customers request visits; employees manage assigned appointments; admins/managers search and operate clinic-wide—overlap detection (30 min buffer) and hourly capacity (max 5 concurrent).
- **Customers, employees, and pets** — CRUD for clinic customers and staff profiles; pet records tied to customers with species/gender validation and soft-delete restore flows.
- **Medical sessions (clinical module)** — Sessions with vaccinations, surgeries, prescriptions, attachments, services, and catalogs—customer read-only `/me/medical` vs staff write `/api/v2/medical` (service wiring in progress).
- **Payments and notifications** — Payment lifecycle (create, process, cancel, refund) and notification inbox plus staff send/summary endpoints; email/SMS senders pluggable via env.
- **Container-ready deploy** — Multi-stage Docker image, golang-migrate on startup, local full stack via `./docker/up-local.sh`, dev profile against cloud RDS + Upstash.

## Platform snapshot

- ~182 Go source files (excluding generated sqlc)
- 9 SQL migrations under db/migrations/
- REST API v2 prefix on most routes (`/api/v2/...`)
- Roles: admin, customer, veterinarian, receptionist (JWT claim `role`)
- Health probe at GET /health (no auth)
- Swagger annotations in handlers; UI when ENABLE_SWAGGER=true (non-production)

## Links

| Resource | URL |
| --- | --- |
| Github | https://github.com/alexisTrejo11/veterinary-clinic-api |
| Demo | https://api.vet-clinic.example.com/health |
| Documentation | https://api.vet-clinic.example.com/swagger/index.html |
| Dockerhub | https://hub.docker.com/r/YOUR_DOCKERHUB_USER/veterinary-clinic-api |

## Veterinary Clinic API — product views

Placeholder assets for portfolio. Replace with screenshots from Swagger UI, admin dashboard, or deployment diagram after AWS go-live.

### API cover

Veterinary Clinic Management REST API

- **Type:** image | **Category:** screenshot
- ![Veterinary Clinic API branding placeholder](https://placehold.co/1200x630/1E3A5F/ffffff?text=Veterinary+Clinic+API)

### OpenAPI / Swagger

Interactive API docs generated with swaggo (development mode)

- **Type:** image | **Category:** demo
- ![Swagger UI placeholder](https://placehold.co/1200x800/2563EB/ffffff?text=Swagger+OpenAPI)

## Additional media

### Target production architecture

Clients → ALB → Docker API on AWS → RDS PostgreSQL + Upstash Redis

### Local Docker stack

compose.local.yml runs Postgres 12, Redis 7, and API with migrations on boot

## Metrics

| Label | Value | Description |
| --- | --- | --- |
| Go module | clinic-vet-api | Monolith binary with internal/core domain packages |
| API version | v2 | Primary prefix /api/v2/ on most route groups |
| Runtime (Docker) | Go 1.24 Alpine | docker/Dockerfile multi-stage build |
| Auth | JWT + Redis | Access/refresh tokens; SessionRepository wiring still TODO in bootstrap |
| Database access | sqlc | Type-safe queries generated from db/queries |
| Deploy status | RDS + Redis live | AWS application tier (ECS/EC2) documented, not yet automated in repo |

## Additional notes

# Overview

> **Audience:** Veterinary clinic operators, pet owners (customer role), veterinarians, receptionists, and platform integrators.

> **Production posture (documented as live):** PostgreSQL on **Amazon RDS** and **Upstash Redis** are configured via `DATABASE_URL` and `REDIS_URL` in `.env`. The API container is the deployable unit described in `docker/README.md`.

> **Critical gaps before calling the API feature-complete:**
> 1. **Role name mismatch:** `router.go` guards many routes with `manager` and `employee`, but JWT/user domain roles are `admin`, `customer`, `veterinarian`, `receptionist`—staff routes may return 403 until aligned.
> 2. **Unmounted handlers:** Pets, payments, and addresses handlers are wired in `Bootstrap()` but **no routes** are registered in `router.go` yet.
> 3. **Medical module:** Routes exist in `router.go` but `MedicalHandler` is `nil` in bootstrap—clinical HTTP API is not active until services are injected.
> 4. **Auth session:** `SessionRepository` is passed as `nil` in `depencies.go`—login/refresh may fail until Redis session store is implemented.
> 5. **Path prefixes:** Public auth uses `/api/v2/auth/*` while logout/refresh use `/auth/*` without the v2 prefix—configure API gateway rules accordingly.
> 6. **Swagger in production:** `main.go` only enables Swagger when not production; use dev/staging for interactive docs or export `docs/swagger.json` to a static site.

