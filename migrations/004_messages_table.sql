-- Write your migrate up statements here

CREATE TABLE IF NOT EXISTS messages (
    id VARCHAR(12) PRIMARY KEY,
    chat_id VARCHAR(12) NOT NULL REFERENCES chats(id),
    user_id VARCHAR(12) NOT NULL REFERENCES users(id),
    content TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_messages_chat_id ON messages(chat_id);
CREATE INDEX idx_messages_user_id ON messages(user_id);
CREATE INDEX idx_messages_chat_id_created_at ON messages(chat_id, created_at);

---- create above / drop below ----

DROP INDEX IF EXISTS idx_messages_chat_id_created_at;
DROP INDEX IF EXISTS idx_messages_chat_id;
DROP INDEX IF EXISTS idx_messages_user_id;

DROP TABLE IF EXISTS messages;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
