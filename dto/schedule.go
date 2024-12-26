package dto

import (
	"time"

	"github.com/snykk/beego-presence-api/models"
)

type ScheduleResponse struct {
	Id           int                 `json:"id"`
	Name         string              `json:"name"`
	DepartmentId *int                `json:"department_id,omitempty"`
	Department   *DepartmentResponse `json:"department,omitempty"` // ForeignKey to Department
	InTime       string              `json:"in_time"`
	OutTime      string              `json:"out_time"`
	Presences    []*PresenceResponse `json:"presences,omitempty"` // Reverse relationship with Presence
	Users        []*UserResponse     `json:"users,omitempty"`
	CreatedAt    time.Time           `json:"created_at"`
	UpdatedAt    time.Time           `json:"updated_at"`
}

func FromScheduleModelToScheduleResponse(s *models.Schedule, isIncludeDepartment, isIncludePresenceList, isIncludeUserList bool) *ScheduleResponse {
	scheduleResponse := &ScheduleResponse{
		Id:           s.Id,
		Name:         s.Name,
		DepartmentId: &s.Department.Id,
		InTime:       s.InTime,
		OutTime:      s.OutTime,
		CreatedAt:    s.CreatedAt,
		UpdatedAt:    s.UpdatedAt,
	}

	if isIncludeDepartment {
		scheduleResponse.DepartmentId = nil
		scheduleResponse.Department = FromDepartmentModelToDepartmentResponse(s.Department, false, false)
	}

	if isIncludePresenceList {
		scheduleResponse.Presences = FromPresenceModelListToPresenceResponseList(s.Presences, false, false)
	}

	if isIncludeUserList {
		scheduleResponse.Users = FromUserModelListToUserResponseList(s.Users, false, false, false)
	}

	return scheduleResponse
}

func FromScheduleModelListToScheduleResponseList(schedules []*models.Schedule, isIncludeDepartment, isIncludePresenceList, isIncludeUserList bool) []*ScheduleResponse {
	var result []*ScheduleResponse

	for _, val := range schedules {
		result = append(result, FromScheduleModelToScheduleResponse(val, isIncludeDepartment, isIncludePresenceList, isIncludeUserList))
	}

	return result
}
