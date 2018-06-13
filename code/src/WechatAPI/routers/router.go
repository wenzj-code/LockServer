package routers

import (
	"WechatAPI/controllers"

	"github.com/astaxie/beego"
)

func init() {
	// beego.Router("/", &controllers.MainController{})
	//开门接口
	beego.Router("/v1/token", &controllers.MainController{}, "get:GetToken")
	beego.Router("/v1/open-door", &controllers.MainController{}, "get:DoorCtrlOpen")
	beego.Router("/v1/get-roominfo", &controllers.MainController{}, "get:GetRoomInfo")

	//APP扫描绑定接口
	beego.Router("/v1/login", &controllers.MainController{}, "get:AppLogin")
	beego.Router("/v1/add-gateway", &controllers.MainController{}, "get:AddGateway")
	beego.Router("/v1/bind-room", &controllers.MainController{}, "get:BindDeviceRoom")

	//模拟推送接收接口
	beego.Router("/test/token", &controllers.MainController{}, "get:TestToken")
	beego.Router("/test/push", &controllers.MainController{}, "post:TestPush")

	//接收设备服务的状态上报接口
	beego.Router("/report/dev-status", &controllers.MainController{}, "get:DoorRecvReport")
}
