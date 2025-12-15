# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

tsks is a terminal-based task management application built in Go with a Bubble Tea TUI frontend and SQLite persistence. Tasks use a versioned/soft-delete pattern where updates create new versions and deletes set `archived_at`.

## Common Commands

```bash
# Install dependencies
make install

# Run tests
make test

# Lint
make lint
make lint/fix

# Format
make fmt

# Run TUI with hot reload (requires air)
air

# Generate sqlc code after modifying queries
make sqlc/generate

# Database migrations
make migrate/up           # Apply migrations
make migrate/down/1       # Rollback last migration
make migrate/create NAME=migration_name  # Create new migration
```

## Architecture

### Entry Points
- `cmd/tui/main.go` - TUI application using Bubble Tea
- `cmd/db/main.go` - CLI for database operations (migrate, seed, list, truncate)

### Database Layer
- **sqlc** generates type-safe Go code from SQL queries
- Queries defined in `sqlc/queries/*.sql`
- Schema derived from `migrations/*.sql`
- Generated code lives in `internal/sqlc/`
- Custom type overrides (e.g., `TaskStatus`) configured in `sqlc.yml`

### Key Patterns
- **Versioned entities**: Tasks have a composite primary key `(id, version)`. Updates insert new rows with the same ID
- **Soft deletes**: `archived_at` timestamp marks deleted records; queries filter by `archived_at IS NULL`
- **Custom types**: `internal/schemas/tasks.go` defines `TaskStatus` enum with `driver.Valuer`/`sql.Scanner` implementations

### TUI Structure
- `internal/tui/model/model.go` - Main Bubble Tea model
- `internal/tui/model/task.go` - Task list item delegate
- Uses charmbracelet/bubbles list component
