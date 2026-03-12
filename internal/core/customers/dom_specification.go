package customers

import (
	"clinic-vet-api/internal/shared/page"
	"strings"
	"time"
)

// CustomerSpecification defines filters and pagination for customer queries.
type CustomerSpecification struct {
	IDs        []CustomerID
	UserIDs    []uint
	IsActive   *bool
	MinAge     *int
	MaxAge     *int
	CreatedFrom *time.Time
	CreatedTo   *time.Time
	SearchTerm *string // matches first/last name or email
	page.Pagination
}

// IsSatisfiedBy is an in-memory predicate useful for testing.
func (s *CustomerSpecification) IsSatisfiedBy(entity any) bool {
	c, ok := entity.(Customer)
	if !ok {
		return false
	}

	if len(s.IDs) > 0 {
		match := false
		for _, id := range s.IDs {
			if id.Value == c.ID.Value {
				match = true
				break
			}
		}
		if !match {
			return false
		}
	}

	if len(s.UserIDs) > 0 {
		match := false
		for _, uid := range s.UserIDs {
			if c.UserID == uid {
				match = true
				break
			}
		}
		if !match {
			return false
		}
	}

	if s.IsActive != nil && c.IsActive != *s.IsActive {
		return false
	}

	age := c.Person.Age()
	if s.MinAge != nil && age < *s.MinAge {
		return false
	}
	if s.MaxAge != nil && age > *s.MaxAge {
		return false
	}

	if s.CreatedFrom != nil && c.CreatedAt.Before(*s.CreatedFrom) {
		return false
	}
	if s.CreatedTo != nil && c.CreatedAt.After(*s.CreatedTo) {
		return false
	}

	if s.SearchTerm != nil && *s.SearchTerm != "" {
		term := strings.ToLower(*s.SearchTerm)
		fullName := strings.ToLower(c.FirstName + " " + c.LastName)
		if !strings.Contains(fullName, term) {
			return false
		}
	}

	return true
}

// WithIDs filters by customer IDs.
func (s *CustomerSpecification) WithIDs(ids ...CustomerID) *CustomerSpecification {
	s.IDs = ids
	return s
}

// WithUserIDs filters by associated user IDs.
func (s *CustomerSpecification) WithUserIDs(ids ...uint) *CustomerSpecification {
	s.UserIDs = ids
	return s
}

// WithIsActive filters by active flag.
func (s *CustomerSpecification) WithIsActive(active bool) *CustomerSpecification {
	s.IsActive = &active
	return s
}

// WithAgeRange filters by age range.
func (s *CustomerSpecification) WithAgeRange(min, max *int) *CustomerSpecification {
	s.MinAge = min
	s.MaxAge = max
	return s
}

// WithCreatedRange filters by created_at range.
func (s *CustomerSpecification) WithCreatedRange(from, to *time.Time) *CustomerSpecification {
	s.CreatedFrom = from
	s.CreatedTo = to
	return s
}

// WithSearchTerm sets the text search term.
func (s *CustomerSpecification) WithSearchTerm(term string) *CustomerSpecification {
	s.SearchTerm = &term
	return s
}

// WithPagination sets pagination details.
func (s *CustomerSpecification) WithPagination(p page.Pagination) *CustomerSpecification {
	s.Pagination = p
	return s
}

