# Card Recognition Research & Analysis

**Date**: 2025-11-15
**Status**: Research Phase - Decision Pending

---

## Problem Statement

The Android app is **NOT able to identify any card** using the current camera-based scanning approach.

---

## Root Cause Analysis

### Technical Issues Identified

#### Issue 1: ImageCapture Not Bound to Lifecycle
```kotlin
// ScanActivity.kt ~line 95
val imageCapture = ImageCapture.Builder().build()
imageCapture.takePicture(...) // FAILS - no camera bound!
```

**Problem**: ImageCapture is created in `captureAndScan()` but never added to `cameraProvider.bindToLifecycle()`, so `takePicture()` fails silently.

#### Issue 2: No Continuous Image Analysis
- Camera preview shows but there's no **ImageAnalysis** use case
- ML Kit barcode/text scanners initialized but never receive images
- Manual flow: tap button → create ImageCapture → attempt takePicture → silent failure

#### Issue 3: Generic OCR Limitations
- Most MTG cards have **no barcodes**
- Generic ML Kit OCR struggles with:
  - Stylized card name fonts
  - Foil/holographic cards
  - Varying lighting conditions
  - Camera angles and focus

---

## Proposed Solutions

### Option 1: Visual Recognition API ⭐ (Recommended)
Send card photo to specialized image recognition service.

**Implementation**: Photo capture → API call → Card identification → Add to inventory

---

## API Research Results

### Commercial Solutions

#### 1. **Ximilar Trading Card Game Identifier** ⭐⭐⭐
- **URL**: https://www.ximilar.com/blog/build-your-own-trading-card-game-identifier-with-our-api/
- **Type**: Commercial REST API
- **Function**: Upload photo → Returns card info as JSON
- **Accuracy**: Built specifically for TCGs including MTG
- **Pricing**: Paid service (pricing TBD)
- **Integration**:
  - Android: Capture photo → Upload to Ximilar
  - Response: Card name, set, details
  - Backend: Use returned data to query Scryfall for full info
- **Pros**:
  - Professional, proven accuracy
  - JSON output ready for integration
  - Handles various card conditions
- **Cons**:
  - Ongoing cost
  - External dependency
  - Requires internet connection

#### 2. **Roboflow MTG Card Scanner API**
- **URL**: https://universe.roboflow.com/mtg-scanner/mtg-card-scanner/model/2
- **Type**: Object detection API
- **Version**: v2 (2023-10-05)
- **Pricing**: Free tier available, paid plans
- **Pros**:
  - Specialized for MTG
  - Quick to integrate
- **Cons**:
  - May require Roboflow account
  - Credit-based system
- **Integration Effort**: Low

---

### Open Source / Self-Hosted Solutions

#### 3. **YamCR (Yet-Another-Magic-Card-Recognizer)** ⭐⭐
- **URL**: https://github.com/ForOhForError/Yet-Another-Magic-Card-Recognizer
- **Type**: Open source perceptual hashing
- **How it works**:
  1. Calculate perceptual hash of card image
  2. Match hash against pre-computed Scryfall database
  3. Return best match
- **Data Source**: Uses Scryfall API/data (which we already use)
- **Pros**:
  - FREE
  - Works offline after initial DB download
  - Well-maintained project
  - High accuracy with perceptual hashing
- **Cons**:
  - Need to implement hashing algorithm
  - Database maintenance required
  - More complex integration
- **Implementation Options**:
  - **Option A**: Run on Go backend
    - Android sends photo → Backend hashes → Returns card
  - **Option B**: Run on Android
    - Implement hashing in Kotlin
    - Store hash database locally
    - Fully offline capable
- **Integration Effort**: Medium-High

#### 4. **MTGScan**
- **URL**: https://github.com/fortierq/mtgscan
- **Type**: Azure OCR + fuzzy search
- **How it works**: Azure OCR → Extract text → Fuzzy match against MTGJSON
- **Data Source**: MTGJSON
- **Pros**:
  - OCR-based approach
  - Blog post with implementation details: https://fortierq.github.io/mtgscan-ocr-azure-flask-celery-socketio/
- **Cons**:
  - Requires Azure account (paid)
  - OCR still has accuracy issues with MTG cards
  - More complex than generic OCR we already tried
- **Integration Effort**: Medium

#### 5. **Other GitHub Projects**
- **mtg_card_detector**: https://github.com/hj3yoo/mtg_card_detector
  - Computer vision + perceptual hashing
  - Real-time video identification
- **MTG-Card-Reader**: https://github.com/TrifectaIII/MTG-Card-Reader
  - Webcam-based reader
  - Set-specific recognition

---

### Scryfall API Investigation

**Finding**: ❌ Scryfall does **NOT** have reverse image search

- **What Scryfall offers**:
  - Text-based card search
  - Card data retrieval (when you know the card)
  - Image URLs for displaying cards

- **What Scryfall does NOT offer**:
  - Image upload for recognition
  - Visual similarity search
  - Computer vision endpoints

