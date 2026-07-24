-- name: FindSessionSurgeryByID :one
SELECT * FROM session_surgeries WHERE id = $1;

-- name: FindSessionSurgeriesBySessionID :many
SELECT * FROM session_surgeries WHERE session_id = $1 ORDER BY created_at;

-- name: CreateSessionSurgery :one
INSERT INTO session_surgeries (
    session_id, procedure_name, anesthesia_type, anesthesia_agent,
    pre_op_notes, intra_op_notes, post_op_notes, duration_minutes, outcome, surgeon_id
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING *;

-- name: UpdateSessionSurgery :one
UPDATE session_surgeries
SET session_id = $2, procedure_name = $3, anesthesia_type = $4, anesthesia_agent = $5,
    pre_op_notes = $6, intra_op_notes = $7, post_op_notes = $8, duration_minutes = $9, outcome = $10, surgeon_id = $11
WHERE id = $1
RETURNING *;

-- name: DeleteSessionSurgeryByID :exec
DELETE FROM session_surgeries WHERE id = $1;

-- name: DeleteSessionSurgeriesBySessionID :exec
DELETE FROM session_surgeries WHERE session_id = $1;
