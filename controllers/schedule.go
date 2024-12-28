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

// ScheduleController handles operations related to schedules
type ScheduleController struct {
	beego.Controller
}

// URLMapping maps HTTP methods to controller functions
// This function binds the URLs for each handler to its corresponding method.
// It helps Beego framework know which function to call for a given route.
func (c *ScheduleController) URLMapping() {
	c.Mapping("GetAll", c.GetAll)   // Maps GET /schedules to GetAll method for retrieving all schedules
	c.Mapping("GetById", c.GetById) // Maps GET /schedules/:id to GetById method for retrieving a specific schedule by ID
	c.Mapping("Create", c.Create)   // Maps POST /schedules to Create method for adding a new schedule
	c.Mapping("Update", c.Update)   // Maps PUT /schedules/:id to Update method for updating an existing schedule by ID
	c.Mapping("Delete", c.Delete)   // Maps DELETE /schedules/:id to Delete method for deleting a specific schedule by ID
}

// @Title Get All Schedules
// @Description Fetch all schedules with optional related data (department, user presence, user list)
// @Accept  json
// @Produce  json
// @Param isIncludeDepartment query bool false "Include department data"
// @Param isIncludeUser query bool false "Include user presence list"
// @Param isIncludeUserList query bool false "Include user list"
// @Success 200 {object} dto.ScheduleResponse "Schedules retrieved successfully"
// @Failure 400 Invalid query parameters
// @Failure 500 Failed to fetch schedules
// @router / [get]
func (c *ScheduleController) GetAll() {
	// Read and parse the query parameters for optional data inclusion.
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

	// Fetch the list of schedules, passing flags for related data.
	schedules, err := models.GetAllSchedules(isIncludePresenceList, isIncludeUserList)
	if err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusInternalServerError, "Failed to fetch schedules", err)
		return
	}

	// Return the fetched schedules in the response.
	helpers.SuccessResponse(c.Ctx.ResponseWriter, http.StatusOK, "Schedules retrieved successfully", dto.FromScheduleModelListToScheduleResponseList(schedules, isIncludeDepartment, isIncludePresenceList, isIncludeUserList))
}

// @Title Get Schedule By ID
// @Description Fetch a schedule by its ID with optional related data (department, user presence, user list)
// @Accept  json
// @Produce  json
// @Param id path int true "Schedule ID"
// @Param isIncludeDepartment query bool false "Include department data"
// @Param isIncludeUser query bool false "Include user presence list"
// @Param isIncludeUserList query bool false "Include user list"
// @Success 200 {object} dto.ScheduleResponse "Schedule retrieved successfully"
// @Failure 400 Invalid query parameters
// @Failure 404 Schedule not found
// @Failure 500 Failed to fetch schedule
// @router /:id [get]
func (c *ScheduleController) GetById() {
	// Read and parse the query parameters for optional data inclusion.
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

	// Fetch the schedule by ID from the database.
	id, _ := c.GetInt(":id")
	schedule, err := models.GetScheduleById(id, isIncludePresenceList, isIncludeUserList)
	if err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusNotFound, "Schedule not found", err)
		return
	}

	// Return the fetched schedule in the response.
	helpers.SuccessResponse(c.Ctx.ResponseWriter, http.StatusOK, "Schedule retrieved successfully", dto.FromScheduleModelToScheduleResponse(schedule, isIncludeDepartment, isIncludePresenceList, isIncludeUserList))
}

// @Title Create Schedule
// @Description Create a new schedule based on the provided data
// @Accept  json
// @Produce  json
// @Param scheduleRequest body dto.ScheduleRequest true "Schedule Data"
// @Success 201 {object} dto.ScheduleResponse "Schedule created successfully"
// @Failure 400 Invalid input data
// @Failure 500 Failed to create schedule
// @router  [post]
func (c *ScheduleController) Create() {
	// Parse the request body into a ScheduleRequest DTO
	var req dto.ScheduleRequest
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusBadRequest, "Invalid input", err)
		return
	}

	// Validate the input data in the request
	if errorsMap, err := helpers.ValidatePayloads(req); err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusBadRequest, constants.ErrValidationMessage, errorsMap)
		return
	}

	// Fetch the department by ID associated with the schedule
	department, err := models.GetDepartmentById(req.DepartmentId, false, false)
	if department == nil && err != nil {
		// Return error if department is not found
		if err == orm.ErrNoRows {
			helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusNotFound, fmt.Sprintf("Failed to fetch department with id %d", req.DepartmentId), fmt.Errorf("department '%d' not found", req.DepartmentId))
			return
		}
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusInternalServerError, fmt.Sprintf("Failed to fetch department with id %d", req.DepartmentId), err)
		return
	}

	// Convert the request to a schedule model and assign the department
	schedule := req.ToScheduleModel()
	schedule.Department = department

	// Create the new schedule in the database
	if err := models.CreateSchedule(schedule); err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusInternalServerError, "Failed to create schedule", err)
		return
	}

	// Return the created schedule in the response.
	helpers.SuccessResponse(c.Ctx.ResponseWriter, http.StatusCreated, "Schedule created successfully", dto.FromScheduleModelToScheduleResponse(schedule, false, false, false))
}

