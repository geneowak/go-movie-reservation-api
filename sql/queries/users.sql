-- name: CreateUser :one
INSERT INTO
    users(
        id,
        email,
        hashed_password,
        created_at,
        updated_at
    )
VALUES
    (uuidv7(), $1, $2, NOW(), NOW())
RETURNING
    *;

-- name: GetUserByEmail :one
SELECT
    *
FROM
    users
WHERE
    email = $1
LIMIT
    1;

-- name: GetUserById :one
SELECT
    *
FROM
    users
WHERE
    id = $1
LIMIT
    1;

-- name: CheckUserId :one
SELECT
    EXISTS(
        SELECT
            1
        FROM
            users
        WHERE
            id = $1
    );

-- name: UpdateUserAdminStatus :one
UPDATE
    users
SET
    is_admin = $1
WHERE
    id = $2
RETURNING
    *;
