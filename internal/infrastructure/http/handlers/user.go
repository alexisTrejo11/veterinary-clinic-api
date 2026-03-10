package handlers

import (
	"clinic-vet-api/internal/core/users"
	"clinic-vet-api/internal/infrastructure/http/handlers/dtos"
	"clinic-vet-api/internal/infrastructure/http/handlers/mappers"
	"clinic-vet-api/internal/shared"
	"clinic-vet-api/internal/shared/errors"
	"clinic-vet-api/internal/shared/http"
	"clinic-vet-api/internal/shared/page"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type UserBaseController struct {
	validator      *validator.Validate
	queryService   users.QueryService
	commandService users.CommandService
	mapper         mappers.UserHandlerMapper
}

func NewUserBaseController(
	validator *validator.Validate,
	queryService users.QueryService,
	commandService users.CommandService,
	mapper mappers.UserHandlerMapper,
) *UserBaseController {
	return &UserBaseController{
		validator:      validator,
		queryService:   queryService,
		commandService: commandService,
		mapper:         mapper,
	}
}

func (ctrl *UserBaseController) GetUserByID(c *gin.Context) {
	userIDUInt, err := http.ParseParamToUInt(c, "id")
	if err != nil {
		http.BadRequest(c, errors.RequestURLParamError(err, "id", c.Param("id")))
		return
	}

	userID := shared.NewUserID(userIDUInt)
	user, err := ctrl.queryService.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		http.ApplicationError(c, err)
		return
	}

	userResponse := ctrl.mapper.ToResponse(user)
	http.Found(c, userResponse, "User")
}

func (ctrl *UserBaseController) GetUserByEmail(c *gin.Context) {
	emailStr := c.Param("email")
	if emailStr == "" {
		http.BadRequest(c, errors.RequestURLParamError(fmt.Errorf("email is required"), "email", ""))
		return
	}

	email, err := users.NewEmail(emailStr)
	if err != nil {
		http.ApplicationError(c, err)
		return
	} 

	user, err := ctrl.queryService.GetByEmail(c.Request.Context(), query)
	if err != nil {
		http.ApplicationError(c, err)
		return
	}

	userResponse := ctrl.mapper.ToResponse(user)
	http.Found(c, userResponse, "User")
}

func (ctrl *UserBaseController) GetUserByPhone(ctx *gin.Context) {
	phone := ctx.Query("phone")
	if phone == "" {
		err := fmt.Errorf("phone is required")
		http.BadRequest(ctx, errors.RequestURLParamError(err, "phone", ""))
		return
	}

	user, err := c.commandService.GetUserByPhone(ctx.Request.Context(), query)
	if err != nil {
		http.ApplicationError(ctx, err)
		return
	}

	userResponse := ctrl.mapper.ToResponse(user)
	http.Found(ctx, userResponse, "User")
}

func (ctrl *UserBaseController) SearchUsers(ctx *gin.Context) {
	role := ctx.Param("role")
	if role == "" {
		http.BadRequest(ctx, errors.RequestURLParamError(errors.New("role is required"), "role", ""))
		return
	}

	var pagination page.PaginationRequest
	if err := http.ShouldBindPageParams(&pagination, ctx, ctrl.validator); err != nil {
		http.BadRequest(ctx, errors.RequestURLQueryError(err, ctx.Request.URL.RawQuery))
		return
	}

	query := query.NewGetUsersByRoleQuery(role, pagination)
	users, err := ctrl.commandService.GetUsersByRole(ctx.Request.Context(), query)
	if err != nil {
		http.ApplicationError(ctx, err)
		return
	}

	userResponses := dtos.usersToResponses(users.Items)
	http.SuccessWithPagination(ctx, userResponses, "Users retrieved successfully", users.Metadata)
}

