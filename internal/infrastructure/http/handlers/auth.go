package handlers

import (
	"clinic-vet-api/internal/core/auth"
	"clinic-vet-api/internal/infrastructure/http/handlers/dtos"
	"clinic-vet-api/internal/infrastructure/http/handlers/mappers"
	"clinic-vet-api/internal/middleware"
	sharedhttp "clinic-vet-api/internal/shared/http"
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type AuthHandler struct {
	service   auth.AuthService
	validator *validator.Validate
	mapper    *mappers.AuthMapper
}

func NewAuthHandler(service auth.AuthService, validator *validator.Validate, mapper *mappers.AuthMapper) *AuthHandler {
	if mapper == nil {
		mapper = mappers.NewAuthMapper()
	}
	return &AuthHandler{service: service, validator: validator, mapper: mapper}
}

// Register POST /auth/register — body: RegisterRequest
func (ctrl *AuthHandler) Register(c *gin.Context) {
	var req dtos.RegisterRequest
	if err := sharedhttp.ShouldBindAndValidateBody(c, &req, ctrl.validator); err != nil {
		sharedhttp.BadRequest(c, err)
		return
	}

	command, err := ctrl.mapper.RegisterRequestToCommand(req)
	if err != nil {
		sharedhttp.BadRequest(c, err)
		return
	}

	msg, err := ctrl.service.Register(c.Request.Context(), command)
	if err != nil {
		sharedhttp.ApplicationError(c, err)
		return
	}

	sharedhttp.SuccessCreated(c, msg, "Registration successful")
}

// ActivateAccount GET /auth/activate?user_id=&code= or POST body (query preferred for email links)
func (ctrl *AuthHandler) ActivateAccount(c *gin.Context) {
	var req dtos.ActivateAccountRequest
	if err := c.ShouldBind(&req); err != nil {
		sharedhttp.BadRequest(c, err)
		return
	}
	if err := ctrl.validator.Struct(&req); err != nil {
		sharedhttp.BadRequest(c, err)
		return
	}

	command, err := ctrl.mapper.ActivateAccountRequestToCommand(req)
	if err != nil {
		sharedhttp.BadRequest(c, err)
		return
	}

	err = ctrl.service.ActivateAccount(c.Request.Context(), command)
	if err != nil {
		sharedhttp.ApplicationError(c, err)
		return
	}

	sharedhttp.Success(c, nil, "Account activated successfully")
}

// Login POST /auth/login — body: LoginRequest
func (ctrl *AuthHandler) Login(c *gin.Context) {
	var req dtos.LoginRequest
	if err := sharedhttp.ShouldBindAndValidateBody(c, &req, ctrl.validator); err != nil {
		sharedhttp.BadRequest(c, err)
		return
	}

	command := ctrl.mapper.LoginRequestToCommand(req, c)
	session, err := ctrl.service.Login(c.Request.Context(), command)
	if err != nil {
		var requires2FA *auth.RequiresTwoFactorError
		if errors.As(err, &requires2FA) {
			sharedhttp.Success(c, ctrl.mapper.RequiresTwoFactorToResponse(requires2FA.UserID), "Two-factor authentication required")
			return
		}
		sharedhttp.ApplicationError(c, err)
		return
	}

	sharedhttp.Success(c, ctrl.mapper.SessionPayloadToResponse(session), "Login successful")
}

// RefreshToken POST /auth/refresh — body: RefreshTokenRequest
func (ctrl *AuthHandler) RefreshToken(c *gin.Context) {
	var req dtos.RefreshTokenRequest
	if err := sharedhttp.ShouldBindAndValidateBody(c, &req, ctrl.validator); err != nil {
		sharedhttp.BadRequest(c, err)
		return
	}

	command := ctrl.mapper.RefreshTokenRequestToCommand(req)
	session, err := ctrl.service.RefreshToken(c.Request.Context(), command)
	if err != nil {
		sharedhttp.ApplicationError(c, err)
		return
	}

	sharedhttp.Success(c, ctrl.mapper.SessionPayloadToResponse(session), "Token refreshed successfully")
}

// Logout POST /auth/logout — body: LogoutRequest
func (ctrl *AuthHandler) Logout(c *gin.Context) {
	var req dtos.LogoutRequest
	if err := sharedhttp.ShouldBindAndValidateBody(c, &req, ctrl.validator); err != nil {
		sharedhttp.BadRequest(c, err)
		return
	}

	command := ctrl.mapper.LogoutRequestToCommand(req)
	err := ctrl.service.Logout(c.Request.Context(), command)
	if err != nil {
		sharedhttp.ApplicationError(c, err)
		return
	}

	sharedhttp.Success(c, nil, "Logout successful")
}

// LogoutAll POST /auth/logout-all — requires auth; user from context
func (ctrl *AuthHandler) LogoutAll(c *gin.Context) {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		sharedhttp.Unauthorized(c, err)
		return
	}

	err = ctrl.service.LogoutAll(c.Request.Context(), auth.LogoutAllCommand{UserID: userID})
	if err != nil {
		sharedhttp.ApplicationError(c, err)
		return
	}

	sharedhttp.Success(c, nil, "Logged out from all devices")
}

