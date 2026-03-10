package appointments

import (
	"clinic-vet-api/internal/core/appointments"
	"clinic-vet-api/internal/core/customers"
	"clinic-vet-api/internal/core/employees"
	"clinic-vet-api/internal/shared/mapper"
	p "clinic-vet-api/internal/shared/page"
	"clinic-vet-api/sqlc"
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

type SqlcAppointmentRepository struct {
	queries *sqlc.Queries
	pgMap   *mapper.SqlcFieldMapper
}

func NewSqlcAppointmentRepository(queries *sqlc.Queries) appointments.AppointmentRepository {
	return &SqlcAppointmentRepository{
		queries: queries,
		pgMap:   mapper.NewSqlcFieldMapper(),
	}
}

func (r *SqlcAppointmentRepository) Find(
	ctx context.Context,
	spec *appointments.AppointmentSpecification,
) (p.Page[appointments.Appointment], error) {
	params := spec.ToSearchParams()
	var service string
	if params.Service != nil {
		service = *params.Service
	}

	var status string
	if params.Status != nil {
		status = *params.Status
	}

	sqlcParams := sqlc.FindAppointmentsBySpecParams{
		Column1: r.pgMap.Primitive.Int32PtrToInt32(params.ID),                              // ID
		Column2: r.pgMap.Primitive.Int32PtrToInt32(params.CustomerID),                      // CustomerID
		Column3: r.pgMap.Primitive.Int32PtrToInt32(params.EmployeeID),                      // EmployeeID
		Column4: r.pgMap.Primitive.Int32PtrToInt32(params.PetID),                           // PetID
		Column5: service,                                                                   // Service
		Column6: status,                                                                    // Status
		Column7: pgtype.Timestamp(r.pgMap.PgTimestamptz.FromTimePtr(params.StartDate)),     // StartDate
		Column8: pgtype.Timestamp(r.pgMap.PgTimestamptz.FromTimePtr(params.EndDate)),       // EndDate
		Column9: pgtype.Timestamp(r.pgMap.PgTimestamptz.FromTimePtr(params.ScheduledDate)), // ScheduledDate
		Limit:   params.Limit,
		Offset:  params.Offset,
	}

	rows, err := r.queries.FindAppointmentsBySpec(ctx, sqlcParams)
	if err != nil {
		return p.Page[appointments.Appointment]{}, err
	}

	countParams := sqlc.CountAppointmentsBySpecParams{
		Column1: sqlcParams.Column1, // ID
		Column2: sqlcParams.Column2, // CustomerID
		Column3: sqlcParams.Column3, // EmployeeID
		Column4: sqlcParams.Column4, // PetID
		Column5: sqlcParams.Column5, // Service
		Column6: sqlcParams.Column6, // Status
		Column7: sqlcParams.Column7, // StartDate
		Column8: sqlcParams.Column8, // EndDate
		Column9: sqlcParams.Column9, // ScheduledDate
	}

	appt := specRowsToEntities(rows)

	total, err := r.queries.CountAppointmentsBySpec(ctx, countParams)
	if err != nil {
		return p.Page[appointments.Appointment]{}, err
	}

	// Fix pagination to have default values if not set
	pagination := p.PaginationRequest{
		PageSize: 10,
		Page:     1,
	}

	return p.NewPage(appt, total, pagination), nil
}

func (r *SqlcAppointmentRepository) FindByID(
	ctx context.Context,
	id appointments.AppointmentID,
) (appointments.Appointment, error) {
	sqlRow, err := r.queries.FindAppointmentByID(ctx, id.Int32())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return appointments.Appointment{}, r.notFoundError("id", id.String())
		}
		return appointments.Appointment{}, r.dbError("select", "failed to get appointment by ID", err)
	}

	return rowToEntity(sqlRow), nil
}

func (r *SqlcAppointmentRepository) FindByIDAndCustomerID(
	ctx context.Context,
	id appointments.AppointmentID,
	customerID customers.CustomerID,
) (appointments.Appointment, error) {
	params := sqlc.FindAppointmentByIDAndCustomerIDParams{
		ID:         id.Int32(),
		CustomerID: customerID.Int32(),
	}

	sqlRow, err := r.queries.FindAppointmentByIDAndCustomerID(ctx, params)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return appointments.Appointment{}, r.notFoundError("id and customer_id", fmt.Sprintf("%s and %s", id.String(), customerID.String()))
		}
		return appointments.Appointment{}, r.dbError("select", "failed to get appointment by ID and CustomerID", err)
	}

	return rowToEntity(sqlRow), nil
}

