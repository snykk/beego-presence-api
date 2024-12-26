// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"github.com/snykk/beego-presence-api/controllers"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {
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
