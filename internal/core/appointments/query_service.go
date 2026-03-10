package appointments

import (
	"context"
	"fmt"
	"time"

	"clinic-vet-api/internal/core/customers"
	"clinic-vet-api/internal/core/employees"
	"clinic-vet-api/internal/core/pets"
	p "clinic-vet-api/internal/shared/page"
)

// =========================================================================
// Query definitions for appointment operations
// =========================================================================

// GetStatsQuery represents a query to get appointment statistics
type GetStatsQuery struct {
	StartDate  *time.Time
	EndDate    *time.Time
	EmployeeID *employees.EmployeeID
	CustomerID *customers.CustomerID
}

type GetBySpecQuery struct {
	spec *AppointmentSpecification
}

// NewGetBySpecQuery creates a new GetBySpecQuery with the given specification
func NewGetBySpecQuery(spec *AppointmentSpecification) GetBySpecQuery {
	return GetBySpecQuery{spec: spec}
}

// =========================================================================
// Query Service for retrieving appointment information
// =========================================================================

// QueryService defines read operations for appointments
type QueryService interface {
	GetByID(ctx context.Context, id AppointmentID) (Appointment, error)
	GetByIDAndCustomerID(ctx context.Context, id AppointmentID, customerID customers.CustomerID) (Appointment, error)
	GetByIDAndEmployeeID(ctx context.Context, id AppointmentID, empID employees.EmployeeID) (Appointment, error)
	GetBySpecfication(ctx context.Context, query GetBySpecQuery) (p.Page[Appointment], error)
	GetByDateRange(ctx context.Context, startDate, endDate time.Time, pagination p.Pagination) (p.Page[Appointment], error)
	GetByCustomerID(ctx context.Context, customerID customers.CustomerID, pagination p.Pagination) (p.Page[Appointment], error)
	GetByEmployeeID(ctx context.Context, empID employees.EmployeeID, pagination p.Pagination) (p.Page[Appointment], error)
	GetByPetID(ctx context.Context, petID pets.PetID, pagination p.Pagination) (p.Page[Appointment], error)
	GetSummary(ctx context.Context, q GetStatsQuery) (AppointmentSummaryStats, error)
}

type queryService struct {
	repository         AppointmentRepository
	customerRepository customers.CustomerRepository
	employeeRepository employees.EmployeeRepository
}

// NewAppointmentQueryHandler creates a new query service instance
func NewAppointmentQueryHandler(
	Repository AppointmentRepository,
	customerRepository customers.CustomerRepository,
	employeeRepository employees.EmployeeRepository,
) QueryService {
	return &queryService{
		repository:         Repository,
		customerRepository: customerRepository,
		employeeRepository: employeeRepository,
	}
}

// =========================================================================
// Query Service implementation
// =========================================================================

func (s *queryService) GetByID(ctx context.Context, id AppointmentID) (Appointment, error) {
	return s.repository.FindByID(ctx, id)
}

func (s *queryService) GetByIDAndCustomerID(
	ctx context.Context,
	id AppointmentID,
	customerID customers.CustomerID,
) (Appointment, error) {
	return s.repository.FindByIDAndCustomerID(ctx, id, customerID)
}

func (s *queryService) GetByIDAndEmployeeID(
	ctx context.Context,
	id AppointmentID,
	empID employees.EmployeeID,
) (Appointment, error) {
	return s.repository.FindByIDAndEmployeeID(ctx, id, empID)
}

func (s *queryService) GetBySpecfication(
	ctx context.Context,
	specQuery GetBySpecQuery,
) (p.Page[Appointment], error) {
	return s.repository.Find(ctx, specQuery.spec)
}

func (s *queryService) GetByDateRange(
	ctx context.Context,
	startDate, endDate time.Time,
	pagination p.Pagination,
) (p.Page[Appointment], error) {
	specification := SpecByDateRange(startDate, endDate).FromPagination(pagination)
	return s.repository.Find(ctx, specification)
}

