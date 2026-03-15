package http

import (
	"clinic-vet-api/internal/infrastructure/http/handlers"
	"clinic-vet-api/internal/middleware"
	"errors"

	"github.com/gin-gonic/gin"
)

type APIRouter struct {
	routerGroup    *gin.RouterGroup
	appHandlers    *appHandlers
	authMiddleware *middleware.AuthMiddleware
}

type appHandlers struct {
	auth    *handlers.AuthHandler
	user    *handlers.UserHandler
	pets    *handlers.PetHandler
	profile *handlers.ProfileHandler
}

func (r *APIRouter) Validate() error {
	if r.appHandlers.auth == nil {
		return errors.New("auth handler is required")
	}
	if r.authMiddleware == nil {
		return errors.New("auth middleware is required")
	}
	if r.routerGroup == nil {
		return errors.New("router group is required")
	}
	return nil
}

func NewAPIRouter(
	appHandlers *appHandlers,
	authMiddleware *middleware.AuthMiddleware,
	routerGroup *gin.RouterGroup,
) (*APIRouter, error) {
	router := &APIRouter{
		appHandlers:    appHandlers,
		authMiddleware: authMiddleware,
		routerGroup:    routerGroup,
	}

	if err := router.Validate(); err != nil {
		return nil, err
	}
	return router, nil
}

// ------------------------------------------------------------
// Auth Routes
// ------------------------------------------------------------

func (r *APIRouter) authRoutes() {

	publicAuthRoutes := r.routerGroup.Group("/auth")
	{
		publicAuthRoutes.POST("/register", r.appHandlers.auth.Register)
		publicAuthRoutes.POST("/login", r.appHandlers.auth.Login)
		publicAuthRoutes.POST("/activate", r.appHandlers.auth.ActivateAccount)
	}

	authenticatedAuthRoutes := r.routerGroup.Group("/auth")
	authenticatedAuthRoutes.Use(r.authMiddleware.Authenticate())
	{
		authenticatedAuthRoutes.POST("/logout", r.appHandlers.auth.Logout)
		authenticatedAuthRoutes.POST("/logout-all", r.appHandlers.auth.LogoutAll)
		authenticatedAuthRoutes.POST("/refresh", r.appHandlers.auth.RefreshToken)
	}

	twoFactorAuthRoutes := r.routerGroup.Group("/auth/2fa")
	twoFactorAuthRoutes.Use(r.authMiddleware.Authenticate())
	{
		twoFactorAuthRoutes.POST("/verify", r.appHandlers.auth.VerifyTwoFactor)
		twoFactorAuthRoutes.POST("/enable", r.appHandlers.auth.EnableTwoFactor)
		twoFactorAuthRoutes.POST("/disable", r.appHandlers.auth.DisableTwoFactor)
	}

	resetPasswordAuthRoutes := r.routerGroup.Group("/auth/reset-password")
	resetPasswordAuthRoutes.Use(r.authMiddleware.Authenticate())
	{
		resetPasswordAuthRoutes.POST("/request", r.appHandlers.auth.RequestResetPassword)
		resetPasswordAuthRoutes.POST("/reset", r.appHandlers.auth.ResetPassword)
	}
}

// ------------------------------------------------------------
// Profile Routes
// ------------------------------------------------------------

func (r *APIRouter) profileRoutes() {
	profileRoutes := r.routerGroup.Group("/profile")
	profileRoutes.Use(r.authMiddleware.Authenticate())
	{
		profileRoutes.GET("/", r.appHandlers.profile.GetProfile)
		profileRoutes.PUT("/", r.appHandlers.profile.UpdateProfile)
	}
}

// ------------------------------------------------------------
// User Routes
// ------------------------------------------------------------

func (r *APIRouter) userRoutes() {
	userRoutes := r.routerGroup.Group("/users")
	userRoutes.Use(r.authMiddleware.Authenticate())
	userRoutes.Use(r.authMiddleware.RequireAnyRole("admin", "manager"))
	{
		userRoutes.GET("/", r.appHandlers.user.SearchUsers)
		userRoutes.GET("/:id", r.appHandlers.user.GetUserByID)
		userRoutes.POST("/", r.appHandlers.user.CreateUser)
		userRoutes.POST("/:id/restore", r.appHandlers.user.RestoreUser)
		userRoutes.POST("/:id/status", r.appHandlers.user.UpdateUserStatus)
		userRoutes.DELETE("/:id", r.appHandlers.user.DeleteUser)
	}
}
