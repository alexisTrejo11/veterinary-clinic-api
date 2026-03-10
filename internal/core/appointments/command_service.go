package appointments

import (
	"context"
	"fmt"

	"clinic-vet-api/internal/core/customers"
	"clinic-vet-api/internal/core/employees"
)

// CommandService defines write operations for appointments
type CommandService interface {
	RequestAppointment(ctx context.Context, cmd RequestByCustomerCommand) error
	CreateAppointment(ctx context.Context, cmd CreateCommand) (Appointment, error)
	CancelAppointment(ctx context.Context, cmd CancelCommand) error
	DeleteAppointment(ctx context.Context, cmd DeleteCommand) error
	ConfirmAppointment(ctx context.Context, cmd ConfirmCommand) error
	MarkAppointmentAsNotAttend(ctx context.Context, cmd NotAttendCommand) error
	RescheduleAppointment(ctx context.Context, cmd RescheduleCommand) error
	CompleteAppointment(ctx context.Context, cmd CompleteCommand) error
	UpdateAppointmentGeneralInfo(ctx context.Context, cmd UpdateCommand) error
}

type commandService struct {
	repository         AppointmentRepository
	customerRepository customers.CustomerRepository
	employeeRepository employees.EmployeeRepository
	domainService      AppointmentDomainService
}

// NewCommandService creates a new command service instance
func NewCommandService(
	repository AppointmentRepository,
	customerRepository customers.CustomerRepository,
	employeeRepository employees.EmployeeRepository,
	domainService AppointmentDomainService,
) CommandService {
	return &commandService{
		repository:         repository,
		customerRepository: customerRepository,
		employeeRepository: employeeRepository,
		domainService:      domainService,
	}
}

// =========================================================================
// Command Service implementation
// =========================================================================

// RequestAppointment creates an appointment request that needs to be approved by a vet.
// It starts as pending and can be confirmed by a vet later.
// Must be performed by customers.
func (s *commandService) RequestAppointment(ctx context.Context, cmd RequestByCustomerCommand) error {
	appointment := Appointment{
		PetID:         cmd.PetID,
		CustomerID:    cmd.CustomerID,
		ScheduledDate: cmd.RequestedDate,
		Status:        AppointmentStatusPending,
		Notes:         cmd.Notes,
		Reason:        cmd.Reason,
		Service:       cmd.Service,
	}

	// Validate appointment business rules
	if err := appointment.ValidateInsert(ctx); err != nil {
		return err
	}

	// Validate capacity (no more than 5 appointments per hour)
	if err := s.domainService.ValidateCapacity(ctx, appointment.ScheduledDate); err != nil {
		return err
	}

	if err := s.repository.Save(ctx, &appointment); err != nil {
		return err
	}

	// TODO: Send notification to staff about new appointment request
	return nil
}

// CreateAppointment creates an appointment that doesn't need approval.
// Must be performed by staff.
func (s *commandService) CreateAppointment(ctx context.Context, cmd CreateCommand) (Appointment, error) {
	status := AppointmentStatusPending
	if cmd.Status != nil {
		status = *cmd.Status
	}

	appointment := Appointment{
		PetID:         cmd.PetID,
		CustomerID:    cmd.CustomerID,
		EmployeeID:    cmd.VetID,
		Service:       cmd.Service,
		ScheduledDate: cmd.Datetime,
		Status:        status,
		Notes:         cmd.Notes,
		Reason:        cmd.Reason,
	}

	// Validate appointment business rules
	if err := appointment.ValidateInsert(ctx); err != nil {
		return Appointment{}, err
	}

	// Validate no overlapping appointments for the assigned employee
	if err := s.domainService.ValidateNoOverlapping(ctx, &appointment); err != nil {
		return Appointment{}, err
	}

	// Validate capacity (no more than 5 appointments per hour)
	if err := s.domainService.ValidateCapacity(ctx, appointment.ScheduledDate); err != nil {
		return Appointment{}, err
	}

	if err := s.repository.Save(ctx, &appointment); err != nil {
		return Appointment{}, err
	}

	// TODO: Send notification to staff about new appointment
	return appointment, nil
}

