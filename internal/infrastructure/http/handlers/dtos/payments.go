package dtos

import (
	"time"

	"clinic-vet-api/internal/shared/page"
)

// PaymentCreateRequest represents the body for creating a new payment record
// @Description Amount, currency (3-letter), method, optional description/due_date/transaction_id/invoice_id, customer_id, optional appointment_id. Internal tracker only (no external gateway).
type PaymentCreateRequest struct {
	Amount      float64   `json:"amount" binding:"required,gt=0"`
	Currency    string    `json:"currency" binding:"required,len=3"`
	Method      string    `json:"method" binding:"required"`
	Description *string   `json:"description,omitempty" binding:"omitempty,max=500"`
	DueDate     *time.Time `json:"due_date,omitempty" binding:"omitempty"`

	// Optional external tracking information
	TransactionID *string `json:"transaction_id,omitempty" binding:"omitempty"`
	InvoiceID     *string `json:"invoice_id,omitempty" binding:"omitempty"`

	// Relations
	CustomerID    uint  `json:"customer_id" binding:"required"`
	AppointmentID *uint `json:"appointment_id,omitempty" binding:"omitempty"`
}

// PaymentUpdateRequest represents the body for updating an existing payment
// @Description All fields optional; only provided fields are updated (amount, currency, method, description, due_date, invoice_id).
type PaymentUpdateRequest struct {
	Amount      *float64   `json:"amount,omitempty" binding:"omitempty,gt=0"`
	Currency    *string    `json:"currency,omitempty" binding:"omitempty,len=3"`
	Method      *string    `json:"method,omitempty" binding:"omitempty"`
	Description *string    `json:"description,omitempty" binding:"omitempty,max=500"`
	DueDate     *time.Time `json:"due_date,omitempty" binding:"omitempty"`

	InvoiceID *string `json:"invoice_id,omitempty" binding:"omitempty"`
}

// PaymentSearchRequest represents filters for listing/searching payments
// @Description Pagination plus optional customer_id, status, method, from/to created, overdue_only.
type PaymentSearchRequest struct {
	page.PaginationRequest

	CustomerID   *uint      `json:"customer_id,omitempty"`
	Status       *string    `json:"status,omitempty"`
	Method       *string    `json:"method,omitempty"`
	FromCreated  *time.Time `json:"from_created_at,omitempty"`
	ToCreated    *time.Time `json:"to_created_at,omitempty"`
	OverdueOnly  bool       `json:"overdue_only,omitempty"`
}

// PaymentResponse represents the HTTP response DTO for a payment
// @Description Payment entity: id, amount, currency, status, method, transaction_id, description, due_date, paid_at, customer_id, timestamps, optional refund/delete info.
type PaymentResponse struct {
	ID uint `json:"id"`

	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`

	Status        string `json:"status"`
	StatusDisplay string `json:"status_display"`

	Method        string `json:"method"`
	MethodDisplay string `json:"method_display"`

	TransactionID *string    `json:"transaction_id,omitempty"`
	Description   *string    `json:"description,omitempty"`
	DueDate       *time.Time `json:"due_date,omitempty"`
	PaidAt        *time.Time `json:"paid_at,omitempty"`
	RefundedAt    *time.Time `json:"refunded_at,omitempty"`
	IsActive      bool       `json:"is_active"`

	CustomerID    uint   `json:"customer_id"`
	InvoiceID     *string `json:"invoice_id,omitempty"`
	FailureReason *string `json:"failure_reason,omitempty"`

	RefundAmount *float64 `json:"refund_amount,omitempty"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

