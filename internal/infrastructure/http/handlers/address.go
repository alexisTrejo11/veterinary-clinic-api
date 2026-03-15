package handlers

import (
	"clinic-vet-api/internal/core/addresses"
	"clinic-vet-api/internal/infrastructure/http/handlers/dtos"
	"clinic-vet-api/internal/infrastructure/http/handlers/mappers"
	"clinic-vet-api/internal/shared/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type AddressHandler struct {
	service   addresses.AddressService
	validator *validator.Validate
	mapper    *mappers.AddressMapper
}

func NewAddressHandler(service addresses.AddressService, validator *validator.Validate, mapper *mappers.AddressMapper) *AddressHandler {
	if mapper == nil {
		mapper = mappers.NewAddressMapper()
	}
	return &AddressHandler{service: service, validator: validator, mapper: mapper}
}

// ------------------------------------------------------------
// Internal handlers
// ------------------------------------------------------------

func (s *AddressHandler) getAddressByIDInternal(ctx *gin.Context, getID EntityIDProvider, getUserID OptionalIdentityProvider) (any, error) {
	entityID, err := getID(ctx)
	if err != nil {
		return nil, err
	}

	userID, err := getUserID(ctx)
	if err != nil {
		return nil, err
	}

	addressID := entityID.(addresses.AddressID)
	address, err := s.service.GetAddressByIDAndUserID(ctx.Request.Context(), addressID, userID)
	return address, nil
}

func (s *AddressHandler) getAddressesByUserIDInternal(ctx *gin.Context, getUserID IdentityProvider) (any, error) {
	userID, err := getUserID(ctx)
	if err != nil {
		return nil, err
	}

	addresses, err := s.service.GetAddressesByUserID(ctx.Request.Context(), userID)
	return addresses, nil
}

func (s *AddressHandler) createAddressInternal(ctx *gin.Context, req dtos.AddressCreateRequest, getID IdentityProvider) (any, error) {
	userID, err := getID(ctx)
	if err != nil {
		return nil, err
	}

	command, err := s.mapper.RequestToCreateCommand(req, userID)
	if err != nil {
		return nil, err
	}

	address, err := s.service.CreateAddress(ctx.Request.Context(), command)
	if err != nil {
		return nil, err
	}
	return address.ID.Value(), nil
}

func (s *AddressHandler) updateAddressInternal(ctx *gin.Context, req dtos.AddressUpdateRequest, getID OptionalIdentityProvider) (any, error) {
	userID, err := getID(ctx)
	if err != nil {
		return nil, err
	}

	command, err := s.mapper.RequestToUpdateCommand(req, userID)
	if err != nil {
		return nil, err
	}

	err = s.service.UpdateAddress(ctx.Request.Context(), command)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *AddressHandler) deleteAddressInternal(ctx *gin.Context, getID OptionalIdentityProvider) (any, error) {
	userID, err := getID(ctx)
	if err != nil {
		return nil, err
	}

	entityID, err := http.ParseParamToUInt(ctx, "id")
	if err != nil {
		return nil, err
	}

	addressID := addresses.NewAddressID(entityID)
	err = s.service.DeleteAddress(ctx.Request.Context(), addressID, userID)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// ------------------------------------------------------------
// Manager level handlers
// ------------------------------------------------------------

func (s *AddressHandler) GetAddressByID(c *gin.Context) {
	logic := func(ctx *gin.Context) (any, error) {
		return s.getAddressByIDInternal(ctx, AddressIDFromParam, UserIDFromContextPtr)
	}
	http.HandleRequestNoBodyStruct(s.validator, logic)(c)
}

func (s *AddressHandler) SearchAddresses(c *gin.Context) {
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

func (s *AddressHandler) GetAddressesByUserID(c *gin.Context) {
	logic := func(ctx *gin.Context) (any, error) {
		return s.getAddressesByUserIDInternal(ctx, UserIDFromParam)
	}
	http.HandleRequestNoBodyStruct(s.validator, logic)(c)
}

func (s *AddressHandler) CreateAddress(c *gin.Context) {
	logic := func(ctx *gin.Context, req dtos.AddressCreateRequest) (any, error) {
		return s.createAddressInternal(ctx, req, UserIDFromParam)
	}
	http.HandleRequest(s.validator, logic)(c)
}

func (s *AddressHandler) UpdateAddress(c *gin.Context) {
	logic := func(ctx *gin.Context, req dtos.AddressUpdateRequest) (any, error) {
		return s.updateAddressInternal(ctx, req, UserIDFromParamOptional)
	}
	http.HandleRequest(s.validator, logic)(c)
}

func (s *AddressHandler) DeleteAddress(c *gin.Context) {
	logic := func(ctx *gin.Context) (any, error) {
		return s.deleteAddressInternal(ctx, UserIDFromParamOptional)
	}
	http.HandleRequestNoBodyStruct(s.validator, logic)(c)
}

// ------------------------------------------------------------
// Customer level handlers
// ------------------------------------------------------------

func (s *AddressHandler) CreateAddressForMe(c *gin.Context) {
	logic := func(ctx *gin.Context, req dtos.AddressCreateRequest) (any, error) {
		return s.createAddressInternal(ctx, req, UserIDFromContext)
	}
	http.HandleRequest(s.validator, logic)(c)
}

func (s *AddressHandler) UpdateMyAddress(c *gin.Context) {
	logic := func(ctx *gin.Context, req dtos.AddressUpdateRequest) (any, error) {
		return s.updateAddressInternal(ctx, req, UserIDFromContextPtr)
	}
	http.HandleRequest(s.validator, logic)(c)
}

func (s *AddressHandler) DeleteMyAddress(c *gin.Context) {
	logic := func(ctx *gin.Context) (any, error) {
		return s.deleteAddressInternal(ctx, UserIDFromContextPtr)
	}
	http.HandleRequestNoBodyStruct(s.validator, logic)(c)
}
