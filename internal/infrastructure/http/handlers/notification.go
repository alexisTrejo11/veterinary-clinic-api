package handlers

import (
	"clinic-vet-api/internal/core/notifications"
	"clinic-vet-api/internal/infrastructure/http/handlers/dtos"
	"clinic-vet-api/internal/infrastructure/http/handlers/mappers"
	"clinic-vet-api/internal/shared"
	"clinic-vet-api/internal/shared/errors"
	"clinic-vet-api/internal/shared/http"
	"clinic-vet-api/internal/shared/page"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// NotificationHandler exposes notification read (monitoring) and manual send (write).
// Customers read only their notifications; staff can read by type/channel/id/summary and send.
type NotificationHandler struct {
	service   notifications.NotificationService
	validator *validator.Validate
}

func NewNotificationHandler(service notifications.NotificationService, validator *validator.Validate) *NotificationHandler {
	return &NotificationHandler{service: service, validator: validator}
}

// ------------------------------------------------------------
// Customer read (my notifications)
// ------------------------------------------------------------

func (h *NotificationHandler) GetMyNotifications(c *gin.Context) {
	logic := func(ctx *gin.Context) (any, error) {
		userID, err := UserIDFromContext(ctx)
		if err != nil {
			return nil, err
		}
		var pageParams page.PaginationRequest
		if err := http.ShouldBindPageParams(&pageParams, ctx, h.validator); err != nil {
			return nil, err
		}
		return h.service.GetNotificationsByUserID(ctx.Request.Context(), userID, pageParams.ToPagination())
	}
	http.HandleRequestNoBodyWithResponder(h.validator, logic, func(c *gin.Context, res any) {
		if res == nil {
			http.Success(c, nil, "No notifications found")
			return
		}
		p := res.(page.Page[notifications.Notification])
		http.Paginated(c, &p, "Notifications")
	})(c)
}

func (h *NotificationHandler) GetMyNotificationByID(c *gin.Context) {
	logic := func(ctx *gin.Context) (any, error) {
		userID, err := UserIDFromContext(ctx)
		if err != nil {
			return nil, err
		}
		id := ctx.Param("id")
		if id == "" {
			return nil, shared.ErrEntityNotFound
		}
		notif, err := h.service.GetNotificationByID(ctx.Request.Context(), id)
		if err != nil {
			return nil, err
		}
		if notif.UserID != userID.String() {
			return nil, errors.ForbiddenError(ctx.Request.Context(), "get_notification", "notification", "notification does not belong to user")
		}
		return notif, nil
	}
	http.HandleGetRequest(h.validator, "Notification", logic)(c)
}

// ------------------------------------------------------------
// Staff read (monitoring: by type, channel, id, summary)
// ------------------------------------------------------------

func (h *NotificationHandler) GetNotificationsByType(c *gin.Context) {
	logic := func(ctx *gin.Context) (any, error) {
		notificationType := ctx.Param("type")
		if notificationType == "" {
			notificationType = ctx.Query("type")
		}
		if notificationType == "" {
			return nil, errors.MissingFieldError(ctx.Request.Context(), "type", "notification type is required", "get_notifications_by_type")
		}
		var pageParams page.PaginationRequest
		if err := http.ShouldBindPageParams(&pageParams, ctx, h.validator); err != nil {
			return nil, err
		}
		return h.service.GetNotificationBySpecification(ctx.Request.Context(), notificationType, pageParams.ToPagination())
	}
	http.HandleRequestNoBodyWithResponder(h.validator, logic, func(c *gin.Context, res any) {
		if res == nil {
			http.Success(c, nil, "No notifications found")
			return
		}
		p := res.(page.Page[notifications.Notification])
		http.Paginated(c, &p, "Notifications")
	})(c)
}

func (h *NotificationHandler) GetNotificationsByChannel(c *gin.Context) {
	logic := func(ctx *gin.Context) (any, error) {
		channel := ctx.Param("channel")
		if channel == "" {
			channel = ctx.Query("channel")
		}
		if channel == "" {
			return nil, errors.MissingFieldError(ctx.Request.Context(), "channel", "channel is required", "get_notifications_by_channel")
		}
		var pageParams page.PaginationRequest
		if err := http.ShouldBindPageParams(&pageParams, ctx, h.validator); err != nil {
			return nil, err
		}
		return h.service.GetNotificationByChannel(ctx.Request.Context(), channel, pageParams.ToPagination())
	}
	http.HandleRequestNoBodyWithResponder(h.validator, logic, func(c *gin.Context, res any) {
		if res == nil {
			http.Success(c, nil, "No notifications found")
			return
		}
		p := res.(page.Page[notifications.Notification])
		http.Paginated(c, &p, "Notifications")
	})(c)
}

func (h *NotificationHandler) GetNotificationByID(c *gin.Context) {
	logic := func(ctx *gin.Context) (any, error) {
		id := ctx.Param("id")
		if id == "" {
			return nil, shared.ErrEntityNotFound
		}
		return h.service.GetNotificationByID(ctx.Request.Context(), id)
	}
	http.HandleGetRequest(h.validator, "Notification", logic)(c)
}

func (h *NotificationHandler) GetNotificationSummary(c *gin.Context) {
	logic := func(ctx *gin.Context) (any, error) {
		senderID := ctx.Query("sender_id")
		if senderID == "" {
			senderID = ctx.Param("sender_id")
		}
		return h.service.GenerateSummary(ctx.Request.Context(), senderID)
	}
	http.HandleGetRequest(h.validator, "NotificationSummary", logic)(c)
}

// ------------------------------------------------------------
// Staff write (manual send only)
// ------------------------------------------------------------

func (h *NotificationHandler) SendNotification(c *gin.Context) {
	logic := func(ctx *gin.Context, req dtos.SendNotificationRequest) (any, error) {
		notif, err := mappers.SendNotificationRequestToNotification(req)
		if err != nil {
			return nil, err
		}
		err = h.service.Send(ctx.Request.Context(), notif)
		if err != nil {
			return nil, err
		}
		return nil, nil
	}
	http.HandleRequestWithResponder(h.validator, logic, func(c *gin.Context, _ any) {
		http.Success(c, nil, "Notification sent successfully")
	})(c)
}
