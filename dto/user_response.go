package dto

import (
	"time"

	"github.com/snykk/beego-presence-api/models"
)

// RegisterResponse represents the structure of a register response
// @Description RegisterResponse represents the structure of a register response
type RegisterResponse struct {
	Id           int       `json:"id" example:"1"`                            // Unique identifier of the user
	Name         string    `json:"name" example:"Najib Fikri"`                // Name of the user
	Email        string    `json:"email" example:"najibfikri@gmail.com"`      // Email of the user
	DepartmentId int       `json:"department_id" example:"1"`                 // ForeignKey to Department
	CreatedAt    time.Time `json:"created_at" example:"2024-12-01T00:00:00Z"` // Time when the user was created
	UpdatedAt    time.Time `json:"updated_at" example:"2024-12-01T00:00:00Z"` // Time when the user was updated
}

func FromUserModelToRegisterResponse(u *models.User) *RegisterResponse {
	return &RegisterResponse{
		Id:           u.Id,
		Name:         u.Name,
		Email:        u.Email,
		DepartmentId: u.Department.Id,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
	}
}

// LoginResponse represents the structure of a login response
// @Description LoginResponse represents the structure of a login response
type LoginResponse struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6Im5hamliZmlrcmlAZ21haWwuY29tIiwiZXhwIjoxNzM1NDY0NTMzLCJpYXQiOjE3MzUzNzgxMzMsImlzcyI6ImJlZWdvLXByZXNlbmNlLWFwaSIsInJvbGUiOiJFTVBMT1lFRSIsInN1YiI6NH0.taLD0Rn2rllB4QW4ArHFbylBs2thl9KGC-wBXBDtJN4"` // JWT token
}

// UserResponse represents the structure of a user response
// @Description UserResponse represents the structure of a user response
type UserResponse struct {
	Id           int                 `json:"id" example:"1"`                             // Unique identifier of the user
	Name         string              `json:"name" example:"Najib Fikri"`                 // Name of the user
	Email        string              `json:"email" example:"najibfikri@gmail.com"`       // Email of the user
	DepartmentId *int                `json:"department_id,omitempty" example:"1"`        // ForeignKey to Department
	Department   *DepartmentResponse `json:"department,omitempty" example:"Engineering"` // Department of the user
	Presences    []*PresenceResponse `json:"presences,omitempty"`                        // Reverse relationship with Presence
	ScheduleId   *int                `json:"schedule_id,omitempty" example:"1"`          // ForeignKey to Schedule
	Schedule     *ScheduleResponse   `json:"schedule,omitempty" example:"Schedule"`      // Schedule of the user
	CreatedAt    time.Time           `json:"created_at" example:"2024-12-01T00:00:00Z"`  // Time when the user was created
	UpdatedAt    time.Time           `json:"updated_at" example:"2024-12-01T00:00:00Z"`  // Time when the user was updated
}

func setScheduleIfNotNull(ms *models.Schedule) *int {
	if ms != nil {
		return &ms.Id
	}
	return nil
}

func FromUserModelToUserResponse(u *models.User, isIncludeDepartment, isIncludePresenceList, isIncludeSchedule bool) *UserResponse {
	userResponse := &UserResponse{
		Id:           u.Id,
		Name:         u.Name,
		Email:        u.Email,
		DepartmentId: &u.Department.Id,
		ScheduleId:   setScheduleIfNotNull(u.Schedule),
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
	}

	if isIncludeDepartment {
		userResponse.DepartmentId = nil
		userResponse.Department = FromDepartmentModelToDepartmentResponse(u.Department, false, false)
	}

	if isIncludePresenceList {
		userResponse.Presences = FromPresenceModelListToPresenceResponseList(u.Presences, false, true)
	}

	if isIncludeSchedule && u.Schedule != nil {
		userResponse.ScheduleId = nil
		userResponse.Schedule = FromScheduleModelToScheduleResponse(u.Schedule, false, false, false)
	}
	return userResponse
}

func FromUserModelListToUserResponseList(users []*models.User, isIncludeDepartment, isIncludePresenceList, isIncludeSchedule bool) []*UserResponse {
	var result []*UserResponse

	for _, val := range users {
		result = append(result, FromUserModelToUserResponse(val, isIncludeDepartment, isIncludePresenceList, isIncludeSchedule))
	}

	return result
}
