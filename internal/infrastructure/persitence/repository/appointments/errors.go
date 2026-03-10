package appointments

import (
	"context"
	"fmt"

	"clinic-vet-api/internal/shared/errors"
)

const (
	OpSelect = "select"
	OpInsert = "insert"
	OpUpdate = "update"
	OpDelete = "delete"
	OpCount  = "count"
	OpSearch = "search"

	TableAppts = "appointments"
	DriverSQL  = "sql"
)

const (
	ErrMsgListAppts       = "failed to list appointments"
	ErrMsgGetAppt         = "failed to get appointment"
	ErrMsgSearchAppts     = "failed to search appointments"
	ErrMsgCountAppts      = "failed to count appointments"
	ErrMsgCreateAppt      = "failed to create appointment"
	ErrMsgUpdateAppt      = "failed to update appointment"
	ErrMsgDeleteAppt      = "failed to delete appointment"
	ErrMsgConvertToDomain = "failed to convert to domain entity"
	ErrMsgNotFound        = "appointment not found"
)

// dbError creates a standardized database operation error
func (r *SqlcAppointmentRepository) dbError(operation, message string, err error) error {
	return errors.DatabaseError(operation, TableAppts, DriverSQL, fmt.Errorf("%s: %v", message, err))
}

// notFoundError creates a standardized entity not found error
func (r *SqlcAppointmentRepository) notFoundError(parameterName, parameterValue string) error {
	return errors.DBNotFoundError(parameterName, parameterValue, OpSelect, TableAppts, DriverSQL)
}

// wrapConversionError wraps domain conversion errors
func (r *SqlcAppointmentRepository) wrapConversionError(err error) error {
	return errors.WrapError(context.Background(), err, OpSelect, TableAppts, DriverSQL, ErrMsgConvertToDomain)
}
