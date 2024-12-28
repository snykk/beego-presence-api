package dto

import (
	"time"

	"github.com/snykk/beego-presence-api/models"
)

// ScheduleResponse represents the structure of a schedule response
// @Description ScheduleResponse represents the structure of a schedule response
type ScheduleResponse struct {
	Id           int                 `json:"id" example:"1"` // Unique identifier of the schedule
	Name         string              `json:"name" example:"Morning Shift"`
	DepartmentId *int                `json:"department_id,omitempty" example:"1"`       // ForeignKey to Department
	Department   *DepartmentResponse `json:"department,omitempty"`                      // Department of the schedule
	InTime       string              `json:"in_time" example:"08:00:00"`                // Time when the schedule starts
	OutTime      string              `json:"out_time" example:"16:00:00"`               // Time when the schedule ends
	Presences    []*PresenceResponse `json:"presences,omitempty"`                       // Reverse relationship with Presence
	Users        []*UserResponse     `json:"users,omitempty"`                           // Reverse relationship with User
	CreatedAt    time.Time           `json:"created_at" example:"2021-01-01T00:00:00Z"` // Time when the schedule was created
	UpdatedAt    time.Time           `json:"updated_at" example:"2021-01-01T00:00:00Z"` // Time when the schedule was updated
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
