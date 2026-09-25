-- name: CreateLocation :one
INSERT INTO
    locations(
        id,
        name,
        address,
        google_map_url,
        created_at,
        updated_at
    )
VALUES
    (uuidv7(), $1, $2, $3, NOW(), NOW())
RETURNING
    *;

-- name: GetLocations :many
SELECT
    *
FROM
    locations;
