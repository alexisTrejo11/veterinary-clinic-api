package addresses

import (
	"context"
	"fmt"

	"clinic-vet-api/internal/shared"
	"clinic-vet-api/internal/shared/errors"
)

// Entity name and fields for error reporting
const (
	addressEntity = "address"
)

const (
	FieldStreet       = "street"
	FieldCity         = "city"
	FieldState        = "state"
	FieldZipCode      = "zip_code"
	FieldCountry      = "country"
	FieldBuildingType = "building_type"
	FieldCustomerID   = "customer_id"
)

// Limits / business constants
const (
	MaxStreetLen  = 255
	MaxCityLen    = 150
	MaxStateLen   = 150
	MaxZipLen     = 20
	MaxCountryLen = 100
)

// ======================================================================
// Validation errors
// ======================================================================

func StreetRequiredError(ctx context.Context, operation string) error {
	return errors.ValidationError(ctx, addressEntity, FieldStreet,
		"street is required", operation)
}

func CityRequiredError(ctx context.Context, operation string) error {
	return errors.ValidationError(ctx, addressEntity, FieldCity,
		"city is required", operation)
}

func StateRequiredError(ctx context.Context, operation string) error {
	return errors.ValidationError(ctx, addressEntity, FieldState,
		"state is required", operation)
}

func ZipRequiredError(ctx context.Context, operation string) error {
	return errors.ValidationError(ctx, addressEntity, FieldZipCode,
		"zip code is required", operation)
}

func CountryRequiredError(ctx context.Context, operation string) error {
	return errors.ValidationError(ctx, addressEntity, FieldCountry,
		"country is required", operation)
}

func InvalidCountryError(ctx context.Context, country Country, operation string) error {
	return errors.InvalidEnumValue(ctx, FieldCountry, string(country),
		"unsupported country", operation)
}

func InvalidBuildingTypeError(ctx context.Context, bt BuildingType, operation string) error {
	return errors.InvalidEnumValue(ctx, FieldBuildingType, string(bt),
		"unsupported building type", operation)
}

// ======================================================================
// Business rule errors
// ======================================================================

func MaxAddressesPerCustomerError(ctx context.Context, customerID shared.UserID, max int, operation string) error {
	return errors.BusinessRuleError(
		ctx,
		fmt.Sprintf("customer %s already has the maximum of %d addresses", customerID.String(), max),
		addressEntity,
		FieldCustomerID,
		operation,
	)
}

func DefaultAddressRequiredError(ctx context.Context, customerID shared.UserID, operation string) error {
	return errors.BusinessRuleError(
		ctx,
		fmt.Sprintf("customer %d must have at least one default address", customerID),
		addressEntity,
		FieldCustomerID,
		operation,
	)
}
