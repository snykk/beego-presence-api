package controllers

import (
	"encoding/json"

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
		c.Data["json"] = map[string]string{"error": err.Error()}
	} else {
		c.Data["json"] = users
	}
	c.ServeJSON()
}

// @router /users/:id [get]
func (c *UserController) GetById() {
	id, _ := c.GetInt(":id")
	user, err := models.GetUserById(id)
	if err != nil {
		c.Data["json"] = map[string]string{"error": "User not found"}
	} else {
		c.Data["json"] = user
	}
	c.ServeJSON()
}

// @router /users [post]
func (c *UserController) Create() {
	var user models.User
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &user); err != nil {
		c.Data["json"] = map[string]string{"error": "Invalid input"}
	} else {
		if err := models.CreateUser(&user); err != nil {
			c.Data["json"] = map[string]string{"error": err.Error()}
		} else {
			c.Data["json"] = user
		}
	}
	c.ServeJSON()
}

// @router /users/:id [put]
func (c *UserController) Update() {
	id, _ := c.GetInt(":id")
	user, err := models.GetUserById(id)
	if err != nil {
		c.Data["json"] = map[string]string{"error": "User not found"}
	} else {
		if err := json.Unmarshal(c.Ctx.Input.RequestBody, &user); err != nil {
			c.Data["json"] = map[string]string{"error": "Invalid input"}
		} else {
			if err := models.UpdateUser(&user); err != nil {
				c.Data["json"] = map[string]string{"error": err.Error()}
			} else {
				c.Data["json"] = user
			}
		}
	}
	c.ServeJSON()
}

// @router /users/:id [delete]
func (c *UserController) Delete() {
	id, _ := c.GetInt(":id")
	if err := models.DeleteUser(id); err != nil {
		c.Data["json"] = map[string]string{"error": err.Error()}
	} else {
		c.Data["json"] = map[string]string{"message": "User deleted"}
	}
	c.ServeJSON()
}
