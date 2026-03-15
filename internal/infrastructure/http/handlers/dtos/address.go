package dtos

import (
	"clinic-vet-api/internal/shared/page"
)

type AddressSearchRequest struct {
	page.PaginationRequest
}

type AddressCreateRequest struct {
	Street              string
	City                string
	State               string
	ZipCode             string
	Country             string
	BuildingType        string
	BuildingOuterNumber string
	BuildingInnerNumber *string
	IsDefault           bool
}

type AddressUpdateRequest struct {
	ID                  uint
	Street              *string
	City                *string
	State               *string
	ZipCode             *string
	Country             *string
	BuildingType        *string
	BuildingOuterNumber *string
	BuildingInnerNumber *string
	IsDefault           *bool
}

type AddressResponse struct {
	ID                  uint
	Street              string
	City                string
	State               string
	ZipCode             string
	Country             string
	BuildingType        string
	BuildingOuterNumber string
	BuildingInnerNumber *string
	IsDefault           bool
}
