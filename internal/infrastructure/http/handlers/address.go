package handlers

import (
	"clinic-vet-api/internal/core/addresses"
	"clinic-vet-api/internal/infrastructure/http/handlers/dtos"
	"clinic-vet-api/internal/infrastructure/http/handlers/mappers"
	"clinic-vet-api/internal/shared/errors"
	"clinic-vet-api/internal/shared/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type BaseAddressHandler struct {
	service   addresses.AddressService
	validator *validator.Validate
	mapper    *mappers.AddressMapper
}

func NewBaseAddressHandler(service addresses.AddressService, validator *validator.Validate) *BaseAddressHandler {
	return &BaseAddressHandler{service: service, validator: validator}
}

func (s *BaseAddressHandler) SearchAddresses(c *gin.Context) {
	var requestBodyData dtos.AddressSearchRequest
	if err := http.ShouldBindAndValidateBody(c, &requestBodyData, s.validator); err != nil {
		http.BadRequest(c, err)
		return
	}

	query, err := s.mapper.RequestToSpecification(requestBodyData)
	if err != nil {
		http.BadRequest(c, err)
		return
	}

	addressPage, err := s.service.GetAddressesBySpecification(c.Request.Context(), query)
	if err != nil {
		http.ApplicationError(c, err)
		return
	}

	responsePage := s.mapper.ToPaginatedResponse(addressPage)
	http.Paginated(c, &responsePage, "Addresses")
}

func (s *BaseAddressHandler) GetAddressesByCustomerID(c *gin.Context) {
	customerIDUInt, err := http.ParseParamToUInt(c, "customer_id")
	if err != nil {
		http.BadRequest(c, errors.RequestURLParamError(err, "customer_id", c.Param("customer_id")))
		return
	}

	addresses, err := s.service.GetAddressesByCustomerID(c.Request.Context(), customerIDUInt)
	if err != nil {
		http.ApplicationError(c, err)
		return
	}

	responsePage := s.mapper.ToCustomerAddressesResponse(addresses)
	http.Success(c, responsePage, "Addresses successfully retrieved")
}

func (s *BaseAddressHandler) CreateAddress(c *gin.Context) {
	var requestBodyData dtos.AddressCreateRequest
	if err := http.ShouldBindAndValidateBody(c, &requestBodyData, s.validator); err != nil {
		http.BadRequest(c, err)
		return
	}

	command, err := s.mapper.RequestToCreateCommand(requestBodyData)
	if err != nil {
		http.BadRequest(c, err)
		return
	}

	address, err := s.service.CreateAddress(c.Request.Context(), command)
	if err != nil {
		http.ApplicationError(c, err)
		return
	}

	http.Created(c, address.ID.Value, "Address")
}

func (s *BaseAddressHandler) UpdateAddress(c *gin.Context) {
	var requestBodyData dtos.AddressUpdateRequest
	if err := http.ShouldBindAndValidateBody(c, &requestBodyData, s.validator); err != nil {
		http.BadRequest(c, err)
		return
	}

	command, err := s.mapper.RequestToUpdateCommand(requestBodyData)
	if err != nil {
		http.BadRequest(c, err)
		return
	}
	err = s.service.UpdateAddress(c.Request.Context(), command)
	if err != nil {
		http.ApplicationError(c, err)
		return
	}

	http.Updated(c, nil, "Address")
}

func (s *BaseAddressHandler) DeleteAddress(c *gin.Context) {
	entityID, err := http.ParseParamToUInt(c, "id")
	if err != nil {
		http.BadRequest(c, err)
		return
	}

	addressID := addresses.NewAddressID(entityID)
	err = s.service.DeleteAddress(c.Request.Context(), addressID)
	if err != nil {
		http.ApplicationError(c, err)
		return
	}

	http.Success(c, nil, "Address")
}

func (s *BaseAddressHandler) GetAddressByID(c *gin.Context) {
	entityID, err := http.ParseParamToUInt(c, "id")
	if err != nil {
		http.BadRequest(c, err)
		return
	}

	customerIDUInt, err := http.ParseParamToUInt(c, "customer_id")
	if err != nil {
		http.BadRequest(c, err)
		return
	}

	addressID := addresses.NewAddressID(entityID)
	address, err := s.service.GetAddressByIDAndCustomerID(c.Request.Context(), addressID, customerIDUInt)
	if err != nil {
		http.ApplicationError(c, err)
		return
	}

	http.Found(c, address, "Address")
}
