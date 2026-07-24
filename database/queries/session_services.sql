-- name: FindSessionServiceByID :one
SELECT * FROM session_services WHERE id = $1;

-- name: FindSessionServicesBySessionID :many
SELECT * FROM session_services WHERE session_id = $1 ORDER BY created_at;

-- name: CreateSessionService :one
INSERT INTO session_services (session_id, service_catalog_id, quantity, price_applied, notes)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: UpdateSessionService :one
UPDATE session_services
SET session_id = $2, service_catalog_id = $3, quantity = $4, price_applied = $5, notes = $6
WHERE id = $1
RETURNING *;

-- name: DeleteSessionServiceByID :exec
DELETE FROM session_services WHERE id = $1;

-- name: DeleteSessionServicesBySessionID :exec
DELETE FROM session_services WHERE session_id = $1;
