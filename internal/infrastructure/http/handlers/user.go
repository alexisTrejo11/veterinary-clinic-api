package handlers

import (
	"clinic-vet-api/internal/core/users"
	"clinic-vet-api/internal/infrastructure/http/handlers/dtos"
	"clinic-vet-api/internal/infrastructure/http/handlers/mappers"
	"clinic-vet-api/internal/shared"
	"clinic-vet-api/internal/shared/errors"
	"clinic-vet-api/internal/shared/http"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type UserHandler struct {
	validator      *validator.Validate
	queryService   users.QueryService
	commandService users.CommandService
	mapper         mappers.UserHandlerMapper
}

func NewUserHandler(
	validator *validator.Validate,
	queryService users.QueryService,
	commandService users.CommandService,
	mapper mappers.UserHandlerMapper,
) *UserHandler {
	return &UserHandler{
		validator:      validator,
		queryService:   queryService,
		commandService: commandService,
		mapper:         mapper,
	}
}

func (ctrl *UserHandler) GetUserByID(c *gin.Context) {
	userIDUInt, err := http.ParseParamToUInt(c, "id")
	if err != nil {
		http.BadRequest(c, errors.RequestURLParamError(err, "id", c.Param("id")))
		return
	}

	userID := shared.NewUserID(userIDUInt)
	user, err := ctrl.queryService.GetByID(c.Request.Context(), userID)
	if err != nil {
		http.ApplicationError(c, err)
		return
	}

	userResponse := ctrl.mapper.ToResponse(user)
	http.Found(c, userResponse, "User")
}

func (ctrl *UserHandler) GetUserByEmail(c *gin.Context) {
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

	user, err := ctrl.queryService.GetByEmail(c.Request.Context(), email)
	if err != nil {
		http.ApplicationError(c, err)
		return
	}

	userResponse := ctrl.mapper.ToResponse(user)
	http.Found(c, userResponse, "User")
}

func (ctrl *UserHandler) GetUserByPhone(ctx *gin.Context) {
	phoneStr := ctx.Query("phone")
	if phoneStr == "" {
		http.BadRequest(ctx, errors.RequestURLParamError(fmt.Errorf("phone is required"), "phone", ""))
		return
	}

	phone, err := users.NewPhoneNumber(phoneStr)
	if err != nil {
		http.ApplicationError(ctx, err)
		return
	}

	user, err := ctrl.queryService.GetByPhone(ctx.Request.Context(), phone)
	if err != nil {
		http.ApplicationError(ctx, err)
		return
	}

	userResponse := ctrl.mapper.ToResponse(user)
	http.Found(ctx, userResponse, "User")
}

func (ctrl *UserHandler) SearchUsers(ctx *gin.Context) {
	var req dtos.UserSearchRequest
	if err := http.ShouldBindAndValidateQuery(ctx, &req, ctrl.validator); err != nil {
		http.BadRequest(ctx, err)
		return
	}

	spec, err := ctrl.mapper.UserSearchRequestToSpecification(req)
	if err != nil {
		http.ApplicationError(ctx, err)
		return
	}

	users, err := ctrl.queryService.GetSpecification(ctx.Request.Context(), *spec)
	if err != nil {
		http.ApplicationError(ctx, err)
		return
	}

	userResponses := ctrl.mapper.ToResponsePage(users)
	http.Paginated(ctx, &userResponses, "Users")
}

func (ctrl *UserHandler) CreateUser(c *gin.Context) {
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

func (ctrl *UserHandler) UpdateUserStatus(c *gin.Context) {
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

func (ctrl *UserHandler) DeleteUser(c *gin.Context) {
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

func (ctrl *UserHandler) RestoreUser(c *gin.Context) {
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

// ------------------------------------------------------------
// Profile Routes
// ------------------------------------------------------------

type ProfileHandler struct {
	queryService   users.QueryService
	commandService users.CommandService
}

func NewProfileHandler(queryService users.QueryService, commandService users.CommandService) *ProfileHandler {
	return &ProfileHandler{queryService: queryService, commandService: commandService}
}

func (h *ProfileHandler) GetProfile(c *gin.Context) {
	userIDStr := c.GetString("userID")

	userID, err := shared.ParseUserIDFromString(userIDStr)
	if err != nil {
		http.ApplicationError(c, err)
		return
	}

	user, err := h.queryService.GetByID(c.Request.Context(), userID)
	if err != nil {
		http.ApplicationError(c, err)
		return
	}

	profile := dtos.ProfileResponse{
		ID:            user.ID.String(),
		Email:         user.Email.Value(),
		Phone:         user.PhoneNumber.Value(),
		Role:          user.Role.DisplayName(),
		Status:        user.Status.DisplayName(),
		Name:          user.Profile.Name,
		Gender:        user.Profile.Gender.DisplayName(),
		DateOfBirth:   user.Profile.DateOfBirth,
		ProfilePicUrl: user.Profile.PhotoURL,
		Bio:           user.Profile.Bio,
	}

	http.Success(c, profile, "Profile retrieved successfully")
}

func (h *ProfileHandler) UpdateProfile(c *gin.Context) {
	userIDStr := c.GetString("userID")

	var requestData dtos.UpdateProfileRequest
	if err := c.ShouldBindJSON(&requestData); err != nil {
		errors.RequestBodyDataError(err)
		return
	}

	genderVO, err := shared.ParseGender(requestData.Gender)
	if err != nil {
		http.ApplicationError(c, err)
		return
	}
	userIDVO, err := shared.ParseUserIDFromString(userIDStr)
	if err != nil {
		http.ApplicationError(c, err)
		return
	}

	cmd := users.UpdateProfileCommand{
		ID:          userIDVO,
		Name:        requestData.Name,
		Gender:      genderVO,
		DateOfBirth: requestData.DateOfBirth,
		PhotoURL:    requestData.PhotoURL,
		Bio:         requestData.Bio,
	}

	if err := h.commandService.UpdateProfileData(c.Request.Context(), cmd); err != nil {
		http.ApplicationError(c, err)
		return
	}

	http.Success(c, nil, "Profile updated successfully")
}
