package dto

import (
	"time"

	"github.com/snykk/beego-presence-api/models"
)

// PresenceResponse represents the structure of a presence response
// @Description PresenceResponse represents the structure of a presence response
type PresenceResponse struct {
	Id         int               `json:"id" example:"1"` // Presence ID
	UserId     *int              `json:"user_id,omitempty" example:"1"`
	User       *UserResponse     `json:"user,omitempty" example:"1"`
	Scheduleid *int              `json:"schedule_id,omitempty" example:"1"`
	Schedule   *ScheduleResponse `json:"schedule,omitempty" example:"1"`
	Type       string            `json:"type" example:"in"`                         // Presence type
	Status     string            `json:"status" example:"ontime"`                   // Presence status
	CreatedAt  time.Time         `json:"created_at" example:"2024-12-01T00:00:00Z"` // Creation timestamp
	UpdatedAt  time.Time         `json:"updated_at" example:"2024-12-02T00:00:00Z"` // Last update timestamp
}

func FromPresenceModelToPresenceResponse(u *models.Presence, isIncludeUser, isIncludeSchedule bool) *PresenceResponse {
	presenceResponse := &PresenceResponse{
		Id:         u.Id,
		UserId:     &u.User.Id,
		Scheduleid: &u.Schedule.Id,
		Type:       u.Type,
		Status:     u.Status,
		CreatedAt:  u.CreatedAt,
		UpdatedAt:  u.UpdatedAt,
	}

	if isIncludeUser {
		presenceResponse.UserId = nil
		presenceResponse.User = FromUserModelToUserResponse(u.User, false, false, false)
	}

	if isIncludeSchedule {
		presenceResponse.Scheduleid = nil
		presenceResponse.Schedule = FromScheduleModelToScheduleResponse(u.Schedule, true, false, false)
	}
	return presenceResponse
}

func FromPresenceModelListToPresenceResponseList(presences []*models.Presence, isIncludeUser, isIncludeSchedule bool) []*PresenceResponse {
	var result []*PresenceResponse

	for _, val := range presences {
		result = append(result, FromPresenceModelToPresenceResponse(val, isIncludeUser, isIncludeSchedule))
	}

	return result
}
