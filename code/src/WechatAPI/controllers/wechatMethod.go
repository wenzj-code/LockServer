package controllers

/*
	该模块主要用来接收微信端的请求设置
*/

import (
	"WechatAPI/DBOpt"
	"WechatAPI/common"
	"WechatAPI/config"
	"io/ioutil"
	"net/http"

	"crypto/sha256"
	"fmt"
	"io"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/astaxie/beego"
)

//WechatController .
type WechatController struct {
	beego.Controller
}

//GetToken 通过APPID+Secrete生成token
func (c *WechatController) GetToken() {
	appid := c.GetString("appid")
	secret := c.GetString("secret")
	log.Info("appid=", appid, ",secret=", secret)
	if appid == "" || secret == "" {
		c.Data["json"] = common.GetErrCodeJSON(10003)
		c.ServeJSON()
		return
	}

	//判断appid,secret是否存在,权限判断
	status, err := DBOpt.GetDataOpt().CheckAppIDSecret(appid, secret)
	if err != nil {
		log.Error("err:", err)
		c.Data["json"] = common.GetErrCodeJSON(10006)
		c.ServeJSON()
		return
	}
	if !status {
		log.Debug("appid,secrete不存在")
		c.Data["json"] = common.GetErrCodeJSON(10001)
		c.ServeJSON()
		return
	}
	//组织应答报文,生成token的方法
	data := make(map[string]interface{})
	hs := sha256.New()
	io.WriteString(hs, secret+time.Now().String())
	token := fmt.Sprintf("%x", string(hs.Sum(nil)))

	//将token保存到Redis缓存
	common.RedisTokenOpt.Set(token, 1, config.GetConfig().RedisTokenTimeOut)

	data["code"] = 0
	data["token"] = token
	data["expired_in"] = config.GetConfig().RedisTokenTimeOut
	c.Data["json"] = data
	c.ServeJSON()
}

//GetRoomInfo 通过设备ID获取房间信息
func (c *WechatController) GetRoomInfo() {
	DeviceID := c.GetString("deviceid")
	Token := c.GetString("token")
	log.Info("DeviceID=", DeviceID, ",Token=", Token)
	if DeviceID == "" || Token == "" {
		c.Data["json"] = common.GetErrCodeJSON(10003)
		c.ServeJSON()
		return
	}

	log.Debug("token:", Token)
	//从Redis里判断该token是否存在，不存在，则没有权限访问
	_, status, err := common.RedisTokenOpt.Get(Token)
	if err != nil {
		log.Error("err:", err)
		c.Data["json"] = common.GetErrCodeJSON(10007)
		c.ServeJSON()
		return
	}
	if !status {
		log.Info("Token数据不存在")
		c.Data["json"] = common.GetErrCodeJSON(10001)
		c.ServeJSON()
		return
	}

	//通过设备ID获取房间信息与对应的酒店appid
	roomnu, appid, err := DBOpt.GetDataOpt().GetRoomInfo(DeviceID)
	if err != nil {
		log.Error("err:", err)
		c.Data["json"] = common.GetErrCodeJSON(10006)
		c.ServeJSON()
		return
	}

	data := make(map[string]interface{})
	if len(roomnu) != 0 {
		data["roomnu"] = roomnu
		data["appid"] = appid
		data["code"] = 0
		c.Data["json"] = data
	} else {
		c.Data["json"] = common.GetErrCodeJSON(10005)
	}
	c.ServeJSON()
	return
}

