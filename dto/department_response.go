package dto

import (
	"time"

	"github.com/snykk/beego-presence-api/models"
)

// DepartmentResponse represents the structure of a department response
// @Description DepartmentResponse represents the structure of a department response
type DepartmentResponse struct {
	Id        int                 `json:"id" example:"1"`                            // Department ID
	Name      string              `json:"name" example:"Human Resources"`            // Department name
	Users     []*UserResponse     `json:"users,omitempty"`                           // List of users in the department
	Schedules []*ScheduleResponse `json:"schedules,omitempty"`                       // List of schedules for the department
	CreatedAt time.Time           `json:"created_at" example:"2023-01-01T00:00:00Z"` // Creation timestamp
	UpdatedAt time.Time           `json:"updated_at" example:"2023-01-02T00:00:00Z"` // Last update timestamp
}

func FromDepartmentModelToDepartmentResponse(d *models.Department, isIncludeUserList, isIncludeScheduleList bool) *DepartmentResponse {
	departmentResponse := &DepartmentResponse{
		Id:        d.Id,
		Name:      d.Name,
		CreatedAt: d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
	}

	if isIncludeUserList {
		departmentResponse.Users = FromUserModelListToUserResponseList(d.Users, false, false, false)
	}

	if isIncludeScheduleList {
		departmentResponse.Schedules = FromScheduleModelListToScheduleResponseList(d.Schedules, false, false, false)
	}

	return departmentResponse
}

func FromDepartmentModelListToDepartmentResponseList(departments []*models.Department, isIncludeUserList, isIncludeScheduleList bool) []*DepartmentResponse {
	var result []*DepartmentResponse

	for _, val := range departments {
		result = append(result, FromDepartmentModelToDepartmentResponse(val, isIncludeUserList, isIncludeScheduleList))
	}

	return result
}