// VerifyTwoFactor POST /auth/2fa/verify — body: VerifyTwoFactorRequest; requires auth, user from context
func (ctrl *AuthHandler) VerifyTwoFactor(c *gin.Context) {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		sharedhttp.Unauthorized(c, err)
		return
	}

	var req dtos.VerifyTwoFactorRequest
	if err := sharedhttp.ShouldBindAndValidateBody(c, &req, ctrl.validator); err != nil {
		sharedhttp.BadRequest(c, err)
		return
	}

	command := ctrl.mapper.VerifyTwoFactorRequestToCommand(req, userID)
	session, err := ctrl.service.VerifyTwoFactor(c.Request.Context(), command)
	if err != nil {
		sharedhttp.ApplicationError(c, err)
		return
	}

	sharedhttp.Success(c, ctrl.mapper.SessionPayloadToResponse(session), "Two factor verified successfully")
}

// EnableTwoFactor POST /auth/2fa/enable — body: EnableTwoFactorRequest; requires auth
func (ctrl *AuthHandler) EnableTwoFactor(c *gin.Context) {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		sharedhttp.Unauthorized(c, err)
		return
	}

	var req dtos.EnableTwoFactorRequest
	if err := sharedhttp.ShouldBindAndValidateBody(c, &req, ctrl.validator); err != nil {
		sharedhttp.BadRequest(c, err)
		return
	}

	command := ctrl.mapper.EnableTwoFactorRequestToCommand(req, userID)
	err = ctrl.service.EnableTwoFactor(c.Request.Context(), command)
	if err != nil {
		sharedhttp.ApplicationError(c, err)
		return
	}

	sharedhttp.Success(c, nil, "Two factor enabled successfully")
}

// DisableTwoFactor POST /auth/2fa/disable — requires auth; user from context only
func (ctrl *AuthHandler) DisableTwoFactor(c *gin.Context) {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		sharedhttp.Unauthorized(c, err)
		return
	}

	err = ctrl.service.DisableTwoFactor(c.Request.Context(), auth.DisableTwoFactorCommand{UserID: userID})
	if err != nil {
		sharedhttp.ApplicationError(c, err)
		return
	}

	sharedhttp.Success(c, nil, "Two factor disabled successfully")
}

// VerifyEmail POST /auth/verify-email — body: VerifyEmailRequest
func (ctrl *AuthHandler) VerifyEmail(c *gin.Context) {
	var req dtos.VerifyEmailRequest
	if err := sharedhttp.ShouldBindAndValidateBody(c, &req, ctrl.validator); err != nil {
		sharedhttp.BadRequest(c, err)
		return
	}

	command := ctrl.mapper.VerifyEmailRequestToCommand(req)
	err := ctrl.service.VerifyEmail(c.Request.Context(), command)
	if err != nil {
		sharedhttp.ApplicationError(c, err)
		return
	}

	sharedhttp.Success(c, nil, "Email verified successfully")
}

// RequestResetPassword POST /auth/request-reset-password — body: RequestResetPasswordRequest
func (ctrl *AuthHandler) RequestResetPassword(c *gin.Context) {
	var req dtos.RequestResetPasswordRequest
	if err := sharedhttp.ShouldBindAndValidateBody(c, &req, ctrl.validator); err != nil {
		sharedhttp.BadRequest(c, err)
		return
	}

	command := ctrl.mapper.RequestResetPasswordRequestToCommand(req)
	err := ctrl.service.RequestResetPassword(c.Request.Context(), command)
	if err != nil {
		sharedhttp.ApplicationError(c, err)
		return
	}

	sharedhttp.Success(c, nil, "If the email exists, a reset link has been sent")
}

// ResetPassword POST /auth/reset-password — body: ResetPasswordRequest
func (ctrl *AuthHandler) ResetPassword(c *gin.Context) {
	var req dtos.ResetPasswordRequest
	if err := sharedhttp.ShouldBindAndValidateBody(c, &req, ctrl.validator); err != nil {
		sharedhttp.BadRequest(c, err)
		return
	}

	command := ctrl.mapper.ResetPasswordRequestToCommand(req)
	err := ctrl.service.ResetPassword(c.Request.Context(), command)
	if err != nil {
		sharedhttp.ApplicationError(c, err)
		return
	}

	sharedhttp.Success(c, nil, "Password reset successfully")
}

// UpdatePassword POST /auth/update-password — body: UpdatePasswordRequest; requires auth
func (ctrl *AuthHandler) UpdatePassword(c *gin.Context) {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		sharedhttp.Unauthorized(c, err)
		return
	}

	var req dtos.UpdatePasswordRequest
	if err := sharedhttp.ShouldBindAndValidateBody(c, &req, ctrl.validator); err != nil {
		sharedhttp.BadRequest(c, err)
		return
	}

	command := ctrl.mapper.UpdatePasswordRequestToCommand(req, userID)
	err = ctrl.service.UpdatePassword(c.Request.Context(), command)
	if err != nil {
		sharedhttp.ApplicationError(c, err)
		return
	}

	sharedhttp.Success(c, nil, "Password updated successfully")
}

// DeleteMyAccount DELETE /auth/me or POST /auth/delete-account — requires auth; user from context
func (ctrl *AuthHandler) DeleteMyAccount(c *gin.Context) {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		sharedhttp.Unauthorized(c, err)
		return
	}

	err = ctrl.service.DeleteAccount(c.Request.Context(), userID)
	if err != nil {
		sharedhttp.ApplicationError(c, err)
		return
	}

	sharedhttp.Success(c, nil, "Account deleted successfully")
}
