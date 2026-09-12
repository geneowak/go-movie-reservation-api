-- +goose Up
CREATE TABLE users(
    id uuid PRIMARY KEY,
    email varchar(255) NOT NULL,
    is_admin bool NOT NULL DEFAULT false,
    hashed_password text NOT NULL,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL
);

-- +goose Down
DROP TABLE users;
