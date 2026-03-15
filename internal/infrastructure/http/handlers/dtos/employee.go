package dtos

import (
	"time"

	"clinic-vet-api/internal/shared/page"
)

// EmployeeResponse represents an employee in HTTP responses
// @Description Employee entity: id, name, gender, date of birth, photo, license, specialty, years experience, is_active, user_id, timestamps.
type EmployeeResponse struct {
	ID              uint      `json:"id"`
	FirstName       string    `json:"first_name"`
	LastName        string    `json:"last_name"`
	Gender          string    `json:"gender"`
	DateOfBirth     time.Time `json:"date_of_birth"`
	Photo           string    `json:"photo,omitempty"`
	LicenseNumber   string    `json:"license_number"`
	Specialty       string    `json:"specialty"`
	YearsExperience int       `json:"years_experience"`
	IsActive        bool      `json:"is_active"`
	UserID          uint      `json:"user_id,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// EmployeeStatsResponse represents aggregate employee statistics
// @Description Total employees, active count, and count per specialty.
type EmployeeStatsResponse struct {
	TotalEmployees  int64            `json:"total_employees"`
	ActiveEmployees int64            `json:"active_employees"`
	Specialties     map[string]int64 `json:"specialties"`
}

// EmployeeSearchRequest represents filters for listing/searching employees
// @Description Pagination plus optional specialty and is_active filter.
type EmployeeSearchRequest struct {
	page.PaginationRequest

	Specialty string `form:"specialty,omitempty" json:"specialty,omitempty"`
	IsActive  *bool  `form:"is_active,omitempty" json:"is_active,omitempty"`
}

// EmployeeScheduleRequest represents schedule data in create/update requests
// @Description Day of week, entry/departure time (0-23), optional break window.
type EmployeeScheduleRequest struct {
	Day           string `json:"day" binding:"required"`                     // e.g. "monday"
	EntryTime     int    `json:"entry_time" binding:"required,gte=0,lte=23"` // 0-23
	DepartureTime int    `json:"departure_time" binding:"required,gte=0,lte=23,gtfield=EntryTime"`
	StartBreak    int    `json:"start_break" binding:"omitempty,gte=0,lte=23"`
	EndBreak      int    `json:"end_break" binding:"omitempty,gte=0,lte=23"`
}

// EmployeeCreateRequest represents the body for creating a new employee
// @Description Name, gender, date of birth, photo, license number, years experience, specialty, schedule, optional is_active.
type EmployeeCreateRequest struct {
	FirstName   string                 `json:"first_name" binding:"required"`
	LastName    string                 `json:"last_name" binding:"required"`
	Gender      string                 `json:"gender" binding:"required,oneof=male female not_specified other"`
	DateOfBirth string                 `json:"date_of_birth" binding:"required,datetime=2006-01-02"`
	Photo       string                 `json:"photo,omitempty" binding:"omitempty,url"`
	LicenseNo   string                 `json:"license_number" binding:"required"`
	YearsExp    int32                  `json:"years_experience" binding:"gte=0,lte=60"`
	IsActive    *bool                  `json:"is_active,omitempty"`
	Specialty   string                 `json:"specialty" binding:"required"`
	Schedule    EmployeeScheduleRequest `json:"schedule" binding:"required,dive"`
}

// EmployeeUpdateRequest represents the body for updating an existing employee
// @Description ID (required) plus optional fields and optional schedule.
type EmployeeUpdateRequest struct {
	ID          uint                   `json:"id" binding:"required"`
	FirstName   *string                `json:"first_name,omitempty"`
	LastName    *string                `json:"last_name,omitempty"`
	Gender      *string                `json:"gender,omitempty" binding:"omitempty,oneof=male female not_specified other"`
	DateOfBirth *string                `json:"date_of_birth,omitempty" binding:"omitempty,datetime=2006-01-02"`
	Photo       *string                `json:"photo,omitempty" binding:"omitempty,url"`
	LicenseNo   *string                `json:"license_number,omitempty"`
	YearsExp    *int32                 `json:"years_experience,omitempty" binding:"omitempty,gte=0,lte=60"`
	IsActive    *bool                  `json:"is_active,omitempty"`
	Specialty   *string                `json:"specialty,omitempty"`
	Schedule    *EmployeeScheduleRequest `json:"schedule,omitempty"`
}

