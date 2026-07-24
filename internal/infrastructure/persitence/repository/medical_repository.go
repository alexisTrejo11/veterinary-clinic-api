package repository

import (
	"clinic-vet-api/database/models"
	"clinic-vet-api/internal/core/medical"
	"clinic-vet-api/internal/shared"
	"clinic-vet-api/internal/shared/mapper"
	"clinic-vet-api/internal/shared/page"
	"clinic-vet-api/database/sqlc"
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"

	customErr "clinic-vet-api/internal/shared/errors"
)

type MedicalSessionSqlcRepository struct {
	queries *sqlc.Queries
	pgMap   *mapper.SqlcFieldMapper
}

func NewMedicalSessionSqlcRepository(queries *sqlc.Queries) medical.MedicalSessionRepository {
	return &MedicalSessionSqlcRepository{
		queries: queries,
		pgMap:   mapper.NewSqlcFieldMapper(),
	}
}

func (r *MedicalSessionSqlcRepository) FindBySpecification(ctx context.Context, spec medical.MedicalSessionSpecification) (page.Page[medical.MedicalSession], error) {
	ids := make([]int32, 0, len(spec.IDs))
	for _, id := range spec.IDs {
		ids = append(ids, id.Int32())
	}
	petIds := make([]int32, 0, len(spec.PetIDs))
	for _, p := range spec.PetIDs {
		petIds = append(petIds, int32(p))
	}
	customerIds := make([]int32, 0, len(spec.CustomerIDs))
	for _, c := range spec.CustomerIDs {
		customerIds = append(customerIds, int32(c))
	}
	employeeIds := make([]int32, 0, len(spec.EmployeeIDs))
	for _, e := range spec.EmployeeIDs {
		employeeIds = append(employeeIds, int32(e))
	}
	clinicServices := make([]string, 0, len(spec.ClinicServices))
	for _, s := range spec.ClinicServices {
		clinicServices = append(clinicServices, string(s))
	}

	var appointmentID pgtype.Int4
	if len(spec.AppointmentIDs) > 0 {
		appointmentID = pgtype.Int4{Int32: int32(spec.AppointmentIDs[0]), Valid: true}
	}
	var isEmergency, isDeleted pgtype.Bool
	if spec.IsEmergency != nil {
		isEmergency = pgtype.Bool{Bool: *spec.IsEmergency, Valid: true}
	}
	if spec.IsDeleted != nil {
		isDeleted = pgtype.Bool{Bool: *spec.IsDeleted, Valid: true}
	}
	var visitDateFrom, visitDateTo, followUpFrom, followUpTo pgtype.Timestamptz
	if spec.VisitDateFrom != nil {
		visitDateFrom = r.pgMap.PgTimestamptz.FromTime(*spec.VisitDateFrom)
	}
	if spec.VisitDateTo != nil {
		visitDateTo = r.pgMap.PgTimestamptz.FromTime(*spec.VisitDateTo)
	}
	if spec.FollowUpFrom != nil {
		followUpFrom = r.pgMap.PgTimestamptz.FromTime(*spec.FollowUpFrom)
	}
	if spec.FollowUpTo != nil {
		followUpTo = r.pgMap.PgTimestamptz.FromTime(*spec.FollowUpTo)
	}

	limit := spec.Pagination.Limit()
	if limit <= 0 {
		limit = int32(page.DefaultPageSize)
	}
	offset := spec.Pagination.Offset()

	params := sqlc.FindMedicalSessionsBySpecParams{
		Ids:            ids,
		PetIds:         petIds,
		CustomerIds:    customerIds,
		EmployeeIds:    employeeIds,
		AppointmentID:  appointmentID,
		ClinicServices: clinicServices,
		IsEmergency:    isEmergency,
		IsDeleted:      isDeleted,
		VisitDateFrom:  visitDateFrom,
		VisitDateTo:    visitDateTo,
		FollowUpFrom:   followUpFrom,
		FollowUpTo:     followUpTo,
		PageLimit:      limit,
		PageOffset:     offset,
	}
	rows, err := r.queries.FindMedicalSessionsBySpec(ctx, params)
	if err != nil {
		return page.Page[medical.MedicalSession]{}, r.dbError(OpSelect, ErrMsgFindMedicalSessionsBySpec, err)
	}
	countParams := sqlc.CountMedicalSessionsBySpecParams{
		Ids:            ids,
		PetIds:         petIds,
		CustomerIds:    customerIds,
		EmployeeIds:    employeeIds,
		AppointmentID:  appointmentID,
		ClinicServices: clinicServices,
		IsEmergency:    isEmergency,
		IsDeleted:      isDeleted,
		VisitDateFrom:  visitDateFrom,
		VisitDateTo:    visitDateTo,
		FollowUpFrom:   followUpFrom,
		FollowUpTo:     followUpTo,
	}
	total, err := r.queries.CountMedicalSessionsBySpec(ctx, countParams)
	if err != nil {
		return page.Page[medical.MedicalSession]{}, r.dbError(OpCount, ErrMsgFindMedicalSessionsBySpec, err)
	}
	items := make([]medical.MedicalSession, len(rows))
	for i := range rows {
		items[i] = r.toSessionEntity(rows[i])
	}
	pagReq := page.PaginationRequest{
		Page:     int32(spec.Pagination.Number),
		PageSize: limit,
	}
	return page.NewPage(items, total, pagReq), nil
}

