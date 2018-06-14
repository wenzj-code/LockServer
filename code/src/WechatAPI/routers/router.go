package routers

/*
	beego框架的路由实现方式
*/
import (
	"WechatAPI/controllers"

	"github.com/astaxie/beego"
)

func init() {
	// beego.Router("/", &controllers.MainController{})
	//开门接口
	beego.Router("/v1/token", &controllers.WechatController{}, "get:GetToken")
	beego.Router("/v1/open-door", &controllers.WechatController{}, "get:DoorCtrlOpen")
	beego.Router("/v1/get-roominfo", &controllers.WechatController{}, "get:GetRoomInfo")

	//APP扫描绑定接口
	beego.Router("/v1/login", &controllers.AppController{}, "get:AppLogin")
	beego.Router("/v1/add-gateway", &controllers.AppController{}, "get:AddGateway")
	beego.Router("/v1/bind-room", &controllers.AppController{}, "get:BindDeviceRoom")

	//模拟推送接收接口
	beego.Router("/test/token", &controllers.TestPushServerController{}, "get:TestToken")
	beego.Router("/test/push", &controllers.TestPushServerController{}, "post:TestPush")

	//接收设备服务的状态上报接口
	beego.Router("/report/dev-status", &controllers.DevStatusController{}, "get:DoorRecvReport")
}
