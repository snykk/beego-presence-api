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

type DepartmentController struct {
	beego.Controller
}

func (c *DepartmentController) GetAll() {
	isIncludeUserList, err := c.GetBool("isIncludeUserList", false) // Default to false if not provided
	if err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusBadRequest, "Invalid value for isIncludeUserList", err)
		return
	}

	isIncludeScheduleList, err := c.GetBool("isIncludeScheduleList", false) // Default to false if not provided
	if err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusBadRequest, "Invalid value for isIncludeScheduleList", err)
		return
	}

	// Fetch departments
	departments, err := models.GetAllDepartments(isIncludeUserList, isIncludeScheduleList)
	if err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusInternalServerError, "Failed to fetch departments", err)
		return
	}

	// Process response with the query params
	helpers.SuccessResponse(
		c.Ctx.ResponseWriter,
		http.StatusOK,
		"Departments retrieved successfully",
		dto.FromDepartmentModelListToDepartmentResponseList(departments, isIncludeUserList, isIncludeScheduleList),
	)
}

// @router /departments/:id [get]
func (c *DepartmentController) GetById() {
	isIncludeUserList, err := c.GetBool("isIncludeUserList", false) // Default to false if not provided
	if err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusBadRequest, "Invalid value for isIncludeUserList", err)
		return
	}

	isIncludeScheduleList, err := c.GetBool("isIncludeScheduleList", false) // Default to false if not provided
	if err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusBadRequest, "Invalid value for isIncludeScheduleList", err)
		return
	}

	id, _ := c.GetInt(":id")
	department, err := models.GetDepartmentById(id, isIncludeUserList, isIncludeScheduleList)
	if err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusNotFound, "Department not found", err)
		return
	}
	helpers.SuccessResponse(c.Ctx.ResponseWriter, http.StatusOK, "Department retrieved successfully", dto.FromDepartmentModelToDepartmentResponse(department, isIncludeUserList, isIncludeScheduleList))
}

// @router /departments [post]
func (c *DepartmentController) Create() {
	var req *dto.DepartmentRequest
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, req); err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusBadRequest, "Invalid input", err)
		return
	}

	if errorsMap, err := helpers.ValidatePayloads(req); err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusBadRequest, constants.ErrValidationMessage, errorsMap)
		return
	}

	department := req.ToDepartmentModel()

	if err := models.CreateDepartment(department); err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusInternalServerError, "Failed to create department", err)
		return
	}
	helpers.SuccessResponse(c.Ctx.ResponseWriter, http.StatusCreated, "Department created successfully", dto.FromDepartmentModelToDepartmentResponse(department, false, false))
}

// @router /departments/:id [put]
func (c *DepartmentController) Update() {
	id, _ := c.GetInt(":id")
	existedDepartment, err := models.GetDepartmentById(id, false, false)
	if existedDepartment == nil && err != nil {
		if err == orm.ErrNoRows {
			helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusNotFound, fmt.Sprintf("Failed to fetch department with id %d", id), fmt.Errorf("department '%d' not found", id))
			return
		}
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusNotFound, fmt.Sprintf("Failed to fetch department with id %d", id), err)
		return
	}

	var req *dto.DepartmentRequest
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, req); err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusBadRequest, "Invalid input", err)
		return
	}

	if errorsMap, err := helpers.ValidatePayloads(req); err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusBadRequest, constants.ErrValidationMessage, errorsMap)
		return
	}

	department := req.ToDepartmentModelWithValue(existedDepartment)

	if err := models.UpdateDepartment(department); err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusInternalServerError, "Failed to update department", err)
		return
	}
	helpers.SuccessResponse(c.Ctx.ResponseWriter, http.StatusOK, "Department updated successfully", dto.FromDepartmentModelToDepartmentResponse(department, false, false))
}

// @router /departments/:id [delete]
func (c *DepartmentController) Delete() {
	id, err := c.GetInt(":id")
	if err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusBadRequest, "Invalid department id", err)
		return
	}

	affectedRows, err := models.DeleteDepartment(id)
	if err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusInternalServerError, "Failed to delete department", err)
		return
	}

	if affectedRows == 0 {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusNotFound, "Department not found", fmt.Errorf("department '%d' not found", id))
		return
	}

	helpers.SuccessResponse(c.Ctx.ResponseWriter, http.StatusOK, "Department deleted successfully", nil)
}
