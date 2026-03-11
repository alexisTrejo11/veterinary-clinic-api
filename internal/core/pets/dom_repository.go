package pets

import (
	"context"

	"clinic-vet-api/internal/shared/page"
)

type PetRepository interface {
	FindBySpecification(ctx context.Context, spec *PetSpecification) (page.Page[Pet], error)
	FindAllByCustomerID(ctx context.Context, customerID uint) ([]Pet, error)
	FindByCustomerID(ctx context.Context, customerID uint, pagination page.Pagination) (page.Page[Pet], error)
	FindByID(ctx context.Context, petID PetID) (Pet, error)
	FindByIDAndCustomerID(ctx context.Context, id PetID, customerID uint) (Pet, error)

	ExistsByID(ctx context.Context, petID PetID) (bool, error)
	ExistsByMicrochip(ctx context.Context, microchip string) (bool, error)
	ExistsDeletedByID(ctx context.Context, petID PetID) (bool, error)
	Save(ctx context.Context, pet Pet) (Pet, error)
	Delete(ctx context.Context, petID PetID, isHard bool) error
	RestoreByID(ctx context.Context, petID PetID) error

	CountByCustomerID(ctx context.Context, customerID uint) (int64, error)
}

type CustomerRepository interface {
	ExistsByIDAndActive(ctx context.Context, customerID uint) (bool, error)
}
