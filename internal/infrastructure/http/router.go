package http

import (
	"clinic-vet-api/internal/infrastructure/http/handlers"
	"errors"
)

type APIRouter struct {
	appHandlers *AppHandlers
	config      *APIConfig
}

type AppHandlers struct {
	address      *handlers.AddressHandler
	auth         *handlers.AuthHandler
	user         *handlers.UserHandler
	pets         *handlers.PetHandler
	profile      *handlers.ProfileHandler
	appointment  *handlers.AppointmentHandler
	customer     *handlers.CustomerHandler
	employee     *handlers.EmployeeHandler
	payment      *handlers.PaymentHandler
	medical      *handlers.MedicalHandler
	notification *handlers.NotificationHandler
	home         *handlers.HomeHandler
}

func (r *APIRouter) Validate() error {
	if r.appHandlers.auth == nil {
		return errors.New("auth handler is required")
	}
	if r.config.AuthMiddleware == nil {
		return errors.New("auth middleware is required")
	}
	if r.config.Router == nil {
		return errors.New("router group is required")
	}
	return nil
}

func (a *AppHandlers) Validate() error {
	if a.address == nil {
		return errors.New("address handler is required")
	}
	if a.auth == nil {
		return errors.New("auth handler is required")
	}
	if a.user == nil {
		return errors.New("user handler is required")
	}
	if a.pets == nil {
		return errors.New("pets handler is required")
	}
	if a.profile == nil {
		return errors.New("profile handler is required")
	}
	if a.appointment == nil {
		return errors.New("appointment handler is required")
	}
	if a.customer == nil {
		return errors.New("customer handler is required")
	}
	if a.employee == nil {
		return errors.New("employee handler is required")
	}
	if a.payment == nil {
		return errors.New("payment handler is required")
	}
	return nil
}

func NewAPIRouter(
	appHandlers *AppHandlers,
	config *APIConfig,
) (*APIRouter, error) {
	router := &APIRouter{
		appHandlers: appHandlers,
		config:      config,
	}

	if err := router.Validate(); err != nil {
		return nil, err
	}
	return router, nil
}

// RegisterRoutes mounts all route groups (auth, profile, users, appointments, etc.).
// Call this after NewAPIRouter so that routes are registered.
func (r *APIRouter) RegisterRoutes() {
	r.authRoutes()
	r.profileRoutes()
	r.userRoutes()
	r.customerRoutes()
	r.employeeRoutes()
	r.appointmentRoutes()
	r.medicalRoutes()
	r.notificationRoutes()
}

// ------------------------------------------------------------
// Auth Routes
// ------------------------------------------------------------

func (r *APIRouter) authRoutes() {

	publicAuthRoutes := r.config.Router.Group("/api/v2/auth")
	{
		publicAuthRoutes.POST("/register", r.appHandlers.auth.Register)
		publicAuthRoutes.POST("/login", r.appHandlers.auth.Login)
		publicAuthRoutes.POST("/activate", r.appHandlers.auth.ActivateAccount)
	}

	authenticatedAuthRoutes := r.config.Router.Group("/api/v2/auth")
	authenticatedAuthRoutes.Use(r.config.AuthMiddleware.Authenticate())
	{
		authenticatedAuthRoutes.POST("/logout", r.appHandlers.auth.Logout)
		authenticatedAuthRoutes.POST("/logout-all", r.appHandlers.auth.LogoutAll)
		authenticatedAuthRoutes.POST("/refresh", r.appHandlers.auth.RefreshToken)
	}

	twoFactorAuthRoutes := r.config.Router.Group("/api/v2/auth/2fa")
	twoFactorAuthRoutes.Use(r.config.AuthMiddleware.Authenticate())
	{
		twoFactorAuthRoutes.POST("/verify", r.appHandlers.auth.VerifyTwoFactor)
		twoFactorAuthRoutes.POST("/enable", r.appHandlers.auth.EnableTwoFactor)
		twoFactorAuthRoutes.POST("/disable", r.appHandlers.auth.DisableTwoFactor)
	}

	resetPasswordAuthRoutes := r.config.Router.Group("/api/v2/auth/reset-password")
	resetPasswordAuthRoutes.Use(r.config.AuthMiddleware.Authenticate())
	{
		resetPasswordAuthRoutes.POST("/request", r.appHandlers.auth.RequestResetPassword)
		resetPasswordAuthRoutes.POST("/reset", r.appHandlers.auth.ResetPassword)
	}
}

