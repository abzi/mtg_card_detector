package models

import "time"

// User represents an anonymous user
type User struct {
	ID        string    `json:"id"`
	DeviceID  string    `json:"device_id"`
	CreatedAt time.Time `json:"created_at"`
	LastSeen  time.Time `json:"last_seen"`
}

// Card represents a Magic: The Gathering card
type Card struct {
	ID              string    `json:"id"`
	ScryfallID      string    `json:"scryfall_id,omitempty"`
	Name            string    `json:"name"`
	SetCode         string    `json:"set_code"`
	CollectorNumber string    `json:"collector_number"`
	ImageURI        string    `json:"image_uri,omitempty"`
	OracleText      string    `json:"oracle_text,omitempty"`
	TypeLine        string    `json:"type_line,omitempty"`
	ManaCost        string    `json:"mana_cost,omitempty"`
	Rarity          string    `json:"rarity,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
}

// InventoryItem represents a card in a user's inventory
type InventoryItem struct {
	ID       int       `json:"id"`
	UserID   string    `json:"user_id"`
	CardID   string    `json:"card_id"`
	Quantity int       `json:"quantity"`
	AddedAt  time.Time `json:"added_at"`
	Card     *Card     `json:"card,omitempty"`
}

// ScanSession represents a scanning session
type ScanSession struct {
	ID              int       `json:"id"`
	UserID          string    `json:"user_id"`
	ScanType        string    `json:"scan_type"`
	CardsScanned    int       `json:"cards_scanned"`
	SuccessfulScans int       `json:"successful_scans"`
	FailedScans     int       `json:"failed_scans"`
	StartedAt       time.Time `json:"started_at"`
	CompletedAt     *time.Time `json:"completed_at,omitempty"`
}

// ScanRequest represents a card scan request
type ScanRequest struct {
	CardName        string `json:"card_name,omitempty"`
	SetCode         string `json:"set_code,omitempty"`
	CollectorNumber string `json:"collector_number,omitempty"`
	Barcode         string `json:"barcode,omitempty"`
}

// BulkScanRequest represents multiple card scans
type BulkScanRequest struct {
	Scans []ScanRequest `json:"scans"`
}

// ScanResponse represents the result of a scan
type ScanResponse struct {
	Success bool   `json:"success"`
	Card    *Card  `json:"card,omitempty"`
	Error   string `json:"error,omitempty"`
}

// BulkScanResponse represents results of bulk scanning
type BulkScanResponse struct {
	SessionID       int            `json:"session_id"`
	TotalScanned    int            `json:"total_scanned"`
	SuccessfulScans int            `json:"successful_scans"`
	FailedScans     int            `json:"failed_scans"`
	Results         []ScanResponse `json:"results"`
}

// AuthResponse represents authentication response
type AuthResponse struct {
	UserID string `json:"user_id"`
	Token  string `json:"token"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}
