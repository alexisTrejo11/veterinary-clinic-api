package http

import (
	"clinic-vet-api/internal/core/addresses"
	"clinic-vet-api/internal/infrastructure/http/handlers"
	"clinic-vet-api/internal/middleware"
	"clinic-vet-api/sqlc"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// ─── Dependencies ────────────────────────────────────────────────────────────

type APIConfig struct {
	Router         *gin.RouterGroup
	Queries        *sqlc.Queries
	Validator      *validator.Validate
	AuthMiddleware *middleware.AuthMiddleware
}

type ApplicationRepositories struct {
	addressRepository addresses.AddressRepository
}

type ApplicationServices struct {
	addressService addresses.AddressService
}

type APIHandlers struct {
	addressHandler handlers.AddressHandler
}