func (r *MedicalSessionSqlcRepository) FindByID(ctx context.Context, id medical.SessionID) (medical.MedicalSession, error) {
	row, err := r.queries.FindMedicalSessionByID(ctx, id.Int32())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return medical.MedicalSession{}, r.notFoundError("id", id.String())
		}
		return medical.MedicalSession{}, r.dbError(OpSelect, ErrMsgFindMedicalSessionByID, err)
	}
	return r.toSessionEntity(row), nil
}

func (r *MedicalSessionSqlcRepository) FindByPetID(ctx context.Context, petID uint, p page.Pagination) (page.Page[medical.MedicalSession], error) {
	limit := p.Limit()
	if limit <= 0 {
		limit = int32(page.DefaultPageSize)
	}
	offset := p.Offset()
	params := sqlc.FindMedicalSessionByPetIDParams{
		PetID:  int32(petID),
		Limit:  limit,
		Offset: offset,
	}
	rows, err := r.queries.FindMedicalSessionByPetID(ctx, params)
	if err != nil {
		return page.Page[medical.MedicalSession]{}, r.dbError(OpSelect, ErrMsgFindMedicalSessionsBySpec, err)
	}
	total, err := r.queries.CountMedicalSessionByPetID(ctx, int32(petID))
	if err != nil {
		return page.Page[medical.MedicalSession]{}, r.dbError(OpCount, ErrMsgCountMedicalSessions, err)
	}
	items := make([]medical.MedicalSession, len(rows))
	for i := range rows {
		items[i] = r.toSessionEntity(rows[i])
	}
	pagReq := page.PaginationRequest{Page: int32(p.Number), PageSize: limit}
	return page.NewPage(items, total, pagReq), nil
}

func (r *MedicalSessionSqlcRepository) FindByCustomerID(ctx context.Context, customerID uint, p page.Pagination) (page.Page[medical.MedicalSession], error) {
	limit := p.Limit()
	if limit <= 0 {
		limit = int32(page.DefaultPageSize)
	}
	offset := p.Offset()
	params := sqlc.FindMedicalSessionByCustomerIDParams{
		CustomerID: int32(customerID),
		Limit:      limit,
		Offset:     offset,
	}
	rows, err := r.queries.FindMedicalSessionByCustomerID(ctx, params)
	if err != nil {
		return page.Page[medical.MedicalSession]{}, r.dbError(OpSelect, ErrMsgFindMedicalSessionsBySpec, err)
	}
	total, err := r.queries.CountMedicalSessionByCustomerID(ctx, int32(customerID))
	if err != nil {
		return page.Page[medical.MedicalSession]{}, r.dbError(OpCount, ErrMsgCountMedicalSessions, err)
	}
	items := make([]medical.MedicalSession, len(rows))
	for i := range rows {
		items[i] = r.toSessionEntity(rows[i])
	}
	pagReq := page.PaginationRequest{Page: int32(p.Number), PageSize: limit}
	return page.NewPage(items, total, pagReq), nil
}

