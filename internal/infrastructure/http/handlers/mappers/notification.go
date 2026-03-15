package mappers

import (
	"clinic-vet-api/internal/core/notifications"
	"clinic-vet-api/internal/infrastructure/http/handlers/dtos"
)

// SendNotificationRequestToNotification builds a Notification from the request for manual send.
func SendNotificationRequestToNotification(req dtos.SendNotificationRequest) (*notifications.Notification, error) {
	channel, err := notifications.ParseNotificationChannel(req.Channel)
	if err != nil {
		return nil, err
	}
	nType, err := notifications.ParseNotificationType(req.Type)
	if err != nil {
		return nil, err
	}
	return &notifications.Notification{
		UserID:    req.UserID,
		UserEmail: req.UserEmail,
		UserPhone: req.UserPhone,
		Title:     req.Title,
		Subject:   req.Subject,
		Message:   req.Message,
		Token:     req.Token,
		NType:     nType,
		Channel:   channel,
	}, nil
}
