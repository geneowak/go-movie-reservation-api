-- +goose Up
CREATE TABLE refresh_tokens(
    token text PRIMARY KEY,
    user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    expires_at timestamp NOT NULL,
    revoked_at timestamp,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL
);
