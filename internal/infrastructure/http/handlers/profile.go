package handlers

import (
	"clinic-vet-api/internal/core/users"
	"clinic-vet-api/internal/infrastructure/http/handlers/dtos"
	"clinic-vet-api/internal/shared"
	"clinic-vet-api/internal/shared/errors"
	"clinic-vet-api/internal/shared/http"

	"github.com/gin-gonic/gin"
)

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

// TODO: Add Handlers to roles
