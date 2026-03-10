package dtos

import (
	"clinic-vet-api/internal/core/appointments"
	"clinic-vet-api/internal/shared/page"
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// =======================================================================
// Request DTOs
// =======================================================================

// AppointmentsByDateRangeRequest represents the request to get appointments by date range
// @Description Request query parameters for fetching appointments within a specified date range
type AppointmentsByDateRangeRequest struct {
	StartDate string `json:"start_date"  form:"start_date" binding:"required"`
	EndDate   string `json:"end_date"  form:"end_date" binding:"required"`
	page.PaginationRequest
}

func (r *AppointmentsByDateRangeRequest) StartTime() time.Time {
	start, err := time.Parse("2006-01-02", r.StartDate)
	if err != nil {
		return time.Time{}
	}
	return start
}

func (r *AppointmentsByDateRangeRequest) EndTime() time.Time {
	end, err := time.Parse("2006-01-02", r.EndDate)
	if err != nil {
		return time.Time{}
	}
	return end
}

// CreateApptRequest represents the request to create an appointment
// @Description Request body for creating a new appointment
type AppointmentCreateRequest struct {
	// ID of the customer making the appointment
	// Required: true
	CustomerID uint `json:"customer_id" binding:"required" example:"123"`

	// ID of the pet for the appointment
	// Required: true
	PetID uint `json:"pet_id" binding:"required" example:"456"`

	// Date and time of the appointment
	// Required: true
	// Format: RFC3339
	ScheduledDate time.Time `json:"scheduled_date" binding:"required" example:"2024-03-15T10:30:00Z"`

	// Reason for the appointment
	// Required: true
	Reason string `json:"reason" binding:"required" example:"Annual checkup"`

	// Service requested
	// Required: true
	Service string `json:"service" binding:"required" example:"Vaccination"`

	// Additional notes for the appointment (optional)
	Notes *string `json:"notes" example:"Patient has allergy to penicillin"`

	// Status of the appointment
	// Required: false (defaults to "scheduled")
	Status *string `json:"status" binding:"required" example:"scheduled"`
}

type AppointmentUpdateGeneralInfoRequest struct {
	// Optional: Reason for updating the appointment
	Reason *string `json:"reason,omitempty" binding:"omitempty,max=500" example:"Updated reason for appointment"`

	// Optional: Additional notes or observations about the appointment
	Notes *string `json:"notes,omitempty" binding:"omitempty,max=1000" example:"Patient requires special handling"`

	// Optional: Service type for the appointment (e.g., "checkup", "vaccination", "surgery")
	Service *string `json:"service,omitempty" binding:"omitempty,max=100" example:"Annual checkup"`
}

type ReassignAppointmentRequest struct {
	// ID of the new employee (vet) to assign the appointment to
	// Required: true
	NewEmployeeID uint `json:"new_employee_id" binding:"required" example:"789"`

	// Optional: Reason for reassigning the appointment
	Reason *string `json:"reason,omitempty" binding:"omitempty,max=500" example:"Original vet is unavailable"`
}

type RescheduleAppointmentRequest struct {
	// @Required Required: New date and time for the appointment (RFC3339 format)
	NewDateTime time.Time `json:"datetime" binding:"required" example:"2024-01-15T10:30:00Z"`

	// Optional: Reason for rescheduling the appointment
	Reason *string `json:"reason,omitempty" binding:"omitempty,max=500" example:"Routine checkup"`
}

type AppointmentSearchRequest struct {
	CustomerID uint   `form:"customer_id" validate:"omitempty,min=1"`
	EmployeeID uint   `form:"vet_id" validate:"omitempty,min=1"`
	PetID      uint   `form:"pet_id" validate:"omitempty,min=1"`
	Service    string `form:"service" validate:"omitempty,clinic_service"`
	Status     string `form:"status" validate:"omitempty,appointment_status"`
	Reason     string `form:"reason" validate:"omitempty,visit_reason"`
	StartDate  string `form:"start_date" validate:"omitempty,datetime=2006-01-02"`
	EndDate    string `form:"end_date" validate:"omitempty,datetime=2006-01-02"`
	HasNotes   *bool  `form:"has_notes" validate:"omitempty"`
	page.PaginationRequest
}

func NewApptSearchRequestFromContext(c *gin.Context) (AppointmentSearchRequest, error) {
	var request AppointmentSearchRequest

	if err := c.ShouldBindQuery(&request); err != nil {
		return AppointmentSearchRequest{}, fmt.Errorf("invalid query parameters: %w", err)
	}

	request.processBooleanParams(c)

	return request, nil
}

func (r *AppointmentSearchRequest) processBooleanParams(c *gin.Context) {
	if hasNotesStr := c.Query("has_notes"); hasNotesStr != "" {
		if hasNotes, err := strconv.ParseBool(hasNotesStr); err == nil {
			r.HasNotes = &hasNotes
		}
	}
}

// =======================================================================
// Response DTOs
// =======================================================================

// AppointmentResponse represents an appointment response
type AppointmentResponse struct {
	ID            uint    `json:"id"`
	PetID         uint    `json:"pet_id"`
	CustomerID    uint    `json:"customer_id"`
	EmployeeID    uint    `json:"vet_id,omitempty"`
	Service       string  `json:"service"`
	Datetime      string  `json:"date_time"`
	ScheduledDate string  `json:"scheduled_date"`
	Status        string  `json:"status"`
	Reason        string  `json:"reason"`
	Notes         *string `json:"notes,omitempty"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
}

// AppointmentStats represents appointment statistics
type AppointmentStats struct {
	TotalAppointments     int                                    `json:"total_appointments"`
	AppointmentsByStatus  map[appointments.AppointmentStatus]int `json:"appointments_by_status"`
	AppointmentsByService map[appointments.ClinicService]int     `json:"appointments_by_service"`
	UpcomingAppointments  int                                    `json:"upcoming_appointments"`
	OverdueAppointments   int                                    `json:"overdue_appointments"`
}
