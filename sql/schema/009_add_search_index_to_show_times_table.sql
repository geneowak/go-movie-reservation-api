-- +goose Up
CREATE INDEX idx_show_times_movie_filter ON show_times(movie_id, end_date, start_date, start_time);

CREATE INDEX idx_movie_genre ON movies(genre);

-- +goose Down
DROP INDEX IF EXISTS idx_show_times_movie_filter;

DROP INDEX IF EXISTS idx_movie_genre;
