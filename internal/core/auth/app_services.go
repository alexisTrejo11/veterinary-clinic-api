package auth

import (
	"clinic-vet-api/internal/core/users"
	"clinic-vet-api/internal/shared"
	"context"
	"errors"
	"fmt"
	"strconv"
)

type Message string

const (
	AdminMessageSuccess    Message = "success registration log in available"
	EmployeeMessageSuccess Message = "success registration. An Email will be sent to the user to activate their account."
	CustomerMessageSuccess Message = "success registration. An Email will be sent to the user to activate their account."
	MessageError           Message = "error registration log in not available"
)

type TokenService interface {
	CreateToken(ctx context.Context, userClaims map[string]any, tokenType TokenType) (Token, error)
	GetTokenClaims(ctx context.Context, token string, tokenType TokenType) (map[string]any, error)
	VerifyToken(ctx context.Context, token string, tokenType TokenType) (bool, error)
	RefreshAccessToken(ctx context.Context, token string) (Token, error)
	RevokeToken(ctx context.Context, token string, tokenType TokenType) error
	SaveToken(ctx context.Context, token Token) error
}

type SessionPayload struct {
	AccessToken  Token
	RefreshToken Token
	User         users.User
}

type AuthService interface {
	Register(ctx context.Context, cmd RegisterCommand) (Message, error)
	ActivateAccount(ctx context.Context, cmd ActivateAccountCommand) error

	Login(ctx context.Context, cmd LoginCommand) (SessionPayload, error)
	RefreshToken(ctx context.Context, cmd RefreshTokenCommand) (SessionPayload, error)

	Logout(ctx context.Context, cmd LogoutCommand) error
	LogoutAll(ctx context.Context, cmd LogoutAllCommand) error

	VerifyTwoFactor(ctx context.Context, cmd VerifyTwoFactorCommand) (SessionPayload, error)

	EnableTwoFactor(ctx context.Context, cmd EnableTwoFactorCommand) error
	DisableTwoFactor(ctx context.Context, cmd DisableTwoFactorCommand) error

	VerifyEmail(ctx context.Context, cmd VerifyEmailCommand) error
	RequestResetPassword(ctx context.Context, cmd RequestResetPasswordCommand) error
	ResetPassword(ctx context.Context, cmd ResetPasswordCommand) error
	UpdatePassword(ctx context.Context, cmd UpdatePasswordCommand) error

	DeleteAccount(ctx context.Context, userID shared.UserID) error
}

// NewAuthService builds the auth service with all dependencies. All external concerns (tokens, sessions, users, 2FA) are behind interfaces.
func NewAuthService(
	tokenService TokenService,
	userCommand users.CommandService,
	userQuery users.QueryService,
	sessionRepository SessionRepository,
	passwordChecker PasswordChecker,
	twoFactorProvider TwoFactorProvider,
) AuthService {
	return &authService{
		tokenService:      tokenService,
		userCommand:       userCommand,
		userQuery:         userQuery,
		sessionRepository: sessionRepository,
		passwordChecker:   passwordChecker,
		twoFactorProvider: twoFactorProvider,
	}
}

type authService struct {
	tokenService      TokenService
	userCommand       users.CommandService
	userQuery         users.QueryService
	sessionRepository SessionRepository
	passwordChecker   PasswordChecker
	twoFactorProvider TwoFactorProvider
}

func (s *authService) Register(ctx context.Context, cmd RegisterCommand) (Message, error) {
	const op = "Register"

	if !cmd.Role.IsValid() {
		return MessageError, InvalidRoleError(ctx, string(cmd.Role), op)
	}

	if err := shared.ValidatePasswordStrength(cmd.Password); err != nil {
		return MessageError, err
	}

	email, err := users.NewEmail(cmd.Email)
	if err != nil {
		return MessageError, err
	}

	var phone *users.PhoneNumber
	if cmd.PhoneNumber != "" {
		p, err := users.NewPhoneNumber(cmd.PhoneNumber)
		if err != nil {
			return MessageError, err
		}
		phone = &p
	}

	status := users.UserStatusPending
	if cmd.Role.IsAdministrative() {
		status = users.UserStatusActive
	} else if cmd.Role.IsEmployee() {
		if !cmd.HasEmployeeData() {
			return MessageError, EmployeeDataRequiredError(ctx, op)
		}
		status = users.UserStatusActive
	}

	createCmd := users.CreateUserCommand{
		Email:         email,
		PhoneNumber:   phone,
		PlainPassword: cmd.Password,
		Role:          cmd.Role,
		Status:        status,
	}

	_, err = s.userCommand.CreateUser(ctx, createCmd)
	if err != nil {
		return MessageError, err
	}

	if cmd.Role.IsEmployee() {
		return EmployeeMessageSuccess, nil
	}
	if cmd.Role.IsCustomer() {
		return CustomerMessageSuccess, nil
	}
	return AdminMessageSuccess, nil
}

