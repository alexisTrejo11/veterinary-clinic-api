package repository

import (
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

const (
	TableVaccineCatalog = "vaccine_catalog"
	TableMedications    = "medications"
	TableServiceCatalog = "service_catalog"
)

// ─── VaccineCatalogRepository ───────────────────────────────────────────────

type VaccineCatalogSqlcRepository struct {
	queries *sqlc.Queries
	pgMap   *mapper.SqlcFieldMapper
}

func NewVaccineCatalogSqlcRepository(queries *sqlc.Queries) medical.VaccineCatalogRepository {
	return &VaccineCatalogSqlcRepository{queries: queries, pgMap: mapper.NewSqlcFieldMapper()}
}

func (r *VaccineCatalogSqlcRepository) FindAll(ctx context.Context, p page.Pagination) (page.Page[medical.VaccineCatalog], error) {
	limit := p.Limit()
	if limit <= 0 {
		limit = int32(page.DefaultPageSize)
	}
	params := sqlc.FindVaccineCatalogAllParams{Limit: limit, Offset: p.Offset()}
	rows, err := r.queries.FindVaccineCatalogAll(ctx, params)
	if err != nil {
		return page.Page[medical.VaccineCatalog]{}, customErr.DatabaseError(OpSelect, TableVaccineCatalog, DriverSQL, err)
	}
	total, err := r.queries.CountVaccineCatalogAll(ctx)
	if err != nil {
		return page.Page[medical.VaccineCatalog]{}, customErr.DatabaseError(OpCount, TableVaccineCatalog, DriverSQL, err)
	}
	items := make([]medical.VaccineCatalog, len(rows))
	for i := range rows {
		items[i] = r.vaccineCatalogToEntity(rows[i])
	}
	return page.NewPage(items, total, page.PaginationRequest{Page: int32(p.Number), PageSize: limit}), nil
}

func (r *VaccineCatalogSqlcRepository) FindByID(
	ctx context.Context,
	id medical.VaccineCatalogID,
) (medical.VaccineCatalog, error) {
	row, err := r.queries.FindVaccineCatalogByID(ctx, id.Int32())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return medical.VaccineCatalog{}, customErr.DBNotFoundError("id", fmt.Sprintf("%d", id.Value()), OpSelect, TableVaccineCatalog, DriverSQL)
		}
		return medical.VaccineCatalog{}, customErr.DatabaseError(OpSelect, TableVaccineCatalog, DriverSQL, err)
	}
	return r.vaccineCatalogToEntity(row), nil
}

func (r *VaccineCatalogSqlcRepository) FindBySpecies(ctx context.Context, species string, p page.Pagination) (page.Page[medical.VaccineCatalog], error) {
	limit := p.Limit()
	if limit <= 0 {
		limit = int32(page.DefaultPageSize)
	}
	speciesParam := r.pgMap.PgText.FromString(species)
	params := sqlc.FindVaccineCatalogBySpeciesParams{Species: speciesParam, Limit: limit, Offset: p.Offset()}
	rows, err := r.queries.FindVaccineCatalogBySpecies(ctx, params)
	if err != nil {
		return page.Page[medical.VaccineCatalog]{}, customErr.DatabaseError(OpSelect, TableVaccineCatalog, DriverSQL, err)
	}
	total, err := r.queries.CountVaccineCatalogBySpecies(ctx, speciesParam)
	if err != nil {
		return page.Page[medical.VaccineCatalog]{}, customErr.DatabaseError(OpCount, TableVaccineCatalog, DriverSQL, err)
	}
	items := make([]medical.VaccineCatalog, len(rows))
	for i := range rows {
		items[i] = r.vaccineCatalogToEntity(rows[i])
	}
	return page.NewPage(items, total, page.PaginationRequest{Page: int32(p.Number), PageSize: limit}), nil
}

func (r *VaccineCatalogSqlcRepository) Save(ctx context.Context, v medical.VaccineCatalog) (medical.VaccineCatalog, error) {
	if v.ID.Value() == 0 {
		return r.createVaccineCatalog(ctx, v)
	}
	return r.updateVaccineCatalog(ctx, v)
}

func (r *VaccineCatalogSqlcRepository) DeleteByID(ctx context.Context, id medical.VaccineCatalogID) error {
	return r.queries.DeleteVaccineCatalogByID(ctx, int32(id.Value()))
}

