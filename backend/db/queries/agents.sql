-- backend/db/queries/agents.sql

-- name: CreateAgent :one
INSERT INTO agents (name, description, system_prompt, model, tools, team_id, created_by)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetAgentByID :one
SELECT * FROM agents WHERE id = $1;

-- name: ListAgentsByTeam :many
SELECT * FROM agents
WHERE team_id = $1
ORDER BY created_at DESC;

-- name: UpdateAgent :one
UPDATE agents
SET name = $2, description = $3, system_prompt = $4, model = $5, tools = $6, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteAgent :exec
DELETE FROM agents WHERE id = $1;
