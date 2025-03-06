-- Write your migrate up statements here

CREATE TABLE IF NOT EXISTS chats (
    id VARCHAR(12) PRIMARY KEY,
    slug VARCHAR(12) UNIQUE NOT NULL,
    name VARCHAR(64) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_chat_slug ON chats(slug);

---- create above / drop below ----

DROP INDEX IF EXISTS idx_chat_slug;

DROP TABLE IF EXISTS chats;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
