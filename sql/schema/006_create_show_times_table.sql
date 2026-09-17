-- +goose Up
CREATE TABLE show_times(
    id uuid PRIMARY KEY,
    start_time time NOT NULL,
    price int NOT NULL,
    description varchar(100),  -- premier, etc..
    price_currency varchar(5) NOT NULL DEFAULT 'UGX',
    movie_id uuid NOT NULL REFERENCES movies(id) ON DELETE CASCADE,
    cinema_id uuid NOT NULL REFERENCES cinemas(id) ON DELETE CASCADE,
    experience_type text NOT NULL,
    start_date timestamp NOT NULL,
    end_date timestamp NOT NULL,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL
);

-- +goose Down
DROP TABLE show_times;
