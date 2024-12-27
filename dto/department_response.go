package dto

import (
	"time"

	"github.com/snykk/beego-presence-api/models"
)

type DepartmentResponse struct {
	Id        int                 `json:"id"`
	Name      string              `json:"name"`
	Users     []*UserResponse     `json:"users,omitempty"`
	Schedules []*ScheduleResponse `json:"schedules,omitempty"`
	CreatedAt time.Time           `json:"created_at"`
	UpdatedAt time.Time           `json:"updated_at"`
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
