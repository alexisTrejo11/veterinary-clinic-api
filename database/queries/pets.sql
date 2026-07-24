-- name: FindPetByID :one
SELECT * FROM pets
WHERE id = $1 AND deleted_at IS NULL;

-- name: FindPetByIDAndCustomerID :one
SELECT * FROM pets
WHERE id = $1 AND customer_id = $2 AND deleted_at IS NULL;

-- name: FindAllPetsByCustomerID :many
SELECT * FROM pets
WHERE customer_id = $1 AND deleted_at IS NULL
ORDER BY created_at DESC;

-- name: FindPetsByCustomerID :many
SELECT * FROM pets
WHERE customer_id = $1 AND deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: FindPetsBySpecies :many
SELECT * FROM pets
WHERE species = $1 AND deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: CountPetsBySpecies :one
SELECT COUNT(*) FROM pets
WHERE species = $1 AND deleted_at IS NULL;

-- name: CountPetsByCustomerID :one
SELECT COUNT(*) FROM pets
WHERE customer_id = $1 AND deleted_at IS NULL;

-- name: ExistsPetByID :one
SELECT COUNT(*) > 0 FROM pets
WHERE id = $1 AND deleted_at IS NULL;

-- name: ExistsDeletedPetByID :one
SELECT COUNT(*) > 0 FROM pets
WHERE id = $1 AND deleted_at IS NOT NULL;

-- name: ExistsPetByMicrochip :one
SELECT COUNT(*) > 0 FROM pets
WHERE microchip = $1 AND deleted_at IS NULL;

-- name: CountPetsBySpecification :one
SELECT COUNT(*) FROM pets
WHERE deleted_at IS NULL
  AND (cardinality(sqlc.arg(ids)::int[]) = 0 OR id = ANY(sqlc.arg(ids)::int[]))
  AND (cardinality(sqlc.arg(customer_ids)::int[]) = 0 OR customer_id = ANY(sqlc.arg(customer_ids)::int[]))
  AND (cardinality(sqlc.arg(species)::text[]) = 0 OR species = ANY(sqlc.arg(species)::text[]))
  AND (cardinality(sqlc.arg(genders)::text[]) = 0 OR gender = ANY(sqlc.arg(genders)::text[]))
  AND (sqlc.narg(is_active)::boolean IS NULL OR is_active = sqlc.narg(is_active)::boolean)
  AND (sqlc.arg(search)::text = '' OR name ILIKE '%' || sqlc.arg(search)::text || '%' OR breed ILIKE '%' || sqlc.arg(search)::text || '%');

-- name: FindPetsBySpecification :many
SELECT * FROM pets
WHERE deleted_at IS NULL
  AND (cardinality(sqlc.arg(ids)::int[]) = 0 OR id = ANY(sqlc.arg(ids)::int[]))
  AND (cardinality(sqlc.arg(customer_ids)::int[]) = 0 OR customer_id = ANY(sqlc.arg(customer_ids)::int[]))
  AND (cardinality(sqlc.arg(species)::text[]) = 0 OR species = ANY(sqlc.arg(species)::text[]))
  AND (cardinality(sqlc.arg(genders)::text[]) = 0 OR gender = ANY(sqlc.arg(genders)::text[]))
  AND (sqlc.narg(is_active)::boolean IS NULL OR is_active = sqlc.narg(is_active)::boolean)
  AND (sqlc.arg(search)::text = '' OR name ILIKE '%' || sqlc.arg(search)::text || '%' OR breed ILIKE '%' || sqlc.arg(search)::text || '%')
ORDER BY created_at DESC
LIMIT sqlc.arg(page_limit) OFFSET sqlc.arg(page_offset);

-- name: CreatePet :one
INSERT INTO pets (
    name, 
    photo, 
    species, 
    breed, 
    age, 
    gender,
    color, 
    microchip, 
    blood_type,
    is_neutered, 
    customer_id,  
    is_active,
    allergies,
    current_medications,
    special_needs,
    feeding_instructions,
    behavioral_notes,
    veterinary_contact,
    emergency_contact_name,
    emergency_contact_phone,
    created_at, 
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
)
RETURNING *;

-- name: UpdatePet :one
UPDATE pets
SET 
    name = $2,
    photo = $3,
    species = $4,
    breed = $5,
    age = $6,
    gender = $7,
    color = $8,
    microchip = $9,
    is_neutered = $10,
    customer_id = $11,
    blood_type = $12,
    is_active = $13,
    allergies = $14,
    current_medications = $15,
    special_needs = $16,
    feeding_instructions = $17,
    behavioral_notes = $18,
    veterinary_contact = $19,
    emergency_contact_name = $20,
    emergency_contact_phone = $21,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: SoftDeletePet :exec
UPDATE pets
SET 
    deleted_at = CURRENT_TIMESTAMP,
    updated_at = CURRENT_TIMESTAMP,
    is_active = FALSE
WHERE id = $1;

-- name: HardDeletePet :exec
DELETE FROM pets WHERE id = $1;


-- name: RestorePet :exec
UPDATE pets
SET 
    deleted_at = NULL,
    updated_at = CURRENT_TIMESTAMP,
    is_active = TRUE
WHERE id = $1;