# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build and Development Commands

### Go Application
**Go version**: 1.26.0 (required, specified in go.mod)

- **Build the project**: `go build ./...`
- **Run the HTTP server**: `go run cmd/server/main.go` (requires `config.yaml` in working directory)
- **Run the Telegram bot**: `go run cmd/bot/main.go` (requires `config.yaml` in working directory)
- **Lint**: `golangci-lint run` (uses `.golangci.yml` version 2 format)
- **Run tests**: `go test ./...`
- **Run single test**: `go test -run <TestName> ./path/to/package`

### Python Text Classifier Service
Located in `text_classifier/` - a FastAPI service for ML-based expense category classification.
- **Makefile location**: Root directory (`Makefile`)
- **Install dependencies**: `make install` (creates venv and installs requirements)
- **Run the service**: `make run` (runs on port 8000 with hot reload)
- **Clean venv**: `make clean`

## Architecture Overview

### Module Structure
The Go module name is `main` (not a domain-based name). All imports use `"main/pkg/..."` paths.

### Entry Points
Two separate applications share the same codebase:

1. **cmd/server** - HTTP REST API server
   - Serves `/health`, `/api/v1/transactions/`, `/api/v1/prefilled-expenses`
   - Uses CORS middleware for local development
   - Port 8081 (hardcoded)

2. **cmd/bot** - Telegram Bot
   - Interactive chat-based expense entry via Telegram Bot API
   - Commands: `/add` (start expense entry), `/summary` (view monthly summary)
   - Session-based Q&A flow defined in `pkg/session/`

### Core Packages

**pkg/config/**
Configuration structs with YAML tags. Key files:
- `config.go` - Root `Config` struct
- `database_config.go` - `DatabaseConfig` and `FeaturesConfig` (note: `save_to_database` YAML key)
- `frequent_expenses_config.go` - `FrequentExpense` for pre-filled expenses
- `bot_config.go` - `TelegramConfig`

**Important**: The `config.yaml` file must be present in the working directory when running either application entry point.

**pkg/storage/**
Database abstraction layer:
- `db.go` - PostgreSQL connection initialization with auto-migration
- `read_db.go` - Query functions for summaries (group by category, claimable, paid_for_family)
- `storage.go` - File-based fallback storage when DB is disabled

Database auto-creates `transactions` table on startup via `createTableIfNotExists()`.

**pkg/session/**
Bot conversation state management:
- `constants.go` - Question order constants (using iota) and question text array
- `session.go` - `UserSession` struct tracking current question and answers
- `handle_answer.go` - Answer validation and state progression
- `default_response.go` - Fallback message handling

Question flow (defined in `constants.go`):
1. QuestionName - "What is the name of the transaction?"
2. QuestionAmount - "How much is the transaction?"
3. QuestionCurrency - "What currency is the transaction in?"
4. QuestionDate - "What is the date of transaction?"
5. QuestionIsClaimable - "Is it claimable?"
6. QuestionPaidForFamily - "Is it paid for the family?"
7. QuestionCategory - "What is the category of transaction?"

Pre-filled expenses can auto-populate Category, PaidForFamily, and Currency, skipping corresponding questions.

**pkg/transaction/**
Core domain model:
- `Transaction` struct - maps to both DB (sql tags) and JSON API (json tags)
- `transaction_categories.go` - Category enum definitions
- `paginated_transaction.go` - Pagination support for API responses

**pkg/bot/**
Telegram bot implementation:
- `Bot` struct manages bot lifecycle and message handlers
- Handles both text messages and inline keyboard callbacks
- Auto-deletes messages for clean chat experience
- Pre-filled expense logic: auto-populates fields based on selected frequent expense

**pkg/handler/**
HTTP handlers for the server:
- `handler.go` - Main transaction CRUD handler
- `transaction.go` - Transaction-specific HTTP handlers
- `prefilled_expenses.go` - Frequent expenses endpoint
- `health_check.go` - Health endpoint

### Key Features

**Pre-filled Expenses**
Defined in `config.yaml` under `frequent_expenses`. When selected:
- Auto-populates `Category` if set
- Auto-populates `PaidForFamily` if set
- Auto-populates `Currency` if set
- Skips corresponding questions in the flow

**Database Toggle**
`FeaturesConfig.SaveToDB` (YAML: `save_to_database`) controls persistence:
- `true`: Saves to PostgreSQL
- `false`: Saves to `responses.txt` file

**Date Input**
Users can enter:
- `t` for today's date
- `DD.MM.YY` format (e.g., `15.04.25`)
- Invalid input defaults to today

## Documentation Requirements

When making changes to the following areas, update the corresponding documentation:

**API Changes**: Update Swagger documentation in `pkg/handler/` (if implemented) and ensure endpoints are documented in comments for future swagger generation.

**Bot Commands or Session Flow**: Update `README.md` - document new commands or changes to the Q&A flow in the Telegram bot section.

**Configuration Changes**: Update `README.md` - document any new config fields in the Configuration section and add to `config.example.yaml`.

**New Features**: Add to `ISSUE_BACKLOG.md` if incomplete, or update `CLAUDE.md` Architecture Overview section if fundamental to the structure.

**Pre-filled Expenses**: If changing the `FrequentExpense` struct or related logic, update both `pkg/config/frequent_expenses_config.go` comments and user-facing documentation.

## CI/CD

GitHub Actions workflow in `.github/workflows/lint.yml`:
- Runs on Go 1.26.0
- Builds with `go build -v ./...`
- Lints with golangci-lint v2.9.0

## External Dependencies

**Go:**
- `github.com/go-telegram-bot-api/telegram-bot-api/v5` - Telegram Bot API
- `github.com/lib/pq` - PostgreSQL driver
- `github.com/rs/cors` - CORS middleware
- `gopkg.in/yaml.v3` - YAML parsing

**Python:**
- FastAPI + uvicorn for ML service
- Hugging Face transformers with `mrm8488/bert-mini-finetuned-expense-category` model
