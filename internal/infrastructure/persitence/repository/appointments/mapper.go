package appointments

import (
	"clinic-vet-api/db/models"
	domain "clinic-vet-api/internal/core/appointments"
	"clinic-vet-api/internal/core/customers"
	"clinic-vet-api/internal/core/employees"
	"clinic-vet-api/internal/core/pets"
	"clinic-vet-api/sqlc"

	"github.com/jackc/pgx/v5/pgtype"
)

func rowToEntity(row sqlc.Appointment) domain.Appointment {
	var employeeID *employees.EmployeeID
	if row.EmployeeID.Valid {
		employeeIDValue := employees.NewEmployeeID(uint(row.EmployeeID.Int32))
		employeeID = &employeeIDValue
	}

	var notes *string
	if row.Notes.Valid {
		notes = &row.Notes.String
	}

	appointment := domain.Appointment{
		CustomerID:    customers.NewCustomerID(uint(row.CustomerID)),
		EmployeeID:    employeeID,
		PetID:         pets.NewPetID(uint(row.PetID)),
		ScheduledDate: row.ScheduledDate.Time,
		Status:        domain.AppointmentStatus(row.Status),
		Service:       domain.ClinicService(row.ClinicService),
		Notes:         notes,
	}
	appointment.SetID(domain.NewAppointmentID(uint(row.ID)))
	appointment.SetTimeStamps(row.CreatedAt.Time, row.UpdatedAt.Time)

	return appointment
}

func specRowsToEntities(rows []sqlc.FindAppointmentsBySpecRow) []domain.Appointment {
	appointments := make([]domain.Appointment, len(rows))
	for i, row := range rows {
		appointments[i] = domain.Appointment{
			CustomerID: customers.NewCustomerID(uint(row.CustomerID)),
			EmployeeID: func() *employees.EmployeeID {
				if row.EmployeeID.Valid {
					id := employees.NewEmployeeID(uint(row.EmployeeID.Int32))
					return &id
				} else {
					return nil
				}
			}(),
			PetID:         pets.NewPetID(uint(row.PetID)),
			ScheduledDate: row.ScheduledDate.Time,
			Status:        domain.AppointmentStatus(row.Status),
			Service:       domain.ClinicService(row.ClinicService),
			Notes: func() *string {
				if row.Notes.Valid {
					return &row.Notes.String
				} else {
					return nil
				}
			}(),
		}
	}
	return appointments
}

func toCreateParams(appt *domain.Appointment) sqlc.CreateAppointmentParams {
	params := sqlc.CreateAppointmentParams{
		CustomerID:    int32(appt.CustomerID.Value),
		PetID:         int32(appt.PetID.Value),
		ScheduledDate: pgtype.Timestamptz{Time: appt.ScheduledDate, Valid: true},
		Status:        models.AppointmentStatus(string(appt.Status)),
		ClinicService: models.ClinicService(string(appt.Service.String())),
	}

	if appt.EmployeeID != nil {
		params.EmployeeID = pgtype.Int4{
			Int32: int32(appt.EmployeeID.Value),
			Valid: true,
		}
	} else {
		params.EmployeeID = pgtype.Int4{Valid: false}
	}

	if appt.Notes != nil {
		params.Notes = pgtype.Text{
			String: *appt.Notes,
			Valid:  true,
		}
	} else {
		params.Notes = pgtype.Text{Valid: false}
	}

	return params
}

func toUpdateParams(appt *domain.Appointment) sqlc.UpdateAppointmentParams {
	return sqlc.UpdateAppointmentParams{
		ID:            int32(appt.ID.Value),
		CustomerID:    int32(appt.CustomerID.Value),
		EmployeeID:    pgtype.Int4{Int32: int32(appt.EmployeeID.Value), Valid: appt.EmployeeID != nil},
		PetID:         int32(appt.PetID.Value),
		ScheduledDate: pgtype.Timestamptz{Time: appt.ScheduledDate, Valid: true},
		Notes:         pgtype.Text{String: *appt.Notes, Valid: appt.Notes != nil},
		Status:        models.AppointmentStatus(string(appt.Status)),
	}
}
