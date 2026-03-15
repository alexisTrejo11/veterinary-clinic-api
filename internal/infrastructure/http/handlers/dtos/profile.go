package dtos

import (
	"time"
)

// UpdateProfileRequest is the body for PUT /profile
// @Description Authenticated user's profile fields to update (name, gender, date of birth, photo URL, bio).
type UpdateProfileRequest struct {
	Name        string     `json:"name" binding:"required" example:"John Doe"`
	Gender      string     `json:"gender" binding:"required" example:"male"`
	DateOfBirth *time.Time `json:"date_of_birth" binding:"omitempty"`
	PhotoURL    string     `json:"photo_url" binding:"omitempty,url"`
	Bio         string     `json:"bio" binding:"omitempty,max=500"`
}

// ProfileResponse is returned by GET /profile
// @Description Authenticated user's profile: id, email, phone, role, status, name, gender, date of birth, profile picture URL, bio.
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
