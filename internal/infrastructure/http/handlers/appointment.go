package handlers

import (
	"errors"

	"clinic-vet-api/internal/core/appointments"
	"clinic-vet-api/internal/core/customers"
	"clinic-vet-api/internal/core/employees"
	"clinic-vet-api/internal/core/pets"
	"clinic-vet-api/internal/infrastructure/http/handlers/dtos"
	"clinic-vet-api/internal/infrastructure/http/handlers/mappers"
	"clinic-vet-api/internal/shared/http"
	"clinic-vet-api/internal/shared/page"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// AppointmentHandler handles appointment HTTP with customer-, employee-, and manager-scoped endpoints.
type AppointmentHandler struct {
	commandService     appointments.CommandService
	queryService       appointments.QueryService
	validator          *validator.Validate
	mapper             *mappers.AppointmentResponseMapper
	customerIDResolver CustomerIDResolver
	employeeIDResolver EmployeeIDResolver
}

func NewAppointmentHandler(
	commandService appointments.CommandService,
	queryService appointments.QueryService,
	validator *validator.Validate,
	mapper *mappers.AppointmentResponseMapper,
	customerIDResolver CustomerIDResolver,
	employeeIDResolver EmployeeIDResolver,
) *AppointmentHandler {
	if mapper == nil {
		mapper = mappers.NewAppointmentResponseMapper()
	}
	return &AppointmentHandler{
		commandService:     commandService,
		queryService:       queryService,
		validator:          validator,
		mapper:             mapper,
		customerIDResolver: customerIDResolver,
		employeeIDResolver: employeeIDResolver,
	}
}

func parseAppointmentIDFromParam(c *gin.Context) (appointments.AppointmentID, error) {
	id, err := http.ParseParamToUInt(c, "id")
	if err != nil {
		return appointments.AppointmentID{}, err
	}
	return appointments.NewAppointmentID(id), nil
}

// ------------------------------------------------------------
// Internal handlers
// ------------------------------------------------------------

func (h *AppointmentHandler) getAppointmentByIDInternal(
	ctx *gin.Context,
	apptID appointments.AppointmentID,
	optCustomerID *uint,
	optEmployeeID *uint,
) (any, error) {
	if optCustomerID != nil {
		cid := customers.NewCustomerID(*optCustomerID)
		return h.queryService.GetByIDAndCustomerID(ctx.Request.Context(), apptID, cid)
	}
	if optEmployeeID != nil {
		eid := employees.NewEmployeeID(*optEmployeeID)
		return h.queryService.GetByIDAndEmployeeID(ctx.Request.Context(), apptID, eid)
	}
	return h.queryService.GetByID(ctx.Request.Context(), apptID)
}

func (h *AppointmentHandler) getAppointmentsByCustomerIDInternal(
	ctx *gin.Context,
	getCustomerID CustomerIDProvider,
	pagination page.Pagination,
) (any, error) {
	customerID, err := getCustomerID(ctx)
	if err != nil {
		return nil, err
	}
	cid := customers.NewCustomerID(customerID)
	return h.queryService.GetByCustomerID(ctx.Request.Context(), cid, pagination)
}

func (h *AppointmentHandler) getAppointmentsByEmployeeIDInternal(
	ctx *gin.Context,
	getEmployeeID EmployeeIDProvider,
	pagination page.Pagination,
) (any, error) {
	employeeID, err := getEmployeeID(ctx)
	if err != nil {
		return nil, err
	}
	eid := employees.NewEmployeeID(employeeID)
	return h.queryService.GetByEmployeeID(ctx.Request.Context(), eid, pagination)
}

func (h *AppointmentHandler) getAppointmentsBySpecificationInternal(
	ctx *gin.Context,
	searchReq dtos.AppointmentSearchRequest,
) (any, error) {
	query, err := h.mapper.RequestToSearchQuery(searchReq)
	if err != nil {
		return nil, err
	}
	return h.queryService.GetBySpecfication(ctx.Request.Context(), query)
}

func (h *AppointmentHandler) requestAppointmentInternal(
	ctx *gin.Context,
	req dtos.AppointmentRequestByCustomerRequest,
	getCustomerID CustomerIDProvider,
) (any, error) {
	customerID, err := getCustomerID(ctx)
	if err != nil {
		return nil, err
	}
	cid := customers.NewCustomerID(customerID)
	command, err := h.mapper.RequestToRequestByCustomerCommand(req, cid)
	if err != nil {
		return nil, err
	}
	err = h.commandService.RequestAppointment(ctx.Request.Context(), command)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (h *AppointmentHandler) createAppointmentInternal(
	ctx *gin.Context,
	req dtos.AppointmentCreateRequest,
	optEmployeeID *uint,
) (any, error) {
	command, err := h.mapper.RequestToCreateCommand(req, optEmployeeID)
	if err != nil {
		return nil, err
	}
	created, err := h.commandService.CreateAppointment(ctx.Request.Context(), command)
	if err != nil {
		return nil, err
	}
	return created.ID.Value(), nil
}

func (h *AppointmentHandler) updateAppointmentInternal(
	ctx *gin.Context,
	req dtos.AppointmentUpdateGeneralInfoRequest,
	apptID uint,
) (any, error) {
	cmd, err := h.mapper.RequestToUpdateCommand(req, apptID)
	if err != nil {
		return nil, err
	}
	err = h.commandService.UpdateAppointmentGeneralInfo(ctx.Request.Context(), cmd)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (h *AppointmentHandler) deleteAppointmentInternal(ctx *gin.Context, apptID uint, isHard bool) (any, error) {
	command := h.mapper.ToDeleteCommand(apptID, isHard)
	err := h.commandService.DeleteAppointment(ctx.Request.Context(), command)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (h *AppointmentHandler) rescheduleInternal(
	ctx *gin.Context,
	req dtos.RescheduleAppointmentRequest,
	apptID uint,
) (any, error) {
	command := h.mapper.RequestToRescheduleCommand(req, apptID)
	err := h.commandService.RescheduleAppointment(ctx.Request.Context(), command)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (h *AppointmentHandler) confirmInternal(ctx *gin.Context, apptID uint, employeeID uint) (any, error) {
	command := h.mapper.ToConfirmCommand(apptID, employeeID)
	err := h.commandService.ConfirmAppointment(ctx.Request.Context(), command)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (h *AppointmentHandler) completeInternal(ctx *gin.Context, apptID uint, optEmployeeID *uint, notes string) (any, error) {
	command := h.mapper.ToCompleteCommand(apptID, optEmployeeID, notes)
	err := h.commandService.CompleteAppointment(ctx.Request.Context(), command)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (h *AppointmentHandler) cancelInternal(ctx *gin.Context, apptID uint, optEmployeeID *uint, reason string) (any, error) {
	command := h.mapper.ToCancelCommand(apptID, optEmployeeID, reason)
	err := h.commandService.CancelAppointment(ctx.Request.Context(), command)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (h *AppointmentHandler) notAttendInternal(ctx *gin.Context, apptID uint, optEmployeeID *uint) (any, error) {
	command := h.mapper.ToNotAttendCommand(apptID, optEmployeeID)
	err := h.commandService.MarkAppointmentAsNotAttend(ctx.Request.Context(), command)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (h *AppointmentHandler) getMyCustomerID(c *gin.Context) (uint, error) {
	if h.customerIDResolver == nil {
		return 0, errors.New("customer resolver not configured")
	}
	return CustomerIDFromContext(c, h.customerIDResolver)
}

func (h *AppointmentHandler) getMyEmployeeID(c *gin.Context) (uint, error) {
	if h.employeeIDResolver == nil {
		return 0, errors.New("employee resolver not configured")
	}
	return EmployeeIDFromContext(c, h.employeeIDResolver)
}

// ------------------------------------------------------------
// Customer handlers (only their appointments)
// ------------------------------------------------------------

func (h *AppointmentHandler) GetMyAppointment(c *gin.Context) {
	logic := func(ctx *gin.Context) (any, error) {
		customerID, err := h.getMyCustomerID(ctx)
		if err != nil {
			return nil, err
		}
		apptID, err := parseAppointmentIDFromParam(ctx)
		if err != nil {
			return nil, err
		}
		return h.getAppointmentByIDInternal(ctx, apptID, &customerID, nil)
	}
	http.HandleGetRequest(h.validator, "Appointment", logic)(c)
}

// GetMyAppointments godoc
// @Summary      Get my appointments (customer)
// @Description  Returns paginated appointments for the authenticated customer. Requires customer role.
// @Tags         appointments
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        page_number  query  int  false  "Page number"  default(1)
// @Param        page_size    query  int  false  "Page size"    default(10)
// @Success      200          {object}  http.APIResponse  "Paginated appointments"
// @Failure      401          {object}  http.APIResponse  "Unauthorized"
// @Failure      500          {object}  http.APIResponse  "Internal server error"
// @Router       /me/appointments [get]
func (h *AppointmentHandler) GetMyAppointments(c *gin.Context) {
	logic := func(ctx *gin.Context) (any, error) {
		var pageParams page.PaginationRequest
		if err := http.ShouldBindPageParams(&pageParams, ctx, h.validator); err != nil {
			return nil, err
		}
		getCustomerID := func(c *gin.Context) (uint, error) { return h.getMyCustomerID(c) }
		return h.getAppointmentsByCustomerIDInternal(ctx, getCustomerID, pageParams.ToPagination())
	}
	http.HandleRequestNoBodyWithResponder(h.validator, logic, func(c *gin.Context, res any) {
		if res == nil {
			http.Success(c, nil, "No appointments found")
			return
		}
		p := res.(page.Page[appointments.Appointment])
		responsePage := h.mapper.ToResponsePage(p)
		http.Paginated(c, responsePage, "Appointments")
	})(c)
}

// RequestAppointment godoc
// @Summary      Request appointment (customer)
// @Description  Submits an appointment request as a customer. Requires customer role.
// @Tags         appointments
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body   body      dtos.AppointmentRequestByCustomerRequest  true  "Appointment request"
// @Success      200    {object}  http.APIResponse                          "Request submitted"
// @Failure      400    {object}  http.APIResponse                          "Validation error"
// @Failure      401    {object}  http.APIResponse                          "Unauthorized"
// @Failure      500    {object}  http.APIResponse                          "Internal server error"
// @Router       /me/appointments [post]
func (h *AppointmentHandler) RequestAppointment(c *gin.Context) {
	logic := func(ctx *gin.Context, req dtos.AppointmentRequestByCustomerRequest) (any, error) {
		getCustomerID := func(c *gin.Context) (uint, error) { return h.getMyCustomerID(c) }
		return h.requestAppointmentInternal(ctx, req, getCustomerID)
	}
	http.HandleRequestWithResponder(h.validator, logic, func(c *gin.Context, _ any) {
		http.Success(c, nil, "Appointment request submitted successfully")
	})(c)
}

// ------------------------------------------------------------
// Employee handlers (only their assigned appointments)
// ------------------------------------------------------------

func (h *AppointmentHandler) GetMyAppointmentAsEmployee(c *gin.Context) {
	logic := func(ctx *gin.Context) (any, error) {
		employeeID, err := h.getMyEmployeeID(ctx)
		if err != nil {
			return nil, err
		}
		apptID, err := parseAppointmentIDFromParam(ctx)
		if err != nil {
			return nil, err
		}
		return h.getAppointmentByIDInternal(ctx, apptID, nil, &employeeID)
	}
	http.HandleGetRequest(h.validator, "Appointment", logic)(c)
}

func (h *AppointmentHandler) GetMyAppointmentsAsEmployee(c *gin.Context) {
	logic := func(ctx *gin.Context) (any, error) {
		var pageParams page.PaginationRequest
		if err := http.ShouldBindPageParams(&pageParams, ctx, h.validator); err != nil {
			return nil, err
		}
		getEmployeeID := func(c *gin.Context) (uint, error) { return h.getMyEmployeeID(c) }
		return h.getAppointmentsByEmployeeIDInternal(ctx, getEmployeeID, pageParams.ToPagination())
	}
	http.HandleRequestNoBodyWithResponder(h.validator, logic, func(c *gin.Context, res any) {
		if res == nil {
			http.Success(c, nil, "No appointments found")
			return
		}
		p := res.(page.Page[appointments.Appointment])
		responsePage := h.mapper.ToResponsePage(p)
		http.Paginated(c, responsePage, "Appointments")
	})(c)
}

func (h *AppointmentHandler) CreateAppointmentAsEmployee(c *gin.Context) {
	logic := func(ctx *gin.Context, req dtos.AppointmentCreateRequest) (any, error) {
		employeeID, err := h.getMyEmployeeID(ctx)
		if err != nil {
			return nil, err
		}
		return h.createAppointmentInternal(ctx, req, &employeeID)
	}
	http.HandleCreateRequest(h.validator, "Appointment", logic)(c)
}

func (h *AppointmentHandler) UpdateMyAppointment(c *gin.Context) {
	logic := func(ctx *gin.Context, req dtos.AppointmentUpdateGeneralInfoRequest) (any, error) {
		apptID, err := parseAppointmentIDFromParam(ctx)
		if err != nil {
			return nil, err
		}
		return h.updateAppointmentInternal(ctx, req, apptID.Value())
	}
	http.HandleUpdateRequest(h.validator, "Appointment", logic)(c)
}

func (h *AppointmentHandler) RescheduleMyAppointment(c *gin.Context) {
	logic := func(ctx *gin.Context, req dtos.RescheduleAppointmentRequest) (any, error) {
		apptID, err := parseAppointmentIDFromParam(ctx)
		if err != nil {
			return nil, err
		}
		return h.rescheduleInternal(ctx, req, apptID.Value())
	}
	http.HandleRequestWithResponder(h.validator, logic, func(c *gin.Context, _ any) {
		http.Success(c, nil, "Appointment rescheduled successfully")
	})(c)
}

func (h *AppointmentHandler) ConfirmMyAppointment(c *gin.Context) {
	logic := func(ctx *gin.Context) (any, error) {
		employeeID, err := h.getMyEmployeeID(ctx)
		if err != nil {
			return nil, err
		}
		apptID, err := parseAppointmentIDFromParam(ctx)
		if err != nil {
			return nil, err
		}
		return h.confirmInternal(ctx, apptID.Value(), employeeID)
	}
	http.HandleRequestNoBodyWithResponder(h.validator, logic, func(c *gin.Context, _ any) {
		http.Success(c, nil, "Appointment confirmed successfully")
	})(c)
}

func (h *AppointmentHandler) CompleteMyAppointment(c *gin.Context) {
	logic := func(ctx *gin.Context) (any, error) {
		employeeID, err := h.getMyEmployeeID(ctx)
		if err != nil {
			return nil, err
		}
		apptID, err := parseAppointmentIDFromParam(ctx)
		if err != nil {
			return nil, err
		}
		notes := ctx.Query("notes")
		return h.completeInternal(ctx, apptID.Value(), &employeeID, notes)
	}
	http.HandleRequestNoBodyWithResponder(h.validator, logic, func(c *gin.Context, _ any) {
		http.Success(c, nil, "Appointment completed successfully")
	})(c)
}

func (h *AppointmentHandler) CancelMyAppointment(c *gin.Context) {
	logic := func(ctx *gin.Context) (any, error) {
		employeeID, err := h.getMyEmployeeID(ctx)
		if err != nil {
			return nil, err
		}
		apptID, err := parseAppointmentIDFromParam(ctx)
		if err != nil {
			return nil, err
		}
		reason := ctx.Query("reason")
		return h.cancelInternal(ctx, apptID.Value(), &employeeID, reason)
	}
	http.HandleRequestNoBodyWithResponder(h.validator, logic, func(c *gin.Context, _ any) {
		http.Success(c, nil, "Appointment cancelled successfully")
	})(c)
}

func (h *AppointmentHandler) NotAttendMyAppointment(c *gin.Context) {
	logic := func(ctx *gin.Context) (any, error) {
		employeeID, err := h.getMyEmployeeID(ctx)
		if err != nil {
			return nil, err
		}
		apptID, err := parseAppointmentIDFromParam(ctx)
		if err != nil {
			return nil, err
		}
		return h.notAttendInternal(ctx, apptID.Value(), &employeeID)
	}
	http.HandleRequestNoBodyWithResponder(h.validator, logic, func(c *gin.Context, _ any) {
		http.Success(c, nil, "Appointment marked as not attended successfully")
	})(c)
}

// ------------------------------------------------------------
// Manager/Admin handlers (all appointments)
// ------------------------------------------------------------

func (h *AppointmentHandler) GetAppointmentByID(c *gin.Context) {
	logic := func(ctx *gin.Context) (any, error) {
		apptID, err := parseAppointmentIDFromParam(ctx)
		if err != nil {
			return nil, err
		}
		return h.getAppointmentByIDInternal(ctx, apptID, nil, nil)
	}
	http.HandleGetRequest(h.validator, "Appointment", logic)(c)
}

// SearchAppointments godoc
// @Summary      Search appointments (manager/admin)
// @Description  Paginated search of all appointments with filters. Requires admin or manager role.
// @Tags         appointments
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body   body      dtos.AppointmentSearchRequest  true  "Search criteria and pagination"
// @Success      200    {object}  http.APIResponse               "Paginated appointments"
// @Failure      400    {object}  http.APIResponse               "Validation error"
// @Failure      401    {object}  http.APIResponse               "Unauthorized"
// @Failure      500    {object}  http.APIResponse               "Internal server error"
// @Router       /appointments [get]
func (h *AppointmentHandler) SearchAppointments(c *gin.Context) {
	searchReq, err := dtos.NewApptSearchRequestFromContext(c)
	if err != nil {
		http.BadRequest(c, err)
		return
	}
	if err := h.validator.Struct(searchReq); err != nil {
		http.BadRequest(c, err)
		return
	}
	logic := func(ctx *gin.Context) (any, error) {
		return h.getAppointmentsBySpecificationInternal(ctx, searchReq)
	}
	http.HandleRequestNoBodyWithResponder(h.validator, logic, func(c *gin.Context, res any) {
		if res == nil {
			http.Success(c, nil, "No appointments found")
			return
		}
		p := res.(page.Page[appointments.Appointment])
		responsePage := h.mapper.ToResponsePage(p)
		http.Paginated(c, responsePage, "Appointments")
	})(c)
}

func (h *AppointmentHandler) GetAppointmentsByCustomerID(c *gin.Context) {
	logic := func(ctx *gin.Context) (any, error) {
		var pageParams page.PaginationRequest
		if err := http.ShouldBindPageParams(&pageParams, ctx, h.validator); err != nil {
			return nil, err
		}
		return h.getAppointmentsByCustomerIDInternal(ctx, CustomerIDFromIDParam, pageParams.ToPagination())
	}
	http.HandleRequestNoBodyWithResponder(h.validator, logic, func(c *gin.Context, res any) {
		if res == nil {
			http.Success(c, nil, "No appointments found")
			return
		}
		p := res.(page.Page[appointments.Appointment])
		responsePage := h.mapper.ToResponsePage(p)
		http.Paginated(c, responsePage, "Appointments")
	})(c)
}

func (h *AppointmentHandler) GetAppointmentsByEmployeeID(c *gin.Context) {
	logic := func(ctx *gin.Context) (any, error) {
		var pageParams page.PaginationRequest
		if err := http.ShouldBindPageParams(&pageParams, ctx, h.validator); err != nil {
			return nil, err
		}
		return h.getAppointmentsByEmployeeIDInternal(ctx, EmployeeIDFromIDParam, pageParams.ToPagination())
	}
	http.HandleRequestNoBodyWithResponder(h.validator, logic, func(c *gin.Context, res any) {
		if res == nil {
			http.Success(c, nil, "No appointments found")
			return
		}
		p := res.(page.Page[appointments.Appointment])
		responsePage := h.mapper.ToResponsePage(p)
		http.Paginated(c, responsePage, "Appointments")
	})(c)
}

func (h *AppointmentHandler) GetAppointmentsByPetID(c *gin.Context) {
	logic := func(ctx *gin.Context) (any, error) {
		petIDRaw, err := http.ParseParamToUInt(ctx, "id")
		if err != nil {
			return nil, err
		}
		var pageParams page.PaginationRequest
		if err := http.ShouldBindPageParams(&pageParams, c, h.validator); err != nil {
			return nil, err
		}
		petID := pets.NewPetID(petIDRaw)
		page, err := h.queryService.GetByPetID(c.Request.Context(), petID, pageParams.ToPagination())
		if err != nil {
			return nil, err
		}
		return page, nil
	}
	http.HandleRequestNoBodyWithResponder(h.validator, logic, func(c *gin.Context, res any) {
		if res == nil {
			http.Success(c, nil, "No appointments found")
			return
		}
		p := res.(page.Page[appointments.Appointment])
		responsePage := h.mapper.ToResponsePage(p)
		http.Paginated(c, responsePage, "Appointments")
	})(c)
}

func (h *AppointmentHandler) CreateAppointment(c *gin.Context) {
	logic := func(ctx *gin.Context, req dtos.AppointmentCreateRequest) (any, error) {
		optEmployeeID, _ := OptionalEmployeeIDFromQuery(ctx)
		return h.createAppointmentInternal(ctx, req, optEmployeeID)
	}
	http.HandleCreateRequest(h.validator, "Appointment", logic)(c)
}

func (h *AppointmentHandler) UpdateAppointment(c *gin.Context) {
	logic := func(ctx *gin.Context, req dtos.AppointmentUpdateGeneralInfoRequest) (any, error) {
		apptID, err := parseAppointmentIDFromParam(ctx)
		if err != nil {
			return nil, err
		}
		return h.updateAppointmentInternal(ctx, req, apptID.Value())
	}
	http.HandleUpdateRequest(h.validator, "Appointment", logic)(c)
}

func (h *AppointmentHandler) DeleteAppointment(c *gin.Context) {
	logic := func(ctx *gin.Context) (any, error) {
		apptID, err := parseAppointmentIDFromParam(ctx)
		if err != nil {
			return nil, err
		}
		isHard := ctx.Query("hard") == "true"
		return h.deleteAppointmentInternal(ctx, apptID.Value(), isHard)
	}
	http.HandleDeleteRequest(h.validator, "Appointment", logic)(c)
}

func (h *AppointmentHandler) RescheduleAppointment(c *gin.Context) {
	logic := func(ctx *gin.Context, req dtos.RescheduleAppointmentRequest) (any, error) {
		apptID, err := parseAppointmentIDFromParam(ctx)
		if err != nil {
			return nil, err
		}
		return h.rescheduleInternal(ctx, req, apptID.Value())
	}
	http.HandleRequestWithResponder(h.validator, logic, func(c *gin.Context, _ any) {
		http.Success(c, nil, "Appointment rescheduled successfully")
	})(c)
}

func (h *AppointmentHandler) ConfirmAppointment(c *gin.Context) {
	logic := func(ctx *gin.Context) (any, error) {
		apptID, err := parseAppointmentIDFromParam(ctx)
		if err != nil {
			return nil, err
		}
		// Manager passes employee_id as query: ?employee_id=1
		employeeID, err := EmployeeIDFromQuery(ctx)
		if err != nil {
			return nil, err
		}
		return h.confirmInternal(ctx, apptID.Value(), employeeID)
	}
	http.HandleRequestNoBodyWithResponder(h.validator, logic, func(c *gin.Context, _ any) {
		http.Success(c, nil, "Appointment confirmed successfully")
	})(c)
}

func (h *AppointmentHandler) CompleteAppointment(c *gin.Context) {
	logic := func(ctx *gin.Context) (any, error) {
		apptID, err := parseAppointmentIDFromParam(ctx)
		if err != nil {
			return nil, err
		}
		optEmployeeID, _ := OptionalEmployeeIDFromQuery(ctx)
		notes := ctx.Query("notes")
		return h.completeInternal(ctx, apptID.Value(), optEmployeeID, notes)
	}
	http.HandleRequestNoBodyWithResponder(h.validator, logic, func(c *gin.Context, _ any) {
		http.Success(c, nil, "Appointment completed successfully")
	})(c)
}

func (h *AppointmentHandler) CancelAppointment(c *gin.Context) {
	logic := func(ctx *gin.Context) (any, error) {
		apptID, err := parseAppointmentIDFromParam(ctx)
		if err != nil {
			return nil, err
		}
		optEmployeeID, _ := OptionalEmployeeIDFromQuery(ctx)
		reason := ctx.Query("reason")
		return h.cancelInternal(ctx, apptID.Value(), optEmployeeID, reason)
	}
	http.HandleRequestNoBodyWithResponder(h.validator, logic, func(c *gin.Context, _ any) {
		http.Success(c, nil, "Appointment cancelled successfully")
	})(c)
}

func (h *AppointmentHandler) NotAttendAppointment(c *gin.Context) {
	logic := func(ctx *gin.Context) (any, error) {
		apptID, err := parseAppointmentIDFromParam(ctx)
		if err != nil {
			return nil, err
		}
		optEmployeeID, _ := OptionalEmployeeIDFromQuery(ctx)
		return h.notAttendInternal(ctx, apptID.Value(), optEmployeeID)
	}
	http.HandleRequestNoBodyWithResponder(h.validator, logic, func(c *gin.Context, _ any) {
		http.Success(c, nil, "Appointment marked as not attended successfully")
	})(c)
}
