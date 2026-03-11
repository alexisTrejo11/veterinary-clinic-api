package dtos

import (
	"clinic-vet-api/internal/shared/page"
	"time"
)

type PetCreateRequest struct {
	Name                string  `json:"name" binding:"required"`
	Photo               *string `json:"photo" binding:"omitempty,url"`
	Species             string  `json:"species" binding:"required"`
	Breed               *string `json:"breed" binding:"omitempty"`
	Age                 *int    `json:"age" binding:"omitempty"`
	Gender              string  `json:"gender" binding:"required"`
	Color               *string `json:"color" binding:"omitempty"`
	MicrochipID         *string `json:"microchip_id" binding:"omitempty"`
	BloodType           *string `json:"blood_type" binding:"omitempty"`
	IsNeutered          *bool   `json:"is_neutered" binding:"omitempty"`
	CustomerID          uint    `json:"customer_id" binding:"required"`
	Allergies           *string `json:"allergies" binding:"omitempty"`
	CurrentMedications  *string `json:"current_medications" binding:"omitempty"`
	SpecialNeeds        *string `json:"special_needs" binding:"omitempty"`
	FeedingInstructions *string `json:"feeding_instructions" binding:"omitempty"`
	BehavioralNotes     *string `json:"behavioral_notes" binding:"omitempty"`
	VeterinaryContact   *string `json:"veterinary_contact" binding:"omitempty"`
	EmergencyName       *string `json:"emergency_contact_name" binding:"omitempty"`
	EmergencyPhone      *string `json:"emergency_contact_phone" binding:"omitempty"`
}

type PetSearchRequest struct {
	page.PaginationRequest
}

// PetResponse represents the HTTP response DTO for a pet.
// It mirrors the fields on the Pet aggregate plus basic metadata.
type PetResponse struct {
	ID                    uint      `json:"id"`
	Name                  string    `json:"name"`
	Species               string    `json:"species"`
	Gender                string    `json:"gender"`
	CustomerID            uint      `json:"customer_id"`
	IsActive              bool      `json:"is_active"`
	Breed                 *string   `json:"breed,omitempty"`
	Age                   *int      `json:"age,omitempty"`
	Photo                 *string   `json:"photo,omitempty"`
	Color                 *string   `json:"color,omitempty"`
	MicrochipID           *string   `json:"microchip_id,omitempty"`
	BloodType             *string   `json:"blood_type,omitempty"`
	IsNeutered            *bool     `json:"is_neutered,omitempty"`
	Allergies             *string   `json:"allergies,omitempty"`
	CurrentMedications    *string   `json:"current_medications,omitempty"`
	SpecialNeeds          *string   `json:"special_needs,omitempty"`
	FeedingInstructions   *string   `json:"feeding_instructions,omitempty"`
	BehavioralNotes       *string   `json:"behavioral_notes,omitempty"`
	VeterinaryContact     *string   `json:"veterinary_contact,omitempty"`
	EmergencyContactName  *string   `json:"emergency_contact_name,omitempty"`
	EmergencyContactPhone *string   `json:"emergency_contact_phone,omitempty"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}
