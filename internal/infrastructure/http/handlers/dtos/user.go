package dtos

import (
	"time"

	"clinic-vet-api/internal/shared/page"
)

// =======================================================================
// Request DTOs
// =======================================================================

// CreateUserRequest represents the request to create a new user
// @Description Request body for creating a new user account
type CreateUserRequest struct {
	// Email address of the user
	// Required: true
	// Example: user@example.com``
	Email string `json:"email" binding:"required,email" example:"user@example.com"`

	// Phone number of the user
	// Required: false
	// Example: +1234567890
	PhoneNumber *string `json:"phone_number,omitempty" binding:"omitempty,min=10" example:"+1234567890"`

	// Plain text password (will be hashed)
	// Required: true
	// Minimum length: 8 characters
	Password string `json:"password" binding:"required,min=8" example:"SecurePass123!"`

	// Role of the user (admin, veterinarian, customer, receptionist)
	// Required: true
	// Example: customer
	Role string `json:"role" binding:"required,oneof=admin veterinarian customer receptionist" example:"customer"`

	// Initial status of the user account
	// Required: false (defaults to "pending")
	// Example: active
	Status *string `json:"status,omitempty" binding:"omitempty,oneof=active inactive pending" example:"active"`

	// User profile information
	Profile *ProfileInfo `json:"profile,omitempty"`
}

// UpdateUserPasswordRequest represents the request to update a user's password
// @Description Request body for updating user password
type UpdateUserPasswordRequest struct {
	// Current password
	// Required: true
	CurrentPassword string `json:"current_password" binding:"required" example:"OldPass123!"`

	// New password
	// Required: true
	// Minimum length: 8 characters
	NewPassword string `json:"new_password" binding:"required,min=8" example:"NewPass123!"`
}

// UpdateUserStatusRequest represents the request to update a user's status
// @Description Request body for updating user account status
type UpdateUserStatusRequest struct {
	// New status (active, inactive, banned)
	// Required: true
	Status string `json:"status" binding:"required,oneof=active inactive banned" example:"active"`
}

// UpdateUserEmailRequest represents the request to update a user's email
// @Description Request body for updating user email address
type UpdateUserEmailRequest struct {
	// New email address
	// Required: true
	Email string `json:"email" binding:"required,email" example:"newemail@example.com"`
}

// UpdateUserPhoneRequest represents the request to update a user's phone number
// @Description Request body for updating user phone number
type UpdateUserPhoneRequest struct {
	// New phone number
	// Required: true
	PhoneNumber string `json:"phone_number" binding:"required,min=10" example:"+1234567890"`
}

// ProfileInfo represents user profile information
// @Description User profile details
type ProfileInfo struct {
	// Full name of the user
	Name string `json:"name,omitempty" example:"John Doe"`

	// Gender of the user
	Gender string `json:"gender,omitempty" example:"male"`

	// URL to user's profile photo
	PhotoURL string `json:"photo_url,omitempty" example:"https://example.com/photo.jpg"`

	// User biography (max 500 characters)
	Bio string `json:"bio,omitempty" binding:"omitempty,max=500" example:"Animal lover and pet owner"`

	// Date of birth
	DateOfBirth *time.Time `json:"date_of_birth,omitempty" example:"1990-01-15T00:00:00Z"`
}

// =======================================================================
// Response DTOs
// =======================================================================

