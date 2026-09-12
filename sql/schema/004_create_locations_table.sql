-- +goose Up
CREATE TABLE locations(
    id uuid PRIMARY KEY,
    name varchar(255) NOT NULL,
    address text NOT NULL,
    google_map_url text NOT NULL,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL
);

-- +goose Down
DROP TABLE locations;
