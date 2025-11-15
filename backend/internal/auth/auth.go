package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/abzi/mtg_card_detector/internal/database"
	"github.com/abzi/mtg_card_detector/internal/models"
)

const (
	TokenExpiration = 365 * 24 * time.Hour // 1 year for anonymous users
)

type Service struct {
	db        *database.DB
	jwtSecret []byte
}

// NewService creates a new auth service
func NewService(db *database.DB, jwtSecret string) *Service {
	return &Service{
		db:        db,
		jwtSecret: []byte(jwtSecret),
	}
}

// GenerateAnonymousUser creates a new anonymous user or retrieves existing one
func (s *Service) GenerateAnonymousUser(deviceID string) (*models.AuthResponse, error) {
	// Check if user with this device ID already exists
	user, err := s.db.GetUserByDeviceID(deviceID)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}

	// Create new user if not exists
	if user == nil {
		user = &models.User{
			ID:        uuid.New().String(),
			DeviceID:  deviceID,
			CreatedAt: time.Now(),
			LastSeen:  time.Now(),
		}
		if err := s.db.CreateUser(user); err != nil {
			return nil, fmt.Errorf("failed to create user: %w", err)
		}
	} else {
		// Update last seen
		if err := s.db.UpdateUserLastSeen(user.ID); err != nil {
			return nil, fmt.Errorf("failed to update last seen: %w", err)
		}
	}

	// Generate JWT token
	token, err := s.GenerateToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &models.AuthResponse{
		UserID: user.ID,
		Token:  token,
	}, nil
}

// GenerateToken generates a JWT token for the user
func (s *Service) GenerateToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(TokenExpiration).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// ValidateToken validates a JWT token and returns the user ID
func (s *Service) ValidateToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.jwtSecret, nil
	})

	if err != nil {
		return "", fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return "", fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("failed to parse claims")
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", fmt.Errorf("invalid user_id claim")
	}

	return userID, nil
}
