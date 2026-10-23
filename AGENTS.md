# AGENTS.md

## Project Overview
This repository contains the backend for Campus Event Management application built with Go, Fiber v3, PostgreSQL, pgx/pgxpool, and Clean Architecture.

## Architecture Guidelines
- **Clean Architecture Layers**:
  - `internal/handler`: HTTP request parsing, validation, and response formatting (Fiber v3 handlers).
  - `internal/service`: Business logic, use-case coordination, validation rules.
  - `internal/repository`: Database queries and data access logic using `pgxpool`.
  - `internal/domain`: Core domain entities, models, and repository interfaces.
  - `internal/config`: Configuration loaded from environment variables.
- Dependency flow: Handler -> Service -> Repository.
- Service and Repository interfaces are declared in the domain or service layer to ensure loose coupling and testability.
- No Fiber-specific context or HTTP types allowed inside Service or Repository layers.
- Parameterized SQL queries must always be used to prevent SQL injection.

## Directory Structure
```
event-management-be/
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── database/
│   │   └── postgres.go
│   ├── domain/
│   │   ├── event.go
│   │   └── errors.go
│   ├── handler/
│   │   ├── event_handler.go
│   │   ├── health_handler.go
│   │   └── response.go
│   ├── repository/
│   │   ├── event_repository.go
│   │   └── postgres/
│   │       └── event_repository.go
│   └── service/
│       ├── event_service.go
│       └── event_service_test.go
├── migrations/
│   └── 000001_create_events_table.up.sql
├── .env.example
├── .gitignore
├── compose.yaml
├── Dockerfile
├── go.mod
├── go.sum
└── README.md
```

## Migration Workflow
Database migrations are plain SQL files located in `migrations/`.
For automatic initialization and beginner-friendliness:
- The SQL migration scripts can be mounted to `/docker-entrypoint-initdb.d/` in PostgreSQL docker container or executed at application startup if configured.
- In this project, migrations are applied automatically during app startup (or via PostgreSQL init scripts for a fresh instance), ensuring zero manual DBeaver table creation.

## Coding Conventions
- Standard Go formatting (`gofmt`).
- Error handling idiomatic to Go (`errors.Is`, custom domain errors).
- UTC/RFC3339 timestamps with `TIMESTAMPTZ` in PostgreSQL.
- Keep scope strictly confined to requested features.
