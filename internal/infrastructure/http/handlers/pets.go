package handlers

import (
	"clinic-vet-api/internal/core/pets"
	"clinic-vet-api/internal/infrastructure/http/handlers/dtos"
	"clinic-vet-api/internal/infrastructure/http/handlers/mappers"
	"clinic-vet-api/internal/shared/errors"
	"clinic-vet-api/internal/shared/http"
	"clinic-vet-api/internal/shared/page"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type BasePetHandler struct {
	service   pets.Service
	validator *validator.Validate
	mapper    *mappers.PetMapper
}

func NewBasePetHandler(service pets.Service, validator *validator.Validate) *BasePetHandler {
	return &BasePetHandler{service: service, validator: validator}
}

func (s *BasePetHandler) SearchPets(c *gin.Context) {
	var requestBodyData dtos.PetSearchRequest
	if err := http.ShouldBindAndValidateBody(c, &requestBodyData, s.validator); err != nil {
		http.BadRequest(c, err)
		return
	}

	query, err := s.mapper.RequestToSpecification(requestBodyData)
	if err != nil {
		http.BadRequest(c, err)
		return
	}

	petPage, err := s.service.GetPetBySpecification(c.Request.Context(), &query)
	if err != nil {
		http.ApplicationError(c, err)
		return
	}

	responsePage := s.mapper.ToPaginatedResponse(petPage)
	http.Paginated(c, &responsePage, "Pets")
}

func (s *BasePetHandler) CreatePet(c *gin.Context, customerID *uint, isActive bool) {
	var requestBodyData dtos.PetCreateRequest
	if err := http.ShouldBindAndValidateBody(c, &requestBodyData, s.validator); err != nil {
		http.BadRequest(c, err)
		return
	}

	command, err := s.mapper.RequestToCreateCommand(requestBodyData, *customerID, isActive)
	if err != nil {
		http.BadRequest(c, err)
		return
	}

	pet, err := s.service.CreatePet(c.Request.Context(), command)
	if err != nil {
		http.ApplicationError(c, err)
		return
	}

	http.Created(c, pet.ID.Value, "Pet")
}

func (s *BasePetHandler) RestorePet(c *gin.Context) {
	petIDUInt, err := http.ParseParamToUInt(c, "id")
	if err != nil {
		http.BadRequest(c, errors.RequestURLParamError(err, "id", c.Param("id")))
		return
	}

	petID := pets.NewPetID(petIDUInt)
	err = s.service.RestorePet(c.Request.Context(), petID)
	if err != nil {
		http.ApplicationError(c, err)
		return
	}

	http.Updated(c, nil, "Pet")
}

func (s *BasePetHandler) DeletePet(c *gin.Context) {
	petIDUInt, err := http.ParseParamToUInt(c, "id")
	if err != nil {
		http.BadRequest(c, errors.RequestURLParamError(err, "id", c.Param("id")))
		return
	}

	// TODO: Include customer ID, maybe as a query parameter or inside gin context
	petID := pets.NewPetID(petIDUInt)
	command := pets.DeletePetCommand{
		PetID:        petID,
		IsHardDelete: false,
	}

	err = s.service.DeletePet(c.Request.Context(), command)
	if err != nil {
		http.ApplicationError(c, err)
		return
	}

	http.Success(c, nil, "Pet")
}

func (s *BasePetHandler) GetPetByID(c *gin.Context) {
	petIDUInt, err := http.ParseParamToUInt(c, "id")
	if err != nil {
		http.BadRequest(c, errors.RequestURLParamError(err, "id", c.Param("id")))
		return
	}

	petID := pets.NewPetID(petIDUInt)
	pet, err := s.service.GetPetByID(c.Request.Context(), petID)
	if err != nil {
		http.ApplicationError(c, err)
		return
	}

	http.Success(c, pet, "Pet")
}

func (s *BasePetHandler) GetPetsByCustomerID(c *gin.Context) {
	customerIDUInt, err := http.ParseParamToUInt(c, "customer_id")
	if err != nil {
		http.BadRequest(c, errors.RequestURLParamError(err, "customer_id", c.Param("customer_id")))
		return
	}

	var pageParams page.PaginationRequest
	if err := http.ShouldBindPageParams(&pageParams, c, s.validator); err != nil {
		http.BadRequest(c, errors.RequestURLQueryError(err, c.Request.URL.RawQuery))
		return
	}
	pagination := pageParams.ToPagination()

	petPage, err := s.service.GetPetsByCustomerID(c.Request.Context(), customerIDUInt, pagination)
	if err != nil {
		http.ApplicationError(c, err)
		return
	}
	responsePage := s.mapper.ToPaginatedResponse(petPage)
	http.Paginated(c, &responsePage, "Pets")
}
