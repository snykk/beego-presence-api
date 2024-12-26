package dto

import (
	"time"

	"github.com/snykk/beego-presence-api/models"
)

type PresenceResponse struct {
	Id         int               `json:"id"`
	UserId     *int              `json:"user_id,omitempty"`
	User       *UserResponse     `json:"user,omitempty"`
	Scheduleid *int              `json:"schedule_id,omitempty"`
	Schedule   *ScheduleResponse `json:"schedule,omitempty"`
	Type       string            `json:"type"`
	Status     string            `json:"status"`
	CreatedAt  time.Time         `json:"created_at"`
	UpdatedAt  time.Time         `json:"updated_at"`
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
