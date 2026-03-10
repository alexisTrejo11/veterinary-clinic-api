package pets

import (
	"clinic-vet-api/internal/shared"
	"strings"
)

type Pet struct {
	shared.Entity[PetID]
	Name        string
	Photo       *string
	Species     PetSpecies
	Breed       *string
	Age         *int
	Gender      *PetGender
	Color       *string
	MicrochipID *string
	BloodType   *string
	IsNeutered  *bool
	CustomerID  uint
	IsActive    bool
}

type PetID struct{ shared.BaseID }

func NewPetID(id uint) PetID {
	return PetID{BaseID: shared.BaseID{Value: id}}
}
func (p *Pet) Activate() error {
	if p.IsActive {
		return nil // Already active
	}
	p.IsActive = true
	p.IncrementVersion()
	return nil
}

func (p *Pet) Deactivate() error {
	if !p.IsActive {
		return nil // Already inactive
	}
	p.IsActive = false
	p.IncrementVersion()
	return nil
}

func (p *Pet) RequiresVaccination() bool {
	// Logic to determine if pet needs vaccination based on age and species
	if p.Age == nil {
		return false
	}

	// Puppies/kittens need more frequent vaccinations
	if *p.Age < 1 {
		return true
	}

	// Adult pets need annual vaccinations
	return true
}

func (p *Pet) LifeStage() string {
	if p.Age == nil {
		return "unknown"
	}

	age := *p.Age
	switch {
	case age < 1:
		return "baby"
	case age < 3:
		return "young"
	case age < 8:
		return "adult"
	default:
		return "senior"
	}
}

// ============================================================================
// Type Definitions
// ============================================================================

// PetSpecies represents the type/species of a pet
type PetSpecies string

// PetGender represents the gender of a pet
type PetGender string

// ============================================================================
// Constants
// ============================================================================

const (
	PetSpeciesDog      PetSpecies = "dog"
	PetSpeciesCat      PetSpecies = "cat"
	PetSpeciesBird     PetSpecies = "bird"
	PetSpeciesRabbit   PetSpecies = "rabbit"
	PetSpeciesHamster  PetSpecies = "hamster"
	PetSpeciesGuinea   PetSpecies = "guinea_pig"
	PetSpeciesFerret   PetSpecies = "ferret"
	PetSpeciesReptile  PetSpecies = "reptile"
	PetSpeciesFish     PetSpecies = "fish"
	PetSpeciesHorse    PetSpecies = "horse"
	PetSpeciesFarm     PetSpecies = "farm_animal"
	PetSpeciesExotic   PetSpecies = "exotic"
	PetSpeciesWildlife PetSpecies = "wildlife"
	PetSpeciesOther    PetSpecies = "other"
	PetSpeciesUnknown  PetSpecies = "unknown"
)

const (
	PetGenderMale         PetGender = "male"
	PetGenderFemale       PetGender = "female"
	PetGenderNeutered     PetGender = "neutered"
	PetGenderSpayed       PetGender = "spayed"
	PetGenderUnknown      PetGender = "unknown"
	PetGenderNotSpecified PetGender = "not_specified"
)

// ============================================================================
// Variables
// ============================================================================

