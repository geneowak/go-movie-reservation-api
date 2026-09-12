-- +goose Up
CREATE TABLE cinemas(
    id uuid PRIMARY KEY,
    location_id uuid NOT NULL REFERENCES locations(id) ON DELETE CASCADE,
    name varchar(255) NOT NULL,
    experience_types text NOT NULL,
    seat_map jsonb NOT NULL,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL
);

-- +goose Down
DROP TABLE cinemas;
