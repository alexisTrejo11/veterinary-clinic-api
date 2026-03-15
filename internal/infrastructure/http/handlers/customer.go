package handlers

import (
	"clinic-vet-api/internal/core/customers"
	"clinic-vet-api/internal/infrastructure/http/handlers/dtos"
	"clinic-vet-api/internal/infrastructure/http/handlers/mappers"
	"clinic-vet-api/internal/shared/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type CustomerHandler struct {
	service   customers.CustomerService
	validator *validator.Validate
	mapper    *mappers.CustomerMapper
}

func NewCustomerHandler(service customers.CustomerService, validator *validator.Validate) *CustomerHandler {
	return &CustomerHandler{service: service, validator: validator}
}

func (s *CustomerHandler) GetCustomerByID(c *gin.Context) {
	customerIDUInt, err := http.ParseParamToUInt(c, "id")
	if err != nil {
		http.BadRequest(c, err)
		return
	}

	customerID := customers.NewCustomerID(customerIDUInt)
	customer, err := s.service.GetCustomerByID(c.Request.Context(), customerID)
	if err != nil {
		http.ApplicationError(c, err)
		return
	}

	customerResponse := s.mapper.ToCustomerResponse(customer)
	http.Found(c, customerResponse, "Customer")
}

func (s *CustomerHandler) SearchCustomers(c *gin.Context) {
	var requestBodyData dtos.CustomerSearchRequest
	if err := http.ShouldBindAndValidateBody(c, &requestBodyData, s.validator); err != nil {
		http.BadRequest(c, err)
		return
	}
}

func (s *CustomerHandler) CreateCustomer(c *gin.Context) {
	var requestBodyData dtos.CustomerCreateRequest
	if err := http.ShouldBindAndValidateBody(c, &requestBodyData, s.validator); err != nil {
		http.BadRequest(c, err)
		return
	}

	command, err := s.mapper.RequestToCreateCommand(requestBodyData)
	if err != nil {
		http.BadRequest(c, err)
		return
	}

	customer, err := s.service.CreateCustomer(c.Request.Context(), command)
	if err != nil {
		http.ApplicationError(c, err)
		return
	}

	http.Created(c, customer.ID.Value(), "Customer")
}

func (s *CustomerHandler) UpdateCustomer(c *gin.Context) {
	var requestBodyData dtos.CustomerUpdateRequest
	if err := http.ShouldBindAndValidateBody(c, &requestBodyData, s.validator); err != nil {
		http.BadRequest(c, err)
		return
	}

	command, err := s.mapper.RequestToUpdateCommand(requestBodyData)
	if err != nil {
		http.BadRequest(c, err)
		return
	}
	err = s.service.UpdateCustomer(c.Request.Context(), command)
	if err != nil {
		http.ApplicationError(c, err)
		return
	}

	http.Updated(c, nil, "Customer")
}

func (s *CustomerHandler) DeleteCustomer(c *gin.Context) {
	entityID, err := http.ParseParamToUInt(c, "id")
	if err != nil {
		http.BadRequest(c, err)
		return
	}

	customerID := customers.NewCustomerID(entityID)
	err = s.service.DeleteCustomer(c.Request.Context(), customerID)
	if err != nil {
		http.ApplicationError(c, err)
		return
	}

	http.Success(c, nil, "Customer")
}

func (s *CustomerHandler) RestoreCustomer(c *gin.Context) {
	entityID, err := http.ParseParamToUInt(c, "id")
	if err != nil {
		http.BadRequest(c, err)
		return
	}

	customerID := customers.NewCustomerID(entityID)
	err = s.service.RestoreCustomer(c.Request.Context(), customerID)
	if err != nil {
		http.ApplicationError(c, err)
		return
	}

	http.Success(c, nil, "Customer")
}
