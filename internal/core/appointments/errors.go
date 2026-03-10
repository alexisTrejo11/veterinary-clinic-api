package appointments

import (
	"clinic-vet-api/internal/shared/errors"
	"context"
	"fmt"

	"net/http"
	"time"
)

// AppointmentErrorCode representa códigos de error específicos para citas
// Missing Errors, found usage or DELETE

type AppointmentErrorCode string

const (
	MinAllowedDaysToSchedule = 2  // Minimum 3 days in advance
	MaxAllowedDaysToSchedule = 29 // Maximum 30 days in advance
	BusinessStartHour        = 8  // 9 AM
	BusinessEndHour          = 17 // 6 PM
)
const (
	AppointmentNotesTooLong           AppointmentErrorCode = "APPOINTMENT_NOTES_TOO_LONG"
	AppointmentInvalidService         AppointmentErrorCode = "APPOINTMENT_INVALID_SERVICE"
	AppointmentInvalidReason          AppointmentErrorCode = "APPOINTMENT_INVALID_REASON"
	AppointmentCannotReschedule       AppointmentErrorCode = "APPOINTMENT_CANNOT_RESCHEDULE"
	AppointmentCannotCancel           AppointmentErrorCode = "APPOINTMENT_CANNOT_CANCEL"
	AppointmentCannotComplete         AppointmentErrorCode = "APPOINTMENT_CANNOT_COMPLETE"
	AppointmentCannotMarkNotPresented AppointmentErrorCode = "APPOINTMENT_CANNOT_MARK_NOT_PRESENTED"
	AppointmentCannotConfirm          AppointmentErrorCode = "APPOINTMENT_CANNOT_CONFIRM"
	AppointmentInvalidTransition      AppointmentErrorCode = "APPOINTMENT_INVALID_TRANSITION"
	AppointmentCannotDelete           AppointmentErrorCode = "APPOINTMENT_CANNOT_DELETE"
	AppointmentScheduledDateInvalid   AppointmentErrorCode = "APPOINTMENT_SCHEDULED_DATE_INVALID"
)

// appointmentValidationError crea un error de validación específico para appointment
func appointmentValidationError(ctx context.Context, code AppointmentErrorCode, field, message, operation string) error {
	return errors.ValidationError(ctx, "appointment", field,
		fmt.Sprintf("Appointment %s: %s", field, message), operation)
}

// appointmentBusinessError crea un error de negocio específico para appointment
func appointmentBusinessError(ctx context.Context, code AppointmentErrorCode, rule, operation string) error {
	return errors.BusinessRuleError(ctx, rule, "appointment", "", operation)
}

func NotesTooLongError(ctx context.Context, operation string) error {
	return appointmentValidationError(ctx, AppointmentNotesTooLong, "notes",
		"notes cannot exceed 1000 characters", operation)
}

func InvalidServiceError(ctx context.Context, service ClinicService, operation string) error {
	return appointmentValidationError(ctx, AppointmentInvalidService, "service",
		fmt.Sprintf("invalid service: %s", service), operation)
}

func InvalidReasonError(ctx context.Context, reason string, operation string) error {
	return appointmentValidationError(ctx, AppointmentInvalidReason, "reason",
		fmt.Sprintf("invalid reason: %s", reason), operation)
}

func CannotRescheduleError(ctx context.Context, currentStatus AppointmentStatus, operation string) error {
	return appointmentBusinessError(ctx, AppointmentCannotReschedule, "reschedule", operation)
}

func CannotCancelError(ctx context.Context, currentStatus AppointmentStatus, operation string) error {
	return appointmentBusinessError(ctx, AppointmentCannotCancel, "cancel", operation)
}

func CannotCompleteError(ctx context.Context, currentStatus AppointmentStatus, operation string) error {
	return appointmentBusinessError(ctx, AppointmentCannotComplete, "complete", operation)
}

func CannotMarkNotPresentedError(ctx context.Context, currentStatus AppointmentStatus, operation string) error {
	return appointmentBusinessError(ctx, AppointmentCannotMarkNotPresented, "mark_not_presented", operation)
}

