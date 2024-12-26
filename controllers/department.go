package controllers

import (
	"encoding/json"

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
		c.Data["json"] = map[string]string{"error": err.Error()}
	} else {
		c.Data["json"] = departments
	}
	c.ServeJSON()
}

// @router /departments/:id [get]
func (c *DepartmentController) GetById() {
	id, _ := c.GetInt(":id")
	department, err := models.GetDepartmentById(id)
	if err != nil {
		c.Data["json"] = map[string]string{"error": "Department not found"}
	} else {
		c.Data["json"] = department
	}
	c.ServeJSON()
}

// @router /departments [post]
func (c *DepartmentController) Create() {
	var department models.Department
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &department); err != nil {
		c.Data["json"] = map[string]string{"error": "Invalid input"}
	} else {
		if err := models.CreateDepartment(&department); err != nil {
			c.Data["json"] = map[string]string{"error": err.Error()}
		} else {
			c.Data["json"] = department
		}
	}
	c.ServeJSON()
}

// @router /departments/:id [put]
func (c *DepartmentController) Update() {
	id, _ := c.GetInt(":id")
	department, err := models.GetDepartmentById(id)
	if err != nil {
		c.Data["json"] = map[string]string{"error": "Department not found"}
	} else {
		if err := json.Unmarshal(c.Ctx.Input.RequestBody, &department); err != nil {
			c.Data["json"] = map[string]string{"error": "Invalid input"}
		} else {
			if err := models.UpdateDepartment(&department); err != nil {
				c.Data["json"] = map[string]string{"error": err.Error()}
			} else {
				c.Data["json"] = department
			}
		}
	}
	c.ServeJSON()
}

// @router /departments/:id [delete]
func (c *DepartmentController) Delete() {
	id, _ := c.GetInt(":id")
	if err := models.DeleteDepartment(id); err != nil {
		c.Data["json"] = map[string]string{"error": err.Error()}
	} else {
		c.Data["json"] = map[string]string{"message": "Department deleted"}
	}
	c.ServeJSON()
}
