-- name: CreateMember :one
INSERT INTO members (email, data)
VALUES ($1, $2)
RETURNING *;

-- name: CreateLocalAuth :exec
INSERT INTO local_auth (member_id, password_hash)
VALUES ($1, $2);

-- name: GetMemberByEmail :one
SELECT * FROM members
WHERE email = $1 LIMIT 1;

-- name: GetMemberById :one
SELECT * FROM members
WHERE id = $1 LIMIT 1;

-- name: GetMemberWithPassword :one
SELECT 
    m.*, 
    la.password_hash 
FROM members m
INNER JOIN local_auth la ON m.id = la.member_id
WHERE m.email = $1 LIMIT 1;
