# MTG Card Detector - Android Client

Android application for scanning Magic: The Gathering cards and managing your card inventory.

## Features

- **Single Scan Mode**: Scan one card at a time and immediately add to inventory
- **Bulk Scan Mode**: Scan multiple cards in a session before submitting
- **Card Recognition**: Uses ML Kit for barcode and text recognition
- **Inventory Management**: View all cards in your collection with images
- **Anonymous Authentication**: Automatic user creation with device-based identification
- **Secure Storage**: Encrypted SharedPreferences for token storage

## Requirements

- Android 7.0 (API 24) or higher
- Camera permission
- Internet connection

## Building the App

### Prerequisites

- Android Studio Hedgehog or newer
- JDK 17
- Android SDK 34

### Steps

1. Open the project in Android Studio:
   ```bash
   cd android
   # Open with Android Studio
   ```

2. Sync Gradle files

3. Configure backend URL in `app/build.gradle`:
   ```gradle
   buildConfigField "String", "API_BASE_URL", "\"http://YOUR_SERVER:8080/api/v1\""
   ```

4. Build the app:
   - For debug: Build > Build Bundle(s) / APK(s) > Build APK(s)
   - For release: Build > Generate Signed Bundle / APK

### Running on Emulator

The default API URL `http://10.0.2.2:8080/api/v1` works for Android emulator connecting to localhost.

### Running on Physical Device

Update the API_BASE_URL to your computer's local network IP address:
```gradle
buildConfigField "String", "API_BASE_URL", "\"http://192.168.1.XXX:8080/api/v1\""
```

## Project Structure

```
app/src/main/java/com/mtgdetector/
├── MTGDetectorApp.kt          # Application class
├── auth/
│   └── AuthManager.kt         # Authentication & token management
├── models/
│   └── Models.kt              # Data models
├── network/
│   ├── ApiService.kt          # Retrofit API interface
│   └── RetrofitClient.kt      # HTTP client configuration
└── ui/
    ├── MainActivity.kt        # Main menu
    ├── ScanActivity.kt        # Camera scanning
    ├── InventoryActivity.kt   # Inventory list
    └── InventoryAdapter.kt    # RecyclerView adapter
```

## Key Technologies

- **Kotlin**: Primary language
- **CameraX**: Camera API
- **ML Kit**: Barcode scanning & text recognition
- **Retrofit**: HTTP client
- **Glide**: Image loading
- **Material Design**: UI components
- **Coroutines**: Async operations
- **EncryptedSharedPreferences**: Secure storage

## Usage

1. **First Launch**: App automatically creates an anonymous user account
2. **Scan Card**: Tap "Single Scan" or "Bulk Scan"
3. **Grant Permission**: Allow camera access when prompted
4. **Scan**: Point camera at MTG card and tap "Scan Card"
5. **View Inventory**: Tap "View Inventory" to see your collection

## Scanning Tips

- Ensure good lighting
- Hold camera steady
- Try to capture the card name clearly
- For best results, scan barcodes if available

## Security

- Auth tokens stored in EncryptedSharedPreferences
- HTTPS recommended for production
- Device ID generated using UUID

## Troubleshooting

### Camera not working
- Check camera permission granted
- Restart the app
- Check device compatibility

### Cards not scanning
- Improve lighting
- Move closer to card
- Ensure card text is visible
- Try manual entry (future feature)

### Network errors
- Check backend server is running
- Verify API_BASE_URL is correct
- Check firewall settings
- Ensure device has internet connection

## Building for Production

1. Update API URL to production server
2. Generate signing key:
   ```bash
   keytool -genkey -v -keystore release.keystore -alias mtg_detector -keyalg RSA -keysize 2048 -validity 10000
   ```
3. Configure signing in `app/build.gradle`
4. Build release APK/AAB
5. Test thoroughly before distribution

## License

See LICENSE file in project root.
