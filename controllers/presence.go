package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/snykk/beego-presence-api/constants"
	"github.com/snykk/beego-presence-api/dto"
	"github.com/snykk/beego-presence-api/helpers"
	"github.com/snykk/beego-presence-api/models"

	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"
)

type PresenceController struct {
	beego.Controller
}

// @router /presences [get]
func (c *PresenceController) GetAll() {
	userRole, ok := c.Ctx.Input.GetData(constants.CtxAuthenticatedUserRole).(string)
	if !ok {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusUnauthorized, "Bad context", errors.New("can't retrieve user role from context"))
		return
	}

	userId, ok := c.Ctx.Input.GetData(constants.CtxAuthenticatedUserId).(int)
	if !ok {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusUnauthorized, "Bad context", errors.New("can't retrieve user role from context"))
		return
	}

	isIncludeUser, err := c.GetBool("isIncludeUser", false)
	if err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusBadRequest, "Invalid value for isIncludeUser", err)
		return
	}

	isIncludeSchedule, err := c.GetBool("isIncludeSchedule", false)
	if err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusBadRequest, "Invalid value for isIncludeSchedule", err)
		return
	}
	var presences []*models.Presence
	if userRole == constants.RoleAdmin {
		presences, err = models.GetAllPresences()
	} else if userRole == constants.RoleEmployee {
		presences, err = models.GetPresencesByUserId(userId)
	}
	if err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusInternalServerError, "Failed to fetch presences", err)
		return
	}

	helpers.SuccessResponse(c.Ctx.ResponseWriter, http.StatusOK, "Presences retrieved successfully", dto.FromPresenceModelListToPresenceResponseList(presences, isIncludeUser, isIncludeSchedule))
}

// @router /presences/:id [get]
func (c *PresenceController) GetById() {
	isIncludeUser, err := c.GetBool("isIncludeUser", false)
	if err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusBadRequest, "Invalid value for isIncludeUser", err)
		return
	}

	isIncludeSchedule, err := c.GetBool("isIncludeSchedule", false)
	if err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusBadRequest, "Invalid value for isIncludeSchedule", err)
		return
	}

	id, _ := c.GetInt(":id")
	presence, err := models.GetPresenceById(id)
	if err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusNotFound, "Presence not found", err)
		return
	}
	helpers.SuccessResponse(c.Ctx.ResponseWriter, http.StatusOK, "Presence retrieved successfully", dto.FromPresenceModelToPresenceResponse(presence, isIncludeUser, isIncludeSchedule))
}

// @router /presences [post]
func (c *PresenceController) Create() {
	userId, ok := c.Ctx.Input.GetData(constants.CtxAuthenticatedUserId).(int)
	if !ok {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusUnauthorized, "Bad context", errors.New("can't retrieve user role from context"))
		return
	}

	var req dto.PresenceCreateRequest
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusBadRequest, "Invalid input", err)
		return
	}

	if errorsMap, err := helpers.ValidatePayloads(req); err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusBadRequest, constants.ErrValidationMessage, errorsMap)
		return
	}

	user, err := models.GetUserById(userId, false)
	if user == nil && err != nil {
		if err == orm.ErrNoRows {
			helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusNotFound, fmt.Sprintf("Failed to fetch user with id %d", userId), fmt.Errorf("user '%d' not found", userId))
			return
		}
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusInternalServerError, fmt.Sprintf("Failed to fetch user with id %d", req.ScheduleId), err)
		return
	}

	// Check if user is assigned to the schedule
	if user.Schedule.Id != req.ScheduleId {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusBadRequest, "User is not assigned to the schedule", fmt.Errorf("user '%d' is not assigned to the schedule '%d'", userId, req.ScheduleId))
		return
	}

	schedule, err := models.GetScheduleById(req.ScheduleId, false, false)
	if schedule == nil && err != nil {
		if err == orm.ErrNoRows {
			helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusNotFound, fmt.Sprintf("Failed to fetch schedule with id %d", req.ScheduleId), fmt.Errorf("schedule '%d' not found", req.ScheduleId))
			return
		}
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusInternalServerError, fmt.Sprintf("Failed to fetch schedule with id %d", req.ScheduleId), err)
		return
	}

	// Check if presence already exists for the user and type
	currentTime := time.Now()
	exists, err := models.CheckPresenceExistsByUserAndType(userId, req.Type, currentTime)
	if err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusInternalServerError, "Failed to check existing presence", err)
		return
	}
	if exists {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusBadRequest, "Presence already exists for this type and day", fmt.Errorf("user '%d' already has a presence for type '%s' today", userId, req.Type))
		return
	}

	presence := req.ToPresenceModelWithValue(user, schedule)

	// Determine status
	status, err := helpers.DeterminePresenceStatus(presence.Type, schedule.InTime, schedule.OutTime, currentTime, constants.PresenceLateThreshold)
	if err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusInternalServerError, "Failed to determine presence status", err)
		return
	}
	presence.Status = status

	if err := models.CreatePresence(presence); err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusInternalServerError, "Failed to create presence", err)
		return
	}
	helpers.SuccessResponse(c.Ctx.ResponseWriter, http.StatusCreated, "Presence created successfully", dto.FromPresenceModelToPresenceResponse(presence, false, false))
}

// @router /presences/:id [put]
func (c *PresenceController) Update() {
	id, _ := c.GetInt(":id")
	presence, err := models.GetPresenceById(id)
	if err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusNotFound, "Presence not found", err)
		return
	}

	var req dto.PresenceUpdateRequest
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusBadRequest, "Invalid input", err)
		return
	}

	if errorsMap, err := helpers.ValidatePayloads(req); err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusBadRequest, constants.ErrValidationMessage, errorsMap)
		return
	}

	schedule, err := models.GetScheduleById(req.ScheduleId, false, false)
	if schedule == nil && err != nil {
		if err == orm.ErrNoRows {
			helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusNotFound, fmt.Sprintf("Failed to fetch schedule with id %d", req.ScheduleId), fmt.Errorf("schedule '%d' not found", req.ScheduleId))
			return
		}
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusInternalServerError, fmt.Sprintf("Failed to fetch schedule with id %d", req.ScheduleId), err)
		return
	}

	user, err := models.GetUserById(req.UserId, false)
	if user == nil && err != nil {
		if err == orm.ErrNoRows {
			helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusNotFound, fmt.Sprintf("Failed to fetch user with id %d", req.UserId), fmt.Errorf("user '%d' not found", req.UserId))
			return
		}
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusInternalServerError, fmt.Sprintf("Failed to fetch user with id %d", req.ScheduleId), err)
		return
	}

	updatedPresence := req.ToPresenceModelWithValue(presence, user, schedule)

	if err := models.UpdatePresence(updatedPresence); err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusInternalServerError, "Failed to update presence", err)
		return
	}
	helpers.SuccessResponse(c.Ctx.ResponseWriter, http.StatusOK, "Presence updated successfully", dto.FromPresenceModelToPresenceResponse(updatedPresence, false, false))
}

// @router /presences/:id [delete]
func (c *PresenceController) Delete() {
	id, err := c.GetInt(":id")
	if err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusBadRequest, "Invalid presence id", err)
		return
	}

	affectedRows, err := models.DeletePresence(id)
	if err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusInternalServerError, "Failed to delete presence", err)
		return
	}

	if affectedRows == 0 {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusNotFound, "Presence not found", fmt.Errorf("presence '%d' not found", id))
		return
	}

	helpers.SuccessResponse(c.Ctx.ResponseWriter, http.StatusOK, "Presence deleted successfully", nil)
}
