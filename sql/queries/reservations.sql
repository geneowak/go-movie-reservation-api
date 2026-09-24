-- name: CreateReservation :one
INSERT INTO
    reservations(
        id,
        show_time_id,
        user_id,
        seat_no,
        reserved_at,
        created_at,
        updated_at
    )
VALUES
    (uuidv7(), $1, $2, $3, NOW(), NOW(), NOW())
RETURNING
    *;

-- name: GetReservationBySeatNo :one
SELECT
    *
FROM
    reservations
WHERE
    show_time_id = $1
    AND seat_no = $2
LIMIT
    1;
