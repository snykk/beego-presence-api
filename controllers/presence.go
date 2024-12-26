package controllers

import (
	"encoding/json"

	"github.com/snykk/beego-presence-api/models"

	beego "github.com/beego/beego/v2/server/web"
)

type PresenceController struct {
	beego.Controller
}

// @router /presences [get]
func (c *PresenceController) GetAll() {
	presences, err := models.GetAllPresences()
	if err != nil {
		c.Data["json"] = map[string]string{"error": err.Error()}
	} else {
		c.Data["json"] = presences
	}
	c.ServeJSON()
}

// @router /presences/:id [get]
func (c *PresenceController) GetById() {
	id, _ := c.GetInt(":id")
	presence, err := models.GetPresenceById(id)
	if err != nil {
		c.Data["json"] = map[string]string{"error": "Presence not found"}
	} else {
		c.Data["json"] = presence
	}
	c.ServeJSON()
}

// @router /presences [post]
func (c *PresenceController) Create() {
	var presence models.Presence
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &presence); err != nil {
		c.Data["json"] = map[string]string{"error": "Invalid input"}
	} else {
		if err := models.CreatePresence(&presence); err != nil {
			c.Data["json"] = map[string]string{"error": err.Error()}
		} else {
			c.Data["json"] = presence
		}
	}
	c.ServeJSON()
}

// @router /presences/:id [put]
func (c *PresenceController) Update() {
	id, _ := c.GetInt(":id")
	presence, err := models.GetPresenceById(id)
	if err != nil {
		c.Data["json"] = map[string]string{"error": "Presence not found"}
	} else {
		if err := json.Unmarshal(c.Ctx.Input.RequestBody, &presence); err != nil {
			c.Data["json"] = map[string]string{"error": "Invalid input"}
		} else {
			if err := models.UpdatePresence(presence); err != nil {
				c.Data["json"] = map[string]string{"error": err.Error()}
			} else {
				c.Data["json"] = presence
			}
		}
	}
	c.ServeJSON()
}

// @router /presences/:id [delete]
func (c *PresenceController) Delete() {
	id, _ := c.GetInt(":id")
	if err := models.DeletePresence(id); err != nil {
		c.Data["json"] = map[string]string{"error": err.Error()}
	} else {
		c.Data["json"] = map[string]string{"message": "Presence deleted"}
	}
	c.ServeJSON()
}
