-- name: GetRecordsByUserID :one
SELECT * FROM records
WHERE user_id = ?;

-- name: RegisterRecord :one
INSERT INTO records (id, user_id, created_at, updated_at, wins, losses, draws)
VALUES (
    ?,
    ?,
    ?,
    ?,
    ?,
    ?,
    ?
)
RETURNING *;

-- name: UpdateRecord :one
UPDATE records
SET
    updated_at = ?,
    wins = ?,
    losses = ?,
    draws = ?
WHERE user_id = ?
RETURNING *;

-- name: ResetRecords :exec
DELETE FROM records;