func (r *VaccineCatalogSqlcRepository) createVaccineCatalog(ctx context.Context, v medical.VaccineCatalog) (medical.VaccineCatalog, error) {
	scheduleDays := make([]int32, 0, len(v.ScheduleDays))
	for _, d := range v.ScheduleDays {
		scheduleDays = append(scheduleDays, int32(d))
	}
	params := sqlc.CreateVaccineCatalogParams{
		Name:          v.Name,
		Manufacturer:  r.pgMap.PgText.FromStringPtr(v.Manufacturer),
		Species:       r.pgMap.PgText.FromStringPtr(v.Species),
		DiseaseTarget: r.pgMap.PgText.FromStringPtr(v.DiseaseTarget),
		TotalDoses:    pgtype.Int4{Int32: int32(v.TotalDoses), Valid: true},
		ScheduleDays:  scheduleDays,
		Notes:         r.pgMap.PgText.FromStringPtr(v.Notes),
		IsActive:      pgtype.Bool{Bool: v.IsActive, Valid: true},
	}
	created, err := r.queries.CreateVaccineCatalog(ctx, params)
	if err != nil {
		return medical.VaccineCatalog{}, customErr.DatabaseError(OpInsert, TableVaccineCatalog, DriverSQL, err)
	}
	return r.vaccineCatalogToEntity(created), nil
}

func (r *VaccineCatalogSqlcRepository) updateVaccineCatalog(ctx context.Context, v medical.VaccineCatalog) (medical.VaccineCatalog, error) {
	scheduleDays := make([]int32, 0, len(v.ScheduleDays))
	for _, d := range v.ScheduleDays {
		scheduleDays = append(scheduleDays, int32(d))
	}
	params := sqlc.UpdateVaccineCatalogParams{
		ID:            int32(v.ID.Value()),
		Name:          v.Name,
		Manufacturer:  r.pgMap.PgText.FromStringPtr(v.Manufacturer),
		Species:       r.pgMap.PgText.FromStringPtr(v.Species),
		DiseaseTarget: r.pgMap.PgText.FromStringPtr(v.DiseaseTarget),
		TotalDoses:    pgtype.Int4{Int32: int32(v.TotalDoses), Valid: true},
		ScheduleDays:  scheduleDays,
		Notes:         r.pgMap.PgText.FromStringPtr(v.Notes),
		IsActive:      pgtype.Bool{Bool: v.IsActive, Valid: true},
	}
	updated, err := r.queries.UpdateVaccineCatalog(ctx, params)
	if err != nil {
		return medical.VaccineCatalog{}, customErr.DatabaseError(OpUpdate, TableVaccineCatalog, DriverSQL, err)
	}
	return r.vaccineCatalogToEntity(updated), nil
}

func (r *VaccineCatalogSqlcRepository) vaccineCatalogToEntity(row sqlc.VaccineCatalog) medical.VaccineCatalog {
	v := medical.VaccineCatalog{
		ID:            medical.NewVaccineCatalogID(uint(row.ID)),
		Name:          row.Name,
		Manufacturer:  r.pgMap.PgText.ToStringPtr(row.Manufacturer),
		Species:       r.pgMap.PgText.ToStringPtr(row.Species),
		DiseaseTarget: r.pgMap.PgText.ToStringPtr(row.DiseaseTarget),
		Notes:         r.pgMap.PgText.ToStringPtr(row.Notes),
		CreatedAt:     r.pgMap.PgTimestamptz.ToTime(row.CreatedAt),
	}
	if row.TotalDoses.Valid {
		v.TotalDoses = int(row.TotalDoses.Int32)
	}
	if len(row.ScheduleDays) > 0 {
		v.ScheduleDays = make([]int, len(row.ScheduleDays))
		for i, d := range row.ScheduleDays {
			v.ScheduleDays[i] = int(d)
		}
	}
	if row.IsActive.Valid {
		v.IsActive = row.IsActive.Bool
	}
	return v
}

// ─── MedicationRepository ───────────────────────────────────────────────────

type MedicationSqlcRepository struct {
	queries *sqlc.Queries
	pgMap   *mapper.SqlcFieldMapper
}

func NewMedicationSqlcRepository(queries *sqlc.Queries) medical.MedicationRepository {
	return &MedicationSqlcRepository{queries: queries, pgMap: mapper.NewSqlcFieldMapper()}
}

