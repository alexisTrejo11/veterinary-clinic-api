package errors

import (
	log "clinic-vet-api/internal/shared/log"
	"context"
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

// BaseApplicationError represents domain-specific errors with structured logging
type BaseApplicationError struct {
	Code       string            `json:"code"`
	Type       string            `json:"type"`
	Message    string            `json:"message"`
	StatusCode int               `json:"-"`
	Details    map[string]string `json:"details,omitempty"`
	Logged     bool              `json:"-"` // Prevent double logging
}

func (e BaseApplicationError) Error() string {
	return e.Message
}

func (e BaseApplicationError) ErrorCode() string {
	return e.Code
}

func (e BaseApplicationError) ErrorType() string {
	return e.Type
}

func (e BaseApplicationError) HTTPStatus() int {
	if e.StatusCode == 0 {
		return http.StatusInternalServerError
	}
	return e.StatusCode
}

func (e BaseApplicationError) DetailMap() map[string]string {
	return e.Details
}

// Log logs the error with structured logging if not already logged
func (e *BaseApplicationError) Log(ctx context.Context, operation string) {
	if e.Logged {
		return
	}

	fields := []zap.Field{
		zap.String("error_code", e.Code),
		zap.String("error_type", e.Type),
		zap.Int("status_code", e.HTTPStatus()),
		zap.String("operation", operation),
	}

	// Add details as separate fields
	for key, value := range e.Details {
		fields = append(fields, zap.String(key, value))
	}

	// Add correlation ID from context if available
	if ctx != nil {
		if correlationID, ok := ctx.Value("correlation_id").(string); ok && correlationID != "" {
			fields = append(fields, zap.String("correlation_id", correlationID))
		}
	}

	switch e.Type {
	case "domain", "validation":
		log.Warn(e.Message, fields...)
	case "application", "business":
		log.Error(e.Message, nil, fields...)
	default:
		log.Error(e.Message, nil, fields...)
	}

	e.Logged = true
}

// BusinessRuleError creates a business rule violation error with logging
func BusinessRuleError(ctx context.Context, rule, entity, field, operation string) error {
	err := &BaseApplicationError{
		Code:    "BUSINESS_RULE_VIOLATION",
		Type:    "domain",
		Message: fmt.Sprintf("Business rule violation: %s", rule),
		Details: map[string]string{
			"rule":      rule,
			"entity":    entity,
			"field":     field,
			"operation": operation,
		},
		StatusCode: http.StatusUnprocessableEntity,
	}

	err.Log(ctx, operation)
	return err
}

// ValidationError creates a validation error with logging
func ValidationError(ctx context.Context, field, value, message, operation string) error {
	err := &BaseApplicationError{
		Code:    "VALIDATION_ERROR",
		Type:    "validation",
		Message: message,
		Details: map[string]string{
			"field":     field,
			"value":     value,
			"operation": operation,
		},
		StatusCode: http.StatusUnprocessableEntity,
	}

	err.Log(ctx, operation)
	return err
}

// InvalidEnumValue creates an invalid enum value error
func InvalidEnumValue(ctx context.Context, field, value, message, operation string) error {
	return ValidationError(ctx, field, value,
		fmt.Sprintf("Invalid enum value for %s: %s. %s", field, value, message),
		operation)
}

// MissingFieldError creates a missing field error
func MissingFieldError(ctx context.Context, field, message, operation string) error {
	return ValidationError(ctx, field, "",
		fmt.Sprintf("Missing required field: %s. %s", field, message),
		operation)
}

// MissingEntity creates a missing entity error
func MissingEntity(ctx context.Context, entity, message, operation string) error {
	return ValidationError(ctx, entity, "",
		fmt.Sprintf("Missing required entity: %s. %s", entity, message),
		operation)
}

// InvalidFieldFormat creates an invalid field format error
func InvalidFieldFormat(ctx context.Context, field, value, message, operation string) error {
	return ValidationError(ctx, field, value,
		fmt.Sprintf("Invalid format for field %s: %s. %s", field, value, message),
		operation)
}

// InvalidFieldValue creates an invalid field value error
func InvalidFieldValue(ctx context.Context, field, value, message, operation string) error {
	return ValidationError(ctx, field, value,
		fmt.Sprintf("Invalid value for field %s: %s. %s", field, value, message),
		operation)
}

// RequiredField creates a required field error
func RequiredField(ctx context.Context, field, message, operation string) error {
	return MissingFieldError(ctx, field, message, operation)
}

// EntityNotFoundError creates an entity not found error
func EntityNotFoundError(ctx context.Context, entityType, entityID, operation string) error {
	err := &BaseApplicationError{
		Code:    "ENTITY_NOT_FOUND",
		Type:    "application",
		Message: fmt.Sprintf("%s with ID %s not found", entityType, entityID),
		Details: map[string]string{
			"entity_type": entityType,
			"entity_id":   entityID,
			"operation":   operation,
		},
		StatusCode: http.StatusNotFound,
	}

	err.Log(ctx, operation)
	return err
}

// ConflictError creates a resource conflict error
func ConflictError(ctx context.Context, resource, message, operation string) error {
	err := &BaseApplicationError{
		Code:       "RESOURCE_CONFLICT",
		Type:       "application",
		Message:    message,
		StatusCode: http.StatusConflict,
		Details: map[string]string{
			"resource":  resource,
			"operation": operation,
		},
	}

	err.Log(ctx, operation)
	return err
}

// WrapError wraps an existing error with domain context
func WrapError(ctx context.Context, err error, message, entity, field, operation string) error {
	if domainErr, ok := err.(*BaseApplicationError); ok {
		// Already a domain error, just add context
		domainErr.Details["wrapped_entity"] = entity
		domainErr.Details["wrapped_field"] = field
		domainErr.Details["wrapped_operation"] = operation
		domainErr.Message = fmt.Sprintf("%s: %s", message, domainErr.Message)
		return domainErr
	}

	// Create new domain error wrapping the original
	return BusinessRuleError(ctx, err.Error(), entity, field, operation)
}

// UnauthorizedError creates an unauthorized access error
func UnauthorizedError(ctx context.Context, operation, resource string) error {
	err := &BaseApplicationError{
		Code:       "UNAUTHORIZED",
		Type:       "security",
		Message:    "Unauthorized access",
		StatusCode: http.StatusUnauthorized,
		Details: map[string]string{
			"operation": operation,
			"resource":  resource,
		},
	}

	err.Log(ctx, operation)
	return err
}

// ForbiddenError creates a forbidden access error
func ForbiddenError(ctx context.Context, operation, resource, reason string) error {
	err := &BaseApplicationError{
		Code:       "FORBIDDEN",
		Type:       "security",
		Message:    "Access forbidden",
		StatusCode: http.StatusForbidden,
		Details: map[string]string{
			"operation": operation,
			"resource":  resource,
			"reason":    reason,
		},
	}

	err.Log(ctx, operation)
	return err
}
