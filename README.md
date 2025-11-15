# MTG Card Detector

A complete Magic: The Gathering card scanner application with Android client and Go backend. Scan physical MTG cards using your phone's camera, automatically identify them, and manage your collection inventory.

![License](https://img.shields.io/badge/license-MIT-blue.svg)
![Go Version](https://img.shields.io/badge/Go-1.21-00ADD8?logo=go)
![Android](https://img.shields.io/badge/Android-7.0+-3DDC84?logo=android)

## Features

### 🎴 Card Scanning
- **Single Scan Mode**: Scan one card at a time
- **Bulk Scan Mode**: Scan multiple cards in a session
- **Smart Recognition**: Combines barcode scanning and OCR
- **Scryfall Integration**: Comprehensive card database

### 📱 Mobile App
- **Camera-based scanning** with CameraX
- **ML Kit integration** for barcode and text recognition
- **Anonymous authentication** - no signup required
- **Inventory management** with card images
- **Offline-first architecture** (future enhancement)

### 🔧 Backend API
- **RESTful API** built with Go
- **SQLite database** for lightweight deployment
- **JWT authentication** for security
- **Rate limiting compliant** with Scryfall API
- **Docker support** for easy deployment

## Quick Start

### Prerequisites

- **Backend**: Go 1.21+, SQLite 3
- **Android**: Android Studio, JDK 17, Android SDK 34
- **Optional**: Docker & Docker Compose

### 1. Backend Setup

```bash
# Clone repository
git clone https://github.com/abzi/mtg_card_detector.git
cd mtg_card_detector/backend

# Install dependencies
go mod download

# Build and run
go build -o server ./cmd/server
./server
```

Server runs on `http://localhost:8080`

**Using Docker:**
```bash
docker-compose up -d
```

### 2. Android Setup

```bash
cd android

# Configure API URL in app/build.gradle
# For emulator: http://10.0.2.2:8080/api/v1
# For device: http://YOUR_IP:8080/api/v1

# Build in Android Studio or:
./gradlew assembleDebug
```

APK location: `android/app/build/outputs/apk/debug/app-debug.apk`

## Architecture

```
┌─────────────────┐
│  Android App    │
│  (Kotlin)       │
└────────┬────────┘
         │ HTTPS/HTTP
         │ REST API
┌────────▼────────┐      ┌──────────────┐
│   Go Backend    │◄────►│   SQLite DB  │
│   (Chi Router)  │      │              │
└────────┬────────┘      └──────────────┘
         │
         │ HTTPS
┌────────▼────────┐
│  Scryfall API   │
│  (Card Data)    │
└─────────────────┘
```

## Project Structure

```
mtg_card_detector/
├── backend/              # Go backend service
│   ├── cmd/server/      # Main application
│   ├── internal/        # Internal packages
│   │   ├── api/        # HTTP handlers
│   │   ├── auth/       # JWT authentication
│   │   ├── database/   # SQLite DAL
│   │   ├── models/     # Data models
│   │   ├── scanner/    # Card recognition
│   │   └── inventory/  # Inventory logic
│   └── migrations/     # Database migrations
│
├── android/             # Android application
│   └── app/src/main/
│       ├── java/com/mtgdetector/
│       │   ├── ui/     # Activities & UI
│       │   ├── network/# Retrofit client
│       │   ├── auth/   # Auth manager
│       │   └── models/ # Data models
│       └── res/        # Resources & layouts
│
├── CLAUDE.md           # Project instructions
├── CLAUDE2.md          # Development plan
├── DEPLOYMENT.md       # Deployment guide
├── Dockerfile          # Backend container
└── docker-compose.yml  # Docker orchestration
```

## API Endpoints

### Public

- `GET /health` - Health check
- `POST /api/v1/auth/anonymous` - Anonymous authentication

### Protected (requires Bearer token)

- `POST /api/v1/cards/scan` - Single card scan
- `POST /api/v1/cards/scan/bulk` - Bulk card scan
- `GET /api/v1/inventory` - Get user inventory
- `GET /api/v1/cards?id=<id>` - Get card details

See [backend/README.md](backend/README.md) for detailed API documentation.

## Technology Stack

### Backend
- **Language**: Go 1.21
- **Framework**: Chi router
- **Database**: SQLite 3
- **Auth**: JWT (golang-jwt)
- **External API**: Scryfall

### Android
- **Language**: Kotlin
- **Min SDK**: 24 (Android 7.0)
- **Target SDK**: 34 (Android 14)
- **Camera**: CameraX
- **ML**: Google ML Kit
- **Networking**: Retrofit + OkHttp
- **Image Loading**: Glide
- **Security**: EncryptedSharedPreferences

## Development

### Running Tests

**Backend:**
```bash
cd backend
go test ./... -v
go test ./... -cover
```

**Android:**
```bash
cd android
./gradlew test
./gradlew connectedAndroidTest
```

### Environment Variables

**Backend:**
```bash
export PORT=8080
export DATABASE_PATH=./data/mtg_cards.db
export JWT_SECRET="your-secure-secret-here"
export MIGRATIONS_PATH=./migrations
```

**Android:**
Update `buildConfigField` in `app/build.gradle`

## Deployment

See [DEPLOYMENT.md](DEPLOYMENT.md) for comprehensive deployment instructions including:
- Docker deployment
- Manual Linux deployment
- Cloud platforms (Heroku, AWS, DigitalOcean)
- Android APK distribution
- Security best practices
- Monitoring setup

### Quick Deploy with Docker

```bash
# Set secure JWT secret
export JWT_SECRET=$(openssl rand -base64 32)

# Start services
docker-compose up -d

# View logs
docker-compose logs -f
```

## Security Features

- ✅ JWT-based authentication
- ✅ Encrypted token storage on Android
- ✅ SQL injection prevention (prepared statements)
- ✅ Input validation
- ✅ HTTPS support (production)
- ✅ CORS configuration
- ✅ Rate limiting compliance

## Roadmap

- [ ] Cloud sync across devices
- [ ] Manual card entry
- [ ] Export inventory (CSV/PDF)
- [ ] Card value tracking
- [ ] Trade/wishlist features
- [ ] Deck building
- [ ] PostgreSQL migration
- [ ] iOS app

## Contributing

Contributions welcome! Please:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see [LICENSE](LICENSE) file for details.

## Acknowledgments

- [Scryfall](https://scryfall.com/) - MTG card database API
- [Google ML Kit](https://developers.google.com/ml-kit) - On-device ML
- [CameraX](https://developer.android.com/camerax) - Camera library

## Support

- **Issues**: [GitHub Issues](https://github.com/abzi/mtg_card_detector/issues)
- **Discussions**: [GitHub Discussions](https://github.com/abzi/mtg_card_detector/discussions)

## Screenshots

_(Add screenshots here once app is built)_

---

**Note**: This app is not affiliated with or endorsed by Wizards of the Coast. Magic: The Gathering is a trademark of Wizards of the Coast LLC.
