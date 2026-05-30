# SilkWave API

A high-performance Go backend for SilkWave — a music platform for artists to sell and distribute their releases directly to fans.

## Tech Stack

| Component | Technology                                                             |
| --------- | ---------------------------------------------------------------------- |
| Language  | Go 1.25                                                                |
| Framework | [Gin](https://github.com/gin-gonic/gin)                                |
| Database  | [ArangoDB](https://www.arangodb.com/)                                  |
| Storage   | [Cloudflare R2](https://developers.cloudflare.com/r2/) (S3-compatible) |
| Auth      | JWT + bcrypt                                                           |
| Logging   | [Zap](https://github.com/uber-go/zap)                                  |

## Project Structure

```
go-silk-wave/
├── cmd/
│   ├── server/          # Main API server entrypoint
│   └── migrate/         # Database migration CLI
├── internal/
│   ├── auth/            # JWT & password services
│   ├── config/          # Environment configuration
│   ├── database/        # ArangoDB client, migrations, schema
│   ├── handlers/        # HTTP request handlers
│   ├── logger/          # Zap logger setup
│   ├── middleware/      # Auth middleware
│   ├── models/          # Data models & DTOs
│   ├── repository/      # Data access layer
│   ├── routes/          # Route definitions
│   └── storage/         # R2/S3 storage client
└── tmp/                 # Build artifacts (gitignored)
```

## Getting Started

### Prerequisites

- Go 1.25+
- ArangoDB 3.x
- Cloudflare R2 bucket (or S3-compatible storage)
- Docker / orbstack

### Environment Variables

Create a `.env` file in the project root:

```env
# Server
GO_ENV=development
SERVER_PORT=8080

# ArangoDB
ARANGO_ENDPOINT=http://localhost:8529
ARANGO_DATABASE=silk_wave
ARANGO_USERNAME=root
ARANGO_PASSWORD=your_password

# Security
JWT_SECRET=your-jwt-secret-min-32-chars
PASSWORD_SECRET=your-password-pepper-secret

# Cloudflare R2 Storage
R2_ACCOUNT_ID=your_account_id
R2_ACCESS_KEY_ID=your_access_key
R2_SECRET_ACCESS_KEY=your_secret_key
R2_BUCKET_NAME=silkwave
```

### Running the Application

```bash
# Install dependencies
go mod download

# Run database migrations
go run cmd/migrate/main.go

# Preview migrations (dry run)
go run cmd/migrate/main.go --dry-run

# Start the server
go run cmd/server/main.go
```

### Building

```bash
# Build binaries
go build -o bin/server ./cmd/server
go build -o bin/migrate ./cmd/migrate

# Run
./bin/migrate
./bin/server
```

## API Endpoints

### Public Routes

| Method | Endpoint                   | Description             |
| ------ | -------------------------- | ----------------------- |
| `GET`  | `/health`                  | Health check            |
| `POST` | `/api/login`               | User login              |
| `POST` | `/api/register`            | User registration       |
| `POST` | `/api/logout`              | User logout             |
| `GET`  | `/api/releases`            | List published releases |
| `GET`  | `/api/artists/:artistSlug` | Get artist by slug      |
| `GET`  | `/api/artists/:artistSlug/releases/:releaseSlug` | Get release by Artist and Release slug |

### Protected Routes (require JWT)

| Method | Endpoint                       | Description          |
| ------ | ------------------------------ | -------------------- |
| `GET`  | `/api/user`                    | Get current user     |
| `POST` | `/api/user/mode`               | Update user mode     |
| `POST` | `/api/settings/email`          | Update email         |
| `POST` | `/api/settings/password`       | Update password      |
| `POST` | `/api/settings/delete-account` | Delete account       |
| `POST` | `/api/register/artist-name`    | Register artist name |

### Release Management (protected)

| Method   | Endpoint                               | Description          |
| -------- | -------------------------------------- | -------------------- |
| `POST`   | `/api/releases/draft`                  | Create draft release |
| `POST`   | `/api/releases/:releaseId/confirm`     | Confirm uploads      |
| `POST`   | `/api/releases/:releaseId/publish`     | Publish release      |
| `POST`   | `/api/releases/:releaseId/archive`     | Archive release      |
| `DELETE` | `/api/releases/:releaseId`             | Delete release       |
| `POST`   | `/api/upload/avatar`                   | Create avatar upload URL |

## Database Schema

Collections managed by the migration system:

- **Users** — User accounts
- **Artists** — Artist profiles
- **UsersArtists** — Edge collection linking users to artists
- **Releases** — Releases across the Draft, Published, and Archived lifecycle states
- **Tracks** — Tracks belonging to a Release
- **Subscriptions** — Artist paid support levels
- **Subscribers** — User subscription relationships

## Storage Structure

Content is stored in R2 with the following structure:

```
artist_content/{artistKey}/releases/{releaseId}/
├── draft/                    # Unpublished content
│   ├── cover.jpg
│   ├── wavs/
│   └── mp3s/
├── cover.jpg                 # Published content
├── wavs/
└── mp3s/

avatars/{userKey}/
└── avatar.png
```

## Development

### Live Reload

Using [air](https://github.com/cosmtrek/air) for hot reloading:

```bash
# Install air
go install github.com/cosmtrek/air@latest

# Run with hot reload
air
```

### Type Generation

TypeScript types are generated using [tygo](https://github.com/gzuidhof/tygo):

```bash
tygo generate
```

## License

Proprietary — All rights reserved.
