-- Write your migrate up statements here

CREATE TYPE attachments_type AS ENUM ('image', 'video', 'file', 'archive', 'archive');

CREATE TABLE IF NOT EXISTS attachments (
    id VARCHAR(12) PRIMARY KEY,
    chat_id VARCHAR(12) NOT NULL REFERENCES chats(id),
    user_id VARCHAR(12) NOT NULL REFERENCES users(id),
    name VARCHAR(128) NOT NULL,
    type attachments_type NOT NULL,
    url VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL  DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_attachments_chat_id ON attachments(chat_id);
CREATE INDEX IF NOT EXISTS idx_attachments_user_id ON attachments(user_id);

---- create above / drop below ----

DROP INDEX IF EXISTS idx_attachments_chat_id;
DROP INDEX IF EXISTS idx_attachments_user_id;

DROP TABLE IF EXISTS attachments;

DROP TYPE IF EXISTS attachments_type;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
