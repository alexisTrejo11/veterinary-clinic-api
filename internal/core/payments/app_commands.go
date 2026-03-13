package payments

import (
	"clinic-vet-api/internal/core/appointments"
	"clinic-vet-api/internal/core/customers"
	"time"
)

type CancelPaymentCommand struct {
	ID PaymentID
}

type ProcessPaymentCommand struct {
	ID PaymentID
}

type RefundPaymentCommand struct {
	ID PaymentID
}

type DeletePaymentCommand struct {
	ID PaymentID
}

type CreatePaymentCommand struct {
	Amount        float64
	Currency      string
	Status        PaymentStatus
	Method        PaymentMethod
	TransactionID string
	Description   string
	DueDate       time.Time
	CustomerID    customers.CustomerID
	AppointmentID appointments.AppointmentID
	InvoiceID     string
}

type UpdatePaymentCommand struct {
	ID            PaymentID
	Amount        float64
	Currency      string
	Status        PaymentStatus
	Method        PaymentMethod
	TransactionID string
	Description   string
	DueDate       time.Time
	CustomerID    customers.CustomerID
	AppointmentID appointments.AppointmentID
	InvoiceID     string
}
