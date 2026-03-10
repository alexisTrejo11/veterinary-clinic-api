package appointments

import "clinic-vet-api/internal/shared"

type AppointmentID struct{ shared.BaseID }

func NewAppointmentID(id uint) AppointmentID {
	return AppointmentID{BaseID: shared.BaseID{Value: id}}
}

type AppointmentSummaryStats struct {
	TotalAppointments int                       `json:"total_appointments"`
	PendingCount      int                       `json:"pending_appointments"`
	ConfirmedCount    int                       `json:"confirmed_appointments"`
	CompletedCount    int                       `json:"completed_appointments"`
	CancelledCount    int                       `json:"cancelled_appointments"`
	NoShowCount       int                       `json:"no_show_appointments"`
	StatusBreakdown   map[AppointmentStatus]int `json:"status_breakdown"`
	ServiceBreakdown  map[ClinicService]int     `json:"service_breakdown"`
	EmergencyCount    int                       `json:"emergency_count"`
	Period            *string                   `json:"period,omitempty"`
}