var (
	ValidPetSpeciess = []PetSpecies{
		PetSpeciesDog,
		PetSpeciesCat,
		PetSpeciesBird,
		PetSpeciesRabbit,
		PetSpeciesHamster,
		PetSpeciesGuinea,
		PetSpeciesFerret,
		PetSpeciesReptile,
		PetSpeciesFish,
		PetSpeciesHorse,
		PetSpeciesFarm,
		PetSpeciesExotic,
		PetSpeciesWildlife,
		PetSpeciesOther,
		PetSpeciesUnknown,
	}

	PetSpeciesMap = map[string]PetSpecies{
		"dog":            PetSpeciesDog,
		"canine":         PetSpeciesDog,
		"puppy":          PetSpeciesDog,
		"k9":             PetSpeciesDog,
		"dogs":           PetSpeciesDog,
		"cat":            PetSpeciesCat,
		"feline":         PetSpeciesCat,
		"kitten":         PetSpeciesCat,
		"kitty":          PetSpeciesCat,
		"cats":           PetSpeciesCat,
		"bird":           PetSpeciesBird,
		"avian":          PetSpeciesBird,
		"parrot":         PetSpeciesBird,
		"canary":         PetSpeciesBird,
		"cockatiel":      PetSpeciesBird,
		"budgie":         PetSpeciesBird,
		"parakeet":       PetSpeciesBird,
		"cockatoo":       PetSpeciesBird,
		"macaw":          PetSpeciesBird,
		"finch":          PetSpeciesBird,
		"chicken":        PetSpeciesBird,
		"duck":           PetSpeciesBird,
		"goose":          PetSpeciesBird,
		"rabbit":         PetSpeciesRabbit,
		"bunny":          PetSpeciesRabbit,
		"hare":           PetSpeciesRabbit,
		"hamster":        PetSpeciesHamster,
		"syrian":         PetSpeciesHamster,
		"dwarf hamster":  PetSpeciesHamster,
		"guinea_pig":     PetSpeciesGuinea,
		"guinea pig":     PetSpeciesGuinea,
		"cavy":           PetSpeciesGuinea,
		"ferret":         PetSpeciesFerret,
		"reptile":        PetSpeciesReptile,
		"snake":          PetSpeciesReptile,
		"lizard":         PetSpeciesReptile,
		"gecko":          PetSpeciesReptile,
		"iguana":         PetSpeciesReptile,
		"turtle":         PetSpeciesReptile,
		"tortoise":       PetSpeciesReptile,
		"chameleon":      PetSpeciesReptile,
		"bearded dragon": PetSpeciesReptile,
		"python":         PetSpeciesReptile,
		"boa":            PetSpeciesReptile,
		"fish":           PetSpeciesFish,
		"goldfish":       PetSpeciesFish,
		"betta":          PetSpeciesFish,
		"tropical":       PetSpeciesFish,
		"aquatic":        PetSpeciesFish,
		"koi":            PetSpeciesFish,
		"angelfish":      PetSpeciesFish,
		"horse":          PetSpeciesHorse,
		"equine":         PetSpeciesHorse,
		"pony":           PetSpeciesHorse,
		"mare":           PetSpeciesHorse,
		"stallion":       PetSpeciesHorse,
		"foal":           PetSpeciesHorse,
		"colt":           PetSpeciesHorse,
		"filly":          PetSpeciesHorse,
		"farm_animal":    PetSpeciesFarm,
		"farm animal":    PetSpeciesFarm,
		"livestock":      PetSpeciesFarm,
		"cow":            PetSpeciesFarm,
		"cattle":         PetSpeciesFarm,
		"pig":            PetSpeciesFarm,
		"sheep":          PetSpeciesFarm,
		"goat":           PetSpeciesFarm,
		"llama":          PetSpeciesFarm,
		"alpaca":         PetSpeciesFarm,
		"exotic":         PetSpeciesExotic,
		"exotic pet":     PetSpeciesExotic,
		"unusual":        PetSpeciesExotic,
		"wildlife":       PetSpeciesWildlife,
		"wild":           PetSpeciesWildlife,
		"rescue":         PetSpeciesWildlife,
		"other":          PetSpeciesOther,
		"mixed":          PetSpeciesOther,
		"crossbreed":     PetSpeciesOther,
		"unknown":        PetSpeciesUnknown,
		"unspecified":    PetSpeciesUnknown,
		"":               PetSpeciesUnknown,
	}

	PetSpeciesDisplayNames = map[PetSpecies]string{
		PetSpeciesDog:      "Dog",
		PetSpeciesCat:      "Cat",
		PetSpeciesBird:     "Bird",
		PetSpeciesRabbit:   "Rabbit",
		PetSpeciesHamster:  "Hamster",
		PetSpeciesGuinea:   "Guinea Pig",
		PetSpeciesFerret:   "Ferret",
		PetSpeciesReptile:  "Reptile",
		PetSpeciesFish:     "Fish",
		PetSpeciesHorse:    "Horse",
		PetSpeciesFarm:     "Farm Animal",
		PetSpeciesExotic:   "Exotic Pet",
		PetSpeciesWildlife: "Wildlife",
		PetSpeciesOther:    "Other",
		PetSpeciesUnknown:  "Unknown",
	}
)

var (
	ValidPetGenders = []PetGender{
		PetGenderMale,
		PetGenderFemale,
		PetGenderNeutered,
		PetGenderSpayed,
		PetGenderUnknown,
		PetGenderNotSpecified,
	}

	petGenderMap = map[string]PetGender{
		"male":          PetGenderMale,
		"m":             PetGenderMale,
		"boy":           PetGenderMale,
		"female":        PetGenderFemale,
		"f":             PetGenderFemale,
		"girl":          PetGenderFemale,
		"neutered":      PetGenderNeutered,
		"neutered_male": PetGenderNeutered,
		"neutered male": PetGenderNeutered,
		"castrated":     PetGenderNeutered,
		"spayed":        PetGenderSpayed,
		"spayed_female": PetGenderSpayed,
		"spayed female": PetGenderSpayed,
		"unknown":       PetGenderUnknown,
		"unk":           PetGenderUnknown,
		"not_specified": PetGenderNotSpecified,
		"not specified": PetGenderNotSpecified,
		"unspecified":   PetGenderNotSpecified,
		"":              PetGenderNotSpecified,
	}

	petGenderDisplayNames = map[PetGender]string{
		PetGenderMale:         "Male",
		PetGenderFemale:       "Female",
		PetGenderNeutered:     "Neutered Male",
		PetGenderSpayed:       "Spayed Female",
		PetGenderUnknown:      "Unknown",
		PetGenderNotSpecified: "Not Specified",
	}
)

