# MTG Card Detector - Backend API

Go backend service for the Magic: The Gathering Card Detector mobile application.

## Features

- Anonymous user authentication with JWT tokens
- Card scanning and identification via Scryfall API
- Inventory management system
- Bulk scanning support
- SQLite database with automatic migrations
- RESTful API with proper error handling and logging

## Prerequisites

- Go 1.21 or higher
- SQLite 3

## Quick Start

### 1. Install Dependencies

```bash
cd backend
go mod download
```

### 2. Build

```bash
go build -o server ./cmd/server
```

### 3. Run

```bash
./server
```

The server will start on port 8080 by default. You can configure this with the `PORT` environment variable.

## Configuration

Environment variables:

- `PORT` - Server port (default: 8080)
- `DATABASE_PATH` - Path to SQLite database file (default: ./data/mtg_cards.db)
- `JWT_SECRET` - Secret key for JWT signing (change in production!)
- `MIGRATIONS_PATH` - Path to migration files (default: ./migrations)

Example:

```bash
export PORT=3000
export JWT_SECRET="your-secure-random-secret-here"
./server
```

## API Endpoints

### Public Endpoints

#### Health Check
```
GET /health
```

#### Anonymous Authentication
```
POST /api/v1/auth/anonymous
Content-Type: application/json

{
  "device_id": "unique-device-identifier"
}

Response:
{
  "user_id": "uuid",
  "token": "jwt-token"
}
```

### Protected Endpoints

All protected endpoints require an `Authorization: Bearer <token>` header.

#### Single Card Scan
```
POST /api/v1/cards/scan
Authorization: Bearer <token>
Content-Type: application/json

{
  "card_name": "Lightning Bolt",
  "set_code": "LEA",
  "collector_number": "161"
}

Response:
{
  "success": true,
  "card": {
    "id": "uuid",
    "name": "Lightning Bolt",
    "set_code": "LEA",
    "collector_number": "161",
    "image_uri": "https://...",
    ...
  }
}
```

#### Bulk Card Scan
```
POST /api/v1/cards/scan/bulk
Authorization: Bearer <token>
Content-Type: application/json

{
  "scans": [
    {"card_name": "Black Lotus", "set_code": "LEA"},
    {"set_code": "M21", "collector_number": "123"}
  ]
}

Response:
{
  "session_id": 1,
  "total_scanned": 2,
  "successful_scans": 2,
  "failed_scans": 0,
  "results": [...]
}
```

#### Get Inventory
```
GET /api/v1/inventory
Authorization: Bearer <token>

Response:
{
  "inventory": [
    {
      "id": 1,
      "user_id": "uuid",
      "card_id": "uuid",
      "quantity": 3,
      "added_at": "2025-11-15T...",
      "card": {...}
    }
  ],
  "count": 1
}
```

#### Get Card Details
```
GET /api/v1/cards?id=<card_id>
Authorization: Bearer <token>
```

## Testing

Run tests:

```bash
cd backend
go test ./...
```

Run tests with coverage:

```bash
go test -cover ./...
```

## Database Schema

The application uses SQLite with the following tables:

- **users** - Anonymous user accounts
- **cards** - MTG card master data (cached from Scryfall)
- **inventory** - User card ownership
- **scan_sessions** - Audit trail of scanning sessions

Migrations are automatically applied on startup.

## Security Features

- JWT-based authentication
- Input validation on all endpoints
- SQL injection prevention via prepared statements
- Foreign key constraints enabled
- Rate limiting compliance with Scryfall API
- CORS configuration for mobile clients

## Production Deployment

1. **Set a secure JWT secret**:
   ```bash
   export JWT_SECRET=$(openssl rand -base64 32)
   ```

2. **Configure database path**:
   ```bash
   export DATABASE_PATH=/var/lib/mtg_card_detector/data.db
   ```

3. **Run behind a reverse proxy** (nginx/caddy) with HTTPS

4. **Set up systemd service** for automatic restart

5. **Configure regular database backups**

## Architecture

```
cmd/server/          - Main application entry point
internal/
  ├── api/          - HTTP handlers and routing
  ├── auth/         - Authentication service
  ├── database/     - Database access layer
  ├── inventory/    - Inventory management
  ├── middleware/   - HTTP middleware (auth, logging)
  ├── models/       - Data models
  └── scanner/      - Card recognition (Scryfall integration)
config/             - Configuration management
migrations/         - Database migrations
```

## License

See LICENSE file in project root.