**Implication**: We still need Scryfall for card data, but need a separate solution for image recognition.

---

## Recommended Approaches

### Approach A: Commercial API (Fast, Reliable) ⭐⭐⭐

**Best for**: Quick deployment, reliable results, minimal development

**Implementation**:
1. Integrate Ximilar or Roboflow API
2. Android: Capture photo → Upload to API
3. Receive card identification (name, set)
4. Query our existing Scryfall integration for full card data
5. Add to inventory

**Timeline**: 1-2 days
**Cost**: Ongoing API fees
**Accuracy**: High (90%+)

---

### Approach B: Perceptual Hashing (Free, Complex) ⭐⭐

**Best for**: No budget, full control, offline capability

**Implementation**:
1. Study YamCR algorithm
2. Implement perceptual hashing in Go backend or Kotlin
3. Download/generate Scryfall card hash database
4. Android: Capture photo → Hash → Match → Return card
5. Add to inventory

**Timeline**: 5-7 days
**Cost**: Free
**Accuracy**: High (85-95%) with good hashing

---

### Approach C: Enhanced Manual Input (Pragmatic) ⭐

**Best for**: Immediate usability, fallback option

**Implementation**:
1. Keep manual text input as primary method
2. Add Scryfall autocomplete as user types
3. Show card preview thumbnail
4. Quick access to recent/popular cards
5. Optional: Keep camera for visual reference only

**Timeline**: 1 day
**Cost**: Free
**Accuracy**: 100% (user validates)

---

## Decision Criteria

| Criteria | Commercial API | Perceptual Hashing | Enhanced Manual |
|----------|---------------|-------------------|-----------------|
| **Cost** | Ongoing fees | Free | Free |
| **Accuracy** | 90%+ | 85-95% | 100% |
| **Speed** | Fast | Fast | Moderate |
| **Internet Required** | Yes | Initial only | Yes (for Scryfall) |
| **Development Time** | 1-2 days | 5-7 days | 1 day |
| **Maintenance** | Low | Medium | Low |
| **User Experience** | Automatic | Automatic | Manual |

---

## Next Steps (Pending Decision)

### If choosing Commercial API:
1. Sign up for Ximilar or Roboflow trial
2. Test accuracy with sample MTG cards
3. Integrate API into Android app
4. Test end-to-end flow
5. Evaluate cost vs accuracy

### If choosing Perceptual Hashing:
1. Clone YamCR repository for reference
2. Study hashing algorithm
3. Decide: Backend (Go) or Android (Kotlin) implementation
4. Download Scryfall bulk data
5. Generate hash database
6. Implement matching algorithm
7. Extensive testing

### If choosing Enhanced Manual:
1. Add Scryfall autocomplete API integration
2. Implement search-as-you-type UI
3. Add card preview thumbnails
4. Test with real usage patterns

---

## Technical Architecture (If using API)

```
┌─────────────────┐
│  Android App    │
│                 │
│  1. Capture     │
│     Photo       │
└────────┬────────┘
         │
         │ HTTPS POST (image)
         ▼
┌─────────────────┐
│  Recognition    │
│  API            │
│  (Ximilar/      │
│   Roboflow)     │
└────────┬────────┘
         │
         │ JSON Response
         │ {name, set, ...}
         ▼
┌─────────────────┐
│  Go Backend     │
│                 │
│  - Validate     │
│  - Query        │◄──────┐
│    Scryfall     │       │
│  - Add to       │       │
│    Inventory    │       │
└────────┬────────┘       │
         │                │
         │         ┌──────┴──────┐
         │         │  Scryfall   │
         │         │  API        │
         │         │  (Card Data)│
         │         └─────────────┘
         ▼
┌─────────────────┐
│  SQLite DB      │
│  (Inventory)    │
└─────────────────┘
```

---

## Additional Research Links

- **Ximilar Blog**: https://www.ximilar.com/blog/build-your-own-trading-card-game-identifier-with-our-api/
- **YamCR GitHub**: https://github.com/ForOhForError/Yet-Another-Magic-Card-Recognizer
- **MTGScan GitHub**: https://github.com/fortierq/mtgscan
- **MTGScan Blog**: https://fortierq.github.io/mtgscan-ocr-azure-flask-celery-socketio/
- **Roboflow Model**: https://universe.roboflow.com/mtg-scanner/mtg-card-scanner/model/2
- **Scryfall API Docs**: https://scryfall.com/docs/api

---

## Open Questions

1. **Budget**: Is there budget for a commercial API (Ximilar/Roboflow)?
2. **Priority**: Speed to market vs cost optimization?
3. **Internet**: Is always-online acceptable for users?
4. **Accuracy**: What's the acceptable recognition accuracy threshold?
5. **Fallback**: Should manual input always be available as backup?

---

**Status**: Awaiting decision on which approach to implement.

**Contact**: Ready to implement once direction is chosen.
