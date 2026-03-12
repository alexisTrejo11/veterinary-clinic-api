package customers

import (
	"clinic-vet-api/internal/shared/page"
	"context"
)

type CustomerRepository interface {
	FindBySpecification(ctx context.Context, spec CustomerSpecification) (page.Page[Customer], error)
	FindByID(ctx context.Context, id CustomerID) (Customer, error)
	FindActive(ctx context.Context, pagination page.Pagination) (page.Page[Customer], error)

	ExistsByID(ctx context.Context, id CustomerID) (bool, error)

	Save(ctx context.Context, customer Customer) (Customer, error)
	SoftDelete(ctx context.Context, id CustomerID) error
	HardDelete(ctx context.Context, id CustomerID) error
	IsDeletedByID(ctx context.Context, id CustomerID) (bool, error)

	ExistsByUserID(ctx context.Context, userID uint) (bool, error)

	CountAll(ctx context.Context) (int64, error)
	CountActive(ctx context.Context) (int64, error)

	// Optional: aggregate stats extension point (pets, appointments, etc.)
}
