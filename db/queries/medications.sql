-- name: FindMedicationByID :one
SELECT * FROM medications WHERE id = $1;

-- name: FindMedicationsAll :many
SELECT * FROM medications WHERE is_active = TRUE ORDER BY name LIMIT $1 OFFSET $2;

-- name: SearchMedications :many
SELECT * FROM medications WHERE is_active = TRUE AND (name ILIKE '%' || $1 || '%' OR active_ingredient ILIKE '%' || $1 || '%') ORDER BY name LIMIT $2 OFFSET $3;

-- name: CountMedicationsAll :one
SELECT COUNT(*) FROM medications WHERE is_active = TRUE;

-- name: CountMedicationsSearch :one
SELECT COUNT(*) FROM medications WHERE is_active = TRUE AND (name ILIKE '%' || $1 || '%' OR active_ingredient ILIKE '%' || $1 || '%');

-- name: CreateMedication :one
INSERT INTO medications (name, active_ingredient, presentation, unit, requires_prescription, species_warnings, is_active)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: UpdateMedication :one
UPDATE medications
SET name = $2, active_ingredient = $3, presentation = $4, unit = $5, requires_prescription = $6, species_warnings = $7, is_active = $8
WHERE id = $1
RETURNING *;

-- name: DeleteMedicationByID :exec
DELETE FROM medications WHERE id = $1;
