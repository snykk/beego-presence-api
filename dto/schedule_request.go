package dto

import "github.com/snykk/beego-presence-api/models"

// ScheduleRequest represents the structure of a presence create request
// @Description ScheduleRequest represents the structure of a presence create request
type ScheduleRequest struct {
	Name         string `json:"name" validate:"required" example:"Morning Shift"`    // Name of the schedule
	DepartmentId int    `json:"department_id" validate:"required,min=1" example:"1"` // ForeignKey to Department
	InTime       string `json:"in_time" validate:"required" example:"08:00:00"`      // Time when the schedule starts
	OutTime      string `json:"out_time" validate:"required" example:"16:00:00"`     // Time when the schedule ends
}

func (s ScheduleRequest) ToScheduleModel() *models.Schedule {
	return &models.Schedule{
		Name:    s.Name,
		InTime:  s.InTime,
		OutTime: s.OutTime,
	}
}

func (s ScheduleRequest) ToScheduleModelWithValue(ms *models.Schedule, md *models.Department) *models.Schedule {
	ms.Name = s.Name
	ms.InTime = s.InTime
	ms.OutTime = s.OutTime
	ms.Department = md
	return ms
}
