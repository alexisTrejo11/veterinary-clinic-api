// Package users defines the User entity and its related functionalities.
package users

import (
	"context"
	"time"

	"clinic-vet-api/internal/core/auth"
	"clinic-vet-api/internal/core/customers"
	"clinic-vet-api/internal/core/employees"
	"clinic-vet-api/internal/shared"
)

// ============================================================================
// Type Definitions
// ============================================================================

// UserID is the unique identifier for a user
type UserID struct{ shared.BaseID }

// NewUserID creates a new UserID
func NewUserID(value uint) UserID {
	return UserID{shared.BaseID{Value: value}}
}

// User represents a user in the system
type User struct {
	shared.Entity[UserID]
	Email          Email
	PhoneNumber    PhoneNumber
	HashedPassword string
	Role           UserRole
	Status         UserStatus
	LastLoginAt    *time.Time
	TwoFactorAuth  auth.TwoFactorAuth
	EmployeeID     *employees.EmployeeID
	CustomerID     *customers.CustomerID
	Profile        Profile
}

// Profile represents user profile information
type Profile struct {
	Name        string
	Gender      shared.PersonGender
	PhotoURL    string
	Bio         string
	DateOfBirth *time.Time
}

// ============================================================================
// User Methods
// ============================================================================

// UpdateEmail updates the user's email address
func (u *User) UpdateEmail(newEmail Email) error {
	if u.Email.String() == newEmail.String() {
		return nil // No change needed
	}
	u.Email = newEmail
	u.IncrementVersion()
	return nil
}

// UpdatePhoneNumber updates the user's phone number
func (u *User) UpdatePhoneNumber(newPhone PhoneNumber) error {
	u.PhoneNumber = newPhone
	u.IncrementVersion()
	return nil
}

// UpdatePassword updates the user's hashed password
func (u *User) UpdatePassword(newPasswordHash string) error {
	u.HashedPassword = newPasswordHash
	u.IncrementVersion()
	return nil
}

// Activate activates the user account
func (u *User) Activate() error {
	ctx := context.Background()
	operation := "ActivateUser"

	if u.Status == UserStatusBanned {
		return CannotActivateBannedUserError(ctx, operation)
	}
	u.Status = UserStatusActive
	u.IncrementVersion()
	return nil
}

// Deactivate deactivates the user account
func (u *User) Deactivate() error {
	u.Status = UserStatusInactive
	u.IncrementVersion()
	return nil
}

// Ban bans the user account
func (u *User) Ban() error {
	u.Status = UserStatusBanned
	u.IncrementVersion()
	return nil
}

// RecordLogin records the user's last login time
func (u *User) RecordLogin() error {
	now := time.Now()
	u.LastLoginAt = &now
	u.IncrementVersion()
	return nil
}

// EnableTwoFactorAuth enables two-factor authentication for the user
func (u *User) EnableTwoFactorAuth(method, secret string, backupCodes []string) error {
	now := time.Now()
	twoFAMethod := auth.TwoFactorMethod(secret)

	twoFA := auth.NewTwoFactorAuth(
		true,
		twoFAMethod,
		secret,
		backupCodes,
		&now,
		&now,
	)
	u.TwoFactorAuth = twoFA
	u.IncrementVersion()
	return nil
}

// DisableTwoFactorAuth disables two-factor authentication for the user
func (u *User) DisableTwoFactorAuth() error {
	u.TwoFactorAuth = auth.NewDisabledTwoFactorAuth()
	u.IncrementVersion()
	return nil
}

// ValidateLogin validates if the user can login
func (u *User) ValidateLogin(is2FA bool) error {
	ctx := context.Background()
	operation := "ValidateLogin"

	if u.Status != UserStatusActive {
		return UserNotActiveError(ctx, operation)
	}

	if u.Status == UserStatusBanned {
		return UserBannedError(ctx, operation)
	}

	if is2FA && !u.TwoFactorAuth.IsEnabled {
		return TwoFactorRequiredError(ctx, operation)
	}

	return nil
}

// AssignEmployee assigns an employee to this user account
func (u *User) AssignEmployee(employeeID employees.EmployeeID) error {
	ctx := context.Background()
	operation := "AssignEmployee"

	if u.EmployeeID != nil {
		return EmployeeAlreadyAssignedError(ctx, operation)
	}

	if u.CustomerID != nil {
		return CannotAssignEmployeeToCustomerError(ctx, operation)
	}

	u.EmployeeID = &employeeID
	u.IncrementVersion()
	return nil
}
