-- name: InsertAndReturnUser :one
INSERT INTO users (id, first_name, last_name, email) VALUES (?, ?, ?, ?) RETURNING *;

-- name: GetUser :one
SELECT * FROM users WHERE id = ?;