// @Title Update Schedule
// @Description Update an existing schedule by ID with the provided data
// @Accept  json
// @Produce  json
// @Param id path int true "Schedule ID"
// @Param scheduleRequest body dto.ScheduleRequest true "Schedule Data"
// @Success 200 {object} dto.ScheduleResponse "Schedule updated successfully"
// @Failure 400 Invalid input data
// @Failure 404 Schedule not found
// @Failure 500 Failed to update schedule
// @router /:id [put]
func (c *ScheduleController) Update() {
	// Fetch the schedule by ID to check if it exists
	id, _ := c.GetInt(":id")
	existedSchedule, err := models.GetScheduleById(id, false, false)
	if existedSchedule == nil && err != nil {
		if err == orm.ErrNoRows {
			helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusNotFound, fmt.Sprintf("Failed to fetch schedule with id %d", id), fmt.Errorf("schedule '%d' not found", id))
			return
		}
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusInternalServerError, fmt.Sprintf("Failed to fetch schedule with id %d", id), err)
		return
	}

	// Parse the request body into a ScheduleRequest DTO for updates
	var req dto.ScheduleRequest
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusBadRequest, "Invalid input", err)
		return
	}

	// Validate the input data
	if errorsMap, err := helpers.ValidatePayloads(req); err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusBadRequest, constants.ErrValidationMessage, errorsMap)
		return
	}

	// Fetch the department associated with the schedule
	department, err := models.GetDepartmentById(req.DepartmentId, false, false)
	if department == nil && err != nil {
		if err == orm.ErrNoRows {
			helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusNotFound, fmt.Sprintf("Failed to fetch department with id %d", req.DepartmentId), fmt.Errorf("department '%d' not found", req.DepartmentId))
			return
		}
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusInternalServerError, fmt.Sprintf("Failed to fetch department with id %d", req.DepartmentId), err)
		return
	}

	// Apply the changes to the existing schedule
	updatedSchedule := req.ToScheduleModelWithValue(existedSchedule, department)

	// Save the updated schedule in the database
	if err := models.UpdateSchedule(updatedSchedule); err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusInternalServerError, "Failed to update schedule", err)
		return
	}

	// Return the updated schedule in the response.
	helpers.SuccessResponse(c.Ctx.ResponseWriter, http.StatusOK, "Schedule updated successfully", dto.FromScheduleModelToScheduleResponse(updatedSchedule, false, false, false))
}

// @Title Delete Schedule
// @Description Delete an existing schedule by ID
// @Accept  json
// @Produce  json
// @Param id path int true "Schedule ID"
// @Success 200 {string} string "Schedule deleted successfully"
// @Failure 400 Invalid schedule ID
// @Failure 404 Schedule not found
// @Failure 500 Failed to delete schedule
// @router /:id [delete]
func (c *ScheduleController) Delete() {
	// Read the schedule ID from the URL parameter
	id, err := c.GetInt(":id")
	if err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusBadRequest, "Invalid schedule id", err)
		return
	}

	// Delete the schedule from the database
	affectedRows, err := models.DeleteSchedule(id)
	if err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusInternalServerError, "Failed to delete schedule", err)
		return
	}

	// If no rows were affected, the schedule was not found
	if affectedRows == 0 {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusNotFound, "Schedule not found", fmt.Errorf("schedule '%d' not found", id))
		return
	}

	// Return success response indicating schedule was deleted.
	helpers.SuccessResponse(c.Ctx.ResponseWriter, http.StatusOK, "Schedule deleted successfully", nil)
}
