package mappers

import (
	"time"

	"clinic-vet-api/internal/core/employees"
	"clinic-vet-api/internal/infrastructure/http/handlers/dtos"
	"clinic-vet-api/internal/shared"
	"clinic-vet-api/internal/shared/page"
)

type EmployeeMapper struct{}

func NewEmployeeMapper() *EmployeeMapper {
	return &EmployeeMapper{}
}

// RequestToCreateCommand maps the HTTP create request DTO to the domain command.
func (m *EmployeeMapper) RequestToCreateCommand(r dtos.EmployeeCreateRequest) (employees.CreateEmployeeCommand, error) {
	gender, err := shared.ParseGender(r.Gender)
	if err != nil {
		return employees.CreateEmployeeCommand{}, err
	}

	dob, err := time.Parse("2006-01-02", r.DateOfBirth)
	if err != nil {
		return employees.CreateEmployeeCommand{}, err
	}

	specialty, err := employees.ParseVetSpecialty(r.Specialty)
	if err != nil {
		return employees.CreateEmployeeCommand{}, err
	}

	isActive := true
	if r.IsActive != nil {
		isActive = *r.IsActive
	}

	schedule := employees.ScheduleDataCommand{
		Day:           r.Schedule.Day,
		EntryTime:     r.Schedule.EntryTime,
		DepartureTime: r.Schedule.DepartureTime,
		StartBreak:    r.Schedule.StartBreak,
		EndBreak:      r.Schedule.EndBreak,
	}

	return employees.CreateEmployeeCommand{
		FirstName:      r.FirstName,
		LastName:       r.LastName,
		Gender:         gender,
		DateOfBirth:    dob,
		Photo:          r.Photo,
		LicenseNumber:  r.LicenseNo,
		YearsExperience: r.YearsExp,
		IsActive:       isActive,
		Specialty:      specialty,
		Schedule:       schedule,
	}, nil
}

// RequestToUpdateCommand maps the HTTP update request DTO to the domain command.
func (m *EmployeeMapper) RequestToUpdateCommand(r dtos.EmployeeUpdateRequest) (employees.UpdateEmployeeCommand, error) {
	id := employees.NewEmployeeID(r.ID)

	var gender *shared.PersonGender
	if r.Gender != nil {
		g, err := shared.ParseGender(*r.Gender)
		if err != nil {
			return employees.UpdateEmployeeCommand{}, err
		}
		gender = &g
	}

	var dob *time.Time
	if r.DateOfBirth != nil {
		d, err := time.Parse("2006-01-02", *r.DateOfBirth)
		if err != nil {
			return employees.UpdateEmployeeCommand{}, err
		}
		dob = &d
	}

	var specialty *employees.VetSpecialty
	if r.Specialty != nil {
		s, err := employees.ParseVetSpecialty(*r.Specialty)
		if err != nil {
			return employees.UpdateEmployeeCommand{}, err
		}
		specialty = &s
	}

	var schedule *employees.ScheduleDataCommand
	if r.Schedule != nil {
		schedule = &employees.ScheduleDataCommand{
			Day:           r.Schedule.Day,
			EntryTime:     r.Schedule.EntryTime,
			DepartureTime: r.Schedule.DepartureTime,
			StartBreak:    r.Schedule.StartBreak,
			EndBreak:      r.Schedule.EndBreak,
		}
	}

	cmd := employees.UpdateEmployeeCommand{
		ID:          &id,
		FirstName:   r.FirstName,
		LastName:    r.LastName,
		Gender:      gender,
		DateOfBirth: dob,
		Photo:       r.Photo,
		YearsExperience: r.YearsExp,
		IsActive:    r.IsActive,
		Specialty:   specialty,
		Schedule:    schedule,
	}

	// License number is a value (not pointer) in the command; only set it when provided.
	if r.LicenseNo != nil {
		cmd.LicenseNumber = *r.LicenseNo
	}

	return cmd, nil
}

// ToEmployeeResponse maps an Employee aggregate to the response DTO.
func (m *EmployeeMapper) ToEmployeeResponse(e employees.Employee) dtos.EmployeeResponse {
	return dtos.EmployeeResponse{
		ID:              e.ID.Value,
		FirstName:       e.FirstName,
		LastName:        e.LastName,
		Gender:          string(e.Gender),
		DateOfBirth:     e.DateOfBirth,
		Photo:           e.Photo,
		LicenseNumber:   e.LicenseNumber,
		Specialty:       e.Specialty.String(),
		YearsExperience: e.YearsExperience,
		IsActive:        e.IsActive,
		UserID:          e.UserID,
		CreatedAt:       e.CreatedAt,
		UpdatedAt:       e.UpdatedAt,
	}
}

// ToEmployeeResponses maps a page of Employee to a page of EmployeeResponse.
func (m *EmployeeMapper) ToEmployeeResponses(p page.Page[employees.Employee]) page.Page[dtos.EmployeeResponse] {
	return page.MapItems(p, m.ToEmployeeResponse)
}

// ToEmployeeStatsResponse maps domain EmployeeStats to the HTTP DTO.
func (m *EmployeeMapper) ToEmployeeStatsResponse(s employees.EmployeeStats) dtos.EmployeeStatsResponse {
	specialties := make(map[string]int64, len(s.Specialties))
	for k, v := range s.Specialties {
		specialties[k.String()] = v
	}
	return dtos.EmployeeStatsResponse{
		TotalEmployees:  s.TotalEmployees,
		ActiveEmployees: s.ActiveEmployees,
		Specialties:     specialties,
	}
}

