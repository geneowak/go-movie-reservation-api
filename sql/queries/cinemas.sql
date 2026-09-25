-- name: CreateCinema :one
INSERT INTO
    cinemas(
        id,
        name,
        location_id,
        experience_types,
        seat_map,
        created_at,
        updated_at
    )
VALUES
    (uuidv7(), $1, $2, $3, $4, NOW(), NOW())
RETURNING
    *;

-- name: CheckCinemaById :one
SELECT
    EXISTS(
        SELECT
            1
        FROM
            cinemas
        WHERE
            id = $1
    );

-- name: UpdateCinemaDetails :one
UPDATE
    cinemas
SET
    name = $1,
    experience_types = $2,
    seat_map = $3
WHERE
    id = $4
RETURNING
    *;
