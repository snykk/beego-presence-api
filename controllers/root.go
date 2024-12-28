package controllers

import (
	"net/http"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/snykk/beego-presence-api/helpers"
)

// RootController handles the root endpoint
type RootController struct {
	beego.Controller
}

// @router / [get]
func (c *RootController) GetRoot() {
	helpers.SuccessResponse(c.Ctx.ResponseWriter, http.StatusOK, "Welcome to Beego Presence API", nil)
}
