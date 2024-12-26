package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/snykk/beego-presence-api/helpers"
	"github.com/snykk/beego-presence-api/models"

	beego "github.com/beego/beego/v2/server/web"
)

type ScheduleController struct {
	beego.Controller
}

// @router /schedules [get]
func (c *ScheduleController) GetAll() {
	schedules, err := models.GetAllSchedules()
	if err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusInternalServerError, "Failed to fetch schedules", err)
		return
	}
	helpers.SuccessResponse(c.Ctx.ResponseWriter, http.StatusOK, "Schedules retrieved successfully", schedules)
}

// @router /schedules/:id [get]
func (c *ScheduleController) GetById() {
	id, _ := c.GetInt(":id")
	schedule, err := models.GetScheduleById(id)
	if err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusNotFound, "Schedule not found", err)
		return
	}
	helpers.SuccessResponse(c.Ctx.ResponseWriter, http.StatusOK, "Schedule retrieved successfully", schedule)
}

// @router /schedules [post]
func (c *ScheduleController) Create() {
	var schedule models.Schedule
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &schedule); err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusBadRequest, "Invalid input", err)
		return
	}

	if err := models.CreateSchedule(&schedule); err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusInternalServerError, "Failed to create schedule", err)
		return
	}
	helpers.SuccessResponse(c.Ctx.ResponseWriter, http.StatusCreated, "Schedule created successfully", schedule)
}

// @router /schedules/:id [put]
func (c *ScheduleController) Update() {
	id, _ := c.GetInt(":id")
	schedule, err := models.GetScheduleById(id)
	if err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusNotFound, "Schedule not found", err)
		return
	}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &schedule); err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusBadRequest, "Invalid input", err)
		return
	}

	if err := models.UpdateSchedule(&schedule); err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusInternalServerError, "Failed to update schedule", err)
		return
	}
	helpers.SuccessResponse(c.Ctx.ResponseWriter, http.StatusOK, "Schedule updated successfully", schedule)
}

// @router /schedules/:id [delete]
func (c *ScheduleController) Delete() {
	id, _ := c.GetInt(":id")
	if err := models.DeleteSchedule(id); err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusInternalServerError, "Failed to delete schedule", err)
		return
	}
	helpers.SuccessResponse(c.Ctx.ResponseWriter, http.StatusOK, "Schedule deleted successfully", nil)
}
