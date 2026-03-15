package handlers

import (
	"clinic-vet-api/internal/core/employees"
	"clinic-vet-api/internal/infrastructure/http/handlers/dtos"
	"clinic-vet-api/internal/infrastructure/http/handlers/mappers"
	"clinic-vet-api/internal/shared/http"
	"clinic-vet-api/internal/shared/page"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type BaseEmployeeHandler struct {
	service   employees.EmployeeService
	validator *validator.Validate
	mapper    *mappers.EmployeeMapper
}

func NewBaseEmployeeHandler(service employees.EmployeeService, validator *validator.Validate) *BaseEmployeeHandler {
	return &BaseEmployeeHandler{
		service:   service,
		validator: validator,
		mapper:    mappers.NewEmployeeMapper(),
	}
}

func (s *BaseEmployeeHandler) GetEmployeeByID(c *gin.Context) {
	employeeIDUInt, err := http.ParseParamToUInt(c, "id")
	if err != nil {
		http.BadRequest(c, err)
		return
	}

	employeeID := employees.NewEmployeeID(employeeIDUInt)
	employee, err := s.service.GetEmployeeByID(c.Request.Context(), employeeID)
	if err != nil {
		http.ApplicationError(c, err)
		return
	}

	employeeResponse := s.mapper.ToEmployeeResponse(employee)
	http.Found(c, employeeResponse, "Employee")
}

func (s *BaseEmployeeHandler) GetActiveEmployees(c *gin.Context) {
	var requestBodyData dtos.EmployeeSearchRequest
	if err := http.ShouldBindAndValidateQuery(c, &requestBodyData, s.validator); err != nil {
		http.BadRequest(c, err)
		return
	}

	pagination := requestBodyData.PaginationRequest.ToPagination()

	employeePage, err := s.service.GetActiveEmployees(c.Request.Context(), pagination)
	if err != nil {
		http.ApplicationError(c, err)
		return
	}

	employeeResponses := s.mapper.ToEmployeeResponses(employeePage)
	http.Paginated(c, &employeeResponses, "Employees")
}

func (s *BaseEmployeeHandler) GetEmployeesBySpecialty(c *gin.Context) {
	var pageRequest page.PaginationRequest
	if err := http.ShouldBindPageParams(&pageRequest, c, s.validator); err != nil {
		if err != nil {
			http.BadRequest(c, err)
			return
		}
	}
	specialty := c.Param("specialty")
	specialtyEnum, err := employees.ParseVetSpecialty(specialty)
	if err != nil {
		http.BadRequest(c, err)
		return
	}

	employeesPage, err := s.service.GetEmployeesBySpecialty(c.Request.Context(), specialtyEnum, pageRequest.ToPagination())
	if err != nil {
		http.ApplicationError(c, err)
		return
	}

	employeeResponses := s.mapper.ToEmployeeResponses(employeesPage)
	http.Paginated(c, &employeeResponses, "Employees")
}

func (s *BaseEmployeeHandler) GetEmployeeStats(c *gin.Context) {
	stats, err := s.service.GetEmployeeStats(c.Request.Context())
	if err != nil {
		http.ApplicationError(c, err)
		return
	}

	employeeStatsResponse := s.mapper.ToEmployeeStatsResponse(stats)
	http.Found(c, employeeStatsResponse, "Employee Stats")
}

func (s *BaseEmployeeHandler) GetEmployeesBySpecification(c *gin.Context) {
	var requestBodyData dtos.EmployeeSearchRequest
	if err := http.ShouldBindAndValidateBody(c, &requestBodyData, s.validator); err != nil {
		http.BadRequest(c, err)
		return
	}
}

func (s *BaseEmployeeHandler) CreateEmployee(c *gin.Context) {
	var requestBodyData dtos.EmployeeCreateRequest
	if err := http.ShouldBindAndValidateBody(c, &requestBodyData, s.validator); err != nil {
		http.BadRequest(c, err)
		return
	}

	cmd, err := s.mapper.RequestToCreateCommand(requestBodyData)
	if err != nil {
		http.BadRequest(c, err)
		return
	}

	employee, err := s.service.CreateEmployee(c.Request.Context(), cmd)
	if err != nil {
		http.ApplicationError(c, err)
		return
	}

	http.Created(c, employee.ID.Value(), "Employee")
}

func (s *BaseEmployeeHandler) UpdateEmployee(c *gin.Context) {
	var requestBodyData dtos.EmployeeUpdateRequest
	if err := http.ShouldBindAndValidateBody(c, &requestBodyData, s.validator); err != nil {
		http.BadRequest(c, err)
		return
	}

	cmd, err := s.mapper.RequestToUpdateCommand(requestBodyData)
	if err != nil {
		http.BadRequest(c, err)
		return
	}

	if err := s.service.UpdateEmployee(c.Request.Context(), cmd); err != nil {
		http.ApplicationError(c, err)
		return
	}

	http.Updated(c, nil, "Employee")
}

func (s *BaseEmployeeHandler) RestoreEmployee(c *gin.Context) {
	employeeIDUInt, err := http.ParseParamToUInt(c, "id")
	if err != nil {
		http.BadRequest(c, err)
		return
	}

	employeeID := employees.NewEmployeeID(employeeIDUInt)
	err = s.service.RestoreEmployee(c.Request.Context(), employeeID)
	if err != nil {
		http.ApplicationError(c, err)
		return
	}
}

func (s *BaseEmployeeHandler) DeleteEmployee(c *gin.Context) {
	employeeIDUInt, err := http.ParseParamToUInt(c, "id")
	if err != nil {
		http.BadRequest(c, err)
		return
	}

	employeeID := employees.NewEmployeeID(employeeIDUInt)
	err = s.service.DeleteEmployee(c.Request.Context(), employeeID)
	if err != nil {
		http.ApplicationError(c, err)
		return
	}
}
