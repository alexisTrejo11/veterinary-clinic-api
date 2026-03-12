package users

import (
	"clinic-vet-api/internal/shared"
	"time"
)

type UpdatePasswordCommand struct {
	ID              shared.UserID
	CurrentPassword string
	NewPassword     string
}

type UpdateEmailCommand struct {
	ID    shared.UserID
	Email Email
}

type UpdatePhoneCommand struct {
	ID    shared.UserID
	Phone PhoneNumber
}

type UpdateUserStatusCommand struct {
	ID     shared.UserID
	Status UserStatus
}

type UpdateProfileCommand struct {
	ID          shared.UserID
	Name        string
	Gender      shared.PersonGender
	PhotoURL    string
	Bio         string
	DateOfBirth *time.Time
}

type CreateUserCommand struct {
	Email         Email
	PhoneNumber   *PhoneNumber
	PlainPassword string
	Role          UserRole
	Status        UserStatus
}

type DeleteUserCommand struct {
	ID           shared.UserID
	IsHardDelete bool
}

type RequestResetPasswordCommand struct {
	Email Email
}

type ResetPasswordCommand struct {
	Email       Email
	Token       string
	NewPassword string
}

// ResetPasswordByCodeCommand is used when password is reset via verification code (no current password).
type ResetPasswordByCodeCommand struct {
	ID          shared.UserID
	NewPassword string
}
