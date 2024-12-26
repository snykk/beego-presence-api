package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/snykk/beego-presence-api/dto"
	"github.com/snykk/beego-presence-api/helpers"
	"github.com/snykk/beego-presence-api/models"

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
	var department *models.Department
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &department); err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusBadRequest, "Invalid input", err)
		return
	}

	if err := models.CreateDepartment(department); err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusInternalServerError, "Failed to create department", err)
		return
	}
	helpers.SuccessResponse(c.Ctx.ResponseWriter, http.StatusCreated, "Department created successfully", dto.FromDepartmentModelToDepartmentResponse(department, false, false))
}

// @router /departments/:id [put]
func (c *DepartmentController) Update() {
	id, _ := c.GetInt(":id")
	department, err := models.GetDepartmentById(id, false, false)
	if err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusNotFound, "Department not found", err)
		return
	}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, department); err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusBadRequest, "Invalid input", err)
		return
	}

	if err := models.UpdateDepartment(department); err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusInternalServerError, "Failed to update department", err)
		return
	}
	helpers.SuccessResponse(c.Ctx.ResponseWriter, http.StatusOK, "Department updated successfully", dto.FromDepartmentModelToDepartmentResponse(department, false, false))
}

// @router /departments/:id [delete]
func (c *DepartmentController) Delete() {
	id, _ := c.GetInt(":id")
	if err := models.DeleteDepartment(id); err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusInternalServerError, "Failed to delete department", err)
		return
	}
	helpers.SuccessResponse(c.Ctx.ResponseWriter, http.StatusNoContent, "Department deleted successfully", nil)
}
