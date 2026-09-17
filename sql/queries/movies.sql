-- name: CreateMovie :one
INSERT INTO
    movies(
        id,
        name,
        description,
        duration_in_mins,
        trailer_url,
        genre,
        pg_rating,
        experience_types,
        created_at,
        updated_at
    )
VALUES
    (
        uuidv7(),
        $1,
        $2,
        $3,
        $4,
        $5,
        $6,
        $7,
        NOW(),
        NOW()
    )
RETURNING
    *;

-- name: GetMovieById :one
SELECT
    *
FROM
    movies
WHERE
    id = $1;
