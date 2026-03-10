package repository

const (
	TableUsers = "users"
	TableAppts = "appointments"

	// User error messages
	ErrMsgFindUser            = "failed to find user"
	ErrMsgFindUserByEmail     = "failed to find user by email"
	ErrMsgFindUserByPhone     = "failed to find user by phone"
	ErrMsgFindUsers           = "failed to list users"
	ErrMsgSearchUsers         = "failed to search users"
	ErrMsgCreateUser          = "failed to create user"
	ErrMsgUpdateUser          = "failed to update user"
	ErrMsgSoftDeleteUser      = "failed to soft delete user"
	ErrMsgHardDeleteUser      = "failed to hard delete user"
	ErrMsgCheckUserExists     = "failed to check if user exists"
	ErrMsgUpdateLastLogin     = "failed to update last login"
	ErrMsgConvertUserToDomain = "failed to convert user to domain entity"

	// Appointment error messages
	ErrMsgFindAppointment        = "failed to find appointment"
	ErrMsgFindAppointmentByID    = "failed to get appointment by ID"
	ErrMsgFindAppointmentsBySpec = "failed to find appointments by specification"
	ErrMsgCreateAppointment      = "failed to create appointment"
	ErrMsgUpdateAppointment      = "failed to update appointment"
	ErrMsgSoftDeleteAppointment  = "failed to soft delete appointment"
	ErrMsgHardDeleteAppointment  = "failed to hard delete appointment"
	ErrMsgRestoreAppointment     = "failed to restore appointment"
	ErrMsgCheckAppointmentExists = "failed to check appointment existence by ID"
	ErrMsgCountAppointments      = "failed to count appointments"
	ErrMsgConvertToAppointment   = "failed to convert to appointment entity"

	// SQL operations
	DriverSQL = "SQLC"
	OpSelect  = "SELECT"
	OpInsert  = "INSERT"
	OpUpdate  = "UPDATE"
	OpDelete  = "DELETE"
	OpCount   = "COUNT"
)
