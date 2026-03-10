package users

import (
	"context"

	cust "clinic-vet-api/internal/core/customers"
	emp "clinic-vet-api/internal/core/employees"
	"clinic-vet-api/internal/shared"
	"clinic-vet-api/internal/shared/page"
)

type UserRepository interface {
	FindByID(ctx context.Context, id shared.UserID) (User, error)
	FindByEmail(ctx context.Context, email Email) (User, error)
	FindByPhone(ctx context.Context, phone PhoneNumber) (User, error)
	FindByEmployeeID(ctx context.Context, employeeID emp.EmployeeID) (User, error)
	FindByCustomerID(ctx context.Context, customerID cust.CustomerID) (User, error)
	FindByOAuthProvider(ctx context.Context, provider string, providerID string) (User, error)
	FindSpecification(ctx context.Context, spec UserSpecification) (page.Page[User], error)

	ExistsByID(ctx context.Context, id shared.UserID) (bool, error)
	ExistsByEmail(ctx context.Context, email Email) (bool, error)
	ExistsByPhone(ctx context.Context, phone PhoneNumber) (bool, error)
	ExistsByEmployeeID(ctx context.Context, employeeID emp.EmployeeID) (bool, error)
	ExistsByCustomerID(ctx context.Context, customerID cust.CustomerID) (bool, error)

	IsDeletedByID(ctx context.Context, id shared.UserID) (bool, error)
	Save(ctx context.Context, user *User) error
	SoftDelete(ctx context.Context, id shared.UserID) error
	HardDelete(ctx context.Context, id shared.UserID) error

	RestoreByID(ctx context.Context, id shared.UserID) error
	Count(ctx context.Context, spec UserSpecification) (int64, error)
}