func (r *SqlcAppointmentRepository) FindByIDAndEmployeeID(
	ctx context.Context,
	id appointments.AppointmentID,
	employeeID employees.EmployeeID,
) (appointments.Appointment, error) {
	params := sqlc.FindAppointmentByIDAndEmployeeIDParams{
		ID:         id.Int32(),
		EmployeeID: r.pgMap.PgInt4.FromInt32(employeeID.Int32()),
	}

	sqlRow, err := r.queries.FindAppointmentByIDAndEmployeeID(ctx, params)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return appointments.Appointment{}, r.notFoundError("id and employee_id", fmt.Sprintf("%s and %s", id.String(), employeeID.String()))
		}
		return appointments.Appointment{}, r.dbError("select", "failed to get appointment by ID and EmployeeID", err)
	}
	return rowToEntity(sqlRow), nil
}

func (r *SqlcAppointmentRepository) ExistsByID(
	ctx context.Context,
	id appointments.AppointmentID,
) (bool, error) {
	exists, err := r.queries.ExistsAppointmentID(ctx, id.Int32())
	if err != nil {
		return false, r.dbError("select", "failed to check appointment existence by ID", err)
	}
	return exists, nil
}

func (r *SqlcAppointmentRepository) Save(
	ctx context.Context,
	appointment *appointments.Appointment,
) error {
	if appointment.ID.IsZero() {
		return r.create(ctx, appointment)
	}
	return r.update(ctx, appointment)
}

func (r *SqlcAppointmentRepository) DeleteByID(
	ctx context.Context,
	id appointments.AppointmentID,
	hardDelete bool,
) error {
	if hardDelete {
		return r.hardDeleteByID(ctx, id)
	}
	return r.softDeleteByID(ctx, id)
}

func (r *SqlcAppointmentRepository) RestoreByID(ctx context.Context, id appointments.AppointmentID) error {
	if err := r.queries.RestoreAppointment(ctx, id.Int32()); err != nil {
		return r.dbError("update", "failed to restore appointment", err)
	}
	return nil
}

func (r *SqlcAppointmentRepository) Count(ctx context.Context, spec *appointments.AppointmentSpecification) (int64, error) {
	params := spec.ToSearchParams()
	var service string
	if params.Service != nil {
		service = *params.Service
	}
	var status string
	if params.Status != nil {
		status = *params.Status
	}
	sqlcParams := sqlc.CountAppointmentsBySpecParams{
		Column1: r.pgMap.Primitive.Int32PtrToInt32(params.ID),                              // ID
		Column2: r.pgMap.Primitive.Int32PtrToInt32(params.CustomerID),                      // CustomerID
		Column3: r.pgMap.Primitive.Int32PtrToInt32(params.EmployeeID),                      // EmployeeID
		Column4: r.pgMap.Primitive.Int32PtrToInt32(params.PetID),                           // PetID
		Column5: service,                                                                   // Service
		Column6: status,                                                                    // Status
		Column7: pgtype.Timestamp(r.pgMap.PgTimestamptz.FromTimePtr(params.StartDate)),     // StartDate
		Column8: pgtype.Timestamp(r.pgMap.PgTimestamptz.FromTimePtr(params.EndDate)),       // EndDate
		Column9: pgtype.Timestamp(r.pgMap.PgTimestamptz.FromTimePtr(params.ScheduledDate)), // ScheduledDate
	}
	total, err := r.queries.CountAppointmentsBySpec(ctx, sqlcParams)
	if err != nil {
		return 0, r.dbError("count", "failed to count appointments", err)
	}
	return total, nil
}

func (r *SqlcAppointmentRepository) create(ctx context.Context, appointment *appointments.Appointment) error {
	params := toCreateParams(appointment)

	_, err := r.queries.CreateAppointment(ctx, params)
	if err != nil {
		return r.dbError("insert", "failed to create appointment", err)
	}

	return nil
}

func (r *SqlcAppointmentRepository) update(ctx context.Context, appointment *appointments.Appointment) error {
	params := toUpdateParams(appointment)

	_, err := r.queries.UpdateAppointment(ctx, params)
	if err != nil {
		return r.dbError("update", fmt.Sprintf("failed to update appointment with ID %d", appointment.ID.Value), err)
	}

	return nil
}

func (r *SqlcAppointmentRepository) softDeleteByID(ctx context.Context, id appointments.AppointmentID) error {
	if err := r.queries.SoftDeleteAppointment(ctx, id.Int32()); err != nil {
		return r.dbError("update", "failed to soft delete appointment", err)
	}
	return nil
}

func (r *SqlcAppointmentRepository) hardDeleteByID(ctx context.Context, id appointments.AppointmentID) error {
	if err := r.queries.HardDeleteAppointment(ctx, id.Int32()); err != nil {
		return r.dbError("delete", "failed to hard delete appointment", err)
	}
	return nil
}
