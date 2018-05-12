package routers

import (
	"WechatAPI/controllers"

	"github.com/astaxie/beego"
)

func init() {
	// beego.Router("/", &controllers.MainController{})
	beego.Router("/v1/token", &controllers.MainController{}, "get:GetToken")
	beego.Router("/v1/open-door", &controllers.MainController{}, "get:DoorCtrlOpen")
	beego.Router("/v1/get-roominfo", &controllers.MainController{}, "get:GetRoomInfo")
}
