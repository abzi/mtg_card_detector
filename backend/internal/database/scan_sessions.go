package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/abzi/mtg_card_detector/internal/models"
)

// CreateScanSession creates a new scan session
func (db *DB) CreateScanSession(session *models.ScanSession) (int, error) {
	query := `INSERT INTO scan_sessions (user_id, scan_type, cards_scanned, successful_scans, failed_scans, started_at)
	          VALUES (?, ?, ?, ?, ?, ?)`
	result, err := db.Exec(query, session.UserID, session.ScanType, session.CardsScanned,
		session.SuccessfulScans, session.FailedScans, session.StartedAt)
	if err != nil {
		return 0, fmt.Errorf("failed to create scan session: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get scan session ID: %w", err)
	}

	return int(id), nil
}

// UpdateScanSession updates an existing scan session
func (db *DB) UpdateScanSession(sessionID int, cardsScanned, successful, failed int) error {
	query := `UPDATE scan_sessions
	          SET cards_scanned = ?, successful_scans = ?, failed_scans = ?, completed_at = ?
	          WHERE id = ?`
	_, err := db.Exec(query, cardsScanned, successful, failed, time.Now(), sessionID)
	if err != nil {
		return fmt.Errorf("failed to update scan session: %w", err)
	}
	return nil
}

// GetScanSession retrieves a scan session by ID
func (db *DB) GetScanSession(sessionID int) (*models.ScanSession, error) {
	query := `SELECT id, user_id, scan_type, cards_scanned, successful_scans, failed_scans, started_at, completed_at
	          FROM scan_sessions WHERE id = ?`

	session := &models.ScanSession{}
	var completedAt sql.NullTime
	err := db.QueryRow(query, sessionID).Scan(&session.ID, &session.UserID, &session.ScanType,
		&session.CardsScanned, &session.SuccessfulScans, &session.FailedScans, &session.StartedAt, &completedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get scan session: %w", err)
	}

	if completedAt.Valid {
		session.CompletedAt = &completedAt.Time
	}

	return session, nil
}
