-- name: FindSessionAttachmentByID :one
SELECT * FROM session_attachments WHERE id = $1;

-- name: FindSessionAttachmentsBySessionID :many
SELECT * FROM session_attachments WHERE session_id = $1 ORDER BY created_at;

-- name: CreateSessionAttachment :one
INSERT INTO session_attachments (session_id, file_type, file_url, description, uploaded_by)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: UpdateSessionAttachment :one
UPDATE session_attachments
SET session_id = $2, file_type = $3, file_url = $4, description = $5, uploaded_by = $6
WHERE id = $1
RETURNING *;

-- name: DeleteSessionAttachmentByID :exec
DELETE FROM session_attachments WHERE id = $1;

-- name: DeleteSessionAttachmentsBySessionID :exec
DELETE FROM session_attachments WHERE session_id = $1;