// ------------------------------------------------------------
// Profile Routes
// ------------------------------------------------------------

func (r *APIRouter) profileRoutes() {
	profileRoutes := r.config.Router.Group("/api/v2/profile")
	profileRoutes.Use(r.config.AuthMiddleware.Authenticate())
	{
		profileRoutes.GET("/", r.appHandlers.profile.GetProfile)
		profileRoutes.PUT("/", r.appHandlers.profile.UpdateProfile)
	}
}

// ------------------------------------------------------------
// User Routes
// ------------------------------------------------------------

func (r *APIRouter) userRoutes() {
	userRoutes := r.config.Router.Group("/api/v2/users")
	userRoutes.Use(r.config.AuthMiddleware.Authenticate())
	userRoutes.Use(r.config.AuthMiddleware.RequireAnyRole("admin", "manager"))
	{
		userRoutes.GET("/", r.appHandlers.user.SearchUsers)
		userRoutes.GET("/:id", r.appHandlers.user.GetUserByID)
		userRoutes.POST("/", r.appHandlers.user.CreateUser)
		userRoutes.POST("/:id/restore", r.appHandlers.user.RestoreUser)
		userRoutes.POST("/:id/status", r.appHandlers.user.UpdateUserStatus)
		userRoutes.DELETE("/:id", r.appHandlers.user.DeleteUser)
	}
}

// ------------------------------------------------------------
// Appointment Routes
// ------------------------------------------------------------

