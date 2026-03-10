package customers

import (
	"clinic-vet-api/internal/shared/page"
	"context"
)

type CustomerRepository interface {
	//FindBySpecification(ctx context.Context, spec specification.CustomerSpecification) (page.Page[Customer], error)
	FindByID(ctx context.Context, id CustomerID) (Customer, error)
	FindActive(ctx context.Context, pagination page.Pagination) (page.Page[Customer], error)

	ExistsByID(ctx context.Context, id CustomerID) (bool, error)

	Save(ctx context.Context, customer *Customer) error
	Update(ctx context.Context, customer *Customer) error
	SoftDelete(ctx context.Context, id CustomerID) error
	HardDelete(ctx context.Context, id CustomerID) error

	CountAll(ctx context.Context) (int64, error)
	CountActive(ctx context.Context) (int64, error)
}
