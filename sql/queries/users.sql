-- name: GetUserByName :one
SELECT * FROM users
WHERE username = ?;
--

-- name: RegisterUser :one
INSERT INTO users (id, username, created_at, updated_at, hashed_password)
VALUES (
    ?,
    ?,
    ?,
    ?,
    ?
)
RETURNING *;
--

-- name: ResetUsers :exec
DELETE FROM users;
--