func CannotConfirmError(ctx context.Context, currentStatus AppointmentStatus, operation string) error {
	return appointmentBusinessError(ctx, AppointmentCannotConfirm, "confirm", operation)
}

func InvalidTransitionError(ctx context.Context, from, to AppointmentStatus, operation string) error {
	return appointmentBusinessError(ctx, AppointmentInvalidTransition, "status_transition", operation)
}

func CannotDeleteError(ctx context.Context, currentStatus AppointmentStatus, operation string) error {
	return appointmentBusinessError(ctx, AppointmentCannotDelete, "delete", operation)
}

func ScheduledDateInvalidError(ctx context.Context, message, operation string) error {
	return appointmentValidationError(ctx, AppointmentScheduledDateInvalid, "scheduled_date",
		message, operation)
}

// Application Errors
const (
	ErrApptNotFound             = "appointment not found"
	ErrMarkAsNotPresentedFailed = "failed to mark appointment as not presented"
	ErrSaveApptFailed           = "failed to save appointment"
	ErrInvalidCommandType       = "invalid command type"
	ErrFindingAppt              = "error finding appointment"
	ErrCannotDeleteCompleted    = "cannot delete completed appointment"
	ErrFailedToDelete           = "failed to delete appointment"
	ErrFailedToCancel           = "failed to cancel appointment"
	ErrUpdateApptFailed         = "failed to update appointment"
	ErrDeleteApptFailed         = "failed to delete appointment"
	ErrMarkNotPresentedFailed   = "failed to mark as not presented"

	SuccessApptUpdated          = "appointment updated successfully"
	SuccessApptDeleted          = "appointment deleted successfully"
	SuccessApptCanceled         = "appointment canceled successfully"
	SuccessMarkedAsNotPresented = "appointment marked as not presented successfully"
)

const (
	AppointmentSchedulingTooEarly    AppointmentErrorCode = "APPOINTMENT_SCHEDULING_TOO_EARLY"
	AppointmentSchedulingTooLate     AppointmentErrorCode = "APPOINTMENT_SCHEDULING_TOO_LATE"
	AppointmentSchedulingWeekend     AppointmentErrorCode = "APPOINTMENT_SCHEDULING_WEEKEND"
	AppointmentSchedulingAfterHours  AppointmentErrorCode = "APPOINTMENT_SCHEDULING_AFTER_HOURS"
	AppointmentSchedulingPastDate    AppointmentErrorCode = "APPOINTMENT_SCHEDULING_PAST_DATE"
	AppointmentSchedulingDateZero    AppointmentErrorCode = "APPOINTMENT_SCHEDULING_DATE_ZERO"
	AppointmentSchedulingUnavailable AppointmentErrorCode = "APPOINTMENT_SCHEDULING_UNAVAILABLE"
)

func AppointmentScheduleDateZeroErr(ctx context.Context) error {
	err := fmt.Errorf("customer ID must be greater than zero")
	return errors.BusinessRuleError(ctx, err.Error(), "Appointment", "schedule_date", "schedule-date validatio")
}

func AppointmentScheduleDateRuleErr(ctx context.Context, ruleViolated string) error {
	return errors.BusinessRuleError(ctx, ruleViolated, "Appointment", "schedule_date", "schedule-date validation")
}

func AppointmentStatusTransitionErr(fromStatus, toStatus string, message string) error {
	return &errors.BaseApplicationError{
		Code:    "BUSINESS_RULE_VIOLATION",
		Type:    "domain",
		Message: message,
		Details: map[string]string{
			"from_status": fromStatus,
			"to_status":   toStatus,
			"entity":      "appointment",
			"field":       "schedule_date",
		},
	}
}

