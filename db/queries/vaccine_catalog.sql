-- name: FindVaccineCatalogByID :one
SELECT * FROM vaccine_catalog WHERE id = $1;

-- name: FindVaccineCatalogAll :many
SELECT * FROM vaccine_catalog WHERE is_active = TRUE ORDER BY name LIMIT $1 OFFSET $2;

-- name: FindVaccineCatalogBySpecies :many
SELECT * FROM vaccine_catalog WHERE is_active = TRUE AND (species IS NULL OR species = $1) ORDER BY name LIMIT $2 OFFSET $3;

-- name: CountVaccineCatalogAll :one
SELECT COUNT(*) FROM vaccine_catalog WHERE is_active = TRUE;

-- name: CountVaccineCatalogBySpecies :one
SELECT COUNT(*) FROM vaccine_catalog WHERE is_active = TRUE AND (species IS NULL OR species = $1);

-- name: CreateVaccineCatalog :one
INSERT INTO vaccine_catalog (name, manufacturer, species, disease_target, total_doses, schedule_days, notes, is_active)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: UpdateVaccineCatalog :one
UPDATE vaccine_catalog
SET name = $2, manufacturer = $3, species = $4, disease_target = $5, total_doses = $6, schedule_days = $7, notes = $8, is_active = $9
WHERE id = $1
RETURNING *;

-- name: DeleteVaccineCatalogByID :exec
DELETE FROM vaccine_catalog WHERE id = $1;
