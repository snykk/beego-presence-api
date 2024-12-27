package dto

import "github.com/snykk/beego-presence-api/models"

type ScheduleRequest struct {
	Name         string `json:"name" validate:"required"`
	DepartmentId int    `json:"department_id" validate:"required,min=1"`
	InTime       string `json:"in_time" validate:"required"`
	OutTime      string `json:"out_time" validate:"required"`
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
