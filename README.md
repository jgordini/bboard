# bboard

A self-hosted feedback portal for collecting feature requests and suggestions.

## Tech Stack

- **Backend:** Go 1.22+
- **Frontend:** React 18, TypeScript
- **Database:** PostgreSQL 17
- **Styling:** SCSS with BEM conventions and utility classes

## Prerequisites

- Go 1.22+
- Node.js 21/22
- Docker
- `air` - Go live reload (`go install github.com/air-verse/air@latest`)
- `godotenv` - Environment loader (`go install github.com/joho/godotenv/cmd/godotenv@latest`)
- `golangci-lint` - Go linter

## Setup

1. Start services:
   ```bash
   docker compose up -d
   ```

2. Copy environment config:
   ```bash
   cp .example.env .env
   ```

3. Run migrations:
   ```bash
   make migrate
   ```

4. Start development server:
   ```bash
   make watch
   ```

5. Open http://localhost:3000

## Development Commands

| Command | Description |
|---------|-------------|
| `make watch` | Hot reload for server and UI |
| `make build` | Build production binaries and assets |
| `make test` | Run all tests (Go + Jest) |
| `make lint` | Lint server and UI code |
| `make migrate` | Run database migrations |

## Local Services

| Service | URL/Port |
|---------|----------|
| App | http://localhost:3000 |
| PostgreSQL (dev) | localhost:5555 |
| PostgreSQL (test) | localhost:5566 |
| MailHog UI | http://localhost:8025 |
| MailHog SMTP | localhost:1025 |

## Project Structure

```
app/
├── handlers/     # HTTP request handlers
├── models/       # Data models (entity, cmd, query, dto)
├── services/     # Business logic
├── pkg/bus/      # Service registry and dispatch
└── cmd/routes.go # All HTTP routes

public/
├── pages/        # Page components (lazy-loaded)
├── components/   # Reusable UI components
├── services/     # API clients
└── assets/styles # SCSS styles

migrations/       # Database migrations (SQL)
```

## Running Tests

```bash
# All tests
make test

# Single Go test
godotenv -f .test.env go test ./app/handlers -v -run TestName

# Single Jest test
npx jest ./public/path/to/file.spec.tsx

# E2E tests
make test-e2e-ui
```