func (r *APIRouter) appointmentRoutes() {
	// ----- Customer: "my" appointments (customer sees only their own) -----
	meAppointments := r.config.Router.Group("/api/v2/me/appointments")
	meAppointments.Use(r.config.AuthMiddleware.Authenticate())
	meAppointments.Use(r.config.AuthMiddleware.RequireAnyRole("customer"))
	{
		meAppointments.GET("/", r.appHandlers.appointment.GetMyAppointments)
		meAppointments.GET("/:id", r.appHandlers.appointment.GetMyAppointment)
		meAppointments.POST("/", r.appHandlers.appointment.RequestAppointment)
	}

	// ----- Employee: "my" assigned appointments (employee sees only their own) -----
	employeeAppointments := r.config.Router.Group("/api/v2/employees/appointments")
	employeeAppointments.Use(r.config.AuthMiddleware.Authenticate())
	employeeAppointments.Use(r.config.AuthMiddleware.RequireAnyRole("employee", "manager"))
	{
		employeeAppointments.GET("/", r.appHandlers.appointment.GetMyAppointmentsAsEmployee)
		employeeAppointments.GET("/:id", r.appHandlers.appointment.GetMyAppointmentAsEmployee)
		employeeAppointments.POST("/", r.appHandlers.appointment.CreateAppointmentAsEmployee)
		employeeAppointments.PUT("/:id", r.appHandlers.appointment.UpdateMyAppointment)
		employeeAppointments.POST("/:id/reschedule", r.appHandlers.appointment.RescheduleMyAppointment)
		employeeAppointments.POST("/:id/confirm", r.appHandlers.appointment.ConfirmMyAppointment)
		employeeAppointments.POST("/:id/complete", r.appHandlers.appointment.CompleteMyAppointment)
		employeeAppointments.POST("/:id/cancel", r.appHandlers.appointment.CancelMyAppointment)
		employeeAppointments.POST("/:id/not-attend", r.appHandlers.appointment.NotAttendMyAppointment)
	}

	// ----- Manager/Admin: all appointments (search, CRUD, by customer/employee/pet) -----
	managerAppointments := r.config.Router.Group("/api/v2/appointments")
	managerAppointments.Use(r.config.AuthMiddleware.Authenticate())
	managerAppointments.Use(r.config.AuthMiddleware.RequireAnyRole("admin", "manager"))
	{
		managerAppointments.GET("/", r.appHandlers.appointment.SearchAppointments)
		managerAppointments.GET("/:id", r.appHandlers.appointment.GetAppointmentByID)
		managerAppointments.POST("/", r.appHandlers.appointment.CreateAppointment)
		managerAppointments.PUT("/:id", r.appHandlers.appointment.UpdateAppointment)
		managerAppointments.DELETE("/:id", r.appHandlers.appointment.DeleteAppointment)
		managerAppointments.POST("/:id/reschedule", r.appHandlers.appointment.RescheduleAppointment)
		managerAppointments.POST("/:id/confirm", r.appHandlers.appointment.ConfirmAppointment)      // ?employee_id=1
		managerAppointments.POST("/:id/complete", r.appHandlers.appointment.CompleteAppointment)    // ?employee_id=1&notes=
		managerAppointments.POST("/:id/cancel", r.appHandlers.appointment.CancelAppointment)        // ?employee_id=1&reason=
		managerAppointments.POST("/:id/not-attend", r.appHandlers.appointment.NotAttendAppointment) // ?employee_id=1
	}

	// Manager: list appointments by customer (use :id to match /customers/:id)
	customersAppointments := r.config.Router.Group("/api/v2/customers/:id/appointments")
	customersAppointments.Use(r.config.AuthMiddleware.Authenticate())
	customersAppointments.Use(r.config.AuthMiddleware.RequireAnyRole("admin", "manager"))
	{
		customersAppointments.GET("/", r.appHandlers.appointment.GetAppointmentsByCustomerID)
	}

	// Manager: list appointments by employee (use :id to match /employees/:id)
	employeesAppointments := r.config.Router.Group("/api/v2/employees/:id/appointments")
	employeesAppointments.Use(r.config.AuthMiddleware.Authenticate())
	employeesAppointments.Use(r.config.AuthMiddleware.RequireAnyRole("admin", "manager"))
	{
		employeesAppointments.GET("/", r.appHandlers.appointment.GetAppointmentsByEmployeeID)
	}

	// Manager: list appointments by pet
	petsAppointments := r.config.Router.Group("/api/v2/pets/:id/appointments")
	petsAppointments.Use(r.config.AuthMiddleware.Authenticate())
	petsAppointments.Use(r.config.AuthMiddleware.RequireAnyRole("admin", "manager"))
	{
		petsAppointments.GET("/", r.appHandlers.appointment.GetAppointmentsByPetID)
	}
}

// ------------------------------------------------------------
// Customer Routes
// ------------------------------------------------------------

func (r *APIRouter) customerRoutes() {
	customerRoutes := r.config.Router.Group("/api/v2/customers")
	customerRoutes.Use(r.config.AuthMiddleware.Authenticate())
	customerRoutes.Use(r.config.AuthMiddleware.RequireAnyRole("admin", "manager"))
	{
		customerRoutes.GET("/", r.appHandlers.customer.SearchCustomers)
		customerRoutes.GET("/:id", r.appHandlers.customer.GetCustomerByID)
		customerRoutes.POST("/", r.appHandlers.customer.CreateCustomer)
		customerRoutes.PUT("/:id", r.appHandlers.customer.UpdateCustomer)
		customerRoutes.DELETE("/:id", r.appHandlers.customer.DeleteCustomer)
		customerRoutes.POST("/:id/restore", r.appHandlers.customer.RestoreCustomer)
	}
}