func (s *queryService) GetByCustomerID(
	ctx context.Context,
	customerID customers.CustomerID,
	pagination p.Pagination,
) (p.Page[Appointment], error) {
	specification := SpecByCustomer(customerID).FromPagination(pagination)
	return s.repository.Find(ctx, specification)
}

func (s *queryService) GetByEmployeeID(
	ctx context.Context,
	empID employees.EmployeeID,
	pagination p.Pagination,
) (p.Page[Appointment], error) {
	specification := SpecByEmployee(empID).FromPagination(pagination)
	return s.repository.Find(ctx, specification)
}

func (s *queryService) GetByPetID(
	ctx context.Context,
	petID pets.PetID,
	pagination p.Pagination,
) (p.Page[Appointment], error) {
	specification := SpecByPet(petID).FromPagination(pagination)
	return s.repository.Find(ctx, specification)
}

func (s *queryService) GetSummary(ctx context.Context, query GetStatsQuery) (AppointmentSummaryStats, error) {
	var appointments []Appointment
	var err error
	maxPage := p.Pagination{
		Number: 1,
		Size:   10000,
	}

	if query.StartDate != nil && query.EndDate != nil {
		specification := SpecByDateRange(*query.StartDate, *query.EndDate).FromPagination(maxPage)
		appointmentsPage, dberr := s.repository.Find(ctx, specification)
		appointments = appointmentsPage.Items
		err = dberr
	} else {
		specification := NewAppointmentSpecification().FromPagination(maxPage)
		appointmentsPage, dberr := s.repository.Find(ctx, specification)
		appointments = appointmentsPage.Items
		err = dberr
	}

	if err != nil {
		return AppointmentSummaryStats{}, err
	}

	// Apply additional filters
	var filteredAppointments []Appointment
	for _, appointment := range appointments {
		includeAppointment := true

		// Filter by vet ID
		if query.EmployeeID != nil {
			if appointment.EmployeeID == nil || appointment.EmployeeID.Value != query.EmployeeID.Value {
				includeAppointment = false
			}
		}

		// Filter by customer ID
		if query.CustomerID != nil && !appointment.CustomerID.Equals(query.CustomerID.Value) {
			includeAppointment = false
		}

		if includeAppointment {
			filteredAppointments = append(filteredAppointments, appointment)
		}
	}

	stats := s.calculateStats(filteredAppointments, query)
	return stats, nil
}

func (s *queryService) calculateStats(appointments []Appointment, query GetStatsQuery) AppointmentSummaryStats {
	totalAppointments := len(appointments)
	pendingCount := 0
	confirmedCount := 0
	completedCount := 0
	cancelledCount := 0
	noShowCount := 0
	emergencyCount := 0

	statusBreakdown := make(map[AppointmentStatus]int)
	serviceBreakdown := make(map[ClinicService]int)

	for _, appointment := range appointments {
		// Count by status
		status := appointment.Status
		statusBreakdown[status]++

		switch status {
		case AppointmentStatusPending:
			pendingCount++
		case AppointmentStatusCompleted:
			completedCount++
		case AppointmentStatusCancelled:
			cancelledCount++
		case AppointmentStatusNotPresented:
			noShowCount++
		}

		// Count by service
		service := appointment.Service
		serviceBreakdown[service]++
	}

	// Generate period string
	var period *string
	if query.StartDate != nil && query.EndDate != nil {
		periodStr := fmt.Sprintf("%s to %s",
			query.StartDate.Format("2006-01-02"),
			query.EndDate.Format("2006-01-02"))
		period = &periodStr
	}

	return AppointmentSummaryStats{
		TotalAppointments: totalAppointments,
		PendingCount:      pendingCount,
		ConfirmedCount:    confirmedCount,
		CompletedCount:    completedCount,
		CancelledCount:    cancelledCount,
		NoShowCount:       noShowCount,
		EmergencyCount:    emergencyCount,
		StatusBreakdown:   statusBreakdown,
		ServiceBreakdown:  serviceBreakdown,
		Period:            period,
	}
}
