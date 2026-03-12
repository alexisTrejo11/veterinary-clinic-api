-- 000009_notification_table.up.sql
-- Notification table

CREATE TABLE IF NOT EXISTS notifications (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    title VARCHAR(255) NOT NULL,
    subject VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    token VARCHAR(255) NOT NULL,
    notification_type VARCHAR(255) NOT NULL,
    channel VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT chk_notification_type CHECK (notification_type IN ('alert', 'reminder', 'info', 'verification_token', 'activation_token', 'password_reset', 'appointment_confirm', 'appointment_remind', 'appointment_cancel', 'payment_receipt', 'payment_reminder', 'welcome', 'security', 'promotional', 'unknown')),
    CONSTRAINT chk_channel CHECK (channel IN ('email', 'sms', 'push', 'in_app', 'voice', 'whatsapp', 'unknown'))
);


CREATE INDEX IF NOT EXISTS idx_notifications_user_id ON notifications (user_id);
CREATE INDEX IF NOT EXISTS idx_notifications_notification_type ON notifications (notification_type);
CREATE INDEX IF NOT EXISTS idx_notifications_channel ON notifications (channel);
CREATE INDEX IF NOT EXISTS idx_notifications_created_at ON notifications (created_at);
CREATE INDEX IF NOT EXISTS idx_notifications_updated_at ON notifications (updated_at);