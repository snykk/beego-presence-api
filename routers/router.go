package routers

import (
	beego "github.com/beego/beego/v2/server/web"
	"github.com/snykk/beego-presence-api/controllers"
	"github.com/snykk/beego-presence-api/middlewares"
)

func init() {
	beego.InsertFilter("/*", beego.BeforeRouter, middlewares.RoleBasedMiddleware()) // Apply to all other routes

	// Define routes
	beego.Router("/auth/regis", &controllers.UserController{}, "post:Register")
	beego.Router("/auth/login", &controllers.UserController{}, "post:Login")

	beego.Router("/users", &controllers.UserController{}, "get:GetAll")
	beego.Router("/users/:id", &controllers.UserController{}, "get:GetById;put:Update;delete:Delete")

	beego.Router("/departments", &controllers.DepartmentController{}, "get:GetAll;post:Create")
	beego.Router("/departments/:id", &controllers.DepartmentController{}, "get:GetById;put:Update;delete:Delete")

	beego.Router("/schedules", &controllers.ScheduleController{}, "get:GetAll;post:Create")
	beego.Router("/schedules/:id", &controllers.ScheduleController{}, "get:GetById;put:Update;delete:Delete")

	beego.Router("/presences", &controllers.PresenceController{}, "get:GetAll;post:Create")
	beego.Router("/presences/:id", &controllers.PresenceController{}, "get:GetById;put:Update;delete:Delete")
}
