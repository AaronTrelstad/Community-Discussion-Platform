-- backend/db/queries/teams.sql

-- name: CreateTeam :one
INSERT INTO teams (name, title, description, created_by)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetTeamByName :one
SELECT * FROM teams WHERE name = $1;

-- name: GetTeamByID :one
SELECT * FROM teams WHERE id = $1;

-- name: ListTeams :many
SELECT * FROM teams
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: UpdateTeam :one
UPDATE teams
SET title = $2, description = $3
WHERE id = $1
RETURNING *;

-- name: DeleteTeam :exec
DELETE FROM teams WHERE id = $1;

-- name: JoinTeam :exec
INSERT INTO team_members (user_id, team_id, role)
VALUES ($1, $2, 'member')
ON CONFLICT (user_id, team_id) DO NOTHING;

-- name: LeaveTeam :exec
DELETE FROM team_members
WHERE user_id = $1 AND team_id = $2;

-- name: GetTeamMember :one
SELECT * FROM team_members
WHERE user_id = $1 AND team_id = $2;

-- name: ListTeamMembers :many
SELECT u.id, u.username, u.avatar_url, tm.role, tm.joined_at
FROM team_members tm
JOIN users u ON u.id = tm.user_id
WHERE tm.team_id = $1
ORDER BY tm.joined_at DESC
LIMIT $2 OFFSET $3;

-- name: GetUserTeams :many
SELECT t.* FROM teams t
JOIN team_members tm ON tm.team_id = t.id
WHERE tm.user_id = $1
ORDER BY tm.joined_at DESC;

-- name: GetTeamMemberCount :one
SELECT COUNT(*) FROM team_members WHERE team_id = $1;

-- name: IsUserMember :one
SELECT EXISTS (
    SELECT 1 FROM team_members
    WHERE user_id = $1 AND team_id = $2
) AS is_member;
