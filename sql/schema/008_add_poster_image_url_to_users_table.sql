-- +goose Up
ALTER TABLE
    movies
ADD
    COLUMN poster_image_url text;

-- +goose Down
ALTER TABLE
    movies DROP COLUMN poster_image_url;