func (s *authService) ActivateAccount(ctx context.Context, cmd ActivateAccountCommand) error {
	const op = "ActivateAccount"

	ok, err := s.tokenService.VerifyToken(ctx, cmd.Code, TokenTypeVerification)
	if err != nil || !ok {
		return InvalidVerificationCodeError(ctx, op)
	}

	claims, err := s.tokenService.GetTokenClaims(ctx, cmd.Code, TokenTypeVerification)
	if err != nil {
		return InvalidVerificationCodeError(ctx, op)
	}

	uidStr, _ := claims[ClaimUserID].(string)
	if uidStr == "" {
		return InvalidVerificationCodeError(ctx, op)
	}

	uid, err := parseUserID(uidStr)
	if err != nil || uid.Value() != cmd.UserID.Value() {
		return InvalidVerificationCodeError(ctx, op)
	}

	user, err := s.userQuery.GetByID(ctx, cmd.UserID)
	if err != nil {
		return UserNotFoundError(ctx, cmd.UserID, op)
	}

	if err := user.Activate(); err != nil {
		return err
	}

	err = s.userCommand.UpdateUserStatus(ctx, users.UpdateUserStatusCommand{
		ID:     cmd.UserID,
		Status: users.UserStatusActive,
	})
	return err
}

func (s *authService) Login(ctx context.Context, cmd LoginCommand) (SessionPayload, error) {
	const op = "Login"

	email, err := users.NewEmail(cmd.Email)
	if err != nil {
		return SessionPayload{}, InvalidCredentialsError(ctx, op)
	}

	user, err := s.userQuery.GetByEmail(ctx, email)
	if err != nil {
		return SessionPayload{}, InvalidCredentialsError(ctx, op)
	}

	if !s.passwordChecker.CheckPassword(user.HashedPassword, cmd.Password) {
		return SessionPayload{}, InvalidCredentialsError(ctx, op)
	}

	needs2FA := user.TwoFactorAuth.IsEnabled
	if err := user.ValidateLogin(needs2FA); err != nil {
		return SessionPayload{}, err
	}

	if needs2FA {
		if cmd.TwoFactorCode == nil || *cmd.TwoFactorCode == "" {
			return SessionPayload{}, &RequiresTwoFactorError{UserID: user.ID}
		}
		valid, err := s.twoFactorProvider.VerifyCode(ctx, user.ID, *cmd.TwoFactorCode)
		if err != nil || !valid {
			return SessionPayload{}, TwoFactorInvalidCodeError(ctx, op)
		}
	}

	return s.createSession(ctx, user, cmd.DeviceInfo, cmd.UserAgent, cmd.IPAddress)
}

func (s *authService) RefreshToken(ctx context.Context, cmd RefreshTokenCommand) (SessionPayload, error) {
	const op = "RefreshToken"

	if cmd.RefreshToken == "" {
		return SessionPayload{}, RefreshTokenRequiredError(ctx, op)
	}

	ok, err := s.tokenService.VerifyToken(ctx, cmd.RefreshToken, TokenTypeRefreshToken)
	if err != nil || !ok {
		return SessionPayload{}, InvalidTokenError(ctx, op)
	}

	_, err = s.sessionRepository.GetByRefreshToken(ctx, cmd.RefreshToken)
	if err != nil {
		return SessionPayload{}, TokenRevokedError(ctx, op)
	}

	claims, err := s.tokenService.GetTokenClaims(ctx, cmd.RefreshToken, TokenTypeRefreshToken)
	if err != nil {
		return SessionPayload{}, InvalidTokenError(ctx, op)
	}

	uidStr, _ := claims[ClaimUserID].(string)
	if uidStr == "" {
		return SessionPayload{}, InvalidTokenError(ctx, op)
	}

	userID, err := parseUserID(uidStr)
	if err != nil {
		return SessionPayload{}, InvalidTokenError(ctx, op)
	}

	user, err := s.userQuery.GetByID(ctx, userID)
	if err != nil {
		return SessionPayload{}, InvalidTokenError(ctx, op)
	}

	if err := user.ValidateLogin(false); err != nil {
		return SessionPayload{}, err
	}

	accessClaims := map[string]any{
		ClaimUserID: uidStr,
		ClaimEmail:  user.Email.String(),
		ClaimRole:   string(user.Role),
	}
	accessToken, err := s.tokenService.CreateToken(ctx, accessClaims, TokenTypeAccessToken)
	if err != nil {
		return SessionPayload{}, err
	}

	refreshToken, err := s.tokenService.CreateToken(ctx, map[string]any{
		ClaimUserID: uidStr,
	}, TokenTypeRefreshToken)
	if err != nil {
		return SessionPayload{}, err
	}

	_ = s.tokenService.RevokeToken(ctx, cmd.RefreshToken, TokenTypeRefreshToken)

	session := JwtSession{
		UserID:       uidStr,
		RefreshToken: refreshToken.Code,
		DeviceInfo:   "",
		UserAgent:    "",
		IPAddress:    "",
		ExpiresAt:    refreshToken.ExpiresAt,
		CreatedAt:    refreshToken.CreatedAt,
	}
	if err := s.sessionRepository.Create(ctx, &session); err != nil {
		return SessionPayload{}, err
	}

	return SessionPayload{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	}, nil
}

