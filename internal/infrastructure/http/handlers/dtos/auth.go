package dtos

// Auth request/response DTOs with validation. Use ShouldBindBodyWithJSON for body, ShouldBindQuery for query.

const (
	maxLenEmail      = 255
	maxLenPassword   = 128
	maxLenCode       = 32
	maxLenRole       = 50
	maxLenName       = 200
	maxLenPhone      = 30
	maxLenBio        = 500
	maxLenDevice     = 200
	maxLenUserAgent  = 512
	maxLenMethod     = 20
)

// RegisterRequest is the body for POST /auth/register
type RegisterRequest struct {
	Role          string  `json:"role" binding:"required,oneof=admin veterinarian customer receptionist"`
	Email         string  `json:"email" binding:"required,email,max=255"`
	Password      string  `json:"password" binding:"required,min=8,max=128"`
	Name          string  `json:"name" binding:"omitempty,max=200"`
	PhoneNumber   string  `json:"phone_number" binding:"omitempty,max=30"`
	DateOfBirth   string  `json:"date_of_birth" binding:"omitempty,max=20"`
	Gender        string  `json:"gender" binding:"omitempty,max=20"`
	PhotoURL      string  `json:"photo_url" binding:"omitempty,max=500"`
	Bio           string  `json:"bio" binding:"omitempty,max=500"`
	EmployeeID    *uint   `json:"employee_id,omitempty"`
	LicenseNumber *string `json:"license_number,omitempty" binding:"omitempty,max=100"`
}

// ActivateAccountRequest for GET /auth/activate (query) or POST body
type ActivateAccountRequest struct {
	UserID string `form:"user_id" json:"user_id" binding:"required,max=20"`
	Code   string `form:"code" json:"code" binding:"required,max=32"`
}

// LoginRequest is the body for POST /auth/login
type LoginRequest struct {
	Email          string  `json:"email" binding:"required,email,max=255"`
	Password       string  `json:"password" binding:"required,max=128"`
	TwoFactorCode  *string `json:"two_factor_code,omitempty" binding:"omitempty,max=16"`
}

// RefreshTokenRequest is the body for POST /auth/refresh
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required,max=2000"`
}

// LogoutRequest is the body for POST /auth/logout
type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required,max=2000"`
}

// VerifyTwoFactorRequest is the body for POST /auth/2fa/verify
type VerifyTwoFactorRequest struct {
	TwoFactorCode string `json:"two_factor_code" binding:"required,max=16"`
}

// EnableTwoFactorRequest is the body for POST /auth/2fa/enable
type EnableTwoFactorRequest struct {
	Method string `json:"method" binding:"required,oneof=totp sms email"`
}

// VerifyEmailRequest is the body for POST /auth/verify-email
type VerifyEmailRequest struct {
	Email string `json:"email" binding:"required,email,max=255"`
	Code  string `json:"code" binding:"required,max=32"`
}

// RequestResetPasswordRequest is the body for POST /auth/request-reset-password
type RequestResetPasswordRequest struct {
	Email string `json:"email" binding:"required,email,max=255"`
}

// ResetPasswordRequest is the body for POST /auth/reset-password
type ResetPasswordRequest struct {
	Email    string `json:"email" binding:"required,email,max=255"`
	Code     string `json:"code" binding:"required,max=32"`
	Password string `json:"password" binding:"required,min=8,max=128"`
}

// UpdatePasswordRequest is the body for POST /auth/update-password
type UpdatePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required,max=128"`
	NewPassword     string `json:"new_password" binding:"required,min=8,max=128"`
}

// SessionResponse returned on login, refresh, verify-2fa
type SessionResponse struct {
	AccessToken  TokenResponse   `json:"access_token"`
	RefreshToken TokenResponse   `json:"refresh_token"`
	User         UserAuthResponse `json:"user"`
}

// TokenResponse for access/refresh in responses
type TokenResponse struct {
	Token     string `json:"token"`
	ExpiresAt string `json:"expires_at"`
	Type      string `json:"type"`
}

// UserAuthResponse minimal user info for auth responses
type UserAuthResponse struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

// RequiresTwoFactorResponse when login needs 2FA step
type RequiresTwoFactorResponse struct {
	RequiresTwoFactor bool   `json:"requires_two_factor"`
	UserID            uint   `json:"user_id"`
	Message           string `json:"message"`
}
