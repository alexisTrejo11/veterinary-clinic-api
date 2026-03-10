package handlers

import (
	"clinic-vet-api/internal/core/appointments"
	"clinic-vet-api/internal/core/customers"
	"clinic-vet-api/internal/core/employees"
	"clinic-vet-api/internal/core/pets"
	"clinic-vet-api/internal/infrastructure/http/handlers/dtos"
	"clinic-vet-api/internal/infrastructure/http/handlers/mappers"
	"clinic-vet-api/internal/shared/errors"
	"clinic-vet-api/internal/shared/http"
	"clinic-vet-api/internal/shared/page"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// ==========================================================================
// Base Controller
// ==========================================================================

// BaseAppointmentController provides base controller that could be consumed by other controllers
type BaseAppointmentController struct {
	commandService appointments.CommandService
	queryService   appointments.QueryService
	validate       *validator.Validate
	mapper         *mappers.AppointmentResponseMapper
}

func NewBaseAppointmentController(
	commandService appointments.CommandService,
	queryService appointments.QueryService,
	validate *validator.Validate,
) *BaseAppointmentController {
	return &BaseAppointmentController{
		commandService: commandService,
		queryService:   queryService,
		validate:       validate,
		mapper:         mappers.NewAppointmentResponseMapper(),
	}
}

type GetByIDExtraArgs struct {
	employeeID *uint
	customerID *uint
}

// Query Handlers

func (ctrl *BaseAppointmentController) GetAppointmentByID(c *gin.Context, args GetByIDExtraArgs) {
	idUint, err := http.ParseParamToUInt(c, "id")
	if err != nil {
		http.BadRequest(c, errors.RequestURLParamError(err, "id", c.Param("id")))
		return
	}

	appointmentId := appointments.NewAppointmentID(idUint)
	appt, err := ctrl.queryService.GetByID(c.Request.Context(), appointmentId)
	if err != nil {
		http.ApplicationError(c, err)
		return
	}

	appointmentResponse := ctrl.mapper.ToResponse(appt)
	http.Found(c, appointmentResponse, "Appointment")
}

func (ctrl *BaseAppointmentController) GetAppointmentsBySpecification(c *gin.Context) {
	searchSpecification, err := dtos.NewApptSearchRequestFromContext(c)
	if err != nil {
		http.BadRequest(c, errors.RequestURLQueryError(err, c.Request.URL.RawQuery))
		return
	}

	if err := ctrl.validate.Struct(searchSpecification); err != nil {
		http.BadRequest(c, errors.InvalidDataError(err))
		return
	}

	searchQuery, err := ctrl.mapper.RequestToSearchQuery(searchSpecification)
	if err != nil {
		http.ApplicationError(c, err)
		return
	}

	appointmentPage, err := ctrl.queryService.GetBySpecfication(c.Request.Context(), searchQuery)
	if err != nil {
		http.ApplicationError(c, err)
		return
	}

	responsePage := ctrl.mapper.ToResponsePage(appointmentPage)
	http.Paginated(c, responsePage, "Appointments")
}

func (ctrl *BaseAppointmentController) FindAppointmentsByDateRange(c *gin.Context) {
	var requestData dtos.AppointmentsByDateRangeRequest
	if err := c.ShouldBindQuery(&requestData); err != nil {
		http.BadRequest(c, errors.RequestBodyDataError(err))
		return
	}

	if err := ctrl.validate.Struct(&requestData); err != nil {
		http.BadRequest(c, errors.InvalidDataError(err))
		return
	}

	appointmentPage, err := ctrl.queryService.GetByDateRange(
		c.Request.Context(),
		requestData.StartTime(),
		requestData.EndTime(),
		requestData.ToPagination(),
	)
	if err != nil {
		http.ApplicationError(c, err)
		return
	}

	responsePage := ctrl.mapper.ToResponsePage(appointmentPage)
	http.Paginated(c, responsePage, "Appointments")
}

func (ctrl *BaseAppointmentController) FindAppointmentsByCustomer(
	c *gin.Context,
	customerID uint,
	//extraArgs customerQueryExtraArgs,
) {
	var pageParams page.PaginationRequest
	if err := http.ShouldBindPageParams(&pageParams, c, ctrl.validate); err != nil {
		http.BadRequest(c, errors.RequestURLQueryError(err, c.Request.URL.RawQuery))
		return
	}

	//query := query.NewFindApptsByCustomerIDQuery(pageParams, customerID, extraArgs.PetID, extraArgs.Status)
	customerIDObj := customers.NewCustomerID(customerID)
	pagintation := pageParams.ToPagination()
	appointmentPage, err := ctrl.queryService.GetByCustomerID(c.Request.Context(), customerIDObj, pagintation)
	if err != nil {
		http.ApplicationError(c, err)
		return
	}

	responsePage := ctrl.mapper.ToResponsePage(appointmentPage)
	http.Paginated(c, responsePage, "Appointments")
}

func (ctrl *BaseAppointmentController) GetAppointmentsByEmployee(c *gin.Context, employeeID uint) {
	var pageParams page.PaginationRequest
	if err := c.ShouldBindQuery(&pageParams); err != nil {
		http.BadRequest(c, errors.RequestURLQueryError(err, c.Request.URL.RawQuery))
		return
	}

	employeeIDObj := employees.NewEmployeeID(employeeID)
	pagination := pageParams.ToPagination()

	appointmentPage, err := ctrl.queryService.GetByEmployeeID(c.Request.Context(), employeeIDObj, pagination)
	if err != nil {
		http.ApplicationError(c, err)
		return
	}

	responsePage := ctrl.mapper.ToResponsePage(appointmentPage)
	http.Paginated(c, responsePage, "Appointments")
}

func (ctrl *BaseAppointmentController) GetAppointmentsByPet(c *gin.Context, petID uint) {
	var pageParams page.PaginationRequest
	if err := c.ShouldBindQuery(&pageParams); err != nil {
		http.BadRequest(c, errors.RequestURLQueryError(err, c.Request.URL.RawQuery))
		return
	}

	petIdObj := pets.NewPetID(petID)
	pagination := pageParams.ToPagination()
	appointmentPage, err := ctrl.queryService.GetByPetID(c.Request.Context(), petIdObj, pagination)
	if err != nil {
		http.ApplicationError(c, err)
		return
	}

	responsePage := ctrl.mapper.ToResponsePage(appointmentPage)
	http.Paginated(c, responsePage, "Appointments")
}

func (ctrl *BaseAppointmentController) GetAppointmentStats(c *gin.Context) {
	// TODO: Implement appointment stats logic or remove this function if not needed
}

// Command Operations

func (ctrl *BaseAppointmentController) CreateAppointment(c *gin.Context, employeeID *uint) {
	var request dtos.AppointmentCreateRequest
	if err := http.ShouldBindAndValidateBody(c, &request, ctrl.validate); err != nil {
		http.BadRequest(c, err)
		return
	}

	command, err := ctrl.mapper.RequestToCreateCommand(request, employeeID)
	if err != nil {
		http.ApplicationError(c, err)
		return
	}

	createdAppt, err := ctrl.commandService.CreateAppointment(c.Request.Context(), command)
	if err != nil {
		http.ApplicationError(c, err)
		return
	}

	http.Created(c, createdAppt.ID, "Appointment")
}

func (ctrl *BaseAppointmentController) UpdateAppointment(c *gin.Context, apptID uint) {
	var requestData dtos.AppointmentUpdateGeneralInfoRequest
	if err := http.ShouldBindAndValidateBody(c, &requestData, ctrl.validate); err != nil {
		http.BadRequest(c, err)
		return
	}

	updateCommand, err := ctrl.mapper.RequestToUpdateCommand(requestData, apptID)
	if err != nil {
		http.ApplicationError(c, err)
		return
	}

	err = ctrl.commandService.UpdateAppointmentGeneralInfo(c.Request.Context(), updateCommand)
	if err != nil {
		http.ApplicationError(c, err)
		return
	}

	http.Updated(c, nil, "Appointment")
}

func (ctrl *BaseAppointmentController) DeleteAppointment(c *gin.Context) {
	entityID, err := http.ParseParamToUInt(c, "id")
	if err != nil {
		http.BadRequest(c, errors.RequestURLParamError(err, "id", c.Param("id")))
		return
	}

	command := ctrl.mapper.ToDeleteCommand(entityID)
	err = ctrl.commandService.DeleteAppointment(c.Request.Context(), command)
	if err != nil {
		http.ApplicationError(c, err)
		return
	}

	http.Success(c, nil, "Appointment deleted successfully")
}

func (ctrl *BaseAppointmentController) RescheduleAppointment(c *gin.Context, employeeID *uint) {
	appointmentID, err := http.ParseParamToUInt(c, "id")
	if err != nil {
		http.BadRequest(c, errors.RequestURLParamError(err, "id", c.Param("id")))
		return
	}

	var requestData dtos.RescheduleAppointmentRequest
	if err := c.ShouldBindJSON(&requestData); err != nil {
		http.BadRequest(c, errors.RequestBodyDataError(err))
		return
	}

	command := ctrl.mapper.RequestToRescheduleCommand(requestData, appointmentID)
	err = ctrl.commandService.RescheduleAppointment(c.Request.Context(), command)
	if err != nil {
		http.ApplicationError(c, err)
		return
	}

	http.Success(c, nil, "Appointment rescheduled successfully")
}

func (ctrl *BaseAppointmentController) NotAttend(c *gin.Context, employeeID *uint) {
	appointmentID, err := http.ParseParamToUInt(c, "id")
	if err != nil {
		http.BadRequest(c, errors.RequestURLParamError(err, "id", c.Param("id")))
		return
	}

	command := ctrl.mapper.ToNotAttendCommand(appointmentID, employeeID)
	err = ctrl.commandService.MarkAppointmentAsNotAttend(c.Request.Context(), command)
	if err != nil {
		http.ApplicationError(c, err)
		return
	}

	http.Success(c, nil, "Appointment marked as not attended successfully")
}

func (ctrl *BaseAppointmentController) ConfirmAppointment(c *gin.Context, employeeID uint) {
	appointmentID, err := http.ParseParamToUInt(c, "id")
	if err != nil {
		http.BadRequest(c, errors.RequestURLParamError(err, "id", c.Param("id")))
		return
	}

	command := ctrl.mapper.ToConfirmCommand(appointmentID, employeeID)
	err = ctrl.commandService.ConfirmAppointment(c.Request.Context(), command)
	if err != nil {
		http.ApplicationError(c, err)
		return
	}

	http.Success(c, nil, "Appointment confirmed successfully")
}

func (ctrl *BaseAppointmentController) CompleteAppointment(c *gin.Context, employeeID *uint) {
	appointmentID, err := http.ParseParamToUInt(c, "id")
	if err != nil {
		http.BadRequest(c, errors.RequestURLParamError(err, "id", c.Param("id")))
		return
	}
	notes := c.Query("notes")

	command := ctrl.mapper.ToCompleteCommand(appointmentID, employeeID, notes)
	err = ctrl.commandService.CompleteAppointment(c.Request.Context(), command)
	if err != nil {
		http.ApplicationError(c, err)
		return
	}

	http.Success(c, nil, "Appointment completed successfully")
}

func (ctrl *BaseAppointmentController) CancelAppointment(c *gin.Context, employeeID *uint) {
	appointmentID, err := http.ParseParamToUInt(c, "id")
	if err != nil {
		http.BadRequest(c, errors.RequestURLParamError(err, "id", c.Param("id")))
		return
	}
	reason := c.Query("reason")

	command := ctrl.mapper.ToCancelCommand(appointmentID, employeeID, reason)
	err = ctrl.commandService.CancelAppointment(c.Request.Context(), command)
	if err != nil {
		http.ApplicationError(c, err)
		return
	}

	http.Success(c, nil, "Appointment cancelled successfully")
}

// ==========================================================================
// Manager Controller
// ==========================================================================

type AppointmentManagerHandler struct {
	commandService appointments.CommandService
	queryService   appointments.QueryService
	validator      *validator.Validate
	baseHandler    *BaseAppointmentController
}

func NewAppointmentManagerHandler(
	commandService appointments.CommandService,
	queryService appointments.QueryService,
	validate *validator.Validate,
) AppointmentManagerHandler {
	baseHandler := NewBaseAppointmentController(commandService, queryService, validate)
	return AppointmentManagerHandler{
		commandService: commandService,
		queryService:   queryService,
		validator:      validate,
		baseHandler:    baseHandler,
	}
}

func (ctrl AppointmentManagerHandler) GetBySpecfificationAppointments(c *gin.Context) {
	ctrl.baseHandler.GetAppointmentsBySpecification(c)
}

func (ctrl AppointmentManagerHandler) GetAppointmentByID(c *gin.Context) {
	ctrl.baseHandler.GetAppointmentByID(c, GetByIDExtraArgs{})
}
func (ctrl AppointmentManagerHandler) CreateAppointment(c *gin.Context) {
	employeeID := c.Param("employeeID")
	if employeeID == "" {
		http.BadRequest(c, errors.RequestURLQueryError(fmt.Errorf("employeeID query param is required"), c.Request.URL.RawQuery))
		return
	}

	employeeIDUint, err := strconv.ParseUint(employeeID, 10, 64)
	if err != nil {
		http.BadRequest(c, errors.RequestURLQueryError(err, c.Request.URL.RawQuery))
		return
	}

	employeeIDUintVal := uint(employeeIDUint)
	ctrl.baseHandler.CreateAppointment(c, &employeeIDUintVal)
}

func (ctrl AppointmentManagerHandler) UpdateAppointment(c *gin.Context) {
	appointmentID, err := http.ParseParamToUInt(c, "id")
	if err != nil {
		http.BadRequest(c, errors.RequestURLQueryError(err, c.Request.URL.RawQuery))
		return
	}
	ctrl.baseHandler.UpdateAppointment(c, appointmentID)
}

func (ctrl AppointmentManagerHandler) DeleteAppointment(c *gin.Context) {
	ctrl.baseHandler.DeleteAppointment(c)
}
