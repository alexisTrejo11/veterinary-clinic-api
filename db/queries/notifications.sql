-- name: CreateNotification :one
INSERT INTO notifications (
    user_id, title, subject, message, token, notification_type, channel
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
)
RETURNING id, user_id, title, subject, message, token, notification_type, channel,
          created_at, updated_at;

-- name: GetNotificationByID :one
SELECT id, user_id, title, subject, message, token, notification_type, channel,
       created_at, updated_at
FROM notifications
WHERE id = $1;

-- name: FindNotificationsByUserID :many
SELECT id, user_id, title, subject, message, token, notification_type, channel,
       created_at, updated_at
FROM notifications
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: CountNotificationsByUserID :one
SELECT COUNT(*) FROM notifications
WHERE user_id = $1;

-- name: FindNotificationsByType :many
SELECT id, user_id, title, subject, message, token, notification_type, channel,
       created_at, updated_at
FROM notifications
WHERE notification_type = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: CountNotificationsByType :one
SELECT COUNT(*) FROM notifications
WHERE notification_type = $1;

-- name: FindNotificationsByChannel :many
SELECT id, user_id, title, subject, message, token, notification_type, channel,
       created_at, updated_at
FROM notifications
WHERE channel = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: CountNotificationsByChannel :one
SELECT COUNT(*) FROM notifications
WHERE channel = $1;
