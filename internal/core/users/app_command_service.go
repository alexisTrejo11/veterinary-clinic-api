package users

import (
	"clinic-vet-api/internal/core/auth"
	"clinic-vet-api/internal/shared"
	"clinic-vet-api/internal/shared/password"
	"context"
)

type CommandService interface {
	CreateUser(ctx context.Context, cmd CreateUserCommand) (User, error)
	UpdatePassword(ctx context.Context, cmd UpdatePasswordCommand) error
	UpdateUserStatus(ctx context.Context, cmd UpdateUserStatusCommand) error
	UpdateProfileData(ctx context.Context, cmd UpdateProfileCommand) error
	RestoreUser(ctx context.Context, id shared.UserID) error
	DeleteUser(ctx context.Context, cmd DeleteUserCommand) error
}

type commandService struct {
	repository      UserRepository
	passwordEncoder password.PasswordEncoder
}

func NewCommandService(
	repository UserRepository, passwordEncoder password.PasswordEncoder,
) CommandService {
	return &commandService{
		repository:      repository,
		passwordEncoder: passwordEncoder,
	}
}

func (s *commandService) CreateUser(ctx context.Context, cmd CreateUserCommand) (User, error) {
	if err := s.validateUniqueEmail(ctx, cmd.Email, "CreateUser"); err != nil {
		return User{}, err
	}

	if cmd.PhoneNumber != nil {
		if err := s.validateUniquePhone(ctx, *cmd.PhoneNumber, "CreateUser"); err != nil {
			return User{}, err
		}
	}

	hashedPassword, err := s.hashPasswordWithValidation(ctx, cmd.PlainPassword)
	if err != nil {
		return User{}, err
	}

	user := s.createUserEntity(cmd, hashedPassword)

	err = s.repository.Save(ctx, &user)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (s *commandService) UpdatePassword(ctx context.Context, cmd UpdatePasswordCommand) error {
	user, err := s.repository.FindByID(ctx, cmd.ID)
	if err != nil {
		return err
	}

	valid := s.passwordEncoder.CheckPassword(user.HashedPassword, cmd.CurrentPassword)
	if !valid {
		return PasswordMismatchError(ctx, "UpdatePassword")
	}

	if err := auth.ValidatePasswordStrength(cmd.NewPassword); err != nil {
		return err
	}

	hashedPassword, err := s.passwordEncoder.HashPassword(cmd.NewPassword)
	if err != nil {
		return err

	}

	user.UpdatePassword(hashedPassword)
	return s.repository.Save(ctx, &user)
}

func (s *commandService) UpdateUserStatus(ctx context.Context, cmd UpdateUserStatusCommand) error {
	user, err := s.repository.FindByID(ctx, cmd.ID)
	if err != nil {
		return err
	}

	switch cmd.Status {
	case UserStatusActive:
		if err := user.Activate(); err != nil {
			return err
		}
	case UserStatusInactive:
		if err := user.Deactivate(); err != nil {
			return err
		}
	case UserStatusBanned:
		if err := user.Ban(); err != nil {
			return err
		}
	default:
		return InvalidUserStatusError(ctx, cmd.Status.DisplayName(), "UpdateUserStatus")
	}

	return s.repository.Save(ctx, &user)
}

func (s *commandService) UpdatePhoneNumber(ctx context.Context, cmd UpdatePhoneCommand) error {
	user, err := s.repository.FindByID(ctx, cmd.ID)
	if err != nil {
		return err
	}

	if err := s.validateUniquePhone(ctx, cmd.Phone, "UpdatePhoneNumber"); err != nil {
		return err
	}

	user.UpdatePhoneNumber(cmd.Phone)
	return s.repository.Save(ctx, &user)
}

func (h *commandService) UpdateEmail(ctx context.Context, cmd UpdateEmailCommand) error {
	user, err := h.repository.FindByID(ctx, cmd.ID)
	if err != nil {
		return err
	}

	if err := h.validateUniqueEmail(ctx, cmd.Email, "UpdateEmail"); err != nil {
		return err
	}

	user.UpdateEmail(cmd.Email)
	return h.repository.Save(ctx, &user)
}

func (s *commandService) RestoreUser(ctx context.Context, id shared.UserID) error {
	if deleted, err := s.repository.IsDeletedByID(ctx, id); err != nil {
		return err
	} else if !deleted {
		return UserNotDeletedError(ctx, id, "RestoreUser")
	}

	return s.repository.RestoreByID(ctx, id)
}

func (s *commandService) DeleteUser(ctx context.Context, cmd DeleteUserCommand) error {
	if _, err := s.repository.FindByID(ctx, cmd.ID); err != nil {
		return err
	}

	if cmd.IsHardDelete {
		return s.repository.HardDelete(ctx, cmd.ID)
	}

	return s.repository.SoftDelete(ctx, cmd.ID)

}

// validateUniqueEmail checks if email is already in use
func (s *commandService) validateUniqueEmail(ctx context.Context, email Email, operation string) error {
	exists, err := s.repository.ExistsByEmail(ctx, email)
	if err != nil {
		return err
	}
	if exists {
		return EmailAlreadyExistsError(ctx, email, operation)
	}
	return nil
}

// validateUniquePhone checks if phone number is already in use
func (s *commandService) validateUniquePhone(ctx context.Context, phone PhoneNumber, operation string) error {
	exists, err := s.repository.ExistsByPhone(ctx, phone)
	if err != nil {
		return err
	}
	if exists {
		return PhoneAlreadyExistsError(ctx, phone, operation)
	}
	return nil
}

// hashPasswordWithValidation validates password strength and returns hashed password
func (s *commandService) hashPasswordWithValidation(ctx context.Context, plainPassword string) (string, error) {
	// Validate password strength
	if err := auth.ValidatePasswordStrength(plainPassword); err != nil {
		return "", err
	}

	// Hash password
	hashedPassword, err := s.passwordEncoder.HashPassword(plainPassword)
	if err != nil {
		return "", err
	}

	return hashedPassword, nil
}

func (s *commandService) UpdateProfileData(ctx context.Context, cmd UpdateProfileCommand) error {
	user, err := s.repository.FindByID(ctx, cmd.ID)
	if err != nil {
		return err
	}

	Profile := Profile{
		Name:        cmd.Name,
		Gender:      cmd.Gender,
		PhotoURL:    cmd.PhotoURL,
		Bio:         cmd.Bio,
		DateOfBirth: cmd.DateOfBirth,
	}

	if err := user.UpdateProfile(Profile); err != nil {
		return err
	}

	return s.repository.Save(ctx, &user)
}

// createUserEntity creates a new User entity from command
func (s *commandService) createUserEntity(cmd CreateUserCommand, hashedPassword string) User {
	phoneNumber := PhoneNumber{}
	if cmd.PhoneNumber != nil {
		phoneNumber = *cmd.PhoneNumber
	}

	return User{
		Entity:         shared.CreateEntity(shared.UserID{}),
		Email:          cmd.Email,
		PhoneNumber:    phoneNumber,
		HashedPassword: hashedPassword,
		Role:           cmd.Role,
		Status:         cmd.Status,
		TwoFactorAuth:  auth.NewDisabledTwoFactorAuth(),
		Profile:        Profile{},
		EmailVerified:  false,
		LoginAttempts:  0,
	}
}
