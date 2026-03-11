package pets

import (
	"context"
	"fmt"

	"clinic-vet-api/internal/shared/errors"
)

// Entity and field names for error reporting
const (
	EntityName = "pet"
)

const (
	FieldName       = "name"
	FieldSpecies    = "species"
	FieldGender     = "gender"
	FieldAge        = "age"
	FieldBreed      = "breed"
	FieldMicrochip  = "microchip"
	FieldColor      = "color"
	FieldPhoto      = "photo"
	FieldCustomerID = "customer_id"
)

// Validation limits (must match domain rules)
const (
	MaxNameLen      = 255
	MaxBreedLen     = 100
	MaxMicrochipLen = 50
	MaxColorLen     = 50
	MaxPhotoURLLen  = 2048
	MaxPetAgeYears  = 50
)

// =========================================================================
// Enum / parse errors
// =========================================================================

// InvalidSpeciesError returns an error for invalid pet species
func InvalidSpeciesError(ctx context.Context, value, operation string) error {
	return errors.InvalidEnumValue(ctx, FieldSpecies, value,
		fmt.Sprintf("allowed values: %v", ValidPetSpeciess), operation)
}

// InvalidGenderError returns an error for invalid pet gender
func InvalidGenderError(ctx context.Context, value, operation string) error {
	return errors.InvalidEnumValue(ctx, FieldGender, value,
		fmt.Sprintf("allowed values: %v", ValidPetGenders), operation)
}

// =========================================================================
// Validation errors (required, length, range)
// =========================================================================

// NameRequiredError returns an error when pet name is empty
func NameRequiredError(ctx context.Context, operation string) error {
	return errors.ValidationError(ctx, EntityName, FieldName,
		"name is required", operation)
}

// NameTooLongError returns an error when pet name exceeds max length
func NameTooLongError(ctx context.Context, length int, operation string) error {
	return errors.ValidationError(ctx, EntityName, FieldName,
		fmt.Sprintf("name too long (%d characters, max %d)", length, MaxNameLen), operation)
}

// SpeciesRequiredError returns an error when species is not set
func SpeciesRequiredError(ctx context.Context, operation string) error {
	return errors.ValidationError(ctx, EntityName, FieldSpecies,
		"species is required", operation)
}

// AgeInvalidError returns an error when age is negative
func AgeInvalidError(ctx context.Context, age int, operation string) error {
	return errors.ValidationError(ctx, EntityName, FieldAge,
		fmt.Sprintf("age cannot be negative, got: %d", age), operation)
}

// AgeUnrealisticError returns an error when age exceeds reasonable max
func AgeUnrealisticError(ctx context.Context, age int, operation string) error {
	return errors.ValidationError(ctx, EntityName, FieldAge,
		fmt.Sprintf("age seems unrealistic (%d years, max %d)", age, MaxPetAgeYears), operation)
}

// BreedTooLongError returns an error when breed string is too long
func BreedTooLongError(ctx context.Context, length int, operation string) error {
	return errors.ValidationError(ctx, EntityName, FieldBreed,
		fmt.Sprintf("breed too long (%d characters, max %d)", length, MaxBreedLen), operation)
}

// MicrochipTooLongError returns an error when microchip string is too long
func MicrochipTooLongError(ctx context.Context, length int, operation string) error {
	return errors.ValidationError(ctx, EntityName, FieldMicrochip,
		fmt.Sprintf("microchip too long (%d characters, max %d)", length, MaxMicrochipLen), operation)
}

// ColorTooLongError returns an error when color string is too long
func ColorTooLongError(ctx context.Context, length int, operation string) error {
	return errors.ValidationError(ctx, EntityName, FieldColor,
		fmt.Sprintf("color too long (%d characters, max %d)", length, MaxColorLen), operation)
}

// PhotoURLTooLongError returns an error when photo URL is too long
func PhotoURLTooLongError(ctx context.Context, length int, operation string) error {
	return errors.ValidationError(ctx, EntityName, FieldPhoto,
		fmt.Sprintf("photo URL too long (%d characters, max %d)", length, MaxPhotoURLLen), operation)
}

// CustomerIDRequiredError returns an error when customer_id is missing
func CustomerIDRequiredError(ctx context.Context, operation string) error {
	return errors.ValidationError(ctx, EntityName, FieldCustomerID,
		"customer ID is required", operation)
}

func CustomerNotActiveError(ctx context.Context, customerID uint, operation string) error {
	return errors.BusinessRuleError(ctx, "Customer must be active to ad pets", "Customer",
		"", operation)
}
