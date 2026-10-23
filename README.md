# Campus Event Management Backend

A Go backend boilerplate for Campus Event Management application built with Go, Fiber v3, PostgreSQL (`pgx`/`pgxpool`), Clean Architecture, Docker Compose, and SQL migrations.

---

## 1. Project Purpose & Tech Stack

- **Go (1.24+)**: High-performance compiled language.
- **Fiber v3 (`github.com/gofiber/fiber/v3`)**: Express-inspired HTTP web framework for Go built on Fasthttp.
- **PostgreSQL 16**: Relational database engine.
- **pgx (`github.com/jackc/pgx/v5`)**: Native PostgreSQL driver and connection pool (`pgxpool`) for Go.
- **Docker & Docker Compose**: Containerized execution and unified service orchestration.

---

## 2. Clean Architecture & Folder Structure

The project follows Clean Architecture principles with strict layer isolation:

```
event-management-be/
├── cmd/
│   └── api/
│       └── main.go          # Application entrypoint & dependency wiring
├── internal/
│   ├── config/
│   │   └── config.go        # Environment-based configuration
│   ├── database/
│   │   └── postgres.go      # PostgreSQL pool setup & migration runner
│   ├── domain/
│   │   ├── errors.go        # Domain error constants
│   │   └── event.go         # Domain models, types & repository interfaces
│   ├── handler/
│   │   ├── event_handler.go      # Fiber HTTP handler for event endpoints
│   │   ├── event_handler_test.go # HTTP route unit tests
│   │   ├── health_handler.go     # Fiber HTTP handler for /health
│   │   └── response.go           # Standard JSON response envelope
│   ├── repository/
│   │   ├── event_repository.go   # Repository interface contract
│   │   └── postgres/
│   │       └── event_repository.go # pgxpool parameterized SQL queries
│   └── service/
│       ├── event_service.go      # Business rules & validation logic
│       └── event_service_test.go # Unit tests for business rules
├── migrations/
│   └── 000001_create_events_table.up.sql # Table creation & constraints
├── .env.example             # Template for environment variables
├── .gitignore               # Ignored files (binaries, .env, etc.)
├── compose.yaml             # Docker Compose orchestration
├── Dockerfile               # Multi-stage Docker build for backend API
├── go.mod
├── go.sum
└── README.md
```

### Layer Responsibilities
- **Handler (`internal/handler`)**: Parses HTTP requests, performs basic schema/format checks, passes DTOs to services, and returns formatted JSON responses with appropriate HTTP status codes. Contains no database logic.
- **Service (`internal/service`)**: Contains business logic, required field validation, and category rules. Coordinates use cases. Does not depend on Fiber or HTTP types.
- **Repository (`internal/repository/postgres`)**: Executes parameterized SQL queries using `pgxpool` to prevent SQL injection.
- **Domain (`internal/domain`)**: Core data structures, entity models, category enums, and repository interface definitions.
- **Config (`internal/config`)**: Loads environment variables with sensible defaults.

---

## 3. Prerequisites

