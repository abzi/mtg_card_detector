package auth

import (
	"os"
	"testing"

	"github.com/abzi/mtg_card_detector/internal/database"
)

func setupTestDB(t *testing.T) *database.DB {
	dbPath := "/tmp/test_auth.db"
	os.Remove(dbPath)

	db, err := database.New(dbPath)
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	if err := db.RunMigrations("../../migrations"); err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	return db
}

func TestGenerateAnonymousUser(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	service := NewService(db, "test-secret")

	deviceID := "test-device-123"
	authResp, err := service.GenerateAnonymousUser(deviceID)
	if err != nil {
		t.Fatalf("Failed to generate anonymous user: %v", err)
	}

	if authResp.UserID == "" {
		t.Error("Expected user ID to be set")
	}

	if authResp.Token == "" {
		t.Error("Expected token to be set")
	}

	// Test retrieving same user with same device ID
	authResp2, err := service.GenerateAnonymousUser(deviceID)
	if err != nil {
		t.Fatalf("Failed to retrieve existing user: %v", err)
	}

	if authResp.UserID != authResp2.UserID {
		t.Error("Expected same user ID for same device ID")
	}
}

func TestValidateToken(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	service := NewService(db, "test-secret")

	// Create user and get token
	authResp, err := service.GenerateAnonymousUser("test-device")
	if err != nil {
		t.Fatalf("Failed to generate user: %v", err)
	}

	// Validate token
	userID, err := service.ValidateToken(authResp.Token)
	if err != nil {
		t.Fatalf("Failed to validate token: %v", err)
	}

	if userID != authResp.UserID {
		t.Errorf("Expected user ID %s, got %s", authResp.UserID, userID)
	}

	// Test invalid token
	_, err = service.ValidateToken("invalid-token")
	if err == nil {
		t.Error("Expected error for invalid token")
	}
}

func TestTokenExpiration(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	service := NewService(db, "test-secret")
	service.jwtSecret = []byte("test-secret")

	// This test verifies token structure, not actual expiration
	// (actual expiration would take too long to test)
	token, err := service.GenerateToken("test-user-id")
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	userID, err := service.ValidateToken(token)
	if err != nil {
		t.Fatalf("Failed to validate fresh token: %v", err)
	}

	if userID != "test-user-id" {
		t.Errorf("Expected user ID test-user-id, got %s", userID)
	}
}
