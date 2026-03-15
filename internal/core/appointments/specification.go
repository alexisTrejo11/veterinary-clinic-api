package appointments

import (
	"clinic-vet-api/internal/core/customers"
	"clinic-vet-api/internal/core/employees"
	"clinic-vet-api/internal/core/pets"
	"clinic-vet-api/internal/shared/page"
	"time"
)

// AppointmentSpecification defines search criteria for appointments
type AppointmentSpecification struct {
	ID            *AppointmentID
	CustomerID    *customers.CustomerID
	EmployeeID    *employees.EmployeeID
	PetID         *pets.PetID
	Service       *ClinicService
	Status        *AppointmentStatus
	StartDate     *time.Time
	EndDate       *time.Time
	ScheduledDate *time.Time
	page.Pagination
}

// SearchParams contains the raw parameters for database queries
type SearchParams struct {
	ID            *int32
	CustomerID    *int32
	EmployeeID    *int32
	PetID         *int32
	Service       *string
	Status        *string
	StartDate     *time.Time
	EndDate       *time.Time
	ScheduledDate *time.Time
	Limit         int32
	Offset        int32
}

// ToSearchParams converts the specification to search parameters
func (s *AppointmentSpecification) ToSearchParams() SearchParams {
	params := SearchParams{
		Limit:  int32(s.Pagination.Size),
		Offset: int32((s.Pagination.Number - 1) * s.Pagination.Size),
	}

	if s.ID != nil {
		id := int32(s.ID.Value())
		params.ID = &id
	}

	if s.CustomerID != nil {
		id := int32(s.CustomerID.Value())
		params.CustomerID = &id
	}

	if s.EmployeeID != nil {
		id := int32(s.EmployeeID.Value())
		params.EmployeeID = &id
	}

	if s.PetID != nil {
		id := int32(s.PetID.Value())
		params.PetID = &id
	}

	if s.Service != nil {
		service := string(*s.Service)
		params.Service = &service
	}

	if s.Status != nil {
		status := string(*s.Status)
		params.Status = &status
	}

	if s.StartDate != nil {
		params.StartDate = s.StartDate
	}

	if s.EndDate != nil {
		params.EndDate = s.EndDate
	}

	if s.ScheduledDate != nil {
		params.ScheduledDate = s.ScheduledDate
	}

	return params
}

// Builder methods (fluent interface)

func (s *AppointmentSpecification) WithID(id AppointmentID) *AppointmentSpecification {
	s.ID = &id
	return s
}

func (s *AppointmentSpecification) WithCustomerID(customerID customers.CustomerID) *AppointmentSpecification {
	s.CustomerID = &customerID
	return s
}

func (s *AppointmentSpecification) WithEmployeeID(employeeID employees.EmployeeID) *AppointmentSpecification {
	s.EmployeeID = &employeeID
	return s
}

func (s *AppointmentSpecification) WithPetID(petID pets.PetID) *AppointmentSpecification {
	s.PetID = &petID
	return s
}

func (s *AppointmentSpecification) WithService(service ClinicService) *AppointmentSpecification {
	s.Service = &service
	return s
}

func (s *AppointmentSpecification) WithStatus(status AppointmentStatus) *AppointmentSpecification {
	s.Status = &status
	return s
}

func (s *AppointmentSpecification) WithDateRange(startDate, endDate time.Time) *AppointmentSpecification {
	s.StartDate = &startDate
	s.EndDate = &endDate
	return s
}

func (s *AppointmentSpecification) WithScheduledDate(date time.Time) *AppointmentSpecification {
	s.ScheduledDate = &date
	return s
}

func (s *AppointmentSpecification) WithPagination(pageNumber, pageSize int, orderBy, sortDir string) *AppointmentSpecification {
	s.Pagination = page.Pagination{
		Number:  pageNumber,
		Size:    pageSize,
		OrderBy: orderBy,
		SortDir: sortDir,
	}
	return s
}

// FromPagination sets pagination from a Pagination object
func (s *AppointmentSpecification) FromPagination(pagination page.Pagination) *AppointmentSpecification {
	s.Pagination = pagination
	return s
}

// Factory functions for creating new specifications

func NewAppointmentSpecification() *AppointmentSpecification {
	return &AppointmentSpecification{
		Pagination: page.Pagination{
			Number: 1,
			Size:   50,
		},
	}
}

func SpecByID(id AppointmentID) *AppointmentSpecification {
	return NewAppointmentSpecification().WithID(id)
}

func SpecByCustomer(customerID customers.CustomerID) *AppointmentSpecification {
	return NewAppointmentSpecification().WithCustomerID(customerID)
}

func SpecByEmployee(employeeID employees.EmployeeID) *AppointmentSpecification {
	return NewAppointmentSpecification().WithEmployeeID(employeeID)
}

func SpecByPet(petID pets.PetID) *AppointmentSpecification {
	return NewAppointmentSpecification().WithPetID(petID)
}

func SpecByService(service ClinicService) *AppointmentSpecification {
	return NewAppointmentSpecification().WithService(service)
}

func SpecByStatus(status AppointmentStatus) *AppointmentSpecification {
	return NewAppointmentSpecification().WithStatus(status)
}

func SpecByDateRange(start, end time.Time) *AppointmentSpecification {
	return NewAppointmentSpecification().WithDateRange(start, end)
}

func SpecByScheduledDate(date time.Time) *AppointmentSpecification {
	return NewAppointmentSpecification().WithScheduledDate(date)
}