// ============================================================================
// PetSpecies Methods
// ============================================================================

func (pt PetSpecies) IsValid() bool {
	for _, validType := range ValidPetSpeciess {
		if pt == validType {
			return true
		}
	}
	return false
}

func ParsePetSpecies(PetSpecies string) (PetSpecies, error) {
	normalized := normalizePetSpeciesInput(PetSpecies)
	if val, exists := PetSpeciesMap[normalized]; exists {
		return val, nil
	}

	return "", InvalidEnumParserError("PetSpecies", PetSpecies)
}

func MustParsePetSpecies(PetSpecies string) PetSpecies {
	parsed, err := ParsePetSpecies(PetSpecies)
	if err != nil {
		panic(err)
	}
	return parsed
}

func (pt PetSpecies) String() string {
	return string(pt)
}

func (pt PetSpecies) DisplayName() string {
	if displayName, exists := PetSpeciesDisplayNames[pt]; exists {
		return displayName
	}
	return "Unknown Type"
}

func (pt PetSpecies) Values() []PetSpecies {
	return ValidPetSpeciess
}

// ============================================================================
// PetGender Methods
// ============================================================================

func (pg PetGender) IsValid() bool {
	_, exists := petGenderMap[string(pg)]
	return exists
}

func ParsePetGender(gender string) (PetGender, error) {
	normalized := normalizePetGenderInput(gender)
	if val, exists := petGenderMap[normalized]; exists {
		return val, nil
	}
	return "", InvalidEnumParserError("PetGender", gender)
}

func MustParsePetGender(gender string) PetGender {
	parsed, err := ParsePetGender(gender)
	if err != nil {
		panic(err)
	}
	return parsed
}

func (pg PetGender) String() string {
	return string(pg)
}

func (pg PetGender) DisplayName() string {
	if displayName, exists := petGenderDisplayNames[pg]; exists {
		return displayName
	}
	return "Unknown Gender"
}

func (pg PetGender) Values() []PetGender {
	return ValidPetGenders
}

// ============================================================================
// Helper Functions
// ============================================================================

func normalizePetSpeciesInput(input string) string {
	input = strings.TrimSpace(strings.ToLower(input))
	input = strings.ReplaceAll(input, " ", "_")
	input = strings.ReplaceAll(input, "-", "_")
	return input
}

func normalizePetGenderInput(input string) string {
	input = strings.TrimSpace(strings.ToLower(input))
	input = strings.ReplaceAll(input, " ", "_")
	return input
}

func GetAllPetSpeciess() []PetSpecies {
	return ValidPetSpeciess
}

func GetAllPetGenders() []PetGender {
	return ValidPetGenders
}

func GetExoticPetSpeciess() []PetSpecies {
	return []PetSpecies{
		PetSpeciesReptile,
		PetSpeciesFerret,
		PetSpeciesExotic,
		PetSpeciesWildlife,
	}
}

func GetPetSpeciessForEmergencyClinic() []PetSpecies {
	// Types commonly seen in emergency clinics
	return []PetSpecies{
		PetSpeciesDog,
		PetSpeciesCat,
		PetSpeciesBird,
		PetSpeciesRabbit,
		PetSpeciesFerret,
		PetSpeciesReptile,
	}
}

func SuggestPetSpeciesFromDescription(description string) PetSpecies {
	description = strings.ToLower(description)

	// Check for specific keywords
	typeKeywords := map[string]PetSpecies{
		"dog":      PetSpeciesDog,
		"puppy":    PetSpeciesDog,
		"canine":   PetSpeciesDog,
		"cat":      PetSpeciesCat,
		"kitten":   PetSpeciesCat,
		"feline":   PetSpeciesCat,
		"bird":     PetSpeciesBird,
		"parrot":   PetSpeciesBird,
		"rabbit":   PetSpeciesRabbit,
		"bunny":    PetSpeciesRabbit,
		"hamster":  PetSpeciesHamster,
		"guinea":   PetSpeciesGuinea,
		"ferret":   PetSpeciesFerret,
		"snake":    PetSpeciesReptile,
		"lizard":   PetSpeciesReptile,
		"turtle":   PetSpeciesReptile,
		"fish":     PetSpeciesFish,
		"goldfish": PetSpeciesFish,
		"horse":    PetSpeciesHorse,
		"pony":     PetSpeciesHorse,
	}

	for keyword, PetSpecies := range typeKeywords {
		if strings.Contains(description, keyword) {
			return PetSpecies
		}
	}

	return PetSpeciesUnknown
}
