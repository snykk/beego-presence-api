package dto

import (
	"time"

	"github.com/snykk/beego-presence-api/models"
)

type RegisterRequest struct {
	Name       string `json:"name"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	Department int    `json:"department"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	Id           int       `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	DepartmentId int       `json:"department_id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func FromUserModelToRegisterResponse(u models.User) RegisterResponse {
	return RegisterResponse{
		Id:           u.Id,
		Name:         u.Name,
		Email:        u.Email,
		DepartmentId: u.Department.Id,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
	}
}

type UserResponse struct {
	Id           int                 `json:"id"`
	Name         string              `json:"name"`
	Email        string              `json:"email"`
	DepartmentId *int                `json:"department_id,omitempty"`
	Department   *DepartmentResponse `json:"department,omitempty"`
	Presences    []*PresenceResponse `json:"presences,omitempty"`
	ScheduleId   *int                `json:"schedule_id,omitempty"`
	Schedule     *ScheduleResponse   `json:"schedule,omitempty"`
	CreatedAt    time.Time           `json:"created_at"`
	UpdatedAt    time.Time           `json:"updated_at"`
}

func FromUserModelToUserResponse(u *models.User, isIncludeDepartment, isIncludePresenceList, isIncludeSchedule bool) *UserResponse {
	userResponse := &UserResponse{
		Id:           u.Id,
		Name:         u.Name,
		Email:        u.Email,
		DepartmentId: &u.Department.Id,
		ScheduleId:   &u.Schedule.Id,
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

	if isIncludeSchedule {
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
