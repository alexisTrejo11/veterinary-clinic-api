package addresses

import (
	"clinic-vet-api/internal/shared"
	"clinic-vet-api/internal/shared/page"
	"context"
	"errors"
)

var (
	MAX_ADDRESS_PER_CUSTOMER = 3
)

type CreateAddressCommand struct {
	UserID              shared.UserID
	Street              string
	City                string
	State               string
	ZipCode             string
	Country             Country
	BuildingType        BuildingType
	BuildingOuterNumber string
	BuildingInnerNumber *string
	IsDefault           bool
}

type UpdateAddressCommand struct {
	ID                  AddressID
	UserID              *shared.UserID
	Street              *string
	City                *string
	State               *string
	ZipCode             *string
	Country             *Country
	BuildingType        *BuildingType
	BuildingOuterNumber *string
	BuildingInnerNumber *string
	IsDefault           *bool
}

type AddressService interface {
	CreateAddress(ctx context.Context, cmd CreateAddressCommand) (Address, error)
	UpdateAddress(ctx context.Context, cmd UpdateAddressCommand) error
	RestoreAddress(ctx context.Context, id AddressID) error
	DeleteAddress(ctx context.Context, id AddressID, userID *shared.UserID) error
	GetAddressByIDAndUserID(ctx context.Context, id AddressID, userID *shared.UserID) (Address, error)
	GetAddressesByUserID(ctx context.Context, userID shared.UserID) ([]Address, error)
	GetAddressesBySpecification(ctx context.Context, spec AddressSpecification) (page.Page[Address], error)
	GetAddressesByIDAndUserID(ctx context.Context, id AddressID, userID shared.UserID) (Address, error)
}

type addressService struct {
	repository AddressRepository
}

func NewAddressService(repository AddressRepository) AddressService {
	return &addressService{repository: repository}
}

func (s *addressService) CreateAddress(ctx context.Context, cmd CreateAddressCommand) (Address, error) {
	const op = "CreateAddress"
	count, err := s.repository.CountByUserID(ctx, cmd.UserID)
	if err != nil {
		return Address{}, err
	}
	if count >= int64(MAX_ADDRESS_PER_CUSTOMER) {
		return Address{}, MaxAddressesPerCustomerError(ctx, cmd.UserID, MAX_ADDRESS_PER_CUSTOMER, op)
	}

	address := Address{
		UserID:              cmd.UserID,
		Street:              cmd.Street,
		City:                cmd.City,
		State:               cmd.State,
		ZipCode:             cmd.ZipCode,
		Country:             cmd.Country,
		BuildingType:        cmd.BuildingType,
		BuildingOuterNumber: cmd.BuildingOuterNumber,
		BuildingInnerNumber: cmd.BuildingInnerNumber,
		IsDefault:           cmd.IsDefault,
	}

	if err := address.Validate(ctx); err != nil {
		return Address{}, err
	}

	created, err := s.repository.Save(ctx, address)
	if err != nil {
		return Address{}, err
	}

	// If this address should be default, clear default from others.
	if created.IsDefault {
		addresses, err := s.repository.GetAllByUserID(ctx, created.UserID)
		if err != nil {
			return Address{}, err
		}
		for i := range addresses {
			if addresses[i].ID.Value() == created.ID.Value() {
				addresses[i].IsDefault = true
			} else if addresses[i].IsDefault {
				addresses[i].IsDefault = false
			}
		}
		if err := s.repository.BulkUpdate(ctx, addresses); err != nil {
			return Address{}, err
		}
	}

	return created, nil
}

func (s *addressService) UpdateAddress(ctx context.Context, cmd UpdateAddressCommand) error {
	const op = "UpdateAddress"
	address, err := s.repository.GetByID(ctx, cmd.ID)
	if err != nil {
		return err
	}

	originalWasDefault := address.IsDefault
	if cmd.Street != nil {
		address.Street = *cmd.Street
	}
	if cmd.City != nil {
		address.City = *cmd.City
	}
	if cmd.State != nil {
		address.State = *cmd.State
	}
	if cmd.ZipCode != nil {
		address.ZipCode = *cmd.ZipCode
	}
	if cmd.Country != nil {
		address.Country = *cmd.Country
	}
	if cmd.BuildingType != nil {
		address.BuildingType = *cmd.BuildingType
	}
	if cmd.BuildingOuterNumber != nil {
		address.BuildingOuterNumber = *cmd.BuildingOuterNumber
	}
	if cmd.BuildingInnerNumber != nil {
		address.BuildingInnerNumber = cmd.BuildingInnerNumber
	}

	// Manage default flag change
	if cmd.IsDefault != nil {
		address.IsDefault = *cmd.IsDefault
	}

	if err := address.Validate(ctx); err != nil {
		return err
	}

	updated, err := s.repository.Save(ctx, address)
	if err != nil {
		return err
	}

	// If this address became default, clear default from others for this customer.
	if cmd.IsDefault != nil && *cmd.IsDefault {
		addresses, err := s.repository.GetAllByUserID(ctx, updated.UserID)
		if err != nil {
			return err
		}
		for i := range addresses {
			if addresses[i].ID.Value() == updated.ID.Value() {
				addresses[i].IsDefault = true
			} else if addresses[i].IsDefault {
				addresses[i].IsDefault = false
			}
		}
		if err := s.repository.BulkUpdate(ctx, addresses); err != nil {
			return err
		}
	}

	// Optional: prevent removing the last default address.
	if cmd.IsDefault != nil && !*cmd.IsDefault && originalWasDefault {
		addresses, err := s.repository.GetAllByUserID(ctx, updated.UserID)
		if err != nil {
			return err
		}
		hasDefault := false
		for _, a := range addresses {
			if a.IsDefault {
				hasDefault = true
				break
			}
		}
		if !hasDefault && len(addresses) > 0 {
			return DefaultAddressRequiredError(ctx, updated.UserID, op)
		}
	}

	return nil
}

func (s *addressService) RestoreAddress(ctx context.Context, id AddressID) error {
	return s.repository.RestoreByID(ctx, id)
}

func AddressNotFoundError(ctx context.Context, id AddressID) error {
	return errors.New("address not found")
}

func (s *addressService) DeleteAddress(ctx context.Context, id AddressID, optUserID *shared.UserID) error {
	if optUserID == nil {
		if address, err := s.repository.GetByID(ctx, id); err != nil {
			return err
		} else {
			return s.repository.Delete(ctx, address.ID)
		}
	}

	address, err := s.repository.GetByIDAndUserID(ctx, id, *optUserID)
	if err != nil {
		return err
	}
	return s.repository.Delete(ctx, address.ID)
}

func (s *addressService) GetAddressByIDAndUserID(ctx context.Context, id AddressID, userID *shared.UserID) (Address, error) {
	if userID == nil {
		return s.repository.GetByID(ctx, id)
	}
	return s.repository.GetByIDAndUserID(ctx, id, *userID)
}

func (s *addressService) GetAddressesByUserID(ctx context.Context, userID shared.UserID) ([]Address, error) {
	return s.repository.GetAllByUserID(ctx, userID)
}

func (s *addressService) GetAddressesBySpecification(ctx context.Context, spec AddressSpecification) (page.Page[Address], error) {
	return s.repository.GetBySpecification(ctx, spec)
}

func (s *addressService) GetAddressesByIDAndUserID(ctx context.Context, id AddressID, userID shared.UserID) (Address, error) {
	return s.repository.GetByIDAndUserID(ctx, id, userID)
}
