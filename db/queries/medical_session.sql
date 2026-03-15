-- name: FindMedicalSessionByID :one
SELECT * FROM medical_sessions
WHERE id = $1 AND deleted_at IS NULL;

-- name: FindMedicalSessionByIDAndPetID :one
SELECT * FROM medical_sessions
WHERE id = $1 AND pet_id = $2 AND deleted_at IS NULL;

-- name: FindMedicalSessionByIDAndCustomerID :one
SELECT * FROM medical_sessions
WHERE id = $1 AND customer_id = $2 AND deleted_at IS NULL;

-- name: FindMedicalSessionByIDAndEmployeeID :one
SELECT * FROM medical_sessions
WHERE id = $1 AND employee_id = $2 AND deleted_at IS NULL;

-- name: FindAllMedicalSession :many
SELECT * FROM medical_sessions
WHERE deleted_at IS NULL
ORDER BY visit_date DESC
LIMIT $1 OFFSET $2;

-- name: CountAllMedicalSession :one
SELECT COUNT(*) FROM medical_sessions
WHERE deleted_at IS NULL;

-- name: FindMedicalSessionByEmployeeID :many
SELECT * FROM medical_sessions
WHERE employee_id = $1 AND deleted_at IS NULL
ORDER BY visit_date DESC
LIMIT $2 OFFSET $3;


-- name: FindMedicalSessionByPetID :many
SELECT * FROM medical_sessions
WHERE pet_id = $1 AND deleted_at IS NULL
ORDER BY visit_date DESC
LIMIT $2 OFFSET $3;


-- name: FindMedicalSessionByCustomerID :many
SELECT * FROM medical_sessions
WHERE customer_id = $1 AND deleted_at IS NULL
ORDER BY visit_date DESC
LIMIT $2 OFFSET $3;


-- name: FindRecentMedicalSessionByPetID :many
SELECT * FROM medical_sessions
WHERE pet_id = $1 AND deleted_at IS NULL
ORDER BY visit_date DESC
LIMIT $2;

-- name: FindMedicalSessionByDateRange :many
SELECT * FROM medical_sessions
WHERE visit_date BETWEEN $1 AND $2
AND deleted_at IS NULL
ORDER BY visit_date DESC
LIMIT $3 OFFSET $4;

-- name: FindMedicalSessionByPetAndDateRange :many
SELECT * FROM medical_sessions
WHERE pet_id = $1
AND visit_date BETWEEN $2 AND $3
AND deleted_at IS NULL
ORDER BY visit_date DESC;

-- name: FindMedicalSessionByDiagnosis :many
SELECT * FROM medical_sessions
WHERE diagnosis ILIKE '%' || $1 || '%'
AND deleted_at IS NULL
ORDER BY visit_date DESC
LIMIT $2 OFFSET $3;

-- name: CountMedicalSessionByDiagnosis :one
SELECT COUNT(*) FROM medical_sessions
WHERE diagnosis ILIKE '%' || $1 || '%'
AND deleted_at IS NULL;

-- name: ExistsMedicalSessionByID :one
SELECT COUNT(*) > 0 FROM medical_sessions
WHERE id = $1 AND deleted_at IS NULL;

-- name: ExistsMedicalSessionByPetAndDate :one
SELECT COUNT(*) > 0 FROM medical_sessions
WHERE pet_id = $1
AND DATE(visit_date) = DATE($2)
AND deleted_at IS NULL;

-- name: SaveMedicalSession :one
INSERT INTO medical_sessions (
    pet_id, customer_id, employee_id, appointment_id,
    visit_date, visit_type, clinic_service,
    diagnosis, treatment, notes, condition,
    weight, temperature, heart_rate, respiratory_rate,
    symptoms, medications, follow_up_date, is_emergency
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19
)
RETURNING *;

-- name: UpdateMedicalSession :one
UPDATE medical_sessions
SET
    pet_id = $2, customer_id = $3, employee_id = $4, appointment_id = $5,
    visit_date = $6, visit_type = $7, clinic_service = $8,
    diagnosis = $9, treatment = $10, notes = $11, condition = $12,
    weight = $13, temperature = $14, heart_rate = $15, respiratory_rate = $16,
    symptoms = $17, medications = $18, follow_up_date = $19, is_emergency = $20,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: SoftDeleteMedicalSession :exec
UPDATE medical_sessions
SET 
    deleted_at = CURRENT_TIMESTAMP,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: HardDeleteMedicalSession :exec
DELETE FROM medical_sessions
WHERE id = $1;

-- name: CountMedicalSessionByPetID :one
SELECT COUNT(*) FROM medical_sessions
WHERE pet_id = $1 AND deleted_at IS NULL;

-- name: CountMedicalSessionByEmployeeID :one
SELECT COUNT(*) FROM medical_sessions
WHERE employee_id = $1 AND deleted_at IS NULL;

