package dto

import "github.com/snykk/beego-presence-api/models"

type PresenceCreateRequest struct {
	ScheduleId int    `json:"schedule_id" validate:"required,min=1"`
	Type       string `json:"type" validate:"required,oneof=in out"`
}

func (p *PresenceCreateRequest) ToPresenceModelWithValue(mu *models.User, ms *models.Schedule) *models.Presence {
	return &models.Presence{
		User:     mu,
		Schedule: ms,
		Type:     p.Type,
	}
}

type PresenceUpdateRequest struct {
	UserId     int    `json:"user_id" validate:"required,min=1"`
	ScheduleId int    `json:"schedule_id" validate:"required,min=1"`
	Type       string `json:"type" validate:"required,oneof=in out"`
	Status     string `json:"status" validate:"required,oneof=ontime late"`
}

func (p *PresenceUpdateRequest) ToPresenceModelWithValue(mp *models.Presence, mu *models.User, ms *models.Schedule) *models.Presence {
	mp.User = mu
	mp.Schedule = ms
	mp.Type = p.Type
	mp.Status = p.Status
	return mp
}