// ------------------------------------------------------------
// Employee Routes
// ------------------------------------------------------------

func (r *APIRouter) employeeRoutes() {
	employeeRoutes := r.config.Router.Group("/api/v2/employees")
	employeeRoutes.Use(r.config.AuthMiddleware.Authenticate())
	employeeRoutes.Use(r.config.AuthMiddleware.RequireAnyRole("admin", "manager"))
	{
		employeeRoutes.GET("/", r.appHandlers.employee.SearchEmployees)
		employeeRoutes.GET("/:id", r.appHandlers.employee.GetEmployeeByID)
		employeeRoutes.POST("/", r.appHandlers.employee.CreateEmployee)
		employeeRoutes.PUT("/:id", r.appHandlers.employee.UpdateEmployee)
		employeeRoutes.DELETE("/:id", r.appHandlers.employee.DeleteEmployee)
		employeeRoutes.POST("/:id/restore", r.appHandlers.employee.RestoreEmployee)
	}
}

// ------------------------------------------------------------
// Medical Routes
// ------------------------------------------------------------

func (r *APIRouter) medicalRoutes() {
	if r.appHandlers.medical == nil {
		return
	}
	m := r.appHandlers.medical

	// Customer: read-only (my sessions, my pets’ data)
	meMedical := r.config.Router.Group("/api/v2/me/medical")
	meMedical.Use(r.config.AuthMiddleware.Authenticate())
	meMedical.Use(r.config.AuthMiddleware.RequireAnyRole("customer"))
	{
		meMedical.GET("/sessions", m.GetMySessions)
		meMedical.GET("/sessions/:id/full", m.GetMySessionFull)
		meMedical.GET("/sessions/:id", m.GetMySessionByID)
		meMedical.GET("/pets/:pet_id/sessions", m.GetMyPetSessions)
		meMedical.GET("/pets/:pet_id/vaccination-summary", m.GetMyPetVaccinationSummary)
		meMedical.GET("/pets/:pet_id/vaccination-history", m.GetMyPetVaccinationHistory)
		meMedical.GET("/pets/:pet_id/prescriptions/active", m.GetMyPetActivePrescriptions)
	}

	// Staff (employee + manager): read and write
	medical := r.config.Router.Group("/api/v2/medical")
	medical.Use(r.config.AuthMiddleware.Authenticate())
	medical.Use(r.config.AuthMiddleware.RequireAnyRole("employee", "manager", "admin"))
	{
		// Sessions (more specific routes first)
		medical.GET("/sessions", m.GetSessionsBySpecification)
		medical.GET("/sessions/stats", m.GetSessionStats)
		medical.GET("/sessions/:id/full", m.GetSessionFull)
		medical.GET("/sessions/:id", m.GetSessionByID)
		medical.POST("/sessions", m.CreateSession)
		medical.PUT("/sessions/:id", m.UpdateSession)
		medical.DELETE("/sessions/:id", m.SoftDeleteSession)
		medical.DELETE("/sessions/:id/hard", m.HardDeleteSession)
		medical.POST("/sessions/:id/restore", m.RestoreSession)
		medical.GET("/customers/:customer_id/sessions", m.GetSessionsByCustomer)
		medical.GET("/pets/:pet_id/sessions", m.GetSessionsByPet)

		// Session extensions
		medical.GET("/sessions/:id/vaccinations", m.GetVaccinationsBySession)
		medical.POST("/vaccinations", m.AddVaccination)
		medical.PUT("/vaccinations/:vaccination_id", m.UpdateVaccination)
		medical.DELETE("/vaccinations/:vaccination_id", m.RemoveVaccination)

		medical.GET("/sessions/:id/surgeries", m.GetSurgeriesBySession)
		medical.POST("/surgeries", m.AddSurgery)
		medical.PUT("/surgeries/:surgery_id", m.UpdateSurgery)
		medical.DELETE("/surgeries/:surgery_id", m.RemoveSurgery)

		medical.GET("/sessions/:id/prescriptions", m.GetPrescriptionsBySession)
		medical.POST("/prescriptions", m.AddPrescription)
		medical.PUT("/prescriptions/:prescription_id", m.UpdatePrescription)
		medical.DELETE("/prescriptions/:prescription_id", m.RemovePrescription)
		medical.GET("/pets/:pet_id/prescriptions/active", m.GetActivePrescriptionsByPet)

		medical.GET("/sessions/:id/attachments", m.GetAttachmentsBySession)
		medical.POST("/attachments", m.AddAttachment)
		medical.DELETE("/attachments/:attachment_id", m.RemoveAttachment)

		medical.GET("/sessions/:id/services", m.GetServicesBySession)
		medical.POST("/session-services", m.AddSessionService)
		medical.DELETE("/session-services/:session_service_id", m.RemoveSessionService)

		medical.GET("/vaccination-history", m.GetVaccinationHistory)
		medical.GET("/pets/:pet_id/vaccination-summary", m.GetPetVaccinationSummary)

		// Catalogs (read + write)
		medical.GET("/catalogs/vaccines", m.ListVaccines)
		medical.GET("/catalogs/vaccines/species/:species", m.ListVaccinesBySpecies)
		medical.GET("/catalogs/vaccines/:id", m.GetVaccineByID)
		medical.POST("/catalogs/vaccines", m.CreateVaccineCatalog)
		medical.DELETE("/catalogs/vaccines/:id", m.DeactivateVaccine)

		medical.GET("/catalogs/medications", m.ListMedications)
		medical.GET("/catalogs/medications/search", m.SearchMedications)
		medical.GET("/catalogs/medications/:id", m.GetMedicationByID)
		medical.POST("/catalogs/medications", m.CreateMedicationCatalog)
		medical.DELETE("/catalogs/medications/:id", m.DeactivateMedication)

		medical.GET("/catalogs/services", m.ListServices)
		medical.GET("/catalogs/services/category/:category", m.ListServicesByCategory)
		medical.GET("/catalogs/services/:id", m.GetServiceByID)
		medical.POST("/catalogs/services", m.CreateServiceCatalog)
		medical.DELETE("/catalogs/services/:id", m.DeactivateServiceCatalog)
	}
}

