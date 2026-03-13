package repository

const (
	TableUsers         = "users"
	TableAppts         = "appointments"
	TablePets          = "pets"
	TableAddresses     = "addresses"
	TableCustomers     = "customers"
	TableEmployees     = "employees"
	TableNotifications = "notifications"
	TablePayments      = "payments"

	// Notification error messages
	ErrMsgFindNotificationByID    = "failed to find notification by ID"
	ErrMsgFindNotificationsByUser = "failed to find notifications by user"
	ErrMsgCreateNotification      = "failed to create notification"

	// Address error messages
	ErrMsgFindAddressByID        = "failed to find address by ID"
	ErrMsgFindAddressByUserID    = "failed to find address by user ID"
	ErrMsgCreateAddress          = "failed to create address"
	ErrMsgUpdateAddress          = "failed to update address"
	ErrMsgSoftDeleteAddress      = "failed to soft delete address"
	ErrMsgRestoreAddress         = "failed to restore address"
	ErrMsgConvertAddressToDomain = "failed to convert address to domain entity"

	// Customer error messages
	ErrMsgFindCustomerByID        = "failed to find customer by ID"
	ErrMsgFindCustomersBySpec     = "failed to find customers by specification"
	ErrMsgCreateCustomer          = "failed to create customer"
	ErrMsgUpdateCustomer          = "failed to update customer"
	ErrMsgConvertCustomerToDomain = "failed to convert customer to domain entity"

	// Employee error messages
	ErrMsgFindEmployeeByID        = "failed to find employee by ID"
	ErrMsgFindEmployeesBySpec     = "failed to find employees by specification"
	ErrMsgCreateEmployee          = "failed to create employee"
	ErrMsgUpdateEmployee          = "failed to update employee"
	ErrMsgConvertEmployeeToDomain = "failed to convert employee to domain entity"

	// Pet error messages
	ErrMsgConvertPetToDomain       = "failed to convert pet to domain entity"
	ErrMsgFindPetByIDAndCustomerID = "failed to find pet by ID and customer ID"
	ErrMsgFindPetsByCustomerID     = "failed to find pets by customer ID"
	ErrMsgFindPetsBySpecies        = "failed to find pets by species"
	ErrMsgFindPetsBySpecification  = "failed to find pets by specification"
	ErrMsgCreatePet                = "failed to create pet"
	ErrMsgUpdatePet                = "failed to update pet"
	ErrMsgSoftDeletePet            = "failed to soft delete pet"
	ErrMsgHardDeletePet            = "failed to hard delete pet"
	ErrMsgRestorePet               = "failed to restore pet"
	ErrMsgCountPetsByCustomerID    = "failed to count pets by customer ID"
	ErrMsgCountPetsBySpecies       = "failed to count pets by species"
	ErrMsgCheckPetExists           = "failed to check pet existence by ID"
	ErrMsgCountPets                = "failed to count pets"
	ErrMsgConvertToPet             = "failed to convert to pet entity"
	ErrMsgFindPetByID              = "failed to find pet by ID"

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

	// Payment error messages
	ErrMsgFindPaymentByID              = "failed to find payment by ID"
	ErrMsgFindPaymentByTransactionID   = "failed to find payment by transaction ID"
	ErrMsgFindPaymentByIDAndCustomerID = "failed to find payment by ID and customer ID"
	ErrMsgFindPaymentsBySpecification  = "failed to find payments by specification"
	ErrMsgCheckPaymentExists           = "failed to check payment existence by ID"
	ErrMsgCountPaymentsByStatus        = "failed to count payments by status"
	ErrMsgCountPaymentsByCustomerID    = "failed to count payments by customer ID"
	ErrMsgCountOverduePayments         = "failed to count overdue payments"
	ErrMsgTotalRevenueByDateRange      = "failed to total revenue by date range"
	ErrMsgCreatePayment                = "failed to create payment"
	ErrMsgUpdatePayment                = "failed to update payment"

	// SQL operations
	DriverSQL = "SQLC"
	OpSelect  = "SELECT"
	OpInsert  = "INSERT"
	OpUpdate  = "UPDATE"
	OpDelete  = "DELETE"
	OpCount   = "COUNT"
)
