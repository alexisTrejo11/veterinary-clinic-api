package handlers

import (
	"clinic-vet-api/internal/core/notifications"

	"github.com/gin-gonic/gin"
)

type NotificationHandler struct {
	service notifications.NotificationService
}

func NewNotificationHandler(service notifications.NotificationService) *NotificationHandler {
	return &NotificationHandler{service: service}
}

func (c *NotificationHandler) SendNotification(ctx *gin.Context) {
}

func (c *NotificationHandler) GetNotificationByUserId(ctx *gin.Context) {
}

func (c *NotificationHandler) GetNotificationById(ctx *gin.Context) {
}

func (c *NotificationHandler) GetNotificationBySpecies(ctx *gin.Context) {
}

func (c *NotificationHandler) GetNotificationByChannel(ctx *gin.Context) {
}