// CancelAppointment cancels an existing appointment
func (s *commandService) CancelAppointment(ctx context.Context, cmd CancelCommand) error {
	appointment, err := s.getApptByEmployeeOpt(ctx, cmd.ID, cmd.EmployeeID)
	if err != nil {
		return err
	}

	if err := appointment.Cancel(ctx); err != nil {
		return err
	}

	if err := s.repository.Save(ctx, &appointment); err != nil {
		return err
	}

	// TODO: Send cancellation notification to customer and staff
	return nil
}

// ConfirmAppointment confirms a pending appointment
func (s *commandService) ConfirmAppointment(ctx context.Context, cmd ConfirmCommand) error {
	appointment, err := s.repository.FindByID(ctx, cmd.ID)
	if err != nil {
		return err
	}

	if err := appointment.Confirm(ctx, cmd.EmployeeID); err != nil {
		return err
	}

	if err := s.repository.Save(ctx, &appointment); err != nil {
		return err
	}

	return nil
}

// DeleteAppointment removes an appointment
func (s *commandService) DeleteAppointment(ctx context.Context, cmd DeleteCommand) error {
	appointment, err := s.repository.FindByID(ctx, cmd.ID)
	if err != nil {
		return err
	}

	if active := appointment.IsActive(); active {
		return fmt.Errorf("cannot delete active appointment")
	}

	return s.repository.DeleteByID(ctx, cmd.ID, cmd.IsHardDelete)
}

// MarkAppointmentAsNotAttend marks an appointment as not attended
func (s *commandService) MarkAppointmentAsNotAttend(ctx context.Context, cmd NotAttendCommand) error {
	appointment, err := s.repository.FindByID(ctx, cmd.ID)
	if err != nil {
		return err
	}

	if err := appointment.MarkAsNotPresented(ctx); err != nil {
		return err
	}

	if err := s.repository.Save(ctx, &appointment); err != nil {
		return err
	}

	return nil
}

// CompleteAppointment marks an appointment as completed
func (s *commandService) CompleteAppointment(ctx context.Context, cmd CompleteCommand) error {
	appointment, err := s.getApptByEmployeeOpt(ctx, cmd.ID, cmd.EmployeeID)
	if err != nil {
		return err
	}

	if err := appointment.Complete(ctx); err != nil {
		return err
	}

	if err := s.repository.Save(ctx, &appointment); err != nil {
		return err
	}

	return nil
}

// UpdateAppointmentGeneralInfo updates general information of an appointment
func (s *commandService) UpdateAppointmentGeneralInfo(ctx context.Context, cmd UpdateCommand) error {
	appointment, err := s.repository.FindByID(ctx, cmd.ID)
	if err != nil {
		return err
	}

	if err := appointment.Update(ctx, cmd.Notes, cmd.Service, cmd.Reason); err != nil {
		return err
	}

	if err := s.repository.Save(ctx, &appointment); err != nil {
		return err
	}

	return nil
}

// RescheduleAppointment changes the scheduled date of an appointment
func (s *commandService) RescheduleAppointment(ctx context.Context, cmd RescheduleCommand) error {
	appointment, err := s.repository.FindByID(ctx, cmd.ID)
	if err != nil {
		return err
	}

	// Validate no overlapping appointments for the employee at new time
	tempAppointment := appointment
	tempAppointment.ScheduledDate = cmd.NewDateTime
	if err := s.domainService.ValidateNoOverlapping(ctx, &tempAppointment); err != nil {
		return err
	}

	// Validate capacity at new time
	if err := s.domainService.ValidateCapacity(ctx, cmd.NewDateTime); err != nil {
		return err
	}

	if err := appointment.Reschedule(ctx, cmd.NewDateTime); err != nil {
		return err
	}

	if err := s.repository.Save(ctx, &appointment); err != nil {
		return err
	}

	return nil
}

func (s *commandService) getApptByEmployeeOpt(
	ctx context.Context,
	appointID AppointmentID,
	empIdOpt *employees.EmployeeID,
) (Appointment, error) {
	if empIdOpt == nil {
		return s.repository.FindByID(ctx, appointID)
	}

	return s.repository.FindByIDAndEmployeeID(ctx, appointID, *empIdOpt)
}
