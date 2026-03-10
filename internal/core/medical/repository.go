package medical

import (
	"clinic-vet-api/internal/core/customers"
	"clinic-vet-api/internal/core/employees"
	"clinic-vet-api/internal/core/pets"
	"context"
	"time"

	p "clinic-vet-api/internal/shared/page"
)

type MedicalSessionRepository interface {
	GetByID(ctx context.Context, id MedSessionID) (*MedicalSession, error)
	GetByIDAndCustomerID(ctx context.Context, id MedSessionID, customerID customers.CustomerID) (*MedicalSession, error)
	GetByIDAndPetID(ctx context.Context, id MedSessionID, petID pets.PetID) (*MedicalSession, error)
	GetByIDAndEmployeeID(ctx context.Context, id MedSessionID, employeeID employees.EmployeeID) (*MedicalSession, error)

	GetByEmployeeID(ctx context.Context, employeeID employees.EmployeeID, pagination p.Pagination) (p.Page[MedicalSession], error)
	GetByPetID(ctx context.Context, petID pets.PetID, pagination p.Pagination) (p.Page[MedicalSession], error)
	GetByCustomerID(ctx context.Context, customerID customers.CustomerID, pagination p.Pagination) (p.Page[MedicalSession], error)
	//GetBySpecification(ctx context.Context, spec specification.MedicalSessionSpecification) (p.Page[MedicalSession], error)
	GetRecentByPetID(ctx context.Context, petID pets.PetID, limit int) ([]MedicalSession, error)
	GetByDateRange(ctx context.Context, startDate, endDate time.Time, pagination p.Pagination) (p.Page[MedicalSession], error)
	GetByPetAndDateRange(ctx context.Context, petID pets.PetID, startDate, endDate time.Time) ([]MedicalSession, error)
	GetByDiagnosis(ctx context.Context, diagnosis string, pagination p.Pagination) (p.Page[MedicalSession], error)

	ExistsByID(ctx context.Context, id MedSessionID) (bool, error)
	ExistsByPetAndDate(ctx context.Context, petID pets.PetID, date time.Time) (bool, error)

	Save(ctx context.Context, medSession *MedicalSession) error
	SoftDelete(ctx context.Context, id MedSessionID) error
	HardDelete(ctx context.Context, id MedSessionID) error
}
