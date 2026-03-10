package medical

import (
	"context"
	"fmt"

	"clinic-vet-api/internal/shared/errors"
)

// Error codes for medical domain
const (
	EntityName = "medical_session"
)

// Field names
const (
	FieldVisitDate       = "visit_date"
	FieldVisitReason     = "visit_reason"
	FieldVisitType       = "visit_type"
	FieldWeight          = "weight"
	FieldTemperature     = "temperature"
	FieldHeartRate       = "heart_rate"
	FieldRespiratoryRate = "respiratory_rate"
	FieldCondition       = "condition"
	FieldMedications     = "medications"
)

// Validation constraints
const (
	MinVisitYear        = 2015
	MaxFutureYears      = 1
	MaxMedicationsCount = 100
)

// =========================================================================
// Validation errors
// =========================================================================

// MissingFieldError creates an error for missing required fields
func MissingFieldError(ctx context.Context, field, operation string) error {
	return errors.ValidationError(ctx, EntityName, field,
		fmt.Sprintf("%s is required", field), operation)
}

// InvalidPositiveValueError creates an error for fields that must be positive
func InvalidPositiveValueError(ctx context.Context, field, operation string) error {
	return errors.ValidationError(ctx, EntityName, field,
		fmt.Sprintf("%s must be a positive value", field), operation)
}

// InvalidEnumError creates an error for invalid enum values
func InvalidEnumError(ctx context.Context, field, value, operation string) error {
	return errors.ValidationError(ctx, EntityName, field,
		fmt.Sprintf("invalid %s: %s", field, value), operation)
}

// VisitDateTooFarInFutureError creates an error when visit date is too far in the future
func VisitDateTooFarInFutureError(ctx context.Context, operation string) error {
	return errors.BusinessRuleError(ctx, "visit_date_future",
		EntityName,
		fmt.Sprintf("visit date cannot be more than %d year in the future", MaxFutureYears),
		operation)
}

// VisitDateTooOldError creates an error when visit date is before minimum allowed year
func VisitDateTooOldError(ctx context.Context, operation string) error {
	return errors.BusinessRuleError(ctx, "visit_date_past",
		EntityName,
		fmt.Sprintf("visit date cannot be before year %d", MinVisitYear),
		operation)
}

// TooManyMedicationsError creates an error when medications list exceeds maximum
func TooManyMedicationsError(ctx context.Context, count int, operation string) error {
	return errors.ValidationError(ctx, EntityName, FieldMedications,
		fmt.Sprintf("medications list exceeds maximum of %d items (got %d)", MaxMedicationsCount, count),
		operation)
}
