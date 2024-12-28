package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/snykk/beego-presence-api/constants"
	"github.com/snykk/beego-presence-api/dto"
	"github.com/snykk/beego-presence-api/helpers"
	"github.com/snykk/beego-presence-api/models"

	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"
)

// UserController handles user-related operations such as retrieving, creating, updating, and deleting users.
type UserController struct {
	beego.Controller
}

// URLMapping maps routes to specific handler functions for the UserController
// This is typically used by the Beego framework to map HTTP methods to controller methods.
func (c *UserController) URLMapping() {
	c.Mapping("GetAll", c.GetAll)   // Maps GET /users to GetAll method for retrieving all users
	c.Mapping("GetById", c.GetById) // Maps GET /users/:id to GetById method for retrieving a specific user by ID
	c.Mapping("Update", c.Update)   // Maps PUT /users/:id to Update method for updating a specific user by ID
	c.Mapping("Delete", c.Delete)   // Maps DELETE /users/:id to Delete method for deleting a specific user by ID
}

// @Title Get All Users
// @Description Fetch all users with optional related data (department, presence list, schedule)
// @Accept  json
// @Produce  json
// @Param isIncludeDepartment query bool false "Include department data"
// @Param isIncludePresenceList query bool false "Include user presence list"
// @Param isIncludeSchedule query bool false "Include user schedule data"
// @Success 200 {object} dto.UserResponse "Users retrieved successfully"
// @Failure 400 Invalid query parameters
// @Failure 500 Failed to fetch users
// @router / [get]
func (c *UserController) GetAll() {
	// Read and parse the query parameters for optional data inclusion.
	isIncludeDepartment, err := c.GetBool("isIncludeDepartment", false)
	if err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusBadRequest, "Invalid value for isIncludeDepartment", err)
		return
	}

	isIncludePresenceList, err := c.GetBool("isIncludePresenceList", false)
	if err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusBadRequest, "Invalid value for isIncludePresenceList", err)
		return
	}

	isIncludeSchedule, err := c.GetBool("isIncludeSchedule", false)
	if err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusBadRequest, "Invalid value for isIncludeSchedule", err)
		return
	}

	// Fetch all users with the optional data inclusion.
	users, err := models.GetAllUsers(isIncludePresenceList)
	if err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusInternalServerError, "Failed to fetch users", err)
		return
	}

	// Return the list of users.
	helpers.SuccessResponse(c.Ctx.ResponseWriter, http.StatusOK, "Users retrieved successfully", dto.FromUserModelListToUserResponseList(users, isIncludeDepartment, isIncludePresenceList, isIncludeSchedule))
}

// @Title Get User By ID
// @Description Fetch a user by their ID with optional related data (department, presence list, schedule)
// @Accept  json
// @Produce  json
// @Param id path int true "User ID"
// @Param isIncludeDepartment query bool false "Include department data"
// @Param isIncludePresenceList query bool false "Include user presence list"
// @Param isIncludeSchedule query bool false "Include user schedule data"
// @Success 200 {object} dto.UserResponse "User retrieved successfully"
// @Failure 400 Invalid query parameters
// @Failure 404 User not found
// @Failure 500 Failed to fetch user
// @router /:id [get]
func (c *UserController) GetById() {
	// Read and parse the query parameters for optional data inclusion.
	isIncludeDepartment, err := c.GetBool("isIncludeDepartment", false)
	if err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusBadRequest, "Invalid value for isIncludeDepartment", err)
		return
	}

	isIncludePresenceList, err := c.GetBool("isIncludePresenceList", false)
	if err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusBadRequest, "Invalid value for isIncludePresenceList", err)
		return
	}

	isIncludeSchedule, err := c.GetBool("isIncludeSchedule", false)
	if err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusBadRequest, "Invalid value for isIncludeSchedule", err)
		return
	}

	// Fetch the user by ID.
	id, _ := c.GetInt(":id")
	user, err := models.GetUserById(id, isIncludePresenceList)
	if err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusNotFound, "User not found", err)
		return
	}

	// Return the user details.
	helpers.SuccessResponse(c.Ctx.ResponseWriter, http.StatusOK, "User retrieved successfully", map[string]interface{}{"users": dto.FromUserModelToUserResponse(user, isIncludeDepartment, isIncludePresenceList, isIncludeSchedule)})
}