// AppointmentSchedulingError crea un error de scheduling específico para Appointment
func AppointmentSchedulingError(ctx context.Context, code AppointmentErrorCode, date, operation string, details map[string]string) error {
	messages := map[AppointmentErrorCode]string{
		AppointmentSchedulingTooEarly:    "Appointment must be scheduled at least %d days in advance",
		AppointmentSchedulingTooLate:     "Appointment cannot be scheduled more than %d days in advance",
		AppointmentSchedulingWeekend:     "Appointments cannot be scheduled on weekends",
		AppointmentSchedulingAfterHours:  "Appointments can only be scheduled during business hours (9 AM to 6 PM)",
		AppointmentSchedulingPastDate:    "Scheduled date cannot be in the past",
		AppointmentSchedulingDateZero:    "Scheduled date is required",
		AppointmentSchedulingUnavailable: "The selected time slot is not available",
	}

	// Mensaje base
	message := messages[code]

	// Agregar detalles específicos al mapa de detalles
	if details == nil {
		details = make(map[string]string)
	}
	details["scheduled_date"] = date
	details["min_advance_days"] = fmt.Sprintf("%d", MinAllowedDaysToSchedule)
	details["max_advance_days"] = fmt.Sprintf("%d", MaxAllowedDaysToSchedule)

	err := &errors.BaseApplicationError{
		Code:       string(code),
		Type:       "business",
		Message:    message,
		Details:    details,
		StatusCode: http.StatusUnprocessableEntity,
	}

	err.Log(ctx, operation)
	return err
}

// Funciones específicas para cada tipo de error de scheduling

func AppointmentTooEarlyError(ctx context.Context, date time.Time, operation string) error {
	return AppointmentSchedulingError(
		ctx,
		AppointmentSchedulingTooEarly,
		date.Format(time.RFC3339),
		operation,
		map[string]string{
			"requested_date":   date.Format(time.RFC3339),
			"min_allowed_date": time.Now().AddDate(0, 0, MinAllowedDaysToSchedule).Format(time.RFC3339),
		},
	)
}

func AppointmentTooLateError(ctx context.Context, date time.Time, operation string) error {
	return AppointmentSchedulingError(
		ctx,
		AppointmentSchedulingTooLate,
		date.Format(time.RFC3339),
		operation,
		map[string]string{
			"requested_date":   date.Format(time.RFC3339),
			"max_allowed_date": time.Now().AddDate(0, 0, MaxAllowedDaysToSchedule).Format(time.RFC3339),
		},
	)
}

func AppointmentWeekendError(ctx context.Context, date time.Time, operation string) error {
	return AppointmentSchedulingError(
		ctx,
		AppointmentSchedulingWeekend,
		date.Format(time.RFC3339),
		operation,
		map[string]string{
			"requested_date": date.Format(time.RFC3339),
			"weekday":        date.Weekday().String(),
		},
	)
}

func AppointmentAfterHoursError(ctx context.Context, date time.Time, operation string) error {
	return AppointmentSchedulingError(
		ctx,
		AppointmentSchedulingAfterHours,
		date.Format(time.RFC3339),
		operation,
		map[string]string{
			"requested_date": date.Format(time.RFC3339),
			"requested_hour": fmt.Sprintf("%d:00", date.Hour()),
			"business_hours": "9:00 - 18:00",
		},
	)
}

func AppointmentPastDateError(ctx context.Context, date time.Time, operation string) error {
	return AppointmentSchedulingError(
		ctx,
		AppointmentSchedulingPastDate,
		date.Format(time.RFC3339),
		operation,
		map[string]string{
			"requested_date": date.Format(time.RFC3339),
			"current_date":   time.Now().Format(time.RFC3339),
		},
	)
}

func AppointmentDateZeroError(ctx context.Context, operation string) error {
	return AppointmentSchedulingError(
		ctx,
		AppointmentSchedulingDateZero,
		"",
		operation,
		map[string]string{
			"current_date": time.Now().Format(time.RFC3339),
		},
	)
}

func AppointmentUnavailableError(ctx context.Context, date time.Time, reason, operation string) error {
	return AppointmentSchedulingError(
		ctx,
		AppointmentSchedulingUnavailable,
		date.Format(time.RFC3339),
		operation,
		map[string]string{
			"requested_date":     date.Format(time.RFC3339),
			"unavailable_reason": reason,
		},
	)
}
