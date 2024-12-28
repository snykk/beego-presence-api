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

// DepartmentController handles operations related to departments management
type DepartmentController struct {
	beego.Controller
}

// URLMapping maps HTTP methods to controller functions
// This function binds the URLs for each handler to its corresponding method.
// It helps Beego framework know which function to call for a given route.
func (c *DepartmentController) URLMapping() {
	c.Mapping("GetAll", c.GetAll)   // Maps GET /departments to GetAll method for retrieving all departments
	c.Mapping("GetById", c.GetById) // Maps GET /departments/:id to GetById method for retrieving a specific department by ID
	c.Mapping("Create", c.Create)   // Maps POST /departments to Create method for adding a new department
	c.Mapping("Update", c.Update)   // Maps PUT /departments/:id to Update method for updating an existing department by ID
	c.Mapping("Delete", c.Delete)   // Maps DELETE /departments/:id to Delete method for deleting a specific department by ID
}

// @Title GetAll
// @Description Retrieve all departments, optionally including user and schedule lists.
// @Produce  json
// @Param   isIncludeUserList		query	bool	false		"Include user list in response"
// @Param   isIncludeScheduleList	query	bool	false		"Include schedule list in response"
// @Success 200 {object} dto.DepartmentResponse "Departments retrieved successfully"
// @Failure 400 Bad request
// @Failure 500 Internal server error
// @router / [get]
func (c *DepartmentController) GetAll() {
	// Read query parameters and handle invalid input
	isIncludeUserList, err := c.GetBool("isIncludeUserList", false)
	if err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusBadRequest, "Invalid value for isIncludeUserList", err)
		return
	}

	isIncludeScheduleList, err := c.GetBool("isIncludeScheduleList", false)
	if err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusBadRequest, "Invalid value for isIncludeScheduleList", err)
		return
	}

	// Fetch departments from the model with additional data as needed
	departments, err := models.GetAllDepartments(isIncludeUserList, isIncludeScheduleList)
	if err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusInternalServerError, "Failed to fetch departments", err)
		return
	}

	// Return success response with the department data
	helpers.SuccessResponse(
		c.Ctx.ResponseWriter,
		http.StatusOK,
		"Departments retrieved successfully",
		dto.FromDepartmentModelListToDepartmentResponseList(departments, isIncludeUserList, isIncludeScheduleList),
	)
}

// @Title GetById
// @Description Retrieve a department by its ID, optionally including user and schedule lists.
// @Produce  json
// @Param   id						path	int		true		"Department ID"
// @Param   isIncludeUserList		query	bool	false		"Include user list in response"
// @Param   isIncludeScheduleList	query	bool	false		"Include schedule list in response"
// @Success 200 {object} dto.DepartmentResponse "Department retrieved successfully"
// @Failure 400 Bad request
// @Failure 404 Department not found
// @Failure 500 Internal server error
// @router /:id [get]
func (c *DepartmentController) GetById() {
	// Read query parameters and handle invalid input
	isIncludeUserList, err := c.GetBool("isIncludeUserList", false)
	if err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusBadRequest, "Invalid value for isIncludeUserList", err)
		return
	}

	isIncludeScheduleList, err := c.GetBool("isIncludeScheduleList", false)
	if err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusBadRequest, "Invalid value for isIncludeScheduleList", err)
		return
	}

	// Fetch department by ID
	id, _ := c.GetInt(":id")
	department, err := models.GetDepartmentById(id, isIncludeUserList, isIncludeScheduleList)
	if err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusNotFound, "Department not found", err)
		return
	}

	// Return success response with department data
	helpers.SuccessResponse(c.Ctx.ResponseWriter, http.StatusOK, "Department retrieved successfully", dto.FromDepartmentModelToDepartmentResponse(department, isIncludeUserList, isIncludeScheduleList))
}

// @Title Create
// @Description Create a new department.
// @Accept  json
// @Produce  json
// @Param   body	body	dto.DepartmentRequest	true		"Department data"
// @Success 201 {object} dto.DepartmentResponse "Department created successfully"
// @Failure 400 Invalid input
// @Failure 500 Internal server error
// @router / [post]
func (c *DepartmentController) Create() {
	// Parse request body to department request object
	var req *dto.DepartmentRequest
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, req); err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusBadRequest, "Invalid input", err)
		return
	}

	// Validate payload for any errors
	if errorsMap, err := helpers.ValidatePayloads(req); err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusBadRequest, constants.ErrValidationMessage, errorsMap)
		return
	}

	// Convert the request object to department model and create it in the database
	department := req.ToDepartmentModel()
	if err := models.CreateDepartment(department); err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusInternalServerError, "Failed to create department", err)
		return
	}

	// Return success response with newly created department data
	helpers.SuccessResponse(c.Ctx.ResponseWriter, http.StatusCreated, "Department created successfully", dto.FromDepartmentModelToDepartmentResponse(department, false, false))
}

// @Title Update
// @Description Update an existing department by ID.
// @Accept  json
// @Produce  json
// @Param   id		path	int	true		"Department ID"
// @Param   body	body	dto.DepartmentRequest	true		"Updated department data"
// @Success 200 {object} dto.DepartmentResponse "Department updated successfully"
// @Failure 400 Invalid input
// @Failure 404 Department not found
// @Failure 500 Internal server error
// @router /:id [put]
func (c *DepartmentController) Update() {
	// Get department ID from path parameter
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

	// Parse request body to department request object
	var req *dto.DepartmentRequest
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, req); err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusBadRequest, "Invalid input", err)
		return
	}

	// Validate the update request data
	if errorsMap, err := helpers.ValidatePayloads(req); err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusBadRequest, constants.ErrValidationMessage, errorsMap)
		return
	}

	// Update department in the database
	department := req.ToDepartmentModelWithValue(existedDepartment)
	if err := models.UpdateDepartment(department); err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusInternalServerError, "Failed to update department", err)
		return
	}

	// Return success response with updated department data
	helpers.SuccessResponse(c.Ctx.ResponseWriter, http.StatusOK, "Department updated successfully", dto.FromDepartmentModelToDepartmentResponse(department, false, false))
}

// @Title Delete
// @Description Delete an existing department by ID.
// @Param   id		path	int	true		"Department ID"
// @Success 200 {string} "Department deleted successfully"
// @Failure 400 Invalid department ID
// @Failure 404 Department not found
// @Failure 500 Internal server error
// @router /:id [delete]
func (c *DepartmentController) Delete() {
	// Get department ID from path parameter
	id, err := c.GetInt(":id")
	if err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusBadRequest, "Invalid department id", err)
		return
	}

	// Delete department from database
	affectedRows, err := models.DeleteDepartment(id)
	if err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusInternalServerError, "Failed to delete department", err)
		return
	}

	// Check if the department was found and deleted
	if affectedRows == 0 {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusNotFound, "Department not found", fmt.Errorf("department '%d' not found", id))
		return
	}

	// Return success response indicating department has been deleted
	helpers.SuccessResponse(c.Ctx.ResponseWriter, http.StatusOK, "Department deleted successfully", nil)
}