// @Title Update User
// @Description Update an existing user's details
// @Accept  json
// @Produce  json
// @Param id path int true "User ID"
// @Param userRequest body dto.UserRequest true "User Data"
// @Success 200 {object} dto.UserResponse "User updated successfully"
// @Failure 400 Invalid input data
// @Failure 404 User not found
// @Failure 500 Failed to update user
// @router /:id [put]
func (c *UserController) Update() {
	// Fetch the authenticated user's ID
	userId, ok := c.Ctx.Input.GetData(constants.CtxAuthenticatedUserId).(int)
	if !ok {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusUnauthorized, "Bad context", errors.New("can't retrieve user id from context"))
		return
	}

	// Fetch the user ID from the URL.
	id, err := c.GetInt(":id")
	if err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	// Ensure the user can only update their own data.
	if userId != id {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusForbidden, "You are not allowed to update another user's data", errors.New("forbidden access"))
		return
	}

	// Fetch the user by ID.
	user, err := models.GetUserById(id, false)
	if err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusNotFound, "User not found", err)
		return
	}

	// Parse the request body to update user data.
	var req dto.UserRequest
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusBadRequest, "Invalid input format", err)
		return
	}

	// Validate the updated data.
	if errorsMap, err := helpers.ValidatePayloads(req); err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusBadRequest, constants.ErrValidationMessage, errorsMap)
		return
	}

	// Fetch the department for the user.
	department, err := models.GetDepartmentById(req.DepartmentId, false, false)
	if department == nil && err != nil {
		if err == orm.ErrNoRows {
			helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusNotFound, fmt.Sprintf("Failed to fetch department with id %d", req.DepartmentId), fmt.Errorf("department '%d' not found", req.DepartmentId))
			return
		}
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusInternalServerError, fmt.Sprintf("Failed to fetch department with id %d", req.DepartmentId), err)
		return
	}

	// Update user model and save to database.
	updatedUser := req.ToUserModel(user, department)
	if err := models.UpdateUser(updatedUser); err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusInternalServerError, "Failed to update user", err)
		return
	}

	helpers.SuccessResponse(c.Ctx.ResponseWriter, http.StatusOK, "User updated successfully", map[string]interface{}{"user": dto.FromUserModelToUserResponse(updatedUser, false, false, false)})
}

// @Title Delete User
// @Description Delete a user by ID
// @Accept  json
// @Produce  json
// @Param id path int true "User ID"
// @Success 200 {object} helpers.SuccessResponse "User deleted successfully"
// @Failure 400 Invalid user ID
// @Failure 404 User not found
// @Failure 500 Failed to delete user
// @router /:id [delete]
func (c *UserController) Delete() {
	// Fetch the authenticated user's ID.
	userId, ok := c.Ctx.Input.GetData(constants.CtxAuthenticatedUserId).(int)
	if !ok {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusUnauthorized, "Bad context", errors.New("can't retrieve user id from context"))
		return
	}

	// Fetch the user ID from the URL.
	id, err := c.GetInt(":id")
	if err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusBadRequest, "Invalid user id", err)
		return
	}

	// Ensure the user can only delete their own data.
	if userId != id {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusForbidden, "You are not allowed to delete another user's data", errors.New("forbidden access"))
		return
	}

	// Attempt to delete the user by ID.
	affectedRows, err := models.DeleteUser(id)
	if err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusInternalServerError, "Failed to delete user", err)
		return
	}

	// If no rows were affected, the user was not found
	if affectedRows == 0 {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusNotFound, "User not found", fmt.Errorf("user '%d' not found", id))
		return
	}

	// Return success response indicating user was deleted.
	helpers.SuccessResponse(c.Ctx.ResponseWriter, http.StatusOK, "User deleted successfully", nil)
}
