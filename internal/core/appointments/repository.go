package appointments

import (
	"context"

	"clinic-vet-api/internal/core/customers"
	"clinic-vet-api/internal/core/employees"
	p "clinic-vet-api/internal/shared/page"
)

type AppointmentRepository interface {
	FindByID(ctx context.Context, id AppointmentID) (Appointment, error)
	FindByIDAndCustomerID(ctx context.Context, id AppointmentID, customerID customers.CustomerID) (Appointment, error)
	FindByIDAndEmployeeID(ctx context.Context, id AppointmentID, employeeID employees.EmployeeID) (Appointment, error)
	Find(ctx context.Context, spec *AppointmentSpecification) (p.Page[Appointment], error)
	Count(ctx context.Context, spec *AppointmentSpecification) (int64, error)
	ExistsByID(ctx context.Context, id AppointmentID) (bool, error)

	Save(ctx context.Context, appointment *Appointment) error
	DeleteByID(ctx context.Context, id AppointmentID, hardDelete bool) error
	RestoreByID(ctx context.Context, id AppointmentID) error
}
