package controllers

import (
	"encoding/json"

	"github.com/snykk/beego-presence-api/models"

	beego "github.com/beego/beego/v2/server/web"
)

type ScheduleController struct {
	beego.Controller
}

// @router /schedules [get]
func (c *ScheduleController) GetAll() {
	schedules, err := models.GetAllSchedules()
	if err != nil {
		c.Data["json"] = map[string]string{"error": err.Error()}
	} else {
		c.Data["json"] = schedules
	}
	c.ServeJSON()
}

// @router /schedules/:id [get]
func (c *ScheduleController) GetById() {
	id, _ := c.GetInt(":id")
	schedule, err := models.GetScheduleById(id)
	if err != nil {
		c.Data["json"] = map[string]string{"error": "Schedule not found"}
	} else {
		c.Data["json"] = schedule
	}
	c.ServeJSON()
}

// @router /schedules [post]
func (c *ScheduleController) Create() {
	var schedule models.Schedule
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &schedule); err != nil {
		c.Data["json"] = map[string]string{"error": "Invalid input"}
	} else {
		if err := models.CreateSchedule(&schedule); err != nil {
			c.Data["json"] = map[string]string{"error": err.Error()}
		} else {
			c.Data["json"] = schedule
		}
	}
	c.ServeJSON()
}

// @router /schedules/:id [put]
func (c *ScheduleController) Update() {
	id, _ := c.GetInt(":id")
	schedule, err := models.GetScheduleById(id)
	if err != nil {
		c.Data["json"] = map[string]string{"error": "Schedule not found"}
	} else {
		if err := json.Unmarshal(c.Ctx.Input.RequestBody, &schedule); err != nil {
			c.Data["json"] = map[string]string{"error": "Invalid input"}
		} else {
			if err := models.UpdateSchedule(&schedule); err != nil {
				c.Data["json"] = map[string]string{"error": err.Error()}
			} else {
				c.Data["json"] = schedule
			}
		}
	}
	c.ServeJSON()
}

// @router /schedules/:id [delete]
func (c *ScheduleController) Delete() {
	id, _ := c.GetInt(":id")
	if err := models.DeleteSchedule(id); err != nil {
		c.Data["json"] = map[string]string{"error": err.Error()}
	} else {
		c.Data["json"] = map[string]string{"message": "Schedule deleted"}
	}
	c.ServeJSON()
}
