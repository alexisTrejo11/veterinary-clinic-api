package handlers

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func NewHomeHandler(router *gin.Engine) *HomeHandler {
	return &HomeHandler{
		router: router,
	}
}

type HomeHandler struct {
	router *gin.Engine
}

// Health check endpoint
// @Summary Health Check
// @Description Check if the API is running
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /health [get]
func (h *HomeHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "ok",
		"timestamp": time.Now().Unix(),
		"service":   "clinic-vet-api",
		"version":   getEnvWithDefault("APP_VERSION", "2.0.0"),
	})
}

func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
