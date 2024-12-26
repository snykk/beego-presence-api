package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/snykk/beego-presence-api/helpers"
	"github.com/snykk/beego-presence-api/models"

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
	helpers.SuccessResponse(c.Ctx.ResponseWriter, http.StatusOK, "User retrieved successfully", user)
}

// @router /users [post]
func (c *UserController) Create() {
	var user models.User
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &user); err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusBadRequest, "Invalid input", err)
		return
	}

	if err := models.CreateUser(&user); err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusInternalServerError, "Failed to create user", err)
		return
	}
	helpers.SuccessResponse(c.Ctx.ResponseWriter, http.StatusCreated, "User created successfully", user)
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
	helpers.SuccessResponse(c.Ctx.ResponseWriter, http.StatusOK, "User updated successfully", user)
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
