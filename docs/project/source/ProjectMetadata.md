---
projectId: "veterinary-clinic-api"
featured: true
name: "Veterinary Clinic API"
language: "Go"
category: "backend"
framework: "Gin"
version: "2.0.0"
repositoryUrl: "https://github.com/alexisTrejo11/veterinary-clinic-api"
liveDemoUrl: "https://api.vet-clinic.example.com/health"
description: "REST API for veterinary clinic operations—customer and pet onboarding, staff scheduling, clinical sessions (vaccinations, prescriptions, surgeries), payments, and multi-channel notifications. Built with hexagonal architecture, PostgreSQL (sqlc), and Redis-backed JWT sessions."
techStack:
  - "Go 1.24"
  - "Gin 1.10"
  - "PostgreSQL 12+ (pgx / sqlc)"
  - "Redis (go-redis) — JWT sessions & token revocation"
  - "golang-migrate"
  - "Swagger (swaggo)"
  - "Twilio SMS"
  - "SMTP email"
  - "Docker multi-stage (Alpine)"
  - "AWS RDS PostgreSQL (production)"
  - "Upstash Redis (production)"
status: "stable"
createdAt: "2025-01-15T00:00:00.000Z"
updatedAt: "2026-06-01T00:00:00.000Z"
---

# Project Metadata

> Portfolio metadata for the Veterinary Clinic API. Replace `api.vet-clinic.example.com` with your production hostname after AWS deploy.

> **Highlight:** Cloud data plane is already in place (RDS + Upstash); the remaining work is packaging the Go binary on AWS (ECS Fargate or EC2 + ALB) using `docker/Dockerfile` and `docker/compose.dev.yml`.

> **Dangerous:** Never commit `.env` with real `DATABASE_PASSWORD`, `JWT_SECRET`, or Twilio tokens. Rotate secrets if they were ever pushed.

> **Missing:** No `LICENSE` file referenced in README template yet—add before public portfolio release.