// UserResponse represents the user data returned in HTTP responses
// @Description User account information (excluding sensitive data)
type UserResponse struct {
	// Unique identifier for the user
	// Example: 123
	ID uint `json:"id" example:"123"`

	// Email address of the user
	// Example: user@example.com
	Email string `json:"email" example:"user@example.com"`

	// Phone number of the user
	// Example: +1234567890
	PhoneNumber string `json:"phone_number,omitempty" example:"+1234567890"`

	// Role of the user (admin, veterinarian, customer, receptionist)
	// Example: customer
	Role string `json:"role" example:"customer"`

	// Status of the user account (active, inactive, banned, pending)
	// Example: active
	Status string `json:"status" example:"active"`

	// ID of the associated employee (for staff users)
	// Example: 456
	EmployeeID *uint `json:"employee_id,omitempty" example:"456"`

	// ID of the associated customer (for customer users)
	// Example: 789
	CustomerID *uint `json:"customer_id,omitempty" example:"789"`

	// User profile information
	Profile ProfileResponse `json:"profile"`

	// OAuth provider name (e.g., "google", "facebook")
	// Example: google
	OAuthProvider string `json:"oauth_provider,omitempty" example:"google"`

	// Whether the email has been verified
	// Example: true
	EmailVerified bool `json:"email_verified" example:"true"`

	// Whether two-factor authentication is enabled
	// Example: false
	TwoFactorEnabled bool `json:"two_factor_enabled" example:"false"`

	// Last login timestamp
	// Example: 2024-03-10T15:30:00Z
	LastLoginAt *time.Time `json:"last_login_at,omitempty" example:"2024-03-10T15:30:00Z"`

	// Account creation timestamp
	// Example: 2024-01-01T00:00:00Z
	CreatedAt time.Time `json:"created_at" example:"2024-01-01T00:00:00Z"`

	// Last update timestamp
	// Example: 2024-03-10T15:30:00Z
	UpdatedAt time.Time `json:"updated_at" example:"2024-03-10T15:30:00Z"`
}

// UserSearchRequest represents query parameters for searching and filtering users.
// It is bound from the request query string using Gin's ShouldBindQuery (form binding).
//
// @Description Query parameters for searching and filtering users (used in GET /users)
type UserSearchRequest struct {
	page.PaginationRequest

	// IDs filters by user IDs (comma-separated, e.g. "1,2,3")
	// @Param ids query string false "Comma-separated user IDs"
	IDs string `form:"ids" example:"1,2,3"`

	// Emails filters by email addresses (comma-separated)
	// @Param emails query string false "Comma-separated emails"
	Emails string `form:"emails" example:"user@example.com"`

	// Roles filters by roles: admin, veterinarian, customer, receptionist (comma-separated)
	// @Param roles query string false "Comma-separated roles (admin,veterinarian,customer,receptionist)"
	Roles string `form:"roles" example:"admin,customer"`

	// Statuses filters by statuses: active, inactive, banned, pending (comma-separated)
	// @Param statuses query string false "Comma-separated statuses (active,inactive,banned,pending)"
	Statuses string `form:"statuses" example:"active,pending"`

	// IsActive filters by active flag: "true" or "false"
	// @Param is_active query string false "Filter by active (true/false)"
	IsActive string `form:"is_active" example:"true"`

	// CreatedAfter filters users created on or after this date (RFC3339 or 2006-01-02)
	// @Param created_after query string false "Created after date (2006-01-02)"
	CreatedAfter string `form:"created_after" example:"2024-01-01"`

	// CreatedBefore filters users created on or before this date
	// @Param created_before query string false "Created before date (2006-01-02)"
	CreatedBefore string `form:"created_before" example:"2024-12-31"`

	// LastLoginAfter filters by last login on or after this date
	// @Param last_login_after query string false "Last login after date (2006-01-02)"
	LastLoginAfter string `form:"last_login_after" example:"2024-01-01"`

	// LastLoginBefore filters by last login on or before this date
	// @Param last_login_before query string false "Last login before date (2006-01-02)"
	LastLoginBefore string `form:"last_login_before" example:"2024-12-31"`

	// SearchTerm searches in email, first name and last name (partial, case-insensitive)
	// @Param search query string false "Search term for email/name"
	SearchTerm string `form:"search" example:"john"`

	// HasTwoFactor filters by two-factor auth enabled: "true" or "false"
	// @Param has_two_factor query string false "Filter by 2FA enabled (true/false)"
	HasTwoFactor string `form:"has_two_factor" example:"false"`
}
