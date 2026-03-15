package mappers

import (
	"clinic-vet-api/internal/core/pets"
	"clinic-vet-api/internal/infrastructure/http/handlers/dtos"
	"clinic-vet-api/internal/shared/page"
)

type PetMapper struct{}

func NewPetMapper() *PetMapper {
	return &PetMapper{}
}

func (m *PetMapper) RequestToCreateCommand(r dtos.PetCreateRequest, customerID uint, isActive bool) (pets.CreatePetCommand, error) {
	species, err := pets.ParsePetSpecies(r.Species)
	if err != nil {
		return pets.CreatePetCommand{}, err
	}

	gender, err := pets.ParsePetGender(r.Gender)
	if err != nil {
		return pets.CreatePetCommand{}, err
	}

	return pets.CreatePetCommand{
		Name:                r.Name,
		Photo:               r.Photo,
		Species:             species,
		Breed:               r.Breed,
		Age:                 r.Age,
		Gender:              gender,
		Color:               r.Color,
		MicrochipID:         r.MicrochipID,
		BloodType:           r.BloodType,
		IsNeutered:          r.IsNeutered,
		CustomerID:          customerID,
		IsActive:            isActive,
		Allergies:           r.Allergies,
		CurrentMedications:  r.CurrentMedications,
		SpecialNeeds:        r.SpecialNeeds,
		FeedingInstructions: r.FeedingInstructions,
		BehavioralNotes:     r.BehavioralNotes,
		VeterinaryContact:   r.VeterinaryContact,
		EmergencyContactName: func() *string {
			return r.EmergencyName
		}(),
		EmergencyContactPhone: func() *string {
			return r.EmergencyPhone
		}(),
	}, nil
}

func (m *PetMapper) RequestToUpdateCommand(r dtos.PetUpdateRequest, petID pets.PetID, optCustomerID *uint) (pets.UpdatePetCommand, error) {
	cmd := pets.UpdatePetCommand{PetID: petID, CustomerID: optCustomerID}
	if r.Name != nil {
		cmd.Name = r.Name
	}
	if r.Photo != nil {
		cmd.Photo = r.Photo
	}
	if r.Species != nil {
		species, err := pets.ParsePetSpecies(*r.Species)
		if err != nil {
			return pets.UpdatePetCommand{}, err
		}
		cmd.Species = &species
	}
	if r.Breed != nil {
		cmd.Breed = r.Breed
	}
	if r.Age != nil {
		cmd.Age = r.Age
	}
	if r.Gender != nil {
		gender, err := pets.ParsePetGender(*r.Gender)
		if err != nil {
			return pets.UpdatePetCommand{}, err
		}
		cmd.Gender = &gender
	}
	if r.Color != nil {
		cmd.Color = r.Color
	}
	if r.MicrochipID != nil {
		cmd.MicrochipID = r.MicrochipID
	}
	if r.BloodType != nil {
		cmd.BloodType = r.BloodType
	}
	if r.IsNeutered != nil {
		cmd.IsNeutered = r.IsNeutered
	}
	if r.IsActive != nil {
		cmd.IsActive = r.IsActive
	}
	return cmd, nil
}

func (m *PetMapper) ToResponse(p pets.Pet) dtos.PetResponse {
	return dtos.PetResponse{
		ID:                  p.ID.Value(),
		Name:                p.Name,
		Species:             p.Species.DisplayName(),
		Gender:              p.Gender.DisplayName(),
		CustomerID:          p.CustomerID,
		IsActive:            p.IsActive,
		Breed:               p.Breed,
		Age:                 p.Age,
		Photo:               p.Photo,
		Color:               p.Color,
		MicrochipID:         p.MicrochipID,
		BloodType:           p.BloodType,
		IsNeutered:          p.IsNeutered,
		Allergies:           p.Allergies,
		CurrentMedications:  p.CurrentMedications,
		SpecialNeeds:        p.SpecialNeeds,
		FeedingInstructions: p.FeedingInstructions,
	}
}

func (m *PetMapper) ToPaginatedResponse(p page.Page[pets.Pet]) page.Page[dtos.PetResponse] {
	return page.MapItems(p, m.ToResponse)
}

func (m *PetMapper) RequestToSpecification(r dtos.PetSearchRequest) (pets.PetSpecification, error) {
	return pets.PetSpecification{
		Pagination: r.PaginationRequest.ToPagination(),
	}, nil
}
