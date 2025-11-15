-- Initial database schema for MTG Card Detector

-- Users table for anonymous authentication
CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY,
    device_id TEXT UNIQUE NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    last_seen DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_device_id ON users(device_id);

-- Cards table - MTG card master data
CREATE TABLE IF NOT EXISTS cards (
    id TEXT PRIMARY KEY,
    scryfall_id TEXT UNIQUE,
    name TEXT NOT NULL,
    set_code TEXT NOT NULL,
    collector_number TEXT NOT NULL,
    image_uri TEXT,
    oracle_text TEXT,
    type_line TEXT,
    mana_cost TEXT,
    rarity TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_cards_name ON cards(name);
CREATE INDEX idx_cards_set_collector ON cards(set_code, collector_number);
CREATE INDEX idx_cards_scryfall_id ON cards(scryfall_id);

-- Inventory table - user card ownership
CREATE TABLE IF NOT EXISTS inventory (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id TEXT NOT NULL,
    card_id TEXT NOT NULL,
    quantity INTEGER DEFAULT 1,
    added_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (card_id) REFERENCES cards(id) ON DELETE CASCADE,
    UNIQUE(user_id, card_id)
);

CREATE INDEX idx_inventory_user_id ON inventory(user_id);
CREATE INDEX idx_inventory_card_id ON inventory(card_id);

-- Scan sessions table - audit trail
CREATE TABLE IF NOT EXISTS scan_sessions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id TEXT NOT NULL,
    scan_type TEXT NOT NULL, -- 'single' or 'bulk'
    cards_scanned INTEGER DEFAULT 0,
    successful_scans INTEGER DEFAULT 0,
    failed_scans INTEGER DEFAULT 0,
    started_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    completed_at DATETIME,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_scan_sessions_user_id ON scan_sessions(user_id);
CREATE INDEX idx_scan_sessions_started_at ON scan_sessions(started_at);
