-- Write your migrate up statements here

CREATE TYPE chat_role AS ENUM ('admin', 'moderator', 'member');

CREATE TABLE IF NOT EXISTS chats_members (
    chat_id VARCHAR(12) NOT NULL REFERENCES chats(id),
    user_id VARCHAR(12) NOT NULL REFERENCES users(id),
    role chat_role NOT NULL DEFAULT 'member',
    joined_at TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (chat_id, user_id)
);

---- create above / drop below ----

DROP TABLE IF EXISTS chats_members;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
