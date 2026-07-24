-- name: FindPaymentByID :one
SELECT id, amount, currency, status, method, med_session_id, transaction_id,
       description, due_date, paid_at, refunded_at, is_active, created_at,
       updated_at, deleted_at, paid_by_customer_id, invoice_id, refund_amount,
       failure_reason
FROM payments
WHERE id = $1 AND deleted_at IS NULL;

-- name: FindPaymentByTransactionID :one
SELECT id, amount, currency, status, method, med_session_id, transaction_id,
       description, due_date, paid_at, refunded_at, is_active, created_at,
       updated_at, deleted_at, paid_by_customer_id, invoice_id, refund_amount,
       failure_reason
FROM payments
WHERE transaction_id = $1 AND deleted_at IS NULL;

-- name: FindPaymentByIDAndCustomerID :one
SELECT id, amount, currency, status, method, med_session_id, transaction_id,
       description, due_date, paid_at, refunded_at, is_active, created_at,
       updated_at, deleted_at, paid_by_customer_id, invoice_id, refund_amount,
       failure_reason
FROM payments
WHERE id = $1 AND paid_by_customer_id = $2 AND deleted_at IS NULL;

-- name: CountPaymentsByCustomerID :one
SELECT COUNT(*) FROM payments
WHERE paid_by_customer_id = $1 AND deleted_at IS NULL;

-- name: CountPaymentsByStatus :one
SELECT COUNT(*) FROM payments
WHERE status = $1 AND deleted_at IS NULL;

-- name: CountOverduePayments :one
SELECT COUNT(*) FROM payments
WHERE due_date < CURRENT_TIMESTAMP 
  AND status NOT IN ('paid', 'refunded', 'cancelled') 
  AND deleted_at IS NULL;

-- name: ExistsPaymentByID :one
SELECT COUNT(*) > 0 FROM payments
WHERE id = $1 AND deleted_at IS NULL;

-- name: ExistsPaymentByTransactionID :one
SELECT COUNT(*) > 0 FROM payments
WHERE transaction_id = $1 AND deleted_at IS NULL;

-- name: ExistsPendingPaymentByCustomerID :one
SELECT COUNT(*) > 0 FROM payments
WHERE paid_by_customer_id = $1 
  AND status = 'pending' 
  AND deleted_at IS NULL;

-- name: CreatePayment :one
INSERT INTO payments (
    amount, currency, status, method, transaction_id, description,
    due_date, paid_at, refunded_at, paid_by_customer_id,
    med_session_id, invoice_id, refund_amount, failure_reason
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
    $11, $12, $13, $14
)
RETURNING id, amount, currency, status, method, med_session_id, transaction_id,
          description, due_date, paid_at, refunded_at, is_active, created_at,
          updated_at, deleted_at, paid_by_customer_id, invoice_id, refund_amount,
          failure_reason;

-- name: UpdatePayment :one
UPDATE payments SET
    amount = $2,
    currency = $3,
    status = $4,
    method = $5,
    transaction_id = $6,
    description = $7,
    due_date = $8,
    paid_at = $9,
    refunded_at = $10,
    paid_by_customer_id = $11,
    med_session_id = $12,
    invoice_id = $13,
    refund_amount = $14,
    failure_reason = $15,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1 AND deleted_at IS NULL
RETURNING id, amount, currency, status, method, med_session_id, transaction_id,
          description, due_date, paid_at, refunded_at, is_active, created_at,
          updated_at, deleted_at, paid_by_customer_id, invoice_id, refund_amount,
          failure_reason;

-- name: SoftDeletePayment :exec
UPDATE payments SET
    is_active = FALSE,
    deleted_at = CURRENT_TIMESTAMP,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: HardDeletePayment :exec
DELETE FROM payments WHERE id = $1;

-- name: TotalRevenueByDateRange :one
SELECT COALESCE(SUM(amount), 0) FROM payments
WHERE status = 'paid'
  AND paid_at BETWEEN $1 AND $2
  AND deleted_at IS NULL;

-- name: FindPaymentsBySpecification :many
SELECT id, amount, currency, status, method, med_session_id, transaction_id,
       description, due_date, paid_at, refunded_at, is_active, created_at,
       updated_at, deleted_at, paid_by_customer_id, invoice_id, refund_amount,
       failure_reason
FROM payments
WHERE deleted_at IS NULL
  AND ($1::int IS NULL OR paid_by_customer_id = $1)
  AND ($2::payment_status IS NULL OR status = $2)
  AND ($3::payment_method IS NULL OR method = $3)
  AND ($4::timestamptz IS NULL OR created_at >= $4)
  AND ($5::timestamptz IS NULL OR created_at <= $5)
  AND (NOT $6::boolean OR (due_date < CURRENT_TIMESTAMP AND status NOT IN ('paid','refunded','cancelled')))
ORDER BY created_at DESC
LIMIT $7 OFFSET $8;

-- name: CountPaymentsBySpecification :one
SELECT COUNT(*) FROM payments
WHERE deleted_at IS NULL
  AND ($1::int IS NULL OR paid_by_customer_id = $1)
  AND ($2::payment_status IS NULL OR status = $2)
  AND ($3::payment_method IS NULL OR method = $3)
  AND ($4::timestamptz IS NULL OR created_at >= $4)
  AND ($5::timestamptz IS NULL OR created_at <= $5)
  AND (NOT $6::boolean OR (due_date < CURRENT_TIMESTAMP AND status NOT IN ('paid','refunded','cancelled')));

