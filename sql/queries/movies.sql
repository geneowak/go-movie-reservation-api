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

-- name: CheckMovieById :one
SELECT
    EXISTS(
        SELECT
            1
        FROM
            movies
        WHERE
            id = $1
    );

-- name: GetShowingMovies :many
SELECT
    movies.*
FROM
    movies
    LEFT JOIN show_times ON show_times.movie_id = movies.id
WHERE
    show_times.start_date >= NOW()
    AND show_times.endtimes <= NOW()
GROUP BY
    movies.id
ORDER BY
    show_times.start_date;
