-- backend/db/queries/events.sql

-- name: CreateEvent :one
INSERT INTO events (run_id, sequence, type, payload)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: ListEventsByRun :many
SELECT * FROM events
WHERE run_id = $1
ORDER BY sequence ASC;

-- name: GetLatestSequence :one
SELECT COALESCE(MAX(sequence), 0) FROM events WHERE run_id = $1;
