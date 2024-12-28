// @APIVersion 1.0.0
// @Title Beego Presence API
// @Description This API provides endpoints for managing presence, schedules, users, and departments.
// @Contact najibfikri13@gmail.com
// @TermsOfServiceUrl http://example.com/terms
// @License MIT
// @LicenseUrl https://github.com/snykk/beego-presence-api/blob/master/LICENSE
// @Schemes http, https
// @Host localhost:8080
// @BasePath /api/v1
// @Name Beego Presence API
// @URL http://localhost:8080/api/v1
package routers

import (
	beego "github.com/beego/beego/v2/server/web"
	"github.com/snykk/beego-presence-api/controllers"
	"github.com/snykk/beego-presence-api/middlewares"
)

func init() {
	// Register root and health endpoints
	beego.Router("/", &controllers.RootController{}, "get:GetRoot")
	beego.Router("/health", &controllers.HealthController{}, "get:CheckHealth")

	// Define namespaces and include controllers
	ns := beego.NewNamespace("/api/v1",
		beego.NSBefore(middlewares.RoleBasedMiddleware()),
		beego.NSNamespace("/auth",
			// Create routes for the UserController in auth endpoint
			beego.NSRouter("/regis", &controllers.AuthController{}, "post:Register"),
			beego.NSRouter("/login", &controllers.AuthController{}, "post:Login"),

			// To generate the swagger documentation for the UserController in auth endpoint
			beego.NSInclude(
				&controllers.AuthController{},
			),
		),
		beego.NSNamespace("/users",
			// Create routes for the UserController in users endpoint
			beego.NSRouter("", &controllers.UserController{}, "get:GetAll"),
			beego.NSRouter("/:id", &controllers.UserController{}, "get:GetById;put:Update;delete:Delete"),

			// To generate the swagger documentation for the UserController in users endpoint
			beego.NSInclude(
				&controllers.UserController{},
			),
		),
		beego.NSNamespace("/departments",
			// Create routes for the DepartmentController
			beego.NSRouter("", &controllers.DepartmentController{}, "get:GetAll;post:Create"),
			beego.NSRouter("/:id", &controllers.DepartmentController{}, "get:GetById;put:Update;delete:Delete"),

			// To generate the swagger documentation for the DepartmentController
			beego.NSInclude(
				&controllers.DepartmentController{},
			),
		),
		beego.NSNamespace("/schedules",
			// Create routes for the ScheduleController
			beego.NSRouter("", &controllers.ScheduleController{}, "get:GetAll;post:Create"),
			beego.NSRouter("/:id", &controllers.ScheduleController{}, "get:GetById;put:Update;delete:Delete"),

			// To generate the swagger documentation for the ScheduleController
			beego.NSInclude(
				&controllers.ScheduleController{},
			),
		),
		beego.NSNamespace("/presences",
			// Create routes for the PresenceController
			beego.NSRouter("", &controllers.PresenceController{}, "get:GetAll;post:Create"),
			beego.NSRouter("/:id", &controllers.PresenceController{}, "get:GetById;put:Update;delete:Delete"),

			// To generate the swagger documentation for the PresenceController
			beego.NSInclude(
				&controllers.PresenceController{},
			),
		),
	)

	// Register namespace
	beego.AddNamespace(ns)
}