- [Go 1.24+](https://go.dev/dl/)
- [Docker & Docker Compose](https://www.docker.com/products/docker-desktop/)
- [DBeaver](https://dbeaver.io/) or your preferred PostgreSQL client (optional, for inspecting tables)
- `curl` or Postman (for testing endpoints)

---

## 4. Environment Setup

Copy `.env.example` to create your local `.env` file:

```bash
# On Linux/macOS
cp .env.example .env

# On Windows (PowerShell)
Copy-Item .env.example .env
```

Default configuration in `.env.example`:

| Variable | Default Value | Description |
| :--- | :--- | :--- |
| `PORT` | `8080` | Port for the HTTP API |
| `ENV` | `development` | Environment mode |
| `DB_HOST` | `postgres` | Hostname of the PostgreSQL service (`localhost` for running outside Docker) |
| `DB_PORT` | `5432` | PostgreSQL port |
| `DB_USER` | `postgres` | PostgreSQL username |
| `DB_PASSWORD` | `postgres` | PostgreSQL password |
| `DB_NAME` | `event_management` | PostgreSQL database name |
| `DB_SSLMODE` | `disable` | SSL mode for Postgres connection |

---

## 5. How to Start and Stop the Project

### Start Services (Build & Run)
Run Docker Compose in the project root:

```bash
docker compose up --build
```

Add `-d` to run in the background (detached mode):

```bash
docker compose up --build -d
```

### Stop Services
To stop containers without deleting database data:

```bash
docker compose down
```

> [!NOTE]
> Named volume `postgres_data` preserves database data across restarts. To perform a complete clean slate reset, run `docker compose down -v`.

---

## 6. How Database Migrations Work

Migrations are placed in the [migrations/](file:///c:/Reza/TEL-U/Code/mentor_search/event-management-be/migrations) directory as `.up.sql` files:
- `000001_create_events_table.up.sql`: Creates `events` table with primary key `id BIGSERIAL`, category `CHECK` constraint (`workshop`, `competition`, `seminar`, `community`), and `TIMESTAMPTZ` columns.

### Migration Application
- **Postgres Initialization**: Migrations are mounted to `/docker-entrypoint-initdb.d` in Docker Compose, ensuring automatic execution for a freshly initialized Postgres container.
- **Application Startup**: The backend service also runs `database.RunMigrations()` on startup, ensuring that existing databases get updated migrations without requiring manual table creation in DBeaver.

---

## 7. Testing the Health Check Endpoint

### `GET /health`

Verifies that the HTTP service is operational and checks database connectivity using `pgxpool.Ping()`.

```bash
curl -X GET http://localhost:8080/health
```

#### Response (`200 OK`)
```json
{
  "success": true,
  "message": "application is healthy",
  "data": {
    "database": "connected",
    "status": "ok"
  }
}
```

---

## 8. Testing Event Creation Endpoint

### `POST /api/v1/events`

Creates a new event record.

Supported categories: `workshop`, `competition`, `seminar`, `community`.

#### Sample Request

```bash
curl -X POST http://localhost:8080/api/v1/events \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Go Backend Workshop",
    "description": "An introductory workshop about backend development.",
    "event_date": "2026-11-15T09:00:00+07:00",
    "category": "workshop"
  }'
```

#### Response (`201 Created`)
```json
{
  "success": true,
  "message": "event created successfully",
  "data": {
    "id": 1,
    "title": "Go Backend Workshop",
    "description": "An introductory workshop about backend development.",
    "event_date": "2026-11-15T02:00:00Z",
    "category": "workshop",
    "created_at": "2026-10-23T12:00:00Z",
    "updated_at": "2026-10-23T12:00:00Z"
  }
}
```

#### Error Examples

1. **Missing Title (`400 Bad Request`)**:
   ```json
   {
     "success": false,
     "message": "validation failed",
     "error": "title is required"
   }
   ```

2. **Invalid Category (`400 Bad Request`)**:
   ```json
   {
     "success": false,
     "message": "validation failed",
     "error": "invalid category: must be workshop, competition, seminar, or community"
   }
   ```

---

## 9. Connecting DBeaver to Local PostgreSQL

When running via Docker Compose, port `5432` is exposed on your host:

1. Open **DBeaver** and click **New Database Connection**.
2. Select **PostgreSQL**.
3. Fill in the connection settings:
   - **Host**: `localhost`
   - **Port**: `5432`
   - **Database**: `event_management`
   - **Username**: `postgres`
   - **Password**: `postgres`
4. Click **Test Connection** and then **Finish**.
5. Navigate to `Schemas -> public -> Tables -> events` to inspect the created events.

---

## 10. Running Checks and Tests Locally

To run tests and code checks without Docker:

```bash
# Format code
gofmt -w .

# Static analysis
go vet ./...

# Run unit tests
go test -v ./...

# Build binary
go build ./...
```

---

## 11. Troubleshooting

- **Database Connection Refused**:
  Ensure the `postgres` container is healthy before the API starts (`docker compose ps`). When running the API directly on the host (outside Docker), set `DB_HOST=localhost` in your `.env`.
- **Port Conflict on 5432 or 8080**:
  If a local PostgreSQL instance is already running on port `5432`, change `DB_PORT` in `.env` (e.g., `DB_PORT=5433`).
