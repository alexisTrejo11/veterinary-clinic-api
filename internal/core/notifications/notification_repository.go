package notifications

import (
	"clinic-vet-api/internal/shared"
	"clinic-vet-api/internal/shared/page"
	"context"
)

// entity.NotificationRepository defines methods for persisting and tracking notifications
//go:generate mockgen -source=repository.go -destination=../../test/mock/notification_repository_mock.go -package=mock

type NotificationRepository interface {
	Create(ctx context.Context, notification *Notification) error
	GetByID(ctx context.Context, id string) (Notification, error)
	ListByUser(ctx context.Context, userID shared.UserID, pagination page.Pagination) (page.Page[Notification], error)
	ListBySpecies(ctx context.Context, notificationType string, pagination page.Pagination) (page.Page[Notification], error)
	ListByChannel(ctx context.Context, channel string, pagination page.Pagination) (page.Page[Notification], error)
}
