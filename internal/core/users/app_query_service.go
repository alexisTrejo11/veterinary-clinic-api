package users

import (
	cust "clinic-vet-api/internal/core/customers"
	emp "clinic-vet-api/internal/core/employees"
	"clinic-vet-api/internal/shared"
	"clinic-vet-api/internal/shared/page"
	"context"
)

type QueryService interface {
	GetByID(ctx context.Context, id shared.UserID) (User, error)
	GetByEmail(ctx context.Context, email Email) (User, error)
	GetByPhone(ctx context.Context, phone PhoneNumber) (User, error)
	GetByEmployeeID(ctx context.Context, employeeID emp.EmployeeID) (User, error)
	GetByCustomerID(ctx context.Context, customerID cust.CustomerID) (User, error)
	GetByOAuthProvider(ctx context.Context, provider string, providerID string) (User, error)
	GetSpecification(ctx context.Context, spec UserSpecification) (page.Page[User], error)
}

type userQueryService struct {
	repository UserRepository
}

func NewUserQueryService(repository UserRepository) QueryService {
	return &userQueryService{repository: repository}
}

func (s *userQueryService) GetByID(ctx context.Context, id shared.UserID) (User, error) {
	return s.repository.FindByID(ctx, id)
}

func (s *userQueryService) GetByEmail(ctx context.Context, email Email) (User, error) {
	return s.repository.FindByEmail(ctx, email)
}

func (s *userQueryService) GetByPhone(ctx context.Context, phone PhoneNumber) (User, error) {
	return s.repository.FindByPhone(ctx, phone)
}

func (s *userQueryService) GetByEmployeeID(ctx context.Context, employeeID emp.EmployeeID) (User, error) {
	return s.repository.FindByEmployeeID(ctx, employeeID)
}

func (s *userQueryService) GetByCustomerID(ctx context.Context, customerID cust.CustomerID) (User, error) {
	return s.repository.FindByCustomerID(ctx, customerID)
}

func (s *userQueryService) GetByOAuthProvider(ctx context.Context, provider string, providerID string) (User, error) {
	return s.repository.FindByOAuthProvider(ctx, provider, providerID)
}

func (s *userQueryService) GetSpecification(ctx context.Context, spec UserSpecification) (page.Page[User], error) {
	return s.repository.FindSpecification(ctx, spec)
}
