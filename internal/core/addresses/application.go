package addresses

import (
	"clinic-vet-api/internal/shared/page"
	"context"
)

var (
	MAX_ADDRESS_PER_CUSTOMER = 3
)

type CreateAddressCommand struct {
	CustomerID          uint
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
	DeleteAddress(ctx context.Context, id AddressID) error
	GetAddressByIDAndCustomerID(ctx context.Context, id AddressID, customerID uint) (Address, error)
	GetAddressesByCustomerID(ctx context.Context, customerID uint) ([]Address, error)
	GetAddressesBySpecification(ctx context.Context, spec AddressSpecification) (page.Page[Address], error)
	GetAddressesByIDAndCustomerID(ctx context.Context, id AddressID, customerID uint) (Address, error)
}

type addressService struct {
	repository AddressRepository
}

func NewAddressService(repository AddressRepository) AddressService {
	return &addressService{repository: repository}
}

func (s *addressService) CreateAddress(ctx context.Context, cmd CreateAddressCommand) (Address, error) {
	const op = "CreateAddress"
	count, err := s.repository.CountByCustomerID(ctx, cmd.CustomerID)
	if err != nil {
		return Address{}, err
	}
	if count >= int64(MAX_ADDRESS_PER_CUSTOMER) {
		return Address{}, MaxAddressesPerCustomerError(ctx, cmd.CustomerID, MAX_ADDRESS_PER_CUSTOMER, op)
	}

	address := Address{
		CustomerID:          cmd.CustomerID,
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
		addresses, err := s.repository.GetAllByCustomerID(ctx, created.CustomerID)
		if err != nil {
			return Address{}, err
		}
		for i := range addresses {
			if addresses[i].ID.Value == created.ID.Value {
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
		addresses, err := s.repository.GetAllByCustomerID(ctx, updated.CustomerID)
		if err != nil {
			return err
		}
		for i := range addresses {
			if addresses[i].ID.Value == updated.ID.Value {
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
		addresses, err := s.repository.GetAllByCustomerID(ctx, updated.CustomerID)
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
			return DefaultAddressRequiredError(ctx, updated.CustomerID, op)
		}
	}

	return nil
}

func (s *addressService) RestoreAddress(ctx context.Context, id AddressID) error {
	return s.repository.RestoreByID(ctx, id)
}

func (s *addressService) DeleteAddress(ctx context.Context, id AddressID) error {
	err := s.repository.Delete(ctx, id)
	return err
}

func (s *addressService) GetAddressByIDAndCustomerID(ctx context.Context, id AddressID, customerID uint) (Address, error) {
	return s.repository.GetByIDAndCustomerID(ctx, id, customerID)
}

func (s *addressService) GetAddressesByCustomerID(ctx context.Context, customerID uint) ([]Address, error) {
	return s.repository.GetAllByCustomerID(ctx, customerID)
}

func (s *addressService) GetAddressesBySpecification(ctx context.Context, spec AddressSpecification) (page.Page[Address], error) {
	return s.repository.GetBySpecification(ctx, spec)
}

func (s *addressService) GetAddressesByIDAndCustomerID(ctx context.Context, id AddressID, customerID uint) (Address, error) {
	return s.repository.GetByIDAndCustomerID(ctx, id, customerID)
}
