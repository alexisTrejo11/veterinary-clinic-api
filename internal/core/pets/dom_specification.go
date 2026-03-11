package pets

import (
	"strconv"
	"strings"

	"clinic-vet-api/internal/shared/page"
)

// PetSpecification holds filter and pagination criteria for querying pets.
type PetSpecification struct {
	IDs          []PetID
	CustomerIDs  []uint
	Species      []PetSpecies
	Genders      []PetGender
	IsActive     *bool
	SearchTerm   *string // name, breed
	page.Pagination
}

// IsSatisfiedBy returns true if the entity matches all specification criteria.
func (s *PetSpecification) IsSatisfiedBy(entity any) bool {
	pet, ok := entity.(*Pet)
	if !ok {
		return false
	}

	if len(s.IDs) > 0 {
		found := false
		for _, id := range s.IDs {
			if id.Value == pet.ID.Value {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	if len(s.CustomerIDs) > 0 {
		found := false
		for _, cid := range s.CustomerIDs {
			if cid == pet.CustomerID {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	if len(s.Species) > 0 {
		found := false
		for _, sp := range s.Species {
			if sp == pet.Species {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	if len(s.Genders) > 0 {
		found := false
		for _, g := range s.Genders {
			if g == pet.Gender {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	if s.IsActive != nil && pet.IsActive != *s.IsActive {
		return false
	}

	if s.SearchTerm != nil && *s.SearchTerm != "" {
		term := strings.ToLower(*s.SearchTerm)
		nameMatch := strings.Contains(strings.ToLower(pet.Name), term)
		breedMatch := false
		if pet.Breed != nil {
			breedMatch = strings.Contains(strings.ToLower(*pet.Breed), term)
		}
		if !nameMatch && !breedMatch {
			return false
		}
	}

	return true
}

// ToSQL returns a WHERE clause and args for use with deleted_at IS NULL.
// Caller must prepend "WHERE deleted_at IS NULL" and append ORDER BY/LIMIT/OFFSET from Pagination.
func (s *PetSpecification) ToSQL() (where string, args []any) {
	var conditions []string
	idx := 1

	if len(s.IDs) > 0 {
		ids := make([]any, len(s.IDs))
		for i, id := range s.IDs {
			ids[i] = id.Value
		}
		conditions = append(conditions, "id = ANY($"+strconv.Itoa(idx)+")")
		args = append(args, ids)
		idx++
	}

	if len(s.CustomerIDs) > 0 {
		conditions = append(conditions, "customer_id = ANY($"+strconv.Itoa(idx)+")")
		args = append(args, s.CustomerIDs)
		idx++
	}

	if len(s.Species) > 0 {
		species := make([]any, len(s.Species))
		for i, sp := range s.Species {
			species[i] = sp.String()
		}
		conditions = append(conditions, "species = ANY($"+strconv.Itoa(idx)+")")
		args = append(args, species)
		idx++
	}

	if len(s.Genders) > 0 {
		genders := make([]any, len(s.Genders))
		for i, g := range s.Genders {
			genders[i] = g.String()
		}
		conditions = append(conditions, "gender = ANY($"+strconv.Itoa(idx)+")")
		args = append(args, genders)
		idx++
	}

	if s.IsActive != nil {
		conditions = append(conditions, "is_active = $"+strconv.Itoa(idx))
		args = append(args, *s.IsActive)
		idx++
	}

	if s.SearchTerm != nil && *s.SearchTerm != "" {
		conditions = append(conditions, "(name ILIKE $"+strconv.Itoa(idx)+" OR breed ILIKE $"+strconv.Itoa(idx)+")")
		args = append(args, "%"+*s.SearchTerm+"%")
	}

	if len(conditions) > 0 {
		where = " AND " + strings.Join(conditions, " AND ")
	}
	return where, args
}

// WithIDs sets the pet IDs filter.
func (s *PetSpecification) WithIDs(ids ...PetID) *PetSpecification {
	s.IDs = ids
	return s
}

// WithCustomerIDs sets the customer IDs filter.
func (s *PetSpecification) WithCustomerIDs(customerIDs ...uint) *PetSpecification {
	s.CustomerIDs = customerIDs
	return s
}

// WithSpecies sets the species filter.
func (s *PetSpecification) WithSpecies(species ...PetSpecies) *PetSpecification {
	s.Species = species
	return s
}

// WithGenders sets the gender filter.
func (s *PetSpecification) WithGenders(genders ...PetGender) *PetSpecification {
	s.Genders = genders
	return s
}

// WithIsActive sets the active filter.
func (s *PetSpecification) WithIsActive(isActive bool) *PetSpecification {
	s.IsActive = &isActive
	return s
}

// WithSearchTerm sets the search term for name/breed.
func (s *PetSpecification) WithSearchTerm(term string) *PetSpecification {
	s.SearchTerm = &term
	return s
}

// WithPagination sets pagination.
func (s *PetSpecification) WithPagination(pageNumber, pageSize int, orderBy, sortDir string) *PetSpecification {
	s.Pagination = page.Pagination{
		Number:  pageNumber,
		Size:    pageSize,
		OrderBy: orderBy,
		SortDir: sortDir,
	}
	return s
}
