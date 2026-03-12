package employees

import (
	"clinic-vet-api/internal/shared/errors"
	"context"
	"fmt"
)

const (
	employeeEntity       = "employee"
	FieldLicenseNumber   = "license_number"
	FieldYearsExperience = "years_experience"
	FieldSpecialty       = "specialty"
	FieldPhoto           = "photo"
)

// InvalidEnumParserError keeps compatibility for enum parse helpers.
func InvalidEnumParserError(enumName, enumValue string) error {
	return errors.InvalidEnumValue(
		context.Background(),
		enumName,
		enumValue,
		fmt.Sprintf("invalid %s value: %s", enumName, enumValue),
		"enum parse",
	)
}

// Validation errors

func LicenseNumberRequiredError(ctx context.Context, operation string) error {
	return errors.ValidationError(ctx, FieldLicenseNumber, "",
		"license number is required", operation)
}

func YearsExperienceNegativeError(ctx context.Context, years int, operation string) error {
	return errors.ValidationError(ctx, FieldYearsExperience, fmt.Sprintf("%d", years),
		"years of experience cannot be negative", operation)
}

func YearsExperienceUnrealisticError(ctx context.Context, years int, operation string) error {
	return errors.ValidationError(ctx, FieldYearsExperience, fmt.Sprintf("%d", years),
		"years of experience seems unrealistic", operation)
}

func InvalidSpecialtyError(ctx context.Context, specialty string, operation string) error {
	return errors.InvalidEnumValue(ctx, FieldSpecialty, specialty,
		"invalid veterinary specialty", operation)
}

func PhotoURLTooLongError(ctx context.Context, length int, operation string) error {
	return errors.ValidationError(ctx, FieldPhoto, fmt.Sprintf("%d", length),
		"photo URL is too long (max 500 characters)", operation)
}

