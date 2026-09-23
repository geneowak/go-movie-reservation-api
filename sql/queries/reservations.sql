-- name: CreateReservation :one
INSERT INTO
    reservations(
        id,
        show_time_id,
        user_id,
        seat,
        reserved_at,
        created_at,
        updated_at
    )
VALUES
(uuidv7(), $1, $2, $3, NOW(), NOW(), NOW())
RETURNING
    *;
