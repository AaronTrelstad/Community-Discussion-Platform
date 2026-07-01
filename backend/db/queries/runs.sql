-- backend/db/queries/runs.sql

-- name: CreateRun :one
INSERT INTO runs (agent_id, team_id, started_by, input)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetRunByID :one
SELECT * FROM runs WHERE id = $1;

-- name: ListRunsByAgent :many
SELECT * FROM runs
WHERE agent_id = $1
ORDER BY started_at DESC
LIMIT $2 OFFSET $3;

-- name: ListRunsByTeam :many
SELECT * FROM runs
WHERE team_id = $1
ORDER BY started_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdateRunStatus :one
UPDATE runs
SET status = $2, output = $3, completed_at = NOW()
WHERE id = $1
RETURNING *;
