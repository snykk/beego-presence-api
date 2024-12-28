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

// PresenceController handles requests related to presences (attendance) management.
type PresenceController struct {
	beego.Controller
}

// URLMapping maps routes to specific handler functions for the PresenceController
// This is typically used by the Beego framework to map HTTP methods to controller methods.
func (c *PresenceController) URLMapping() {
	c.Mapping("GetAll", c.GetAll)   // Maps GET /presences to GetAll method for retrieving all presences or a user's presences based on the role
	c.Mapping("GetById", c.GetById) // Maps GET /presences/:id to GetById method for retrieving a specific presence by ID
	c.Mapping("Create", c.Create)   // Maps POST /presences to Create method for creating a new presence entry for a user (employee only)
	c.Mapping("Update", c.Update)   // Maps PUT /presences/:id to Update method for updating an existing presence entry by ID (admin only)
	c.Mapping("Delete", c.Delete)   // Maps DELETE /presences/:id to Delete method for deleting a specific presence entry by ID (admin only)
}

// @Title GetAll
// @Description Retrieve all presences or the presences of a specific user based on the role.
// @Param isIncludeUser query bool false "Include user data in the response"
// @Param isIncludeSchedule query bool false "Include schedule data in the response"
// @Success 200 {object} dto.PresenceResponseList "Success"
// @Failure 400 Bad Request
// @Failure 401 Unauthorized
// @Failure 500 Internal Server Error
// @router / [get]
func (c *PresenceController) GetAll() {
	// Retrieve user role and ID from the context
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

	// Parse query parameters
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
		// Admin can fetch all presences
		presences, err = models.GetAllPresences()
	} else if userRole == constants.RoleEmployee {
		// Employees can only fetch their own presences
		presences, err = models.GetPresencesByUserId(userId)
	}
	if err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusInternalServerError, "Failed to fetch presences", err)
		return
	}

	// Return success response with presence data
	helpers.SuccessResponse(c.Ctx.ResponseWriter, http.StatusOK, "Presences retrieved successfully", dto.FromPresenceModelListToPresenceResponseList(presences, isIncludeUser, isIncludeSchedule))
}

// @Title GetById
// @Description Retrieve a specific presence by ID.
// @Param id path int true "Presence ID"
// @Param isIncludeUser query bool false "Include user data in the response"
// @Param isIncludeSchedule query bool false "Include schedule data in the response"
// @Success 200 {object} dto.PresenceResponse "Success"
// @Failure 400 Bad Request
// @Failure 401 Unauthorized
// @Failure 404 Not Found
// @Failure 500 Internal Server Error
// @router /:id [get]
func (c *PresenceController) GetById() {
	// Parse query parameters
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

	// Get presence ID from URL
	id, _ := c.GetInt(":id")
	presence, err := models.GetPresenceById(id)
	if err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusNotFound, "Presence not found", err)
		return
	}

	// Return success response with the presence data
	helpers.SuccessResponse(c.Ctx.ResponseWriter, http.StatusOK, "Presence retrieved successfully", dto.FromPresenceModelToPresenceResponse(presence, isIncludeUser, isIncludeSchedule))
}

// @Title Create
// @Description Create a new presence entry for a user based on the schedule.
// @Param presence body dto.PresenceCreateRequest true "Presence data"
// @Success 201 {object} dto.PresenceResponse "Created"
// @Failure 400 Bad Request
// @Failure 401 Unauthorized
// @Failure 404 Not Found
// @Failure 500 Internal Server Error
// @router / [post]
func (c *PresenceController) Create() {
	// Retrieve authenticated user ID from the context
	userId, ok := c.Ctx.Input.GetData(constants.CtxAuthenticatedUserId).(int)
	if !ok {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusUnauthorized, "Bad context", errors.New("can't retrieve user role from context"))
		return
	}

	// Parse the request body to get presence data
	var req dto.PresenceCreateRequest
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusBadRequest, "Invalid input", err)
		return
	}

	// Validate the presence data
	if errorsMap, err := helpers.ValidatePayloads(req); err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusBadRequest, constants.ErrValidationMessage, errorsMap)
		return
	}

	// Fetch user details
	user, err := models.GetUserById(userId, false)
	if user == nil && err != nil {
		if err == orm.ErrNoRows {
			helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusNotFound, fmt.Sprintf("Failed to fetch user with id %d", userId), fmt.Errorf("user '%d' not found", userId))
			return
		}
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusInternalServerError, fmt.Sprintf("Failed to fetch user with id %d", req.ScheduleId), err)
		return
	}

	// Check if the user is assigned to the specified schedule
	if user.Schedule.Id != req.ScheduleId {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusBadRequest, "User is not assigned to the schedule", fmt.Errorf("user '%d' is not assigned to the schedule '%d'", userId, req.ScheduleId))
		return
	}

	// Fetch the schedule details
	schedule, err := models.GetScheduleById(req.ScheduleId, false, false)
	if schedule == nil && err != nil {
		if err == orm.ErrNoRows {
			helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusNotFound, fmt.Sprintf("Failed to fetch schedule with id %d", req.ScheduleId), fmt.Errorf("schedule '%d' not found", req.ScheduleId))
			return
		}
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusInternalServerError, fmt.Sprintf("Failed to fetch schedule with id %d", req.ScheduleId), err)
		return
	}

	// Check if the presence already exists for the user and type
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

	// Create the new presence entry
	presence := req.ToPresenceModelWithValue(user, schedule)

	// Determine the status of the presence (e.g., late, on time)
	status, err := helpers.DeterminePresenceStatus(presence.Type, schedule.InTime, schedule.OutTime, currentTime, constants.PresenceLateThreshold)
	if err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusInternalServerError, "Failed to determine presence status", err)
		return
	}
	presence.Status = status

	// Save the presence to the database
	if err := models.CreatePresence(presence); err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusInternalServerError, "Failed to create presence", err)
		return
	}

	// Return success response
	helpers.SuccessResponse(c.Ctx.ResponseWriter, http.StatusCreated, "Presence created successfully", dto.FromPresenceModelToPresenceResponse(presence, false, false))
}

