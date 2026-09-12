-- +goose Up
CREATE TABLE reservations(
    id uuid PRIMARY KEY,
    show_time_id uuid NOT NULL REFERENCES show_times(id) ON DELETE CASCADE,
    user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    seat varchar(100) NOT NULL,
    STATUS varchar(100) NOT NULL DEFAULT 'available',
    reserved_at timestamp,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL
);

-- +goose Down
DROP TABLE reservations;
