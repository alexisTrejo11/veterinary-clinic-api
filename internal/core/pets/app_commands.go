package pets

type CreatePetCommand struct {
	Name                  string
	Photo                 *string
	Species               PetSpecies
	Breed                 *string
	Age                   *int
	Gender                PetGender
	Color                 *string
	MicrochipID           *string
	BloodType             *string
	IsNeutered            *bool
	CustomerID            uint
	IsActive              bool
	Allergies             *string
	CurrentMedications    *string
	SpecialNeeds          *string
	FeedingInstructions   *string
	BehavioralNotes       *string
	VeterinaryContact     *string
	EmergencyContactName  *string
	EmergencyContactPhone *string
}

type UpdatePetCommand struct {
	PetID       PetID
	Name        *string
	Photo       *string
	Species     *PetSpecies
	Breed       *string
	Age         *int
	Gender      *PetGender
	Color       *string
	BloodType   *string
	MicrochipID *string
	IsNeutered  *bool
	CustomerID  *uint
	IsActive    *bool
}

type DeletePetCommand struct {
	PetID         PetID
	OptCustomerID *uint
	IsHardDelete  bool
}
