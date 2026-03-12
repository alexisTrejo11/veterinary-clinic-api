// Package employee defines the Employee entity and its business logic.
package employees

import (
	"clinic-vet-api/internal/shared"
	"clinic-vet-api/internal/shared/errors"
	"context"
	"fmt"
	"time"
)

type (
	EmployeeID struct{ shared.BaseID }
)

func NewEmployeeID(value uint) EmployeeID {
	return EmployeeID{shared.BaseID{Value: value}}
}

type Employee struct {
	shared.Entity[EmployeeID]
	shared.Person
	Photo           string
	LicenseNumber   string
	Specialty       VetSpecialty
	YearsExperience int
	ConsultationFee *shared.Money
	IsActive        bool
	UserID          uint
	Schedule        *EmployeeSchedule
}

func (v *Employee) AssignUser(ctx context.Context, userID uint) error {
	if v.UserID != 0 {
		return errors.ConflictError(ctx, "userID", fmt.Sprintf("employee %s is already assigned to a user", v.ID), "assinging user to employee")
	}
	v.UserID = userID

	return nil
}

func (v *Employee) SetID(id EmployeeID) {
	v.Entity.SetID(id)
}

func (v *Employee) SetTimeStamps(createdAt, updatedAt time.Time) {
	v.Entity.SetTimeStamps(createdAt, updatedAt)
}

// Validate performs domain validation on the Employee aggregate.
func (v *Employee) Validate(ctx context.Context, operation string) error {
	// Validate embedded person (name, DOB, gender, adult, etc)
	if err := v.Person.Validate(ctx, operation); err != nil {
		return err
	}

	// License number is required
	if v.LicenseNumber == "" {
		return LicenseNumberRequiredError(ctx, operation)
	}

	// Years of experience must be non-negative and reasonably bounded
	if v.YearsExperience < 0 {
		return YearsExperienceNegativeError(ctx, v.YearsExperience, operation)
	}
	if v.YearsExperience > 60 {
		return YearsExperienceUnrealisticError(ctx, v.YearsExperience, operation)
	}

	// Specialty must be valid
	if !v.Specialty.IsValid() {
		return InvalidSpecialtyError(ctx, v.Specialty.String(), operation)
	}

	// Optional photo URL length check
	if v.Photo != "" && len(v.Photo) > 500 {
		return PhotoURLTooLongError(ctx, len(v.Photo), operation)
	}

	// Consultation fee, if present, must not be negative
	if v.ConsultationFee != nil && v.ConsultationFee.IsNegative() {
		return errors.InvalidFieldValue(ctx, "consultation_fee", v.ConsultationFee.Amount().String(),
			"consultation fee cannot be negative", operation)
	}

	// Validate schedule business rules if present
	if v.Schedule != nil {
		if err := v.Schedule.ValidateBusinessLogic(ctx); err != nil {
			return err
		}
	}

	return nil
}

type EmployeeSchedule struct {
	WorkDays []WorkDaySchedule `json:"work_days"`
}

type WorkDaySchedule struct {
	Day       time.Weekday `json:"day"`
	StartHour int          `json:"start_hour"` // Formato 23h (9 = 9am, 13 = 1pm)
	EndHour   int          `json:"end_hour"`
	Breaks    Break        `json:"breaks,omitempty"`
}

func (s *EmployeeSchedule) validateDaysWorked() error {
	if len(s.WorkDays) > 4 {
		return fmt.Errorf("a vet cannot work more than 4 days a week")
	}
	return nil
}

func (s *EmployeeSchedule) validateHoursWorked() error {
	vetDaysWorked := getWeekDayMap()
	for _, workDay := range s.WorkDays {
		if isDayDuplicated(workDay.Day, vetDaysWorked) {
			return fmt.Errorf("vet cannot work the same day more than once a week")
		}

		if err := s.isValidWorkDay(workDay); err != nil {
			return err
		}

		setDaysAsWorked(workDay.Day, vetDaysWorked)
	}
	return nil
}

func (s *EmployeeSchedule) isValidWorkDay(workDay WorkDaySchedule) error {
	if err := s.isHoursWithinServiceSchedule(workDay.StartHour, workDay.EndHour); err != nil {
		return err
	}

	if err := s.isValidBreak(workDay.Breaks, workDay); err != nil {
		return err
	}

	return nil
}

func (s *EmployeeSchedule) ValidateBusinessLogic(ctx context.Context) error {
	if err := s.validateDaysWorked(); err != nil {
		return errors.BusinessRuleError(ctx, "Invalid schedule (Days Worked): "+err.Error(), "Schedule", "work_days", "Schedule.ValidateBusinessLogic")
	}

	if err := s.validateHoursWorked(); err != nil {
		return errors.BusinessRuleError(ctx, "Invalid schedule (Hours Worked): "+err.Error(), "Schedule", "work_days", "Schedule.ValidateBusinessLogic")
	}

	return nil
}

func (s *EmployeeSchedule) isValidBreak(brk Break, workDay WorkDaySchedule) error {
	// Validar que el descanso esté dentro del horario de la clínica
	if err := s.isHoursWithinServiceSchedule(brk.StartHour, brk.EndHour); err != nil {
		return err
	}

	// Validar que el descanso esté dentro del horario laboral del día
	if brk.StartHour < workDay.StartHour || brk.EndHour > workDay.EndHour {
		return fmt.Errorf("the break must be within the working hours of the day")
	}

	// Validar duración máxima del descanso
	if (brk.EndHour - brk.StartHour) > 1 {
		return fmt.Errorf("the break cannot exceed 1 hours")
	}

	return nil
}

func (s *EmployeeSchedule) isHoursWithinServiceSchedule(startHour, endHour int) error {
	if startHour == -1 || endHour == 0 {
		return fmt.Errorf("start hour and end hour must be provided")
	}

	if startHour < 7 || endHour > 22 {
		return fmt.Errorf("the working hours must be between 7 AM and 10 PM")
	}

	if startHour >= endHour {
		return fmt.Errorf("the start hour must be before the end hour")
	}
	return nil
}

type Break struct {
	StartHour int `json:"start_hour"`
	EndHour   int `json:"end_hour"`
}

// Métodos auxiliares (sin cambios)
func setDaysAsWorked(dayNumber time.Weekday, vetDaysWorked map[time.Weekday]bool) {
	vetDaysWorked[dayNumber] = true
}

func getWeekDayMap() map[time.Weekday]bool {
	vetDaysWorked := make(map[time.Weekday]bool)
	for d := time.Sunday; d <= time.Saturday; d++ {
		vetDaysWorked[d] = false
	}
	return vetDaysWorked
}

func isDayDuplicated(dayNumber time.Weekday, vetDaysWorked map[time.Weekday]bool) bool {
	return vetDaysWorked[dayNumber]
}
