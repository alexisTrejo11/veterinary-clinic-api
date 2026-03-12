package notifications

import (
	"context"
)

type EmailSender interface {
	Send(ctx context.Context, notification *Notification) error
}

type SMSender interface {
	Send(ctx context.Context, notification *Notification) error
}