func (r *MedicationSqlcRepository) FindAll(ctx context.Context, p page.Pagination) (page.Page[medical.Medication], error) {
	limit := p.Limit()
	if limit <= 0 {
		limit = int32(page.DefaultPageSize)
	}
	params := sqlc.FindMedicationsAllParams{Limit: limit, Offset: p.Offset()}
	rows, err := r.queries.FindMedicationsAll(ctx, params)
	if err != nil {
		return page.Page[medical.Medication]{}, customErr.DatabaseError(OpSelect, TableMedications, DriverSQL, err)
	}
	total, err := r.queries.CountMedicationsAll(ctx)
	if err != nil {
		return page.Page[medical.Medication]{}, customErr.DatabaseError(OpCount, TableMedications, DriverSQL, err)
	}
	items := make([]medical.Medication, len(rows))
	for i := range rows {
		items[i] = r.medicationToEntity(rows[i])
	}
	return page.NewPage(items, total, page.PaginationRequest{Page: int32(p.Number), PageSize: limit}), nil
}

func (r *MedicationSqlcRepository) FindByID(ctx context.Context, id medical.MedicationID) (medical.Medication, error) {
	row, err := r.queries.FindMedicationByID(ctx, int32(id.Value()))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return medical.Medication{}, customErr.DBNotFoundError("id", fmt.Sprintf("%d", id.Value()), OpSelect, TableMedications, DriverSQL)
		}
		return medical.Medication{}, customErr.DatabaseError(OpSelect, TableMedications, DriverSQL, err)
	}
	return r.medicationToEntity(row), nil
}

func (r *MedicationSqlcRepository) Search(ctx context.Context, term string, p page.Pagination) (page.Page[medical.Medication], error) {
	limit := p.Limit()
	if limit <= 0 {
		limit = int32(page.DefaultPageSize)
	}
	termParam := r.pgMap.PgText.FromString(term)
	params := sqlc.SearchMedicationsParams{
		Column1: termParam,
		Limit:   limit,
		Offset:  p.Offset(),
	}
	rows, err := r.queries.SearchMedications(ctx, params)
	if err != nil {
		return page.Page[medical.Medication]{}, customErr.DatabaseError(OpSelect, TableMedications, DriverSQL, err)
	}
	total, err := r.queries.CountMedicationsSearch(ctx, termParam)
	if err != nil {
		return page.Page[medical.Medication]{}, customErr.DatabaseError(OpCount, TableMedications, DriverSQL, err)
	}
	items := make([]medical.Medication, len(rows))
	for i := range rows {
		items[i] = r.medicationToEntity(rows[i])
	}
	return page.NewPage(items, total, page.PaginationRequest{Page: int32(p.Number), PageSize: limit}), nil
}

func (r *MedicationSqlcRepository) Save(ctx context.Context, m medical.Medication) (medical.Medication, error) {
	if m.ID.Value() == 0 {
		return r.createMedication(ctx, m)
	}
	return r.updateMedication(ctx, m)
}

func (r *MedicationSqlcRepository) DeleteByID(ctx context.Context, id medical.MedicationID) error {
	return r.queries.DeleteMedicationByID(ctx, int32(id.Value()))
}

func (r *MedicationSqlcRepository) createMedication(ctx context.Context, m medical.Medication) (medical.Medication, error) {
	params := sqlc.CreateMedicationParams{
		Name:                 m.Name,
		ActiveIngredient:     r.pgMap.PgText.FromStringPtr(m.ActiveIngredient),
		Presentation:         r.pgMap.PgText.FromStringPtr(m.Presentation),
		Unit:                 r.pgMap.PgText.FromStringPtr(m.Unit),
		RequiresPrescription: pgtype.Bool{Bool: m.RequiresPrescription, Valid: true},
		SpeciesWarnings:      r.pgMap.PgText.FromStringPtr(m.SpeciesWarnings),
		IsActive:             pgtype.Bool{Bool: m.IsActive, Valid: true},
	}
	created, err := r.queries.CreateMedication(ctx, params)
	if err != nil {
		return medical.Medication{}, customErr.DatabaseError(OpInsert, TableMedications, DriverSQL, err)
	}
	return r.medicationToEntity(created), nil
}

func (r *MedicationSqlcRepository) updateMedication(ctx context.Context, m medical.Medication) (medical.Medication, error) {
	params := sqlc.UpdateMedicationParams{
		ID:                   int32(m.ID.Value()),
		Name:                 m.Name,
		ActiveIngredient:     r.pgMap.PgText.FromStringPtr(m.ActiveIngredient),
		Presentation:         r.pgMap.PgText.FromStringPtr(m.Presentation),
		Unit:                 r.pgMap.PgText.FromStringPtr(m.Unit),
		RequiresPrescription: pgtype.Bool{Bool: m.RequiresPrescription, Valid: true},
		SpeciesWarnings:      r.pgMap.PgText.FromStringPtr(m.SpeciesWarnings),
		IsActive:             pgtype.Bool{Bool: m.IsActive, Valid: true},
	}
	updated, err := r.queries.UpdateMedication(ctx, params)
	if err != nil {
		return medical.Medication{}, customErr.DatabaseError(OpUpdate, TableMedications, DriverSQL, err)
	}
	return r.medicationToEntity(updated), nil
}

