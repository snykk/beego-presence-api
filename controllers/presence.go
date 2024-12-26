package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/snykk/beego-presence-api/dto"
	"github.com/snykk/beego-presence-api/helpers"
	"github.com/snykk/beego-presence-api/models"

	beego "github.com/beego/beego/v2/server/web"
)

type PresenceController struct {
	beego.Controller
}

// @router /presences [get]
func (c *PresenceController) GetAll() {
	isIncludeUser, err := c.GetBool("isIncludeUser", false) // Default to false if not provided
	if err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusBadRequest, "Invalid value for isIncludeUser", err)
		return
	}

	isIncludeSchedule, err := c.GetBool("isIncludeSchedule", false) // Default to false if not provided
	if err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusBadRequest, "Invalid value for isIncludeSchedule", err)
		return
	}

	presences, err := models.GetAllPresences()
	if err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusInternalServerError, "Failed to fetch presences", err)
		return
	}
	helpers.SuccessResponse(c.Ctx.ResponseWriter, http.StatusOK, "Presences retrieved successfully", dto.FromPresenceModelListToPresenceResponseList(presences, isIncludeUser, isIncludeSchedule))
}

// @router /presences/:id [get]
func (c *PresenceController) GetById() {
	isIncludeUser, err := c.GetBool("isIncludeUser", false) // Default to false if not provided
	if err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusBadRequest, "Invalid value for isIncludeUser", err)
		return
	}

	isIncludeSchedule, err := c.GetBool("isIncludeSchedule", false) // Default to false if not provided
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
	var presence *models.Presence
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, presence); err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusBadRequest, "Invalid input", err)
		return
	}

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

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, presence); err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusBadRequest, "Invalid input", err)
		return
	}

	if err := models.UpdatePresence(presence); err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusInternalServerError, "Failed to update presence", err)
		return
	}
	helpers.SuccessResponse(c.Ctx.ResponseWriter, http.StatusOK, "Presence updated successfully", dto.FromPresenceModelToPresenceResponse(presence, false, false))
}

// @router /presences/:id [delete]
func (c *PresenceController) Delete() {
	id, _ := c.GetInt(":id")
	if err := models.DeletePresence(id); err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusInternalServerError, "Failed to delete presence", err)
		return
	}
	helpers.SuccessResponse(c.Ctx.ResponseWriter, http.StatusNoContent, "Presence deleted successfully", nil)
}
