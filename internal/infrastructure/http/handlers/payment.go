package handlers

import (
	"clinic-vet-api/internal/core/customers"
	"clinic-vet-api/internal/core/payments"
	"clinic-vet-api/internal/infrastructure/http/handlers/dtos"
	"clinic-vet-api/internal/infrastructure/http/handlers/mappers"
	"clinic-vet-api/internal/shared/errors"
	"clinic-vet-api/internal/shared/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type PaymentHandler struct {
	queryService   payments.QueryService
	commandService payments.CommandService
	validator      *validator.Validate
	mapper         *mappers.PaymentMapper
}

func NewPaymentHandler(
	queryService payments.QueryService,
	commandService payments.CommandService,
	validator *validator.Validate,
	mapper *mappers.PaymentMapper,
) *PaymentHandler {
	return &PaymentHandler{
		queryService:   queryService,
		commandService: commandService,
		validator:      validator,
		mapper:         mapper,
	}
}

// GetPaymentByID godoc
// @Summary      Get payment by ID
// @Description  Returns a single payment by ID.
// @Tags         payments
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Payment ID"
// @Success      200  {object}  http.APIResponse{data=dtos.PaymentResponse}  "Payment found"
// @Failure      400  {object}  http.APIResponse  "Invalid ID"
// @Failure      401  {object}  http.APIResponse  "Unauthorized"
// @Failure      404  {object}  http.APIResponse  "Payment not found"
// @Failure      500  {object}  http.APIResponse  "Internal server error"
// @Router       /payments/{id} [get]
func (s *PaymentHandler) GetPaymentByID(c *gin.Context) {
	paymentIDUInt, err := http.ParseParamToUInt(c, "id")
	if err != nil {
		http.BadRequest(c, errors.RequestURLParamError(err, "id", c.Param("id")))
		return
	}

	paymentID := payments.NewPaymentID(paymentIDUInt)
	payment, err := s.queryService.GetPaymentByID(c.Request.Context(), paymentID)
	if err != nil {
		http.ApplicationError(c, err)
		return
	}

	paymentResponse := s.mapper.ToPaymentResponse(payment)
	http.Found(c, paymentResponse, "Payment")
}

func (s *PaymentHandler) GetPaymentByTransactionID(c *gin.Context) {
	transactionID := c.Param("transaction_id")
	payment, err := s.queryService.GetPaymentByTransactionID(c.Request.Context(), transactionID)
	if err != nil {
		http.ApplicationError(c, err)
		return
	}

	paymentResponse := s.mapper.ToPaymentResponse(payment)
	http.Found(c, paymentResponse, "Payment")
}

func (s *PaymentHandler) GetPaymentsBySpecification(c *gin.Context) {
	var requestBodyData dtos.PaymentSearchRequest
	if err := http.ShouldBindAndValidateBody(c, &requestBodyData, s.validator); err != nil {
		http.BadRequest(c, err)
		return
	}

	query, err := s.mapper.RequestToSpecification(requestBodyData)
	if err != nil {
		http.BadRequest(c, err)
		return
	}

	paymentPage, err := s.queryService.GetPaymentsBySpecification(c.Request.Context(), query)
	if err != nil {
		http.ApplicationError(c, err)
		return
	}

	responsePage := s.mapper.ToPaginatedResponse(paymentPage)
	http.Paginated(c, &responsePage, "Payments")
}

func (s *PaymentHandler) GetPaymentsByCustomerID(c *gin.Context) {
	customerIDUInt, err := http.ParseParamToUInt(c, "id")
	if err != nil {
		http.BadRequest(c, errors.RequestURLParamError(err, "id", c.Param("id")))
		return
	}

	customerID := customers.NewCustomerID(customerIDUInt)
	paymentPage, err := s.queryService.GetPaymentsByCustomerID(c.Request.Context(), customerID)
	if err != nil {
		http.ApplicationError(c, err)
		return
	}

	responsePage := s.mapper.ToPaginatedResponse(paymentPage)
	http.Paginated(c, &responsePage, "Payments")
}

