-- +goose Up
CREATE TABLE movies(
    id uuid PRIMARY KEY,
    name varchar(255) NOT NULL,
    description text NOT NULL,
    duration_in_mins int NOT NULL,
    trailer_url text NOT NULL,
    genre text NOT NULL,
    pg_rating varchar(255) NOT NULL,
    experience_types text NOT NULL,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL
);

-- +goose Down
DROP TABLE movies;
