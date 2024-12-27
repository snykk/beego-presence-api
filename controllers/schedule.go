package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/snykk/beego-presence-api/constants"
	"github.com/snykk/beego-presence-api/dto"
	"github.com/snykk/beego-presence-api/helpers"
	"github.com/snykk/beego-presence-api/models"

	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"
)

type ScheduleController struct {
	beego.Controller
}

// @router /schedules [get]
func (c *ScheduleController) GetAll() {
	isIncludeDepartment, err := c.GetBool("isIncludeDepartment", false) // Default to false if not provided
	if err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusBadRequest, "Invalid value for isIncludeDepartment", err)
		return
	}

	isIncludePresenceList, err := c.GetBool("isIncludeUser", false) // Default to false if not provided
	if err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusBadRequest, "Invalid value for isIncludeUser", err)
		return
	}

	isIncludeUserList, err := c.GetBool("isIncludeUserList", false) // Default to false if not provided
	if err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusBadRequest, "Invalid value for isIncludeUserList", err)
		return
	}

	schedules, err := models.GetAllSchedules(isIncludePresenceList, isIncludeUserList)
	if err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusInternalServerError, "Failed to fetch schedules", err)
		return
	}
	helpers.SuccessResponse(c.Ctx.ResponseWriter, http.StatusOK, "Schedules retrieved successfully", dto.FromScheduleModelListToScheduleResponseList(schedules, isIncludeDepartment, isIncludePresenceList, isIncludeUserList))
}

// @router /schedules/:id [get]
func (c *ScheduleController) GetById() {
	isIncludeDepartment, err := c.GetBool("isIncludeDepartment", false) // Default to false if not provided
	if err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusBadRequest, "Invalid value for isIncludeDepartment", err)
		return
	}

	isIncludePresenceList, err := c.GetBool("isIncludeUser", false) // Default to false if not provided
	if err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusBadRequest, "Invalid value for isIncludeUser", err)
		return
	}

	isIncludeUserList, err := c.GetBool("isIncludeUserList", false) // Default to false if not provided
	if err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusBadRequest, "Invalid value for isIncludeUserList", err)
		return
	}

	id, _ := c.GetInt(":id")
	schedule, err := models.GetScheduleById(id, isIncludePresenceList, isIncludeUserList)
	if err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusNotFound, "Schedule not found", err)
		return
	}
	helpers.SuccessResponse(c.Ctx.ResponseWriter, http.StatusOK, "Schedule retrieved successfully", dto.FromScheduleModelToScheduleResponse(schedule, isIncludeDepartment, isIncludePresenceList, isIncludeUserList))
}

// @router /schedules [post]
func (c *ScheduleController) Create() {
	var req dto.ScheduleRequest
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusBadRequest, "Invalid input", err)
		return
	}

	if errorsMap, err := helpers.ValidatePayloads(req); err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusBadRequest, constants.ErrValidationMessage, errorsMap)
		return
	}

	department, err := models.GetDepartmentById(req.DepartmentId, false, false)
	if department == nil && err != nil {
		if err == orm.ErrNoRows {
			helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusNotFound, fmt.Sprintf("Failed to fetch department with id %d", req.DepartmentId), fmt.Errorf("department '%d' not found", req.DepartmentId))
			return
		}
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusNotFound, fmt.Sprintf("Failed to fetch department with id %d", req.DepartmentId), err)
		return
	}

	schedule := req.ToScheduleModel()
	schedule.Department = department

	if err := models.CreateSchedule(schedule); err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusInternalServerError, "Failed to create schedule", err)
		return
	}
	helpers.SuccessResponse(c.Ctx.ResponseWriter, http.StatusCreated, "Schedule created successfully", dto.FromScheduleModelToScheduleResponse(schedule, false, false, false))
}

// @router /schedules/:id [put]
func (c *ScheduleController) Update() {
	id, _ := c.GetInt(":id")
	var req dto.ScheduleRequest
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusBadRequest, "Invalid input", err)
		return
	}

	if errorsMap, err := helpers.ValidatePayloads(req); err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusBadRequest, constants.ErrValidationMessage, errorsMap)
		return
	}

	existedSchedule, err := models.GetScheduleById(id, false, false)
	if existedSchedule == nil && err != nil {
		if err == orm.ErrNoRows {
			helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusNotFound, fmt.Sprintf("Failed to fetch schedule with id %d", id), fmt.Errorf("schedule '%d' not found", id))
			return
		}
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusNotFound, fmt.Sprintf("Failed to fetch schedule with id %d", id), err)
		return
	}

	department, err := models.GetDepartmentById(req.DepartmentId, false, false)
	if department == nil && err != nil {
		if err == orm.ErrNoRows {
			helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusNotFound, fmt.Sprintf("Failed to fetch department with id %d", req.DepartmentId), fmt.Errorf("department '%d' not found", req.DepartmentId))
			return
		}
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusNotFound, fmt.Sprintf("Failed to fetch department with id %d", req.DepartmentId), err)
		return
	}

	updatedSchedule := req.ToScheduleModelWithValue(existedSchedule, department)

	if err := models.UpdateSchedule(updatedSchedule); err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusInternalServerError, "Failed to update schedule", err)
		return
	}
	helpers.SuccessResponse(c.Ctx.ResponseWriter, http.StatusOK, "Schedule updated successfully", dto.FromScheduleModelToScheduleResponse(updatedSchedule, false, false, false))
}

// @router /schedules/:id [delete]
func (c *ScheduleController) Delete() {
	id, err := c.GetInt(":id")
	if err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusBadRequest, "Invalid schedule id", err)
		return
	}

	affectedRows, err := models.DeleteSchedule(id)
	if err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusInternalServerError, "Failed to delete schedule", err)
		return
	}

	if affectedRows == 0 {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusNotFound, "Schedule not found", fmt.Errorf("schedule '%d' not found", id))
		return
	}

	helpers.SuccessResponse(c.Ctx.ResponseWriter, http.StatusOK, "Schedule deleted successfully", nil)
}
