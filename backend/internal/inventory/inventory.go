package inventory

import (
	"fmt"
	"time"

	"github.com/abzi/mtg_card_detector/internal/database"
	"github.com/abzi/mtg_card_detector/internal/models"
	"github.com/abzi/mtg_card_detector/internal/scanner"
)

type Service struct {
	db      *database.DB
	scanner *scanner.Service
}

// NewService creates a new inventory service
func NewService(db *database.DB, scanner *scanner.Service) *Service {
	return &Service{
		db:      db,
		scanner: scanner,
	}
}

// ProcessSingleScan processes a single card scan and adds to inventory
func (s *Service) ProcessSingleScan(userID string, req *models.ScanRequest) (*models.ScanResponse, error) {
	// Create scan session
	session := &models.ScanSession{
		UserID:    userID,
		ScanType:  "single",
		StartedAt: time.Now(),
	}
	sessionID, err := s.db.CreateScanSession(session)
	if err != nil {
		return nil, fmt.Errorf("failed to create scan session: %w", err)
	}

	// Scan the card
	card, err := s.scanner.ScanCard(req)
	if err != nil {
		// Update session with failure
		s.db.UpdateScanSession(sessionID, 1, 0, 1)
		return &models.ScanResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	// Add to inventory
	if err := s.db.AddToInventory(userID, card.ID, 1); err != nil {
		s.db.UpdateScanSession(sessionID, 1, 0, 1)
		return &models.ScanResponse{
			Success: false,
			Error:   fmt.Sprintf("failed to add to inventory: %v", err),
		}, nil
	}

	// Update session with success
	s.db.UpdateScanSession(sessionID, 1, 1, 0)

	return &models.ScanResponse{
		Success: true,
		Card:    card,
	}, nil
}

// ProcessBulkScan processes multiple card scans
func (s *Service) ProcessBulkScan(userID string, req *models.BulkScanRequest) (*models.BulkScanResponse, error) {
	// Create scan session
	session := &models.ScanSession{
		UserID:    userID,
		ScanType:  "bulk",
		StartedAt: time.Now(),
	}
	sessionID, err := s.db.CreateScanSession(session)
	if err != nil {
		return nil, fmt.Errorf("failed to create scan session: %w", err)
	}

	results := make([]models.ScanResponse, 0, len(req.Scans))
	successful := 0
	failed := 0

	for _, scanReq := range req.Scans {
		card, err := s.scanner.ScanCard(&scanReq)
		if err != nil {
			failed++
			results = append(results, models.ScanResponse{
				Success: false,
				Error:   err.Error(),
			})
			continue
		}

		// Add to inventory
		if err := s.db.AddToInventory(userID, card.ID, 1); err != nil {
			failed++
			results = append(results, models.ScanResponse{
				Success: false,
				Error:   fmt.Sprintf("failed to add to inventory: %v", err),
			})
			continue
		}

		successful++
		results = append(results, models.ScanResponse{
			Success: true,
			Card:    card,
		})
	}

	// Update scan session
	s.db.UpdateScanSession(sessionID, len(req.Scans), successful, failed)

	return &models.BulkScanResponse{
		SessionID:       sessionID,
		TotalScanned:    len(req.Scans),
		SuccessfulScans: successful,
		FailedScans:     failed,
		Results:         results,
	}, nil
}

// GetInventory retrieves user's inventory
func (s *Service) GetInventory(userID string) ([]models.InventoryItem, error) {
	return s.db.GetUserInventory(userID)
}

// GetInventoryStats retrieves inventory statistics
func (s *Service) GetInventoryStats(userID string) (map[string]interface{}, error) {
	count, err := s.db.GetInventoryCount(userID)
	if err != nil {
		return nil, err
	}

	inventory, err := s.db.GetUserInventory(userID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"total_cards":       count,
		"unique_cards":      len(inventory),
		"last_updated":      time.Now(),
	}, nil
}
