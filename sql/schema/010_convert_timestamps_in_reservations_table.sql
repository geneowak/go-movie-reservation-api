-- +goose Up
ALTER TABLE
    reservations
ALTER COLUMN
    reserved_at TYPE timestamptz USING reserved_at AT TIME ZONE 'UTC',
ALTER COLUMN
    created_at TYPE timestamptz USING created_at AT TIME ZONE 'UTC',
ALTER COLUMN
    updated_at TYPE timestamptz USING updated_at AT TIME ZONE 'UTC';

-- +goose Down
ALTER TABLE
    cinemas
ALTER COLUMN
    reserved_at TYPE timestamp,
ALTER COLUMN
    created_at TYPE timestamp,
ALTER COLUMN
    updated_at TYPE timestamp;