func (s *authService) Logout(ctx context.Context, cmd LogoutCommand) error {
	const op = "Logout"

	if cmd.RefreshToken == "" {
		return nil
	}

	session, err := s.sessionRepository.GetByRefreshToken(ctx, cmd.RefreshToken)
	if err != nil {
		return nil
	}

	userID, err := parseUserID(session.UserID)
	if err != nil {
		return nil
	}

	_ = s.tokenService.RevokeToken(ctx, cmd.RefreshToken, TokenTypeRefreshToken)
	return s.sessionRepository.DeleteUserSession(ctx, userID.Value(), cmd.RefreshToken)
}

func (s *authService) LogoutAll(ctx context.Context, cmd LogoutAllCommand) error {
	const op = "LogoutAll"

	_, err := s.userQuery.GetByID(ctx, cmd.UserID)
	if err != nil {
		return UserNotFoundError(ctx, cmd.UserID, op)
	}

	return s.sessionRepository.DeleteAllUserSessions(ctx, cmd.UserID.Value())
}

func (s *authService) VerifyTwoFactor(ctx context.Context, cmd VerifyTwoFactorCommand) (SessionPayload, error) {
	const op = "VerifyTwoFactor"

	user, err := s.userQuery.GetByID(ctx, cmd.UserID)
	if err != nil {
		return SessionPayload{}, InvalidCredentialsError(ctx, op)
	}

	if !user.TwoFactorAuth.IsEnabled {
		return SessionPayload{}, TwoFactorNotEnabledError(ctx, op)
	}

	valid, err := s.twoFactorProvider.VerifyCode(ctx, cmd.UserID, cmd.TwoFactorCode)
	if err != nil || !valid {
		return SessionPayload{}, TwoFactorInvalidCodeError(ctx, op)
	}

	return s.createSession(ctx, user, "", "", "")
}

func (s *authService) EnableTwoFactor(ctx context.Context, cmd EnableTwoFactorCommand) error {
	const op = "EnableTwoFactor"

	user, err := s.userQuery.GetByID(ctx, cmd.UserID)
	if err != nil {
		return UserNotFoundError(ctx, cmd.UserID, op)
	}

	if user.TwoFactorAuth.IsEnabled {
		return TwoFactorAlreadyEnabledError(ctx, op)
	}

	_, _, err = s.twoFactorProvider.GenerateSecret(ctx, cmd.UserID, cmd.Method)
	if err != nil {
		return err
	}

	return nil
}

func (s *authService) DisableTwoFactor(ctx context.Context, cmd DisableTwoFactorCommand) error {
	const op = "DisableTwoFactor"

	user, err := s.userQuery.GetByID(ctx, cmd.UserID)
	if err != nil {
		return UserNotFoundError(ctx, cmd.UserID, op)
	}

	if !user.TwoFactorAuth.IsEnabled {
		return TwoFactorNotEnabledError(ctx, op)
	}

	return s.twoFactorProvider.Disable(ctx, cmd.UserID)
}

func (s *authService) VerifyEmail(ctx context.Context, cmd VerifyEmailCommand) error {
	_ = cmd
	return nil
}

func (s *authService) RequestResetPassword(ctx context.Context, cmd RequestResetPasswordCommand) error {
	_, err := s.userQuery.GetByEmail(ctx, users.NewEmailNoErr(cmd.Email))
	if err != nil {
		return nil
	}

	// Token creation and email sending are infrastructure concerns; interface can be added later.
	return nil
}

