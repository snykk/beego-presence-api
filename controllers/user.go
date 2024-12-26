package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/snykk/beego-presence-api/constants"
	"github.com/snykk/beego-presence-api/dto"
	"github.com/snykk/beego-presence-api/helpers"
	"github.com/snykk/beego-presence-api/models"
	"golang.org/x/crypto/bcrypt"

	beego "github.com/beego/beego/v2/server/web"
)

type UserController struct {
	beego.Controller
}

// @router /users [get]
func (c *UserController) GetAll() {
	users, err := models.GetAllUsers()
	if err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusInternalServerError, "Failed to fetch users", err)
		return
	}
	helpers.SuccessResponse(c.Ctx.ResponseWriter, http.StatusOK, "Users retrieved successfully", users)
}

// @router /users/:id [get]
func (c *UserController) GetById() {
	id, _ := c.GetInt(":id")
	user, err := models.GetUserById(id)
	if err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusNotFound, "User not found", err)
		return
	}
	helpers.SuccessResponse(c.Ctx.ResponseWriter, http.StatusOK, "User retrieved successfully", map[string]interface{}{"users": user})
}

// @router /auth/regis [post]
func (c *UserController) Register() {
	var req dto.RegisterRequest
	// Unmarshal request body ke dalam RegisterRequest struct
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusBadRequest, "Invalid input", err)
		return
	}

	department, err := models.GetDepartmentById(req.Department)
	if err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusBadRequest, "Invalid department ID", err)
		return
	}

	hashedPassword, err := helpers.HashPassword(req.Password)
	if err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusInternalServerError, "Failed to hash password", err)
		return
	}

	user := models.User{
		Name:       req.Name,
		Email:      req.Email,
		Password:   hashedPassword,
		Department: &department,
		Role:       constants.RoleStudent, // default registered user is STUDENT
	}

	if err := models.CreateUser(&user); err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusInternalServerError, "Failed to register user", err)
		return
	}

	helpers.SuccessResponse(c.Ctx.ResponseWriter, http.StatusCreated, "User registered successfully", map[string]interface{}{"users": dto.FromUserModelToRegisterResponse(user)})
}

// @router /auth/login [post]
func (c *UserController) Login() {
	var credentials dto.LoginRequest
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &credentials); err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusBadRequest, "Invalid input", err)
		return
	}

	user, err := models.GetUserByEmail(credentials.Email)
	if err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusUnauthorized, "Invalid credentials", nil)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
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
	id, _ := c.GetInt(":id")
	user, err := models.GetUserById(id)
	if err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusNotFound, "User not found", err)
		return
	}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &user); err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusBadRequest, "Invalid input", err)
		return
	}

	if err := models.UpdateUser(&user); err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusInternalServerError, "Failed to update user", err)
		return
	}
	helpers.SuccessResponse(c.Ctx.ResponseWriter, http.StatusOK, "User updated successfully", map[string]interface{}{"users": user})
}

// @router /users/:id [delete]
func (c *UserController) Delete() {
	id, _ := c.GetInt(":id")
	if err := models.DeleteUser(id); err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusInternalServerError, "Failed to delete user", err)
		return
	}
	helpers.SuccessResponse(c.Ctx.ResponseWriter, http.StatusOK, "User deleted successfully", nil)
}