// @Title Update
// @Description Update an existing presence entry by ID.
// @Param id path int true "Presence ID"
// @Param presence body dto.PresenceUpdateRequest true "Updated presence data"
// @Success 200 {object} dto.PresenceResponse "Success"
// @Failure 400 Bad Request
// @Failure 401 Unauthorized
// @Failure 404 Not Found
// @Failure 500 Internal Server Error
// @router /:id [put]
func (c *PresenceController) Update() {
	// Get the presence ID from the URL path
	id, _ := c.GetInt(":id")
	presence, err := models.GetPresenceById(id)
	if err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusNotFound, "Presence not found", err)
		return
	}

	// Parse the request body to get updated presence data
	var req dto.PresenceUpdateRequest
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusBadRequest, "Invalid input", err)
		return
	}

	// Validate the updated data
	if errorsMap, err := helpers.ValidatePayloads(req); err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusBadRequest, constants.ErrValidationMessage, errorsMap)
		return
	}

	// Fetch the schedule details
	schedule, err := models.GetScheduleById(req.ScheduleId, false, false)
	if schedule == nil && err != nil {
		if err == orm.ErrNoRows {
			helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusNotFound, fmt.Sprintf("Failed to fetch schedule with id %d", req.ScheduleId), fmt.Errorf("schedule '%d' not found", req.ScheduleId))
			return
		}
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusInternalServerError, fmt.Sprintf("Failed to fetch schedule with id %d", req.ScheduleId), err)
		return
	}

	// Fetch the user details
	user, err := models.GetUserById(req.UserId, false)
	if user == nil && err != nil {
		if err == orm.ErrNoRows {
			helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusNotFound, fmt.Sprintf("Failed to fetch user with id %d", req.UserId), fmt.Errorf("user '%d' not found", req.UserId))
			return
		}
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusInternalServerError, fmt.Sprintf("Failed to fetch user with id %d", req.ScheduleId), err)
		return
	}

	// Create the updated presence model
	updatedPresence := req.ToPresenceModelWithValue(presence, user, schedule)

	// Update the presence in the database
	if err := models.UpdatePresence(updatedPresence); err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusInternalServerError, "Failed to update presence", err)
		return
	}

	// Return success response
	helpers.SuccessResponse(c.Ctx.ResponseWriter, http.StatusOK, "Presence updated successfully", dto.FromPresenceModelToPresenceResponse(updatedPresence, false, false))
}

// @Title Delete
// @Description Delete a specific presence entry by ID.
// @Param id path int true "Presence ID"
// @Success 200 {object} helpers.SuccessResponse "Success"
// @Failure 400 Bad Request
// @Failure 401 Unauthorized
// @Failure 404 Not Found
// @Failure 500 Internal Server Error
// @router /:id [delete]
func (c *PresenceController) Delete() {
	// Get the presence ID from the URL path
	id, err := c.GetInt(":id")
	if err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusBadRequest, "Invalid presence id", err)
		return
	}

	// Delete the presence from the database
	affectedRows, err := models.DeletePresence(id)
	if err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusInternalServerError, "Failed to delete presence", err)
		return
	}

	// Check if any rows were affected (i.e., the presence was deleted)
	if affectedRows == 0 {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusNotFound, "Presence not found", fmt.Errorf("presence '%d' not found", id))
		return
	}

	// Return success response
	helpers.SuccessResponse(c.Ctx.ResponseWriter, http.StatusOK, "Presence deleted successfully", nil)
}