func (r *MedicationSqlcRepository) medicationToEntity(row sqlc.Medication) medical.Medication {
	m := medical.Medication{
		ID:               medical.NewMedicationID(uint(row.ID)),
		Name:             row.Name,
		ActiveIngredient: r.pgMap.PgText.ToStringPtr(row.ActiveIngredient),
		Presentation:     r.pgMap.PgText.ToStringPtr(row.Presentation),
		Unit:             r.pgMap.PgText.ToStringPtr(row.Unit),
		SpeciesWarnings:  r.pgMap.PgText.ToStringPtr(row.SpeciesWarnings),
		CreatedAt:        r.pgMap.PgTimestamptz.ToTime(row.CreatedAt),
	}
	if row.RequiresPrescription.Valid {
		m.RequiresPrescription = row.RequiresPrescription.Bool
	}
	if row.IsActive.Valid {
		m.IsActive = row.IsActive.Bool
	}
	return m
}

// ─── ServiceCatalogRepository ───────────────────────────────────────────────

type ServiceCatalogSqlcRepository struct {
	queries *sqlc.Queries
	pgMap   *mapper.SqlcFieldMapper
}

func NewServiceCatalogSqlcRepository(queries *sqlc.Queries) medical.ServiceCatalogRepository {
	return &ServiceCatalogSqlcRepository{queries: queries, pgMap: mapper.NewSqlcFieldMapper()}
}

func (r *ServiceCatalogSqlcRepository) FindAll(ctx context.Context, p page.Pagination) (page.Page[medical.ServiceCatalog], error) {
	limit := p.Limit()
	if limit <= 0 {
		limit = int32(page.DefaultPageSize)
	}
	params := sqlc.FindServiceCatalogAllParams{Limit: limit, Offset: p.Offset()}
	rows, err := r.queries.FindServiceCatalogAll(ctx, params)
	if err != nil {
		return page.Page[medical.ServiceCatalog]{}, customErr.DatabaseError(OpSelect, TableServiceCatalog, DriverSQL, err)
	}
	total, err := r.queries.CountServiceCatalogAll(ctx)
	if err != nil {
		return page.Page[medical.ServiceCatalog]{}, customErr.DatabaseError(OpCount, TableServiceCatalog, DriverSQL, err)
	}
	items := make([]medical.ServiceCatalog, len(rows))
	for i := range rows {
		items[i] = r.serviceCatalogToEntity(rows[i])
	}
	return page.NewPage(items, total, page.PaginationRequest{Page: int32(p.Number), PageSize: limit}), nil
}

func (r *ServiceCatalogSqlcRepository) FindByID(ctx context.Context, id medical.ServiceCatalogID) (medical.ServiceCatalog, error) {
	row, err := r.queries.FindServiceCatalogByID(ctx, int32(id.Value()))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return medical.ServiceCatalog{}, customErr.DBNotFoundError("id", fmt.Sprintf("%d", id.Value()), OpSelect, TableServiceCatalog, DriverSQL)
		}
		return medical.ServiceCatalog{}, customErr.DatabaseError(OpSelect, TableServiceCatalog, DriverSQL, err)
	}
	return r.serviceCatalogToEntity(row), nil
}

func (r *ServiceCatalogSqlcRepository) FindByCategory(ctx context.Context, cat medical.ServiceCategory, p page.Pagination) (page.Page[medical.ServiceCatalog], error) {
	limit := p.Limit()
	if limit <= 0 {
		limit = int32(page.DefaultPageSize)
	}
	params := sqlc.FindServiceCatalogByCategoryParams{Category: string(cat), Limit: limit, Offset: p.Offset()}
	rows, err := r.queries.FindServiceCatalogByCategory(ctx, params)
	if err != nil {
		return page.Page[medical.ServiceCatalog]{}, customErr.DatabaseError(OpSelect, TableServiceCatalog, DriverSQL, err)
	}
	total, err := r.queries.CountServiceCatalogByCategory(ctx, string(cat))
	if err != nil {
		return page.Page[medical.ServiceCatalog]{}, customErr.DatabaseError(OpCount, TableServiceCatalog, DriverSQL, err)
	}
	items := make([]medical.ServiceCatalog, len(rows))
	for i := range rows {
		items[i] = r.serviceCatalogToEntity(rows[i])
	}
	return page.NewPage(items, total, page.PaginationRequest{Page: int32(p.Number), PageSize: limit}), nil
}

