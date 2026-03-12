package notifications

import (
	"time"
)

type Notification struct {
	ID        string
	UserID    string
	UserEmail string
	UserPhone string
	Title     string
	Subject   string
	Message   string
	Token     string
	CreatedAt time.Time
	NType     NotificationType
	Channel   NotificationChannel
}

func NewActivateAccountNotification(userID, userEmail string) Notification {
	return Notification{
		UserID:    userID,
		UserEmail: userEmail,
		Subject:   "Activate your account",
		Message:   "Please click the link below to activate your account.",
		NType:     NotificationTypeActivationToken,
		Channel:   NotificationChannelEmail,
	}
}