//DoorCtrlOpen 开门
func (c *WechatController) DoorCtrlOpen() {
	roomnu := c.GetString("roomnu")
	appid := c.GetString("appid")
	method, err := c.GetInt("method")
	if err != nil {
		log.Error("err:", err)
		c.Data["json"] = common.GetErrCodeJSON(10003)
		c.ServeJSON()
		return
	}

	Token := c.GetString("token")
	log.Info("DoorCtrlOpen DeviceID=", roomnu, ",Token=", Token, ",appid:", appid)
	if roomnu == "" || appid == "" || Token == "" {
		c.Data["json"] = common.GetErrCodeJSON(10003)
		c.ServeJSON()
		return
	}

	//从Redis里判断该token是否存在，不存在，则没有权限访问
	_, status, err := common.RedisTokenOpt.Get(Token)
	if err != nil {
		log.Error("err:", err)
		c.Data["json"] = common.GetErrCodeJSON(10007)
		c.ServeJSON()
		return
	}
	if !status {
		log.Info("Token数据不存在")
		c.Data["json"] = common.GetErrCodeJSON(10001)
		c.ServeJSON()
		return
	}

	//通过房间号与酒店appid获取设备id信息
	DeviceID, err := DBOpt.GetDataOpt().GetDeviceID(roomnu, appid)
	if err != nil {
		log.Error("err:", err)
		c.Data["json"] = common.GetErrCodeJSON(10006)
		c.ServeJSON()
		return
	}
	if len(DeviceID) == 0 {
		log.Error("房间数据不存在:", roomnu, ",userid:", appid)
		c.Data["json"] = common.GetErrCodeJSON(10004)
		c.ServeJSON()
		return
	}
	log.Debug("DeviceID:", DeviceID)

	//DeviceID := DeviceIDAgent[4:]

	// agentid, err := DBOpt.GetDataOpt().GetAgentIDAPPID(appid)
	// if err != nil {
	// 	log.Error("err:", err)
	// 	c.Data["json"] = common.GetErrCodeJSON(10006)
	// 	c.ServeJSON()
	// 	return
	// }

	// agentidStr := fmt.Sprintf("%04d", agentid)

	//通过设备ID获取网关ID与在线状态
	gatewayID, gwOnline, devOnline, err := DBOpt.GetDataOpt().CheckGatewayOnline(DeviceID)
	if err != nil {
		log.Error("err:", err)
		c.Data["json"] = common.GetErrCodeJSON(10006)
		c.ServeJSON()
		return
	}

	//用Redis获取该网关连接到哪台服务器，并且或者所在连接的服务器地址
	serverIP, isExist, err := common.RedisServerListOpt.Get(gatewayID)
	if err != nil {
		log.Error("err:", err)
		c.Data["json"] = common.GetErrCodeJSON(10007)
		c.ServeJSON()
		return
	}
	if !isExist {
		log.Error("err:", err)
		c.Data["json"] = common.GetErrCodeJSON(10008)
		c.ServeJSON()
		return
	}

	//目前网关心跳只是一人空包，没有网关ＩＤ，无法做到网关是否线
	gwOnline = true
	devOnline = true
	var errcode int
	if gwOnline {
		//网关在线
		errcode = 0
		if devOnline {
			errcode = 0
		} else {
			//设备不在线
			errcode = 10009
		}
	} else {
		//网关不在线
		errcode = 10008
	}
	c.Data["json"] = common.GetErrCodeJSON(errcode)
	c.ServeJSON()
	if errcode != 0 {
		log.Info("网关或者设备不在线：gw=", gatewayID, ",deviceID=", DeviceID)
		return
	}

	//向设备服务请求开门
	httpServerIP := fmt.Sprintf("http://%s/dev-ctrl?gwid=%s&deviceid=%s", serverIP, gatewayID, DeviceID)
	log.Debug("httpServerIP:", httpServerIP)
	resp, err := http.Get(httpServerIP)
	if err != nil {
		log.Error("err:", err)
		c.Data["json"] = common.GetErrCodeJSON(10000)
		c.ServeJSON()
		return
	}
	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("err:", err)
		c.Data["json"] = common.GetErrCodeJSON(10000)
		c.ServeJSON()
		return
	}

	if method == 1 {
		//保存开门信息
		err = DBOpt.GetDataOpt().WechatOpenMethod(DeviceID)
	} else {
		err = DBOpt.GetDataOpt().WechatOpenMethod(DeviceID)
	}
	if err != nil {
		log.Error("err:", err)
	}

	return
}