func (s *authService) ResetPassword(ctx context.Context, cmd ResetPasswordCommand) error {
	const op = "ResetPassword"

	ok, err := s.tokenService.VerifyToken(ctx, cmd.Code, TokenTypePasswordReset)
	if err != nil || !ok {
		return InvalidVerificationCodeError(ctx, op)
	}

	claims, err := s.tokenService.GetTokenClaims(ctx, cmd.Code, TokenTypePasswordReset)
	if err != nil {
		return InvalidVerificationCodeError(ctx, op)
	}

	uidStr, _ := claims[ClaimUserID].(string)
	if uidStr == "" {
		return InvalidVerificationCodeError(ctx, op)
	}

	userID, err := parseUserID(uidStr)
	if err != nil {
		return InvalidVerificationCodeError(ctx, op)
	}

	if err := shared.ValidatePasswordStrength(cmd.Password); err != nil {
		return err
	}

	return s.userCommand.ResetPasswordByCode(ctx, users.ResetPasswordByCodeCommand{
		ID:          userID,
		NewPassword: cmd.Password,
	})
}

func (s *authService) UpdatePassword(ctx context.Context, cmd UpdatePasswordCommand) error {
	const op = "UpdatePassword"

	user, err := s.userQuery.GetByID(ctx, cmd.UserID)
	if err != nil {
		return UserNotFoundError(ctx, cmd.UserID, op)
	}

	if !s.passwordChecker.CheckPassword(user.HashedPassword, cmd.CurrentPassword) {
		return PasswordMismatchError(ctx, op)
	}

	if err := shared.ValidatePasswordStrength(cmd.NewPassword); err != nil {
		return err
	}

	return s.userCommand.UpdatePassword(ctx, users.UpdatePasswordCommand{
		ID:              cmd.UserID,
		CurrentPassword: cmd.CurrentPassword,
		NewPassword:     cmd.NewPassword,
	})
}

func (s *authService) DeleteAccount(ctx context.Context, userID shared.UserID) error {
	const op = "DeleteAccount"

	_, err := s.userQuery.GetByID(ctx, userID)
	if err != nil {
		return UserNotFoundError(ctx, userID, op)
	}

	_ = s.sessionRepository.DeleteAllUserSessions(ctx, userID.Value())
	return s.userCommand.DeleteUser(ctx, users.DeleteUserCommand{
		ID:           userID,
		IsHardDelete: false,
	})
}

func (s *authService) registerCustomer(ctx context.Context, cmd RegisterCommand) (Message, error) {
	return s.Register(ctx, cmd)
}

func (s *authService) registerEmployee(ctx context.Context, cmd RegisterCommand) (Message, error) {
	return s.Register(ctx, cmd)
}

func (s *authService) registerAdmin(ctx context.Context, cmd RegisterCommand) (Message, error) {
	return s.Register(ctx, cmd)
}

func (s *authService) createSession(ctx context.Context, user users.User, deviceInfo, userAgent, ipAddress string) (SessionPayload, error) {
	uidStr := user.ID.String()
	claims := map[string]any{
		ClaimUserID: uidStr,
		ClaimEmail:  user.Email.String(),
		ClaimRole:   string(user.Role),
	}

	accessToken, err := s.tokenService.CreateToken(ctx, claims, TokenTypeAccessToken)
	if err != nil {
		return SessionPayload{}, err
	}

	refreshToken, err := s.tokenService.CreateToken(ctx, map[string]any{ClaimUserID: uidStr}, TokenTypeRefreshToken)
	if err != nil {
		return SessionPayload{}, err
	}

	if err := s.tokenService.SaveToken(ctx, refreshToken); err != nil {
		return SessionPayload{}, err
	}

	session := JwtSession{
		UserID:       uidStr,
		RefreshToken: refreshToken.Code,
		DeviceInfo:   deviceInfo,
		UserAgent:    userAgent,
		IPAddress:    ipAddress,
		ExpiresAt:    refreshToken.ExpiresAt,
		CreatedAt:    refreshToken.CreatedAt,
	}
	if err := s.sessionRepository.Create(ctx, &session); err != nil {
		return SessionPayload{}, err
	}

	return SessionPayload{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	}, nil
}

func parseUserID(s string) (shared.UserID, error) {
	if s == "" {
		return shared.UserID{}, errors.New("user id is empty")
	}
	v, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return shared.UserID{}, fmt.Errorf("invalid user id: %w", err)
	}
	return shared.NewUserID(uint(v)), nil
}