func (r *MedicalSessionSqlcRepository) ExistsByID(ctx context.Context, id medical.SessionID) (bool, error) {
	exists, err := r.queries.ExistsMedicalSessionByID(ctx, id.Int32())
	if err != nil {
		return false, r.dbError(OpSelect, ErrMsgCheckMedicalSessionExists, err)
	}
	return exists, nil
}

func (r *MedicalSessionSqlcRepository) Save(ctx context.Context, session medical.MedicalSession) (medical.MedicalSession, error) {
	if session.ID.Value() == 0 {
		return r.createSession(ctx, session)
	}
	return r.updateSession(ctx, session)
}

func (r *MedicalSessionSqlcRepository) SoftDelete(ctx context.Context, id medical.SessionID) error {
	if err := r.queries.SoftDeleteMedicalSession(ctx, id.Int32()); err != nil {
		return r.dbError(OpDelete, ErrMsgSoftDeleteMedicalSession, err)
	}
	return nil
}

func (r *MedicalSessionSqlcRepository) HardDelete(ctx context.Context, id medical.SessionID) error {
	if err := r.queries.HardDeleteMedicalSession(ctx, id.Int32()); err != nil {
		return r.dbError(OpDelete, ErrMsgHardDeleteMedicalSession, err)
	}
	return nil
}

func (r *MedicalSessionSqlcRepository) RestoreByID(ctx context.Context, id medical.SessionID) error {
	if err := r.queries.RestoreMedicalSessionByID(ctx, id.Int32()); err != nil {
		return r.dbError(OpUpdate, ErrMsgRestoreMedicalSession, err)
	}
	return nil
}

func (r *MedicalSessionSqlcRepository) IsDeletedByID(ctx context.Context, id medical.SessionID) (bool, error) {
	val, err := r.queries.IsDeletedMedicalSessionByID(ctx, id.Int32())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, r.notFoundError("id", id.String())
		}
		return false, r.dbError(OpSelect, ErrMsgCheckMedicalSessionExists, err)
	}
	deleted, _ := val.(bool)
	return deleted, nil
}

func (r *MedicalSessionSqlcRepository) CountAll(ctx context.Context) (int64, error) {
	count, err := r.queries.CountAllMedicalSessionsIncludingDeleted(ctx)
	if err != nil {
		return 0, r.dbError(OpCount, ErrMsgCountMedicalSessions, err)
	}
	return count, nil
}

func (r *MedicalSessionSqlcRepository) CountActive(ctx context.Context) (int64, error) {
	count, err := r.queries.CountAllMedicalSession(ctx)
	if err != nil {
		return 0, r.dbError(OpCount, ErrMsgCountMedicalSessions, err)
	}
	return count, nil
}

func (r *MedicalSessionSqlcRepository) CountByClinicService(ctx context.Context) (map[medical.ClinicService]int64, error) {
	rows, err := r.queries.CountMedicalSessionsByClinicService(ctx)
	if err != nil {
		return nil, r.dbError(OpCount, ErrMsgCountMedicalSessions, err)
	}
	out := make(map[medical.ClinicService]int64, len(rows))
	for _, row := range rows {
		out[medical.ClinicService(row.ClinicService)] = row.Count
	}
	return out, nil
}