// ------------------------------------------------------------
// Notification Routes
// ------------------------------------------------------------

func (r *APIRouter) notificationRoutes() {
	if r.appHandlers.notification == nil {
		return
	}
	n := r.appHandlers.notification

	// Customer: read-only (my notifications)
	meNotif := r.config.Router.Group("/api/v2/me/notifications")
	meNotif.Use(r.config.AuthMiddleware.Authenticate())
	meNotif.Use(r.config.AuthMiddleware.RequireAnyRole("customer", "employee", "manager", "admin"))
	{
		meNotif.GET("", n.GetMyNotifications)
		meNotif.GET("/:id", n.GetMyNotificationByID)
	}

	// Staff: monitoring (by type, channel, id, summary) + manual send
	notif := r.config.Router.Group("/api/v2/notifications")
	notif.Use(r.config.AuthMiddleware.Authenticate())
	notif.Use(r.config.AuthMiddleware.RequireAnyRole("employee", "manager", "admin"))
	{
		notif.GET("/type/:type", n.GetNotificationsByType)
		notif.GET("/channel/:channel", n.GetNotificationsByChannel)
		notif.GET("/summary", n.GetNotificationSummary)
		notif.GET("/:id", n.GetNotificationByID)
		notif.POST("", n.SendNotification)
	}
}

// ------------------------------------------------------------
// Home Routes
// ------------------------------------------------------------

func (r *APIRouter) homeRoutes() {
	homeRoutes := r.config.Router.Group("")
	homeRoutes.GET("/health", r.appHandlers.home.HealthCheck)
}
