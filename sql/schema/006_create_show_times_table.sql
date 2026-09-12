-- +goose Up
CREATE TABLE show_times(
    id uuid PRIMARY KEY,
    time text NOT NULL,
    movie_id uuid NOT NULL REFERENCES movies(id) ON DELETE CASCADE,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL
);

-- +goose Down
DROP TABLE show_times;
