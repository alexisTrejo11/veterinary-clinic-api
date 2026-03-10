package shared

import (
	"clinic-vet-api/internal/shared/errors"
	"context"
	"time"
)

type Person struct {
	FirstName   string
	LastName    string
	DateOfBirth time.Time
	Gender      PersonGender
}

func PersonValidationError(ctx context.Context, code, field, value, operation string) error {
	return errors.ValidationError(ctx, code, field, value, operation)

}

func (p *Person) Validate(ctx context.Context, operation string) error {
	PersonDOBRequired := "PERSON_DATA_FIELD_REQUIRED"
	PersonGenderInvalid := "INVALID_GENDER"
	PersonDOBFuture := "PERSON_DOB_IN_FUTURE"
	PersonUnderage := "PERSON_UNDERAGE"

	if p.FirstName == "" {
		return PersonValidationError(ctx,
			PersonDOBRequired, "first_name", "", operation)
	}

	if p.LastName == "" {
		return PersonValidationError(ctx,
			PersonDOBRequired, "last_name", "", operation)
	}

	if p.DateOfBirth.IsZero() {
		return PersonValidationError(ctx,
			PersonDOBRequired, "date_of_birth", "", operation)
	}

	if p.DateOfBirth.After(time.Now()) {
		return PersonValidationError(ctx,
			PersonDOBFuture, "date_of_birth", p.DateOfBirth.Format(time.RFC3339), operation)
	}

	// Validar mayoría de edad (18 años)
	minAgeDate := time.Now().AddDate(-18, 0, 0)
	if p.DateOfBirth.After(minAgeDate) {
		return PersonValidationError(ctx,
			PersonUnderage, "date_of_birth", p.DateOfBirth.Format(time.RFC3339), operation)
	}

	if !p.Gender.IsValid() {
		return PersonValidationError(ctx,
			PersonGenderInvalid, "gender", string(p.Gender), operation)
	}

	return nil
}

func (p *Person) Age() int {
	now := time.Now()
	age := now.Year() - p.DateOfBirth.Year()

	if now.YearDay() < p.DateOfBirth.YearDay() {
		age--
	}
	return age
}

func (p *Person) IsAdult() bool {
	return p.Age() >= 18
}
