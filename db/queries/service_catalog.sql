-- name: FindServiceCatalogByID :one
SELECT * FROM service_catalog WHERE id = $1;

-- name: FindServiceCatalogAll :many
SELECT * FROM service_catalog WHERE is_active = TRUE ORDER BY name LIMIT $1 OFFSET $2;

-- name: FindServiceCatalogByCategory :many
SELECT * FROM service_catalog WHERE is_active = TRUE AND category = $1 ORDER BY name LIMIT $2 OFFSET $3;

-- name: CountServiceCatalogAll :one
SELECT COUNT(*) FROM service_catalog WHERE is_active = TRUE;

-- name: CountServiceCatalogByCategory :one
SELECT COUNT(*) FROM service_catalog WHERE is_active = TRUE AND category = $1;

-- name: CreateServiceCatalog :one
INSERT INTO service_catalog (name, category, description, base_price, duration_min, requires_fasting, is_active)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: UpdateServiceCatalog :one
UPDATE service_catalog
SET name = $2, category = $3, description = $4, base_price = $5, duration_min = $6, requires_fasting = $7, is_active = $8
WHERE id = $1
RETURNING *;

-- name: DeleteServiceCatalogByID :exec
DELETE FROM service_catalog WHERE id = $1;
