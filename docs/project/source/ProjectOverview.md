---
problemStatement:
  problemTitle: "Disconnected tools for small and mid-size vet clinics"
  problemDescription: "Veterinary practices often run reception desks, clinical records, pet owner portals, and billing on separate spreadsheets or legacy software. That leads to double bookings, lost vaccination history, weak access control on pet medical data, and no single API for a future mobile app or admin dashboard."
  problemList:
    - "Staff cannot see a unified schedule across veterinarians and receptionists"
    - "Pet owners lack self-service booking and read-only access to their pets' clinical history"
    - "Clinical data (sessions, vaccines, prescriptions) scattered across paper or ad-hoc files"
    - "Payments and appointment status not linked in one transactional system"
    - "No standardized REST contract for integrators or a future client portal"

solution:
  solutionTitle: "One Go API for end-to-end clinic workflows"
  solutionList:
    - title: "JWT auth with role-aware routes"
      description: "Register, activate, login, refresh, logout, 2FA hooks, and password reset with tokens stored in Redis; Gin middleware enforces Bearer JWT on protected routes."
    - title: "Appointments with domain rules"
      description: "Customers request visits; employees manage assigned appointments; admins/managers search and operate clinic-wide—overlap detection (30 min buffer) and hourly capacity (max 5 concurrent)."
    - title: "Customers, employees, and pets"
      description: "CRUD for clinic customers and staff profiles; pet records tied to customers with species/gender validation and soft-delete restore flows."
    - title: "Medical sessions (clinical module)"
      description: "Sessions with vaccinations, surgeries, prescriptions, attachments, services, and catalogs—customer read-only `/me/medical` vs staff write `/api/v2/medical` (service wiring in progress)."
    - title: "Payments and notifications"
      description: "Payment lifecycle (create, process, cancel, refund) and notification inbox plus staff send/summary endpoints; email/SMS senders pluggable via env."
    - title: "Container-ready deploy"
      description: "Multi-stage Docker image, golang-migrate on startup, local full stack via `./docker/up-local.sh`, dev profile against cloud RDS + Upstash."

keyMetrics:
  metricsTitle: "Platform snapshot"
  metricsList:
    - "~182 Go source files (excluding generated sqlc)"
    - "9 SQL migrations under db/migrations/"
    - "REST API v2 prefix on most routes (`/api/v2/...`)"
    - "Roles: admin, customer, veterinarian, receptionist (JWT claim `role`)"
    - "Health probe at GET /health (no auth)"
    - "Swagger annotations in handlers; UI when ENABLE_SWAGGER=true (non-production)"

links:
  github: "https://github.com/alexisTrejo11/veterinary-clinic-api"
  demo: "https://api.vet-clinic.example.com/health"
  documentation: "https://api.vet-clinic.example.com/swagger/index.html"
  dockerHub: "https://hub.docker.com/r/YOUR_DOCKERHUB_USER/veterinary-clinic-api"

mediaGallery:
  title: "Veterinary Clinic API — product views"
  description: "Placeholder assets for portfolio. Replace with screenshots from Swagger UI, admin dashboard, or deployment diagram after AWS go-live."
  items:
    - type: "image"
      url: "https://placehold.co/1200x630/1E3A5F/ffffff?text=Veterinary+Clinic+API"
      thumbnail: "https://placehold.co/400x210/1E3A5F/ffffff?text=Vet+API"
      title: "API cover"
      description: "Veterinary Clinic Management REST API"
      alt: "Veterinary Clinic API branding placeholder"
      category: "screenshot"
    - type: "image"
      url: "https://placehold.co/1200x800/2563EB/ffffff?text=Swagger+OpenAPI"
      thumbnail: "https://placehold.co/400x267/2563EB/ffffff?text=OpenAPI"
      title: "OpenAPI / Swagger"
      description: "Interactive API docs generated with swaggo (development mode)"
      alt: "Swagger UI placeholder"
      category: "demo"

mediaItems:
  - type: "image"
    url: "https://placehold.co/800x500/059669/ffffff?text=Architecture"
    thumbnail: "https://placehold.co/320x200/059669/ffffff?text=Arch"
    title: "Target production architecture"
    description: "Clients → ALB → Docker API on AWS → RDS PostgreSQL + Upstash Redis"
    alt: "Architecture diagram placeholder"
    category: "architecture"
  - type: "image"
    url: "https://placehold.co/800x500/DC2626/ffffff?text=Docker+Local"
    thumbnail: "https://placehold.co/320x200/DC2626/ffffff?text=Docker"
    title: "Local Docker stack"
    description: "compose.local.yml runs Postgres 12, Redis 7, and API with migrations on boot"
    alt: "Docker local stack placeholder"
    category: "diagram"

metrics:
  - label: "Go module"
    value: "clinic-vet-api"
    description: "Monolith binary with internal/core domain packages"
  - label: "API version"
    value: "v2"
    description: "Primary prefix /api/v2/ on most route groups"
  - label: "Runtime (Docker)"
    value: "Go 1.24 Alpine"
    description: "docker/Dockerfile multi-stage build"
  - label: "Auth"
    value: "JWT + Redis"
    description: "Access/refresh tokens; SessionRepository wiring still TODO in bootstrap"
  - label: "Database access"
    value: "sqlc"
    description: "Type-safe queries generated from db/queries"
  - label: "Deploy status"
    value: "RDS + Redis live"
    description: "AWS application tier (ECS/EC2) documented, not yet automated in repo"
---

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
