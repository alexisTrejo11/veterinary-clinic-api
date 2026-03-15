package customers

import (
	"clinic-vet-api/internal/core/pets"
	pet "clinic-vet-api/internal/core/pets"
	"clinic-vet-api/internal/shared"
	"clinic-vet-api/internal/shared/errors"
	"context"
)

type Customer struct {
	shared.Entity[CustomerID]
	shared.Person
	PhotoURL string
	UserID   uint
	IsActive bool
	Pets     []pets.Pet
}

// CustomerID is a unique identifier for a customer user
type CustomerID struct{ shared.IntegerID }

func NewCustomerID(value uint) CustomerID {
	return CustomerID{shared.NewBaseID(value)}
}

func (o *Customer) UpdatePhoto(ctx context.Context, newPhoto string) error {
	const operation = "update_photo"

	if newPhoto != "" && len(newPhoto) > 500 {
		return PhotoURLLongError(ctx, len(newPhoto), operation)
	}

	o.PhotoURL = newPhoto
	o.IncrementVersion()
	return nil
}

func (o *Customer) UpdateFullName(ctx context.Context, firstName, lastName string) error {
	operation := "UPDATE_NAME"
	o.FirstName = firstName
	o.LastName = lastName

	if err := o.Person.Validate(ctx, operation); err != nil {
		return err
	}

	o.IncrementVersion()
	return nil
}

func (o *Customer) UpdateGender(ctx context.Context, newGender shared.PersonGender) error {
	const operation = "update_gender"
	o.Gender = newGender

	if err := o.Person.Validate(ctx, operation); err != nil {
		return err
	}

	o.IncrementVersion()
	return nil
}

func (o *Customer) AssociateWithUser(ctx context.Context, userID uint) error {
	const operation = "associate_user"

	if o.UserID != 0 && o.UserID == userID {
		return UserAlreadyAssociatedError(ctx, userID, operation)
	}

	o.UserID = userID
	o.IncrementVersion()
	return nil
}

func (o *Customer) RemoveUserAssociation(ctx context.Context) error {
	const operation = "remove_user_association"

	if o.UserID == 0 {
		return UserNotAssociatedError(ctx, operation)
	}

	o.UserID = 0
	o.IncrementVersion()
	return nil
}

func (o *Customer) Activate(ctx context.Context) error {
	if o.IsActive {
		return nil // Already active
	}

	o.IsActive = true
	o.IncrementVersion()
	return nil
}

func (o *Customer) Deactivate(ctx context.Context) error {
	const operation = "deactivate_customer"

	if !o.IsActive {
		return nil // Already inactive
	}

	if !o.CanBeDeactivated() {
		activePets := 0
		for _, pet := range o.Pets {
			if pet.IsActive {
				activePets++
			}
		}
		return CannotDeactivateWithPetsError(ctx, activePets, operation)
	}

	o.IsActive = false
	o.IncrementVersion()
	return nil
}

func (o *Customer) AddPet(ctx context.Context, newPet *pet.Pet) error {
	const operation = "add_pet"
	if newPet == nil {
		return errors.MissingEntity(ctx, "pet", "Pet cannot be nil", operation)
	}

	for _, existingPet := range o.Pets {
		if existingPet.ID.Value() == newPet.ID.Value() {
			return PetAlreadyExistsError(ctx, newPet.ID, operation)
		}
	}

	o.Pets = append(o.Pets, *newPet)
	o.IncrementVersion()
	return nil
}

func (o *Customer) RemovePet(ctx context.Context, petID pets.PetID) error {
	const operation = "remove_pet"

	for i, existingPet := range o.Pets {
		if existingPet.ID.Value() == petID.Value() {
			// Remove the pet
			o.Pets = append(o.Pets[:i], o.Pets[i+1:]...)
			o.IncrementVersion()
			return nil
		}
	}

	return PetNotFoundError(ctx, petID, operation)
}

func (o *Customer) HasActivePets() bool {
	for _, pet := range o.Pets {
		if pet.IsActive {
			return true
		}
	}
	return false
}

func (o *Customer) CanBeDeactivated() bool {
	return !o.HasActivePets()
}
