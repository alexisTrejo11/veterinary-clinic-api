package dtos

import (
	"clinic-vet-api/internal/shared/page"
)

// AddressSearchRequest carries pagination for address listing
// @Description Pagination and optional filters for searching addresses.
type AddressSearchRequest struct {
	page.PaginationRequest
}

// AddressCreateRequest is the body for creating an address
// @Description Street, city, state, zip, country, building type/numbers, and whether it is the default address.
type AddressCreateRequest struct {
	Street              string  `json:"street"`
	City                string  `json:"city"`
	State               string  `json:"state"`
	ZipCode             string  `json:"zip_code"`
	Country             string  `json:"country"`
	BuildingType        string  `json:"building_type"`
	BuildingOuterNumber string  `json:"building_outer_number"`
	BuildingInnerNumber *string `json:"building_inner_number,omitempty"`
	IsDefault           bool    `json:"is_default"`
}

// AddressUpdateRequest is the body for updating an address
// @Description Optional fields to update; only provided fields are changed.
type AddressUpdateRequest struct {
	ID                  uint    `json:"id"`
	Street              *string `json:"street,omitempty"`
	City                *string `json:"city,omitempty"`
	State               *string `json:"state,omitempty"`
	ZipCode             *string `json:"zip_code,omitempty"`
	Country             *string `json:"country,omitempty"`
	BuildingType        *string `json:"building_type,omitempty"`
	BuildingOuterNumber *string `json:"building_outer_number,omitempty"`
	BuildingInnerNumber *string `json:"building_inner_number,omitempty"`
	IsDefault           *bool   `json:"is_default,omitempty"`
}

// AddressResponse is returned when reading an address
// @Description Address entity: id, street, city, state, zip, country, building details, is_default.
type AddressResponse struct {
	ID                  uint    `json:"id"`
	Street              string  `json:"street"`
	City                string  `json:"city"`
	State               string  `json:"state"`
	ZipCode             string  `json:"zip_code"`
	Country             string  `json:"country"`
	BuildingType        string  `json:"building_type"`
	BuildingOuterNumber string  `json:"building_outer_number"`
	BuildingInnerNumber *string `json:"building_inner_number,omitempty"`
	IsDefault           bool    `json:"is_default"`
}
