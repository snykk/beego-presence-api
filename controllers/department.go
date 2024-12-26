package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/snykk/beego-presence-api/helpers"
	"github.com/snykk/beego-presence-api/models"

	beego "github.com/beego/beego/v2/server/web"
)

type DepartmentController struct {
	beego.Controller
}

// @router /departments [get]
func (c *DepartmentController) GetAll() {
	departments, err := models.GetAllDepartments()
	if err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusInternalServerError, "Failed to fetch departments", err)
		return
	}
	helpers.SuccessResponse(c.Ctx.ResponseWriter, http.StatusOK, "Departments retrieved successfully", departments)
}

// @router /departments/:id [get]
func (c *DepartmentController) GetById() {
	id, _ := c.GetInt(":id")
	department, err := models.GetDepartmentById(id)
	if err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusNotFound, "Department not found", err)
		return
	}
	helpers.SuccessResponse(c.Ctx.ResponseWriter, http.StatusOK, "Department retrieved successfully", department)
}

// @router /departments [post]
func (c *DepartmentController) Create() {
	var department models.Department
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &department); err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusBadRequest, "Invalid input", err)
		return
	}

	if err := models.CreateDepartment(&department); err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusInternalServerError, "Failed to create department", err)
		return
	}
	helpers.SuccessResponse(c.Ctx.ResponseWriter, http.StatusCreated, "Department created successfully", department)
}

// @router /departments/:id [put]
func (c *DepartmentController) Update() {
	id, _ := c.GetInt(":id")
	department, err := models.GetDepartmentById(id)
	if err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusNotFound, "Department not found", err)
		return
	}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &department); err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusBadRequest, "Invalid input", err)
		return
	}

	if err := models.UpdateDepartment(&department); err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusInternalServerError, "Failed to update department", err)
		return
	}
	helpers.SuccessResponse(c.Ctx.ResponseWriter, http.StatusOK, "Department updated successfully", department)
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