-- name: CountMedicalSessionByCustomerID :one
SELECT COUNT(*) FROM medical_sessions
WHERE customer_id = $1 AND deleted_at IS NULL;

-- name: CountMedicalSessionByDateRange :one
SELECT COUNT(*) FROM medical_sessions
WHERE visit_date BETWEEN $1 AND $2
AND deleted_at IS NULL;

-- name: FindMedicalSessionsBySpec :many
SELECT * FROM medical_sessions
WHERE (cardinality(sqlc.arg(ids)::int[]) = 0 OR id = ANY(sqlc.arg(ids)::int[]))
  AND (cardinality(sqlc.arg(pet_ids)::int[]) = 0 OR pet_id = ANY(sqlc.arg(pet_ids)::int[]))
  AND (cardinality(sqlc.arg(customer_ids)::int[]) = 0 OR customer_id = ANY(sqlc.arg(customer_ids)::int[]))
  AND (cardinality(sqlc.arg(employee_ids)::int[]) = 0 OR employee_id = ANY(sqlc.arg(employee_ids)::int[]))
  AND (sqlc.narg(appointment_id)::int IS NULL OR appointment_id = sqlc.narg(appointment_id)::int)
  AND (cardinality(sqlc.arg(clinic_services)::varchar[]) = 0 OR clinic_service = ANY(sqlc.arg(clinic_services)::varchar[]))
  AND (sqlc.narg(is_emergency)::boolean IS NULL OR is_emergency = sqlc.narg(is_emergency)::boolean)
  AND (sqlc.narg(is_deleted)::boolean IS NULL OR (deleted_at IS NOT NULL) = sqlc.narg(is_deleted)::boolean)
  AND (sqlc.narg(visit_date_from)::timestamptz IS NULL OR visit_date >= sqlc.narg(visit_date_from)::timestamptz)
  AND (sqlc.narg(visit_date_to)::timestamptz IS NULL OR visit_date <= sqlc.narg(visit_date_to)::timestamptz)
  AND (sqlc.narg(follow_up_from)::timestamptz IS NULL OR follow_up_date >= sqlc.narg(follow_up_from)::timestamptz)
  AND (sqlc.narg(follow_up_to)::timestamptz IS NULL OR follow_up_date <= sqlc.narg(follow_up_to)::timestamptz)
ORDER BY visit_date DESC
LIMIT sqlc.arg(page_limit) OFFSET sqlc.arg(page_offset);

-- name: CountMedicalSessionsBySpec :one
SELECT COUNT(*) FROM medical_sessions
WHERE (cardinality(sqlc.arg(ids)::int[]) = 0 OR id = ANY(sqlc.arg(ids)::int[]))
  AND (cardinality(sqlc.arg(pet_ids)::int[]) = 0 OR pet_id = ANY(sqlc.arg(pet_ids)::int[]))
  AND (cardinality(sqlc.arg(customer_ids)::int[]) = 0 OR customer_id = ANY(sqlc.arg(customer_ids)::int[]))
  AND (cardinality(sqlc.arg(employee_ids)::int[]) = 0 OR employee_id = ANY(sqlc.arg(employee_ids)::int[]))
  AND (sqlc.narg(appointment_id)::int IS NULL OR appointment_id = sqlc.narg(appointment_id)::int)
  AND (cardinality(sqlc.arg(clinic_services)::varchar[]) = 0 OR clinic_service = ANY(sqlc.arg(clinic_services)::varchar[]))
  AND (sqlc.narg(is_emergency)::boolean IS NULL OR is_emergency = sqlc.narg(is_emergency)::boolean)
  AND (sqlc.narg(is_deleted)::boolean IS NULL OR (deleted_at IS NOT NULL) = sqlc.narg(is_deleted)::boolean)
  AND (sqlc.narg(visit_date_from)::timestamptz IS NULL OR visit_date >= sqlc.narg(visit_date_from)::timestamptz)
  AND (sqlc.narg(visit_date_to)::timestamptz IS NULL OR visit_date <= sqlc.narg(visit_date_to)::timestamptz)
  AND (sqlc.narg(follow_up_from)::timestamptz IS NULL OR follow_up_date >= sqlc.narg(follow_up_from)::timestamptz)
  AND (sqlc.narg(follow_up_to)::timestamptz IS NULL OR follow_up_date <= sqlc.narg(follow_up_to)::timestamptz);

-- name: RestoreMedicalSessionByID :exec
UPDATE medical_sessions
SET deleted_at = NULL, updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: IsDeletedMedicalSessionByID :one
SELECT deleted_at IS NOT NULL FROM medical_sessions WHERE id = $1;

-- name: CountAllMedicalSessionsIncludingDeleted :one
SELECT COUNT(*) FROM medical_sessions;

-- name: CountMedicalSessionsByClinicService :many
SELECT clinic_service, COUNT(*) AS count
FROM medical_sessions
WHERE deleted_at IS NULL
GROUP BY clinic_service;