# Boilerplate Go API

A production-ready Go REST API boilerplate with JWT authentication, Swagger documentation, and Docker support.

## Features

- **Clean Architecture**: Handler → Service → Repository layers
- **JWT Authentication**: HS256 tokens with configurable expiration
- **Swagger Documentation**: Auto-generated API docs at `/swagger`
- **Docker Ready**: Multi-stage builds with Alpine and Distroless variants
- **Graceful Shutdown**: Proper signal handling and connection draining
- **Security**: Rate limiting, secure headers, request size limits
- **Hot Reload**: Air configuration for development

## Quick Start

### Prerequisites

- Go 1.22+
- Docker (optional)

### Run Locally

```bash
# Clone the repository
git clone https://github.com/muflihunaf/boilerplate-go.git
cd boilerplate-go

# Copy environment file
cp env.example .env

# Run the application
make run
```

The API will be available at `http://localhost:8080`.

### Run with Docker

```bash
# Build and run
make docker-build
make docker-run
```

## API Documentation

Swagger UI is available at: `http://localhost:8080/swagger/index.html`

## Project Structure

```
.
├── cmd/api/            # Application entry point
├── internal/
│   ├── app/            # Application bootstrap
│   ├── config/         # Configuration loading
│   ├── handler/        # HTTP handlers
│   ├── middleware/     # Custom middleware
│   ├── repository/     # Data access layer
│   ├── server/         # HTTP server setup
│   └── service/        # Business logic
├── pkg/
│   ├── jwt/            # JWT token service
│   ├── response/       # Standard API responses
│   └── validator/      # Input validation
├── docs/               # Generated Swagger docs
├── Dockerfile          # Multi-stage Docker build
├── Makefile            # Build automation
└── env.example         # Environment template
```

## Configuration

All configuration is via environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `APP_ENV` | Environment (development/production) | `development` |
| `PORT` | HTTP server port | `8080` |
| `LOG_LEVEL` | Log level (debug/info/warn/error) | `info` |
| `JWT_SECRET` | JWT signing secret (required in production) | - |
| `JWT_EXPIRATION` | Token expiration in seconds | `86400` |
| `JWT_ISSUER` | Token issuer | `boilerplate-go` |
| `READ_TIMEOUT` | HTTP read timeout (seconds) | `15` |
| `WRITE_TIMEOUT` | HTTP write timeout (seconds) | `15` |
| `IDLE_TIMEOUT` | HTTP idle timeout (seconds) | `60` |

> ⚠️ In production, `JWT_SECRET` must be set and be at least 32 characters.

## Available Commands

```bash
make help           # Show all available commands
make run            # Run the application
make dev            # Run with hot reload (requires air)
make build          # Build binary
make test           # Run tests
make lint           # Run linter
make swagger        # Generate Swagger docs
make docker-build   # Build Docker image
```

## API Endpoints

### Public Routes

| Method | Path | Description |
|--------|------|-------------|
| GET | `/health` | Health check |
| GET | `/ready` | Readiness check |
| GET | `/swagger/*` | Swagger documentation |
| POST | `/api/v1/auth/login` | User login |
| POST | `/api/v1/auth/register` | User registration |

### Protected Routes (require JWT)

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/v1/me` | Get current user |
| GET | `/api/v1/users` | List all users |
| POST | `/api/v1/users` | Create user |
| GET | `/api/v1/users/{id}` | Get user by ID |
| PUT | `/api/v1/users/{id}` | Update user |
| DELETE | `/api/v1/users/{id}` | Delete user |

## Authentication

Include the JWT token in the Authorization header:

```
Authorization: Bearer <your-jwt-token>
```

## Response Format

All responses follow this format:

```json
{
  "success": true,
  "data": { ... },
  "error": null,
  "meta": null
}
```

Error responses:

```json
{
  "success": false,
  "data": null,
  "error": {
    "code": "UNAUTHORIZED",
    "message": "invalid email or password"
  }
}
```

## Development

### Hot Reload

```bash
# Install air
go install github.com/cosmtrek/air@latest

# Run with hot reload
make dev
```

### Generate Swagger Docs

```bash
make swagger
```

## Deployment

### Docker (Recommended)

```bash
# Production build (Alpine)
make docker-build

# Minimal build (Distroless)
make docker-build-distroless
```

### Environment Variables

For production, ensure these are set:

```bash
APP_ENV=production
JWT_SECRET=your-very-long-secret-key-at-least-32-chars
```

## License

MIT License - see [LICENSE](LICENSE) for details.

