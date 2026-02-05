-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE members (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	email TEXT UNIQUE NOT NULL,
	data JSONB NOT NULL DEFAULT '{}'::jsonb,
	email_verified_at TIMESTAMPTZ,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE local_auth (
	member_id UUID PRIMARY KEY REFERENCES members(id) ON DELETE CASCADE,
	password_HASH TEXT NOT NULL
);

CREATE TABLE sessions (
	token TEXT PRIMARY KEY,
	data BYTEA NOT NULL,
	expiry TIMESTAMPTZ NOT NULL
);

CREATE INDEX sessions_expiry_idx ON sessions (expiry);
CREATE INDEX members_data_gin ON members USING gin (data);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP INDEX members_data_gin;
DROP INDEX sessions_expiry_idx;
DROP TABLE sessions;
DROP TABLE local_auth;
DROP TABLE members;
-- +goose StatementEnd
