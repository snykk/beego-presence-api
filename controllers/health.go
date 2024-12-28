package controllers

import (
	"net/http"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/snykk/beego-presence-api/helpers"
)

// HealthController handles the health check endpoint
type HealthController struct {
	beego.Controller
}

// @router /health [get]
func (c *HealthController) CheckHealth() {
	helpers.SuccessResponse(c.Ctx.ResponseWriter, http.StatusOK, "Healthy...", nil)
}
