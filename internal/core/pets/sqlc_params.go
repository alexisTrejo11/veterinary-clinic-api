package pets

import (
	"clinic-vet-api/sqlc"

	"github.com/jackc/pgx/v5/pgtype"
)

func textFromPtr(s *string) pgtype.Text {
	if s == nil {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: *s, Valid: true}
}

func int2FromPtr(i *int) pgtype.Int2 {
	if i == nil {
		return pgtype.Int2{Valid: false}
	}
	return pgtype.Int2{Int16: int16(*i), Valid: true}
}

func boolFromPtr(b *bool) pgtype.Bool {
	if b == nil {
		return pgtype.Bool{Valid: false}
	}
	return pgtype.Bool{Bool: *b, Valid: true}
}

// ToCreateParams maps a Pet domain entity to sqlc CreatePetParams.
func (p Pet) ToCreateParams() sqlc.CreatePetParams {
	return sqlc.CreatePetParams{
		Name:                  p.Name,
		Photo:                 textFromPtr(p.Photo),
		Species:               p.Species.String(),
		Breed:                 textFromPtr(p.Breed),
		Age:                   int2FromPtr(p.Age),
		Gender:                pgtype.Text{String: p.Gender.String(), Valid: p.Gender.String() != ""},
		Color:                 textFromPtr(p.Color),
		Microchip:             textFromPtr(p.MicrochipID),
		BloodType:             textFromPtr(p.BloodType),
		IsNeutered:            boolFromPtr(p.IsNeutered),
		CustomerID:            int32(p.CustomerID),
		IsActive:              p.IsActive,
		Allergies:             textFromPtr(p.Allergies),
		CurrentMedications:    textFromPtr(p.CurrentMedications),
		SpecialNeeds:          textFromPtr(p.SpecialNeeds),
		FeedingInstructions:   textFromPtr(p.FeedingInstructions),
		BehavioralNotes:       textFromPtr(p.BehavioralNotes),
		VeterinaryContact:     textFromPtr(p.VeterinaryContact),
		EmergencyContactName:  textFromPtr(p.EmergencyContactName),
		EmergencyContactPhone: textFromPtr(p.EmergencyContactPhone),
	}
}

// ToUpdateParams maps a Pet domain entity to sqlc UpdatePetParams.
func (p Pet) ToUpdateParams() sqlc.UpdatePetParams {
	return sqlc.UpdatePetParams{
		ID:                    p.ID.Int32(),
		Name:                  p.Name,
		Photo:                 textFromPtr(p.Photo),
		Species:               p.Species.String(),
		Breed:                 textFromPtr(p.Breed),
		Age:                   int2FromPtr(p.Age),
		Gender:                pgtype.Text{String: p.Gender.String(), Valid: p.Gender.String() != ""},
		Color:                 textFromPtr(p.Color),
		Microchip:             textFromPtr(p.MicrochipID),
		IsNeutered:            boolFromPtr(p.IsNeutered),
		CustomerID:            int32(p.CustomerID),
		BloodType:             textFromPtr(p.BloodType),
		IsActive:              p.IsActive,
		Allergies:             textFromPtr(p.Allergies),
		CurrentMedications:    textFromPtr(p.CurrentMedications),
		SpecialNeeds:          textFromPtr(p.SpecialNeeds),
		FeedingInstructions:   textFromPtr(p.FeedingInstructions),
		BehavioralNotes:       textFromPtr(p.BehavioralNotes),
		VeterinaryContact:     textFromPtr(p.VeterinaryContact),
		EmergencyContactName:  textFromPtr(p.EmergencyContactName),
		EmergencyContactPhone: textFromPtr(p.EmergencyContactPhone),
	}
}

