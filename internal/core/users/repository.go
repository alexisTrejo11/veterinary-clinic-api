package users

import (
	"context"
	"time"

	cust "clinic-vet-api/internal/core/customers"
	emp "clinic-vet-api/internal/core/employees"
	"clinic-vet-api/internal/shared/page"
)

type UserRepository interface {
	FindSpecification(ctx context.Context, spec UserSpecification) (page.Page[User], error)
	FindByID(ctx context.Context, id UserID) (User, error)
	FindByEmail(ctx context.Context, email string) (User, error)
	FindByPhone(ctx context.Context, phone string) (User, error)

	FindAll(ctx context.Context, pageInput page.Pagination) (page.Page[User], error)
	FindByRole(ctx context.Context, role string, pageInput page.Pagination) (page.Page[User], error)
	FindActive(ctx context.Context, pageInput page.Pagination) (page.Page[User], error)
	FindInactive(ctx context.Context, pageInput page.Pagination) (page.Page[User], error)

	FindByEmployeeID(ctx context.Context, employeeID emp.EmployeeID) (User, error)
	FindByCustomerID(ctx context.Context, customerID cust.CustomerID) (User, error)
	FindRecentlyLoggedIn(ctx context.Context, since time.Time, pageInput page.Pagination) (page.Page[User], error)
	FindByOAuthProvider(ctx context.Context, provider string, providerID string) (User, error)

	ExistsByID(ctx context.Context, id UserID) (bool, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	ExistsByPhone(ctx context.Context, phone string) (bool, error)
	ExistsByEmployeeID(ctx context.Context, employeeID emp.EmployeeID) (bool, error)
	ExistsByCustomerID(ctx context.Context, customerID cust.CustomerID) (bool, error)

	Save(ctx context.Context, user *User) error
	SoftDelete(ctx context.Context, id UserID) error
	HardDelete(ctx context.Context, id UserID) error

	UpdateLastLogin(ctx context.Context, id UserID) error
	UpdatePassword(ctx context.Context, id UserID, hashedPassword string) error
	UpdateStatus(ctx context.Context, id UserID, status UserStatus) error

	CountAll(ctx context.Context) (int64, error)
	CountByRole(ctx context.Context, role string) (int64, error)
	CountByStatus(ctx context.Context, status UserStatus) (int64, error)
	CountActive(ctx context.Context) (int64, error)
}
