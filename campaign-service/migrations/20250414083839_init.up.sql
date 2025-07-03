CREATE TABLE IF NOT EXISTS campaigns (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    master_id INT,
    invite_code VARCHAR(255) DEFAULT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (name, master_id)
);

CREATE TABLE IF NOT EXISTS campaign_characters (
    char_id INT NOT NULL UNIQUE,
    player_id INT NOT NULL,
    campaign_id INT NOT NULL REFERENCES campaigns(id),
    PRIMARY KEY (campaign_id, player_id, char_id)
);