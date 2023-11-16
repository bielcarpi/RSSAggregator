ALTER TABLE users ADD COLUMN api_key VARCHAR(64) NOT NULL UNIQUE DEFAULT (
    encode(sha256(random()::text::bytea), 'hex')
);