CREATE TABLE IF NOT EXISTS campaigns (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    master_id INT REFERENCES users(id),
    invite_code VARCHAR(255) DEFAULT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (name, master_id)
);

CREATE TABLE IF NOT EXISTS players (
    campaign_id INT REFERENCES campaigns(id),
    player_id INT REFERENCES users(id),
    -- character_id INT REFERENCES characters(id),
    joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (campaign_id, player_id)
);