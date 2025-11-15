package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/abzi/mtg_card_detector/internal/models"
)

// CreateUser creates a new anonymous user
func (db *DB) CreateUser(user *models.User) error {
	query := `INSERT INTO users (id, device_id, created_at, last_seen) VALUES (?, ?, ?, ?)`
	_, err := db.Exec(query, user.ID, user.DeviceID, user.CreatedAt, user.LastSeen)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

// GetUserByID retrieves a user by ID
func (db *DB) GetUserByID(id string) (*models.User, error) {
	query := `SELECT id, device_id, created_at, last_seen FROM users WHERE id = ?`
	user := &models.User{}
	err := db.QueryRow(query, id).Scan(&user.ID, &user.DeviceID, &user.CreatedAt, &user.LastSeen)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

// GetUserByDeviceID retrieves a user by device ID
func (db *DB) GetUserByDeviceID(deviceID string) (*models.User, error) {
	query := `SELECT id, device_id, created_at, last_seen FROM users WHERE device_id = ?`
	user := &models.User{}
	err := db.QueryRow(query, deviceID).Scan(&user.ID, &user.DeviceID, &user.CreatedAt, &user.LastSeen)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by device ID: %w", err)
	}
	return user, nil
}

// UpdateUserLastSeen updates the last seen timestamp
func (db *DB) UpdateUserLastSeen(userID string) error {
	query := `UPDATE users SET last_seen = ? WHERE id = ?`
	_, err := db.Exec(query, time.Now(), userID)
	if err != nil {
		return fmt.Errorf("failed to update user last seen: %w", err)
	}
	return nil
}
