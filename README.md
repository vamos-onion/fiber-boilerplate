# Fiber Boilerplate

A production-ready Go web service boilerplate built with Fiber framework, featuring OpenAPI-based code generation, PostgreSQL integration, and a clean, modular architecture.

## Features

- **High-Performance Web Framework**: Built on [Fiber](https://gofiber.io/) (powered by fasthttp)
- **OpenAPI-First Development**: API definitions drive code generation with `oapi-codegen`
- **Type-Safe Database Access**: SQL code generation with [sqlc](https://sqlc.dev/)
- **JWT Authentication**: Secure token-based authentication with automatic validation
- **Session Management**: Two-tier caching strategy (in-memory + database)
- **Graceful Shutdown**: Proper cleanup and connection handling
- **Structured Logging**: Using zerolog for performance and clarity
- **CORS Support**: Configurable cross-origin resource sharing
- **Auto-generated API Validation**: Request/response validation via OpenAPI middleware

## Technology Stack

### Core
- **Go**: 1.24.4
- **Web Framework**: [Fiber v2](https://github.com/gofiber/fiber) - Express-inspired web framework
- **Database**: PostgreSQL with [sqlx](https://github.com/jmoiron/sqlx)
- **Cache**: Ristretto (high-performance) / go-cache (simple)
- **Logging**: [zerolog](https://github.com/rs/zerolog) - Zero allocation JSON logger

### Authentication & Security
- **JWT**: [golang-jwt/jwt](https://github.com/golang-jwt/jwt) v5
- **Password/Token Storage**: Secure environment variable management

### Code Generation
- **OpenAPI Generator**: [oapi-codegen](https://github.com/oapi-codegen/oapi-codegen) v2
- **SQL Code Generator**: [sqlc](https://sqlc.dev/)

### Utilities
- **Environment Variables**: [caarlos0/env](https://github.com/caarlos0/env)
- **Dotenv**: [joho/godotenv](https://github.com/joho/godotenv)
- **Null Types**: [guregu/null](https://github.com/guregu/null)

## Prerequisites

Before you begin, ensure you have the following installed:

### Required
- **Go**: 1.24.4 or higher
- **PostgreSQL**: 12 or higher
- **Node.js & npm**: For OpenAPI tools (swagger-cli)
- **Make**: For build automation

### Development Tools
- **sqlc**: SQL code generator
  ```bash
  go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
  ```

- **oapi-codegen**: OpenAPI code generator
  ```bash
  go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
  ```

- **swagger-cli**: OpenAPI validation and bundling
  ```bash
  npm install -g @apidevtools/swagger-cli
  ```

## Installation

### 1. Clone the Repository
```bash
git clone <repository-url>
cd fiber-boilerplate
```

### 2. Install Go Dependencies
```bash
go mod download
```

### 3. Set Up Database

#### Create Database
```bash
createdb playground
```

#### Run Schema Migration
```bash
psql -U postgres -d playground -f database/V0__init.sql
```

Or connect to your PostgreSQL instance and run:
```sql
\i database/V0__init.sql
```

### 4. Configure Environment Variables

Copy the example environment file:
```bash
cp .env.example .env
```

Edit `.env` with your configuration:
```bash
# Application Configuration
PORT=8080
ENV=local  # Options: local, development, production

# PostgreSQL Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password_here
DB_NAME=playground
DB_MAX_IDLE_CONNS=5
DB_MAX_OPEN_CONNS=100
DB_MAX_LIFETIME=300

# Redis Configuration (Optional)
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

# JWT Configuration (REQUIRED)
JWT_SECRET=your_secure_random_secret_here

# CORS Configuration (Optional)
CORS_ENABLED=true
CORS_ALLOW_ORIGINS=*
CORS_ALLOW_METHODS=GET,HEAD,PUT,PATCH,POST,DELETE,OPTIONS

# Graceful Shutdown (Optional)
GRACEFUL_TIMEOUT=10  # seconds
LOG_REQUESTS_ENABLED=true
```

**Important**: Change `JWT_SECRET` to a secure random string in production!

### 5. Generate Code

Generate OpenAPI and sqlc code:
```bash
make openapi  # Validates, bundles, and generates OpenAPI code
make sqlc     # Generates database models from SQL
```

### 6. Build the Application
```bash
make build
```

The compiled binary will be in `out/fiber-boilerplate`.

## Usage

### Running the Application

#### Development Mode
```bash
make run-main
```

Or run the binary directly:
```bash
./out/fiber-boilerplate
```

#### Production Mode
Set `ENV=production` in your `.env` file or environment:
```bash
ENV=production ./out/fiber-boilerplate
```

### Testing the API

Once the server is running (default: `http://localhost:8080`):

#### Health Check
```bash
curl http://localhost:8080/api/ping
```

Expected response:
```json
{
  "code": 200,
  "message": "OK",
  "data": {
    "message": "pong"
  }
}
```

#### Create User (Requires JWT)
```bash
curl -X POST http://localhost:8080/api/appuser/create \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "John Doe",
    "birthday": "1990-01-01",
    "gender": "M"
  }'
```

#### List Users (Requires JWT)
```bash
curl -X GET "http://localhost:8080/api/appuser/list?limit=10&page=1" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Generating JWT Tokens

For development, you have two options:

#### Option 1: Use the Provided Utility
Generate a JWT token with the included utility (valid for 7 days):
```bash
go run generate_jwt.go
```

This generates a token with:
- `uuid`: Sample user UUID (12956e54-503d-46f1-8b9b-7cf304fba601)
- `exp`: 7 days from now (auto-calculated)
- `iat`: Current timestamp (auto-calculated)

#### Option 2: Create Tokens Manually
Use online tools like [jwt.io](https://jwt.io/) or create tokens programmatically.

JWT Payload Format:
```json
{
  "uuid": "your-user-uuid-here",
  "exp": 1735689600,
  "iat": 1735344000
}
```

**Required fields:**
- `uuid` (string): User UUID
- `exp` (integer): Unix timestamp when token expires

**Optional fields:**
- `iat` (integer): Unix timestamp when token was issued

**Signing:**
- Algorithm: HS256
- Secret: Your `JWT_SECRET` value from `.env`

## Development Workflow

### 1. Modifying the API

#### Edit OpenAPI Specifications
```bash
# Main entry point
vim api/index.yaml

# Schemas
vim api/schemas.yaml

# Endpoints
vim api/v1/create_appuser.yaml
```

#### Regenerate Code
```bash
make openapi
```

This will:
1. Validate OpenAPI specs
2. Bundle into `out/openapi.yaml`
3. Generate Go code in `internal/generated/serviceapi/`

### 2. Modifying Database Schema

#### Edit SQL Schema
```bash
vim database/V0__init.sql
```

#### Edit SQL Queries
```bash
vim database/queries/appuser.sql
```

#### Regenerate Models
```bash
make sqlc
```

### 3. Building and Testing

```bash
# Run tests
make test

# Build
make build

# Run all (test + build)
make all

# Clean build artifacts
make clean
```

## API Endpoints

All API endpoints are prefixed with `/api` and defined in OpenAPI specifications.

### Public Endpoints
- `GET /api/ping` - Health check (no authentication required)

### Protected Endpoints (Requires JWT)
- `POST /api/appuser/create` - Create a new user
- `GET /api/appuser/list` - List users with pagination and filtering
- `PUT /api/appuser/update` - Update an existing user

### Query Parameters for List

- `uuid` (string) - Filter by user UUID
- `name` (string) - Filter by name
- `gender` (string) - Filter by gender (M/F)
- `withdraw` (boolean) - Filter by withdrawal status
- `limit` (integer) - Items per page (max: 1000)
- `page` (integer) - Page number (1-based)
- `sorting[key]` (string) - Sort field
- `sorting[dir]` (string) - Sort direction (asc/desc)

## Project Structure

```
fiber-boilerplate/
├── cmd/
│   └── main.go                 # Application entry point
├── internal/
│   ├── app/
│   │   ├── app.go              # Application startup and shutdown
│   │   ├── config/             # Configuration management
│   │   ├── handlers/           # HTTP request handlers
│   │   │   └── v1/             # API v1 handlers
│   │   ├── middleware/         # Fiber middleware
│   │   └── router/             # Route definitions
│   ├── pkg/
│   │   ├── cache/              # Cache implementations
│   │   ├── database/           # Database drivers
│   │   ├── logging/            # Logging utilities
│   │   ├── session/            # Session management
│   │   ├── setting/            # Runtime settings
│   │   └── util/               # Utility functions
│   ├── models/                 # Generated database models (sqlc)
│   ├── generated/              # Generated OpenAPI code
│   └── defs/                   # Error definitions
├── api/                        # OpenAPI specifications
│   ├── index.yaml              # Main spec entry point
│   ├── schemas.yaml            # Data schemas
│   ├── parameters.yaml         # Reusable parameters
│   └── v1/                     # API v1 endpoint specs
├── database/
│   ├── V0__init.sql            # Database schema
│   └── queries/                # SQL query definitions
│       └── appuser.sql
├── sqlc_conf/
│   ├── sqlc.yaml               # sqlc configuration
│   └── overrides.yaml          # Field name overrides
├── build/                      # Makefile modules
│   ├── common.mk
│   ├── go.mk
│   └── openapi.mk
├── out/                        # Build artifacts (gitignored)
├── .env.example                # Environment template
├── go.mod                      # Go dependencies
├── Makefile                    # Build automation
├── CLAUDE.md                   # AI assistant guidelines
├── LICENSE                     # MIT License
└── README.md                   # This file
```

## Middleware Pipeline

Requests flow through middleware in this order:

1. **Request ID** - Unique request tracking
2. **Logger** - Request/response logging
3. **Recover** - Panic recovery
4. **ETag** - HTTP caching
5. **Compress** - Response compression
6. **Pprof** - Profiling (local/dev only)
7. **CORS** - Cross-origin resource sharing
8. **OpenAPI Validation** - Request validation and auth requirement detection
9. **JWT Authentication (keyauth)** - Bearer token extraction and validation
10. **Session** - User session loading from database

## Database Schema Conventions

All tables follow a standard pattern:

```sql
CREATE TABLE example (
    id          bigserial   PRIMARY KEY,
    uuid        uuid        NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    created_at  timestamptz NOT NULL DEFAULT now(),
    modified_at timestamptz NOT NULL DEFAULT now(),
    -- domain-specific fields --
);

-- Auto-update modified_at on every UPDATE
CREATE TRIGGER tr_example_update_modified_at
    BEFORE UPDATE ON example
    FOR EACH ROW EXECUTE PROCEDURE fn_set_modified_at();
```

## Query Patterns

### Direct Query
```go
// In Fiber handlers, use ctx.Context() to convert *fiber.Ctx to context.Context
result, err := models.Appuser.CreateAppuser(nil, ctx.Context(), params)
```

### Transaction Query
```go
// Get context.Context from Fiber context
tx, qctx, err := models.SQL.BeginxContext(ctx.Context())
if err != nil {
    return err
}
defer tx.Rollback()

qtx := models.New(tx)
result, err := models.Appuser.UpdateAppuser(qtx, qctx, params)
if err != nil {
    return err
}

return tx.Commit()
```

## Configuration Reference

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | 8080 | HTTP server port |
| `ENV` | local | Environment (local/development/production) |
| `JWT_SECRET` | **required** | JWT signing secret |
| `GRACEFUL_TIMEOUT` | 10 | Graceful shutdown timeout (seconds) |
| `LOG_REQUESTS_ENABLED` | true | Enable request logging |
| `DB_HOST` | localhost | PostgreSQL host |
| `DB_PORT` | 5432 | PostgreSQL port |
| `DB_USER` | postgres | Database user |
| `DB_PASSWORD` | "" | Database password |
| `DB_NAME` | playground | Database name |
| `DB_MAX_IDLE_CONNS` | 5 | Max idle connections |
| `DB_MAX_OPEN_CONNS` | 100 | Max open connections |
| `DB_MAX_LIFETIME` | 300 | Connection lifetime (seconds) |
| `REDIS_HOST` | localhost | Redis host |
| `REDIS_PORT` | 6379 | Redis port |
| `REDIS_PASSWORD` | "" | Redis password |
| `CORS_ENABLED` | true | Enable CORS |
| `CORS_ALLOW_ORIGINS` | * | Allowed origins (comma-separated) |
| `CORS_ALLOW_METHODS` | GET,HEAD,PUT,... | Allowed HTTP methods |
| `CORS_ALLOW_HEADERS` | "" | Allowed headers |
| `CORS_ALLOW_CREDENTIALS` | false | Allow credentials |
| `CORS_EXPOSE_HEADERS` | "" | Exposed headers |
| `CORS_MAX_AGE` | 0 | Preflight cache duration (seconds) |

## Build Commands

```bash
make all                # Run tests and build
make build              # Compile application
make run-main           # Run the application
make test               # Run tests
make tidy               # Tidy Go modules
make clean              # Remove build artifacts
make openapi            # OpenAPI workflow (validate + bundle + generate)
make openapi-validate   # Validate OpenAPI specs
make openapi-bundle     # Bundle OpenAPI specs
make openapi-generate   # Generate Go code from OpenAPI
make sqlc               # Generate database models
```

## Troubleshooting

### Database Connection Issues
- Verify PostgreSQL is running: `pg_isready`
- Check connection settings in `.env`
- Ensure database exists: `psql -l | grep playground`

### OpenAPI Generation Fails
- Install swagger-cli: `npm install -g @apidevtools/swagger-cli`
- Validate specs manually: `make openapi-validate`

### SQLC Generation Fails
- Install sqlc: `go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest`
- Check `database/queries/*.sql` syntax

### JWT Authentication Fails
- Ensure `JWT_SECRET` is set in `.env`
- Verify token format (Bearer scheme)
- Check token expiration (`exp` claim)

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Support

For issues and questions:
- Open an issue on GitHub
- Check existing documentation in `CLAUDE.md`
- Review OpenAPI specs in `api/` directory

## Acknowledgments

- [Fiber](https://gofiber.io/) - Web framework
- [sqlc](https://sqlc.dev/) - SQL code generator
- [oapi-codegen](https://github.com/oapi-codegen/oapi-codegen) - OpenAPI code generator
- [zerolog](https://github.com/rs/zerolog) - Logging library
