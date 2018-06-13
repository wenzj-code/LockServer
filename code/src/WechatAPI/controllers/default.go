package controllers

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

//MainController .
type MainController struct {
	beego.Controller
}

//GetToken 通过APPID+Secrete生成token
func (c *MainController) GetToken() {
	appid := c.GetString("appid")
	secret := c.GetString("secret")
	log.Info("appid=", appid, ",secret=", secret)
	if appid == "" || secret == "" {
		c.Data["json"] = common.GetErrCodeJSON(10003)
		c.ServeJSON()
		return
	}

	status, err := DBOpt.GetDataOpt().CheckAppIDSecret(appid, secret)
	if err != nil {
		log.Error("err:", err)
		c.Data["json"] = common.GetErrCodeJSON(10006)
		c.ServeJSON()
		return
	}
	if !status {
		c.Data["json"] = common.GetErrCodeJSON(10001)
		c.ServeJSON()
		return
	}
	//组织应答报文
	data := make(map[string]interface{})
	hs := sha256.New()
	io.WriteString(hs, secret+time.Now().String())
	token := fmt.Sprintf("%x", string(hs.Sum(nil)))

	common.RedisTokenOpt.Set(token, 1, config.GetConfig().RedisTokenTimeOut)

	data["code"] = 0
	data["token"] = token
	data["expired_in"] = config.GetConfig().RedisTokenTimeOut
	c.Data["json"] = data
	c.ServeJSON()
}

//GetRoomInfo 通过设备ID获取房间信息
func (c *MainController) GetRoomInfo() {
	DeviceID := c.GetString("deviceid")
	Token := c.GetString("token")
	log.Info("DeviceID=", DeviceID, ",Token=", Token)
	if DeviceID == "" || Token == "" {
		c.Data["json"] = common.GetErrCodeJSON(10003)
		c.ServeJSON()
		return
	}

	log.Debug("token:", Token)
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
func (c *MainController) DoorCtrlOpen() {
	roomnu := c.GetString("roomnu")
	appid := c.GetString("appid")

	Token := c.GetString("token")
	log.Info("DoorCtrlOpen DeviceID=", roomnu, ",Token=", Token, ",appid:", appid)
	if roomnu == "" || appid == "" || Token == "" {
		c.Data["json"] = common.GetErrCodeJSON(10003)
		c.ServeJSON()
		return
	}

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

	gatewayID, gwOnline, devOnline, err := DBOpt.GetDataOpt().CheckGatewayOnline(DeviceID)
	if err != nil {
		log.Error("err:", err)
		c.Data["json"] = common.GetErrCodeJSON(10006)
		c.ServeJSON()
		return
	}

	data := make(map[string]interface{})

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

	err = DBOpt.GetDataOpt().WechatOpenMethod(DeviceID)
	if err != nil {
		log.Error("err:", err)
	}

	if gwOnline {
		//网关在线
		data["code"] = common.GetErrCodeJSON(0)
		if devOnline {
			//设备在线
			data["code"] = common.GetErrCodeJSON(0)
		} else {
			data["code"] = common.GetErrCodeJSON(10009)
		}
	} else {
		data["code"] = common.GetErrCodeJSON(10008)
	}
	c.Data["json"] = data
	c.ServeJSON()
	return
}