func (r *MedicalSessionSqlcRepository) createSession(ctx context.Context, s medical.MedicalSession) (medical.MedicalSession, error) {
	params := sqlc.SaveMedicalSessionParams{
		PetID:         int32(s.PetID),
		CustomerID:    int32(s.CustomerID),
		EmployeeID:    int32(s.EmployeeID),
		VisitDate:     r.pgMap.PgTimestamptz.FromTime(s.VisitDate),
		VisitType:     s.VisitType,
		ClinicService: models.ClinicService(s.ClinicService),
		Diagnosis:     r.pgMap.PgText.FromStringPtr(s.Diagnosis),
		Treatment:     r.pgMap.PgText.FromStringPtr(s.Treatment),
		Notes:         r.pgMap.PgText.FromStringPtr(s.Notes),
		Condition:     r.pgMap.PgText.FromStringPtr(s.Condition),
		Symptoms:      r.pgMap.PgText.FromStringPtr(s.Symptoms),
		Medications:   r.pgMap.PgText.FromStringPtr(s.Medications),
		FollowUpDate:  r.pgMap.PgTimestamptz.FromTimePtr(s.FollowUpDate),
		IsEmergency:   pgtype.Bool{Bool: s.IsEmergency, Valid: true},
	}
	if s.AppointmentID != nil {
		params.AppointmentID = pgtype.Int4{Int32: int32(*s.AppointmentID), Valid: true}
	}
	if s.Vitals.WeightKg != nil {
		params.Weight = r.pgMap.PgNumeric.FromDecimal(shared.NewDecimalFromFloat(*s.Vitals.WeightKg))
	}
	if s.Vitals.TemperatureC != nil {
		params.Temperature = r.pgMap.PgNumeric.FromDecimal(shared.NewDecimalFromFloat(*s.Vitals.TemperatureC))
	}
	if s.Vitals.HeartRate != nil {
		params.HeartRate = pgtype.Int4{Int32: int32(*s.Vitals.HeartRate), Valid: true}
	}
	if s.Vitals.RespiratoryRate != nil {
		params.RespiratoryRate = pgtype.Int4{Int32: int32(*s.Vitals.RespiratoryRate), Valid: true}
	}
	created, err := r.queries.SaveMedicalSession(ctx, params)
	if err != nil {
		return medical.MedicalSession{}, r.dbError(OpInsert, ErrMsgSaveMedicalSession, err)
	}
	return r.toSessionEntity(created), nil
}

func (r *MedicalSessionSqlcRepository) updateSession(ctx context.Context, s medical.MedicalSession) (medical.MedicalSession, error) {
	params := sqlc.UpdateMedicalSessionParams{
		ID:            int32(s.ID.Value()),
		PetID:         int32(s.PetID),
		CustomerID:    int32(s.CustomerID),
		EmployeeID:    int32(s.EmployeeID),
		VisitDate:     r.pgMap.PgTimestamptz.FromTime(s.VisitDate),
		VisitType:     s.VisitType,
		ClinicService: models.ClinicService(s.ClinicService),
		Diagnosis:     r.pgMap.PgText.FromStringPtr(s.Diagnosis),
		Treatment:     r.pgMap.PgText.FromStringPtr(s.Treatment),
		Notes:         r.pgMap.PgText.FromStringPtr(s.Notes),
		Condition:     r.pgMap.PgText.FromStringPtr(s.Condition),
		Symptoms:      r.pgMap.PgText.FromStringPtr(s.Symptoms),
		Medications:   r.pgMap.PgText.FromStringPtr(s.Medications),
		FollowUpDate:  r.pgMap.PgTimestamptz.FromTimePtr(s.FollowUpDate),
		IsEmergency:   pgtype.Bool{Bool: s.IsEmergency, Valid: true},
	}
	if s.AppointmentID != nil {
		params.AppointmentID = pgtype.Int4{Int32: int32(*s.AppointmentID), Valid: true}
	}
	if s.Vitals.WeightKg != nil {
		params.Weight = r.pgMap.PgNumeric.FromDecimal(shared.NewDecimalFromFloat(*s.Vitals.WeightKg))
	}
	if s.Vitals.TemperatureC != nil {
		params.Temperature = r.pgMap.PgNumeric.FromDecimal(shared.NewDecimalFromFloat(*s.Vitals.TemperatureC))
	}
	if s.Vitals.HeartRate != nil {
		params.HeartRate = pgtype.Int4{Int32: int32(*s.Vitals.HeartRate), Valid: true}
	}
	if s.Vitals.RespiratoryRate != nil {
		params.RespiratoryRate = pgtype.Int4{Int32: int32(*s.Vitals.RespiratoryRate), Valid: true}
	}
	updated, err := r.queries.UpdateMedicalSession(ctx, params)
	if err != nil {
		return medical.MedicalSession{}, r.dbError(OpUpdate, ErrMsgUpdateMedicalSession, err)
	}
	return r.toSessionEntity(updated), nil
}

