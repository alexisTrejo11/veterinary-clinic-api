package users

import (
	"context"
	"strings"
)

// ============================================================================
// UserRole Enum
// ============================================================================

// UserRole represents the role of a user in the system
type UserRole string

const (
	UserRoleAdmin        UserRole = "admin"
	UserRoleVeterinarian UserRole = "veterinarian"
	UserRoleCustomer     UserRole = "customer"
	UserRoleReceptionist UserRole = "receptionist"
)

var (
	ValidUserRoles = []UserRole{
		UserRoleAdmin,
		UserRoleVeterinarian,
		UserRoleCustomer,
		UserRoleReceptionist,
	}

	userRoleMap = map[string]UserRole{
		"admin":        UserRoleAdmin,
		"veterinarian": UserRoleVeterinarian,
		"customer":     UserRoleCustomer,
		"receptionist": UserRoleReceptionist,
	}

	userRoleDisplayNames = map[UserRole]string{
		UserRoleAdmin:        "Administrator",
		UserRoleVeterinarian: "Veterinarian",
		UserRoleCustomer:     "Pet customer",
		UserRoleReceptionist: "Receptionist",
	}
)

// IsValid checks if the user role is valid
func (r UserRole) IsValid() bool {
	_, exists := userRoleMap[string(r)]
	return exists
}

// ParseUserRole parses a string into a UserRole
func ParseUserRole(role string) (UserRole, error) {
	ctx := context.Background()
	operation := "ParseUserRole"

	normalized := strings.TrimSpace(strings.ToLower(role))
	if val, exists := userRoleMap[normalized]; exists {
		return val, nil
	}
	return "", InvalidUserRoleError(ctx, role, operation)
}

// MustParseUserRole parses a string into a UserRole, panics on error
func MustParseUserRole(role string) UserRole {
	parsed, err := ParseUserRole(role)
	if err != nil {
		panic(err)
	}
	return parsed
}

// String returns the string representation
func (r UserRole) String() string {
	return string(r)
}

// DisplayName returns the human-readable name
func (r UserRole) DisplayName() string {
	if displayName, exists := userRoleDisplayNames[r]; exists {
		return displayName
	}
	return "Unknown Role"
}

// Values returns all valid user roles
func (r UserRole) Values() []UserRole {
	return ValidUserRoles
}

// IsAdministrative checks if the role is administrative
func (r UserRole) IsAdministrative() bool {
	return r == UserRoleAdmin
}

// IsStaff checks if the role is staff (admin, vet, or receptionist)
func (r UserRole) IsStaff() bool {
	return r == UserRoleAdmin || r == UserRoleVeterinarian || r == UserRoleReceptionist
}

func (r UserRole) IsEmployee() bool {
	return r == UserRoleVeterinarian || r == UserRoleReceptionist
}

func (r UserRole) IsCustomer() bool {
	return r == UserRoleCustomer
}

// ============================================================================
// UserStatus Enum
// ============================================================================

// UserStatus represents the status of a user account
type UserStatus string

const (
	UserStatusActive   UserStatus = "active"
	UserStatusInactive UserStatus = "inactive"
	UserStatusBanned   UserStatus = "banned"
	UserStatusPending  UserStatus = "pending"
	UserStatusDeleted  UserStatus = "deleted"
)

var (
	ValidUserStatuses = []UserStatus{
		UserStatusActive,
		UserStatusInactive,
		UserStatusBanned,
		UserStatusPending,
		UserStatusDeleted,
	}

	userStatusMap = map[string]UserStatus{
		"active":   UserStatusActive,
		"inactive": UserStatusInactive,
		"banned":   UserStatusBanned,
		"pending":  UserStatusPending,
		"deleted":  UserStatusDeleted,
	}

	userStatusDisplayNames = map[UserStatus]string{
		UserStatusActive:   "Active",
		UserStatusInactive: "Inactive",
		UserStatusBanned:   "Banned",
		UserStatusPending:  "Pending",
		UserStatusDeleted:  "Deleted",
	}
)

// IsValid checks if the user status is valid
func (s UserStatus) IsValid() bool {
	_, exists := userStatusMap[string(s)]
	return exists
}

// ParseUserStatus parses a string into a UserStatus
func ParseUserStatus(status string) (UserStatus, error) {
	ctx := context.Background()
	operation := "ParseUserStatus"

	normalized := strings.TrimSpace(strings.ToLower(status))
	if val, exists := userStatusMap[normalized]; exists {
		return val, nil
	}
	return "", InvalidUserStatusError(ctx, status, operation)
}

// MustParseUserStatus parses a string into a UserStatus, panics on error
func MustParseUserStatus(status string) UserStatus {
	parsed, err := ParseUserStatus(status)
	if err != nil {
		panic(err)
	}
	return parsed
}

// String returns the string representation
func (s UserStatus) String() string {
	return string(s)
}

// DisplayName returns the human-readable name
func (s UserStatus) DisplayName() string {
	if displayName, exists := userStatusDisplayNames[s]; exists {
		return displayName
	}
	return "Unknown Status"
}

// Values returns all valid user statuses
func (s UserStatus) Values() []UserStatus {
	return ValidUserStatuses
}

// IsActive checks if the status is active
func (s UserStatus) IsActive() bool {
	return s == UserStatusActive
}

// CanLogin checks if the user can login
func (s UserStatus) CanLogin() bool {
	return s == UserStatusActive
}

// ============================================================================
// Helper Functions
// ============================================================================

// GetAllUserRoles returns all valid user roles
func GetAllUserRoles() []UserRole {
	return ValidUserRoles
}

// GetAllUserStatuses returns all valid user statuses
func GetAllUserStatuses() []UserStatus {
	return ValidUserStatuses
}