func (s *PaymentHandler) CreatePayment(c *gin.Context) {
	var requestBodyData dtos.PaymentCreateRequest
	if err := http.ShouldBindAndValidateBody(c, &requestBodyData, s.validator); err != nil {
		http.BadRequest(c, err)
		return
	}

	command, err := s.mapper.RequestToCreateCommand(requestBodyData)
	if err != nil {
		http.BadRequest(c, err)
		return
	}

	payment, err := s.commandService.CreatePayment(c.Request.Context(), command)
	if err != nil {
		http.ApplicationError(c, err)
		return
	}

	http.Created(c, payment.ID.Value(), "Payment")
}

func (s *PaymentHandler) UpdatePayment(c *gin.Context) {
	paymentIDUInt, err := http.ParseParamToUInt(c, "id")
	if err != nil {
		http.BadRequest(c, errors.RequestURLParamError(err, "id", c.Param("id")))
		return
	}

	var requestBodyData dtos.PaymentUpdateRequest
	if err := http.ShouldBindAndValidateBody(c, &requestBodyData, s.validator); err != nil {
		http.BadRequest(c, err)
		return
	}

	paymentID := payments.NewPaymentID(paymentIDUInt)
	current, err := s.queryService.GetPaymentByID(c.Request.Context(), paymentID)
	if err != nil {
		http.ApplicationError(c, err)
		return
	}

	command, err := s.mapper.RequestToUpdateCommand(paymentIDUInt, requestBodyData, current)
	if err != nil {
		http.BadRequest(c, err)
		return
	}

	if err := s.commandService.UpdatePayment(c.Request.Context(), command); err != nil {
		http.ApplicationError(c, err)
		return
	}

	http.Updated(c, nil, "Payment")
}

func (s *PaymentHandler) DeletePayment(c *gin.Context) {
	paymentIDUInt, err := http.ParseParamToUInt(c, "id")
	if err != nil {
		http.BadRequest(c, errors.RequestURLParamError(err, "id", c.Param("id")))
		return
	}

	paymentID := payments.NewPaymentID(paymentIDUInt)
	command := payments.DeletePaymentCommand{
		ID: paymentID,
	}

	err = s.commandService.DeletePayment(c.Request.Context(), command)
	if err != nil {
		http.ApplicationError(c, err)
		return
	}

	http.Success(c, nil, "Payment")
}

func (s *PaymentHandler) CancelPayment(c *gin.Context) {
	paymentIDUInt, err := http.ParseParamToUInt(c, "id")
	if err != nil {
		http.BadRequest(c, errors.RequestURLParamError(err, "id", c.Param("id")))
		return
	}

	paymentID := payments.NewPaymentID(paymentIDUInt)
	command := payments.CancelPaymentCommand{
		ID: paymentID,
	}

	err = s.commandService.CancelPayment(c.Request.Context(), command)
	if err != nil {
		http.ApplicationError(c, err)
		return
	}

	http.Success(c, nil, "Payment")
}

func (s *PaymentHandler) ProcessPayment(c *gin.Context) {
	paymentIDUInt, err := http.ParseParamToUInt(c, "id")
	if err != nil {
		http.BadRequest(c, errors.RequestURLParamError(err, "id", c.Param("id")))
		return
	}

	paymentID := payments.NewPaymentID(paymentIDUInt)
	command := payments.ProcessPaymentCommand{
		ID: paymentID,
	}

	err = s.commandService.ProcessPayment(c.Request.Context(), command)
	if err != nil {
		http.ApplicationError(c, err)
		return
	}

	http.Success(c, nil, "Payment")
}

func (s *PaymentHandler) RefundPayment(c *gin.Context) {
	paymentIDUInt, err := http.ParseParamToUInt(c, "id")
	if err != nil {
		http.BadRequest(c, errors.RequestURLParamError(err, "id", c.Param("id")))
		return
	}

	paymentID := payments.NewPaymentID(paymentIDUInt)
	command := payments.RefundPaymentCommand{
		ID: paymentID,
	}

	err = s.commandService.RefundPayment(c.Request.Context(), command)
	if err != nil {
		http.ApplicationError(c, err)
		return
	}

	http.Success(c, nil, "Payment")
}
