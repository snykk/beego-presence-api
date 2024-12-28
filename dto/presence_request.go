package dto

import "github.com/snykk/beego-presence-api/models"

// PresenceCreateRequest represents the structure of a presence create request
// @Description PresenceCreateRequest represents the structure of a presence create request
type PresenceCreateRequest struct {
	ScheduleId int    `json:"schedule_id" validate:"required,min=1" example:"1"`  // Schedule ID
	Type       string `json:"type" validate:"required,oneof=in out" example:"in"` // Presence type
}

func (p *PresenceCreateRequest) ToPresenceModelWithValue(mu *models.User, ms *models.Schedule) *models.Presence {
	return &models.Presence{
		User:     mu,
		Schedule: ms,
		Type:     p.Type,
	}
}

// PresenceUpdateRequest represents the structure of a presence update request
// @Description PresenceUpdateRequest represents the structure of a presence update request
type PresenceUpdateRequest struct {
	UserId     int    `json:"user_id" validate:"required,min=1" example:"1"`                 // User ID
	ScheduleId int    `json:"schedule_id" validate:"required,min=1" example:"1"`             // Schedule ID
	Type       string `json:"type" validate:"required,oneof=in out" example:"in"`            // Presence type
	Status     string `json:"status" validate:"required,oneof=ontime late" example:"ontime"` // Presence status
}

func (p *PresenceUpdateRequest) ToPresenceModelWithValue(mp *models.Presence, mu *models.User, ms *models.Schedule) *models.Presence {
	mp.User = mu
	mp.Schedule = ms
	mp.Type = p.Type
	mp.Status = p.Status
	return mp
}
