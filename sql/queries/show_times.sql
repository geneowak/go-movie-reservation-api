-- name: CreateShowTime :one
INSERT INTO
    show_times(
        id,
        start_time,
        price,
        movie_id,
        cinema_id,
        experience_type,
        start_date,
        end_date,
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

-- name: GetShowTimeDetails :one
SELECT
    show_times.*,
    to_jsonb(c.*) AS cinema
FROM
    show_times
    INNER JOIN cinemas c ON c.id = show_times.cinema_id
WHERE
    show_times.id = $1;
