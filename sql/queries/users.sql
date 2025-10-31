-- name: CreateUser :exec
INSERT INTO users (email, hwid) VALUES (?, ?);

-- name: GetUserByEmail :one
SELECT
    id,
    email,
    hwid,
    created_at,
    updated_at
FROM users
WHERE email = ?;

-- name: GetUserByHWID :one
SELECT
    id,
    email,
    hwid,
    created_at,
    updated_at
FROM users
WHERE hwid = ?;

-- name: UpdateUserHWID :exec
UPDATE users
SET hwid = ?
WHERE email = ?;

-- name: ResetUserHWID :exec
UPDATE users
SET hwid = NULL
WHERE email = ?;
