package mappers

import (
	"clinic-vet-api/internal/core/users"
	"clinic-vet-api/internal/infrastructure/http/handlers/dtos"
)

type ProfileMapper struct {
}

type UserHandlerMapper struct {
	profile ProfileMapper
}

func (m *ProfileMapper) ToProfileResponse(user users.User) dtos.ProfileResponse {
	return dtos.ProfileResponse{
		ID:            user.ID.String(),
		Email:         user.Email.Value(),
		Phone:         user.PhoneNumber.Value(),
		Role:          user.Role.DisplayName(),
		Status:        user.Status.DisplayName(),
		Name:          user.Profile.Name,
		Gender:        user.Profile.Gender.DisplayName(),
		DateOfBirth:   user.Profile.DateOfBirth,
		ProfilePicUrl: user.Profile.PhotoURL,
		Bio:           user.Profile.Bio,
	}
}

func (m *UserHandlerMapper) ToCreateUserCommand(req dtos.CreateUserRequest) (users.CreateUserCommand, error) {
	var status users.UserStatus
	if req.Status != nil {
		statusObj, err := users.ParseUserStatus(*req.Status)
		if err != nil {
			return users.CreateUserCommand{}, err
		}
		status = statusObj
	}

	role, err := users.ParseUserRole(req.Role)
	if err != nil {
		return users.CreateUserCommand{}, err
	}

	email, err := users.NewEmail(req.Email) // Validate email format early
	if err != nil {
		return users.CreateUserCommand{}, err
	}

	var phoneNumber *users.PhoneNumber
	if req.PhoneNumber != nil {
		phone, err := users.NewPhoneNumber(*req.PhoneNumber) // Validate phone number format early
		if err != nil {
			return users.CreateUserCommand{}, err
		}
		phoneNumber = &phone
	}

	return users.CreateUserCommand{
		Email:         email,
		PhoneNumber:   phoneNumber,
		PlainPassword: req.Password,
		Role:          role,
		Status:        status,
	}, nil

}

func (m *UserHandlerMapper) ToResponse(user users.User) dtos.UserResponse {
	return dtos.UserResponse{
		ID:            user.ID.Value,
		Email:         user.Email.Value(),
		PhoneNumber:   user.PhoneNumber.Value(),
		Role:          user.Role.DisplayName(),
		Status:        user.Status.DisplayName(),
		EmployeeID:    &user.EmployeeID.Value,
		CustomerID:    &user.CustomerID.Value,
		OAuthProvider: user.OAuthProvider,
		EmailVerified: user.EmailVerified,
		CreatedAt:     user.CreatedAt,
		UpdatedAt:     user.UpdatedAt,
		Profile:       m.profile.ToProfileResponse(user),
	}
}
