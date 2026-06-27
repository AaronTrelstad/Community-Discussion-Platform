-- name: CreateCommunity :one
INSERT INTO communities (name, title, description, created_by)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetCommunityByName :one
SELECT * FROM communities WHERE name = $1;

-- name: GetCommunityByID :one
SELECT * FROM communities WHERE id = $1;

-- name: ListCommunities :many
SELECT * FROM communities
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: UpdateCommunity :one
UPDATE communities
SET title = $2, description = $3
WHERE id = $1
RETURNING *;

-- name: DeleteCommunity :exec
DELETE FROM communities WHERE id = $1;

-- name: JoinCommunity :exec
INSERT INTO community_members (user_id, community_id, role)
VALUES ($1, $2, 'member')
ON CONFLICT (user_id, community_id) DO NOTHING;

-- name: LeaveCommunity :exec
DELETE FROM community_members
WHERE user_id = $1 AND community_id = $2;

-- name: GetCommunityMember :one
SELECT * FROM community_members
WHERE user_id = $1 AND community_id = $2;

-- name: ListCommunityMembers :many
SELECT u.id, u.username, u.avatar_url, cm.role, cm.joined_at
FROM community_members cm
JOIN users u ON u.id = cm.user_id
WHERE cm.community_id = $1
ORDER BY cm.joined_at DESC
LIMIT $2 OFFSET $3;

-- name: GetUserCommunities :many
SELECT c.* FROM communities c
JOIN community_members cm ON cm.community_id = c.id
WHERE cm.user_id = $1
ORDER BY cm.joined_at DESC;

-- name: GetCommunityMemberCount :one
SELECT COUNT(*) FROM community_members WHERE community_id = $1;

-- name: IsUserMember :one
SELECT EXISTS (
    SELECT 1 FROM community_members
    WHERE user_id = $1 AND community_id = $2
) AS is_member;
