package scanner

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/abzi/mtg_card_detector/internal/database"
	"github.com/abzi/mtg_card_detector/internal/models"
	"github.com/google/uuid"
)

const (
	ScryfallAPIBase = "https://api.scryfall.com"
	RateLimit       = time.Millisecond * 100 // Scryfall rate limit: 10 requests per second
)

type ScryfallCard struct {
	ID              string                 `json:"id"`
	Name            string                 `json:"name"`
	SetCode         string                 `json:"set"`
	CollectorNumber string                 `json:"collector_number"`
	ImageURIs       map[string]string      `json:"image_uris,omitempty"`
	CardFaces       []map[string]interface{} `json:"card_faces,omitempty"`
	OracleText      string                 `json:"oracle_text,omitempty"`
	TypeLine        string                 `json:"type_line"`
	ManaCost        string                 `json:"mana_cost,omitempty"`
	Rarity          string                 `json:"rarity"`
}

type Service struct {
	db         *database.DB
	httpClient *http.Client
	lastCall   time.Time
}

// NewService creates a new scanner service
func NewService(db *database.DB) *Service {
	return &Service{
		db: db,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		lastCall: time.Time{},
	}
}

// ScanCard identifies a card from scan data
func (s *Service) ScanCard(req *models.ScanRequest) (*models.Card, error) {
	// Try to find card in local database first
	var card *models.Card
	var err error

	if req.SetCode != "" && req.CollectorNumber != "" {
		card, err = s.db.GetCardBySetAndNumber(req.SetCode, req.CollectorNumber)
		if err != nil {
			return nil, fmt.Errorf("database error: %w", err)
		}
		if card != nil {
			return card, nil
		}
	}

	// If not found locally, query Scryfall API
	if req.SetCode != "" && req.CollectorNumber != "" {
		card, err = s.fetchFromScryfallBySetNumber(req.SetCode, req.CollectorNumber)
	} else if req.CardName != "" {
		card, err = s.fetchFromScryfallByName(req.CardName, req.SetCode)
	} else {
		return nil, fmt.Errorf("insufficient scan data")
	}

	if err != nil {
		return nil, err
	}

	// Store the card in local database
	if card != nil {
		if err := s.db.CreateCard(card); err != nil {
			// Ignore duplicate errors, card might have been added by another request
			if !strings.Contains(err.Error(), "UNIQUE constraint failed") {
				return nil, fmt.Errorf("failed to store card: %w", err)
			}
		}
	}

	return card, nil
}

// fetchFromScryfallBySetNumber fetches card data from Scryfall by set and collector number
func (s *Service) fetchFromScryfallBySetNumber(setCode, collectorNumber string) (*models.Card, error) {
	s.rateLimit()

	url := fmt.Sprintf("%s/cards/%s/%s", ScryfallAPIBase,
		url.PathEscape(strings.ToLower(setCode)),
		url.PathEscape(collectorNumber))

	scryfallCard, err := s.makeRequest(url)
	if err != nil {
		return nil, err
	}

	return s.convertScryfallCard(scryfallCard), nil
}

// fetchFromScryfallByName fetches card data from Scryfall by name
func (s *Service) fetchFromScryfallByName(name, setCode string) (*models.Card, error) {
	s.rateLimit()

	searchURL := fmt.Sprintf("%s/cards/named?fuzzy=%s", ScryfallAPIBase, url.QueryEscape(name))
	if setCode != "" {
		searchURL += fmt.Sprintf("&set=%s", url.QueryEscape(setCode))
	}

	scryfallCard, err := s.makeRequest(searchURL)
	if err != nil {
		return nil, err
	}

	return s.convertScryfallCard(scryfallCard), nil
}

// makeRequest makes HTTP request to Scryfall API
func (s *Service) makeRequest(url string) (*ScryfallCard, error) {
	resp, err := s.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("card not found")
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("scryfall API error: %d - %s", resp.StatusCode, string(body))
	}

	var scryfallCard ScryfallCard
	if err := json.NewDecoder(resp.Body).Decode(&scryfallCard); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &scryfallCard, nil
}

// convertScryfallCard converts Scryfall API response to internal Card model
func (s *Service) convertScryfallCard(sc *ScryfallCard) *models.Card {
	card := &models.Card{
		ID:              uuid.New().String(),
		ScryfallID:      sc.ID,
		Name:            sc.Name,
		SetCode:         strings.ToUpper(sc.SetCode),
		CollectorNumber: sc.CollectorNumber,
		OracleText:      sc.OracleText,
		TypeLine:        sc.TypeLine,
		ManaCost:        sc.ManaCost,
		Rarity:          sc.Rarity,
		CreatedAt:       time.Now(),
	}

	// Get image URI
	if sc.ImageURIs != nil {
		if normal, ok := sc.ImageURIs["normal"]; ok {
			card.ImageURI = normal
		} else if large, ok := sc.ImageURIs["large"]; ok {
			card.ImageURI = large
		}
	} else if len(sc.CardFaces) > 0 {
		// For double-faced cards, use the first face
		if faceImages, ok := sc.CardFaces[0]["image_uris"].(map[string]interface{}); ok {
			if normal, ok := faceImages["normal"].(string); ok {
				card.ImageURI = normal
			}
		}
	}

	return card
}

// rateLimit ensures we don't exceed Scryfall's rate limit
func (s *Service) rateLimit() {
	if !s.lastCall.IsZero() {
		elapsed := time.Since(s.lastCall)
		if elapsed < RateLimit {
			time.Sleep(RateLimit - elapsed)
		}
	}
	s.lastCall = time.Now()
}
