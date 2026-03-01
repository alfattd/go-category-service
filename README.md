# Go Category Service

A production-ready RESTful microservice for managing categories, built with Go. Publishes domain events to RabbitMQ on every state change.

![Go](https://img.shields.io/badge/go-1.21+-00ADD8?logo=go)
![Docker](https://img.shields.io/badge/docker-ready-blue)
![License](https://img.shields.io/badge/license-MIT-green)

## Table of Contents

- [Architecture](#architecture)
- [Requirements](#requirements)
- [Quick Start](#quick-start)
- [Configuration](#configuration)
- [API Reference](#api-reference)
- [Running Tests](#running-tests)
- [Project Structure](#project-structure)

---

## Architecture

```
┌─────────────┐     HTTP      ┌─────────────────────────────────────────┐
│   Client    │ ────────────► │              Handler Layer              │
└─────────────┘               │         (command.go / query.go)         │
                              └───────────────────┬─────────────────────┘
                                                  │
                               ┌──────────────────▼──────────────────────┐
                               │              Service Layer              │
                               │         (command.go / query.go)         │
                               └──────────┬──────────────────┬───────────┘
                                          │                  │
                    ┌─────────────────────▼──┐   ┌───────────▼──────────────┐
                    │   Repository Layer     │   │     Event Publisher      │
                    │   (PostgreSQL / pq)    │   │     (RabbitMQ AMQP)      │
                    └────────────────────────┘   └──────────────────────────┘
```

The service follows **Clean Architecture** principles with a clear separation between domain, service, repository, and handler layers. All cross-layer contracts are defined as interfaces in the `domain` package, making each layer independently testable.

---

## Requirements

- Go 1.21+
- Docker & Docker Compose
- Docker network (shared with other services)
- [RabbitMQ](https://github.com/alfattd/rabbitmq) (run separately)

---

## Quick Start

### 1. Create Docker Network

```bash
docker network create net
```

### 2. Start RabbitMQ

Follow the setup guide in the [RabbitMQ repo](https://github.com/alfattd/rabbitmq), then come back here.

### 3. Clone & Configure

```bash
git clone https://github.com/alfattd/category-service.git
cd category-service
cp .env.example .env
```

Edit `.env` if needed (defaults work out of the box for local dev):

```env
APP_PORT=80
SERVICE_NAME=category-service
SERVICE_VERSION=dev

RABBITMQ_URL=amqp://rabbitmq:rabbitmq@rabbitmq:5672/

DB_HOST=postgres
DB_PORT=5432
DB_NAME=postgres
DB_USER=postgres
DB_PASSWORD=password
DB_SSLMODE=disable

NETWORK=net
```

### 4. Start Services

```bash
docker compose up -d --build
```

This starts PostgreSQL, runs migrations automatically, then starts the app. The service is now running at `http://localhost:80`.

---

## Configuration

| Variable | Description | Default |
|---|---|---|
| `APP_PORT` | HTTP server port | `80` |
| `SERVICE_NAME` | Service identifier (used in logs) | — |
| `SERVICE_VERSION` | Service version (used in `/version` endpoint) | — |
| `RABBITMQ_URL` | Full AMQP connection string | — |
| `DB_HOST` | PostgreSQL host | — |
| `DB_PORT` | PostgreSQL port | `5432` |
| `DB_NAME` | Database name | — |
| `DB_USER` | Database user | — |
| `DB_PASSWORD` | Database password | — |
| `DB_SSLMODE` | SSL mode (`disable` / `require`) | `disable` |
| `NETWORK` | Docker network name | `net` |

---

## API Reference

Base URL: `http://localhost`

### Health & Info

| Method | Path | Description |
|---|---|---|
| `GET` | `/health` | Liveness check |
| `GET` | `/version` | Service name & version |

### Categories

| Method | Path | Description |
|---|---|---|
| `GET` | `/categories` | List all categories |
| `POST` | `/categories` | Create a category |
| `GET` | `/categories/{id}` | Get category by ID |
| `PUT` | `/categories/{id}` | Update a category |
| `DELETE` | `/categories/{id}` | Delete a category |

#### Create Category

```bash
curl -X POST http://localhost/categories \
  -H "Content-Type: application/json" \
  -d '{"name": "Electronics"}'
```

```json
{
  "message": "category created",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "Electronics",
    "created_at": "2026-01-01T00:00:00Z",
    "updated_at": "2026-01-01T00:00:00Z"
  }
}
```

#### Error Responses

| HTTP Status | Meaning |
|---|---|
| `400 Bad Request` | Invalid input (e.g. empty name) |
| `404 Not Found` | Category not found |
| `409 Conflict` | Category name already exists |
| `500 Internal Server Error` | Unexpected server error |

---

## Running Tests

### Unit Tests

Tests for service and handler layers using mocks — no external dependencies needed.

```bash
make test-unit
```

### Integration Tests

Tests for the repository layer against a real PostgreSQL instance via [Testcontainers](https://testcontainers.com/). Requires Docker.

```bash
make test-integration
```

### Coverage Report

```bash
make test-coverage
```

Generates `app/coverage.html` with a full visual report.

---

## Project Structure

```
.
├── app/
│   ├── cmd/server/         # Entrypoint
│   ├── internal/
│   │   ├── config/         # App-level config (loads from env)
│   │   ├── domain/         # Domain models, interfaces, errors
│   │   ├── handler/        # HTTP handlers (command & query)
│   │   ├── mocks/          # Testify mocks for all interfaces
│   │   ├── pkg/
│   │   │   ├── config/     # Base config helpers
│   │   │   ├── database/   # PostgreSQL connection
│   │   │   ├── logger/     # slog-based structured logger
│   │   │   ├── middleware/ # RequestID, logging, recovery
│   │   │   ├── rabbitmq/   # AMQP publisher with retry & confirm mode
│   │   │   ├── requestid/  # Context-based request ID
│   │   │   └── system/     # Health & version endpoints
│   │   ├── repository/     # PostgreSQL repository implementation
│   │   ├── server/         # HTTP server setup & dependency wiring
│   │   └── service/        # Business logic (command & query)
│   ├── Dockerfile
│   └── go.mod
├── postgres/
│   └── migrations/         # SQL migration files
├── compose.yml
├── Makefile
└── .env.example
```

## RabbitMQ Events

On every state change, the service publishes an event to the `category_events` queue:

| Event Type | Trigger |
|---|---|
| `category_created` | `POST /categories` |
| `category_updated` | `PUT /categories/{id}` |
| `category_deleted` | `DELETE /categories/{id}` |

Event payload:
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "Electronics",
  "type": "category_created"
}
```

> **Note:** Publish failures are logged but do not fail the HTTP response. The service guarantees at-least-once delivery via broker confirm mode with exponential backoff retry (up to 3 attempts).
