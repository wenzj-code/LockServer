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
	beego.Router("/v1/setting-card-password", &controllers.WechatController{}, "get:SettingCardPassword")
	beego.Router("/v1/cancel-card-password", &controllers.WechatController{}, "get:CancleCardPassword")

	//模拟推送接收接口
	beego.Router("/test/token", &controllers.TestPushServerController{}, "get:TestToken")
	beego.Router("/test/push", &controllers.TestPushServerController{}, "post:TestPush")

	//接收设备服务的状态上报接口
	beego.Router("/report/door-ctrl-rsp", &controllers.DevStatusController{}, "get:DoorCtrlRsp")
	beego.Router("/report/dev-setting-password-status", &controllers.DevStatusController{}, "get:SettingCardlRsp")
	beego.Router("/report/dev-cancel-password-status", &controllers.DevStatusController{}, "get:CancelCardlRsp")
	beego.Router("/report/card-openlock-record", &controllers.DevStatusController{}, "get:CardDoorOpenlRsp")
}
