package payments

import (
	"clinic-vet-api/internal/shared"
	"context"
	"fmt"
	"time"
)

type CommandService interface {
	CancelPayment(ctx context.Context, cmd CancelPaymentCommand) error
	ProcessPayment(ctx context.Context, cmd ProcessPaymentCommand) error
	RefundPayment(ctx context.Context, cmd RefundPaymentCommand) error
	DeletePayment(ctx context.Context, cmd DeletePaymentCommand) error
	CreatePayment(ctx context.Context, cmd CreatePaymentCommand) (Payment, error)
	UpdatePayment(ctx context.Context, cmd UpdatePaymentCommand) error
}

type commandService struct {
	repository PaymentRepository
}

func NewCommandService(repository PaymentRepository) CommandService {
	return &commandService{repository: repository}
}

func (s *commandService) CancelPayment(ctx context.Context, cmd CancelPaymentCommand) error {
	payment, err := s.repository.GetByID(ctx, cmd.ID)
	if err != nil {
		return err
	}

	// For now we use a generic reason as this is an internal tracker.
	reason := "Cancelled by clinic staff"
	if err := payment.Cancel(ctx, reason); err != nil {
		return err
	}

	return s.repository.Save(ctx, &payment)
}

func (s *commandService) ProcessPayment(ctx context.Context, cmd ProcessPaymentCommand) error {
	payment, err := s.repository.GetByID(ctx, cmd.ID)
	if err != nil {
		return err
	}

	// Generate an internal transaction ID since we do not integrate
	// with external payment gateways.
	transactionID := fmt.Sprintf("INT-%d-%d", payment.ID.Value, time.Now().UnixNano())

	if err := payment.Pay(ctx, transactionID); err != nil {
		return err
	}

	return s.repository.Save(ctx, &payment)
}

func (s *commandService) RefundPayment(ctx context.Context, cmd RefundPaymentCommand) error {
	payment, err := s.repository.GetByID(ctx, cmd.ID)
	if err != nil {
		return err
	}

	if err := payment.Refund(ctx); err != nil {
		return err
	}

	return s.repository.Save(ctx, &payment)
}

func (s *commandService) DeletePayment(ctx context.Context, cmd DeletePaymentCommand) error {
	// Use soft delete for safety; this keeps a history of payments.
	if err := s.repository.SoftDelete(ctx, cmd.ID); err != nil {
		return err
	}
	return nil
}

func (s *commandService) CreatePayment(ctx context.Context, cmd CreatePaymentCommand) (Payment, error) {
	operation := "CreatePayment"

	if cmd.Amount <= 0 {
		money := shared.NewMoney(shared.NewDecimalFromFloat(cmd.Amount), cmd.Currency)
		return Payment{}, AmountInvalidError(ctx, money, operation)
	}

	amountDecimal := shared.NewDecimalFromFloat(cmd.Amount)
	amount := shared.NewMoney(amountDecimal, cmd.Currency)

	if !cmd.Method.IsValid() {
		return Payment{}, MethodInvalidError(ctx, cmd.Method, operation)
	}

	if cmd.CustomerID.IsZero() {
		return Payment{}, CustomerIDRequiredError(ctx, operation)
	}

	if !cmd.DueDate.IsZero() && cmd.DueDate.Before(time.Now()) {
		return Payment{}, DueDatePastError(ctx, operation)
	}

	entity := Payment{
		Amount:           amount,
		Status:           cmd.Status,
		Method:           cmd.Method,
		TransactionID:    nil,
		Description:      nil,
		DueDate:          nil,
		PaidAt:           nil,
		RefundedAt:       nil,
		IsActive:         true,
		PaidFromCustomer: cmd.CustomerID,
		InvoiceID:        nil,
	}

	if cmd.TransactionID != "" {
		entity.TransactionID = &cmd.TransactionID
	}
	if cmd.Description != "" {
		entity.Description = &cmd.Description
	}
	if !cmd.DueDate.IsZero() {
		due := cmd.DueDate
		entity.DueDate = &due
	}
	if cmd.InvoiceID != "" {
		inv := cmd.InvoiceID
		entity.InvoiceID = &inv
	}
	if !cmd.AppointmentID.IsZero() {
		appt := cmd.AppointmentID
		entity.AppointmentID = &appt
	}

	if entity.Status == "" {
		entity.Status = PaymentStatusPending
	}

	if err := s.repository.Save(ctx, &entity); err != nil {
		return Payment{}, err
	}

	return entity, nil
}

func (s *commandService) UpdatePayment(ctx context.Context, cmd UpdatePaymentCommand) error {
	operation := "UpdatePayment"

	payment, err := s.repository.GetByID(ctx, cmd.ID)
	if err != nil {
		return err
	}

	// Do not allow updates once the payment has been processed.
	if !payment.IsPending() {
		return AlreadyProcessedError(ctx, payment.Status, operation)
	}

	var amount *shared.Money
	if cmd.Amount > 0 {
		dec := shared.NewDecimalFromFloat(cmd.Amount)
		m := shared.NewMoney(dec, cmd.Currency)
		amount = &m
	}

	var method *PaymentMethod
	if cmd.Method != "" {
		m := cmd.Method
		method = &m
	}

	var description *string
	if cmd.Description != "" {
		desc := cmd.Description
		description = &desc
	}

	var dueDate *time.Time
	if !cmd.DueDate.IsZero() {
		d := cmd.DueDate
		dueDate = &d
	}

	if err := payment.Update(ctx, amount, method, description, dueDate); err != nil {
		return err
	}

	// Update relations if provided.
	if !cmd.CustomerID.IsZero() {
		payment.PaidFromCustomer = cmd.CustomerID
	}
	if !cmd.AppointmentID.IsZero() {
		appt := cmd.AppointmentID
		payment.AppointmentID = &appt
	}
	if cmd.InvoiceID != "" {
		inv := cmd.InvoiceID
		payment.InvoiceID = &inv
	}

	return s.repository.Save(ctx, &payment)
}