// TODO: Move to Auth
/*
func (ctrl *UserBaseController) SendResetPasswordEmail(c *gin.Context) {
	var req dtos.RequestEmailRequest
	if err := http.ShouldBindAndValidateQuery(c, &req, ctrl.validator); err != nil {
		http.BadRequest(c, err)
		return
	}

	command, err := req.ToCommand()
	if err != nil {
		http.ApplicationError(c, err)
		return
	}

	err = ctrl.commandService.RequestResetPassword(c.Request.Context(), command)
	if err != nil {
		http.ApplicationError(c, err)
		return
	}

	http.Success(c, nil, "")
}

func (ctrl *UserBaseController) ResetPassword(c *gin.Context) {
	var req dtos.ResetPasswordRequest
	if err := http.ShouldBindAndValidateQuery(c, &req, ctrl.validator); err != nil {
		http.BadRequest(c, err)
		return
	}

	command, err := req.ToCommand()
	if err != nil {
		http.ApplicationError(c, err)
		return
	}

	err := ctrl.commandService.ResetPassword(c.Request.Context(), command)
	if err != nil {
		http.ApplicationError(c, err)
		return
	}

	http.Success(c, nil, "")
}

func (ctrl *UserBaseController) UpdatePassword(c *gin.Context) {
	id, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		http.BadRequest(c, autherror.UnauthorizedCTXError())
		return
	}

	var req dtos.UpdatePasswordRequest
	if err := http.ShouldBindAndValidateQuery(c, &req, ctrl.validator); err != nil {
		http.BadRequest(c, err)
		return
	}

	cmd, err := req.ToCommand(id.Value())
	if err != nil {
		http.ApplicationError(c, err)
		return
	}

	err := ctrl.commandService.UpdatePassword(c.Request.Context(), cmd)
	if err != nil {
		http.ApplicationError(c, err)
		return
	}

	http.Success(c, nil, "")
}

func (ctrl *UserBaseController) UpdateEmail(c *gin.Context) {

}

func (ctrl *UserBaseController) UpdatePhoneNumber(c *gin.Context) {

}

func (ctrl *UserBaseController) ActivateAccount(c *gin.Context) {
	userID, err := http.ParseParamToUInt(c, "id")
	if err != nil && userID == 0 {
		http.BadRequest(c, http.MissingRequiredFieldError("user id is required"))
		return
	}

	code := c.Query("code")
	if code == "" {
		http.BadRequest(c, http.MissingRequiredFieldError("activation code is required"))
		return
	}

	command, err := users.NewActivateAccountCommand(code, userID)
	if err != nil {
		http.ApplicationError(c, err)
		return
	}

	err := ctrl.commandService.ActivateAccount(c.Request.Context(), command)
	if err != nil {
		http.ApplicationError(c, err)
		return
	}

	http.Success(c, nil, "")
}
*/

func (ctrl *UserBaseController) CreateUser(c *gin.Context) {
	var requestData dtos.CreateUserRequest
	if err := http.ShouldBindAndValidateBody(c, &requestData, ctrl.validator); err != nil {
		http.BadRequest(c, err)
		return
	}

	createUserCommand, err := ctrl.mapper.ToCreateUserCommand(requestData)
	if err != nil {
		http.ApplicationError(c, err)
		return
	}

	user, err := ctrl.commandService.CreateUser(c.Request.Context(), createUserCommand)
	if err != nil {
		http.ApplicationError(c, err)
		return
	}

	http.Created(c, user.ID, "User created successfully")
}

func (ctrl *UserBaseController) UpdateUserStatus(c *gin.Context) {
	userID, err := http.ParseParamToUInt(c, "id")
	if err != nil {
		http.BadRequest(c, errors.RequestURLParamError(err, "id", c.Param("id")))
		return
	}

	statusStr := c.Param("status")
	if statusStr == "" {
		err := fmt.Errorf("status is required")
		http.BadRequest(c, errors.RequestURLParamError(err, "status", ""))
	}

	status, err := users.ParseUserStatus(statusStr)
	if err != nil {
		http.ApplicationError(c, err)
		return
	}

	command := users.UpdateUserStatusCommand{
		ID:     shared.NewUserID(userID),
		Status: status,
	}

	err = ctrl.commandService.UpdateUserStatus(c.Request.Context(), command)
	if err != nil {
		http.ApplicationError(c, err)
		return
	}

	http.Success(c, nil, "User Succesfully Unbaned")
}

func (ctrl *UserBaseController) DeleteUser(c *gin.Context) {
	userIDUint, err := http.ParseParamToUInt(c, "id")
	if err != nil {
		http.BadRequest(c, errors.RequestURLParamError(err, "id", c.Param("id")))
		return
	}

	userID := shared.NewUserID(userIDUint)
	deleteCommand := users.DeleteUserCommand{
		ID:           userID,
		IsHardDelete: false,
	}

	err = ctrl.commandService.DeleteUser(c.Request.Context(), deleteCommand)
	if err != nil {
		http.ApplicationError(c, err)
		return
	}

	http.Success(c, nil, "User Succesfully Deleted")
}

func (ctrl *UserBaseController) RestoreUser(c *gin.Context) {
	userIDUint, err := http.ParseParamToUInt(c, "id")
	if err != nil {
		http.BadRequest(c, errors.RequestURLParamError(err, "id", c.Param("id")))
		return
	}

	userID := shared.NewUserID(userIDUint)
	err = ctrl.commandService.RestoreUser(c.Request.Context(), userID)
	if err != nil {
		http.ApplicationError(c, err)
		return
	}

	http.Success(c, nil, "User Succesfully Restore")
}
