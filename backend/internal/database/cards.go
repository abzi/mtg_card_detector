package database

import (
	"database/sql"
	"fmt"

	"github.com/abzi/mtg_card_detector/internal/models"
)

// CreateCard inserts a new card into the database
func (db *DB) CreateCard(card *models.Card) error {
	query := `INSERT INTO cards (id, scryfall_id, name, set_code, collector_number, image_uri, oracle_text, type_line, mana_cost, rarity, created_at)
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := db.Exec(query, card.ID, card.ScryfallID, card.Name, card.SetCode, card.CollectorNumber,
		card.ImageURI, card.OracleText, card.TypeLine, card.ManaCost, card.Rarity, card.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create card: %w", err)
	}
	return nil
}

// GetCardByID retrieves a card by ID
func (db *DB) GetCardByID(id string) (*models.Card, error) {
	query := `SELECT id, scryfall_id, name, set_code, collector_number, image_uri, oracle_text, type_line, mana_cost, rarity, created_at
	          FROM cards WHERE id = ?`
	card := &models.Card{}
	err := db.QueryRow(query, id).Scan(&card.ID, &card.ScryfallID, &card.Name, &card.SetCode, &card.CollectorNumber,
		&card.ImageURI, &card.OracleText, &card.TypeLine, &card.ManaCost, &card.Rarity, &card.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get card: %w", err)
	}
	return card, nil
}

// GetCardBySetAndNumber retrieves a card by set code and collector number
func (db *DB) GetCardBySetAndNumber(setCode, collectorNumber string) (*models.Card, error) {
	query := `SELECT id, scryfall_id, name, set_code, collector_number, image_uri, oracle_text, type_line, mana_cost, rarity, created_at
	          FROM cards WHERE set_code = ? AND collector_number = ?`
	card := &models.Card{}
	err := db.QueryRow(query, setCode, collectorNumber).Scan(&card.ID, &card.ScryfallID, &card.Name, &card.SetCode,
		&card.CollectorNumber, &card.ImageURI, &card.OracleText, &card.TypeLine, &card.ManaCost, &card.Rarity, &card.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get card by set and number: %w", err)
	}
	return card, nil
}

// GetCardByScryfallID retrieves a card by Scryfall ID
func (db *DB) GetCardByScryfallID(scryfallID string) (*models.Card, error) {
	query := `SELECT id, scryfall_id, name, set_code, collector_number, image_uri, oracle_text, type_line, mana_cost, rarity, created_at
	          FROM cards WHERE scryfall_id = ?`
	card := &models.Card{}
	err := db.QueryRow(query, scryfallID).Scan(&card.ID, &card.ScryfallID, &card.Name, &card.SetCode,
		&card.CollectorNumber, &card.ImageURI, &card.OracleText, &card.TypeLine, &card.ManaCost, &card.Rarity, &card.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get card by scryfall ID: %w", err)
	}
	return card, nil
}

// SearchCardsByName searches for cards by name (partial match)
func (db *DB) SearchCardsByName(name string, limit int) ([]models.Card, error) {
	query := `SELECT id, scryfall_id, name, set_code, collector_number, image_uri, oracle_text, type_line, mana_cost, rarity, created_at
	          FROM cards WHERE name LIKE ? ORDER BY name LIMIT ?`
	rows, err := db.Query(query, "%"+name+"%", limit)
	if err != nil {
		return nil, fmt.Errorf("failed to search cards: %w", err)
	}
	defer rows.Close()

	var cards []models.Card
	for rows.Next() {
		var card models.Card
		err := rows.Scan(&card.ID, &card.ScryfallID, &card.Name, &card.SetCode, &card.CollectorNumber,
			&card.ImageURI, &card.OracleText, &card.TypeLine, &card.ManaCost, &card.Rarity, &card.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan card: %w", err)
		}
		cards = append(cards, card)
	}

	return cards, nil
}
