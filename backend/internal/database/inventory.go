package database

import (
	"database/sql"
	"fmt"

	"github.com/abzi/mtg_card_detector/internal/models"
)

// AddToInventory adds a card to user's inventory or increments quantity
func (db *DB) AddToInventory(userID, cardID string, quantity int) error {
	query := `INSERT INTO inventory (user_id, card_id, quantity)
	          VALUES (?, ?, ?)
	          ON CONFLICT(user_id, card_id)
	          DO UPDATE SET quantity = quantity + ?`
	_, err := db.Exec(query, userID, cardID, quantity, quantity)
	if err != nil {
		return fmt.Errorf("failed to add to inventory: %w", err)
	}
	return nil
}

// GetUserInventory retrieves all cards in user's inventory
func (db *DB) GetUserInventory(userID string) ([]models.InventoryItem, error) {
	query := `SELECT i.id, i.user_id, i.card_id, i.quantity, i.added_at,
	                 c.id, c.scryfall_id, c.name, c.set_code, c.collector_number,
	                 c.image_uri, c.oracle_text, c.type_line, c.mana_cost, c.rarity, c.created_at
	          FROM inventory i
	          JOIN cards c ON i.card_id = c.id
	          WHERE i.user_id = ?
	          ORDER BY i.added_at DESC`

	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get inventory: %w", err)
	}
	defer rows.Close()

	var items []models.InventoryItem
	for rows.Next() {
		var item models.InventoryItem
		item.Card = &models.Card{}
		err := rows.Scan(&item.ID, &item.UserID, &item.CardID, &item.Quantity, &item.AddedAt,
			&item.Card.ID, &item.Card.ScryfallID, &item.Card.Name, &item.Card.SetCode, &item.Card.CollectorNumber,
			&item.Card.ImageURI, &item.Card.OracleText, &item.Card.TypeLine, &item.Card.ManaCost,
			&item.Card.Rarity, &item.Card.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan inventory item: %w", err)
		}
		items = append(items, item)
	}

	return items, nil
}

// RemoveFromInventory removes a card from inventory or decrements quantity
func (db *DB) RemoveFromInventory(userID, cardID string, quantity int) error {
	// First check current quantity
	var currentQty int
	err := db.QueryRow(`SELECT quantity FROM inventory WHERE user_id = ? AND card_id = ?`, userID, cardID).Scan(&currentQty)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("card not found in inventory")
		}
		return fmt.Errorf("failed to check inventory: %w", err)
	}

	if currentQty <= quantity {
		// Remove completely
		_, err = db.Exec(`DELETE FROM inventory WHERE user_id = ? AND card_id = ?`, userID, cardID)
	} else {
		// Decrement quantity
		_, err = db.Exec(`UPDATE inventory SET quantity = quantity - ? WHERE user_id = ? AND card_id = ?`, quantity, userID, cardID)
	}

	if err != nil {
		return fmt.Errorf("failed to remove from inventory: %w", err)
	}
	return nil
}

// GetInventoryCount returns total number of cards in user's inventory
func (db *DB) GetInventoryCount(userID string) (int, error) {
	var count int
	err := db.QueryRow(`SELECT COALESCE(SUM(quantity), 0) FROM inventory WHERE user_id = ?`, userID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get inventory count: %w", err)
	}
	return count, nil
}
