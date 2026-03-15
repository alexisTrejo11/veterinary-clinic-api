package appointments

import (
	"context"
	"fmt"
	"time"

	p "clinic-vet-api/internal/shared/page"
)

// AppointmentDomainService handles complex business rules that require repository access
type AppointmentDomainService interface {
	ValidateNoOverlapping(ctx context.Context, appointment *Appointment) error
	ValidateCapacity(ctx context.Context, scheduledDate time.Time) error
}

type appointmentDomainService struct {
	repository AppointmentRepository
}

// NewAppointmentDomainService creates a new domain service instance
func NewAppointmentDomainService(repository AppointmentRepository) AppointmentDomainService {
	return &appointmentDomainService{
		repository: repository,
	}
}

// ValidateNoOverlapping ensures the employee doesn't have overlapping appointments
func (s *appointmentDomainService) ValidateNoOverlapping(ctx context.Context, appointment *Appointment) error {
	if appointment.EmployeeID == nil {
		// No employee assigned yet, skip validation
		return nil
	}

	// Check appointments for the same employee within ±30 minutes
	startTime := appointment.ScheduledDate.Add(-30 * time.Minute)
	endTime := appointment.ScheduledDate.Add(30 * time.Minute)

	spec := NewAppointmentSpecification().
		WithEmployeeID(*appointment.EmployeeID).
		WithDateRange(startTime, endTime).
		FromPagination(p.Pagination{Number: 1, Size: 100})

	page, err := s.repository.Find(ctx, spec)
	if err != nil {
		return fmt.Errorf("failed to check overlapping appointments: %w", err)
	}

	// Check for conflicts, excluding the current appointment if it exists
	for _, existing := range page.Items {
		// Skip if it's the same appointment (update case)
		if !appointment.ID.IsZero() && existing.ID.Equals(appointment.ID.Value()) {
			continue
		}

		// Skip cancelled or completed appointments
		if existing.Status == AppointmentStatusCancelled ||
			existing.Status == AppointmentStatusNotPresented {
			continue
		}

		// Check if dates are too close (within 30 minutes)
		timeDiff := existing.ScheduledDate.Sub(appointment.ScheduledDate)
		if timeDiff < 0 {
			timeDiff = -timeDiff
		}

		if timeDiff < 30*time.Minute {
			return fmt.Errorf("employee already has an appointment scheduled at %s",
				existing.ScheduledDate.Format("2006-01-02 15:04"))
		}
	}

	return nil
}

// ValidateCapacity ensures no more than 5 appointments are scheduled at the same time
func (s *appointmentDomainService) ValidateCapacity(ctx context.Context, scheduledDate time.Time) error {
	const maxConcurrentAppointments = 5

	// Round to hour for capacity check
	hourStart := time.Date(
		scheduledDate.Year(),
		scheduledDate.Month(),
		scheduledDate.Day(),
		scheduledDate.Hour(),
		0, 0, 0,
		scheduledDate.Location(),
	)
	hourEnd := hourStart.Add(1 * time.Hour)

	spec := NewAppointmentSpecification().
		WithDateRange(hourStart, hourEnd).
		FromPagination(p.Pagination{Number: 1, Size: 100})

	page, err := s.repository.Find(ctx, spec)
	if err != nil {
		return fmt.Errorf("failed to check appointment capacity: %w", err)
	}

	// Count active appointments (not cancelled or no-show)
	activeCount := 0
	for _, appointment := range page.Items {
		if appointment.Status != AppointmentStatusCancelled &&
			appointment.Status != AppointmentStatusNotPresented {
			activeCount++
		}
	}

	if activeCount >= maxConcurrentAppointments {
		return fmt.Errorf("maximum capacity reached for %s (limit: %d appointments per hour)",
			hourStart.Format("2006-01-02 15:00"), maxConcurrentAppointments)
	}

	return nil
}
