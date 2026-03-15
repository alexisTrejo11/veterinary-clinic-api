package mappers

import (
	"clinic-vet-api/internal/core/addresses"
	"clinic-vet-api/internal/infrastructure/http/handlers/dtos"
	"clinic-vet-api/internal/shared"
	"clinic-vet-api/internal/shared/page"
)

type AddressMapper struct{}

func NewAddressMapper() *AddressMapper {
	return &AddressMapper{}
}

// TODO: Add filters
func (m *AddressMapper) RequestToSpecification(r dtos.AddressSearchRequest) (addresses.AddressSpecification, error) {
	return addresses.AddressSpecification{
		Pagination: r.PaginationRequest.ToPagination(),
	}, nil
}

func (m *AddressMapper) RequestToCreateCommand(r dtos.AddressCreateRequest, userID shared.UserID) (addresses.CreateAddressCommand, error) {
	country, err := addresses.ParseCountry(r.Country)
	if err != nil {
		return addresses.CreateAddressCommand{}, err
	}
	buildingType, err := addresses.ParseBuildingType(r.BuildingType)
	if err != nil {
		return addresses.CreateAddressCommand{}, err
	}
	return addresses.CreateAddressCommand{
		UserID:              userID,
		Street:              r.Street,
		City:                r.City,
		State:               r.State,
		ZipCode:             r.ZipCode,
		Country:             country,
		BuildingType:        buildingType,
		BuildingOuterNumber: r.BuildingOuterNumber,
		BuildingInnerNumber: r.BuildingInnerNumber,
		IsDefault:           r.IsDefault,
	}, nil
}

func (m *AddressMapper) RequestToUpdateCommand(r dtos.AddressUpdateRequest, userID *shared.UserID) (addresses.UpdateAddressCommand, error) {
	var country *addresses.Country
	if r.Country != nil {
		countryObj, err := addresses.ParseCountry(*r.Country)
		if err != nil {
			return addresses.UpdateAddressCommand{}, err
		}
		country = &countryObj
	}

	var buildingType *addresses.BuildingType
	if r.BuildingType != nil {
		buildingTypeObj, err := addresses.ParseBuildingType(*r.BuildingType)
		if err != nil {
			return addresses.UpdateAddressCommand{}, err
		}
		buildingType = &buildingTypeObj
	}

	return addresses.UpdateAddressCommand{
		ID:                  addresses.NewAddressID(r.ID),
		UserID:              userID,
		Street:              r.Street,
		City:                r.City,
		State:               r.State,
		ZipCode:             r.ZipCode,
		Country:             country,
		BuildingType:        buildingType,
		BuildingOuterNumber: r.BuildingOuterNumber,
		BuildingInnerNumber: r.BuildingInnerNumber,
		IsDefault:           r.IsDefault,
	}, nil
}

func (m *AddressMapper) ToResponse(a addresses.Address) dtos.AddressResponse {
	return dtos.AddressResponse{
		ID:                  a.ID.Value(),
		Street:              a.Street,
		City:                a.City,
		State:               a.State,
		ZipCode:             a.ZipCode,
		Country:             a.Country.DisplayName(),
		BuildingType:        a.BuildingType.DisplayName(),
		BuildingOuterNumber: a.BuildingOuterNumber,
		BuildingInnerNumber: a.BuildingInnerNumber,
		IsDefault:           a.IsDefault,
	}
}

func (m *AddressMapper) ToPaginatedResponse(a page.Page[addresses.Address]) page.Page[dtos.AddressResponse] {
	return page.MapItems(a, m.ToResponse)
}

func (m *AddressMapper) ToCustomerAddressesResponse(a []addresses.Address) []dtos.AddressResponse {
	addressesResponses := make([]dtos.AddressResponse, len(a))
	for i, address := range a {
		addressesResponses[i] = m.ToResponse(address)
	}
	return addressesResponses
}