func (r *ServiceCatalogSqlcRepository) Save(ctx context.Context, s medical.ServiceCatalog) (medical.ServiceCatalog, error) {
	if s.ID.Value() == 0 {
		return r.createServiceCatalog(ctx, s)
	}
	return r.updateServiceCatalog(ctx, s)
}

func (r *ServiceCatalogSqlcRepository) DeleteByID(ctx context.Context, id medical.ServiceCatalogID) error {
	return r.queries.DeleteServiceCatalogByID(ctx, int32(id.Value()))
}

func (r *ServiceCatalogSqlcRepository) createServiceCatalog(ctx context.Context, s medical.ServiceCatalog) (medical.ServiceCatalog, error) {
	params := sqlc.CreateServiceCatalogParams{
		Name:            s.Name,
		Category:        string(s.Category),
		Description:     r.pgMap.PgText.FromStringPtr(s.Description),
		RequiresFasting: pgtype.Bool{Bool: s.RequiresFasting, Valid: true},
		IsActive:        pgtype.Bool{Bool: s.IsActive, Valid: true},
	}
	if s.BasePrice != nil {
		params.BasePrice = r.pgMap.PgNumeric.FromDecimal(shared.NewDecimalFromFloat(*s.BasePrice))
	} else {
		params.BasePrice = pgtype.Numeric{Valid: false}
	}
	if s.DurationMinutes != nil {
		params.DurationMin = pgtype.Int4{Int32: int32(*s.DurationMinutes), Valid: true}
	} else {
		params.DurationMin = pgtype.Int4{Valid: false}
	}
	created, err := r.queries.CreateServiceCatalog(ctx, params)
	if err != nil {
		return medical.ServiceCatalog{}, customErr.DatabaseError(OpInsert, TableServiceCatalog, DriverSQL, err)
	}
	return r.serviceCatalogToEntity(created), nil
}

func (r *ServiceCatalogSqlcRepository) updateServiceCatalog(ctx context.Context, s medical.ServiceCatalog) (medical.ServiceCatalog, error) {
	params := sqlc.UpdateServiceCatalogParams{
		ID:              int32(s.ID.Value()),
		Name:            s.Name,
		Category:        string(s.Category),
		Description:     r.pgMap.PgText.FromStringPtr(s.Description),
		RequiresFasting: pgtype.Bool{Bool: s.RequiresFasting, Valid: true},
		IsActive:        pgtype.Bool{Bool: s.IsActive, Valid: true},
	}
	if s.BasePrice != nil {
		params.BasePrice = r.pgMap.PgNumeric.FromDecimal(shared.NewDecimalFromFloat(*s.BasePrice))
	} else {
		params.BasePrice = pgtype.Numeric{Valid: false}
	}
	if s.DurationMinutes != nil {
		params.DurationMin = pgtype.Int4{Int32: int32(*s.DurationMinutes), Valid: true}
	} else {
		params.DurationMin = pgtype.Int4{Valid: false}
	}
	updated, err := r.queries.UpdateServiceCatalog(ctx, params)
	if err != nil {
		return medical.ServiceCatalog{}, customErr.DatabaseError(OpUpdate, TableServiceCatalog, DriverSQL, err)
	}
	return r.serviceCatalogToEntity(updated), nil
}

func (r *ServiceCatalogSqlcRepository) serviceCatalogToEntity(row sqlc.ServiceCatalog) medical.ServiceCatalog {
	s := medical.ServiceCatalog{
		ID:              medical.NewServiceCatalogID(uint(row.ID)),
		Name:            row.Name,
		Category:        medical.ServiceCategory(row.Category),
		Description:     r.pgMap.PgText.ToStringPtr(row.Description),
		RequiresFasting: row.RequiresFasting.Bool,
		IsActive:        row.IsActive.Bool,
		CreatedAt:       r.pgMap.PgTimestamptz.ToTime(row.CreatedAt),
	}
	if row.BasePrice.Valid {
		f := r.pgMap.PgNumeric.ToDecimal(row.BasePrice).Float64()
		s.BasePrice = &f
	}
	if row.DurationMin.Valid {
		d := int(row.DurationMin.Int32)
		s.DurationMinutes = &d
	}
	return s
}
