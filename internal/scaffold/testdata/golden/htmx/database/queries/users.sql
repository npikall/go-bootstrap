-- name: ListUsers :many
SELECT id, name FROM users ORDER BY id;

-- name: CreateUser :exec
INSERT INTO users (name) VALUES ($1);
