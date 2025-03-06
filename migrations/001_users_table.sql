-- Write your migrate up statements here

CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(12) PRIMARY KEY,
    email VARCHAR(128) UNIQUE NOT NULL,
    slug VARCHAR(32) UNIQUE NOT NULL,
    avatar_url VARCHAR(255),
    full_name VARCHAR(64) NOT NULL,
    bio VARCHAR(500),
    pass_hash VARCHAR(64) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_users_slug ON users(slug);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

---- create above / drop below ----

DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_users_slug;

DROP TABLE IF EXISTS users;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
