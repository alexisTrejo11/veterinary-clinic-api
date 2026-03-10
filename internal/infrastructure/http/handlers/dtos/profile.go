package dtos

import (
	"time"
)

type UpdateProfileRequest struct {
	Name        string     `json:"name" binding:"required"`
	Gender      string     `json:"gender" binding:"required"`
	DateOfBirth *time.Time `json:"date_of_birth" binding:"omitempty"`
	PhotoURL    string     `json:"photo_url" binding:"omitempty,url"`
	Bio         string     `json:"bio" binding:"omitempty,max=500"`
}

type ProfileResponse struct {
	ID            string     `json:"id"`
	Email         string     `json:"email"`
	Phone         string     `json:"phone,omitempty"`
	Role          string     `json:"role"`
	Status        string     `json:"status"`
	Name          string     `json:"name"`
	Gender        string     `json:"gender"`
	DateOfBirth   *time.Time `json:"date_of_birth"`
	ProfilePicUrl string     `json:"profile_pic_url"`
	Bio           string     `json:"bio"`
}
