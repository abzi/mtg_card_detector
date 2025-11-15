# MTG Card Detector App - Development Plan

## Project Overview
A Magic The Gathering card scanner application with Android client and Go backend for card identification and inventory management.

## Architecture Summary
- **Client**: Android APK (Kotlin/Java) with camera-based scanning
- **Backend**: Go binary with SQLite database
- **Authentication**: Anonymous signin
- **Core Features**: Single scan, bulk scan, inventory management

---

## Phase 1: Project Foundation & Setup (Days 1-2)

### 1.1 Repository & Environment Setup
- [x] Initialize Git repository structure
- [ ] Set up Go module (`go.mod`)
- [ ] Create Android project structure
- [ ] Configure `.gitignore` for Go and Android
- [ ] Set up development environment documentation

### 1.2 Database Schema Design
- [ ] Design SQLite schema:
  - `users` table (anonymous user tracking)
  - `cards` table (MTG card master data)
  - `inventory` table (user card ownership)
  - `scan_sessions` table (scan history/audit)
- [ ] Create migration scripts
- [ ] Seed database with MTG card dataset

### 1.3 API Contract Definition
- [ ] Define REST API endpoints:
  - `POST /api/v1/auth/anonymous` - Generate anonymous user
  - `POST /api/v1/cards/scan` - Submit single card scan
  - `POST /api/v1/cards/scan/bulk` - Submit bulk scans
  - `GET /api/v1/inventory` - Retrieve user inventory
  - `GET /api/v1/cards/{id}` - Get card details
- [ ] Document request/response schemas (OpenAPI/Swagger)

---

## Phase 2: Backend Development (Days 3-7)

### 2.1 Core Go Application Structure
- [ ] Set up project layout:
  ```
  backend/
  ├── cmd/server/           # Main application
  ├── internal/
  │   ├── api/             # HTTP handlers
  │   ├── auth/            # Authentication logic
  │   ├── database/        # DB connections & queries
  │   ├── models/          # Data models
  │   ├── scanner/         # Card recognition logic
  │   └── inventory/       # Inventory management
  ├── migrations/          # SQL migrations
  └── config/              # Configuration
  ```

### 2.2 Database Layer
- [ ] Implement SQLite connection pool
- [ ] Create database access layer (DAL) with interfaces
- [ ] Implement CRUD operations for all tables
- [ ] Add transaction support for bulk operations
- [ ] Create indexes for performance

### 2.3 Authentication System
- [ ] Implement anonymous user generation (UUID-based)
- [ ] Create middleware for request authentication
- [ ] Token/session management (JWT or session tokens)
- [ ] Device fingerprinting for user persistence

### 2.4 Card Recognition System
- [ ] Research MTG card identification methods:
  - Option A: Barcode/Set number scanning
  - Option B: Image recognition (OCR for card name)
  - Option C: Scryfall API integration
- [ ] Implement card lookup service
- [ ] Create card validation logic
- [ ] Handle card variants and editions

### 2.5 Inventory Management
- [ ] Implement add-to-inventory logic
- [ ] Handle bulk inventory updates (transactions)
- [ ] Create inventory query/filter endpoints
- [ ] Add duplicate card handling
- [ ] Implement inventory statistics

### 2.6 API Implementation
- [ ] Set up HTTP router (chi/gorilla/gin)
- [ ] Implement all REST endpoints
- [ ] Add request validation
- [ ] Error handling middleware
- [ ] Logging middleware
- [ ] CORS configuration for mobile client

### 2.7 Testing
- [ ] Unit tests for business logic (80%+ coverage)
- [ ] Integration tests for database layer
- [ ] API endpoint tests
- [ ] Load testing for bulk operations

---

## Phase 3: Android Client Development (Days 8-14)

### 3.1 Project Setup
- [ ] Create Android Studio project (Kotlin)
- [ ] Set up dependencies:
  - CameraX for scanning
  - ML Kit or ZXing for barcode/text recognition
  - Retrofit for API calls
  - Room for local caching (optional)
  - Jetpack Compose or XML layouts
- [ ] Configure build variants (dev/prod)

### 3.2 Authentication Flow
- [ ] Implement anonymous signin on first launch
- [ ] Store user token securely (EncryptedSharedPreferences)
- [ ] Handle token refresh/expiration
- [ ] Device ID generation

### 3.3 Camera & Scanning UI
- [ ] Implement camera permission handling
- [ ] Create camera preview screen
- [ ] Add viewfinder/scan area overlay
- [ ] Implement single scan mode:
  - Capture image/barcode
  - Show card preview
  - Confirm/retry options
- [ ] Implement bulk scan mode:
  - Continuous scanning
  - Queue management
  - Progress indicator
  - Batch submission

### 3.4 Card Recognition
- [ ] Integrate ML Kit/ZXing for barcode scanning
- [ ] Implement OCR for card name extraction (fallback)
- [ ] Add client-side validation
- [ ] Handle scan failures gracefully

### 3.5 Networking Layer
- [ ] Create Retrofit API service
- [ ] Implement API call wrappers
- [ ] Add offline queue for scans (optional)
- [ ] Handle network errors/retry logic
- [ ] Request/response logging (debug mode)

