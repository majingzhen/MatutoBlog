# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

MatutoBlog is a Go-based blog system with the following architecture:
- **Backend**: Go + Gin web framework + GORM + MySQL
- **Frontend**: HTML templates with static assets
- **Authentication**: JWT tokens + session-based admin auth
- **Database**: MySQL with GORM models using custom table prefixes (`m_` for models, `p_` for legacy tables)

## Development Commands

### Build and Run
```bash
# Run in development mode
make run
# or
go run main.go

# Build production binary
make build

# Build development version with debug info
make build-dev

# Run built development version
make run-dev
```

### Testing
```bash
# Run all tests
make test

# Run with coverage report
make test-coverage

# Run unit tests only
make test-unit

# Run integration tests only
make test-integration
```

### Code Quality
```bash
# Format code
make fmt

# Run linter (requires golangci-lint)
make lint

# Install dependencies
make deps

# Update dependencies
make deps-update
```

### Documentation
```bash
# Generate Swagger API docs (requires swag)
make swagger
# Access at: http://localhost:8080/swagger/index.html
```

### Database
```bash
# Run database migrations
make migrate
```

### Docker
```bash
# Build Docker image
make docker-build

# Run Docker container
make docker-run
```

## Architecture Overview

### Key Components
- **Entry Point**: `main.go` - Initializes config, logger, database, storage, and routes
- **Configuration**: `config/` - Viper-based config with YAML + env vars
- **API Layer**: `internal/api/` - Controllers, middleware, and routing
- **Models**: `internal/models/` - GORM models with custom table names
- **Database**: `internal/database/` - MySQL connection and table initialization
- **Utilities**: `pkg/` - Common utilities, pagination, storage abstraction

### Database Schema
- Uses custom table prefixes (`m_` for current models, `p_` for legacy)
- Main entities: Article, Category, Tag, Comment, Attachment
- Supports both HTML and Markdown content types
- Article-Category and Article-Tag many-to-many relationships

### Authentication Flow
- **Admin Routes** (`/admin/*`): Session-based authentication with middleware
- **API Routes** (`/api/*`): JWT-based authentication
- Frontend serves both public blog pages and admin management interface

### Configuration System
- Primary config: `config/config.yaml`
- Environment variable override support
- Sensible defaults for development in `config/config.go`
- Database, JWT, storage, and CORS settings

### Template System
- HTML templates in `web/templates/`
- Static assets in `web/static/`
- Admin interface templates in `web/admin/`
- Template loading via custom renderer

### Storage System
- Local file storage in `./uploads/` (configurable)
- Cloud storage abstraction ready (not yet implemented)

## Important Notes

### Running Tests
- No test files currently exist in the codebase
- Test structure planned for `test/unit/` and `test/integration/`

### Database Setup
- Requires MySQL database named `matuto_blog`
- Tables auto-created on startup via `internal/database/init_tables.go`
- Database schema available in `doc/blog_optimized.sql`

### Development Workflow
1. Start with `make run` for development
2. Use `make fmt` before commits
3. Run `make lint` for code quality (requires golangci-lint)
4. Generate API docs with `make swagger` (requires swag tool)

### Module Dependencies
- Main module: `matuto-blog`
- Key dependencies: Gin, GORM, Viper, JWT, Logrus
- Import paths use module prefix: `matuto-blog/internal/...`