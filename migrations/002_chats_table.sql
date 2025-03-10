-- Write your migrate up statements here

CREATE TYPE chat_type AS ENUM ('private', 'group');

CREATE TABLE IF NOT EXISTS chats (
    id VARCHAR(12) PRIMARY KEY,
    slug VARCHAR(12) UNIQUE NOT NULL,
    name VARCHAR(64) NOT NULL,
    type chat_type NOT NULL DEFAULT 'private',
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_chat_slug ON chats(slug);

---- create above / drop below ----

DROP INDEX IF EXISTS idx_chat_slug;

DROP TABLE IF EXISTS chats;

DROP TYPE IF EXISTS chat_type;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