func (r *MedicalSessionSqlcRepository) toSessionEntity(row sqlc.MedicalSession) medical.MedicalSession {
	s := medical.MedicalSession{
		ID:            medical.NewSessionID(uint(row.ID)),
		PetID:         uint(row.PetID),
		CustomerID:    uint(row.CustomerID),
		EmployeeID:    uint(row.EmployeeID),
		ClinicService: medical.ClinicService(row.ClinicService),
		VisitType:     row.VisitType,
		VisitDate:     r.pgMap.PgTimestamptz.ToTime(row.VisitDate),
		IsEmergency:   row.IsEmergency.Bool,
		Diagnosis:     r.pgMap.PgText.ToStringPtr(row.Diagnosis),
		Treatment:     r.pgMap.PgText.ToStringPtr(row.Treatment),
		Notes:         r.pgMap.PgText.ToStringPtr(row.Notes),
		Condition:     r.pgMap.PgText.ToStringPtr(row.Condition),
		Symptoms:      r.pgMap.PgText.ToStringPtr(row.Symptoms),
		Medications:   r.pgMap.PgText.ToStringPtr(row.Medications),
		FollowUpDate:  r.pgMap.PgTimestamptz.ToTimePtr(row.FollowUpDate),
		CreatedAt:     r.pgMap.PgTimestamptz.ToTime(row.CreatedAt),
		UpdatedAt:     r.pgMap.PgTimestamptz.ToTime(row.UpdatedAt),
		DeletedAt:     r.pgMap.PgTimestamptz.ToTimePtr(row.DeletedAt),
	}
	if row.AppointmentID.Valid {
		apt := uint(row.AppointmentID.Int32)
		s.AppointmentID = &apt
	}
	if row.Weight.Valid {
		f := r.pgMap.PgNumeric.ToDecimal(row.Weight).Float64()
		s.Vitals.WeightKg = &f
	}
	if row.Temperature.Valid {
		f := r.pgMap.PgNumeric.ToDecimal(row.Temperature).Float64()
		s.Vitals.TemperatureC = &f
	}
	if row.HeartRate.Valid {
		hr := int(row.HeartRate.Int32)
		s.Vitals.HeartRate = &hr
	}
	if row.RespiratoryRate.Valid {
		rr := int(row.RespiratoryRate.Int32)
		s.Vitals.RespiratoryRate = &rr
	}
	return s
}

func (r *MedicalSessionSqlcRepository) dbError(operation, message string, err error) error {
	return customErr.DatabaseError(operation, TableMedicalSessions, DriverSQL, fmt.Errorf("%s: %w", message, err))
}

func (r *MedicalSessionSqlcRepository) notFoundError(parameterName, parameterValue string) error {
	return customErr.DBNotFoundError(parameterName, parameterValue, OpSelect, TableMedicalSessions, DriverSQL)
}
