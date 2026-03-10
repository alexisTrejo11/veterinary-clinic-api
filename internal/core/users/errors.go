package users

import (
	"context"
	"fmt"

	"clinic-vet-api/internal/shared"
	"clinic-vet-api/internal/shared/errors"
)

// Error codes for users domain
const (
	EntityName = "user"
)

// Field names
const (
	FieldEmail     = "email"
	FieldPhone     = "phone_number"
	FieldFirstName = "first_name"
	FieldLastName  = "last_name"
	FieldAge       = "age"
	FieldPassword  = "password"
	FieldRole      = "role"
	FieldStatus    = "status"
)

// =========================================================================
// User business rule errors
// =========================================================================

func InvalidProfileError(ctx context.Context, message string) error {
	return errors.BusinessRuleError(ctx, "invalid_profile", EntityName, message, "Profile validation")
}

// CannotActivateBannedUserError creates an error when trying to activate a banned user
func CannotActivateBannedUserError(ctx context.Context, operation string) error {
	return errors.BusinessRuleError(ctx, "activation_banned",
		EntityName,
		"cannot activate a banned user account",
		operation)
}

// UserNotActiveError creates an error when user is not in active status
func UserNotActiveError(ctx context.Context, operation string) error {
	return errors.BusinessRuleError(ctx, "user_not_active",
		EntityName,
		"user account is not active",
		operation)
}

// UserBannedError creates an error when user is banned
func UserBannedError(ctx context.Context, operation string) error {
	return errors.BusinessRuleError(ctx, "user_banned",
		EntityName,
		"user account is banned",
		operation)
}

// TwoFactorRequiredError creates an error when 2FA is required but not enabled
func TwoFactorRequiredError(ctx context.Context, operation string) error {
	return errors.BusinessRuleError(ctx, "2fa_required",
		EntityName,
		"two-factor authentication is required to login",
		operation)
}

// EmployeeAlreadyAssignedError creates an error when employee is already assigned
func EmployeeAlreadyAssignedError(ctx context.Context, operation string) error {
	return errors.BusinessRuleError(ctx, "employee_already_assigned",
		EntityName,
		"employee is already assigned to this user",
		operation)
}

// CannotAssignEmployeeToCustomerError creates an error when trying to assign employee to a customer
func CannotAssignEmployeeToCustomerError(ctx context.Context, operation string) error {
	return errors.BusinessRuleError(ctx, "customer_employee_conflict",
		EntityName,
		"cannot assign employee to a customer user",
		operation)
}

// =========================================================================
// Value object validation errors
// =========================================================================

// EmptyEmailError creates an error when email is empty
func EmptyEmailError(ctx context.Context, operation string) error {
	return errors.ValidationError(ctx, EntityName, FieldEmail,
		"email cannot be empty",
		operation)
}

// InvalidEmailFormatError creates an error when email format is invalid
func InvalidEmailFormatError(ctx context.Context, email, operation string) error {
	return errors.ValidationError(ctx, EntityName, FieldEmail,
		fmt.Sprintf("invalid email format: %s", email),
		operation)
}

// EmptyPhoneNumberError creates an error when phone number is empty
func EmptyPhoneNumberError(ctx context.Context, operation string) error {
	return errors.ValidationError(ctx, EntityName, FieldPhone,
		"phone number cannot be empty",
		operation)
}

// PhoneNumberTooShortError creates an error when phone number is too short
func PhoneNumberTooShortError(ctx context.Context, operation string) error {
	return errors.ValidationError(ctx, EntityName, FieldPhone,
		"phone number must have at least 10 digits",
		operation)
}

// EmptyFirstNameError creates an error when first name is empty
func EmptyFirstNameError(ctx context.Context, operation string) error {
	return errors.ValidationError(ctx, EntityName, FieldFirstName,
		"first name cannot be empty",
		operation)
}

// EmptyLastNameError creates an error when last name is empty
func EmptyLastNameError(ctx context.Context, operation string) error {
	return errors.ValidationError(ctx, EntityName, FieldLastName,
		"last name cannot be empty",
		operation)
}

// NegativeAgeError creates an error when age values are negative
func NegativeAgeError(ctx context.Context, operation string) error {
	return errors.ValidationError(ctx, EntityName, FieldAge,
		"age years and months must be non-negative",
		operation)
}

// UnrealisticAgeError creates an error when age is too high
func UnrealisticAgeError(ctx context.Context, maxYears int, operation string) error {
	return errors.ValidationError(ctx, EntityName, FieldAge,
		fmt.Sprintf("age must be less than or equal to %d years", maxYears),
		operation)
}

// =========================================================================
// Enum errors
// =========================================================================

// InvalidUserRoleError creates an error for invalid user role
func InvalidUserRoleError(ctx context.Context, role, operation string) error {
	return errors.ValidationError(ctx, EntityName, FieldRole,
		fmt.Sprintf("invalid user role: %s", role),
		operation)
}

// InvalidUserStatusError creates an error for invalid user status
func InvalidUserStatusError(ctx context.Context, status string, operation string) error {
	return errors.ValidationError(ctx, EntityName, FieldStatus,
		fmt.Sprintf("invalid user status: %s", status),
		operation)
}

// =========================================================================
// Service errors
// =========================================================================

func PasswordMismatchError(ctx context.Context, operation string) error {
	return errors.ValidationError(ctx, EntityName, FieldPassword, "current password is incorrect", operation)
}

func UserNotDeletedError(ctx context.Context, id shared.UserID, operation string) error {
	return errors.BusinessRuleError(ctx, "user_not_deleted", EntityName, fmt.Sprintf("user with ID %s is not deleted", id), operation)
}

func EmailAlreadyExistsError(ctx context.Context, email Email, operation string) error {
	return errors.BusinessRuleError(ctx, "email_already_exists", EntityName, fmt.Sprintf("email %s is already in use", email), operation)
}

func PhoneAlreadyExistsError(ctx context.Context, phone PhoneNumber, operation string) error {
	return errors.BusinessRuleError(ctx, "phone_already_exists", EntityName, fmt.Sprintf("phone number %s is already in use", phone), operation)
}
