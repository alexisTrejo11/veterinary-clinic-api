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

func (m *PetMapper) ToResponse(p pets.Pet) dtos.PetResponse {
	return dtos.PetResponse{
		ID:                  p.ID.Value,
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
