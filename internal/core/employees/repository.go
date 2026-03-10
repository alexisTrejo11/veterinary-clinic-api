package employees

import (
	p "clinic-vet-api/internal/shared/page"
	"context"
)

type EmployeeRepository interface {
	//FindBySpecification(ctx context.Context, spec specification.EmployeeSearchSpecification) (p.Page[Employee], error)
	FindByID(ctx context.Context, id EmployeeID) (Employee, error)
	FindByUserID(ctx context.Context, userID uint) (Employee, error)
	FindActive(ctx context.Context, pagination p.Pagination) (p.Page[Employee], error)
	FindAll(ctx context.Context, pagination p.Pagination) (p.Page[Employee], error)
	FindBySpeciality(ctx context.Context, speciality VetSpecialty, pagination p.Pagination) (p.Page[Employee], error)

	ExistsByID(ctx context.Context, id EmployeeID) (bool, error)
	ExistsByUserID(ctx context.Context, userID uint) (bool, error)

	Save(ctx context.Context, employee *Employee) error
	Update(ctx context.Context, employee *Employee) error
	SoftDelete(ctx context.Context, id EmployeeID) error
	HardDelete(ctx context.Context, id EmployeeID) error

	CountBySpeciality(ctx context.Context, speciality VetSpecialty) (int64, error)
	CountActive(ctx context.Context) (int64, error)
}
