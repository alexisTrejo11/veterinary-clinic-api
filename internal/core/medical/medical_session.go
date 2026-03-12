package medical

import (
	"context"
	"time"

	"clinic-vet-api/internal/core/appointments"
	"clinic-vet-api/internal/core/customers"
	"clinic-vet-api/internal/core/employees"
	"clinic-vet-api/internal/core/pets"
	"clinic-vet-api/internal/shared"
)

type MedSessionID struct{ shared.BaseID }
type DewormID struct{ shared.BaseID }

func NewDewormID(value uint) DewormID {
	return DewormID{shared.BaseID{Value: value}}
}

type MedicalSession struct {
	shared.Entity[MedSessionID]
	CustomerID  customers.CustomerID
	VisitReason appointments.VisitReason
	VisitType   appointments.VisitType
	VisitDate   time.Time
	Notes       *string
	EmployeeID  employees.EmployeeID
	SessionSummary
}

type SessionSummary struct {
	PetID           pets.PetID
	Weight          *shared.Decimal
	HeartRate       *int
	RespiratoryRate *int
	Temperature     *shared.Decimal
	Diagnosis       string
	Treatment       string
	Condition       appointments.PetCondition
	Medications     []string
	FollowUpDate    *time.Time
	Symptoms        []string
}

func NewMedSessionID(value uint) MedSessionID {
	return MedSessionID{shared.BaseID{Value: value}}
}

// Validate validates the medical session business rules
func (ms *MedicalSession) Validate(ctx context.Context) error {
	operation := "ValidateMedicalSession"

	if err := ms.validateRequiredFields(ctx, operation); err != nil {
		return err
	}

	if err := ms.validateEnums(ctx, operation); err != nil {
		return err
	}

	if err := ms.validateVitals(ctx, operation); err != nil {
		return err
	}

	if err := ms.validateVisitDateConstraints(ctx, operation); err != nil {
		return err
	}

	if err := ms.validateMedications(ctx, operation); err != nil {
		return err
	}

	return nil
}

// validateRequiredFields validates that all required fields are present
func (ms *MedicalSession) validateRequiredFields(ctx context.Context, operation string) error {
	if ms.VisitDate.IsZero() {
		return MissingFieldError(ctx, FieldVisitDate, operation)
	}

	if ms.VisitReason == "" {
		return MissingFieldError(ctx, FieldVisitReason, operation)
	}

	if ms.VisitType == "" {
		return MissingFieldError(ctx, FieldVisitType, operation)
	}

	return nil
}

// validateEnums validates that enum fields have valid values
func (ms *MedicalSession) validateEnums(ctx context.Context, operation string) error {
	if !ms.VisitReason.IsValid() {
		return InvalidEnumError(ctx, FieldVisitReason, string(ms.VisitReason), operation)
	}

	if !ms.VisitType.IsValid() {
		return InvalidEnumError(ctx, FieldVisitType, string(ms.VisitType), operation)
	}

	if ms.Condition != "" && !ms.Condition.IsValid() {
		return InvalidEnumError(ctx, FieldCondition, string(ms.Condition), operation)
	}

	return nil
}

// validateVitals validates that vital signs have valid values
func (ms *MedicalSession) validateVitals(ctx context.Context, operation string) error {
	if ms.Weight != nil && (ms.Weight.IsNegative() || ms.Weight.IsZero()) {
		return InvalidPositiveValueError(ctx, FieldWeight, operation)
	}

	if ms.Temperature != nil && (ms.Temperature.IsNegative() || ms.Temperature.IsZero()) {
		return InvalidPositiveValueError(ctx, FieldTemperature, operation)
	}

	if ms.HeartRate != nil && *ms.HeartRate <= 0 {
		return InvalidPositiveValueError(ctx, FieldHeartRate, operation)
	}

	if ms.RespiratoryRate != nil && *ms.RespiratoryRate <= 0 {
		return InvalidPositiveValueError(ctx, FieldRespiratoryRate, operation)
	}

	return nil
}

// validateVisitDateConstraints validates visit date business rules
func (ms *MedicalSession) validateVisitDateConstraints(ctx context.Context, operation string) error {
	now := time.Now()
	maxFutureDate := now.AddDate(MaxFutureYears, 0, 0)
	minDate := time.Date(MinVisitYear, 1, 1, 0, 0, 0, 0, time.UTC)

	if ms.VisitDate.After(maxFutureDate) {
		return VisitDateTooFarInFutureError(ctx, operation)
	}

	if ms.VisitDate.Before(minDate) {
		return VisitDateTooOldError(ctx, operation)
	}

	return nil
}

// validateMedications validates medications list
func (ms *MedicalSession) validateMedications(ctx context.Context, operation string) error {
	if ms.Medications != nil && len(ms.Medications) > MaxMedicationsCount {
		return TooManyMedicationsError(ctx, len(ms.Medications), operation)
	}

	return nil
}
