// Package base contains the base schema to apply in each entity
package shared

import (
	"fmt"
	"strings"
	"time"
)

type Entity[T any] struct {
	ID        T
	CreatedAt time.Time
	UpdatedAt time.Time
	Version   int
}

// NewEntity creates a new base entity
func NewEntity[T any](id T, createdAt, updatedAt time.Time, version int) Entity[T] {
	return Entity[T]{
		ID:        id,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		Version:   version,
	}
}

func CreateEntity[T any](id T) Entity[T] {
	return Entity[T]{
		ID:        id,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Version:   1,
	}
}

// IncrementVersion updates the version and updatedAt timestamp
func (e *Entity[T]) IncrementVersion() {
	e.Version++
	e.UpdatedAt = time.Now()
}

func (e *Entity[T]) SetTimeStamps(createAt, updateAt time.Time) {
	e.CreatedAt = createAt
	e.UpdatedAt = updateAt
}

func (e *Entity[T]) SetID(id T) {
	e.ID = id
}

// PersonGender represents the gender of a person
type PersonGender string

type (
	// EmployeeID is a unique identifier for an employee user
	EmployeeID struct{ BaseID }
	// UserID is a unique identifier for a user
	UserID struct{ BaseID }
)

const (
	GenderMale         PersonGender = "male"
	GenderFemale       PersonGender = "female"
	GenderNotSpecified PersonGender = "not_specified"
	GenderOther        PersonGender = "other"
)

var (
	ValidPersonGenders = []PersonGender{
		GenderMale,
		GenderFemale,
		GenderNotSpecified,
		GenderOther,
	}

	personGenderMap = map[string]PersonGender{
		"male":          GenderMale,
		"female":        GenderFemale,
		"not_specified": GenderNotSpecified,
		"not specified": GenderNotSpecified,
		"other":         GenderOther,
		"":              GenderNotSpecified,
	}

	personGenderDisplayNames = map[PersonGender]string{
		GenderMale:         "Male",
		GenderFemale:       "Female",
		GenderNotSpecified: "Not Specified",
		GenderOther:        "Other",
	}
)

// ============================================================================
// PersonGender Methods
// ============================================================================

func (g PersonGender) IsValid() bool {
	_, exists := personGenderMap[string(g)]
	return exists
}

func ParseGender(gender string) (PersonGender, error) {
	normalized := normalizeGenderInput(gender)
	if val, exists := personGenderMap[normalized]; exists {
		return val, nil
	}
	return GenderNotSpecified, fmt.Errorf("invalid gender: %s", gender)
}

func MustParseGender(gender string) PersonGender {
	parsed, err := ParseGender(gender)
	if err != nil {
		panic(err)
	}
	return parsed
}

func (g PersonGender) String() string {
	return string(g)
}

func (g PersonGender) DisplayName() string {
	if displayName, exists := personGenderDisplayNames[g]; exists {
		return displayName
	}
	return "Unknown Gender"
}

func (g PersonGender) Values() []PersonGender {
	return ValidPersonGenders
}

func GetAllGenders() []PersonGender {
	return ValidPersonGenders
}

func normalizeGenderInput(input string) string {
	input = strings.TrimSpace(strings.ToLower(input))
	input = strings.ReplaceAll(input, " ", "_")
	return input
}

func NormalizeInput(input string) string {
	return fmt.Sprintf("%s", input)
}
