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
	"golang.org/x/crypto/bcrypt"

	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"
)

type UserController struct {
	beego.Controller
}

// @router /users [get]
func (c *UserController) GetAll() {
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

	users, err := models.GetAllUsers(isIncludePresenceList)
	if err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusInternalServerError, "Failed to fetch users", err)
		return
	}

	helpers.SuccessResponse(c.Ctx.ResponseWriter, http.StatusOK, "Users retrieved successfully", dto.FromUserModelListToUserResponseList(users, isIncludeDepartment, isIncludePresenceList, isIncludeSchedule))
}

// @router /users/:id [get]
func (c *UserController) GetById() {
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

	id, _ := c.GetInt(":id")
	user, err := models.GetUserById(id, isIncludePresenceList)
	if err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusNotFound, "User not found", err)
		return
	}
	helpers.SuccessResponse(c.Ctx.ResponseWriter, http.StatusOK, "User retrieved successfully", map[string]interface{}{"users": dto.FromUserModelToUserResponse(user, isIncludeDepartment, isIncludePresenceList, isIncludeSchedule)})
}

// @router /auth/regis [post]
func (c *UserController) Register() {
	var req dto.RegisterRequest
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusBadRequest, "Invalid input", err)
		return
	}

	if errorsMap, err := helpers.ValidatePayloads(req); err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusBadRequest, constants.ErrValidationMessage, errorsMap)
		return
	}

	department, err := models.GetDepartmentById(req.DepartmentId, false, false)
	if err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusBadRequest, "Invalid department ID", err)
		return
	}

	hashedPassword, err := helpers.HashPassword(req.Password)
	if err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusInternalServerError, "Failed to hash password", err)
		return
	}

	user := req.ToUserModel(department)
	user.Password = hashedPassword

	if err := models.CreateUser(user); err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusInternalServerError, "Failed to register user", err)
		return
	}

	helpers.SuccessResponse(c.Ctx.ResponseWriter, http.StatusCreated, "User registered successfully", map[string]interface{}{"users": dto.FromUserModelToRegisterResponse(user)})
}

// @router /auth/login [post]
func (c *UserController) Login() {
	var req dto.LoginRequest
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusBadRequest, "Invalid input", err)
		return
	}

	if errorsMap, err := helpers.ValidatePayloads(req); err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusBadRequest, constants.ErrValidationMessage, errorsMap)
		return
	}

	user, err := models.GetUserByEmail(req.Email)
	if err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusUnauthorized, "Invalid credentials", nil)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusUnauthorized, "Invalid credentials", nil)
		return
	}

	token, err := helpers.GenerateJWT(user.Id, user.Email, user.Role)
	if err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusInternalServerError, "Failed to generate token", err)
		return
	}

	helpers.SuccessResponse(c.Ctx.ResponseWriter, http.StatusOK, "Login successful", map[string]string{"token": token})
}

// @router /users/:id [put]
func (c *UserController) Update() {
	userId, ok := c.Ctx.Input.GetData(constants.CtxAuthenticatedUserId).(int)
	if !ok {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusUnauthorized, "Bad context", errors.New("can't retrieve user id from context"))
		return
	}

	id, err := c.GetInt(":id")
	if err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	if userId != id {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusForbidden, "You are not allowed to update another user's data", errors.New("forbidden access"))
		return
	}

	user, err := models.GetUserById(id, false)
	if err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusNotFound, "User not found", err)
		return
	}

	var req dto.UserRequest
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusBadRequest, "Invalid input format", err)
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

	updatedUser := req.ToUserModel(user, department)

	// Menyimpan perubahan ke database
	if err := models.UpdateUser(updatedUser); err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusInternalServerError, "Failed to update user", err)
		return
	}

	helpers.SuccessResponse(c.Ctx.ResponseWriter, http.StatusOK, "User updated successfully", map[string]interface{}{
		"user": dto.FromUserModelToUserResponse(updatedUser, false, false, false),
	})
}

// @router /users/:id [delete]
func (c *UserController) Delete() {
	userId, ok := c.Ctx.Input.GetData(constants.CtxAuthenticatedUserId).(int)
	if !ok {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusUnauthorized, "Bad context", errors.New("can't retrieve user id from context"))
		return
	}

	id, err := c.GetInt(":id")
	if err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusBadRequest, "Invalid user id", err)
		return
	}

	if userId != id {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusForbidden, "You are not allowed to delete another user's data", errors.New("forbidden access"))
		return
	}

	affectedRows, err := models.DeleteUser(id)
	if err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusInternalServerError, "Failed to delete user", err)
		return
	}

	if affectedRows == 0 {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusNotFound, "User not found", fmt.Errorf("user '%d' not found", id))
		return
	}

	helpers.SuccessResponse(c.Ctx.ResponseWriter, http.StatusOK, "User deleted successfully", nil)
}