### 3.6 Inventory UI
- [ ] Create inventory list screen
- [ ] Implement card detail view
- [ ] Add search/filter functionality
- [ ] Display card images (Scryfall integration)
- [ ] Show scan history

### 3.7 UX Enhancements
- [ ] Add haptic feedback for successful scans
- [ ] Implement loading states
- [ ] Create error screens with retry
- [ ] Add animations for card additions
- [ ] Tutorial/onboarding flow

### 3.8 Testing
- [ ] Unit tests for ViewModels/business logic
- [ ] UI tests for critical flows
- [ ] Camera integration testing
- [ ] Network layer testing (mock responses)

---

## Phase 4: Integration & Testing (Days 15-17)

### 4.1 End-to-End Testing
- [ ] Test complete scan-to-inventory flow
- [ ] Verify bulk scan performance (100+ cards)
- [ ] Test offline/online transitions
- [ ] Validate error handling across stack
- [ ] Test on multiple Android devices/versions

### 4.2 Performance Optimization
- [ ] Profile backend API response times
- [ ] Optimize database queries
- [ ] Reduce APK size
- [ ] Optimize camera preview performance
- [ ] Test under poor network conditions

### 4.3 Security Review
- [ ] Validate API authentication
- [ ] Test for common vulnerabilities (SQL injection, etc.)
- [ ] Secure data transmission (HTTPS)
- [ ] Review token storage security

---

## Phase 5: Deployment & Release (Days 18-20)

### 5.1 Backend Deployment
- [ ] Create production build configuration
- [ ] Set up hosting (VPS/cloud provider)
- [ ] Configure systemd service
- [ ] Set up HTTPS (Let's Encrypt)
- [ ] Configure backup strategy for SQLite
- [ ] Set up logging/monitoring

### 5.2 Android Release
- [ ] Generate signed APK/AAB
- [ ] Create app icon and assets
- [ ] Prepare release notes
- [ ] Test release build
- [ ] Distribute via:
  - Direct APK download
  - Or Google Play Store (requires setup)

### 5.3 Documentation
- [ ] API documentation (hosted Swagger)
- [ ] User guide/FAQ
- [ ] Developer setup instructions
- [ ] Deployment guide

---

## Phase 6: Post-Launch (Ongoing)

### 6.1 Monitoring & Maintenance
- [ ] Monitor API usage and errors
- [ ] Track scan success rates
- [ ] User feedback collection
- [ ] Bug fixes and updates

### 6.2 Future Enhancements
- [ ] Cloud sync across devices
- [ ] Export inventory (CSV/Excel)
- [ ] Card value tracking
- [ ] Trade/wishlist features
- [ ] Social features (share collection)

---

## Critical Path & Dependencies

### Must-Have for MVP:
1. Anonymous authentication ✓
2. Single card scanning ✓
3. Card identification ✓
4. Basic inventory storage ✓
5. Inventory viewing ✓

### Nice-to-Have:
1. Bulk scanning optimization
2. Offline support
3. Advanced inventory filtering
4. Card images/details

---

## Technology Stack Details

### Backend
- **Language**: Go 1.21+
- **Database**: SQLite 3
- **HTTP Router**: Chi or Gin
- **ORM/Query**: sqlc or GORM
- **Testing**: testify, httptest

### Android
- **Language**: Kotlin
- **Min SDK**: 24 (Android 7.0)
- **Target SDK**: 34 (Android 14)
- **Camera**: CameraX
- **Recognition**: ML Kit Barcode + ML Kit Text Recognition
- **Networking**: Retrofit + OkHttp
- **UI**: Jetpack Compose (modern) or XML
- **Architecture**: MVVM with ViewModel

### Card Data Source
- **Scryfall API**: Free MTG card database
- **Alternative**: MTG JSON

---

## Risk Mitigation

### Technical Risks:
1. **Card Recognition Accuracy**:
   - Mitigation: Combine barcode + OCR, manual entry fallback
2. **Bulk Scan Performance**:
   - Mitigation: Client-side queuing, batch API calls
3. **Database Scalability**:
   - Mitigation: Proper indexing, consider migration path to PostgreSQL

### Timeline Risks:
1. **Android Camera Integration Complexity**:
   - Buffer: 2-3 extra days allocated
2. **Card Data Sourcing**:
   - Mitigation: Use Scryfall API, pre-seed database

---

## Development Best Practices

1. **Version Control**: Feature branches, meaningful commits
2. **Code Review**: Self-review before commits
3. **Testing**: Write tests alongside features
4. **Documentation**: Inline comments, README updates
5. **Error Handling**: Comprehensive logging, user-friendly messages

---

## Estimated Timeline: 20 Days (4 Weeks)

- **Week 1**: Foundation + Backend Core
- **Week 2**: Backend Completion + Android Setup
- **Week 3**: Android Development
- **Week 4**: Integration, Testing, Deployment

**Note**: Timeline assumes single developer, full-time work. Adjust for part-time or team scenarios.

---

## Next Immediate Steps

1. Set up Go module and project structure
2. Design and create SQLite database schema
3. Set up Android project in Android Studio
4. Create API contract documentation
5. Research MTG card identification approach (Scryfall API integration)

---

*Plan Version: 1.0*
*Last Updated: 2025-11-15*
