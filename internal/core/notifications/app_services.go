package notifications

import (
	"clinic-vet-api/internal/shared"
	"clinic-vet-api/internal/shared/page"
	"context"
	"fmt"
)

type NotificationService interface {
	Send(ctx context.Context, notif *Notification) error

	GetNotificationsByUserID(ctx context.Context, userID shared.UserID, pagination page.Pagination) (page.Page[Notification], error)
	GetNotificationBySpecies(ctx context.Context, notificationType string, pagination page.Pagination) (page.Page[Notification], error)
	GetNotificationByChannel(ctx context.Context, channel string, pagination page.Pagination) (page.Page[Notification], error)
	GenerateSummary(ctx context.Context, senderID string) ([]Notification, error)
	GetNotificationByID(ctx context.Context, id string) (Notification, error)
}

type notificationService struct {
	repository  NotificationRepository
	emailSender EmailSender
	smsSender   SMSender
}

func NewNotificationService(repository NotificationRepository, emailSender EmailSender, smsSender SMSender) NotificationService {
	return &notificationService{
		repository:  repository,
		emailSender: emailSender,
		smsSender:   smsSender,
	}
}

func (ns *notificationService) Send(ctx context.Context, notif *Notification) error {
	switch notif.Channel {
	case NotificationChannelEmail:
		return ns.emailSender.Send(ctx, notif)
	case NotificationChannelSMS:
		return ns.smsSender.Send(ctx, notif)
	default:
		return fmt.Errorf("unsupported notification channel: %s", notif.Channel)
	}
}

func (s *notificationService) GetNotificationsByUserID(ctx context.Context, userID shared.UserID, pagination page.Pagination) (page.Page[Notification], error) {
	notificationPage, err := s.repository.FindByUser(ctx, userID, pagination)
	if err != nil {
		return page.Page[Notification]{}, err
	}

	return notificationPage, nil
}

func (s *notificationService) GetNotificationByID(ctx context.Context, id string) (Notification, error) {
	return s.repository.FindByID(ctx, id)
}

func (s *notificationService) GetNotificationBySpecies(ctx context.Context, notificationType string, pagination page.Pagination) (page.Page[Notification], error) {
	notificationPage, err := s.repository.FindBySpecies(ctx, notificationType, pagination)
	if err != nil {
		return page.Page[Notification]{}, err
	}

	return notificationPage, nil
}

func (s *notificationService) GetNotificationByChannel(ctx context.Context, channel string, pagination page.Pagination) (page.Page[Notification], error) {
	notificationPage, err := s.repository.FindByChannel(ctx, channel, pagination)
	if err != nil {
		return page.Page[Notification]{}, err
	}

	return notificationPage, nil
}

func (s *notificationService) GenerateSummary(ctx context.Context, senderID string) ([]Notification, error) {
	return []Notification{}, nil
}
