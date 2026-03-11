package mappers

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"clinic-vet-api/internal/core/users"
	"clinic-vet-api/internal/infrastructure/http/handlers/dtos"
	"clinic-vet-api/internal/shared"
	"clinic-vet-api/internal/shared/page"
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

func (m *UserHandlerMapper) ToResponsePage(userPage page.Page[users.User]) page.Page[dtos.UserResponse] {
	return page.MapItems(userPage, m.ToResponse)
}

// UserSearchRequestToSpecification maps a user search request DTO (from query params) to a UserSpecification.
// Parses comma-separated ids, emails, roles, statuses; optional date ranges; and pagination.
// Returns an error if any required parsing fails (e.g. invalid role or status).
func (m *UserHandlerMapper) UserSearchRequestToSpecification(req dtos.UserSearchRequest) (*users.UserSpecification, error) {
	spec := &users.UserSpecification{}

	// IDs (comma-separated)
	if req.IDs != "" {
		parts := splitAndTrim(req.IDs, ",")
		ids := make([]shared.UserID, 0, len(parts))
		for _, p := range parts {
			id, err := shared.ParseUserIDFromString(p)
			if err != nil {
				return nil, err
			}
			ids = append(ids, id)
		}
		spec.WithIDs(ids...)
	}

	// Emails (comma-separated)
	if req.Emails != "" {
		parts := splitAndTrim(req.Emails, ",")
		emails := make([]users.Email, 0, len(parts))
		for _, p := range parts {
			email, err := users.NewEmail(p)
			if err != nil {
				return nil, err
			}
			emails = append(emails, email)
		}
		spec.WithEmails(emails...)
	}

	// Roles (comma-separated)
	if req.Roles != "" {
		parts := splitAndTrim(req.Roles, ",")
		roles := make([]users.UserRole, 0, len(parts))
		for _, p := range parts {
			role, err := users.ParseUserRole(p)
			if err != nil {
				return nil, err
			}
			roles = append(roles, role)
		}
		spec.WithRoles(roles...)
	}

	// Statuses (comma-separated)
	if req.Statuses != "" {
		parts := splitAndTrim(req.Statuses, ",")
		statuses := make([]users.UserStatus, 0, len(parts))
		for _, p := range parts {
			status, err := users.ParseUserStatus(p)
			if err != nil {
				return nil, err
			}
			statuses = append(statuses, status)
		}
		spec.WithStatuses(statuses...)
	}

	// IsActive (optional boolean)
	if v := parseBoolPtr(req.IsActive); v != nil {
		spec.WithIsActive(*v)
	}

	// Created date range
	if req.CreatedAfter != "" || req.CreatedBefore != "" {
		var after, before *time.Time
		if req.CreatedAfter != "" {
			t, err := parseDate(req.CreatedAfter)
			if err != nil {
				return nil, err
			}
			after = &t
		}
		if req.CreatedBefore != "" {
			t, err := parseDate(req.CreatedBefore)
			if err != nil {
				return nil, err
			}
			before = &t
		}
		spec.WithCreatedDateRange(after, before)
	}

	// Last login date range
	if req.LastLoginAfter != "" || req.LastLoginBefore != "" {
		var after, before *time.Time
		if req.LastLoginAfter != "" {
			t, err := parseDate(req.LastLoginAfter)
			if err != nil {
				return nil, err
			}
			after = &t
		}
		if req.LastLoginBefore != "" {
			t, err := parseDate(req.LastLoginBefore)
			if err != nil {
				return nil, err
			}
			before = &t
		}
		spec.WithLastLoginRange(after, before)
	}

	// Search term
	if req.SearchTerm != "" {
		spec.WithSearchTerm(strings.TrimSpace(req.SearchTerm))
	}

	// HasTwoFactor (optional boolean)
	if v := parseBoolPtr(req.HasTwoFactor); v != nil {
		spec.WithTwoFactor(*v)
	}

	// Pagination (apply defaults so spec always has valid pagination)
	pagi := req.PaginationRequest.WithDefaults().ToPagination()
	spec.WithPagination(pagi.Number, pagi.Size, pagi.OrderBy, pagi.SortDir)

	return spec, nil
}

func splitAndTrim(s, sep string) []string {
	parts := strings.Split(s, sep)
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if t := strings.TrimSpace(p); t != "" {
			out = append(out, t)
		}
	}
	return out
}

func parseBoolPtr(s string) *bool {
	s = strings.TrimSpace(strings.ToLower(s))
	if s == "" {
		return nil
	}
	v, err := strconv.ParseBool(s)
	if err != nil {
		return nil
	}
	return &v
}

func parseDate(s string) (time.Time, error) {
	s = strings.TrimSpace(s)
	formats := []string{time.RFC3339, "2006-01-02", "2006-01-02T15:04:05Z07:00"}
	for _, f := range formats {
		if t, err := time.Parse(f, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("invalid date format: %q (use 2006-01-02 or RFC3339)", s)
}
